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
    public partial class AlarmMessageObjectPushType : IMessage
    {
        public const byte VariableSpec = 0x12;

        public byte LengthSpec { get; }
        public SyntaxIdType SyntaxId { get; }
        public byte NumberOfValues { get; }
        public uint EventId { get; }
        public State EventState { get; }
        public State LocalState { get; }
        public State AckStateGoing { get; }
        public State AckStateComing { get; }
        public System.Collections.Generic.List<AssociatedValueType> AssociatedValues { get; }

        public AlarmMessageObjectPushType(byte lengthSpec, SyntaxIdType syntaxId, byte numberOfValues, uint eventId, State eventState, State localState, State ackStateGoing, State ackStateComing, System.Collections.Generic.List<AssociatedValueType> associatedValues)
        {
            LengthSpec = lengthSpec;
            SyntaxId = syntaxId;
            NumberOfValues = numberOfValues;
            EventId = eventId;
            EventState = eventState;
            LocalState = localState;
            AckStateGoing = ackStateGoing;
            AckStateComing = ackStateComing;
            AssociatedValues = associatedValues;
        }

        public static AlarmMessageObjectPushType StaticParse(ReadBuffer readBuffer)
        {
            var variableSpec = readBuffer.ReadByte("variableSpec", 8);
            if (!Equals(variableSpec, (byte) (0x12)))
                throw new ParseException($"Expected constant {VariableSpec} for 'variableSpec' but got {variableSpec}");
            var lengthSpec = readBuffer.ReadByte("lengthSpec", 8);
            var syntaxId = (SyntaxIdType) readBuffer.ReadByte("syntaxId", 8);
            var numberOfValues = readBuffer.ReadByte("numberOfValues", 8);
            var eventId = readBuffer.ReadUint("eventId", 32);
            var eventState = State.StaticParse(readBuffer);
            var localState = State.StaticParse(readBuffer);
            var ackStateGoing = State.StaticParse(readBuffer);
            var ackStateComing = State.StaticParse(readBuffer);
            var associatedValues = new System.Collections.Generic.List<AssociatedValueType>();
            var _associatedValuesCnt = (int) (numberOfValues);
            for (var _associatedValuesI = 0; _associatedValuesI < _associatedValuesCnt; _associatedValuesI++)
            {
                associatedValues.Add(AssociatedValueType.StaticParse(readBuffer));
            }
            return new AlarmMessageObjectPushType(lengthSpec, syntaxId, numberOfValues, eventId, eventState, localState, ackStateGoing, ackStateComing, associatedValues);
        }

        public void Serialize(WriteBuffer writeBuffer)
        {
            writeBuffer.WriteByte("variableSpec", 8, VariableSpec);
            writeBuffer.WriteByte("lengthSpec", 8, LengthSpec);
            writeBuffer.WriteByte("syntaxId", 8, (byte) SyntaxId);
            writeBuffer.WriteByte("numberOfValues", 8, NumberOfValues);
            writeBuffer.WriteUint("eventId", 32, EventId);
            EventState.Serialize(writeBuffer);
            LocalState.Serialize(writeBuffer);
            AckStateGoing.Serialize(writeBuffer);
            AckStateComing.Serialize(writeBuffer);
            foreach (var _e in AssociatedValues)
            {
                _e.Serialize(writeBuffer);
            }
        }

        public int GetLengthInBits()
        {
            var lengthInBits = 0;
            lengthInBits += 8;
            lengthInBits += 8;
            lengthInBits += 8;
            lengthInBits += 8;
            lengthInBits += 32;
            lengthInBits += EventState.GetLengthInBits();
            lengthInBits += LocalState.GetLengthInBits();
            lengthInBits += AckStateGoing.GetLengthInBits();
            lengthInBits += AckStateComing.GetLengthInBits();
            lengthInBits += AssociatedValues.Sum(_e => _e.GetLengthInBits());
            return lengthInBits;
        }

        public int GetLengthInBytes() => GetLengthInBits() / 8;
    }
}
