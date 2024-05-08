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
from plc4py.protocols.umas.readwrite.VariableWriteRequestReference import (
    VariableWriteRequestReference,
)
from plc4py.spi.generation.ReadBuffer import ReadBuffer
from plc4py.spi.generation.WriteBuffer import WriteBuffer
from plc4py.utils.GenericTypes import ByteOrder
from typing import Any
from typing import ClassVar
from typing import List
import math


@dataclass
class UmasPDUWriteVariableRequest(UmasPDUItem):
    crc: int
    variable_count: int
    variables: List[VariableWriteRequestReference]
    # Arguments.
    byte_length: int
    # Accessors for discriminator values.
    umas_function_key: ClassVar[int] = 0x23
    umas_request_function_key: ClassVar[int] = 0

    def serialize_umas_pdu_item_child(self, write_buffer: WriteBuffer):
        write_buffer.push_context("UmasPDUWriteVariableRequest")

        # Simple Field (crc)
        write_buffer.write_unsigned_int(self.crc, bit_length=32, logical_name="crc")

        # Simple Field (variableCount)
        write_buffer.write_unsigned_byte(
            self.variable_count, bit_length=8, logical_name="variableCount"
        )

        # Array Field (variables)
        write_buffer.write_complex_array(self.variables, logical_name="variables")

        write_buffer.pop_context("UmasPDUWriteVariableRequest")

    def length_in_bytes(self) -> int:
        return int(math.ceil(float(self.length_in_bits() / 8.0)))

    def length_in_bits(self) -> int:
        length_in_bits: int = super().length_in_bits()
        _value: UmasPDUWriteVariableRequest = self

        # Simple field (crc)
        length_in_bits += 32

        # Simple field (variableCount)
        length_in_bits += 8

        # Array field
        if self.variables is not None:
            for element in self.variables:
                length_in_bits += element.length_in_bits()

        return length_in_bits

    @staticmethod
    def static_parse_builder(
        read_buffer: ReadBuffer, umas_request_function_key: int, byte_length: int
    ):
        read_buffer.push_context("UmasPDUWriteVariableRequest")

        crc: int = read_buffer.read_unsigned_int(
            logical_name="crc",
            bit_length=32,
            byte_order=ByteOrder.LITTLE_ENDIAN,
            umas_request_function_key=umas_request_function_key,
            byte_length=byte_length,
        )

        variable_count: int = read_buffer.read_unsigned_byte(
            logical_name="variable_count",
            bit_length=8,
            byte_order=ByteOrder.LITTLE_ENDIAN,
            umas_request_function_key=umas_request_function_key,
            byte_length=byte_length,
        )

        variables: List[Any] = read_buffer.read_array_field(
            logical_name="variables",
            read_function=VariableWriteRequestReference.static_parse,
            count=variable_count,
            byte_order=ByteOrder.LITTLE_ENDIAN,
            umas_request_function_key=umas_request_function_key,
            byte_length=byte_length,
        )

        read_buffer.pop_context("UmasPDUWriteVariableRequest")
        # Create the instance
        return UmasPDUWriteVariableRequestBuilder(crc, variable_count, variables)

    def equals(self, o: object) -> bool:
        if self == o:
            return True

        if not isinstance(o, UmasPDUWriteVariableRequest):
            return False

        that: UmasPDUWriteVariableRequest = UmasPDUWriteVariableRequest(o)
        return (
            (self.crc == that.crc)
            and (self.variable_count == that.variable_count)
            and (self.variables == that.variables)
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
class UmasPDUWriteVariableRequestBuilder:
    crc: int
    variable_count: int
    variables: List[VariableWriteRequestReference]

    def build(self, byte_length: int, pairing_key) -> UmasPDUWriteVariableRequest:
        umas_pdu_write_variable_request: UmasPDUWriteVariableRequest = (
            UmasPDUWriteVariableRequest(
                byte_length, pairing_key, self.crc, self.variable_count, self.variables
            )
        )
        return umas_pdu_write_variable_request
