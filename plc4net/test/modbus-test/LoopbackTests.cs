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
using System.Diagnostics;
using System.Threading.Tasks;
using org.apache.plc4net.spi.transports;
using org.apache.plc4net.types;
using Xunit;
using static org.apache.plc4net.drivers.modbus.test.ConnectionTestSupport;

namespace org.apache.plc4net.drivers.modbus.test
{
    /// <summary>
    /// The Modbus TCP driver over real sockets: <see cref="ModbusTcpDriver"/> and the TCP
    /// transport against <see cref="LoopbackModbusSlave"/>. The slave's tables are indexed
    /// by wire address, and a tag addresses the register one higher.
    /// </summary>
    public class LoopbackTests
    {
        private static ModbusTcpConnection Connect(LoopbackModbusSlave slave, string parameters = "")
        {
            var driver = new ModbusTcpDriver(new DefaultTransportManager());
            return Assert.IsType<ModbusTcpConnection>(driver.Connect(
                $"modbus-tcp://127.0.0.1:{slave.Port}?default-unit-identifier=1&request-timeout-ms=3000{parameters}"));
        }

        [Fact]
        public async Task Registers_and_coils_can_be_read_and_written()
        {
            using var slave = new LoopbackModbusSlave();
            slave.Registers[100] = 0x1234;
            slave.Registers[101] = 0x5678;
            slave.Coils[7] = true;
            using var connection = Connect(slave);

            Assert.Equal((ushort)0x1234, ValueOf(await ReadAsync(connection, "holding-register:101:UINT")).GetUshort());
            Assert.Equal((ushort)0x5678, ValueOf(await ReadAsync(connection, "input-register:102:UINT")).GetUshort());
            Assert.True(ValueOf(await ReadAsync(connection, "coil:8")).GetBool());
            Assert.False(ValueOf(await ReadAsync(connection, "discrete-input:9")).GetBool());

            var register = await WriteAsync<ushort>(connection, "holding-register:6:UINT", 0xABCD);
            var coil = await WriteAsync(connection, "coil:10", true);
            Assert.Equal(PlcResponseCode.Ok, register.GetResponseCode("v"));
            Assert.Equal(PlcResponseCode.Ok, coil.GetResponseCode("v"));
            Assert.Equal((ushort)0xABCD, slave.Registers[5]);
            Assert.True(slave.Coils[9]);
            Assert.Equal((ushort)0xABCD, ValueOf(await ReadAsync(connection, "holding-register:6:UINT")).GetUshort());
        }

        [Fact]
        public async Task A_connection_can_be_opened_through_the_driver_manager_the_way_a_user_does()
        {
            using var slave = new LoopbackModbusSlave();
            slave.Registers[9] = 0x0123;
            PlcDriverManager.Instance.RegisterDriver(new ModbusTcpDriver(new DefaultTransportManager()));

            using var connection = PlcDriverManager.Instance.GetConnection(
                $"modbus-tcp://127.0.0.1:{slave.Port}?request-timeout-ms=3000");

            Assert.IsType<ModbusTcpConnection>(connection);
            Assert.Equal((short)0x0123, ValueOf(await ReadAsync(connection, "holding-register:10")).GetShort());
        }

        [Fact]
        public async Task A_signed_register_is_read_and_written_as_twos_complement()
        {
            using var slave = new LoopbackModbusSlave();
            using var connection = Connect(slave);

            var written = await WriteAsync<short>(connection, "holding-register:4", -1234);
            var read = await ReadAsync(connection, "holding-register:4");

            Assert.Equal(PlcResponseCode.Ok, written.GetResponseCode("v"));
            Assert.Equal((ushort)(65536 - 1234), slave.Registers[3]);
            Assert.Equal((short)-1234, ValueOf(read).GetShort());
        }

        [Fact]
        public async Task A_device_exception_is_reported_and_the_connection_carries_on()
        {
            using var slave = new LoopbackModbusSlave();
            slave.Registers[1] = 0x0042;
            using var connection = Connect(slave);

            var failed = await ReadAsync(connection, $"holding-register:{LoopbackModbusSlave.FirstIllegalAddress + 51}");
            var next = await ReadAsync(connection, "holding-register:2");

            Assert.Equal(PlcResponseCode.InvalidAddress, failed.GetResponseCode("v"));
            Assert.Equal((short)0x0042, ValueOf(next).GetShort());
        }

        [Fact]
        public async Task A_response_forwarded_one_byte_at_a_time_is_assembled()
        {
            using var slave = new LoopbackModbusSlave();
            slave.Registers[LoopbackModbusSlave.SlowAddress] = 0x4AFE;
            using var connection = Connect(slave);

            var response = await ReadAsync(connection, $"holding-register:{LoopbackModbusSlave.SlowAddress + 1}");

            Assert.Equal((short)0x4AFE, ValueOf(response).GetShort());
        }

        [Fact]
        public async Task A_late_answer_to_a_timed_out_request_is_skipped_by_the_next_one()
        {
            using var slave = new LoopbackModbusSlave();
            slave.Registers[LoopbackModbusSlave.LateAddress] = 0x1111;
            slave.Registers[2] = 0x2222;
            using var connection = Connect(slave, "&request-timeout-ms=100"); // the later parameter wins

            var timedOut = await ReadAsync(connection, $"holding-register:{LoopbackModbusSlave.LateAddress + 1}");
            Assert.Equal(PlcResponseCode.RequestTimeout, timedOut.GetResponseCode("v"));

            // The first request's answer arrives now, ahead of the answer to this one.
            await Task.Delay(LoopbackModbusSlave.LateDelayMs + 200);
            var next = await ReadAsync(connection, "holding-register:3");

            Assert.Equal((short)0x2222, ValueOf(next).GetShort());
        }

        [Fact]
        public async Task A_device_that_disappears_fails_the_read_instead_of_hanging()
        {
            using var slave = new LoopbackModbusSlave();
            slave.Registers[1] = 0x0001;
            using var connection = Connect(slave, "&request-timeout-ms=1000");
            Assert.Equal((short)1, ValueOf(await ReadAsync(connection, "holding-register:2")).GetShort());

            slave.DropClients();
            var clock = Stopwatch.StartNew();
            var response = await ReadAsync(connection, "holding-register:2");

            Assert.NotEqual(PlcResponseCode.Ok, response.GetResponseCode("v"));
            Assert.True(clock.Elapsed < TimeSpan.FromSeconds(10), $"took {clock.Elapsed}");
        }
    }
}