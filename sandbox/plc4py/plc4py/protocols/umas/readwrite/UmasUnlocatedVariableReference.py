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
    address: int
    unknown4: int
    string_length: int
    string_type: int
    value: str

    def serialize(self, write_buffer: WriteBuffer):
        write_buffer.push_context("UmasUnlocatedVariableReference")

        # Simple Field (dataType)
        write_buffer.write_unsigned_byte(
            self.data_type, bit_length=8, logical_name="dataType"
        )

        # Simple Field (address)
        write_buffer.write_unsigned_int(
            self.address, bit_length=32, logical_name="address"
        )

        # Simple Field (unknown4)
        write_buffer.write_unsigned_short(
            self.unknown4, bit_length=16, logical_name="unknown4"
        )

        # Simple Field (stringLength)
        write_buffer.write_unsigned_short(
            self.string_length, bit_length=16, logical_name="stringLength"
        )

        # Simple Field (stringType)
        write_buffer.write_unsigned_byte(
            self.string_type, bit_length=8, logical_name="stringType"
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

        # Simple field (address)
        length_in_bits += 32

        # Simple field (unknown4)
        length_in_bits += 16

        # Simple field (stringLength)
        length_in_bits += 16

        # Simple field (stringType)
        length_in_bits += 8

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

        address: int = read_buffer.read_unsigned_int(
            logical_name="address", bit_length=32
        )

        unknown4: int = read_buffer.read_unsigned_short(
            logical_name="unknown4", bit_length=16
        )

        string_length: int = read_buffer.read_unsigned_short(
            logical_name="stringLength", bit_length=16
        )

        string_type: int = read_buffer.read_unsigned_byte(
            logical_name="stringType", bit_length=8
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
                data_type, address, unknown4, string_length, string_type, value
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
            and (self.address == that.address)
            and (self.unknown4 == that.unknown4)
            and (self.string_length == that.string_length)
            and (self.string_type == that.string_type)
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
