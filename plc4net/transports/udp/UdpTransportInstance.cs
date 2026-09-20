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
using org.apache.plc4net.spi.generation;
using org.apache.plc4net.spi.transports;

namespace org.apache.plc4net.transports.udp
{
    /// <summary>
    /// UDP transport: a bound socket, connected to one remote endpoint, with a
    /// background receive loop feeding a ring buffer.
    /// </summary>
    /// <remarks>
    /// Each <see cref="Write"/> sends one datagram; each datagram received is appended
    /// to the buffer whole, in arrival order. There is no stream, so a codec that
    /// wants message boundaries relies on the protocol's own length field
    /// (KNXnet/IP's total-length header) — the same as over TCP.
    /// </remarks>
    public class UdpTransportInstance : BaseTransportInstance, IAsyncTransportInstance
    {
        private readonly Socket _socket;
        private readonly RingBuffer _ringBuffer;
        private readonly object _readLock = new object();
        private readonly object _writeLock = new object();
        private readonly CancellationTokenSource _shutdown = new CancellationTokenSource();
        private readonly Task _receiveLoop;

        private const int MaxConsecutiveTransientErrors = 20;
        private const int TransientBackoffMs = 50;

        private int _open = 1;
        private int _released;

        private volatile Action _dataListener;
        private volatile Action<Exception> _disconnectListener;
        private int _listenerThreadId;

        public UdpTransportInstance(IPEndPoint remoteAddress, UdpTransportConfiguration configuration)
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
                socket = new Socket(remoteAddress.AddressFamily, SocketType.Dgram, ProtocolType.Udp);
                if (configuration.Broadcast)
                {
                    socket.EnableBroadcast = true;
                }

                var bindAddress = string.IsNullOrEmpty(configuration.LocalAddress)
                    ? (remoteAddress.AddressFamily == AddressFamily.InterNetworkV6
                        ? IPAddress.IPv6Any
                        : IPAddress.Any)
                    : IPAddress.Parse(configuration.LocalAddress);
                socket.Bind(new IPEndPoint(bindAddress, configuration.LocalPort));

                // "Connecting" a datagram socket fixes the destination for Send and
                // filters inbound datagrams to that peer - exactly the KNXnet/IP model
                // (one gateway per connection).
                socket.Connect(remoteAddress);

                _socket = socket;
                RemoteAddress = remoteAddress;
                LocalAddress = socket.LocalEndPoint as IPEndPoint;
            }
            catch (Exception e)
            {
                socket?.Dispose();
                throw new TransportException(
                    $"Failed to open a UDP socket to {remoteAddress.Address}:{remoteAddress.Port} - {e.Message}", e);
            }

            _receiveLoop = Task.Run(RunReceiveLoopAsync);
        }

        public IPEndPoint RemoteAddress { get; }

        /// <summary>The local endpoint the socket bound to. KNXnet/IP advertises this
        /// to the gateway in its HPAI fields, so the gateway knows where to reply.</summary>
        public IPEndPoint LocalAddress { get; }

        public override bool IsOpen => Volatile.Read(ref _open) == 1;

        internal bool ResourcesReleased => Volatile.Read(ref _released) == 1;

        /// <summary>Optional diagnostic output: written and read bytes as annotated hex.</summary>
        public TextWriter DiagnosticOutput { get; set; }

        public override int GetNumBytesAvailable()
        {
            lock (_readLock)
            {
                return _ringBuffer.AvailableForReading;
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
            WriteHex("UDP RECV", result);
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
                    // One datagram per call. A short Send on a datagram socket means the
                    // datagram did not fit; treat it as an error rather than fragmenting.
                    var sent = _socket.Send(bytes, 0, bytes.Length, SocketFlags.None);
                    if (sent != bytes.Length)
                    {
                        throw new TransportException(
                            $"UDP datagram truncated: sent {sent} of {bytes.Length} bytes.");
                    }
                    WriteHex("UDP SEND", bytes);
                }
                catch (ObjectDisposedException) when (Volatile.Read(ref _open) == 0)
                {
                    // A concurrent Close() disposed the socket mid-write: normal shutdown.
                }
                catch (SocketException e)
                {
                    throw new TransportException("Failed to send datagram", e);
                }
            }
        }

        public override void Close()
        {
            Volatile.Write(ref _open, 0);
            ReleaseResources(fromReceiveLoop: false);
        }

        private void ReleaseResources(bool fromReceiveLoop)
        {
            if (Interlocked.CompareExchange(ref _released, 1, 0) != 0)
            {
                return;
            }

            _shutdown.Cancel();
            try
            {
                _socket.Dispose();
            }
            catch (ObjectDisposedException)
            {
            }
            finally
            {
                var insideLoop = fromReceiveLoop
                    || Volatile.Read(ref _listenerThreadId) == Environment.CurrentManagedThreadId;
                if (!insideLoop && _receiveLoop != null && !_receiveLoop.IsCompleted)
                {
                    _receiveLoop.Wait(TimeSpan.FromSeconds(1));
                }
                _shutdown.Dispose();
            }
        }

        public void RegisterDataListener(Action listener)
        {
            _dataListener = listener;

            // A datagram may have arrived and been buffered before the listener was
            // registered; deliver one wake-up so the codec drains what is waiting.
            bool hasBufferedData;
            lock (_readLock)
            {
                hasBufferedData = _ringBuffer.AvailableForReading > 0;
            }
            if (listener != null && hasBufferedData)
            {
                NotifyData();
            }
        }

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
                // A misbehaving listener must not take down the receive loop.
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

        private async Task RunReceiveLoopAsync()
        {
            var datagram = new byte[65535];
            var transientErrors = 0;
            try
            {
                while (Volatile.Read(ref _open) == 1)
                {
                    int received;
                    try
                    {
                        received = await _socket
                            .ReceiveAsync(new ArraySegment<byte>(datagram), SocketFlags.None)
                            .ConfigureAwait(false);
                        transientErrors = 0;
                    }
                    catch (SocketException e) when (IsTransientReceiveError(e.SocketErrorCode))
                    {
                        // ICMP port/host-unreachable from the peer (ConnectionReset on
                        // Windows, ConnectionRefused on Linux), or an oversized datagram
                        // (MessageSize). The gateway may just not be up yet; back off and
                        // keep listening rather than tearing the transport down. Escalate
                        // only if the errors do not stop.
                        if (++transientErrors > MaxConsecutiveTransientErrors)
                        {
                            throw;
                        }
                        await Task.Delay(TransientBackoffMs, _shutdown.Token).ConfigureAwait(false);
                        continue;
                    }

                    if (received <= 0)
                    {
                        continue;
                    }

                    var accepted = false;
                    lock (_readLock)
                    {
                        if (received <= _ringBuffer.Capacity)
                        {
                            if (_ringBuffer.RemainingForWriting < received)
                            {
                                // The codec has not kept up. Discard everything buffered
                                // rather than drop bytes from the middle of a frame - a
                                // partial frame left behind desyncs the length-prefixed
                                // codec permanently. The next datagram starts clean.
                                _ringBuffer.Clear();
                            }
                            _ringBuffer.Write(datagram, 0, received);
                            accepted = true;
                        }
                    }

                    if (accepted)
                    {
                        NotifyData();
                    }
                    // else: a datagram larger than the whole buffer can never be
                    // delivered intact; drop it and stay alive.
                }
            }
            catch (OperationCanceledException)
            {
                // Close() cancelled us.
            }
            catch (ObjectDisposedException)
            {
                // Close() disposed the socket while we were awaiting a receive.
            }
            catch (Exception e)
            {
                if (Volatile.Read(ref _open) == 1)
                {
                    Volatile.Write(ref _open, 0);
                    NotifyDisconnect(e);
                }
            }
            finally
            {
                Volatile.Write(ref _open, 0);
                ReleaseResources(fromReceiveLoop: true);
            }
        }

        private static bool IsTransientReceiveError(SocketError code) =>
            code == SocketError.ConnectionReset       // Windows: ICMP port-unreachable
            || code == SocketError.ConnectionRefused  // Linux: ICMP port-unreachable
            || code == SocketError.MessageSize        // datagram larger than the buffer
            || code == SocketError.HostUnreachable
            || code == SocketError.NetworkUnreachable;

        private void WriteHex(string label, byte[] bytes)
        {
            var diag = DiagnosticOutput;
            if (diag == null)
            {
                return;
            }
            diag.WriteLine($"┌ {label} ({bytes.Length} bytes)");
            for (var i = 0; i < bytes.Length; i += 16)
            {
                diag.Write("│ ");
                for (var j = i; j < Math.Min(i + 16, bytes.Length); j++)
                {
                    diag.Write($"{bytes[j]:X2} ");
                }
                diag.WriteLine();
            }
            diag.WriteLine("└");
            diag.Flush();
        }
    }
}
