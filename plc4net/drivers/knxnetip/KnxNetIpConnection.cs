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
using System.Collections.Concurrent;
using System.Collections.Generic;
using System.Globalization;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;
using Microsoft.Extensions.Logging;
using org.apache.plc4net.api.metadata;
using org.apache.plc4net.api.value;
using org.apache.plc4net.drivers.knxnetip.readwrite.model;
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
using org.apache.plc4net.transports.udp;
using org.apache.plc4net.types;

namespace org.apache.plc4net.drivers.knxnetip
{
    /// <summary>
    /// A KNXnet/IP tunnelling connection.
    /// </summary>
    /// <remarks>
    /// <para>Handshake, mirroring the Java driver: SEARCH_REQUEST → SEARCH_RESPONSE
    /// (learn the gateway), CONNECT_REQUEST → CONNECT_RESPONSE (get a channel id and
    /// an assigned client KNX address), then a CONNECTIONSTATE_REQUEST heartbeat every
    /// 60 s. <see cref="Close"/> sends a DISCONNECT_REQUEST.</para>
    /// <para>Read sends a <c>GroupValueRead</c> and waits for the matching
    /// <c>GroupValueResponse</c> some device on the bus sends back. Write sends a
    /// <c>GroupValueWrite</c> and waits for the gateway's tunnelling ack. Decoding /
    /// encoding needs a datapoint type - supply it as a <c>:DPT…</c> suffix on the
    /// tag; without one only raw byte payloads work.</para>
    /// <para>The plc4net SPI has no subscription path, so bus monitoring is a plain
    /// callback: <see cref="RegisterGroupValueListener"/>.</para>
    /// </remarks>
    public sealed class KnxNetIpConnection : ConnectionBase, PlcReader, PlcWriter
    {
        private const long HeartbeatIntervalMs = 60_000;

        private readonly int _groupAddressLevels;
        private readonly int _requestTimeoutMs;

        private KnxNetIpMessageCodec _codec;
        private IAsyncTransportInstance _asyncTransport;
        private Timer _heartbeat;

        private byte _communicationChannelId;
        private KnxAddress _clientKnxAddress = new KnxAddress(0, 0, 0);
        private int _sequenceCounter;

        // Handshake / heartbeat slots - one outstanding at a time.
        private volatile TaskCompletionSource<SearchResponse> _pendingSearch;
        private volatile TaskCompletionSource<ConnectionResponse> _pendingConnect;
        private volatile TaskCompletionSource<ConnectionStateResponse> _pendingConnState;

        // Correlated by tunnelling sequence counter / destination group address.
        private readonly ConcurrentDictionary<byte, TaskCompletionSource<TunnelingResponse>> _pendingAcks
            = new ConcurrentDictionary<byte, TaskCompletionSource<TunnelingResponse>>();
        private readonly ConcurrentDictionary<string, TaskCompletionSource<byte[]>> _pendingReads
            = new ConcurrentDictionary<string, TaskCompletionSource<byte[]>>();

        // A registered :DPT… hint per group address, so inbound telegrams decode.
        private readonly ConcurrentDictionary<string, KnxDatapointType> _dptHints
            = new ConcurrentDictionary<string, KnxDatapointType>();
        private readonly List<Action<KnxGroupValueEvent>> _listeners = new List<Action<KnxGroupValueEvent>>();
        private readonly object _listenerLock = new object();

        public KnxNetIpConnection(ConnectionString connectionString, ITransportInstance transport)
            : base(connectionString, transport)
        {
            _groupAddressLevels = connectionString.GetIntParameter("knx-group-address-num-levels", 3);
            if (_groupAddressLevels is < 1 or > 3)
            {
                throw new PlcConnectionException(
                    $"knx-group-address-num-levels must be 1, 2 or 3, but was {_groupAddressLevels}.");
            }
            _requestTimeoutMs = connectionString.GetIntParameter("request-timeout", 3000);
        }

        /// <summary>The KNX address the gateway assigned this client, e.g. <c>1.1.255</c>.</summary>
        public string ClientKnxAddress =>
            $"{_clientKnxAddress.MainGroup}.{_clientKnxAddress.MiddleGroup}.{_clientKnxAddress.SubGroup}";

        public override IPlcConnectionMetadata PlcConnectionMetadata => new DefaultPlcConnectionMetadata
        {
            CanRead = true,
            CanWrite = true,
            CanSubscribe = false,
        };

        public override IPlcTag Parse(string tagQuery) => KnxNetIpTag.Parse(tagQuery);

        public override IPlcReadRequestBuilder ReadRequestBuilder
            => new DefaultPlcReadRequestBuilder(this, Parse);

        public override IPlcWriteRequestBuilder WriteRequestBuilder
            => new DefaultPlcWriteRequestBuilder(this, Parse);

        public override IPlcSubscriptionRequestBuilder SubscriptionRequestBuilder => null;

        public override IPlcUnsubscriptionRequestBuilder UnsubscriptionRequestBuilder => null;

        // ── handshake ──────────────────────────────────────────

        protected override void OnConnect()
        {
            if (!(TransportInstance is IAsyncTransportInstance asyncTransport))
            {
                throw new PlcConnectionException(
                    "The KNXnet/IP driver needs a transport that pushes inbound data (the 'udp' transport).");
            }
            _asyncTransport = asyncTransport;

            var local = ResolveLocalEndpoint();
            _codec = new KnxNetIpMessageCodec(TransportInstance, HandleIncomingMessage);
            asyncTransport.RegisterDataListener(() =>
            {
                try
                {
                    _codec.ProcessIncomingData();
                }
                catch (Exception e)
                {
                    Logger.LogWarning(e, "Error processing inbound KNXnet/IP data");
                }
            });

            var hpai = new HPAIDiscoveryEndpoint(HostProtocolCode.IPV4_UDP, local.Address, local.Port);

            // SEARCH: identify the gateway.
            _pendingSearch = new TaskCompletionSource<SearchResponse>(TaskCreationOptions.RunContinuationsAsynchronously);
            Send(new SearchRequest(hpai));
            Await(_pendingSearch.Task, "SEARCH_RESPONSE");

            // CONNECT: open a tunnelling channel.
            _pendingConnect = new TaskCompletionSource<ConnectionResponse>(TaskCreationOptions.RunContinuationsAsynchronously);
            Send(new ConnectionRequest(
                new HPAIDiscoveryEndpoint(HostProtocolCode.IPV4_UDP, local.Address, local.Port),
                new HPAIDataEndpoint(HostProtocolCode.IPV4_UDP, local.Address, local.Port),
                new ConnectionRequestInformationTunnelConnection(KnxLayer.TUNNEL_LINK_LAYER)));
            var connectionResponse = Await(_pendingConnect.Task, "CONNECT_RESPONSE");

            if (connectionResponse.Status != Status.NO_ERROR)
            {
                throw new PlcConnectionException(
                    $"KNXnet/IP gateway refused the connection: {connectionResponse.Status}.");
            }
            _communicationChannelId = connectionResponse.CommunicationChannelId;
            if (connectionResponse.ConnectionResponseDataBlock is ConnectionResponseDataBlockTunnelConnection block)
            {
                _clientKnxAddress = block.KnxAddress;
            }
            Logger.LogInformation(
                "KNXnet/IP tunnelling connection established (channel {Channel}, client address {Address})",
                _communicationChannelId, ClientKnxAddress);

            _heartbeat = new Timer(_ => SendHeartbeat(), null, HeartbeatIntervalMs, HeartbeatIntervalMs);
        }

        private (IPAddress Address, ushort Port) ResolveLocalEndpoint()
        {
            if (!(TransportInstance is UdpTransportInstance udp) || udp.LocalAddress == null)
            {
                throw new PlcConnectionException("The UDP transport has no bound local endpoint.");
            }
            var addr = udp.LocalAddress.Address;
            // A socket bound to 0.0.0.0 cannot be handed back to the gateway; fall
            // back to the loopback-safe "same subnet as the gateway" address the OS
            // would route through. For the common single-homed case the bound
            // address is already the right one once we resolve Any -> a real NIC.
            if (Equals(addr, System.Net.IPAddress.Any))
            {
                addr = LocalRoutableAddress(udp.RemoteAddress.Address) ?? System.Net.IPAddress.Loopback;
            }
            return (new IPAddress(addr.GetAddressBytes()), (ushort) udp.LocalAddress.Port);
        }

        private static System.Net.IPAddress LocalRoutableAddress(System.Net.IPAddress remote)
        {
            try
            {
                using var probe = new System.Net.Sockets.Socket(
                    remote.AddressFamily, System.Net.Sockets.SocketType.Dgram, System.Net.Sockets.ProtocolType.Udp);
                probe.Connect(remote, 9);
                return (probe.LocalEndPoint as System.Net.IPEndPoint)?.Address;
            }
            catch
            {
                return null;
            }
        }

        private void SendHeartbeat()
        {
            try
            {
                var local = ResolveLocalEndpoint();
                _pendingConnState = new TaskCompletionSource<ConnectionStateResponse>(
                    TaskCreationOptions.RunContinuationsAsynchronously);
                Send(new ConnectionStateRequest(_communicationChannelId,
                    new HPAIControlEndpoint(HostProtocolCode.IPV4_UDP, local.Address, local.Port)));
                var response = Await(_pendingConnState.Task, "CONNECTIONSTATE_RESPONSE");
                if (response.Status != Status.NO_ERROR)
                {
                    Logger.LogWarning("KNXnet/IP heartbeat returned status {Status}", response.Status);
                }
            }
            catch (Exception e)
            {
                Logger.LogWarning(e, "KNXnet/IP heartbeat failed");
            }
        }

        // ── inbound dispatch ───────────────────────────────────

        private void HandleIncomingMessage(KnxNetIpMessage message)
        {
            switch (message)
            {
                case SearchResponse searchResponse:
                    _pendingSearch?.TrySetResult(searchResponse);
                    break;
                case ConnectionResponse connectionResponse:
                    _pendingConnect?.TrySetResult(connectionResponse);
                    break;
                case ConnectionStateResponse connectionStateResponse:
                    _pendingConnState?.TrySetResult(connectionStateResponse);
                    break;
                case TunnelingResponse tunnelingResponse:
                    if (_pendingAcks.TryRemove(
                            tunnelingResponse.TunnelingResponseDataBlock.SequenceCounter, out var ack))
                    {
                        ack.TrySetResult(tunnelingResponse);
                    }
                    break;
                case TunnelingRequest tunnelingRequest:
                    HandleTunnelingRequest(tunnelingRequest);
                    break;
                case DisconnectRequest _:
                    Logger.LogInformation("KNXnet/IP gateway sent DISCONNECT_REQUEST");
                    break;
            }
        }

        private void HandleTunnelingRequest(TunnelingRequest tunnelingRequest)
        {
            // Ack first: a missed ack makes the gateway resend the same indication.
            var seq = ExtractSequenceCounter(tunnelingRequest);
            try
            {
                Send(new TunnelingResponse(
                    new TunnelingResponseDataBlock(_communicationChannelId, seq, Status.NO_ERROR)));
            }
            catch (Exception e)
            {
                Logger.LogWarning(e, "Failed to ack a TunnelingRequest");
            }

            if (!(tunnelingRequest.Cemi is LDataInd { DataFrame: LDataExtended ext })
                || !(ext.Apdu is ApduDataContainer container))
            {
                return;
            }

            switch (container.DataApdu)
            {
                case ApduDataGroupValueWrite write:
                    DispatchGroupValue(ext.SourceAddress, ext.DestinationAddress,
                        write.DataFirstByte, write.Data);
                    break;
                case ApduDataGroupValueResponse response:
                    CompletePendingRead(ext.DestinationAddress, response.DataFirstByte, response.Data);
                    DispatchGroupValue(ext.SourceAddress, ext.DestinationAddress,
                        response.DataFirstByte, response.Data);
                    break;
            }
        }

        private static byte ExtractSequenceCounter(TunnelingRequest request)
        {
            // The generated TunnelingRequest does not expose its data block as a
            // property. Re-serialise: after the 6-byte KNXnet/IP header comes the
            // TunnelingRequestDataBlock - structLen(1), channel(1), sequence(1),
            // reserved(1) - so the sequence counter is byte 8.
            var buffer = new WriteBuffer();
            request.Serialize(buffer);
            var bytes = buffer.GetBytes();
            return bytes.Length > 8 ? bytes[8] : (byte) 0;
        }

        private void DispatchGroupValue(KnxAddress source, byte[] destination, sbyte firstByte, byte[] rest)
        {
            var groupAddress = DecodeGroupAddress(destination);
            var payload = new byte[1 + (rest?.Length ?? 0)];
            payload[0] = (byte) firstByte;
            if (rest is { Length: > 0 })
            {
                Array.Copy(rest, 0, payload, 1, rest.Length);
            }

            IPlcValue value = new PlcRawByteArray(payload);
            if (_dptHints.TryGetValue(groupAddress, out var dpt))
            {
                try
                {
                    value = KnxDatapoint.StaticParse(new ReadBuffer(payload), dpt);
                }
                catch (Exception e)
                {
                    Logger.LogWarning(e, "Failed to decode a telegram for {GroupAddress} as {Dpt}",
                        groupAddress, dpt);
                }
            }

            var sourceAddress =
                $"{source.MainGroup}.{source.MiddleGroup}.{source.SubGroup}";
            var evt = new KnxGroupValueEvent(sourceAddress, groupAddress, payload, value);

            Action<KnxGroupValueEvent>[] snapshot;
            lock (_listenerLock)
            {
                snapshot = _listeners.ToArray();
            }
            foreach (var listener in snapshot)
            {
                try
                {
                    listener(evt);
                }
                catch (Exception e)
                {
                    Logger.LogWarning(e, "A KNX group-value listener threw");
                }
            }
        }

        private void CompletePendingRead(byte[] destination, sbyte firstByte, byte[] rest)
        {
            var groupAddress = DecodeGroupAddress(destination);
            if (!_pendingReads.TryRemove(groupAddress, out var pending))
            {
                return;
            }
            var payload = new byte[1 + (rest?.Length ?? 0)];
            payload[0] = (byte) firstByte;
            if (rest is { Length: > 0 })
            {
                Array.Copy(rest, 0, payload, 1, rest.Length);
            }
            pending.TrySetResult(payload);
        }

        private string DecodeGroupAddress(byte[] wire)
        {
            var raw = (wire[0] << 8) | wire[1];
            return _groupAddressLevels switch
            {
                3 => $"{(raw >> 11) & 0x1F}/{(raw >> 8) & 0x07}/{raw & 0xFF}",
                2 => $"{(raw >> 11) & 0x1F}/{raw & 0x7FF}",
                _ => (raw & 0xFFFF).ToString(CultureInfo.InvariantCulture),
            };
        }

        // ── bus monitor hook ───────────────────────────────────

        /// <summary>
        /// Registers a callback for every group telegram the gateway forwards. The
        /// callback runs on the transport receive thread, so it should hand off
        /// rather than block. Pass a <c>:DPT…</c>-suffixed address to
        /// <see cref="RegisterDatapointHint"/> first to have the value decoded.
        /// </summary>
        public void RegisterGroupValueListener(Action<KnxGroupValueEvent> listener)
        {
            if (listener == null)
            {
                throw new ArgumentNullException(nameof(listener));
            }
            lock (_listenerLock)
            {
                _listeners.Add(listener);
            }
        }

        public void RemoveGroupValueListener(Action<KnxGroupValueEvent> listener)
        {
            lock (_listenerLock)
            {
                _listeners.Remove(listener);
            }
        }

        /// <summary>
        /// Tells the connection the datapoint type of a group address, so inbound
        /// telegrams for it decode to a typed value. A <c>1/2/3:DPT9.001</c> tag
        /// registers its own hint when read; this is for addresses only monitored.
        /// </summary>
        public void RegisterDatapointHint(string groupAddress, string dptId)
        {
            var dpt = new KnxNetIpTag(2, "0", null, "0", dptId).ResolveDatapointType();
            if (dpt.HasValue)
            {
                _dptHints[groupAddress] = dpt.Value;
            }
        }

        // ── PlcReader ──────────────────────────────────────────

        public Task<IPlcReadResponse> Read(DefaultPlcReadRequest request, CancellationToken cancellationToken = default)
        {
            var results = new Dictionary<string, PlcResponseItem<IPlcValue>>();
            foreach (var name in request.TagNames)
            {
                cancellationToken.ThrowIfCancellationRequested();
                if (!(request.GetTagByName(name) is KnxNetIpTag tag))
                {
                    results[name] = new DefaultPlcTagErrorItem<IPlcValue>(PlcResponseCode.NotFound);
                    continue;
                }
                try
                {
                    results[name] = ReadTag(tag);
                }
                catch (TimeoutException)
                {
                    results[name] = new DefaultPlcTagErrorItem<IPlcValue>(PlcResponseCode.NotFound);
                }
                catch (Exception e)
                {
                    Logger.LogWarning(e, "KNX read of {Tag} failed", tag.AddressString);
                    results[name] = new DefaultPlcTagErrorItem<IPlcValue>(PlcResponseCode.InternalError);
                }
            }
            return Task.FromResult<IPlcReadResponse>(new DefaultPlcReadResponse(request, results));
        }

        private PlcResponseItem<IPlcValue> ReadTag(KnxNetIpTag tag)
        {
            var wire = tag.ToWireAddress(_groupAddressLevels);
            var groupAddress = DecodeGroupAddress(wire);
            var dpt = tag.ResolveDatapointType();
            if (dpt.HasValue)
            {
                _dptHints[groupAddress] = dpt.Value;
            }

            var responseSlot = new TaskCompletionSource<byte[]>(TaskCreationOptions.RunContinuationsAsynchronously);
            _pendingReads[groupAddress] = responseSlot;
            try
            {
                SendTunnelling(wire, new ApduDataContainer(false, 0, new ApduDataGroupValueRead()));
                var payload = Await(responseSlot.Task, "GroupValueResponse");
                if (dpt.HasValue)
                {
                    return new DefaultPlcResponseItem<IPlcValue>(
                        PlcResponseCode.Ok, KnxDatapoint.StaticParse(new ReadBuffer(payload), dpt.Value));
                }
                return new DefaultPlcResponseItem<IPlcValue>(PlcResponseCode.Ok, new PlcRawByteArray(payload));
            }
            finally
            {
                _pendingReads.TryRemove(groupAddress, out _);
            }
        }

        // ── PlcWriter ──────────────────────────────────────────

        public Task<IPlcWriteResponse> Write(DefaultPlcWriteRequest request, CancellationToken cancellationToken = default)
        {
            var codes = new Dictionary<string, PlcResponseCode>();
            foreach (var name in request.TagNames)
            {
                cancellationToken.ThrowIfCancellationRequested();
                if (!(request.GetTagByName(name) is KnxNetIpTag tag))
                {
                    codes[name] = PlcResponseCode.NotFound;
                    continue;
                }
                try
                {
                    var payload = EncodePayload(tag, request.GetValue(name));
                    SendTunnelling(tag.ToWireAddress(_groupAddressLevels),
                        new ApduDataContainer(false, 0,
                            new ApduDataGroupValueWrite((sbyte) payload[0], payload.Skip(1).ToArray())));
                    codes[name] = PlcResponseCode.Ok;
                }
                catch (Exception e)
                {
                    Logger.LogWarning(e, "KNX write of {Tag} failed", tag.AddressString);
                    codes[name] = PlcResponseCode.InternalError;
                }
            }
            return Task.FromResult<IPlcWriteResponse>(new DefaultPlcWriteResponse(request, codes));
        }

        private byte[] EncodePayload(KnxNetIpTag tag, object value)
        {
            var dpt = tag.ResolveDatapointType();
            if (dpt.HasValue)
            {
                var plcValue = AsPlcValue(value);
                var buffer = new WriteBuffer();
                KnxDatapoint.StaticSerialize(buffer, plcValue, dpt.Value);
                var bytes = buffer.GetBytes();
                return bytes.Length == 0 ? new byte[] { 0 } : bytes;
            }

            // No datapoint type: a raw single byte (0..63 in the APDU's 6 bits) or a
            // byte array, matching the Java no-ETS path.
            return value switch
            {
                byte b => new[] { b },
                byte[] bytes when bytes.Length > 0 => bytes,
                bool flag => new byte[] { (byte) (flag ? 1 : 0) },
                _ => throw new PlcException(
                    $"Without a :DPT… suffix on '{tag.AddressString}', only a byte, bool or byte[] can be written."),
            };
        }

        private static IPlcValue AsPlcValue(object value) => value switch
        {
            IPlcValue plcValue => plcValue,
            bool b => new PlcBOOL(b),
            byte b => new PlcUSINT(b),
            sbyte sb => new PlcSINT(sb),
            short s => new PlcINT(s),
            ushort us => new PlcUINT(us),
            int i => new PlcDINT(i),
            uint ui => new PlcUDINT(ui),
            long l => new PlcLINT(l),
            float f => new PlcREAL(f),
            double d => new PlcLREAL(d),
            string str => new PlcSTRING(str),
            _ => throw new PlcException($"Cannot convert {value?.GetType().Name ?? "null"} to a KNX value."),
        };

        // ── frame plumbing ─────────────────────────────────────

        private void SendTunnelling(byte[] destinationGroupAddress, Apdu apdu)
        {
            var seq = (byte) (Interlocked.Increment(ref _sequenceCounter) & 0xFF);
            var ackSlot = new TaskCompletionSource<TunnelingResponse>(TaskCreationOptions.RunContinuationsAsynchronously);
            _pendingAcks[seq] = ackSlot;

            // Frame values known to work against real gateways (see the Java driver):
            // frameType true, not repeated, LOW priority, group addressing, hop count 6.
            var request = new TunnelingRequest(
                new TunnelingRequestDataBlock(_communicationChannelId, seq),
                new LDataReq(0, new List<CEMIAdditionalInformation>(),
                    new LDataExtended(true, false, CEMIPriority.LOW, false, false,
                        true, 6, 0, _clientKnxAddress, destinationGroupAddress, apdu)));

            try
            {
                Send(request);
            }
            catch
            {
                _pendingAcks.TryRemove(seq, out _);
                throw;
            }

            var ack = Await(ackSlot.Task, "TunnelingResponse");
            if (ack.TunnelingResponseDataBlock.Status != Status.NO_ERROR)
            {
                throw new PlcException(
                    $"KNXnet/IP tunnelling request rejected: {ack.TunnelingResponseDataBlock.Status}.");
            }
        }

        private void Send(KnxNetIpMessage message)
        {
            try
            {
                _codec.Send(message);
            }
            catch (MessageCodecException e)
            {
                throw new PlcException($"Failed to send {message.GetType().Name}", e);
            }
        }

        private T Await<T>(Task<T> task, string what)
        {
            if (!task.Wait(_requestTimeoutMs))
            {
                throw new TimeoutException($"No {what} within {_requestTimeoutMs} ms.");
            }
            return task.Result;
        }

        // ── lifecycle ──────────────────────────────────────────

        public override void Close()
        {
            try
            {
                _heartbeat?.Dispose();
                _heartbeat = null;

                if (_communicationChannelId != 0 && _asyncTransport is { IsOpen: true })
                {
                    var local = ResolveLocalEndpoint();
                    Send(new DisconnectRequest(_communicationChannelId,
                        new HPAIControlEndpoint(HostProtocolCode.IPV4_UDP, local.Address, local.Port)));
                }
            }
            catch (Exception e)
            {
                Logger.LogDebug(e, "Error while sending DISCONNECT_REQUEST");
            }
            finally
            {
                _asyncTransport?.RemoveDataListener();
                base.Close();
            }
        }
    }
}
