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
    /// A Modbus RTU connection: wraps each request PDU in a generated
    /// <see cref="ModbusRtuADU"/>, which also computes and checks the CRC.
    /// </summary>
    /// <remarks>
    /// Connection string format:
    /// <c>modbus-rtu:serial://COM1?default-unit-identifier=1&amp;baud-rate=19200&amp;parity=Even</c>
    /// </remarks>
    public sealed class ModbusRtuConnection : ModbusConnectionBase
    {
        private const int PollIntervalMs = 2;

        // Half-duplex bus: one outstanding request per connection.
        private readonly SemaphoreSlim _exchangeGate = new SemaphoreSlim(1, 1);
        private readonly ModbusRtuMessageCodec _codec;
        private readonly int _requestTimeoutMs;

        private byte _expectedFunction;
        private ModbusRtuADU? _response;

        public ModbusRtuConnection(ConnectionString connectionString, ITransportInstance transport)
            : base(connectionString, transport)
        {
            var address = IntParameter(connectionString, 1, "default-unit-identifier", "unit-identifier");
            if (address < 1 || address > 247)
            {
                throw new PlcConnectionException(
                    $"The Modbus RTU default-unit-identifier must be between 1 and 247, but was {address}.");
            }
            UnitIdentifier = (byte)address;
            _requestTimeoutMs = RequestTimeoutMs(connectionString);

            _codec = new ModbusRtuMessageCodec(transport, OnMessage);

            if (transport is IAsyncTransportInstance asyncTransport)
            {
                asyncTransport.RegisterDisconnectListener(ex =>
                    Logger.LogWarning(ex, "The Modbus RTU serial port disconnected"));
            }
        }

        /// <summary>The slave address that every request carries and every response must carry.</summary>
        public byte UnitIdentifier { get; }

        protected override async Task<ModbusPDU> ExchangeAsync(
            ModbusPDU request, CancellationToken cancellationToken)
        {
            await _exchangeGate.WaitAsync(cancellationToken).ConfigureAwait(false);
            try
            {
                // A frame left in the buffer by an earlier failure would be read as
                // if it answered this request.
                DiscardBufferedBytes();
                _response = null;
                _expectedFunction = request.FunctionFlag;

                _codec.Send(new ModbusRtuADU(UnitIdentifier, request));
                var deadline = Environment.TickCount64 + _requestTimeoutMs;

                // The echo of a read request that a half-duplex adapter sends back is skipped
                // by the codec like any other noise. (The response to a write is byte-for-byte
                // its request, so an echo of a write cannot be told from the response.)
                while (true)
                {
                    _codec.ProcessIncomingData();
                    if (_response != null)
                    {
                        return _response.Pdu;
                    }

                    if (!TransportInstance.IsOpen)
                    {
                        throw new ModbusDriverException("The Modbus RTU connection was closed.");
                    }
                    if (Environment.TickCount64 >= deadline)
                    {
                        throw new TimeoutException(
                            $"No Modbus RTU response within {_requestTimeoutMs} ms.");
                    }
                    await Task.Delay(PollIntervalMs, cancellationToken).ConfigureAwait(false);
                }
            }
            catch (MessageCodecException ex)
            {
                // The transport failed while sending or receiving. (A frame that fails
                // its CRC is skipped by the codec and never gets here: it ends as a
                // timeout.)
                DiscardBufferedBytes();
                throw new ModbusDriverException($"Modbus RTU exchange failed: {ex.Message}", ex);
            }
            catch (TimeoutException)
            {
                DiscardBufferedBytes();
                throw;
            }
            finally
            {
                _response = null;
                _exchangeGate.Release();
            }
        }

        private void OnMessage(ModbusRtuADU adu)
        {
            // Something else on the bus, or a late answer to an earlier request: not ours,
            // so keep waiting for the one that is.
            if (adu.Address != UnitIdentifier)
            {
                Logger.LogDebug(
                    "Modbus RTU: ignored a frame from address {Got}, waiting for {Want}",
                    adu.Address, UnitIdentifier);
                return;
            }

            // A normal response and an exception response both carry the function code of
            // the request they answer.
            if (adu.Pdu.FunctionFlag != _expectedFunction)
            {
                Logger.LogDebug(
                    "Modbus RTU: ignored a response to function 0x{Got:X2}, waiting for 0x{Want:X2}",
                    adu.Pdu.FunctionFlag, _expectedFunction);
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