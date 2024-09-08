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
class UmasPDUPlcStatusResponse(UmasPDUItem):
    not_used1: int
    not_used2: int
    number_of_blocks: int
    blocks: List[int]
    # Arguments.
    byte_length: int
    # Accessors for discriminator values.
    umas_function_key: ClassVar[int] = 0xFE
    umas_request_function_key: ClassVar[int] = 0x04

    def serialize_umas_pduitem_child(self, write_buffer: WriteBuffer):
        write_buffer.push_context("UmasPDUPlcStatusResponse")

        # Simple Field (notUsed1)
        write_buffer.write_unsigned_byte(
            self.not_used1, bit_length=8, logical_name="notUsed1"
        )

        # Simple Field (notUsed2)
        write_buffer.write_unsigned_short(
            self.not_used2, bit_length=16, logical_name="notUsed2"
        )

        # Simple Field (numberOfBlocks)
        write_buffer.write_unsigned_byte(
            self.number_of_blocks, bit_length=8, logical_name="numberOfBlocks"
        )

        # Array Field (blocks)
        write_buffer.write_simple_array(
            self.blocks, write_buffer.write_unsigned_int, logical_name="blocks"
        )

        write_buffer.pop_context("UmasPDUPlcStatusResponse")

    def length_in_bytes(self) -> int:
        return int(math.ceil(float(self.length_in_bits() / 8.0)))

    def length_in_bits(self) -> int:
        length_in_bits: int = super().length_in_bits()
        _value: UmasPDUPlcStatusResponse = self

        # Simple field (notUsed1)
        length_in_bits += 8

        # Simple field (notUsed2)
        length_in_bits += 16

        # Simple field (numberOfBlocks)
        length_in_bits += 8

        # Array field
        if self.blocks is not None:
            length_in_bits += 32 * len(self.blocks)

        return length_in_bits

    @staticmethod
    def static_parse_builder(
        read_buffer: ReadBuffer, umas_request_function_key: int, byte_length: int
    ):
        read_buffer.push_context("UmasPDUPlcStatusResponse")

        if isinstance(umas_request_function_key, str):
            umas_request_function_key = int(umas_request_function_key)
        if isinstance(byte_length, str):
            byte_length = int(byte_length)

        not_used1: int = read_buffer.read_unsigned_byte(
            logical_name="not_used1",
            bit_length=8,
            byte_order=ByteOrder.LITTLE_ENDIAN,
            umas_request_function_key=umas_request_function_key,
            byte_length=byte_length,
        )

        not_used2: int = read_buffer.read_unsigned_short(
            logical_name="not_used2",
            bit_length=16,
            byte_order=ByteOrder.LITTLE_ENDIAN,
            umas_request_function_key=umas_request_function_key,
            byte_length=byte_length,
        )

        number_of_blocks: int = read_buffer.read_unsigned_byte(
            logical_name="number_of_blocks",
            bit_length=8,
            byte_order=ByteOrder.LITTLE_ENDIAN,
            umas_request_function_key=umas_request_function_key,
            byte_length=byte_length,
        )

        blocks: List[Any] = read_buffer.read_array_field(
            logical_name="blocks",
            read_function=read_buffer.read_unsigned_int,
            count=number_of_blocks,
            byte_order=ByteOrder.LITTLE_ENDIAN,
            umas_request_function_key=umas_request_function_key,
            byte_length=byte_length,
        )

        read_buffer.pop_context("UmasPDUPlcStatusResponse")
        # Create the instance
        return UmasPDUPlcStatusResponseBuilder(
            not_used1, not_used2, number_of_blocks, blocks
        )

    def equals(self, o: object) -> bool:
        if self == o:
            return True

        if not isinstance(o, UmasPDUPlcStatusResponse):
            return False

        that: UmasPDUPlcStatusResponse = UmasPDUPlcStatusResponse(o)
        return (
            (self.not_used1 == that.not_used1)
            and (self.not_used2 == that.not_used2)
            and (self.number_of_blocks == that.number_of_blocks)
            and (self.blocks == that.blocks)
            and super().equals(that)
            and True
        )

    def hash_code(self) -> int:
        return hash(self)

    def __str__(self) -> str:
        # TODO:- Implement a generic python object to probably json convertor or something.
        return ""


@dataclass
class UmasPDUPlcStatusResponseBuilder:
    not_used1: int
    not_used2: int
    number_of_blocks: int
    blocks: List[int]

    def build(self, byte_length: int, pairing_key) -> UmasPDUPlcStatusResponse:
        umas_pduplc_status_response: UmasPDUPlcStatusResponse = (
            UmasPDUPlcStatusResponse(
                byte_length,
                pairing_key,
                self.not_used1,
                self.not_used2,
                self.number_of_blocks,
                self.blocks,
            )
        )
        return umas_pduplc_status_response
