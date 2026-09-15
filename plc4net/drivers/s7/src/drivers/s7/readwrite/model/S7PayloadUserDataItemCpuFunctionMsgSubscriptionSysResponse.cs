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
    public partial class S7PayloadUserDataItemCpuFunctionMsgSubscriptionSysResponse : S7PayloadUserDataItem
    {
        public override byte CpuFunctionGroup => (byte) (0x04);
        public override byte CpuFunctionType => (byte) (0x08);
        public override byte CpuSubfunction => (byte) (0x02);

        public byte Result { get; }
        public byte Reserved01 { get; }

        public S7PayloadUserDataItemCpuFunctionMsgSubscriptionSysResponse(DataTransportErrorCode returnCode, DataTransportSize transportSize, ushort dataLength, byte result, byte reserved01) : base(returnCode, transportSize, dataLength)
        {
            Result = result;
            Reserved01 = reserved01;
        }

        public static S7PayloadUserDataItemCpuFunctionMsgSubscriptionSysResponse StaticParse(ReadBuffer readBuffer, byte cpuFunctionGroup, byte cpuFunctionType, byte cpuSubfunction, DataTransportErrorCode returnCode, DataTransportSize transportSize, ushort dataLength)
        {
            var result = readBuffer.ReadByte("result", 8);
            var reserved01 = readBuffer.ReadByte("reserved01", 8);
            return new S7PayloadUserDataItemCpuFunctionMsgSubscriptionSysResponse(returnCode, transportSize, dataLength, result, reserved01);
        }

        protected override void SerializeChild(WriteBuffer writeBuffer)
        {
            writeBuffer.WriteByte("result", 8, Result);
            writeBuffer.WriteByte("reserved01", 8, Reserved01);
        }

        protected override int GetLengthInBitsChild()
        {
            var lengthInBits = 0;
            lengthInBits += 8;
            lengthInBits += 8;
            return lengthInBits;
        }

    }
}
