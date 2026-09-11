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
    public partial class S7ParameterModeTransition : S7Parameter
    {
        public override byte ParameterType => (byte) (0x01);
        public override byte MessageType => (byte) (0x07);

        public byte Method { get; }
        public byte CpuFunctionType { get; }
        public byte CpuFunctionGroup { get; }
        public byte CurrentMode { get; }
        public byte SequenceNumber { get; }

        public S7ParameterModeTransition(byte method, byte cpuFunctionType, byte cpuFunctionGroup, byte currentMode, byte sequenceNumber)
        {
            Method = method;
            CpuFunctionType = cpuFunctionType;
            CpuFunctionGroup = cpuFunctionGroup;
            CurrentMode = currentMode;
            SequenceNumber = sequenceNumber;
        }

        public static new S7ParameterModeTransition StaticParse(ReadBuffer readBuffer, byte messageType)
        {
            {
                var reserved = readBuffer.ReadUshort("reserved", 16);
                if (!Equals(reserved, (ushort) (0x0010))) { /* mspec reserved: value differs from the spec default */ }
            }
            var itemLength = readBuffer.ReadByte("itemLength", 8);
            var method = readBuffer.ReadByte("method", 8);
            var cpuFunctionType = readBuffer.ReadByte("cpuFunctionType", 4);
            var cpuFunctionGroup = readBuffer.ReadByte("cpuFunctionGroup", 4);
            var currentMode = readBuffer.ReadByte("currentMode", 8);
            var sequenceNumber = readBuffer.ReadByte("sequenceNumber", 8);
            return new S7ParameterModeTransition(method, cpuFunctionType, cpuFunctionGroup, currentMode, sequenceNumber);
        }

        protected override void SerializeChild(WriteBuffer writeBuffer)
        {
            writeBuffer.WriteUshort("reserved", 16, (ushort) (0x0010));
            writeBuffer.WriteByte("itemLength", 8, (byte) ((GetLengthInBytes() - 2)));
            writeBuffer.WriteByte("method", 8, Method);
            writeBuffer.WriteByte("cpuFunctionType", 4, CpuFunctionType);
            writeBuffer.WriteByte("cpuFunctionGroup", 4, CpuFunctionGroup);
            writeBuffer.WriteByte("currentMode", 8, CurrentMode);
            writeBuffer.WriteByte("sequenceNumber", 8, SequenceNumber);
        }

        protected override int GetLengthInBitsChild()
        {
            var lengthInBits = 0;
            lengthInBits += 16;
            lengthInBits += 8;
            lengthInBits += 8;
            lengthInBits += 4;
            lengthInBits += 4;
            lengthInBits += 8;
            lengthInBits += 8;
            return lengthInBits;
        }

    }
}
