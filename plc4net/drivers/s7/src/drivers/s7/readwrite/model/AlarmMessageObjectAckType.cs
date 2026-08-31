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
    public partial class AlarmMessageObjectAckType : IMessage
    {
        public const byte VariableSpec = 0x12;
        public const byte Length = 0x08;

        public SyntaxIdType SyntaxId { get; }
        public byte NumberOfValues { get; }
        public uint EventId { get; }
        public State AckStateGoing { get; }
        public State AckStateComing { get; }

        public AlarmMessageObjectAckType(SyntaxIdType syntaxId, byte numberOfValues, uint eventId, State ackStateGoing, State ackStateComing)
        {
            SyntaxId = syntaxId;
            NumberOfValues = numberOfValues;
            EventId = eventId;
            AckStateGoing = ackStateGoing;
            AckStateComing = ackStateComing;
        }

        public static AlarmMessageObjectAckType StaticParse(ReadBuffer readBuffer)
        {
            var variableSpec = readBuffer.ReadByte("variableSpec", 8);
            if (!Equals(variableSpec, (byte) (0x12)))
                throw new ParseException($"Expected constant {VariableSpec} for 'variableSpec' but got {variableSpec}");
            var length = readBuffer.ReadByte("length", 8);
            if (!Equals(length, (byte) (0x08)))
                throw new ParseException($"Expected constant {Length} for 'length' but got {length}");
            var syntaxId = (SyntaxIdType) readBuffer.ReadByte("syntaxId", 8);
            var numberOfValues = readBuffer.ReadByte("numberOfValues", 8);
            var eventId = readBuffer.ReadUint("eventId", 32);
            var ackStateGoing = State.StaticParse(readBuffer);
            var ackStateComing = State.StaticParse(readBuffer);
            return new AlarmMessageObjectAckType(syntaxId, numberOfValues, eventId, ackStateGoing, ackStateComing);
        }

        public void Serialize(WriteBuffer writeBuffer)
        {
            writeBuffer.WriteByte("variableSpec", 8, VariableSpec);
            writeBuffer.WriteByte("length", 8, Length);
            writeBuffer.WriteByte("syntaxId", 8, (byte) SyntaxId);
            writeBuffer.WriteByte("numberOfValues", 8, NumberOfValues);
            writeBuffer.WriteUint("eventId", 32, EventId);
            AckStateGoing.Serialize(writeBuffer);
            AckStateComing.Serialize(writeBuffer);
        }

        public int GetLengthInBits()
        {
            var lengthInBits = 0;
            lengthInBits += 8;
            lengthInBits += 8;
            lengthInBits += 8;
            lengthInBits += 8;
            lengthInBits += 32;
            lengthInBits += AckStateGoing.GetLengthInBits();
            lengthInBits += AckStateComing.GetLengthInBits();
            return lengthInBits;
        }

        public int GetLengthInBytes() => GetLengthInBits() / 8;
    }
}
