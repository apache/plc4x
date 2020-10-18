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
	internalModel "plc4x.apache.org/plc4go-modbus-driver/0.8.0/internal/plc4go/model"
	"plc4x.apache.org/plc4go-modbus-driver/0.8.0/internal/plc4go/spi"
	"plc4x.apache.org/plc4go-modbus-driver/0.8.0/pkg/plc4go"
	"plc4x.apache.org/plc4go-modbus-driver/0.8.0/pkg/plc4go/model"
)

type ConnectionMetadata struct {
	model.PlcConnectionMetadata
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

type Connection struct {
	fieldHandler spi.PlcFieldHandler
	valueHandler spi.PlcValueHandler
	plc4go.PlcConnection
}

func NewConnection(fieldHandler spi.PlcFieldHandler, valueHandler spi.PlcValueHandler) Connection {
	return Connection{
		fieldHandler: fieldHandler,
		valueHandler: valueHandler,
	}
}

func (m Connection) Connect() <-chan plc4go.PlcConnectionConnectResult {
	// TODO: Implement ...
	ch := make(chan plc4go.PlcConnectionConnectResult)
	go func() {
		ch <- plc4go.NewPlcConnectionConnectResult(m, nil)
	}()
	return ch
}

func (m Connection) Close() <-chan plc4go.PlcConnectionCloseResult {
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
	panic("implement me")
}

func (m Connection) GetMetadata() model.PlcConnectionMetadata {
	return ConnectionMetadata{}
}

func (m Connection) ReadRequestBuilder() model.PlcReadRequestBuilder {
	return internalModel.NewDefaultPlcReadRequestBuilder(m.fieldHandler, NewModbusReader())
}

func (m Connection) WriteRequestBuilder() model.PlcWriteRequestBuilder {
	return internalModel.NewDefaultPlcWriteRequestBuilder(m.fieldHandler, m.valueHandler, NewModbusWriter())
}

func (m Connection) SubscriptionRequestBuilder() model.PlcSubscriptionRequestBuilder {
	panic("implement me")
}

func (m Connection) UnsubscriptionRequestBuilder() model.PlcUnsubscriptionRequestBuilder {
	panic("implement me")
}
