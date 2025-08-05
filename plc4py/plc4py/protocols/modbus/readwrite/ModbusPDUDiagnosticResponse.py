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
class ModbusPDUDiagnosticResponse(ModbusPDU):
    sub_function: int
    data: int
    # Accessors for discriminator values.
    error_flag: ClassVar[bool] = False
    function_flag: ClassVar[int] = 0x08
    response: ClassVar[bool] = True

    def serialize_modbus_pdu_child(self, write_buffer: WriteBuffer):
        write_buffer.push_context("ModbusPDUDiagnosticResponse")

        # Simple Field (subFunction)
        write_buffer.write_unsigned_short(
            self.sub_function, bit_length=16, logical_name="subFunction"
        )

        # Simple Field (data)
        write_buffer.write_unsigned_short(self.data, bit_length=16, logical_name="data")

        write_buffer.pop_context("ModbusPDUDiagnosticResponse")

    def length_in_bytes(self) -> int:
        return int(math.ceil(float(self.length_in_bits() / 8.0)))

    def length_in_bits(self) -> int:
        length_in_bits: int = super().length_in_bits()
        _value: ModbusPDUDiagnosticResponse = self

        # Simple field (subFunction)
        length_in_bits += 16

        # Simple field (data)
        length_in_bits += 16

        return length_in_bits

    @staticmethod
    def static_parse_builder(read_buffer: ReadBuffer, response: bool):
        read_buffer.push_context("ModbusPDUDiagnosticResponse")

        if isinstance(response, str):
            response = bool(strtobool(response))

        sub_function: int = read_buffer.read_unsigned_short(
            logical_name="sub_function", bit_length=16, response=response
        )

        data: int = read_buffer.read_unsigned_short(
            logical_name="data", bit_length=16, response=response
        )

        read_buffer.pop_context("ModbusPDUDiagnosticResponse")
        # Create the instance
        return ModbusPDUDiagnosticResponseBuilder(sub_function, data)

    def equals(self, o: object) -> bool:
        if self == o:
            return True

        if not isinstance(o, ModbusPDUDiagnosticResponse):
            return False

        that: ModbusPDUDiagnosticResponse = ModbusPDUDiagnosticResponse(o)
        return (
            (self.sub_function == that.sub_function)
            and (self.data == that.data)
            and super().equals(that)
            and True
        )

    def hash_code(self) -> int:
        return hash(self)

    def __str__(self) -> str:
        # TODO:- Implement a generic python object to probably json convertor or something.
        return ""


@dataclass
class ModbusPDUDiagnosticResponseBuilder:
    sub_function: int
    data: int

    def build(
        self,
    ) -> ModbusPDUDiagnosticResponse:
        modbus_pdudiagnostic_response: ModbusPDUDiagnosticResponse = (
            ModbusPDUDiagnosticResponse(self.sub_function, self.data)
        )
        return modbus_pdudiagnostic_response
