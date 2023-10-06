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

package transports

import (
	"github.com/apache/plc4x/plc4go/pkg/api"
	"github.com/apache/plc4x/plc4go/pkg/api/config"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/options/converter"
	"github.com/apache/plc4x/plc4go/spi/transports/serial"
	"github.com/apache/plc4x/plc4go/spi/transports/tcp"
	"github.com/apache/plc4x/plc4go/spi/transports/udp"
)

func RegisterTcpTransport(driverManager plc4go.PlcDriverManager, _options ...config.WithOption) {
	driverManager.(spi.TransportAware).RegisterTransport(tcp.NewTransport(converter.WithOptionToInternal(_options...)...))
}

func RegisterUdpTransport(driverManager plc4go.PlcDriverManager, _options ...config.WithOption) {
	driverManager.(spi.TransportAware).RegisterTransport(udp.NewTransport(converter.WithOptionToInternal(_options...)...))
}

func RegisterSerialTransport(driverManager plc4go.PlcDriverManager, _options ...config.WithOption) {
	driverManager.(spi.TransportAware).RegisterTransport(serial.NewTransport(converter.WithOptionToInternal(_options...)...))
}
