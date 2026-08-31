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
    public partial class SubItem : IMessage
    {
        public byte BytesToRead { get; }
        public ushort DbNumber { get; }
        public ushort StartAddress { get; }

        public SubItem(byte bytesToRead, ushort dbNumber, ushort startAddress)
        {
            BytesToRead = bytesToRead;
            DbNumber = dbNumber;
            StartAddress = startAddress;
        }

        public static SubItem StaticParse(ReadBuffer readBuffer)
        {
            var bytesToRead = readBuffer.ReadByte("bytesToRead", 8);
            var dbNumber = readBuffer.ReadUshort("dbNumber", 16);
            var startAddress = readBuffer.ReadUshort("startAddress", 16);
            return new SubItem(bytesToRead, dbNumber, startAddress);
        }

        public void Serialize(WriteBuffer writeBuffer)
        {
            writeBuffer.WriteByte("bytesToRead", 8, BytesToRead);
            writeBuffer.WriteUshort("dbNumber", 16, DbNumber);
            writeBuffer.WriteUshort("startAddress", 16, StartAddress);
        }

        public int GetLengthInBits()
        {
            var lengthInBits = 0;
            lengthInBits += 8;
            lengthInBits += 16;
            lengthInBits += 16;
            return lengthInBits;
        }

        public int GetLengthInBytes() => GetLengthInBits() / 8;
    }
}
