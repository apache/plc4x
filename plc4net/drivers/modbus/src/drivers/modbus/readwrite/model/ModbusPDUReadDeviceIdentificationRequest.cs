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
    public partial class ModbusPDUReadDeviceIdentificationRequest : ModbusPDU
    {
        public const byte MeiType = 0x0E;

        public override bool ErrorFlag => false;
        public override byte FunctionFlag => (byte) (0x2B);
        public override bool Response => false;

        public ModbusDeviceInformationLevel Level { get; }
        public byte ObjectId { get; }

        public ModbusPDUReadDeviceIdentificationRequest(ModbusDeviceInformationLevel level, byte objectId)
        {
            Level = level;
            ObjectId = objectId;
        }

        public static new ModbusPDUReadDeviceIdentificationRequest StaticParse(ReadBuffer readBuffer, bool response)
        {
            var meiType = readBuffer.ReadByte("meiType", 8);
            if (!Equals(meiType, (byte) (0x0E)))
                throw new ParseException($"Expected constant {MeiType} for 'meiType' but got {meiType}");
            var level = (ModbusDeviceInformationLevel) readBuffer.ReadByte("level", 8);
            var objectId = readBuffer.ReadByte("objectId", 8);
            return new ModbusPDUReadDeviceIdentificationRequest(level, objectId);
        }

        protected override void SerializeChild(WriteBuffer writeBuffer)
        {
            writeBuffer.WriteByte("meiType", 8, MeiType);
            writeBuffer.WriteByte("level", 8, (byte) Level);
            writeBuffer.WriteByte("objectId", 8, ObjectId);
        }

        protected override int GetLengthInBitsChild()
        {
            var lengthInBits = 0;
            lengthInBits += 8;
            lengthInBits += 8;
            lengthInBits += 8;
            return lengthInBits;
        }

    }
}
