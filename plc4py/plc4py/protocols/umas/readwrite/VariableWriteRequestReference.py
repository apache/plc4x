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
from typing import Any
from typing import List
import math


@dataclass
class VariableWriteRequestReference:
    is_array: int
    data_size_index: int
    block: int
    base_offset: int
    offset: int
    array_length: int
    record_data: List[int]

    def serialize(self, write_buffer: WriteBuffer):
        write_buffer.push_context("VariableWriteRequestReference")

        # Simple Field (isArray)
        write_buffer.write_unsigned_byte(
            self.is_array, bit_length=4, logical_name="isArray"
        )

        # Simple Field (dataSizeIndex)
        write_buffer.write_unsigned_byte(
            self.data_size_index, bit_length=4, logical_name="dataSizeIndex"
        )

        # Simple Field (block)
        write_buffer.write_unsigned_short(
            self.block, bit_length=16, logical_name="block"
        )

        # Simple Field (baseOffset)
        write_buffer.write_unsigned_short(
            self.base_offset, bit_length=16, logical_name="baseOffset"
        )

        # Simple Field (offset)
        write_buffer.write_unsigned_short(
            self.offset, bit_length=16, logical_name="offset"
        )

        # Optional Field (arrayLength) (Can be skipped, if the value is null)
        if self.is_array:
            write_buffer.write_unsigned_short(
                self.array_length, logical_name="arrayLength"
            )

        # Array Field (recordData)
        write_buffer.write_byte_array(self.record_data, logical_name="record_data")

        write_buffer.pop_context("VariableWriteRequestReference")

    def length_in_bytes(self) -> int:
        return int(math.ceil(float(self.length_in_bits() / 8.0)))

    def length_in_bits(self) -> int:
        length_in_bits: int = 0
        _value: VariableWriteRequestReference = self

        # Simple field (isArray)
        length_in_bits += 4

        # Simple field (dataSizeIndex)
        length_in_bits += 4

        # Simple field (block)
        length_in_bits += 16

        # Simple field (baseOffset)
        length_in_bits += 16

        # Simple field (offset)
        length_in_bits += 16

        # Optional Field (arrayLength)
        if self.is_array:
            length_in_bits += 16

        # Array field
        if self.record_data is not None:
            length_in_bits += 8 * len(self.record_data)

        return length_in_bits

    @staticmethod
    def static_parse(read_buffer: ReadBuffer, **kwargs):
        return VariableWriteRequestReference.static_parse_context(read_buffer)

    @staticmethod
    def static_parse_context(read_buffer: ReadBuffer):
        read_buffer.push_context("VariableWriteRequestReference")

        is_array: int = read_buffer.read_unsigned_byte(
            logical_name="is_array", bit_length=4
        )

        data_size_index: int = read_buffer.read_unsigned_byte(
            logical_name="data_size_index", bit_length=4
        )

        block: int = read_buffer.read_unsigned_short(
            logical_name="block", bit_length=16
        )

        base_offset: int = read_buffer.read_unsigned_short(
            logical_name="base_offset", bit_length=16
        )

        offset: int = read_buffer.read_unsigned_short(
            logical_name="offset", bit_length=16
        )

        array_length: int = None
        if is_array:
            array_length = read_buffer.read_unsigned_short(logical_name="array_length")

        record_data: List[Any] = read_buffer.read_array_field(
            logical_name="recordData",
            read_function=read_buffer.read_byte,
            count=(
                data_size_index * array_length
                if bool((is_array) == ((1)))
                else data_size_index
            ),
        )

        read_buffer.pop_context("VariableWriteRequestReference")
        # Create the instance
        _variable_write_request_reference: VariableWriteRequestReference = (
            VariableWriteRequestReference(
                is_array,
                data_size_index,
                block,
                base_offset,
                offset,
                array_length,
                record_data,
            )
        )
        return _variable_write_request_reference

    def equals(self, o: object) -> bool:
        if self == o:
            return True

        if not isinstance(o, VariableWriteRequestReference):
            return False

        that: VariableWriteRequestReference = VariableWriteRequestReference(o)
        return (
            (self.is_array == that.is_array)
            and (self.data_size_index == that.data_size_index)
            and (self.block == that.block)
            and (self.base_offset == that.base_offset)
            and (self.offset == that.offset)
            and (self.array_length == that.array_length)
            and (self.record_data == that.record_data)
            and True
        )

    def hash_code(self) -> int:
        return hash(self)

    def __str__(self) -> str:
        # TODO:- Implement a generic python object to probably json convertor or something.
        return ""
