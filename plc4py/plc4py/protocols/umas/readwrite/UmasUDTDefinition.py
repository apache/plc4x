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
class UmasUDTDefinition:
    data_type: int
    offset: int
    unknown5: int
    unknown4: int
    value: str

    def serialize(self, write_buffer: WriteBuffer):
        write_buffer.push_context("UmasUDTDefinition")

        # Simple Field (dataType)
        write_buffer.write_unsigned_short(
            self.data_type, bit_length=16, logical_name="dataType"
        )

        # Simple Field (offset)
        write_buffer.write_unsigned_short(
            self.offset, bit_length=16, logical_name="offset"
        )

        # Simple Field (unknown5)
        write_buffer.write_unsigned_short(
            self.unknown5, bit_length=16, logical_name="unknown5"
        )

        # Simple Field (unknown4)
        write_buffer.write_unsigned_short(
            self.unknown4, bit_length=16, logical_name="unknown4"
        )
        # Manual Field (value)
        write_buffer.write_manual(
            write_function=lambda: StaticHelper.serialize_terminated_string(
                write_buffer, self.value, 1
            ),
            logical_name="value",
        )

        write_buffer.pop_context("UmasUDTDefinition")

    def length_in_bytes(self) -> int:
        return int(math.ceil(float(self.length_in_bits() / 8.0)))

    def length_in_bits(self) -> int:
        length_in_bits: int = 0
        _value: UmasUDTDefinition = self

        # Simple field (dataType)
        length_in_bits += 16

        # Simple field (offset)
        length_in_bits += 16

        # Simple field (unknown5)
        length_in_bits += 16

        # Simple field (unknown4)
        length_in_bits += 16

        # Manual Field (value)
        length_in_bits += self.string_length * int(8)

        return length_in_bits

    @staticmethod
    def static_parse(read_buffer: ReadBuffer, **kwargs):
        return UmasUDTDefinition.static_parse_context(read_buffer)

    @staticmethod
    def static_parse_context(read_buffer: ReadBuffer):
        read_buffer.push_context("UmasUDTDefinition")

        data_type: int = read_buffer.read_unsigned_short(
            logical_name="data_type", bit_length=16
        )

        offset: int = read_buffer.read_unsigned_short(
            logical_name="offset", bit_length=16
        )

        unknown5: int = read_buffer.read_unsigned_short(
            logical_name="unknown5", bit_length=16
        )

        unknown4: int = read_buffer.read_unsigned_short(
            logical_name="unknown4", bit_length=16
        )

        value = read_buffer.read_manual(
            read_function=lambda: StaticHelper.parse_terminated_string(read_buffer, 1),
            logical_name="value",
        )

        read_buffer.pop_context("UmasUDTDefinition")
        # Create the instance
        _umas_udtdefinition: UmasUDTDefinition = UmasUDTDefinition(
            data_type, offset, unknown5, unknown4, value
        )
        return _umas_udtdefinition

    def equals(self, o: object) -> bool:
        if self == o:
            return True

        if not isinstance(o, UmasUDTDefinition):
            return False

        that: UmasUDTDefinition = UmasUDTDefinition(o)
        return (
            (self.data_type == that.data_type)
            and (self.offset == that.offset)
            and (self.unknown5 == that.unknown5)
            and (self.unknown4 == that.unknown4)
            and (self.value == that.value)
            and True
        )

    def hash_code(self) -> int:
        return hash(self)

    def __str__(self) -> str:
        # TODO:- Implement a generic python object to probably json convertor or something.
        return ""
