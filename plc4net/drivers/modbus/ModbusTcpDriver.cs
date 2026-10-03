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
        }

        public override string ProtocolCode => "modbus-tcp";

        public override string ProtocolName => "Modbus TCP";

        public override string DefaultTransportCode => "tcp";

        protected override string[] SupportedTransportCodes => new[] { "tcp" };

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