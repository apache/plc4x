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
    range: int
    hardware_id: int
    block_no: int
    BLANK: int = 0x0000
    # Accessors for discriminator values.
    umas_function_key: ClassVar[int] = 0x26
    umas_request_function_key: ClassVar[int] = 0

    def serialize_umas_pdu_item_child(self, write_buffer: WriteBuffer):
        write_buffer.push_context("UmasPDUReadUnlocatedVariableNamesRequest")

        # Simple Field (range)
        write_buffer.write_unsigned_short(
            self.range, bit_length=16, logical_name="range"
        )

        # Simple Field (hardwareId)
        write_buffer.write_unsigned_long(
            self.hardware_id, bit_length=40, logical_name="hardwareId"
        )

        # Simple Field (blockNo)
        write_buffer.write_unsigned_short(
            self.block_no, bit_length=16, logical_name="blockNo"
        )

        # Const Field (blank)
        write_buffer.write_unsigned_int(self.BLANK, logical_name="blank")

        write_buffer.pop_context("UmasPDUReadUnlocatedVariableNamesRequest")

    def length_in_bytes(self) -> int:
        return int(math.ceil(float(self.length_in_bits() / 8.0)))

    def length_in_bits(self) -> int:
        length_in_bits: int = super().length_in_bits()
        _value: UmasPDUReadUnlocatedVariableNamesRequest = self

        # Simple field (range)
        length_in_bits += 16

        # Simple field (hardwareId)
        length_in_bits += 40

        # Simple field (blockNo)
        length_in_bits += 16

        # Const Field (blank)
        length_in_bits += 32

        return length_in_bits

    @staticmethod
    def static_parse_builder(read_buffer: ReadBuffer, umas_request_function_key: int):
        read_buffer.push_context("UmasPDUReadUnlocatedVariableNamesRequest")

        range: int = read_buffer.read_unsigned_short(
            logical_name="range",
            bit_length=16,
            byte_order=ByteOrder.LITTLE_ENDIAN,
            umas_request_function_key=umas_request_function_key,
        )

        hardware_id: int = read_buffer.read_unsigned_long(
            logical_name="hardwareId",
            bit_length=40,
            byte_order=ByteOrder.LITTLE_ENDIAN,
            umas_request_function_key=umas_request_function_key,
        )

        block_no: int = read_buffer.read_unsigned_short(
            logical_name="blockNo",
            bit_length=16,
            byte_order=ByteOrder.LITTLE_ENDIAN,
            umas_request_function_key=umas_request_function_key,
        )

        BLANK: int = read_buffer.read_unsigned_int(
            logical_name="blank",
            byte_order=ByteOrder.LITTLE_ENDIAN,
            umas_request_function_key=umas_request_function_key,
        )

        read_buffer.pop_context("UmasPDUReadUnlocatedVariableNamesRequest")
        # Create the instance
        return UmasPDUReadUnlocatedVariableNamesRequestBuilder(
            range, hardware_id, block_no
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
            (self.range == that.range)
            and (self.hardware_id == that.hardware_id)
            and (self.block_no == that.block_no)
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
class UmasPDUReadUnlocatedVariableNamesRequestBuilder:
    range: int
    hardware_id: int
    block_no: int

    def build(self, pairing_key) -> UmasPDUReadUnlocatedVariableNamesRequest:
        umas_pdu_read_unlocated_variable_names_request: UmasPDUReadUnlocatedVariableNamesRequest = UmasPDUReadUnlocatedVariableNamesRequest(
            pairing_key, self.range, self.hardware_id, self.block_no
        )
        return umas_pdu_read_unlocated_variable_names_request
