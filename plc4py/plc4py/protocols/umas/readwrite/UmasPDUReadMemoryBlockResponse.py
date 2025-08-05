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
from plc4py.protocols.umas.readwrite.UmasPDUItem import UmasPDUItem
from plc4py.spi.generation.ReadBuffer import ReadBuffer
from plc4py.spi.generation.WriteBuffer import WriteBuffer
from plc4py.utils.GenericTypes import ByteOrder
from typing import Any
from typing import ClassVar
from typing import List
import math


@dataclass
class UmasPDUReadMemoryBlockResponse(UmasPDUItem):
    range: int
    number_of_bytes: int
    block: List[int]
    # Arguments.
    byte_length: int
    # Accessors for discriminator values.
    umas_function_key: ClassVar[int] = 0xFE
    umas_request_function_key: ClassVar[int] = 0x20

    def serialize_umas_pduitem_child(self, write_buffer: WriteBuffer):
        write_buffer.push_context("UmasPDUReadMemoryBlockResponse")

        # Simple Field (range)
        write_buffer.write_unsigned_byte(self.range, bit_length=8, logical_name="range")

        # Simple Field (numberOfBytes)
        write_buffer.write_unsigned_short(
            self.number_of_bytes, bit_length=16, logical_name="numberOfBytes"
        )

        # Array Field (block)
        write_buffer.write_simple_array(
            self.block, write_buffer.write_unsigned_byte, logical_name="block"
        )

        write_buffer.pop_context("UmasPDUReadMemoryBlockResponse")

    def length_in_bytes(self) -> int:
        return int(math.ceil(float(self.length_in_bits() / 8.0)))

    def length_in_bits(self) -> int:
        length_in_bits: int = super().length_in_bits()
        _value: UmasPDUReadMemoryBlockResponse = self

        # Simple field (range)
        length_in_bits += 8

        # Simple field (numberOfBytes)
        length_in_bits += 16

        # Array field
        if self.block is not None:
            length_in_bits += 8 * len(self.block)

        return length_in_bits

    @staticmethod
    def static_parse_builder(
        read_buffer: ReadBuffer, umas_request_function_key: int, byte_length: int
    ):
        read_buffer.push_context("UmasPDUReadMemoryBlockResponse")

        if isinstance(umas_request_function_key, str):
            umas_request_function_key = int(umas_request_function_key)
        if isinstance(byte_length, str):
            byte_length = int(byte_length)

        range: int = read_buffer.read_unsigned_byte(
            logical_name="range",
            bit_length=8,
            byte_order=ByteOrder.LITTLE_ENDIAN,
            umas_request_function_key=umas_request_function_key,
            byte_length=byte_length,
        )

        number_of_bytes: int = read_buffer.read_unsigned_short(
            logical_name="number_of_bytes",
            bit_length=16,
            byte_order=ByteOrder.LITTLE_ENDIAN,
            umas_request_function_key=umas_request_function_key,
            byte_length=byte_length,
        )

        block: List[Any] = read_buffer.read_array_field(
            logical_name="block",
            read_function=read_buffer.read_unsigned_byte,
            count=number_of_bytes,
            byte_order=ByteOrder.LITTLE_ENDIAN,
            umas_request_function_key=umas_request_function_key,
            byte_length=byte_length,
        )

        read_buffer.pop_context("UmasPDUReadMemoryBlockResponse")
        # Create the instance
        return UmasPDUReadMemoryBlockResponseBuilder(range, number_of_bytes, block)

    def equals(self, o: object) -> bool:
        if self == o:
            return True

        if not isinstance(o, UmasPDUReadMemoryBlockResponse):
            return False

        that: UmasPDUReadMemoryBlockResponse = UmasPDUReadMemoryBlockResponse(o)
        return (
            (self.range == that.range)
            and (self.number_of_bytes == that.number_of_bytes)
            and (self.block == that.block)
            and super().equals(that)
            and True
        )

    def hash_code(self) -> int:
        return hash(self)

    def __str__(self) -> str:
        # TODO:- Implement a generic python object to probably json convertor or something.
        return ""


@dataclass
class UmasPDUReadMemoryBlockResponseBuilder:
    range: int
    number_of_bytes: int
    block: List[int]

    def build(self, byte_length: int, pairing_key) -> UmasPDUReadMemoryBlockResponse:
        umas_pduread_memory_block_response: UmasPDUReadMemoryBlockResponse = (
            UmasPDUReadMemoryBlockResponse(
                byte_length, pairing_key, self.range, self.number_of_bytes, self.block
            )
        )
        return umas_pduread_memory_block_response
