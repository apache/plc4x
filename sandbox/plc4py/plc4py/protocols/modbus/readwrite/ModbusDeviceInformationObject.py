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
from plc4py.spi.generation.ReadBuffer import ReadBuffer
from plc4py.spi.generation.WriteBuffer import WriteBuffer
from typing import List
import math


@dataclass
class ModbusDeviceInformationObject(PlcMessage):
    object_id: int
    data: List[int]

    def __post_init__(self):
        super().__init__()

    def serialize(self, write_buffer: WriteBuffer):
        write_buffer.push_context("ModbusDeviceInformationObject")

        # Simple Field (objectId)
        write_buffer.write_unsigned_byte(self.object_id, logical_name="objectId")

        # Implicit Field (object_length) (Used for parsing, but its value is not stored as it's implicitly given by the objects content)
        object_length: int = int(len(self.data))
        write_buffer.write_unsigned_byte(object_length, logical_name="objectLength")

        # Array Field (data)
        write_buffer.write_byte_array(self.data, logical_name="data")

        write_buffer.pop_context("ModbusDeviceInformationObject")

    def length_in_bytes(self) -> int:
        return int(math.ceil(float(self.get_length_in_bits() / 8.0)))

    def get_length_in_bits(self) -> int:
        length_in_bits: int = 0
        _value: ModbusDeviceInformationObject = self

        # Simple field (objectId)
        length_in_bits += 8

        # Implicit Field (objectLength)
        length_in_bits += 8

        # Array field
        if self.data != None:
            length_in_bits += 8 * len(self.data)

        return length_in_bits

    def static_parse(self, read_buffer: ReadBuffer, args):
        return self.static_parse_context(read_buffer)

    @staticmethod
    def static_parse_context(read_buffer: ReadBuffer):
        read_buffer.push_context("ModbusDeviceInformationObject")

        self.object_id = read_simple_field("objectId", read_unsigned_short)

        object_length: int = read_implicit_field("objectLength", read_unsigned_short)

        self.data = read_buffer.read_byte_array("data", int(object_length))

        read_buffer.pop_context("ModbusDeviceInformationObject")
        # Create the instance
        _modbus_device_information_object: ModbusDeviceInformationObject = (
            ModbusDeviceInformationObject(object_id, data)
        )
        return _modbus_device_information_object

    def equals(self, o: object) -> bool:
        if self == o:
            return True

        if not isinstance(o, ModbusDeviceInformationObject):
            return False

        that: ModbusDeviceInformationObject = ModbusDeviceInformationObject(o)
        return (self.object_id == that.object_id) and (self.data == that.data) and True

    def hash_code(self) -> int:
        return hash(self)

    def __str__(self) -> str:
        write_buffer_box_based: WriteBufferBoxBased = WriteBufferBoxBased(True, True)
        try:
            write_buffer_box_based.writeSerializable(self)
        except SerializationException as e:
            raise RuntimeException(e)

        return "\n" + str(write_buffer_box_based.get_box()) + "\n"
