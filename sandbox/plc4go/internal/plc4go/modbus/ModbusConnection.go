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

type ModbusConnectionMetadata struct {
	model.PlcConnectionMetadata
}

func (m ModbusConnectionMetadata) CanRead() bool {
	return true
}

func (m ModbusConnectionMetadata) CanWrite() bool {
	return true
}

func (m ModbusConnectionMetadata) CanSubscribe() bool {
	return false
}

type ModbusConnection struct {
	fieldHandler spi.PlcFieldHandler
	plc4go.PlcConnection
}

func NewModbusConnection(fieldHandler spi.PlcFieldHandler) ModbusConnection {
	return ModbusConnection{
		fieldHandler: fieldHandler,
	}
}

func (m ModbusConnection) Connect() <-chan plc4go.PlcConnectionConnectResult {
	panic("implement me")
}

func (m ModbusConnection) Close() <-chan plc4go.PlcConnectionCloseResult {
	panic("implement me")
}

func (m ModbusConnection) IsConnected() bool {
	panic("implement me")
}

func (m ModbusConnection) Ping() <-chan plc4go.PlcConnectionPingResult {
	panic("implement me")
}

func (m ModbusConnection) GetMetadata() model.PlcConnectionMetadata {
	return ModbusConnectionMetadata{}
}

func (m ModbusConnection) ReadRequestBuilder() model.PlcReadRequestBuilder {
	// TODO: Pass a real reader in here ...
	return internalModel.NewDefaultPlcReadRequestBuilder(m.fieldHandler, nil)
}

func (m ModbusConnection) WriteRequestBuilder() model.PlcWriteRequestBuilder {
	panic("implement me")
}

func (m ModbusConnection) SubscriptionRequestBuilder() model.PlcSubscriptionRequestBuilder {
	panic("implement me")
}

func (m ModbusConnection) UnsubscriptionRequestBuilder() model.PlcUnsubscriptionRequestBuilder {
	panic("implement me")
}
