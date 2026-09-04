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
using System.Collections.Generic;
using System.Globalization;
using System.Text;
using System.Threading;
using System.Threading.Tasks;
using Microsoft.Extensions.Logging;
using org.apache.plc4net.api.metadata;
using org.apache.plc4net.api.value;
using org.apache.plc4net.drivers.s7.messages;
using org.apache.plc4net.drivers.s7.readwrite.model;
using org.apache.plc4net.exceptions;
using org.apache.plc4net.spi.generation;
using org.apache.plc4net.messages;
using org.apache.plc4net.model;
using org.apache.plc4net.spi.drivers;
using org.apache.plc4net.spi.drivers.functions;
using org.apache.plc4net.spi.drivers.messages;
using org.apache.plc4net.spi.drivers.messages.items;
using org.apache.plc4net.spi.model.values;
using org.apache.plc4net.spi.transports;
using org.apache.plc4net.transports.cotp;
using org.apache.plc4net.types;

namespace org.apache.plc4net.drivers.s7
{
    /// <summary>
    /// An S7 connection using the S7comm protocol over COTP/TPKT.
    /// </summary>
    public class S7Connection : ConnectionBase, PlcReader, PlcWriter
    {
        // TSAP encoding mirrors the Java SPI3 reference (S7TsapIdEncoder):
        //   encodeS7TsapId(deviceGroup, rack, slot) = (deviceGroup << 8) | (rack << 4) | (slot & 0x0F)
        // Each of rack and slot occupies 4 bits of the low byte, so values outside 0..15
        // cannot be represented and are rejected rather than silently truncated.
        // Device groups (Java DeviceGroup): PG_OR_PC = 0x01, OS = 0x02, OTHERS = 0x03.
        private const int DefaultLocalDeviceGroup = 0x03;   // OTHERS, matching Java's local default
        private const int DefaultRemoteDeviceGroup = 0x01;  // PG_OR_PC, matching Java's remote default
        private const int MaxRack = 0x0F;
        private const int MaxSlot = 0x0F;
        private const int DefaultRequestTimeoutMs = 5000;

        private readonly int _requestTimeoutMs;
        private int _pduRef = 1;
        private readonly int _remoteRack;
        private readonly int _remoteSlot;
        private readonly int _localRack;
        private readonly int _localSlot;
        private readonly int _localDeviceGroup;
        private readonly int _remoteDeviceGroup;
        private readonly int _explicitRemoteTsap;
        private readonly int _explicitLocalTsap;

        /// <summary>PDU length negotiated in Setup Communication (bytes). Set in OnConnect.</summary>
        public int NegotiatedPduLength { get; private set; }

        public S7Connection(ConnectionString connectionString, ITransportInstance transport)
            : base(connectionString, transport)
        {
            // Record configuration only. The COTP handshake happens in OnConnect(), which
            // DriverBase invokes after construction, so this constructor performs no I/O.
            // Rack/slot defaults match Java S7CotpTransportConfiguration: local 1/1,
            // remote 0/0.
            _remoteRack = ParseRackSlotParameter(connectionString, "remote-rack", 0);
            _remoteSlot = ParseRackSlotParameter(connectionString, "remote-slot", 0);
            _localRack = ParseRackSlotParameter(connectionString, "local-rack", 1);
            _localSlot = ParseRackSlotParameter(connectionString, "local-slot", 1);

            // Device group for the derived TSAP, matching Java's remote-device-group /
            // local-device-group. S7-1200 / S7-1500 often need remote-device-group=OTHERS
            // (0x03) rather than the S7-300/400 default PG_OR_PC (0x01); if that still
            // fails, configure an explicit S7 connection in TIA and pass remote-tsap.
            _localDeviceGroup = ParseDeviceGroupParameter(
                connectionString, "local-device-group", DefaultLocalDeviceGroup);
            _remoteDeviceGroup = ParseDeviceGroupParameter(
                connectionString, "remote-device-group", DefaultRemoteDeviceGroup);

            // Explicit TSAP overrides, honoured the way Java's getLocalTsap()/
            // getRemoteTsap() do: a non-zero value replaces the rack/slot derivation,
            // for CPUs/CPs whose TSAP is not rack/slot-expressible.
            _explicitRemoteTsap = ParseHexParameter(connectionString, "remote-tsap", 0);
            _explicitLocalTsap = ParseHexParameter(connectionString, "local-tsap", 0);

            _requestTimeoutMs = connectionString.GetIntParameter("request-timeout", DefaultRequestTimeoutMs);

            ValidateRackSlot("remote-rack", _remoteRack, MaxRack);
            ValidateRackSlot("remote-slot", _remoteSlot, MaxSlot);
            ValidateRackSlot("local-rack", _localRack, MaxRack);
            ValidateRackSlot("local-slot", _localSlot, MaxSlot);

            // Deliberate divergence from Java, documented here because it is visible:
            // Java's S7TsapIdEncoder accepts any int and silently masks (slot & 0x0F) or
            // lets rack overflow into the device-group byte, so remote-rack=16 in Java
            // silently means rack 0. We reject unrepresentable values instead, because
            // connecting to the wrong CPU beats an error message nobody asked for.
        }

        private static void ValidateRackSlot(string name, int value, int max)
        {
            if (value < 0 || value > max)
            {
                throw new PlcConnectionException(
                    $"{name} must be between 0 and {max}, but was {value}.");
            }
        }

        /// <summary>
        /// Reads a decimal parameter, failing on unparseable input instead of silently
        /// falling back to the default (Java's Integer.parseInt fails too; the shared
        /// GetIntParameter falls back, which is how a typo like 'remote-slot=1a' would
        /// silently become 0).
        /// </summary>
        private static int ParseRackSlotParameter(ConnectionString connectionString, string name, int defaultValue)
        {
            var raw = connectionString.GetParameter(name);
            if (raw == null) return defaultValue;
            if (!int.TryParse(raw, NumberStyles.Integer, CultureInfo.InvariantCulture, out var value))
            {
                throw new PlcConnectionException(
                    $"{name} must be a decimal integer, but was '{raw}'.");
            }
            return value;
        }

        /// <summary>
        /// Reads a device-group parameter: the Java names (<c>PG_OR_PC</c>, <c>OS</c>,
        /// <c>OTHERS</c>) or a raw value (<c>1</c>, <c>0x03</c>).
        /// </summary>
        private static int ParseDeviceGroupParameter(
            ConnectionString connectionString, string name, int defaultValue)
        {
            var raw = connectionString.GetParameter(name);
            if (raw == null) return defaultValue;
            switch (raw.Trim().ToUpperInvariant())
            {
                case "PG_OR_PC": case "PG": return 0x01;
                case "OS": return 0x02;
                case "OTHERS": case "OTHER": return 0x03;
            }
            var hex = raw.StartsWith("0x", StringComparison.OrdinalIgnoreCase);
            if (int.TryParse(hex ? raw.Substring(2) : raw,
                    hex ? NumberStyles.HexNumber : NumberStyles.Integer,
                    CultureInfo.InvariantCulture, out var value)
                && value >= 0 && value <= 0xFF)
            {
                return value;
            }
            throw new PlcConnectionException(
                $"{name} must be PG_OR_PC / OS / OTHERS or a byte value, but was '{raw}'.");
        }

        /// <summary>
        /// Reads a TSAP parameter, which users write in hex (0x0300). Accepts decimal
        /// too; anything unparseable fails loudly rather than silently becoming 0.
        /// </summary>
        private static int ParseHexParameter(ConnectionString connectionString, string name, int defaultValue)
        {
            var raw = connectionString.GetParameter(name);
            if (raw == null) return defaultValue;
            var hex = raw.StartsWith("0x", StringComparison.OrdinalIgnoreCase) ? raw.Substring(2) : raw;
            if (!int.TryParse(hex, NumberStyles.HexNumber, CultureInfo.InvariantCulture, out var value))
            {
                throw new PlcConnectionException(
                    $"{name} must be a hex TSAP like 0x0300, but was '{raw}'.");
            }
            return value;
        }

        /// <summary>
        /// Performs the COTP CR/CC handshake so that the S7 driver can exchange Data
        /// Transfer frames. Modbus needs no equivalent step, which is why this hook is
        /// virtual on the base rather than mandatory.
        /// </summary>
        protected override void OnConnect()
        {
            // S7 requires COTP: the CR/CC handshake is what establishes the ISO-on-TCP
            // session, and there is no framing path that works without it. The S7Driver
            // declares only the cotp transport, but a driver could still be handed a
            // plain transport through a future path — fail loudly rather than handing
            // back a connection that silently cannot talk.
            if (!(TransportInstance is CotpTransportInstance cotp))
            {
                throw new PlcConnectionException(
                    "The S7 driver requires the 'cotp' transport, but got '" +
                    TransportInstance.GetType().Name + "'.");
            }

            // An explicit remote-tsap/local-tsap parameter wins over the rack/slot
            // derivation, mirroring Java's getRemoteTsap()/getLocalTsap().
            var localTsap = _explicitLocalTsap != 0
                ? _explicitLocalTsap
                : EncodeTsap(_localDeviceGroup, _localRack, _localSlot);
            var remoteTsap = _explicitRemoteTsap != 0
                ? _explicitRemoteTsap
                : EncodeTsap(_remoteDeviceGroup, _remoteRack, _remoteSlot);

            cotp.Open(
                localTsapHi: (byte)((localTsap >> 8) & 0xFF),
                localTsapLo: (byte)(localTsap & 0xFF),
                remoteTsapHi: (byte)((remoteTsap >> 8) & 0xFF),
                remoteTsapLo: (byte)(remoteTsap & 0xFF));

            // S7 Setup Communication - negotiate PDU length and parallel-job limits.
            // A real CPU ignores Read / Write Var until this has completed.
            PerformSetupCommunication();
        }

        private void PerformSetupCommunication()
        {
            var pduRef = (ushort)(Interlocked.Increment(ref _pduRef) & 0xFFFF);
            TransportInstance.Write(S7Constants.BuildSetupCommunication(pduRef));

            S7Response response;
            try
            {
                response = ParseResponse(
                    ReadOneS7MessageAsync(CancellationToken.None).GetAwaiter().GetResult());
            }
            catch (TimeoutException e)
            {
                throw new PlcConnectionException(
                    "No S7 Setup Communication response - the CPU may not permit PUT/GET " +
                    "access, or the rack/slot / TSAP is wrong.", e);
            }
            catch (S7DriverException e)
            {
                throw new PlcConnectionException(e.Message, e);
            }

            if (response.HasError)
            {
                throw new PlcConnectionException(
                    $"S7 Setup Communication rejected: error 0x{response.ErrorClass:X2}{response.ErrorCode:X2}"
                    + DescribeS7Error(response.ErrorClass, response.ErrorCode) + ".");
            }

            if (response.Parameter is not S7ParameterSetupCommunication setup)
            {
                throw new PlcConnectionException(
                    "S7 Setup Communication: the response carried no setup parameter.");
            }

            NegotiatedPduLength = setup.PduLength;
            Logger.LogInformation(
                "S7 Setup Communication complete: negotiated PDU length {PduLength} bytes",
                NegotiatedPduLength);
        }

        /// <summary>
        /// The parts of an S7 acknowledgement the driver acts on. A CPU answers a
        /// Job with either an AckData (ROSCTR 0x03, carries the payload) or - when it
        /// refuses the request before looking at the items - a bare Ack (ROSCTR
        /// 0x02). Both put the same 2-byte error-info field first, so both land here.
        /// </summary>
        private readonly struct S7Response
        {
            public S7Response(ushort tpduReference, byte errorClass, byte errorCode,
                S7Parameter? parameter, S7Payload? payload)
            {
                TpduReference = tpduReference;
                ErrorClass = errorClass;
                ErrorCode = errorCode;
                Parameter = parameter;
                Payload = payload;
            }

            public ushort TpduReference { get; }
            public byte ErrorClass { get; }
            public byte ErrorCode { get; }
            public S7Parameter? Parameter { get; }
            public S7Payload? Payload { get; }
            public bool HasError => ErrorClass != 0x00 || ErrorCode != 0x00;
        }

        /// <summary>
        /// Parses one S7 message and reduces it to its error info, parameter and
        /// payload. Throws only when the bytes are not a parseable acknowledgement -
        /// a CPU-reported error is data the caller maps to a response code, not an
        /// exception.
        /// </summary>
        private static S7Response ParseResponse(byte[] bytes)
        {
            S7Message message;
            try
            {
                message = S7Message.StaticParse(new ReadBuffer(bytes));
            }
            catch (Exception e)
            {
                throw new S7DriverException("S7 response could not be parsed.", e);
            }

            return message switch
            {
                S7MessageResponseData d => new S7Response(
                    d.TpduReference, d.ErrorClass, d.ErrorCode, d.Parameter, d.Payload),
                S7MessageResponse a => new S7Response(
                    a.TpduReference, a.ErrorClass, a.ErrorCode, a.Parameter, a.Payload),
                _ => throw new S7DriverException(
                    $"S7 response is a {message.GetType().Name}, not an acknowledgement."),
            };
        }

        // Serialises the send/receive pair — the PLC4X contract is one
        // outstanding request per connection, and the raw-poll receive path is
        // not safe against two callers interleaving on the wire.
        private readonly SemaphoreSlim _io = new SemaphoreSlim(1, 1);

        public override void Close()
        {
            try { base.Close(); }
            finally { _io.Dispose(); }
        }

        /// <summary>
        /// Writes a request and reads its acknowledgement, correlating on the
        /// TPDU reference so a stale response left by an earlier timed-out call
        /// is skipped rather than returned. Drains the transport on any failure.
        /// </summary>
        private async Task<S7Response> SendAndReceiveS7(
            byte[] request, ushort pduRef, CancellationToken ct)
        {
            await _io.WaitAsync(ct).ConfigureAwait(false);
            try
            {
                TransportInstance.Write(request);
                var deadline = Environment.TickCount64 + _requestTimeoutMs;
                while (true)
                {
                    var response = ParseResponse(
                        await ReadOneS7MessageAsync(ct).ConfigureAwait(false));
                    if (response.TpduReference == pduRef)
                    {
                        return response;
                    }
                    Logger.LogDebug(
                        "S7: discarded a stale response (tpdu {Got}, want {Want})",
                        response.TpduReference, pduRef);
                    if (Environment.TickCount64 >= deadline)
                    {
                        throw new TimeoutException(
                            "S7: no response with the expected TPDU reference.");
                    }
                }
            }
            catch
            {
                DrainTransport(); // leave no partial frame for the next call
                throw;
            }
            finally
            {
                _io.Release();
            }
        }

        private void DrainTransport()
        {
            var n = TransportInstance.GetNumBytesAvailable();
            if (n > 0)
            {
                TransportInstance.Read(n);
            }
        }

        /// <summary>
        /// Whether a Read/Write of these tags fits one negotiated PDU. Estimates
        /// the S7 header (12) + parameter (2) + per-item (4 header + data rounded
        /// up to an even byte count) for the response, and the request side.
        /// </summary>
        private bool FitsOnePdu(List<S7Tag> tags)
        {
            var response = 14;
            var request = 12 + tags.Count * 12;
            foreach (var t in tags)
            {
                var dataBytes = t.BitOffset >= 0 ? 1 : (t.DataTypeSize + 1) & ~1;
                response += 4 + dataBytes;
            }
            return response <= NegotiatedPduLength && request <= NegotiatedPduLength;
        }

        /// <summary>
        /// Maps an S7 header-level error (errorClass / errorCode) to a response code,
        /// mirroring the Java driver's mapping. Class 0x81 code 0x04 is the CPU
        /// refusing PUT/GET access; an S7-300 reports the same refusal as 0x83/0x04;
        /// class 0x85 is a supply/access error.
        /// </summary>
        private static PlcResponseCode MapHeaderError(byte errorClass, byte errorCode)
        {
            if (errorClass == 0x81 && errorCode == 0x04) return PlcResponseCode.AccessDenied;
            if (errorClass == 0x83 && errorCode == 0x04) return PlcResponseCode.AccessDenied;
            if (errorClass == 0x85) return PlcResponseCode.AccessDenied;
            return PlcResponseCode.InternalError;
        }

        /// <summary>
        /// A human-readable hint for the S7 header errors a misconfigured CPU most
        /// often returns, appended to connection-level exception messages.
        /// </summary>
        private static string DescribeS7Error(byte errorClass, byte errorCode)
        {
            if (errorClass == 0x81 && errorCode == 0x04 || errorClass == 0x83 && errorCode == 0x04)
            {
                return " (the CPU refused access: tick \"Permit access with PUT/GET communication "
                     + "from remote partner\" in the CPU Protection & Security settings, and make "
                     + "sure the data block has \"Optimized block access\" unticked)";
            }
            if (errorClass == 0x85)
            {
                return " (access error: the address is out of range or the block does not exist)";
            }
            return string.Empty;
        }

        private static IPlcValue ParseValue(DataTransportSize transportSize, byte[] data)
        {
            if (transportSize == DataTransportSize.BIT)
            {
                return new PlcBOOL(data.Length > 0 && data[0] != 0x00);
            }
            return data.Length switch
            {
                1 => new PlcBYTE(data[0]),
                2 => new PlcUINT((ushort)((data[0] << 8) | data[1])),
                4 => new PlcUDINT((uint)((data[0] << 24) | (data[1] << 16) | (data[2] << 8) | data[3])),
                _ => new PlcRawByteArray(data),
            };
        }

        /// <summary>
        /// Reads exactly one S7 message off the transport, framed by the length in its
        /// own header (10-byte Job header, 12-byte AckData header), so a second buffered
        /// response is not swallowed.
        /// </summary>
        private async Task<byte[]> ReadOneS7MessageAsync(CancellationToken cancellationToken)
        {
            var deadline = Environment.TickCount64 + _requestTimeoutMs;
            await WaitForBytesAsync(12, deadline, cancellationToken).ConfigureAwait(false);

            var header = TransportInstance.PeekReadableBytes(12);
            if (header[0] != S7Constants.ProtocolId)
            {
                // Drain the junk so the next read is not stuck behind it.
                TransportInstance.Read(TransportInstance.GetNumBytesAvailable());
                throw new S7DriverException(
                    $"Not an S7 message: first byte 0x{header[0]:X2} (expected 0x32).");
            }

            // ROSCTR 0x02 (Ack) and 0x03 (AckData) both carry a 2-byte error-info
            // field (errorClass, errorCode) after the 10-byte common header; a Job
            // (0x01) and UserData (0x07) do not. Treating a bare Ack as a 10-byte
            // header truncates those two bytes and desyncs every following frame.
            var headerLength =
                header[1] == S7Constants.Ack || header[1] == S7Constants.AckData ? 12 : 10;
            var paramLength = (header[6] << 8) | header[7];
            var dataLength = (header[8] << 8) | header[9];
            var total = headerLength + paramLength + dataLength;

            await WaitForBytesAsync(total, deadline, cancellationToken).ConfigureAwait(false);
            return TransportInstance.Read(total);
        }

        private async Task WaitForBytesAsync(int count, long deadline, CancellationToken cancellationToken)
        {
            while (TransportInstance.GetNumBytesAvailable() < count)
            {
                if (!TransportInstance.IsOpen)
                {
                    throw new S7DriverException(
                        "Transport closed while waiting for an S7 response.");
                }
                if (Environment.TickCount64 >= deadline)
                {
                    throw new TimeoutException(
                        $"S7 response incomplete: {TransportInstance.GetNumBytesAvailable()}/{count} bytes " +
                        $"within {_requestTimeoutMs} ms.");
                }
                await Task.Delay(5, cancellationToken).ConfigureAwait(false);
            }
        }

        private static int EncodeTsap(int deviceGroup, int rack, int slot)
        {
            return (deviceGroup << 8) | (rack << 4) | (slot & 0x0F);
        }

        public override IPlcConnectionMetadata PlcConnectionMetadata => _metadata;

        private readonly DefaultPlcConnectionMetadata _metadata
            = new DefaultPlcConnectionMetadata
            {
                CanRead = true,
                CanWrite = true,
                CanSubscribe = false
            };

        public override IPlcTag Parse(string tagQuery)
        {
            return S7Tag.Parse(tagQuery);
        }

        public override IPlcReadRequestBuilder ReadRequestBuilder
            => new DefaultPlcReadRequestBuilder(this, Parse);

        public override IPlcWriteRequestBuilder WriteRequestBuilder
            => new DefaultPlcWriteRequestBuilder(this, Parse);
        public override IPlcSubscriptionRequestBuilder SubscriptionRequestBuilder => null;
        public override IPlcUnsubscriptionRequestBuilder UnsubscriptionRequestBuilder => null;

        // ── PlcReader ──────────────────────────────────────────

        public async Task<IPlcReadResponse> Read(
            DefaultPlcReadRequest request, CancellationToken cancellationToken = default)
        {
            // Build aligned lists: tagNames[i] corresponds to tags[i].
            var tagNames = new List<string>();
            var tags = new List<S7Tag>();
            foreach (var name in request.TagNames)
            {
                var tag = request.GetTagByName(name) as S7Tag;
                if (tag != null)
                {
                    tagNames.Add(name);
                    tags.Add(tag);
                }
            }

            var results = new Dictionary<string, PlcResponseItem<IPlcValue>>();

            // The driver does not split a request across PDUs yet; a read whose
            // response would exceed the negotiated size fails with a clear
            // message rather than a truncated / rejected frame.
            if (NegotiatedPduLength > 0 && !FitsOnePdu(tags))
            {
                foreach (var name in tagNames)
                {
                    results[name] = new DefaultPlcTagErrorItem<IPlcValue>(PlcResponseCode.InternalError);
                }
                Logger.LogWarning(
                    "S7 read of {Count} tag(s) would exceed the negotiated {Pdu}-byte PDU; "
                    + "split it into smaller requests.", tags.Count, NegotiatedPduLength);
                return new DefaultPlcReadResponse(request, results);
            }

            try
            {
                var pduRef = (ushort)(Interlocked.Increment(ref _pduRef) & 0xFFFF);
                var response = await SendAndReceiveS7(
                    S7Constants.BuildReadRequest(pduRef, tags), pduRef, cancellationToken)
                    .ConfigureAwait(false);

                if (response.HasError)
                {
                    var code = MapHeaderError(response.ErrorClass, response.ErrorCode);
                    Logger.LogWarning(
                        "S7 CPU refused the read: error 0x{ErrorClass:X2}{ErrorCode:X2}{Hint}",
                        response.ErrorClass, response.ErrorCode,
                        DescribeS7Error(response.ErrorClass, response.ErrorCode));
                    foreach (var name in tagNames)
                        results[name] = new DefaultPlcTagErrorItem<IPlcValue>(code);
                    return new DefaultPlcReadResponse(request, results);
                }

                if (response.Payload is not S7PayloadReadVarResponse payload)
                {
                    foreach (var name in tagNames)
                        results[name] = new DefaultPlcTagErrorItem<IPlcValue>(PlcResponseCode.InternalError);
                    return new DefaultPlcReadResponse(request, results);
                }

                for (var i = 0; i < tagNames.Count; i++)
                {
                    var name = tagNames[i];
                    if (i >= payload.Items.Count)
                    {
                        results[name] = new DefaultPlcTagErrorItem<IPlcValue>(PlcResponseCode.InternalError);
                        continue;
                    }
                    var item = payload.Items[i];
                    results[name] = item.ReturnCode == DataTransportErrorCode.OK
                        ? new DefaultPlcResponseItem<IPlcValue>(
                            PlcResponseCode.Ok, ParseValue(item.TransportSize, item.Data))
                        : new DefaultPlcTagErrorItem<IPlcValue>(MapReturnCode(item.ReturnCode));
                }

                return new DefaultPlcReadResponse(request, results);
            }
            catch (OperationCanceledException)
            {
                throw; // Propagate cancellation — do not mask it as a protocol error.
            }
            catch (Exception ex)
            {
                Logger.LogWarning(ex,
                    "S7 Read failed for {TagCount} tag(s) over {Transport}",
                    tagNames.Count,
                    TransportInstance.GetType().Name);
                foreach (var name in tagNames)
                {
                    if (!results.ContainsKey(name))
                        results[name] = new DefaultPlcTagErrorItem<IPlcValue>(
                            PlcResponseCode.InternalError);
                }
                return new DefaultPlcReadResponse(request, results);
            }
        }

        // ── PlcWriter ──────────────────────────────────────────

        public async Task<IPlcWriteResponse> Write(
            DefaultPlcWriteRequest request, CancellationToken cancellationToken = default)
        {
            var tagNames = new List<string>();
            var items = new List<(S7Tag tag, byte[] data)>();
            foreach (var name in request.TagNames)
            {
                var tag = request.GetTagByName(name) as S7Tag;
                if (tag == null) continue;
                var raw = request.GetValue(name);
                if (raw == null) continue;
                var data = EncodeWriteValue(raw, tag);
                tagNames.Add(name);
                items.Add((tag, data));
            }

            var codes = new Dictionary<string, PlcResponseCode>();
            try
            {
                var pduRef = (ushort)(Interlocked.Increment(ref _pduRef) & 0xFFFF);
                var response = await SendAndReceiveS7(
                    S7Constants.BuildWriteRequest(pduRef, items), pduRef, cancellationToken)
                    .ConfigureAwait(false);

                if (response.HasError)
                {
                    var code = MapHeaderError(response.ErrorClass, response.ErrorCode);
                    Logger.LogWarning(
                        "S7 CPU refused the write: error 0x{ErrorClass:X2}{ErrorCode:X2}{Hint}",
                        response.ErrorClass, response.ErrorCode,
                        DescribeS7Error(response.ErrorClass, response.ErrorCode));
                    foreach (var name in tagNames)
                        codes[name] = code;
                    return new DefaultPlcWriteResponse(request, codes);
                }

                if (response.Payload is not S7PayloadWriteVarResponse payload)
                {
                    foreach (var name in tagNames)
                        codes[name] = PlcResponseCode.InternalError;
                    return new DefaultPlcWriteResponse(request, codes);
                }

                for (var i = 0; i < tagNames.Count; i++)
                {
                    var name = tagNames[i];
                    if (i >= payload.Items.Count)
                    {
                        codes[name] = PlcResponseCode.InternalError;
                        continue;
                    }
                    var returnCode = payload.Items[i].ReturnCode;
                    codes[name] = returnCode == DataTransportErrorCode.OK
                        ? PlcResponseCode.Ok
                        : MapReturnCode(returnCode);
                }
            }
            catch (OperationCanceledException)
            {
                throw;
            }
            catch (Exception ex)
            {
                Logger.LogWarning(ex,
                    "S7 Write failed for {TagCount} tag(s) over {Transport}",
                    tagNames.Count,
                    TransportInstance.GetType().Name);
                foreach (var name in tagNames)
                {
                    if (!codes.ContainsKey(name))
                        codes[name] = PlcResponseCode.InternalError;
                }
            }

            return new DefaultPlcWriteResponse(request, codes);
        }

        /// <summary>
        /// Encodes a raw value (IPlcValue, or a primitive .NET type such as
        /// byte / ushort / int / bool / string) into the big-endian byte
        /// representation S7 Write Var expects — exactly <c>tag.DataTypeSize</c>
        /// bytes for a scalar, so the payload width always matches the address
        /// item the CPU is told to expect.
        /// </summary>
        internal static byte[] EncodeWriteValue(object raw, S7Tag tag)
        {
            // Bit tag: a single 0/1 byte.
            if (tag.BitOffset >= 0)
            {
                var bit = raw switch
                {
                    bool bo => bo,
                    IPlcValue pv => pv.IsBool() ? pv.GetBool() : pv.GetLong() != 0,
                    _ => System.Convert.ToInt64(raw, System.Globalization.CultureInfo.InvariantCulture) != 0,
                };
                return new[] { (byte)(bit ? 0x01 : 0x00) };
            }

            var width = tag.DataTypeSize; // 1, 2 or 4

            // REAL: only a 4-byte tag can hold one.
            var isFloat = raw is float || (raw is IPlcValue fv && fv.IsFloat());
            if (isFloat)
            {
                if (width != 4)
                {
                    throw new S7DriverException(
                        $"S7 Write: a REAL needs a 4-byte tag, this one is {width}.");
                }
                var f = raw is float ff ? ff : ((IPlcValue)raw).GetFloat();
                var bits = BitConverter.GetBytes(f);
                if (BitConverter.IsLittleEndian) System.Array.Reverse(bits);
                return bits;
            }

            if (raw is string str)
            {
                return System.Text.Encoding.ASCII.GetBytes(str);
            }

            // Integer: coerce to `width` big-endian bytes, range-checked so a
            // caller passing an int for a byte tag fails loudly instead of
            // putting 4 bytes on the wire for a 1-byte address item.
            long value = raw switch
            {
                IPlcValue pv => pv.IsLong() ? pv.GetLong() : pv.GetInt(),
                _ => System.Convert.ToInt64(raw, System.Globalization.CultureInfo.InvariantCulture),
            };
            long maxUnsigned = width switch { 1 => byte.MaxValue, 2 => ushort.MaxValue, _ => uint.MaxValue };
            long minSigned = width switch { 1 => sbyte.MinValue, 2 => short.MinValue, _ => int.MinValue };
            if (value > maxUnsigned || value < minSigned)
            {
                throw new S7DriverException(
                    $"S7 Write: value {value} does not fit the {width}-byte tag.");
            }
            var u = unchecked((ulong)value);
            var bytes = new byte[width];
            for (var i = 0; i < width; i++)
            {
                bytes[width - 1 - i] = (byte)(u >> (8 * i));
            }
            return bytes;
        }

        private static PlcResponseCode MapReturnCode(DataTransportErrorCode code) => code switch
        {
            DataTransportErrorCode.ACCESS_DENIED => PlcResponseCode.AccessDenied,
            DataTransportErrorCode.INVALID_ADDRESS => PlcResponseCode.InvalidAddress,
            DataTransportErrorCode.DATA_TYPE_NOT_SUPPORTED => PlcResponseCode.InvalidDatatype,
            DataTransportErrorCode.NOT_FOUND => PlcResponseCode.NotFound,
            _ => PlcResponseCode.InternalError,
        };
    }
}
