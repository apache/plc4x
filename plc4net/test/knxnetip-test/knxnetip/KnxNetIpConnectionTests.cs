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
using System.Diagnostics;
using System.Net;
using System.Threading;
using System.Threading.Tasks;
using org.apache.plc4net.drivers.knxnetip;
using org.apache.plc4net.spi.drivers;
using org.apache.plc4net.spi.drivers.messages;
using org.apache.plc4net.transports.udp;
using org.apache.plc4net.types;
using Xunit;

namespace org.apache.plc4net.test.knxnetip
{
    /// <summary>
    /// The KNXnet/IP connection state machine against a scripted gateway on a real
    /// loopback UDP socket. Every frame is parsed and serialised by the generated
    /// model, so a round trip here exercises the whole codec.
    /// </summary>
    public class KnxNetIpConnectionTests : IDisposable
    {
        private readonly FakeKnxGateway _gateway = new();

        public void Dispose() => _gateway.Dispose();

        private KnxNetIpConnection Connect()
        {
            var cs = ConnectionString.Parse(
                $"knxnet-ip:udp://127.0.0.1:{_gateway.Port}?request-timeout=2000");
            var transport = new UdpTransportInstance(
                new IPEndPoint(IPAddress.Loopback, _gateway.Port), new UdpTransportConfiguration());
            var connection = new KnxNetIpConnection(cs, transport);
            connection.Connect();
            return connection;
        }

        private static bool WaitFor(Func<bool> condition)
        {
            var sw = Stopwatch.StartNew();
            while (sw.ElapsedMilliseconds < 2000)
            {
                if (condition())
                {
                    return true;
                }
                Thread.Sleep(10);
            }
            return condition();
        }

        [Fact]
        public void The_handshake_completes_and_adopts_the_assigned_client_address()
        {
            using var connection = Connect();
            Assert.True(connection.IsConnected);
            Assert.Equal(FakeKnxGateway.AssignedClientAddress, connection.ClientKnxAddress);
        }

        [Fact]
        public async Task A_bare_group_write_reaches_the_gateway_as_a_group_value_write()
        {
            using var connection = Connect();

            var builder = (DefaultPlcWriteRequestBuilder) connection.WriteRequestBuilder;
            builder.AddTag("light", "2/2/12", (byte) 1);
            var response = (DefaultPlcWriteResponse) await connection.Write(
                (DefaultPlcWriteRequest) builder.Build());

            Assert.Equal(PlcResponseCode.Ok, response.GetResponseCode("light"));
            Assert.True(WaitFor(() => _gateway.ReceivedWrites.Count == 1));
            Assert.Equal("2/2/12", _gateway.ReceivedWrites[0].GroupAddress);
            Assert.Equal(new byte[] { 0x01 }, _gateway.ReceivedWrites[0].Payload);
        }

        [Fact]
        public async Task A_typed_write_encodes_through_the_datapoint_codec()
        {
            using var connection = Connect();

            var builder = (DefaultPlcWriteRequestBuilder) connection.WriteRequestBuilder;
            builder.AddTag("switch", "1/0/1:DPT1.001", true);
            var response = (DefaultPlcWriteResponse) await connection.Write(
                (DefaultPlcWriteRequest) builder.Build());

            Assert.Equal(PlcResponseCode.Ok, response.GetResponseCode("switch"));
            Assert.True(WaitFor(() => _gateway.ReceivedWrites.Count == 1));
            Assert.Equal(new byte[] { 0x01 }, _gateway.ReceivedWrites[0].Payload);
        }

        [Fact]
        public async Task The_first_tunnelling_request_uses_sequence_counter_zero()
        {
            using var connection = Connect();

            var builder = (DefaultPlcWriteRequestBuilder) connection.WriteRequestBuilder;
            builder.AddTag("a", "2/2/12", (byte) 1);
            builder.AddTag("b", "2/2/13", (byte) 1);
            await connection.Write((DefaultPlcWriteRequest) builder.Build());

            Assert.True(WaitFor(() => _gateway.TunnellingSequenceNumbers.Count >= 2));
            // KNX 03/08/04 §2.6: the first request after CONNECT carries counter 0.
            Assert.Equal(new byte[] { 0, 1 }, _gateway.TunnellingSequenceNumbers.GetRange(0, 2).ToArray());
        }

        [Fact]
        public async Task A_tag_whose_level_count_does_not_match_the_connection_is_rejected()
        {
            using var connection = Connect();  // default: 3 levels

            var builder = (DefaultPlcWriteRequestBuilder) connection.WriteRequestBuilder;
            builder.AddTag("two-level", "1/2", (byte) 1);
            var response = (DefaultPlcWriteResponse) await connection.Write(
                (DefaultPlcWriteRequest) builder.Build());

            Assert.Equal(PlcResponseCode.InvalidAddress, response.GetResponseCode("two-level"));
        }

        [Fact]
        public async Task A_raw_write_larger_than_six_bits_is_rejected()
        {
            using var connection = Connect();

            var builder = (DefaultPlcWriteRequestBuilder) connection.WriteRequestBuilder;
            builder.AddTag("big", "2/2/12", (byte) 200);
            var response = (DefaultPlcWriteResponse) await connection.Write(
                (DefaultPlcWriteRequest) builder.Build());

            Assert.Equal(PlcResponseCode.InternalError, response.GetResponseCode("big"));
            Assert.Empty(_gateway.ReceivedWrites);
        }

        [Fact]
        public async Task A_read_correlates_the_group_value_response_from_the_bus()
        {
            using var connection = Connect();
            // DPT 5.001 (percent) - one data byte after the APDU first byte.
            _gateway.NextReadResponsePayload = new byte[] { 0x00, 0x80 };

            var builder = (DefaultPlcReadRequestBuilder) connection.ReadRequestBuilder;
            builder.AddTagAddress("dim", "2/2/12:DPT5.001");
            var response = (DefaultPlcReadResponse) await connection.Read(
                (DefaultPlcReadRequest) builder.Build());

            Assert.Equal(PlcResponseCode.Ok, response.GetResponseCode("dim"));
            Assert.Equal((byte) 0x80, response.GetValue("dim").GetByte());
        }

        [Fact]
        public async Task A_read_with_no_device_answering_reports_not_found()
        {
            using var connection = Connect();
            _gateway.NextReadResponsePayload = null; // stay silent

            var builder = (DefaultPlcReadRequestBuilder) connection.ReadRequestBuilder;
            builder.AddTagAddress("dim", "2/2/13:DPT5.001");
            var response = (DefaultPlcReadResponse) await connection.Read(
                (DefaultPlcReadRequest) builder.Build());

            Assert.Equal(PlcResponseCode.NotFound, response.GetResponseCode("dim"));
        }

        [Fact]
        public void A_bus_listener_sees_an_unsolicited_group_write()
        {
            using var connection = Connect();
            connection.RegisterDatapointHint("2/2/12", "DPT1.001");

            var events = new List<KnxGroupValueEvent>();
            connection.RegisterGroupValueListener(e =>
            {
                lock (events)
                {
                    events.Add(e);
                }
            });

            _gateway.PushGroupWrite(new byte[] { 0x12, 0x0C }, new byte[] { 0x01 });

            Assert.True(WaitFor(() =>
            {
                lock (events)
                {
                    return events.Count == 1;
                }
            }));
            Assert.Equal("2/2/12", events[0].GroupAddress);
            Assert.Equal(KnxGroupValueEventType.Write, events[0].EventType);
            Assert.Equal("1.1.5", events[0].SourceAddress);
            Assert.True(events[0].Value.GetBool());
        }
    }
}
