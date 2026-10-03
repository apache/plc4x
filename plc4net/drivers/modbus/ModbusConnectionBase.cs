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

            foreach (var name in request.TagNames)
            {
                cancellationToken.ThrowIfCancellationRequested();

                if (!(request.GetTagByName(name) is ModbusTag tag))
                {
                    results[name] = new DefaultPlcTagErrorItem<IPlcValue>(PlcResponseCode.InvalidAddress);
                    continue;
                }

                try
                {
                    var value = await ReadTagAsync(tag, cancellationToken).ConfigureAwait(false);
                    results[name] = new DefaultPlcResponseItem<IPlcValue>(PlcResponseCode.Ok, value);
                }
                catch (OperationCanceledException)
                {
                    throw; // Cancellation is the caller's decision, not a per-tag error.
                }
                catch (Exception ex)
                {
                    Logger.LogWarning(ex, "Modbus read of tag '{TagName}' failed", name);
                    results[name] = new DefaultPlcTagErrorItem<IPlcValue>(MapException(ex));
                }
            }

            return new DefaultPlcReadResponse(request, results);
        }

        private async Task<IPlcValue> ReadTagAsync(ModbusTag tag, CancellationToken cancellationToken)
        {
            switch (tag.Type)
            {
                case ModbusTag.TagType.Coil:
                    {
                        var response = await SendAsync<ModbusPDUReadCoilsResponse>(
                            new ModbusPDUReadCoilsRequest(tag.WireAddress, 1), cancellationToken).ConfigureAwait(false);
                        return new PlcBOOL(FirstBit(response.Value));
                    }
                case ModbusTag.TagType.DiscreteInput:
                    {
                        var response = await SendAsync<ModbusPDUReadDiscreteInputsResponse>(
                            new ModbusPDUReadDiscreteInputsRequest(tag.WireAddress, 1), cancellationToken).ConfigureAwait(false);
                        return new PlcBOOL(FirstBit(response.Value));
                    }
                case ModbusTag.TagType.HoldingRegister:
                    {
                        var response = await SendAsync<ModbusPDUReadHoldingRegistersResponse>(
                            new ModbusPDUReadHoldingRegistersRequest(tag.WireAddress, 1), cancellationToken).ConfigureAwait(false);
                        return RegisterValue(response.Value, tag.DataType);
                    }
                case ModbusTag.TagType.InputRegister:
                    {
                        var response = await SendAsync<ModbusPDUReadInputRegistersResponse>(
                            new ModbusPDUReadInputRegistersRequest(tag.WireAddress, 1), cancellationToken).ConfigureAwait(false);
                        return RegisterValue(response.Value, tag.DataType);
                    }
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

        /// <summary>Decodes the first register of a response with the generated <c>DataItem</c>.</summary>
        private static IPlcValue RegisterValue(byte[] data, ModbusDataType dataType)
        {
            if (data.Length < 2)
            {
                throw new ModbusDriverException(
                    $"The Modbus response carries {data.Length} byte(s); a register is 2.");
            }
            return DataItem.StaticParse(new ReadBuffer(new[] { data[0], data[1] }), dataType, 0);
        }

        /// <summary>
        /// Encodes a value for a register with the generated <c>DataItem</c>, so the value
        /// has to fit the tag's data type: an <c>INT</c> tag takes -32768..32767, a
        /// <c>UINT</c> or <c>WORD</c> tag 0..65535.
        /// </summary>
        /// <exception cref="OverflowException">The value does not fit the data type.</exception>
        /// <exception cref="FormatException">The value is not a whole number.</exception>
        /// <exception cref="InvalidCastException">The value is not a number.</exception>
        private static ushort RegisterWord(object value, ModbusDataType dataType)
        {
            var number = ToInteger(value);
            IPlcValue plcValue;
            switch (dataType)
            {
                case ModbusDataType.INT:
                    plcValue = new PlcINT(checked((short)number));
                    break;
                case ModbusDataType.UINT:
                    plcValue = new PlcUINT(checked((ushort)number));
                    break;
                case ModbusDataType.WORD:
                    plcValue = new PlcWORD(checked((ushort)number));
                    break;
                default:
                    throw new ModbusDriverException($"A register cannot be written as {dataType}.");
            }

            var buffer = new WriteBuffer();
            DataItem.StaticSerialize(buffer, plcValue, dataType, 0);
            var bytes = buffer.GetBytes();
            return (ushort)((bytes[0] << 8) | bytes[1]);
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
                        // A coil is written as 0xFF00 (on) or 0x0000 (off).
                        var on = value is bool b ? b : Convert.ToBoolean(value, CultureInfo.InvariantCulture);
                        await SendAsync<ModbusPDUWriteSingleCoilResponse>(
                            new ModbusPDUWriteSingleCoilRequest(tag.WireAddress, on ? (ushort)0xFF00 : (ushort)0x0000),
                            cancellationToken).ConfigureAwait(false);
                        return PlcResponseCode.Ok;
                    }
                case ModbusTag.TagType.HoldingRegister:
                    {
                        // Converted before anything is sent: a value that does not fit the
                        // tag's type never reaches the device.
                        var register = RegisterWord(value, tag.DataType);
                        await SendAsync<ModbusPDUWriteSingleRegisterResponse>(
                            new ModbusPDUWriteSingleRegisterRequest(tag.WireAddress, register),
                            cancellationToken).ConfigureAwait(false);
                        return PlcResponseCode.Ok;
                    }
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