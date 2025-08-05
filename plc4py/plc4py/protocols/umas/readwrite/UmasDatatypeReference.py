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
class UmasDatatypeReference:
    data_size: int
    unknown1: int
    class_identifier: int
    data_type: int
    string_length: int
    value: str

    def serialize(self, write_buffer: WriteBuffer):
        write_buffer.push_context("UmasDatatypeReference")

        # Simple Field (dataSize)
        write_buffer.write_unsigned_short(
            self.data_size, bit_length=16, logical_name="dataSize"
        )

        # Simple Field (unknown1)
        write_buffer.write_unsigned_short(
            self.unknown1, bit_length=16, logical_name="unknown1"
        )

        # Simple Field (classIdentifier)
        write_buffer.write_unsigned_byte(
            self.class_identifier, bit_length=8, logical_name="classIdentifier"
        )

        # Simple Field (dataType)
        write_buffer.write_unsigned_byte(
            self.data_type, bit_length=8, logical_name="dataType"
        )

        # Simple Field (stringLength)
        write_buffer.write_unsigned_byte(
            self.string_length, bit_length=8, logical_name="stringLength"
        )
        # Manual Field (value)
        write_buffer.write_manual(
            write_function=lambda: StaticHelper.serialize_terminated_string(
                write_buffer, self.value, self.string_length
            ),
            logical_name="value",
        )

        write_buffer.pop_context("UmasDatatypeReference")

    def length_in_bytes(self) -> int:
        return int(math.ceil(float(self.length_in_bits() / 8.0)))

    def length_in_bits(self) -> int:
        length_in_bits: int = 0
        _value: UmasDatatypeReference = self

        # Simple field (dataSize)
        length_in_bits += 16

        # Simple field (unknown1)
        length_in_bits += 16

        # Simple field (classIdentifier)
        length_in_bits += 8

        # Simple field (dataType)
        length_in_bits += 8

        # Simple field (stringLength)
        length_in_bits += 8

        # Manual Field (value)
        length_in_bits += self.string_length * int(8)

        return length_in_bits

    @staticmethod
    def static_parse(read_buffer: ReadBuffer, **kwargs):
        return UmasDatatypeReference.static_parse_context(read_buffer)

    @staticmethod
    def static_parse_context(read_buffer: ReadBuffer):
        read_buffer.push_context("UmasDatatypeReference")

        data_size: int = read_buffer.read_unsigned_short(
            logical_name="data_size", bit_length=16
        )

        unknown1: int = read_buffer.read_unsigned_short(
            logical_name="unknown1", bit_length=16
        )

        class_identifier: int = read_buffer.read_unsigned_byte(
            logical_name="class_identifier", bit_length=8
        )

        data_type: int = read_buffer.read_unsigned_byte(
            logical_name="data_type", bit_length=8
        )

        string_length: int = read_buffer.read_unsigned_byte(
            logical_name="string_length", bit_length=8
        )

        value = read_buffer.read_manual(
            read_function=lambda: StaticHelper.parse_terminated_string(
                read_buffer, string_length
            ),
            logical_name="value",
        )

        read_buffer.pop_context("UmasDatatypeReference")
        # Create the instance
        _umas_datatype_reference: UmasDatatypeReference = UmasDatatypeReference(
            data_size, unknown1, class_identifier, data_type, string_length, value
        )
        return _umas_datatype_reference

    def equals(self, o: object) -> bool:
        if self == o:
            return True

        if not isinstance(o, UmasDatatypeReference):
            return False

        that: UmasDatatypeReference = UmasDatatypeReference(o)
        return (
            (self.data_size == that.data_size)
            and (self.unknown1 == that.unknown1)
            and (self.class_identifier == that.class_identifier)
            and (self.data_type == that.data_type)
            and (self.string_length == that.string_length)
            and (self.value == that.value)
            and True
        )

    def hash_code(self) -> int:
        return hash(self)

    def __str__(self) -> str:
        # TODO:- Implement a generic python object to probably json convertor or something.
        return ""
