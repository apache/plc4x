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
using System.Threading;
using org.apache.plc4net.spi.transports;

namespace org.apache.plc4net.transports.test
{
    /// <summary>
    /// In-memory transport instance for testing. Data flows in two directions:
    /// the test harness <see cref="InjectTestData"/> into the read buffer (simulating
    /// bytes arriving from the wire), and the driver writes into the write buffer
    /// (from which the harness reads with <see cref="GetWrittenData"/>).
    /// </summary>
    /// <remarks>
    /// Mirrors the Java SPI3 <c>TestTransportInstance</c>. Thread-safe: the read
    /// and write paths each have their own lock, and tests can inject data from a
    /// background thread while the driver reads on its own.
    /// </remarks>
    public class TestTransportInstance : BaseTransportInstance
    {
        private readonly RingBuffer _readBuffer;
        private readonly RingBuffer _writeBuffer;
        private readonly object _readLock = new object();
        private readonly object _writeLock = new object();
        private volatile bool _open = true;

        /// <summary>
        /// Creates an instance backed by the given configuration.
        /// </summary>
        public TestTransportInstance(TestTransportConfiguration configuration)
            : base(configuration)
        {
            _readBuffer = new RingBuffer(configuration.ReceiveBufferSize);
            _writeBuffer = new RingBuffer(configuration.ReceiveBufferSize);
        }

        // ── test-facing API ──────────────────────────────────────────

        /// <summary>
        /// Feeds bytes into the read buffer, simulating data arriving from the
        /// remote endpoint. The transport's data listener (if any) is notified.
        /// </summary>
        /// <returns>The number of bytes actually written.</returns>
        public int InjectTestData(byte[] data)
        {
            if (data == null || data.Length == 0)
            {
                return 0;
            }

            int written;
            lock (_readLock)
            {
                written = Math.Min(data.Length, _readBuffer.RemainingForWriting);
                _readBuffer.Write(data, 0, written);
            }

            return written;
        }

        /// <summary>
        /// Returns up to <paramref name="numBytes"/> bytes that the driver has
        /// written to the transport. Consumes them.
        /// </summary>
        public byte[] GetWrittenData(int numBytes)
        {
            lock (_writeLock)
            {
                var available = _writeBuffer.AvailableForReading;
                var toRead = Math.Min(numBytes, available);
                if (toRead <= 0)
                {
                    return Array.Empty<byte>();
                }
                return _writeBuffer.Read(toRead);
            }
        }

        /// <summary>Returns and consumes all bytes the driver has written.</summary>
        public byte[] GetAllWrittenData()
        {
            lock (_writeLock)
            {
                var available = _writeBuffer.AvailableForReading;
                return available <= 0
                    ? Array.Empty<byte>()
                    : _writeBuffer.Read(available);
            }
        }

        /// <summary>
        /// Returns the count of unread bytes the driver has written.
        /// </summary>
        public int GetNumBytesWritten()
        {
            lock (_writeLock)
            {
                return _writeBuffer.AvailableForReading;
            }
        }

        /// <summary>
        /// Waits for at least <paramref name="expectedBytes"/> to be written,
        /// then returns them. Times out with a <see cref="TransportException"/>
        /// if the deadline passes without enough data.
        /// </summary>
        public byte[] WaitForWrittenData(int expectedBytes, int timeoutMs)
        {
            var deadline = Environment.TickCount64 + timeoutMs;
            lock (_writeLock)
            {
                while (_writeBuffer.AvailableForReading < expectedBytes)
                {
                    var remaining = deadline - Environment.TickCount64;
                    if (remaining <= 0)
                    {
                        throw new TransportException(
                            $"Timeout waiting for {expectedBytes} bytes; only " +
                            $"{_writeBuffer.AvailableForReading} arrived after {timeoutMs} ms.");
                    }

                    // Woken by the Monitor.PulseAll in Write/Close, so a write is
                    // observed the moment it lands. The old form polled every 1 ms
                    // on a thread-pool thread; under parallel-test load on a CI
                    // runner the writer's own continuation could be starved past
                    // the deadline, which is what made the caller flaky.
                    Monitor.Wait(_writeLock, (int)remaining);
                }

                return _writeBuffer.Read(expectedBytes);
            }
        }

        // ── ITransportInstance ────────────────────────────────────────

        public override bool IsOpen => _open;

        public override int GetNumBytesAvailable()
        {
            lock (_readLock)
            {
                return _open ? _readBuffer.AvailableForReading : 0;
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
                EnsureOpen();
                if (_readBuffer.AvailableForReading < numBytes)
                {
                    throw new TransportException(
                        $"Requested {numBytes} bytes but only {_readBuffer.AvailableForReading} available.");
                }
                return _readBuffer.Peek(numBytes);
            }
        }

        public override byte[] Read(int numBytes)
        {
            if (numBytes <= 0)
            {
                return Array.Empty<byte>();
            }

            lock (_readLock)
            {
                EnsureOpen();
                if (_readBuffer.AvailableForReading < numBytes)
                {
                    throw new TransportException(
                        $"Requested {numBytes} bytes but only {_readBuffer.AvailableForReading} available.");
                }
                return _readBuffer.Read(numBytes);
            }
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
                if (bytes.Length > _writeBuffer.RemainingForWriting)
                {
                    throw new TransportException(
                        $"Cannot write {bytes.Length} bytes; only {_writeBuffer.RemainingForWriting} free.");
                }
                _writeBuffer.Write(bytes);
                Monitor.PulseAll(_writeLock);
            }
        }

        public override void Close()
        {
            if (!_open)
                return;

            // Only hold _writeLock; _readLock callers check _open after acquiring
            // their lock, so setting _open = false under _writeLock is sufficient.
            lock (_writeLock)
            {
                _open = false;
                Monitor.PulseAll(_writeLock);
            }
        }

        }
}
