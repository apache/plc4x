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
    public abstract partial class Apdu : IMessage
    {
        public abstract byte Control { get; }

        public bool Numbered { get; }
        public byte Counter { get; }

        protected Apdu(bool numbered, byte counter)
        {
            Numbered = numbered;
            Counter = counter;
        }

        public static Apdu StaticParse(ReadBuffer readBuffer, byte dataLength)
        {
            var control = readBuffer.ReadByte("control", 1);
            var numbered = readBuffer.ReadBit("numbered");
            var counter = readBuffer.ReadByte("counter", 4);
            if (Equals(control, (byte) (1)))
            {
                return ApduControlContainer.StaticParse(readBuffer, dataLength, numbered, counter);
            }
            if (Equals(control, (byte) (0)))
            {
                return ApduDataContainer.StaticParse(readBuffer, dataLength, numbered, counter);
            }
            throw new ParseException("No matching subtype found for Apdu");
        }

        public void Serialize(WriteBuffer writeBuffer)
        {
            writeBuffer.WriteByte("control", 1, Control);
            writeBuffer.WriteBit("numbered", Numbered);
            writeBuffer.WriteByte("counter", 4, Counter);
            SerializeChild(writeBuffer);
        }

        protected abstract void SerializeChild(WriteBuffer writeBuffer);

        public int GetLengthInBits()
        {
            var lengthInBits = 0;
            lengthInBits += 1;
            lengthInBits += 1;
            lengthInBits += 4;
            lengthInBits += GetLengthInBitsChild();
            return lengthInBits;
        }

        protected abstract int GetLengthInBitsChild();

        public int GetLengthInBytes() => (GetLengthInBits() + 7) / 8;
    }
}
