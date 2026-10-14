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
using System.Linq;
using System.Threading;
using System.Threading.Tasks;
using Microsoft.Extensions.Logging;
using org.apache.plc4net.api.metadata;
using org.apache.plc4net.api.value;
using org.apache.plc4net.drivers.s7.readwrite.model;
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
using org.apache.plc4net.transports.cotp;
using org.apache.plc4net.types;

namespace org.apache.plc4net.drivers.s7
{
    /// <summary>
    /// An S7 connection: S7comm over the ISO-on-TCP session that the COTP transport sets up.
    /// </summary>
    /// <remarks>
    /// <para>
    /// Every message, and every value, is encoded and decoded by a type generated from
    /// <c>s7.mspec</c>; nothing here handles bytes itself. A read or write of several tags is cut
    /// into jobs that each fit the PDU length the CPU granted, and the jobs go out one after the
    /// other: the connection keeps one request outstanding.
    /// </para>
    /// <para>
    /// Connection string parameters, with plc4j's names: <c>remote-rack</c> and <c>remote-slot</c>
    /// (0 and 0), <c>local-rack</c> and <c>local-slot</c> (1 and 1), <c>remote-device-group</c> and
    /// <c>local-device-group</c> (<c>PG_OR_PC</c>, <c>OS</c> or <c>OTHERS</c>; <c>PG_OR_PC</c> and
    /// <c>OTHERS</c>), <c>remote-tsap</c> and <c>local-tsap</c> (hexadecimal, replacing what the
    /// rack, slot and device group give), <c>pdu-size</c> (1024), <c>max-amq-caller</c> and
    /// <c>max-amq-callee</c> (8) and <c>read-timeout-ms</c> (10000).
    /// </para>
    /// </remarks>
    public sealed class S7Connection : ConnectionBase, PlcReader, PlcWriter
    {
        private static readonly DefaultPlcConnectionMetadata Metadata = new DefaultPlcConnectionMetadata
        {
            CanRead = true,
            CanWrite = true,
            CanSubscribe = false
        };

        private const int PollIntervalMs = 2;
        private const int DefaultReadTimeoutMs = 10000;
        private const int DefaultPduSize = 1024;
        private const int DefaultMaxAmq = 8;

        // The device groups of plc4j's DeviceGroup: PG_OR_PC = 0x01, OS = 0x02, OTHERS = 0x03.
        private const int DefaultLocalDeviceGroup = 0x03;
        private const int DefaultRemoteDeviceGroup = 0x01;

        // A TSAP holds the rack in four bits and the slot in four: a larger number cannot be
        // represented and is rejected rather than truncated.
        private const int MaxRackOrSlot = 0x0F;

        // What a job costs besides its items: the ten header bytes and the two of the parameter
        // (function code, item count); an acknowledgement has two more header bytes. An S7-ANY
        // item in the request is twelve bytes, a data item four bytes plus its data, which is
        // padded to an even length.
        private const int JobSize = 12;
        private const int AckSize = 14;
        private const int AddressItemSize = 12;
        private const int DataItemHeaderSize = 4;

        private readonly int _remoteRack;
        private readonly int _remoteSlot;
        private readonly int _localRack;
        private readonly int _localSlot;
        private readonly int _remoteDeviceGroup;
        private readonly int _localDeviceGroup;
        private readonly int _remoteTsap;
        private readonly int _localTsap;
        private readonly int _requestedPduSize;
        private readonly int _maxAmqCaller;
        private readonly int _maxAmqCallee;
        private readonly int _readTimeoutMs;

        // The CPU takes one job at a time from this connection; the gate keeps two callers from
        // interleaving their jobs on the wire.
        private readonly SemaphoreSlim _exchangeGate = new SemaphoreSlim(1, 1);
        private readonly S7MessageCodec _codec;

        private int _tpduReference;
        private ushort _expectedTpduReference;
        private S7Response? _response;

        public S7Connection(ConnectionString connectionString, ITransportInstance transport)
            : base(connectionString, transport)
        {
            // Record configuration only. The COTP handshake and Setup Communication happen in
            // OnConnect, which the driver invokes after construction: no I/O here.
            _remoteRack = NumberParameter(connectionString, "remote-rack", 0, 0, MaxRackOrSlot);
            _remoteSlot = NumberParameter(connectionString, "remote-slot", 0, 0, MaxRackOrSlot);
            _localRack = NumberParameter(connectionString, "local-rack", 1, 0, MaxRackOrSlot);
            _localSlot = NumberParameter(connectionString, "local-slot", 1, 0, MaxRackOrSlot);
            _remoteDeviceGroup = DeviceGroupParameter(connectionString, "remote-device-group", DefaultRemoteDeviceGroup);
            _localDeviceGroup = DeviceGroupParameter(connectionString, "local-device-group", DefaultLocalDeviceGroup);
            _remoteTsap = TsapParameter(connectionString, "remote-tsap");
            _localTsap = TsapParameter(connectionString, "local-tsap");
            _requestedPduSize = NumberParameter(connectionString, "pdu-size", DefaultPduSize, 1, ushort.MaxValue);
            _maxAmqCaller = NumberParameter(connectionString, "max-amq-caller", DefaultMaxAmq, 1, ushort.MaxValue);
            _maxAmqCallee = NumberParameter(connectionString, "max-amq-callee", DefaultMaxAmq, 1, ushort.MaxValue);
            _readTimeoutMs = NumberParameter(connectionString, "read-timeout-ms", DefaultReadTimeoutMs, 1, int.MaxValue);

            _codec = new S7MessageCodec(transport, OnMessage, _requestedPduSize);
        }

        /// <summary>The PDU length the CPU granted in Setup Communication, in bytes; 0 before that.</summary>
        public int NegotiatedPduLength { get; private set; }

        // ── Connection string ──────────────────────────────────

        /// <summary>
        /// Reads a whole-number parameter. A value that is present but not a number, or out of
        /// range, is an error rather than a silent default: a mistyped slot would address the
        /// wrong CPU.
        /// </summary>
        private static int NumberParameter(
            ConnectionString connectionString, string name, int defaultValue, int min, int max)
        {
            var raw = connectionString.GetParameter(name);
            if (raw == null)
            {
                return defaultValue;
            }
            if (!int.TryParse(raw, NumberStyles.Integer, CultureInfo.InvariantCulture, out var value))
            {
                throw new PlcConnectionException($"{name} must be a whole number, but was '{raw}'.");
            }
            if (value < min || value > max)
            {
                throw new PlcConnectionException($"{name} must be between {min} and {max}, but was {value}.");
            }
            return value;
        }

        /// <summary>A device group: plc4j's names (<c>PG_OR_PC</c>, <c>OS</c>, <c>OTHERS</c>) or a byte value.</summary>
        private static int DeviceGroupParameter(ConnectionString connectionString, string name, int defaultValue)
        {
            var raw = connectionString.GetParameter(name);
            if (raw == null)
            {
                return defaultValue;
            }
            switch (raw.Trim().ToUpperInvariant())
            {
                case "PG_OR_PC":
                case "PG":
                    return 0x01;
                case "OS":
                    return 0x02;
                case "OTHERS":
                case "OTHER":
                    return 0x03;
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
                $"{name} must be PG_OR_PC, OS, OTHERS or a byte value, but was '{raw}'.");
        }

        /// <summary>A TSAP, written in hexadecimal (<c>0x0300</c>); 0 when the parameter is absent.</summary>
        private static int TsapParameter(ConnectionString connectionString, string name)
        {
            var raw = connectionString.GetParameter(name);
            if (raw == null)
            {
                return 0;
            }
            var hex = raw.StartsWith("0x", StringComparison.OrdinalIgnoreCase) ? raw.Substring(2) : raw;
            if (!int.TryParse(hex, NumberStyles.HexNumber, CultureInfo.InvariantCulture, out var value)
                || value < 0 || value > 0xFFFF)
            {
                throw new PlcConnectionException($"{name} must be a hexadecimal TSAP such as 0x0300, but was '{raw}'.");
            }
            return value;
        }

        /// <summary>
        /// The TSAP of plc4j's S7TsapIdEncoder: the device group in the high byte, the rack in
        /// the high nibble of the low byte and the slot in the low nibble.
        /// </summary>
        internal static int EncodeTsap(int deviceGroup, int rack, int slot)
            => (deviceGroup << 8) | (rack << 4) | (slot & 0x0F);

        // ── Connect ────────────────────────────────────────────

        protected override void OnConnect()
        {
            // S7 needs the COTP session: its handshake is part of the protocol, and nothing
            // works without it.
            if (!(TransportInstance is CotpTransportInstance cotp))
            {
                throw new PlcConnectionException(
                    "The S7 driver requires the 'cotp' transport, but got '" +
                    TransportInstance.GetType().Name + "'.");
            }

            // An explicit TSAP replaces the one that the rack, slot and device group give.
            var localTsap = _localTsap != 0 ? _localTsap : EncodeTsap(_localDeviceGroup, _localRack, _localSlot);
            var remoteTsap = _remoteTsap != 0 ? _remoteTsap : EncodeTsap(_remoteDeviceGroup, _remoteRack, _remoteSlot);

            cotp.Open(
                localTsapHi: (byte)((localTsap >> 8) & 0xFF),
                localTsapLo: (byte)(localTsap & 0xFF),
                remoteTsapHi: (byte)((remoteTsap >> 8) & 0xFF),
                remoteTsapLo: (byte)(remoteTsap & 0xFF));

            // A CPU ignores Read Var and Write Var until Setup Communication has run.
            PerformSetupCommunication();

            // Say why the link went away when it does, instead of leaving the caller with a bare
            // timeout on its next request.
            if (TransportInstance is IAsyncTransportInstance asyncTransport)
            {
                asyncTransport.RegisterDisconnectListener(ex => Logger.LogWarning(ex, "The S7 transport disconnected"));
            }
        }

        private void PerformSetupCommunication()
        {
            var reference = NextTpduReference();
            var request = new S7MessageRequest(
                reference,
                new S7ParameterSetupCommunication(
                    (ushort)_maxAmqCaller, (ushort)_maxAmqCallee, (ushort)_requestedPduSize),
                null);

            S7Response response;
            try
            {
                response = ExchangeAsync(request, reference, CancellationToken.None).GetAwaiter().GetResult();
            }
            catch (TimeoutException ex)
            {
                throw new PlcConnectionException(
                    "No answer to S7 Setup Communication: the CPU may not permit PUT/GET access, " +
                    "or the rack, slot or TSAP is wrong.", ex);
            }
            catch (S7DriverException ex)
            {
                throw new PlcConnectionException(ex.Message, ex);
            }

            if (response.HasError)
            {
                throw new PlcConnectionException(
                    $"S7 Setup Communication was refused: error 0x{response.ErrorClass:X2}{response.ErrorCode:X2}"
                    + DescribeHeaderError(response.ErrorClass, response.ErrorCode) + ".");
            }
            if (!(response.Parameter is S7ParameterSetupCommunication setup))
            {
                throw new PlcConnectionException("The answer to S7 Setup Communication carried no setup parameter.");
            }

            if (setup.PduLength == 0 || setup.PduLength > _requestedPduSize)
            {
                throw new PlcConnectionException(
                    $"S7 Setup Communication returned an invalid PDU length {setup.PduLength}; " +
                    $"the requested maximum was {_requestedPduSize}.");
            }

            NegotiatedPduLength = setup.PduLength;
            _codec.MaximumMessageSize = NegotiatedPduLength;
            Logger.LogInformation("S7 Setup Communication done, PDU length {PduLength} bytes", NegotiatedPduLength);
        }

        public override IPlcConnectionMetadata PlcConnectionMetadata => Metadata;

        public override IPlcTag Parse(string tagQuery) => S7Tag.Parse(tagQuery);

        public override IPlcReadRequestBuilder? ReadRequestBuilder
            => new DefaultPlcReadRequestBuilder(this, Parse);

        public override IPlcWriteRequestBuilder? WriteRequestBuilder
            => new DefaultPlcWriteRequestBuilder(this, Parse);

        public override IPlcSubscriptionRequestBuilder? SubscriptionRequestBuilder => null;

        public override IPlcUnsubscriptionRequestBuilder? UnsubscriptionRequestBuilder => null;

        // ── Exchange ───────────────────────────────────────────

        /// <summary>The parts of an acknowledgement that the driver acts on.</summary>
        private sealed class S7Response
        {
            public S7Response(byte errorClass, byte errorCode, S7Parameter? parameter, S7Payload? payload)
            {
                ErrorClass = errorClass;
                ErrorCode = errorCode;
                Parameter = parameter;
                Payload = payload;
            }

            public byte ErrorClass { get; }
            public byte ErrorCode { get; }
            public S7Parameter? Parameter { get; }
            public S7Payload? Payload { get; }

            /// <summary>The CPU refused the whole job: the error class and code of the header.</summary>
            public bool HasError => ErrorClass != 0x00 || ErrorCode != 0x00;
        }

        private ushort NextTpduReference() => (ushort)(Interlocked.Increment(ref _tpduReference) & 0xFFFF);

        private void OnMessage(S7Message message)
        {
            ushort reference;
            S7Response response;
            switch (message)
            {
                // A CPU answers a job with an acknowledgement with data, or - when it refuses the
                // job before it looks at the items - with a bare acknowledgement. Both carry the
                // error class and code first.
                case S7MessageResponseData data:
                    reference = data.TpduReference;
                    response = new S7Response(data.ErrorClass, data.ErrorCode, data.Parameter, data.Payload);
                    break;
                case S7MessageResponse ack:
                    reference = ack.TpduReference;
                    response = new S7Response(ack.ErrorClass, ack.ErrorCode, ack.Parameter, ack.Payload);
                    break;
                default:
                    Logger.LogDebug("S7: ignored a {Type}, which no job of ours is waiting for", message.GetType().Name);
                    return;
            }

            // The answer to an earlier job that timed out arrives late, with that job's reference.
            if (reference != _expectedTpduReference)
            {
                Logger.LogDebug(
                    "S7: discarded a stale answer (TPDU reference {Got}, waiting for {Want})",
                    reference, _expectedTpduReference);
                return;
            }
            _response = response;
        }

        /// <summary>
        /// Sends a job and returns the acknowledgement that carries its TPDU reference. Throws
        /// <see cref="TimeoutException"/> when none arrives within the read timeout.
        /// </summary>
        private async Task<S7Response> ExchangeAsync(
            S7Message request, ushort reference, CancellationToken cancellationToken)
        {
            await _exchangeGate.WaitAsync(cancellationToken).ConfigureAwait(false);
            try
            {
                _expectedTpduReference = reference;
                _response = null;

                _codec.Send(request);
                var deadline = Environment.TickCount64 + _readTimeoutMs;

                while (true)
                {
                    _codec.ProcessIncomingData();
                    if (_response != null)
                    {
                        return _response;
                    }

                    if (!TransportInstance.IsOpen)
                    {
                        throw new S7DriverException("The S7 connection was closed.");
                    }
                    if (Environment.TickCount64 >= deadline)
                    {
                        throw new TimeoutException(
                            $"No S7 answer to TPDU reference {reference} within {_readTimeoutMs} ms.");
                    }
                    await Task.Delay(PollIntervalMs, cancellationToken).ConfigureAwait(false);
                }
            }
            catch (MessageCodecException ex)
            {
                DiscardBufferedBytes();
                throw new S7DriverException($"S7 exchange failed: {ex.Message}", ex);
            }
            catch (TimeoutException)
            {
                // Whatever is half-received now would be misread as the start of the next answer.
                DiscardBufferedBytes();
                throw;
            }
            finally
            {
                _response = null;
                _exchangeGate.Release();
            }
        }

        private void DiscardBufferedBytes()
        {
            try
            {
                var available = TransportInstance.GetNumBytesAvailable();
                if (available > 0)
                {
                    TransportInstance.Read(available);
                }
            }
            catch (TransportException)
            {
                // The transport is gone; there is nothing left to discard, and this must not hide
                // the failure that led here.
            }
        }

        protected override void Dispose(bool disposing)
        {
            if (disposing)
            {
                _exchangeGate.Dispose();
            }
            base.Dispose(disposing);
        }

        // ── Jobs ───────────────────────────────────────────────

        /// <summary>
        /// Cuts a list of tags into jobs that each fit one PDU, as lists of indexes into
        /// <paramref name="tags"/>: the request has to fit, and so does the answer. A write needs
        /// room for the data in the request, a read for the data in the answer.
        /// </summary>
        internal static List<List<int>> PlanJobs(IReadOnlyList<S7Tag> tags, bool write, int pduLength)
        {
            var jobs = new List<List<int>>();
            var job = new List<int>();
            var requestSize = JobSize;
            var answerSize = AckSize;

            for (var i = 0; i < tags.Count; i++)
            {
                // The generated model pads an odd-sized item only when a following item makes it
                // non-final. The prospective item therefore adds padding for the current last
                // item, while it is itself counted without padding until another item follows.
                var previousPadding = job.Count > 0 && tags[job[job.Count - 1]].SizeInBytes % 2 != 0 ? 1 : 0;
                var data = DataItemHeaderSize + tags[i].SizeInBytes;
                var requestItem = write ? AddressItemSize + data + previousPadding : AddressItemSize;
                var answerItem = write ? 1 : data + previousPadding;

                if (job.Count > 0
                    && (requestSize + requestItem > pduLength
                        || answerSize + answerItem > pduLength))
                {
                    jobs.Add(job);
                    job = new List<int>();
                    requestSize = JobSize;
                    answerSize = AckSize;
                }

                job.Add(i);
                requestSize += requestItem;
                answerSize += answerItem;
            }

            if (job.Count > 0)
            {
                jobs.Add(job);
            }
            return jobs;
        }

        private int PduLength => NegotiatedPduLength > 0 ? NegotiatedPduLength : _requestedPduSize;

        /// <summary>
        /// The S7-ANY address of a tag, built the way plc4j builds it: the tag's own data type
        /// names the transport size, and the CPU is asked for one element.
        /// </summary>
        internal static S7AddressAny ToAddress(S7Tag tag)
            => new S7AddressAny(
                tag.DataType,
                1,
                (ushort)tag.DbNumber,
                tag.Area,
                (ushort)tag.ByteOffset,
                tag.BitOffset);

        // ── PlcReader ──────────────────────────────────────────

        public async Task<IPlcReadResponse> Read(
            DefaultPlcReadRequest request, CancellationToken cancellationToken = default)
        {
            var results = new Dictionary<string, PlcResponseItem<IPlcValue>>();
            var work = new List<(string Name, S7Tag Tag)>();

            foreach (var name in request.TagNames)
            {
                if (request.GetTagByName(name) is S7Tag tag)
                {
                    work.Add((name, tag));
                }
                else
                {
                    results[name] = new DefaultPlcTagErrorItem<IPlcValue>(PlcResponseCode.InvalidAddress);
                }
            }

            foreach (var job in PlanJobs(work.Select(w => w.Tag).ToList(), false, PduLength))
            {
                cancellationToken.ThrowIfCancellationRequested();

                var items = job.Select(i => work[i]).ToList();
                try
                {
                    await ReadJobAsync(items, results, cancellationToken).ConfigureAwait(false);
                }
                catch (OperationCanceledException)
                {
                    throw; // Cancellation is the caller's decision, not a per-tag error.
                }
                catch (Exception ex)
                {
                    Logger.LogWarning(ex, "S7 read of {Count} tag(s) failed", items.Count);
                    var code = MapException(ex);
                    foreach (var (name, _) in items)
                    {
                        results[name] = new DefaultPlcTagErrorItem<IPlcValue>(code);
                    }
                }
            }

            return new DefaultPlcReadResponse(request, results);
        }

        private async Task ReadJobAsync(
            List<(string Name, S7Tag Tag)> items,
            Dictionary<string, PlcResponseItem<IPlcValue>> results,
            CancellationToken cancellationToken)
        {
            var reference = NextTpduReference();
            var request = new S7MessageRequest(
                reference,
                new S7ParameterReadVarRequest(items
                    .Select(i => (S7VarRequestParameterItem)new S7VarRequestParameterItemAddress(ToAddress(i.Tag)))
                    .ToList()),
                null);

            var response = await ExchangeAsync(request, reference, cancellationToken).ConfigureAwait(false);

            if (response.HasError)
            {
                LogRefusal("read", response);
                var refused = MapHeaderError(response.ErrorClass, response.ErrorCode);
                foreach (var (name, _) in items)
                {
                    results[name] = new DefaultPlcTagErrorItem<IPlcValue>(refused);
                }
                return;
            }

            if (!(response.Payload is S7PayloadReadVarResponse payload))
            {
                throw new S7DriverException(
                    $"The answer to a read carries {response.Payload?.GetType().Name ?? "no payload"}.");
            }

            for (var i = 0; i < items.Count; i++)
            {
                var (name, tag) = items[i];
                if (i >= payload.Items.Count)
                {
                    results[name] = new DefaultPlcTagErrorItem<IPlcValue>(PlcResponseCode.InternalError);
                    continue;
                }

                var item = payload.Items[i];
                if (item.ReturnCode != DataTransportErrorCode.OK)
                {
                    results[name] = new DefaultPlcTagErrorItem<IPlcValue>(MapReturnCode(item.ReturnCode));
                    continue;
                }

                try
                {
                    results[name] = new DefaultPlcResponseItem<IPlcValue>(PlcResponseCode.Ok, ParseValue(tag, item.Data));
                }
                catch (Exception ex) when (ex is S7DriverException || ex is ParseException)
                {
                    Logger.LogWarning(ex, "S7 read of tag '{TagName}' returned a value that cannot be read", name);
                    results[name] = new DefaultPlcTagErrorItem<IPlcValue>(PlcResponseCode.InternalError);
                }
            }
        }

        /// <summary>
        /// Decodes the first element of an item's data with the generated <c>DataItem</c>; a CPU
        /// may answer with more than the one element that was asked for.
        /// </summary>
        private static IPlcValue ParseValue(S7Tag tag, byte[] data)
        {
            if (data.Length < tag.SizeInBytes)
            {
                throw new S7DriverException(
                    $"The S7 answer carries {data.Length} byte(s) for {tag}; a {tag.DataType} takes {tag.SizeInBytes}.");
            }
            return DataItem.StaticParse(new ReadBuffer(data), tag.DataType.GetDataProtocolId(), ControllerType.ANY, 0);
        }

        // ── PlcWriter ──────────────────────────────────────────

        public async Task<IPlcWriteResponse> Write(
            DefaultPlcWriteRequest request, CancellationToken cancellationToken = default)
        {
            var codes = new Dictionary<string, PlcResponseCode>();
            var work = new List<(string Name, S7Tag Tag, byte[] Data)>();

            foreach (var name in request.TagNames)
            {
                var value = request.GetValue(name);
                if (!(request.GetTagByName(name) is S7Tag tag))
                {
                    codes[name] = PlcResponseCode.InvalidAddress;
                    continue;
                }
                if (value == null)
                {
                    codes[name] = PlcResponseCode.InvalidDatatype;
                    continue;
                }

                // Encoded before anything is sent: a value that does not fit the tag's type never
                // reaches the CPU.
                try
                {
                    work.Add((name, tag, EncodeValue(tag, value)));
                }
                catch (Exception ex) when (ex is FormatException || ex is InvalidCastException || ex is OverflowException)
                {
                    Logger.LogWarning(ex, "S7 write of tag '{TagName}': the value does not fit {Type}", name, tag.DataType);
                    codes[name] = PlcResponseCode.InvalidDatatype;
                }
            }

            foreach (var job in PlanJobs(work.Select(w => w.Tag).ToList(), true, PduLength))
            {
                cancellationToken.ThrowIfCancellationRequested();

                var items = job.Select(i => work[i]).ToList();
                try
                {
                    await WriteJobAsync(items, codes, cancellationToken).ConfigureAwait(false);
                }
                catch (OperationCanceledException)
                {
                    throw; // Cancellation is the caller's decision, not a per-tag error.
                }
                catch (Exception ex)
                {
                    Logger.LogWarning(ex, "S7 write of {Count} tag(s) failed", items.Count);
                    var code = MapException(ex);
                    foreach (var (name, _, _) in items)
                    {
                        codes[name] = code;
                    }
                }
            }

            return new DefaultPlcWriteResponse(request, codes);
        }

        private async Task WriteJobAsync(
            List<(string Name, S7Tag Tag, byte[] Data)> items,
            Dictionary<string, PlcResponseCode> codes,
            CancellationToken cancellationToken)
        {
            var reference = NextTpduReference();
            var request = new S7MessageRequest(
                reference,
                new S7ParameterWriteVarRequest(items
                    .Select(i => (S7VarRequestParameterItem)new S7VarRequestParameterItemAddress(ToAddress(i.Tag)))
                    .ToList()),
                new S7PayloadWriteVarRequest(items
                    .Select(i => new S7VarPayloadDataItem(
                        DataTransportErrorCode.OK, i.Tag.DataType.GetDataTransportSize(), i.Data))
                    .ToList()));

            var response = await ExchangeAsync(request, reference, cancellationToken).ConfigureAwait(false);

            if (response.HasError)
            {
                LogRefusal("write", response);
                var refused = MapHeaderError(response.ErrorClass, response.ErrorCode);
                foreach (var (name, _, _) in items)
                {
                    codes[name] = refused;
                }
                return;
            }

            if (!(response.Payload is S7PayloadWriteVarResponse payload))
            {
                throw new S7DriverException(
                    $"The answer to a write carries {response.Payload?.GetType().Name ?? "no payload"}.");
            }

            for (var i = 0; i < items.Count; i++)
            {
                codes[items[i].Name] = i >= payload.Items.Count
                    ? PlcResponseCode.InternalError
                    : payload.Items[i].ReturnCode == DataTransportErrorCode.OK
                        ? PlcResponseCode.Ok
                        : MapReturnCode(payload.Items[i].ReturnCode);
            }
        }

        /// <summary>
        /// Encodes a value for a tag with the generated <c>DataItem</c>, so the value has to fit
        /// the tag's data type: an <c>INT</c> tag takes -32768..32767, a <c>BYTE</c> tag 0..255.
        /// </summary>
        /// <exception cref="OverflowException">The value does not fit the data type.</exception>
        /// <exception cref="FormatException">The value is not a whole number where one is needed.</exception>
        /// <exception cref="InvalidCastException">The value is not of a kind the data type takes.</exception>
        internal static byte[] EncodeValue(S7Tag tag, object value)
        {
            IPlcValue plcValue;
            switch (tag.DataType)
            {
                case TransportSize.BOOL:
                    plcValue = new PlcBOOL(ToBoolean(value));
                    break;
                case TransportSize.REAL:
                    plcValue = new PlcREAL(ToSingle(value));
                    break;
                case TransportSize.BYTE:
                    plcValue = new PlcBYTE(checked((byte)ToInteger(value)));
                    break;
                case TransportSize.SINT:
                    plcValue = new PlcSINT(checked((sbyte)ToInteger(value)));
                    break;
                case TransportSize.USINT:
                    plcValue = new PlcUSINT(checked((byte)ToInteger(value)));
                    break;
                case TransportSize.WORD:
                    plcValue = new PlcWORD(checked((ushort)ToInteger(value)));
                    break;
                case TransportSize.INT:
                    plcValue = new PlcINT(checked((short)ToInteger(value)));
                    break;
                case TransportSize.UINT:
                    plcValue = new PlcUINT(checked((ushort)ToInteger(value)));
                    break;
                case TransportSize.DWORD:
                    plcValue = new PlcDWORD(checked((uint)ToInteger(value)));
                    break;
                case TransportSize.DINT:
                    plcValue = new PlcDINT(checked((int)ToInteger(value)));
                    break;
                case TransportSize.UDINT:
                    plcValue = new PlcUDINT(checked((uint)ToInteger(value)));
                    break;
                default:
                    throw new S7DriverException($"A tag cannot be written as {tag.DataType}.");
            }

            var buffer = new WriteBuffer();
            DataItem.StaticSerialize(buffer, plcValue, tag.DataType.GetDataProtocolId(), ControllerType.ANY, 0);
            return buffer.GetBytes();
        }

        private static bool ToBoolean(object value)
        {
            switch (value)
            {
                case bool b:
                    return b;
                case IPlcValue plc:
                    return plc.GetBool();
                default:
                    return Convert.ToBoolean(value, CultureInfo.InvariantCulture);
            }
        }

        private static long ToInteger(object value)
        {
            switch (value)
            {
                case bool _:
                    throw new InvalidCastException("A boolean cannot be written to a number.");
                case IPlcValue plc:
                    return plc.GetLong();
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

        private static float ToSingle(object value)
        {
            switch (value)
            {
                case bool _:
                    throw new InvalidCastException("A boolean cannot be written to a REAL.");
                case float f:
                    return f;
                case IPlcValue plc:
                    return plc.GetFloat();
                case string s:
                    return float.Parse(s, NumberStyles.Float, CultureInfo.InvariantCulture);
                default:
                    var number = Convert.ToDouble(value, CultureInfo.InvariantCulture);
                    var single = (float)number;
                    if (float.IsInfinity(single) && !double.IsInfinity(number))
                    {
                        throw new OverflowException($"{number.ToString(CultureInfo.InvariantCulture)} does not fit a REAL.");
                    }
                    return single;
            }
        }

        // ── Response codes ─────────────────────────────────────

        private void LogRefusal(string what, S7Response response)
            => Logger.LogWarning(
                "The S7 CPU refused the {What}: error 0x{ErrorClass:X2}{ErrorCode:X2}{Hint}",
                what, response.ErrorClass, response.ErrorCode,
                DescribeHeaderError(response.ErrorClass, response.ErrorCode));

        /// <summary>
        /// Maps the error class and code of a refused job to a response code. Class 0x81 code 0x04
        /// is a CPU that refuses PUT/GET access; an S7-300 says so with 0x83 0x04; class 0x85 is
        /// an access error.
        /// </summary>
        internal static PlcResponseCode MapHeaderError(byte errorClass, byte errorCode)
        {
            if ((errorClass == 0x81 || errorClass == 0x83) && errorCode == 0x04)
            {
                return PlcResponseCode.AccessDenied;
            }
            if (errorClass == 0x85)
            {
                return PlcResponseCode.AccessDenied;
            }
            return PlcResponseCode.InternalError;
        }

        /// <summary>A hint for the errors that a CPU which is not set up for S7 access most often gives.</summary>
        private static string DescribeHeaderError(byte errorClass, byte errorCode)
        {
            if ((errorClass == 0x81 || errorClass == 0x83) && errorCode == 0x04)
            {
                return " (the CPU refused access: tick \"Permit access with PUT/GET communication " +
                       "from remote partner\" in the CPU's protection settings, compile and download the " +
                       "hardware configuration, and untick \"Optimized block access\" for the data block)";
            }
            if (errorClass == 0x85)
            {
                return " (access error: the address is out of range or the block does not exist)";
            }
            return string.Empty;
        }

        /// <summary>Maps the return code of one item to a response code.</summary>
        internal static PlcResponseCode MapReturnCode(DataTransportErrorCode code)
        {
            switch (code)
            {
                case DataTransportErrorCode.ACCESS_DENIED:
                    return PlcResponseCode.AccessDenied;
                case DataTransportErrorCode.INVALID_ADDRESS:
                    return PlcResponseCode.InvalidAddress;
                case DataTransportErrorCode.DATA_TYPE_NOT_SUPPORTED:
                    return PlcResponseCode.InvalidDatatype;
                case DataTransportErrorCode.NOT_FOUND:
                    return PlcResponseCode.NotFound;
                default:
                    return PlcResponseCode.InternalError;
            }
        }

        /// <summary>Maps a failed exchange to a response code.</summary>
        internal static PlcResponseCode MapException(Exception exception)
        {
            switch (exception)
            {
                case TimeoutException _:
                    return PlcResponseCode.RequestTimeout;
                case FormatException _:
                case InvalidCastException _:
                case OverflowException _:
                    return PlcResponseCode.InvalidDatatype;
                default:
                    return PlcResponseCode.InternalError;
            }
        }
    }
}
