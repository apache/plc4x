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
from plc4py.protocols.modbus.readwrite.ModbusPDU import ModbusPDU
from plc4py.protocols.modbus.readwrite.ModbusPDU import ModbusPDUBuilder
from plc4py.protocols.modbus.readwrite.ModbusPDUReadFileRecordRequestItem import (
    ModbusPDUReadFileRecordRequestItem,
)
from plc4py.spi.generation.ReadBuffer import ReadBuffer
from plc4py.spi.generation.WriteBuffer import WriteBuffer
from sys import getsizeof
from typing import List
import math


@dataclass
class ModbusPDUReadFileRecordRequest(PlcMessage, ModbusPDU):
    items: List[ModbusPDUReadFileRecordRequestItem]
    # Accessors for discriminator values.
    error_flag: bool = False
    function_flag: int = 0x14
    response: bool = False

    def __post_init__(self):
        super().__init__()

    def serialize_modbus_pdu_child(self, write_buffer: WriteBuffer):
        write_buffer.push_context("ModbusPDUReadFileRecordRequest")

        # Implicit Field (byte_count) (Used for parsing, but its value is not stored as it's implicitly given by the objects content)
        byte_count: int = int(getsizeof(self.items))
        write_buffer.write_unsigned_byte(byte_count, logical_name="byteCount")

        # Array Field (items)
        write_buffer.write_complex_array(self.items, logical_name="items")

        write_buffer.pop_context("ModbusPDUReadFileRecordRequest")

    def length_in_bytes(self) -> int:
        return int(math.ceil(float(self.get_length_in_bits() / 8.0)))

    def get_length_in_bits(self) -> int:
        length_in_bits: int = super().get_length_in_bits()
        _value: ModbusPDUReadFileRecordRequest = self

        # Implicit Field (byteCount)
        length_in_bits += 8

        # Array field
        if self.items != None:
            for element in self.items:
                length_in_bits += element.get_length_in_bits()

        return length_in_bits

    @staticmethod
    def static_parse_builder(read_buffer: ReadBuffer, response: bool):
        read_buffer.push_context("ModbusPDUReadFileRecordRequest")

        byte_count: int = read_implicit_field("byteCount", read_unsigned_short)

        self.items = read_length_array_field(
            "items",
            DataReaderComplexDefault(
                ModbusPDUReadFileRecordRequestItem.static_parse(read_buffer),
                read_buffer,
            ),
            byte_count,
        )

        read_buffer.pop_context("ModbusPDUReadFileRecordRequest")
        # Create the instance
        return ModbusPDUReadFileRecordRequestBuilder(items)

    def equals(self, o: object) -> bool:
        if self == o:
            return True

        if not isinstance(o, ModbusPDUReadFileRecordRequest):
            return False

        that: ModbusPDUReadFileRecordRequest = ModbusPDUReadFileRecordRequest(o)
        return (self.items == that.items) and super().equals(that) and True

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
class ModbusPDUReadFileRecordRequestBuilder(ModbusPDUBuilder):
    items: List[ModbusPDUReadFileRecordRequestItem]

    def __post_init__(self):
        pass

    def build(
        self,
    ) -> ModbusPDUReadFileRecordRequest:
        modbus_pdu_read_file_record_request: ModbusPDUReadFileRecordRequest = (
            ModbusPDUReadFileRecordRequest(self.items)
        )
        return modbus_pdu_read_file_record_request
