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

package drivers

import (
	"github.com/apache/plc4x/plc4go/internal/ads"
	"github.com/apache/plc4x/plc4go/internal/bacnetip"
	"github.com/apache/plc4x/plc4go/internal/cbus"
	"github.com/apache/plc4x/plc4go/internal/eip"
	"github.com/apache/plc4x/plc4go/internal/knxnetip"
	modbus2 "github.com/apache/plc4x/plc4go/internal/modbus"
	"github.com/apache/plc4x/plc4go/internal/s7"
	"github.com/apache/plc4x/plc4go/pkg/api"
	"github.com/apache/plc4x/plc4go/pkg/api/transports"
)

func RegisterAdsDriver(driverManager plc4go.PlcDriverManager) {
	driverManager.RegisterDriver(ads.NewDriver())
	transports.RegisterTcpTransport(driverManager)
}

func RegisterBacnetDriver(driverManager plc4go.PlcDriverManager) {
	driverManager.RegisterDriver(bacnetip.NewDriver())
	transports.RegisterUdpTransport(driverManager)
}

func RegisterCBusDriver(driverManager plc4go.PlcDriverManager) {
	driverManager.RegisterDriver(cbus.NewDriver())
	transports.RegisterTcpTransport(driverManager)
}

func RegisterEipDriver(driverManager plc4go.PlcDriverManager) {
	driverManager.RegisterDriver(eip.NewDriver())
	transports.RegisterTcpTransport(driverManager)
}

func RegisterKnxDriver(driverManager plc4go.PlcDriverManager) {
	driverManager.RegisterDriver(knxnetip.NewDriver())
	transports.RegisterUdpTransport(driverManager)
}

func RegisterModbusTcpDriver(driverManager plc4go.PlcDriverManager) {
	driverManager.RegisterDriver(modbus2.NewModbusTcpDriver())
	transports.RegisterTcpTransport(driverManager)
}

func RegisterModbusRtuDriver(driverManager plc4go.PlcDriverManager) {
	driverManager.RegisterDriver(modbus2.NewModbusRtuDriver())
	transports.RegisterSerialTransport(driverManager)
	transports.RegisterTcpTransport(driverManager)
}

func RegisterModbusAsciiDriver(driverManager plc4go.PlcDriverManager) {
	driverManager.RegisterDriver(modbus2.NewModbusAsciiDriver())
	transports.RegisterSerialTransport(driverManager)
	transports.RegisterTcpTransport(driverManager)
}

func RegisterS7Driver(driverManager plc4go.PlcDriverManager) {
	driverManager.RegisterDriver(s7.NewDriver())
	transports.RegisterTcpTransport(driverManager)
}
