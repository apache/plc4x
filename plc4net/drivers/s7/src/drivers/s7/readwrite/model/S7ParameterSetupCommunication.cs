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
    public partial class S7ParameterSetupCommunication : S7Parameter
    {
        public override byte ParameterType => (byte) (0xF0);
        public override byte MessageType => default(byte);

        public ushort MaxAmqCaller { get; }
        public ushort MaxAmqCallee { get; }
        public ushort PduLength { get; }

        public S7ParameterSetupCommunication(ushort maxAmqCaller, ushort maxAmqCallee, ushort pduLength)
        {
            MaxAmqCaller = maxAmqCaller;
            MaxAmqCallee = maxAmqCallee;
            PduLength = pduLength;
        }

        public static new S7ParameterSetupCommunication StaticParse(ReadBuffer readBuffer, byte messageType)
        {
            {
                var reserved = readBuffer.ReadByte("reserved", 8);
                if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
            }
            var maxAmqCaller = readBuffer.ReadUshort("maxAmqCaller", 16);
            var maxAmqCallee = readBuffer.ReadUshort("maxAmqCallee", 16);
            var pduLength = readBuffer.ReadUshort("pduLength", 16);
            return new S7ParameterSetupCommunication(maxAmqCaller, maxAmqCallee, pduLength);
        }

        protected override void SerializeChild(WriteBuffer writeBuffer)
        {
            writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
            writeBuffer.WriteUshort("maxAmqCaller", 16, MaxAmqCaller);
            writeBuffer.WriteUshort("maxAmqCallee", 16, MaxAmqCallee);
            writeBuffer.WriteUshort("pduLength", 16, PduLength);
        }

        protected override int GetLengthInBitsChild()
        {
            var lengthInBits = 0;
            lengthInBits += 8;
            lengthInBits += 16;
            lengthInBits += 16;
            lengthInBits += 16;
            return lengthInBits;
        }

    }
}
