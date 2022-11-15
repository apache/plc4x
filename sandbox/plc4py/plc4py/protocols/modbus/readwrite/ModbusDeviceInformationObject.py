#
# Licensed to the Apache Software Foundation (ASF) under one
# or more contributor license agreements.  See the NOTICE file
# distributed with this work for additional information
# regarding copyright ownership.  The ASF licenses this file
# to you under the Apache License, Version 2.0 (the
# "License") you may not use this file except in compliance
# with the License.  You may obtain a copy of the License at
#
#   https:#www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing,
# software distributed under the License is distributed on an
# "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
# KIND, either express or implied.  See the License for the
# specific language governing permissions and limitations
# under the License.
#

# Code generated by code-generation. DO NOT EDIT.
from abc import ABC, abstractmethod
from dataclasses import dataclass


from ctypes import c_byte
from ctypes import c_uint8
from plc4py.api.messages.PlcMessage import PlcMessage
import math

    
@dataclass
class ModbusDeviceInformationObject(PlcMessage):
    objectId: c_uint8
    data: []c_byte



    def __post_init__(self):
        super().__init__( )



    def getObjectId(self) -> c_uint8:
        return self.objectId

    def getData(self) -> []c_byte:
        return self.data


    def serialize(self, writeBuffer: WriteBuffer):
        positionAware: PositionAware = writeBuffer
        startPos: int = positionAware.getPos()
        writeBuffer.pushContext("ModbusDeviceInformationObject")

        # Simple Field (objectId)
        writeSimpleField("objectId", objectId, writeUnsignedShort(writeBuffer, 8))

        # Implicit Field (objectLength) (Used for parsing, but its value is not stored as it's implicitly given by the objects content)
        c_uint8 objectLength = (c_uint8) (COUNT(getData()))
        writeImplicitField("objectLength", objectLength, writeUnsignedShort(writeBuffer, 8))

        # Array Field (data)
        writeByteArrayField("data", data, writeByteArray(writeBuffer, 8))

        writeBuffer.popContext("ModbusDeviceInformationObject")


    def getLengthInBytes(self) -> int:
        return int(math.ceil(float(self.getLengthInBits() / 8.0)))

    def getLengthInBits(self) -> int:
        lengthInBits: int = 0
        _value: ModbusDeviceInformationObject = self

        # Simple field (objectId)
        lengthInBits += 8

        # Implicit Field (objectLength)
        lengthInBits += 8

        # Array field
        if data is not None):
            lengthInBits += 8 * data.length


        return lengthInBits


    def staticParse(readBuffer: ReadBuffer , args) -> ModbusDeviceInformationObject:
        positionAware: PositionAware = readBuffer
        return staticParse(readBuffer)


    @staticmethod
    def staticParseContext(readBuffer: ReadBuffer) -> ModbusDeviceInformationObject:
        readBuffer.pullContext("ModbusDeviceInformationObject")
        positionAware: PositionAware = readBuffer
        startPos: int = positionAware.getPos()
        curPos: int = 0

        objectId: c_uint8 = readSimpleField("objectId", readUnsignedShort(readBuffer, 8))

        objectLength: c_uint8 = readImplicitField("objectLength", readUnsignedShort(readBuffer, 8))

        data: byte[] = readBuffer.readByteArray("data", Math.toIntExact(objectLength))

        readBuffer.closeContext("ModbusDeviceInformationObject")
        # Create the instance
        _modbusDeviceInformationObject: ModbusDeviceInformationObject = ModbusDeviceInformationObject(objectId, data )
        return _modbusDeviceInformationObject


    def equals(self, o: object) -> bool:
        if self == o:
            return True

        if not isinstance(o, ModbusDeviceInformationObject):
            return False

        that: ModbusDeviceInformationObject = ModbusDeviceInformationObject(o)
        return (getObjectId() == that.getObjectId()) && (getData() == that.getData()) && True

    def hashCode(self) -> int:
        return hash(getObjectId(), getData() )

    def __str__(self) -> str:
        writeBufferBoxBased: WriteBufferBoxBased = WriteBufferBoxBased(True, True)
        try:
            writeBufferBoxBased.writeSerializable(self)
        except SerializationException as e:
            raise RuntimeException(e)

        return "\n" + str(writeBufferBoxBased.getBox()) + "\n"




