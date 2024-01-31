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
class VariableRequestReference:
    data_type: int
    block: int
    base_offset: int
    offset: int
    UNKNOWN1: int = 0x01

    def serialize(self, write_buffer: WriteBuffer):
        write_buffer.push_context("VariableRequestReference")

        # Simple Field (dataType)
        write_buffer.write_unsigned_byte(
            self.data_type, bit_length=8, logical_name="dataType"
        )

        # Simple Field (block)
        write_buffer.write_unsigned_short(
            self.block, bit_length=16, logical_name="block"
        )

        # Const Field (unknown1)
        write_buffer.write_unsigned_byte(self.UNKNOWN1, logical_name="unknown1")

        # Simple Field (baseOffset)
        write_buffer.write_unsigned_short(
            self.base_offset, bit_length=16, logical_name="baseOffset"
        )

        # Simple Field (offset)
        write_buffer.write_unsigned_byte(
            self.offset, bit_length=8, logical_name="offset"
        )

        write_buffer.pop_context("VariableRequestReference")

    def length_in_bytes(self) -> int:
        return int(math.ceil(float(self.length_in_bits() / 8.0)))

    def length_in_bits(self) -> int:
        length_in_bits: int = 0
        _value: VariableRequestReference = self

        # Simple field (dataType)
        length_in_bits += 8

        # Simple field (block)
        length_in_bits += 16

        # Const Field (unknown1)
        length_in_bits += 8

        # Simple field (baseOffset)
        length_in_bits += 16

        # Simple field (offset)
        length_in_bits += 8

        return length_in_bits

    @staticmethod
    def static_parse(read_buffer: ReadBuffer, **kwargs):
        return VariableRequestReference.static_parse_context(read_buffer)

    @staticmethod
    def static_parse_context(read_buffer: ReadBuffer):
        read_buffer.push_context("VariableRequestReference")

        data_type: int = read_buffer.read_unsigned_byte(
            logical_name="dataType", bit_length=8
        )

        block: int = read_buffer.read_unsigned_short(
            logical_name="block", bit_length=16
        )

        UNKNOWN1: int = read_buffer.read_unsigned_byte(logical_name="unknown1")

        base_offset: int = read_buffer.read_unsigned_short(
            logical_name="baseOffset", bit_length=16
        )

        offset: int = read_buffer.read_unsigned_byte(
            logical_name="offset", bit_length=8
        )

        read_buffer.pop_context("VariableRequestReference")
        # Create the instance
        _variable_request_reference: VariableRequestReference = (
            VariableRequestReference(data_type, block, base_offset, offset)
        )
        return _variable_request_reference

    def equals(self, o: object) -> bool:
        if self == o:
            return True

        if not isinstance(o, VariableRequestReference):
            return False

        that: VariableRequestReference = VariableRequestReference(o)
        return (
            (self.data_type == that.data_type)
            and (self.block == that.block)
            and (self.base_offset == that.base_offset)
            and (self.offset == that.offset)
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
