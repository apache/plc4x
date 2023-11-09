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
import math
    
@dataclass
class ModbusPDUWriteSingleRegisterResponse(ModbusPDU):
    address: int
    value: int
    # Accessors for discriminator values.
    error_flag: bool = False
    function_flag: int = 0x06
    response: bool = True



    def serialize_modbus_pdu_child(self, write_buffer: WriteBuffer):
        write_buffer.push_context("ModbusPDUWriteSingleRegisterResponse")

        # Simple Field (address)
        write_buffer.write_unsigned_short(self.address, logical_name="address")

        # Simple Field (value)
        write_buffer.write_unsigned_short(self.value, logical_name="value")

        write_buffer.pop_context("ModbusPDUWriteSingleRegisterResponse")


    def length_in_bytes(self) -> int:
        return int(math.ceil(float(self.length_in_bits() / 8.0)))

    def length_in_bits(self) -> int:
        length_in_bits: int = super().length_in_bits()
        _value: ModbusPDUWriteSingleRegisterResponse = self

        # Simple field (address)
        length_in_bits += 16

        # Simple field (value)
        length_in_bits += 16

        return length_in_bits


    @staticmethod
    def static_parse_builder(read_buffer: ReadBuffer, response: bool):
        read_buffer.push_context("ModbusPDUWriteSingleRegisterResponse")

        address: int = read_buffer.read_unsigned_int(logical_name="address")  

        value: int = read_buffer.read_unsigned_int(logical_name="value")  

        read_buffer.pop_context("ModbusPDUWriteSingleRegisterResponse")
        # Create the instance
        return ModbusPDUWriteSingleRegisterResponseBuilder(address, value )


    def equals(self, o: object) -> bool:
        if self == o:
            return True

        if not isinstance(o, ModbusPDUWriteSingleRegisterResponse):
            return False

        that: ModbusPDUWriteSingleRegisterResponse = ModbusPDUWriteSingleRegisterResponse(o)
        return (self.address == that.address) and (self.value == that.value) and super().equals(that) and True

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
class ModbusPDUWriteSingleRegisterResponseBuilder(ModbusPDUBuilder):
    address: int
    value: int

    def build(self,) -> ModbusPDUWriteSingleRegisterResponse:
        modbus_pdu_write_single_register_response: ModbusPDUWriteSingleRegisterResponse = ModbusPDUWriteSingleRegisterResponse(self.address, self.value )
        return modbus_pdu_write_single_register_response



