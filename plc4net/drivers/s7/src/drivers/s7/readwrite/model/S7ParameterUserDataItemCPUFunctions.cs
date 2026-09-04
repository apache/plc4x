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

// Code generated from the mspec by plc4net-code-gen. DO NOT EDIT.

using System;
using System.Linq;
using org.apache.plc4net.spi.drivers;
using org.apache.plc4net.spi.generation;

namespace org.apache.plc4net.drivers.s7.readwrite.model
{
    public partial class S7ParameterUserDataItemCPUFunctions : S7ParameterUserDataItem
    {
        public override byte ItemType => (byte) (0x12);

        public byte Method { get; }
        public byte CpuFunctionType { get; }
        public byte CpuFunctionGroup { get; }
        public byte CpuSubfunction { get; }
        public byte SequenceNumber { get; }
        public byte? DataUnitReferenceNumber { get; }
        public byte? LastDataUnit { get; }
        public ushort? ErrorCode { get; }

        public S7ParameterUserDataItemCPUFunctions(byte method, byte cpuFunctionType, byte cpuFunctionGroup, byte cpuSubfunction, byte sequenceNumber, byte? dataUnitReferenceNumber, byte? lastDataUnit, ushort? errorCode)
        {
            Method = method;
            CpuFunctionType = cpuFunctionType;
            CpuFunctionGroup = cpuFunctionGroup;
            CpuSubfunction = cpuSubfunction;
            SequenceNumber = sequenceNumber;
            DataUnitReferenceNumber = dataUnitReferenceNumber;
            LastDataUnit = lastDataUnit;
            ErrorCode = errorCode;
        }

        public static new S7ParameterUserDataItemCPUFunctions StaticParse(ReadBuffer readBuffer)
        {
            var itemLength = readBuffer.ReadByte("itemLength", 8);
            var method = readBuffer.ReadByte("method", 8);
            var cpuFunctionType = readBuffer.ReadByte("cpuFunctionType", 4);
            var cpuFunctionGroup = readBuffer.ReadByte("cpuFunctionGroup", 4);
            var cpuSubfunction = readBuffer.ReadByte("cpuSubfunction", 8);
            var sequenceNumber = readBuffer.ReadByte("sequenceNumber", 8);
            byte? dataUnitReferenceNumber = null;
            if (((cpuFunctionType == 8) || ((cpuFunctionType == 0) && (cpuFunctionGroup == 2))))
            {
                dataUnitReferenceNumber = readBuffer.ReadByte("dataUnitReferenceNumber", 8);
            }
            byte? lastDataUnit = null;
            if (((cpuFunctionType == 8) || ((cpuFunctionType == 0) && (cpuFunctionGroup == 2))))
            {
                lastDataUnit = readBuffer.ReadByte("lastDataUnit", 8);
            }
            ushort? errorCode = null;
            if (((cpuFunctionType == 8) || ((cpuFunctionType == 0) && (cpuFunctionGroup == 2))))
            {
                errorCode = readBuffer.ReadUshort("errorCode", 16);
            }
            return new S7ParameterUserDataItemCPUFunctions(method, cpuFunctionType, cpuFunctionGroup, cpuSubfunction, sequenceNumber, dataUnitReferenceNumber, lastDataUnit, errorCode);
        }

        protected override void SerializeChild(WriteBuffer writeBuffer)
        {
            writeBuffer.WriteByte("itemLength", 8, (byte) ((GetLengthInBytes() - 2)));
            writeBuffer.WriteByte("method", 8, Method);
            writeBuffer.WriteByte("cpuFunctionType", 4, CpuFunctionType);
            writeBuffer.WriteByte("cpuFunctionGroup", 4, CpuFunctionGroup);
            writeBuffer.WriteByte("cpuSubfunction", 8, CpuSubfunction);
            writeBuffer.WriteByte("sequenceNumber", 8, SequenceNumber);
            if (DataUnitReferenceNumber != null)
            {
                writeBuffer.WriteByte("dataUnitReferenceNumber", 8, DataUnitReferenceNumber.Value);
            }
            if (LastDataUnit != null)
            {
                writeBuffer.WriteByte("lastDataUnit", 8, LastDataUnit.Value);
            }
            if (ErrorCode != null)
            {
                writeBuffer.WriteUshort("errorCode", 16, ErrorCode.Value);
            }
        }

        protected override int GetLengthInBitsChild()
        {
            var lengthInBits = 0;
            lengthInBits += 8;
            lengthInBits += 8;
            lengthInBits += 4;
            lengthInBits += 4;
            lengthInBits += 8;
            lengthInBits += 8;
            lengthInBits += (DataUnitReferenceNumber != null ? 8 : 0);
            lengthInBits += (LastDataUnit != null ? 8 : 0);
            lengthInBits += (ErrorCode != null ? 16 : 0);
            return lengthInBits;
        }

    }
}
