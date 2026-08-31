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
    public abstract partial class S7PayloadUserDataItem : IMessage
    {
        public abstract byte CpuFunctionGroup { get; }
        public abstract byte CpuFunctionType { get; }
        public abstract byte CpuSubfunction { get; }

        public DataTransportErrorCode ReturnCode { get; }
        public DataTransportSize TransportSize { get; }
        public ushort DataLength { get; }

        protected S7PayloadUserDataItem(DataTransportErrorCode returnCode, DataTransportSize transportSize, ushort dataLength)
        {
            ReturnCode = returnCode;
            TransportSize = transportSize;
            DataLength = dataLength;
        }

        public static S7PayloadUserDataItem StaticParse(ReadBuffer readBuffer, byte cpuFunctionGroup, byte cpuFunctionType, byte cpuSubfunction)
        {
            var returnCode = (DataTransportErrorCode) readBuffer.ReadByte("returnCode", 8);
            var transportSize = (DataTransportSize) readBuffer.ReadByte("transportSize", 8);
            var dataLength = readBuffer.ReadUshort("dataLength", 16);
            if (Equals(cpuFunctionGroup, (byte) (0x02)) && Equals(cpuFunctionType, (byte) (0x00)) && Equals(cpuSubfunction, (byte) (0x01)))
            {
                return S7PayloadUserDataItemCyclicServicesPush.StaticParse(readBuffer, cpuFunctionGroup, cpuFunctionType, cpuSubfunction, returnCode, transportSize, dataLength);
            }
            if (Equals(cpuFunctionGroup, (byte) (0x02)) && Equals(cpuFunctionType, (byte) (0x00)) && Equals(cpuSubfunction, (byte) (0x05)))
            {
                return S7PayloadUserDataItemCyclicServicesChangeDrivenPush.StaticParse(readBuffer, cpuFunctionGroup, cpuFunctionType, cpuSubfunction, returnCode, transportSize, dataLength);
            }
            if (Equals(cpuFunctionGroup, (byte) (0x02)) && Equals(cpuFunctionType, (byte) (0x04)) && Equals(cpuSubfunction, (byte) (0x01)))
            {
                return S7PayloadUserDataItemCyclicServicesSubscribeRequest.StaticParse(readBuffer, cpuFunctionGroup, cpuFunctionType, cpuSubfunction, returnCode, transportSize, dataLength);
            }
            if (Equals(cpuFunctionGroup, (byte) (0x02)) && Equals(cpuFunctionType, (byte) (0x04)) && Equals(cpuSubfunction, (byte) (0x04)))
            {
                return S7PayloadUserDataItemCyclicServicesUnsubscribeRequest.StaticParse(readBuffer, cpuFunctionGroup, cpuFunctionType, cpuSubfunction, returnCode, transportSize, dataLength);
            }
            if (Equals(cpuFunctionGroup, (byte) (0x02)) && Equals(cpuFunctionType, (byte) (0x08)) && Equals(cpuSubfunction, (byte) (0x01)) && Equals(dataLength, (ushort) (0x00)))
            {
                return S7PayloadUserDataItemCyclicServicesSubscribeEmptyResponse.StaticParse(readBuffer, cpuFunctionGroup, cpuFunctionType, cpuSubfunction, returnCode, transportSize, dataLength);
            }
            if (Equals(cpuFunctionGroup, (byte) (0x02)) && Equals(cpuFunctionType, (byte) (0x08)) && Equals(cpuSubfunction, (byte) (0x01)))
            {
                return S7PayloadUserDataItemCyclicServicesSubscribeResponse.StaticParse(readBuffer, cpuFunctionGroup, cpuFunctionType, cpuSubfunction, returnCode, transportSize, dataLength);
            }
            if (Equals(cpuFunctionGroup, (byte) (0x02)) && Equals(cpuFunctionType, (byte) (0x08)) && Equals(cpuSubfunction, (byte) (0x04)))
            {
                return S7PayloadUserDataItemCyclicServicesUnsubscribeResponse.StaticParse(readBuffer, cpuFunctionGroup, cpuFunctionType, cpuSubfunction, returnCode, transportSize, dataLength);
            }
            if (Equals(cpuFunctionGroup, (byte) (0x02)) && Equals(cpuFunctionType, (byte) (0x08)) && Equals(cpuSubfunction, (byte) (0x05)) && Equals(dataLength, (ushort) (0x00)))
            {
                return S7PayloadUserDataItemCyclicServicesErrorResponse.StaticParse(readBuffer, cpuFunctionGroup, cpuFunctionType, cpuSubfunction, returnCode, transportSize, dataLength);
            }
            if (Equals(cpuFunctionGroup, (byte) (0x02)) && Equals(cpuFunctionType, (byte) (0x08)) && Equals(cpuSubfunction, (byte) (0x05)))
            {
                return S7PayloadUserDataItemCyclicServicesChangeDrivenSubscribeResponse.StaticParse(readBuffer, cpuFunctionGroup, cpuFunctionType, cpuSubfunction, returnCode, transportSize, dataLength);
            }
            if (Equals(cpuFunctionGroup, (byte) (0x03)) && Equals(cpuFunctionType, (byte) (0x04)) && Equals(cpuSubfunction, (byte) (0x01)))
            {
                return S7PayloadUserDataItemCpuFunctionListBlocksRequest.StaticParse(readBuffer, cpuFunctionGroup, cpuFunctionType, cpuSubfunction, returnCode, transportSize, dataLength);
            }
            if (Equals(cpuFunctionGroup, (byte) (0x03)) && Equals(cpuFunctionType, (byte) (0x08)) && Equals(cpuSubfunction, (byte) (0x01)))
            {
                return S7PayloadUserDataItemCpuFunctionListBlocksResponse.StaticParse(readBuffer, cpuFunctionGroup, cpuFunctionType, cpuSubfunction, returnCode, transportSize, dataLength);
            }
            if (Equals(cpuFunctionGroup, (byte) (0x03)) && Equals(cpuFunctionType, (byte) (0x04)) && Equals(cpuSubfunction, (byte) (0x02)))
            {
                return S7PayloadUserDataItemCpuFunctionListBlocksOfTypeRequest.StaticParse(readBuffer, cpuFunctionGroup, cpuFunctionType, cpuSubfunction, returnCode, transportSize, dataLength);
            }
            if (Equals(cpuFunctionGroup, (byte) (0x03)) && Equals(cpuFunctionType, (byte) (0x08)) && Equals(cpuSubfunction, (byte) (0x02)))
            {
                return S7PayloadUserDataItemCpuFunctionListBlocksOfTypeResponse.StaticParse(readBuffer, cpuFunctionGroup, cpuFunctionType, cpuSubfunction, returnCode, transportSize, dataLength);
            }
            if (Equals(cpuFunctionGroup, (byte) (0x03)) && Equals(cpuFunctionType, (byte) (0x04)) && Equals(cpuSubfunction, (byte) (0x03)))
            {
                return S7PayloadUserDataItemCpuFunctionGetBlockInfoRequest.StaticParse(readBuffer, cpuFunctionGroup, cpuFunctionType, cpuSubfunction, returnCode, transportSize, dataLength);
            }
            if (Equals(cpuFunctionGroup, (byte) (0x03)) && Equals(cpuFunctionType, (byte) (0x08)) && Equals(cpuSubfunction, (byte) (0x03)))
            {
                return S7PayloadUserDataItemCpuFunctionGetBlockInfoResponse.StaticParse(readBuffer, cpuFunctionGroup, cpuFunctionType, cpuSubfunction, returnCode, transportSize, dataLength);
            }
            if (Equals(cpuFunctionGroup, (byte) (0x04)) && Equals(cpuFunctionType, (byte) (0x00)) && Equals(cpuSubfunction, (byte) (0x03)))
            {
                return S7PayloadDiagnosticMessage.StaticParse(readBuffer, cpuFunctionGroup, cpuFunctionType, cpuSubfunction, returnCode, transportSize, dataLength);
            }
            if (Equals(cpuFunctionGroup, (byte) (0x04)) && Equals(cpuFunctionType, (byte) (0x00)) && Equals(cpuSubfunction, (byte) (0x05)))
            {
                return S7PayloadAlarm8.StaticParse(readBuffer, cpuFunctionGroup, cpuFunctionType, cpuSubfunction, returnCode, transportSize, dataLength);
            }
            if (Equals(cpuFunctionGroup, (byte) (0x04)) && Equals(cpuFunctionType, (byte) (0x00)) && Equals(cpuSubfunction, (byte) (0x06)))
            {
                return S7PayloadNotify.StaticParse(readBuffer, cpuFunctionGroup, cpuFunctionType, cpuSubfunction, returnCode, transportSize, dataLength);
            }
            if (Equals(cpuFunctionGroup, (byte) (0x04)) && Equals(cpuFunctionType, (byte) (0x00)) && Equals(cpuSubfunction, (byte) (0x0c)))
            {
                return S7PayloadAlarmAckInd.StaticParse(readBuffer, cpuFunctionGroup, cpuFunctionType, cpuSubfunction, returnCode, transportSize, dataLength);
            }
            if (Equals(cpuFunctionGroup, (byte) (0x04)) && Equals(cpuFunctionType, (byte) (0x00)) && Equals(cpuSubfunction, (byte) (0x11)))
            {
                return S7PayloadAlarmSQ.StaticParse(readBuffer, cpuFunctionGroup, cpuFunctionType, cpuSubfunction, returnCode, transportSize, dataLength);
            }
            if (Equals(cpuFunctionGroup, (byte) (0x04)) && Equals(cpuFunctionType, (byte) (0x00)) && Equals(cpuSubfunction, (byte) (0x12)))
            {
                return S7PayloadAlarmS.StaticParse(readBuffer, cpuFunctionGroup, cpuFunctionType, cpuSubfunction, returnCode, transportSize, dataLength);
            }
            if (Equals(cpuFunctionGroup, (byte) (0x04)) && Equals(cpuFunctionType, (byte) (0x00)) && Equals(cpuSubfunction, (byte) (0x13)))
            {
                return S7PayloadAlarmSC.StaticParse(readBuffer, cpuFunctionGroup, cpuFunctionType, cpuSubfunction, returnCode, transportSize, dataLength);
            }
            if (Equals(cpuFunctionGroup, (byte) (0x04)) && Equals(cpuFunctionType, (byte) (0x00)) && Equals(cpuSubfunction, (byte) (0x16)))
            {
                return S7PayloadNotify8.StaticParse(readBuffer, cpuFunctionGroup, cpuFunctionType, cpuSubfunction, returnCode, transportSize, dataLength);
            }
            if (Equals(cpuFunctionGroup, (byte) (0x04)) && Equals(cpuFunctionType, (byte) (0x04)) && Equals(cpuSubfunction, (byte) (0x01)) && Equals(dataLength, (ushort) (0x00)))
            {
                return S7PayloadUserDataItemCpuFunctionReadSzlNoDataRequest.StaticParse(readBuffer, cpuFunctionGroup, cpuFunctionType, cpuSubfunction, returnCode, transportSize, dataLength);
            }
            if (Equals(cpuFunctionGroup, (byte) (0x04)) && Equals(cpuFunctionType, (byte) (0x04)) && Equals(cpuSubfunction, (byte) (0x01)))
            {
                return S7PayloadUserDataItemCpuFunctionReadSzlRequest.StaticParse(readBuffer, cpuFunctionGroup, cpuFunctionType, cpuSubfunction, returnCode, transportSize, dataLength);
            }
            if (Equals(cpuFunctionGroup, (byte) (0x04)) && Equals(cpuFunctionType, (byte) (0x08)) && Equals(cpuSubfunction, (byte) (0x01)))
            {
                return S7PayloadUserDataItemCpuFunctionReadSzlResponse.StaticParse(readBuffer, cpuFunctionGroup, cpuFunctionType, cpuSubfunction, returnCode, transportSize, dataLength);
            }
            if (Equals(cpuFunctionGroup, (byte) (0x04)) && Equals(cpuFunctionType, (byte) (0x04)) && Equals(cpuSubfunction, (byte) (0x02)))
            {
                return S7PayloadUserDataItemCpuFunctionMsgSubscriptionRequest.StaticParse(readBuffer, cpuFunctionGroup, cpuFunctionType, cpuSubfunction, returnCode, transportSize, dataLength);
            }
            if (Equals(cpuFunctionGroup, (byte) (0x04)) && Equals(cpuFunctionType, (byte) (0x08)) && Equals(cpuSubfunction, (byte) (0x02)) && Equals(dataLength, (ushort) (0x00)))
            {
                return S7PayloadUserDataItemCpuFunctionMsgSubscriptionResponse.StaticParse(readBuffer, cpuFunctionGroup, cpuFunctionType, cpuSubfunction, returnCode, transportSize, dataLength);
            }
            if (Equals(cpuFunctionGroup, (byte) (0x04)) && Equals(cpuFunctionType, (byte) (0x08)) && Equals(cpuSubfunction, (byte) (0x02)) && Equals(dataLength, (ushort) (0x02)))
            {
                return S7PayloadUserDataItemCpuFunctionMsgSubscriptionSysResponse.StaticParse(readBuffer, cpuFunctionGroup, cpuFunctionType, cpuSubfunction, returnCode, transportSize, dataLength);
            }
            if (Equals(cpuFunctionGroup, (byte) (0x04)) && Equals(cpuFunctionType, (byte) (0x08)) && Equals(cpuSubfunction, (byte) (0x02)) && Equals(dataLength, (ushort) (0x05)))
            {
                return S7PayloadUserDataItemCpuFunctionMsgSubscriptionAlarmResponse.StaticParse(readBuffer, cpuFunctionGroup, cpuFunctionType, cpuSubfunction, returnCode, transportSize, dataLength);
            }
            if (Equals(cpuFunctionGroup, (byte) (0x04)) && Equals(cpuFunctionType, (byte) (0x04)) && Equals(cpuSubfunction, (byte) (0x0b)))
            {
                return S7PayloadUserDataItemCpuFunctionAlarmAckRequest.StaticParse(readBuffer, cpuFunctionGroup, cpuFunctionType, cpuSubfunction, returnCode, transportSize, dataLength);
            }
            if (Equals(cpuFunctionGroup, (byte) (0x04)) && Equals(cpuFunctionType, (byte) (0x08)) && Equals(cpuSubfunction, (byte) (0x0b)) && Equals(dataLength, (ushort) (0x00)))
            {
                return S7PayloadUserDataItemCpuFunctionAlarmAckErrorResponse.StaticParse(readBuffer, cpuFunctionGroup, cpuFunctionType, cpuSubfunction, returnCode, transportSize, dataLength);
            }
            if (Equals(cpuFunctionGroup, (byte) (0x04)) && Equals(cpuFunctionType, (byte) (0x08)) && Equals(cpuSubfunction, (byte) (0x0b)))
            {
                return S7PayloadUserDataItemCpuFunctionAlarmAckResponse.StaticParse(readBuffer, cpuFunctionGroup, cpuFunctionType, cpuSubfunction, returnCode, transportSize, dataLength);
            }
            if (Equals(cpuFunctionGroup, (byte) (0x04)) && Equals(cpuFunctionType, (byte) (0x04)) && Equals(cpuSubfunction, (byte) (0x13)))
            {
                return S7PayloadUserDataItemCpuFunctionAlarmQueryRequest.StaticParse(readBuffer, cpuFunctionGroup, cpuFunctionType, cpuSubfunction, returnCode, transportSize, dataLength);
            }
            if (Equals(cpuFunctionGroup, (byte) (0x04)) && Equals(cpuFunctionType, (byte) (0x08)) && Equals(cpuSubfunction, (byte) (0x13)))
            {
                return S7PayloadUserDataItemCpuFunctionAlarmQueryResponse.StaticParse(readBuffer, cpuFunctionGroup, cpuFunctionType, cpuSubfunction, returnCode, transportSize, dataLength);
            }
            if (Equals(cpuFunctionGroup, (byte) (0x07)) && Equals(cpuFunctionType, (byte) (0x04)) && Equals(cpuSubfunction, (byte) (0x01)))
            {
                return S7PayloadUserDataItemClkRequest.StaticParse(readBuffer, cpuFunctionGroup, cpuFunctionType, cpuSubfunction, returnCode, transportSize, dataLength);
            }
            if (Equals(cpuFunctionGroup, (byte) (0x07)) && Equals(cpuFunctionType, (byte) (0x08)) && Equals(cpuSubfunction, (byte) (0x01)))
            {
                return S7PayloadUserDataItemClkResponse.StaticParse(readBuffer, cpuFunctionGroup, cpuFunctionType, cpuSubfunction, returnCode, transportSize, dataLength);
            }
            if (Equals(cpuFunctionGroup, (byte) (0x07)) && Equals(cpuFunctionType, (byte) (0x04)) && Equals(cpuSubfunction, (byte) (0x03)))
            {
                return S7PayloadUserDataItemClkFRequest.StaticParse(readBuffer, cpuFunctionGroup, cpuFunctionType, cpuSubfunction, returnCode, transportSize, dataLength);
            }
            if (Equals(cpuFunctionGroup, (byte) (0x07)) && Equals(cpuFunctionType, (byte) (0x08)) && Equals(cpuSubfunction, (byte) (0x03)))
            {
                return S7PayloadUserDataItemClkFResponse.StaticParse(readBuffer, cpuFunctionGroup, cpuFunctionType, cpuSubfunction, returnCode, transportSize, dataLength);
            }
            if (Equals(cpuFunctionGroup, (byte) (0x07)) && Equals(cpuFunctionType, (byte) (0x04)) && Equals(cpuSubfunction, (byte) (0x04)))
            {
                return S7PayloadUserDataItemClkSetRequest.StaticParse(readBuffer, cpuFunctionGroup, cpuFunctionType, cpuSubfunction, returnCode, transportSize, dataLength);
            }
            if (Equals(cpuFunctionGroup, (byte) (0x07)) && Equals(cpuFunctionType, (byte) (0x08)) && Equals(cpuSubfunction, (byte) (0x04)))
            {
                return S7PayloadUserDataItemClkSetResponse.StaticParse(readBuffer, cpuFunctionGroup, cpuFunctionType, cpuSubfunction, returnCode, transportSize, dataLength);
            }
            throw new ParseException("No matching subtype found for S7PayloadUserDataItem");
        }

        public void Serialize(WriteBuffer writeBuffer)
        {
            writeBuffer.WriteByte("returnCode", 8, (byte) ReturnCode);
            writeBuffer.WriteByte("transportSize", 8, (byte) TransportSize);
            writeBuffer.WriteUshort("dataLength", 16, DataLength);
            SerializeChild(writeBuffer);
        }

        protected abstract void SerializeChild(WriteBuffer writeBuffer);

        public int GetLengthInBits()
        {
            var lengthInBits = 0;
            lengthInBits += 8;
            lengthInBits += 8;
            lengthInBits += 16;
            lengthInBits += GetLengthInBitsChild();
            return lengthInBits;
        }

        protected abstract int GetLengthInBitsChild();

        public int GetLengthInBytes() => GetLengthInBits() / 8;
    }
}
