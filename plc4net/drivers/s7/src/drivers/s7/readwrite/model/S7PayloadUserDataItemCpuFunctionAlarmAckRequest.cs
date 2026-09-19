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
    public partial class S7PayloadUserDataItemCpuFunctionAlarmAckRequest : S7PayloadUserDataItem
    {
        public const byte FunctionId = 0x09;

        public override byte CpuFunctionGroup => (byte) (0x04);
        public override byte CpuFunctionType => (byte) (0x04);
        public override byte CpuSubfunction => (byte) (0x0b);

        public System.Collections.Generic.List<AlarmMessageObjectAckType> MessageObjects { get; }

        public S7PayloadUserDataItemCpuFunctionAlarmAckRequest(DataTransportErrorCode returnCode, DataTransportSize transportSize, ushort dataLength, System.Collections.Generic.List<AlarmMessageObjectAckType> messageObjects) : base(returnCode, transportSize, dataLength)
        {
            MessageObjects = messageObjects;
        }

        public static S7PayloadUserDataItemCpuFunctionAlarmAckRequest StaticParse(ReadBuffer readBuffer, byte cpuFunctionGroup, byte cpuFunctionType, byte cpuSubfunction, DataTransportErrorCode returnCode, DataTransportSize transportSize, ushort dataLength)
        {
            var functionId = readBuffer.ReadByte("functionId", 8);
            if (!Equals(functionId, (byte) (0x09)))
                throw new ParseException($"Expected constant {FunctionId} for 'functionId' but got {functionId}");
            var numberOfObjects = readBuffer.ReadByte("numberOfObjects", 8);
            var messageObjects = new System.Collections.Generic.List<AlarmMessageObjectAckType>();
            var _messageObjectsCnt = (int) (numberOfObjects);
            for (var _messageObjectsI = 0; _messageObjectsI < _messageObjectsCnt; _messageObjectsI++)
            {
                messageObjects.Add(AlarmMessageObjectAckType.StaticParse(readBuffer));
            }
            return new S7PayloadUserDataItemCpuFunctionAlarmAckRequest(returnCode, transportSize, dataLength, messageObjects);
        }

        protected override void SerializeChild(WriteBuffer writeBuffer)
        {
            writeBuffer.WriteByte("functionId", 8, FunctionId);
            writeBuffer.WriteByte("numberOfObjects", 8, (byte) (MessageObjects.Count));
            foreach (var _e in MessageObjects)
            {
                _e.Serialize(writeBuffer);
            }
        }

        protected override int GetLengthInBitsChild()
        {
            var lengthInBits = 0;
            lengthInBits += 8;
            lengthInBits += 8;
            lengthInBits += MessageObjects.Sum(_e => _e.GetLengthInBits());
            return lengthInBits;
        }

    }
}
