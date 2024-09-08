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
from plc4py.protocols.umas.readwrite.ModbusPDU import ModbusPDU
from plc4py.protocols.umas.readwrite.UmasPDUItem import UmasPDUItem
from plc4py.spi.generation.ReadBuffer import ReadBuffer
from plc4py.spi.generation.WriteBuffer import WriteBuffer
from typing import ClassVar
import math


@dataclass
class UmasPDU(ModbusPDU):
    item: UmasPDUItem
    # Arguments.
    umas_request_function_key: int
    byte_length: int
    # Accessors for discriminator values.
    error_flag: ClassVar[bool] = False
    function_flag: ClassVar[int] = 0x5A

    def serialize_modbus_pdu_child(self, write_buffer: WriteBuffer):
        write_buffer.push_context("UmasPDU")

        # Simple Field (item)
        write_buffer.write_serializable(self.item, logical_name="item")

        write_buffer.pop_context("UmasPDU")

    def length_in_bytes(self) -> int:
        return int(math.ceil(float(self.length_in_bits() / 8.0)))

    def length_in_bits(self) -> int:
        length_in_bits: int = super().length_in_bits()
        _value: UmasPDU = self

        # Simple field (item)
        length_in_bits += self.item.length_in_bits()

        return length_in_bits

    @staticmethod
    def static_parse_builder(
        read_buffer: ReadBuffer, umas_request_function_key: int, byte_length: int
    ):
        read_buffer.push_context("UmasPDU")

        if isinstance(umas_request_function_key, str):
            umas_request_function_key = int(umas_request_function_key)
        if isinstance(byte_length, str):
            byte_length = int(byte_length)

        item: UmasPDUItem = read_buffer.read_complex(
            read_function=UmasPDUItem.static_parse,
            logical_name="item",
            umas_request_function_key=umas_request_function_key,
            byte_length=byte_length,
        )

        read_buffer.pop_context("UmasPDU")
        # Create the instance
        return UmasPDUBuilder(item)

    def equals(self, o: object) -> bool:
        if self == o:
            return True

        if not isinstance(o, UmasPDU):
            return False

        that: UmasPDU = UmasPDU(o)
        return (self.item == that.item) and super().equals(that) and True

    def hash_code(self) -> int:
        return hash(self)

    def __str__(self) -> str:
        # TODO:- Implement a generic python object to probably json convertor or something.
        return ""


@dataclass
class UmasPDUBuilder:
    item: UmasPDUItem

    def build(
        self,
        umas_request_function_key: int,
        byte_length: int,
    ) -> UmasPDU:
        umas_pdu: UmasPDU = UmasPDU(umas_request_function_key, byte_length, self.item)
        return umas_pdu
