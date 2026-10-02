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
using System.Collections.Concurrent;
using System.Collections.Generic;
using System.Diagnostics;
using System.IO;
using System.Linq;
using System.Text;
using System.Threading;
using System.Threading.Tasks;
using org.apache.plc4net.spi.transports;
using org.apache.plc4net.transports.cotp;
using Xunit;

namespace org.apache.plc4net.spi.test.transports
{
    /// <summary>
    /// An inner transport that serves pre-loaded bytes and records everything written,
    /// so COTP framing can be asserted byte for byte without a socket.
    /// </summary>
    internal class ScriptedTransportConfiguration : ITransportConfiguration
    {
    }

    internal class ScriptedTransportInstance : BaseTransportInstance
    {
        private readonly List<byte> _inbound = new List<byte>();
        private readonly object _sync = new object();

        public ScriptedTransportInstance() : base(new ScriptedTransportConfiguration()) { }

        public List<byte[]> Written { get; } = new List<byte[]>();

        public bool Closed { get { lock (_sync) { return _closed; } } }
        private bool _closed;
        public int WrittenCount { get { lock (_sync) { return Written.Count; } } }

        public override bool IsOpen => !Closed;

        public void Inject(params byte[] bytes)
        {
            lock (_sync) { _inbound.AddRange(bytes); }
        }

        public override int GetNumBytesAvailable()
        {
            lock (_sync) { return _inbound.Count; }
        }

        public override byte[] PeekReadableBytes(int numBytes)
        {
            lock (_sync)
            {
                if (numBytes > _inbound.Count)
                {
                    throw new TransportException(
                        $"Requested {numBytes} bytes but only {_inbound.Count} available.");
                }
                return _inbound.Take(numBytes).ToArray();
            }
        }

        public override byte[] Read(int numBytes)
        {
            lock (_sync)
            {
                if (numBytes > _inbound.Count)
                {
                    throw new TransportException(
                        $"Requested {numBytes} bytes but only {_inbound.Count} available.");
                }
                var result = _inbound.Take(numBytes).ToArray();
                _inbound.RemoveRange(0, numBytes);
                return result;
            }
        }

        public override void Write(byte[] bytes)
        {
            lock (_sync) { Written.Add(bytes); }
        }

        public override void Close()
        {
            lock (_sync) { _closed = true; }
        }
    }

    /// <summary>
    /// An inner transport whose writes can be made to hang like a socket send against a peer
    /// that has stopped reading. A stalled write only returns once the transport is closed
    /// underneath it, which is how closing a real socket releases a blocked send.
    /// </summary>
    internal class StalledWriteTransportInstance : ScriptedTransportInstance
    {
        private readonly ManualResetEventSlim _writeStarted = new ManualResetEventSlim(false);
        private readonly ManualResetEventSlim _released = new ManualResetEventSlim(false);
        private volatile bool _stallWrites;

        public void StallWrites() => _stallWrites = true;

        public bool WaitForStalledWrite(TimeSpan timeout) => _writeStarted.Wait(timeout);

        public override void Write(byte[] bytes)
        {
            if (!_stallWrites)
            {
                base.Write(bytes);
                return;
            }

            _writeStarted.Set();
            _released.Wait();
            throw new TransportException("Transport closed while a write was blocked.");
        }

        public override void Close()
        {
            base.Close();
            _released.Set();
        }
    }

    /// <summary>
    /// An inner transport whose writes take a moment, so that two writers racing for the same
    /// COTP instance reliably overlap instead of one finishing before the other has started.
    /// </summary>
    internal class SlowWriteTransportInstance : ScriptedTransportInstance
    {
        private volatile bool _slow;

        public void SlowWrites() => _slow = true;

        public override void Write(byte[] bytes)
        {
            if (_slow) Thread.Sleep(5);
            base.Write(bytes);
        }
    }

    /// <summary>
    /// A diagnostic sink that runs a callback the first time a line containing a marker is
    /// written. It lets a test act at an exact point inside Open() without depending on
    /// thread timing.
    /// </summary>
    internal sealed class CallbackTextWriter : TextWriter
    {
        private readonly string _marker;
        private readonly Action _onMarker;
        private bool _fired;

        public CallbackTextWriter(string marker, Action onMarker)
        {
            _marker = marker;
            _onMarker = onMarker;
        }

        public override Encoding Encoding => Encoding.UTF8;

        public override void WriteLine(string? value)
        {
            if (_fired || value == null || !value.Contains(_marker)) return;
            _fired = true;
            _onMarker();
        }
    }

    public class CotpTransportInstanceTests
    {
        // A well-formed Connection Confirm: TPKT(4) + LI + 0xD0 + DST-REF(2) + SRC-REF(2)
        // + class(1). Captured shape from an S7-1500 exchange.
        internal static byte[] ConnectionConfirm() => new byte[]
        {
            0x03, 0x00, 0x00, 0x0B,             // TPKT: version, reserved, length = 11
            0x06,                               // LI = 6
            0xD0,                               // CC
            0x00, 0x01,                         // DST-REF
            0x00, 0x02,                         // SRC-REF
            0x00                                // class 0
        };

        private static (CotpTransportInstance cotp, ScriptedTransportInstance inner) Handshaken()
        {
            var inner = new ScriptedTransportInstance();
            inner.Inject(ConnectionConfirm());
            var cotp = new CotpTransportInstance(inner);
            cotp.Open(0x01, 0x00, 0x03, 0x01);
            return (cotp, inner);
        }

        // A call that hangs must fail its test rather than hang it, so awaiting one is bounded.
        private static async Task AssertCompletes(Task task, string failureMessage)
        {
            var finished = await Task.WhenAny(task, Task.Delay(TimeSpan.FromSeconds(2)));
            Assert.True(ReferenceEquals(finished, task), failureMessage);
            await task;
        }

        private static async Task<T> ResultWithin<T>(Task<T> task, string failureMessage)
        {
            await AssertCompletes(task, failureMessage);
            return await task;
        }

        [Fact]
        public void The_constructor_rejects_a_null_inner_transport()
        {
            var ex = Assert.Throws<ArgumentNullException>(() => new CotpTransportInstance(null!));

            Assert.Equal("inner", ex.ParamName);
        }

        [Fact]
        public void Open_sends_a_connection_request_carrying_both_tsaps()
        {
            var inner = new ScriptedTransportInstance();
            inner.Inject(ConnectionConfirm());
            var cotp = new CotpTransportInstance(inner);

            cotp.Open(localTsapHi: 0x01, localTsapLo: 0x00,
                      remoteTsapHi: 0x03, remoteTsapLo: 0x01);

            var cr = Assert.Single(inner.Written);

            // TPKT header, then the COTP CR.
            Assert.Equal(0x03, cr[0]);
            Assert.Equal(22, cr.Length);            // TPKT(4) + CR(18)
            Assert.Equal(22, (cr[2] << 8) | cr[3]); // TPKT length field covers the whole frame

            Assert.Equal(0x11, cr[4]);              // LI = 17
            Assert.Equal(0xE0, cr[5]);              // CR PDU type

            // Calling TSAP: parameter code 0xC1, length 2, then the local TSAP.
            Assert.Equal(0xC1, cr[14]);
            Assert.Equal(0x02, cr[15]);
            Assert.Equal(0x01, cr[16]);
            Assert.Equal(0x00, cr[17]);

            // Called TSAP: parameter code 0xC2, length 2, then the remote TSAP.
            Assert.Equal(0xC2, cr[18]);
            Assert.Equal(0x02, cr[19]);
            Assert.Equal(0x03, cr[20]);
            Assert.Equal(0x01, cr[21]);
        }

        [Fact]
        public void Open_is_idempotent()
        {
            var (cotp, inner) = Handshaken();

            cotp.Open(0x01, 0x00, 0x03, 0x01);

            // Still exactly one CR — the second call must not re-handshake.
            Assert.Single(inner.Written);
        }

        [Fact]
        public void Open_rejects_a_response_that_is_not_a_connection_confirm()
        {
            var inner = new ScriptedTransportInstance();
            // Same framing as a CC but the PDU type says Disconnect Request (0x80).
            inner.Inject(0x03, 0x00, 0x00, 0x0B, 0x06, 0x80, 0x00, 0x01, 0x00, 0x02, 0x00);
            var cotp = new CotpTransportInstance(inner);

            var ex = Assert.Throws<TransportException>(
                () => cotp.Open(0x01, 0x00, 0x03, 0x01));
            Assert.Contains("0x80", ex.Message);
        }

        [Fact]
        public void Open_reports_a_closed_inner_transport()
        {
            var inner = new ScriptedTransportInstance();
            inner.Close();
            var cotp = new CotpTransportInstance(inner);

            var sw = Stopwatch.StartNew();
            var ex = Assert.Throws<TransportException>(() => cotp.Open(0x01, 0x00, 0x03, 0x01));
            sw.Stop();

            Assert.Contains("not open", ex.Message, StringComparison.OrdinalIgnoreCase);

            // The timeout path would throw a TransportException too ("No COTP Connection
            // Confirm received within timeout"), so the elapsed time is what pins the
            // guard: it must fail fast, not burn the 5s wait.
            Assert.True(sw.Elapsed < TimeSpan.FromSeconds(1),
                $"expected a fast closed-transport failure, took {sw.Elapsed.TotalSeconds:F1}s");
        }

        [Fact]
        public void Read_before_Open_is_rejected()
        {
            var cotp = new CotpTransportInstance(new ScriptedTransportInstance());

            var ex = Assert.Throws<TransportException>(() => cotp.Read(1));
            Assert.Contains("handshake", ex.Message, StringComparison.OrdinalIgnoreCase);
        }

        [Fact]
        public void Write_before_Open_is_rejected()
        {
            var cotp = new CotpTransportInstance(new ScriptedTransportInstance());

            var ex = Assert.Throws<TransportException>(() => cotp.Write(new byte[] { 0x32 }));
            Assert.Contains("handshake", ex.Message, StringComparison.OrdinalIgnoreCase);
        }

        [Fact]
        public void Write_wraps_the_payload_in_a_data_transfer_frame()
        {
            var (cotp, inner) = Handshaken();
            inner.Written.Clear();

            cotp.Write(new byte[] { 0x32, 0x01, 0x00, 0x00 });

            var frame = Assert.Single(inner.Written);
            Assert.Equal(11, frame.Length);             // TPKT(4) + COTP DT(3) + payload(4)
            Assert.Equal(11, (frame[2] << 8) | frame[3]);
            Assert.Equal(0x02, frame[4]);               // LI
            Assert.Equal(0xF0, frame[5]);               // DT PDU type
            Assert.Equal(0x80, frame[6]);               // EOT
            Assert.Equal(new byte[] { 0x32, 0x01, 0x00, 0x00 }, frame.Skip(7).ToArray());
        }

        [Fact]
        public void Read_strips_the_data_transfer_header()
        {
            var (cotp, inner) = Handshaken();
            // A DT frame carrying four payload bytes.
            inner.Inject(0x03, 0x00, 0x00, 0x0B, 0x02, 0xF0, 0x80, 0x32, 0x03, 0x00, 0x00);

            Assert.Equal(4, cotp.GetNumBytesAvailable());
            Assert.Equal(new byte[] { 0x32, 0x03, 0x00, 0x00 }, cotp.Read(4));
        }

        [Fact]
        public void Read_rejects_a_frame_with_an_unknown_pdu_type()
        {
            var (cotp, inner) = Handshaken();
            // A frame whose PDU type is neither DT (0xF0) nor a COTP control frame is a
            // protocol violation and must not be handed to the S7 layer as payload.
            inner.Inject(0x03, 0x00, 0x00, 0x0B, 0x02, 0x01, 0x80, 0x32, 0x03, 0x00, 0x00);

            var ex = Assert.Throws<TransportException>(() => cotp.GetNumBytesAvailable());
            Assert.Contains("0x01", ex.Message);
        }

        [Theory]
        [InlineData(0x00, 0x00)]
        [InlineData(0x00, 0x03)]
        public void Read_rejects_a_tpkt_frame_length_smaller_than_its_header(byte lengthHigh, byte lengthLow)
        {
            var (cotp, inner) = Handshaken();
            inner.Inject(0x03, 0x00, lengthHigh, lengthLow);

            var ex = Assert.Throws<TransportException>(() => cotp.GetNumBytesAvailable());

            Assert.Contains("TPKT frame length", ex.Message);
        }

        [Fact]
        public void Read_rejects_a_tpkt_frame_with_an_invalid_version()
        {
            var (cotp, inner) = Handshaken();
            inner.Inject(0x04, 0x00, 0x00, 0x0B, 0x02, 0xF0, 0x80, 0x32, 0x03, 0x00, 0x00);

            var ex = Assert.Throws<TransportException>(() => cotp.GetNumBytesAvailable());

            Assert.Contains("Invalid TPKT header", ex.Message);
        }

        [Fact]
        public void Read_rejects_a_data_transfer_frame_with_an_invalid_length_indicator()
        {
            var (cotp, inner) = Handshaken();
            inner.Inject(0x03, 0x00, 0x00, 0x0B, 0x03, 0xF0, 0x80, 0x32, 0x03, 0x00, 0x00);

            var ex = Assert.Throws<TransportException>(() => cotp.GetNumBytesAvailable());

            Assert.Contains("length indicator", ex.Message);
        }

        [Fact]
        public void A_stray_confirm_after_the_handshake_is_skipped_not_fatal()
        {
            var (cotp, inner) = Handshaken();
            // A late Connection Confirm can legitimately precede the S7 response; it is
            // consumed and skipped, and the following DT frame still delivers its payload.
            // Disconnect and Error TPDUs are different: they end the session (see
            // A_session_ending_tpdu_from_the_peer_closes_the_transport).
            inner.Inject(ConnectionConfirm());
            inner.Inject(0x03, 0x00, 0x00, 0x0B, 0x02, 0xF0, 0x80, 0x32, 0x03, 0x00, 0x00);

            Assert.Equal(4, cotp.GetNumBytesAvailable());
            Assert.Equal(new byte[] { 0x32, 0x03, 0x00, 0x00 }, cotp.Read(4));
        }

        [Fact]
        public void Concurrent_Open_calls_send_only_one_connection_request()
        {
            var inner = new ScriptedTransportInstance();
            inner.Inject(ConnectionConfirm());
            var cotp = new CotpTransportInstance(inner);

            using var start = new ManualResetEventSlim(false);
            var threads = Enumerable.Range(0, 4)
                .Select(_ => new Thread(() =>
                {
                    start.Wait();
                    try { cotp.Open(0x01, 0x00, 0x03, 0x01); }
                    catch (TransportException) { /* ignored */ }
                }))
                .ToList();

            threads.ForEach(t => t.Start());
            start.Set();
            threads.ForEach(t => t.Join());

            // Exactly one Connection Request on the wire, however many callers raced.
            Assert.Single(inner.Written);
        }

        [Fact]
        public void Concurrent_Open_waits_for_the_owner_handshake_to_complete()
        {
            var inner = new ScriptedTransportInstance();
            var cotp = new CotpTransportInstance(inner)
            {
                HandshakeTimeout = TimeSpan.FromSeconds(2)
            };
            Exception? ownerFailure = null;
            Exception? followerFailure = null;
            using var ownerDone = new ManualResetEventSlim(false);
            using var followerEntered = new ManualResetEventSlim(false);
            using var followerDone = new ManualResetEventSlim(false);

            var owner = new Thread(() =>
            {
                try { cotp.Open(0x01, 0x00, 0x03, 0x01); }
                catch (Exception ex) { ownerFailure = ex; }
                finally { ownerDone.Set(); }
            });
            var follower = new Thread(() =>
            {
                followerEntered.Set();
                try { cotp.Open(0x01, 0x00, 0x03, 0x01); }
                catch (Exception ex) { followerFailure = ex; }
                finally { followerDone.Set(); }
            });

            owner.Start();
            Assert.True(SpinWait.SpinUntil(() => inner.WrittenCount == 1, TimeSpan.FromSeconds(1)));
            follower.Start();
            Assert.True(followerEntered.Wait(TimeSpan.FromSeconds(1)));
            Assert.False(followerDone.Wait(TimeSpan.FromMilliseconds(100)));

            inner.Inject(ConnectionConfirm());
            Assert.True(ownerDone.Wait(TimeSpan.FromSeconds(2)));
            Assert.True(followerDone.Wait(TimeSpan.FromSeconds(2)));
            owner.Join();
            follower.Join();

            Assert.Null(ownerFailure);
            Assert.Null(followerFailure);
            Assert.True(cotp.IsOpen);
            Assert.Single(inner.Written);
        }

        [Fact]
        public void Concurrent_Open_propagates_the_owner_handshake_failure()
        {
            var inner = new ScriptedTransportInstance();
            var cotp = new CotpTransportInstance(inner)
            {
                HandshakeTimeout = TimeSpan.FromSeconds(2)
            };
            Exception? ownerFailure = null;
            Exception? followerFailure = null;
            using var ownerDone = new ManualResetEventSlim(false);
            using var followerEntered = new ManualResetEventSlim(false);
            using var followerDone = new ManualResetEventSlim(false);

            var owner = new Thread(() =>
            {
                try { cotp.Open(0x01, 0x00, 0x03, 0x01); }
                catch (Exception ex) { ownerFailure = ex; }
                finally { ownerDone.Set(); }
            });
            var follower = new Thread(() =>
            {
                followerEntered.Set();
                try { cotp.Open(0x01, 0x00, 0x03, 0x01); }
                catch (Exception ex) { followerFailure = ex; }
                finally { followerDone.Set(); }
            });

            owner.Start();
            Assert.True(SpinWait.SpinUntil(() => inner.WrittenCount == 1, TimeSpan.FromSeconds(1)));
            follower.Start();
            Assert.True(followerEntered.Wait(TimeSpan.FromSeconds(1)));
            Assert.False(followerDone.Wait(TimeSpan.FromMilliseconds(100)));

            inner.Inject(0x04, 0x00, 0x00, 0x0B, 0x06, 0xD0, 0x00, 0x01, 0x00, 0x02, 0x00);
            Assert.True(ownerDone.Wait(TimeSpan.FromSeconds(2)));
            Assert.True(followerDone.Wait(TimeSpan.FromSeconds(2)));
            owner.Join();
            follower.Join();

            var ownerTransportFailure = Assert.IsType<TransportException>(ownerFailure);
            var followerTransportFailure = Assert.IsType<TransportException>(followerFailure);
            Assert.Equal(ownerTransportFailure.Message, followerTransportFailure.Message);
            Assert.True(cotp.IsOpen);
            Assert.Single(inner.Written);
        }

        [Fact]
        public void Concurrent_Open_unblocks_the_waiter_when_the_transport_closes()
        {
            var inner = new ScriptedTransportInstance();
            var cotp = new CotpTransportInstance(inner)
            {
                HandshakeTimeout = TimeSpan.FromSeconds(2)
            };
            Exception? ownerFailure = null;
            Exception? followerFailure = null;
            using var ownerDone = new ManualResetEventSlim(false);
            using var followerEntered = new ManualResetEventSlim(false);
            using var followerDone = new ManualResetEventSlim(false);

            var owner = new Thread(() =>
            {
                try { cotp.Open(0x01, 0x00, 0x03, 0x01); }
                catch (Exception ex) { ownerFailure = ex; }
                finally { ownerDone.Set(); }
            });
            var follower = new Thread(() =>
            {
                followerEntered.Set();
                try { cotp.Open(0x01, 0x00, 0x03, 0x01); }
                catch (Exception ex) { followerFailure = ex; }
                finally { followerDone.Set(); }
            });

            owner.Start();
            Assert.True(SpinWait.SpinUntil(() => inner.WrittenCount == 1, TimeSpan.FromSeconds(1)));
            follower.Start();
            Assert.True(followerEntered.Wait(TimeSpan.FromSeconds(1)));
            Assert.False(followerDone.Wait(TimeSpan.FromMilliseconds(100)));
            cotp.Close();

            Assert.True(ownerDone.Wait(TimeSpan.FromSeconds(2)));
            Assert.True(followerDone.Wait(TimeSpan.FromSeconds(2)));
            owner.Join();
            follower.Join();

            Assert.IsType<TransportException>(ownerFailure);
            Assert.IsType<TransportException>(followerFailure);
            Assert.True(inner.Closed);
        }

        [Fact]
        public void Open_fails_when_Close_wins_the_race_before_the_handshake_is_published()
        {
            var inner = new ScriptedTransportInstance();
            inner.Inject(ConnectionConfirm());
            var cotp = new CotpTransportInstance(inner);

            // The diagnostic hook fires once the Confirm has been validated, immediately
            // before Open() takes the lock to publish the handshake as complete: the window
            // a concurrent Close() can hit. Running Close() from the hook makes that
            // interleaving deterministic.
            cotp.DiagnosticOutput = new CallbackTextWriter("Connection Confirm (CC)", () => cotp.Close());

            var ex = Assert.Throws<TransportException>(() => cotp.Open(0x01, 0x00, 0x03, 0x01));

            Assert.Contains("closed", ex.Message, StringComparison.OrdinalIgnoreCase);
            Assert.False(cotp.IsOpen);
            // A Confirm that arrives after Close() must not resurrect the ended session.
            var writeEx = Assert.Throws<TransportException>(() => cotp.Write(new byte[] { 0x32 }));
            Assert.Contains("handshake", writeEx.Message, StringComparison.OrdinalIgnoreCase);
        }

        [Fact]
        public async Task Close_releases_a_write_stalled_on_the_inner_transport()
        {
            var inner = new StalledWriteTransportInstance();
            inner.Inject(ConnectionConfirm());
            var cotp = new CotpTransportInstance(inner);
            cotp.Open(0x01, 0x00, 0x03, 0x01);
            inner.StallWrites();

            Exception? writeFailure = null;
            using var writeDone = new ManualResetEventSlim(false);
            var writer = new Thread(() =>
            {
                try { cotp.Write(new byte[] { 0x32, 0x01 }); }
                catch (Exception ex) { writeFailure = ex; }
                finally { writeDone.Set(); }
            })
            { IsBackground = true };

            try
            {
                writer.Start();
                Assert.True(inner.WaitForStalledWrite(TimeSpan.FromSeconds(2)),
                    "the write never reached the inner transport");

                // Closing the inner transport is the only thing that can release a send that
                // is blocked on a peer that stopped reading, so Close() must not queue behind
                // the writer.
                await AssertCompletes(Task.Run(() => cotp.Close()), "Close() hung behind the stalled write");
                Assert.True(writeDone.Wait(TimeSpan.FromSeconds(2)), "the stalled write was never released");
                Assert.IsType<TransportException>(writeFailure);
                Assert.True(inner.Closed);
            }
            finally
            {
                // Never leave the writer parked when an assertion above failed.
                inner.Close();
                writer.Join(TimeSpan.FromSeconds(5));
            }
        }

        [Fact]
        public async Task A_stalled_write_does_not_block_reads()
        {
            var inner = new StalledWriteTransportInstance();
            inner.Inject(ConnectionConfirm());
            var cotp = new CotpTransportInstance(inner);
            cotp.Open(0x01, 0x00, 0x03, 0x01);
            inner.StallWrites();

            Exception? writeFailure = null;
            var writer = new Thread(() =>
            {
                try { cotp.Write(new byte[] { 0x32, 0x01 }); }
                catch (Exception ex) { writeFailure = ex; }
            })
            { IsBackground = true };

            try
            {
                writer.Start();
                Assert.True(inner.WaitForStalledWrite(TimeSpan.FromSeconds(2)),
                    "the write never reached the inner transport");

                // A driver polls GetNumBytesAvailable()/Read() while a request is in flight; a
                // send stuck on a stalled peer must not freeze those calls.
                inner.Inject(0x03, 0x00, 0x00, 0x0B, 0x02, 0xF0, 0x80, 0x32, 0x03, 0x00, 0x00);
                var available = await ResultWithin(
                    Task.Run(() => cotp.GetNumBytesAvailable()),
                    "GetNumBytesAvailable() blocked behind the stalled write");
                Assert.Equal(4, available);

                var payload = await ResultWithin(
                    Task.Run(() => cotp.Read(4)),
                    "Read() blocked behind the stalled write");
                Assert.Equal(new byte[] { 0x32, 0x03, 0x00, 0x00 }, payload);
            }
            finally
            {
                inner.Close();
                writer.Join(TimeSpan.FromSeconds(5));
            }

            // Closing the inner transport is what ends the stalled write, and it must end with
            // the transport's own failure rather than something unrelated.
            Assert.IsType<TransportException>(writeFailure);
        }

        [Fact]
        public void Concurrent_writes_keep_the_fragments_of_each_message_contiguous()
        {
            var inner = new SlowWriteTransportInstance();
            inner.Inject(ConnectionConfirm());
            var cotp = new CotpTransportInstance(inner);
            cotp.Open(0x01, 0x00, 0x03, 0x01);
            inner.Written.Clear();
            inner.SlowWrites();

            // The default TPDU size of 1024 leaves 1021 payload bytes per DT frame, so each of
            // these messages spans three fragments.
            var messages = new[]
            {
                Enumerable.Repeat((byte)0xAA, 2500).ToArray(),
                Enumerable.Repeat((byte)0xBB, 2500).ToArray()
            };
            var failures = new ConcurrentBag<Exception>();
            using var start = new ManualResetEventSlim(false);
            var threads = messages
                .Select(message => new Thread(() =>
                {
                    start.Wait();
                    try { cotp.Write(message); }
                    catch (Exception ex) { failures.Add(ex); }
                }))
                .ToList();

            threads.ForEach(t => t.Start());
            start.Set();
            threads.ForEach(t => t.Join());

            Assert.Empty(failures);
            Assert.Equal(6, inner.Written.Count);
            for (var i = 0; i < inner.Written.Count; i++)
            {
                var frame = inner.Written[i];
                var message = frame[7];

                // One message per frame, and frames 0-2 and 3-5 belong together: a second writer
                // slipping a fragment in between would split a message on the wire.
                Assert.All(frame.Skip(7), b => Assert.Equal(message, b));
                Assert.Equal(inner.Written[i / 3 * 3][7], message);

                // TPDU numbers follow the order on the wire, and only the last fragment of a
                // message carries EOT.
                Assert.Equal(i, frame[6] & 0x7F);
                Assert.Equal(i % 3 == 2, (frame[6] & 0x80) != 0);
            }
            Assert.NotEqual(inner.Written[0][7], inner.Written[3][7]);
        }

        [Fact]
        public void Concurrent_reads_do_not_interfere()
        {
            var (cotp, inner) = Handshaken();
            // Two DT frames, four payload bytes each.
            inner.Inject(0x03, 0x00, 0x00, 0x0B, 0x02, 0xF0, 0x80, 0x11, 0x22, 0x33, 0x44);
            inner.Inject(0x03, 0x00, 0x00, 0x0B, 0x02, 0xF0, 0x80, 0xAA, 0xBB, 0xCC, 0xDD);

            var results = new ConcurrentBag<byte[]>();
            using var start = new ManualResetEventSlim(false);
            var threads = Enumerable.Range(0, 2)
                .Select(_ => new Thread(() =>
                {
                    start.Wait();
                    try { results.Add(cotp.Read(4)); }
                    catch (TransportException) { /* ignored */ }
                }))
                .ToList();

            threads.ForEach(t => t.Start());
            start.Set();
            threads.ForEach(t => t.Join());

            // Each reader must get exactly its own frame's payload, not a duplicated or
            // merged buffer.
            Assert.Equal(2, results.Count);
            Assert.Contains(results, r => r.SequenceEqual(new byte[] { 0x11, 0x22, 0x33, 0x44 }));
            Assert.Contains(results, r => r.SequenceEqual(new byte[] { 0xAA, 0xBB, 0xCC, 0xDD }));
        }

        [Fact]
        public void Open_accepts_a_confirm_that_does_not_echo_the_src_ref()
        {
            // Java parity: the reference checks only the PDU type, and some ISO-on-TCP
            // peers reply with a DST-REF of 0x0000. The strict echo check was removed.
            var inner = new ScriptedTransportInstance();
            inner.Inject(0x03, 0x00, 0x00, 0x0B, 0x06, 0xD0, 0x00, 0x00, 0x00, 0x02, 0x00);
            var cotp = new CotpTransportInstance(inner);

            cotp.Open(0x01, 0x00, 0x03, 0x01);

            Assert.True(cotp.IsOpen);
        }

        [Fact]
        public void Open_rejects_a_frame_with_a_bad_tpkt_version()
        {
            var inner = new ScriptedTransportInstance();
            // TPKT version byte is 0x04 instead of 0x03.
            inner.Inject(0x04, 0x00, 0x00, 0x0B, 0x06, 0xD0, 0x00, 0x01, 0x00, 0x02, 0x00);
            var cotp = new CotpTransportInstance(inner);

            var ex = Assert.Throws<TransportException>(() => cotp.Open(0x01, 0x00, 0x03, 0x01));
            Assert.Contains("Invalid TPKT header", ex.Message);
        }

        [Fact]
        public void Open_rejects_a_connection_confirm_with_a_nonzero_tpkt_reserved_byte()
        {
            var inner = new ScriptedTransportInstance();
            inner.Inject(0x03, 0x01, 0x00, 0x0B, 0x06, 0xD0, 0x00, 0x01, 0x00, 0x02, 0x00);
            var cotp = new CotpTransportInstance(inner);

            var ex = Assert.Throws<TransportException>(() => cotp.Open(0x01, 0x00, 0x03, 0x01));
            Assert.Contains("Invalid TPKT header", ex.Message);
        }

        [Fact]
        public void Open_rejects_a_connection_confirm_with_an_invalid_length_indicator()
        {
            var inner = new ScriptedTransportInstance();
            inner.Inject(0x03, 0x00, 0x00, 0x0B, 0x00, 0xD0, 0x00, 0x01, 0x00, 0x02, 0x00);
            var cotp = new CotpTransportInstance(inner);

            var ex = Assert.Throws<TransportException>(() => cotp.Open(0x01, 0x00, 0x03, 0x01));
            Assert.Contains("length indicator", ex.Message);
        }

        [Fact]
        public void Open_rejects_a_connection_confirm_whose_length_indicator_exceeds_the_frame()
        {
            var inner = new ScriptedTransportInstance();
            inner.Inject(0x03, 0x00, 0x00, 0x0B, 0x09, 0xD0, 0x00, 0x01, 0x00, 0x02, 0x00);
            var cotp = new CotpTransportInstance(inner);

            var ex = Assert.Throws<TransportException>(() => cotp.Open(0x01, 0x00, 0x03, 0x01));
            Assert.Contains("requires", ex.Message);
        }

        [Fact]
        public void Open_reports_a_peer_that_closes_during_the_handshake()
        {
            var inner = new ScriptedTransportInstance();
            var cotp = new CotpTransportInstance(inner)
            {
                HandshakeTimeout = TimeSpan.FromMilliseconds(2000)
            };

            Exception? captured = null;
            using var done = new ManualResetEventSlim(false);
            var openThread = new Thread(() =>
            {
                try { cotp.Open(0x01, 0x00, 0x03, 0x01); }
                catch (Exception ex) { captured = ex; }
                finally { done.Set(); }
            });

            openThread.Start();

            // Wait until the Open thread has actually sent the Connection Request and is
            // inside its poll loop — closing before that would fail the CR write instead
            // of exercising the wait loop's closed-transport detection. This also removes
            // the timing dependency that could flake under CI load.
            var crDeadline = Environment.TickCount64 + 500;
            while (inner.WrittenCount == 0 && Environment.TickCount64 < crDeadline)
            {
                Thread.Sleep(5);
            }
            Assert.NotEmpty(inner.Written);
            inner.Close();

            Assert.True(done.Wait(TimeSpan.FromSeconds(4)), "Open() did not finish in time");
            openThread.Join();

            var ex = Assert.IsType<TransportException>(captured);
            Assert.Contains("closed", ex.Message, StringComparison.OrdinalIgnoreCase);
            Assert.True(captured == null || ex.Message.Contains("closed"),
                $"expected a closed-transport error, got: {captured?.Message}");
        }

        [Fact]
        public void Write_splits_a_payload_across_multiple_DT_frames()
        {
            var (cotp, inner) = Handshaken();
            inner.Written.Clear();

            // Default TPDU size is 1024; minus the 3-byte DT header, max
            // payload per frame is 1021. A 1024-byte payload spans 2 frames:
            // frame 0 = 1021 bytes (no EOT), frame 1 = 3 bytes (EOT).
            cotp.Write(new byte[1024]);

            Assert.Equal(2, inner.Written.Count);
            Assert.Equal(0x00, inner.Written[0][6] & 0x7F);   // TPDU-NR = 0
            Assert.Equal(0x00, inner.Written[0][6] & 0x80);   // no EOT
            Assert.Equal(0x01, inner.Written[1][6] & 0x7F);   // TPDU-NR = 1
            Assert.Equal(0x80, inner.Written[1][6] & 0x80);   // EOT set
        }

        [Fact]
        public void Write_that_fits_in_one_frame_sends_EOT()
        {
            var (cotp, inner) = Handshaken();
            inner.Written.Clear();

            cotp.Write(new byte[10]);

            Assert.Single(inner.Written);
            Assert.Equal(0x80, inner.Written[0][6] & 0x80);   // EOT
        }

        [Fact]
        public void Open_waits_for_the_frame_body_after_the_header()
        {
            var inner = new ScriptedTransportInstance();
            var cotp = new CotpTransportInstance(inner)
            {
                // The header arrives at ~1400ms of this window; the body at ~2500ms. A
                // reused deadline would expire at 2000ms and abort; the separate frame-body
                // deadline starts fresh at ~1400ms and accepts the body. The ~600ms of
                // slack on each side is what keeps the test stable on a loaded CI runner.
                HandshakeTimeout = TimeSpan.FromMilliseconds(2000)
            };

            Exception? captured = null;
            using var done = new ManualResetEventSlim(false);
            var openThread = new Thread(() =>
            {
                try { cotp.Open(0x01, 0x00, 0x03, 0x01); }
                catch (Exception ex) { captured = ex; }
                finally { done.Set(); }
            });

            openThread.Start();

            // Deliver the TPKT header well inside the first deadline, then the body only
            // after that deadline has passed. The frame-body wait gets its own fresh
            // deadline, so Open still succeeds. If the deadline were reused (the original
            // bug), the body arriving after it expired would abort the handshake with
            // "Incomplete COTP CC frame".
            Thread.Sleep(1400);
            inner.Inject(0x03, 0x00, 0x00, 0x0B);                       // TPKT header
            Thread.Sleep(1100);                                         // past the header deadline
            inner.Inject(0x06, 0xD0, 0x00, 0x01, 0x00, 0x02, 0x00);     // frame body

            Assert.True(done.Wait(TimeSpan.FromSeconds(6)), "Open() did not finish in time");
            openThread.Join();

            Assert.Null(captured);
            Assert.True(cotp.IsOpen);
        }

        // The TPDUs that end a COTP session, as s7.mspec defines them (TPKT-wrapped, with the
        // LI byte first): Disconnect Request 0x80, Disconnect Confirm 0xC0 and TPDU Error 0x70.
        public static IEnumerable<object[]> SessionEndingTpdus() => new[]
        {
            new object[] { 0x80, TpktFrame.Wrap(new byte[] { 0x06, 0x80, 0x00, 0x01, 0x00, 0x02, 0x00 }) },
            new object[] { 0xC0, TpktFrame.Wrap(new byte[] { 0x05, 0xC0, 0x00, 0x01, 0x00, 0x02 }) },
            new object[] { 0x70, TpktFrame.Wrap(new byte[] { 0x04, 0x70, 0x00, 0x01, 0x05 }) }
        };

        [Theory]
        [MemberData(nameof(SessionEndingTpdus))]
        public void A_session_ending_tpdu_from_the_peer_closes_the_transport(int pduType, byte[] frame)
        {
            var (cotp, inner) = Handshaken();
            // The session-ending TPDU, followed by a DT frame that must NOT be reached.
            inner.Inject(frame);
            inner.Inject(0x03, 0x00, 0x00, 0x0B, 0x02, 0xF0, 0x80, 0x32, 0x03, 0x00, 0x00);

            var ex = Assert.Throws<TransportException>(() => cotp.GetNumBytesAvailable());
            Assert.Contains("disconnected", ex.Message, StringComparison.OrdinalIgnoreCase);
            Assert.Contains($"0x{pduType:X2}", ex.Message);

            // Surfacing the teardown is not enough: the transport must stop claiming a live
            // session, or the next Write() goes out on a connection the peer already ended.
            Assert.True(inner.Closed);
            Assert.False(cotp.IsOpen);

            // Every later call names the teardown. Blaming a missing Open() would send the
            // caller looking for a mistake in their own call order.
            AssertTeardownReported(() => cotp.Write(new byte[] { 0x32 }), pduType);
            AssertTeardownReported(() => cotp.Read(1), pduType);
            AssertTeardownReported(() => cotp.PeekReadableBytes(1), pduType);
            AssertTeardownReported(() => cotp.GetNumBytesAvailable(), pduType);
            AssertTeardownReported(() => cotp.Open(0x01, 0x00, 0x03, 0x01), pduType);
        }

        [Theory]
        [MemberData(nameof(SessionEndingTpdus))]
        public void A_response_delivered_before_a_peer_teardown_stays_readable(int pduType, byte[] frame)
        {
            var (cotp, inner) = Handshaken();
            // A complete DT frame carrying four payload bytes, then the session-ending TPDU.
            inner.Inject(0x03, 0x00, 0x00, 0x0B, 0x02, 0xF0, 0x80, 0x32, 0x03, 0x00, 0x00);
            inner.Inject(frame);

            // The response is reported first; the teardown queued behind it does not hide it.
            Assert.Equal(4, cotp.GetNumBytesAvailable());

            // The teardown has been consumed by now, so the session no longer claims to be live.
            Assert.True(inner.Closed);
            Assert.False(cotp.IsOpen);
            AssertTeardownReported(() => cotp.Write(new byte[] { 0x32 }), pduType);

            // The bytes the peer delivered before ending the session are still handed out.
            Assert.Equal(new byte[] { 0x32, 0x03, 0x00, 0x00 }, cotp.PeekReadableBytes(4));
            Assert.Equal(new byte[] { 0x32, 0x03, 0x00, 0x00 }, cotp.Read(4));

            // Only once they are gone does the teardown surface.
            AssertTeardownReported(() => cotp.GetNumBytesAvailable(), pduType);
            AssertTeardownReported(() => cotp.Read(1), pduType);
        }

        [Fact]
        public void A_read_for_more_than_was_delivered_before_a_peer_teardown_reports_the_teardown()
        {
            var (cotp, inner) = Handshaken();
            inner.Inject(0x03, 0x00, 0x00, 0x0B, 0x02, 0xF0, 0x80, 0x32, 0x03, 0x00, 0x00);
            inner.Inject(TpktFrame.Wrap(new byte[] { 0x06, 0x80, 0x00, 0x01, 0x00, 0x02, 0x00 }));

            // Four bytes arrived and then the peer ended the session: the other four never
            // will, so a caller waiting for them is told why instead of "only 4 available".
            AssertTeardownReported(() => cotp.Read(8), 0x80);

            // The bytes that did arrive are not lost by the failed read.
            Assert.Equal(new byte[] { 0x32, 0x03, 0x00, 0x00 }, cotp.Read(4));
        }

        private static void AssertTeardownReported(Action call, int pduType)
        {
            var ex = Assert.Throws<TransportException>(call);

            Assert.Contains("disconnected", ex.Message, StringComparison.OrdinalIgnoreCase);
            Assert.Contains($"0x{pduType:X2}", ex.Message);
            Assert.DoesNotContain("Open() first", ex.Message);
        }

        [Fact]
        public void Write_fragments_at_the_negotiated_tpdu_size()
        {
            var inner = new ScriptedTransportInstance();
            // Confirm carrying a TPDU-size parameter (0xC0, length 1, value 0x09 =
            // 512 bytes) — smaller than the 1024 the CR asked for.
            inner.Inject(0x03, 0x00, 0x00, 0x0E, 0x09, 0xD0, 0x00, 0x01, 0x00, 0x02,
                         0x00, 0xC0, 0x01, 0x09);
            var cotp = new CotpTransportInstance(inner);
            cotp.Open(0x01, 0x00, 0x03, 0x01);
            inner.Written.Clear();

            // TPDU size 512 − 3 (DT header) = 509 bytes per frame.
            // A 510-byte payload spans 2 frames.
            cotp.Write(new byte[510]);

            Assert.Equal(2, inner.Written.Count);
            Assert.Equal(0x00, inner.Written[0][6] & 0x80);   // first: no EOT
            Assert.Equal(0x80, inner.Written[1][6] & 0x80);   // second: EOT
        }

        [Fact]
        public void Open_ignores_parameters_outside_the_connection_confirm_length_indicator()
        {
            var inner = new ScriptedTransportInstance();
            // LI=6 ends after the fixed CC header. The apparent C0 parameter is trailing
            // data and must not shrink the negotiated TPDU size to 64 bytes.
            inner.Inject(0x03, 0x00, 0x00, 0x0E, 0x06, 0xD0, 0x00, 0x01, 0x00, 0x02,
                         0x00, 0xC0, 0x01, 0x06);
            var cotp = new CotpTransportInstance(inner);
            cotp.Open(0x01, 0x00, 0x03, 0x01);
            inner.Written.Clear();

            cotp.Write(new byte[64]);

            Assert.Single(inner.Written);
        }

        [Fact]
        public void Open_caps_a_larger_confirmed_tpdu_size_to_the_requested_size()
        {
            var inner = new ScriptedTransportInstance();
            // The CR asks for 1024 bytes; a larger CC must not raise that limit.
            inner.Inject(0x03, 0x00, 0x00, 0x0E, 0x09, 0xD0, 0x00, 0x01, 0x00, 0x02,
                         0x00, 0xC0, 0x01, 0x0D);
            var cotp = new CotpTransportInstance(inner);
            cotp.Open(0x01, 0x00, 0x03, 0x01);
            inner.Written.Clear();

            cotp.Write(new byte[1024]);

            Assert.Equal(2, inner.Written.Count);
            Assert.Equal(1028, inner.Written[0].Length); // TPKT(4) + DT header(3) + 1021 payload
            Assert.Equal(0x00, inner.Written[0][6] & 0x80);
            Assert.Equal(0x80, inner.Written[1][6] & 0x80);
        }

        [Fact]
        public void An_incomplete_frame_is_left_for_the_next_read()
        {
            var (cotp, inner) = Handshaken();
            // Header announces 11 bytes but only 8 have arrived.
            inner.Inject(0x03, 0x00, 0x00, 0x0B, 0x02, 0xF0, 0x80, 0x32);

            Assert.Equal(0, cotp.GetNumBytesAvailable());

            // The rest arrives; now the whole payload is readable.
            inner.Inject(0x03, 0x00, 0x00);
            Assert.Equal(4, cotp.GetNumBytesAvailable());
        }

        [Fact]
        public void Close_propagates_to_the_inner_transport()
        {
            var (cotp, inner) = Handshaken();

            cotp.Close();

            Assert.True(inner.Closed);
            Assert.False(cotp.IsOpen);
        }
    }
}
