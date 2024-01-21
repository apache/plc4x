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
from typing import Any
from typing import ClassVar
from typing import List
import math


@dataclass
class UmasPDUPlcStatusResponse(UmasPDUItem):
    not_used: int
    number_of_blocks: int
    blocks: List[int]
    # Accessors for discriminator values.
    umas_function_key: ClassVar[int] = 0xFE
    umas_request_function_key: ClassVar[int] = 0x04

    def serialize_umas_pdu_item_child(self, write_buffer: WriteBuffer):
        write_buffer.push_context("UmasPDUPlcStatusResponse")

        # Simple Field (notUsed)
        write_buffer.write_unsigned_int(self.not_used, logical_name="notUsed")

        # Simple Field (numberOfBlocks)
        write_buffer.write_unsigned_byte(
            self.number_of_blocks, logical_name="numberOfBlocks"
        )

        # Array Field (blocks)
        write_buffer.write_simple_array(
            self.blocks, write_unsigned_int, logical_name="blocks"
        )

        write_buffer.pop_context("UmasPDUPlcStatusResponse")

    def length_in_bytes(self) -> int:
        return int(math.ceil(float(self.length_in_bits() / 8.0)))

    def length_in_bits(self) -> int:
        length_in_bits: int = super().length_in_bits()
        _value: UmasPDUPlcStatusResponse = self

        # Simple field (notUsed)
        length_in_bits += 24

        # Simple field (numberOfBlocks)
        length_in_bits += 8

        # Array field
        if self.blocks is not None:
            length_in_bits += 32 * len(self.blocks)

        return length_in_bits

    @staticmethod
    def static_parse_builder(read_buffer: ReadBuffer, umas_request_function_key: int):
        read_buffer.push_context("UmasPDUPlcStatusResponse")

        not_used: int = read_buffer.read_unsigned_int(
            logical_name="notUsed",
            bit_length=24,
            umas_request_function_key=umas_request_function_key,
        )

        number_of_blocks: int = read_buffer.read_unsigned_byte(
            logical_name="numberOfBlocks",
            bit_length=8,
            umas_request_function_key=umas_request_function_key,
        )

        blocks: List[Any] = read_buffer.read_array_field(
            logical_name="blocks",
            read_function=read_buffer.read_unsigned_int,
            count=number_of_blocks,
            umas_request_function_key=umas_request_function_key,
        )

        read_buffer.pop_context("UmasPDUPlcStatusResponse")
        # Create the instance
        return UmasPDUPlcStatusResponseBuilder(not_used, number_of_blocks, blocks)

    def equals(self, o: object) -> bool:
        if self == o:
            return True

        if not isinstance(o, UmasPDUPlcStatusResponse):
            return False

        that: UmasPDUPlcStatusResponse = UmasPDUPlcStatusResponse(o)
        return (
            (self.not_used == that.not_used)
            and (self.number_of_blocks == that.number_of_blocks)
            and (self.blocks == that.blocks)
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
class UmasPDUPlcStatusResponseBuilder:
    not_used: int
    number_of_blocks: int
    blocks: List[int]

    def build(self, pairing_key) -> UmasPDUPlcStatusResponse:
        umas_pdu_plc_status_response: UmasPDUPlcStatusResponse = (
            UmasPDUPlcStatusResponse(
                pairing_key, self.not_used, self.number_of_blocks, self.blocks
            )
        )
        return umas_pdu_plc_status_response
