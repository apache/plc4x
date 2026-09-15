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
    public partial class S7PayloadUserDataItemCpuFunctionAlarmAckResponse : S7PayloadUserDataItem
    {
        public override byte CpuFunctionGroup => (byte) (0x04);
        public override byte CpuFunctionType => (byte) (0x08);
        public override byte CpuSubfunction => (byte) (0x0b);

        public byte FunctionId { get; }
        public System.Collections.Generic.List<byte> MessageObjects { get; }

        public S7PayloadUserDataItemCpuFunctionAlarmAckResponse(DataTransportErrorCode returnCode, DataTransportSize transportSize, ushort dataLength, byte functionId, System.Collections.Generic.List<byte> messageObjects) : base(returnCode, transportSize, dataLength)
        {
            FunctionId = functionId;
            MessageObjects = messageObjects;
        }

        public static S7PayloadUserDataItemCpuFunctionAlarmAckResponse StaticParse(ReadBuffer readBuffer, byte cpuFunctionGroup, byte cpuFunctionType, byte cpuSubfunction, DataTransportErrorCode returnCode, DataTransportSize transportSize, ushort dataLength)
        {
            var functionId = readBuffer.ReadByte("functionId", 8);
            var numberOfObjects = readBuffer.ReadByte("numberOfObjects", 8);
            var messageObjects = new System.Collections.Generic.List<byte>();
            var _messageObjectsCnt = (int) (numberOfObjects);
            for (var _messageObjectsI = 0; _messageObjectsI < _messageObjectsCnt; _messageObjectsI++)
            {
                messageObjects.Add(readBuffer.ReadByte("messageObjects", 8));
            }
            return new S7PayloadUserDataItemCpuFunctionAlarmAckResponse(returnCode, transportSize, dataLength, functionId, messageObjects);
        }

        protected override void SerializeChild(WriteBuffer writeBuffer)
        {
            writeBuffer.WriteByte("functionId", 8, FunctionId);
            writeBuffer.WriteByte("numberOfObjects", 8, (byte) (MessageObjects.Count));
            foreach (var _e in MessageObjects)
            {
                writeBuffer.WriteByte("messageObjects", 8, _e);
            }
        }

        protected override int GetLengthInBitsChild()
        {
            var lengthInBits = 0;
            lengthInBits += 8;
            lengthInBits += 8;
            lengthInBits += (MessageObjects.Count * 8);
            return lengthInBits;
        }

    }
}
