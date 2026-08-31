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

        public ModbusConnection(ConnectionString connectionString, ITransportInstance transport)
            : base(connectionString, transport)
        {
            var unitId = connectionString.GetIntParameter("unit-identifier", 1);
            if (unitId < 0 || unitId > 255)
                throw new PlcConnectionException(
                    $"unit-identifier must be in [0,255], got {unitId}.");
            UnitId = (byte)unitId;
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
                        PlcResponseCode.InternalError);
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
                    codes[name] = PlcResponseCode.InternalError;
                }
                catch (Exception ex)
                {
                    Logger.LogWarning(ex,
                        "Modbus Write failed for tag '{TagName}'",
                        name);
                    codes[name] = PlcResponseCode.InternalError;
                }
            }

            return new DefaultPlcWriteResponse(request, codes);
        }

        // ── Byte-level I/O ─────────────────────────────────────

        private async Task<byte[]> SendAndReceive(
            byte[] pduData, CancellationToken cancellationToken = default)
        {
            var txId = (ushort)(Interlocked.Increment(ref _nextTransactionId) & 0xFFFF);

            // Build MBAP header.
            var pduLen = 1 + pduData.Length;
            var totalLen = 6 + 1 + pduData.Length;
            var frame = new byte[totalLen];
            frame[0] = (byte)(txId >> 8);
            frame[1] = (byte)(txId & 0xFF);
            frame[2] = 0; frame[3] = 0;
            frame[4] = (byte)(pduLen >> 8);
            frame[5] = (byte)(pduLen & 0xFF);
            frame[6] = UnitId;
            Array.Copy(pduData, 0, frame, 7, pduData.Length);

            TransportInstance.Write(frame);

            // Poll for response honouring cancellation.
            var deadline = Environment.TickCount64 + 5000;
            while (Environment.TickCount64 < deadline)
            {
                cancellationToken.ThrowIfCancellationRequested();

                var available = TransportInstance.GetNumBytesAvailable();
                if (available >= 8) break;
                await Task.Delay(2, cancellationToken).ConfigureAwait(false);
            }

            cancellationToken.ThrowIfCancellationRequested();

            var available2 = TransportInstance.GetNumBytesAvailable();
            if (available2 < 8)
                throw new TimeoutException("No Modbus response received within timeout.");

            var header = TransportInstance.PeekReadableBytes(6);
            var respLen = (ushort)((header[4] << 8) | header[5]);
            var totalRespLen = 6 + respLen;

            // Re-query availability — data may have arrived since the initial check.
            var availNow = TransportInstance.GetNumBytesAvailable();
            if (availNow < totalRespLen)
                throw new Exception($"Incomplete response: need {totalRespLen}, have {availNow}.");

            var frameBytes = TransportInstance.Read(totalRespLen);

            // Validate transaction ID matches the request.
            var respTxId = (ushort)((frameBytes[0] << 8) | frameBytes[1]);
            if (respTxId != txId)
                throw new Exception(
                    $"Transaction ID mismatch: sent {txId}, received {respTxId}.");

            var functionCode = frameBytes[7];

            // Check for error response.
            if ((functionCode & 0x80) != 0)
            {
                var errorCode = frameBytes.Length > 8 ? frameBytes[8] : (byte)0;
                throw new Exception(
                    $"Modbus error: function={functionCode & 0x7F:X2}, code={errorCode:X2}");
            }

            var dataLen = totalRespLen - 8;
            if (dataLen <= 0) return Array.Empty<byte>();

            var result = new byte[dataLen];
            Array.Copy(frameBytes, 8, result, 0, dataLen);
            return result;
        }
    }
}
