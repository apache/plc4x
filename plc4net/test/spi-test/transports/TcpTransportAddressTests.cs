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
using org.apache.plc4net.spi.transports;
using Xunit;

namespace org.apache.plc4net.spi.test.transports
{
    public class TcpTransportAddressTests
    {
        [Fact]
        public void Parses_host_and_port()
        {
            var (host, port, driverConfig) =
                org.apache.plc4net.transports.tcp.TcpTransport.ParseAddress("192.168.0.1:102", -1);

            Assert.Equal("192.168.0.1", host);
            Assert.Equal(102, port);
            Assert.Equal("", driverConfig);
        }

        [Fact]
        public void Host_without_port_uses_the_default()
        {
            var (host, port, _) =
                org.apache.plc4net.transports.tcp.TcpTransport.ParseAddress("192.168.0.1", 502);

            Assert.Equal(502, port);
        }

        [Fact]
        public void Host_without_port_and_no_default_throws()
        {
            Assert.Throws<TransportException>(
                () => org.apache.plc4net.transports.tcp.TcpTransport.ParseAddress("192.168.0.1", -1));
        }

        [Fact]
        public void Parses_IPv6_literal()
        {
            var (host, port, _) =
                org.apache.plc4net.transports.tcp.TcpTransport.ParseAddress("[::1]:502", -1);

            Assert.Equal("::1", host);
            Assert.Equal(502, port);
        }

        [Fact]
        public void Extracts_the_driver_config()
        {
            var (host, port, driverConfig) =
                org.apache.plc4net.transports.tcp.TcpTransport.ParseAddress("localhost:4840/milo", -1);

            Assert.Equal("localhost", host);
            Assert.Equal(4840, port);
            Assert.Equal("/milo", driverConfig);
        }

        [Fact]
        public void Empty_address_throws()
        {
            Assert.Throws<TransportException>(
                () => org.apache.plc4net.transports.tcp.TcpTransport.ParseAddress("", 502));
        }

        [Fact]
        public void Port_out_of_range_throws()
        {
            Assert.Throws<TransportException>(
                () => org.apache.plc4net.transports.tcp.TcpTransport.ParseAddress("host:99999", -1));
        }
    }
}
