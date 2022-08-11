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

package eip

import (
	"fmt"
	"github.com/apache/plc4x/plc4go/pkg/api"
	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/eip/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/default"
	internalModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/plcerrors"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type Connection struct {
	_default.DefaultConnection
	messageCodec  spi.MessageCodec
	configuration Configuration
	driverContext DriverContext
	tm            *spi.RequestTransactionManager
	sessionHandle uint32
	senderContext []uint8
	connectionId  string
	tracer        *spi.Tracer
}

func NewConnection(messageCodec spi.MessageCodec, configuration Configuration, driverContext DriverContext, fieldHandler spi.PlcFieldHandler, tm *spi.RequestTransactionManager, options map[string][]string) *Connection {
	connection := &Connection{
		messageCodec:  messageCodec,
		configuration: configuration,
		driverContext: driverContext,
		tm:            tm,
	}
	if traceEnabledOption, ok := options["traceEnabled"]; ok {
		if len(traceEnabledOption) == 1 {
			connection.tracer = spi.NewTracer(connection.connectionId)
		}
	}
	connection.DefaultConnection = _default.NewDefaultConnection(connection,
		_default.WithPlcFieldHandler(fieldHandler),
		_default.WithPlcValueHandler(NewValueHandler()),
	)
	return connection
}

func (m *Connection) GetConnectionId() string {
	return m.connectionId
}

func (m *Connection) IsTraceEnabled() bool {
	return m.tracer != nil
}

func (m *Connection) GetTracer() *spi.Tracer {
	return m.tracer
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
			ch <- _default.NewDefaultPlcConnectionConnectResult(m, err)
		}

		// For testing purposes we can skip the waiting for a complete connection
		if !m.driverContext.awaitSetupComplete {
			go m.setupConnection(ch)
			log.Warn().Msg("Connection used in an unsafe way. !!!DON'T USE IN PRODUCTION!!!")
			// Here we write directly and don't wait till the connection is "really" connected
			// Note: we can't use fireConnected here as it's guarded against m.driverContext.awaitSetupComplete
			ch <- _default.NewDefaultPlcConnectionConnectResult(m, err)
			m.SetConnected(true)
			return
		}

		m.setupConnection(ch)
	}()
	return ch
}

func (m *Connection) Close() <-chan plc4go.PlcConnectionCloseResult {
	result := make(chan plc4go.PlcConnectionCloseResult)
	go func() {
		log.Debug().Msg("Sending UnregisterSession EIP Packet")
		_ = m.messageCodec.SendRequest(
			readWriteModel.NewEipDisconnectRequest(m.sessionHandle, 0, make([]byte, 8), 0),
			func(message spi.Message) bool {
				return true
			},
			func(message spi.Message) error {
				return nil
			},
			func(err error) error {
				return nil
			},
			m.GetTtl(),
		) //Unregister gets no response
		log.Debug().Msgf("Unregistred Session %d", m.sessionHandle)
	}()
	return result
}

func (m *Connection) setupConnection(ch chan plc4go.PlcConnectionConnectResult) {
	log.Debug().Msg("Sending EIP Connection Request")
	if err := m.messageCodec.SendRequest(
		readWriteModel.NewEipConnectionRequest(0, 0, make([]byte, 8), 0),
		func(message spi.Message) bool {
			eipPacket := message.(readWriteModel.EipPacket)
			if eipPacket == nil {
				return false
			}
			eipPacketConnectionRequest := eipPacket.(readWriteModel.EipConnectionRequest)
			return eipPacketConnectionRequest != nil
		},
		func(message spi.Message) error {
			eipPacket := message.(readWriteModel.EipPacket)
			if eipPacket.GetStatus() == 0 {
				m.sessionHandle = eipPacket.GetSessionHandle()
				m.senderContext = eipPacket.GetSenderContext()
				log.Debug().Msgf("Got assigned with Session %d", m.sessionHandle)
				// Send an event that connection setup is complete.
				m.fireConnected(ch)
			} else {

			}
			return nil
		},
		func(err error) error {
			// If this is a timeout, do a check if the connection requires a reconnection
			if _, isTimeout := err.(plcerrors.TimeoutError); isTimeout {
				log.Warn().Msg("Timeout during Connection establishing, closing channel...")
				m.Close()
			}
			m.fireConnectionError(errors.Wrap(err, "got error processing request"), ch)
			return nil
		},
		m.GetTtl(),
	); err != nil {
		m.fireConnectionError(errors.Wrap(err, "Error during sending of EIP Connection Request"), ch)
	}
}

func (m *Connection) fireConnectionError(err error, ch chan<- plc4go.PlcConnectionConnectResult) {
	if m.driverContext.awaitSetupComplete {
		ch <- _default.NewDefaultPlcConnectionConnectResult(nil, errors.Wrap(err, "Error during connection"))
	} else {
		log.Error().Err(err).Msg("awaitSetupComplete set to false and we got a error during connect")
	}
}

func (m *Connection) fireConnected(ch chan<- plc4go.PlcConnectionConnectResult) {
	if m.driverContext.awaitSetupComplete {
		ch <- _default.NewDefaultPlcConnectionConnectResult(m, nil)
	} else {
		log.Info().Msg("Successfully connected")
	}
	m.SetConnected(true)
}

func (m *Connection) GetMetadata() apiModel.PlcConnectionMetadata {
	return _default.DefaultConnectionMetadata{
		ProvidesReading: true,
		ProvidesWriting: true,
	}
}

func (m *Connection) ReadRequestBuilder() apiModel.PlcReadRequestBuilder {
	return internalModel.NewDefaultPlcReadRequestBuilder(m.GetPlcFieldHandler(), NewReader(m.messageCodec, m.tm, m.configuration, &m.sessionHandle))
}

func (m *Connection) WriteRequestBuilder() apiModel.PlcWriteRequestBuilder {
	return internalModel.NewDefaultPlcWriteRequestBuilder(
		m.GetPlcFieldHandler(), m.GetPlcValueHandler(), NewWriter(m.messageCodec, m.tm, m.configuration, &m.sessionHandle, &m.senderContext))
}

func (m *Connection) String() string {
	return fmt.Sprintf("eip.Connection")
}
