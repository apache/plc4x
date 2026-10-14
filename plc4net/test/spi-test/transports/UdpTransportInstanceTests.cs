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
    public class UdpTransportInstanceTests
    {
        private const int SettleMillis = 2000;

        private static bool WaitFor(Func<bool> condition)
        {
            var deadline = Stopwatch.StartNew();
            while (deadline.ElapsedMilliseconds < SettleMillis)
            {
                if (condition()) return true;
                Thread.Sleep(5);
            }
            return condition();
        }

        private static UdpTransportInstance Connect(UdpClient peer)
        {
            var peerAddress = (IPEndPoint)peer.Client.LocalEndPoint!;
            return new UdpTransportInstance(
                new IPEndPoint(IPAddress.Loopback, peerAddress.Port), new UdpTransportConfiguration());
        }

        [Fact]
        public void Reads_a_datagram_and_notifies_the_listener()
        {
            using var peer = new UdpClient(new IPEndPoint(IPAddress.Loopback, 0));
            using var instance = Connect(peer);
            using var notified = new ManualResetEventSlim();
            instance.RegisterDataListener(() => notified.Set());
            var bytes = new byte[] { 0x01, 0x02, 0x03 };

            peer.Send(bytes, bytes.Length, instance.LocalAddress);

            Assert.True(notified.Wait(SettleMillis), "the data listener never fired");
            Assert.True(WaitFor(() => instance.GetNumBytesAvailable() == bytes.Length));
            Assert.Equal(bytes, instance.Read(bytes.Length));
        }

        [Fact]
        public void Writes_one_datagram_to_the_connected_peer()
        {
            using var peer = new UdpClient(new IPEndPoint(IPAddress.Loopback, 0));
            using var instance = Connect(peer);
            peer.Client.ReceiveTimeout = SettleMillis;
            var bytes = new byte[] { 0xAA, 0xBB };

            instance.Write(bytes);

            var source = new IPEndPoint(IPAddress.Any, 0);
            Assert.Equal(bytes, peer.Receive(ref source));
        }

        [Fact]
        public void Close_is_idempotent_and_releases_the_socket()
        {
            using var peer = new UdpClient(new IPEndPoint(IPAddress.Loopback, 0));
            using var instance = Connect(peer);

            instance.Close();
            instance.Close();

            Assert.False(instance.IsOpen);
            Assert.True(WaitFor(() => instance.ResourcesReleased));
        }

        [Theory]
        [InlineData("-1")]
        [InlineData("65536")]
        public void Configuration_rejects_a_local_port_outside_the_socket_range(string localPort)
        {
            var transport = new UdpTransport();

            var exception = Assert.Throws<TransportException>(() => transport.CreateConfiguration(
                new System.Collections.Generic.Dictionary<string, string> { ["local-port"] = localPort }));

            Assert.Contains("local-port", exception.Message);
        }

        [Fact]
        public void Configuration_accepts_prefixed_options_and_clamps_the_receive_buffer()
        {
            var transport = new UdpTransport(3671);
            var configuration = Assert.IsType<UdpTransportConfiguration>(transport.CreateConfiguration(
                new System.Collections.Generic.Dictionary<string, string>
                {
                    ["udp.local-address"] = "127.0.0.1",
                    ["udp.local-port"] = "12345",
                    ["udp.receive-buffer-size"] = "1"
                }));

            Assert.Equal(3671, configuration.DefaultPort);
            Assert.Equal("127.0.0.1", configuration.LocalAddress);
            Assert.Equal(12345, configuration.LocalPort);
            Assert.Equal(65535, configuration.ReceiveBufferSize);
        }
    }
}