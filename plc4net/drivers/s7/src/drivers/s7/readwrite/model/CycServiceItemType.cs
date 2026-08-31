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
    public abstract partial class CycServiceItemType : IMessage
    {
        public const byte FunctionId = 0x12;


        public byte ByteLength { get; }
        public byte SyntaxId { get; }

        protected CycServiceItemType(byte byteLength, byte syntaxId)
        {
            ByteLength = byteLength;
            SyntaxId = syntaxId;
        }

        public static CycServiceItemType StaticParse(ReadBuffer readBuffer)
        {
            var functionId = readBuffer.ReadByte("functionId", 8);
            if (!Equals(functionId, (byte) (0x12)))
                throw new ParseException($"Expected constant {FunctionId} for 'functionId' but got {functionId}");
            var byteLength = readBuffer.ReadByte("byteLength", 8);
            var syntaxId = readBuffer.ReadByte("syntaxId", 8);
            if (Equals(syntaxId, (byte) (0x10)))
            {
                return CycServiceItemAnyType.StaticParse(readBuffer, byteLength, syntaxId);
            }
            if (Equals(syntaxId, (byte) (0xb0)))
            {
                return CycServiceItemDbReadType.StaticParse(readBuffer, byteLength, syntaxId);
            }
            throw new ParseException("No matching subtype found for CycServiceItemType");
        }

        public void Serialize(WriteBuffer writeBuffer)
        {
            writeBuffer.WriteByte("functionId", 8, FunctionId);
            writeBuffer.WriteByte("byteLength", 8, ByteLength);
            writeBuffer.WriteByte("syntaxId", 8, SyntaxId);
            SerializeChild(writeBuffer);
        }

        protected abstract void SerializeChild(WriteBuffer writeBuffer);

        public int GetLengthInBits()
        {
            var lengthInBits = 0;
            lengthInBits += 8;
            lengthInBits += 8;
            lengthInBits += 8;
            lengthInBits += GetLengthInBitsChild();
            return lengthInBits;
        }

        protected abstract int GetLengthInBitsChild();

        public int GetLengthInBytes() => GetLengthInBits() / 8;
    }
}
