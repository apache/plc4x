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
using System.Threading;
using System.Threading.Tasks;
using Microsoft.Extensions.Logging;
using org.apache.plc4net.api.metadata;
using org.apache.plc4net.api.value;
using org.apache.plc4net.exceptions;
using org.apache.plc4net.drivers.modbus.messages;
using org.apache.plc4net.messages;
using org.apache.plc4net.model;
using org.apache.plc4net.spi.drivers;
using org.apache.plc4net.spi.drivers.functions;
using org.apache.plc4net.spi.drivers.messages;
using org.apache.plc4net.spi.drivers.messages.items;
using org.apache.plc4net.spi.model.values;
using org.apache.plc4net.spi.transports;
using org.apache.plc4net.types;

namespace org.apache.plc4net.drivers.modbus
{
    public class ModbusConnection : ConnectionBase, PlcReader, PlcWriter
    {
        private int _nextTransactionId = 1;

        private readonly int _responseTimeoutMs;

        public ModbusConnection(ConnectionString connectionString, ITransportInstance transport)
            : base(connectionString, transport)
        {
            var unitId = connectionString.GetIntParameter("unit-identifier", 1);
            if (unitId < 0 || unitId > 255)
                throw new PlcConnectionException(
                    $"unit-identifier must be in [0,255], got {unitId}.");
            UnitId = (byte)unitId;
            _responseTimeoutMs = connectionString.GetIntParameter("request-timeout", 5000);

            if (transport is IAsyncTransportInstance async)
            {
                async.RegisterDisconnectListener(ex =>
                    Logger.LogWarning(ex, "Modbus TCP transport disconnected"));
            }
        }

        public byte UnitId { get; }

        public override IPlcConnectionMetadata PlcConnectionMetadata => _metadata;

        private readonly DefaultPlcConnectionMetadata _metadata
            = new DefaultPlcConnectionMetadata
            {
                CanRead = true,
                CanWrite = true,
                CanSubscribe = false
            };

        public override IPlcTag Parse(string tagQuery) => ModbusTag.Parse(tagQuery);

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
            var results = new Dictionary<string, PlcResponseItem<IPlcValue>>();

            cancellationToken.ThrowIfCancellationRequested();

            foreach (var name in request.TagNames)
            {
                var tag = request.GetTagByName(name) as ModbusTag;
                if (tag == null)
                {
                    results[name] = new DefaultPlcTagErrorItem<IPlcValue>(PlcResponseCode.NotFound);
                    continue;
                }

                try
                {
                    byte[] response;
                    switch (tag.Type)
                    {
                        case ModbusTag.TagType.Coil:
                            response = await SendAndReceive(
                                ModbusPDU.BuildReadBitsRequest(
                                    ModbusFunctionCodes.ReadCoils, tag.Address, 1),
                                cancellationToken).ConfigureAwait(false);
                            var bits = ModbusPDU.ParseReadBitsResponse(response, 1);
                            results[name] = new DefaultPlcResponseItem<IPlcValue>(
                                PlcResponseCode.Ok, new PlcBOOL(bits[0]));
                            break;

                        case ModbusTag.TagType.HoldingRegister:
                            response = await SendAndReceive(
                                ModbusPDU.BuildReadRegistersRequest(
                                    ModbusFunctionCodes.ReadHoldingRegisters, tag.Address, 1),
                                cancellationToken).ConfigureAwait(false);
                            var regs = ModbusPDU.ParseReadRegistersResponse(response, 1);
                            results[name] = new DefaultPlcResponseItem<IPlcValue>(
                                PlcResponseCode.Ok, new PlcUINT(regs[0]));
                            break;

                        case ModbusTag.TagType.DiscreteInput:
                            response = await SendAndReceive(
                                ModbusPDU.BuildReadBitsRequest(
                                    ModbusFunctionCodes.ReadDiscreteInputs, tag.Address, 1),
                                cancellationToken).ConfigureAwait(false);
                            var discreteBits = ModbusPDU.ParseReadBitsResponse(response, 1);
                            results[name] = new DefaultPlcResponseItem<IPlcValue>(
                                PlcResponseCode.Ok, new PlcBOOL(discreteBits[0]));
                            break;

                        case ModbusTag.TagType.InputRegister:
                            response = await SendAndReceive(
                                ModbusPDU.BuildReadRegistersRequest(
                                    ModbusFunctionCodes.ReadInputRegisters, tag.Address, 1),
                                cancellationToken).ConfigureAwait(false);
                            var inputRegs = ModbusPDU.ParseReadRegistersResponse(response, 1);
                            results[name] = new DefaultPlcResponseItem<IPlcValue>(
                                PlcResponseCode.Ok, new PlcUINT(inputRegs[0]));
                            break;

                        default:
                            results[name] = new DefaultPlcTagErrorItem<IPlcValue>(
                                PlcResponseCode.AccessDenied);
                            break;
                    }
                }
                catch (OperationCanceledException)
                {
                    throw; // Propagate cancellation — do not mask as per-tag error.
                }
                catch (Exception ex)
                {
                    Logger.LogWarning(ex,
                        "Modbus Read failed for tag '{TagName}'",
                        name);
                    results[name] = new DefaultPlcTagErrorItem<IPlcValue>(
                        MapModbusException(ex));
                }
            }

            return new DefaultPlcReadResponse(request, results);
        }

        // ── PlcWriter ──────────────────────────────────────────

        public async Task<IPlcWriteResponse> Write(
            DefaultPlcWriteRequest request, CancellationToken cancellationToken = default)
        {
            var codes = new Dictionary<string, PlcResponseCode>();

            foreach (var name in request.TagNames)
            {
                var tag = (ModbusTag)request.GetTagByName(name);
                var value = request.GetValue(name);
                if (tag == null || value == null)
                {
                    codes[name] = PlcResponseCode.NotFound;
                    continue;
                }

                try
                {
                    switch (tag.Type)
                    {
                        case ModbusTag.TagType.Coil:
                            var boolVal = value is bool b ? b : Convert.ToBoolean(value);
                            await SendAndReceive(
                                ModbusPDU.BuildWriteSingleCoilRequest(tag.Address, boolVal),
                                cancellationToken).ConfigureAwait(false);
                            codes[name] = PlcResponseCode.Ok;
                            break;

                        case ModbusTag.TagType.HoldingRegister:
                            var ushortVal = value is ushort u ? u :
                                value is short s ? checked((ushort)s) : Convert.ToUInt16(value, CultureInfo.InvariantCulture);
                            await SendAndReceive(
                                ModbusPDU.BuildWriteSingleRegisterRequest(tag.Address, ushortVal),
                                cancellationToken).ConfigureAwait(false);
                            codes[name] = PlcResponseCode.Ok;
                            break;

                        default:
                            codes[name] = PlcResponseCode.AccessDenied;
                            break;
                    }
                }
                catch (OperationCanceledException)
                {
                    throw; // Propagate cancellation — do not mask as a per-tag error.
                }
                catch (Exception ex)
                {
                    Logger.LogWarning(ex,
                        "Modbus Write failed for tag '{TagName}'",
                        name);
                    codes[name] = MapModbusException(ex);
                }
            }

            return new DefaultPlcWriteResponse(request, codes);
        }

        // ── Byte-level I/O ─────────────────────────────────────

        // Unit id + a 252-byte PDU. A length field outside [2, 253] is garbled.
        private const int MaxModbusResponseLen = 253;

        // Modbus is strict request/response with one outstanding transaction per
        // connection. The gate serialises callers so two concurrent Read/Write
        // calls cannot interleave on the wire and swallow each other's frames.
        private readonly SemaphoreSlim _io = new SemaphoreSlim(1, 1);

        public override void Close()
        {
            try { base.Close(); }
            finally { _io.Dispose(); }
        }

        private async Task<byte[]> SendAndReceive(
            byte[] pduData, CancellationToken cancellationToken = default)
        {
            await _io.WaitAsync(cancellationToken).ConfigureAwait(false);
            try
            {
                return await SendAndReceiveLocked(pduData, cancellationToken).ConfigureAwait(false);
            }
            catch (Exception ex) when (ex is not ModbusDriverException)
            {
                // A timeout or a torn read may have stranded a partial frame;
                // do not let the next call mis-parse it. A ModbusDriverException
                // means the response was read in full — leave the buffer alone.
                DrainReadBuffer();
                throw;
            }
            finally
            {
                _io.Release();
            }
        }

        private async Task<byte[]> SendAndReceiveLocked(
            byte[] pduData, CancellationToken cancellationToken)
        {
            var txId = (ushort)(Interlocked.Increment(ref _nextTransactionId) & 0xFFFF);

            var pduLen = 1 + pduData.Length;
            var frame = new byte[6 + pduLen];
            frame[0] = (byte)(txId >> 8);
            frame[1] = (byte)(txId & 0xFF);
            frame[2] = 0; frame[3] = 0;
            frame[4] = (byte)(pduLen >> 8);
            frame[5] = (byte)(pduLen & 0xFF);
            frame[6] = UnitId;
            Array.Copy(pduData, 0, frame, 7, pduData.Length);

            TransportInstance.Write(frame);
            var deadline = Environment.TickCount64 + _responseTimeoutMs;

            // Loop so a stale response left by an earlier timed-out request is
            // skipped (wrong transaction id) rather than mis-parsed as ours, and
            // a garbled length field resyncs by dropping a byte.
            while (true)
            {
                await WaitForBytes(8, deadline, cancellationToken).ConfigureAwait(false);

                var header = TransportInstance.PeekReadableBytes(6);
                var respLen = (header[4] << 8) | header[5];
                if (respLen < 2 || respLen > MaxModbusResponseLen)
                {
                    TransportInstance.Read(1); // resync
                    continue;
                }

                var totalRespLen = 6 + respLen;
                // The frame may still be arriving (a TCP↔serial gateway forwards
                // it byte by byte); wait for all of it, not just the header.
                await WaitForBytes(totalRespLen, deadline, cancellationToken).ConfigureAwait(false);
                var frameBytes = TransportInstance.Read(totalRespLen);

                var respTxId = (ushort)((frameBytes[0] << 8) | frameBytes[1]);
                if (respTxId != txId)
                {
                    Logger.LogDebug(
                        "Modbus: discarded a stale response (tx {Got}, want {Want})",
                        respTxId, txId);
                    continue;
                }

                var functionCode = frameBytes[7];
                if ((functionCode & 0x80) != 0)
                {
                    var errorCode = frameBytes.Length > 8 ? frameBytes[8] : (byte)0;
                    var name = Enum.IsDefined(typeof(ModbusErrorCode), errorCode)
                        ? ((ModbusErrorCode)errorCode).ToString() : "unknown";
                    throw new ModbusDriverException(
                        $"Modbus exception: function 0x{functionCode & 0x7F:X2}, "
                        + $"code 0x{errorCode:X2} ({name}).", errorCode);
                }

                var dataLen = totalRespLen - 8;
                if (dataLen <= 0) return Array.Empty<byte>();
                var result = new byte[dataLen];
                Array.Copy(frameBytes, 8, result, 0, dataLen);
                return result;
            }
        }

        private async Task WaitForBytes(int count, long deadline, CancellationToken ct)
        {
            while (TransportInstance.GetNumBytesAvailable() < count)
            {
                ct.ThrowIfCancellationRequested();
                if (Environment.TickCount64 >= deadline)
                {
                    throw new TimeoutException(
                        $"No complete Modbus response within {_responseTimeoutMs} ms "
                        + $"({TransportInstance.GetNumBytesAvailable()} of {count} bytes).");
                }
                await Task.Delay(2, ct).ConfigureAwait(false);
            }
        }

        private void DrainReadBuffer()
        {
            var n = TransportInstance.GetNumBytesAvailable();
            if (n > 0)
            {
                TransportInstance.Read(n);
            }
        }

        /// <summary>
        /// Maps a failure from <see cref="SendAndReceive"/> to a response code,
        /// mirroring plc4j's ModbusTcpConnection.getErrorCode().
        /// </summary>
        internal static PlcResponseCode MapModbusException(Exception ex)
        {
            switch (ex)
            {
                case TimeoutException:
                    return PlcResponseCode.RequestTimeout;
                case ModbusDriverException m when m.ExceptionCode != 0:
                    return m.ExceptionCode switch
                    {
                        (byte)ModbusErrorCode.IllegalFunction => PlcResponseCode.Unsupported,
                        (byte)ModbusErrorCode.IllegalDataAddress => PlcResponseCode.InvalidAddress,
                        (byte)ModbusErrorCode.IllegalDataValue => PlcResponseCode.InvalidDatatype,
                        (byte)ModbusErrorCode.SlaveDeviceBusy => PlcResponseCode.RequestTimeout,
                        _ => PlcResponseCode.InternalError,
                    };
                default:
                    return PlcResponseCode.InternalError;
            }
        }
    }
}
