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

using System.Text;
using System.Threading.Tasks;
using org.apache.plc4net.types;
using Xunit;
using static org.apache.plc4net.drivers.modbus.test.ConnectionTestSupport;

namespace org.apache.plc4net.drivers.modbus.test
{
    public class AsciiConnectionTests
    {
        [Fact]
        public async Task A_register_read_uses_ascii_hex_and_lrc()
        {
            var transport = NewTransport();
            var connection = NewAscii(transport);
            var inject = InjectAfterRequest(
                transport, Encoding.ASCII.GetBytes(":0103021234B4\r\n"));

            var response = await ReadAsync(connection, "holding-register:1:UINT");
            await inject;

            Assert.Equal((ushort)0x1234, ValueOf(response).GetUshort());
            Assert.Equal(":010300000001FB\r\n",
                Encoding.ASCII.GetString(transport.GetAllWrittenData()));
        }

        [Fact]
        public async Task A_bad_lrc_is_ignored_until_timeout()
        {
            var transport = NewTransport();
            var connection = NewAscii(transport, "?request-timeout-ms=50");
            var inject = InjectAfterRequest(
                transport, Encoding.ASCII.GetBytes(":010302123400\r\n"));

            var response = await ReadAsync(connection, "holding-register:1:UINT");
            await inject;

            Assert.Equal(PlcResponseCode.RequestTimeout, response.GetResponseCode("v"));
        }
    }
}
