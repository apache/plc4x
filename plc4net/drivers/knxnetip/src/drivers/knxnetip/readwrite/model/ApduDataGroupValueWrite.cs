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
    public partial class ApduDataGroupValueWrite : ApduData
    {
        public override byte ApciType => (byte) (0x2);

        public sbyte DataFirstByte { get; }
        public byte[] Data { get; }

        public ApduDataGroupValueWrite(sbyte dataFirstByte, byte[] data)
        {
            DataFirstByte = dataFirstByte;
            Data = data;
        }

        public static new ApduDataGroupValueWrite StaticParse(ReadBuffer readBuffer, byte dataLength)
        {
            var dataFirstByte = readBuffer.ReadSbyte("dataFirstByte", 6);
            var data = readBuffer.ReadByteArray("data", (int) (((dataLength < 1) ? 0 : (dataLength - 1))) * 8);
            return new ApduDataGroupValueWrite(dataFirstByte, data);
        }

        protected override void SerializeChild(WriteBuffer writeBuffer)
        {
            writeBuffer.WriteSbyte("dataFirstByte", 6, DataFirstByte);
            writeBuffer.WriteByteArray("data", Data);
        }

        protected override int GetLengthInBitsChild()
        {
            var lengthInBits = 0;
            lengthInBits += 6;
            lengthInBits += (Data.Length * 8);
            return lengthInBits;
        }

    }
}
