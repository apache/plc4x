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

from abc import ABC
from abc import abstractmethod
from plc4py.api.exceptions.exceptions import ParseException
from plc4py.api.exceptions.exceptions import PlcRuntimeException
from plc4py.api.exceptions.exceptions import SerializationException
from plc4py.api.messages.PlcMessage import PlcMessage
from plc4py.spi.generation.ReadBuffer import ReadBuffer
from plc4py.spi.generation.WriteBuffer import WriteBuffer
import math


@dataclass
class ModbusPDU(ABC, PlcMessage):
    # Arguments.
    umas_request_function_key: int
    byte_length: int

    # Abstract accessors for discriminator values.
    @property
    def error_flag(self) -> bool:
        pass

    @property
    def function_flag(self) -> int:
        pass

    @abstractmethod
    def serialize_modbus_pdu_child(self, write_buffer: WriteBuffer) -> None:
        pass

    def serialize(self, write_buffer: WriteBuffer):
        write_buffer.push_context("ModbusPDU")

        # Discriminator Field (errorFlag) (Used as input to a switch field)
        write_buffer.write_bit(self.error_flag, logical_name="error_flag", bit_length=1)

        # Discriminator Field (functionFlag) (Used as input to a switch field)
        write_buffer.write_unsigned_byte(
            self.function_flag, logical_name="function_flag", bit_length=7
        )

        # Switch field (Serialize the sub-type)
        self.serialize_modbus_pdu_child(write_buffer)

        write_buffer.pop_context("ModbusPDU")

    def length_in_bytes(self) -> int:
        return int(math.ceil(float(self.length_in_bits() / 8.0)))

    def length_in_bits(self) -> int:
        length_in_bits: int = 0
        _value: ModbusPDU = self

        # Discriminator Field (errorFlag)
        length_in_bits += 1

        # Discriminator Field (functionFlag)
        length_in_bits += 7

        # Length of subtype elements will be added by sub-type...

        return length_in_bits

    @staticmethod
    def static_parse(read_buffer: ReadBuffer, **kwargs):

        if kwargs is None:
            raise PlcRuntimeException(
                "Wrong number of arguments, expected 2, but got None"
            )

        umas_request_function_key: int = 0
        if isinstance(kwargs.get("umas_request_function_key"), int):
            umas_request_function_key = int(kwargs.get("umas_request_function_key"))
        elif isinstance(kwargs.get("umas_request_function_key"), str):
            umas_request_function_key = int(
                str(kwargs.get("umas_request_function_key"))
            )
        else:
            raise PlcRuntimeException(
                "Argument 0 expected to be of type int or a string which is parseable but was "
                + kwargs.get("umas_request_function_key").getClass().getName()
            )

        byte_length: int = 0
        if isinstance(kwargs.get("byte_length"), int):
            byte_length = int(kwargs.get("byte_length"))
        elif isinstance(kwargs.get("byte_length"), str):
            byte_length = int(str(kwargs.get("byte_length")))
        else:
            raise PlcRuntimeException(
                "Argument 1 expected to be of type int or a string which is parseable but was "
                + kwargs.get("byte_length").getClass().getName()
            )

        return ModbusPDU.static_parse_context(
            read_buffer, umas_request_function_key, byte_length
        )

    @staticmethod
    def static_parse_context(
        read_buffer: ReadBuffer, umas_request_function_key: int, byte_length: int
    ):
        read_buffer.push_context("ModbusPDU")

        if isinstance(umas_request_function_key, str):
            umas_request_function_key = int(umas_request_function_key)
        if isinstance(byte_length, str):
            byte_length = int(byte_length)

        error_flag: bool = read_buffer.read_bit(
            logical_name="error_flag",
            bit_length=1,
            umas_request_function_key=umas_request_function_key,
            byte_length=byte_length,
        )

        function_flag: int = read_buffer.read_unsigned_byte(
            logical_name="function_flag",
            bit_length=7,
            umas_request_function_key=umas_request_function_key,
            byte_length=byte_length,
        )

        # Switch Field (Depending on the discriminator values, passes the instantiation to a sub-type)
        builder: ModbusPDUBuilder = None
        from plc4py.protocols.umas.readwrite.ModbusPDUError import ModbusPDUError

        if error_flag == bool(True):

            builder = ModbusPDUError.static_parse_builder(
                read_buffer, umas_request_function_key, byte_length
            )
        from plc4py.protocols.umas.readwrite.UmasPDU import UmasPDU

        if error_flag == bool(False) and function_flag == int(0x5A):

            builder = UmasPDU.static_parse_builder(
                read_buffer, umas_request_function_key, byte_length
            )
        if builder is None:
            raise ParseException(
                "Unsupported case for discriminated type"
                + " parameters ["
                + "errorFlag="
                + str(error_flag)
                + " "
                + "functionFlag="
                + str(function_flag)
                + "]"
            )

        read_buffer.pop_context("ModbusPDU")
        # Create the instance
        _modbus_pdu: ModbusPDU = builder.build(umas_request_function_key, byte_length)
        return _modbus_pdu

    def equals(self, o: object) -> bool:
        if self == o:
            return True

        if not isinstance(o, ModbusPDU):
            return False

        that: ModbusPDU = ModbusPDU(o)
        return True

    def hash_code(self) -> int:
        return hash(self)

    def __str__(self) -> str:
        # TODO:- Implement a generic python object to probably json convertor or something.
        return ""


@dataclass
class ModbusPDUBuilder:
    def build(self, umas_request_function_key: int, byte_length: int) -> ModbusPDU:
        pass
