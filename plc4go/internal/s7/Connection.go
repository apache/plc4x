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

package s7

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/rs/zerolog"

	"github.com/apache/plc4x/plc4go/pkg/api"
	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/s7/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/default"
	"github.com/apache/plc4x/plc4go/spi/errors"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/tracer"
	"github.com/apache/plc4x/plc4go/spi/transactions"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

type TpduGenerator struct {
	currentTpduId uint16
	lock          sync.Mutex
}

func (t *TpduGenerator) getAndIncrement() uint16 {
	t.lock.Lock()
	defer t.lock.Unlock()
	// If we've reached the max value for a 16 bit transaction identifier, reset back to 1
	if t.currentTpduId >= 0xFFFF {
		t.currentTpduId = 1
	}
	result := t.currentTpduId
	t.currentTpduId += 1
	return result
}

type Connection struct {
	_default.DefaultConnection

	tpduGenerator TpduGenerator
	messageCodec  spi.MessageCodec
	configuration Configuration
	driverContext DriverContext
	tm            transactions.RequestTransactionManager

	// subscriber is created lazily and shared: the alarm subscription and cyclic jobs are
	// PLC-side state bound to this connection.
	subscriber      *Subscriber
	subscriberMutex sync.Mutex

	connectionId string
	tracer       tracer.Tracer

	wg sync.WaitGroup // use to track spawned go routines

	log      zerolog.Logger
	_options []options.WithOption // Used to pass them downstream
}

var (
	_ spi.TransportInstanceExposer = (*Connection)(nil)
)

func NewConnection(messageCodec spi.MessageCodec, configuration Configuration, driverContext DriverContext, tagHandler spi.PlcTagHandler, tm transactions.RequestTransactionManager, connectionOptions map[string][]string, _options ...options.WithOption) *Connection {
	customLogger := options.ExtractCustomLoggerOrDefaultToGlobal(_options...)
	connection := &Connection{
		tpduGenerator: TpduGenerator{currentTpduId: 10},
		messageCodec:  messageCodec,
		configuration: configuration,
		driverContext: driverContext,
		tm:            tm,
		log:           customLogger,
		_options:      _options,
	}
	if traceEnabledOption, ok := connectionOptions["traceEnabled"]; ok {
		if len(traceEnabledOption) == 1 {
			connection.tracer = tracer.NewTracer(connection.connectionId, _options...)
		}
	}
	connection.DefaultConnection = _default.NewDefaultConnection(connection,
		_default.WithPlcTagHandler(tagHandler),
		_default.WithPlcValueHandler(NewValueHandler()),
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
	err := c.messageCodec.Connect(ctx)
	if err != nil {
		return errors.Wrap(err, "Error during message codec setup")
	}

	// Route unsolicited UserData pushes (cyclic data, alarm indications) to the subscriber.
	// The channel is closed when the codec disconnects, which ends this goroutine.
	if codec, ok := c.messageCodec.(*MessageCodec); ok {
		c.wg.Go(func() {
			c.dispatchUnsolicitedUserData(codec.GetUnsolicitedUserData())
		})
	}

	// Only on active connections we do a connection
	if c.driverContext.PassiveMode {
		c.log.Info().Msg("S7 Driver running in PASSIVE mode.")
		return nil
	}

	// For testing purposes we can skip the waiting for a complete connection
	if !c.driverContext.awaitSetupComplete {
		c.wg.Go(func() {
			if err := c.setupConnection(ctx); err != nil {
				c.log.Error().Err(err).Msg("Error during connection setup")
			}
		})
		c.log.Warn().Msg("Connection used in an unsafe way. !!!DON'T USE IN PRODUCTION!!!")
		// Here we write directly and don't wait till the connection is "really" connected
		// Note: we can't use fireConnected here as it's guarded against c.driverContext.awaitSetupComplete
		c.SetConnected(true)
		return nil
	}

	// Only the TCP transport supports login.
	c.log.Info().Msg("S7 Driver running in ACTIVE mode.")

	if err := c.setupConnection(ctx); err != nil {
		return errors.Wrap(err, "Error during connection setup")
	}
	return nil
}

func (c *Connection) setupConnection(ctx context.Context) error {
	c.log.Debug().Msg("Sending COTP Connection Request")
	// Open the session on ISO Transport Protocol first.
	cotpConnectionResult := make(chan readWriteModel.COTPPacketConnectionResponse, 1)
	cotpConnectionErrorChan := make(chan error, 1)
	if err := c.messageCodec.SendRequest(ctx, "setup_connection", readWriteModel.NewTPKTPacket(c.createCOTPConnectionRequest()), func(message spi.Message) bool {
		tpktPacket := message.(readWriteModel.TPKTPacket)
		if tpktPacket == nil {
			return false
		}
		cotpPacketConnectionResponse := tpktPacket.GetPayload().(readWriteModel.COTPPacketConnectionResponse)
		return cotpPacketConnectionResponse != nil
	}, func(message spi.Message) error {
		tpktPacket := message.(readWriteModel.TPKTPacket)
		cotpPacketConnectionResponse := tpktPacket.GetPayload().(readWriteModel.COTPPacketConnectionResponse)
		cotpConnectionResult <- cotpPacketConnectionResponse
		return nil
	}, func(err error) error {
		// If this is a timeout, do a check if the connection requires a reconnection
		var timeoutError utils.TimeoutError
		if errors.As(err, &timeoutError) {
			c.log.Warn().Msg("Timeout during Connection establishing, closing channel...")
			c.Close()
		}
		cotpConnectionErrorChan <- errors.Wrap(err, "got error processing request")
		return nil
	}); err != nil {
		return errors.Wrap(err, "Error during sending of COTP Connection Request")
	}
	select {
	case cotpPacketConnectionResponse := <-cotpConnectionResult:
		c.log.Debug().Msg("Got COTP Connection Response")
		c.log.Debug().Msg("Sending S7 Connection Request")

		// Send an S7 login message.
		s7ConnectionResult := make(chan readWriteModel.S7ParameterSetupCommunication, 1)
		s7ConnectionErrorChan := make(chan error, 1)
		if err := c.messageCodec.SendRequest(ctx, "setup_connection_connection_request", c.createS7ConnectionRequest(cotpPacketConnectionResponse), func(message spi.Message) bool {
			tpktPacket, ok := message.(readWriteModel.TPKTPacket)
			if !ok {
				return false
			}
			cotpPacketData, ok := tpktPacket.GetPayload().(readWriteModel.COTPPacketData)
			if !ok {
				return false
			}
			messageResponseData, ok := cotpPacketData.GetPayload().(readWriteModel.S7MessageResponseData)
			if !ok {
				return false
			}
			_, ok = messageResponseData.GetParameter().(readWriteModel.S7ParameterSetupCommunication)
			return ok
		}, func(message spi.Message) error {
			tpktPacket := message.(readWriteModel.TPKTPacket)
			cotpPacketData := tpktPacket.GetPayload().(readWriteModel.COTPPacketData)
			messageResponseData := cotpPacketData.GetPayload().(readWriteModel.S7MessageResponseData)
			setupCommunication := messageResponseData.GetParameter().(readWriteModel.S7ParameterSetupCommunication)
			s7ConnectionResult <- setupCommunication
			return nil
		}, func(err error) error {
			// If this is a timeout, do a check if the connection requires a reconnection
			var timeoutError utils.TimeoutError
			if errors.As(err, &timeoutError) {
				c.log.Warn().Msg("Timeout during Connection establishing, closing channel...")
				if err := c.Close(); err != nil {
					c.log.Error().Err(err).Msg("Error during closing of connection")
				}
			}
			s7ConnectionErrorChan <- errors.Wrap(err, "got error processing request")
			return nil
		}); err != nil {
			return errors.Wrap(err, "Error during sending of S7 Connection Request")
		}
		select {
		case setupCommunication := <-s7ConnectionResult:
			c.log.Debug().Msg("Got S7 Connection Response")
			c.log.Debug().Msg("Sending identify remote Request")
			// Save some data from the response.
			c.driverContext.MaxAmqCaller = setupCommunication.GetMaxAmqCaller()
			c.driverContext.MaxAmqCallee = setupCommunication.GetMaxAmqCallee()
			c.driverContext.PduSize = setupCommunication.GetPduLength()

			// Update the number of concurrent requests to the negotiated number.
			// I have never seen anything else than equal values for caller and
			// callee, but if they were different, we're only limiting the outgoing
			// requests.
			c.tm.SetNumberOfConcurrentRequests(int(c.driverContext.MaxAmqCallee))

			// Detect the controller type (unless it was pinned by the user) and derive
			// whether the device supports the S7Comm UserData services. Probe failures are
			// not connection failures - plain read/write stays usable (e.g. LOGO devices).
			c.identifyRemote(ctx)

			// Send an event that connection setup is complete.
			c.SetConnected(true)
			return nil
		case err := <-s7ConnectionErrorChan:
			return errors.Wrap(err, "Error receiving S7 Connection Request")
		}
	case err := <-cotpConnectionErrorChan:
		return errors.Wrap(err, "Error receiving of COTP Connection Request")
	}
}

// identifyRemote determines the controller type and UserData-service capability. If the user
// pinned the controller type it is trusted and no SZL probe runs; otherwise the modern
// COMPONENT_IDENTIFICATION SZL is tried first with the legacy MODULE_IDENTIFICATION as
// fallback. Failure of both just disables the UserData services (browse/alarms/cyclic).
func (c *Connection) identifyRemote(ctx context.Context) {
	if c.driverContext.ControllerType != readWriteModel.ControllerType_ANY {
		c.driverContext.UserDataServicesSupported = supportsUserDataServices(c.driverContext.ControllerType)
		c.log.Info().
			Stringer("controllerType", c.driverContext.ControllerType).
			Bool("userDataServices", c.driverContext.UserDataServicesSupported).
			Msg("Controller type pinned, skipping SZL probe")
		return
	}
	c.log.Debug().Msg("Sending S7 Identification Request")
	articleNumber, controllerType, err := c.trySzlProbe(ctx, szlIdComponentIdentification(), 0x0001)
	if err != nil {
		c.log.Debug().Err(err).Msg("COMPONENT_IDENTIFICATION SZL probe failed, trying MODULE_IDENTIFICATION")
		articleNumber, controllerType, err = c.trySzlProbe(ctx, szlIdModuleIdentification(), 0x0000)
	}
	if err != nil {
		// Common on LOGO and other non-SZL devices.
		c.log.Info().Err(err).Msg("SZL probe yielded no usable identification; " +
			"S7Comm UserData services disabled for this device")
		c.driverContext.UserDataServicesSupported = false
		return
	}
	c.driverContext.ArticleNumber = articleNumber
	c.driverContext.ControllerType = controllerType
	c.driverContext.UserDataServicesSupported = true
	c.log.Info().
		Str("articleNumber", articleNumber).
		Stringer("controllerType", controllerType).
		Msg("SZL probe ok")
}

func (c *Connection) trySzlProbe(ctx context.Context, szlId readWriteModel.SzlId, szlIndex uint16) (string, readWriteModel.ControllerType, error) {
	tpduId := c.tpduGenerator.getAndIncrement()
	message, err := c.sendUserData(ctx, tpduId, buildSzlRequest(tpduId, szlId, szlIndex), "setup_connection_identify_remote_message")
	if err != nil {
		return "", 0, err
	}
	return parseSzlProbeResponse(message)
}

// sendUserData sends a single S7 UserData request and waits (blocking) for the response
// with the matching tpdu reference.
func (c *Connection) sendUserData(ctx context.Context, tpduId uint16, request readWriteModel.TPKTPacket, interactionInfo string) (readWriteModel.S7Message, error) {
	resultChan := make(chan readWriteModel.S7Message, 1)
	errChan := make(chan error, 1)
	if err := c.messageCodec.SendRequest(ctx, interactionInfo, request, func(message spi.Message) bool {
		tpktPacket, ok := message.(readWriteModel.TPKTPacket)
		if !ok {
			return false
		}
		cotpPacketData, ok := tpktPacket.GetPayload().(readWriteModel.COTPPacketData)
		if !ok {
			return false
		}
		messageUserData, ok := cotpPacketData.GetPayload().(readWriteModel.S7MessageUserData)
		if !ok {
			return false
		}
		return messageUserData.GetTpduReference() == tpduId
	}, func(message spi.Message) error {
		tpktPacket := message.(readWriteModel.TPKTPacket)
		cotpPacketData := tpktPacket.GetPayload().(readWriteModel.COTPPacketData)
		resultChan <- cotpPacketData.GetPayload()
		return nil
	}, func(err error) error {
		errChan <- err
		return nil
	}); err != nil {
		return nil, errors.Wrapf(err, "error sending %s request", interactionInfo)
	}
	select {
	case message := <-resultChan:
		return message, nil
	case err := <-errChan:
		return nil, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (c *Connection) createS7ConnectionRequest(cotpPacketConnectionResponse readWriteModel.COTPPacketConnectionResponse) readWriteModel.TPKTPacket {
	for _, parameter := range cotpPacketConnectionResponse.GetParameters() {
		switch parameter := parameter.(type) {
		case readWriteModel.COTPParameterCalledTsap:
			c.driverContext.CalledTsapId = parameter.GetTsapId()
		case readWriteModel.COTPParameterCallingTsap:
			if parameter.GetTsapId() != c.driverContext.CallingTsapId {
				c.driverContext.CallingTsapId = parameter.GetTsapId()
				c.log.Warn().Uint16("callingTsapId", c.driverContext.CallingTsapId).Msg("Switching calling TSAP id to")
			}
		case readWriteModel.COTPParameterTpduSize:
			c.driverContext.CotpTpduSize = parameter.GetTpduSize()
		default:
			c.log.Warn().Type("v", parameter).Msg("Got unknown parameter type")
		}
	}

	s7ParameterSetupCommunication := readWriteModel.NewS7ParameterSetupCommunication(
		c.driverContext.MaxAmqCaller, c.driverContext.MaxAmqCallee, c.driverContext.PduSize,
	)
	s7Message := readWriteModel.NewS7MessageRequest(0, s7ParameterSetupCommunication, nil)
	cotpPacketData := readWriteModel.NewCOTPPacketData(nil, s7Message, true, 1)
	return readWriteModel.NewTPKTPacket(cotpPacketData)
}

func (c *Connection) createCOTPConnectionRequest() readWriteModel.COTPPacket {
	return readWriteModel.NewCOTPPacketConnectionRequest(
		[]readWriteModel.COTPParameter{
			readWriteModel.NewCOTPParameterCallingTsap(c.driverContext.CallingTsapId),
			readWriteModel.NewCOTPParameterCalledTsap(c.driverContext.CalledTsapId),
			readWriteModel.NewCOTPParameterTpduSize(c.driverContext.CotpTpduSize),
		},
		nil,
		0x0000,
		0x000F,
		readWriteModel.COTPProtocolClass_CLASS_0,
	)
}

func (c *Connection) GetMetadata() apiModel.PlcConnectionMetadata {
	return &_default.DefaultConnectionMetadata{
		ConnectionAttributes: map[string]string{
			"article-number":  c.driverContext.ArticleNumber,
			"controller-type": c.driverContext.ControllerType.String(),
			"pdu-size":        strconv.Itoa(int(c.driverContext.PduSize)),
			"max-amq-caller":  strconv.Itoa(int(c.driverContext.MaxAmqCaller)),
			"max-amq-callee":  strconv.Itoa(int(c.driverContext.MaxAmqCallee)),
		},
		ProvidesReading:     true,
		ProvidesWriting:     true,
		ProvidesSubscribing: c.driverContext.UserDataServicesSupported,
		ProvidesBrowsing:    c.driverContext.UserDataServicesSupported,
	}
}

// Ping does a real round trip: S7 has no dedicated ping, so a one-byte read of %M0 is
// issued and any S7 answer (including a PLC error) counts as alive - only transport
// failures and timeouts report an error.
func (c *Connection) Ping(ctx context.Context) error {
	if !c.IsConnected() {
		return errors.New("not connected")
	}
	tpduId := c.tpduGenerator.getAndIncrement()
	pingRequest := readWriteModel.NewTPKTPacket(readWriteModel.NewCOTPPacketData(
		nil,
		readWriteModel.NewS7MessageRequest(
			tpduId,
			readWriteModel.NewS7ParameterReadVarRequest([]readWriteModel.S7VarRequestParameterItem{
				readWriteModel.NewS7VarRequestParameterItemAddress(
					readWriteModel.NewS7AddressAny(readWriteModel.TransportSize_BYTE, 1, 0, readWriteModel.MemoryArea_FLAGS_MARKERS, 0, 0),
				),
			}),
			nil,
		),
		true,
		uint8(tpduId),
	))
	successChan := make(chan struct{}, 1)
	errChan := make(chan error, 1)
	if err := c.messageCodec.SendRequest(ctx, "ping", pingRequest, func(message spi.Message) bool {
		tpktPacket, ok := message.(readWriteModel.TPKTPacket)
		if !ok {
			return false
		}
		cotpPacketData, ok := tpktPacket.GetPayload().(readWriteModel.COTPPacketData)
		if !ok {
			return false
		}
		payload := cotpPacketData.GetPayload()
		return payload != nil && payload.GetTpduReference() == tpduId
	}, func(message spi.Message) error {
		select {
		case successChan <- struct{}{}:
		default:
		}
		return nil
	}, func(err error) error {
		select {
		case errChan <- err:
		default:
		}
		return nil
	}); err != nil {
		return errors.Wrap(err, "error sending ping request")
	}
	select {
	case <-successChan:
		return nil
	case err := <-errChan:
		return errors.Wrap(err, "got error while waiting for ping response")
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (c *Connection) dispatchUnsolicitedUserData(channel <-chan readWriteModel.S7MessageUserData) {
	for message := range channel {
		c.subscriberMutex.Lock()
		subscriber := c.subscriber
		c.subscriberMutex.Unlock()
		if subscriber == nil || !subscriber.handleUserDataPush(message) {
			c.log.Debug().Stringer("message", message).Msg("Unsolicited user data message not handled by any subscription")
		}
	}
}

func (c *Connection) getSubscriber() *Subscriber {
	c.subscriberMutex.Lock()
	defer c.subscriberMutex.Unlock()
	if c.subscriber == nil {
		c.subscriber = NewSubscriber(c, append(c._options, options.WithCustomLogger(c.log))...)
	}
	return c.subscriber
}

func (c *Connection) SubscriptionRequestBuilder() apiModel.PlcSubscriptionRequestBuilder {
	return spiModel.NewDefaultPlcSubscriptionRequestBuilder(
		c.GetPlcTagHandler(), c.GetPlcValueHandler(), c.getSubscriber())
}

func (c *Connection) UnsubscriptionRequestBuilder() apiModel.PlcUnsubscriptionRequestBuilder {
	return spiModel.NewDefaultPlcUnsubscriptionRequestBuilder()
}

func (c *Connection) Close() error {
	// Cancel PLC-side subscription state before tearing the connection down (bounded).
	c.subscriberMutex.Lock()
	subscriber := c.subscriber
	c.subscriberMutex.Unlock()
	if subscriber != nil {
		drainCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		subscriber.drain(drainCtx)
		cancel()
	}
	return c.DefaultConnection.Close()
}

func (c *Connection) BrowseRequestBuilder() apiModel.PlcBrowseRequestBuilder {
	return spiModel.NewDefaultPlcBrowseRequestBuilder(
		c.GetPlcTagHandler(),
		NewBrowser(c, append(c._options, options.WithCustomLogger(c.log))...),
	)
}

func (c *Connection) ReadRequestBuilder() apiModel.PlcReadRequestBuilder {
	return spiModel.NewDefaultPlcReadRequestBuilder(
		c.GetPlcTagHandler(),
		NewReader(
			&c.tpduGenerator,
			c.messageCodec,
			c.tm,
			&c.driverContext,
			append(c._options, options.WithCustomLogger(c.log))...,
		),
	)
}

func (c *Connection) WriteRequestBuilder() apiModel.PlcWriteRequestBuilder {
	return spiModel.NewDefaultPlcWriteRequestBuilder(
		c.GetPlcTagHandler(), c.GetPlcValueHandler(), NewWriter(
			&c.tpduGenerator,
			c.messageCodec,
			c.tm,
			&c.driverContext,
			append(c._options, options.WithCustomLogger(c.log))...,
		))
}

func (c *Connection) String() string {
	return fmt.Sprintf("s7.Connection")
}
