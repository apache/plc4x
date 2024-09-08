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
from plc4py.spi.generation.ReadBuffer import ReadBuffer
from plc4py.spi.generation.WriteBuffer import WriteBuffer
import math


@dataclass
class PlcMemoryBlockIdent:
    block_type: int
    folio: int
    status: int
    memory_length: int

    def serialize(self, write_buffer: WriteBuffer):
        write_buffer.push_context("PlcMemoryBlockIdent")

        # Simple Field (blockType)
        write_buffer.write_unsigned_byte(
            self.block_type, bit_length=8, logical_name="blockType"
        )

        # Simple Field (folio)
        write_buffer.write_unsigned_byte(self.folio, bit_length=8, logical_name="folio")

        # Simple Field (status)
        write_buffer.write_unsigned_short(
            self.status, bit_length=16, logical_name="status"
        )

        # Simple Field (memoryLength)
        write_buffer.write_unsigned_int(
            self.memory_length, bit_length=32, logical_name="memoryLength"
        )

        write_buffer.pop_context("PlcMemoryBlockIdent")

    def length_in_bytes(self) -> int:
        return int(math.ceil(float(self.length_in_bits() / 8.0)))

    def length_in_bits(self) -> int:
        length_in_bits: int = 0
        _value: PlcMemoryBlockIdent = self

        # Simple field (blockType)
        length_in_bits += 8

        # Simple field (folio)
        length_in_bits += 8

        # Simple field (status)
        length_in_bits += 16

        # Simple field (memoryLength)
        length_in_bits += 32

        return length_in_bits

    @staticmethod
    def static_parse(read_buffer: ReadBuffer, **kwargs):
        return PlcMemoryBlockIdent.static_parse_context(read_buffer)

    @staticmethod
    def static_parse_context(read_buffer: ReadBuffer):
        read_buffer.push_context("PlcMemoryBlockIdent")

        block_type: int = read_buffer.read_unsigned_byte(
            logical_name="block_type", bit_length=8
        )

        folio: int = read_buffer.read_unsigned_byte(logical_name="folio", bit_length=8)

        status: int = read_buffer.read_unsigned_short(
            logical_name="status", bit_length=16
        )

        memory_length: int = read_buffer.read_unsigned_int(
            logical_name="memory_length", bit_length=32
        )

        read_buffer.pop_context("PlcMemoryBlockIdent")
        # Create the instance
        _plc_memory_block_ident: PlcMemoryBlockIdent = PlcMemoryBlockIdent(
            block_type, folio, status, memory_length
        )
        return _plc_memory_block_ident

    def equals(self, o: object) -> bool:
        if self == o:
            return True

        if not isinstance(o, PlcMemoryBlockIdent):
            return False

        that: PlcMemoryBlockIdent = PlcMemoryBlockIdent(o)
        return (
            (self.block_type == that.block_type)
            and (self.folio == that.folio)
            and (self.status == that.status)
            and (self.memory_length == that.memory_length)
            and True
        )

    def hash_code(self) -> int:
        return hash(self)

    def __str__(self) -> str:
        # TODO:- Implement a generic python object to probably json convertor or something.
        return ""
