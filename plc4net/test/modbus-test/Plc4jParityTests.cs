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
using System.Threading.Tasks;
using org.apache.plc4net.types;
using Xunit;
using static org.apache.plc4net.drivers.modbus.test.ConnectionTestSupport;

namespace org.apache.plc4net.drivers.modbus.test
{
    /// <summary>
    /// The requests and answers of plc4j's own driver test data,
    /// <c>protocols/modbus/src/test/resources/protocols/modbus/tcp/ManualFactoryModbusTCPDriverTest-testsuite.xml</c>,
    /// for the tags this driver supports. The same tag has to put the same PDU on the wire
    /// and decode the same answer; the bytes below are copied from that file (the transaction
    /// identifier and the unit identifier are the connection's own).
    /// </summary>
    public class Plc4jParityTests
    {
        private static byte[] Pdu(string hex) => Convert.FromHexString(hex);

        [Theory]
        // Read holding-register:3:WORD   (request 000500000006 01 0300020001, answer 000500000005 01 0302A5B8)
        [InlineData("holding-register:3:WORD", "0300020001", "0302A5B8", 42424)]
        // Read holding-register:12:INT   (-2424)
        [InlineData("holding-register:12:INT", "03000B0001", "0302F688", -2424)]
        // Read holding-register:13:UINT
        [InlineData("holding-register:13:UINT", "03000C0001", "0302A5B8", 42424)]
        public async Task A_read_puts_the_same_pdu_on_the_wire_and_decodes_the_same_answer(
            string tag, string requestPdu, string answerPdu, int expected)
        {
            var transport = NewTransport();
            transport.InjectTestData(ModbusFrames.Tcp(1, 1, Pdu(answerPdu)));
            var connection = NewTcp(transport);

            var response = await ReadAsync(connection, tag);

            var value = ValueOf(response);
            var actual = expected < 0 ? value.GetShort() : (long)value.GetUshort();
            Assert.Equal(expected, actual);
            Assert.Equal(ModbusFrames.Tcp(1, 1, Pdu(requestPdu)), transport.GetAllWrittenData());
        }

        [Theory]
        // Write holding-register:3:WORD  (request 00060000 0006 01 060002A5B8, the answer echoes it)
        [InlineData("holding-register:3:WORD", 42424, "060002A5B8")]
        // Write holding-register:12:INT
        [InlineData("holding-register:12:INT", -2424, "06000BF688")]
        // Write holding-register:13:UINT
        [InlineData("holding-register:13:UINT", 42424, "06000CA5B8")]
        public async Task A_write_puts_the_same_pdu_on_the_wire(string tag, int value, string requestPdu)
        {
            var transport = NewTransport();
            transport.InjectTestData(ModbusFrames.Tcp(1, 1, Pdu(requestPdu)));
            var connection = NewTcp(transport);

            var response = await WriteAsync(connection, tag, value);

            Assert.Equal(PlcResponseCode.Ok, response.GetResponseCode("v"));
            Assert.Equal(ModbusFrames.Tcp(1, 1, Pdu(requestPdu)), transport.GetAllWrittenData());
        }
    }
}