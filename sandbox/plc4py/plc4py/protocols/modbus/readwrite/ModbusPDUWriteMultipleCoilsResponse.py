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

from plc4py.api.messages.PlcMessage import PlcMessage
from plc4py.protocols.modbus.readwrite.ModbusPDU import ModbusPDU
from plc4py.protocols.modbus.readwrite.ModbusPDU import ModbusPDUBuilder
from plc4py.spi.generation.ReadBuffer import ReadBuffer
from plc4py.spi.generation.WriteBuffer import WriteBuffer
import math


@dataclass
class ModbusPDUWriteMultipleCoilsResponse(PlcMessage, ModbusPDU):
    starting_address: int
    quantity: int
    # Accessors for discriminator values.
    error_flag: bool = False
    function_flag: int = 0x0F
    response: bool = True

    def __post_init__(self):
        super().__init__()

    def serialize_modbus_pdu_child(self, write_buffer: WriteBuffer):
        write_buffer.push_context("ModbusPDUWriteMultipleCoilsResponse")

        # Simple Field (startingAddress)
        write_buffer.write_unsigned_short(
            self.starting_address, logical_name="startingAddress"
        )

        # Simple Field (quantity)
        write_buffer.write_unsigned_short(self.quantity, logical_name="quantity")

        write_buffer.pop_context("ModbusPDUWriteMultipleCoilsResponse")

    def length_in_bytes(self) -> int:
        return int(math.ceil(float(self.get_length_in_bits() / 8.0)))

    def get_length_in_bits(self) -> int:
        length_in_bits: int = super().get_length_in_bits()
        _value: ModbusPDUWriteMultipleCoilsResponse = self

        # Simple field (startingAddress)
        length_in_bits += 16

        # Simple field (quantity)
        length_in_bits += 16

        return length_in_bits

    @staticmethod
    def static_parse_builder(read_buffer: ReadBuffer, response: bool):
        read_buffer.push_context("ModbusPDUWriteMultipleCoilsResponse")

        self.starting_address = read_simple_field("startingAddress", read_unsigned_int)

        self.quantity = read_simple_field("quantity", read_unsigned_int)

        read_buffer.pop_context("ModbusPDUWriteMultipleCoilsResponse")
        # Create the instance
        return ModbusPDUWriteMultipleCoilsResponseBuilder(starting_address, quantity)

    def equals(self, o: object) -> bool:
        if self == o:
            return True

        if not isinstance(o, ModbusPDUWriteMultipleCoilsResponse):
            return False

        that: ModbusPDUWriteMultipleCoilsResponse = ModbusPDUWriteMultipleCoilsResponse(
            o
        )
        return (
            (self.starting_address == that.starting_address)
            and (self.quantity == that.quantity)
            and super().equals(that)
            and True
        )

    def hash_code(self) -> int:
        return hash(self)

    def __str__(self) -> str:
        write_buffer_box_based: WriteBufferBoxBased = WriteBufferBoxBased(True, True)
        try:
            write_buffer_box_based.writeSerializable(self)
        except SerializationException as e:
            raise RuntimeException(e)

        return "\n" + str(write_buffer_box_based.get_box()) + "\n"


@dataclass
class ModbusPDUWriteMultipleCoilsResponseBuilder(ModbusPDUBuilder):
    startingAddress: int
    quantity: int

    def __post_init__(self):
        pass

    def build(
        self,
    ) -> ModbusPDUWriteMultipleCoilsResponse:
        modbus_pdu_write_multiple_coils_response: ModbusPDUWriteMultipleCoilsResponse = ModbusPDUWriteMultipleCoilsResponse(
            self.starting_address, self.quantity
        )
        return modbus_pdu_write_multiple_coils_response
