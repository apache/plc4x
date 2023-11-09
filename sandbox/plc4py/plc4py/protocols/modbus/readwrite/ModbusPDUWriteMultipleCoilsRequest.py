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
from plc4py.protocols.modbus.readwrite.ModbusPDU import ModbusPDUBuilder
from plc4py.spi.generation.ReadBuffer import ReadBuffer
from plc4py.spi.generation.WriteBuffer import WriteBuffer
from typing import Any
from typing import List
import math
    
@dataclass
class ModbusPDUWriteMultipleCoilsRequest(ModbusPDU):
    starting_address: int
    quantity: int
    value: List[int]
    # Accessors for discriminator values.
    error_flag: bool = False
    function_flag: int = 0x0F
    response: bool = False



    def serialize_modbus_pdu_child(self, write_buffer: WriteBuffer):
        write_buffer.push_context("ModbusPDUWriteMultipleCoilsRequest")

        # Simple Field (startingAddress)
        write_buffer.write_unsigned_short(self.starting_address, logical_name="startingAddress")

        # Simple Field (quantity)
        write_buffer.write_unsigned_short(self.quantity, logical_name="quantity")

        # Implicit Field (byte_count) (Used for parsing, but its value is not stored as it's implicitly given by the objects content)
        byte_count: int = (int(len(self.value)))
        write_buffer.write_unsigned_byte(byte_count, logical_name="byteCount")

        # Array Field (value)
        write_buffer.write_byte_array(self.value, logical_name="value")

        write_buffer.pop_context("ModbusPDUWriteMultipleCoilsRequest")


    def length_in_bytes(self) -> int:
        return int(math.ceil(float(self.length_in_bits() / 8.0)))

    def length_in_bits(self) -> int:
        length_in_bits: int = super().length_in_bits()
        _value: ModbusPDUWriteMultipleCoilsRequest = self

        # Simple field (startingAddress)
        length_in_bits += 16

        # Simple field (quantity)
        length_in_bits += 16

        # Implicit Field (byteCount)
        length_in_bits += 8

        # Array field
        if self.value is not None:
            length_in_bits += 8 * len(self.value)


        return length_in_bits


    @staticmethod
    def static_parse_builder(read_buffer: ReadBuffer, response: bool):
        read_buffer.push_context("ModbusPDUWriteMultipleCoilsRequest")

        starting_address: int = read_buffer.read_unsigned_int(logical_name="startingAddress")  

        quantity: int = read_buffer.read_unsigned_int(logical_name="quantity")  

        byte_count: int = read_buffer.read_unsigned_short(logical_name="byteCount")

        value: List[Any] = read_buffer.read_array_field(logical_name="value", read_function=read_buffer.read_byte, count=byte_count)

        read_buffer.pop_context("ModbusPDUWriteMultipleCoilsRequest")
        # Create the instance
        return ModbusPDUWriteMultipleCoilsRequestBuilder(starting_address, quantity, value )


    def equals(self, o: object) -> bool:
        if self == o:
            return True

        if not isinstance(o, ModbusPDUWriteMultipleCoilsRequest):
            return False

        that: ModbusPDUWriteMultipleCoilsRequest = ModbusPDUWriteMultipleCoilsRequest(o)
        return (self.starting_address == that.starting_address) and (self.quantity == that.quantity) and (self.value == that.value) and super().equals(that) and True

    def hash_code(self) -> int:
        return hash(self)

    def __str__(self) -> str:
        pass
        #write_buffer_box_based: WriteBufferBoxBased = WriteBufferBoxBased(True, True)
        #try:
        #    write_buffer_box_based.writeSerializable(self)
        #except SerializationException as e:
        #    raise PlcRuntimeException(e)

        #return "\n" + str(write_buffer_box_based.get_box()) + "\n"


@dataclass
class ModbusPDUWriteMultipleCoilsRequestBuilder(ModbusPDUBuilder):
    starting_address: int
    quantity: int
    value: List[int]

    def build(self,) -> ModbusPDUWriteMultipleCoilsRequest:
        modbus_pdu_write_multiple_coils_request: ModbusPDUWriteMultipleCoilsRequest = ModbusPDUWriteMultipleCoilsRequest(self.starting_address, self.quantity, self.value )
        return modbus_pdu_write_multiple_coils_request



