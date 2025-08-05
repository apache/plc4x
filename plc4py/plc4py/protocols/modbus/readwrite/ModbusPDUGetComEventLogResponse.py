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
from typing import Any
from typing import ClassVar
from typing import List
import math


@dataclass
class ModbusPDUGetComEventLogResponse(ModbusPDU):
    status: int
    event_count: int
    message_count: int
    events: List[int]
    # Accessors for discriminator values.
    error_flag: ClassVar[bool] = False
    function_flag: ClassVar[int] = 0x0C
    response: ClassVar[bool] = True

    def serialize_modbus_pdu_child(self, write_buffer: WriteBuffer):
        write_buffer.push_context("ModbusPDUGetComEventLogResponse")

        # Implicit Field (byte_count) (Used for parsing, but its value is not stored as it's implicitly given by the objects content)
        byte_count: int = int(len(self.events)) + int(6)
        write_buffer.write_unsigned_byte(byte_count, logical_name="byte_count")

        # Simple Field (status)
        write_buffer.write_unsigned_short(
            self.status, bit_length=16, logical_name="status"
        )

        # Simple Field (eventCount)
        write_buffer.write_unsigned_short(
            self.event_count, bit_length=16, logical_name="eventCount"
        )

        # Simple Field (messageCount)
        write_buffer.write_unsigned_short(
            self.message_count, bit_length=16, logical_name="messageCount"
        )

        # Array Field (events)
        write_buffer.write_byte_array(self.events, logical_name="events")

        write_buffer.pop_context("ModbusPDUGetComEventLogResponse")

    def length_in_bytes(self) -> int:
        return int(math.ceil(float(self.length_in_bits() / 8.0)))

    def length_in_bits(self) -> int:
        length_in_bits: int = super().length_in_bits()
        _value: ModbusPDUGetComEventLogResponse = self

        # Implicit Field (byteCount)
        length_in_bits += 8

        # Simple field (status)
        length_in_bits += 16

        # Simple field (eventCount)
        length_in_bits += 16

        # Simple field (messageCount)
        length_in_bits += 16

        # Array field
        if self.events is not None:
            length_in_bits += 8 * len(self.events)

        return length_in_bits

    @staticmethod
    def static_parse_builder(read_buffer: ReadBuffer, response: bool):
        read_buffer.push_context("ModbusPDUGetComEventLogResponse")

        if isinstance(response, str):
            response = bool(strtobool(response))

        byte_count: int = read_buffer.read_unsigned_byte(
            logical_name="byte_count", response=response
        )

        status: int = read_buffer.read_unsigned_short(
            logical_name="status", bit_length=16, response=response
        )

        event_count: int = read_buffer.read_unsigned_short(
            logical_name="event_count", bit_length=16, response=response
        )

        message_count: int = read_buffer.read_unsigned_short(
            logical_name="message_count", bit_length=16, response=response
        )

        events: List[Any] = read_buffer.read_array_field(
            logical_name="events",
            read_function=read_buffer.read_byte,
            count=byte_count - int(6),
            response=response,
        )

        read_buffer.pop_context("ModbusPDUGetComEventLogResponse")
        # Create the instance
        return ModbusPDUGetComEventLogResponseBuilder(
            status, event_count, message_count, events
        )

    def equals(self, o: object) -> bool:
        if self == o:
            return True

        if not isinstance(o, ModbusPDUGetComEventLogResponse):
            return False

        that: ModbusPDUGetComEventLogResponse = ModbusPDUGetComEventLogResponse(o)
        return (
            (self.status == that.status)
            and (self.event_count == that.event_count)
            and (self.message_count == that.message_count)
            and (self.events == that.events)
            and super().equals(that)
            and True
        )

    def hash_code(self) -> int:
        return hash(self)

    def __str__(self) -> str:
        # TODO:- Implement a generic python object to probably json convertor or something.
        return ""


@dataclass
class ModbusPDUGetComEventLogResponseBuilder:
    status: int
    event_count: int
    message_count: int
    events: List[int]

    def build(
        self,
    ) -> ModbusPDUGetComEventLogResponse:
        modbus_pduget_com_event_log_response: ModbusPDUGetComEventLogResponse = (
            ModbusPDUGetComEventLogResponse(
                self.status, self.event_count, self.message_count, self.events
            )
        )
        return modbus_pduget_com_event_log_response
