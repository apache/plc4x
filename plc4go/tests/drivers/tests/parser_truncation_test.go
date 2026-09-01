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

package tests

import (
	"testing"

	abethIO "github.com/apache/plc4x/plc4go/protocols/abeth/readwrite"
	adsIO "github.com/apache/plc4x/plc4go/protocols/ads/readwrite"
	bacnetIO "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite"
	df1IO "github.com/apache/plc4x/plc4go/protocols/df1/readwrite"
	eipIO "github.com/apache/plc4x/plc4go/protocols/eip/readwrite"
	firmataIO "github.com/apache/plc4x/plc4go/protocols/firmata/readwrite"
	iecIO "github.com/apache/plc4x/plc4go/protocols/iec608705104/readwrite"
	knxIO "github.com/apache/plc4x/plc4go/protocols/knxnetip/readwrite"
	modbusIO "github.com/apache/plc4x/plc4go/protocols/modbus/readwrite"
	s7IO "github.com/apache/plc4x/plc4go/protocols/s7/readwrite"
	slmpIO "github.com/apache/plc4x/plc4go/protocols/slmp/readwrite"
	"github.com/apache/plc4x/plc4go/spi/testutils"
)

// The remote end decides where a message ends, so every prefix of a well formed message is
// something a driver can be handed. Running the reference messages of each protocol through
// RunParserTruncationTestsuite requires all of those prefixes to fail as errors rather than as
// panics - see that function for why truncation in particular is worth exercising.
//
// The corpora are the same ones the parser-serializer tests use, so a protocol earns this coverage
// simply by having reference messages.
func TestParserTruncation(t *testing.T) {
	tests := []struct {
		name   string
		path   string
		parser testutils.Parser
	}{
		{"abeth", "assets/testing/protocols/abeth/ParserSerializerTestsuite.xml", abethIO.AbethParserHelper{}},
		{"ads", "assets/testing/protocols/ads/ParserSerializerTestsuite.xml", adsIO.AdsParserHelper{}},
		{"bacnet", "assets/testing/protocols/bacnet/ParserSerializerTestsuite.xml", bacnetIO.BacnetipParserHelper{}},
		{"df1", "assets/testing/protocols/df1/ParserSerializerTestsuite.xml", df1IO.Df1ParserHelper{}},
		{"eip big endian", "assets/testing/protocols/eip/ParserSerializerTestsuiteBigEndian.xml", eipIO.EipParserHelper{}},
		{"eip little endian", "assets/testing/protocols/eip/ParserSerializerTestsuiteLittleEndian.xml", eipIO.EipParserHelper{}},
		{"firmata", "assets/testing/protocols/firmata/ParserSerializerTestsuite.xml", firmataIO.FirmataParserHelper{}},
		{"iec 60870-5-104", "assets/testing/protocols/iec608705104/IEC-60870-5-104.xml", iecIO.Iec608705104ParserHelper{}},
		{"knxnetip", "assets/testing/protocols/knxnetip/ParserSerializerTestsuite.xml", knxIO.KnxnetipParserHelper{}},
		{"modbus ascii", "assets/testing/protocols/modbus/ascii/ParserSerializerTestsuite.xml", modbusIO.ModbusParserHelper{}},
		{"modbus rtu", "assets/testing/protocols/modbus/rtu/ParserSerializerTestsuite.xml", modbusIO.ModbusParserHelper{}},
		{"modbus tcp", "assets/testing/protocols/modbus/tcp/ParserSerializerTestsuite.xml", modbusIO.ModbusParserHelper{}},
		{"s7", "assets/testing/protocols/s7/ParserSerializerTestsuite.xml", s7IO.S7ParserHelper{}},
		{"slmp", "assets/testing/protocols/slmp/ParserSerializerTestsuite.xml", slmpIO.SlmpParserHelper{}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			testutils.RunParserTruncationTestsuite(t, test.path, test.parser)
		})
	}
}
