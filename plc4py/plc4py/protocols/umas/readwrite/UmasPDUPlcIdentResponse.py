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
from plc4py.protocols.umas.readwrite.PlcMemoryBlockIdent import PlcMemoryBlockIdent
from plc4py.protocols.umas.readwrite.UmasPDUItem import UmasPDUItem
from plc4py.spi.generation.ReadBuffer import ReadBuffer
from plc4py.spi.generation.WriteBuffer import WriteBuffer
from plc4py.utils.GenericTypes import ByteOrder
from typing import Any
from typing import ClassVar
from typing import List
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
    unknown1: int
    hostname_length: int
    hostname: str
    number_of_memory_banks: int
    memory_idents: List[PlcMemoryBlockIdent]
    # Arguments.
    byte_length: int
    # Accessors for discriminator values.
    umas_function_key: ClassVar[int] = 0xFE
    umas_request_function_key: ClassVar[int] = 0x02

    def serialize_umas_pduitem_child(self, write_buffer: WriteBuffer):
        write_buffer.push_context("UmasPDUPlcIdentResponse")

        # Simple Field (range)
        write_buffer.write_unsigned_short(
            self.range, bit_length=16, logical_name="range"
        )

        # Simple Field (ident)
        write_buffer.write_unsigned_int(self.ident, bit_length=32, logical_name="ident")

        # Simple Field (model)
        write_buffer.write_unsigned_short(
            self.model, bit_length=16, logical_name="model"
        )

        # Simple Field (comVersion)
        write_buffer.write_unsigned_short(
            self.com_version, bit_length=16, logical_name="comVersion"
        )

        # Simple Field (comPatch)
        write_buffer.write_unsigned_short(
            self.com_patch, bit_length=16, logical_name="comPatch"
        )

        # Simple Field (intVersion)
        write_buffer.write_unsigned_short(
            self.int_version, bit_length=16, logical_name="intVersion"
        )

        # Simple Field (hardwareVersion)
        write_buffer.write_unsigned_short(
            self.hardware_version, bit_length=16, logical_name="hardwareVersion"
        )

        # Simple Field (crashCode)
        write_buffer.write_unsigned_int(
            self.crash_code, bit_length=32, logical_name="crashCode"
        )

        # Simple Field (unknown1)
        write_buffer.write_unsigned_short(
            self.unknown1, bit_length=16, logical_name="unknown1"
        )

        # Simple Field (hostnameLength)
        write_buffer.write_unsigned_byte(
            self.hostname_length, bit_length=8, logical_name="hostnameLength"
        )

        # Simple Field (hostname)
        write_buffer.write_str(self.hostname, bit_length=-1, logical_name="hostname")

        # Simple Field (numberOfMemoryBanks)
        write_buffer.write_unsigned_byte(
            self.number_of_memory_banks,
            bit_length=8,
            logical_name="numberOfMemoryBanks",
        )

        # Array Field (memoryIdents)
        write_buffer.write_complex_array(
            self.memory_idents, logical_name="memory_idents"
        )

        write_buffer.pop_context("UmasPDUPlcIdentResponse")

    def length_in_bytes(self) -> int:
        return int(math.ceil(float(self.length_in_bits() / 8.0)))

    def length_in_bits(self) -> int:
        length_in_bits: int = super().length_in_bits()
        _value: UmasPDUPlcIdentResponse = self

        # Simple field (range)
        length_in_bits += 16

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

        # Simple field (unknown1)
        length_in_bits += 16

        # Simple field (hostnameLength)
        length_in_bits += 8

        # Simple field (hostname)
        length_in_bits += self.hostname_length * int(8)

        # Simple field (numberOfMemoryBanks)
        length_in_bits += 8

        # Array field
        if self.memory_idents is not None:
            for element in self.memory_idents:
                length_in_bits += element.length_in_bits()

        return length_in_bits

    @staticmethod
    def static_parse_builder(
        read_buffer: ReadBuffer, umas_request_function_key: int, byte_length: int
    ):
        read_buffer.push_context("UmasPDUPlcIdentResponse")

        if isinstance(umas_request_function_key, str):
            umas_request_function_key = int(umas_request_function_key)
        if isinstance(byte_length, str):
            byte_length = int(byte_length)

        range: int = read_buffer.read_unsigned_short(
            logical_name="range",
            bit_length=16,
            byte_order=ByteOrder.LITTLE_ENDIAN,
            umas_request_function_key=umas_request_function_key,
            byte_length=byte_length,
        )

        ident: int = read_buffer.read_unsigned_int(
            logical_name="ident",
            bit_length=32,
            byte_order=ByteOrder.LITTLE_ENDIAN,
            umas_request_function_key=umas_request_function_key,
            byte_length=byte_length,
        )

        model: int = read_buffer.read_unsigned_short(
            logical_name="model",
            bit_length=16,
            byte_order=ByteOrder.LITTLE_ENDIAN,
            umas_request_function_key=umas_request_function_key,
            byte_length=byte_length,
        )

        com_version: int = read_buffer.read_unsigned_short(
            logical_name="com_version",
            bit_length=16,
            byte_order=ByteOrder.LITTLE_ENDIAN,
            umas_request_function_key=umas_request_function_key,
            byte_length=byte_length,
        )

        com_patch: int = read_buffer.read_unsigned_short(
            logical_name="com_patch",
            bit_length=16,
            byte_order=ByteOrder.LITTLE_ENDIAN,
            umas_request_function_key=umas_request_function_key,
            byte_length=byte_length,
        )

        int_version: int = read_buffer.read_unsigned_short(
            logical_name="int_version",
            bit_length=16,
            byte_order=ByteOrder.LITTLE_ENDIAN,
            umas_request_function_key=umas_request_function_key,
            byte_length=byte_length,
        )

        hardware_version: int = read_buffer.read_unsigned_short(
            logical_name="hardware_version",
            bit_length=16,
            byte_order=ByteOrder.LITTLE_ENDIAN,
            umas_request_function_key=umas_request_function_key,
            byte_length=byte_length,
        )

        crash_code: int = read_buffer.read_unsigned_int(
            logical_name="crash_code",
            bit_length=32,
            byte_order=ByteOrder.LITTLE_ENDIAN,
            umas_request_function_key=umas_request_function_key,
            byte_length=byte_length,
        )

        unknown1: int = read_buffer.read_unsigned_short(
            logical_name="unknown1",
            bit_length=16,
            byte_order=ByteOrder.LITTLE_ENDIAN,
            umas_request_function_key=umas_request_function_key,
            byte_length=byte_length,
        )

        hostname_length: int = read_buffer.read_unsigned_byte(
            logical_name="hostname_length",
            bit_length=8,
            byte_order=ByteOrder.LITTLE_ENDIAN,
            umas_request_function_key=umas_request_function_key,
            byte_length=byte_length,
        )

        hostname: str = read_buffer.read_str(
            logical_name="hostname",
            bit_length=hostname_length * int(8),
            byte_order=ByteOrder.LITTLE_ENDIAN,
            umas_request_function_key=umas_request_function_key,
            byte_length=byte_length,
        )

        number_of_memory_banks: int = read_buffer.read_unsigned_byte(
            logical_name="number_of_memory_banks",
            bit_length=8,
            byte_order=ByteOrder.LITTLE_ENDIAN,
            umas_request_function_key=umas_request_function_key,
            byte_length=byte_length,
        )

        memory_idents: List[Any] = read_buffer.read_array_field(
            logical_name="memoryIdents",
            read_function=PlcMemoryBlockIdent.static_parse,
            count=number_of_memory_banks,
            byte_order=ByteOrder.LITTLE_ENDIAN,
            umas_request_function_key=umas_request_function_key,
            byte_length=byte_length,
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
            unknown1,
            hostname_length,
            hostname,
            number_of_memory_banks,
            memory_idents,
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
            and (self.unknown1 == that.unknown1)
            and (self.hostname_length == that.hostname_length)
            and (self.hostname == that.hostname)
            and (self.number_of_memory_banks == that.number_of_memory_banks)
            and (self.memory_idents == that.memory_idents)
            and super().equals(that)
            and True
        )

    def hash_code(self) -> int:
        return hash(self)

    def __str__(self) -> str:
        # TODO:- Implement a generic python object to probably json convertor or something.
        return ""


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
    unknown1: int
    hostname_length: int
    hostname: str
    number_of_memory_banks: int
    memory_idents: List[PlcMemoryBlockIdent]

    def build(self, byte_length: int, pairing_key) -> UmasPDUPlcIdentResponse:
        umas_pduplc_ident_response: UmasPDUPlcIdentResponse = UmasPDUPlcIdentResponse(
            byte_length,
            pairing_key,
            self.range,
            self.ident,
            self.model,
            self.com_version,
            self.com_patch,
            self.int_version,
            self.hardware_version,
            self.crash_code,
            self.unknown1,
            self.hostname_length,
            self.hostname,
            self.number_of_memory_banks,
            self.memory_idents,
        )
        return umas_pduplc_ident_response
