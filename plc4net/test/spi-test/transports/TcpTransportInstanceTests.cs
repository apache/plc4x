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
using org.apache.plc4net.transports.tcp;
using Xunit;

namespace org.apache.plc4net.spi.test.transports
{
    /// <summary>
    /// Connection lifecycle tests against a real loopback socket.
    /// </summary>
    /// <remarks>
    /// Everything else in this suite is a pure unit test over parsing or buffering.
    /// Those are the parts that are easy to test, and they are not where the defects
    /// were: the read loop, the locks and the teardown paths only misbehave once there
    /// is a peer that can hang up. A local <see cref="TcpListener"/> is enough to reach
    /// them and costs a few milliseconds.
    /// </remarks>
    public class TcpTransportInstanceTests : IDisposable
    {
        private const int SettleMillis = 2000;

        private readonly TcpListener _listener;
        private readonly IPEndPoint _endPoint;

        public TcpTransportInstanceTests()
        {
            // Port 0 lets the OS pick a free one, so parallel test runs cannot collide.
            _listener = new TcpListener(IPAddress.Loopback, 0);
            _listener.Start();
            _endPoint = (IPEndPoint) _listener.LocalEndpoint;
        }

        public void Dispose()
        {
            _listener.Stop();
        }

        private TcpTransportInstance Connect(out Socket serverSide)
        {
            var accepting = _listener.AcceptSocketAsync();
            var instance = new TcpTransportInstance(_endPoint, new TcpTransportConfiguration());
            Assert.True(accepting.Wait(SettleMillis), "the listener never accepted the connection");
            serverSide = accepting.Result;
            return instance;
        }

        /// <summary>Spins until <paramref name="condition"/> holds, or gives up.</summary>
        private static bool WaitFor(Func<bool> condition)
        {
            var deadline = Stopwatch.StartNew();
            while (deadline.ElapsedMilliseconds < SettleMillis)
            {
                if (condition())
                {
                    return true;
                }
                Thread.Sleep(5);
            }
            return condition();
        }

        [Fact]
        public void Reads_what_the_peer_sends()
        {
            using (var instance = Connect(out var serverSide))
            {
                var notified = new ManualResetEventSlim();
                instance.RegisterDataListener(() => notified.Set());

                serverSide.Send(new byte[] { 0x01, 0x02, 0x03, 0x04 });

                Assert.True(notified.Wait(SettleMillis), "the data listener never fired");
                Assert.True(WaitFor(() => instance.GetNumBytesAvailable() >= 4));
                Assert.Equal(new byte[] { 0x01, 0x02, 0x03, 0x04 }, instance.Read(4));

                serverSide.Dispose();
            }
        }

        [Fact]
        public void Writes_reach_the_peer()
        {
            using (var instance = Connect(out var serverSide))
            {
                instance.Write(new byte[] { 0xAA, 0xBB });

                var received = new byte[2];
                Assert.Equal(2, serverSide.Receive(received));
                Assert.Equal(new byte[] { 0xAA, 0xBB }, received);

                serverSide.Dispose();
            }
        }

        [Fact]
        public void A_peer_close_notifies_and_releases_the_socket()
        {
            using (var instance = Connect(out var serverSide))
            {
                var disconnected = new ManualResetEventSlim();
                instance.RegisterDisconnectListener(_ => disconnected.Set());

                serverSide.Shutdown(SocketShutdown.Both);
                serverSide.Dispose();

                Assert.True(disconnected.Wait(SettleMillis), "the disconnect listener never fired");
                Assert.True(WaitFor(() => !instance.IsOpen));

                // The point of the test: the peer hanging up has to release our end too.
                Assert.True(WaitFor(() => instance.ResourcesReleased), "the socket was never released");
            }
        }

        [Fact]
        public void Buffered_data_survives_a_peer_close()
        {
            using (var instance = Connect(out var serverSide))
            {
                serverSide.Send(new byte[] { 0x10, 0x20 });
                Assert.True(WaitFor(() => instance.GetNumBytesAvailable() >= 2));

                serverSide.Shutdown(SocketShutdown.Both);
                serverSide.Dispose();
                Assert.True(WaitFor(() => instance.ResourcesReleased));

                // Only the codec knows where a frame ends, so whatever arrived before the
                // peer went away must still be drainable.
                Assert.Equal(new byte[] { 0x10, 0x20 }, instance.Read(2));
            }
        }

        [Fact]
        public void Close_from_inside_the_disconnect_listener_completes()
        {
            using (var instance = Connect(out var serverSide))
            {
                var closed = new ManualResetEventSlim();
                instance.RegisterDisconnectListener(_ =>
                {
                    instance.Close();
                    closed.Set();
                });

                serverSide.Shutdown(SocketShutdown.Both);
                serverSide.Dispose();

                // The listener runs on the read loop's own thread, so this Close() must not
                // sit waiting for that loop to finish.
                Assert.True(closed.Wait(SettleMillis), "Close() from the disconnect listener did not return");
                Assert.True(instance.ResourcesReleased);
            }
        }

        [Fact]
        public void Close_is_idempotent()
        {
            using (var instance = Connect(out var serverSide))
            {
                instance.Close();
                Assert.False(instance.IsOpen);
                Assert.True(instance.ResourcesReleased);

                // Repeat calls, and a Dispose() on top, must all be harmless.
                instance.Close();
                instance.Dispose();
                instance.Dispose();

                serverSide.Dispose();
            }
        }

        [Fact]
        public void Writing_to_a_closed_connection_is_rejected()
        {
            using (var instance = Connect(out var serverSide))
            {
                instance.Close();

                Assert.Throws<TransportException>(() => instance.Write(new byte[] { 0x01 }));

                serverSide.Dispose();
            }
        }
    }
}
