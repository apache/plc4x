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
from ctypes import c_uint16
from ctypes import c_uint8
from plc4py.api.messages.PlcMessage import PlcMessage
from plc4py.protocols.modbus.readwrite.ModbusPDU import ModbusPDU
from plc4py.protocols.modbus.readwrite.ModbusPDU import ModbusPDUBuilder
from plc4py.spi.generation.WriteBuffer import WriteBuffer
import math


@dataclass
class ModbusPDUDiagnosticResponse(PlcMessage, ModbusPDU):
    sub_function: c_uint16
    data: c_uint16
    # Accessors for discriminator values.
    error_flag: c_bool = False
    function_flag: c_uint8 = 0x08
    response: c_bool = True

    def __post_init__(self):
        super().__init__()

    def serialize_modbus_pdu_child(self, write_buffer: WriteBuffer):
        write_buffer.push_context("ModbusPDUDiagnosticResponse")

        # Simple Field (subFunction)
        write_buffer.write_unsigned_short(self.sub_function, logical_name="subFunction")

        # Simple Field (data)
        write_buffer.write_unsigned_short(self.data, logical_name="data")

        write_buffer.pop_context("ModbusPDUDiagnosticResponse")

    def length_in_bytes(self) -> int:
        return int(math.ceil(float(self.get_length_in_bits() / 8.0)))

    def get_length_in_bits(self) -> int:
        length_in_bits: int = super().get_length_in_bits()
        _value: ModbusPDUDiagnosticResponse = self

        # Simple field (subFunction)
        length_in_bits += 16

        # Simple field (data)
        length_in_bits += 16

        return length_in_bits

    @staticmethod
    def static_parse_builder(read_buffer: ReadBuffer, response: c_bool):
        read_buffer.pull_context("ModbusPDUDiagnosticResponse")
        cur_pos: int = 0

        sub_function: c_uint16 = read_simple_field("subFunction", read_unsigned_int)

        data: c_uint16 = read_simple_field("data", read_unsigned_int)

        read_buffer.close_context("ModbusPDUDiagnosticResponse")
        # Create the instance
        return ModbusPDUDiagnosticResponseBuilder(sub_function, data)

    def equals(self, o: object) -> bool:
        if self == o:
            return True

        if not isinstance(o, ModbusPDUDiagnosticResponse):
            return False

        that: ModbusPDUDiagnosticResponse = ModbusPDUDiagnosticResponse(o)
        return (
            (self.sub_function == that.sub_function)
            and (self.data == that.data)
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
class ModbusPDUDiagnosticResponseBuilder(ModbusPDUBuilder):
    subFunction: c_uint16
    data: c_uint16

    def __post_init__(self):
        pass

    def build(
        self,
    ) -> ModbusPDUDiagnosticResponse:
        modbus_pdu_diagnostic_response: ModbusPDUDiagnosticResponse = (
            ModbusPDUDiagnosticResponse(self.sub_function, self.data)
        )
        return modbus_pdu_diagnostic_response
