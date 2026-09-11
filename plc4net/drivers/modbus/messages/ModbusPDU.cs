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

namespace org.apache.plc4net.drivers.modbus.messages
{
    /// <summary>
    /// Modbus function codes.
    /// </summary>
    public static class ModbusFunctionCodes
    {
        public const byte ReadCoils = 0x01;
        public const byte ReadDiscreteInputs = 0x02;
        public const byte ReadHoldingRegisters = 0x03;
        public const byte ReadInputRegisters = 0x04;
        public const byte WriteSingleCoil = 0x05;
        public const byte WriteSingleRegister = 0x06;
        public const byte WriteMultipleCoils = 0x0F;
        public const byte WriteMultipleRegisters = 0x10;
        public const byte MaskWriteRegister = 0x16;
        public const byte ReadWriteMultipleRegisters = 0x17;

        /// <summary>Modbus error offset added to the function code.</summary>
        public const byte ErrorOffset = 0x80;
    }

    /// <summary>
    /// Well-known Modbus exception codes.
    /// </summary>
    public enum ModbusErrorCode : byte
    {
        IllegalFunction = 0x01,
        IllegalDataAddress = 0x02,
        IllegalDataValue = 0x03,
        SlaveDeviceFailure = 0x04,
        Acknowledge = 0x05,
        SlaveDeviceBusy = 0x06,
        MemoryParityError = 0x08,
    }

    /// <summary>
    /// Builds Modbus PDUs for read and write requests.
    /// </summary>
    public static class ModbusPDU
    {
        /// <summary>Builds a Read Coils or Read Discrete Inputs request.</summary>
        public static byte[] BuildReadBitsRequest(byte functionCode, ushort startAddress, ushort quantity)
            => BuildReadRequest(functionCode, startAddress, quantity);

        /// <summary>Builds a Read Holding/Input Registers request.</summary>
        public static byte[] BuildReadRegistersRequest(byte functionCode, ushort startAddress, ushort quantity)
            => BuildReadRequest(functionCode, startAddress, quantity);

        /// <summary>Shared implementation: all Modbus read function codes use the same PDU layout.</summary>
        private static byte[] BuildReadRequest(byte functionCode, ushort startAddress, ushort quantity)
        {
            return new byte[]
            {
                functionCode,
                (byte)(startAddress >> 8),
                (byte)(startAddress & 0xFF),
                (byte)(quantity >> 8),
                (byte)(quantity & 0xFF)
            };
        }

        /// <summary>Builds a Write Single Coil request.</summary>
        public static byte[] BuildWriteSingleCoilRequest(ushort address, bool value)
        {
            return new byte[]
            {
                ModbusFunctionCodes.WriteSingleCoil,
                (byte)(address >> 8),
                (byte)(address & 0xFF),
                (byte)(value ? 0xFF : 0x00),
                0x00
            };
        }

        /// <summary>Builds a Write Single Register request.</summary>
        public static byte[] BuildWriteSingleRegisterRequest(ushort address, ushort value)
        {
            return new byte[]
            {
                ModbusFunctionCodes.WriteSingleRegister,
                (byte)(address >> 8),
                (byte)(address & 0xFF),
                (byte)(value >> 8),
                (byte)(value & 0xFF)
            };
        }

        /// <summary>Builds a Write Multiple Registers request.</summary>
        public static byte[] BuildWriteMultipleRegistersRequest(
            ushort startAddress, ushort[] values)
        {
            if (values == null)
                throw new ArgumentNullException(nameof(values));
            var totalByteCount = values.Length * 2;
            if (totalByteCount > byte.MaxValue)
                throw new ArgumentException(
                    $"Too many values: {values.Length} registers requires {totalByteCount} bytes (max {byte.MaxValue}).");
            var byteCount = (byte)totalByteCount;
            var result = new List<byte>
            {
                ModbusFunctionCodes.WriteMultipleRegisters,
                (byte)(startAddress >> 8),
                (byte)(startAddress & 0xFF),
                (byte)(values.Length >> 8),
                (byte)(values.Length & 0xFF),
                byteCount
            };
            foreach (var v in values)
            {
                result.Add((byte)(v >> 8));
                result.Add((byte)(v & 0xFF));
            }
            return result.ToArray();
        }

        /// <summary>Parses a Read Coils/Discrete Inputs response.</summary>
        public static bool[] ParseReadBitsResponse(byte[] data, ushort quantity)
        {
            if (data == null || data.Length < 1)
                throw new ArgumentException("Response data must contain at least a byte count.");
            var byteCount = data[0];
            if (data.Length < 1 + byteCount)
                throw new ArgumentException(
                    $"Byte count {byteCount} exceeds available data ({data.Length - 1}).");
            var result = new bool[quantity];
            var maxBits = Math.Min(quantity, byteCount * 8);
            for (var i = 0; i < maxBits; i++)
            {
                var b = data[1 + i / 8];
                result[i] = ((b >> (i % 8)) & 0x01) == 0x01;
            }
            return result;
        }

        /// <summary>Parses a Read Holding/Input Registers response.</summary>
        public static ushort[] ParseReadRegistersResponse(byte[] data, ushort quantity)
        {
            if (data == null || data.Length < 1)
                throw new ArgumentException("Response data must contain at least a byte count.");
            var byteCount = data[0];
            if (data.Length < 1 + byteCount)
                throw new ArgumentException(
                    $"Byte count {byteCount} exceeds available data ({data.Length - 1}).");
            var result = new ushort[quantity];
            var maxRegs = Math.Min(quantity, byteCount / 2);
            for (var i = 0; i < maxRegs; i++)
            {
                result[i] = (ushort)((data[1 + i * 2] << 8) | data[2 + i * 2]);
            }
            return result;
        }

        /// <summary>Parses a Write Single Coil/Register response.</summary>
        public static bool ParseWriteSingleResponse(byte[] data)
        {
            // Write responses echo the request; a non-error response means success.
            return data.Length >= 4;
        }

        /// <summary>Parses a Write Multiple Registers response.</summary>
        public static int ParseWriteMultipleResponse(byte[] data)
        {
            // Response: FC + startAddr(2) + quantity(2)
            if (data.Length >= 4)
            {
                return (data[2] << 8) | data[3]; // quantity written
            }
            return 0;
        }
    }
}
