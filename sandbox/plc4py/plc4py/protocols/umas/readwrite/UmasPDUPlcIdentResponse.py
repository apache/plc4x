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
from plc4py.spi.generation.ReadBuffer import ReadBuffer
from plc4py.spi.generation.WriteBuffer import WriteBuffer
from typing import ClassVar
import math


@dataclass
class UmasPDUPlcIdentResponse(UmasPDUItem):
    range: int
    ident: int
    model: int
    com_version: int
    com_patch: int
    int_version: int
    hardware_version: int
    crash_code: int
    string_length: int
    string_value: str
    # Accessors for discriminator values.
    umas_function_key: ClassVar[int] = 0xFE
    umas_request_function_key: ClassVar[int] = 0x02

    def serialize_umas_pdu_item_child(self, write_buffer: WriteBuffer):
        write_buffer.push_context("UmasPDUPlcIdentResponse")

        # Simple Field (range)
        write_buffer.write_unsigned_byte(self.range, logical_name="range")

        # Simple Field (ident)
        write_buffer.write_unsigned_int(self.ident, logical_name="ident")

        # Simple Field (model)
        write_buffer.write_unsigned_short(self.model, logical_name="model")

        # Simple Field (comVersion)
        write_buffer.write_unsigned_short(self.com_version, logical_name="comVersion")

        # Simple Field (comPatch)
        write_buffer.write_unsigned_short(self.com_patch, logical_name="comPatch")

        # Simple Field (intVersion)
        write_buffer.write_unsigned_short(self.int_version, logical_name="intVersion")

        # Simple Field (hardwareVersion)
        write_buffer.write_unsigned_short(
            self.hardware_version, logical_name="hardwareVersion"
        )

        # Simple Field (crashCode)
        write_buffer.write_unsigned_int(self.crash_code, logical_name="crashCode")

        # Simple Field (stringLength)
        write_buffer.write_unsigned_int(self.string_length, logical_name="stringLength")

        # Simple Field (stringValue)
        write_buffer.write_str(self.string_value, logical_name="stringValue")

        write_buffer.pop_context("UmasPDUPlcIdentResponse")

    def length_in_bytes(self) -> int:
        return int(math.ceil(float(self.length_in_bits() / 8.0)))

    def length_in_bits(self) -> int:
        length_in_bits: int = super().length_in_bits()
        _value: UmasPDUPlcIdentResponse = self

        # Simple field (range)
        length_in_bits += 8

        # Simple field (ident)
        length_in_bits += 32

        # Simple field (model)
        length_in_bits += 16

        # Simple field (comVersion)
        length_in_bits += 16

        # Simple field (comPatch)
        length_in_bits += 16

        # Simple field (intVersion)
        length_in_bits += 16

        # Simple field (hardwareVersion)
        length_in_bits += 16

        # Simple field (crashCode)
        length_in_bits += 32

        # Simple field (stringLength)
        length_in_bits += 32

        # Simple field (stringValue)
        length_in_bits += self.string_length * int(8)

        return length_in_bits

    @staticmethod
    def static_parse_builder(read_buffer: ReadBuffer, umas_request_function_key: int):
        read_buffer.push_context("UmasPDUPlcIdentResponse")

        range: int = read_buffer.read_unsigned_byte(
            logical_name="range",
            bit_length=8,
            umas_request_function_key=umas_request_function_key,
        )

        ident: int = read_buffer.read_unsigned_int(
            logical_name="ident",
            bit_length=32,
            umas_request_function_key=umas_request_function_key,
        )

        model: int = read_buffer.read_unsigned_short(
            logical_name="model",
            bit_length=16,
            umas_request_function_key=umas_request_function_key,
        )

        com_version: int = read_buffer.read_unsigned_short(
            logical_name="comVersion",
            bit_length=16,
            umas_request_function_key=umas_request_function_key,
        )

        com_patch: int = read_buffer.read_unsigned_short(
            logical_name="comPatch",
            bit_length=16,
            umas_request_function_key=umas_request_function_key,
        )

        int_version: int = read_buffer.read_unsigned_short(
            logical_name="intVersion",
            bit_length=16,
            umas_request_function_key=umas_request_function_key,
        )

        hardware_version: int = read_buffer.read_unsigned_short(
            logical_name="hardwareVersion",
            bit_length=16,
            umas_request_function_key=umas_request_function_key,
        )

        crash_code: int = read_buffer.read_unsigned_int(
            logical_name="crashCode",
            bit_length=32,
            umas_request_function_key=umas_request_function_key,
        )

        string_length: int = read_buffer.read_unsigned_int(
            logical_name="stringLength",
            bit_length=32,
            umas_request_function_key=umas_request_function_key,
        )

        string_value: str = read_buffer.read_str(
            logical_name="stringValue",
            bit_length=string_length * int(8),
            umas_request_function_key=umas_request_function_key,
        )

        read_buffer.pop_context("UmasPDUPlcIdentResponse")
        # Create the instance
        return UmasPDUPlcIdentResponseBuilder(
            range,
            ident,
            model,
            com_version,
            com_patch,
            int_version,
            hardware_version,
            crash_code,
            string_length,
            string_value,
        )

    def equals(self, o: object) -> bool:
        if self == o:
            return True

        if not isinstance(o, UmasPDUPlcIdentResponse):
            return False

        that: UmasPDUPlcIdentResponse = UmasPDUPlcIdentResponse(o)
        return (
            (self.range == that.range)
            and (self.ident == that.ident)
            and (self.model == that.model)
            and (self.com_version == that.com_version)
            and (self.com_patch == that.com_patch)
            and (self.int_version == that.int_version)
            and (self.hardware_version == that.hardware_version)
            and (self.crash_code == that.crash_code)
            and (self.string_length == that.string_length)
            and (self.string_value == that.string_value)
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
class UmasPDUPlcIdentResponseBuilder:
    range: int
    ident: int
    model: int
    com_version: int
    com_patch: int
    int_version: int
    hardware_version: int
    crash_code: int
    string_length: int
    string_value: str

    def build(self, pairing_key) -> UmasPDUPlcIdentResponse:
        umas_pdu_plc_ident_response: UmasPDUPlcIdentResponse = UmasPDUPlcIdentResponse(
            pairing_key,
            self.range,
            self.ident,
            self.model,
            self.com_version,
            self.com_patch,
            self.int_version,
            self.hardware_version,
            self.crash_code,
            self.string_length,
            self.string_value,
        )
        return umas_pdu_plc_ident_response