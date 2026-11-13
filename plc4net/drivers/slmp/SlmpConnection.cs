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
using org.apache.plc4net.spi.drivers.subscriptions;
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
        private const int MaxRandomReadPoints = 192;
        private const int MaxMultiBlockReadBlocks = 120;
        private const int MaxMultiBlockReadPoints = 960;
        private const ushort DefaultMonitoringTimer = 16;

        private static readonly DefaultPlcConnectionMetadata Metadata = new DefaultPlcConnectionMetadata
        {
            CanRead = true,
            CanWrite = true,
            CanSubscribe = true
        };

        private readonly SemaphoreSlim _exchangeGate = new SemaphoreSlim(1, 1);
        private readonly SlmpMessageCodec _codec;
        private readonly PollingSubscriptionManager _subscriptions;
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
            var subscriptionInterval = IntParameter(connectionString, 1000, "subscription-polling-interval-ms");
            if (subscriptionInterval < 1)
                throw new PlcConnectionException("subscription-polling-interval-ms must be positive.");
            _subscriptions = new PollingSubscriptionManager(this, Parse,
                TimeSpan.FromMilliseconds(subscriptionInterval),
                ex => Logger.LogWarning(ex, "SLMP subscription observer or read failed"));

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
        public override IPlcSubscriptionRequestBuilder? SubscriptionRequestBuilder => _subscriptions.SubscriptionRequestBuilder;
        public override IPlcUnsubscriptionRequestBuilder? UnsubscriptionRequestBuilder => _subscriptions.UnsubscriptionRequestBuilder;

        public override void Close()
        {
            _subscriptions.DisposeAsync().AsTask().GetAwaiter().GetResult();
            base.Close();
        }

        public async Task<IPlcReadResponse> Read(
            DefaultPlcReadRequest request, CancellationToken cancellationToken = default)
        {
            cancellationToken.ThrowIfCancellationRequested();
            var results = new Dictionary<string, PlcResponseItem<IPlcValue>>();
            var names = request.TagNames.ToArray();
            var randomItems = new List<ReadItem>();
            var multiBlockItems = new List<ReadItem>();
            var batchItems = new List<ReadItem>();
            foreach (var name in names)
            {
                if (request.GetTagByName(name) is not SlmpTag tag)
                {
                    results[name] = new DefaultPlcTagErrorItem<IPlcValue>(PlcResponseCode.InvalidAddress);
                    continue;
                }
                var item = new ReadItem(name, tag);
                if (!tag.IsBitDevice && !tag.ExplicitRange && tag.Quantity == 1)
                    randomItems.Add(item);
                else if (tag.ExplicitRange && item.BlockPoints < MaxMultiBlockReadPoints &&
                         (!tag.IsBitDevice || tag.Quantity % 16 == 0))
                    multiBlockItems.Add(item);
                else
                    batchItems.Add(item);
            }

            while (multiBlockItems.Count >= 2)
            {
                var group = new List<ReadItem>();
                var points = 0;
                foreach (var item in multiBlockItems.OrderBy(item => item.Tag.IsBitDevice))
                {
                    if (group.Count == MaxMultiBlockReadBlocks || points + item.BlockPoints > MaxMultiBlockReadPoints)
                        continue;
                    group.Add(item);
                    points += item.BlockPoints;
                }
                if (group.Count < 2)
                {
                    batchItems.Add(multiBlockItems[0]);
                    multiBlockItems.RemoveAt(0);
                    continue;
                }
                await ReadMultiBlockGroup(group.ToArray(), results, cancellationToken).ConfigureAwait(false);
                var selected = group.ToHashSet();
                multiBlockItems.RemoveAll(selected.Contains);
            }

            // Random Read saves a round-trip only when at least two points can share a frame.
            // SLMP subcommand 0000 permits at most 192 word and double-word points combined.
            while (randomItems.Count >= 2)
            {
                var words = randomItems.Where(item => item.Tag.NumberOfPoints == 1)
                    .Take(MaxRandomReadPoints).ToArray();
                var doubleWords = randomItems.Where(item => item.Tag.NumberOfPoints == 2)
                    .Take(MaxRandomReadPoints - words.Length).ToArray();
                var group = words.Concat(doubleWords).ToArray();
                if (group.Length < 2)
                    break;
                await ReadRandomGroup(group, words.Length, doubleWords.Length, results, cancellationToken)
                    .ConfigureAwait(false);
                var selected = group.ToHashSet();
                randomItems.RemoveAll(selected.Contains);
            }
            var remainingBatchItems = batchItems.Concat(randomItems).Concat(multiBlockItems)
                .ToDictionary(item => item.Name);

            foreach (var name in names)
            {
                if (!remainingBatchItems.TryGetValue(name, out var item))
                    continue;
                cancellationToken.ThrowIfCancellationRequested();
                await ReadBatchItem(item, results, cancellationToken).ConfigureAwait(false);
            }

            // Preserve request order even though Random Read returns all words before dwords.
            var ordered = names.Where(results.ContainsKey)
                .ToDictionary(name => name, name => results[name]);
            return new DefaultPlcReadResponse(request, ordered);
        }

        private async Task ReadMultiBlockGroup(
            ReadItem[] items, Dictionary<string, PlcResponseItem<IPlcValue>> results,
            CancellationToken cancellationToken)
        {
            try
            {
                var words = items.Where(item => !item.Tag.IsBitDevice).ToArray();
                var bits = items.Where(item => item.Tag.IsBitDevice).ToArray();
                var data = new SlmpMultiBlockReadRequest(
                    checked((byte)words.Length), checked((byte)bits.Length),
                    words.Select(ToDeviceBlock).ToList(), bits.Select(ToDeviceBlock).ToList());
                var frame = new SlmpRequestFrame3E(_monitoringTimer, 0x0406, 0, data);
                var response = await ExchangeAsync(frame, cancellationToken).ConfigureAwait(false);
                if (response.EndCode != 0)
                {
                    Logger.LogWarning("SLMP returned end code 0x{EndCode:X4} for Multi-block Read", response.EndCode);
                    SetReadErrors(items, results, PlcResponseCode.RemoteError);
                    return;
                }
                var expectedBytes = checked(items.Sum(item => item.BlockPoints) * 2);
                if (response.ResponseData.Length != expectedBytes)
                {
                    SetReadErrors(items, results, PlcResponseCode.InvalidData);
                    return;
                }
                var offset = 0;
                foreach (var item in words.Concat(bits))
                {
                    var length = item.BlockPoints * 2;
                    var payload = response.ResponseData.AsSpan(offset, length).ToArray();
                    var value = item.Tag.IsBitDevice ? DecodeBitBlock(item.Tag, payload) : Decode(item.Tag, payload);
                    results[item.Name] = new DefaultPlcResponseItem<IPlcValue>(PlcResponseCode.Ok, value);
                    offset += length;
                }
            }
            catch (OperationCanceledException) { throw; }
            catch (TimeoutException ex)
            {
                Logger.LogWarning(ex, "SLMP Multi-block Read timed out");
                SetReadErrors(items, results, PlcResponseCode.RequestTimeout);
            }
            catch (Exception ex)
            {
                Logger.LogWarning(ex, "SLMP Multi-block Read failed");
                SetReadErrors(items, results, PlcResponseCode.InternalError);
            }
        }

        private static SlmpDeviceBlock ToDeviceBlock(ReadItem item) => new SlmpDeviceBlock(
            (uint)item.Tag.DeviceNumber, item.Tag.DeviceCode, checked((ushort)item.BlockPoints));

        private async Task ReadRandomGroup(
            ReadItem[] items, int wordCount, int doubleWordCount,
            Dictionary<string, PlcResponseItem<IPlcValue>> results,
            CancellationToken cancellationToken)
        {
            try
            {
                var words = items.Take(wordCount)
                    .Select(item => new SlmpDeviceSpec((uint)item.Tag.DeviceNumber, item.Tag.DeviceCode)).ToList();
                var doubleWords = items.Skip(wordCount).Take(doubleWordCount)
                    .Select(item => new SlmpDeviceSpec((uint)item.Tag.DeviceNumber, item.Tag.DeviceCode)).ToList();
                var data = new SlmpRandomReadRequest(
                    checked((byte)wordCount), checked((byte)doubleWordCount), words, doubleWords);
                var frame = new SlmpRequestFrame3E(_monitoringTimer, 0x0403, 0, data);
                var response = await ExchangeAsync(frame, cancellationToken).ConfigureAwait(false);
                if (response.EndCode != 0)
                {
                    Logger.LogWarning("SLMP returned end code 0x{EndCode:X4} for Random Read", response.EndCode);
                    SetReadErrors(items, results, PlcResponseCode.RemoteError);
                    return;
                }
                var expectedBytes = checked(wordCount * 2 + doubleWordCount * 4);
                if (response.ResponseData.Length != expectedBytes)
                {
                    SetReadErrors(items, results, PlcResponseCode.InvalidData);
                    return;
                }
                var offset = 0;
                foreach (var item in items)
                {
                    var length = item.Tag.NumberOfPoints * 2;
                    results[item.Name] = new DefaultPlcResponseItem<IPlcValue>(PlcResponseCode.Ok,
                        Decode(item.Tag, response.ResponseData.AsSpan(offset, length).ToArray()));
                    offset += length;
                }
            }
            catch (OperationCanceledException)
            {
                throw;
            }
            catch (TimeoutException ex)
            {
                Logger.LogWarning(ex, "SLMP Random Read timed out");
                SetReadErrors(items, results, PlcResponseCode.RequestTimeout);
            }
            catch (Exception ex)
            {
                Logger.LogWarning(ex, "SLMP Random Read failed");
                SetReadErrors(items, results, PlcResponseCode.InternalError);
            }
        }

        private async Task ReadBatchItem(
            ReadItem item, Dictionary<string, PlcResponseItem<IPlcValue>> results,
            CancellationToken cancellationToken)
        {
            var name = item.Name;
            var tag = item.Tag;
            try
            {
                var data = new SlmpReadRequest(
                    (uint)tag.DeviceNumber, tag.DeviceCode, checked((ushort)tag.NumberOfPoints));
                var subCommand = tag.IsBitDevice ? (ushort)1 : (ushort)0;
                var frame = new SlmpRequestFrame3E(_monitoringTimer, 0x0401, subCommand, data);
                var response = await ExchangeAsync(frame, cancellationToken).ConfigureAwait(false);
                if (response.EndCode != 0)
                {
                    Logger.LogWarning("SLMP returned end code 0x{EndCode:X4} for tag {TagName}",
                        response.EndCode, name);
                    results[name] = new DefaultPlcTagErrorItem<IPlcValue>(PlcResponseCode.RemoteError);
                    return;
                }
                var expectedBytes = tag.IsBitDevice ? (tag.NumberOfPoints + 1) / 2 : checked(tag.NumberOfPoints * 2);
                if (response.ResponseData.Length != expectedBytes ||
                    tag.IsBitDevice && !IsValidBitPayload(tag.NumberOfPoints, response.ResponseData))
                {
                    results[name] = new DefaultPlcTagErrorItem<IPlcValue>(PlcResponseCode.InvalidData);
                    return;
                }
                results[name] = new DefaultPlcResponseItem<IPlcValue>(
                    PlcResponseCode.Ok, Decode(tag, response.ResponseData));
            }
            catch (OperationCanceledException) { throw; }
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

        private static void SetReadErrors(IEnumerable<ReadItem> items,
            Dictionary<string, PlcResponseItem<IPlcValue>> results, PlcResponseCode code)
        {
            foreach (var item in items)
                results[item.Name] = new DefaultPlcTagErrorItem<IPlcValue>(code);
        }

        private sealed record ReadItem(string Name, SlmpTag Tag)
        {
            public int BlockPoints => Tag.IsBitDevice ? (Tag.Quantity + 15) / 16 : Tag.NumberOfPoints;
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
                    var subCommand = tag.IsBitDevice ? (ushort)1 : (ushort)0;
                    var frame = new SlmpRequestFrame3E(_monitoringTimer, 0x1401, subCommand, write);
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
                if (value is bool[] boolArray && tag.IsBitDevice)
                {
                    if (boolArray.Length != tag.Quantity)
                    {
                        return null;
                    }
                    values.AddRange(boolArray.Cast<object>());
                }
                else if (value is object[] array)
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
                if (value is Array || value is IPlcValue plcValue && plcValue.IsList())
                {
                    return null;
                }
                values.Add(value);
            }

            if (tag.IsBitDevice)
            {
                var packed = new byte[(values.Count + 1) / 2];
                for (var i = 0; i < values.Count; i++)
                {
                    var item = values[i];
                    bool bit;
                    if (item is IPlcValue plcValue && plcValue.IsBool())
                    {
                        bit = plcValue.GetBool();
                    }
                    else if (item is bool boolean)
                    {
                        bit = boolean;
                    }
                    else
                    {
                        return null;
                    }
                    if (bit)
                    {
                        // The first addressed point occupies the high nibble.
                        packed[i / 2] |= (byte)(1 << (i % 2 == 0 ? 4 : 0));
                    }
                }
                return packed;
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
            if (tag.IsBitDevice)
            {
                if (!tag.ExplicitRange)
                {
                    return new PlcBOOL(ReadBitPoint(data, 0));
                }

                var bitValues = new List<IPlcValue>(tag.Quantity);
                for (var i = 0; i < tag.Quantity; i++)
                {
                    bitValues.Add(new PlcBOOL(ReadBitPoint(data, i)));
                }
                return new PlcList(bitValues);
            }

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

        private static IPlcValue DecodeBitBlock(SlmpTag tag, byte[] data)
        {
            var values = new List<IPlcValue>(tag.Quantity);
            for (var i = 0; i < tag.Quantity; i++)
                values.Add(new PlcBOOL(((data[i / 8] >> (i % 8)) & 1) != 0));
            return new PlcList(values);
        }

        private static bool IsValidBitPayload(int numberOfPoints, byte[] data)
        {
            for (var i = 0; i < numberOfPoints; i++)
            {
                if (ReadBitNibble(data, i) > 1)
                {
                    return false;
                }
            }
            return numberOfPoints % 2 == 0 || (data[^1] & 0x0F) == 0;
        }

        private static bool ReadBitPoint(byte[] data, int point)
            => ReadBitNibble(data, point) == 1;

        private static int ReadBitNibble(byte[] data, int point)
            => (data[point / 2] >> (point % 2 == 0 ? 4 : 0)) & 0x0F;

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
                Close();
                _exchangeGate.Dispose();
            }
            base.Dispose(disposing);
        }
    }
}
