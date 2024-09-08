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
from typing import ClassVar
import math


@dataclass
class UmasPDUReadMemoryBlockRequest(UmasPDUItem):
    range: int
    block_number: int
    offset: int
    unknown_object1: int
    number_of_bytes: int
    # Arguments.
    byte_length: int
    # Accessors for discriminator values.
    umas_function_key: ClassVar[int] = 0x20
    umas_request_function_key: ClassVar[int] = 0

    def serialize_umas_pduitem_child(self, write_buffer: WriteBuffer):
        write_buffer.push_context("UmasPDUReadMemoryBlockRequest")

        # Simple Field (range)
        write_buffer.write_unsigned_byte(self.range, bit_length=8, logical_name="range")

        # Simple Field (blockNumber)
        write_buffer.write_unsigned_short(
            self.block_number, bit_length=16, logical_name="blockNumber"
        )

        # Simple Field (offset)
        write_buffer.write_unsigned_short(
            self.offset, bit_length=16, logical_name="offset"
        )

        # Simple Field (unknownObject1)
        write_buffer.write_unsigned_short(
            self.unknown_object1, bit_length=16, logical_name="unknownObject1"
        )

        # Simple Field (numberOfBytes)
        write_buffer.write_unsigned_short(
            self.number_of_bytes, bit_length=16, logical_name="numberOfBytes"
        )

        write_buffer.pop_context("UmasPDUReadMemoryBlockRequest")

    def length_in_bytes(self) -> int:
        return int(math.ceil(float(self.length_in_bits() / 8.0)))

    def length_in_bits(self) -> int:
        length_in_bits: int = super().length_in_bits()
        _value: UmasPDUReadMemoryBlockRequest = self

        # Simple field (range)
        length_in_bits += 8

        # Simple field (blockNumber)
        length_in_bits += 16

        # Simple field (offset)
        length_in_bits += 16

        # Simple field (unknownObject1)
        length_in_bits += 16

        # Simple field (numberOfBytes)
        length_in_bits += 16

        return length_in_bits

    @staticmethod
    def static_parse_builder(
        read_buffer: ReadBuffer, umas_request_function_key: int, byte_length: int
    ):
        read_buffer.push_context("UmasPDUReadMemoryBlockRequest")

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

        block_number: int = read_buffer.read_unsigned_short(
            logical_name="block_number",
            bit_length=16,
            byte_order=ByteOrder.LITTLE_ENDIAN,
            umas_request_function_key=umas_request_function_key,
            byte_length=byte_length,
        )

        offset: int = read_buffer.read_unsigned_short(
            logical_name="offset",
            bit_length=16,
            byte_order=ByteOrder.LITTLE_ENDIAN,
            umas_request_function_key=umas_request_function_key,
            byte_length=byte_length,
        )

        unknown_object1: int = read_buffer.read_unsigned_short(
            logical_name="unknown_object1",
            bit_length=16,
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

        read_buffer.pop_context("UmasPDUReadMemoryBlockRequest")
        # Create the instance
        return UmasPDUReadMemoryBlockRequestBuilder(
            range, block_number, offset, unknown_object1, number_of_bytes
        )

    def equals(self, o: object) -> bool:
        if self == o:
            return True

        if not isinstance(o, UmasPDUReadMemoryBlockRequest):
            return False

        that: UmasPDUReadMemoryBlockRequest = UmasPDUReadMemoryBlockRequest(o)
        return (
            (self.range == that.range)
            and (self.block_number == that.block_number)
            and (self.offset == that.offset)
            and (self.unknown_object1 == that.unknown_object1)
            and (self.number_of_bytes == that.number_of_bytes)
            and super().equals(that)
            and True
        )

    def hash_code(self) -> int:
        return hash(self)

    def __str__(self) -> str:
        # TODO:- Implement a generic python object to probably json convertor or something.
        return ""


@dataclass
class UmasPDUReadMemoryBlockRequestBuilder:
    range: int
    block_number: int
    offset: int
    unknown_object1: int
    number_of_bytes: int

    def build(self, byte_length: int, pairing_key) -> UmasPDUReadMemoryBlockRequest:
        umas_pduread_memory_block_request: UmasPDUReadMemoryBlockRequest = (
            UmasPDUReadMemoryBlockRequest(
                byte_length,
                pairing_key,
                self.range,
                self.block_number,
                self.offset,
                self.unknown_object1,
                self.number_of_bytes,
            )
        )
        return umas_pduread_memory_block_request
