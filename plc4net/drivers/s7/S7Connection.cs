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
using org.apache.plc4net.exceptions;
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
        private const int LocalDeviceGroup = 0x03;   // OTHERS, matching Java's local default
        private const int RemoteDeviceGroup = 0x01;  // PG_OR_PC, matching Java's remote default
        private const int MaxRack = 0x0F;
        private const int MaxSlot = 0x0F;

        private int _pduRef = 1;
        private readonly int _remoteRack;
        private readonly int _remoteSlot;
        private readonly int _localRack;
        private readonly int _localSlot;
        private readonly int _explicitRemoteTsap;
        private readonly int _explicitLocalTsap;

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

            // Explicit TSAP overrides, honoured the way Java's getLocalTsap()/
            // getRemoteTsap() do: a non-zero value replaces the rack/slot derivation,
            // for CPUs/CPs whose TSAP is not rack/slot-expressible.
            _explicitRemoteTsap = ParseHexParameter(connectionString, "remote-tsap", 0);
            _explicitLocalTsap = ParseHexParameter(connectionString, "local-tsap", 0);

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
                : EncodeTsap(LocalDeviceGroup, _localRack, _localSlot);
            var remoteTsap = _explicitRemoteTsap != 0
                ? _explicitRemoteTsap
                : EncodeTsap(RemoteDeviceGroup, _remoteRack, _remoteSlot);

            cotp.Open(
                localTsapHi: (byte)((localTsap >> 8) & 0xFF),
                localTsapLo: (byte)(localTsap & 0xFF),
                remoteTsapHi: (byte)((remoteTsap >> 8) & 0xFF),
                remoteTsapLo: (byte)(remoteTsap & 0xFF));
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
                var s7Payload = S7Constants.BuildReadRequest(pduRef, tags);

                // Write and read through the transport.
                TransportInstance.Write(s7Payload);

                // Wait for response.
                var deadline = Environment.TickCount64 + 5000;
                while (Environment.TickCount64 < deadline)
                {
                    var available = TransportInstance.GetNumBytesAvailable();
                    if (available >= 12) break; // S7 header minimum
                    await Task.Delay(5, cancellationToken).ConfigureAwait(false);
                }

                cancellationToken.ThrowIfCancellationRequested();

                var avail = TransportInstance.GetNumBytesAvailable();
                if (avail < 12)
                throw new TimeoutException("No S7 response received.");

            // Read and parse the S7 response.
            var responseBytes = TransportInstance.Read(avail);

            var paramOffset = 10; // after S7 header

            // Function code
            var funcCode = responseBytes[paramOffset];
            if (funcCode != S7Constants.ReadVar)
            {
                foreach (var name in tagNames)
                    results[name] = new DefaultPlcTagErrorItem<IPlcValue>(PlcResponseCode.InternalError);
                return new DefaultPlcReadResponse(request, results);
            }

            var itemCount = responseBytes[paramOffset + 1];
            var dataOffset = 10 + ((responseBytes[6] << 8) | responseBytes[7]); // after param part

            var paramIdx = paramOffset + 2; // after funcCode + itemCount
            var dataIdx = dataOffset + 2; // after funcCode + itemCount in data part

            for (int i = 0; i < itemCount && i < tagNames.Count; i++)
            {
                // Bounds check: each parameter item is 12 bytes.
                if (paramIdx + 12 > responseBytes.Length) break;

                var name = tagNames[i];
                var returnCode = responseBytes[paramIdx]; // 0xFF = success

                if (returnCode == 0xFF)
                {
                    // Read data: returnCode(1) + transportSize(1) + length(2) + data
                    if (dataIdx + 4 <= responseBytes.Length)
                    {
                        var dataRetCode = responseBytes[dataIdx];
                        var dataLen = (responseBytes[dataIdx + 2] << 8) | responseBytes[dataIdx + 3];

                        if (dataRetCode == 0xFF && dataIdx + 4 + dataLen <= responseBytes.Length)
                        {
                            var value = ParseS7Value(responseBytes, dataIdx + 4, dataLen);
                            results[name] = new DefaultPlcResponseItem<IPlcValue>(
                                PlcResponseCode.Ok, value);
                        }
                        else
                        {
                            results[name] = new DefaultPlcTagErrorItem<IPlcValue>(
                                PlcResponseCode.InternalError);
                        }
                    }
                    else
                    {
                        // Truncated response — data item header missing.
                        results[name] = new DefaultPlcTagErrorItem<IPlcValue>(
                            PlcResponseCode.InternalError);
                    }
                }
                else
                {
                    results[name] = new DefaultPlcTagErrorItem<IPlcValue>(
                        PlcResponseCode.InternalError);
                }

                paramIdx += 12; // each parameter item is 12 bytes
                if (dataIdx + 4 <= responseBytes.Length)
                    dataIdx += 4 + ((responseBytes[dataIdx + 2] << 8) | responseBytes[dataIdx + 3]);
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
                var s7Payload = S7Constants.BuildWriteRequest(pduRef, items);

                TransportInstance.Write(s7Payload);

                // Wait for response.
                var deadline = Environment.TickCount64 + 5000;
                while (Environment.TickCount64 < deadline)
                {
                    var available = TransportInstance.GetNumBytesAvailable();
                    if (available >= 12) break;
                    await Task.Delay(5, cancellationToken).ConfigureAwait(false);
                }

                cancellationToken.ThrowIfCancellationRequested();

                var avail = TransportInstance.GetNumBytesAvailable();
                if (avail < 12)
                {
                    // No response — mark all tags as error.
                    foreach (var name in tagNames)
                        codes[name] = PlcResponseCode.InternalError;
                    return new DefaultPlcWriteResponse(request, codes);
                }

                var responseBytes = TransportInstance.Read(avail);

                var funcCode = responseBytes[10];
                if (funcCode != S7Constants.WriteVar)
                {
                    foreach (var name in tagNames)
                        codes[name] = PlcResponseCode.InternalError;
                    return new DefaultPlcWriteResponse(request, codes);
                }

                var itemCount = responseBytes[11];
                var dataOffset = 10 + ((responseBytes[6] << 8) | responseBytes[7]);
                var dataIdx = dataOffset + 2; // funcCode + itemCount in data part

                for (int i = 0; i < itemCount && i < tagNames.Count; i++)
                {
                    var name = tagNames[i];
                    if (dataIdx < responseBytes.Length)
                    {
                        var retCode = responseBytes[dataIdx];
                        codes[name] = retCode == 0xFF
                            ? PlcResponseCode.Ok
                            : PlcResponseCode.InternalError;
                        // Advance past: return code (1) + transport size (1) + length (2) + data
                        if (dataIdx + 4 <= responseBytes.Length)
                        {
                            var dataLen = (responseBytes[dataIdx + 2] << 8) | responseBytes[dataIdx + 3];
                            dataIdx += 4 + dataLen;
                        }
                        else dataIdx++;
                    }
                    else
                    {
                        codes[name] = PlcResponseCode.InternalError;
                    }
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

        private static IPlcValue ParseS7Value(byte[] data, int offset, int length)
        {
            if (length == 1)
                return new PlcBYTE(data[offset]);
            if (length == 2)
                return new PlcUINT((ushort)((data[offset] << 8) | data[offset + 1]));
            if (length >= 4)
                return new PlcUDINT((uint)((data[offset] << 24) | (data[offset + 1] << 16)
                    | (data[offset + 2] << 8) | data[offset + 3]));
            return new PlcBYTE(0);
        }
    }
}
