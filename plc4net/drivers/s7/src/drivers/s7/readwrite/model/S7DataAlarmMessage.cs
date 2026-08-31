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
    public abstract partial class S7DataAlarmMessage : IMessage
    {
        public const byte FunctionId = 0x00;
        public const byte NumberMessageObj = 0x01;

        public abstract byte CpuFunctionType { get; }

        public static S7DataAlarmMessage StaticParse(ReadBuffer readBuffer, byte cpuFunctionType)
        {
            var functionId = readBuffer.ReadByte("functionId", 8);
            if (!Equals(functionId, (byte) (0x00)))
                throw new ParseException($"Expected constant {FunctionId} for 'functionId' but got {functionId}");
            var numberMessageObj = readBuffer.ReadByte("numberMessageObj", 8);
            if (!Equals(numberMessageObj, (byte) (0x01)))
                throw new ParseException($"Expected constant {NumberMessageObj} for 'numberMessageObj' but got {numberMessageObj}");
            if (Equals(cpuFunctionType, (byte) (0x04)))
            {
                return S7MessageObjectRequest.StaticParse(readBuffer, cpuFunctionType);
            }
            if (Equals(cpuFunctionType, (byte) (0x08)))
            {
                return S7MessageObjectResponse.StaticParse(readBuffer, cpuFunctionType);
            }
            throw new ParseException("No matching subtype found for S7DataAlarmMessage");
        }

        public void Serialize(WriteBuffer writeBuffer)
        {
            writeBuffer.WriteByte("functionId", 8, FunctionId);
            writeBuffer.WriteByte("numberMessageObj", 8, NumberMessageObj);
            SerializeChild(writeBuffer);
        }

        protected abstract void SerializeChild(WriteBuffer writeBuffer);

        public int GetLengthInBits()
        {
            var lengthInBits = 0;
            lengthInBits += 8;
            lengthInBits += 8;
            lengthInBits += GetLengthInBitsChild();
            return lengthInBits;
        }

        protected abstract int GetLengthInBitsChild();

        public int GetLengthInBytes() => GetLengthInBits() / 8;
    }
}
