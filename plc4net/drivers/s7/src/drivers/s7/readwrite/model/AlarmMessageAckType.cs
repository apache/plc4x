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
    public partial class AlarmMessageAckType : IMessage
    {
        public byte FunctionId { get; }
        public byte NumberOfObjects { get; }
        public System.Collections.Generic.List<AlarmMessageObjectAckType> MessageObjects { get; }

        public AlarmMessageAckType(byte functionId, byte numberOfObjects, System.Collections.Generic.List<AlarmMessageObjectAckType> messageObjects)
        {
            FunctionId = functionId;
            NumberOfObjects = numberOfObjects;
            MessageObjects = messageObjects;
        }

        public static AlarmMessageAckType StaticParse(ReadBuffer readBuffer)
        {
            var functionId = readBuffer.ReadByte("functionId", 8);
            var numberOfObjects = readBuffer.ReadByte("numberOfObjects", 8);
            var messageObjects = new System.Collections.Generic.List<AlarmMessageObjectAckType>();
            var _messageObjectsCnt = (int) (numberOfObjects);
            for (var _messageObjectsI = 0; _messageObjectsI < _messageObjectsCnt; _messageObjectsI++)
            {
                messageObjects.Add(AlarmMessageObjectAckType.StaticParse(readBuffer));
            }
            return new AlarmMessageAckType(functionId, numberOfObjects, messageObjects);
        }

        public void Serialize(WriteBuffer writeBuffer)
        {
            writeBuffer.WriteByte("functionId", 8, FunctionId);
            writeBuffer.WriteByte("numberOfObjects", 8, NumberOfObjects);
            foreach (var _e in MessageObjects)
            {
                _e.Serialize(writeBuffer);
            }
        }

        public int GetLengthInBits()
        {
            var lengthInBits = 0;
            lengthInBits += 8;
            lengthInBits += 8;
            lengthInBits += MessageObjects.Sum(_e => _e.GetLengthInBits());
            return lengthInBits;
        }

        public int GetLengthInBytes() => GetLengthInBits() / 8;
    }
}
