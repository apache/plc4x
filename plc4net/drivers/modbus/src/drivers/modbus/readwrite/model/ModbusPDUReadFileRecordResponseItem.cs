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
    public partial class ModbusPDUReadFileRecordResponseItem : IMessage
    {
        public byte ReferenceType { get; }
        public byte[] Data { get; }

        public ModbusPDUReadFileRecordResponseItem(byte referenceType, byte[] data)
        {
            ReferenceType = referenceType;
            Data = data;
        }

        public static ModbusPDUReadFileRecordResponseItem StaticParse(ReadBuffer readBuffer)
        {
            var dataLength = readBuffer.ReadByte("dataLength", 8);
            var referenceType = readBuffer.ReadByte("referenceType", 8);
            var data = readBuffer.ReadByteArray("data", (int) ((dataLength - 1)) * 8);
            return new ModbusPDUReadFileRecordResponseItem(referenceType, data);
        }

        public void Serialize(WriteBuffer writeBuffer)
        {
            writeBuffer.WriteByte("dataLength", 8, (byte) ((Data.Length + 1)));
            writeBuffer.WriteByte("referenceType", 8, ReferenceType);
            writeBuffer.WriteByteArray("data", Data);
        }

        public int GetLengthInBits()
        {
            var lengthInBits = 0;
            lengthInBits += 8;
            lengthInBits += 8;
            lengthInBits += (Data.Length * 8);
            return lengthInBits;
        }

        public int GetLengthInBytes() => GetLengthInBits() / 8;
    }
}
