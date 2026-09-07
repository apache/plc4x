/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package opcua

import (
	"context"
	"slices"
	"sync"
	"time"

	"github.com/rs/zerolog"

	plc4go "github.com/apache/plc4x/plc4go/pkg/api"
	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/spi"
	_default "github.com/apache/plc4x/plc4go/spi/default"
	"github.com/apache/plc4x/plc4go/spi/errors"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/tracer"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

//go:generate go tool plc4xGenerator -type=Connection
type Connection struct {
	_default.DefaultConnection

	messageCodec *MessageCodec
	subscribers  []*Subscriber

	configuration Configuration `stringer:"true"`
	driverContext DriverContext `stringer:"true"`

	channel *SecureChannel

	connectEvent      chan struct{}
	connectTimeout    time.Duration // TODO: do we need to have that in general, where to get that from
	disconnectEvent   chan struct{}
	disconnectTimeout time.Duration // TODO: do we need to have that in general, where to get that from

	connectionId string
	tracer       tracer.Tracer

	wg sync.WaitGroup // use to track spawned go routines

	log      zerolog.Logger       `ignore:"true"`
	_options []options.WithOption `ignore:"true"` // Used to pass them downstream
}

var (
	_ spi.TransportInstanceExposer = (*Connection)(nil)
)

func NewConnection(messageCodec *MessageCodec, configuration Configuration, driverContext DriverContext, tagHandler spi.PlcTagHandler, connectionOptions map[string][]string, _options ...options.WithOption) *Connection {
	customLogger := options.ExtractCustomLoggerOrDefaultToGlobal(_options...)
	connection := &Connection{
		messageCodec:      messageCodec,
		configuration:     configuration,
		driverContext:     driverContext,
		channel:           NewSecureChannel(customLogger, driverContext, configuration),
		connectEvent:      make(chan struct{}),
		connectTimeout:    5 * time.Second,
		disconnectEvent:   make(chan struct{}),
		disconnectTimeout: 5 * time.Second,
		log:               customLogger,
		_options:          _options,
	}
	if traceEnabledOption, ok := connectionOptions["traceEnabled"]; ok {
		if len(traceEnabledOption) == 1 {
			connection.tracer = tracer.NewTracer(connection.connectionId, _options...)
		}
	}
	connection.DefaultConnection = _default.NewDefaultConnection(
		connection,
		append(_options,
			_default.WithPlcTagHandler(tagHandler),
		)...,
	)
	return connection
}

func (c *Connection) GetConnectionId() string {
	return c.connectionId
}

func (c *Connection) IsTraceEnabled() bool {
	return c.tracer != nil
}

func (c *Connection) GetTracer() tracer.Tracer {
	return c.tracer
}

func (c *Connection) GetConnection() plc4go.PlcConnection {
	return c
}

func (c *Connection) GetMessageCodec() spi.MessageCodec {
	return c.messageCodec
}

func (c *Connection) Connect(ctx context.Context) error {
	c.log.Trace().Msg("Connecting")
	c.log.Trace().Msg("connecting codec")
	if err := c.messageCodec.Connect(ctx); err != nil {
		return errors.Wrap(err, "Error connecting codec")
	}

	if c.driverContext.fireDiscoverEvent {
		c.log.Trace().Msg("calling onDiscover")
		c.channel.onDiscover(ctx, c.messageCodec)
	} else {
		c.log.Trace().Msg("we don't wait for session discover")
	}

	// For testing purposes we can skip the waiting for a complete connection
	if !c.driverContext.awaitSetupComplete {
		c.wg.Go(func() {
			if err := c.setupConnection(ctx); err != nil {
				c.log.Error().Err(err).Msg("Error during setup")
			}
		})
		c.log.Warn().Msg("Connection used in an unsafe way. !!!DON'T USE IN PRODUCTION!!!")
		// Here we write directly and don't wait till the connection is "really" connected
		// Note: we can't use fireConnected here as it's guarded against m.driverContext.awaitSetupComplete
		c.SetConnected(true)
		return nil
	}

	if err := c.setupConnection(ctx); err != nil {
		return errors.Wrap(err, "Error setting up connection")
	}
	return nil
}

func (c *Connection) Close() error {
	ctx := context.TODO()
	ctx, cancelFunc := utils.WithNamedTimeout(ctx, "connection close timeout", 5*time.Second)
	defer cancelFunc()

	c.channel.onDisconnect(ctx, c)
	disconnectTimeout := time.NewTimer(c.disconnectTimeout)
	select {
	case <-c.disconnectEvent:
		c.log.Info().Msg("disconnected")
	case <-disconnectTimeout.C:
		return errors.Errorf("timeout after %s", c.disconnectTimeout)
	}
	return nil
}

func (c *Connection) GetMetadata() apiModel.PlcConnectionMetadata {
	return &_default.DefaultConnectionMetadata{
		ProvidesReading:     true,
		ProvidesWriting:     true,
		ProvidesSubscribing: true,
		ProvidesBrowsing:    false,
	}
}

func (c *Connection) ReadRequestBuilder() apiModel.PlcReadRequestBuilder {
	return spiModel.NewDefaultPlcReadRequestBuilder(
		c.GetPlcTagHandler(),
		NewReader(
			c,
			append(c._options, options.WithCustomLogger(c.log))...,
		),
	)
}

func (c *Connection) WriteRequestBuilder() apiModel.PlcWriteRequestBuilder {
	return spiModel.NewDefaultPlcWriteRequestBuilder(c.GetPlcTagHandler(), c.GetPlcValueHandler(), NewWriter(c))
}

func (c *Connection) SubscriptionRequestBuilder() apiModel.PlcSubscriptionRequestBuilder {
	return spiModel.NewDefaultPlcSubscriptionRequestBuilder(
		c.GetPlcTagHandler(),
		c.GetPlcValueHandler(),
		NewSubscriber(
			c.addSubscriber,
			c,
			append(c._options, options.WithCustomLogger(c.log))...,
		),
	)
}

func (c *Connection) UnsubscriptionRequestBuilder() apiModel.PlcUnsubscriptionRequestBuilder {
	// TODO: where do we get the unsubscriber from
	return nil
}

func (c *Connection) addSubscriber(subscriber *Subscriber) {
	if slices.Contains(c.subscribers, subscriber) {
		c.log.Debug().Interface("subscriber", subscriber).Msg("Subscriber already added")
		return
	}
	c.subscribers = append(c.subscribers, subscriber)
}

func (c *Connection) setupConnection(ctx context.Context) error {
	c.log.Trace().Msg("setup connection")

	c.log.Debug().Msg("Opcua Driver running in ACTIVE mode.")
	if err := c.channel.onConnect(ctx, c); err != nil {
		return errors.Wrap(err, "Error during connection setup")
	}

	connectTimeout := time.NewTimer(c.connectTimeout)
	select {
	case <-c.connectEvent:
		c.log.Info().Msg("connected")
		c.SetConnected(true)
		c.log.Trace().Msg("Connect fired")
	case <-connectTimeout.C:
		return errors.Errorf("timeout after %s", c.connectTimeout)
	}

	c.log.Trace().Msg("connection setup done")
	return nil
}

func (c *Connection) fireConnectionError(err error, ch chan<- error) {
	c.log.Trace().Err(err).Msg("fire connection error")
	if c.driverContext.awaitSetupComplete {
		ch <- errors.Wrap(err, "Error during connection")
	} else {
		c.log.Error().Err(err).Msg("awaitSetupComplete set to false and we got a error during connect")
	}
	if err := c.messageCodec.Disconnect(); err != nil {
		c.log.Debug().Err(err).Msg("Error disconnecting message codec on connection error")
	}
	c.SetConnected(false)
	select {
	case c.disconnectEvent <- struct{}{}:
	default:
	}
}

func (c *Connection) fireConnected(ch chan<- struct{}) {
	c.log.Trace().Msg("fire connected")
	if c.driverContext.awaitSetupComplete {
		ch <- struct{}{}
	} else {
		c.log.Info().Msg("Successfully connected")
	}
	c.SetConnected(true)
	select {
	case c.connectEvent <- struct{}{}:
	default:
	}
}
