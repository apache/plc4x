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
using System.Net;
using System.Net.Sockets;
using System.Threading;
using System.Threading.Tasks;
using org.apache.plc4net.drivers.knxnetip.readwrite.model;
using org.apache.plc4net.spi.generation;
using IPAddress = System.Net.IPAddress;

namespace org.apache.plc4net.test.knxnetip
{
    /// <summary>
    /// A minimal KNXnet/IP gateway for driver tests: answers the SEARCH / CONNECT /
    /// CONNECTIONSTATE handshake with frames from the shared test suite, acks every
    /// tunnelling request, and - for a <c>GroupValueRead</c> - sends back a scripted
    /// <c>GroupValueResponse</c>.
    /// </summary>
    internal sealed class FakeKnxGateway : IDisposable
    {
        // Channel id 0x02 and client address 15.15.254 (0xFFFE) - the values in the
        // shared "Connect Response" vector.
        internal const byte ChannelId = 0x02;
        internal const string AssignedClientAddress = "15.15.254";

        private readonly Socket _socket;
        private readonly CancellationTokenSource _cts = new CancellationTokenSource();
        private readonly Task _loop;

        public FakeKnxGateway()
        {
            _socket = new Socket(AddressFamily.InterNetwork, SocketType.Dgram, ProtocolType.Udp);
            _socket.Bind(new IPEndPoint(IPAddress.Loopback, 0));
            Port = ((IPEndPoint) _socket.LocalEndPoint).Port;
            _loop = Task.Run(RunAsync);
        }

        public int Port { get; }

        /// <summary>Group telegrams the driver sent us (write path), newest last.</summary>
        public List<(string GroupAddress, byte[] Payload)> ReceivedWrites { get; } = new();

        /// <summary>Every tunnelling sequence counter the driver sent, in order.</summary>
        public List<byte> TunnellingSequenceNumbers { get; } = new();

        /// <summary>The scripted answer to the next GroupValueRead, or null to stay silent.</summary>
        public byte[] NextReadResponsePayload { get; set; }

        /// <summary>The driver's source endpoint, learned from the first datagram.</summary>
        private volatile EndPoint _driverEndpoint;

        private async Task RunAsync()
        {
            var buffer = new byte[2048];
            var remote = (EndPoint) new IPEndPoint(IPAddress.Any, 0);
            while (!_cts.IsCancellationRequested)
            {
                SocketReceiveFromResult received;
                try
                {
                    received = await _socket.ReceiveFromAsync(
                        new ArraySegment<byte>(buffer), SocketFlags.None, remote).ConfigureAwait(false);
                }
                catch (ObjectDisposedException)
                {
                    return;
                }
                catch (SocketException)
                {
                    continue;
                }

                _driverEndpoint = received.RemoteEndPoint;
                var frame = new byte[received.ReceivedBytes];
                Array.Copy(buffer, frame, frame.Length);
                Handle(frame, received.RemoteEndPoint);
            }
        }

        private void Handle(byte[] frame, EndPoint from)
        {
            var msgType = (ushort) ((frame[2] << 8) | frame[3]);
            switch (msgType)
            {
                case 0x0201: // SEARCH_REQUEST
                    SendTo(Hex(
                        "06100202004c0801c0a82a0b0e5736010200ffff000000082d409852e000170c000a" +
                        "b327553647697261204b4e582f49502d5363686e6974747374656c6c6500000000" +
                        "000802020103010401"), from);
                    break;
                case 0x0205: // CONNECT_REQUEST
                    SendTo(Hex("06100206001402000801c0a82a0b0e570404fffe"), from);
                    break;
                case 0x0207: // CONNECTIONSTATE_REQUEST
                    SendTo(Hex("0610020800080200"), from);
                    break;
                case 0x0209: // DISCONNECT_REQUEST
                    SendTo(Hex("0610020a00080200"), from);
                    break;
                case 0x0420: // TUNNELING_REQUEST
                    HandleTunnelling(frame, from);
                    break;
            }
        }

        private void HandleTunnelling(byte[] frame, EndPoint from)
        {
            var seq = frame[8];
            lock (TunnellingSequenceNumbers)
            {
                TunnellingSequenceNumbers.Add(seq);
            }

            // ack
            var ack = new TunnelingResponse(new TunnelingResponseDataBlock(ChannelId, seq, Status.NO_ERROR));
            SendTo(Serialize(ack), from);

            var message = KnxNetIpMessage.StaticParse(new ReadBuffer(frame));
            if (message is not TunnelingRequest { Cemi: LDataReq { DataFrame: LDataExtended ext } })
            {
                return;
            }
            if (ext.Apdu is not ApduDataContainer container)
            {
                return;
            }

            var groupAddress = DecodeThreeLevel(ext.DestinationAddress);
            switch (container.DataApdu)
            {
                case ApduDataGroupValueWrite write:
                    var payload = new byte[1 + write.Data.Length];
                    payload[0] = (byte) write.DataFirstByte;
                    Array.Copy(write.Data, 0, payload, 1, write.Data.Length);
                    lock (ReceivedWrites)
                    {
                        ReceivedWrites.Add((groupAddress, payload));
                    }
                    break;
                case ApduDataGroupValueRead when NextReadResponsePayload != null:
                    SendGroupValueResponse(ext.DestinationAddress, NextReadResponsePayload, from);
                    break;
            }
        }

        /// <summary>Pushes an unsolicited GroupValueWrite telegram (bus-monitor path).</summary>
        public void PushGroupWrite(byte[] destinationAddress, byte[] payload)
        {
            var driver = _driverEndpoint
                         ?? throw new InvalidOperationException("No driver endpoint recorded yet - connect first.");
            var apdu = new ApduDataContainer(false, 0,
                new ApduDataGroupValueWrite((sbyte) payload[0], payload[1..]));
            SendTo(Serialize(new TunnelingRequest(
                new TunnelingRequestDataBlock(ChannelId, 0x7F),
                new LDataInd(0, new List<CEMIAdditionalInformation>(),
                    new LDataExtended(true, false, CEMIPriority.LOW, false, false, true, 6, 0,
                        new KnxAddress(1, 1, 5), destinationAddress, apdu)))), driver);
        }

        private void SendGroupValueResponse(byte[] destinationAddress, byte[] payload, EndPoint from)
        {
            var apdu = new ApduDataContainer(false, 0,
                new ApduDataGroupValueResponse((sbyte) payload[0], payload[1..]));
            SendTo(Serialize(new TunnelingRequest(
                new TunnelingRequestDataBlock(ChannelId, 0x40),
                new LDataInd(0, new List<CEMIAdditionalInformation>(),
                    new LDataExtended(true, false, CEMIPriority.LOW, false, false, true, 6, 0,
                        new KnxAddress(1, 1, 5), destinationAddress, apdu)))), from);
        }

        private static string DecodeThreeLevel(byte[] wire)
        {
            var raw = (wire[0] << 8) | wire[1];
            return string.Create(CultureInfo.InvariantCulture,
                $"{(raw >> 11) & 0x1F}/{(raw >> 8) & 0x07}/{raw & 0xFF}");
        }

        private static byte[] Serialize(KnxNetIpMessage message)
        {
            var buffer = new WriteBuffer();
            message.Serialize(buffer);
            return buffer.GetBytes();
        }

        private void SendTo(byte[] bytes, EndPoint to) => _socket.SendTo(bytes, to);

        private static byte[] Hex(string hex)
        {
            var bytes = new byte[hex.Length / 2];
            for (var i = 0; i < bytes.Length; i++)
            {
                bytes[i] = byte.Parse(hex.Substring(i * 2, 2), NumberStyles.HexNumber, CultureInfo.InvariantCulture);
            }
            return bytes;
        }

        public void Dispose()
        {
            _cts.Cancel();
            _socket.Dispose();
            try
            {
                _loop.Wait(TimeSpan.FromSeconds(1));
            }
            catch
            {
                // best effort
            }
            _cts.Dispose();
        }
    }
}
