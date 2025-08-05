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
from plc4py.utils.GenericTypes import ByteOrder
from typing import ClassVar
import math


@dataclass
class UmasPDUReadUnlocatedVariableNamesRequest(UmasPDUItem):
    record_type: int
    index: int
    hardware_id: int
    block_no: int
    offset: int
    # Arguments.
    byte_length: int
    BLANK: int = 0x00
    # Accessors for discriminator values.
    umas_function_key: ClassVar[int] = 0x26
    umas_request_function_key: ClassVar[int] = 0

    def serialize_umas_pduitem_child(self, write_buffer: WriteBuffer):
        write_buffer.push_context("UmasPDUReadUnlocatedVariableNamesRequest")

        # Simple Field (recordType)
        write_buffer.write_unsigned_short(
            self.record_type, bit_length=16, logical_name="recordType"
        )

        # Simple Field (index)
        write_buffer.write_unsigned_byte(self.index, bit_length=8, logical_name="index")

        # Simple Field (hardwareId)
        write_buffer.write_unsigned_int(
            self.hardware_id, bit_length=32, logical_name="hardwareId"
        )

        # Simple Field (blockNo)
        write_buffer.write_unsigned_short(
            self.block_no, bit_length=16, logical_name="blockNo"
        )

        # Simple Field (offset)
        write_buffer.write_unsigned_short(
            self.offset, bit_length=16, logical_name="offset"
        )

        # Const Field (blank)
        write_buffer.write_unsigned_short(self.BLANK, logical_name="blank")

        write_buffer.pop_context("UmasPDUReadUnlocatedVariableNamesRequest")

    def length_in_bytes(self) -> int:
        return int(math.ceil(float(self.length_in_bits() / 8.0)))

    def length_in_bits(self) -> int:
        length_in_bits: int = super().length_in_bits()
        _value: UmasPDUReadUnlocatedVariableNamesRequest = self

        # Simple field (recordType)
        length_in_bits += 16

        # Simple field (index)
        length_in_bits += 8

        # Simple field (hardwareId)
        length_in_bits += 32

        # Simple field (blockNo)
        length_in_bits += 16

        # Simple field (offset)
        length_in_bits += 16

        # Const Field (blank)
        length_in_bits += 16

        return length_in_bits

    @staticmethod
    def static_parse_builder(
        read_buffer: ReadBuffer, umas_request_function_key: int, byte_length: int
    ):
        read_buffer.push_context("UmasPDUReadUnlocatedVariableNamesRequest")

        if isinstance(umas_request_function_key, str):
            umas_request_function_key = int(umas_request_function_key)
        if isinstance(byte_length, str):
            byte_length = int(byte_length)

        record_type: int = read_buffer.read_unsigned_short(
            logical_name="record_type",
            bit_length=16,
            byte_order=ByteOrder.LITTLE_ENDIAN,
            umas_request_function_key=umas_request_function_key,
            byte_length=byte_length,
        )

        index: int = read_buffer.read_unsigned_byte(
            logical_name="index",
            bit_length=8,
            byte_order=ByteOrder.LITTLE_ENDIAN,
            umas_request_function_key=umas_request_function_key,
            byte_length=byte_length,
        )

        hardware_id: int = read_buffer.read_unsigned_int(
            logical_name="hardware_id",
            bit_length=32,
            byte_order=ByteOrder.LITTLE_ENDIAN,
            umas_request_function_key=umas_request_function_key,
            byte_length=byte_length,
        )

        block_no: int = read_buffer.read_unsigned_short(
            logical_name="block_no",
            bit_length=16,
            byte_order=ByteOrder.LITTLE_ENDIAN,
            umas_request_function_key=umas_request_function_key,
            byte_length=byte_length,
        )

        offset: int = read_buffer.read_unsigned_short(
            logical_name="offset",
            bit_length=16,
            byte_order=ByteOrder.LITTLE_ENDIAN,
            umas_request_function_key=umas_request_function_key,
            byte_length=byte_length,
        )

        BLANK: int = read_buffer.read_unsigned_short(
            logical_name="blank",
            byte_order=ByteOrder.LITTLE_ENDIAN,
            umas_request_function_key=umas_request_function_key,
            byte_length=byte_length,
        )

        read_buffer.pop_context("UmasPDUReadUnlocatedVariableNamesRequest")
        # Create the instance
        return UmasPDUReadUnlocatedVariableNamesRequestBuilder(
            record_type, index, hardware_id, block_no, offset
        )

    def equals(self, o: object) -> bool:
        if self == o:
            return True

        if not isinstance(o, UmasPDUReadUnlocatedVariableNamesRequest):
            return False

        that: UmasPDUReadUnlocatedVariableNamesRequest = (
            UmasPDUReadUnlocatedVariableNamesRequest(o)
        )
        return (
            (self.record_type == that.record_type)
            and (self.index == that.index)
            and (self.hardware_id == that.hardware_id)
            and (self.block_no == that.block_no)
            and (self.offset == that.offset)
            and super().equals(that)
            and True
        )

    def hash_code(self) -> int:
        return hash(self)

    def __str__(self) -> str:
        # TODO:- Implement a generic python object to probably json convertor or something.
        return ""


@dataclass
class UmasPDUReadUnlocatedVariableNamesRequestBuilder:
    record_type: int
    index: int
    hardware_id: int
    block_no: int
    offset: int

    def build(
        self, byte_length: int, pairing_key
    ) -> UmasPDUReadUnlocatedVariableNamesRequest:
        umas_pduread_unlocated_variable_names_request: (
            UmasPDUReadUnlocatedVariableNamesRequest
        ) = UmasPDUReadUnlocatedVariableNamesRequest(
            byte_length,
            pairing_key,
            self.record_type,
            self.index,
            self.hardware_id,
            self.block_no,
            self.offset,
        )
        return umas_pduread_unlocated_variable_names_request
