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
    public partial class State : IMessage
    {
        public bool SIG_8 { get; }
        public bool SIG_7 { get; }
        public bool SIG_6 { get; }
        public bool SIG_5 { get; }
        public bool SIG_4 { get; }
        public bool SIG_3 { get; }
        public bool SIG_2 { get; }
        public bool SIG_1 { get; }

        public State(bool sIG_8, bool sIG_7, bool sIG_6, bool sIG_5, bool sIG_4, bool sIG_3, bool sIG_2, bool sIG_1)
        {
            SIG_8 = sIG_8;
            SIG_7 = sIG_7;
            SIG_6 = sIG_6;
            SIG_5 = sIG_5;
            SIG_4 = sIG_4;
            SIG_3 = sIG_3;
            SIG_2 = sIG_2;
            SIG_1 = sIG_1;
        }

        public static State StaticParse(ReadBuffer readBuffer)
        {
            var sIG_8 = readBuffer.ReadBit("SIG_8");
            var sIG_7 = readBuffer.ReadBit("SIG_7");
            var sIG_6 = readBuffer.ReadBit("SIG_6");
            var sIG_5 = readBuffer.ReadBit("SIG_5");
            var sIG_4 = readBuffer.ReadBit("SIG_4");
            var sIG_3 = readBuffer.ReadBit("SIG_3");
            var sIG_2 = readBuffer.ReadBit("SIG_2");
            var sIG_1 = readBuffer.ReadBit("SIG_1");
            return new State(sIG_8, sIG_7, sIG_6, sIG_5, sIG_4, sIG_3, sIG_2, sIG_1);
        }

        public void Serialize(WriteBuffer writeBuffer)
        {
            writeBuffer.WriteBit("SIG_8", SIG_8);
            writeBuffer.WriteBit("SIG_7", SIG_7);
            writeBuffer.WriteBit("SIG_6", SIG_6);
            writeBuffer.WriteBit("SIG_5", SIG_5);
            writeBuffer.WriteBit("SIG_4", SIG_4);
            writeBuffer.WriteBit("SIG_3", SIG_3);
            writeBuffer.WriteBit("SIG_2", SIG_2);
            writeBuffer.WriteBit("SIG_1", SIG_1);
        }

        public int GetLengthInBits()
        {
            var lengthInBits = 0;
            lengthInBits += 1;
            lengthInBits += 1;
            lengthInBits += 1;
            lengthInBits += 1;
            lengthInBits += 1;
            lengthInBits += 1;
            lengthInBits += 1;
            lengthInBits += 1;
            return lengthInBits;
        }

        public int GetLengthInBytes() => GetLengthInBits() / 8;
    }
}
