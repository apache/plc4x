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
from plc4py.spi.generation.ReadBuffer import ReadBuffer
from plc4py.spi.generation.WriteBuffer import WriteBuffer
from plc4py.utils.GenericTypes import ByteOrder
from typing import ClassVar
import math


@dataclass
class UmasPDUProjectInfoRequest(UmasPDUItem):
    subcode: int
    # Arguments.
    byte_length: int
    # Accessors for discriminator values.
    umas_function_key: ClassVar[int] = 0x03
    umas_request_function_key: ClassVar[int] = 0

    def serialize_umas_pduitem_child(self, write_buffer: WriteBuffer):
        write_buffer.push_context("UmasPDUProjectInfoRequest")

        # Simple Field (subcode)
        write_buffer.write_unsigned_byte(
            self.subcode, bit_length=8, logical_name="subcode"
        )

        write_buffer.pop_context("UmasPDUProjectInfoRequest")

    def length_in_bytes(self) -> int:
        return int(math.ceil(float(self.length_in_bits() / 8.0)))

    def length_in_bits(self) -> int:
        length_in_bits: int = super().length_in_bits()
        _value: UmasPDUProjectInfoRequest = self

        # Simple field (subcode)
        length_in_bits += 8

        return length_in_bits

    @staticmethod
    def static_parse_builder(
        read_buffer: ReadBuffer, umas_request_function_key: int, byte_length: int
    ):
        read_buffer.push_context("UmasPDUProjectInfoRequest")

        if isinstance(umas_request_function_key, str):
            umas_request_function_key = int(umas_request_function_key)
        if isinstance(byte_length, str):
            byte_length = int(byte_length)

        subcode: int = read_buffer.read_unsigned_byte(
            logical_name="subcode",
            bit_length=8,
            byte_order=ByteOrder.LITTLE_ENDIAN,
            umas_request_function_key=umas_request_function_key,
            byte_length=byte_length,
        )

        read_buffer.pop_context("UmasPDUProjectInfoRequest")
        # Create the instance
        return UmasPDUProjectInfoRequestBuilder(subcode)

    def equals(self, o: object) -> bool:
        if self == o:
            return True

        if not isinstance(o, UmasPDUProjectInfoRequest):
            return False

        that: UmasPDUProjectInfoRequest = UmasPDUProjectInfoRequest(o)
        return (self.subcode == that.subcode) and super().equals(that) and True

    def hash_code(self) -> int:
        return hash(self)

    def __str__(self) -> str:
        # TODO:- Implement a generic python object to probably json convertor or something.
        return ""


@dataclass
class UmasPDUProjectInfoRequestBuilder:
    subcode: int

    def build(self, byte_length: int, pairing_key) -> UmasPDUProjectInfoRequest:
        umas_pduproject_info_request: UmasPDUProjectInfoRequest = (
            UmasPDUProjectInfoRequest(byte_length, pairing_key, self.subcode)
        )
        return umas_pduproject_info_request
