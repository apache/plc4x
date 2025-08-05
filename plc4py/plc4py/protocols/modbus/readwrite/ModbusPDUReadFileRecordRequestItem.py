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
from plc4py.spi.generation.ReadBuffer import ReadBuffer
from plc4py.spi.generation.WriteBuffer import WriteBuffer
import math


@dataclass
class ModbusPDUReadFileRecordRequestItem:
    reference_type: int
    file_number: int
    record_number: int
    record_length: int

    def serialize(self, write_buffer: WriteBuffer):
        write_buffer.push_context("ModbusPDUReadFileRecordRequestItem")

        # Simple Field (referenceType)
        write_buffer.write_unsigned_byte(
            self.reference_type, bit_length=8, logical_name="referenceType"
        )

        # Simple Field (fileNumber)
        write_buffer.write_unsigned_short(
            self.file_number, bit_length=16, logical_name="fileNumber"
        )

        # Simple Field (recordNumber)
        write_buffer.write_unsigned_short(
            self.record_number, bit_length=16, logical_name="recordNumber"
        )

        # Simple Field (recordLength)
        write_buffer.write_unsigned_short(
            self.record_length, bit_length=16, logical_name="recordLength"
        )

        write_buffer.pop_context("ModbusPDUReadFileRecordRequestItem")

    def length_in_bytes(self) -> int:
        return int(math.ceil(float(self.length_in_bits() / 8.0)))

    def length_in_bits(self) -> int:
        length_in_bits: int = 0
        _value: ModbusPDUReadFileRecordRequestItem = self

        # Simple field (referenceType)
        length_in_bits += 8

        # Simple field (fileNumber)
        length_in_bits += 16

        # Simple field (recordNumber)
        length_in_bits += 16

        # Simple field (recordLength)
        length_in_bits += 16

        return length_in_bits

    @staticmethod
    def static_parse(read_buffer: ReadBuffer, **kwargs):
        return ModbusPDUReadFileRecordRequestItem.static_parse_context(read_buffer)

    @staticmethod
    def static_parse_context(read_buffer: ReadBuffer):
        read_buffer.push_context("ModbusPDUReadFileRecordRequestItem")

        reference_type: int = read_buffer.read_unsigned_byte(
            logical_name="reference_type", bit_length=8
        )

        file_number: int = read_buffer.read_unsigned_short(
            logical_name="file_number", bit_length=16
        )

        record_number: int = read_buffer.read_unsigned_short(
            logical_name="record_number", bit_length=16
        )

        record_length: int = read_buffer.read_unsigned_short(
            logical_name="record_length", bit_length=16
        )

        read_buffer.pop_context("ModbusPDUReadFileRecordRequestItem")
        # Create the instance
        _modbus_pduread_file_record_request_item: ModbusPDUReadFileRecordRequestItem = (
            ModbusPDUReadFileRecordRequestItem(
                reference_type, file_number, record_number, record_length
            )
        )
        return _modbus_pduread_file_record_request_item

    def equals(self, o: object) -> bool:
        if self == o:
            return True

        if not isinstance(o, ModbusPDUReadFileRecordRequestItem):
            return False

        that: ModbusPDUReadFileRecordRequestItem = ModbusPDUReadFileRecordRequestItem(o)
        return (
            (self.reference_type == that.reference_type)
            and (self.file_number == that.file_number)
            and (self.record_number == that.record_number)
            and (self.record_length == that.record_length)
            and True
        )

    def hash_code(self) -> int:
        return hash(self)

    def __str__(self) -> str:
        # TODO:- Implement a generic python object to probably json convertor or something.
        return ""
