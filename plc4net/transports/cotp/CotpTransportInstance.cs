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
using System.Threading.Tasks;
using org.apache.plc4net.spi.transports;

namespace org.apache.plc4net.transports.cotp
{
    /// <summary>
    /// Wraps a TCP transport instance, adding TPKT (RFC 1006) framing and
    /// COTP protocol support for S7 communication (connection establishment
    /// via CR/CC handshake, Data Transfer via DT frames).
    /// </summary>
    public class CotpTransportInstance : BaseTransportInstance, IAsyncTransportInstance
    {
        private readonly ITransportInstance _inner;
        private readonly object _handshakeLock = new object();

        // Serializes Write() so the DT fragments of one message stay contiguous on the wire.
        // Deliberately not _handshakeLock: the socket write can block indefinitely against a
        // peer that has stopped reading, and neither Close() nor the polling reads may queue
        // behind it. The lock order is _writeLock then _handshakeLock; nothing takes them the
        // other way round.
        private readonly object _writeLock = new object();
        private byte[] _leftover = Array.Empty<byte>();
        private const int RequestedTpduSize = 1024;

        // Volatile so the pre-handshake fast path in Read/Write/Peek/GetNumBytesAvailable
        // can be read without the lock; Open() sets it only after the handshake completes.
        private volatile bool _handshakeDone;

        // Set while Open() is waiting for the Confirm. Read under _handshakeLock only;
        // a concurrent Close() clears it (and closes the inner transport) so the poll
        // loops below notice and fail fast instead of burning the whole timeout.
        private volatile bool _handshakeInProgress;
        private TaskCompletionSource<bool>? _handshakeCompletion;

        // Why the peer ended the session (Disconnect Request, Disconnect Confirm or TPDU
        // Error); null while it is alive. Read and written under _handshakeLock. Payload the
        // peer delivered before that stays in _leftover and remains readable. Once none is
        // left, every call reports this instead of a missing handshake.
        private string? _peerTeardown;

        // Negotiated COTP TPDU size, taken from the Confirm's TPDU-size parameter. The
        // Connection Request asks for 1024 bytes; a PLC may confirm a smaller value
        // (S7-1200/1500 commonly confirm 480), and Data Transfer frames must respect it.
        private int _tpduSize = RequestedTpduSize;

        /// <summary>
        /// How long Open() waits for the Connection Confirm before giving up. Injectable
        /// for tests that exercise split header/body arrival without a five-second wait.
        /// </summary>
        public TimeSpan HandshakeTimeout { get; set; } = TimeSpan.FromSeconds(5);

        /// <summary>
        /// Optional diagnostic output. When set, every raw frame sent or received
        /// (CR, CC, DT payloads) is written here as annotated hex so that a
        /// verification tool or diagnostic harness can capture the wire-level
        /// conversation without external packet capture.
        /// </summary>
        public TextWriter? DiagnosticOutput { get; set; }

        // Tracks buffered bytes from fully-consumed frames for accurate
        // availability reporting across partial-DT-header scenarios.
        private int _bufferedFramePayloads;

        public CotpTransportInstance(ITransportInstance inner)
            : base((inner ?? throw new ArgumentNullException(nameof(inner))).Configuration)
        {
            _inner = inner;
            DriverConfig = inner.DriverConfig;
        }

        public override bool IsOpen => _inner.IsOpen;

        public override int GetNumBytesAvailable()
        {
            // Fast path: before the handshake there is nothing to report, and taking the
            // lock here would block behind a concurrent Open() polling for the CC for up
            // to the full handshake timeout.
            if (!_handshakeDone) return 0;

            lock (_handshakeLock)
            {
                // Re-check under the lock: a Close() that ran between the fast path and
                // here has cleared the flag and closed the inner transport, so draining
                // now would serve stale bytes from a dead session.
                if (!_handshakeDone) return 0;

                // Sum of: leftover bytes from prior partial reads + buffered frame
                // payloads from fully-consumed frames. The inner transport count is
                // only used to detect frame-completeness, not to report phantom bytes.
                // The drain mutates _leftover, so it must run under the same lock that
                // Read()/PeekReadableBytes() hold when they mutate it.
                DrainFramesIfAvailable();
                var available = _leftover.Length + _bufferedFramePayloads;

                // Nothing left to hand out and nothing more can arrive: say why, instead of an
                // endless 0 that looks like a slow peer.
                if (available == 0 && _peerTeardown != null)
                    throw new TransportException(_peerTeardown);

                return available;
            }
        }

        public override byte[] PeekReadableBytes(int numBytes)
        {
            if (numBytes <= 0) return Array.Empty<byte>();

            // Fast path, lock-free: fail immediately instead of blocking behind a
            // concurrent Open() for the full handshake timeout.
            if (!_handshakeDone)
                throw new TransportException("COTP handshake has not been performed. Call Open() first.");

            lock (_handshakeLock)
            {
                if (!_handshakeDone)
                    throw new TransportException("COTP handshake has not been performed. Call Open() first.");

                DrainFramesIfAvailable();

                if (_leftover.Length < numBytes)
                {
                    // A session the peer has ended will never deliver the missing bytes, so
                    // say so rather than suggest they may still arrive.
                    if (_peerTeardown != null)
                        throw new TransportException(_peerTeardown);

                    throw new TransportException(
                        $"Requested {numBytes} bytes but only {_leftover.Length} available.");
                }

                var result = new byte[numBytes];
                Array.Copy(_leftover, 0, result, 0, numBytes);
                return result;
            }
        }

        public override byte[] Read(int numBytes)
        {
            if (numBytes <= 0) return Array.Empty<byte>();

            // Fast path, lock-free: fail immediately instead of blocking behind a
            // concurrent Open() for the full handshake timeout.
            if (!_handshakeDone)
                throw new TransportException("COTP handshake has not been performed. Call Open() first.");

            lock (_handshakeLock)
            {
                if (!_handshakeDone)
                    throw new TransportException("COTP handshake has not been performed. Call Open() first.");

                DrainFramesIfAvailable();

                if (_leftover.Length < numBytes)
                {
                    // A session the peer has ended will never deliver the missing bytes, so
                    // say so rather than suggest they may still arrive.
                    if (_peerTeardown != null)
                        throw new TransportException(_peerTeardown);

                    throw new TransportException(
                        $"Requested {numBytes} bytes but only {_leftover.Length} available.");
                }

                var result = new byte[numBytes];
                Array.Copy(_leftover, 0, result, 0, numBytes);

                // Shift leftover.
                var remaining = _leftover.Length - numBytes;
                if (remaining > 0)
                {
                    var newLeftover = new byte[remaining];
                    Array.Copy(_leftover, numBytes, newLeftover, 0, remaining);
                    _leftover = newLeftover;
                }
                else
                {
                    _leftover = Array.Empty<byte>();
                }

                return result;
            }
        }

        // ── COTP connection establishment ────────────────────────

        /// <summary>
        /// Performs the COTP Connection Request / Confirm handshake against
        /// the inner TCP transport so the S7 driver can send Data Transfer
        /// frames afterwards.
        /// </summary>
        /// <param name="localTsapHi">High byte of the local TSAP (e.g. 0x01).</param>
        /// <param name="localTsapLo">Low byte of the local TSAP (e.g. 0x00).</param>
        /// <param name="remoteTsapHi">High byte of the remote TSAP (computed from rack/slot).</param>
        /// <param name="remoteTsapLo">Low byte of the remote TSAP.</param>
        public void Open(byte localTsapHi, byte localTsapLo, byte remoteTsapHi, byte remoteTsapLo)
        {
            // Only the state transition and the CR write need the lock. The polling below
            // runs WITHOUT the lock, so a concurrent Close() — the natural way to abort a
            // hanging handshake — is not blocked for the full timeout. Close() clears
            // _handshakeInProgress and closes the inner transport; the poll loops notice
            // via IsOpen and fail fast.
            TaskCompletionSource<bool>? existingHandshake = null;
            TaskCompletionSource<bool>? ownedHandshake = null;
            lock (_handshakeLock)
            {
                // A session the peer has ended stays ended. Without this the check below would
                // report success for a connection that is gone.
                if (_peerTeardown != null) throw new TransportException(_peerTeardown);
                if (_handshakeDone) return;
                if (_handshakeInProgress)
                {
                    existingHandshake = _handshakeCompletion;
                }
                else
                {
                    if (!_inner.IsOpen) throw new TransportException("Inner transport is not open.");
                    _handshakeInProgress = true;
                    ownedHandshake = new TaskCompletionSource<bool>(
                        TaskCreationOptions.RunContinuationsAsynchronously);
                    _handshakeCompletion = ownedHandshake;
                }
            }

            // A concurrent caller must observe the owner handshake's result. Returning
            // early would let it issue a write against a session that is not ready yet.
            if (existingHandshake != null)
            {
                existingHandshake.Task.GetAwaiter().GetResult();
                return;
            }

            try
            {
                // ── Build COTP Connection Request (CR) ──
                //    LI(1) + CR(0xE0) + DST-REF(2) + SRC-REF(2) + class(1)
                //    + TPDU-size(3) + calling-TSAP(4) + called-TSAP(4) = 18 bytes, LI = 17.
                var cr = new byte[18];
                cr[0] = 0x11;   // LI = 17 (everything after this byte)
                cr[1] = 0xE0;   // CR
                cr[2] = 0x00; cr[3] = 0x00;  // DST-REF = 0 (initial)
                cr[4] = 0x00; cr[5] = 0x01;  // SRC-REF = 1
                cr[6] = 0x00;   // class 0, no options
                cr[7] = 0xC0; cr[8] = 0x01; cr[9] = 0x0A;  // TPDU size = 1024 (requested)
                cr[10] = 0xC1; cr[11] = 0x02;               // calling TSAP
                cr[12] = localTsapHi; cr[13] = localTsapLo;
                cr[14] = 0xC2; cr[15] = 0x02;               // called TSAP
                cr[16] = remoteTsapHi; cr[17] = remoteTsapLo;

                // Send CR through the inner transport — no DT header wrapping.
                var crFrame = TpktFrame.Wrap(cr);
                _inner.Write(crFrame);
                WriteHex(crFrame, "SEND  COTP Connection Request (CR) TPKT frame");

                // ── Read COTP Connection Confirm (CC) ──
                // Wait for the TPKT header to arrive. A peer that rejects the CR commonly
                // closes the TCP connection right away; detect that instead of burning the
                // full timeout, so the caller sees "connection closed" not "timed out".
                var timeoutMs = (int)HandshakeTimeout.TotalMilliseconds;
                var deadline = Environment.TickCount64 + timeoutMs;
                while (Environment.TickCount64 < deadline)
                {
                    if (!_inner.IsOpen)
                        throw new TransportException(
                            "Inner transport closed while waiting for the COTP Connection Confirm.");
                    var available = _inner.GetNumBytesAvailable();
                    if (available >= 4) break;
                    System.Threading.Thread.Sleep(5);
                }

                if (_inner.GetNumBytesAvailable() < 4)
                    throw new TransportException("No COTP Connection Confirm received within timeout.");

                // Peek TPKT header to learn the frame size.
                var tpktHeader = _inner.PeekReadableBytes(4);
                if (tpktHeader[0] != TpktFrame.Version || tpktHeader[1] != 0)
                {
                    throw new TransportException(
                        $"Invalid TPKT header: version 0x{tpktHeader[0]:X2}, reserved 0x{tpktHeader[1]:X2}.");
                }
                var frameLen = ((tpktHeader[2] << 8) | tpktHeader[3]);
                if (frameLen < 4 + 7) // TPKT(4) + LI(1) + CC(1) + DST-REF(2) + SRC-REF(2) + class(1)
                    throw new TransportException(
                        $"COTP CC frame too short: {frameLen} bytes.");

                // Wait for the complete frame — separate deadline so a slow header arrival
                // does not eat the body's waiting budget (the timeout-reuse bug).
                var frameDeadline = Environment.TickCount64 + timeoutMs;
                while (Environment.TickCount64 < frameDeadline)
                {
                    if (!_inner.IsOpen)
                        throw new TransportException(
                            "Inner transport closed while waiting for the COTP Connection Confirm.");
                    if (_inner.GetNumBytesAvailable() >= frameLen) break;
                    System.Threading.Thread.Sleep(2);
                }

                if (_inner.GetNumBytesAvailable() < frameLen)
                    throw new TransportException(
                        $"Incomplete COTP CC frame: have {_inner.GetNumBytesAvailable()}, need {frameLen}.");

                var frame = _inner.Read(frameLen);

                // The TPKT version byte is 0x03; anything else means the stream is not
                // ISO-on-TCP (or the peer is sending garbage) and the frame must not be
                // interpreted as a Confirm.
                if (frame[0] != 0x03)
                {
                    throw new TransportException(
                        $"Unexpected TPKT version 0x{frame[0]:X2} in the Connection Confirm (expected 0x03).");
                }

                var cc = TpktFrame.Unwrap(frame);

                // Validate the CC: LI ≥ 6 and PDU type 0xD0. Java parity: the reference
                // checks only the PDU type (some ISO-on-TCP peers reply with a DST-REF
                // that does not echo our SRC-REF). The CR always sends SRC-REF 1, so a
                // stale Confirm from a timed-out first attempt is indistinguishable from
                // a fresh one — accepted, matching Java's behaviour.
                if (cc.Length < 7)
                {
                    throw new TransportException("COTP Connection Confirm is shorter than its fixed header.");
                }
                if (cc[0] < 0x06)
                {
                    throw new TransportException(
                        $"Invalid COTP Connection Confirm length indicator 0x{cc[0]:X2}; expected at least 0x06.");
                }
                var ccLength = cc[0] + 1;
                if (ccLength > cc.Length)
                {
                    throw new TransportException(
                        $"COTP Connection Confirm length indicator requires {ccLength} bytes, but only {cc.Length} are available.");
                }
                if (cc[1] != 0xD0)
                {
                    throw new TransportException(
                        $"Expected COTP CC (0xD0) but received PDU type 0x{cc[1]:X2}.");
                }

                // Negotiate the TPDU size from the CC's TPDU-size parameter (0xC0), which
                // may confirm something smaller than the 1024 bytes requested.
                _tpduSize = Math.Min(RequestedTpduSize, ParseTpduSize(cc, ccLength));

                WriteHex(frame, $"RECV  COTP Connection Confirm (CC) — negotiated TPDU size = {_tpduSize} bytes");

                TaskCompletionSource<bool>? completedHandshake;
                lock (_handshakeLock)
                {
                    // A Close() that ran after the Confirm was parsed has already failed this
                    // handshake and closed the inner transport. Publishing success now would
                    // bring a closed session back for this caller while every waiter was told
                    // the handshake failed, so only the owner of the still-current handshake
                    // may complete it.
                    if (!_handshakeInProgress || !ReferenceEquals(_handshakeCompletion, ownedHandshake))
                    {
                        throw new TransportException(
                            "COTP transport was closed while the handshake was in progress.");
                    }

                    _handshakeDone = true;
                    _handshakeInProgress = false;
                    completedHandshake = _handshakeCompletion;
                    _handshakeCompletion = null;
                }
                completedHandshake?.TrySetResult(true);
            }
            catch (Exception ex)
            {
                // A failed handshake leaves the instance open for a retry: only clear the
                // in-progress flag, never mark the handshake done. A Close() that already
                // cleared this handshake (and failed its waiters) owns that state now, so it
                // is left alone.
                TaskCompletionSource<bool>? failedHandshake = null;
                lock (_handshakeLock)
                {
                    if (ReferenceEquals(_handshakeCompletion, ownedHandshake))
                    {
                        _handshakeInProgress = false;
                        failedHandshake = _handshakeCompletion;
                        _handshakeCompletion = null;
                    }
                }
                failedHandshake?.TrySetException(ex);
                throw;
            }
        }

        /// <summary>
        /// Reads the TPDU-size parameter (code 0xC0) from a Connection Confirm. The
        /// value is the ISO 8073 TPDU-size index; absent a parameter, keep the 1024
        /// bytes the Connection Request asked for.
        /// </summary>
        private static int ParseTpduSize(byte[] cc, int ccLength)
        {
            // Parameter list starts after LI(1) + PDU type(1) + DST-REF(2) + SRC-REF(2)
            // + class(1); each entry is [code][length][value].
            var offset = 7;
            while (offset + 1 < ccLength)
            {
                var code = cc[offset];
                var len = cc[offset + 1];
                if (offset + 2 + len > ccLength) break;

                if (code == 0xC0 && len == 1)
                {
                    return cc[offset + 2] switch
                    {
                        0x06 => 64,
                        0x07 => 128,
                        0x08 => 256,
                        0x09 => 512,
                        0x0A => 1024,
                        0x0B => 2048,
                        0x0C => 4096,
                        0x0D => 8192,
                        _ => RequestedTpduSize
                    };
                }

                offset += 2 + len;
            }
            return RequestedTpduSize;
        }

        // ── Diagnostic helpers ──────────────────────────────────

        private void WriteHex(ReadOnlySpan<byte> bytes, string label)
        {
            var diag = DiagnosticOutput;
            if (diag == null) return;
            diag.WriteLine($"┌ {label} ({bytes.Length} bytes)");
            for (int i = 0; i < bytes.Length; i += 16)
            {
                var end = Math.Min(i + 16, bytes.Length);
                diag.Write("│ ");
                for (int j = i; j < end; j++)
                    diag.Write($"{bytes[j]:X2} ");
                diag.WriteLine();
            }
            diag.WriteLine("└");
            diag.Flush();
        }

        private void WriteHex(string label, params byte[][] parts)
        {
            var diag = DiagnosticOutput;
            if (diag == null) return;
            int total = 0;
            foreach (var p in parts) total += p.Length;
            diag.WriteLine($"┌ {label} ({total} bytes, {parts.Length} parts)");
            for (int p = 0; p < parts.Length; p++)
            {
                if (p > 0) diag.WriteLine("│ ---");
                var bytes = parts[p];
                for (int i = 0; i < bytes.Length; i += 16)
                {
                    var end = Math.Min(i + 16, bytes.Length);
                    diag.Write("│ ");
                    for (int j = i; j < end; j++)
                        diag.Write($"{bytes[j]:X2} ");
                    diag.WriteLine();
                }
            }
            diag.WriteLine("└");
            diag.Flush();
        }

        // ── Data Transfer ────────────────────────────────────────

        // TPDU number for DT frames. Starts at 0 after the handshake
        // and increments for every DT frame sent (or fragment thereof),
        // wrapping at 128 per ISO 8073.
        private int _currentTpduNr;

        public override void Write(byte[] bytes)
        {
            if (bytes == null || bytes.Length == 0) return;

            // Fast path, lock-free.
            if (!_handshakeDone)
                throw new TransportException("COTP handshake has not been performed. Call Open() first.");

            // _writeLock keeps the fragments of one message contiguous. _handshakeLock is held
            // only to validate the session and read the negotiated TPDU size, never across the
            // socket write: a send blocked on a stalled peer would otherwise freeze Close() and
            // every polling read behind it.
            lock (_writeLock)
            {
                int tpduSize;
                lock (_handshakeLock)
                {
                    if (_peerTeardown != null)
                        throw new TransportException(_peerTeardown);

                    if (!_handshakeDone)
                        throw new TransportException("COTP handshake has not been performed. Call Open() first.");

                    tpduSize = _tpduSize;
                }

                // The COTP DT header is 3 bytes (LI + 0xF0 + TPDU-NR/EOT).
                // The 16-bit TPKT length field caps the frame at 65535;
                // the negotiated TPDU size caps the DT TPDU. The tighter wins.
                var maxPayload = Math.Min(TpktFrame.MaxPayloadSize - 3, tpduSize - 3);
                var offset = 0;

                while (offset < bytes.Length)
                {
                    var isLast = offset + maxPayload >= bytes.Length;
                    var chunk = new byte[Math.Min(maxPayload, bytes.Length - offset)];
                    Array.Copy(bytes, offset, chunk, 0, chunk.Length);

                    var nr = _currentTpduNr;
                    _currentTpduNr = (_currentTpduNr + 1) & 0x7F; // wrap at 128

                    var cotpHeader = new byte[3];
                    cotpHeader[0] = 0x02;        // LI
                    cotpHeader[1] = 0xF0;        // DT pdu type
                    cotpHeader[2] = (byte)(isLast ? (nr | 0x80) : nr); // EOT on last

                    var combined = new byte[3 + chunk.Length];
                    Array.Copy(cotpHeader, 0, combined, 0, 3);
                    Array.Copy(chunk, 0, combined, 3, chunk.Length);

                    _inner.Write(TpktFrame.Wrap(combined));

                    var label = isLast
                        ? $"SEND  COTP DT (TPDU-NR={nr}, EOT)"
                        : $"SEND  COTP DT (TPDU-NR={nr}, fragment)";
                    WriteHex(label, cotpHeader, chunk);

                    offset += chunk.Length;
                }
            }
        }

        public override void Close()
        {
            // Close the inner transport first, with no COTP lock held: closing the socket is the
            // only thing that releases a Write() blocked against a stalled peer, and
            // TcpTransportInstance.Close() takes no locks for the same reason. Doing it under
            // _handshakeLock would let a stalled writer hold this call hostage.
            _inner.Close();

            TaskCompletionSource<bool>? pendingHandshake;
            lock (_handshakeLock)
            {
                _leftover = Array.Empty<byte>();
                _bufferedFramePayloads = 0;

                // Clearing the flags keeps this object's state internally consistent: once
                // closed it no longer claims a live COTP session. Clearing
                // _handshakeInProgress also lets a handshake that is still polling notice
                // (via IsOpen) and fail fast instead of waiting out its timeout. Reopening
                // is not possible today because the inner transport is fixed at
                // construction, so no test covers that path — Open() would fail on the
                // closed inner transport first.
                _handshakeDone = false;
                _handshakeInProgress = false;
                pendingHandshake = _handshakeCompletion;
                _handshakeCompletion = null;
            }
            pendingHandshake?.TrySetException(
                new TransportException("COTP transport was closed while the handshake was in progress."));
        }

        // ── IAsyncTransportInstance ──────────────────────────────
        //
        // COTP itself has no read loop; these forward to the inner transport so
        // a driver over COTP (S7) can still learn why the socket dropped. The
        // data-listener side is unused by the S7 driver's polling model.

        public void RegisterDataListener(Action listener)
        {
            (_inner as IAsyncTransportInstance)?.RegisterDataListener(listener);
        }

        public void RemoveDataListener()
        {
            (_inner as IAsyncTransportInstance)?.RemoveDataListener();
        }

        public void RegisterDisconnectListener(Action<Exception?> listener)
        {
            (_inner as IAsyncTransportInstance)?.RegisterDisconnectListener(ex =>
            {
                // A peer COTP/TCP drop invalidates the handshake state too.
                _handshakeDone = false;
                listener(ex);
            });
        }

        public void RemoveDisconnectListener()
        {
            (_inner as IAsyncTransportInstance)?.RemoveDisconnectListener();
        }

        /// <summary>
        /// Consumes complete TPKT frames from the inner transport, strips the
        /// COTP DT header, and appends the S7 payload to the leftover buffer.
        /// Only frames whose full body has arrived are consumed; incomplete
        /// frames are left in the inner transport.
        /// </summary>
        private void DrainFramesIfAvailable()
        {
            // The peer has ended the session: whatever it sent after that is not ours to read.
            if (_peerTeardown != null) return;

            while (true)
            {
                var available = _inner.GetNumBytesAvailable();
                if (available < TpktFrame.HeaderSize) return;

                var header = _inner.PeekReadableBytes(TpktFrame.HeaderSize);
                var totalLength = ((header[2] << 8) | header[3]);

                if (header[0] != TpktFrame.Version || header[1] != 0)
                {
                    throw new TransportException(
                        $"Invalid TPKT header: version 0x{header[0]:X2}, reserved 0x{header[1]:X2}.");
                }

                if (totalLength < TpktFrame.HeaderSize)
                {
                    throw new TransportException(
                        $"Invalid TPKT frame length {totalLength}; it must include the {TpktFrame.HeaderSize}-byte header.");
                }

                // Frame not yet complete — wait for more data.
                if (available < totalLength) return;

                // Read full frame and unwrap.
                var frame = _inner.Read(totalLength);
                var payload = TpktFrame.Unwrap(frame);

                // A frame too short to carry LI + PDU type cannot be interpreted; consume
                // it and continue (some S7-compatible stacks emit empty TPKT keepalives).
                // Treating it as fatal would kill a request a valid frame may be following.
                if (payload.Length < 3) continue;

                var pduType = payload[1];

                // A late Connection Confirm can legitimately precede an S7 response;
                // consume and skip it.
                if (pduType == 0xD0)
                {
                    continue;
                }

                // A Disconnect Request (0x80), a Disconnect Confirm (0xC0) or a TPDU Error
                // (0x70) means the peer is tearing the session down — most commonly a CPU
                // stop or a connection rejection. A DR ends the connection by definition,
                // a DC is not valid on a class 0 connection at all, and after an ER the
                // connection is released too (RFC 905 6.22 recommends it for the receiver,
                // and the sender releases it itself when its error timer expires).
                // Swallowing it made every read stall for the full timeout and report a
                // misleading "no response" while IsConnected stayed true. So end the session
                // here: the transport is closed, which stops it claiming a live session and
                // makes a later Write() fail instead of going out on a connection the peer has
                // ended. Callers learn why once no payload the peer delivered first is left to
                // hand out; the frames queued behind this one are not read.
                if (pduType == 0x80 || pduType == 0xC0 || pduType == 0x70)
                {
                    var teardown = pduType switch
                    {
                        0x80 => "Disconnect Request",
                        0xC0 => "Disconnect Confirm",
                        _ => "TPDU Error"
                    };
                    EndSessionByPeer(
                        $"COTP peer disconnected the session ({teardown}, PDU type 0x{pduType:X2}).");
                    return;
                }

                // Anything else that is not a Data Transfer frame is a protocol violation.
                if (pduType != 0xF0)
                {
                    throw new TransportException(
                        $"Expected COTP DT frame (0xF0) but received PDU type 0x{pduType:X2}.");
                }

                if (payload[0] != 0x02)
                {
                    throw new TransportException(
                        $"Invalid COTP DT length indicator 0x{payload[0]:X2}; expected 0x02.");
                }

                // Diagnostic: log the DT frame before stripping its header.
                if (DiagnosticOutput != null)
                {
                    var dtHeader = new byte[3];
                    Array.Copy(payload, 0, dtHeader, 0, 3);
                    var s7Part = new byte[payload.Length - 3];
                    Array.Copy(payload, 3, s7Part, 0, s7Part.Length);
                    WriteHex("RECV  COTP Data Transfer (DT)", dtHeader, s7Part);
                }

                // Strip COTP DT header (3 bytes: LI + 0xF0 + TPDU-NR/EOT).
                var s7Payload = new byte[payload.Length - 3];
                Array.Copy(payload, 3, s7Payload, 0, s7Payload.Length);

                // Append to leftover buffer.
                var newLeftover = new byte[_leftover.Length + s7Payload.Length];
                Buffer.BlockCopy(_leftover, 0, newLeftover, 0, _leftover.Length);
                Buffer.BlockCopy(s7Payload, 0, newLeftover, _leftover.Length, s7Payload.Length);
                _leftover = newLeftover;
                _bufferedFramePayloads = 0; // consolidated into leftover
            }
        }

        /// <summary>
        /// Records that the peer ended the session and closes the inner transport. Unlike
        /// <see cref="Close"/> it keeps <c>_leftover</c>: a response the peer sent before ending
        /// the session is still the caller's to read. Must be called with
        /// <c>_handshakeLock</c> held.
        /// </summary>
        /// <remarks>
        /// <see cref="Close"/> closes the inner transport before taking the lock, so that a
        /// stalled Write() cannot hold it up. Here the lock is already held, which is safe for
        /// the same reason: closing the inner transport never waits for the writer.
        /// <c>TcpTransportInstance.Close()</c> takes no locks and does not wait for its read
        /// loop, and a stalled Write() holds only <c>_writeLock</c>, never <c>_handshakeLock</c>.
        /// </remarks>
        private void EndSessionByPeer(string reason)
        {
            _peerTeardown = reason;
            _inner.Close();
        }
    }
}
