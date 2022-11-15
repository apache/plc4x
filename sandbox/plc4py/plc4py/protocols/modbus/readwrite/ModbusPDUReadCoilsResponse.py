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
from ctypes import c_byte
from ctypes import c_uint8
from plc4py.api.messages.PlcMessage import PlcMessage
import math

    
@dataclass
class ModbusPDUReadCoilsResponse(PlcMessage,ModbusPDU):
    value: []c_byte

    # Accessors for discriminator values.
    def getErrorFlag(self) -> c_bool:
        return (c_bool) False
    def getFunctionFlag(self) -> c_uint8:
        return (c_uint8) 0x01
    def getResponse(self) -> c_bool:
        return (c_bool) True


    def __post_init__(self):
        super().__init__( )



    def getValue(self) -> []c_byte:
        return self.value


    def serializeModbusPDUChild(self, writeBuffer: WriteBuffer):
        positionAware: PositionAware = writeBuffer
        startPos: int = positionAware.getPos()
        writeBuffer.pushContext("ModbusPDUReadCoilsResponse")

        # Implicit Field (byteCount) (Used for parsing, but its value is not stored as it's implicitly given by the objects content)
        c_uint8 byteCount = (c_uint8) (COUNT(getValue()))
        writeImplicitField("byteCount", byteCount, writeUnsignedShort(writeBuffer, 8))

        # Array Field (value)
        writeByteArrayField("value", value, writeByteArray(writeBuffer, 8))

        writeBuffer.popContext("ModbusPDUReadCoilsResponse")


    def getLengthInBytes(self) -> int:
        return int(math.ceil(float(self.getLengthInBits() / 8.0)))

    def getLengthInBits(self) -> int:
        lengthInBits: int = super().getLengthInBits()
        _value: ModbusPDUReadCoilsResponse = self

        # Implicit Field (byteCount)
        lengthInBits += 8

        # Array field
        if value is not None):
            lengthInBits += 8 * value.length


        return lengthInBits


    @staticmethod
    def staticParseBuilder(readBuffer: ReadBuffer, response: c_bool) -> ModbusPDUReadCoilsResponseBuilder:
        readBuffer.pullContext("ModbusPDUReadCoilsResponse")
        positionAware: PositionAware = readBuffer
        startPos: int = positionAware.getPos()
        curPos: int = 0

        byteCount: c_uint8 = readImplicitField("byteCount", readUnsignedShort(readBuffer, 8))

        value: byte[] = readBuffer.readByteArray("value", Math.toIntExact(byteCount))

        readBuffer.closeContext("ModbusPDUReadCoilsResponse")
        # Create the instance
        return ModbusPDUReadCoilsResponseBuilder(value )


    def equals(self, o: object) -> bool:
        if self == o:
            return True

        if not isinstance(o, ModbusPDUReadCoilsResponse):
            return False

        that: ModbusPDUReadCoilsResponse = ModbusPDUReadCoilsResponse(o)
        return (getValue() == that.getValue()) && super().equals(that) && True

    def hashCode(self) -> int:
        return hash(super().hashCode(), getValue() )

    def __str__(self) -> str:
        writeBufferBoxBased: WriteBufferBoxBased = WriteBufferBoxBased(True, True)
        try:
            writeBufferBoxBased.writeSerializable(self)
        except SerializationException as e:
            raise RuntimeException(e)

        return "\n" + str(writeBufferBoxBased.getBox()) + "\n"


class ModbusPDUReadCoilsResponseBuilder(ModbusPDUModbusPDUBuilder: value: []c_bytedef ModbusPDUReadCoilsResponseBuilder( []c_byte value ):        self.value = value


        def build(self,
        ) -> ModbusPDUReadCoilsResponse:
        modbusPDUReadCoilsResponse: ModbusPDUReadCoilsResponse = ModbusPDUReadCoilsResponse(
            value
)
        return modbusPDUReadCoilsResponse



