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
class ModbusPDUWriteMultipleCoilsResponse(ModbusPDU):
    starting_address: int
    quantity: int
    # Accessors for discriminator values.
    error_flag: ClassVar[bool] = False
    function_flag: ClassVar[int] = 0x0F
    response: ClassVar[bool] = True

    def serialize_modbus_pdu_child(self, write_buffer: WriteBuffer):
        write_buffer.push_context("ModbusPDUWriteMultipleCoilsResponse")

        # Simple Field (startingAddress)
        write_buffer.write_unsigned_short(
            self.starting_address, bit_length=16, logical_name="startingAddress"
        )

        # Simple Field (quantity)
        write_buffer.write_unsigned_short(
            self.quantity, bit_length=16, logical_name="quantity"
        )

        write_buffer.pop_context("ModbusPDUWriteMultipleCoilsResponse")

    def length_in_bytes(self) -> int:
        return int(math.ceil(float(self.length_in_bits() / 8.0)))

    def length_in_bits(self) -> int:
        length_in_bits: int = super().length_in_bits()
        _value: ModbusPDUWriteMultipleCoilsResponse = self

        # Simple field (startingAddress)
        length_in_bits += 16

        # Simple field (quantity)
        length_in_bits += 16

        return length_in_bits

    @staticmethod
    def static_parse_builder(read_buffer: ReadBuffer, response: bool):
        read_buffer.push_context("ModbusPDUWriteMultipleCoilsResponse")

        if isinstance(response, str):
            response = bool(strtobool(response))

        starting_address: int = read_buffer.read_unsigned_short(
            logical_name="starting_address", bit_length=16, response=response
        )

        quantity: int = read_buffer.read_unsigned_short(
            logical_name="quantity", bit_length=16, response=response
        )

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
        # TODO:- Implement a generic python object to probably json convertor or something.
        return ""


@dataclass
class ModbusPDUWriteMultipleCoilsResponseBuilder:
    starting_address: int
    quantity: int

    def build(
        self,
    ) -> ModbusPDUWriteMultipleCoilsResponse:
        modbus_pduwrite_multiple_coils_response: ModbusPDUWriteMultipleCoilsResponse = (
            ModbusPDUWriteMultipleCoilsResponse(self.starting_address, self.quantity)
        )
        return modbus_pduwrite_multiple_coils_response
