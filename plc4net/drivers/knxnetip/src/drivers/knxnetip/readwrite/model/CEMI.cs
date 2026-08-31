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
    public abstract partial class CEMI : IMessage
    {
        public abstract byte MessageCode { get; }

        public static CEMI StaticParse(ReadBuffer readBuffer, ushort size)
        {
            var messageCode = readBuffer.ReadByte("messageCode", 8);
            if (Equals(messageCode, (byte) (0x2B)))
            {
                return LBusmonInd.StaticParse(readBuffer, size);
            }
            if (Equals(messageCode, (byte) (0x11)))
            {
                return LDataReq.StaticParse(readBuffer, size);
            }
            if (Equals(messageCode, (byte) (0x29)))
            {
                return LDataInd.StaticParse(readBuffer, size);
            }
            if (Equals(messageCode, (byte) (0x2E)))
            {
                return LDataCon.StaticParse(readBuffer, size);
            }
            if (Equals(messageCode, (byte) (0x10)))
            {
                return LRawReq.StaticParse(readBuffer, size);
            }
            if (Equals(messageCode, (byte) (0x2D)))
            {
                return LRawInd.StaticParse(readBuffer, size);
            }
            if (Equals(messageCode, (byte) (0x2F)))
            {
                return LRawCon.StaticParse(readBuffer, size);
            }
            if (Equals(messageCode, (byte) (0x13)))
            {
                return LPollDataReq.StaticParse(readBuffer, size);
            }
            if (Equals(messageCode, (byte) (0x25)))
            {
                return LPollDataCon.StaticParse(readBuffer, size);
            }
            if (Equals(messageCode, (byte) (0x41)))
            {
                return TDataConnectedReq.StaticParse(readBuffer, size);
            }
            if (Equals(messageCode, (byte) (0x89)))
            {
                return TDataConnectedInd.StaticParse(readBuffer, size);
            }
            if (Equals(messageCode, (byte) (0x4A)))
            {
                return TDataIndividualReq.StaticParse(readBuffer, size);
            }
            if (Equals(messageCode, (byte) (0x94)))
            {
                return TDataIndividualInd.StaticParse(readBuffer, size);
            }
            if (Equals(messageCode, (byte) (0xFC)))
            {
                return MPropReadReq.StaticParse(readBuffer, size);
            }
            if (Equals(messageCode, (byte) (0xFB)))
            {
                return MPropReadCon.StaticParse(readBuffer, size);
            }
            if (Equals(messageCode, (byte) (0xF6)))
            {
                return MPropWriteReq.StaticParse(readBuffer, size);
            }
            if (Equals(messageCode, (byte) (0xF5)))
            {
                return MPropWriteCon.StaticParse(readBuffer, size);
            }
            if (Equals(messageCode, (byte) (0xF7)))
            {
                return MPropInfoInd.StaticParse(readBuffer, size);
            }
            if (Equals(messageCode, (byte) (0xF8)))
            {
                return MFuncPropCommandReq.StaticParse(readBuffer, size);
            }
            if (Equals(messageCode, (byte) (0xF9)))
            {
                return MFuncPropStateReadReq.StaticParse(readBuffer, size);
            }
            if (Equals(messageCode, (byte) (0xFA)))
            {
                return MFuncPropCon.StaticParse(readBuffer, size);
            }
            if (Equals(messageCode, (byte) (0xF1)))
            {
                return MResetReq.StaticParse(readBuffer, size);
            }
            if (Equals(messageCode, (byte) (0xF0)))
            {
                return MResetInd.StaticParse(readBuffer, size);
            }
            throw new ParseException("No matching subtype found for CEMI");
        }

        public void Serialize(WriteBuffer writeBuffer)
        {
            writeBuffer.WriteByte("messageCode", 8, MessageCode);
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
