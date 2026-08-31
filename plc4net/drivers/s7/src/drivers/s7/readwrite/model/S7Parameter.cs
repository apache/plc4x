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
    public abstract partial class S7Parameter : IMessage
    {
        public abstract byte ParameterType { get; }
        public abstract byte MessageType { get; }

        public static S7Parameter StaticParse(ReadBuffer readBuffer, byte messageType)
        {
            var parameterType = readBuffer.ReadByte("parameterType", 8);
            if (Equals(parameterType, (byte) (0xF0)))
            {
                return S7ParameterSetupCommunication.StaticParse(readBuffer, messageType);
            }
            if (Equals(parameterType, (byte) (0x04)) && Equals(messageType, (byte) (0x01)))
            {
                return S7ParameterReadVarRequest.StaticParse(readBuffer, messageType);
            }
            if (Equals(parameterType, (byte) (0x04)) && Equals(messageType, (byte) (0x03)))
            {
                return S7ParameterReadVarResponse.StaticParse(readBuffer, messageType);
            }
            if (Equals(parameterType, (byte) (0x05)) && Equals(messageType, (byte) (0x01)))
            {
                return S7ParameterWriteVarRequest.StaticParse(readBuffer, messageType);
            }
            if (Equals(parameterType, (byte) (0x05)) && Equals(messageType, (byte) (0x03)))
            {
                return S7ParameterWriteVarResponse.StaticParse(readBuffer, messageType);
            }
            if (Equals(parameterType, (byte) (0x00)) && Equals(messageType, (byte) (0x07)))
            {
                return S7ParameterUserData.StaticParse(readBuffer, messageType);
            }
            if (Equals(parameterType, (byte) (0x01)) && Equals(messageType, (byte) (0x07)))
            {
                return S7ParameterModeTransition.StaticParse(readBuffer, messageType);
            }
            throw new ParseException("No matching subtype found for S7Parameter");
        }

        public void Serialize(WriteBuffer writeBuffer)
        {
            writeBuffer.WriteByte("parameterType", 8, ParameterType);
            SerializeChild(writeBuffer);
        }

        protected abstract void SerializeChild(WriteBuffer writeBuffer);

        public int GetLengthInBits()
        {
            var lengthInBits = 0;
            lengthInBits += 8;
            lengthInBits += GetLengthInBitsChild();
            return lengthInBits;
        }

        protected abstract int GetLengthInBitsChild();

        public int GetLengthInBytes() => GetLengthInBits() / 8;
    }
}
