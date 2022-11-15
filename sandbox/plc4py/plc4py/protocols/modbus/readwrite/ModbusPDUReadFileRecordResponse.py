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


from ctypes import c_bool
from ctypes import c_uint8
from plc4py.api.messages.PlcMessage import PlcMessage
from plc4py.protocols.modbus.readwrite.ModbusPDUReadFileRecordResponseItem import ModbusPDUReadFileRecordResponseItem
import math

    
@dataclass
class ModbusPDUReadFileRecordResponse(PlcMessage,ModbusPDU):
    items: []ModbusPDUReadFileRecordResponseItem

    # Accessors for discriminator values.
    def getErrorFlag(self) -> c_bool:
        return (c_bool) False
    def getFunctionFlag(self) -> c_uint8:
        return (c_uint8) 0x14
    def getResponse(self) -> c_bool:
        return (c_bool) True


    def __post_init__(self):
        super().__init__( )



    def getItems(self) -> []ModbusPDUReadFileRecordResponseItem:
        return self.items


    def serializeModbusPDUChild(self, writeBuffer: WriteBuffer):
        positionAware: PositionAware = writeBuffer
        startPos: int = positionAware.getPos()
        writeBuffer.pushContext("ModbusPDUReadFileRecordResponse")

        # Implicit Field (byteCount) (Used for parsing, but its value is not stored as it's implicitly given by the objects content)
        c_uint8 byteCount = (c_uint8) (ARRAY_SIZE_IN_BYTES(getItems()))
        writeImplicitField("byteCount", byteCount, writeUnsignedShort(writeBuffer, 8))

        # Array Field (items)
        writeComplexTypeArrayField("items", items, writeBuffer)

        writeBuffer.popContext("ModbusPDUReadFileRecordResponse")


    def getLengthInBytes(self) -> int:
        return int(math.ceil(float(self.getLengthInBits() / 8.0)))

    def getLengthInBits(self) -> int:
        lengthInBits: int = super().getLengthInBits()
        _value: ModbusPDUReadFileRecordResponse = self

        # Implicit Field (byteCount)
        lengthInBits += 8

        # Array field
        if items is not None):
            for element in items:
                lengthInBits += element.getLengthInBits()



        return lengthInBits


    @staticmethod
    def staticParseBuilder(readBuffer: ReadBuffer, response: c_bool) -> ModbusPDUReadFileRecordResponseBuilder:
        readBuffer.pullContext("ModbusPDUReadFileRecordResponse")
        positionAware: PositionAware = readBuffer
        startPos: int = positionAware.getPos()
        curPos: int = 0

        byteCount: c_uint8 = readImplicitField("byteCount", readUnsignedShort(readBuffer, 8))

        items: []ModbusPDUReadFileRecordResponseItem = readLengthArrayField("items", DataReaderComplexDefault<>(() -> ModbusPDUReadFileRecordResponseItem.staticParse(readBuffer), readBuffer), byteCount)

        readBuffer.closeContext("ModbusPDUReadFileRecordResponse")
        # Create the instance
        return ModbusPDUReadFileRecordResponseBuilder(items )


    def equals(self, o: object) -> bool:
        if self == o:
            return True

        if not isinstance(o, ModbusPDUReadFileRecordResponse):
            return False

        that: ModbusPDUReadFileRecordResponse = ModbusPDUReadFileRecordResponse(o)
        return (getItems() == that.getItems()) && super().equals(that) && True

    def hashCode(self) -> int:
        return hash(super().hashCode(), getItems() )

    def __str__(self) -> str:
        writeBufferBoxBased: WriteBufferBoxBased = WriteBufferBoxBased(True, True)
        try:
            writeBufferBoxBased.writeSerializable(self)
        except SerializationException as e:
            raise RuntimeException(e)

        return "\n" + str(writeBufferBoxBased.getBox()) + "\n"


class ModbusPDUReadFileRecordResponseBuilder(ModbusPDUModbusPDUBuilder: items: []ModbusPDUReadFileRecordResponseItemdef ModbusPDUReadFileRecordResponseBuilder( []ModbusPDUReadFileRecordResponseItem items ):        self.items = items


        def build(self,
        ) -> ModbusPDUReadFileRecordResponse:
        modbusPDUReadFileRecordResponse: ModbusPDUReadFileRecordResponse = ModbusPDUReadFileRecordResponse(
            items
)
        return modbusPDUReadFileRecordResponse



