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
    public partial class ModbusPDUWriteFileRecordRequestItem : IMessage
    {
        public byte ReferenceType { get; }
        public ushort FileNumber { get; }
        public ushort RecordNumber { get; }
        public byte[] RecordData { get; }

        public ModbusPDUWriteFileRecordRequestItem(byte referenceType, ushort fileNumber, ushort recordNumber, byte[] recordData)
        {
            ReferenceType = referenceType;
            FileNumber = fileNumber;
            RecordNumber = recordNumber;
            RecordData = recordData;
        }

        public static ModbusPDUWriteFileRecordRequestItem StaticParse(ReadBuffer readBuffer)
        {
            var referenceType = readBuffer.ReadByte("referenceType", 8);
            var fileNumber = readBuffer.ReadUshort("fileNumber", 16);
            var recordNumber = readBuffer.ReadUshort("recordNumber", 16);
            var recordLength = readBuffer.ReadUshort("recordLength", 16);
            var recordData = readBuffer.ReadByteArray("recordData", (int) ((recordLength * 2)) * 8);
            return new ModbusPDUWriteFileRecordRequestItem(referenceType, fileNumber, recordNumber, recordData);
        }

        public void Serialize(WriteBuffer writeBuffer)
        {
            writeBuffer.WriteByte("referenceType", 8, ReferenceType);
            writeBuffer.WriteUshort("fileNumber", 16, FileNumber);
            writeBuffer.WriteUshort("recordNumber", 16, RecordNumber);
            writeBuffer.WriteUshort("recordLength", 16, (ushort) ((RecordData.Length / 2)));
            writeBuffer.WriteByteArray("recordData", RecordData);
        }

        public int GetLengthInBits()
        {
            var lengthInBits = 0;
            lengthInBits += 8;
            lengthInBits += 16;
            lengthInBits += 16;
            lengthInBits += 16;
            lengthInBits += (RecordData.Length * 8);
            return lengthInBits;
        }

        public int GetLengthInBytes() => GetLengthInBits() / 8;
    }
}
