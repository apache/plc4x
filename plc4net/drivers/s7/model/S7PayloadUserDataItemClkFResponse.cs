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
    public partial class S7PayloadUserDataItemClkFResponse : S7PayloadUserDataItem
    {
        public override byte CpuFunctionGroup => (byte) (0x07);
        public override byte CpuFunctionType => (byte) (0x08);
        public override byte CpuSubfunction => (byte) (0x03);

        public byte Res { get; }
        public byte Year1 { get; }
        public DateAndTime TimeStamp { get; }

        public S7PayloadUserDataItemClkFResponse(DataTransportErrorCode returnCode, DataTransportSize transportSize, ushort dataLength, byte res, byte year1, DateAndTime timeStamp) : base(returnCode, transportSize, dataLength)
        {
            Res = res;
            Year1 = year1;
            TimeStamp = timeStamp;
        }

        public static S7PayloadUserDataItemClkFResponse StaticParse(ReadBuffer readBuffer, byte cpuFunctionGroup, byte cpuFunctionType, byte cpuSubfunction, DataTransportErrorCode returnCode, DataTransportSize transportSize, ushort dataLength)
        {
            var res = readBuffer.ReadByte("res", 8);
            var year1 = readBuffer.ReadByte("year1", 8);
            var timeStamp = DateAndTime.StaticParse(readBuffer);
            return new S7PayloadUserDataItemClkFResponse(returnCode, transportSize, dataLength, res, year1, timeStamp);
        }

        protected override void SerializeChild(WriteBuffer writeBuffer)
        {
            writeBuffer.WriteByte("res", 8, Res);
            writeBuffer.WriteByte("year1", 8, Year1);
            TimeStamp.Serialize(writeBuffer);
        }

        protected override int GetLengthInBitsChild()
        {
            var lengthInBits = 0;
            lengthInBits += 8;
            lengthInBits += 8;
            lengthInBits += TimeStamp.GetLengthInBits();
            return lengthInBits;
        }

    }
}
