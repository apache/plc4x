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
from plc4py.protocols.modbus.readwrite.ModbusPDU import ModbusPDU
from plc4py.spi.generation.ReadBuffer import ReadBuffer
from plc4py.spi.generation.WriteBuffer import WriteBuffer
from plc4py.utils.ConnectionStringHandling import strtobool
from typing import ClassVar
import math


@dataclass
class ModbusPDUGetComEventCounterResponse(ModbusPDU):
    status: int
    event_count: int
    # Accessors for discriminator values.
    error_flag: ClassVar[bool] = False
    function_flag: ClassVar[int] = 0x0B
    response: ClassVar[bool] = True

    def serialize_modbus_pdu_child(self, write_buffer: WriteBuffer):
        write_buffer.push_context("ModbusPDUGetComEventCounterResponse")

        # Simple Field (status)
        write_buffer.write_unsigned_short(
            self.status, bit_length=16, logical_name="status"
        )

        # Simple Field (eventCount)
        write_buffer.write_unsigned_short(
            self.event_count, bit_length=16, logical_name="eventCount"
        )

        write_buffer.pop_context("ModbusPDUGetComEventCounterResponse")

    def length_in_bytes(self) -> int:
        return int(math.ceil(float(self.length_in_bits() / 8.0)))

    def length_in_bits(self) -> int:
        length_in_bits: int = super().length_in_bits()
        _value: ModbusPDUGetComEventCounterResponse = self

        # Simple field (status)
        length_in_bits += 16

        # Simple field (eventCount)
        length_in_bits += 16

        return length_in_bits

    @staticmethod
    def static_parse_builder(read_buffer: ReadBuffer, response: bool):
        read_buffer.push_context("ModbusPDUGetComEventCounterResponse")

        if isinstance(response, str):
            response = bool(strtobool(response))

        status: int = read_buffer.read_unsigned_short(
            logical_name="status", bit_length=16, response=response
        )

        event_count: int = read_buffer.read_unsigned_short(
            logical_name="event_count", bit_length=16, response=response
        )

        read_buffer.pop_context("ModbusPDUGetComEventCounterResponse")
        # Create the instance
        return ModbusPDUGetComEventCounterResponseBuilder(status, event_count)

    def equals(self, o: object) -> bool:
        if self == o:
            return True

        if not isinstance(o, ModbusPDUGetComEventCounterResponse):
            return False

        that: ModbusPDUGetComEventCounterResponse = ModbusPDUGetComEventCounterResponse(
            o
        )
        return (
            (self.status == that.status)
            and (self.event_count == that.event_count)
            and super().equals(that)
            and True
        )

    def hash_code(self) -> int:
        return hash(self)

    def __str__(self) -> str:
        # TODO:- Implement a generic python object to probably json convertor or something.
        return ""


@dataclass
class ModbusPDUGetComEventCounterResponseBuilder:
    status: int
    event_count: int

    def build(
        self,
    ) -> ModbusPDUGetComEventCounterResponse:
        modbus_pduget_com_event_counter_response: (
            ModbusPDUGetComEventCounterResponse
        ) = ModbusPDUGetComEventCounterResponse(self.status, self.event_count)
        return modbus_pduget_com_event_counter_response
