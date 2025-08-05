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
from plc4py.protocols.umas.readwrite.UmasMemoryBlock import UmasMemoryBlock
from plc4py.spi.generation.ReadBuffer import ReadBuffer
from plc4py.spi.generation.WriteBuffer import WriteBuffer
from typing import ClassVar
import math


@dataclass
class UmasMemoryBlockBasicInfo(UmasMemoryBlock):
    range: int
    not_sure: int
    index: int
    hardware_id: int
    # Accessors for discriminator values.
    block_number: ClassVar[int] = 0x30
    offset: ClassVar[int] = 0x00

    def serialize_umas_memory_block_child(self, write_buffer: WriteBuffer):
        write_buffer.push_context("UmasMemoryBlockBasicInfo")

        # Simple Field (range)
        write_buffer.write_unsigned_short(
            self.range, bit_length=16, logical_name="range"
        )

        # Simple Field (notSure)
        write_buffer.write_unsigned_short(
            self.not_sure, bit_length=16, logical_name="notSure"
        )

        # Simple Field (index)
        write_buffer.write_unsigned_byte(self.index, bit_length=8, logical_name="index")

        # Simple Field (hardwareId)
        write_buffer.write_unsigned_int(
            self.hardware_id, bit_length=32, logical_name="hardwareId"
        )

        write_buffer.pop_context("UmasMemoryBlockBasicInfo")

    def length_in_bytes(self) -> int:
        return int(math.ceil(float(self.length_in_bits() / 8.0)))

    def length_in_bits(self) -> int:
        length_in_bits: int = super().length_in_bits()
        _value: UmasMemoryBlockBasicInfo = self

        # Simple field (range)
        length_in_bits += 16

        # Simple field (notSure)
        length_in_bits += 16

        # Simple field (index)
        length_in_bits += 8

        # Simple field (hardwareId)
        length_in_bits += 32

        return length_in_bits

    @staticmethod
    def static_parse_builder(read_buffer: ReadBuffer, block_number: int, offset: int):
        read_buffer.push_context("UmasMemoryBlockBasicInfo")

        if isinstance(block_number, str):
            block_number = int(block_number)
        if isinstance(offset, str):
            offset = int(offset)

        range: int = read_buffer.read_unsigned_short(
            logical_name="range",
            bit_length=16,
            block_number=block_number,
            offset=offset,
        )

        not_sure: int = read_buffer.read_unsigned_short(
            logical_name="not_sure",
            bit_length=16,
            block_number=block_number,
            offset=offset,
        )

        index: int = read_buffer.read_unsigned_byte(
            logical_name="index", bit_length=8, block_number=block_number, offset=offset
        )

        hardware_id: int = read_buffer.read_unsigned_int(
            logical_name="hardware_id",
            bit_length=32,
            block_number=block_number,
            offset=offset,
        )

        read_buffer.pop_context("UmasMemoryBlockBasicInfo")
        # Create the instance
        return UmasMemoryBlockBasicInfoBuilder(range, not_sure, index, hardware_id)

    def equals(self, o: object) -> bool:
        if self == o:
            return True

        if not isinstance(o, UmasMemoryBlockBasicInfo):
            return False

        that: UmasMemoryBlockBasicInfo = UmasMemoryBlockBasicInfo(o)
        return (
            (self.range == that.range)
            and (self.not_sure == that.not_sure)
            and (self.index == that.index)
            and (self.hardware_id == that.hardware_id)
            and super().equals(that)
            and True
        )

    def hash_code(self) -> int:
        return hash(self)

    def __str__(self) -> str:
        # TODO:- Implement a generic python object to probably json convertor or something.
        return ""


@dataclass
class UmasMemoryBlockBasicInfoBuilder:
    range: int
    not_sure: int
    index: int
    hardware_id: int

    def build(
        self,
    ) -> UmasMemoryBlockBasicInfo:
        umas_memory_block_basic_info: UmasMemoryBlockBasicInfo = (
            UmasMemoryBlockBasicInfo(
                self.range, self.not_sure, self.index, self.hardware_id
            )
        )
        return umas_memory_block_basic_info
