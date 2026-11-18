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
using org.apache.plc4net.transports.tcp;

namespace org.apache.plc4net.drivers.modbus
{
    /// <summary>
    /// The Modbus TCP driver, registered under the <c>modbus-tcp</c> protocol code.
    /// </summary>
    /// <remarks>
    /// Connection string format:
    /// <c>modbus-tcp://10.0.0.9:502?default-unit-identifier=1&amp;request-timeout-ms=5000</c>
    /// </remarks>
    public class ModbusTcpDriver : DriverBase
    {
        /// <summary>The well-known Modbus TCP port.</summary>
        public const int DefaultPort = 502;

        public ModbusTcpDriver(ITransportManager transportManager)
            : base(transportManager)
        {
            RegisterTransport(new TcpTransport(DefaultPort));
            RegisterTransport(new org.apache.plc4net.transports.udp.UdpTransport(502));
            RegisterTransport(new org.apache.plc4net.transports.tcp.TlsTransport());
        }

        public override string ProtocolCode => "modbus-tcp";

        public override string ProtocolName => "Modbus TCP";

        public override string DefaultTransportCode => "tcp";

        /// <summary>Probes one explicitly supplied endpoint using basic device identification.</summary>
        public async System.Threading.Tasks.Task<System.Collections.Generic.IReadOnlyDictionary<byte, byte[]>> DiscoverAsync(
            string connectionString, System.Threading.CancellationToken cancellationToken = default)
        {
            cancellationToken.ThrowIfCancellationRequested();
            using var connection = (ModbusTcpConnection)Connect(connectionString);
            return await connection.IdentifyAsync(cancellationToken).ConfigureAwait(false);
        }

        protected override string[] SupportedTransportCodes => new[] { "tcp", "udp", "tls" };

        /// <summary>Identifies explicitly supplied endpoints sequentially; failures remain visible per endpoint.</summary>
        public async System.Threading.Tasks.Task<System.Collections.Generic.IReadOnlyList<(string Endpoint,
            System.Collections.Generic.IReadOnlyDictionary<byte, byte[]>? Objects, System.Exception? Error)>> DiscoverAsync(
            System.Collections.Generic.IEnumerable<string> endpoints, System.Threading.CancellationToken cancellationToken = default)
        {
            var results = new System.Collections.Generic.List<(string,
                System.Collections.Generic.IReadOnlyDictionary<byte, byte[]>?, System.Exception?)>();
            foreach (var endpoint in endpoints)
            {
                cancellationToken.ThrowIfCancellationRequested();
                try { results.Add((endpoint, await DiscoverAsync(endpoint, cancellationToken).ConfigureAwait(false), null)); }
                catch (System.OperationCanceledException) { throw; }
                catch (System.Exception ex) { results.Add((endpoint, null, ex)); }
            }
            return results;
        }

        protected override ConnectionBase CreateConnection(
            ConnectionString connectionString,
            ITransportInstance transportInstance,
            IPlcAuthentication? authentication)
        {
            var connection = new ModbusTcpConnection(connectionString, transportInstance);
            connection.Authentication = authentication;
            return connection;
        }
    }
}
