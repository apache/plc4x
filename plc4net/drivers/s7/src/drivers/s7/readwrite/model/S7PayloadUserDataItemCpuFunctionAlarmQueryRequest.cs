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
    public partial class S7PayloadUserDataItemCpuFunctionAlarmQueryRequest : S7PayloadUserDataItem
    {
        public const byte FunctionId = 0x00;
        public const byte NumberMessageObj = 0x01;
        public const byte VariableSpec = 0x12;
        public const byte Length = 0x08;

        public override byte CpuFunctionGroup => (byte) (0x04);
        public override byte CpuFunctionType => (byte) (0x04);
        public override byte CpuSubfunction => (byte) (0x13);

        public SyntaxIdType SyntaxId { get; }
        public QueryType QueryType { get; }
        public AlarmType AlarmType { get; }

        public S7PayloadUserDataItemCpuFunctionAlarmQueryRequest(DataTransportErrorCode returnCode, DataTransportSize transportSize, ushort dataLength, SyntaxIdType syntaxId, QueryType queryType, AlarmType alarmType) : base(returnCode, transportSize, dataLength)
        {
            SyntaxId = syntaxId;
            QueryType = queryType;
            AlarmType = alarmType;
        }

        public static S7PayloadUserDataItemCpuFunctionAlarmQueryRequest StaticParse(ReadBuffer readBuffer, byte cpuFunctionGroup, byte cpuFunctionType, byte cpuSubfunction, DataTransportErrorCode returnCode, DataTransportSize transportSize, ushort dataLength)
        {
            var functionId = readBuffer.ReadByte("functionId", 8);
            if (!Equals(functionId, (byte) (0x00)))
                throw new ParseException($"Expected constant {FunctionId} for 'functionId' but got {functionId}");
            var numberMessageObj = readBuffer.ReadByte("numberMessageObj", 8);
            if (!Equals(numberMessageObj, (byte) (0x01)))
                throw new ParseException($"Expected constant {NumberMessageObj} for 'numberMessageObj' but got {numberMessageObj}");
            var variableSpec = readBuffer.ReadByte("variableSpec", 8);
            if (!Equals(variableSpec, (byte) (0x12)))
                throw new ParseException($"Expected constant {VariableSpec} for 'variableSpec' but got {variableSpec}");
            var length = readBuffer.ReadByte("length", 8);
            if (!Equals(length, (byte) (0x08)))
                throw new ParseException($"Expected constant {Length} for 'length' but got {length}");
            var syntaxId = (SyntaxIdType) readBuffer.ReadByte("syntaxId", 8);
            {
                var reserved = readBuffer.ReadByte("reserved", 8);
                if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
            }
            var queryType = (QueryType) readBuffer.ReadByte("queryType", 8);
            {
                var reserved = readBuffer.ReadByte("reserved", 8);
                if (!Equals(reserved, (byte) (0x34))) { /* mspec reserved: value differs from the spec default */ }
            }
            var alarmType = (AlarmType) readBuffer.ReadByte("alarmType", 8);
            return new S7PayloadUserDataItemCpuFunctionAlarmQueryRequest(returnCode, transportSize, dataLength, syntaxId, queryType, alarmType);
        }

        protected override void SerializeChild(WriteBuffer writeBuffer)
        {
            writeBuffer.WriteByte("functionId", 8, FunctionId);
            writeBuffer.WriteByte("numberMessageObj", 8, NumberMessageObj);
            writeBuffer.WriteByte("variableSpec", 8, VariableSpec);
            writeBuffer.WriteByte("length", 8, Length);
            writeBuffer.WriteByte("syntaxId", 8, (byte) SyntaxId);
            writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
            writeBuffer.WriteByte("queryType", 8, (byte) QueryType);
            writeBuffer.WriteByte("reserved", 8, (byte) (0x34));
            writeBuffer.WriteByte("alarmType", 8, (byte) AlarmType);
        }

        protected override int GetLengthInBitsChild()
        {
            var lengthInBits = 0;
            lengthInBits += 8;
            lengthInBits += 8;
            lengthInBits += 8;
            lengthInBits += 8;
            lengthInBits += 8;
            lengthInBits += 8;
            lengthInBits += 8;
            lengthInBits += 8;
            lengthInBits += 8;
            return lengthInBits;
        }

    }
}
