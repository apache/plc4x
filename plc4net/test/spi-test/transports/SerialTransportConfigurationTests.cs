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

using System.IO.Ports;
using org.apache.plc4net.spi.drivers;
using org.apache.plc4net.spi.transports;
using org.apache.plc4net.transports.serial;
using Xunit;

namespace org.apache.plc4net.spi.test.transports
{
    public class SerialTransportConfigurationTests
    {
        private static SerialTransportConfiguration Configure(string connectionString)
        {
            var parsed = ConnectionString.Parse(connectionString);
            return Assert.IsType<SerialTransportConfiguration>(
                new SerialTransport().CreateConfiguration(parsed.Parameters));
        }

        [Fact]
        public void Defaults_match_the_modbus_rtu_convention()
        {
            var config = Configure("modbus-rtu:serial://COM1");

            Assert.Equal(19200, config.BaudRate);
            Assert.Equal(8, config.DataBits);
            Assert.Equal(Parity.Even, config.Parity);
            Assert.Equal(StopBits.One, config.StopBits);
            Assert.Equal(Handshake.None, config.Handshake);
        }

        [Fact]
        public void Prefixed_connection_string_parameters_are_applied()
        {
            var config = Configure(
                "modbus-rtu:serial://COM1?serial.baud-rate=9600&serial.data-bits=7" +
                "&serial.parity=Even&serial.stop-bits=One");

            Assert.Equal(9600, config.BaudRate);
            Assert.Equal(7, config.DataBits);
            Assert.Equal(Parity.Even, config.Parity);
            Assert.Equal(StopBits.One, config.StopBits);
        }

        [Fact]
        public void Prefixed_parameters_take_precedence_over_unprefixed_parameters()
        {
            var config = Configure(
                "modbus-rtu:serial://COM1?baud-rate=19200&serial.baud-rate=9600");

            Assert.Equal(9600, config.BaudRate);
        }

        [Theory]
        [InlineData("serial.baud-rate=0")]
        [InlineData("serial.data-bits=4")]
        [InlineData("serial.data-bits=9")]
        [InlineData("serial.receive-buffer-size=0")]
        public void Invalid_transport_settings_fail_before_opening_the_port(string parameter)
        {
            Assert.Throws<TransportException>(
                () => Configure("modbus-rtu:serial://COM1?" + parameter));
        }

        [Theory]
        [InlineData("serial.parity=99")]
        [InlineData("serial.stop-bits=99")]
        [InlineData("serial.handshake=99")]
        public void Undefined_enum_values_fall_back_to_the_safe_default(string parameter)
        {
            var config = Configure("modbus-rtu:serial://COM1?" + parameter);

            Assert.Equal(Parity.Even, config.Parity);
            Assert.Equal(StopBits.One, config.StopBits);
            Assert.Equal(Handshake.None, config.Handshake);
        }

        [Fact]
        public void Stop_bits_none_is_rejected_before_opening_the_port()
        {
            Assert.Throws<TransportException>(
                () => Configure("modbus-rtu:serial://COM1?serial.stop-bits=None"));
        }
    }
}
