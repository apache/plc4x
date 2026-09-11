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
using System.Linq;
using System.Threading;
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

        public ScriptedTransportInstance() : base(new ScriptedTransportConfiguration()) { }

        public List<byte[]> Written { get; } = new List<byte[]>();

        public bool Closed { get; private set; }

        public override bool IsOpen => !Closed;

        public void Inject(params byte[] bytes) => _inbound.AddRange(bytes);

        public override int GetNumBytesAvailable() => _inbound.Count;

        public override byte[] PeekReadableBytes(int numBytes)
        {
            if (numBytes > _inbound.Count)
            {
                throw new TransportException(
                    $"Requested {numBytes} bytes but only {_inbound.Count} available.");
            }
            return _inbound.Take(numBytes).ToArray();
        }

        public override byte[] Read(int numBytes)
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

        public override void Write(byte[] bytes) => Written.Add(bytes);

        public override void Close() => Closed = true;
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

        [Fact]
        public void A_stray_confirm_after_the_handshake_is_skipped_not_fatal()
        {
            var (cotp, inner) = Handshaken();
            // A late Confirm (or a DR/ER teardown) can legitimately precede the S7
            // response; it is consumed and skipped, and the following DT frame still
            // delivers its payload.
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
            Assert.Contains("TPKT version", ex.Message);
        }

        [Fact]
        public void Open_reports_a_peer_that_closes_during_the_handshake()
        {
            var inner = new ScriptedTransportInstance();
            var cotp = new CotpTransportInstance(inner)
            {
                HandshakeTimeout = TimeSpan.FromMilliseconds(2000)
            };

            Exception captured = null;
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
            while (inner.Written.Count == 0 && Environment.TickCount64 < crDeadline)
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

            Exception captured = null;
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

        [Fact]
        public void A_disconnect_request_from_the_peer_is_surfaced_not_swallowed()
        {
            var (cotp, inner) = Handshaken();
            // A Disconnect Request, followed by a DT frame that must NOT be reached.
            inner.Inject(0x03, 0x00, 0x00, 0x0B, 0x02, 0x80, 0x00, 0x00);
            inner.Inject(0x03, 0x00, 0x00, 0x0B, 0x02, 0xF0, 0x80, 0x32, 0x03, 0x00, 0x00);

            var ex = Assert.Throws<TransportException>(() => cotp.GetNumBytesAvailable());
            Assert.Contains("disconnected", ex.Message, StringComparison.OrdinalIgnoreCase);
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
