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

namespace org.apache.plc4net.drivers.knxnetip.readwrite.model
{
    public abstract partial class LDataFrame : IMessage
    {
        public abstract bool Polling { get; }
        public abstract bool NotAckFrame { get; }

        public bool FrameType { get; }
        public bool NotRepeated { get; }
        public CEMIPriority Priority { get; }
        public bool AcknowledgeRequested { get; }
        public bool ErrorFlag { get; }

        protected LDataFrame(bool frameType, bool notRepeated, CEMIPriority priority, bool acknowledgeRequested, bool errorFlag)
        {
            FrameType = frameType;
            NotRepeated = notRepeated;
            Priority = priority;
            AcknowledgeRequested = acknowledgeRequested;
            ErrorFlag = errorFlag;
        }

        public static LDataFrame StaticParse(ReadBuffer readBuffer)
        {
            var frameType = readBuffer.ReadBit("frameType");
            var polling = readBuffer.ReadBit("polling");
            var notRepeated = readBuffer.ReadBit("notRepeated");
            var notAckFrame = readBuffer.ReadBit("notAckFrame");
            var priority = (CEMIPriority) readBuffer.ReadByte("priority", 2);
            var acknowledgeRequested = readBuffer.ReadBit("acknowledgeRequested");
            var errorFlag = readBuffer.ReadBit("errorFlag");
            if (Equals(notAckFrame, true) && Equals(polling, false))
            {
                return LDataExtended.StaticParse(readBuffer, frameType, notRepeated, priority, acknowledgeRequested, errorFlag);
            }
            if (Equals(notAckFrame, true) && Equals(polling, true))
            {
                return LPollData.StaticParse(readBuffer, frameType, notRepeated, priority, acknowledgeRequested, errorFlag);
            }
            if (Equals(notAckFrame, false))
            {
                return LDataFrameACK.StaticParse(readBuffer, frameType, notRepeated, priority, acknowledgeRequested, errorFlag);
            }
            throw new ParseException("No matching subtype found for LDataFrame");
        }

        public void Serialize(WriteBuffer writeBuffer)
        {
            writeBuffer.WriteBit("frameType", FrameType);
            writeBuffer.WriteBit("polling", Polling);
            writeBuffer.WriteBit("notRepeated", NotRepeated);
            writeBuffer.WriteBit("notAckFrame", NotAckFrame);
            writeBuffer.WriteByte("priority", 2, (byte) Priority);
            writeBuffer.WriteBit("acknowledgeRequested", AcknowledgeRequested);
            writeBuffer.WriteBit("errorFlag", ErrorFlag);
            SerializeChild(writeBuffer);
        }

        protected abstract void SerializeChild(WriteBuffer writeBuffer);

        public int GetLengthInBits()
        {
            var lengthInBits = 0;
            lengthInBits += 1;
            lengthInBits += 1;
            lengthInBits += 1;
            lengthInBits += 1;
            lengthInBits += 2;
            lengthInBits += 1;
            lengthInBits += 1;
            lengthInBits += GetLengthInBitsChild();
            return lengthInBits;
        }

        protected abstract int GetLengthInBitsChild();

        public int GetLengthInBytes() => GetLengthInBits() / 8;
    }
}
