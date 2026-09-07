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

package abeth

import (
	"context"
	"fmt"
	"runtime/debug"
	"slices"
	"sync"

	"github.com/rs/zerolog"

	plc4go "github.com/apache/plc4x/plc4go/pkg/api"
	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/abeth/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	_default "github.com/apache/plc4x/plc4go/spi/default"
	"github.com/apache/plc4x/plc4go/spi/errors"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/tracer"
	"github.com/apache/plc4x/plc4go/spi/transactions"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

// connectionRequestSenderContext is the sender context the connect handshake sends, kept byte for
// byte from plc4j's AbEthConnection (which kept it from the legacy driver). The PLC echoes it back
// and nothing looks at it, but real PLCs have been observed to be picky about the handshake.
var connectionRequestSenderContext = []uint8{0x00, 0x04, 0x00, 0x05, 0x00, 0x00, 0x00, 0x00}

// Connection is a connection to an Allen-Bradley PLC speaking CIP encapsulation over TCP, the port
// of plc4j's AbEthConnection.
//
// The protocol is read-only in plc4x: after a CIPEncapsulationConnectionRequest / -Response
// handshake that yields a session handle, the only thing the driver ever sends is a DF1 "protected
// typed logical read". plc4j's AbEthDriver implements canRead and nothing else, so neither does
// this driver - what it does add on top of reading is subscriptions, which plc4j gets for free from
// PollingSubscriptionConnectionBase and this connection gets from the SPI's polling subscriber.
type Connection struct {
	_default.DefaultConnection

	configuration Configuration
	driverContext DriverContext
	messageCodec  spi.MessageCodec
	tm            transactions.RequestTransactionManager
	session       *session
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
	driverContext DriverContext,
	tagHandler spi.PlcTagHandler,
	tm transactions.RequestTransactionManager,
	connectionOptions map[string][]string,
	_options ...options.WithOption,
) *Connection {
	customLogger := options.ExtractCustomLoggerOrDefaultToGlobal(_options...)
	// Clipping the option list is what keeps the reader and the polling subscriber from racing: both
	// hand their own logger option downstream with append, and on a slice that still has spare
	// capacity those two appends write into the very same backing array.
	_options = slices.Clip(_options)
	connection := &Connection{
		configuration: configuration,
		driverContext: driverContext,
		messageCodec:  messageCodec,
		tm:            tm,
		session:       newSession(),
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

// Connect opens the transport and runs the connect handshake, which is the only way to learn the
// session handle every later packet has to carry.
func (c *Connection) Connect(ctx context.Context) error {
	c.log.Trace().Msg("Connecting")
	if err := c.messageCodec.Connect(ctx); err != nil {
		return errors.Wrap(err, "error connecting message codec")
	}
	c.startStrayMessageDrain()

	// For testing purposes we can skip waiting for a complete connection.
	if !c.driverContext.awaitSetupComplete {
		// The caller's ctx is canceled as soon as GetConnection returns, so the background
		// handshake must not inherit that cancellation.
		setupCtx := context.WithoutCancel(ctx)
		c.wg.Go(func() {
			setupCtx, cancel := utils.WithNamedTimeout(setupCtx, "abeth setup timeout", c.configuration.requestTimeout)
			defer cancel()
			if err := c.handshake(setupCtx); err != nil {
				c.log.Error().Err(err).Msg("error during the ab-eth connect handshake")
			}
		})
		c.log.Warn().Msg("Connection used in an unsafe way. !!!DON'T USE IN PRODUCTION!!!")
		// Here we mark connected without waiting for the handshake to have happened.
		c.SetConnected(true)
		return nil
	}

	if err := c.handshake(ctx); err != nil {
		// The caller gets an error instead of a connection and will never call Close on it, so the
		// drain has to be reeled back in here or its goroutine outlives the failed attempt.
		c.stopStrayMessageDrain()
		if disconnectErr := c.messageCodec.Disconnect(); disconnectErr != nil {
			c.log.Debug().Err(disconnectErr).Msg("error disconnecting after a failed handshake")
		}
		return errors.Wrap(err, "error during the ab-eth connect handshake")
	}
	return nil
}

// handshake sends the CIPEncapsulationConnectionRequest and remembers the session handle the PLC
// answers with. Ported from plc4j's AbEthConnection.doHandshake.
func (c *Connection) handshake(ctx context.Context) error {
	requestCtx, cancelRequest := context.WithTimeout(ctx, c.configuration.requestTimeout)
	defer cancelRequest()

	request := readWriteModel.NewCIPEncapsulationConnectionRequest(0, 0, connectionRequestSenderContext, 0)
	message, err := sendRequestAndWait(requestCtx, c.messageCodec, "connect", request,
		func(message spi.Message) bool {
			_, ok := message.(readWriteModel.CIPEncapsulationConnectionResponse)
			return ok
		})
	if err != nil {
		return errors.Wrap(err, "error exchanging the connection request")
	}
	response, ok := message.(readWriteModel.CIPEncapsulationConnectionResponse)
	if !ok {
		return errors.Errorf("expected a CIPEncapsulationConnectionResponse, got %T", message)
	}
	c.session.setSessionHandle(response.GetSessionHandle())
	c.log.Info().
		Uint32("sessionHandle", response.GetSessionHandle()).
		Uint32("status", response.GetStatus()).
		Msg("Connected to an ab-eth host")
	c.SetConnected(true)
	return nil
}

// Close stops the pollers and drops the transport. ab-eth has no disconnect packet - plc4j's
// AbEthConnection just closes the socket - so there is nothing to send here.
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

// startStrayMessageDrain keeps the codec's default incoming message channel empty. ab-eth is
// strictly request/response and every response is matched to its request by transaction counter, so
// anything the codec can't match - a response which arrived after its request timed out, a
// duplicate, an unsolicited packet - is pushed into that 100 slot buffer and would stay there for
// the life of the connection. Once the buffer is full the codec logs a warning per further packet,
// so on a flaky link a handful of late responses turn into permanent log noise. Draining costs one
// goroutine and turns that into a single debug line per stray packet.
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
					Msg("Discarding a packet no request was waiting for")
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
		// ab-eth is read-only in plc4x: the mspec has no write command at all, and plc4j's
		// AbEthDriver overrides canRead and nothing else. Subscribing works because it is emulated
		// by polling the read path, the same way plc4j's PollingSubscriptionConnectionBase does it.
		ProvidesReading:     true,
		ProvidesWriting:     false,
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
			c.session,
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
	return fmt.Sprintf("abeth.Connection{%s, %s}", c.configuration, c.session)
}
