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

package drivers

import (
	"github.com/apache/plc4x/plc4go/internal/plc4go/ads"
	"github.com/apache/plc4x/plc4go/internal/plc4go/eip"
	"github.com/apache/plc4x/plc4go/internal/plc4go/knxnetip"
	"github.com/apache/plc4x/plc4go/internal/plc4go/modbus"
	"github.com/apache/plc4x/plc4go/pkg/plc4go"
	"github.com/apache/plc4x/plc4go/pkg/plc4go/transports"
)

func RegisterAdsDriver(driverManager plc4go.PlcDriverManager) {
	driverManager.RegisterDriver(ads.NewDriver())
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

func RegisterModbusDriver(driverManager plc4go.PlcDriverManager) {
	driverManager.RegisterDriver(modbus.NewDriver())
	transports.RegisterTcpTransport(driverManager)
}
