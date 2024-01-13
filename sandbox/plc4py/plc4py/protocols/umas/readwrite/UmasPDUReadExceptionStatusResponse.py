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
from plc4py.protocols.umas.readwrite.UmasPDU import UmasPDU
from plc4py.protocols.umas.readwrite.UmasPDU import UmasPDUBuilder
from plc4py.spi.generation.ReadBuffer import ReadBuffer
from plc4py.spi.generation.WriteBuffer import WriteBuffer
import math


@dataclass
class UmasPDUReadExceptionStatusResponse(UmasPDU):
    value: int
    # Accessors for discriminator values.
    error_flag: bool = False
    function_flag: int = 0x07
    response: bool = True

    def serialize_umas_pdu_child(self, write_buffer: WriteBuffer):
        write_buffer.push_context("UmasPDUReadExceptionStatusResponse")

        # Simple Field (value)
        write_buffer.write_unsigned_byte(self.value, logical_name="value")

        write_buffer.pop_context("UmasPDUReadExceptionStatusResponse")

    def length_in_bytes(self) -> int:
        return int(math.ceil(float(self.length_in_bits() / 8.0)))

    def length_in_bits(self) -> int:
        length_in_bits: int = super().length_in_bits()
        _value: UmasPDUReadExceptionStatusResponse = self

        # Simple field (value)
        length_in_bits += 8

        return length_in_bits

    @staticmethod
    def static_parse_builder(read_buffer: ReadBuffer, response: bool):
        read_buffer.push_context("UmasPDUReadExceptionStatusResponse")

        value: int = read_buffer.read_unsigned_byte(
            logical_name="value", bit_length=8, response=response
        )

        read_buffer.pop_context("UmasPDUReadExceptionStatusResponse")
        # Create the instance
        return UmasPDUReadExceptionStatusResponseBuilder(value)

    def equals(self, o: object) -> bool:
        if self == o:
            return True

        if not isinstance(o, UmasPDUReadExceptionStatusResponse):
            return False

        that: UmasPDUReadExceptionStatusResponse = UmasPDUReadExceptionStatusResponse(o)
        return (self.value == that.value) and super().equals(that) and True

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
class UmasPDUReadExceptionStatusResponseBuilder(UmasPDUBuilder):
    value: int

    def build(
        self,
    ) -> UmasPDUReadExceptionStatusResponse:
        umas_pdu_read_exception_status_response: UmasPDUReadExceptionStatusResponse = (
            UmasPDUReadExceptionStatusResponse(self.value)
        )
        return umas_pdu_read_exception_status_response
