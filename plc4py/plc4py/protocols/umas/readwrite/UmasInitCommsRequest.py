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
class UmasInitCommsRequest(UmasPDUItem):
    unknown_object: int
    # Arguments.
    byte_length: int
    # Accessors for discriminator values.
    umas_function_key: ClassVar[int] = 0x01
    umas_request_function_key: ClassVar[int] = 0

    def serialize_umas_pduitem_child(self, write_buffer: WriteBuffer):
        write_buffer.push_context("UmasInitCommsRequest")

        # Simple Field (unknownObject)
        write_buffer.write_unsigned_byte(
            self.unknown_object, bit_length=8, logical_name="unknownObject"
        )

        write_buffer.pop_context("UmasInitCommsRequest")

    def length_in_bytes(self) -> int:
        return int(math.ceil(float(self.length_in_bits() / 8.0)))

    def length_in_bits(self) -> int:
        length_in_bits: int = super().length_in_bits()
        _value: UmasInitCommsRequest = self

        # Simple field (unknownObject)
        length_in_bits += 8

        return length_in_bits

    @staticmethod
    def static_parse_builder(
        read_buffer: ReadBuffer, umas_request_function_key: int, byte_length: int
    ):
        read_buffer.push_context("UmasInitCommsRequest")

        if isinstance(umas_request_function_key, str):
            umas_request_function_key = int(umas_request_function_key)
        if isinstance(byte_length, str):
            byte_length = int(byte_length)

        unknown_object: int = read_buffer.read_unsigned_byte(
            logical_name="unknown_object",
            bit_length=8,
            byte_order=ByteOrder.LITTLE_ENDIAN,
            umas_request_function_key=umas_request_function_key,
            byte_length=byte_length,
        )

        read_buffer.pop_context("UmasInitCommsRequest")
        # Create the instance
        return UmasInitCommsRequestBuilder(unknown_object)

    def equals(self, o: object) -> bool:
        if self == o:
            return True

        if not isinstance(o, UmasInitCommsRequest):
            return False

        that: UmasInitCommsRequest = UmasInitCommsRequest(o)
        return (
            (self.unknown_object == that.unknown_object)
            and super().equals(that)
            and True
        )

    def hash_code(self) -> int:
        return hash(self)

    def __str__(self) -> str:
        # TODO:- Implement a generic python object to probably json convertor or something.
        return ""


@dataclass
class UmasInitCommsRequestBuilder:
    unknown_object: int

    def build(self, byte_length: int, pairing_key) -> UmasInitCommsRequest:
        umas_init_comms_request: UmasInitCommsRequest = UmasInitCommsRequest(
            byte_length, pairing_key, self.unknown_object
        )
        return umas_init_comms_request
