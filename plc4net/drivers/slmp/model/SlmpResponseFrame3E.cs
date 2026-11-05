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

namespace org.apache.plc4net.drivers.slmp.readwrite.model
{
    public partial class SlmpResponseFrame3E : SlmpMessage
    {
        public const byte Network = 0x00;
        public const byte PcNumber = 0xFF;
        public const ushort RequestDestModuleIoNo = 0x03FF;
        public const byte RequestDestModuleStation = 0x00;

        public override byte SubHeader => (byte) (0xD0);

        public ushort EndCode { get; }
        public byte[] ResponseData { get; }

        public SlmpResponseFrame3E(ushort endCode, byte[] responseData)
        {
            EndCode = endCode;
            ResponseData = responseData;
        }

        public static new SlmpResponseFrame3E StaticParse(ReadBuffer readBuffer)
        {
            var network = readBuffer.ReadByte("network", 8);
            if (!Equals(network, (byte) (0x00)))
                throw new ParseException($"Expected constant {Network} for 'network' but got {network}");
            var pcNumber = readBuffer.ReadByte("pcNumber", 8);
            if (!Equals(pcNumber, (byte) (0xFF)))
                throw new ParseException($"Expected constant {PcNumber} for 'pcNumber' but got {pcNumber}");
            var requestDestModuleIoNo = readBuffer.ReadUshort("requestDestModuleIoNo", 16);
            if (!Equals(requestDestModuleIoNo, (ushort) (0x03FF)))
                throw new ParseException($"Expected constant {RequestDestModuleIoNo} for 'requestDestModuleIoNo' but got {requestDestModuleIoNo}");
            var requestDestModuleStation = readBuffer.ReadByte("requestDestModuleStation", 8);
            if (!Equals(requestDestModuleStation, (byte) (0x00)))
                throw new ParseException($"Expected constant {RequestDestModuleStation} for 'requestDestModuleStation' but got {requestDestModuleStation}");
            var responseDataLength = readBuffer.ReadUshort("responseDataLength", 16);
            var endCode = readBuffer.ReadUshort("endCode", 16);
            var responseData = readBuffer.ReadByteArray("responseData", (int) ((responseDataLength - 2)) * 8);
            return new SlmpResponseFrame3E(endCode, responseData);
        }

        protected override void SerializeChild(WriteBuffer writeBuffer)
        {
            writeBuffer.WriteByte("network", 8, Network);
            writeBuffer.WriteByte("pcNumber", 8, PcNumber);
            writeBuffer.WriteUshort("requestDestModuleIoNo", 16, RequestDestModuleIoNo);
            writeBuffer.WriteByte("requestDestModuleStation", 8, RequestDestModuleStation);
            writeBuffer.WriteUshort("responseDataLength", 16, (ushort) ((2 + ResponseData.Length)));
            writeBuffer.WriteUshort("endCode", 16, EndCode);
            writeBuffer.WriteByteArray("responseData", ResponseData);
        }

        protected override int GetLengthInBitsChild()
        {
            var lengthInBits = 0;
            lengthInBits += 8;
            lengthInBits += 8;
            lengthInBits += 16;
            lengthInBits += 8;
            lengthInBits += 16;
            lengthInBits += 16;
            lengthInBits += (ResponseData.Length * 8);
            return lengthInBits;
        }

    }
}
