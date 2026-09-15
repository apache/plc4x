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
    public partial class KnxGroupAddress2Level : KnxGroupAddress
    {
        public override byte NumLevels => (byte) (2);

        public byte MainGroup { get; }
        public ushort SubGroup { get; }

        public KnxGroupAddress2Level(byte mainGroup, ushort subGroup)
        {
            MainGroup = mainGroup;
            SubGroup = subGroup;
        }

        public static new KnxGroupAddress2Level StaticParse(ReadBuffer readBuffer, byte numLevels)
        {
            var mainGroup = readBuffer.ReadByte("mainGroup", 5);
            var subGroup = readBuffer.ReadUshort("subGroup", 11);
            return new KnxGroupAddress2Level(mainGroup, subGroup);
        }

        protected override void SerializeChild(WriteBuffer writeBuffer)
        {
            writeBuffer.WriteByte("mainGroup", 5, MainGroup);
            writeBuffer.WriteUshort("subGroup", 11, SubGroup);
        }

        protected override int GetLengthInBitsChild()
        {
            var lengthInBits = 0;
            lengthInBits += 5;
            lengthInBits += 11;
            return lengthInBits;
        }

    }
}
