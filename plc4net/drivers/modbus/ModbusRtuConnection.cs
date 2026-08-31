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
        }

        public byte SlaveAddress => _slaveAddress;

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
                        PlcResponseCode.InternalError);
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
                    codes[name] = PlcResponseCode.InternalError;
                }
                catch (Exception ex)
                {
                    Logger.LogWarning(ex,
                        "Modbus RTU Write failed for tag '{TagName}'", name);
                    codes[name] = PlcResponseCode.InternalError;
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

        private async Task<byte[]> SendAndReceive(
            byte[] pdu, CancellationToken ct)
        {
            var frame = BuildRtuFrame(pdu);
            TransportInstance.Write(frame);

            // Wait for the response.  Modbus RTU is half-duplex;
            // the slave responds with its own address + function
            // code + data + CRC.
            var deadline = Environment.TickCount64 + 2000;
            while (Environment.TickCount64 < deadline)
            {
                var available = TransportInstance.GetNumBytesAvailable();
                if (available >= 4) break; // addr + func + 2-byte CRC minimum
                await Task.Delay(5, ct).ConfigureAwait(false);
            }

            ct.ThrowIfCancellationRequested();

            var avail = TransportInstance.GetNumBytesAvailable();
            if (avail < 4)
                throw new TimeoutException(
                    "No Modbus RTU response received within timeout.");

            var response = TransportInstance.Read(avail);

            if (response.Length < 4
                || !ModbusCRC.Validate(response, response.Length))
                throw new ModbusDriverException(
                    "Modbus RTU response failed CRC validation.");

            if (response[0] != _slaveAddress)
                throw new ModbusDriverException(
                    $"Modbus RTU address mismatch: expected {_slaveAddress}, got {response[0]}.");

            // Check for Modbus exception (function code with high bit set).
            if ((response[1] & 0x80) != 0)
                throw new ModbusDriverException(
                    $"Modbus RTU exception: function 0x{(response[1] & 0x7F):X2}, code 0x{response[2]:X2}.");

            // Return PDU without address and CRC.
            var pduLen = response.Length - 3; // addr(1) + crc(2)
            var result = new byte[pduLen];
            Array.Copy(response, 1, result, 0, pduLen);
            return result;
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
                if (pdu.Length >= 3)
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
