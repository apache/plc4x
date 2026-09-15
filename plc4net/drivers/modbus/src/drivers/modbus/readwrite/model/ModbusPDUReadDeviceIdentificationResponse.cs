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
    public partial class ModbusPDUReadDeviceIdentificationResponse : ModbusPDU
    {
        public const byte MeiType = 0x0E;

        public override bool ErrorFlag => false;
        public override byte FunctionFlag => (byte) (0x2B);
        public override bool Response => true;

        public ModbusDeviceInformationLevel Level { get; }
        public bool IndividualAccess { get; }
        public ModbusDeviceInformationConformityLevel ConformityLevel { get; }
        public ModbusDeviceInformationMoreFollows MoreFollows { get; }
        public byte NextObjectId { get; }
        public System.Collections.Generic.List<ModbusDeviceInformationObject> Objects { get; }

        public ModbusPDUReadDeviceIdentificationResponse(ModbusDeviceInformationLevel level, bool individualAccess, ModbusDeviceInformationConformityLevel conformityLevel, ModbusDeviceInformationMoreFollows moreFollows, byte nextObjectId, System.Collections.Generic.List<ModbusDeviceInformationObject> objects)
        {
            Level = level;
            IndividualAccess = individualAccess;
            ConformityLevel = conformityLevel;
            MoreFollows = moreFollows;
            NextObjectId = nextObjectId;
            Objects = objects;
        }

        public static new ModbusPDUReadDeviceIdentificationResponse StaticParse(ReadBuffer readBuffer, bool response)
        {
            var meiType = readBuffer.ReadByte("meiType", 8);
            if (!Equals(meiType, (byte) (0x0E)))
                throw new ParseException($"Expected constant {MeiType} for 'meiType' but got {meiType}");
            var level = (ModbusDeviceInformationLevel) readBuffer.ReadByte("level", 8);
            var individualAccess = readBuffer.ReadBit("individualAccess");
            var conformityLevel = (ModbusDeviceInformationConformityLevel) readBuffer.ReadByte("conformityLevel", 7);
            var moreFollows = (ModbusDeviceInformationMoreFollows) readBuffer.ReadByte("moreFollows", 8);
            var nextObjectId = readBuffer.ReadByte("nextObjectId", 8);
            var numberOfObjects = readBuffer.ReadByte("numberOfObjects", 8);
            var objects = new System.Collections.Generic.List<ModbusDeviceInformationObject>();
            var _objectsCnt = (int) (numberOfObjects);
            for (var _objectsI = 0; _objectsI < _objectsCnt; _objectsI++)
            {
                objects.Add(ModbusDeviceInformationObject.StaticParse(readBuffer));
            }
            return new ModbusPDUReadDeviceIdentificationResponse(level, individualAccess, conformityLevel, moreFollows, nextObjectId, objects);
        }

        protected override void SerializeChild(WriteBuffer writeBuffer)
        {
            writeBuffer.WriteByte("meiType", 8, MeiType);
            writeBuffer.WriteByte("level", 8, (byte) Level);
            writeBuffer.WriteBit("individualAccess", IndividualAccess);
            writeBuffer.WriteByte("conformityLevel", 7, (byte) ConformityLevel);
            writeBuffer.WriteByte("moreFollows", 8, (byte) MoreFollows);
            writeBuffer.WriteByte("nextObjectId", 8, NextObjectId);
            writeBuffer.WriteByte("numberOfObjects", 8, (byte) (Objects.Count));
            foreach (var _e in Objects)
            {
                _e.Serialize(writeBuffer);
            }
        }

        protected override int GetLengthInBitsChild()
        {
            var lengthInBits = 0;
            lengthInBits += 8;
            lengthInBits += 8;
            lengthInBits += 1;
            lengthInBits += 7;
            lengthInBits += 8;
            lengthInBits += 8;
            lengthInBits += 8;
            lengthInBits += Objects.Sum(_e => _e.GetLengthInBits());
            return lengthInBits;
        }

    }
}
