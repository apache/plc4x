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

package iec608705104

import (
	"context"
	"fmt"
	"runtime/debug"
	"slices"
	"sync"
	"sync/atomic"

	"github.com/rs/zerolog"

	plc4go "github.com/apache/plc4x/plc4go/pkg/api"
	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/iec608705104/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	_default "github.com/apache/plc4x/plc4go/spi/default"
	"github.com/apache/plc4x/plc4go/spi/errors"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/tracer"
)

const (
	// The U-format control field values, as the mspec's type switch on APDU.command spells them.
	commandTestFrameActivation         = uint16(0x43)
	commandTestFrameConfirmation       = uint16(0x83)
	commandStopDataTransferActivation  = uint16(0x13)
	commandStartDataTransferActivation = uint16(0x07)
	commandSupervisoryFormat           = uint16(0x01)
	// sequenceNumberMask is the width of a send or receive sequence number: 15 bits, because the
	// low bit of the control field octet pair is the format discriminator.
	sequenceNumberMask = uint16(0x7FFF)
)

// Connection is a connection to an IEC 60870-5-104 controlled station (an RTU or a substation
// gateway). It is the port of plc4j's Iec60870Connection.
//
// The protocol is push driven, and that shapes everything:
//
//  1. The controlling station opens the TCP connection and exchanges a test frame with the
//     controlled station to prove the link works.
//  2. It sends STARTDT (start data transfer) and waits for the confirmation, after which the
//     controlled station is free to emit unsolicited I-format ASDUs.
//  3. Each ASDU carries one or more information objects, each of which is one point; they are
//     decoded and handed to the subscriptions whose tags cover them.
//
// There is no read path and no write path - not because they were left out, but because the
// protocol as this driver speaks it has neither: a value arrives when the station decides to send
// it. plc4j's Iec60870Connection says the same in its own javadoc ("No reads or writes are exposed -
// only subscriptions") and implements exactly that.
type Connection struct {
	_default.DefaultConnection

	configuration Configuration
	messageCodec  spi.MessageCodec
	options       map[string][]string

	connectionId string
	tracer       tracer.Tracer

	// handshakeComplete says whether STARTDT has been confirmed. Until it is, the station sends
	// nothing and the connection is not usable.
	handshakeComplete atomic.Bool

	// ackMutex guards the sequence number bookkeeping below.
	ackMutex sync.Mutex
	// unacknowledged is how many I-format frames have arrived since the last acknowledgement.
	unacknowledged uint16
	// nextExpectedSequenceNo is the send sequence number the station is expected to use next, which
	// is exactly what an S-format acknowledgement has to carry.
	nextExpectedSequenceNo uint16

	// subscribersMutex guards the subscriber list, which Subscribe writes while the
	// incoming-message worker reads it.
	subscribersMutex sync.RWMutex
	subscribers      []*Subscriber

	// lifecycleMutex guards the worker's cancel function so Close can be called concurrently with,
	// and without, a preceding Connect.
	lifecycleMutex sync.Mutex
	workerCancel   context.CancelFunc

	wg sync.WaitGroup // use to track spawned go routines

	log      zerolog.Logger
	_options []options.WithOption // Used to pass them downstream
}

var (
	_ spi.TransportInstanceExposer = (*Connection)(nil)
)

func NewConnection(configuration Configuration, messageCodec spi.MessageCodec, connectionOptions map[string][]string, tagHandler spi.PlcTagHandler, _options ...options.WithOption) *Connection {
	customLogger := options.ExtractCustomLoggerOrDefaultToGlobal(_options...)
	// Clipping the option list keeps the subscription builders from racing: each hands its own
	// logger option downstream with append, and on a slice which still has spare capacity those
	// appends write into the very same backing array.
	_options = slices.Clip(_options)
	connection := &Connection{
		configuration: configuration,
		messageCodec:  messageCodec,
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

// IsConnected reports a connection as usable only once data transfer has been started. Until STARTDT
// is confirmed the station stays silent, so a subscription registered before then would never
// deliver anything - saying "connected" would be a lie a caller cannot see through.
func (c *Connection) IsConnected() bool {
	return c.handshakeComplete.Load() && c.DefaultConnection.IsConnected()
}

////////////////////////////////////////////////////////////////////////////////////////////////////
// Connect / close
////////////////////////////////////////////////////////////////////////////////////////////////////

// Connect opens the transport, starts listening for what the station pushes and then runs the two
// handshake round trips.
func (c *Connection) Connect(ctx context.Context) error {
	c.log.Trace().Msg("Connecting")
	if err := c.DefaultConnection.Connect(ctx); err != nil {
		return errors.Wrap(err, "error connecting default connection")
	}
	c.startIncomingWorker()
	if err := c.handshake(ctx); err != nil {
		if closeErr := c.Close(); closeErr != nil {
			c.log.Debug().Err(closeErr).Msg("error closing connection after a failed handshake")
		}
		return errors.Wrap(err, "error during the IEC 60870-5-104 connect handshake")
	}
	c.handshakeComplete.Store(true)
	c.log.Info().Msg("IEC 60870-5-104 data transfer started")
	return nil
}

// handshake proves the link with a test frame and then starts data transfer, the same two round
// trips plc4j's Iec60870Connection.doHandshake runs.
func (c *Connection) handshake(ctx context.Context) error {
	if _, err := c.exchange(ctx, readWriteModel.NewAPDUUFormatTestFrameActivation(commandTestFrameActivation),
		func(message spi.Message) bool {
			_, ok := message.(readWriteModel.APDUUFormatTestFrameConfirmation)
			return ok
		}); err != nil {
		return errors.Wrap(err, "error exchanging the test frame")
	}
	if _, err := c.exchange(ctx, readWriteModel.NewAPDUUFormatStartDataTransferActivation(commandStartDataTransferActivation),
		func(message spi.Message) bool {
			_, ok := message.(readWriteModel.APDUUFormatStartDataTransferConfirmation)
			return ok
		}); err != nil {
		return errors.Wrap(err, "error starting data transfer")
	}
	return nil
}

// exchange sends one control frame and waits for the confirmation which answers it. Every path out
// of here either hands back a message or an error - a handshake step which silently returned nothing
// would leave Connect believing it succeeded.
func (c *Connection) exchange(ctx context.Context, request readWriteModel.APDU, accepts func(spi.Message) bool) (spi.Message, error) {
	requestCtx, cancelRequest := context.WithTimeout(ctx, c.configuration.requestTimeout)
	defer cancelRequest()

	responseChan := make(chan spi.Message, 1)
	errChan := make(chan error, 1)
	if err := c.messageCodec.SendRequest(requestCtx, "handshake", request,
		accepts,
		func(message spi.Message) error {
			select {
			case responseChan <- message:
			default:
				c.log.Warn().Msg("failed to hand on the handshake response")
			}
			return nil
		},
		func(err error) error {
			select {
			case errChan <- err:
			default:
				c.log.Warn().Msg("failed to hand on the handshake error")
			}
			return nil
		}); err != nil {
		return nil, errors.Wrapf(err, "error sending %T", request)
	}

	select {
	case response := <-responseChan:
		return response, nil
	case err := <-errChan:
		return nil, errors.Wrap(err, "error waiting for the confirmation")
	case <-requestCtx.Done():
		return nil, requestCtx.Err()
	}
}

func (c *Connection) Close() error {
	c.log.Trace().Msg("Closing")
	// Tell the station to stop sending before the transport goes away. A station whose peer just
	// vanishes keeps the session (and its send buffer) alive until its own t1 timer fires; STOPDT is
	// the orderly way out, and it costs one frame. plc4j never sends it. It is best effort by
	// design: a connection which is already broken is exactly the case where this fails.
	if c.handshakeComplete.Swap(false) {
		stopCtx, cancelStop := context.WithTimeout(context.Background(), c.configuration.requestTimeout)
		stopDataTransfer := readWriteModel.NewAPDUUFormatStopDataTransferActivation(commandStopDataTransferActivation)
		if err := c.messageCodec.Send(stopCtx, "close", stopDataTransfer); err != nil {
			c.log.Debug().Err(err).Msg("error sending the stop-data-transfer activation")
		}
		cancelStop()
	}

	c.lifecycleMutex.Lock()
	workerCancel := c.workerCancel
	c.workerCancel = nil
	c.lifecycleMutex.Unlock()
	if workerCancel != nil {
		workerCancel()
	}
	err := c.DefaultConnection.Close()
	c.wg.Wait()

	c.subscribersMutex.Lock()
	subscribers := c.subscribers
	c.subscribers = nil
	c.subscribersMutex.Unlock()
	for _, subscriber := range subscribers {
		subscriber.Close()
	}
	return err
}

// startIncomingWorker pumps the messages the codec couldn't match to an expectation - which after
// the handshake is every frame the station sends - into the acknowledgement bookkeeping and on to
// the subscribers.
func (c *Connection) startIncomingWorker() {
	workerCtx, workerCancel := context.WithCancel(context.Background())
	c.lifecycleMutex.Lock()
	c.workerCancel = workerCancel
	c.lifecycleMutex.Unlock()

	incomingMessageChannel := c.messageCodec.GetDefaultIncomingMessageChannel()
	c.wg.Go(func() {
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
			case <-workerCtx.Done():
				c.log.Info().Msg("Ending incoming message transfer")
				return
			case message := <-incomingMessageChannel:
				c.handleIncomingMessage(workerCtx, message)
			}
		}
	})
}

////////////////////////////////////////////////////////////////////////////////////////////////////
// Incoming messages
////////////////////////////////////////////////////////////////////////////////////////////////////

// handleIncomingMessage sorts one frame the station sent.
func (c *Connection) handleIncomingMessage(ctx context.Context, message spi.Message) {
	switch typed := message.(type) {
	case readWriteModel.APDUIFormat:
		c.handleDataFrame(ctx, typed)
	case readWriteModel.APDUUFormatTestFrameActivation:
		// The station probes the link whenever its own idle timer fires. Answering is not optional:
		// a station which doesn't get its confirmation closes the connection.
		confirmation := readWriteModel.NewAPDUUFormatTestFrameConfirmation(commandTestFrameConfirmation)
		if err := c.messageCodec.Send(ctx, "testFrame", confirmation); err != nil {
			c.log.Error().Err(err).Msg("Failed to confirm a test frame")
		}
	case readWriteModel.APDUSFormat:
		// The station acknowledging frames we sent. This driver sends no I-format frames at all, so
		// there is nothing outstanding for it to acknowledge.
		c.log.Trace().Uint16("receiveSequenceNo", typed.GetReceiveSequenceNo()>>1).Msg("Supervisory frame")
	default:
		c.log.Debug().Type("message", message).Msg("Unexpected IEC 60870-5-104 frame")
	}
}

// handleDataFrame acknowledges an I-format frame when one is due and hands its ASDU to the
// subscribers.
func (c *Connection) handleDataFrame(ctx context.Context, apdu readWriteModel.APDUIFormat) {
	// The send sequence number of an I-format frame lives in the control field the mspec models as
	// 'command', shifted left by one because the low bit is the format discriminator.
	sendSequenceNo := (apdu.GetCommand() >> 1) & sequenceNumberMask
	if nextExpected, due := c.noteReceivedDataFrame(sendSequenceNo); due {
		// An S-format acknowledgement carries the sequence number of the frame we expect next, in
		// the same shifted-by-one encoding.
		//
		// plc4j gets this wrong twice over: it answers with the *receive* sequence number of the
		// incoming frame (which is the station acknowledging our sends, nothing to do with what we
		// have received) and it doesn't shift, so the number it sends is off by a factor of two.
		// A station which checks its acknowledgements - and it must, that is what the window is for -
		// stops sending after k frames and eventually drops the connection.
		acknowledgement := readWriteModel.NewAPDUSFormat(commandSupervisoryFormat, nextExpected<<1)
		if err := c.messageCodec.Send(ctx, "acknowledge", acknowledgement); err != nil {
			c.log.Error().Err(err).Msg("Failed to acknowledge received data frames")
		}
	}

	asdu := apdu.GetAsdu()
	if asdu == nil {
		c.log.Debug().Msg("An I-format frame without an ASDU")
		return
	}
	if !decodableAsdu(asdu) {
		// A sequence-of-objects ASDU (structure qualifier set) carries one information object
		// address followed by that many payloads with no addresses of their own. The mspec reads an
		// address in front of every object regardless, so the generated model lines the payloads up
		// against the wrong addresses - and where the arithmetic happens to work out it does so
		// without any error at all. Publishing that would put a reading of one point under the
		// address of another, which in grid telemetry is worse than publishing nothing, so the whole
		// ASDU is refused. The frame is still acknowledged: it did arrive.
		c.log.Warn().
			Uint16("asduAddress", asdu.GetAsduAddressField()).
			Stringer("typeIdentification", asdu.GetTypeIdentification()).
			Int("numberOfObjects", len(asdu.GetInformationObjects())).
			Msg("Dropping a sequence-of-objects ASDU: the generated model cannot address its objects")
		return
	}
	for _, subscriber := range c.activeSubscribers() {
		subscriber.publish(asdu)
	}
}

// decodableAsdu reports whether the objects of an ASDU can be trusted to belong to the addresses the
// model read them under. With the structure qualifier clear every object carries its own address, so
// they always can; with it set they only can when there is a single object, which is the one case in
// which both layouts are the same bytes.
func decodableAsdu(asdu readWriteModel.ASDU) bool {
	return !asdu.GetStructureQualifier() || len(asdu.GetInformationObjects()) <= 1
}

// noteReceivedDataFrame records an I-format frame and reports the sequence number to acknowledge
// plus whether an acknowledgement is due. The window is the standard's 'w' parameter.
func (c *Connection) noteReceivedDataFrame(sendSequenceNo uint16) (uint16, bool) {
	c.ackMutex.Lock()
	defer c.ackMutex.Unlock()
	c.nextExpectedSequenceNo = (sendSequenceNo + 1) & sequenceNumberMask
	c.unacknowledged++
	if c.unacknowledged < c.configuration.ackThreshold {
		return c.nextExpectedSequenceNo, false
	}
	c.unacknowledged = 0
	return c.nextExpectedSequenceNo, true
}

////////////////////////////////////////////////////////////////////////////////////////////////////
// Subscribers
////////////////////////////////////////////////////////////////////////////////////////////////////

func (c *Connection) addSubscriber(subscriber *Subscriber) {
	c.subscribersMutex.Lock()
	defer c.subscribersMutex.Unlock()
	if slices.Contains(c.subscribers, subscriber) {
		return
	}
	c.subscribers = append(c.subscribers, subscriber)
}

func (c *Connection) removeSubscriber(subscriber *Subscriber) {
	c.subscribersMutex.Lock()
	defer c.subscribersMutex.Unlock()
	for i, existing := range c.subscribers {
		if existing == subscriber {
			c.subscribers = append(c.subscribers[:i], c.subscribers[i+1:]...)
			return
		}
	}
}

// activeSubscribers snapshots the subscriber list so events can be delivered without holding the
// lock while calling into consumer code, which may subscribe or unsubscribe from inside a callback.
func (c *Connection) activeSubscribers() []*Subscriber {
	c.subscribersMutex.RLock()
	defer c.subscribersMutex.RUnlock()
	return append([]*Subscriber(nil), c.subscribers...)
}

////////////////////////////////////////////////////////////////////////////////////////////////////
// Metadata and builders
////////////////////////////////////////////////////////////////////////////////////////////////////

// GetMetadata advertises subscribing and nothing else, which is all this protocol offers: the
// controlled station pushes what its configuration says it should, and there is no request a
// controlling station can make for the current value of a point (nor, in this driver, any command it
// can send). plc4j's Iec60870514PlcDriver says the same by overriding canSubscribe alone.
func (c *Connection) GetMetadata() apiModel.PlcConnectionMetadata {
	return &_default.DefaultConnectionMetadata{
		ProvidesReading:     false,
		ProvidesWriting:     false,
		ProvidesSubscribing: true,
		ProvidesBrowsing:    false,
	}
}

func (c *Connection) SubscriptionRequestBuilder() apiModel.PlcSubscriptionRequestBuilder {
	return spiModel.NewDefaultPlcSubscriptionRequestBuilder(
		c.GetPlcTagHandler(),
		c.GetPlcValueHandler(),
		NewSubscriber(c, append(c._options, options.WithCustomLogger(c.log))...),
	)
}

// UnsubscriptionRequestBuilder hands out the default builder: the handles carry the subscriber they
// belong to, so nothing driver specific is needed to take them back.
func (c *Connection) UnsubscriptionRequestBuilder() apiModel.PlcUnsubscriptionRequestBuilder {
	return spiModel.NewDefaultPlcUnsubscriptionRequestBuilder()
}

func (c *Connection) String() string {
	return fmt.Sprintf("iec608705104.Connection{requestTimeout: %s, ackThreshold: %d}",
		c.configuration.requestTimeout, c.configuration.ackThreshold)
}
