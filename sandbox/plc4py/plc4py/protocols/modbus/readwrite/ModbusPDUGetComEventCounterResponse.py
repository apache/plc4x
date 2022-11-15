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
from ctypes import c_uint16
from ctypes import c_uint8
from plc4py.api.messages.PlcMessage import PlcMessage
import math

    
@dataclass
class ModbusPDUGetComEventCounterResponse(PlcMessage,ModbusPDU):
    status: c_uint16
    eventCount: c_uint16

    # Accessors for discriminator values.
    def getErrorFlag(self) -> c_bool:
        return (c_bool) False
    def getFunctionFlag(self) -> c_uint8:
        return (c_uint8) 0x0B
    def getResponse(self) -> c_bool:
        return (c_bool) True


    def __post_init__(self):
        super().__init__( )



    def getStatus(self) -> c_uint16:
        return self.status

    def getEventCount(self) -> c_uint16:
        return self.eventCount


    def serializeModbusPDUChild(self, writeBuffer: WriteBuffer):
        positionAware: PositionAware = writeBuffer
        startPos: int = positionAware.getPos()
        writeBuffer.pushContext("ModbusPDUGetComEventCounterResponse")

        # Simple Field (status)
        writeSimpleField("status", status, writeUnsignedInt(writeBuffer, 16))

        # Simple Field (eventCount)
        writeSimpleField("eventCount", eventCount, writeUnsignedInt(writeBuffer, 16))

        writeBuffer.popContext("ModbusPDUGetComEventCounterResponse")


    def getLengthInBytes(self) -> int:
        return int(math.ceil(float(self.getLengthInBits() / 8.0)))

    def getLengthInBits(self) -> int:
        lengthInBits: int = super().getLengthInBits()
        _value: ModbusPDUGetComEventCounterResponse = self

        # Simple field (status)
        lengthInBits += 16

        # Simple field (eventCount)
        lengthInBits += 16

        return lengthInBits


    @staticmethod
    def staticParseBuilder(readBuffer: ReadBuffer, response: c_bool) -> ModbusPDUGetComEventCounterResponseBuilder:
        readBuffer.pullContext("ModbusPDUGetComEventCounterResponse")
        positionAware: PositionAware = readBuffer
        startPos: int = positionAware.getPos()
        curPos: int = 0

        status: c_uint16 = readSimpleField("status", readUnsignedInt(readBuffer, 16))

        eventCount: c_uint16 = readSimpleField("eventCount", readUnsignedInt(readBuffer, 16))

        readBuffer.closeContext("ModbusPDUGetComEventCounterResponse")
        # Create the instance
        return ModbusPDUGetComEventCounterResponseBuilder(status, eventCount )


    def equals(self, o: object) -> bool:
        if self == o:
            return True

        if not isinstance(o, ModbusPDUGetComEventCounterResponse):
            return False

        that: ModbusPDUGetComEventCounterResponse = ModbusPDUGetComEventCounterResponse(o)
        return (getStatus() == that.getStatus()) && (getEventCount() == that.getEventCount()) && super().equals(that) && True

    def hashCode(self) -> int:
        return hash(super().hashCode(), getStatus(), getEventCount() )

    def __str__(self) -> str:
        writeBufferBoxBased: WriteBufferBoxBased = WriteBufferBoxBased(True, True)
        try:
            writeBufferBoxBased.writeSerializable(self)
        except SerializationException as e:
            raise RuntimeException(e)

        return "\n" + str(writeBufferBoxBased.getBox()) + "\n"


class ModbusPDUGetComEventCounterResponseBuilder(ModbusPDUModbusPDUBuilder: status: c_uint16 eventCount: c_uint16def ModbusPDUGetComEventCounterResponseBuilder( c_uint16 status, c_uint16 eventCount ):        self.status = status
        self.eventCount = eventCount


        def build(self,
        ) -> ModbusPDUGetComEventCounterResponse:
        modbusPDUGetComEventCounterResponse: ModbusPDUGetComEventCounterResponse = ModbusPDUGetComEventCounterResponse(
            status, 
            eventCount
)
        return modbusPDUGetComEventCounterResponse



