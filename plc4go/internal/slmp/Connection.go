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

package slmp

import (
	"context"
	"fmt"
	"runtime/debug"
	"slices"
	"sync"

	"github.com/rs/zerolog"

	plc4go "github.com/apache/plc4x/plc4go/pkg/api"
	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/spi"
	_default "github.com/apache/plc4x/plc4go/spi/default"
	"github.com/apache/plc4x/plc4go/spi/errors"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/tracer"
	"github.com/apache/plc4x/plc4go/spi/transactions"
)

// Connection is a connection to a MELSEC-Q/L or iQ-F/FX5 device speaking SLMP 3E binary over TCP,
// the port of plc4j's SlmpConnection.
//
// SLMP has no session: there is no connect handshake and no disconnect packet, so connecting is
// opening the socket and closing is dropping it. What the connection does add on top of Batch Read
// and Batch Write is subscriptions, which plc4j gets from PollingSubscriptionConnectionBase and this
// connection gets from the SPI's polling subscriber.
type Connection struct {
	_default.DefaultConnection

	configuration Configuration
	messageCodec  spi.MessageCodec
	tm            transactions.RequestTransactionManager
	subscriber    _default.DefaultPollingSubscriber
	options       map[string][]string

	connectionId string
	tracer       tracer.Tracer

	// lifecycleMutex guards the stray-message drain's handles so Close can be called concurrently
	// with, and without, a preceding Connect.
	lifecycleMutex sync.Mutex
	drainCancel    context.CancelFunc
	drainDone      chan struct{}

	wg sync.WaitGroup // use to track spawned go routines

	log      zerolog.Logger
	_options []options.WithOption // Used to pass them downstream
}

var (
	_ spi.TransportInstanceExposer                  = (*Connection)(nil)
	_ _default.DefaultPollingSubscriberRequirements = (*Connection)(nil)
)

func NewConnection(
	messageCodec spi.MessageCodec,
	configuration Configuration,
	tagHandler spi.PlcTagHandler,
	tm transactions.RequestTransactionManager,
	connectionOptions map[string][]string,
	_options ...options.WithOption,
) *Connection {
	customLogger := options.ExtractCustomLoggerOrDefaultToGlobal(_options...)
	// Clipping the option list is what keeps the reader, the writer and the polling subscriber from
	// racing: each hands its own logger option downstream with append, and on a slice that still has
	// spare capacity those appends write into the very same backing array.
	_options = slices.Clip(_options)
	connection := &Connection{
		configuration: configuration,
		messageCodec:  messageCodec,
		tm:            tm,
		options:       connectionOptions,
		log:           customLogger,
		_options:      _options,
	}
	if traceEnabledOption, ok := connectionOptions["traceEnabled"]; ok {
		if len(traceEnabledOption) == 1 {
			connection.tracer = tracer.NewTracer(connection.connectionId, _options...)
		}
	}
	connection.DefaultConnection = _default.NewDefaultConnection(connection,
		append(_options,
			_default.WithPlcTagHandler(tagHandler),
			_default.WithPlcValueHandler(NewValueHandler(_options...)),
		)...,
	)
	// The subscriber emulates subscriptions by polling the read path, which is exactly what plc4j's
	// PollingSubscriptionConnectionBase does for this driver.
	connection.subscriber = _default.NewDefaultPollingSubscriber(connection, _options...)
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

////////////////////////////////////////////////////////////////////////////////////////////////////
// Connect / close
////////////////////////////////////////////////////////////////////////////////////////////////////

// Connect opens the transport. SLMP has no handshake - plc4j's SlmpConnection.onConnect only starts
// the receive loop - so a connected socket is a connected device as far as the protocol is
// concerned.
func (c *Connection) Connect(ctx context.Context) error {
	c.log.Trace().Msg("Connecting")
	if err := c.messageCodec.Connect(ctx); err != nil {
		return errors.Wrap(err, "error connecting message codec")
	}
	c.startStrayMessageDrain()
	c.SetConnected(true)
	c.log.Info().Msg("SLMP connection established")
	return nil
}

// Close stops the pollers and drops the transport. SLMP has no disconnect packet, so there is
// nothing to send here.
func (c *Connection) Close() error {
	c.log.Trace().Msg("Closing")
	c.stopStrayMessageDrain()
	if c.subscriber != nil {
		c.subscriber.Close()
	}
	err := c.DefaultConnection.Close()
	c.wg.Wait()
	return err
}

// startStrayMessageDrain keeps the codec's default incoming message channel empty. A 3E response
// which arrives while no request is in flight - a late answer to a request that already timed out,
// say - can't be matched to anything and would otherwise sit in that 100 slot buffer for the life of
// the connection; once the buffer is full the codec logs a warning per further frame, so on a flaky
// link a handful of late responses turn into permanent log noise.
//
// Note what the drain does *not* fix: a late response that arrives while the *next* request is in
// flight is taken as that request's answer, because a 3E frame carries no correlation id. See
// acceptsAnyResponseFrame.
func (c *Connection) startStrayMessageDrain() {
	// A second Connect on the same connection must not leave the first drain running.
	c.stopStrayMessageDrain()

	drainCtx, drainCancel := context.WithCancel(context.Background())
	drainDone := make(chan struct{})
	c.lifecycleMutex.Lock()
	c.drainCancel = drainCancel
	c.drainDone = drainDone
	c.lifecycleMutex.Unlock()

	incomingMessageChannel := c.messageCodec.GetDefaultIncomingMessageChannel()
	c.wg.Go(func() {
		// Closing this last, after the recover below, is what lets stopStrayMessageDrain wait for
		// the goroutine to really be gone.
		defer close(drainDone)
		defer func() {
			if err := recover(); err != nil {
				c.log.Error().
					Str("stack", string(debug.Stack())).
					Interface("err", err).
					Msg("panic-ed")
			}
		}()
		for {
			select {
			case <-drainCtx.Done():
				c.log.Trace().Msg("Ending the stray message drain")
				return
			case message := <-incomingMessageChannel:
				c.log.Debug().
					Type("message", message).
					Msg("Discarding a frame no request was waiting for")
			}
		}
	})
}

// stopStrayMessageDrain ends the drain and waits for its goroutine to be gone, so that a caller
// which stops the drain can rely on nothing reading the channel any more. Calling it without a
// running drain, or twice, is harmless.
func (c *Connection) stopStrayMessageDrain() {
	c.lifecycleMutex.Lock()
	drainCancel, drainDone := c.drainCancel, c.drainDone
	c.drainCancel, c.drainDone = nil, nil
	c.lifecycleMutex.Unlock()
	if drainCancel == nil {
		return
	}
	drainCancel()
	<-drainDone
}

////////////////////////////////////////////////////////////////////////////////////////////////////
// Metadata and builders
////////////////////////////////////////////////////////////////////////////////////////////////////

func (c *Connection) GetMetadata() apiModel.PlcConnectionMetadata {
	return &_default.DefaultConnectionMetadata{
		// Exactly what plc4j's SlmpConnection derives: it implements onRead and onWrite, and it
		// inherits onSubscribe from PollingSubscriptionConnectionBase, so reading, writing and
		// subscribing are advertised. There is no browser on either side.
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
			c.tm,
			c.configuration,
			append(c._options, options.WithCustomLogger(c.log))...,
		),
	)
}

func (c *Connection) WriteRequestBuilder() apiModel.PlcWriteRequestBuilder {
	return spiModel.NewDefaultPlcWriteRequestBuilder(
		c.GetPlcTagHandler(),
		c.GetPlcValueHandler(),
		NewWriter(
			c.messageCodec,
			c.tm,
			c.configuration,
			append(c._options, options.WithCustomLogger(c.log))...,
		),
	)
}

// PollingReadRequestBuilder feeds the polling subscriber, which turns every poll cycle into an
// ordinary read request.
func (c *Connection) PollingReadRequestBuilder() apiModel.PlcReadRequestBuilder {
	return c.ReadRequestBuilder()
}

func (c *Connection) SubscriptionRequestBuilder() apiModel.PlcSubscriptionRequestBuilder {
	return spiModel.NewDefaultPlcSubscriptionRequestBuilder(
		c.GetPlcTagHandler(),
		c.GetPlcValueHandler(),
		c.subscriber,
	)
}

// UnsubscriptionRequestBuilder hands out the default builder: the handles carry the subscriber they
// belong to, so nothing driver specific is needed to take them back.
func (c *Connection) UnsubscriptionRequestBuilder() apiModel.PlcUnsubscriptionRequestBuilder {
	return spiModel.NewDefaultPlcUnsubscriptionRequestBuilder()
}

func (c *Connection) String() string {
	return fmt.Sprintf("slmp.Connection{%s}", c.configuration)
}
