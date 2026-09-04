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
    public partial class ApduDataExtAuthorizeRequest : ApduDataExt
    {
        public override byte ExtApciType => (byte) (0x11);

        public byte Level { get; }
        public byte[] Data { get; }

        public ApduDataExtAuthorizeRequest(byte level, byte[] data)
        {
            Level = level;
            Data = data;
        }

        public static new ApduDataExtAuthorizeRequest StaticParse(ReadBuffer readBuffer, byte length)
        {
            var level = readBuffer.ReadByte("level", 8);
            var data = readBuffer.ReadByteArray("data", (int) (4) * 8);
            return new ApduDataExtAuthorizeRequest(level, data);
        }

        protected override void SerializeChild(WriteBuffer writeBuffer)
        {
            writeBuffer.WriteByte("level", 8, Level);
            writeBuffer.WriteByteArray("data", Data);
        }

        protected override int GetLengthInBitsChild()
        {
            var lengthInBits = 0;
            lengthInBits += 8;
            lengthInBits += (Data.Length * 8);
            return lengthInBits;
        }

    }
}
