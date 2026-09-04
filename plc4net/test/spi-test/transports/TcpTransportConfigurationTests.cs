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

using org.apache.plc4net.spi.drivers;
using org.apache.plc4net.transports.tcp;
using Xunit;

namespace org.apache.plc4net.spi.test.transports
{
    /// <summary>
    /// Covers how the TCP transport reads its options out of a connection string.
    /// </summary>
    /// <remarks>
    /// The option names and the prefixing rule both come from the Java side: the
    /// transport declares its options via <c>@ConfigurationParameter</c> and
    /// <c>ConfigurationFactory.createPrefixedConfiguration</c> builds the prefix from
    /// the transport code, then strips it before matching. So an option Java names
    /// "tcp-no-delay" is addressed as "tcp.tcp-no-delay" in a connection string - the
    /// repetition is the transport code meeting an option name that happens to start
    /// with the same word, not a typo.
    /// </remarks>
    public class TcpTransportConfigurationTests
    {
        private static TcpTransportConfiguration Configure(string connectionString)
        {
            var parsed = ConnectionString.Parse(connectionString);
            return Assert.IsType<TcpTransportConfiguration>(
                new TcpTransport().CreateConfiguration(parsed.Parameters));
        }

        [Fact]
        public void Defaults_apply_when_no_parameters_are_given()
        {
            var config = Configure("modbus-tcp://10.0.0.9:502");

            Assert.True(config.TcpNoDelay);
            Assert.True(config.KeepAlive);
            Assert.Equal(5000, config.ConnectTimeout);
            Assert.Equal(81920, config.ReceiveBufferSize);
            Assert.Equal(-1, config.DefaultPort);
        }

        [Fact]
        public void The_transport_prefixed_form_is_read()
        {
            var config = Configure("modbus-tcp://10.0.0.9:502?tcp.tcp-no-delay=false");

            Assert.False(config.TcpNoDelay);
        }

        [Fact]
        public void The_unprefixed_form_is_read_too()
        {
            // Not something Java accepts, but harmless and convenient for drivers that
            // only ever speak over one transport.
            var config = Configure("modbus-tcp://10.0.0.9:502?tcp-no-delay=false");

            Assert.False(config.TcpNoDelay);
        }

        [Fact]
        public void The_prefixed_form_wins_over_the_unprefixed_one()
        {
            var config = Configure(
                "modbus-tcp://10.0.0.9:502?tcp-no-delay=true&tcp.tcp-no-delay=false");

            Assert.False(config.TcpNoDelay);
        }

        [Theory]
        [InlineData("?tcp.connect-timeout=1234")]
        [InlineData("?connect-timeout=1234")]
        public void Numeric_options_are_read_in_either_form(string parameters)
        {
            var config = Configure("modbus-tcp://10.0.0.9:502" + parameters);

            Assert.Equal(1234, config.ConnectTimeout);
        }

        [Fact]
        public void Several_options_are_read_from_one_connection_string()
        {
            var config = Configure(
                "s7:tcp://10.0.0.5:102?tcp.connect-timeout=250&tcp.keep-alive=false" +
                "&tcp.send-buffer-size=4096&tcp.local-address=127.0.0.1&tcp.local-port=9000");

            Assert.Equal(250, config.ConnectTimeout);
            Assert.False(config.KeepAlive);
            Assert.Equal(4096, config.SendBufferSize);
            Assert.Equal("127.0.0.1", config.LocalAddress);
            Assert.Equal(9000, config.LocalPort);
        }

        [Fact]
        public void An_unparseable_value_falls_back_to_the_default()
        {
            // A bad value must not take down the connection attempt before the driver
            // has had a chance to report anything more useful.
            var config = Configure("modbus-tcp://10.0.0.9:502?tcp.connect-timeout=soon");

            Assert.Equal(5000, config.ConnectTimeout);
        }

        [Fact]
        public void A_zero_receive_buffer_size_falls_back_to_the_default()
        {
            // A non-positive value would crash the RingBuffer constructor, which
            // requires a positive capacity.
            var config = Configure("modbus-tcp://10.0.0.9:502?tcp.receive-buffer-size=0");

            Assert.Equal(81920, config.ReceiveBufferSize);
        }

        [Fact]
        public void A_negative_receive_buffer_size_falls_back_to_the_default()
        {
            var config = Configure("modbus-tcp://10.0.0.9:502?tcp.receive-buffer-size=-100");

            Assert.Equal(81920, config.ReceiveBufferSize);
        }

        [Fact]
        public void Parameters_of_other_transports_are_ignored()
        {
            var config = Configure("modbus-tcp://10.0.0.9:502?serial.baud-rate=9600");

            Assert.Equal(5000, config.ConnectTimeout);
            Assert.True(config.TcpNoDelay);
        }
    }
}
