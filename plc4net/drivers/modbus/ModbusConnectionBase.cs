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
using System.Collections;
using System.Collections.Generic;
using System.Globalization;
using System.Threading;
using System.Threading.Tasks;
using Microsoft.Extensions.Logging;
using org.apache.plc4net.api.metadata;
using org.apache.plc4net.api.value;
using org.apache.plc4net.drivers.modbus.readwrite.model;
using org.apache.plc4net.exceptions;
using org.apache.plc4net.messages;
using org.apache.plc4net.model;
using org.apache.plc4net.spi.drivers;
using org.apache.plc4net.spi.drivers.functions;
using org.apache.plc4net.spi.drivers.messages;
using org.apache.plc4net.spi.drivers.messages.items;
using org.apache.plc4net.spi.generation;
using org.apache.plc4net.spi.model.values;
using org.apache.plc4net.spi.transports;
using org.apache.plc4net.types;

namespace org.apache.plc4net.drivers.modbus
{
    /// <summary>
    /// The Modbus logic that the TCP and RTU connections share: turning a tag into
    /// a request PDU, turning the response PDU into a value or a response code.
    /// A subclass adds the ADU framing and the exchange over its transport.
    /// </summary>
    /// <remarks>
    /// Every PDU, and every register value, is encoded and decoded by a type generated from
    /// <c>modbus.mspec</c>; nothing here handles bytes itself. Reads and writes handle one
    /// value per tag.
    /// </remarks>
    public abstract class ModbusConnectionBase : ConnectionBase, PlcReader, PlcWriter
    {
        private static readonly DefaultPlcConnectionMetadata Metadata = new DefaultPlcConnectionMetadata
        {
            CanRead = true,
            CanWrite = true,
            CanSubscribe = false
        };

        protected ModbusConnectionBase(ConnectionString connectionString, ITransportInstance transport)
            : base(connectionString, transport)
        {
        }

        /// <summary>
        /// Reads an integer connection-string parameter. The names are tried in order, so
        /// plc4j's current name comes first and an older spelling still works; a value that
        /// is present but not a number is an error rather than a silent default, because a
        /// mistyped unit identifier would address the wrong device.
        /// </summary>
        protected static int IntParameter(ConnectionString connectionString, int defaultValue, params string[] names)
        {
            foreach (var name in names)
            {
                var raw = connectionString.GetParameter(name);
                if (raw == null)
                {
                    continue;
                }
                if (!int.TryParse(raw, NumberStyles.Integer, CultureInfo.InvariantCulture, out var value))
                {
                    throw new PlcConnectionException(
                        $"The connection string parameter '{name}' must be a whole number, but was '{raw}'.");
                }
                return value;
            }
            return defaultValue;
        }

        /// <summary>The request timeout in milliseconds; plc4j's <c>request-timeout-ms</c>, 5000 unless stated.</summary>
        protected static int RequestTimeoutMs(ConnectionString connectionString)
        {
            var timeout = IntParameter(connectionString, 5000, "request-timeout-ms", "request-timeout");
            if (timeout < 1)
            {
                throw new PlcConnectionException($"request-timeout-ms must be positive, but was {timeout}.");
            }
            return timeout;
        }

        public override IPlcConnectionMetadata PlcConnectionMetadata => Metadata;

        public override IPlcTag Parse(string tagQuery) => ModbusTag.Parse(tagQuery);

        public override IPlcReadRequestBuilder? ReadRequestBuilder
            => new DefaultPlcReadRequestBuilder(this, Parse);

        public override IPlcWriteRequestBuilder? WriteRequestBuilder
            => new DefaultPlcWriteRequestBuilder(this, Parse);

        public override IPlcSubscriptionRequestBuilder? SubscriptionRequestBuilder => null;

        public override IPlcUnsubscriptionRequestBuilder? UnsubscriptionRequestBuilder => null;

        /// <summary>
        /// Sends one request PDU and returns the response PDU that answers it.
        /// Implementations enforce one outstanding request per connection and throw
        /// <see cref="TimeoutException"/> when nothing arrives within the request timeout.
        /// </summary>
        protected abstract Task<ModbusPDU> ExchangeAsync(
            ModbusPDU request, CancellationToken cancellationToken);

        // ── PlcReader ──────────────────────────────────────────

        public async Task<IPlcReadResponse> Read(
            DefaultPlcReadRequest request, CancellationToken cancellationToken = default)
        {
            var results = new Dictionary<string, PlcResponseItem<IPlcValue>>();
            var work = new List<(string Name, ModbusTag Tag)>();
            foreach (var name in request.TagNames)
            {
                if (request.GetTagByName(name) is ModbusTag tag)
                {
                    work.Add((name, tag));
                }
                else
                {
                    results[name] = new DefaultPlcTagErrorItem<IPlcValue>(PlcResponseCode.InvalidAddress);
                }
            }

            for (var index = 0; index < work.Count;)
            {
                cancellationToken.ThrowIfCancellationRequested();
                var block = PlanReadBlock(work, index);
                try
                {
                    if (block.Count == 1)
                    {
                        var item = work[index];
                        var value = await ReadTagAsync(item.Tag, cancellationToken).ConfigureAwait(false);
                        results[item.Name] = new DefaultPlcResponseItem<IPlcValue>(PlcResponseCode.Ok, value);
                    }
                    else
                    {
                        await ReadRegisterBlockAsync(work, index, block.Count, block.Registers,
                            results, cancellationToken).ConfigureAwait(false);
                    }
                }
                catch (OperationCanceledException)
                {
                    throw; // Cancellation is the caller's decision, not a per-tag error.
                }
                catch (Exception ex)
                {
                    Logger.LogWarning(ex, "Modbus read of {Count} tag(s) failed", block.Count);
                    for (var failed = 0; failed < block.Count; failed++)
                    {
                        results[work[index + failed].Name] =
                            new DefaultPlcTagErrorItem<IPlcValue>(MapException(ex));
                    }
                }
                index += block.Count;
            }

            return new DefaultPlcReadResponse(request, results);
        }

        private static (int Count, ushort Registers) PlanReadBlock(
            List<(string Name, ModbusTag Tag)> work, int start)
        {
            var first = work[start].Tag;
            if (first.IsBit || first.Type == ModbusTag.TagType.ExtendedRegister)
            {
                return (1, first.RegisterCount);
            }
            var count = 1;
            var registers = first.RegisterCount;
            var nextAddress = first.WireAddress + registers;
            while (start + count < work.Count)
            {
                var next = work[start + count].Tag;
                if (next.Type != first.Type || next.ByteOrder != first.ByteOrder
                    || next.WireAddress != nextAddress || registers + next.RegisterCount > 125)
                {
                    break;
                }
                registers += next.RegisterCount;
                nextAddress += next.RegisterCount;
                count++;
            }
            return (count, checked((ushort)registers));
        }

        private async Task ReadRegisterBlockAsync(
            List<(string Name, ModbusTag Tag)> work, int start, int count, ushort registers,
            Dictionary<string, PlcResponseItem<IPlcValue>> results, CancellationToken cancellationToken)
        {
            var first = work[start].Tag;
            byte[] data;
            if (first.Type == ModbusTag.TagType.HoldingRegister)
            {
                data = (await SendAsync<ModbusPDUReadHoldingRegistersResponse>(
                    new ModbusPDUReadHoldingRegistersRequest(first.WireAddress, registers),
                    cancellationToken).ConfigureAwait(false)).Value;
            }
            else
            {
                data = (await SendAsync<ModbusPDUReadInputRegistersResponse>(
                    new ModbusPDUReadInputRegistersRequest(first.WireAddress, registers),
                    cancellationToken).ConfigureAwait(false)).Value;
            }
            var offset = 0;
            for (var i = 0; i < count; i++)
            {
                var item = work[start + i];
                var bytes = item.Tag.RegisterCount * 2;
                results[item.Name] = new DefaultPlcResponseItem<IPlcValue>(
                    PlcResponseCode.Ok, RegisterValue(Slice(data, offset, bytes), item.Tag));
                offset += bytes;
            }
        }

        private async Task<IPlcValue> ReadTagAsync(ModbusTag tag, CancellationToken cancellationToken)
        {
            switch (tag.Type)
            {
                case ModbusTag.TagType.Coil:
                    {
                        var response = await SendAsync<ModbusPDUReadCoilsResponse>(
                            new ModbusPDUReadCoilsRequest(tag.WireAddress, tag.Quantity), cancellationToken).ConfigureAwait(false);
                        return BitValue(response.Value, tag);
                    }
                case ModbusTag.TagType.DiscreteInput:
                    {
                        var response = await SendAsync<ModbusPDUReadDiscreteInputsResponse>(
                            new ModbusPDUReadDiscreteInputsRequest(tag.WireAddress, tag.Quantity), cancellationToken).ConfigureAwait(false);
                        return BitValue(response.Value, tag);
                    }
                case ModbusTag.TagType.HoldingRegister:
                    {
                        var response = await SendAsync<ModbusPDUReadHoldingRegistersResponse>(
                            new ModbusPDUReadHoldingRegistersRequest(tag.WireAddress, tag.RegisterCount), cancellationToken).ConfigureAwait(false);
                        return RegisterValue(response.Value, tag);
                    }
                case ModbusTag.TagType.InputRegister:
                    {
                        var response = await SendAsync<ModbusPDUReadInputRegistersResponse>(
                            new ModbusPDUReadInputRegistersRequest(tag.WireAddress, tag.RegisterCount), cancellationToken).ConfigureAwait(false);
                        return RegisterValue(response.Value, tag);
                    }
                case ModbusTag.TagType.ExtendedRegister:
                    return await ReadExtendedRegisterAsync(tag, cancellationToken).ConfigureAwait(false);
                default:
                    throw new ModbusDriverException($"Unsupported tag type for a read: {tag.Type}.");
            }
        }

        // A device may answer with more data than the one element asked for; the
        // first element is the one that was requested.
        private static bool FirstBit(byte[] data)
        {
            if (data.Length < 1)
            {
                throw new ModbusDriverException("The Modbus response carries no data.");
            }
            return (data[0] & 0x01) != 0;
        }

        private static IPlcValue BitValue(byte[] data, ModbusTag tag)
        {
            if (data.Length * 8 < tag.Quantity)
            {
                throw new ModbusDriverException("The Modbus bit response is shorter than requested.");
            }
            if (!tag.ExplicitRange) return new PlcBOOL(FirstBit(data));
            var values = new List<IPlcValue>(tag.Quantity);
            for (var i = 0; i < tag.Quantity; i++)
            {
                values.Add(new PlcBOOL((data[i / 8] & (1 << (i % 8))) != 0));
            }
            return new PlcList(values);
        }

        /// <summary>Decodes one register-backed value with the generated <c>DataItem</c>.</summary>
        private static IPlcValue RegisterValue(byte[] data, ModbusTag tag)
        {
            var byteCount = tag.RegisterCount * 2;
            if (data.Length < byteCount)
            {
                throw new ModbusDriverException(
                    $"The Modbus response carries {data.Length} byte(s); {tag.DataType} needs {byteCount}.");
            }
            var canonical = Reorder(data, byteCount, tag.ByteOrder);
            var values = new List<IPlcValue>(tag.Quantity);
            var buffer = new ReadBuffer(canonical);
            if (tag.Quantity == 1 && IsEightBit(tag.DataType))
            {
                buffer.ReadByte("padding", 8);
            }
            for (var i = 0; i < tag.Quantity; i++)
            {
                values.Add(DataItem.StaticParse(buffer, tag.DataType, tag.StringLength));
            }
            return tag.ExplicitRange ? new PlcList(values) : values[0];
        }

        private async Task<IPlcValue> ReadExtendedRegisterAsync(
            ModbusTag tag, CancellationToken cancellationToken)
        {
            var items = FileReadItems(tag.WireAddress, tag.RegisterCount);
            var response = await SendAsync<ModbusPDUReadFileRecordResponse>(
                new ModbusPDUReadFileRecordRequest(items), cancellationToken).ConfigureAwait(false);
            if (response.Items.Count != items.Count)
            {
                throw new ModbusDriverException("The Modbus file-record response has the wrong item count.");
            }
            var bytes = new List<byte>();
            foreach (var item in response.Items) bytes.AddRange(item.Data);
            return RegisterValue(bytes.ToArray(), tag);
        }

        private static List<ModbusPDUReadFileRecordRequestItem> FileReadItems(ushort address, ushort count)
        {
            var items = new List<ModbusPDUReadFileRecordRequestItem>();
            var current = (int)address;
            var remaining = (int)count;
            while (remaining > 0)
            {
                var record = current % 10000;
                var take = Math.Min(remaining, 10000 - record);
                items.Add(new ModbusPDUReadFileRecordRequestItem(
                    0x06, checked((ushort)(current / 10000 + 1)), checked((ushort)record), checked((ushort)take)));
                current += take;
                remaining -= take;
            }
            return items;
        }

        /// <summary>
        /// Encodes a value for a register with the generated <c>DataItem</c>, so the value
        /// has to fit the tag's data type: an <c>INT</c> tag takes -32768..32767, a
        /// <c>UINT</c> or <c>WORD</c> tag 0..65535.
        /// </summary>
        /// <exception cref="OverflowException">The value does not fit the data type.</exception>
        /// <exception cref="FormatException">The value is not a whole number.</exception>
        /// <exception cref="InvalidCastException">The value is not a number.</exception>
        private static byte[] RegisterBytes(object value, ModbusTag tag)
        {
            var objects = ValuesOf(value, tag.Quantity);
            var buffer = new WriteBuffer();
            if (tag.Quantity == 1 && IsEightBit(tag.DataType))
            {
                buffer.WriteByte("padding", 8, 0);
            }
            foreach (var item in objects)
            {
                DataItem.StaticSerialize(buffer, PlcValueOf(item, tag.DataType), tag.DataType, tag.StringLength);
            }
            var encoded = buffer.GetBytes();
            var byteCount = tag.RegisterCount * 2;
            Array.Resize(ref encoded, byteCount);
            return Reorder(encoded, byteCount, tag.ByteOrder);
        }

        private static IPlcValue PlcValueOf(object value, ModbusDataType dataType)
        {
            IPlcValue plcValue;
            switch (dataType)
            {
                case ModbusDataType.BYTE: plcValue = new PlcBYTE(checked((byte)ToUnsigned(value))); break;
                case ModbusDataType.WORD: plcValue = new PlcWORD(checked((ushort)ToUnsigned(value))); break;
                case ModbusDataType.DWORD: plcValue = new PlcDWORD(checked((uint)ToUnsigned(value))); break;
                case ModbusDataType.LWORD: plcValue = new PlcLWORD(ToUnsigned(value)); break;
                case ModbusDataType.SINT: plcValue = new PlcSINT(checked((sbyte)ToInteger(value))); break;
                case ModbusDataType.INT: plcValue = new PlcINT(checked((short)ToInteger(value))); break;
                case ModbusDataType.DINT: plcValue = new PlcDINT(checked((int)ToInteger(value))); break;
                case ModbusDataType.LINT: plcValue = new PlcLINT(ToInteger(value)); break;
                case ModbusDataType.USINT: plcValue = new PlcUSINT(checked((byte)ToUnsigned(value))); break;
                case ModbusDataType.UINT: plcValue = new PlcUINT(checked((ushort)ToUnsigned(value))); break;
                case ModbusDataType.UDINT: plcValue = new PlcUDINT(checked((uint)ToUnsigned(value))); break;
                case ModbusDataType.ULINT: plcValue = new PlcULINT(ToUnsigned(value)); break;
                case ModbusDataType.REAL: plcValue = new PlcREAL(Convert.ToSingle(value, CultureInfo.InvariantCulture)); break;
                case ModbusDataType.LREAL: plcValue = new PlcLREAL(Convert.ToDouble(value, CultureInfo.InvariantCulture)); break;
                case ModbusDataType.STRING: plcValue = new PlcSTRING(Convert.ToString(value, CultureInfo.InvariantCulture) ?? string.Empty); break;
                case ModbusDataType.WSTRING: plcValue = new PlcWSTRING(Convert.ToString(value, CultureInfo.InvariantCulture) ?? string.Empty); break;
                default:
                    throw new ModbusDriverException($"A register cannot be written as {dataType}.");
            }
            return plcValue;
        }

        private static List<object> ValuesOf(object value, ushort quantity)
        {
            if (quantity == 1) return new List<object> { value };
            if (!(value is IEnumerable enumerable) || value is string)
            {
                throw new InvalidCastException($"A selection of {quantity} values requires an enumerable value.");
            }
            var values = new List<object>();
            foreach (var item in enumerable)
            {
                if (item == null) throw new InvalidCastException("A Modbus array cannot contain null.");
                values.Add(item);
            }
            if (values.Count != quantity)
            {
                throw new InvalidCastException($"The value has {values.Count} elements; the tag requires {quantity}.");
            }
            return values;
        }

        private async Task<PlcResponseCode> WriteExtendedRegisterAsync(
            ModbusTag tag, object value, CancellationToken cancellationToken)
        {
            var bytes = RegisterBytes(value, tag);
            var items = new List<ModbusPDUWriteFileRecordRequestItem>();
            var current = (int)tag.WireAddress;
            var byteOffset = 0;
            while (byteOffset < bytes.Length)
            {
                var record = current % 10000;
                var words = Math.Min((bytes.Length - byteOffset) / 2, 10000 - record);
                items.Add(new ModbusPDUWriteFileRecordRequestItem(
                    0x06, checked((ushort)(current / 10000 + 1)), checked((ushort)record),
                    Slice(bytes, byteOffset, words * 2)));
                current += words;
                byteOffset += words * 2;
            }
            await SendAsync<ModbusPDUWriteFileRecordResponse>(
                new ModbusPDUWriteFileRecordRequest(items), cancellationToken).ConfigureAwait(false);
            return PlcResponseCode.Ok;
        }

        private static bool IsEightBit(ModbusDataType dataType)
            => dataType == ModbusDataType.BYTE || dataType == ModbusDataType.SINT || dataType == ModbusDataType.USINT;

        // Each layout is its own inverse, so the same operation converts from and to canonical
        // big-endian bytes.
        private static byte[] Reorder(byte[] source, int count, ModbusByteOrder byteOrder)
        {
            var result = Slice(source, 0, count);
            if (byteOrder == ModbusByteOrder.LittleEndian)
            {
                Array.Reverse(result);
            }
            else if (byteOrder == ModbusByteOrder.BigEndianByteSwap)
            {
                for (var i = 0; i < result.Length; i += 2)
                {
                    var value = result[i]; result[i] = result[i + 1]; result[i + 1] = value;
                }
            }
            else if (byteOrder == ModbusByteOrder.LittleEndianByteSwap)
            {
                for (int left = 0, right = result.Length - 2; left < right; left += 2, right -= 2)
                {
                    var first = result[left]; var second = result[left + 1];
                    result[left] = result[right]; result[left + 1] = result[right + 1];
                    result[right] = first; result[right + 1] = second;
                }
            }
            return result;
        }

        private static byte[] Slice(byte[] source, int offset, int count)
        {
            var result = new byte[count];
            Array.Copy(source, offset, result, 0, count);
            return result;
        }

        private static long ToInteger(object value)
        {
            switch (value)
            {
                case bool _:
                    throw new InvalidCastException("A boolean cannot be written to a register.");
                case float f:
                    return WholeNumber(f);
                case double d:
                    return WholeNumber(d);
                case decimal m:
                    return WholeNumber((double)m);
                case ulong u:
                    return checked((long)u);
                case string s:
                    return long.Parse(s, NumberStyles.Integer, CultureInfo.InvariantCulture);
                default:
                    return Convert.ToInt64(value, CultureInfo.InvariantCulture);
            }
        }

        private static ulong ToUnsigned(object value)
        {
            if (value is bool)
            {
                throw new InvalidCastException("A boolean cannot be written to a register.");
            }
            if (value is float f)
            {
                return checked((ulong)WholeNumber(f));
            }
            if (value is double d)
            {
                return checked((ulong)WholeNumber(d));
            }
            if (value is decimal m)
            {
                return checked((ulong)WholeNumber((double)m));
            }
            if (value is string s)
            {
                return ulong.Parse(s, NumberStyles.Integer, CultureInfo.InvariantCulture);
            }
            return Convert.ToUInt64(value, CultureInfo.InvariantCulture);
        }

        private static long WholeNumber(double value)
        {
            if (double.IsNaN(value) || double.IsInfinity(value) || value != Math.Floor(value))
            {
                throw new FormatException($"{value.ToString(CultureInfo.InvariantCulture)} is not a whole number.");
            }
            return checked((long)value);
        }

        // ── PlcWriter ──────────────────────────────────────────

        public async Task<IPlcWriteResponse> Write(
            DefaultPlcWriteRequest request, CancellationToken cancellationToken = default)
        {
            var codes = new Dictionary<string, PlcResponseCode>();

            foreach (var name in request.TagNames)
            {
                cancellationToken.ThrowIfCancellationRequested();

                var value = request.GetValue(name);
                if (!(request.GetTagByName(name) is ModbusTag tag))
                {
                    codes[name] = PlcResponseCode.InvalidAddress;
                    continue;
                }
                if (value == null)
                {
                    codes[name] = PlcResponseCode.InvalidDatatype;
                    continue;
                }

                try
                {
                    codes[name] = await WriteTagAsync(tag, value, cancellationToken).ConfigureAwait(false);
                }
                catch (OperationCanceledException)
                {
                    throw; // Cancellation is the caller's decision, not a per-tag error.
                }
                catch (Exception ex)
                {
                    Logger.LogWarning(ex, "Modbus write of tag '{TagName}' failed", name);
                    codes[name] = MapException(ex);
                }
            }

            return new DefaultPlcWriteResponse(request, codes);
        }

        private async Task<PlcResponseCode> WriteTagAsync(
            ModbusTag tag, object value, CancellationToken cancellationToken)
        {
            switch (tag.Type)
            {
                case ModbusTag.TagType.Coil:
                    {
                        var values = ValuesOf(value, tag.Quantity);
                        if (tag.Quantity == 1)
                        {
                            var on = Convert.ToBoolean(values[0], CultureInfo.InvariantCulture);
                            await SendAsync<ModbusPDUWriteSingleCoilResponse>(
                                new ModbusPDUWriteSingleCoilRequest(tag.WireAddress, on ? (ushort)0xFF00 : (ushort)0x0000),
                                cancellationToken).ConfigureAwait(false);
                        }
                        else
                        {
                            var packed = new byte[(tag.Quantity + 7) / 8];
                            for (var i = 0; i < values.Count; i++)
                            {
                                if (Convert.ToBoolean(values[i], CultureInfo.InvariantCulture)) packed[i / 8] |= (byte)(1 << (i % 8));
                            }
                            await SendAsync<ModbusPDUWriteMultipleCoilsResponse>(
                                new ModbusPDUWriteMultipleCoilsRequest(tag.WireAddress, tag.Quantity, packed),
                                cancellationToken).ConfigureAwait(false);
                        }
                        return PlcResponseCode.Ok;
                    }
                case ModbusTag.TagType.HoldingRegister:
                    {
                        // Converted before anything is sent: a value that does not fit the
                        // tag's type never reaches the device.
                        var bytes = RegisterBytes(value, tag);
                        if (tag.RegisterCount == 1)
                        {
                            var register = (ushort)((bytes[0] << 8) | bytes[1]);
                            await SendAsync<ModbusPDUWriteSingleRegisterResponse>(
                                new ModbusPDUWriteSingleRegisterRequest(tag.WireAddress, register),
                                cancellationToken).ConfigureAwait(false);
                        }
                        else
                        {
                            await SendAsync<ModbusPDUWriteMultipleHoldingRegistersResponse>(
                                new ModbusPDUWriteMultipleHoldingRegistersRequest(
                                    tag.WireAddress, tag.RegisterCount, bytes), cancellationToken).ConfigureAwait(false);
                        }
                        return PlcResponseCode.Ok;
                    }
                case ModbusTag.TagType.ExtendedRegister:
                    return await WriteExtendedRegisterAsync(tag, value, cancellationToken).ConfigureAwait(false);
                default:
                    // Discrete inputs and input registers are read-only.
                    return PlcResponseCode.AccessDenied;
            }
        }

        // ── Exchange ───────────────────────────────────────────

        /// <summary>
        /// Exchanges a request and checks that the answer is the response type that
        /// belongs to it. An exception response from the device becomes a
        /// <see cref="ModbusDriverException"/> carrying its exception code.
        /// </summary>
        private async Task<TResponse> SendAsync<TResponse>(
            ModbusPDU request, CancellationToken cancellationToken) where TResponse : ModbusPDU
        {
            var response = await ExchangeAsync(request, cancellationToken).ConfigureAwait(false);

            if (response is ModbusPDUError error)
            {
                var code = (byte)error.ExceptionCode;
                var name = Enum.IsDefined(typeof(ModbusErrorCode), error.ExceptionCode)
                    ? error.ExceptionCode.ToString()
                    : "unknown";
                throw new ModbusDriverException(
                    $"Modbus exception response: function 0x{request.FunctionFlag:X2}, code 0x{code:X2} ({name}).",
                    code);
            }

            return response as TResponse
                ?? throw new ModbusDriverException(
                    $"Unexpected Modbus response {response.GetType().Name}; expected {typeof(TResponse).Name}.");
        }

        /// <summary>Maps a failed exchange to a response code.</summary>
        internal static PlcResponseCode MapException(Exception exception)
        {
            switch (exception)
            {
                case TimeoutException _:
                    return PlcResponseCode.RequestTimeout;
                case ModbusDriverException modbus when modbus.ExceptionCode != 0:
                    switch (modbus.ExceptionCode)
                    {
                        case (byte)ModbusErrorCode.ILLEGAL_FUNCTION:
                            return PlcResponseCode.Unsupported;
                        case (byte)ModbusErrorCode.ILLEGAL_DATA_ADDRESS:
                            return PlcResponseCode.InvalidAddress;
                        case (byte)ModbusErrorCode.ILLEGAL_DATA_VALUE:
                            return PlcResponseCode.InvalidDatatype;
                        case (byte)ModbusErrorCode.ACKNOWLEDGE:
                            // The device has accepted the request and will process it, as
                            // plc4j also reports it. The codes plc4j answers with dedicated
                            // remote-busy/remote-error values have no counterpart in
                            // PlcResponseCode and fall through to InternalError.
                            return PlcResponseCode.Ok;
                        case (byte)ModbusErrorCode.SLAVE_DEVICE_BUSY:
                            return PlcResponseCode.RequestTimeout;
                        default:
                            return PlcResponseCode.InternalError;
                    }
                case FormatException _:
                case InvalidCastException _:
                case OverflowException _:
                    // The value to write cannot be represented in the tag's type.
                    return PlcResponseCode.InvalidDatatype;
                default:
                    return PlcResponseCode.InternalError;
            }
        }
    }
}
