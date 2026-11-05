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
    public partial class S7PayloadUserDataItemCpuFunctionMsgSubscriptionRequest : S7PayloadUserDataItem
    {
        public override byte CpuFunctionGroup => (byte) (0x04);
        public override byte CpuFunctionType => (byte) (0x04);
        public override byte CpuSubfunction => (byte) (0x02);

        public byte Subscription { get; }
        public string MagicKey { get; }
        public AlarmStateType? Alarmtype { get; }
        public byte? Reserve { get; }

        public S7PayloadUserDataItemCpuFunctionMsgSubscriptionRequest(DataTransportErrorCode returnCode, DataTransportSize transportSize, ushort dataLength, byte subscription, string magicKey, AlarmStateType? alarmtype, byte? reserve) : base(returnCode, transportSize, dataLength)
        {
            Subscription = subscription;
            MagicKey = magicKey;
            Alarmtype = alarmtype;
            Reserve = reserve;
        }

        public static S7PayloadUserDataItemCpuFunctionMsgSubscriptionRequest StaticParse(ReadBuffer readBuffer, byte cpuFunctionGroup, byte cpuFunctionType, byte cpuSubfunction, DataTransportErrorCode returnCode, DataTransportSize transportSize, ushort dataLength)
        {
            var subscription = readBuffer.ReadByte("subscription", 8);
            {
                var reserved = readBuffer.ReadByte("reserved", 8);
                if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
            }
            var magicKey = readBuffer.ReadString("magicKey", 64, System.Text.Encoding.UTF8);
            AlarmStateType? alarmtype = null;
            if ((dataLength >= 12))
            {
                alarmtype = (AlarmStateType) readBuffer.ReadByte("alarmtype", 8);
            }
            byte? reserve = null;
            if ((dataLength >= 12))
            {
                reserve = readBuffer.ReadByte("reserve", 8);
            }
            return new S7PayloadUserDataItemCpuFunctionMsgSubscriptionRequest(returnCode, transportSize, dataLength, subscription, magicKey, alarmtype, reserve);
        }

        protected override void SerializeChild(WriteBuffer writeBuffer)
        {
            writeBuffer.WriteByte("subscription", 8, Subscription);
            writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
            writeBuffer.WriteString("magicKey", 64, "UTF8", MagicKey);
            if (Alarmtype != null)
            {
                writeBuffer.WriteByte("alarmtype", 8, (byte) Alarmtype.Value);
            }
            if (Reserve != null)
            {
                writeBuffer.WriteByte("reserve", 8, Reserve.Value);
            }
        }

        protected override int GetLengthInBitsChild()
        {
            var lengthInBits = 0;
            lengthInBits += 8;
            lengthInBits += 8;
            lengthInBits += 64;
            lengthInBits += (Alarmtype != null ? 8 : 0);
            lengthInBits += (Reserve != null ? 8 : 0);
            return lengthInBits;
        }

    }
}
