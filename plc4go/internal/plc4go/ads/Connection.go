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

package ads

import (
	"fmt"
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi"
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi/default"
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi/interceptors"
	internalModel "github.com/apache/plc4x/plc4go/internal/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi/transports"
	"github.com/apache/plc4x/plc4go/pkg/plc4go"
	apiModel "github.com/apache/plc4x/plc4go/pkg/plc4go/model"
)

type Connection struct {
	_default.DefaultConnection
	messageCodec       spi.MessageCodec
	requestInterceptor interceptors.RequestInterceptor
	configuration      Configuration
	reader             *Reader
	writer             *Writer
}

func NewConnection(messageCodec spi.MessageCodec, configuration Configuration, fieldHandler spi.PlcFieldHandler) (*Connection, error) {
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
		reader: &reader,
		writer: &writer,
	}
	connection.DefaultConnection = _default.NewDefaultConnection(connection,
		_default.WithPlcFieldHandler(fieldHandler),
		_default.WithPlcValueHandler(NewValueHandler()),
	)
	return connection, nil
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
	}
}

func (m *Connection) ReadRequestBuilder() apiModel.PlcReadRequestBuilder {
	return internalModel.NewDefaultPlcReadRequestBuilder(m.GetPlcFieldHandler(), m.reader)
}

func (m *Connection) WriteRequestBuilder() apiModel.PlcWriteRequestBuilder {
	return internalModel.NewDefaultPlcWriteRequestBuilder(m.GetPlcFieldHandler(), m.GetPlcValueHandler(), m.writer)
}

func (m *Connection) SubscriptionRequestBuilder() apiModel.PlcSubscriptionRequestBuilder {
	panic("implement me")
}

func (m *Connection) UnsubscriptionRequestBuilder() apiModel.PlcUnsubscriptionRequestBuilder {
	panic("implement me")
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
