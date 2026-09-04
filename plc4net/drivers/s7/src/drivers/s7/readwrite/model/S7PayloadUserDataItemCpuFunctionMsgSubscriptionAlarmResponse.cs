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
    public partial class S7PayloadUserDataItemCpuFunctionMsgSubscriptionAlarmResponse : S7PayloadUserDataItem
    {
        public override byte CpuFunctionGroup => (byte) (0x04);
        public override byte CpuFunctionType => (byte) (0x08);
        public override byte CpuSubfunction => (byte) (0x02);

        public byte Result { get; }
        public byte Reserved01 { get; }
        public AlarmStateType AlarmType { get; }
        public byte Reserved02 { get; }
        public byte Reserved03 { get; }

        public S7PayloadUserDataItemCpuFunctionMsgSubscriptionAlarmResponse(DataTransportErrorCode returnCode, DataTransportSize transportSize, ushort dataLength, byte result, byte reserved01, AlarmStateType alarmType, byte reserved02, byte reserved03) : base(returnCode, transportSize, dataLength)
        {
            Result = result;
            Reserved01 = reserved01;
            AlarmType = alarmType;
            Reserved02 = reserved02;
            Reserved03 = reserved03;
        }

        public static S7PayloadUserDataItemCpuFunctionMsgSubscriptionAlarmResponse StaticParse(ReadBuffer readBuffer, byte cpuFunctionGroup, byte cpuFunctionType, byte cpuSubfunction, DataTransportErrorCode returnCode, DataTransportSize transportSize, ushort dataLength)
        {
            var result = readBuffer.ReadByte("result", 8);
            var reserved01 = readBuffer.ReadByte("reserved01", 8);
            var alarmType = (AlarmStateType) readBuffer.ReadByte("alarmType", 8);
            var reserved02 = readBuffer.ReadByte("reserved02", 8);
            var reserved03 = readBuffer.ReadByte("reserved03", 8);
            return new S7PayloadUserDataItemCpuFunctionMsgSubscriptionAlarmResponse(returnCode, transportSize, dataLength, result, reserved01, alarmType, reserved02, reserved03);
        }

        protected override void SerializeChild(WriteBuffer writeBuffer)
        {
            writeBuffer.WriteByte("result", 8, Result);
            writeBuffer.WriteByte("reserved01", 8, Reserved01);
            writeBuffer.WriteByte("alarmType", 8, (byte) AlarmType);
            writeBuffer.WriteByte("reserved02", 8, Reserved02);
            writeBuffer.WriteByte("reserved03", 8, Reserved03);
        }

        protected override int GetLengthInBitsChild()
        {
            var lengthInBits = 0;
            lengthInBits += 8;
            lengthInBits += 8;
            lengthInBits += 8;
            lengthInBits += 8;
            lengthInBits += 8;
            return lengthInBits;
        }

    }
}
