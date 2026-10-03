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
using System.Threading;
using System.Threading.Tasks;
using org.apache.plc4net.api;
using org.apache.plc4net.api.value;
using org.apache.plc4net.spi.drivers;
using org.apache.plc4net.spi.drivers.functions;
using org.apache.plc4net.spi.drivers.messages;
using org.apache.plc4net.transports.test;
using org.apache.plc4net.types;
using Xunit;

namespace org.apache.plc4net.drivers.modbus.test
{
    /// <summary>Shared helpers for the connection tests.</summary>
    internal static class ConnectionTestSupport
    {
        public static TestTransportInstance NewTransport()
            => new TestTransportInstance(new TestTransportConfiguration());

        public static ModbusTcpConnection NewTcp(TestTransportInstance transport, string parameters = "")
            => new ModbusTcpConnection(
                ConnectionString.Parse("modbus-tcp://10.0.0.9:502" + parameters), transport);

        public static ModbusRtuConnection NewRtu(TestTransportInstance transport, string parameters = "")
            => new ModbusRtuConnection(
                ConnectionString.Parse("modbus-rtu:serial://COM1" + parameters), transport);

        /// <summary>Reads one tag; the response is keyed "v".</summary>
        public static async Task<DefaultPlcReadResponse> ReadAsync(
            IPlcConnection connection, string address, CancellationToken cancellationToken = default)
        {
            var builder = connection.ReadRequestBuilder!;
            builder.AddTagAddress("v", address);
            var request = (DefaultPlcReadRequest)builder.Build();
            return (DefaultPlcReadResponse)await ((PlcReader)connection).Read(request, cancellationToken);
        }

        /// <summary>Writes one value to one tag; the response is keyed "v".</summary>
        public static async Task<DefaultPlcWriteResponse> WriteAsync<T>(
            IPlcConnection connection, string address, T value, CancellationToken cancellationToken = default)
        {
            var builder = (DefaultPlcWriteRequestBuilder)connection.WriteRequestBuilder!;
            builder.AddTag<T>("v", address, value);
            var request = (DefaultPlcWriteRequest)builder.Build();
            return (DefaultPlcWriteResponse)await ((PlcWriter)connection).Write(request, cancellationToken);
        }

        public static IPlcValue ValueOf(DefaultPlcReadResponse response)
        {
            Assert.Equal(PlcResponseCode.Ok, response.GetResponseCode("v"));
            var value = response.GetValue("v");
            Assert.NotNull(value);
            return value!;
        }

        /// <summary>
        /// Injects the given frames once the driver has written its request. The RTU
        /// connection discards whatever is buffered before it writes, so a response
        /// pre-loaded into the transport would be thrown away.
        /// </summary>
        public static Task InjectAfterRequest(TestTransportInstance transport, params byte[][] frames)
            => InjectAfterRequest(transport, 0, frames);

        /// <summary>
        /// As above, but waits <paramref name="delayMs"/> after the request before injecting, so
        /// whatever was buffered earlier is seen on its own first.
        /// </summary>
        public static Task InjectAfterRequest(TestTransportInstance transport, int delayMs, params byte[][] frames)
        {
            return Task.Run(async () =>
            {
                var deadline = Environment.TickCount64 + 5000;
                while (transport.GetNumBytesWritten() == 0 && Environment.TickCount64 < deadline)
                {
                    await Task.Delay(1);
                }
                if (delayMs > 0)
                {
                    await Task.Delay(delayMs);
                }
                foreach (var frame in frames)
                {
                    transport.InjectTestData(frame);
                }
            });
        }
    }
}