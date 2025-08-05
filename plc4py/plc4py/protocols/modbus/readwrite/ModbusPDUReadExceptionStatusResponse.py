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
from plc4py.protocols.modbus.readwrite.ModbusPDU import ModbusPDU
from plc4py.spi.generation.ReadBuffer import ReadBuffer
from plc4py.spi.generation.WriteBuffer import WriteBuffer
from plc4py.utils.ConnectionStringHandling import strtobool
from typing import ClassVar
import math


@dataclass
class ModbusPDUReadExceptionStatusResponse(ModbusPDU):
    value: int
    # Accessors for discriminator values.
    error_flag: ClassVar[bool] = False
    function_flag: ClassVar[int] = 0x07
    response: ClassVar[bool] = True

    def serialize_modbus_pdu_child(self, write_buffer: WriteBuffer):
        write_buffer.push_context("ModbusPDUReadExceptionStatusResponse")

        # Simple Field (value)
        write_buffer.write_unsigned_byte(self.value, bit_length=8, logical_name="value")

        write_buffer.pop_context("ModbusPDUReadExceptionStatusResponse")

    def length_in_bytes(self) -> int:
        return int(math.ceil(float(self.length_in_bits() / 8.0)))

    def length_in_bits(self) -> int:
        length_in_bits: int = super().length_in_bits()
        _value: ModbusPDUReadExceptionStatusResponse = self

        # Simple field (value)
        length_in_bits += 8

        return length_in_bits

    @staticmethod
    def static_parse_builder(read_buffer: ReadBuffer, response: bool):
        read_buffer.push_context("ModbusPDUReadExceptionStatusResponse")

        if isinstance(response, str):
            response = bool(strtobool(response))

        value: int = read_buffer.read_unsigned_byte(
            logical_name="value", bit_length=8, response=response
        )

        read_buffer.pop_context("ModbusPDUReadExceptionStatusResponse")
        # Create the instance
        return ModbusPDUReadExceptionStatusResponseBuilder(value)

    def equals(self, o: object) -> bool:
        if self == o:
            return True

        if not isinstance(o, ModbusPDUReadExceptionStatusResponse):
            return False

        that: ModbusPDUReadExceptionStatusResponse = (
            ModbusPDUReadExceptionStatusResponse(o)
        )
        return (self.value == that.value) and super().equals(that) and True

    def hash_code(self) -> int:
        return hash(self)

    def __str__(self) -> str:
        # TODO:- Implement a generic python object to probably json convertor or something.
        return ""


@dataclass
class ModbusPDUReadExceptionStatusResponseBuilder:
    value: int

    def build(
        self,
    ) -> ModbusPDUReadExceptionStatusResponse:
        modbus_pduread_exception_status_response: (
            ModbusPDUReadExceptionStatusResponse
        ) = ModbusPDUReadExceptionStatusResponse(self.value)
        return modbus_pduread_exception_status_response
