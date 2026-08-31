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
    public partial class AlarmMessageObjectQueryType : IMessage
    {
        public const byte VariableSpec = 0x12;

        public byte LengthDataset { get; }
        public State EventState { get; }
        public State AckStateGoing { get; }
        public State AckStateComing { get; }
        public DateAndTime TimeComing { get; }
        public AssociatedValueType ValueComing { get; }
        public DateAndTime TimeGoing { get; }
        public AssociatedValueType ValueGoing { get; }

        public AlarmMessageObjectQueryType(byte lengthDataset, State eventState, State ackStateGoing, State ackStateComing, DateAndTime timeComing, AssociatedValueType valueComing, DateAndTime timeGoing, AssociatedValueType valueGoing)
        {
            LengthDataset = lengthDataset;
            EventState = eventState;
            AckStateGoing = ackStateGoing;
            AckStateComing = ackStateComing;
            TimeComing = timeComing;
            ValueComing = valueComing;
            TimeGoing = timeGoing;
            ValueGoing = valueGoing;
        }

        public static AlarmMessageObjectQueryType StaticParse(ReadBuffer readBuffer)
        {
            var lengthDataset = readBuffer.ReadByte("lengthDataset", 8);
            {
                var reserved = readBuffer.ReadUshort("reserved", 16);
                if (!Equals(reserved, (ushort) (0x0000))) { /* mspec reserved: value differs from the spec default */ }
            }
            var variableSpec = readBuffer.ReadByte("variableSpec", 8);
            if (!Equals(variableSpec, (byte) (0x12)))
                throw new ParseException($"Expected constant {VariableSpec} for 'variableSpec' but got {variableSpec}");
            var eventState = State.StaticParse(readBuffer);
            var ackStateGoing = State.StaticParse(readBuffer);
            var ackStateComing = State.StaticParse(readBuffer);
            var timeComing = DateAndTime.StaticParse(readBuffer);
            var valueComing = AssociatedValueType.StaticParse(readBuffer);
            var timeGoing = DateAndTime.StaticParse(readBuffer);
            var valueGoing = AssociatedValueType.StaticParse(readBuffer);
            return new AlarmMessageObjectQueryType(lengthDataset, eventState, ackStateGoing, ackStateComing, timeComing, valueComing, timeGoing, valueGoing);
        }

        public void Serialize(WriteBuffer writeBuffer)
        {
            writeBuffer.WriteByte("lengthDataset", 8, LengthDataset);
            writeBuffer.WriteUshort("reserved", 16, (ushort) (0x0000));
            writeBuffer.WriteByte("variableSpec", 8, VariableSpec);
            EventState.Serialize(writeBuffer);
            AckStateGoing.Serialize(writeBuffer);
            AckStateComing.Serialize(writeBuffer);
            TimeComing.Serialize(writeBuffer);
            ValueComing.Serialize(writeBuffer);
            TimeGoing.Serialize(writeBuffer);
            ValueGoing.Serialize(writeBuffer);
        }

        public int GetLengthInBits()
        {
            var lengthInBits = 0;
            lengthInBits += 8;
            lengthInBits += 16;
            lengthInBits += 8;
            lengthInBits += EventState.GetLengthInBits();
            lengthInBits += AckStateGoing.GetLengthInBits();
            lengthInBits += AckStateComing.GetLengthInBits();
            lengthInBits += TimeComing.GetLengthInBits();
            lengthInBits += ValueComing.GetLengthInBits();
            lengthInBits += TimeGoing.GetLengthInBits();
            lengthInBits += ValueGoing.GetLengthInBits();
            return lengthInBits;
        }

        public int GetLengthInBytes() => GetLengthInBits() / 8;
    }
}
