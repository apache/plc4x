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
    /// <para>
    /// The transport SPI is synchronous: opening a connection returns a connected
    /// instance, so the constructor has to block its caller until the connect is done.
    /// It does that in exactly one place, <see cref="Connect"/>. The connect itself is
    /// async I/O with a cancellation token that fires when the connect timeout expires,
    /// so an attempt that runs out of time is actually aborted rather than abandoned.
    /// The work runs on the thread pool, so no synchronization context on the calling
    /// thread can ever be needed to complete it.
    /// </para>
    /// <para>
    /// <see cref="Close"/> never waits for the read loop. Closing the socket is what
    /// ends the loop, and it finishes on its own thread.
    /// </para>
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

        // Set by Close() before it touches the socket: a write that fails while this
        // is 1 is failing during our own shutdown, so the write path treats it as a
        // normal close rather than an error, whether the failure came from the
        // shutdown itself or from a peer that happened to go away at the same time.
        private int _closing;

        private volatile Action? _dataListener;
        private volatile Action<Exception?>? _disconnectListener;

        public TcpTransportInstance(IPEndPoint remoteAddress, TcpTransportConfiguration configuration)
            : base(configuration)
        {
            if (remoteAddress == null)
            {
                throw new ArgumentNullException(nameof(remoteAddress));
            }

            _ringBuffer = new RingBuffer(configuration.ReceiveBufferSize);
            _socket = Connect(remoteAddress, configuration);
            RemoteAddress = remoteAddress;
            LocalAddress = _socket.LocalEndPoint as IPEndPoint;

            // Started last so a throw during setup cannot leak a running loop.
            _readLoop = Task.Run(RunReadLoopAsync);
        }

        public IPEndPoint RemoteAddress { get; }

        public IPEndPoint? LocalAddress { get; }

        /// <summary>
        /// Opens the connection and waits for it, but no longer than the configured
        /// connect timeout.
        /// </summary>
        private static Socket Connect(IPEndPoint remoteAddress, TcpTransportConfiguration configuration)
        {
            CancellationTokenSource? deadline = null;
            try
            {
                deadline = new CancellationTokenSource(configuration.ConnectTimeoutMillis);
                var token = deadline.Token;
                return Task.Run(() => ConnectAsync(remoteAddress, configuration, token), token)
                    .GetAwaiter().GetResult();
            }
            catch (OperationCanceledException) when (deadline?.IsCancellationRequested == true)
            {
                throw new TransportException(
                    $"Failed to connect to {remoteAddress.Address}:{remoteAddress.Port} - " +
                    $"timed out after {configuration.ConnectTimeout} ms.");
            }
            catch (Exception e) when (!(e is TransportException))
            {
                throw new TransportException(
                    $"Failed to connect to {remoteAddress.Address}:{remoteAddress.Port} - {e.Message}", e);
            }
            finally
            {
                deadline?.Dispose();
            }
        }

        /// <summary>
        /// Creates a socket, applies the configured options and connects it. Cancelling
        /// <paramref name="cancellationToken"/> aborts the connect and disposes the socket.
        /// </summary>
        internal static async Task<Socket> ConnectAsync(
            IPEndPoint remoteAddress, TcpTransportConfiguration configuration, CancellationToken cancellationToken)
        {
            var socket = new Socket(remoteAddress.AddressFamily, SocketType.Stream, ProtocolType.Tcp);
            try
            {
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

                await socket.ConnectAsync(remoteAddress, cancellationToken).ConfigureAwait(false);
                return socket;
            }
            catch
            {
                socket.Dispose();
                throw;
            }
        }

        public override bool IsOpen => Volatile.Read(ref _open) == 1 && _socket.Connected;

        /// <summary>
        /// Whether the socket has been shut down and disposed. Exposed so tests can tell
        /// "the peer went away" apart from "our end was actually released" - the two used
        /// to be the same flag, which is how the release came to be skipped.
        /// </summary>
        internal bool ResourcesReleased => Volatile.Read(ref _released) == 1;

        /// <summary>
        /// Completes when the read loop has finished. Exposed so tests can wait for the
        /// end of a shutdown that <see cref="Close"/> deliberately does not wait for.
        /// </summary>
        internal Task ReadLoopCompleted => _readLoop;

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
                catch (SocketException) when (Volatile.Read(ref _closing) == 1)
                {
                    // Close() shuts the socket down before disposing it, so a write
                    // blocked on a full send buffer wakes up with SocketException
                    // (an abort/interrupt style error, or a broken pipe) rather than
                    // ObjectDisposedException. _closing is set before the shutdown
                    // starts, so seeing it 1 means our own close caused this: normal
                    // shutdown, same as the disposed case above. A peer failure that
                    // races a concurrent close gets reported the same way; a peer
                    // reset without any local close still surfaces below.
                }
                catch (SocketException e)
                {
                    throw new TransportException("Failed to write data", e);
                }
            }
        }

        public override void Close()
        {
            Interlocked.Exchange(ref _closing, 1);
            Volatile.Write(ref _open, 0);
            ReleaseResources();
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
        /// <para>
        /// Takes no locks and does not wait for the read loop. Closing the socket is what
        /// unblocks the loop, and taking <c>_writeLock</c> first would deadlock against a
        /// blocked writer. The loop ends on its own thread; a listener callback that is
        /// still running when this returns simply finishes first.
        /// </para>
        /// </remarks>
        private void ReleaseResources()
        {
            if (Interlocked.CompareExchange(ref _released, 1, 0) != 0)
            {
                return;
            }

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
                // The read loop may still be reading the token. It treats the resulting
                // ObjectDisposedException as the normal end of a shutdown.
                _shutdown.Dispose();
            }
        }

        public void RegisterDataListener(Action listener) => _dataListener = listener;

        public void RemoveDataListener() => _dataListener = null;

        public void RegisterDisconnectListener(Action<Exception?> listener) => _disconnectListener = listener;

        public void RemoveDisconnectListener() => _disconnectListener = null;

        private void NotifyDisconnect(Exception? cause)
        {
            var listener = _disconnectListener;
            if (listener == null)
            {
                return;
            }
            try
            {
                listener(cause);
            }
            catch
            {
                // A misbehaving listener must not take down the read loop.
            }
        }

        private void NotifyData()
        {
            var listener = _dataListener;
            if (listener == null)
            {
                return;
            }
            try
            {
                listener();
            }
            catch
            {
                // Same reasoning as NotifyDisconnect.
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
                        .ReceiveAsync(buffer.AsMemory(0, toRead), SocketFlags.None, _shutdown.Token)
                        .ConfigureAwait(false);

                    if (bytesRead == 0)
                    {
                        // Orderly shutdown by the remote end. When our own Close() got
                        // here first, the zero is just its Shutdown coming back and
                        // there is nothing to report.
                        if (Interlocked.Exchange(ref _open, 0) == 1)
                        {
                            NotifyDisconnect(null);
                        }
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
                // Close() disposed the socket or the token while we were using them.
            }
            catch (Exception e)
            {
                // Only a failure while still open counts as a disconnect.
                if (Interlocked.Exchange(ref _open, 0) == 1)
                {
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
                ReleaseResources();
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