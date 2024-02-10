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
from plc4py.protocols.umas import StaticHelper
from plc4py.spi.generation.ReadBuffer import ReadBuffer
from plc4py.spi.generation.WriteBuffer import WriteBuffer
import math


@dataclass
class UmasUnlocatedVariableReference:
    data_type: int
    unknown1: int
    block: int
    offset: int
    unknown5: int
    unknown4: int
    string_length: int
    value: str

    def serialize(self, write_buffer: WriteBuffer):
        write_buffer.push_context("UmasUnlocatedVariableReference")

        # Simple Field (dataType)
        write_buffer.write_unsigned_byte(
            self.data_type, bit_length=8, logical_name="dataType"
        )

        # Simple Field (unknown1)
        write_buffer.write_unsigned_byte(
            self.unknown1, bit_length=8, logical_name="unknown1"
        )

        # Simple Field (block)
        write_buffer.write_unsigned_short(
            self.block, bit_length=16, logical_name="block"
        )

        # Simple Field (offset)
        write_buffer.write_unsigned_byte(
            self.offset, bit_length=8, logical_name="offset"
        )

        # Simple Field (unknown5)
        write_buffer.write_unsigned_byte(
            self.unknown5, bit_length=8, logical_name="unknown5"
        )

        # Simple Field (unknown4)
        write_buffer.write_unsigned_short(
            self.unknown4, bit_length=16, logical_name="unknown4"
        )

        # Simple Field (stringLength)
        write_buffer.write_unsigned_short(
            self.string_length, bit_length=16, logical_name="stringLength"
        )
        # Manual Field (value)
        write_buffer.write_manual(
            write_function=lambda: StaticHelper.serialize_terminated_string(
                write_buffer, self.value, self.string_length
            ),
            logical_name="value",
        )

        write_buffer.pop_context("UmasUnlocatedVariableReference")

    def length_in_bytes(self) -> int:
        return int(math.ceil(float(self.length_in_bits() / 8.0)))

    def length_in_bits(self) -> int:
        length_in_bits: int = 0
        _value: UmasUnlocatedVariableReference = self

        # Simple field (dataType)
        length_in_bits += 8

        # Simple field (unknown1)
        length_in_bits += 8

        # Simple field (block)
        length_in_bits += 16

        # Simple field (offset)
        length_in_bits += 8

        # Simple field (unknown5)
        length_in_bits += 8

        # Simple field (unknown4)
        length_in_bits += 16

        # Simple field (stringLength)
        length_in_bits += 16

        # Manual Field (value)
        length_in_bits += self.string_length * int(8)

        return length_in_bits

    @staticmethod
    def static_parse(read_buffer: ReadBuffer, **kwargs):
        return UmasUnlocatedVariableReference.static_parse_context(read_buffer)

    @staticmethod
    def static_parse_context(read_buffer: ReadBuffer):
        read_buffer.push_context("UmasUnlocatedVariableReference")

        data_type: int = read_buffer.read_unsigned_byte(
            logical_name="dataType", bit_length=8
        )

        unknown1: int = read_buffer.read_unsigned_byte(
            logical_name="unknown1", bit_length=8
        )

        block: int = read_buffer.read_unsigned_short(
            logical_name="block", bit_length=16
        )

        offset: int = read_buffer.read_unsigned_byte(
            logical_name="offset", bit_length=8
        )

        unknown5: int = read_buffer.read_unsigned_byte(
            logical_name="unknown5", bit_length=8
        )

        unknown4: int = read_buffer.read_unsigned_short(
            logical_name="unknown4", bit_length=16
        )

        string_length: int = read_buffer.read_unsigned_short(
            logical_name="stringLength", bit_length=16
        )

        value = read_buffer.read_manual(
            read_function=lambda: StaticHelper.parse_terminated_string(
                read_buffer, string_length
            ),
            logical_name="value",
        )

        read_buffer.pop_context("UmasUnlocatedVariableReference")
        # Create the instance
        _umas_unlocated_variable_reference: UmasUnlocatedVariableReference = (
            UmasUnlocatedVariableReference(
                data_type,
                unknown1,
                block,
                offset,
                unknown5,
                unknown4,
                string_length,
                value,
            )
        )
        return _umas_unlocated_variable_reference

    def equals(self, o: object) -> bool:
        if self == o:
            return True

        if not isinstance(o, UmasUnlocatedVariableReference):
            return False

        that: UmasUnlocatedVariableReference = UmasUnlocatedVariableReference(o)
        return (
            (self.data_type == that.data_type)
            and (self.unknown1 == that.unknown1)
            and (self.block == that.block)
            and (self.offset == that.offset)
            and (self.unknown5 == that.unknown5)
            and (self.unknown4 == that.unknown4)
            and (self.string_length == that.string_length)
            and (self.value == that.value)
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
