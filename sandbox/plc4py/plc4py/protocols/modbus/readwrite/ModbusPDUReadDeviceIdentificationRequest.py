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
from plc4py.protocols.modbus.readwrite.ModbusDeviceInformationLevel import (
    ModbusDeviceInformationLevel,
)
from plc4py.protocols.modbus.readwrite.ModbusPDU import ModbusPDU
from plc4py.protocols.modbus.readwrite.ModbusPDU import ModbusPDUBuilder
from plc4py.spi.generation.ReadBuffer import ReadBuffer
from plc4py.spi.generation.WriteBuffer import WriteBuffer
import math


@dataclass
class ModbusPDUReadDeviceIdentificationRequest(PlcMessage, ModbusPDU):
    level: ModbusDeviceInformationLevel
    object_id: int
    MEITYPE: int = 0x0E
    # Accessors for discriminator values.
    error_flag: bool = False
    function_flag: int = 0x2B
    response: bool = False

    def __post_init__(self):
        super().__init__()

    def serialize_modbus_pdu_child(self, write_buffer: WriteBuffer):
        write_buffer.push_context("ModbusPDUReadDeviceIdentificationRequest")

        # Const Field (meiType)
        write_buffer.write_unsigned_byte(self.mei_type.value, logical_name="meiType")

        # Simple Field (level)
        write_buffer.DataWriterEnumDefault(
            ModbusDeviceInformationLevel.value,
            ModbusDeviceInformationLevel.name,
            write_unsigned_byte,
        )(self.level, logical_name="level")

        # Simple Field (objectId)
        write_buffer.write_unsigned_byte(self.object_id, logical_name="objectId")

        write_buffer.pop_context("ModbusPDUReadDeviceIdentificationRequest")

    def length_in_bytes(self) -> int:
        return int(math.ceil(float(self.get_length_in_bits() / 8.0)))

    def get_length_in_bits(self) -> int:
        length_in_bits: int = super().get_length_in_bits()
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

        self.mei_type: int = read_const_field(
            "meiType",
            read_unsigned_short,
            ModbusPDUReadDeviceIdentificationRequest.MEITYPE,
        )

        self.level = read_enum_field(
            "level",
            "ModbusDeviceInformationLevel",
            DataReaderEnumDefault(
                ModbusDeviceInformationLevel.enumForValue, read_unsigned_short
            ),
        )

        self.object_id = read_simple_field("objectId", read_unsigned_short)

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
        write_buffer_box_based: WriteBufferBoxBased = WriteBufferBoxBased(True, True)
        try:
            write_buffer_box_based.writeSerializable(self)
        except SerializationException as e:
            raise RuntimeException(e)

        return "\n" + str(write_buffer_box_based.get_box()) + "\n"


@dataclass
class ModbusPDUReadDeviceIdentificationRequestBuilder(ModbusPDUBuilder):
    level: ModbusDeviceInformationLevel
    objectId: int

    def __post_init__(self):
        pass

    def build(
        self,
    ) -> ModbusPDUReadDeviceIdentificationRequest:
        modbus_pdu_read_device_identification_request: ModbusPDUReadDeviceIdentificationRequest = ModbusPDUReadDeviceIdentificationRequest(
            self.level, self.object_id
        )
        return modbus_pdu_read_device_identification_request
