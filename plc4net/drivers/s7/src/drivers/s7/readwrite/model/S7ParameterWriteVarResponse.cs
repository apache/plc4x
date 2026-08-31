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
    public partial class S7ParameterWriteVarResponse : S7Parameter
    {
        public override byte ParameterType => (byte) (0x05);
        public override byte MessageType => (byte) (0x03);

        public byte NumItems { get; }

        public S7ParameterWriteVarResponse(byte numItems)
        {
            NumItems = numItems;
        }

        public static new S7ParameterWriteVarResponse StaticParse(ReadBuffer readBuffer, byte messageType)
        {
            var numItems = readBuffer.ReadByte("numItems", 8);
            return new S7ParameterWriteVarResponse(numItems);
        }

        protected override void SerializeChild(WriteBuffer writeBuffer)
        {
            writeBuffer.WriteByte("numItems", 8, NumItems);
        }

        protected override int GetLengthInBitsChild()
        {
            var lengthInBits = 0;
            lengthInBits += 8;
            return lengthInBits;
        }

    }
}
