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
from plc4py.protocols.modbus.readwrite.ModbusDeviceInformationLevel import (
    ModbusDeviceInformationLevel,
)
from plc4py.protocols.modbus.readwrite.ModbusPDU import ModbusPDU
from plc4py.spi.generation.ReadBuffer import ReadBuffer
from plc4py.spi.generation.WriteBuffer import WriteBuffer
from plc4py.utils.ConnectionStringHandling import strtobool
from typing import ClassVar
import math


@dataclass
class ModbusPDUReadDeviceIdentificationRequest(ModbusPDU):
    level: ModbusDeviceInformationLevel
    object_id: int
    MEI_TYPE: int = 0x0E
    # Accessors for discriminator values.
    error_flag: ClassVar[bool] = False
    function_flag: ClassVar[int] = 0x2B
    response: ClassVar[bool] = False

    def serialize_modbus_pdu_child(self, write_buffer: WriteBuffer):
        write_buffer.push_context("ModbusPDUReadDeviceIdentificationRequest")

        # Const Field (meiType)
        write_buffer.write_unsigned_byte(self.MEI_TYPE, logical_name="meiType")

        # Simple Field (level)
        write_buffer.write_unsigned_byte(self.level, logical_name="level")

        # Simple Field (objectId)
        write_buffer.write_unsigned_byte(
            self.object_id, bit_length=8, logical_name="objectId"
        )

        write_buffer.pop_context("ModbusPDUReadDeviceIdentificationRequest")

    def length_in_bytes(self) -> int:
        return int(math.ceil(float(self.length_in_bits() / 8.0)))

    def length_in_bits(self) -> int:
        length_in_bits: int = super().length_in_bits()
        _value: ModbusPDUReadDeviceIdentificationRequest = self

        # Const Field (meiType)
        length_in_bits += 8

        # Simple field (level)
        length_in_bits += 8

        # Simple field (objectId)
        length_in_bits += 8

        return length_in_bits

    @staticmethod
    def static_parse_builder(read_buffer: ReadBuffer, response: bool):
        read_buffer.push_context("ModbusPDUReadDeviceIdentificationRequest")

        if isinstance(response, str):
            response = bool(strtobool(response))

        MEI_TYPE: int = read_buffer.read_unsigned_byte(
            logical_name="mei_type", response=response
        )

        level: ModbusDeviceInformationLevel = read_buffer.read_enum(
            read_function=ModbusDeviceInformationLevel,
            bit_length=8,
            logical_name="level",
            response=response,
        )

        object_id: int = read_buffer.read_unsigned_byte(
            logical_name="object_id", bit_length=8, response=response
        )

        read_buffer.pop_context("ModbusPDUReadDeviceIdentificationRequest")
        # Create the instance
        return ModbusPDUReadDeviceIdentificationRequestBuilder(level, object_id)

    def equals(self, o: object) -> bool:
        if self == o:
            return True

        if not isinstance(o, ModbusPDUReadDeviceIdentificationRequest):
            return False

        that: ModbusPDUReadDeviceIdentificationRequest = (
            ModbusPDUReadDeviceIdentificationRequest(o)
        )
        return (
            (self.level == that.level)
            and (self.object_id == that.object_id)
            and super().equals(that)
            and True
        )

    def hash_code(self) -> int:
        return hash(self)

    def __str__(self) -> str:
        # TODO:- Implement a generic python object to probably json convertor or something.
        return ""


@dataclass
class ModbusPDUReadDeviceIdentificationRequestBuilder:
    level: ModbusDeviceInformationLevel
    object_id: int

    def build(
        self,
    ) -> ModbusPDUReadDeviceIdentificationRequest:
        modbus_pduread_device_identification_request: (
            ModbusPDUReadDeviceIdentificationRequest
        ) = ModbusPDUReadDeviceIdentificationRequest(self.level, self.object_id)
        return modbus_pduread_device_identification_request
