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
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

using System;
using System.Collections.Generic;
using System.Globalization;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;
using Microsoft.Extensions.Logging;
using org.apache.plc4net.api.metadata;
using org.apache.plc4net.api.value;
using org.apache.plc4net.drivers.slmp.readwrite.model;
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

namespace org.apache.plc4net.drivers.slmp
{
    /// <summary>Executes SLMP 3E binary batch reads and writes.</summary>
    public sealed class SlmpConnection : ConnectionBase, PlcReader, PlcWriter
    {
        private const int PollIntervalMs = 2;
        private const int DefaultRequestTimeoutMs = 5000;
        private const ushort DefaultMonitoringTimer = 16;

        private static readonly DefaultPlcConnectionMetadata Metadata = new DefaultPlcConnectionMetadata
        {
            CanRead = true,
            CanWrite = true,
            CanSubscribe = false
        };

        private readonly SemaphoreSlim _exchangeGate = new SemaphoreSlim(1, 1);
        private readonly SlmpMessageCodec _codec;
        private readonly int _requestTimeoutMs;
        private readonly ushort _monitoringTimer;
        private SlmpResponseFrame3E? _response;
        private bool _unusable;

        public SlmpConnection(ConnectionString connectionString, ITransportInstance transport)
            : base(connectionString, transport)
        {
            _requestTimeoutMs = IntParameter(connectionString, DefaultRequestTimeoutMs, "request-timeout-ms");
            if (_requestTimeoutMs < 1)
            {
                throw new PlcConnectionException("request-timeout-ms must be positive.");
            }

            var monitoringTimer = IntParameter(connectionString, DefaultMonitoringTimer, "monitoring-timer");
            if (monitoringTimer < 0 || monitoringTimer > ushort.MaxValue)
            {
                throw new PlcConnectionException("monitoring-timer must be between 0 and 65535.");
            }
            _monitoringTimer = (ushort)monitoringTimer;
            _codec = new SlmpMessageCodec(transport, OnMessage);

            if (transport is IAsyncTransportInstance asyncTransport)
            {
                asyncTransport.RegisterDisconnectListener(ex =>
                    Logger.LogWarning(ex, "The SLMP transport disconnected"));
            }
        }

        public override IPlcConnectionMetadata PlcConnectionMetadata => Metadata;
        public override IPlcTag Parse(string tagQuery) => SlmpTag.Parse(tagQuery);
        public override IPlcReadRequestBuilder? ReadRequestBuilder => new DefaultPlcReadRequestBuilder(this, Parse);
        public override IPlcWriteRequestBuilder? WriteRequestBuilder => new SlmpWriteRequestBuilder(this, Parse);
        public override IPlcSubscriptionRequestBuilder? SubscriptionRequestBuilder => null;
        public override IPlcUnsubscriptionRequestBuilder? UnsubscriptionRequestBuilder => null;

        public async Task<IPlcReadResponse> Read(
            DefaultPlcReadRequest request, CancellationToken cancellationToken = default)
        {
            var results = new Dictionary<string, PlcResponseItem<IPlcValue>>();
            foreach (var name in request.TagNames)
            {
                cancellationToken.ThrowIfCancellationRequested();
                if (request.GetTagByName(name) is not SlmpTag tag)
                {
                    results[name] = new DefaultPlcTagErrorItem<IPlcValue>(PlcResponseCode.InvalidAddress);
                    continue;
                }

                try
                {
                    var data = new SlmpReadRequest(
                        (uint)tag.DeviceNumber, tag.DeviceCode, checked((ushort)tag.NumberOfPoints));
                    var frame = new SlmpRequestFrame3E(_monitoringTimer, 0x0401, 0, data);
                    var response = await ExchangeAsync(frame, cancellationToken).ConfigureAwait(false);
                    if (response.EndCode != 0)
                    {
                        Logger.LogWarning("SLMP returned end code 0x{EndCode:X4} for tag {TagName}",
                            response.EndCode, name);
                        results[name] = new DefaultPlcTagErrorItem<IPlcValue>(PlcResponseCode.RemoteError);
                        continue;
                    }

                    var expectedBytes = checked(tag.NumberOfPoints * 2);
                    if (response.ResponseData.Length != expectedBytes)
                    {
                        results[name] = new DefaultPlcTagErrorItem<IPlcValue>(PlcResponseCode.InvalidData);
                        continue;
                    }
                    results[name] = new DefaultPlcResponseItem<IPlcValue>(
                        PlcResponseCode.Ok, Decode(tag, response.ResponseData));
                }
                catch (OperationCanceledException)
                {
                    throw;
                }
                catch (TimeoutException ex)
                {
                    Logger.LogWarning(ex, "SLMP read of tag '{TagName}' timed out", name);
                    results[name] = new DefaultPlcTagErrorItem<IPlcValue>(PlcResponseCode.RequestTimeout);
                }
                catch (Exception ex)
                {
                    Logger.LogWarning(ex, "SLMP read of tag '{TagName}' failed", name);
                    results[name] = new DefaultPlcTagErrorItem<IPlcValue>(PlcResponseCode.InternalError);
                }
            }
            return new DefaultPlcReadResponse(request, results);
        }

        public async Task<IPlcWriteResponse> Write(
            DefaultPlcWriteRequest request, CancellationToken cancellationToken = default)
        {
            var results = new Dictionary<string, PlcResponseCode>();
            foreach (var name in request.TagNames)
            {
                cancellationToken.ThrowIfCancellationRequested();
                if (request.GetTagByName(name) is not SlmpTag tag)
                {
                    results[name] = PlcResponseCode.InvalidAddress;
                    continue;
                }

                try
                {
                    var data = Encode(tag, request.GetValue(name));
                    if (data == null)
                    {
                        results[name] = PlcResponseCode.InvalidData;
                        continue;
                    }

                    var write = new SlmpWriteRequest(
                        (uint)tag.DeviceNumber, tag.DeviceCode, checked((ushort)tag.NumberOfPoints), data);
                    var frame = new SlmpRequestFrame3E(_monitoringTimer, 0x1401, 0, write);
                    var response = await ExchangeAsync(frame, cancellationToken).ConfigureAwait(false);
                    if (response.EndCode != 0)
                    {
                        Logger.LogWarning("SLMP returned end code 0x{EndCode:X4} for write tag {TagName}",
                            response.EndCode, name);
                        results[name] = PlcResponseCode.RemoteError;
                    }
                    else
                    {
                        // A successful 3E Batch Write response contains only the end code.
                        results[name] = response.ResponseData.Length == 0
                            ? PlcResponseCode.Ok
                            : PlcResponseCode.RemoteError;
                    }
                }
                catch (OperationCanceledException)
                {
                    throw;
                }
                catch (TimeoutException ex)
                {
                    Logger.LogWarning(ex, "SLMP write of tag '{TagName}' timed out", name);
                    results[name] = PlcResponseCode.RequestTimeout;
                }
                catch (Exception ex)
                {
                    Logger.LogWarning(ex, "SLMP write of tag '{TagName}' failed", name);
                    results[name] = PlcResponseCode.InternalError;
                }
            }
            return new DefaultPlcWriteResponse(request, results);
        }

        private static byte[]? Encode(SlmpTag tag, object? value)
        {
            if (value == null)
            {
                return null;
            }
            var values = new List<object>(tag.Quantity);
            if (tag.ExplicitRange)
            {
                if (value is object[] array)
                {
                    if (array.Length != tag.Quantity || array.Any(item => item == null))
                    {
                        return null;
                    }
                    values.AddRange(array);
                }
                else if (value is IPlcValue list && list.IsList() && list.GetLength() == tag.Quantity)
                {
                    for (var i = 0; i < tag.Quantity; i++)
                    {
                        values.Add(list.GetIndex(i));
                    }
                }
                else
                {
                    return null;
                }
            }
            else
            {
                if (value is object[] || value is IPlcValue plcValue && plcValue.IsList())
                {
                    return null;
                }
                values.Add(value);
            }

            var buffer = new WriteBuffer();
            buffer.PushByteOrder(ByteOrder.LITTLE_ENDIAN);
            try
            {
                foreach (var item in values)
                {
                    WriteOne(buffer, tag.DataType, item);
                }
                var data = buffer.GetBytes();
                return data.Length == tag.NumberOfPoints * 2 ? data : null;
            }
            catch (Exception ex) when (ex is FormatException || ex is InvalidCastException ||
                                       ex is OverflowException || ex is ArgumentException ||
                                       ex is PlcIncompatibleDatatypeException)
            {
                return null;
            }
            finally
            {
                buffer.PopByteOrder();
            }
        }

        private static void WriteOne(WriteBuffer buffer, SlmpDataType type, object value)
        {
            var plcValue = value as IPlcValue;
            if (plcValue != null && !IsNumeric(plcValue))
            {
                throw new PlcIncompatibleDatatypeException("SLMP values must be numeric.");
            }
            switch (type)
            {
                case SlmpDataType.WORD:
                case SlmpDataType.UINT:
                    buffer.WriteUshort("value", 16, plcValue != null
                        ? checked((ushort)plcValue.GetInt())
                        : Convert.ToUInt16(value, CultureInfo.InvariantCulture));
                    break;
                case SlmpDataType.INT:
                    buffer.WriteShort("value", 16, plcValue != null
                        ? plcValue.GetShort()
                        : Convert.ToInt16(value, CultureInfo.InvariantCulture));
                    break;
                case SlmpDataType.DINT:
                    buffer.WriteInt("value", 32, plcValue != null
                        ? plcValue.GetInt()
                        : Convert.ToInt32(value, CultureInfo.InvariantCulture));
                    break;
                case SlmpDataType.UDINT:
                    buffer.WriteUint("value", 32, plcValue != null
                        ? checked((uint)plcValue.GetLong())
                        : Convert.ToUInt32(value, CultureInfo.InvariantCulture));
                    break;
                case SlmpDataType.REAL:
                    buffer.WriteFloat("value", 32, plcValue != null
                        ? plcValue.GetFloat()
                        : Convert.ToSingle(value, CultureInfo.InvariantCulture));
                    break;
                default:
                    throw new ArgumentOutOfRangeException(nameof(type), type, "Unsupported SLMP data type.");
            }
        }

        private static bool IsNumeric(IPlcValue value)
            => value.IsByte() || value.IsUshort() || value.IsUint() || value.IsUlong() ||
               value.IsSbyte() || value.IsShort() || value.IsInt() || value.IsLong() ||
               value.IsFloat() || value.IsDouble();

        private async Task<SlmpResponseFrame3E> ExchangeAsync(
            SlmpRequestFrame3E request, CancellationToken cancellationToken)
        {
            await _exchangeGate.WaitAsync(cancellationToken).ConfigureAwait(false);
            try
            {
                if (_unusable)
                {
                    throw new InvalidOperationException(
                        "The SLMP connection cannot be reused after a timed-out or cancelled exchange.");
                }
                if (!TransportInstance.IsOpen)
                {
                    throw new InvalidOperationException("The SLMP transport is closed.");
                }

                _response = null;
                _codec.Send(request);
                var deadline = Environment.TickCount64 + _requestTimeoutMs;
                while (true)
                {
                    cancellationToken.ThrowIfCancellationRequested();
                    _codec.ProcessIncomingData();
                    if (_response != null)
                    {
                        return _response;
                    }
                    if (!TransportInstance.IsOpen)
                    {
                        throw new InvalidOperationException("The SLMP transport disconnected.");
                    }
                    if (Environment.TickCount64 >= deadline)
                    {
                        _unusable = true;
                        throw new TimeoutException(
                            $"No SLMP response arrived within {_requestTimeoutMs} ms.");
                    }
                    await Task.Delay(PollIntervalMs, cancellationToken).ConfigureAwait(false);
                }
            }
            catch (OperationCanceledException)
            {
                // 3E frames have no transaction identifier. A late answer could otherwise
                // be mistaken for the response to the next request on this connection.
                _unusable = true;
                throw;
            }
            catch (MessageCodecException ex)
            {
                _unusable = true;
                throw new InvalidOperationException($"SLMP exchange failed: {ex.Message}", ex);
            }
            finally
            {
                _response = null;
                _exchangeGate.Release();
            }
        }

        private void OnMessage(SlmpMessage message)
        {
            if (_response != null)
            {
                Logger.LogWarning("Discarding an unsolicited SLMP message while a response is pending");
                return;
            }
            if (message is SlmpResponseFrame3E response)
            {
                _response = response;
            }
        }

        private static IPlcValue Decode(SlmpTag tag, byte[] data)
        {
            var buffer = new ReadBuffer(data);
            buffer.PushByteOrder(ByteOrder.LITTLE_ENDIAN);
            try
            {
                if (!tag.ExplicitRange)
                {
                    return DecodeOne(buffer, tag.DataType);
                }

                var values = new List<IPlcValue>(tag.Quantity);
                for (var i = 0; i < tag.Quantity; i++)
                {
                    values.Add(DecodeOne(buffer, tag.DataType));
                }
                return new PlcList(values);
            }
            finally
            {
                buffer.PopByteOrder();
            }
        }

        private static IPlcValue DecodeOne(ReadBuffer buffer, SlmpDataType type) => type switch
        {
            SlmpDataType.WORD => new PlcWORD(buffer.ReadUshort("value", 16)),
            SlmpDataType.INT => new PlcINT(buffer.ReadShort("value", 16)),
            SlmpDataType.UINT => new PlcUINT(buffer.ReadUshort("value", 16)),
            SlmpDataType.DINT => new PlcDINT(buffer.ReadInt("value", 32)),
            SlmpDataType.UDINT => new PlcUDINT(buffer.ReadUint("value", 32)),
            SlmpDataType.REAL => new PlcREAL(buffer.ReadFloat("value", 32)),
            _ => throw new ArgumentOutOfRangeException(nameof(type), type, "Unsupported SLMP data type.")
        };

        private static int IntParameter(ConnectionString connectionString, int defaultValue, string name)
        {
            var raw = connectionString.GetParameter(name);
            if (raw == null)
            {
                return defaultValue;
            }
            if (!int.TryParse(raw, NumberStyles.Integer, CultureInfo.InvariantCulture, out var value))
            {
                throw new PlcConnectionException(
                    $"The connection string parameter '{name}' must be a whole number, but was '{raw}'.");
            }
            return value;
        }

        protected override void Dispose(bool disposing)
        {
            if (disposing)
            {
                _exchangeGate.Dispose();
            }
            base.Dispose(disposing);
        }
    }
}
