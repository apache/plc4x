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
using System.IO;
using System.Net;
using System.Net.Sockets;
using System.Threading;
using System.Threading.Tasks;
using org.apache.plc4net.spi.transports;

namespace org.apache.plc4net.transports.tcp
{
    /// <summary>
    /// TCP transport: one background read loop per connection feeding a ring buffer.
    /// </summary>
    /// <remarks>
    /// The Java SPI3 transport uses a virtual thread per connection blocked in a socket
    /// read. .NET has no virtual threads, but async I/O gives the same result: the read
    /// loop awaits <c>ReceiveAsync</c> and holds no thread while idle. The observable
    /// contract is the same — the loop fills the ring buffer and then invokes the data
    /// listener.
    /// </remarks>
    public class TcpTransportInstance : BaseTransportInstance, IAsyncTransportInstance
    {
        private readonly Socket _socket;
        private readonly RingBuffer _ringBuffer;
        private readonly object _readLock = new object();
        private readonly object _writeLock = new object();
        private readonly CancellationTokenSource _shutdown = new CancellationTokenSource();
        private readonly Task _readLoop;

        // Whether the connection is usable. Cleared by Close() and by the read loop
        // when the peer goes away. Deliberately separate from _released: the peer can
        // end the connection without anyone having released our end of it yet.
        private int _open = 1;

        // One-shot guard for ReleaseResources(). Every path out of the connection ends
        // there, and only the first one through does the work.
        private int _released;

        private volatile Action _dataListener;
        private volatile Action<Exception> _disconnectListener;

        // Managed thread currently running a listener callback, or 0. ReleaseResources()
        // consults this to recognise being called from inside the read loop.
        private int _listenerThreadId;

        public TcpTransportInstance(IPEndPoint remoteAddress, TcpTransportConfiguration configuration)
            : base(configuration)
        {
            if (remoteAddress == null)
            {
                throw new ArgumentNullException(nameof(remoteAddress));
            }

            _ringBuffer = new RingBuffer(configuration.ReceiveBufferSize);

            Socket socket = null;
            try
            {
                socket = new Socket(remoteAddress.AddressFamily, SocketType.Stream, ProtocolType.Tcp);

                if (!string.IsNullOrEmpty(configuration.LocalAddress))
                {
                    socket.Bind(new IPEndPoint(IPAddress.Parse(configuration.LocalAddress), configuration.LocalPort));
                }

                socket.NoDelay = configuration.TcpNoDelay;
                socket.SetSocketOption(SocketOptionLevel.Socket, SocketOptionName.KeepAlive, configuration.KeepAlive);
                if (configuration.SendBufferSize > 0)
                {
                    socket.SendBufferSize = configuration.SendBufferSize;
                }
                if (configuration.ReceiveBufferSize > 0)
                {
                    socket.ReceiveBufferSize = configuration.ReceiveBufferSize;
                }

                Connect(socket, remoteAddress, configuration.ConnectTimeout);

                _socket = socket;
                RemoteAddress = remoteAddress;
                LocalAddress = socket.LocalEndPoint as IPEndPoint;
            }
            catch (Exception e)
            {
                socket?.Dispose();
                throw new TransportException(
                    $"Failed to connect to {remoteAddress.Address}:{remoteAddress.Port} - {e.Message}", e);
            }

            // Started last so a throw during setup cannot leak a running loop.
            _readLoop = Task.Run(RunReadLoopAsync);
        }

        public IPEndPoint RemoteAddress { get; }

        public IPEndPoint LocalAddress { get; }

        /// <summary>
        /// Connects with a bounded wait. ConnectAsync has no timeout overload, so the
        /// timeout is applied by racing it and disposing the socket on expiry — closing
        /// the socket is what aborts an in-flight connect.
        /// </summary>
        private static void Connect(Socket socket, IPEndPoint remoteAddress, int connectTimeoutMillis)
        {
            var connectTask = socket.ConnectAsync(remoteAddress);
            if (!connectTask.Wait(connectTimeoutMillis))
            {
                throw new TimeoutException(
                    $"Connection to {remoteAddress.Address}:{remoteAddress.Port} timed out after {connectTimeoutMillis} ms.");
            }
            // Surface any connect error (Wait already completed, so this cannot block).
            connectTask.GetAwaiter().GetResult();
        }

        public override bool IsOpen => Volatile.Read(ref _open) == 1 && _socket.Connected;

        /// <summary>
        /// Whether the socket has been shut down and disposed. Exposed so tests can tell
        /// "the peer went away" apart from "our end was actually released" - the two used
        /// to be the same flag, which is how the release came to be skipped.
        /// </summary>
        internal bool ResourcesReleased => Volatile.Read(ref _released) == 1;

        /// <summary>
        /// Optional diagnostic output. When set, raw bytes written and read
        /// are logged as annotated hex.
        /// </summary>
        public TextWriter? DiagnosticOutput { get; set; }

        public override int GetNumBytesAvailable()
        {
            lock (_readLock)
            {
                return IsOpen || _ringBuffer.AvailableForReading > 0
                    ? _ringBuffer.AvailableForReading
                    : 0;
            }
        }

        public override byte[] PeekReadableBytes(int numBytes)
        {
            if (numBytes <= 0)
            {
                return Array.Empty<byte>();
            }
            lock (_readLock)
            {
                if (_ringBuffer.AvailableForReading < numBytes)
                {
                    throw new TransportException(
                        $"Requested {numBytes} bytes but only {_ringBuffer.AvailableForReading} available");
                }
                return _ringBuffer.Peek(numBytes);
            }
        }

        public override byte[] Read(int numBytes)
        {
            if (numBytes <= 0)
            {
                return Array.Empty<byte>();
            }
            byte[] result;
            lock (_readLock)
            {
                if (_ringBuffer.AvailableForReading < numBytes)
                {
                    throw new TransportException(
                        $"Requested {numBytes} bytes but only {_ringBuffer.AvailableForReading} available");
                }
                result = _ringBuffer.Read(numBytes);
            }
            WriteHex("TCP RECV", result);
            return result;
        }

        public override void Write(byte[] bytes)
        {
            if (bytes == null || bytes.Length == 0)
            {
                return;
            }

            lock (_writeLock)
            {
                EnsureOpen();
                try
                {
                    var offset = 0;
                    while (offset < bytes.Length)
                    {
                        var sent = _socket.Send(bytes, offset, bytes.Length - offset, SocketFlags.None);
                        if (sent == 0)
                        {
                            // A blocking Send only returns 0 once the connection is gone.
                            // Looping again would spin forever, since offset cannot move.
                            throw new TransportException(
                                $"Connection closed after writing {offset} of {bytes.Length} bytes.");
                        }
                        offset += sent;
                    }
                    WriteHex("TCP SEND", bytes);
                }
                catch (ObjectDisposedException) when (Volatile.Read(ref _open) == 0)
                {
                    // A concurrent Close() disposed the socket mid-write: normal shutdown.
                }
                catch (SocketException e)
                {
                    throw new TransportException("Failed to write data", e);
                }
            }
        }

        public override void Close()
        {
            Volatile.Write(ref _open, 0);
            ReleaseResources(fromReadLoop: false);
        }

        /// <summary>
        /// Shuts the socket down and disposes it, exactly once.
        /// </summary>
        /// <remarks>
        /// Reached from three directions: an explicit <see cref="Close"/>, a
        /// <see cref="BaseTransportInstance.Dispose()"/>, or the read loop finishing —
        /// which happens on its own whenever the peer hangs up. Keeping the guard here
        /// rather than on the open flag is what makes the last of those release
        /// anything: the loop clears the open flag as it leaves, so a Close() arriving
        /// afterwards would otherwise see an already-closed connection and return
        /// without ever disposing the socket.
        /// </remarks>
        /// <param name="fromReadLoop">
        /// True when the read loop itself is the caller, in which case waiting for that
        /// loop to finish would only burn the timeout.
        /// </param>
        private void ReleaseResources(bool fromReadLoop)
        {
            if (Interlocked.CompareExchange(ref _released, 1, 0) != 0)
            {
                return;
            }

            // Deliberately takes no locks: closing the socket is what unblocks the read
            // loop. Taking _writeLock first would deadlock against a blocked writer.
            _shutdown.Cancel();
            try
            {
                _socket.Shutdown(SocketShutdown.Both);
            }
            catch (SocketException)
            {
                // Already torn down by the peer; nothing to do.
            }
            catch (ObjectDisposedException)
            {
            }
            finally
            {
                _socket.Dispose();
                // Do not wait on the read loop when we are standing inside it - it cannot
                // finish while we are, so the wait would just burn its timeout. That
                // happens two ways: the loop itself calling in on its way out, and a
                // listener callback (which the loop invokes synchronously) calling
                // Close(). Task.CurrentId cannot answer either: the read loop is an async
                // method, and its continuations after the first await do not run as the
                // Task returned by Task.Run, so Task.CurrentId is null there.
                var insideReadLoop = fromReadLoop
                    || Volatile.Read(ref _listenerThreadId) == Environment.CurrentManagedThreadId;
                if (!insideReadLoop && _readLoop != null && !_readLoop.IsCompleted)
                {
                    _readLoop.Wait(TimeSpan.FromSeconds(1));
                }
                _shutdown.Dispose();
            }
        }

        public void RegisterDataListener(Action listener) => _dataListener = listener;

        public void RemoveDataListener() => _dataListener = null;

        public void RegisterDisconnectListener(Action<Exception> listener) => _disconnectListener = listener;

        public void RemoveDisconnectListener() => _disconnectListener = null;

        private void NotifyDisconnect(Exception cause)
        {
            var listener = _disconnectListener;
            if (listener == null)
            {
                return;
            }
            Volatile.Write(ref _listenerThreadId, Environment.CurrentManagedThreadId);
            try
            {
                listener(cause);
            }
            catch
            {
                // A misbehaving listener must not take down the read loop.
            }
            finally
            {
                Volatile.Write(ref _listenerThreadId, 0);
            }
        }

        private void NotifyData()
        {
            var listener = _dataListener;
            if (listener == null)
            {
                return;
            }
            Volatile.Write(ref _listenerThreadId, Environment.CurrentManagedThreadId);
            try
            {
                listener();
            }
            catch
            {
                // Same reasoning as NotifyDisconnect.
            }
            finally
            {
                Volatile.Write(ref _listenerThreadId, 0);
            }
        }

        private async Task RunReadLoopAsync()
        {
            var buffer = new byte[8192];
            try
            {
                while (Volatile.Read(ref _open) == 1)
                {
                    int free;
                    lock (_readLock)
                    {
                        free = _ringBuffer.RemainingForWriting;
                    }
                    if (free == 0)
                    {
                        // Backpressure: the codec has not drained yet. Wait and retry
                        // rather than disconnect — only the codec knows frame boundaries.
                        await Task.Delay(1, _shutdown.Token).ConfigureAwait(false);
                        continue;
                    }

                    var toRead = Math.Min(buffer.Length, free);
                    var bytesRead = await _socket
                        .ReceiveAsync(new ArraySegment<byte>(buffer, 0, toRead), SocketFlags.None)
                        .ConfigureAwait(false);

                    if (bytesRead == 0)
                    {
                        // Orderly shutdown by the remote end.
                        Volatile.Write(ref _open, 0);
                        NotifyDisconnect(null);
                        return;
                    }

                    lock (_readLock)
                    {
                        _ringBuffer.Write(buffer, 0, bytesRead);
                    }

                    // Notify outside the lock; the bytes are already buffered.
                    NotifyData();
                }
            }
            catch (OperationCanceledException)
            {
                // Close() cancelled us: normal shutdown.
            }
            catch (ObjectDisposedException)
            {
                // Close() disposed the socket while we were awaiting a receive.
            }
            catch (Exception e)
            {
                // Only a failure while still open counts as a disconnect.
                if (Volatile.Read(ref _open) == 1)
                {
                    Volatile.Write(ref _open, 0);
                    NotifyDisconnect(e);
                }
            }
            finally
            {
                // However this loop ends, our end of the connection is finished with.
                // When Close() started the teardown this is a no-op; when the peer did,
                // this is the only thing that releases the socket. Clear _open so that
                // IsOpen does not read _socket.Connected on a disposed socket, even if a
                // future catch block forgets to.
                Volatile.Write(ref _open, 0);
                ReleaseResources(fromReadLoop: true);
            }
        }

        private void WriteHex(string label, byte[] bytes)
        {
            var diag = DiagnosticOutput;
            if (diag == null) return;
            diag.WriteLine($"┌ {label} ({bytes.Length} bytes)");
            for (int i = 0; i < bytes.Length; i += 16)
            {
                diag.Write("│ ");
                for (int j = i; j < Math.Min(i + 16, bytes.Length); j++)
                    diag.Write($"{bytes[j]:X2} ");
                diag.WriteLine();
            }
            diag.WriteLine("└");
            diag.Flush();
        }
    }
}
