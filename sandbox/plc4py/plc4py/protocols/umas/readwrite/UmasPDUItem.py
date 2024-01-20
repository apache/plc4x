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
import math


@dataclass
class UmasPDUItem(ABC, PlcMessage):
    pairing_key: int

    # Abstract accessors for discriminator values.
    @property
    def umas_function_key(self) -> int:
        pass

    @property
    def umas_request_function_key(self) -> int:
        pass

    @abstractmethod
    def serialize_umas_pdu_item_child(self, write_buffer: WriteBuffer) -> None:
        pass

    def serialize(self, write_buffer: WriteBuffer):
        write_buffer.push_context("UmasPDUItem")

        # Simple Field (pairingKey)
        write_buffer.write_unsigned_byte(self.pairing_key, logical_name="pairingKey")

        # Discriminator Field (umasFunctionKey) (Used as input to a switch field)
        write_buffer.write_unsigned_byte(
            self.umas_function_key,
            logical_name="umasFunctionKey",
            bit_length=8,
        )

        # Switch field (Serialize the sub-type)
        self.serialize_umas_pdu_item_child(write_buffer)

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
                "Wrong number of arguments, expected 1, but got None"
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

        return UmasPDUItem.static_parse_context(read_buffer, umas_request_function_key)

    @staticmethod
    def static_parse_context(read_buffer: ReadBuffer, umas_request_function_key: int):
        read_buffer.push_context("UmasPDUItem")

        pairing_key: int = read_buffer.read_unsigned_byte(
            logical_name="pairingKey",
            bit_length=8,
            umas_request_function_key=umas_request_function_key,
        )

        umas_function_key: int = read_buffer.read_unsigned_byte(
            logical_name="umasFunctionKey",
            bit_length=8,
            umas_request_function_key=umas_request_function_key,
        )

        # Switch Field (Depending on the discriminator values, passes the instantiation to a sub-type)
        builder: UmasPDUItemBuilder = None
        from plc4py.protocols.umas.readwrite.UmasInitCommsRequest import (
            UmasInitCommsRequest,
        )

        if umas_function_key == int(0x01):
            builder = UmasInitCommsRequest.static_parse_builder(
                read_buffer, umas_request_function_key
            )
        from plc4py.protocols.umas.readwrite.UmasPDUPlcIdentRequest import (
            UmasPDUPlcIdentRequest,
        )

        if umas_function_key == int(0x02):
            builder = UmasPDUPlcIdentRequest.static_parse_builder(
                read_buffer, umas_request_function_key
            )
        from plc4py.protocols.umas.readwrite.UmasInitCommsResponse import (
            UmasInitCommsResponse,
        )

        if umas_function_key == int(0xFE) and umas_request_function_key == int(0x01):
            builder = UmasInitCommsResponse.static_parse_builder(
                read_buffer, umas_request_function_key
            )
        from plc4py.protocols.umas.readwrite.UmasPDUPlcIdentResponse import (
            UmasPDUPlcIdentResponse,
        )

        if umas_function_key == int(0xFE) and umas_request_function_key == int(0x02):
            builder = UmasPDUPlcIdentResponse.static_parse_builder(
                read_buffer, umas_request_function_key
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
        _umas_pdu_item: UmasPDUItem = builder.build(pairing_key)
        return _umas_pdu_item

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
        pass
        # write_buffer_box_based: WriteBufferBoxBased = WriteBufferBoxBased(True, True)
        # try:
        #    write_buffer_box_based.writeSerializable(self)
        # except SerializationException as e:
        #    raise PlcRuntimeException(e)

        # return "\n" + str(write_buffer_box_based.get_box()) + "\n"


@dataclass
class UmasPDUItemBuilder:
    pairing_key: int

    def build(
        self,
    ) -> UmasPDUItem:
        pass
