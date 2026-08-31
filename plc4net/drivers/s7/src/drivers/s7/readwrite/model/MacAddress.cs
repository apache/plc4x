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
    public partial class MacAddress : IMessage
    {
        public byte[] Address { get; }

        public MacAddress(byte[] address)
        {
            Address = address;
        }

        public static MacAddress StaticParse(ReadBuffer readBuffer)
        {
            var address = readBuffer.ReadByteArray("address", (int) (6) * 8);
            return new MacAddress(address);
        }

        public void Serialize(WriteBuffer writeBuffer)
        {
            writeBuffer.WriteByteArray("address", Address);
        }

        public int GetLengthInBits()
        {
            var lengthInBits = 0;
            lengthInBits += (Address.Length * 8);
            return lengthInBits;
        }

        public int GetLengthInBytes() => GetLengthInBits() / 8;
    }
}
