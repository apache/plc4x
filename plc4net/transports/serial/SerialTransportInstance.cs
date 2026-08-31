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
using System.IO.Ports;
using System.Threading;
using System.Threading.Tasks;
using org.apache.plc4net.spi.transports;

namespace org.apache.plc4net.transports.serial
{
    /// <summary>
    /// A live serial connection.  Reads land in a ring buffer served by an
    /// async read loop; writes go directly to the port.
    /// </summary>
    public class SerialTransportInstance : BaseTransportInstance, IAsyncTransportInstance
    {
        private readonly SerialPort _port;
        private readonly RingBuffer _readBuffer;
        private readonly CancellationTokenSource _readCts = new CancellationTokenSource();

        private Action? _dataListener;
        private Action<Exception>? _disconnectListener;
        private volatile bool _closed;

        public SerialTransportInstance(
            string portName,
            SerialTransportConfiguration config)
            : base(config)
        {
            if (string.IsNullOrWhiteSpace(portName))
                throw new ArgumentException("Port name is required.", nameof(portName));

            _port = new SerialPort(portName)
            {
                BaudRate = config.BaudRate,
                DataBits = config.DataBits,
                Parity = config.Parity,
                StopBits = config.StopBits,
                Handshake = config.Handshake,
                ReadTimeout = config.ReadTimeout,
                WriteTimeout = config.WriteTimeout
            };
            _readBuffer = new RingBuffer(config.ReceiveBufferSize);

            // Open the port and start the read loop immediately —
            // the transport instance IS the connection, so its
            // lifecycle is bound to the serial port's.
            _port.Open();
            _ = Task.Run(() => ReadLoopAsync(_readCts.Token));
        }

        // ── async read loop ─────────────────────────────────────

        private async Task ReadLoopAsync(CancellationToken ct)
        {
            var buf = new byte[4096];
            while (!ct.IsCancellationRequested)
            {
                try
                {
                    // SerialPort.BaseStream.ReadAsync does not respect
                    // CancellationToken natively; the token cancels the
                    // polling task.
                    var bytesRead = await _port.BaseStream.ReadAsync(
                        buf, 0, buf.Length, ct).ConfigureAwait(false);
                    if (bytesRead == 0) break;

                    lock (_readBuffer)
                    {
                        var toWrite = Math.Min(bytesRead,
                            _readBuffer.RemainingForWriting);
                        if (toWrite > 0)
                            _readBuffer.Write(buf, 0, toWrite);
                    }
                    _dataListener?.Invoke();
                }
                catch (OperationCanceledException)
                {
                    break;
                }
                catch (ObjectDisposedException)
                {
                    break;
                }
                catch (InvalidOperationException)
                {
                    break;
                }
                catch (Exception ex)
                {
                    _disconnectListener?.Invoke(ex);
                    break;
                }
            }
        }

        // ── ITransportInstance ───────────────────────────────────

        public override bool IsOpen => !_closed && _port.IsOpen;

        public override int GetNumBytesAvailable()
        {
            lock (_readBuffer) { return _readBuffer.AvailableForReading; }
        }

        public override byte[] PeekReadableBytes(int numBytes)
        {
            if (numBytes <= 0) return Array.Empty<byte>();
            lock (_readBuffer)
            {
                return _readBuffer.Peek(numBytes);
            }
        }

        public override byte[] Read(int numBytes)
        {
            if (numBytes <= 0) return Array.Empty<byte>();
            lock (_readBuffer)
            {
                return _readBuffer.Read(numBytes);
            }
        }

        public override void Write(byte[] bytes)
        {
            if (bytes == null || bytes.Length == 0) return;
            if (_closed) throw new TransportException("Serial port is closed.");
            _port.Write(bytes, 0, bytes.Length);
        }

        public override void Close()
        {
            if (_closed) return;
            _closed = true;
            _readCts.Cancel();
            try { _port.Close(); } catch { /* best-effort */ }
            _readCts.Dispose();
            try { _port.Dispose(); } catch { }
        }

        // ── IAsyncTransportInstance ──────────────────────────────

        public void RegisterDataListener(Action listener)
        {
            _dataListener = listener;
        }

        public void RemoveDataListener()
        {
            _dataListener = null;
        }

        public void RegisterDisconnectListener(Action<Exception> listener)
        {
            _disconnectListener = listener;
        }

        public void RemoveDisconnectListener()
        {
            _disconnectListener = null;
        }
    }
}
