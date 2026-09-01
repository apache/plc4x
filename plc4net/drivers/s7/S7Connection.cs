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

            S7MessageResponseData response;
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
        /// Parses one S7 message and asserts it is an error-free AckData response.
        /// </summary>
        private static S7MessageResponseData ParseResponse(byte[] bytes)
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

            if (message is not S7MessageResponseData response)
            {
                throw new S7DriverException(
                    $"S7 response is a {message.GetType().Name}, not an AckData message.");
            }
            if (response.ErrorClass != 0x00 || response.ErrorCode != 0x00)
            {
                throw new S7DriverException(
                    $"S7 request rejected: error 0x{response.ErrorClass:X2}{response.ErrorCode:X2}.");
            }
            return response;
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

            var headerLength = header[1] == S7Constants.AckData ? 12 : 10;
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
            try
            {
                var pduRef = (ushort)(Interlocked.Increment(ref _pduRef) & 0xFFFF);
                TransportInstance.Write(S7Constants.BuildReadRequest(pduRef, tags));

                var response = ParseResponse(
                    await ReadOneS7MessageAsync(cancellationToken).ConfigureAwait(false));

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
                TransportInstance.Write(S7Constants.BuildWriteRequest(pduRef, items));

                var response = ParseResponse(
                    await ReadOneS7MessageAsync(cancellationToken).ConfigureAwait(false));

                if (response.Payload is not S7PayloadWriteVarResponse payload)
                {
                    foreach (var name in tagNames)
                        codes[name] = PlcResponseCode.InternalError;
                    return new DefaultPlcWriteResponse(request, codes);
                }

                for (var i = 0; i < tagNames.Count; i++)
                {
                    var name = tagNames[i];
                    codes[name] = i < payload.Items.Count
                                  && payload.Items[i].ReturnCode == DataTransportErrorCode.OK
                        ? PlcResponseCode.Ok
                        : PlcResponseCode.InternalError;
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
        /// representation that S7 Write Var expects on the wire.
        /// </summary>
        private static byte[] EncodeWriteValue(object raw, S7Tag tag)
        {
            // If the builder already wrapped it in a PlcValue, decode from there.
            if (raw is IPlcValue pv)
            {
                if (tag.DataTypeSize == 1 && tag.BitOffset >= 0)
                {
                    bool b = pv.IsBool() ? pv.GetBool() :
                             pv.IsByte() ? pv.GetByte() != 0 : false;
                    return new[] { (byte)(b ? 0x01 : 0x00) };
                }
                if (pv.IsBool())   return new[] { (byte)(pv.GetBool() ? 1 : 0) };
                if (pv.IsByte())   return new[] { pv.GetByte() };
                if (pv.IsUshort()) { var v = pv.GetUshort(); return new[] { (byte)(v >> 8), (byte)v }; }
                if (pv.IsUint())   { var v = pv.GetUint();   return new[] { (byte)(v >> 24), (byte)(v >> 16), (byte)(v >> 8), (byte)v }; }
                if (pv.IsInt())    { var v = unchecked((uint)pv.GetInt()); return new[] { (byte)(v >> 24), (byte)(v >> 16), (byte)(v >> 8), (byte)v }; }
                if (pv.IsShort())  { var v = unchecked((ushort)pv.GetShort()); return new[] { (byte)(v >> 8), (byte)v }; }
            }

            // Raw .NET primitives (the common case through the builder API).
            switch (raw)
            {
                case bool b:      return new[] { (byte)(b ? 1 : 0) };
                case byte by:     return new[] { by };
                case ushort us:   return new[] { (byte)(us >> 8), (byte)us };
                case short s:     return new[] { (byte)(unchecked((ushort)s) >> 8), (byte)unchecked((ushort)s) };
                case uint ui:     return new[] { (byte)(ui >> 24), (byte)(ui >> 16), (byte)(ui >> 8), (byte)ui };
                case int i:       return new[] { (byte)(unchecked((uint)i) >> 24), (byte)(unchecked((uint)i) >> 16), (byte)(unchecked((uint)i) >> 8), (byte)unchecked((uint)i) };
                case float f:     { var bits = BitConverter.GetBytes(f);
                                    if (BitConverter.IsLittleEndian) Array.Reverse(bits);
                                    return bits; }
                case string str:  return System.Text.Encoding.ASCII.GetBytes(str);
            }

            throw new NotSupportedException(
                $"S7 Write: cannot encode value of type {raw.GetType().Name}.");
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
