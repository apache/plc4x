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
    # Accessors for discriminator values.
    umas_function_key: ClassVar[int] = 0x20
    umas_request_function_key: ClassVar[int] = 0

    def serialize_umas_pdu_item_child(self, write_buffer: WriteBuffer):
        write_buffer.push_context("UmasPDUReadMemoryBlockRequest")

        # Simple Field (range)
        write_buffer.write_unsigned_byte(self.range, logical_name="range")

        # Simple Field (blockNumber)
        write_buffer.write_unsigned_short(self.block_number, logical_name="blockNumber")

        # Simple Field (offset)
        write_buffer.write_unsigned_short(self.offset, logical_name="offset")

        # Simple Field (unknownObject1)
        write_buffer.write_unsigned_short(
            self.unknown_object1, logical_name="unknownObject1"
        )

        # Simple Field (numberOfBytes)
        write_buffer.write_unsigned_short(
            self.number_of_bytes, logical_name="numberOfBytes"
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
    def static_parse_builder(read_buffer: ReadBuffer, umas_request_function_key: int):
        read_buffer.push_context("UmasPDUReadMemoryBlockRequest")

        range: int = read_buffer.read_unsigned_byte(
            logical_name="range",
            bit_length=8,
            byte_order=ByteOrder.LITTLE_ENDIAN,
            umas_request_function_key=umas_request_function_key,
        )

        block_number: int = read_buffer.read_unsigned_short(
            logical_name="blockNumber",
            bit_length=16,
            byte_order=ByteOrder.LITTLE_ENDIAN,
            umas_request_function_key=umas_request_function_key,
        )

        offset: int = read_buffer.read_unsigned_short(
            logical_name="offset",
            bit_length=16,
            byte_order=ByteOrder.LITTLE_ENDIAN,
            umas_request_function_key=umas_request_function_key,
        )

        unknown_object1: int = read_buffer.read_unsigned_short(
            logical_name="unknownObject1",
            bit_length=16,
            byte_order=ByteOrder.LITTLE_ENDIAN,
            umas_request_function_key=umas_request_function_key,
        )

        number_of_bytes: int = read_buffer.read_unsigned_short(
            logical_name="numberOfBytes",
            bit_length=16,
            byte_order=ByteOrder.LITTLE_ENDIAN,
            umas_request_function_key=umas_request_function_key,
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
        pass
        # write_buffer_box_based: WriteBufferBoxBased = WriteBufferBoxBased(True, True)
        # try:
        #    write_buffer_box_based.writeSerializable(self)
        # except SerializationException as e:
        #    raise PlcRuntimeException(e)

        # return "\n" + str(write_buffer_box_based.get_box()) + "\n"


@dataclass
class UmasPDUReadMemoryBlockRequestBuilder:
    range: int
    block_number: int
    offset: int
    unknown_object1: int
    number_of_bytes: int

    def build(self, pairing_key) -> UmasPDUReadMemoryBlockRequest:
        umas_pdu_read_memory_block_request: UmasPDUReadMemoryBlockRequest = (
            UmasPDUReadMemoryBlockRequest(
                pairing_key,
                self.range,
                self.block_number,
                self.offset,
                self.unknown_object1,
                self.number_of_bytes,
            )
        )
        return umas_pdu_read_memory_block_request
