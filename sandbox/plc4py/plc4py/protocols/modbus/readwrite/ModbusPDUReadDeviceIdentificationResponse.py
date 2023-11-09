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
from plc4py.protocols.modbus.readwrite.ModbusDeviceInformationConformityLevel import (
    ModbusDeviceInformationConformityLevel,
)
from plc4py.protocols.modbus.readwrite.ModbusDeviceInformationLevel import (
    ModbusDeviceInformationLevel,
)
from plc4py.protocols.modbus.readwrite.ModbusDeviceInformationMoreFollows import (
    ModbusDeviceInformationMoreFollows,
)
from plc4py.protocols.modbus.readwrite.ModbusDeviceInformationObject import (
    ModbusDeviceInformationObject,
)
from plc4py.protocols.modbus.readwrite.ModbusPDU import ModbusPDU
from plc4py.protocols.modbus.readwrite.ModbusPDU import ModbusPDUBuilder
from plc4py.spi.generation.ReadBuffer import ReadBuffer
from plc4py.spi.generation.WriteBuffer import WriteBuffer
from typing import Any
from typing import List
import math


@dataclass
class ModbusPDUReadDeviceIdentificationResponse(ModbusPDU):
    level: ModbusDeviceInformationLevel
    individual_access: bool
    conformity_level: ModbusDeviceInformationConformityLevel
    more_follows: ModbusDeviceInformationMoreFollows
    next_object_id: int
    objects: List[ModbusDeviceInformationObject]
    MEI_TYPE: int = 0x0E
    # Accessors for discriminator values.
    error_flag: bool = False
    function_flag: int = 0x2B
    response: bool = True

    def serialize_modbus_pdu_child(self, write_buffer: WriteBuffer):
        write_buffer.push_context("ModbusPDUReadDeviceIdentificationResponse")

        # Const Field (meiType)
        write_buffer.write_unsigned_byte(self.MEI_TYPE, logical_name="meiType")

        # Simple Field (level)
        write_buffer.write_unsigned_byte(self.level, logical_name="level")

        # Simple Field (individualAccess)
        write_buffer.write_bit(self.individual_access, logical_name="individualAccess")

        # Simple Field (conformityLevel)
        write_buffer.write_unsigned_byte(
            self.conformity_level, logical_name="conformityLevel"
        )

        # Simple Field (moreFollows)
        write_buffer.write_unsigned_byte(self.more_follows, logical_name="moreFollows")

        # Simple Field (nextObjectId)
        write_buffer.write_unsigned_byte(
            self.next_object_id, logical_name="nextObjectId"
        )

        # Implicit Field (number_of_objects) (Used for parsing, but its value is not stored as it's implicitly given by the objects content)
        number_of_objects: int = int(len(self.objects))
        write_buffer.write_unsigned_byte(
            number_of_objects, logical_name="numberOfObjects"
        )

        # Array Field (objects)
        write_buffer.write_complex_array(self.objects, logical_name="objects")

        write_buffer.pop_context("ModbusPDUReadDeviceIdentificationResponse")

    def length_in_bytes(self) -> int:
        return int(math.ceil(float(self.length_in_bits() / 8.0)))

    def length_in_bits(self) -> int:
        length_in_bits: int = super().length_in_bits()
        _value: ModbusPDUReadDeviceIdentificationResponse = self

        # Const Field (meiType)
        length_in_bits += 8

        # Simple field (level)
        length_in_bits += 8

        # Simple field (individualAccess)
        length_in_bits += 1

        # Simple field (conformityLevel)
        length_in_bits += 7

        # Simple field (moreFollows)
        length_in_bits += 8

        # Simple field (nextObjectId)
        length_in_bits += 8

        # Implicit Field (numberOfObjects)
        length_in_bits += 8

        # Array field
        if self.objects is not None:
            for element in self.objects:
                length_in_bits += element.length_in_bits()

        return length_in_bits

    @staticmethod
    def static_parse_builder(read_buffer: ReadBuffer, response: bool):
        read_buffer.push_context("ModbusPDUReadDeviceIdentificationResponse")

        MEI_TYPE: int = read_buffer.read_unsigned_short(logical_name="meiType")

        level: ModbusDeviceInformationLevel = read_buffer.read_complex(
            read_function=ModbusDeviceInformationLevel, logical_name="level"
        )

        individual_access: bool = read_buffer.read_bit(logical_name="individualAccess")

        conformity_level: ModbusDeviceInformationConformityLevel = (
            read_buffer.read_complex(
                read_function=ModbusDeviceInformationConformityLevel,
                logical_name="conformityLevel",
            )
        )

        more_follows: ModbusDeviceInformationMoreFollows = read_buffer.read_complex(
            read_function=ModbusDeviceInformationMoreFollows, logical_name="moreFollows"
        )

        next_object_id: int = read_buffer.read_unsigned_short(
            logical_name="nextObjectId"
        )

        number_of_objects: int = read_buffer.read_unsigned_short(
            logical_name="numberOfObjects"
        )

        objects: List[Any] = read_buffer.read_array_field(
            logical_name="objects",
            read_function=ModbusDeviceInformationObject.static_parse,
            count=number_of_objects,
        )

        read_buffer.pop_context("ModbusPDUReadDeviceIdentificationResponse")
        # Create the instance
        return ModbusPDUReadDeviceIdentificationResponseBuilder(
            level,
            individual_access,
            conformity_level,
            more_follows,
            next_object_id,
            objects,
        )

    def equals(self, o: object) -> bool:
        if self == o:
            return True

        if not isinstance(o, ModbusPDUReadDeviceIdentificationResponse):
            return False

        that: ModbusPDUReadDeviceIdentificationResponse = (
            ModbusPDUReadDeviceIdentificationResponse(o)
        )
        return (
            (self.level == that.level)
            and (self.individual_access == that.individual_access)
            and (self.conformity_level == that.conformity_level)
            and (self.more_follows == that.more_follows)
            and (self.next_object_id == that.next_object_id)
            and (self.objects == that.objects)
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
class ModbusPDUReadDeviceIdentificationResponseBuilder(ModbusPDUBuilder):
    level: ModbusDeviceInformationLevel
    individual_access: bool
    conformity_level: ModbusDeviceInformationConformityLevel
    more_follows: ModbusDeviceInformationMoreFollows
    next_object_id: int
    objects: List[ModbusDeviceInformationObject]

    def build(
        self,
    ) -> ModbusPDUReadDeviceIdentificationResponse:
        modbus_pdu_read_device_identification_response: ModbusPDUReadDeviceIdentificationResponse = ModbusPDUReadDeviceIdentificationResponse(
            self.level,
            self.individual_access,
            self.conformity_level,
            self.more_follows,
            self.next_object_id,
            self.objects,
        )
        return modbus_pdu_read_device_identification_response
