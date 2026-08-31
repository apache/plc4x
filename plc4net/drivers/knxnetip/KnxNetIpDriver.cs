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
using org.apache.plc4net.transports.udp;

namespace org.apache.plc4net.drivers.knxnetip
{
    /// <summary>
    /// The KNXnet/IP driver. Registers under the "knxnetip" protocol code with UDP
    /// tunnelling as the transport.
    /// </summary>
    /// <remarks>
    /// Connection string: <c>knxnetip:udp://192.168.1.10:3671?knx-group-address-num-levels=3</c>.
    /// The gateway port defaults to 3671. <c>knx-group-address-num-levels</c>
    /// (1, 2 or 3; default 3) fixes the bit layout for every group address on the
    /// connection.
    /// </remarks>
    public class KnxNetIpDriver : DriverBase
    {
        public KnxNetIpDriver(ITransportManager transportManager)
            : base(transportManager)
        {
            RegisterTransport(new UdpTransport());
        }

        public override string ProtocolCode => "knxnetip";

        public override string ProtocolName => "KNXnet/IP";

        public override string DefaultTransportCode => "udp";

        protected override string[] SupportedTransportCodes => new[] { "udp" };

        protected override ConnectionBase CreateConnection(
            ConnectionString connectionString,
            ITransportInstance transportInstance,
            IPlcAuthentication authentication)
        {
            var connection = new KnxNetIpConnection(connectionString, transportInstance);
            connection.Authentication = authentication;
            return connection;
        }
    }
}
