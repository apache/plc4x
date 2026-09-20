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
        private readonly Task _readLoopTask;

        private Action? _dataListener;
        private Action<Exception>? _disconnectListener;
        private volatile bool _closed;

        // Set to the current thread's id for the duration of a synchronous listener
        // callback invoked from the read loop, so Close() can tell "a listener called
        // Close() from inside the loop" (must not wait on itself) apart from every
        // other caller (should wait for the loop to actually unwind).
        private volatile int _listenerThreadId;

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
            // lifecycle is bound to the serial port's. If Open() fails
            // (port missing / busy / denied) dispose what was built so
            // far rather than leak the CTS and the half-open port.
            try
            {
                _port.Open();
            }
            catch
            {
                _readCts.Dispose();
                _port.Dispose();
                throw;
            }
            _readLoopTask = Task.Run(() => ReadLoopAsync(_readCts.Token));
        }

        // ── async read loop ─────────────────────────────────────

        private async Task ReadLoopAsync(CancellationToken ct)
        {
            var buf = new byte[4096];
            var idleDelayMs = 1;
            while (!ct.IsCancellationRequested)
            {
                try
                {
                    // Backpressure: if the codec has not drained the ring buffer
                    // there is no room. Wait and retry rather than truncate the
                    // read — dropping bytes here misaligns every following frame.
                    int free;
                    lock (_readBuffer) { free = _readBuffer.RemainingForWriting; }
                    if (free == 0)
                    {
                        await Task.Delay(1, ct).ConfigureAwait(false);
                        continue;
                    }

                    var toRead = Math.Min(buf.Length, free);
                    var available = _port.BytesToRead;
                    if (available == 0)
                    {
                        // The natural blocking wait here would be
                        // _port.BaseStream.ReadAsync(ct), but it does not observe
                        // CancellationToken on a real SerialPort (confirmed on real
                        // hardware — see the commit that introduced this poll), which
                        // left Close() unable to ever stop this loop. Poll BytesToRead
                        // instead, backing off while the line is idle so an
                        // hours-long-idle connection is not spinning at full poll rate;
                        // reset to a tight poll the moment bytes actually show up.
                        // Reverting to the ReadAsync(ct) form would silently
                        // reintroduce that hardware-verified bug.
                        await Task.Delay(idleDelayMs, ct).ConfigureAwait(false);
                        idleDelayMs = Math.Min(idleDelayMs * 2, 20);
                        continue;
                    }
                    idleDelayMs = 1;

                    var bytesRead = _port.Read(buf, 0, Math.Min(toRead, available));
                    if (bytesRead == 0)
                    {
                        // Not observed in practice: the guards above only ever request
                        // a positive count already confirmed available via BytesToRead,
                        // and SerialPort.Read throws TimeoutException rather than
                        // returning 0 for a positive count. Retry rather than end the
                        // loop in case some platform ever does return 0 here.
                        continue;
                    }

                    lock (_readBuffer)
                    {
                        _readBuffer.Write(buf, 0, bytesRead);
                    }

                    _listenerThreadId = Environment.CurrentManagedThreadId;
                    try
                    {
                        _dataListener?.Invoke();
                    }
                    catch (Exception listenerEx)
                    {
                        // A throwing listener must not kill the read loop.
                        _disconnectListener?.Invoke(listenerEx);
                    }
                    finally
                    {
                        _listenerThreadId = 0;
                    }
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
                    // Close() already set _closed before tearing down the port, so a
                    // fault that only happened because Close() disposed out from under
                    // us is expected shutdown noise, not a real disconnect to report.
                    if (!_closed)
                    {
                        _listenerThreadId = Environment.CurrentManagedThreadId;
                        try
                        {
                            _disconnectListener?.Invoke(ex);
                        }
                        finally
                        {
                            _listenerThreadId = 0;
                        }
                    }
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

            // Do not wait on the read loop when we are standing inside a listener
            // callback it invoked synchronously — the loop cannot finish while that
            // callback (running on the loop's own thread) is still on the stack, so
            // the wait would just burn its timeout for nothing. Otherwise, give the
            // loop a bounded window to actually unwind before returning: without this,
            // a caller could treat the port as fully released while the background
            // task is still mid-iteration on it — a known SerialPort Close()-vs-
            // concurrent-Read() hazard (dotnet/runtime#20362, dotnet/corefx#36040).
            if (Environment.CurrentManagedThreadId != _listenerThreadId)
            {
                try { _readLoopTask.Wait(TimeSpan.FromSeconds(1)); } catch { /* best-effort */ }
            }
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
