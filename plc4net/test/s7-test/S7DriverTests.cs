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

using org.apache.plc4net.exceptions;
using org.apache.plc4net.spi.transports;
using Xunit;

namespace org.apache.plc4net.drivers.s7.test
{
    /// <summary>The driver: what it is registered as, and which connection strings it takes.</summary>
    public class S7DriverTests
    {
        private static S7Driver NewDriver(out DefaultTransportManager transports)
        {
            transports = new DefaultTransportManager();
            return new S7Driver(transports);
        }

        [Fact]
        public void The_driver_is_s7_over_cotp()
        {
            var driver = NewDriver(out var transports);

            Assert.Equal("s7", driver.ProtocolCode);
            Assert.Equal("Siemens S7", driver.ProtocolName);
            Assert.Equal("cotp", driver.DefaultTransportCode);
            Assert.NotNull(transports.GetTransport("cotp"));
        }

        [Fact]
        public void The_iso_on_tcp_port_is_102()
        {
            Assert.Equal(102, S7Driver.DefaultPort);
        }

        [Theory]
        [InlineData("s7:tcp://127.0.0.1")] // S7 needs the COTP session; a plain connection could never speak it
        [InlineData("modbus-tcp://127.0.0.1")]
        [InlineData("")]
        public void A_connection_string_that_the_driver_cannot_serve_is_rejected(string connectionString)
        {
            var driver = NewDriver(out _);

            Assert.Throws<PlcConnectionException>(() => driver.Connect(connectionString));
        }
    }
}