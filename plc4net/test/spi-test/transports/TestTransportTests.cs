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

using System.Threading;
using org.apache.plc4net.spi.transports;
using org.apache.plc4net.transports.test;
using Xunit;

namespace org.apache.plc4net.spi.test.transports
{
    /// <summary>
    /// Tests for the in-memory test transport, which drives the read/write path
    /// without a real network.
    /// </summary>
    public class TestTransportTests
    {
        private static TestTransportInstance CreateInstance(int bufferSize = 1024)
        {
            var config = new TestTransportConfiguration { ReceiveBufferSize = bufferSize };
            return (TestTransportInstance)new TestTransport().CreateTransportInstance("", config);
        }

        [Fact]
        public void A_new_instance_is_open()
        {
            using var instance = CreateInstance();
            Assert.True(instance.IsOpen);
        }

        [Fact]
        public void Close_marks_the_instance_closed()
        {
            using var instance = CreateInstance();
            instance.Close();
            Assert.False(instance.IsOpen);
        }

        [Fact]
        public void Close_is_idempotent()
        {
            using var instance = CreateInstance();
            instance.Close();
            instance.Close();  // must not throw
        }

        [Fact]
        public void Injected_data_can_be_read()
        {
            using var instance = CreateInstance();
            var data = new byte[] { 0x01, 0x02, 0x03 };

            instance.InjectTestData(data);

            Assert.Equal(3, instance.GetNumBytesAvailable());
            Assert.Equal(data, instance.PeekReadableBytes(3));
            Assert.Equal(data, instance.Read(3));
            Assert.Equal(0, instance.GetNumBytesAvailable());
        }

        [Fact]
        public void Peek_does_not_consume()
        {
            using var instance = CreateInstance();
            instance.InjectTestData(new byte[] { 0xAB, 0xCD });

            Assert.Equal(new byte[] { 0xAB, 0xCD }, instance.PeekReadableBytes(2));
            Assert.Equal(2, instance.GetNumBytesAvailable());  // still there
        }

        [Fact]
        public void Write_is_captured_for_the_harness()
        {
            using var instance = CreateInstance();
            instance.Write(new byte[] { 0x10, 0x20 });

            Assert.Equal(2, instance.GetNumBytesWritten());
            Assert.Equal(new byte[] { 0x10, 0x20 }, instance.GetAllWrittenData());
            Assert.Equal(0, instance.GetNumBytesWritten());
        }

        [Fact]
        public void Reading_past_what_is_available_throws()
        {
            using var instance = CreateInstance();
            instance.InjectTestData(new byte[] { 0xFF });

            Assert.Throws<TransportException>(() => instance.Read(10));
        }

        [Fact]
        public void WaitForWrittenData_returns_when_enough_bytes_arrive()
        {
            using var instance = CreateInstance();

            // A dedicated thread, not Task.Run: a fire-and-forget task whose
            // continuation resumes after Task.Delay depends on a free thread-pool
            // thread, and under parallel-test load on a CI runner that resumption
            // was starved past the deadline — the write never landed and the wait
            // timed out on zero bytes. A real thread always runs.
            var writer = new Thread(() =>
            {
                Thread.Sleep(50);
                instance.Write(new byte[] { 0x50, 0x60, 0x70 });
            })
            { IsBackground = true, Name = "test-transport-writer" };
            writer.Start();

            // WaitForWrittenData is woken by the write, so the real wait is ~50 ms;
            // the timeout is only a ceiling for a slow runner, never the happy path.
            var data = instance.WaitForWrittenData(3, 10_000);

            writer.Join();
            Assert.Equal(new byte[] { 0x50, 0x60, 0x70 }, data);
        }

        [Fact]
        public void WaitForWrittenData_throws_when_the_deadline_passes()
        {
            using var instance = CreateInstance();
            instance.Write(new byte[] { 0x01 });   // one byte; the call asks for three

            var ex = Assert.Throws<TransportException>(
                () => instance.WaitForWrittenData(3, 50));
            Assert.Contains("only 1", ex.Message);
        }
    }
}
