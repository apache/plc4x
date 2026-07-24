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

package bacnetip

import (
	"context"
	"fmt"
	"runtime/debug"
	"slices"
	"sync"
	"time"

	"github.com/rs/zerolog"

	plc4go "github.com/apache/plc4x/plc4go/pkg/api"
	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	_default "github.com/apache/plc4x/plc4go/spi/default"
	"github.com/apache/plc4x/plc4go/spi/errors"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/tracer"
	"github.com/apache/plc4x/plc4go/spi/transactions"
)

type Connection struct {
	_default.DefaultConnection

	invokeIdGenerator InvokeIdGenerator
	messageCodec      spi.MessageCodec
	configuration     Configuration
	routedDest        *routedDestination // non-nil when the target device is behind a BACnet router
	driverContext     DriverContext
	subscribers       []*Subscriber
	subscribersMu     sync.Mutex
	tm                transactions.RequestTransactionManager

	connectionId string
	tracer       tracer.Tracer

	wg sync.WaitGroup // use to track spawned go routines

	log      zerolog.Logger
	_options []options.WithOption // Used to pass them downstream
}

var (
	_ spi.TransportInstanceExposer = (*Connection)(nil)
)

func NewConnection(messageCodec spi.MessageCodec, tagHandler spi.PlcTagHandler, tm transactions.RequestTransactionManager, connectionOptions map[string][]string, _options ...options.WithOption) *Connection {
	customLogger := options.ExtractCustomLoggerOrDefaultToGlobal(_options...)
	configuration, err := ParseFromOptions(customLogger, connectionOptions)
	if err != nil {
		customLogger.Warn().Err(err).Msg("invalid driver options; falling back to defaults")
		configuration = createDefaultConfiguration()
	}
	routedDest, err := routedDestinationFromConfiguration(configuration)
	if err != nil {
		// Fail closed on the routing options: silently ignoring them would
		// unicast routed requests at the router without a destination
		// specifier, which the router (correctly) cannot forward.
		customLogger.Error().Err(err).Msg("invalid routed-destination options; connection will address the local segment only")
		routedDest = nil
	}
	connection := &Connection{
		invokeIdGenerator: InvokeIdGenerator{currentInvokeId: 0},
		messageCodec:      messageCodec,
		configuration:     configuration,
		driverContext:     NewDriverContext(configuration),
		routedDest:        routedDest,
		tm:                tm,
		log:               customLogger,
		_options:          _options,
	}
	if traceEnabledOption, ok := connectionOptions["traceEnabled"]; ok {
		if len(traceEnabledOption) == 1 {
			connection.tracer = tracer.NewTracer(connection.connectionId, _options...)
		}
	}
	connection.DefaultConnection = _default.NewDefaultConnection(connection,
		_default.WithPlcTagHandler(tagHandler),
		_default.WithPlcValueHandler(NewValueHandler()),
		options.WithCustomLogger(customLogger),
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

func (c *Connection) Connect(ctx context.Context) error {
	c.log.Trace().Msg("Connecting")
	if err := c.DefaultConnection.Connect(ctx); err != nil {
		return errors.Wrap(err, "Error connecting default connection")
	}
	c.wg.Go(func() {
		defer func() {
			if err := recover(); err != nil {
				c.log.Error().
					Str("stack", string(debug.Stack())).
					Interface("err", err).
					Msg("panic-ed")
			}
		}()
		for c.IsConnected() {
			c.log.Trace().Msg("Polling data")
			c.passToDefaultIncomingMessageChannel()
		}
		c.log.Info().Msg("Ending incoming message transfer")
	})
	return nil
}

func (c *Connection) Close() error {
	err := c.DefaultConnection.Close()
	c.wg.Wait()
	return err
}

func (c *Connection) passToDefaultIncomingMessageChannel() {
	incomingMessageChannel := c.messageCodec.GetDefaultIncomingMessageChannel()
	// Block (with a short timeout so the Connect loop can re-check IsConnected
	// for shutdown) rather than busy-spinning with a default case. The previous
	// non-blocking select pegged a CPU core and starved the codec's receive
	// worker, so request responses were never read from the socket.
	select {
	case message := <-incomingMessageChannel:
		c.routeIncomingMessage(message)
	case <-time.After(100 * time.Millisecond):
		c.log.Trace().Msg("no incoming message")
	}
}

// routeIncomingMessage inspects an unsolicited (non-request-matched) BVLC
// frame and dispatches it to the appropriate handler. The only message we
// route today is COV (Confirmed/Unconfirmed) — everything else is logged at
// debug level for visibility.
func (c *Connection) routeIncomingMessage(message spi.Message) {
	bvlc, ok := message.(readWriteModel.BVLC)
	if !ok {
		c.log.Debug().Type("message", message).Msg("non-BVLC incoming message")
		return
	}
	npduRetriever, ok := bvlc.(interface{ GetNpdu() readWriteModel.NPDU })
	if !ok {
		c.log.Debug().Msg("BVLC without an NPDU")
		return
	}
	apdu := npduRetriever.GetNpdu().GetApdu()
	switch apdu := apdu.(type) {
	case readWriteModel.APDUUnconfirmedRequest:
		switch sr := apdu.GetServiceRequest().(type) {
		case readWriteModel.BACnetUnconfirmedServiceRequestUnconfirmedCOVNotification:
			for _, s := range c.subscribers {
				s.HandleUnconfirmedCOVNotification(sr)
			}
		default:
			c.log.Debug().Type("serviceRequest", sr).Msg("unhandled unconfirmed service request")
		}
	case readWriteModel.APDUConfirmedRequest:
		switch sr := apdu.GetServiceRequest().(type) {
		case readWriteModel.BACnetConfirmedServiceRequestConfirmedCOVNotification:
			for _, s := range c.subscribers {
				s.HandleConfirmedCOVNotification(sr)
			}
			// TODO: send APDUSimpleAck back to the publisher so it doesn't retry.
		default:
			c.log.Debug().Type("serviceRequest", sr).Msg("unhandled confirmed service request")
		}
	default:
		c.log.Debug().Type("apdu", apdu).Msg("unhandled APDU")
	}
}

func (c *Connection) GetConnection() plc4go.PlcConnection {
	return c
}

func (c *Connection) GetMessageCodec() spi.MessageCodec {
	return c.messageCodec
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
			&c.invokeIdGenerator,
			c.messageCodec,
			c.tm,
			c.driverContext,
			c.routedDest,
			append(c._options, options.WithCustomLogger(c.log))...,
		),
	)
}

func (c *Connection) WriteRequestBuilder() apiModel.PlcWriteRequestBuilder {
	return spiModel.NewDefaultPlcWriteRequestBuilder(
		c.GetPlcTagHandler(),
		c.GetPlcValueHandler(),
		NewWriter(
			&c.invokeIdGenerator,
			c.messageCodec,
			c.tm,
			c.driverContext,
			c.routedDest,
			append(c._options, options.WithCustomLogger(c.log))...,
		),
	)
}

func (c *Connection) SubscriptionRequestBuilder() apiModel.PlcSubscriptionRequestBuilder {
	return spiModel.NewDefaultPlcSubscriptionRequestBuilder(
		c.GetPlcTagHandler(),
		c.GetPlcValueHandler(),
		NewSubscriber(
			c,
			append(c._options, options.WithCustomLogger(c.log))...,
		),
	)
}

func (c *Connection) UnsubscriptionRequestBuilder() apiModel.PlcUnsubscriptionRequestBuilder {
	// The default request implementation dispatches each handle's
	// Unsubscribe back through the embedded Subscriber, so we don't need
	// to pass our own here — the SubscriptionHandles created by Subscribe
	// already carry the Subscriber reference.
	return spiModel.NewDefaultPlcUnsubscriptionRequestBuilder()
}

func (c *Connection) addSubscriber(subscriber *Subscriber) {
	c.subscribersMu.Lock()
	defer c.subscribersMu.Unlock()
	if slices.Contains(c.subscribers, subscriber) {
		c.log.Debug().Interface("subscriber", subscriber).Msg("Subscriber already added")
		return
	}
	c.subscribers = append(c.subscribers, subscriber)
}

// ActiveSubscriptionCount reports how many COV subscription handles are
// currently registered across this connection's subscribers. The connection
// cache consults it (as an optional capability) to exempt subscription-
// carrying connections from idle reaping: their server-side COV
// registrations and refresh timers live on this connection and would be
// destroyed by a reap.
func (c *Connection) ActiveSubscriptionCount() int {
	c.subscribersMu.Lock()
	defer c.subscribersMu.Unlock()
	count := 0
	for _, s := range c.subscribers {
		count += s.activeHandleCount()
	}
	return count
}

func (c *Connection) String() string {
	return fmt.Sprintf("bacnetip.Connection")
}

type InvokeIdGenerator struct {
	currentInvokeId uint8
	lock            sync.Mutex
}

func (t *InvokeIdGenerator) getAndIncrement() uint8 {
	t.lock.Lock()
	defer t.lock.Unlock()
	result := t.currentInvokeId
	t.currentInvokeId += 1
	return result
}
