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

package ads

import (
	"fmt"
	"math"
	"sync/atomic"

	"github.com/apache/plc4x/plc4go/pkg/api"
	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/protocols/ads/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/default"
	"github.com/apache/plc4x/plc4go/spi/interceptors"
	internalModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/transports"
)

type Connection struct {
	_default.DefaultConnection
	messageCodec       spi.MessageCodec
	requestInterceptor interceptors.RequestInterceptor
	configuration      Configuration
	reader             *Reader
	writer             *Writer
	connectionId       string
	invokeId           uint32
	dataTypeTable      map[string]model.AdsDataTypeTableEntry
	symbolTable        map[string]model.AdsSymbolTableEntry
	tracer             *spi.Tracer
}

/*

	// First read the device info
	deviceInfoResponseChanel := make(chan model.AdsReadDeviceInfoResponse)
	go func() {
		deviceInfoRequest := model.NewAdsReadDeviceInfoRequest(
			m.targetAmsNetId, uint16(model.DefaultAmsPorts_RUNTIME_SYSTEM_01), m.sourceAmsNetId,
			800, 0, m.transactionIdentifier)
		if err := m.messageCodec.SendRequest(
			context.TODO(),
			model.NewAmsTCPPacket(deviceInfoRequest),
			func(message spi.Message) bool {
				amsTcpPacket, ok := message.(model.AmsTCPPacket)
				if !ok {
					return false
				}
				return amsTcpPacket.GetUserdata().GetInvokeId() == deviceInfoRequest.GetInvokeId()
			},
			func(message spi.Message) error {
				amsTcpPacket := message.(model.AmsTCPPacket)
				deviceInfoResponse := amsTcpPacket.GetUserdata().(model.AdsReadDeviceInfoResponse)
				deviceInfoResponseChanel <- deviceInfoResponse
				close(deviceInfoResponseChanel)
				return nil
			},
			func(err error) error {
				return nil
			},
			time.Second); err != nil {
			// TODO: Return an error
		} else {
			close(deviceInfoResponseChanel)
		}
	}()
	deviceInfoResponse := ReadWithTimeout(deviceInfoResponseChanel)
	if deviceInfoResponse == nil {
		return apiModel.PlcResponseCode_NOT_FOUND, []apiModel.PlcBrowseFoundField{}
	}

*/

func NewConnection(messageCodec spi.MessageCodec, configuration Configuration, tagHandler spi.PlcTagHandler, options map[string][]string) (*Connection, error) {
	reader := *NewReader(
		messageCodec,
		configuration.targetAmsNetId,
		configuration.targetAmsPort,
		configuration.sourceAmsNetId,
		configuration.sourceAmsPort,
	)
	writer := *NewWriter(
		messageCodec,
		configuration.targetAmsNetId,
		configuration.targetAmsPort,
		configuration.sourceAmsNetId,
		configuration.sourceAmsPort,
		&reader,
	)
	connection := &Connection{
		messageCodec: messageCodec,
		requestInterceptor: interceptors.NewSingleItemRequestInterceptor(
			internalModel.NewDefaultPlcReadRequest,
			internalModel.NewDefaultPlcWriteRequest,
			internalModel.NewDefaultPlcReadResponse,
			internalModel.NewDefaultPlcWriteResponse,
		),
		configuration: configuration,
		reader:        &reader,
		writer:        &writer,
		invokeId:      0,
	}
	if traceEnabledOption, ok := options["traceEnabled"]; ok {
		if len(traceEnabledOption) == 1 {
			connection.tracer = spi.NewTracer(connection.connectionId)
		}
	}
	connection.DefaultConnection = _default.NewDefaultConnection(connection,
		_default.WithPlcTagHandler(tagHandler),
		_default.WithPlcValueHandler(NewValueHandler()),
	)
	return connection, nil
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

func (m *Connection) GetMetadata() apiModel.PlcConnectionMetadata {
	return _default.DefaultConnectionMetadata{
		ProvidesReading:     true,
		ProvidesWriting:     true,
		ProvidesSubscribing: true,
		ProvidesBrowsing:    true,
	}
}

func (m *Connection) ReadRequestBuilder() apiModel.PlcReadRequestBuilder {
	return internalModel.NewDefaultPlcReadRequestBuilder(m.GetPlcTagHandler(), m.reader)
}

func (m *Connection) WriteRequestBuilder() apiModel.PlcWriteRequestBuilder {
	return internalModel.NewDefaultPlcWriteRequestBuilder(m.GetPlcTagHandler(), m.GetPlcValueHandler(), m.writer)
}

func (m *Connection) SubscriptionRequestBuilder() apiModel.PlcSubscriptionRequestBuilder {
	panic("implement me")
}

func (m *Connection) UnsubscriptionRequestBuilder() apiModel.PlcUnsubscriptionRequestBuilder {
	panic("implement me")
}

func (m *Connection) BrowseRequestBuilder() apiModel.PlcBrowseRequestBuilder {
	return internalModel.NewDefaultPlcBrowseRequestBuilder(m.GetPlcTagHandler(), m)
}

func (m *Connection) GetTransportInstance() transports.TransportInstance {
	if mc, ok := m.messageCodec.(spi.TransportInstanceExposer); ok {
		return mc.GetTransportInstance()
	}
	return nil
}

func (m *Connection) String() string {
	return fmt.Sprintf("ads.Connection{}")
}

func (m *Connection) getInvokeId() uint32 {
	// Calculate a new transaction identifier
	transactionIdentifier := atomic.AddUint32(&m.invokeId, 1)
	if transactionIdentifier > math.MaxUint8 {
		transactionIdentifier = 1
		atomic.StoreUint32(&m.invokeId, 1)
	}
	return transactionIdentifier
}
