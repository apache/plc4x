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
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/tracer"
	"github.com/apache/plc4x/plc4go/spi/transactions"
	"github.com/rs/zerolog"
	"runtime/debug"
	"strings"
	"sync"

	"github.com/apache/plc4x/plc4go/pkg/api"
	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/s7/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/default"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/utils"
	"github.com/pkg/errors"
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

	connectionId string
	tracer       tracer.Tracer

	log      zerolog.Logger
	_options []options.WithOption // Used to pass them downstream
}

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

func (c *Connection) ConnectWithContext(ctx context.Context) <-chan plc4go.PlcConnectionConnectResult {
	c.log.Trace().Msg("Connecting")
	ch := make(chan plc4go.PlcConnectionConnectResult, 1)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				ch <- _default.NewDefaultPlcConnectionConnectResult(nil, errors.Errorf("panic-ed %v. Stack: %s", err, debug.Stack()))
			}
		}()
		err := c.messageCodec.ConnectWithContext(ctx)
		if err != nil {
			ch <- _default.NewDefaultPlcConnectionConnectResult(c, err)
		}

		// Only on active connections we do a connection
		if c.driverContext.PassiveMode {
			c.log.Info().Msg("S7 Driver running in PASSIVE mode.")
			ch <- _default.NewDefaultPlcConnectionConnectResult(c, nil)
			return
		}

		// For testing purposes we can skip the waiting for a complete connection
		if !c.driverContext.awaitSetupComplete {
			go c.setupConnection(ctx, ch)
			c.log.Warn().Msg("Connection used in an unsafe way. !!!DON'T USE IN PRODUCTION!!!")
			// Here we write directly and don't wait till the connection is "really" connected
			// Note: we can't use fireConnected here as it's guarded against c.driverContext.awaitSetupComplete
			ch <- _default.NewDefaultPlcConnectionConnectResult(c, err)
			c.SetConnected(true)
			return
		}

		// Only the TCP transport supports login.
		c.log.Info().Msg("S7 Driver running in ACTIVE mode.")

		c.setupConnection(ctx, ch)
	}()
	return ch
}

func (c *Connection) setupConnection(ctx context.Context, ch chan plc4go.PlcConnectionConnectResult) {
	c.log.Debug().Msg("Sending COTP Connection Request")
	// Open the session on ISO Transport Protocol first.
	cotpConnectionResult := make(chan readWriteModel.COTPPacketConnectionResponse, 1)
	cotpConnectionErrorChan := make(chan error, 1)
	if err := c.messageCodec.SendRequest(ctx, readWriteModel.NewTPKTPacket(c.createCOTPConnectionRequest()), func(message spi.Message) bool {
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
		if _, isTimeout := err.(utils.TimeoutError); isTimeout {
			c.log.Warn().Msg("Timeout during Connection establishing, closing channel...")
			c.Close()
		}
		cotpConnectionErrorChan <- errors.Wrap(err, "got error processing request")
		return nil
	}, c.GetTtl()); err != nil {
		c.fireConnectionError(errors.Wrap(err, "Error during sending of COTP Connection Request"), ch)
	}
	select {
	case cotpPacketConnectionResponse := <-cotpConnectionResult:
		c.log.Debug().Msg("Got COTP Connection Response")
		c.log.Debug().Msg("Sending S7 Connection Request")

		// Send an S7 login message.
		s7ConnectionResult := make(chan readWriteModel.S7ParameterSetupCommunication, 1)
		s7ConnectionErrorChan := make(chan error, 1)
		if err := c.messageCodec.SendRequest(ctx, c.createS7ConnectionRequest(cotpPacketConnectionResponse), func(message spi.Message) bool {
			tpktPacket, ok := message.(readWriteModel.TPKTPacketExactly)
			if !ok {
				return false
			}
			cotpPacketData, ok := tpktPacket.GetPayload().(readWriteModel.COTPPacketDataExactly)
			if !ok {
				return false
			}
			messageResponseData, ok := cotpPacketData.GetPayload().(readWriteModel.S7MessageResponseDataExactly)
			if !ok {
				return false
			}
			_, ok = messageResponseData.GetParameter().(readWriteModel.S7ParameterSetupCommunicationExactly)
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
			if _, isTimeout := err.(utils.TimeoutError); isTimeout {
				c.log.Warn().Msg("Timeout during Connection establishing, closing channel...")
				c.Close()
			}
			s7ConnectionErrorChan <- errors.Wrap(err, "got error processing request")
			return nil
		}, c.GetTtl()); err != nil {
			c.fireConnectionError(errors.Wrap(err, "Error during sending of S7 Connection Request"), ch)
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

			// If the controller type is explicitly set, were finished with the login
			// process. If it's set to ANY, we have to query the serial number information
			// in order to detect the type of PLC.
			if c.driverContext.ControllerType != ControllerType_ANY {
				// Send an event that connection setup is complete.
				c.fireConnected(ch)
				return
			}

			// Prepare a message to request the remote to identify itself.
			c.log.Debug().Msg("Sending S7 Identification Request")
			s7IdentificationResult := make(chan readWriteModel.S7PayloadUserData, 1)
			s7IdentificationErrorChan := make(chan error, 1)
			if err := c.messageCodec.SendRequest(ctx, c.createIdentifyRemoteMessage(), func(message spi.Message) bool {
				tpktPacket, ok := message.(readWriteModel.TPKTPacketExactly)
				if !ok {
					return false
				}
				cotpPacketData, ok := tpktPacket.GetPayload().(readWriteModel.COTPPacketDataExactly)
				if !ok {
					return false
				}
				messageUserData, ok := cotpPacketData.GetPayload().(readWriteModel.S7MessageUserDataExactly)
				if !ok {
					return false
				}
				_, ok = messageUserData.GetPayload().(readWriteModel.S7PayloadUserDataExactly)
				return ok
			}, func(message spi.Message) error {
				tpktPacket := message.(readWriteModel.TPKTPacket)
				cotpPacketData := tpktPacket.GetPayload().(readWriteModel.COTPPacketData)
				messageUserData := cotpPacketData.GetPayload().(readWriteModel.S7MessageUserData)
				s7IdentificationResult <- messageUserData.GetPayload().(readWriteModel.S7PayloadUserData)
				return nil
			}, func(err error) error {
				// If this is a timeout, do a check if the connection requires a reconnection
				if _, isTimeout := err.(utils.TimeoutError); isTimeout {
					c.log.Warn().Msg("Timeout during Connection establishing, closing channel...")
					c.Close()
				}
				s7IdentificationErrorChan <- errors.Wrap(err, "got error processing request")
				return nil
			}, c.GetTtl()); err != nil {
				c.fireConnectionError(errors.Wrap(err, "Error during sending of identify remote Request"), ch)
			}
			select {
			case payloadUserData := <-s7IdentificationResult:
				c.log.Debug().Msg("Got S7 Identification Response")
				c.extractControllerTypeAndFireConnected(payloadUserData, ch)
			case err := <-s7IdentificationErrorChan:
				c.fireConnectionError(errors.Wrap(err, "Error receiving identify remote Request"), ch)
			}
		case err := <-s7ConnectionErrorChan:
			c.fireConnectionError(errors.Wrap(err, "Error receiving S7 Connection Request"), ch)
		}
	case err := <-cotpConnectionErrorChan:
		c.fireConnectionError(errors.Wrap(err, "Error receiving of COTP Connection Request"), ch)
	}
}

func (c *Connection) fireConnectionError(err error, ch chan<- plc4go.PlcConnectionConnectResult) {
	if c.driverContext.awaitSetupComplete {
		ch <- _default.NewDefaultPlcConnectionConnectResult(nil, errors.Wrap(err, "Error during connection"))
	} else {
		c.log.Error().Err(err).Msg("awaitSetupComplete set to false and we got a error during connect")
	}
}

func (c *Connection) fireConnected(ch chan<- plc4go.PlcConnectionConnectResult) {
	if c.driverContext.awaitSetupComplete {
		ch <- _default.NewDefaultPlcConnectionConnectResult(c, nil)
	} else {
		c.log.Info().Msg("Successfully connected")
	}
	c.SetConnected(true)
}

func (c *Connection) extractControllerTypeAndFireConnected(payloadUserData readWriteModel.S7PayloadUserData, ch chan<- plc4go.PlcConnectionConnectResult) {
	// TODO: how do we handle the case if there no items at all? Should we assume it a successful or failure...
	// TODO ... opposed to the java implementation we treat it as a failure
	for _, item := range payloadUserData.GetItems() {
		switch readSzlResponseItem := item.(type) {
		case readWriteModel.S7PayloadUserDataItemCpuFunctionReadSzlResponse:
			for _, readSzlResponseItemItem := range readSzlResponseItem.GetItems() {
				_ = readSzlResponseItemItem
				/* TODO: broken by mspec changes from carcia
				if readSzlResponseItemItem.GetItemIndex() != 0x0001 {
					continue
				}
				articleNumber := string(readSzlResponseItemItem.GetMlfb())
				*/
				articleNumber := "broken at the moment"
				var controllerType ControllerType
				if !strings.HasPrefix(articleNumber, "6ES7 ") {
					controllerType = ControllerType_ANY
				}
				blankIndex := strings.Index(articleNumber, " ")
				model := articleNumber[blankIndex+1 : blankIndex+2]
				switch model {
				case "2":
					controllerType = ControllerType_S7_1200
				case "5":
					controllerType = ControllerType_S7_1500
				case "3":
					controllerType = ControllerType_S7_300
				case "4":
					controllerType = ControllerType_S7_400
				default:
					c.log.Info().Str("articleNumber", articleNumber).Msg("Looking up unknown article number")
					controllerType = ControllerType_ANY
				}
				c.driverContext.ControllerType = controllerType

				// Send an event that connection setup is complete.
				c.fireConnected(ch)
				return
			}
		}
	}
	c.fireConnectionError(errors.New("Coudln't find the required information"), ch)
}

func (c *Connection) createIdentifyRemoteMessage() readWriteModel.TPKTPacket {
	identifyRemoteMessage := readWriteModel.NewS7MessageUserData(
		1,
		readWriteModel.NewS7ParameterUserData(
			[]readWriteModel.S7ParameterUserDataItem{
				readWriteModel.NewS7ParameterUserDataItemCPUFunctions(
					0x11,
					0x4,
					0x4,
					0x01,
					0x00,
					nil,
					nil,
					nil,
				),
			},
		),
		readWriteModel.NewS7PayloadUserData(
			[]readWriteModel.S7PayloadUserDataItem{
				readWriteModel.NewS7PayloadUserDataItemCpuFunctionReadSzlRequest(
					readWriteModel.NewSzlId(
						readWriteModel.SzlModuleTypeClass_CPU,
						0x00,
						readWriteModel.SzlSublist_MODULE_IDENTIFICATION,
					),
					0x0000,
					readWriteModel.DataTransportErrorCode_OK,
					readWriteModel.DataTransportSize_OCTET_STRING,
					4,
				),
			},
			nil,
		),
	)
	cotpPacketData := readWriteModel.NewCOTPPacketData(true, 2, nil, identifyRemoteMessage, 0)
	return readWriteModel.NewTPKTPacket(cotpPacketData)
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
	cotpPacketData := readWriteModel.NewCOTPPacketData(true, 1, nil, s7Message, 0)
	return readWriteModel.NewTPKTPacket(cotpPacketData)
}

func (c *Connection) createCOTPConnectionRequest() readWriteModel.COTPPacket {
	return readWriteModel.NewCOTPPacketConnectionRequest(
		0x0000,
		0x000F,
		readWriteModel.COTPProtocolClass_CLASS_0,
		[]readWriteModel.COTPParameter{
			readWriteModel.NewCOTPParameterCallingTsap(c.driverContext.CallingTsapId, 0),
			readWriteModel.NewCOTPParameterCalledTsap(c.driverContext.CalledTsapId, 0),
			readWriteModel.NewCOTPParameterTpduSize(c.driverContext.CotpTpduSize, 0),
		},
		nil,
		0,
	)
}

func (c *Connection) GetMetadata() apiModel.PlcConnectionMetadata {
	return &_default.DefaultConnectionMetadata{
		ProvidesReading: true,
		ProvidesWriting: true,
	}
}

func (c *Connection) ReadRequestBuilder() apiModel.PlcReadRequestBuilder {
	return spiModel.NewDefaultPlcReadRequestBuilder(
		c.GetPlcTagHandler(),
		NewReader(
			&c.tpduGenerator,
			c.messageCodec,
			c.tm,
			append(c._options, options.WithCustomLogger(c.log))...,
		),
	)
}

func (c *Connection) WriteRequestBuilder() apiModel.PlcWriteRequestBuilder {
	return spiModel.NewDefaultPlcWriteRequestBuilder(
		c.GetPlcTagHandler(), c.GetPlcValueHandler(), NewWriter(&c.tpduGenerator, c.messageCodec, c.tm))
}

func (c *Connection) String() string {
	return fmt.Sprintf("s7.Connection")
}
