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
from plc4py.protocols.modbus.readwrite.ModbusErrorCode import ModbusErrorCode
from plc4py.protocols.modbus.readwrite.ModbusPDU import ModbusPDU
from plc4py.protocols.modbus.readwrite.ModbusPDU import ModbusPDUBuilder
from plc4py.spi.generation.ReadBuffer import ReadBuffer
from plc4py.spi.generation.WriteBuffer import WriteBuffer
import math


@dataclass
class ModbusPDUError(PlcMessage, ModbusPDU):
    exception_code: ModbusErrorCode
    # Accessors for discriminator values.
    error_flag: bool = True
    function_flag: int = 0
    response: bool = False

    def __post_init__(self):
        super().__init__()

    def serialize_modbus_pdu_child(self, write_buffer: WriteBuffer):
        write_buffer.push_context("ModbusPDUError")

        # Simple Field (exceptionCode)
        write_buffer.DataWriterEnumDefault(
            ModbusErrorCode.value, ModbusErrorCode.name, write_unsigned_byte
        )(self.exception_code, logical_name="exceptionCode")

        write_buffer.pop_context("ModbusPDUError")

    def length_in_bytes(self) -> int:
        return int(math.ceil(float(self.get_length_in_bits() / 8.0)))

    def get_length_in_bits(self) -> int:
        length_in_bits: int = super().get_length_in_bits()
        _value: ModbusPDUError = self

        # Simple field (exceptionCode)
        length_in_bits += 8

        return length_in_bits

    @staticmethod
    def static_parse_builder(read_buffer: ReadBuffer, response: bool):
        read_buffer.push_context("ModbusPDUError")

        self.exception_code = read_enum_field(
            "exceptionCode",
            "ModbusErrorCode",
            DataReaderEnumDefault(ModbusErrorCode.enumForValue, read_unsigned_short),
        )

        read_buffer.pop_context("ModbusPDUError")
        # Create the instance
        return ModbusPDUErrorBuilder(exception_code)

    def equals(self, o: object) -> bool:
        if self == o:
            return True

        if not isinstance(o, ModbusPDUError):
            return False

        that: ModbusPDUError = ModbusPDUError(o)
        return (
            (self.exception_code == that.exception_code)
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
class ModbusPDUErrorBuilder(ModbusPDUBuilder):
    exceptionCode: ModbusErrorCode

    def __post_init__(self):
        pass

    def build(
        self,
    ) -> ModbusPDUError:
        modbus_pdu_error: ModbusPDUError = ModbusPDUError(self.exception_code)
        return modbus_pdu_error
