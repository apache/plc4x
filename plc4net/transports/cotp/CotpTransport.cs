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

using System.Collections.Generic;
using org.apache.plc4net.spi.transports;
using org.apache.plc4net.transports.tcp;

namespace org.apache.plc4net.transports.cotp
{
    /// <summary>
    /// The "cotp" transport — wraps a TCP connection with TPKT (RFC 1006)
    /// framing and COTP (ISO 8073 class 0) connection handling for use by the
    /// S7 driver.
    /// </summary>
    /// <remarks>
    /// The transport instances created here do the COTP work themselves instead
    /// of leaving it to the driver: <c>CotpTransportInstance.Open</c> exchanges
    /// the Connection Request and Confirm for the TSAPs the driver supplies and
    /// negotiates the TPDU size, <c>Write</c> splits a payload into Data
    /// Transfer TPDUs that fit that size, and <c>Read</c> returns the payload of
    /// the received Data Transfer TPDUs with the TPKT and COTP headers
    /// stripped. The driver only deals in the payload bytes above COTP.
    /// </remarks>
    public class CotpTransport : ITransport
    {
        private readonly TcpTransport _tcp;

        public CotpTransport() : this(-1)
        {
        }

        /// <param name="defaultPort">Port to use when the connection string omits one (102 for S7).</param>
        public CotpTransport(int defaultPort)
        {
            _tcp = new TcpTransport(defaultPort);
        }

        public string TransportCode => "cotp";

        public string TransportName => "COTP (ISO-on-TCP) Transport";

        public ITransportConfiguration CreateConfiguration(
            IReadOnlyDictionary<string, string> parameters)
        {
            return _tcp.CreateConfiguration(parameters);
        }

        public ITransportInstance CreateTransportInstance(
            string transportConfig, ITransportConfiguration configuration)
        {
            var tcpInstance = _tcp.CreateTransportInstance(transportConfig, configuration);
            return new CotpTransportInstance(tcpInstance);
        }
    }
}
