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
using org.apache.plc4net.exceptions;
using org.apache.plc4net.spi.drivers;
using org.apache.plc4net.spi.transports;
using org.apache.plc4net.transports.test;
using Xunit;

namespace org.apache.plc4net.drivers.modbus.test
{
    /// <summary>The drivers: protocol codes, transport selection and connection setup.</summary>
    public class DriverTests
    {
        /// <summary>The test transport under a code the drivers do not list, remembering what it opened.</summary>
        private sealed class RecordingTransport : ITransport
        {
            private readonly TestTransport _inner = new TestTransport();

            public TestTransportInstance? Last { get; private set; }

            public string TransportCode => "test";

            public string TransportName => "Recording test transport";

            public ITransportConfiguration CreateConfiguration(IReadOnlyDictionary<string, string> parameters)
                => _inner.CreateConfiguration(parameters);

            public ITransportInstance CreateTransportInstance(
                string transportConfig, ITransportConfiguration configuration)
            {
                var instance = (TestTransportInstance)_inner.CreateTransportInstance(transportConfig, configuration);
                Last = instance;
                return instance;
            }
        }

        private static (ITransportManager Manager, RecordingTransport Transport) NewManager()
        {
            var transport = new RecordingTransport();
            return (new DefaultTransportManager(new ITransport[] { transport }), transport);
        }

        [Fact]
        public void The_TCP_driver_is_modbus_tcp_over_tcp()
        {
            var driver = new ModbusTcpDriver(new DefaultTransportManager());

            Assert.Equal("modbus-tcp", driver.ProtocolCode);
            Assert.Equal("tcp", driver.DefaultTransportCode);
            Assert.Equal(502, ModbusTcpDriver.DefaultPort);
        }

        [Fact]
        public void The_RTU_driver_is_modbus_rtu_over_serial()
        {
            var driver = new ModbusRtuDriver(new DefaultTransportManager());

            Assert.Equal("modbus-rtu", driver.ProtocolCode);
            Assert.Equal("serial", driver.DefaultTransportCode);
        }

        [Fact]
        public void The_drivers_register_their_transports_with_the_manager()
        {
            var manager = new DefaultTransportManager();

            _ = new ModbusTcpDriver(manager);
            _ = new ModbusRtuDriver(manager);

            Assert.Contains("tcp", manager.GetTransportCodes());
            Assert.Contains("serial", manager.GetTransportCodes());
        }

        [Fact]
        public void Scanning_the_assembly_registers_both_drivers_with_the_driver_manager()
        {
            var count = PlcDriverManager.Instance.ScanAndRegisterDrivers(
                typeof(ModbusTcpDriver).Assembly, new DefaultTransportManager());

            Assert.Equal(2, count);
            Assert.IsType<ModbusTcpDriver>(PlcDriverManager.Instance.GetDriverByCode("modbus-tcp"));
            Assert.IsType<ModbusRtuDriver>(PlcDriverManager.Instance.GetDriverByCode("modbus-rtu"));
        }

        [Fact]
        public void Connecting_builds_a_TCP_connection_with_the_requested_unit_identifier()
        {
            var (manager, _) = NewManager();
            var driver = new ModbusTcpDriver(manager);

            using var connection = driver.Connect(
                "modbus-tcp:test://device?allow-unsupported-transport=true&default-unit-identifier=9");

            var tcp = Assert.IsType<ModbusTcpConnection>(connection);
            Assert.Equal((byte)9, tcp.UnitIdentifier);
            Assert.True(connection.IsConnected);
            Assert.True(connection.PlcConnectionMetadata.CanRead);
            Assert.True(connection.PlcConnectionMetadata.CanWrite);
            Assert.False(connection.PlcConnectionMetadata.CanSubscribe);
        }

        [Fact]
        public void Connecting_builds_an_RTU_connection_with_the_requested_slave_address()
        {
            var (manager, _) = NewManager();
            var driver = new ModbusRtuDriver(manager);

            using var connection = driver.Connect(
                "modbus-rtu:test://device?allow-unsupported-transport=true&default-unit-identifier=21");

            var rtu = Assert.IsType<ModbusRtuConnection>(connection);
            Assert.Equal((byte)21, rtu.UnitIdentifier);
        }

        [Fact]
        public void Subscriptions_are_not_offered()
        {
            var (manager, _) = NewManager();
            var driver = new ModbusTcpDriver(manager);

            using var connection = driver.Connect("modbus-tcp:test://device?allow-unsupported-transport=true");

            Assert.Null(connection.SubscriptionRequestBuilder);
            Assert.Null(connection.UnsubscriptionRequestBuilder);
            Assert.NotNull(connection.ReadRequestBuilder);
            Assert.NotNull(connection.WriteRequestBuilder);
        }

        [Fact]
        public void A_transport_the_driver_does_not_list_is_rejected_unless_allowed()
        {
            var (manager, _) = NewManager();
            var driver = new ModbusTcpDriver(manager);

            var exception = Assert.Throws<PlcConnectionException>(
                () => driver.Connect("modbus-tcp:test://device"));

            Assert.Contains("not supported", exception.Message);
        }

        [Fact]
        public void A_driver_refuses_another_protocols_connection_string()
        {
            var driver = new ModbusTcpDriver(new DefaultTransportManager());

            Assert.Throws<PlcConnectionException>(() => driver.Connect("modbus-rtu:serial://COM1"));
        }

        [Theory]
        [InlineData("modbus-tcp", "default-unit-identifier=300")]
        [InlineData("modbus-rtu", "default-unit-identifier=0")]
        [InlineData("modbus-tcp", "request-timeout-ms=never")]
        public void A_bad_unit_identifier_fails_the_connect_and_closes_the_transport(
            string protocol, string parameter)
        {
            var (manager, transport) = NewManager();
            DriverBaseFor(protocol, manager, out var connect);

            Assert.Throws<PlcConnectionException>(
                () => connect($"{protocol}:test://device?allow-unsupported-transport=true&{parameter}"));

            Assert.NotNull(transport.Last);
            Assert.False(transport.Last!.IsOpen); // not leaked
        }

        private static void DriverBaseFor(
            string protocol, ITransportManager manager, out System.Func<string, object> connect)
        {
            if (protocol == "modbus-tcp")
            {
                var driver = new ModbusTcpDriver(manager);
                connect = cs => driver.Connect(cs);
            }
            else
            {
                var driver = new ModbusRtuDriver(manager);
                connect = cs => driver.Connect(cs);
            }
        }
    }
}