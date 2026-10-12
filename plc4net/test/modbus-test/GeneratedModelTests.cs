//
// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//      https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
//

using System;
using System.Globalization;
using System.Linq;
using org.apache.plc4net.drivers.modbus.readwrite.model;
using org.apache.plc4net.spi.drivers;
using org.apache.plc4net.spi.generation;
using org.apache.plc4net.tools.codegen;
using Xunit;

namespace org.apache.plc4net.drivers.modbus.test
{
    /// <summary>
    /// Runs the ParserSerializerTestsuite vectors that plc4j and plc4go validate
    /// against through the checked-in, generated Modbus model. A vector passes
    /// when its bytes parse and serialize back byte-identical.
    /// </summary>
    public class GeneratedModelTests
    {
        private static readonly string[] Protocols = { "tcp", "rtu", "ascii" };

        public static TheoryData<string, string> Vectors()
        {
            var data = new TheoryData<string, string>();
            foreach (var protocol in Protocols)
            {
                foreach (var testCase in LoadSuite(protocol).TestCases)
                {
                    data.Add(protocol, testCase.Name);
                }
            }
            return data;
        }

        [Theory]
        [MemberData(nameof(Vectors))]
        public void Vector_round_trips_byte_identical(string protocol, string name)
        {
            var testCase = LoadSuite(protocol).TestCases.Single(t => t.Name == name);
            var expected = testCase.GetRawBytes();

            var adu = Parse(testCase);

            var writeBuffer = new WriteBuffer();
            adu.Serialize(writeBuffer);

            Assert.Equal(ToHex(expected), ToHex(writeBuffer.GetBytes()));
        }

        [Fact]
        public void The_shared_suites_are_actually_loaded()
        {
            // Guards against a path or layout change turning the theory above into
            // zero cases, which xunit reports as a pass.
            Assert.True(LoadSuite("tcp").TestCases.Count >= 7);
            Assert.True(LoadSuite("rtu").TestCases.Count >= 2);
            Assert.True(LoadSuite("ascii").TestCases.Count >= 2);
        }

        [Fact]
        public void A_corrupted_RTU_frame_fails_the_checksum()
        {
            var testCase = LoadSuite("rtu").TestCases.First(t => t.Name.Contains("Response"));
            var bytes = testCase.GetRawBytes();
            bytes[^1] ^= 0xFF; // flip the last CRC byte

            var exception = Assert.Throws<ParseException>(() =>
                ModbusADU.StaticParse(new ReadBuffer(bytes), DriverType.MODBUS_RTU, true));
            Assert.Contains("crc", exception.Message, StringComparison.OrdinalIgnoreCase);
        }

        [Fact]
        public void The_discriminated_type_resolves_to_the_expected_concrete_classes()
        {
            var testCase = LoadSuite("tcp").TestCases.Single(t => t.Name == "Read Input Registers Request");

            var adu = Assert.IsType<ModbusTcpADU>(Parse(testCase));
            var pdu = Assert.IsType<ModbusPDUReadInputRegistersRequest>(adu.Pdu);
            Assert.Equal((ushort)2258, pdu.StartingAddress);
        }

        private static ModbusADU Parse(ParserSerializerTestcase testCase)
        {
            Assert.Equal("ModbusADU", testCase.RootType);
            var driverType = Enum.Parse<DriverType>(testCase.ParserArguments["driverType"], ignoreCase: true);
            var response = bool.Parse(testCase.ParserArguments["response"]);
            return ModbusADU.StaticParse(new ReadBuffer(testCase.GetRawBytes()), driverType, response);
        }

        private static ParserSerializerTestsuite LoadSuite(string protocol)
            => ParserSerializerTestsuiteRunner.Load(RepoPaths.ParserSerializerTestsuite(protocol));

        private static string ToHex(byte[] bytes)
            => string.Concat(bytes.Select(b => b.ToString("x2", CultureInfo.InvariantCulture)));
    }
}