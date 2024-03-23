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
from plc4py.protocols.umas.readwrite.UmasUDTDefinition import UmasUDTDefinition
from plc4py.spi.generation.ReadBuffer import ReadBuffer
from plc4py.spi.generation.WriteBuffer import WriteBuffer
from typing import Any
from typing import List
import math


@dataclass
class UmasPDUReadUmasUDTDefinitionResponse:
    range: int
    unknown1: int
    no_of_records: int
    records: List[UmasUDTDefinition]

    def serialize(self, write_buffer: WriteBuffer):
        write_buffer.push_context("UmasPDUReadUmasUDTDefinitionResponse")

        # Simple Field (range)
        write_buffer.write_unsigned_byte(self.range, bit_length=8, logical_name="range")

        # Simple Field (unknown1)
        write_buffer.write_unsigned_int(
            self.unknown1, bit_length=32, logical_name="unknown1"
        )

        # Simple Field (noOfRecords)
        write_buffer.write_unsigned_short(
            self.no_of_records, bit_length=16, logical_name="noOfRecords"
        )

        # Array Field (records)
        write_buffer.write_complex_array(self.records, logical_name="records")

        write_buffer.pop_context("UmasPDUReadUmasUDTDefinitionResponse")

    def length_in_bytes(self) -> int:
        return int(math.ceil(float(self.length_in_bits() / 8.0)))

    def length_in_bits(self) -> int:
        length_in_bits: int = 0
        _value: UmasPDUReadUmasUDTDefinitionResponse = self

        # Simple field (range)
        length_in_bits += 8

        # Simple field (unknown1)
        length_in_bits += 32

        # Simple field (noOfRecords)
        length_in_bits += 16

        # Array field
        if self.records is not None:
            for element in self.records:
                length_in_bits += element.length_in_bits()

        return length_in_bits

    @staticmethod
    def static_parse(read_buffer: ReadBuffer, **kwargs):
        return UmasPDUReadUmasUDTDefinitionResponse.static_parse_context(read_buffer)

    @staticmethod
    def static_parse_context(read_buffer: ReadBuffer):
        read_buffer.push_context("UmasPDUReadUmasUDTDefinitionResponse")

        range: int = read_buffer.read_unsigned_byte(logical_name="range", bit_length=8)

        unknown1: int = read_buffer.read_unsigned_int(
            logical_name="unknown1", bit_length=32
        )

        no_of_records: int = read_buffer.read_unsigned_short(
            logical_name="noOfRecords", bit_length=16
        )

        records: List[Any] = read_buffer.read_array_field(
            logical_name="records",
            read_function=UmasUDTDefinition.static_parse,
            count=no_of_records,
        )

        read_buffer.pop_context("UmasPDUReadUmasUDTDefinitionResponse")
        # Create the instance
        _umas_pdu_read_umas_udt_definition_response: (
            UmasPDUReadUmasUDTDefinitionResponse
        ) = UmasPDUReadUmasUDTDefinitionResponse(
            range, unknown1, no_of_records, records
        )
        return _umas_pdu_read_umas_udt_definition_response

    def equals(self, o: object) -> bool:
        if self == o:
            return True

        if not isinstance(o, UmasPDUReadUmasUDTDefinitionResponse):
            return False

        that: UmasPDUReadUmasUDTDefinitionResponse = (
            UmasPDUReadUmasUDTDefinitionResponse(o)
        )
        return (
            (self.range == that.range)
            and (self.unknown1 == that.unknown1)
            and (self.no_of_records == that.no_of_records)
            and (self.records == that.records)
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
