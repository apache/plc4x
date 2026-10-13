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
using System.Collections.Generic;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;
using org.apache.plc4net.api;
using org.apache.plc4net.api.value;
using org.apache.plc4net.spi.drivers;
using org.apache.plc4net.spi.drivers.functions;
using org.apache.plc4net.spi.drivers.messages;
using org.apache.plc4net.transports.cotp;
using org.apache.plc4net.transports.test;
using org.apache.plc4net.types;
using Xunit;

namespace org.apache.plc4net.drivers.s7.test
{
    /// <summary>Shared helpers for the connection tests.</summary>
    internal static class ConnectionTestSupport
    {
        public const string S7Uri = "s7://192.168.0.1?remote-rack=0&remote-slot=1";

        public static TestTransportInstance NewTransport()
            => new TestTransportInstance(new TestTransportConfiguration());

        /// <summary>
        /// A transport that already holds what a CPU answers to the connection: the COTP confirm
        /// and the answer to Setup Communication, which carries TPDU reference 1.
        /// </summary>
        public static TestTransportInstance Ready(ushort pduLength = 240)
        {
            var transport = NewTransport();
            transport.InjectTestData(S7Frames.ConnectionConfirm());
            transport.InjectTestData(S7Frames.SetupAck(1, pduLength));
            return transport;
        }

        /// <summary>Connects over the transport; what the driver wrote while connecting stays in it.</summary>
        public static S7Connection Connect(TestTransportInstance transport, string uri = S7Uri)
        {
            var connection = new S7Connection(ConnectionString.Parse(uri), new CotpTransportInstance(transport));
            connection.Connect();
            return connection;
        }

        /// <summary>Connects and forgets what was written, so that a test sees only its own requests.</summary>
        public static S7Connection ConnectQuietly(TestTransportInstance transport, string uri = S7Uri)
        {
            var connection = Connect(transport, uri);
            transport.GetAllWrittenData();
            return connection;
        }

        /// <summary>Cuts the bytes the driver wrote into the TPKT frames they are made of.</summary>
        public static List<byte[]> Frames(byte[] written)
        {
            var frames = new List<byte[]>();
            var at = 0;
            while (at + 4 <= written.Length)
            {
                var length = (written[at + 2] << 8) | written[at + 3];
                frames.Add(written[at..(at + length)]);
                at += length;
            }
            Assert.Equal(written.Length, at);
            return frames;
        }

        /// <summary>
        /// A frame with the TPDU number of its COTP Data Transfer header set to 0. The transport
        /// numbers the frames it sends 0, 1, 2 and so on, which is not what these tests are about.
        /// </summary>
        public static byte[] Normalized(byte[] frame)
        {
            var copy = (byte[])frame.Clone();
            if (copy.Length > 6 && copy[5] == 0xF0)
            {
                copy[6] = 0x80;
            }
            return copy;
        }

        /// <summary>What the driver has written since the last time, with the Data Transfer frames normalized.</summary>
        public static byte[] Written(TestTransportInstance transport)
            => Frames(transport.GetAllWrittenData()).SelectMany(Normalized).ToArray();

        /// <summary>
        /// Plays the CPU for a test whose connection sends several requests: whenever the driver has
        /// written a whole frame, it takes the frame and injects the next answer, a null for none.
        /// An answer that was there before its request would be taken for a stale one and thrown
        /// away, as it would be if a real CPU could answer a request that it has not got yet.
        /// </summary>
        public sealed class FakeCpu
        {
            private readonly List<byte[]> _requests = new List<byte[]>();

            public FakeCpu(TestTransportInstance transport, params byte[]?[] answers)
                : this(transport, 0, answers)
            {
            }

            /// <param name="transport">The transport the driver writes to and reads from.</param>
            /// <param name="delayMs">How long the CPU takes to answer, so that whatever was buffered
            /// before is seen on its own first.</param>
            /// <param name="answers">The answers, in the order of the requests.</param>
            public FakeCpu(TestTransportInstance transport, int delayMs, params byte[]?[] answers)
            {
                Running = Task.Run(async () =>
                {
                    var pending = new List<byte>();
                    var next = 0;
                    var deadline = Environment.TickCount64 + 8000;
                    while (next < answers.Length && Environment.TickCount64 < deadline)
                    {
                        pending.AddRange(transport.GetAllWrittenData());
                        while (pending.Count >= 4 && pending.Count >= ((pending[2] << 8) | pending[3]))
                        {
                            var length = (pending[2] << 8) | pending[3];
                            var frame = pending.GetRange(0, length).ToArray();
                            pending.RemoveRange(0, length);
                            lock (_requests)
                            {
                                _requests.Add(Normalized(frame));
                            }
                            var answer = next < answers.Length ? answers[next++] : null;
                            if (answer != null)
                            {
                                if (delayMs > 0)
                                {
                                    await Task.Delay(delayMs);
                                }
                                transport.InjectTestData(answer);
                            }
                        }
                        await Task.Delay(1);
                    }
                });
            }

            /// <summary>Completes when every answer has been given.</summary>
            public Task Running { get; }

            /// <summary>The requests that the driver sent, in order, with the Data Transfer frames normalized.</summary>
            public IReadOnlyList<byte[]> Requests
            {
                get
                {
                    lock (_requests)
                    {
                        return _requests.ToList();
                    }
                }
            }
        }

        /// <summary>Reads the given tags in one request.</summary>
        public static async Task<DefaultPlcReadResponse> ReadAsync(
            IPlcConnection connection, params (string Name, string Address)[] tags)
        {
            var builder = connection.ReadRequestBuilder!;
            foreach (var (name, address) in tags)
            {
                builder.AddTagAddress(name, address);
            }
            var request = (DefaultPlcReadRequest)builder.Build();
            return (DefaultPlcReadResponse)await ((PlcReader)connection).Read(request);
        }

        /// <summary>Reads one tag; the response is keyed "v".</summary>
        public static Task<DefaultPlcReadResponse> ReadAsync(IPlcConnection connection, string address)
            => ReadAsync(connection, ("v", address));

        /// <summary>Writes the given values in one request.</summary>
        public static async Task<DefaultPlcWriteResponse> WriteAsync(
            IPlcConnection connection, params (string Name, string Address, object Value)[] values)
        {
            var builder = (DefaultPlcWriteRequestBuilder)connection.WriteRequestBuilder!;
            foreach (var (name, address, value) in values)
            {
                builder.AddTag<object>(name, address, value);
            }
            var request = (DefaultPlcWriteRequest)builder.Build();
            return (DefaultPlcWriteResponse)await ((PlcWriter)connection).Write(request);
        }

        /// <summary>Writes one value; the response is keyed "v".</summary>
        public static Task<DefaultPlcWriteResponse> WriteAsync(IPlcConnection connection, string address, object value)
            => WriteAsync(connection, ("v", address, value));

        public static IPlcValue ValueOf(DefaultPlcReadResponse response, string name = "v")
        {
            Assert.Equal(PlcResponseCode.Ok, response.GetResponseCode(name));
            var value = response.GetValue(name);
            Assert.NotNull(value);
            return value!;
        }
    }
}
