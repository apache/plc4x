/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package tests

import (
	"github.com/apache/plc4x/plc4go/internal/modbus"
	"github.com/apache/plc4x/plc4go/internal/spi/testutils"
	"github.com/apache/plc4x/plc4go/internal/spi/utils"
	modbusIO "github.com/apache/plc4x/plc4go/protocols/modbus/readwrite"
	modbusModel "github.com/apache/plc4x/plc4go/protocols/modbus/readwrite/model"
	_ "github.com/apache/plc4x/plc4go/tests/initializetest"
	"testing"
)

func TestModbusDriver(t *testing.T) {
	options := []testutils.WithOption{testutils.WithRootTypeParser(func(readBufferByteBased utils.ReadBufferByteBased) (interface{}, error) {
		return modbusModel.ModbusTcpADUParse(readBufferByteBased, modbusModel.DriverType_MODBUS_TCP, false)
	})}
	testutils.RunDriverTestsuiteWithOptions(t, modbus.NewModbusTcpDriver(), "assets/testing/protocols/modbus/tcp/DriverTestsuite.xml", modbusIO.ModbusXmlParserHelper{}, options)
}
