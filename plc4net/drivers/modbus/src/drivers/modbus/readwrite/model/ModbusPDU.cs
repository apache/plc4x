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

namespace org.apache.plc4net.drivers.modbus.readwrite.model
{
    public abstract partial class ModbusPDU : IMessage
    {
        public abstract bool ErrorFlag { get; }
        public abstract byte FunctionFlag { get; }
        public abstract bool Response { get; }

        public static ModbusPDU StaticParse(ReadBuffer readBuffer, bool response)
        {
            var errorFlag = readBuffer.ReadBit("errorFlag");
            var functionFlag = readBuffer.ReadByte("functionFlag", 7);
            if (Equals(errorFlag, true))
            {
                return ModbusPDUError.StaticParse(readBuffer, response);
            }
            if (Equals(errorFlag, false) && Equals(functionFlag, (byte) (0x02)) && Equals(response, false))
            {
                return ModbusPDUReadDiscreteInputsRequest.StaticParse(readBuffer, response);
            }
            if (Equals(errorFlag, false) && Equals(functionFlag, (byte) (0x02)) && Equals(response, true))
            {
                return ModbusPDUReadDiscreteInputsResponse.StaticParse(readBuffer, response);
            }
            if (Equals(errorFlag, false) && Equals(functionFlag, (byte) (0x01)) && Equals(response, false))
            {
                return ModbusPDUReadCoilsRequest.StaticParse(readBuffer, response);
            }
            if (Equals(errorFlag, false) && Equals(functionFlag, (byte) (0x01)) && Equals(response, true))
            {
                return ModbusPDUReadCoilsResponse.StaticParse(readBuffer, response);
            }
            if (Equals(errorFlag, false) && Equals(functionFlag, (byte) (0x05)) && Equals(response, false))
            {
                return ModbusPDUWriteSingleCoilRequest.StaticParse(readBuffer, response);
            }
            if (Equals(errorFlag, false) && Equals(functionFlag, (byte) (0x05)) && Equals(response, true))
            {
                return ModbusPDUWriteSingleCoilResponse.StaticParse(readBuffer, response);
            }
            if (Equals(errorFlag, false) && Equals(functionFlag, (byte) (0x0F)) && Equals(response, false))
            {
                return ModbusPDUWriteMultipleCoilsRequest.StaticParse(readBuffer, response);
            }
            if (Equals(errorFlag, false) && Equals(functionFlag, (byte) (0x0F)) && Equals(response, true))
            {
                return ModbusPDUWriteMultipleCoilsResponse.StaticParse(readBuffer, response);
            }
            if (Equals(errorFlag, false) && Equals(functionFlag, (byte) (0x04)) && Equals(response, false))
            {
                return ModbusPDUReadInputRegistersRequest.StaticParse(readBuffer, response);
            }
            if (Equals(errorFlag, false) && Equals(functionFlag, (byte) (0x04)) && Equals(response, true))
            {
                return ModbusPDUReadInputRegistersResponse.StaticParse(readBuffer, response);
            }
            if (Equals(errorFlag, false) && Equals(functionFlag, (byte) (0x03)) && Equals(response, false))
            {
                return ModbusPDUReadHoldingRegistersRequest.StaticParse(readBuffer, response);
            }
            if (Equals(errorFlag, false) && Equals(functionFlag, (byte) (0x03)) && Equals(response, true))
            {
                return ModbusPDUReadHoldingRegistersResponse.StaticParse(readBuffer, response);
            }
            if (Equals(errorFlag, false) && Equals(functionFlag, (byte) (0x06)) && Equals(response, false))
            {
                return ModbusPDUWriteSingleRegisterRequest.StaticParse(readBuffer, response);
            }
            if (Equals(errorFlag, false) && Equals(functionFlag, (byte) (0x06)) && Equals(response, true))
            {
                return ModbusPDUWriteSingleRegisterResponse.StaticParse(readBuffer, response);
            }
            if (Equals(errorFlag, false) && Equals(functionFlag, (byte) (0x10)) && Equals(response, false))
            {
                return ModbusPDUWriteMultipleHoldingRegistersRequest.StaticParse(readBuffer, response);
            }
            if (Equals(errorFlag, false) && Equals(functionFlag, (byte) (0x10)) && Equals(response, true))
            {
                return ModbusPDUWriteMultipleHoldingRegistersResponse.StaticParse(readBuffer, response);
            }
            if (Equals(errorFlag, false) && Equals(functionFlag, (byte) (0x17)) && Equals(response, false))
            {
                return ModbusPDUReadWriteMultipleHoldingRegistersRequest.StaticParse(readBuffer, response);
            }
            if (Equals(errorFlag, false) && Equals(functionFlag, (byte) (0x17)) && Equals(response, true))
            {
                return ModbusPDUReadWriteMultipleHoldingRegistersResponse.StaticParse(readBuffer, response);
            }
            if (Equals(errorFlag, false) && Equals(functionFlag, (byte) (0x16)) && Equals(response, false))
            {
                return ModbusPDUMaskWriteHoldingRegisterRequest.StaticParse(readBuffer, response);
            }
            if (Equals(errorFlag, false) && Equals(functionFlag, (byte) (0x16)) && Equals(response, true))
            {
                return ModbusPDUMaskWriteHoldingRegisterResponse.StaticParse(readBuffer, response);
            }
            if (Equals(errorFlag, false) && Equals(functionFlag, (byte) (0x18)) && Equals(response, false))
            {
                return ModbusPDUReadFifoQueueRequest.StaticParse(readBuffer, response);
            }
            if (Equals(errorFlag, false) && Equals(functionFlag, (byte) (0x18)) && Equals(response, true))
            {
                return ModbusPDUReadFifoQueueResponse.StaticParse(readBuffer, response);
            }
            if (Equals(errorFlag, false) && Equals(functionFlag, (byte) (0x14)) && Equals(response, false))
            {
                return ModbusPDUReadFileRecordRequest.StaticParse(readBuffer, response);
            }
            if (Equals(errorFlag, false) && Equals(functionFlag, (byte) (0x14)) && Equals(response, true))
            {
                return ModbusPDUReadFileRecordResponse.StaticParse(readBuffer, response);
            }
            if (Equals(errorFlag, false) && Equals(functionFlag, (byte) (0x15)) && Equals(response, false))
            {
                return ModbusPDUWriteFileRecordRequest.StaticParse(readBuffer, response);
            }
            if (Equals(errorFlag, false) && Equals(functionFlag, (byte) (0x15)) && Equals(response, true))
            {
                return ModbusPDUWriteFileRecordResponse.StaticParse(readBuffer, response);
            }
            if (Equals(errorFlag, false) && Equals(functionFlag, (byte) (0x07)) && Equals(response, false))
            {
                return ModbusPDUReadExceptionStatusRequest.StaticParse(readBuffer, response);
            }
            if (Equals(errorFlag, false) && Equals(functionFlag, (byte) (0x07)) && Equals(response, true))
            {
                return ModbusPDUReadExceptionStatusResponse.StaticParse(readBuffer, response);
            }
            if (Equals(errorFlag, false) && Equals(functionFlag, (byte) (0x08)) && Equals(response, false))
            {
                return ModbusPDUDiagnosticRequest.StaticParse(readBuffer, response);
            }
            if (Equals(errorFlag, false) && Equals(functionFlag, (byte) (0x08)) && Equals(response, true))
            {
                return ModbusPDUDiagnosticResponse.StaticParse(readBuffer, response);
            }
            if (Equals(errorFlag, false) && Equals(functionFlag, (byte) (0x0B)) && Equals(response, false))
            {
                return ModbusPDUGetComEventCounterRequest.StaticParse(readBuffer, response);
            }
            if (Equals(errorFlag, false) && Equals(functionFlag, (byte) (0x0B)) && Equals(response, true))
            {
                return ModbusPDUGetComEventCounterResponse.StaticParse(readBuffer, response);
            }
            if (Equals(errorFlag, false) && Equals(functionFlag, (byte) (0x0C)) && Equals(response, false))
            {
                return ModbusPDUGetComEventLogRequest.StaticParse(readBuffer, response);
            }
            if (Equals(errorFlag, false) && Equals(functionFlag, (byte) (0x0C)) && Equals(response, true))
            {
                return ModbusPDUGetComEventLogResponse.StaticParse(readBuffer, response);
            }
            if (Equals(errorFlag, false) && Equals(functionFlag, (byte) (0x11)) && Equals(response, false))
            {
                return ModbusPDUReportServerIdRequest.StaticParse(readBuffer, response);
            }
            if (Equals(errorFlag, false) && Equals(functionFlag, (byte) (0x11)) && Equals(response, true))
            {
                return ModbusPDUReportServerIdResponse.StaticParse(readBuffer, response);
            }
            if (Equals(errorFlag, false) && Equals(functionFlag, (byte) (0x2B)) && Equals(response, false))
            {
                return ModbusPDUReadDeviceIdentificationRequest.StaticParse(readBuffer, response);
            }
            if (Equals(errorFlag, false) && Equals(functionFlag, (byte) (0x2B)) && Equals(response, true))
            {
                return ModbusPDUReadDeviceIdentificationResponse.StaticParse(readBuffer, response);
            }
            throw new ParseException("No matching subtype found for ModbusPDU");
        }

        public void Serialize(WriteBuffer writeBuffer)
        {
            writeBuffer.WriteBit("errorFlag", ErrorFlag);
            writeBuffer.WriteByte("functionFlag", 7, FunctionFlag);
            SerializeChild(writeBuffer);
        }

        protected abstract void SerializeChild(WriteBuffer writeBuffer);

        public int GetLengthInBits()
        {
            var lengthInBits = 0;
            lengthInBits += 1;
            lengthInBits += 7;
            lengthInBits += GetLengthInBitsChild();
            return lengthInBits;
        }

        protected abstract int GetLengthInBitsChild();

        public int GetLengthInBytes() => GetLengthInBits() / 8;
    }
}
