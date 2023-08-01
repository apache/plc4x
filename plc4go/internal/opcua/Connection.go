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
	"sync"

	"github.com/apache/plc4x/plc4go/pkg/api"
	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/default"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/tracer"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

//go:generate go run ../../tools/plc4xgenerator/gen.go -type=Connection
type Connection struct {
	_default.DefaultConnection
	messageCodec *MessageCodec
	subscribers  []*Subscriber

	configuration Configuration `stringer:"true"`
	driverContext DriverContext `stringer:"true"`

	handlerWaitGroup sync.WaitGroup

	connectionId string
	tracer       tracer.Tracer

	log      zerolog.Logger       `ignore:"true"`
	_options []options.WithOption `ignore:"true"` // Used to pass them downstream
}

func NewConnection(messageCodec *MessageCodec, configuration Configuration, driverContext DriverContext, tagHandler spi.PlcTagHandler, connectionOptions map[string][]string, _options ...options.WithOption) *Connection {
	customLogger := options.ExtractCustomLoggerOrDefaultToGlobal(_options...)
	connection := &Connection{
		messageCodec:  messageCodec,
		configuration: configuration,
		driverContext: driverContext,

		log:      customLogger,
		_options: _options,
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
	ch := make(chan plc4go.PlcConnectionConnectResult)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				c.fireConnectionError(errors.Errorf("panic-ed %v. Stack:\n%s", err, debug.Stack()), ch)
			}
		}()
		if err := c.messageCodec.ConnectWithContext(ctx); err != nil {
			c.fireConnectionError(errors.Wrap(err, "Error connecting codec"), ch)
			return
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
		c.log.Trace().Msg("Waiting for handlers to stop")
		c.handlerWaitGroup.Wait()
		c.log.Trace().Msg("handlers stopped, dispatching result")
		results <- result
	}()
	return results
}

func (c *Connection) GetMetadata() apiModel.PlcConnectionMetadata {
	return _default.DefaultConnectionMetadata{
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
			c.messageCodec,
			append(c._options, options.WithCustomLogger(c.log))...,
		),
	)
}

func (c *Connection) WriteRequestBuilder() apiModel.PlcWriteRequestBuilder {
	return spiModel.NewDefaultPlcWriteRequestBuilder(c.GetPlcTagHandler(), c.GetPlcValueHandler(), NewWriter(c.messageCodec))
}

func (c *Connection) SubscriptionRequestBuilder() apiModel.PlcSubscriptionRequestBuilder {
	return spiModel.NewDefaultPlcSubscriptionRequestBuilder(
		c.GetPlcTagHandler(),
		c.GetPlcValueHandler(),
		NewSubscriber(
			c.addSubscriber,
			c.messageCodec,
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
			c.log.Debug().Msgf("Subscriber %v already added", subscriber)
			return
		}
	}
	c.subscribers = append(c.subscribers, subscriber)
}

func (c *Connection) setupConnection(ctx context.Context, ch chan plc4go.PlcConnectionConnectResult) {
	c.log.Trace().Msg("Connection setup done")
	c.fireConnected(ch)
	c.log.Trace().Msg("Connect fired")
	c.startSubscriptionHandler()
	c.log.Trace().Msg("subscription handler started")
}

func (c *Connection) startSubscriptionHandler() {
	c.log.Debug().Msg("Starting TODO handler")
	c.handlerWaitGroup.Add(1)
	go func() {
		salLogger := c.log.With().Str("handlerType", "TODO").Logger()
		defer c.handlerWaitGroup.Done()
		defer func() {
			if err := recover(); err != nil {
				salLogger.Error().Msgf("panic-ed %v. Stack:\n%s", err, debug.Stack())
			}
		}()
		salLogger.Debug().Msg("TODO handler started")
		for c.IsConnected() {
			// TODO: dispatch subs
			/*
				for monitoredSal := range c.messageCodec.monitoredSALs {
					salLogger.Trace().Msgf("got a SAL\n%v", monitoredSal)
					handled := false
					for _, subscriber := range c.subscribers {
						if ok := subscriber.handleMonitoredSAL(monitoredSal); ok {
							salLogger.Debug().Msgf("\n%v handled\n%s", subscriber, monitoredSal)
							handled = true
						}
					}
					if !handled {
						salLogger.Debug().Msgf("SAL was not handled:\n%s", monitoredSal)
					}
				}*/
		}
		salLogger.Info().Msg("handler ended")
	}()
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
}

func (c *Connection) fireConnected(ch chan<- plc4go.PlcConnectionConnectResult) {
	c.log.Trace().Msg("fire connected")
	if c.driverContext.awaitSetupComplete {
		ch <- _default.NewDefaultPlcConnectionConnectResult(c, nil)
	} else {
		c.log.Info().Msg("Successfully connected")
	}
	c.SetConnected(true)
}
