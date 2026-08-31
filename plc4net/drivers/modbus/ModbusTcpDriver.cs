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
using org.apache.plc4net.transports.tcp;

namespace org.apache.plc4net.drivers.modbus
{
    /// <summary>
    /// The Modbus TCP driver. Registers itself under the "modbus-tcp" protocol code.
    /// </summary>
    public class ModbusTcpDriver : DriverBase
    {
        public ModbusTcpDriver(ITransportManager transportManager)
            : base(transportManager)
        {
            RegisterTransport(new TcpTransport());
        }

        public override string ProtocolCode => "modbus-tcp";

        public override string ProtocolName => "Modbus TCP";

        public override string DefaultTransportCode => "tcp";

        protected override string[] SupportedTransportCodes => new[] { "tcp" };

        protected override ConnectionBase CreateConnection(
            ConnectionString connectionString,
            ITransportInstance transportInstance,
            IPlcAuthentication authentication)
        {
            var connection = new ModbusConnection(connectionString, transportInstance);
            connection.Authentication = authentication;
            return connection;
        }
    }
}
