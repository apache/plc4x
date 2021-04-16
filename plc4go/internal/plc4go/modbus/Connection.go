//
// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
//

package modbus

import (
	"fmt"
	readWriteModel "github.com/apache/plc4x/plc4go/internal/plc4go/modbus/readwrite/model"
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi"
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi/interceptors"
	internalModel "github.com/apache/plc4x/plc4go/internal/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi/transports"
	"github.com/apache/plc4x/plc4go/pkg/plc4go"
	apiModel "github.com/apache/plc4x/plc4go/pkg/plc4go/model"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"time"
)

type ConnectionMetadata struct {
}

func (m ConnectionMetadata) GetConnectionAttributes() map[string]string {
	return map[string]string{}
}

func (m ConnectionMetadata) CanRead() bool {
	return true
}

func (m ConnectionMetadata) CanWrite() bool {
	return true
}

func (m ConnectionMetadata) CanSubscribe() bool {
	return false
}

func (m ConnectionMetadata) CanBrowse() bool {
	return false
}

// TODO: maybe we can use a DefaultConnection struct here with delegates
type Connection struct {
	unitIdentifier     uint8
	messageCodec       spi.MessageCodec
	options            map[string][]string
	fieldHandler       spi.PlcFieldHandler
	valueHandler       spi.PlcValueHandler
	requestInterceptor internalModel.RequestInterceptor
}

func NewConnection(unitIdentifier uint8, messageCodec spi.MessageCodec, options map[string][]string, fieldHandler spi.PlcFieldHandler) Connection {
	return Connection{
		unitIdentifier:     unitIdentifier,
		messageCodec:       messageCodec,
		options:            options,
		fieldHandler:       fieldHandler,
		valueHandler:       NewValueHandler(),
		requestInterceptor: interceptors.NewSingleItemRequestInterceptor(),
	}
}

func (m Connection) Connect() <-chan plc4go.PlcConnectionConnectResult {
	log.Trace().Msg("Connecting")
	ch := make(chan plc4go.PlcConnectionConnectResult)
	go func() {
		err := m.messageCodec.Connect()
		ch <- plc4go.NewPlcConnectionConnectResult(m, err)
	}()
	return ch
}

func (m Connection) BlockingClose() {
	log.Trace().Msg("Closing blocked")
	closeResults := m.Close()
	select {
	case <-closeResults:
		return
	case <-time.After(time.Second * 5):
		return
	}
}

func (m Connection) Close() <-chan plc4go.PlcConnectionCloseResult {
	log.Trace().Msg("Close")
	// TODO: Implement ...
	ch := make(chan plc4go.PlcConnectionCloseResult)
	go func() {
		ch <- plc4go.NewPlcConnectionCloseResult(m, nil)
	}()
	return ch
}

func (m Connection) IsConnected() bool {
	panic("implement me")
}

func (m Connection) Ping() <-chan plc4go.PlcConnectionPingResult {
	log.Trace().Msg("Pinging")
	result := make(chan plc4go.PlcConnectionPingResult)
	go func() {
		diagnosticRequestPdu := readWriteModel.NewModbusPDUDiagnosticRequest(0, 0x42)
		pingRequest := readWriteModel.NewModbusTcpADU(1, m.unitIdentifier, diagnosticRequestPdu)
		if err := m.messageCodec.SendRequest(
			pingRequest,
			func(message interface{}) bool {
				responseAdu := readWriteModel.CastModbusTcpADU(message)
				return responseAdu.TransactionIdentifier == 1 && responseAdu.UnitIdentifier == m.unitIdentifier
			},
			func(message interface{}) error {
				log.Trace().Msgf("Received Message")
				if message != nil {
					// If we got a valid response (even if it will probably contain an error, we know the remote is available)
					log.Trace().Msg("got valid response")
					result <- plc4go.NewPlcConnectionPingResult(nil)
				} else {
					log.Trace().Msg("got no response")
					result <- plc4go.NewPlcConnectionPingResult(errors.New("no response"))
				}
				return nil
			},
			func(err error) error {
				log.Trace().Msgf("Received Error")
				result <- plc4go.NewPlcConnectionPingResult(errors.Wrap(err, "got error processing request"))
				return nil
			},
			time.Second*1); err != nil {
			result <- plc4go.NewPlcConnectionPingResult(err)
		}
	}()
	return result
}

func (m Connection) GetMetadata() apiModel.PlcConnectionMetadata {
	return ConnectionMetadata{}
}

func (m Connection) ReadRequestBuilder() apiModel.PlcReadRequestBuilder {
	return internalModel.NewDefaultPlcReadRequestBuilderWithInterceptor(m.fieldHandler,
		NewReader(m.unitIdentifier, m.messageCodec), m.requestInterceptor)
}

func (m Connection) WriteRequestBuilder() apiModel.PlcWriteRequestBuilder {
	return internalModel.NewDefaultPlcWriteRequestBuilder(
		m.fieldHandler, m.valueHandler, NewWriter(m.unitIdentifier, m.messageCodec))
}

func (m Connection) SubscriptionRequestBuilder() apiModel.PlcSubscriptionRequestBuilder {
	panic("implement me")
}

func (m Connection) UnsubscriptionRequestBuilder() apiModel.PlcUnsubscriptionRequestBuilder {
	panic("implement me")
}

func (m Connection) BrowseRequestBuilder() apiModel.PlcBrowseRequestBuilder {
	panic("implement me")
}

func (m Connection) GetTransportInstance() transports.TransportInstance {
	if mc, ok := m.messageCodec.(spi.TransportInstanceExposer); ok {
		return mc.GetTransportInstance()
	}
	return nil
}

func (m Connection) GetPlcFieldHandler() spi.PlcFieldHandler {
	return m.fieldHandler
}

func (m Connection) GetPlcValueHandler() spi.PlcValueHandler {
	return m.valueHandler
}

func (m Connection) String() string {
	return fmt.Sprintf("modbus.Connection{unitIdentifier: %d}", m.unitIdentifier)
}
