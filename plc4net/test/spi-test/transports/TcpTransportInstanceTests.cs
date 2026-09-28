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
using System.Threading.Tasks;
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
            return Connect(new TcpTransportConfiguration(), out serverSide);
        }

        private TcpTransportInstance Connect(TcpTransportConfiguration configuration, out Socket serverSide)
        {
            var accepting = _listener.AcceptSocketAsync();
            var instance = new TcpTransportInstance(_endPoint, configuration);
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

        /// <summary>Whether <paramref name="task"/> completes within the given time.</summary>
        private static async Task<bool> CompletesWithin(Task task, int millis)
        {
            return await Task.WhenAny(task, Task.Delay(millis)) == task;
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
        public async Task Close_does_not_wait_for_a_data_listener_that_is_still_running()
        {
            using (var instance = Connect(out var serverSide))
            {
                var entered = new ManualResetEventSlim();
                var release = new ManualResetEventSlim();
                instance.RegisterDataListener(() =>
                {
                    entered.Set();
                    // Keeps the read loop inside the callback for as long as the test needs.
                    release.Wait(10 * SettleMillis);
                });

                serverSide.Send(new byte[] { 0x01 });
                Assert.True(entered.Wait(SettleMillis), "the data listener never fired");

                // The read loop is now stuck in the callback. Close() has to return without
                // waiting for it: closing the socket is what ends the loop, not the other
                // way round. A bounded wait inside Close() shows up here as the wait.
                var watch = Stopwatch.StartNew();
                instance.Close();
                watch.Stop();

                Assert.True(instance.ResourcesReleased);
                Assert.True(watch.ElapsedMilliseconds < 500,
                    $"Close() waited {watch.ElapsedMilliseconds} ms for the read loop");

                release.Set();
                Assert.True(await CompletesWithin(instance.ReadLoopCompleted, SettleMillis), "the read loop never ended");

                serverSide.Dispose();
            }
        }

        [Fact]
        public async Task Closing_our_own_end_is_not_reported_as_a_disconnect()
        {
            using (var instance = Connect(out var serverSide))
            {
                var disconnected = new ManualResetEventSlim();
                instance.RegisterDisconnectListener(_ => disconnected.Set());

                instance.Close();

                // Our own Shutdown comes back to the read loop as a zero-byte read, which
                // looks just like the peer hanging up. Wait for the loop to end so that a
                // notification, if one were coming, has had every chance to arrive.
                Assert.True(await CompletesWithin(instance.ReadLoopCompleted, SettleMillis), "the read loop never ended");
                Assert.False(disconnected.IsSet, "a local Close() was reported as a disconnect");

                serverSide.Dispose();
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

        [Fact]
        public void A_refused_connection_is_reported_as_a_transport_error()
        {
            // A port that was free a moment ago: bind it, note the number, let it go.
            var probe = new TcpListener(IPAddress.Loopback, 0);
            probe.Start();
            var closedPort = ((IPEndPoint) probe.LocalEndpoint).Port;
            probe.Stop();

            var e = Assert.Throws<TransportException>(() => new TcpTransportInstance(
                new IPEndPoint(IPAddress.Loopback, closedPort), new TcpTransportConfiguration()));

            Assert.Contains($"Failed to connect to 127.0.0.1:{closedPort}", e.Message);
            Assert.IsType<SocketException>(e.InnerException);
        }

        [Fact]
        public async Task A_connect_with_a_cancelled_token_is_never_attempted()
        {
            using (var cancelled = new CancellationTokenSource())
            {
                cancelled.Cancel();

                await Assert.ThrowsAnyAsync<OperationCanceledException>(() =>
                    TcpTransportInstance.ConnectAsync(_endPoint, new TcpTransportConfiguration(), cancelled.Token));

                // The token is what aborts an attempt that runs out of time, so it has to
                // reach the connect itself: nothing may have arrived at the listener.
                Assert.False(_listener.Pending(), "a connection was made although the token was already cancelled");
            }
        }

        [Fact]
        public async Task Close_silences_a_write_blocked_on_a_full_send_buffer()
        {
            // Keep the kernel from absorbing the payload: the explicit send buffer
            // sizes the transport's side, and the shrunken listener receive buffer
            // carries over to the accepted side on Windows (Linux sizes accepted
            // sockets with its own defaults instead). Either way the buffers stay
            // small next to 2000 chunks, and the peer never reads, so repeated
            // writes wedge inside Socket.Send.
            _listener.Server.ReceiveBufferSize = 1024;

            using (var instance = Connect(new TcpTransportConfiguration
            {
                SendBufferSize = 1024,
            }, out var serverSide))
            {
                var chunk = new byte[64 * 1024];
                long written = 0;
                using var stop = new CancellationTokenSource();
                var writer = Task.Run(() =>
                {
                    for (var i = 0; i < 2000 && !stop.IsCancellationRequested; i++)
                    {
                        instance.Write(chunk);
                        Interlocked.Add(ref written, chunk.Length);
                    }
                });

                // Wait until the writer has wedged inside a blocked Send: the peer
                // never reads, so progress stops while Write is still in flight.
                // Progress is not required to have been made first: Linux honors the
                // small buffers so strictly that the very first chunk may block.
                var wedge = Stopwatch.StartNew();
                long last = -1;
                var stalled = 0;
                var wedged = false;
                while (wedge.ElapsedMilliseconds < 5000)
                {
                    var current = Interlocked.Read(ref written);
                    if (current == last)
                    {
                        stalled += 50;
                        if (stalled >= 500)
                        {
                            wedged = true;
                            break;
                        }
                    }
                    else
                    {
                        stalled = 0;
                        last = current;
                    }
                    Thread.Sleep(50);
                }

                try
                {
                    Assert.True(wedged, "the writer never wedged inside a blocked Send");
                    Assert.True(Interlocked.Read(ref written) < 2000L * chunk.Length,
                        "every write completed although the peer never reads, so nothing was in flight to interrupt");
                }
                finally
                {
                    // Closing now aborts the blocked Send. The caller asked for the
                    // close, so the write it interrupts must end silently instead of
                    // surfacing a write error, and the writer must stop of its own
                    // accord. This has to run even when the asserts above fail, so the
                    // writer cannot keep looping into a disposed instance.
                    stop.Cancel();
                    instance.Close();
                }

                Assert.True(await CompletesWithin(writer, SettleMillis), "the in-flight write never finished");

                // A send failure (the fix absent) surfaces as a TransportException
                // wrapping the socket error, two levels below the task's own wrapper.
                // A write refused at the entry guard because the close raced the
                // loop-top cancellation check faults without a socket error; that is
                // our own close at work, not a failed send, so it must not fail
                // this test.
                Assert.True(writer.Exception?.InnerException?.InnerException is not SocketException,
                    $"the in-flight write failed with {writer.Exception?.InnerException?.InnerException}");

                serverSide.Dispose();
            }
        }
    }
}