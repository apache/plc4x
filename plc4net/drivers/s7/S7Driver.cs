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

using org.apache.plc4net.api;
using org.apache.plc4net.api.authentication;
using org.apache.plc4net.spi.drivers;
using org.apache.plc4net.spi.transports;
using org.apache.plc4net.transports.cotp;
using org.apache.plc4net.transports.tcp;

namespace org.apache.plc4net.drivers.s7
{
    /// <summary>
    /// The Siemens S7 driver. Registers under the "s7" protocol code with
    /// COTP as the default transport.
    /// </summary>
    public class S7Driver : DriverBase
    {
        public S7Driver(ITransportManager transportManager)
            : base(transportManager)
        {
            RegisterTransport(new TcpTransport());
            RegisterTransport(new CotpTransport());
        }

        public override string ProtocolCode => "s7";

        public override string ProtocolName => "Siemens S7";

        public override string DefaultTransportCode => "cotp";

        // S7 speaks ISO-on-TCP (TPKT + COTP), not raw TCP: the CR/CC handshake is part of
        // the protocol, not an optional enhancement. Java SPI3 declares only cotp for the
        // same reason. The TcpTransport stays registered because CotpTransport wraps it.
        protected override string[] SupportedTransportCodes => new[] { "cotp" };

        protected override ConnectionBase CreateConnection(
            ConnectionString connectionString,
            ITransportInstance transportInstance,
            IPlcAuthentication authentication)
        {
            var connection = new S7Connection(connectionString, transportInstance);
            connection.Authentication = authentication;
            return connection;
        }
    }
}
