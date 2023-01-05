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

from ctypes import c_bool
from ctypes import c_byte
from ctypes import c_uint16
from ctypes import c_uint8
from plc4py.api.messages.PlcMessage import PlcMessage
from plc4py.protocols.modbus.readwrite.ModbusPDU import ModbusPDU
from plc4py.protocols.modbus.readwrite.ModbusPDU import ModbusPDUBuilder
from typing import List
import math


@dataclass
class ModbusPDUReadWriteMultipleHoldingRegistersRequest(PlcMessage, ModbusPDU):
    read_starting_address: c_uint16
    read_quantity: c_uint16
    write_starting_address: c_uint16
    write_quantity: c_uint16
    value: List[c_byte]
    # Accessors for discriminator values.
    error_flag: c_bool = False
    function_flag: c_uint8 = 0x17
    response: c_bool = False

    def __post_init__(self):
        super().__init__()

    def serialize_modbus_pdu_child(self, write_buffer: WriteBuffer):
        position_aware: PositionAware = write_buffer
        start_pos: int = position_aware.get_pos()
        write_buffer.push_context("ModbusPDUReadWriteMultipleHoldingRegistersRequest")

        # Simple Field (readStartingAddress)
        write_simple_field(
            "readStartingAddress",
            self.read_starting_address,
            write_unsigned_int(write_buffer, 16),
        )

        # Simple Field (readQuantity)
        write_simple_field(
            "readQuantity", self.read_quantity, write_unsigned_int(write_buffer, 16)
        )

        # Simple Field (writeStartingAddress)
        write_simple_field(
            "writeStartingAddress",
            self.write_starting_address,
            write_unsigned_int(write_buffer, 16),
        )

        # Simple Field (writeQuantity)
        write_simple_field(
            "writeQuantity", self.write_quantity, write_unsigned_int(write_buffer, 16)
        )

        # Implicit Field (byte_count) (Used for parsing, but its value is not stored as it's implicitly given by the objects content)
        byte_count: c_uint8 = c_uint8((COUNT(self.value())))
        write_implicit_field(
            "byteCount", byte_count, write_unsigned_short(write_buffer, 8)
        )

        # Array Field (value)
        write_byte_array_field("value", self.value, writeByteArray(write_buffer, 8))

        write_buffer.pop_context("ModbusPDUReadWriteMultipleHoldingRegistersRequest")

    def length_in_bytes(self) -> int:
        return int(math.ceil(float(self.get_length_in_bits() / 8.0)))

    def get_length_in_bits(self) -> int:
        length_in_bits: int = super().get_length_in_bits()
        _value: ModbusPDUReadWriteMultipleHoldingRegistersRequest = self

        # Simple field (readStartingAddress)
        length_in_bits += 16

        # Simple field (readQuantity)
        length_in_bits += 16

        # Simple field (writeStartingAddress)
        length_in_bits += 16

        # Simple field (writeQuantity)
        length_in_bits += 16

        # Implicit Field (byteCount)
        length_in_bits += 8

        # Array field
        if self.value is not None:
            length_in_bits += 8 * self.value.length

        return length_in_bits

    @staticmethod
    def static_parse_builder(read_buffer: ReadBuffer, response: c_bool):
        read_buffer.pull_context("ModbusPDUReadWriteMultipleHoldingRegistersRequest")
        position_aware: PositionAware = read_buffer
        start_pos: int = position_aware.get_pos()
        cur_pos: int = 0

        read_starting_address: c_uint16 = read_simple_field(
            "readStartingAddress", read_unsigned_int(read_buffer, 16)
        )

        read_quantity: c_uint16 = read_simple_field(
            "readQuantity", read_unsigned_int(read_buffer, 16)
        )

        write_starting_address: c_uint16 = read_simple_field(
            "writeStartingAddress", read_unsigned_int(read_buffer, 16)
        )

        write_quantity: c_uint16 = read_simple_field(
            "writeQuantity", read_unsigned_int(read_buffer, 16)
        )

        byte_count: c_uint8 = read_implicit_field(
            "byteCount", read_unsigned_short(read_buffer, 8)
        )

        value: List[c_byte] = read_buffer.read_byte_array("value", int(byteCount))

        read_buffer.close_context("ModbusPDUReadWriteMultipleHoldingRegistersRequest")
        # Create the instance
        return ModbusPDUReadWriteMultipleHoldingRegistersRequestBuilder(
            read_starting_address,
            read_quantity,
            write_starting_address,
            write_quantity,
            value,
        )

    def equals(self, o: object) -> bool:
        if self == o:
            return True

        if not isinstance(o, ModbusPDUReadWriteMultipleHoldingRegistersRequest):
            return False

        that: ModbusPDUReadWriteMultipleHoldingRegistersRequest = (
            ModbusPDUReadWriteMultipleHoldingRegistersRequest(o)
        )
        return (
            (self.read_starting_address == that.read_starting_address)
            and (self.read_quantity == that.read_quantity)
            and (self.write_starting_address == that.write_starting_address)
            and (self.write_quantity == that.write_quantity)
            and (self.value == that.value)
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
class ModbusPDUReadWriteMultipleHoldingRegistersRequestBuilder(ModbusPDUBuilder):
    readStartingAddress: c_uint16
    readQuantity: c_uint16
    writeStartingAddress: c_uint16
    writeQuantity: c_uint16
    value: List[c_byte]

    def __post_init__(self):
        pass

    def build(
        self,
    ) -> ModbusPDUReadWriteMultipleHoldingRegistersRequest:
        modbus_pdu_read_write_multiple_holding_registers_request: ModbusPDUReadWriteMultipleHoldingRegistersRequest = ModbusPDUReadWriteMultipleHoldingRegistersRequest(
            self.read_starting_address,
            self.read_quantity,
            self.write_starting_address,
            self.write_quantity,
            self.value,
        )
        return modbus_pdu_read_write_multiple_holding_registers_request
