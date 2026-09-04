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
using System.Threading;
using System.Threading.Tasks;
using Microsoft.Extensions.Logging;
using org.apache.plc4net.api.metadata;
using org.apache.plc4net.api.value;
using org.apache.plc4net.drivers.modbus.messages;
using org.apache.plc4net.messages;
using org.apache.plc4net.model;
using org.apache.plc4net.spi.drivers;
using org.apache.plc4net.spi.drivers.functions;
using org.apache.plc4net.spi.drivers.messages;
using org.apache.plc4net.spi.drivers.messages.items;
using org.apache.plc4net.spi.transports;
using org.apache.plc4net.types;

namespace org.apache.plc4net.drivers.modbus
{
    /// <summary>
    /// A Modbus RTU connection over a serial transport.
    /// </summary>
    /// <remarks>
    /// Connection string format:
    /// <c>modbus-rtu:serial://COM1?unit-identifier=1&amp;baud-rate=19200&amp;parity=Even</c>
    /// </remarks>
    public class ModbusRtuConnection : ConnectionBase, PlcReader, PlcWriter
    {
        private byte _slaveAddress = 1;
        private readonly int _responseTimeoutMs;

        // One outstanding transaction per connection — serialise callers so two
        // Read/Write calls cannot interleave on the half-duplex bus.
        private readonly SemaphoreSlim _io = new SemaphoreSlim(1, 1);

        public ModbusRtuConnection(
            ConnectionString connectionString,
            ITransportInstance transport)
            : base(connectionString, transport)
        {
            var addr = connectionString.GetIntParameter("unit-identifier", 1);
            if (addr < 0 || addr > 247)
                throw new ArgumentOutOfRangeException(
                    nameof(connectionString),
                    $"Modbus RTU address must be 1–247, got {addr}.");
            _slaveAddress = (byte)addr;
            _responseTimeoutMs = connectionString.GetIntParameter("request-timeout", 1000);
        }

        public byte SlaveAddress => _slaveAddress;

        public override void Close()
        {
            try { base.Close(); }
            finally { _io.Dispose(); }
        }

        // ── IPlcConnection ──────────────────────────────────

        public override IPlcConnectionMetadata PlcConnectionMetadata => _metadata;
        private readonly DefaultPlcConnectionMetadata _metadata
            = new DefaultPlcConnectionMetadata
            { CanRead = true, CanWrite = true, CanSubscribe = false };

        public override IPlcTag Parse(string tagQuery) => ModbusTag.Parse(tagQuery);

        public override IPlcReadRequestBuilder ReadRequestBuilder
            => new DefaultPlcReadRequestBuilder(this, Parse);

        public override IPlcWriteRequestBuilder WriteRequestBuilder
            => new DefaultPlcWriteRequestBuilder(this, Parse);

        public override IPlcSubscriptionRequestBuilder SubscriptionRequestBuilder => null;
        public override IPlcUnsubscriptionRequestBuilder UnsubscriptionRequestBuilder => null;

        // ── PlcReader ───────────────────────────────────────

        public async Task<IPlcReadResponse> Read(
            DefaultPlcReadRequest request,
            CancellationToken cancellationToken = default)
        {
            var results = new Dictionary<string, PlcResponseItem<IPlcValue>>();

            foreach (var name in request.TagNames)
            {
                var tag = request.GetTagByName(name) as ModbusTag;
                if (tag == null)
                {
                    results[name] = new DefaultPlcTagErrorItem<IPlcValue>(
                        PlcResponseCode.InternalError);
                    continue;
                }

                try
                {
                    cancellationToken.ThrowIfCancellationRequested();
                    var pdu = BuildReadPdu(tag);
                    var response = await SendAndReceive(pdu, cancellationToken)
                        .ConfigureAwait(false);
                    var value = ParseReadResponse(response, tag);
                    results[name] = new DefaultPlcResponseItem<IPlcValue>(
                        PlcResponseCode.Ok, value);
                }
                catch (OperationCanceledException)
                {
                    throw;
                }
                catch (Exception ex)
                {
                    Logger.LogWarning(ex,
                        "Modbus RTU Read failed for tag '{TagName}'", name);
                    results[name] = new DefaultPlcTagErrorItem<IPlcValue>(
                        ModbusConnection.MapModbusException(ex));
                }
            }

            return new DefaultPlcReadResponse(request, results);
        }

        // ── PlcWriter ───────────────────────────────────────

        public async Task<IPlcWriteResponse> Write(
            DefaultPlcWriteRequest request,
            CancellationToken cancellationToken = default)
        {
            var codes = new Dictionary<string, PlcResponseCode>();

            foreach (var name in request.TagNames)
            {
                var tag = request.GetTagByName(name) as ModbusTag;
                if (tag == null)
                {
                    codes[name] = PlcResponseCode.InternalError;
                    continue;
                }

                try
                {
                    cancellationToken.ThrowIfCancellationRequested();
                    var raw = request.GetValue(name);
                    var pdu = BuildWritePdu(tag, raw);
                    await SendAndReceive(pdu, cancellationToken)
                        .ConfigureAwait(false);
                    codes[name] = PlcResponseCode.Ok;
                }
                catch (OperationCanceledException)
                {
                    throw;
                }
                catch (Exception ex)
                {
                    Logger.LogWarning(ex,
                        "Modbus RTU Write failed for tag '{TagName}'", name);
                    codes[name] = ModbusConnection.MapModbusException(ex);
                }
            }

            return new DefaultPlcWriteResponse(request, codes);
        }

        // ── RTU framing ─────────────────────────────────────

        private byte[] BuildRtuFrame(byte[] pdu)
        {
            var frame = new byte[1 + pdu.Length + 2];
            frame[0] = _slaveAddress;
            Array.Copy(pdu, 0, frame, 1, pdu.Length);
            var crc = ModbusCRC.Compute(frame, 0, frame.Length - 2);
            frame[frame.Length - 2] = crc[0];
            frame[frame.Length - 1] = crc[1];
            return frame;
        }

        private const int MaxRtuFrameLen = 256;

        private async Task<byte[]> SendAndReceive(byte[] pdu, CancellationToken ct)
        {
            await _io.WaitAsync(ct).ConfigureAwait(false);
            try
            {
                return await SendAndReceiveLocked(pdu, ct).ConfigureAwait(false);
            }
            catch (Exception ex) when (ex is not ModbusDriverException)
            {
                DrainReadBuffer();
                throw;
            }
            finally
            {
                _io.Release();
            }
        }

        private async Task<byte[]> SendAndReceiveLocked(byte[] pdu, CancellationToken ct)
        {
            var request = BuildRtuFrame(pdu);

            // A frame left in the buffer by an earlier timeout would be read as
            // if it answered this request.
            DrainReadBuffer();
            TransportInstance.Write(request);
            var deadline = Environment.TickCount64 + _responseTimeoutMs;

            // Some 2-wire RS-485 adapters echo the transmitter onto the
            // receiver, so the request bytes arrive ahead of the answer. Only a
            // read is safe to de-echo: a write's response is byte-identical to
            // its request, so a lone match there is the response, not an echo.
            var isRead = pdu.Length > 0 && pdu[0] is >= ModbusFunctionCodes.ReadCoils
                                                   and <= ModbusFunctionCodes.ReadInputRegisters;
            if (isRead)
            {
                var settle = Environment.TickCount64 + 40;
                while (Environment.TickCount64 < settle
                       && TransportInstance.GetNumBytesAvailable() < request.Length)
                {
                    ct.ThrowIfCancellationRequested();
                    await Task.Delay(2, ct).ConfigureAwait(false);
                }
                if (TransportInstance.GetNumBytesAvailable() >= request.Length
                    && TransportInstance.PeekReadableBytes(request.Length)
                        .AsSpan().SequenceEqual(request))
                {
                    TransportInstance.Read(request.Length);
                    Logger.LogDebug(
                        "Modbus RTU: discarded a {N}-byte transmitter echo.", request.Length);
                }
            }

            // Read the header (addr + function [+ byte count / exception code])
            // and work out the exact frame length from the function code.
            await WaitForBytes(3, deadline, ct).ConfigureAwait(false);
            var head = TransportInstance.PeekReadableBytes(3);
            if (head[0] != _slaveAddress)
                throw new ModbusDriverException(
                    $"Modbus RTU: response address {head[0]}, expected {_slaveAddress}.");

            var fc = head[1];
            var total = (fc & 0x80) != 0
                ? 5 // addr + fc + code + crc
                : fc switch
                {
                    ModbusFunctionCodes.ReadCoils or ModbusFunctionCodes.ReadDiscreteInputs
                        or ModbusFunctionCodes.ReadHoldingRegisters
                        or ModbusFunctionCodes.ReadInputRegisters
                        => 3 + head[2] + 2, // addr + fc + byteCount + data + crc
                    _ => 8,                 // write single/multiple echo: addr + fc + addr(2) + val(2) + crc
                };
            if (total is < 5 or > MaxRtuFrameLen)
                throw new ModbusDriverException($"Modbus RTU: implausible frame length {total}.");

            await WaitForBytes(total, deadline, ct).ConfigureAwait(false);
            var frame = TransportInstance.Read(total);

            if (!ModbusCRC.Validate(frame, frame.Length))
            {
                DrainReadBuffer(); // the length guess was wrong / the line is noisy
                throw new ModbusDriverException("Modbus RTU: response failed CRC validation.");
            }

            if ((frame[1] & 0x80) != 0)
            {
                var code = frame[2];
                var enumName = Enum.IsDefined(typeof(ModbusErrorCode), code)
                    ? ((ModbusErrorCode)code).ToString() : "unknown";
                throw new ModbusDriverException(
                    $"Modbus RTU exception: function 0x{frame[1] & 0x7F:X2}, "
                    + $"code 0x{code:X2} ({enumName}).", code);
            }

            var result = new byte[frame.Length - 3]; // strip addr(1) + crc(2)
            Array.Copy(frame, 1, result, 0, result.Length);
            return result;
        }

        private async Task WaitForBytes(int count, long deadline, CancellationToken ct)
        {
            while (TransportInstance.GetNumBytesAvailable() < count)
            {
                ct.ThrowIfCancellationRequested();
                if (Environment.TickCount64 >= deadline)
                {
                    throw new TimeoutException(
                        $"No complete Modbus RTU response within {_responseTimeoutMs} ms "
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

        // ── PDU helpers (shared logic with Modbus TCP) ──────

        private static byte[] BuildReadPdu(ModbusTag tag)
        {
            return tag.Type switch
            {
                ModbusTag.TagType.Coil
                    => ModbusPDU.BuildReadBitsRequest(
                        ModbusFunctionCodes.ReadCoils,
                        tag.Address, 1),
                ModbusTag.TagType.DiscreteInput
                    => ModbusPDU.BuildReadBitsRequest(
                        ModbusFunctionCodes.ReadDiscreteInputs,
                        tag.Address, 1),
                ModbusTag.TagType.HoldingRegister
                    => ModbusPDU.BuildReadRegistersRequest(
                        ModbusFunctionCodes.ReadHoldingRegisters,
                        tag.Address, 1),
                ModbusTag.TagType.InputRegister
                    => ModbusPDU.BuildReadRegistersRequest(
                        ModbusFunctionCodes.ReadInputRegisters,
                        tag.Address, 1),
                _ => throw new ModbusDriverException(
                    $"Unsupported tag type for read: {tag.Type}")
            };
        }

        private static byte[] BuildWritePdu(ModbusTag tag, object rawValue)
        {
            if (tag.Type == ModbusTag.TagType.Coil)
            {
                var v = Convert.ToBoolean(rawValue);
                return new byte[] {
                    ModbusFunctionCodes.WriteSingleCoil,
                    (byte)(tag.Address >> 8),
                    (byte)(tag.Address & 0xFF),
                    (byte)(v ? 0xFF : 0x00),
                    0x00
                };
            }

            if (tag.Type == ModbusTag.TagType.HoldingRegister)
            {
                var v = Convert.ToUInt16(rawValue);
                return new byte[] {
                    ModbusFunctionCodes.WriteSingleRegister,
                    (byte)(tag.Address >> 8),
                    (byte)(tag.Address & 0xFF),
                    (byte)(v >> 8),
                    (byte)(v & 0xFF)
                };
            }

            throw new ModbusDriverException(
                $"Modbus RTU write not supported for tag type {tag.Type}.");
        }

        private static IPlcValue ParseReadResponse(
            byte[] pdu, ModbusTag tag)
        {
            if (pdu.Length == 0)
                throw new ModbusDriverException(
                    "Empty Modbus RTU response PDU.");

            var funcCode = pdu[0];

            // Read Coils / Discrete Inputs: func + byte count + data
            if (funcCode is ModbusFunctionCodes.ReadCoils
                           or ModbusFunctionCodes.ReadDiscreteInputs)
            {
                if (pdu.Length >= 3)
                    return new org.apache.plc4net.spi.model.values.PlcBOOL(
                        (pdu[2] & 0x01) != 0);
            }

            // Read Holding / Input Registers: func + byte count + data (2 bytes)
            if (funcCode is ModbusFunctionCodes.ReadHoldingRegisters
                           or ModbusFunctionCodes.ReadInputRegisters)
            {
                if (pdu.Length >= 4)
                {
                    var val = (ushort)((pdu[2] << 8) | pdu[3]);
                    return new org.apache.plc4net.spi.model.values.PlcUINT(val);
                }
            }

            throw new ModbusDriverException(
                $"Unexpected Modbus function code in response: 0x{funcCode:X2}");
        }
    }
}
