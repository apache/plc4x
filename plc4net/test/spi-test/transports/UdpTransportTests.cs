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
using System.Net;
using System.Net.Sockets;
using System.Threading;
using org.apache.plc4net.spi.transports;
using org.apache.plc4net.transports.udp;
using Xunit;

namespace org.apache.plc4net.spi.test.transports
{
    /// <summary>
    /// Address parsing plus a round trip against a real loopback UDP socket acting
    /// as the peer.
    /// </summary>
    public class UdpTransportTests : IDisposable
    {
        private const int SettleMillis = 2000;

        private readonly Socket _peer;
        private readonly IPEndPoint _peerEndPoint;

        public UdpTransportTests()
        {
            _peer = new Socket(AddressFamily.InterNetwork, SocketType.Dgram, ProtocolType.Udp);
            _peer.Bind(new IPEndPoint(IPAddress.Loopback, 0));
            _peerEndPoint = (IPEndPoint) _peer.LocalEndPoint;
        }

        public void Dispose() => _peer.Dispose();

        private static bool WaitFor(Func<bool> condition)
        {
            var deadline = Stopwatch.StartNew();
            while (deadline.ElapsedMilliseconds < SettleMillis)
            {
                if (condition())
                {
                    return true;
                }
                Thread.Sleep(10);
            }
            return condition();
        }

        [Theory]
        [InlineData("192.168.1.10:3671", "192.168.1.10", 3671)]
        [InlineData("192.168.1.10", "192.168.1.10", 3671)]          // default port
        [InlineData("gw.example:3671/knx", "gw.example", 3671)]     // driver config stripped
        public void Parses_address_forms(string address, string host, int port)
        {
            var (h, p, _) = UdpTransport.ParseAddress(address, 3671);
            Assert.Equal(host, h);
            Assert.Equal(port, p);
        }

        [Fact]
        public void An_address_without_a_port_and_no_default_is_rejected()
        {
            Assert.Throws<TransportException>(() => UdpTransport.ParseAddress("192.168.1.10", -1));
        }

        [Fact]
        public void Datagrams_round_trip_through_the_buffer()
        {
            using var instance = new UdpTransportInstance(_peerEndPoint, new UdpTransportConfiguration());

            // client -> peer
            instance.Write(new byte[] { 0x06, 0x10, 0x02, 0x03, 0x00, 0x06 });
            var rx = new byte[64];
            var from = (EndPoint) new IPEndPoint(IPAddress.Any, 0);
            var n = _peer.ReceiveFrom(rx, ref from);
            Assert.Equal(6, n);
            Assert.Equal(0x10, rx[1]);

            // peer -> client, appended to the ring buffer and surfaced through Read
            _peer.SendTo(new byte[] { 0x06, 0x10, 0x02, 0x06, 0x00, 0x08, 0xAA, 0xBB }, from);
            Assert.True(WaitFor(() => instance.GetNumBytesAvailable() >= 8),
                "the datagram never reached the transport buffer");
            var frame = instance.Read(8);
            Assert.Equal(new byte[] { 0x06, 0x10, 0x02, 0x06, 0x00, 0x08, 0xAA, 0xBB }, frame);
        }

        [Fact]
        public void The_data_listener_fires_when_a_datagram_arrives()
        {
            using var instance = new UdpTransportInstance(_peerEndPoint, new UdpTransportConfiguration());
            var fired = 0;
            instance.RegisterDataListener(() => Interlocked.Increment(ref fired));

            instance.Write(new byte[] { 1 });
            var rx = new byte[8];
            var from = (EndPoint) new IPEndPoint(IPAddress.Any, 0);
            _peer.ReceiveFrom(rx, ref from);
            _peer.SendTo(new byte[] { 9, 9, 9 }, from);

            Assert.True(WaitFor(() => Volatile.Read(ref fired) > 0), "the data listener never fired");
        }

        [Fact]
        public void Close_releases_the_socket_and_stops_the_loop()
        {
            var instance = new UdpTransportInstance(_peerEndPoint, new UdpTransportConfiguration());
            Assert.True(instance.IsOpen);

            instance.Close();

            Assert.False(instance.IsOpen);
            Assert.True(WaitFor(() => instance.ResourcesReleased), "resources were never released");
            Assert.Throws<TransportException>(() => instance.Write(new byte[] { 1 }));
        }

        [Fact]
        public void The_local_endpoint_is_exposed_for_the_knx_hpai()
        {
            using var instance = new UdpTransportInstance(_peerEndPoint, new UdpTransportConfiguration());
            Assert.NotNull(instance.LocalAddress);
            Assert.NotEqual(0, instance.LocalAddress.Port);
        }
    }
}
