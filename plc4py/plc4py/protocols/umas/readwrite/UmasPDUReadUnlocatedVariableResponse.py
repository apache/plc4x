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

import math
from dataclasses import dataclass
from typing import Any, ClassVar, List

from plc4py.api.exceptions.exceptions import (
    PlcRuntimeException,
    SerializationException,
)
from plc4py.api.messages.PlcMessage import PlcMessage
from plc4py.protocols.umas.readwrite.UmasPDUItem import UmasPDUItem
from plc4py.spi.generation.ReadBuffer import ReadBuffer
from plc4py.spi.generation.WriteBuffer import WriteBuffer
from plc4py.utils.GenericTypes import ByteOrder


@dataclass
class UmasPDUReadUnlocatedVariableResponse(UmasPDUItem):
    block: List[int]
    # Arguments.
    byte_length: int
    # Accessors for discriminator values.
    umas_function_key: ClassVar[int] = 0xFE
    umas_request_function_key: ClassVar[int] = 0x26

    def serialize_umas_pduitem_child(self, write_buffer: WriteBuffer):
        write_buffer.push_context("UmasPDUReadUnlocatedVariableResponse")

        # Array Field (block)
        write_buffer.write_simple_array(
            self.block, write_buffer.write_unsigned_byte, logical_name="block"
        )

        write_buffer.pop_context("UmasPDUReadUnlocatedVariableResponse")

    def length_in_bytes(self) -> int:
        return int(math.ceil(float(self.length_in_bits() / 8.0)))

    def length_in_bits(self) -> int:
        length_in_bits: int = super().length_in_bits()
        _value: UmasPDUReadUnlocatedVariableResponse = self

        # Array field
        if self.block is not None:
            length_in_bits += 8 * len(self.block)

        return length_in_bits

    @staticmethod
    def static_parse_builder(
        read_buffer: ReadBuffer,
        umas_request_function_key: int,
        byte_length: int,
    ):
        read_buffer.push_context("UmasPDUReadUnlocatedVariableResponse")

        block: List[Any] = read_buffer.read_array_field(
            logical_name="block",
            read_function=read_buffer.read_unsigned_byte,
            count=byte_length - int(2),
            byte_order=ByteOrder.LITTLE_ENDIAN,
            umas_request_function_key=umas_request_function_key,
            byte_length=byte_length,
        )

        read_buffer.pop_context("UmasPDUReadUnlocatedVariableResponse")
        # Create the instance
        return UmasPDUReadUnlocatedVariableResponseBuilder(block)

    def equals(self, o: object) -> bool:
        if self == o:
            return True

        if not isinstance(o, UmasPDUReadUnlocatedVariableResponse):
            return False

        that: UmasPDUReadUnlocatedVariableResponse = (
            UmasPDUReadUnlocatedVariableResponse(o)
        )
        return (self.block == that.block) and super().equals(that) and True

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
class UmasPDUReadUnlocatedVariableResponseBuilder:
    block: List[int]

    def build(
        self, byte_length: int, pairing_key
    ) -> UmasPDUReadUnlocatedVariableResponse:
        umas_pduread_unlocated_variable_response: UmasPDUReadUnlocatedVariableResponse = UmasPDUReadUnlocatedVariableResponse(
            byte_length, pairing_key, self.block
        )
        return umas_pduread_unlocated_variable_response
