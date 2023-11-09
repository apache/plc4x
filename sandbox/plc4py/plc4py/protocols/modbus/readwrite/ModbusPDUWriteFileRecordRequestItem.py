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
from typing import Any
from typing import List
import math
    
@dataclass
class ModbusPDUWriteFileRecordRequestItem():
    reference_type: int
    file_number: int
    record_number: int
    record_data: List[int]



    def serialize(self, write_buffer: WriteBuffer):
        write_buffer.push_context("ModbusPDUWriteFileRecordRequestItem")

        # Simple Field (referenceType)
        write_buffer.write_unsigned_byte(self.reference_type, logical_name="referenceType")

        # Simple Field (fileNumber)
        write_buffer.write_unsigned_short(self.file_number, logical_name="fileNumber")

        # Simple Field (recordNumber)
        write_buffer.write_unsigned_short(self.record_number, logical_name="recordNumber")

        # Implicit Field (record_length) (Used for parsing, but its value is not stored as it's implicitly given by the objects content)
        record_length: int = (int(len(self.record_data))/ int(2))
        write_buffer.write_unsigned_short(record_length, logical_name="recordLength")

        # Array Field (recordData)
        write_buffer.write_byte_array(self.record_data, logical_name="recordData")

        write_buffer.pop_context("ModbusPDUWriteFileRecordRequestItem")


    def length_in_bytes(self) -> int:
        return int(math.ceil(float(self.length_in_bits() / 8.0)))

    def length_in_bits(self) -> int:
        length_in_bits: int = 0
        _value: ModbusPDUWriteFileRecordRequestItem = self

        # Simple field (referenceType)
        length_in_bits += 8

        # Simple field (fileNumber)
        length_in_bits += 16

        # Simple field (recordNumber)
        length_in_bits += 16

        # Implicit Field (recordLength)
        length_in_bits += 16

        # Array field
        if self.record_data is not None:
            length_in_bits += 8 * len(self.record_data)


        return length_in_bits


    @staticmethod
    def static_parse(read_buffer: ReadBuffer, **kwargs):
        return ModbusPDUWriteFileRecordRequestItem.static_parse_context(read_buffer)


    @staticmethod
    def static_parse_context(read_buffer: ReadBuffer):
        read_buffer.push_context("ModbusPDUWriteFileRecordRequestItem")

        reference_type: int = read_buffer.read_unsigned_short(logical_name="referenceType")  

        file_number: int = read_buffer.read_unsigned_int(logical_name="fileNumber")  

        record_number: int = read_buffer.read_unsigned_int(logical_name="recordNumber")  

        record_length: int = read_buffer.read_unsigned_int(logical_name="recordLength")

        record_data: List[Any] = read_buffer.read_array_field(logical_name="recordData", read_function=read_buffer.read_byte, count=record_length* int(2))

        read_buffer.pop_context("ModbusPDUWriteFileRecordRequestItem")
        # Create the instance
        _modbus_pdu_write_file_record_request_item: ModbusPDUWriteFileRecordRequestItem = ModbusPDUWriteFileRecordRequestItem(reference_type, file_number, record_number, record_data )
        return _modbus_pdu_write_file_record_request_item


    def equals(self, o: object) -> bool:
        if self == o:
            return True

        if not isinstance(o, ModbusPDUWriteFileRecordRequestItem):
            return False

        that: ModbusPDUWriteFileRecordRequestItem = ModbusPDUWriteFileRecordRequestItem(o)
        return (self.reference_type == that.reference_type) and (self.file_number == that.file_number) and (self.record_number == that.record_number) and (self.record_data == that.record_data) and True

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




