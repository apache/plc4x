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

package modbus

import (
	"fmt"
	"github.com/apache/plc4x/plc4go/pkg/api"
	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/modbus/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/default"
	"github.com/apache/plc4x/plc4go/spi/interceptors"
	internalModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"time"
)

type Connection struct {
	_default.DefaultConnection
	unitIdentifier     uint8
	messageCodec       spi.MessageCodec
	options            map[string][]string
	requestInterceptor interceptors.RequestInterceptor

	connectionId string
	tracer       *spi.Tracer
}

func NewConnection(unitIdentifier uint8, messageCodec spi.MessageCodec, options map[string][]string, fieldHandler spi.PlcFieldHandler) *Connection {
	connection := &Connection{
		unitIdentifier: unitIdentifier,
		messageCodec:   messageCodec,
		options:        options,
		requestInterceptor: interceptors.NewSingleItemRequestInterceptor(
			internalModel.NewDefaultPlcReadRequest,
			internalModel.NewDefaultPlcWriteRequest,
			internalModel.NewDefaultPlcReadResponse,
			internalModel.NewDefaultPlcWriteResponse,
		),
	}
	if traceEnabledOption, ok := options["traceEnabled"]; ok {
		if len(traceEnabledOption) == 1 {
			connection.tracer = spi.NewTracer(connection.connectionId)
		}
	}
	connection.DefaultConnection = _default.NewDefaultConnection(connection,
		_default.WithDefaultTtl(time.Second*5),
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

func (m *Connection) Ping() <-chan plc4go.PlcConnectionPingResult {
	log.Trace().Msg("Pinging")
	result := make(chan plc4go.PlcConnectionPingResult)
	go func() {
		diagnosticRequestPdu := readWriteModel.NewModbusPDUDiagnosticRequest(0, 0x42)
		pingRequest := readWriteModel.NewModbusTcpADU(1, m.unitIdentifier, diagnosticRequestPdu, false)
		if err := m.messageCodec.SendRequest(
			pingRequest,
			func(message spi.Message) bool {
				responseAdu := readWriteModel.CastModbusTcpADU(message)
				return responseAdu.GetTransactionIdentifier() == 1 && responseAdu.GetUnitIdentifier() == m.unitIdentifier
			},
			func(message spi.Message) error {
				log.Trace().Msgf("Received Message")
				if message != nil {
					// If we got a valid response (even if it will probably contain an error, we know the remote is available)
					log.Trace().Msg("got valid response")
					result <- _default.NewDefaultPlcConnectionPingResult(nil)
				} else {
					log.Trace().Msg("got no response")
					result <- _default.NewDefaultPlcConnectionPingResult(errors.New("no response"))
				}
				return nil
			},
			func(err error) error {
				log.Trace().Msgf("Received Error")
				result <- _default.NewDefaultPlcConnectionPingResult(errors.Wrap(err, "got error processing request"))
				return nil
			},
			time.Second*1); err != nil {
			result <- _default.NewDefaultPlcConnectionPingResult(err)
		}
	}()
	return result
}

func (m *Connection) GetMetadata() apiModel.PlcConnectionMetadata {
	return _default.DefaultConnectionMetadata{
		ProvidesReading: true,
		ProvidesWriting: true,
	}
}

func (m *Connection) ReadRequestBuilder() apiModel.PlcReadRequestBuilder {
	return internalModel.NewDefaultPlcReadRequestBuilderWithInterceptor(
		m.GetPlcFieldHandler(),
		NewReader(m.unitIdentifier, m.messageCodec),
		m.requestInterceptor,
	)
}

func (m *Connection) WriteRequestBuilder() apiModel.PlcWriteRequestBuilder {
	return internalModel.NewDefaultPlcWriteRequestBuilderWithInterceptor(
		m.GetPlcFieldHandler(),
		m.GetPlcValueHandler(),
		NewWriter(m.unitIdentifier, m.messageCodec),
		m.requestInterceptor,
	)
}

func (m *Connection) String() string {
	return fmt.Sprintf("modbus.Connection{unitIdentifier: %d}", m.unitIdentifier)
}
