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
from plc4py.protocols.umas.readwrite.UmasPDUItem import UmasPDUItem
from plc4py.protocols.umas.readwrite.UmasUnlocatedVariableReference import (
    UmasUnlocatedVariableReference,
)
from plc4py.spi.generation.ReadBuffer import ReadBuffer
from plc4py.spi.generation.WriteBuffer import WriteBuffer
from plc4py.utils.GenericTypes import ByteOrder
from typing import Any
from typing import ClassVar
from typing import List
import math


@dataclass
class UmasPDUReadUnlocatedVariableNamesResponse(UmasPDUItem):
    range: int
    no_of_records: int
    no_of_records_null: int
    records: List[UmasUnlocatedVariableReference]
    # Accessors for discriminator values.
    umas_function_key: ClassVar[int] = 0xFE
    umas_request_function_key: ClassVar[int] = 0x26

    def serialize_umas_pdu_item_child(self, write_buffer: WriteBuffer):
        write_buffer.push_context("UmasPDUReadUnlocatedVariableNamesResponse")

        # Simple Field (range)
        write_buffer.write_unsigned_int(self.range, bit_length=32, logical_name="range")

        # Simple Field (noOfRecords)
        write_buffer.write_unsigned_short(
            self.no_of_records, bit_length=16, logical_name="noOfRecords"
        )

        # Simple Field (noOfRecordsNull)
        write_buffer.write_unsigned_byte(
            self.no_of_records_null, bit_length=8, logical_name="noOfRecordsNull"
        )

        # Array Field (records)
        write_buffer.write_complex_array(self.records, logical_name="records")

        write_buffer.pop_context("UmasPDUReadUnlocatedVariableNamesResponse")

    def length_in_bytes(self) -> int:
        return int(math.ceil(float(self.length_in_bits() / 8.0)))

    def length_in_bits(self) -> int:
        length_in_bits: int = super().length_in_bits()
        _value: UmasPDUReadUnlocatedVariableNamesResponse = self

        # Simple field (range)
        length_in_bits += 32

        # Simple field (noOfRecords)
        length_in_bits += 16

        # Simple field (noOfRecordsNull)
        length_in_bits += 8

        # Array field
        if self.records is not None:
            for element in self.records:
                length_in_bits += element.length_in_bits()

        return length_in_bits

    @staticmethod
    def static_parse_builder(read_buffer: ReadBuffer, umas_request_function_key: int):
        read_buffer.push_context("UmasPDUReadUnlocatedVariableNamesResponse")

        range: int = read_buffer.read_unsigned_int(
            logical_name="range",
            bit_length=32,
            byte_order=ByteOrder.LITTLE_ENDIAN,
            umas_request_function_key=umas_request_function_key,
        )

        no_of_records: int = read_buffer.read_unsigned_short(
            logical_name="noOfRecords",
            bit_length=16,
            byte_order=ByteOrder.LITTLE_ENDIAN,
            umas_request_function_key=umas_request_function_key,
        )

        no_of_records_null: int = read_buffer.read_unsigned_byte(
            logical_name="noOfRecordsNull",
            bit_length=8,
            byte_order=ByteOrder.LITTLE_ENDIAN,
            umas_request_function_key=umas_request_function_key,
        )

        records: List[Any] = read_buffer.read_array_field(
            logical_name="records",
            read_function=UmasUnlocatedVariableReference.static_parse,
            count=no_of_records,
            byte_order=ByteOrder.LITTLE_ENDIAN,
            umas_request_function_key=umas_request_function_key,
        )

        read_buffer.pop_context("UmasPDUReadUnlocatedVariableNamesResponse")
        # Create the instance
        return UmasPDUReadUnlocatedVariableNamesResponseBuilder(
            range, no_of_records, no_of_records_null, records
        )

    def equals(self, o: object) -> bool:
        if self == o:
            return True

        if not isinstance(o, UmasPDUReadUnlocatedVariableNamesResponse):
            return False

        that: UmasPDUReadUnlocatedVariableNamesResponse = (
            UmasPDUReadUnlocatedVariableNamesResponse(o)
        )
        return (
            (self.range == that.range)
            and (self.no_of_records == that.no_of_records)
            and (self.no_of_records_null == that.no_of_records_null)
            and (self.records == that.records)
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
class UmasPDUReadUnlocatedVariableNamesResponseBuilder:
    range: int
    no_of_records: int
    no_of_records_null: int
    records: List[UmasUnlocatedVariableReference]

    def build(self, pairing_key) -> UmasPDUReadUnlocatedVariableNamesResponse:
        umas_pdu_read_unlocated_variable_names_response: UmasPDUReadUnlocatedVariableNamesResponse = UmasPDUReadUnlocatedVariableNamesResponse(
            pairing_key,
            self.range,
            self.no_of_records,
            self.no_of_records_null,
            self.records,
        )
        return umas_pdu_read_unlocated_variable_names_response
