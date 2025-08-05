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
from plc4py.protocols.umas.readwrite.ModbusErrorCode import ModbusErrorCode
from plc4py.protocols.umas.readwrite.ModbusPDU import ModbusPDU
from plc4py.spi.generation.ReadBuffer import ReadBuffer
from plc4py.spi.generation.WriteBuffer import WriteBuffer
from typing import ClassVar
import math


@dataclass
class ModbusPDUError(ModbusPDU):
    exception_code: ModbusErrorCode
    # Arguments.
    umas_request_function_key: int
    byte_length: int
    # Accessors for discriminator values.
    error_flag: ClassVar[bool] = True
    function_flag: ClassVar[int] = 0

    def serialize_modbus_pdu_child(self, write_buffer: WriteBuffer):
        write_buffer.push_context("ModbusPDUError")

        # Simple Field (exceptionCode)
        write_buffer.write_unsigned_byte(
            self.exception_code, logical_name="exceptionCode"
        )

        write_buffer.pop_context("ModbusPDUError")

    def length_in_bytes(self) -> int:
        return int(math.ceil(float(self.length_in_bits() / 8.0)))

    def length_in_bits(self) -> int:
        length_in_bits: int = super().length_in_bits()
        _value: ModbusPDUError = self

        # Simple field (exceptionCode)
        length_in_bits += 8

        return length_in_bits

    @staticmethod
    def static_parse_builder(
        read_buffer: ReadBuffer, umas_request_function_key: int, byte_length: int
    ):
        read_buffer.push_context("ModbusPDUError")

        if isinstance(umas_request_function_key, str):
            umas_request_function_key = int(umas_request_function_key)
        if isinstance(byte_length, str):
            byte_length = int(byte_length)

        exception_code: ModbusErrorCode = read_buffer.read_enum(
            read_function=ModbusErrorCode,
            bit_length=8,
            logical_name="exception_code",
            umas_request_function_key=umas_request_function_key,
            byte_length=byte_length,
        )

        read_buffer.pop_context("ModbusPDUError")
        # Create the instance
        return ModbusPDUErrorBuilder(exception_code)

    def equals(self, o: object) -> bool:
        if self == o:
            return True

        if not isinstance(o, ModbusPDUError):
            return False

        that: ModbusPDUError = ModbusPDUError(o)
        return (
            (self.exception_code == that.exception_code)
            and super().equals(that)
            and True
        )

    def hash_code(self) -> int:
        return hash(self)

    def __str__(self) -> str:
        # TODO:- Implement a generic python object to probably json convertor or something.
        return ""


@dataclass
class ModbusPDUErrorBuilder:
    exception_code: ModbusErrorCode

    def build(
        self,
        umas_request_function_key: int,
        byte_length: int,
    ) -> ModbusPDUError:
        modbus_pduerror: ModbusPDUError = ModbusPDUError(
            umas_request_function_key, byte_length, self.exception_code
        )
        return modbus_pduerror
