/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
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
	"fmt"
	readWriteModel "github.com/apache/plc4x/plc4go/internal/plc4go/s7/readwrite/model"
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi"
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi/default"
	internalModel "github.com/apache/plc4x/plc4go/internal/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi/plcerrors"
	"github.com/apache/plc4x/plc4go/pkg/plc4go"
	apiModel "github.com/apache/plc4x/plc4go/pkg/plc4go/model"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"reflect"
	"strings"
	"sync"
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
	tm            *spi.RequestTransactionManager
}

func NewConnection(messageCodec spi.MessageCodec, configuration Configuration, driverContext DriverContext, fieldHandler spi.PlcFieldHandler, tm *spi.RequestTransactionManager) *Connection {
	connection := &Connection{
		tpduGenerator: TpduGenerator{currentTpduId: 10},
		messageCodec:  messageCodec,
		configuration: configuration,
		driverContext: driverContext,
		tm:            tm,
	}
	connection.DefaultConnection = _default.NewDefaultConnection(connection,
		_default.WithPlcFieldHandler(fieldHandler),
		_default.WithPlcValueHandler(NewValueHandler()),
	)
	return connection
}

func (m *Connection) GetConnection() plc4go.PlcConnection {
	return m
}

func (m *Connection) GetMessageCodec() spi.MessageCodec {
	return m.messageCodec
}

func (m *Connection) Connect() <-chan plc4go.PlcConnectionConnectResult {
	log.Trace().Msg("Connecting")
	ch := make(chan plc4go.PlcConnectionConnectResult)
	go func() {
		err := m.messageCodec.Connect()
		if err != nil {
			ch <- plc4go.NewPlcConnectionConnectResult(m, err)
		}

		// Only on active connections we do a connection
		if m.driverContext.PassiveMode {
			log.Info().Msg("S7 Driver running in PASSIVE mode.")
			ch <- plc4go.NewPlcConnectionConnectResult(m, nil)
			return
		}

		// For testing purposes we can skip the waiting for a complete connection
		if !m.driverContext.awaitSetupComplete {
			go m.setupConnection(ch)
			log.Warn().Msg("Connection used in an unsafe way. !!!DON'T USE IN PRODUCTION!!!")
			// Here we write directly and don't wait till the connection is "really" connected
			// Note: we can't use fireConnected here as it's guarded against m.driverContext.awaitSetupComplete
			ch <- plc4go.NewPlcConnectionConnectResult(m, err)
			m.SetConnected(true)
			return
		}

		// Only the TCP transport supports login.
		log.Info().Msg("S7 Driver running in ACTIVE mode.")

		m.setupConnection(ch)
	}()
	return ch
}

func (m *Connection) setupConnection(ch chan plc4go.PlcConnectionConnectResult) {
	log.Debug().Msg("Sending COTP Connection Request")
	// Open the session on ISO Transport Protocol first.
	cotpConnectionResult := make(chan *readWriteModel.COTPPacketConnectionResponse)
	cotpConnectionErrorChan := make(chan error)
	if err := m.messageCodec.SendRequest(
		readWriteModel.NewTPKTPacket(m.createCOTPConnectionRequest()),
		func(message interface{}) bool {
			tpktPacket := readWriteModel.CastTPKTPacket(message)
			if tpktPacket == nil {
				return false
			}
			cotpPacketConnectionResponse := readWriteModel.CastCOTPPacketConnectionResponse(tpktPacket.Payload)
			return cotpPacketConnectionResponse != nil
		},
		func(message interface{}) error {
			tpktPacket := readWriteModel.CastTPKTPacket(message)
			cotpPacketConnectionResponse := readWriteModel.CastCOTPPacketConnectionResponse(tpktPacket.Payload)
			cotpConnectionResult <- cotpPacketConnectionResponse
			return nil
		},
		func(err error) error {
			// If this is a timeout, do a check if the connection requires a reconnection
			if _, isTimeout := err.(plcerrors.TimeoutError); isTimeout {
				log.Warn().Msg("Timeout during Connection establishing, closing channel...")
				m.Close()
			}
			cotpConnectionErrorChan <- errors.Wrap(err, "got error processing request")
			return nil
		},
		m.GetTtl(),
	); err != nil {
		m.fireConnectionError(errors.Wrap(err, "Error during sending of COTP Connection Request"), ch)
	}
	select {
	case cotpPacketConnectionResponse := <-cotpConnectionResult:
		log.Debug().Msg("Got COTP Connection Response")
		log.Debug().Msg("Sending S7 Connection Request")

		// Send an S7 login message.
		s7ConnectionResult := make(chan *readWriteModel.S7ParameterSetupCommunication)
		s7ConnectionErrorChan := make(chan error)
		if err := m.messageCodec.SendRequest(
			m.createS7ConnectionRequest(cotpPacketConnectionResponse),
			func(message interface{}) bool {
				tpktPacket := readWriteModel.CastTPKTPacket(message)
				if tpktPacket == nil {
					return false
				}
				cotpPacketData := readWriteModel.CastCOTPPacketData(tpktPacket.Payload)
				if cotpPacketData == nil {
					return false
				}
				messageResponseData := readWriteModel.CastS7MessageResponseData(cotpPacketData.Parent.Payload)
				if messageResponseData == nil {
					return false
				}
				parameterSetupCommunication := readWriteModel.CastS7ParameterSetupCommunication(messageResponseData.Parent.Parameter)
				return parameterSetupCommunication != nil
			},
			func(message interface{}) error {
				tpktPacket := readWriteModel.CastTPKTPacket(message)
				cotpPacketData := readWriteModel.CastCOTPPacketData(tpktPacket.Payload)
				messageResponseData := readWriteModel.CastS7MessageResponseData(cotpPacketData.Parent.Payload)
				setupCommunication := readWriteModel.CastS7ParameterSetupCommunication(messageResponseData.Parent.Parameter)
				s7ConnectionResult <- setupCommunication
				return nil
			},
			func(err error) error {
				// If this is a timeout, do a check if the connection requires a reconnection
				if _, isTimeout := err.(plcerrors.TimeoutError); isTimeout {
					log.Warn().Msg("Timeout during Connection establishing, closing channel...")
					m.Close()
				}
				s7ConnectionErrorChan <- errors.Wrap(err, "got error processing request")
				return nil
			},
			m.GetTtl(),
		); err != nil {
			m.fireConnectionError(errors.Wrap(err, "Error during sending of S7 Connection Request"), ch)
		}
		select {
		case setupCommunication := <-s7ConnectionResult:
			log.Debug().Msg("Got S7 Connection Response")
			log.Debug().Msg("Sending identify remote Request")
			// Save some data from the response.
			m.driverContext.MaxAmqCaller = setupCommunication.MaxAmqCaller
			m.driverContext.MaxAmqCallee = setupCommunication.MaxAmqCallee
			m.driverContext.PduSize = setupCommunication.PduLength

			// Update the number of concurrent requests to the negotiated number.
			// I have never seen anything else than equal values for caller and
			// callee, but if they were different, we're only limiting the outgoing
			// requests.
			m.tm.SetNumberOfConcurrentRequests(int(m.driverContext.MaxAmqCallee))

			// If the controller type is explicitly set, were finished with the login
			// process. If it's set to ANY, we have to query the serial number information
			// in order to detect the type of PLC.
			if m.driverContext.ControllerType != ControllerType_ANY {
				// Send an event that connection setup is complete.
				m.fireConnected(ch)
				return
			}

			// Prepare a message to request the remote to identify itself.
			log.Debug().Msg("Sending S7 Identification Request")
			s7IdentificationResult := make(chan *readWriteModel.S7PayloadUserData)
			s7IdentificationErrorChan := make(chan error)
			if err := m.messageCodec.SendRequest(
				m.createIdentifyRemoteMessage(),
				func(message interface{}) bool {
					tpktPacket := readWriteModel.CastTPKTPacket(message)
					if tpktPacket == nil {
						return false
					}
					cotpPacketData := readWriteModel.CastCOTPPacketData(tpktPacket.Payload)
					if cotpPacketData == nil {
						return false
					}
					messageUserData := readWriteModel.CastS7MessageUserData(cotpPacketData.Parent.Payload)
					if messageUserData == nil {
						return false
					}
					return readWriteModel.CastS7PayloadUserData(messageUserData.Parent.Payload) != nil
				},
				func(message interface{}) error {
					tpktPacket := readWriteModel.CastTPKTPacket(message)
					cotpPacketData := readWriteModel.CastCOTPPacketData(tpktPacket.Payload)
					messageUserData := readWriteModel.CastS7MessageUserData(cotpPacketData.Parent.Payload)
					s7IdentificationResult <- readWriteModel.CastS7PayloadUserData(messageUserData.Parent.Payload)
					return nil
				},
				func(err error) error {
					// If this is a timeout, do a check if the connection requires a reconnection
					if _, isTimeout := err.(plcerrors.TimeoutError); isTimeout {
						log.Warn().Msg("Timeout during Connection establishing, closing channel...")
						m.Close()
					}
					s7IdentificationErrorChan <- errors.Wrap(err, "got error processing request")
					return nil
				},
				m.GetTtl(),
			); err != nil {
				m.fireConnectionError(errors.Wrap(err, "Error during sending of identify remote Request"), ch)
			}
			select {
			case payloadUserData := <-s7IdentificationResult:
				log.Debug().Msg("Got S7 Identification Response")
				m.extractControllerTypeAndFireConnected(payloadUserData, ch)
			case err := <-s7IdentificationErrorChan:
				m.fireConnectionError(errors.Wrap(err, "Error receiving identify remote Request"), ch)
			}
		case err := <-s7ConnectionErrorChan:
			m.fireConnectionError(errors.Wrap(err, "Error receiving S7 Connection Request"), ch)
		}
	case err := <-cotpConnectionErrorChan:
		m.fireConnectionError(errors.Wrap(err, "Error receiving of COTP Connection Request"), ch)
	}
}

func (m *Connection) fireConnectionError(err error, ch chan<- plc4go.PlcConnectionConnectResult) {
	if m.driverContext.awaitSetupComplete {
		ch <- plc4go.NewPlcConnectionConnectResult(nil, errors.Wrap(err, "Error during connection"))
	} else {
		log.Error().Err(err).Msg("awaitSetupComplete set to false and we got a error during connect")
	}
}

func (m *Connection) fireConnected(ch chan<- plc4go.PlcConnectionConnectResult) {
	if m.driverContext.awaitSetupComplete {
		ch <- plc4go.NewPlcConnectionConnectResult(m, nil)
	} else {
		log.Info().Msg("Successfully connected")
	}
	m.SetConnected(true)
}

func (m *Connection) extractControllerTypeAndFireConnected(payloadUserData *readWriteModel.S7PayloadUserData, ch chan<- plc4go.PlcConnectionConnectResult) {
	// TODO: how do we handle the case if there no items at all? Should we assume it a successful or failure...
	// TODO ... opposed to the java implementation we treat it as a failure
	for _, item := range payloadUserData.Items {
		switch item.Child.(type) {
		case *readWriteModel.S7PayloadUserDataItemCpuFunctionReadSzlResponse:
			readSzlResponseItem := item.Child.(*readWriteModel.S7PayloadUserDataItemCpuFunctionReadSzlResponse)
			for _, readSzlResponseItemItem := range readSzlResponseItem.Items {
				if readSzlResponseItemItem.ItemIndex != 0x0001 {
					continue
				}
				articleNumber := string(readSzlResponseItemItem.Mlfb)
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
					log.Info().Msgf("Looking up unknown article number %s", articleNumber)
					controllerType = ControllerType_ANY
				}
				m.driverContext.ControllerType = controllerType

				// Send an event that connection setup is complete.
				m.fireConnected(ch)
				return
			}
		}
	}
	m.fireConnectionError(errors.New("Coudln't find the required information"), ch)
}

func (m *Connection) createIdentifyRemoteMessage() *readWriteModel.TPKTPacket {
	identifyRemoteMessage := readWriteModel.NewS7MessageUserData(
		1,
		readWriteModel.NewS7ParameterUserData(
			[]*readWriteModel.S7ParameterUserDataItem{
				readWriteModel.NewS7ParameterUserDataItemCPUFunctions(
					0x11,
					0x4,
					0x4,
					0x01,
					0x00,
					nil,
					nil,
					nil),
			},
		),
		readWriteModel.NewS7PayloadUserData(
			[]*readWriteModel.S7PayloadUserDataItem{
				readWriteModel.NewS7PayloadUserDataItemCpuFunctionReadSzlRequest(
					readWriteModel.NewSzlId(
						readWriteModel.SzlModuleTypeClass_CPU,
						0x00,
						readWriteModel.SzlSublist_MODULE_IDENTIFICATION,
					),
					0x0000,
					readWriteModel.DataTransportErrorCode_OK,
					readWriteModel.DataTransportSize_OCTET_STRING,
				),
			},
		),
	)
	cotpPacketData := readWriteModel.NewCOTPPacketData(true, 2, nil, identifyRemoteMessage)
	return readWriteModel.NewTPKTPacket(cotpPacketData)
}

func (m *Connection) createS7ConnectionRequest(cotpPacketConnectionResponse *readWriteModel.COTPPacketConnectionResponse) *readWriteModel.TPKTPacket {
	for _, parameter := range cotpPacketConnectionResponse.Parent.Parameters {
		switch parameter.Child.(type) {
		case *readWriteModel.COTPParameterCalledTsap:
			cotpParameterCalledTsap := parameter.Child.(*readWriteModel.COTPParameterCalledTsap)
			m.driverContext.CalledTsapId = cotpParameterCalledTsap.TsapId
		case *readWriteModel.COTPParameterCallingTsap:
			cotpParameterCallingTsap := parameter.Child.(*readWriteModel.COTPParameterCallingTsap)
			if cotpParameterCallingTsap.TsapId != m.driverContext.CallingTsapId {
				m.driverContext.CallingTsapId = cotpParameterCallingTsap.TsapId
				log.Warn().Msgf("Switching calling TSAP id to '%x'", m.driverContext.CallingTsapId)
			}
		case *readWriteModel.COTPParameterTpduSize:
			cotpParameterTpduSize := parameter.Child.(*readWriteModel.COTPParameterTpduSize)
			m.driverContext.CotpTpduSize = cotpParameterTpduSize.TpduSize
		default:
			log.Warn().Msgf("Got unknown parameter type '%v'", reflect.TypeOf(parameter))
		}
	}

	s7ParameterSetupCommunication := readWriteModel.NewS7ParameterSetupCommunication(
		m.driverContext.MaxAmqCaller, m.driverContext.MaxAmqCallee, m.driverContext.PduSize,
	)
	s7Message := readWriteModel.NewS7MessageRequest(0, s7ParameterSetupCommunication, nil)
	cotpPacketData := readWriteModel.NewCOTPPacketData(true, 1, nil, s7Message)
	return readWriteModel.NewTPKTPacket(cotpPacketData)
}

func (m *Connection) createCOTPConnectionRequest() *readWriteModel.COTPPacket {
	return readWriteModel.NewCOTPPacketConnectionRequest(
		0x0000,
		0x000F,
		readWriteModel.COTPProtocolClass_CLASS_0,
		[]*readWriteModel.COTPParameter{
			readWriteModel.NewCOTPParameterCalledTsap(m.driverContext.CalledTsapId),
			readWriteModel.NewCOTPParameterCallingTsap(m.driverContext.CallingTsapId),
			readWriteModel.NewCOTPParameterTpduSize(m.driverContext.CotpTpduSize),
		},
		nil,
	)
}

func (m *Connection) GetMetadata() apiModel.PlcConnectionMetadata {
	return _default.DefaultConnectionMetadata{
		ProvidesReading: true,
		ProvidesWriting: true,
	}
}

func (m *Connection) ReadRequestBuilder() apiModel.PlcReadRequestBuilder {
	return internalModel.NewDefaultPlcReadRequestBuilder(m.GetPlcFieldHandler(), NewReader(&m.tpduGenerator, m.messageCodec, m.tm))
}

func (m *Connection) WriteRequestBuilder() apiModel.PlcWriteRequestBuilder {
	return internalModel.NewDefaultPlcWriteRequestBuilder(
		m.GetPlcFieldHandler(), m.GetPlcValueHandler(), NewWriter(&m.tpduGenerator, m.messageCodec, m.tm))
}

func (m *Connection) String() string {
	return fmt.Sprintf("s7.Connection")
}
