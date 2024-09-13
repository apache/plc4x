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
	"runtime/debug"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/apache/plc4x/plc4go/pkg/api"
	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/default"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/tracer"
)

//go:generate plc4xGenerator -type=Connection
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

	log      zerolog.Logger       `ignore:"true"`
	_options []options.WithOption `ignore:"true"` // Used to pass them downstream
}

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

func (c *Connection) ConnectWithContext(ctx context.Context) <-chan plc4go.PlcConnectionConnectResult {
	c.log.Trace().Msg("Connecting")
	ch := make(chan plc4go.PlcConnectionConnectResult, 1)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				c.fireConnectionError(errors.Errorf("panic-ed %v. Stack:\n%s", err, debug.Stack()), ch)
			}
		}()
		c.log.Trace().Msg("connecting codec")
		if err := c.messageCodec.ConnectWithContext(ctx); err != nil {
			c.fireConnectionError(errors.Wrap(err, "Error connecting codec"), ch)
			return
		}

		if c.driverContext.fireDiscoverEvent {
			c.log.Trace().Msg("calling onDiscover")
			c.channel.onDiscover(ctx, c.messageCodec)
		} else {
			c.log.Trace().Msg("we don't wait for session discover")
		}

		// For testing purposes we can skip the waiting for a complete connection
		if !c.driverContext.awaitSetupComplete {
			go c.setupConnection(ctx, ch)
			c.log.Warn().Msg("Connection used in an unsafe way. !!!DON'T USE IN PRODUCTION!!!")
			// Here we write directly and don't wait till the connection is "really" connected
			// Note: we can't use fireConnected here as it's guarded against m.driverContext.awaitSetupComplete
			ch <- _default.NewDefaultPlcConnectionConnectResult(c, nil)
			c.SetConnected(true)
			return
		}

		c.setupConnection(ctx, ch)
	}()
	return ch
}

func (c *Connection) Close() <-chan plc4go.PlcConnectionCloseResult {
	results := make(chan plc4go.PlcConnectionCloseResult, 1)
	go func() {
		result := <-c.DefaultConnection.Close()
		c.channel.onDisconnect(context.Background(), c)
		disconnectTimeout := time.NewTimer(c.disconnectTimeout)
		select {
		case <-c.disconnectEvent:
			c.log.Info().Msg("disconnected")
			results <- result
		case <-disconnectTimeout.C:
			results <- _default.NewDefaultPlcConnectionCloseResult(c, errors.Errorf("timeout after %s", c.disconnectTimeout))
		}
	}()
	return results
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
	for _, sub := range c.subscribers {
		if sub == subscriber {
			c.log.Debug().Stringer("subscriber", subscriber).Msg("Subscriber already added")
			return
		}
	}
	c.subscribers = append(c.subscribers, subscriber)
}

func (c *Connection) setupConnection(ctx context.Context, ch chan plc4go.PlcConnectionConnectResult) {
	c.log.Trace().Msg("setup connection")

	c.log.Debug().Msg("Opcua Driver running in ACTIVE mode.")
	c.channel.onConnect(ctx, c, ch)

	connectTimeout := time.NewTimer(c.connectTimeout)
	select {
	case <-c.connectEvent:
		c.log.Info().Msg("connected")
		c.fireConnected(ch)
		c.log.Trace().Msg("Connect fired")
	case <-connectTimeout.C:
		c.fireConnectionError(errors.Errorf("timeout after %s", c.connectTimeout), ch)
		c.log.Trace().Msg("connection error fired")
		return
	}

	c.log.Trace().Msg("connection setup done")
}

func (c *Connection) fireConnectionError(err error, ch chan<- plc4go.PlcConnectionConnectResult) {
	c.log.Trace().Err(err).Msg("fire connection error")
	if c.driverContext.awaitSetupComplete {
		ch <- _default.NewDefaultPlcConnectionConnectResult(nil, errors.Wrap(err, "Error during connection"))
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

func (c *Connection) fireConnected(ch chan<- plc4go.PlcConnectionConnectResult) {
	c.log.Trace().Msg("fire connected")
	if c.driverContext.awaitSetupComplete {
		ch <- _default.NewDefaultPlcConnectionConnectResult(c, nil)
	} else {
		c.log.Info().Msg("Successfully connected")
	}
	c.SetConnected(true)
	select {
	case c.connectEvent <- struct{}{}:
	default:
	}
}
