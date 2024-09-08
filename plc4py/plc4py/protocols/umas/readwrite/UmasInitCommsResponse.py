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
class UmasInitCommsResponse(UmasPDUItem):
    max_frame_size: int
    firmware_version: int
    not_sure: int
    internal_code: int
    hostname_length: int
    hostname: str
    # Arguments.
    byte_length: int
    # Accessors for discriminator values.
    umas_function_key: ClassVar[int] = 0xFE
    umas_request_function_key: ClassVar[int] = 0x01

    def serialize_umas_pduitem_child(self, write_buffer: WriteBuffer):
        write_buffer.push_context("UmasInitCommsResponse")

        # Simple Field (maxFrameSize)
        write_buffer.write_unsigned_short(
            self.max_frame_size, bit_length=16, logical_name="maxFrameSize"
        )

        # Simple Field (firmwareVersion)
        write_buffer.write_unsigned_short(
            self.firmware_version, bit_length=16, logical_name="firmwareVersion"
        )

        # Simple Field (notSure)
        write_buffer.write_unsigned_int(
            self.not_sure, bit_length=32, logical_name="notSure"
        )

        # Simple Field (internalCode)
        write_buffer.write_unsigned_int(
            self.internal_code, bit_length=32, logical_name="internalCode"
        )

        # Simple Field (hostnameLength)
        write_buffer.write_unsigned_byte(
            self.hostname_length, bit_length=8, logical_name="hostnameLength"
        )

        # Simple Field (hostname)
        write_buffer.write_str(self.hostname, bit_length=-1, logical_name="hostname")

        write_buffer.pop_context("UmasInitCommsResponse")

    def length_in_bytes(self) -> int:
        return int(math.ceil(float(self.length_in_bits() / 8.0)))

    def length_in_bits(self) -> int:
        length_in_bits: int = super().length_in_bits()
        _value: UmasInitCommsResponse = self

        # Simple field (maxFrameSize)
        length_in_bits += 16

        # Simple field (firmwareVersion)
        length_in_bits += 16

        # Simple field (notSure)
        length_in_bits += 32

        # Simple field (internalCode)
        length_in_bits += 32

        # Simple field (hostnameLength)
        length_in_bits += 8

        # Simple field (hostname)
        length_in_bits += self.hostname_length * int(8)

        return length_in_bits

    @staticmethod
    def static_parse_builder(
        read_buffer: ReadBuffer, umas_request_function_key: int, byte_length: int
    ):
        read_buffer.push_context("UmasInitCommsResponse")

        if isinstance(umas_request_function_key, str):
            umas_request_function_key = int(umas_request_function_key)
        if isinstance(byte_length, str):
            byte_length = int(byte_length)

        max_frame_size: int = read_buffer.read_unsigned_short(
            logical_name="max_frame_size",
            bit_length=16,
            byte_order=ByteOrder.LITTLE_ENDIAN,
            umas_request_function_key=umas_request_function_key,
            byte_length=byte_length,
        )

        firmware_version: int = read_buffer.read_unsigned_short(
            logical_name="firmware_version",
            bit_length=16,
            byte_order=ByteOrder.LITTLE_ENDIAN,
            umas_request_function_key=umas_request_function_key,
            byte_length=byte_length,
        )

        not_sure: int = read_buffer.read_unsigned_int(
            logical_name="not_sure",
            bit_length=32,
            byte_order=ByteOrder.LITTLE_ENDIAN,
            umas_request_function_key=umas_request_function_key,
            byte_length=byte_length,
        )

        internal_code: int = read_buffer.read_unsigned_int(
            logical_name="internal_code",
            bit_length=32,
            byte_order=ByteOrder.LITTLE_ENDIAN,
            umas_request_function_key=umas_request_function_key,
            byte_length=byte_length,
        )

        hostname_length: int = read_buffer.read_unsigned_byte(
            logical_name="hostname_length",
            bit_length=8,
            byte_order=ByteOrder.LITTLE_ENDIAN,
            umas_request_function_key=umas_request_function_key,
            byte_length=byte_length,
        )

        hostname: str = read_buffer.read_str(
            logical_name="hostname",
            bit_length=hostname_length * int(8),
            byte_order=ByteOrder.LITTLE_ENDIAN,
            umas_request_function_key=umas_request_function_key,
            byte_length=byte_length,
        )

        read_buffer.pop_context("UmasInitCommsResponse")
        # Create the instance
        return UmasInitCommsResponseBuilder(
            max_frame_size,
            firmware_version,
            not_sure,
            internal_code,
            hostname_length,
            hostname,
        )

    def equals(self, o: object) -> bool:
        if self == o:
            return True

        if not isinstance(o, UmasInitCommsResponse):
            return False

        that: UmasInitCommsResponse = UmasInitCommsResponse(o)
        return (
            (self.max_frame_size == that.max_frame_size)
            and (self.firmware_version == that.firmware_version)
            and (self.not_sure == that.not_sure)
            and (self.internal_code == that.internal_code)
            and (self.hostname_length == that.hostname_length)
            and (self.hostname == that.hostname)
            and super().equals(that)
            and True
        )

    def hash_code(self) -> int:
        return hash(self)

    def __str__(self) -> str:
        # TODO:- Implement a generic python object to probably json convertor or something.
        return ""


@dataclass
class UmasInitCommsResponseBuilder:
    max_frame_size: int
    firmware_version: int
    not_sure: int
    internal_code: int
    hostname_length: int
    hostname: str

    def build(self, byte_length: int, pairing_key) -> UmasInitCommsResponse:
        umas_init_comms_response: UmasInitCommsResponse = UmasInitCommsResponse(
            byte_length,
            pairing_key,
            self.max_frame_size,
            self.firmware_version,
            self.not_sure,
            self.internal_code,
            self.hostname_length,
            self.hostname,
        )
        return umas_init_comms_response
