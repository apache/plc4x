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

from abc import ABC
from abc import abstractmethod
from plc4py.api.exceptions.exceptions import ParseException
from plc4py.api.exceptions.exceptions import PlcRuntimeException
from plc4py.api.exceptions.exceptions import SerializationException
from plc4py.api.messages.PlcMessage import PlcMessage
from plc4py.spi.generation.ReadBuffer import ReadBuffer
from plc4py.spi.generation.WriteBuffer import WriteBuffer
from plc4py.utils.GenericTypes import ByteOrder
import math


@dataclass
class UmasPDUItem(ABC, PlcMessage):
    pairing_key: int
    # Arguments.
    byte_length: int

    # Abstract accessors for discriminator values.
    @property
    def umas_function_key(self) -> int:
        pass

    @property
    def umas_request_function_key(self) -> int:
        pass

    @abstractmethod
    def serialize_umas_pduitem_child(self, write_buffer: WriteBuffer) -> None:
        pass

    def serialize(self, write_buffer: WriteBuffer):
        write_buffer.push_context("UmasPDUItem")

        # Simple Field (pairingKey)
        write_buffer.write_unsigned_byte(
            self.pairing_key, bit_length=8, logical_name="pairingKey"
        )

        # Discriminator Field (umasFunctionKey) (Used as input to a switch field)
        write_buffer.write_unsigned_byte(
            self.umas_function_key,
            logical_name="umas_function_key",
            bit_length=8,
            byte_order=ByteOrder.LITTLE_ENDIAN,
        )

        # Switch field (Serialize the sub-type)
        self.serialize_umas_pduitem_child(write_buffer)

        write_buffer.pop_context("UmasPDUItem")

    def length_in_bytes(self) -> int:
        return int(math.ceil(float(self.length_in_bits() / 8.0)))

    def length_in_bits(self) -> int:
        length_in_bits: int = 0
        _value: UmasPDUItem = self

        # Simple field (pairingKey)
        length_in_bits += 8

        # Discriminator Field (umasFunctionKey)
        length_in_bits += 8

        # Length of subtype elements will be added by sub-type...

        return length_in_bits

    @staticmethod
    def static_parse(read_buffer: ReadBuffer, **kwargs):

        if kwargs is None:
            raise PlcRuntimeException(
                "Wrong number of arguments, expected 2, but got None"
            )

        umas_request_function_key: int = 0
        if isinstance(kwargs.get("umas_request_function_key"), int):
            umas_request_function_key = int(kwargs.get("umas_request_function_key"))
        elif isinstance(kwargs.get("umas_request_function_key"), str):
            umas_request_function_key = int(
                str(kwargs.get("umas_request_function_key"))
            )
        else:
            raise PlcRuntimeException(
                "Argument 0 expected to be of type int or a string which is parseable but was "
                + kwargs.get("umas_request_function_key").getClass().getName()
            )

        byte_length: int = 0
        if isinstance(kwargs.get("byte_length"), int):
            byte_length = int(kwargs.get("byte_length"))
        elif isinstance(kwargs.get("byte_length"), str):
            byte_length = int(str(kwargs.get("byte_length")))
        else:
            raise PlcRuntimeException(
                "Argument 1 expected to be of type int or a string which is parseable but was "
                + kwargs.get("byte_length").getClass().getName()
            )

        return UmasPDUItem.static_parse_context(
            read_buffer, umas_request_function_key, byte_length
        )

    @staticmethod
    def static_parse_context(
        read_buffer: ReadBuffer, umas_request_function_key: int, byte_length: int
    ):
        read_buffer.push_context("UmasPDUItem")

        if isinstance(umas_request_function_key, str):
            umas_request_function_key = int(umas_request_function_key)
        if isinstance(byte_length, str):
            byte_length = int(byte_length)

        pairing_key: int = read_buffer.read_unsigned_byte(
            logical_name="pairing_key",
            bit_length=8,
            byte_order=ByteOrder.LITTLE_ENDIAN,
            umas_request_function_key=umas_request_function_key,
            byte_length=byte_length,
        )

        umas_function_key: int = read_buffer.read_unsigned_byte(
            logical_name="umas_function_key",
            bit_length=8,
            byte_order=ByteOrder.LITTLE_ENDIAN,
            umas_request_function_key=umas_request_function_key,
            byte_length=byte_length,
        )

        # Switch Field (Depending on the discriminator values, passes the instantiation to a sub-type)
        builder: UmasPDUItemBuilder = None
        from plc4py.protocols.umas.readwrite.UmasInitCommsRequest import (
            UmasInitCommsRequest,
        )

        if umas_function_key == int(0x01):

            builder = UmasInitCommsRequest.static_parse_builder(
                read_buffer, umas_request_function_key, byte_length
            )
        from plc4py.protocols.umas.readwrite.UmasPDUPlcIdentRequest import (
            UmasPDUPlcIdentRequest,
        )

        if umas_function_key == int(0x02):

            builder = UmasPDUPlcIdentRequest.static_parse_builder(
                read_buffer, umas_request_function_key, byte_length
            )
        from plc4py.protocols.umas.readwrite.UmasPDUProjectInfoRequest import (
            UmasPDUProjectInfoRequest,
        )

        if umas_function_key == int(0x03):

            builder = UmasPDUProjectInfoRequest.static_parse_builder(
                read_buffer, umas_request_function_key, byte_length
            )
        from plc4py.protocols.umas.readwrite.UmasPDUPlcStatusRequest import (
            UmasPDUPlcStatusRequest,
        )

        if umas_function_key == int(0x04):

            builder = UmasPDUPlcStatusRequest.static_parse_builder(
                read_buffer, umas_request_function_key, byte_length
            )
        from plc4py.protocols.umas.readwrite.UmasPDUReadMemoryBlockRequest import (
            UmasPDUReadMemoryBlockRequest,
        )

        if umas_function_key == int(0x20):

            builder = UmasPDUReadMemoryBlockRequest.static_parse_builder(
                read_buffer, umas_request_function_key, byte_length
            )
        from plc4py.protocols.umas.readwrite.UmasPDUReadVariableRequest import (
            UmasPDUReadVariableRequest,
        )

        if umas_function_key == int(0x22):

            builder = UmasPDUReadVariableRequest.static_parse_builder(
                read_buffer, umas_request_function_key, byte_length
            )
        from plc4py.protocols.umas.readwrite.UmasPDUWriteVariableRequest import (
            UmasPDUWriteVariableRequest,
        )

        if umas_function_key == int(0x23):

            builder = UmasPDUWriteVariableRequest.static_parse_builder(
                read_buffer, umas_request_function_key, byte_length
            )
        from plc4py.protocols.umas.readwrite.UmasPDUReadUnlocatedVariableNamesRequest import (
            UmasPDUReadUnlocatedVariableNamesRequest,
        )

        if umas_function_key == int(0x26):

            builder = UmasPDUReadUnlocatedVariableNamesRequest.static_parse_builder(
                read_buffer, umas_request_function_key, byte_length
            )
        from plc4py.protocols.umas.readwrite.UmasPDUErrorResponse import (
            UmasPDUErrorResponse,
        )

        if umas_function_key == int(0xFD):

            builder = UmasPDUErrorResponse.static_parse_builder(
                read_buffer, umas_request_function_key, byte_length
            )
        from plc4py.protocols.umas.readwrite.UmasInitCommsResponse import (
            UmasInitCommsResponse,
        )

        if umas_function_key == int(0xFE) and umas_request_function_key == int(0x01):

            builder = UmasInitCommsResponse.static_parse_builder(
                read_buffer, umas_request_function_key, byte_length
            )
        from plc4py.protocols.umas.readwrite.UmasPDUPlcIdentResponse import (
            UmasPDUPlcIdentResponse,
        )

        if umas_function_key == int(0xFE) and umas_request_function_key == int(0x02):

            builder = UmasPDUPlcIdentResponse.static_parse_builder(
                read_buffer, umas_request_function_key, byte_length
            )
        from plc4py.protocols.umas.readwrite.UmasPDUPlcStatusResponse import (
            UmasPDUPlcStatusResponse,
        )

        if umas_function_key == int(0xFE) and umas_request_function_key == int(0x04):

            builder = UmasPDUPlcStatusResponse.static_parse_builder(
                read_buffer, umas_request_function_key, byte_length
            )
        from plc4py.protocols.umas.readwrite.UmasPDUReadMemoryBlockResponse import (
            UmasPDUReadMemoryBlockResponse,
        )

        if umas_function_key == int(0xFE) and umas_request_function_key == int(0x20):

            builder = UmasPDUReadMemoryBlockResponse.static_parse_builder(
                read_buffer, umas_request_function_key, byte_length
            )
        from plc4py.protocols.umas.readwrite.UmasPDUReadVariableResponse import (
            UmasPDUReadVariableResponse,
        )

        if umas_function_key == int(0xFE) and umas_request_function_key == int(0x22):

            builder = UmasPDUReadVariableResponse.static_parse_builder(
                read_buffer, umas_request_function_key, byte_length
            )
        from plc4py.protocols.umas.readwrite.UmasPDUWriteVariableResponse import (
            UmasPDUWriteVariableResponse,
        )

        if umas_function_key == int(0xFE) and umas_request_function_key == int(0x23):

            builder = UmasPDUWriteVariableResponse.static_parse_builder(
                read_buffer, umas_request_function_key, byte_length
            )
        from plc4py.protocols.umas.readwrite.UmasPDUReadUnlocatedVariableResponse import (
            UmasPDUReadUnlocatedVariableResponse,
        )

        if umas_function_key == int(0xFE) and umas_request_function_key == int(0x26):

            builder = UmasPDUReadUnlocatedVariableResponse.static_parse_builder(
                read_buffer, umas_request_function_key, byte_length
            )
        if builder is None:
            raise ParseException(
                "Unsupported case for discriminated type"
                + " parameters ["
                + "umasFunctionKey="
                + str(umas_function_key)
                + " "
                + "umasRequestFunctionKey="
                + str(umas_request_function_key)
                + "]"
            )

        read_buffer.pop_context("UmasPDUItem")
        # Create the instance
        _umas_pduitem: UmasPDUItem = builder.build(pairing_key, byte_length)
        return _umas_pduitem

    def equals(self, o: object) -> bool:
        if self == o:
            return True

        if not isinstance(o, UmasPDUItem):
            return False

        that: UmasPDUItem = UmasPDUItem(o)
        return (self.pairing_key == that.pairing_key) and True

    def hash_code(self) -> int:
        return hash(self)

    def __str__(self) -> str:
        # TODO:- Implement a generic python object to probably json convertor or something.
        return ""


@dataclass
class UmasPDUItemBuilder:
    pairing_key: int

    def build(self, byte_length: int) -> UmasPDUItem:
        pass
