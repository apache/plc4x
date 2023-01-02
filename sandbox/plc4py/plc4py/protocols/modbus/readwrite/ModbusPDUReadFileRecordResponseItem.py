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

from ctypes import c_byte
from ctypes import c_uint8
from plc4py.api.messages.PlcMessage import PlcMessage
from typing import List
import math


@dataclass
class ModbusPDUReadFileRecordResponseItem(PlcMessage):
    reference_type: c_uint8
    data: List[c_byte]

    def __post_init__(self):
        super().__init__()

    def serialize(self, write_buffer: WriteBuffer):
        position_aware: PositionAware = write_buffer
        start_pos: int = position_aware.get_pos()
        write_buffer.push_context("ModbusPDUReadFileRecordResponseItem")

        # Implicit Field (data_length) (Used for parsing, but its value is not stored as it's implicitly given by the objects content)
        data_length: c_uint8 = c_uint8(((COUNT(self.data())) + (1)))
        write_implicit_field(
            "dataLength", data_length, write_unsigned_short(write_buffer, 8)
        )

        # Simple Field (referenceType)
        write_simple_field(
            "referenceType", self.reference_type, write_unsigned_short(write_buffer, 8)
        )

        # Array Field (data)
        write_byte_array_field("data", self.data, writeByteArray(write_buffer, 8))

        write_buffer.pop_context("ModbusPDUReadFileRecordResponseItem")

    def length_in_bytes(self) -> int:
        return int(math.ceil(float(self.get_length_in_bits() / 8.0)))

    def get_length_in_bits(self) -> int:
        length_in_bits: int = 0
        _value: ModbusPDUReadFileRecordResponseItem = self

        # Implicit Field (dataLength)
        length_in_bits += 8

        # Simple field (referenceType)
        length_in_bits += 8

        # Array field
        if self.data is not None:
            length_in_bits += 8 * self.data.length

        return length_in_bits

    def static_parse(read_buffer: ReadBuffer, args):
        position_aware: PositionAware = read_buffer
        return staticParse(read_buffer)

    @staticmethod
    def static_parse_context(read_buffer: ReadBuffer):
        read_buffer.pull_context("ModbusPDUReadFileRecordResponseItem")
        position_aware: PositionAware = read_buffer
        start_pos: int = position_aware.get_pos()
        cur_pos: int = 0

        data_length: c_uint8 = read_implicit_field(
            "dataLength", read_unsigned_short(read_buffer, 8)
        )

        reference_type: c_uint8 = read_simple_field(
            "referenceType", read_unsigned_short(read_buffer, 8)
        )

        data: List[c_byte] = read_buffer.read_byte_array(
            "data", int((dataLength) - (1))
        )

        read_buffer.close_context("ModbusPDUReadFileRecordResponseItem")
        # Create the instance
        _modbus_pdu_read_file_record_response_item: ModbusPDUReadFileRecordResponseItem = ModbusPDUReadFileRecordResponseItem(
            reference_type, data
        )
        return _modbus_pdu_read_file_record_response_item

    def equals(self, o: object) -> bool:
        if self == o:
            return True

        if not isinstance(o, ModbusPDUReadFileRecordResponseItem):
            return False

        that: ModbusPDUReadFileRecordResponseItem = ModbusPDUReadFileRecordResponseItem(
            o
        )
        return (
            (self.reference_type == that.reference_type)
            and (self.data == that.data)
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
