#
# Licensed to the Apache Software Foundation (ASF) under one
# or more contributor license agreements.  See the NOTICE file
# distributed with this work for additional information
# regarding copyright ownership.  The ASF licenses this file
# to you under the Apache License, Version 2.0 (the
# "License"); you may not use this file except in compliance
# with the License.  You may obtain a copy of the License at
#
#     https://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing,
# software distributed under the License is distributed on an
# "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
# KIND, either express or implied.  See the License for the
# specific language governing permissions and limitations
# under the License.
#

from dataclasses import dataclass

from plc4py.api.exceptions.exceptions import PlcRuntimeException
from plc4py.api.exceptions.exceptions import SerializationException
from plc4py.api.messages.PlcMessage import PlcMessage
from plc4py.protocols.modbus.readwrite.ModbusPDU import ModbusPDU
from plc4py.spi.generation.ReadBuffer import ReadBuffer
from plc4py.spi.generation.WriteBuffer import WriteBuffer
from typing import ClassVar
import math


@dataclass
class ModbusPDUReadExceptionStatusRequest(ModbusPDU):
    # Accessors for discriminator values.
    error_flag: ClassVar[bool] = False
    function_flag: ClassVar[int] = 0x07
    response: ClassVar[bool] = False

    def serialize_modbus_pdu_child(self, write_buffer: WriteBuffer):
        write_buffer.push_context("ModbusPDUReadExceptionStatusRequest")

        write_buffer.pop_context("ModbusPDUReadExceptionStatusRequest")

    def length_in_bytes(self) -> int:
        return int(math.ceil(float(self.length_in_bits() / 8.0)))

    def length_in_bits(self) -> int:
        length_in_bits: int = super().length_in_bits()
        _value: ModbusPDUReadExceptionStatusRequest = self

        return length_in_bits

    @staticmethod
    def static_parse_builder(read_buffer: ReadBuffer, response: bool):
        read_buffer.push_context("ModbusPDUReadExceptionStatusRequest")

        read_buffer.pop_context("ModbusPDUReadExceptionStatusRequest")
        # Create the instance
        return ModbusPDUReadExceptionStatusRequestBuilder()

    def equals(self, o: object) -> bool:
        if self == o:
            return True

        if not isinstance(o, ModbusPDUReadExceptionStatusRequest):
            return False

        that: ModbusPDUReadExceptionStatusRequest = ModbusPDUReadExceptionStatusRequest(
            o
        )
        return super().equals(that) and True

    def hash_code(self) -> int:
        return hash(self)

    def __str__(self) -> str:
        pass
        # write_buffer_box_based: WriteBufferBoxBased = WriteBufferBoxBased(True, True)
        # try:
        #    write_buffer_box_based.writeSerializable(self)
        # except SerializationException as e:
        #    raise PlcRuntimeException(e)

        # return "\n" + str(write_buffer_box_based.get_box()) + "\n"


@dataclass
class ModbusPDUReadExceptionStatusRequestBuilder:

    def build(
        self,
    ) -> ModbusPDUReadExceptionStatusRequest:
        modbus_pduread_exception_status_request: ModbusPDUReadExceptionStatusRequest = (
            ModbusPDUReadExceptionStatusRequest()
        )
        return modbus_pduread_exception_status_request
