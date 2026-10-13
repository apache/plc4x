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

using org.apache.plc4net.api.authentication;
using org.apache.plc4net.spi.drivers;
using org.apache.plc4net.spi.transports;
using org.apache.plc4net.transports.cotp;

namespace org.apache.plc4net.drivers.s7
{
    /// <summary>
    /// The Siemens S7 driver, registered under the <c>s7</c> protocol code.
    /// </summary>
    /// <remarks>
    /// Connection string format: <c>s7://192.168.0.1?remote-rack=0&amp;remote-slot=1</c>. S7 runs
    /// over ISO-on-TCP, so the <c>cotp</c> transport is the only one and the default.
    /// </remarks>
    public class S7Driver : DriverBase
    {
        /// <summary>The ISO-on-TCP port (RFC 1006).</summary>
        public const int DefaultPort = 102;

        public S7Driver(ITransportManager transportManager)
            : base(transportManager)
        {
            RegisterTransport(new CotpTransport(DefaultPort));
        }

        public override string ProtocolCode => "s7";

        public override string ProtocolName => "Siemens S7";

        public override string DefaultTransportCode => "cotp";

        // The COTP handshake is part of the protocol, not an optional extra, so a plain transport
        // would never get to speak S7.
        protected override string[] SupportedTransportCodes => new[] { "cotp" };

        protected override ConnectionBase CreateConnection(
            ConnectionString connectionString,
            ITransportInstance transportInstance,
            IPlcAuthentication? authentication)
        {
            var connection = new S7Connection(connectionString, transportInstance);
            connection.Authentication = authentication;
            return connection;
        }
    }
}