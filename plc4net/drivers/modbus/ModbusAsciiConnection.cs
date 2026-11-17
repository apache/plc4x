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
    /// <summary>A serial Modbus ASCII connection with one outstanding request.</summary>
    public sealed class ModbusAsciiConnection : ModbusConnectionBase
    {
        private const int PollIntervalMs = 2;
        private readonly SemaphoreSlim _exchangeGate = new SemaphoreSlim(1, 1);
        private readonly ModbusAsciiMessageCodec _codec;
        private readonly int _requestTimeoutMs;
        private byte _expectedFunction;
        private byte _expectedUnitIdentifier;
        private ModbusAsciiADU? _response;

        public ModbusAsciiConnection(ConnectionString connectionString, ITransportInstance transport)
            : base(connectionString, transport)
        {
            var address = IntParameter(connectionString, 1, "default-unit-identifier", "unit-identifier");
            if (address < 1 || address > 247)
            {
                throw new PlcConnectionException(
                    $"The Modbus ASCII default-unit-identifier must be between 1 and 247, but was {address}.");
            }
            UnitIdentifier = (byte)address;
            _requestTimeoutMs = RequestTimeoutMs(connectionString);
            _codec = new ModbusAsciiMessageCodec(transport, OnMessage);
            if (transport is IAsyncTransportInstance asyncTransport)
            {
                asyncTransport.RegisterDisconnectListener(ex =>
                    Logger.LogWarning(ex, "The Modbus ASCII serial port disconnected"));
            }
        }

        public byte UnitIdentifier { get; }

        protected override async Task<ModbusPDU> ExchangeAsync(
            ModbusPDU request, byte? unitIdentifier, CancellationToken cancellationToken)
        {
            await _exchangeGate.WaitAsync(cancellationToken).ConfigureAwait(false);
            try
            {
                DiscardBufferedBytes();
                _response = null;
                _expectedFunction = request.FunctionFlag;
                _expectedUnitIdentifier = unitIdentifier ?? UnitIdentifier;
                if (_expectedUnitIdentifier < 1 || _expectedUnitIdentifier > 247)
                    throw new ModbusDriverException($"Modbus ASCII unit-id must be between 1 and 247, but was {_expectedUnitIdentifier}.");
                _codec.Send(new ModbusAsciiADU(_expectedUnitIdentifier, request));
                var deadline = Environment.TickCount64 + _requestTimeoutMs;
                while (true)
                {
                    _codec.ProcessIncomingData();
                    if (_response != null) return _response.Pdu;
                    if (!TransportInstance.IsOpen)
                    {
                        throw new ModbusDriverException("The Modbus ASCII connection was closed.");
                    }
                    if (Environment.TickCount64 >= deadline)
                    {
                        throw new TimeoutException(
                            $"No Modbus ASCII response within {_requestTimeoutMs} ms.");
                    }
                    await Task.Delay(PollIntervalMs, cancellationToken).ConfigureAwait(false);
                }
            }
            catch (MessageCodecException ex)
            {
                DiscardBufferedBytes();
                throw new ModbusDriverException($"Modbus ASCII exchange failed: {ex.Message}", ex);
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

        private void OnMessage(ModbusAsciiADU adu)
        {
            if (adu.Address != _expectedUnitIdentifier || adu.Pdu.FunctionFlag != _expectedFunction)
            {
                Logger.LogDebug(
                    "Modbus ASCII ignored address {Address}, function 0x{Function:X2}",
                    adu.Address, adu.Pdu.FunctionFlag);
                return;
            }
            _response = adu;
        }

        private void DiscardBufferedBytes()
        {
            try
            {
                var available = TransportInstance.GetNumBytesAvailable();
                if (available > 0) TransportInstance.Read(available);
            }
            catch (TransportException)
            {
                // A disconnected transport has nothing useful left to discard.
            }
        }

        protected override void Dispose(bool disposing)
        {
            if (disposing) _exchangeGate.Dispose();
            base.Dispose(disposing);
        }
    }
}
