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
using System.Text;
using org.apache.plc4net.drivers.modbus.readwrite.model;
using org.apache.plc4net.spi.drivers;
using org.apache.plc4net.spi.generation;
using org.apache.plc4net.spi.transports;

namespace org.apache.plc4net.drivers.modbus
{
    /// <summary>Frames binary Modbus ADUs as ':' + uppercase hex + CRLF.</summary>
    internal sealed class ModbusAsciiMessageCodec : MessageCodecBase<ModbusAsciiADU>
    {
        public ModbusAsciiMessageCodec(ITransportInstance transport, Action<ModbusAsciiADU> handler)
            : base("Modbus ASCII", transport, handler)
        {
        }

        protected override int GetMinimumHeaderSize() => 1;

        protected override int CalculateTotalMessageSize(byte[] header, int availableBytes)
        {
            if (header[0] != (byte)':')
            {
                throw new MessageCodecException("A Modbus ASCII frame must start with ':'.");
            }
            var bytes = TransportInstance.PeekReadableBytes(availableBytes);
            for (var i = 1; i + 1 < bytes.Length; i++)
            {
                if (bytes[i] == (byte)'\r' && bytes[i + 1] == (byte)'\n')
                {
                    var length = i + 2;
                    try
                    {
                        ModbusAsciiADU.StaticParse(
                            new ReadBuffer(Decode(bytes, length)), DriverType.MODBUS_ASCII, true);
                    }
                    catch (ParseException)
                    {
                        throw new MessageCodecException("Invalid Modbus ASCII frame or LRC.");
                    }
                    return length;
                }
            }
            return -1;
        }

        public override void Send(ModbusAsciiADU message)
        {
            try
            {
                var buffer = new WriteBuffer();
                message.Serialize(buffer);
                var line = ":" + Convert.ToHexString(buffer.GetBytes()) + "\r\n";
                TransportInstance.Write(Encoding.ASCII.GetBytes(line));
            }
            catch (Exception ex) when (ex is ParseException || ex is TransportException)
            {
                throw new MessageCodecException("Failed to send Modbus ASCII message", ex);
            }
        }

        protected override ModbusAsciiADU ParseMessage(ReadBuffer readBuffer)
        {
            var lineBytes = readBuffer.ReadByteArray("line", readBuffer.GetTotalBytes() * 8);
            var binary = Decode(lineBytes, lineBytes.Length);
            return ModbusAsciiADU.StaticParse(
                new ReadBuffer(binary), DriverType.MODBUS_ASCII, true);
        }

        private static byte[] Decode(byte[] lineBytes, int length)
        {
            var line = Encoding.ASCII.GetString(lineBytes, 0, length);
            if (line.Length < 7 || line[0] != ':' || !line.EndsWith("\r\n", StringComparison.Ordinal)
                || ((line.Length - 3) & 1) != 0)
            {
                throw new ParseException("Malformed Modbus ASCII frame.");
            }
            try
            {
                return Convert.FromHexString(line.Substring(1, line.Length - 3));
            }
            catch (FormatException ex)
            {
                throw new ParseException($"Malformed Modbus ASCII hex payload: {ex.Message}");
            }
        }
    }
}
