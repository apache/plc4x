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
using System.Threading;
using System.Threading.Tasks;
using Microsoft.Extensions.Logging;
using org.apache.plc4net.drivers.modbus.readwrite.model;
using org.apache.plc4net.exceptions;
using org.apache.plc4net.spi.drivers;
using org.apache.plc4net.spi.transports;

namespace org.apache.plc4net.drivers.modbus
{
    /// <summary>
    /// A Modbus TCP connection: wraps each request PDU in a generated
    /// <see cref="ModbusTcpADU"/>, and matches the response by its transaction identifier.
    /// </summary>
    public sealed class ModbusTcpConnection : ModbusConnectionBase
    {
        private const int PollIntervalMs = 2;

        // Modbus is strict request/response with one outstanding transaction per
        // connection; the gate keeps two concurrent calls from interleaving.
        private readonly SemaphoreSlim _exchangeGate = new SemaphoreSlim(1, 1);
        private readonly ModbusTcpMessageCodec _codec;
        private readonly int _requestTimeoutMs;

        private ushort _lastTransactionId;
        private ushort _expectedTransactionId;
        private ModbusTcpADU? _response;

        public ModbusTcpConnection(ConnectionString connectionString, ITransportInstance transport)
            : base(connectionString, transport)
        {
            var unitId = IntParameter(connectionString, 1, "default-unit-identifier", "unit-identifier");
            if (unitId < 0 || unitId > 255)
            {
                throw new PlcConnectionException(
                    $"default-unit-identifier must be between 0 and 255, but was {unitId}.");
            }
            UnitIdentifier = (byte)unitId;
            _requestTimeoutMs = RequestTimeoutMs(connectionString);

            _codec = new ModbusTcpMessageCodec(transport, OnMessage);

            if (transport is IAsyncTransportInstance asyncTransport)
            {
                asyncTransport.RegisterDisconnectListener(ex =>
                    Logger.LogWarning(ex, "The Modbus TCP transport disconnected"));
            }
        }

        /// <summary>The unit identifier that every request carries.</summary>
        public byte UnitIdentifier { get; }

        protected override async Task<ModbusPDU> ExchangeAsync(
            ModbusPDU request, CancellationToken cancellationToken)
        {
            await _exchangeGate.WaitAsync(cancellationToken).ConfigureAwait(false);
            try
            {
                var transactionId = NextTransactionId(_lastTransactionId);
                _lastTransactionId = transactionId;
                _expectedTransactionId = transactionId;
                _response = null;

                _codec.Send(new ModbusTcpADU(transactionId, UnitIdentifier, request));
                var deadline = Environment.TickCount64 + _requestTimeoutMs;

                // A response left over from an earlier request that timed out carries
                // another transaction identifier; OnMessage skips it.
                while (true)
                {
                    _codec.ProcessIncomingData();
                    if (_response != null)
                    {
                        return _response.Pdu;
                    }

                    if (!TransportInstance.IsOpen)
                    {
                        throw new ModbusDriverException("The Modbus TCP connection was closed.");
                    }
                    if (Environment.TickCount64 >= deadline)
                    {
                        throw new TimeoutException(
                            $"No Modbus TCP response to transaction {transactionId} within {_requestTimeoutMs} ms.");
                    }
                    await Task.Delay(PollIntervalMs, cancellationToken).ConfigureAwait(false);
                }
            }
            catch (MessageCodecException ex)
            {
                DiscardBufferedBytes();
                throw new ModbusDriverException($"Modbus TCP exchange failed: {ex.Message}", ex);
            }
            catch (TimeoutException)
            {
                // Whatever is half-received now would be misread as the start of the
                // next response.
                DiscardBufferedBytes();
                throw;
            }
            finally
            {
                _response = null;
                _exchangeGate.Release();
            }
        }

        /// <summary>
        /// The transaction identifier that follows <paramref name="last"/>. Identifiers run from
        /// 1 to 0xFFFE and then start again at 1, as in plc4j; 0 and 0xFFFF are never sent.
        /// </summary>
        internal static ushort NextTransactionId(ushort last) => last >= 0xFFFE ? (ushort)1 : (ushort)(last + 1);

        private void OnMessage(ModbusTcpADU adu)
        {
            if (adu.TransactionIdentifier != _expectedTransactionId)
            {
                Logger.LogDebug(
                    "Modbus TCP: discarded a stale response (transaction {Got}, waiting for {Want})",
                    adu.TransactionIdentifier, _expectedTransactionId);
                return;
            }
            _response = adu;
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
                // The transport is gone; there is nothing left to discard, and this
                // must not hide the failure that led here.
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
    }
}