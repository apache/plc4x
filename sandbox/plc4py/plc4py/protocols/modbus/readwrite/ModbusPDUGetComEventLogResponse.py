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

from plc4py.api.messages.PlcMessage import PlcMessage
from plc4py.protocols.modbus.readwrite.ModbusPDU import ModbusPDU
from plc4py.protocols.modbus.readwrite.ModbusPDU import ModbusPDUBuilder
from plc4py.spi.generation.ReadBuffer import ReadBuffer
from plc4py.spi.generation.WriteBuffer import WriteBuffer
from typing import List
import math


@dataclass
class ModbusPDUGetComEventLogResponse(PlcMessage, ModbusPDU):
    status: int
    event_count: int
    message_count: int
    events: List[int]
    # Accessors for discriminator values.
    error_flag: bool = False
    function_flag: int = 0x0C
    response: bool = True

    def __post_init__(self):
        super().__init__()

    def serialize_modbus_pdu_child(self, write_buffer: WriteBuffer):
        write_buffer.push_context("ModbusPDUGetComEventLogResponse")

        # Implicit Field (byte_count) (Used for parsing, but its value is not stored as it's implicitly given by the objects content)
        byte_count: int = int(len(self.events)) + int(6)
        write_buffer.write_unsigned_byte(byte_count, logical_name="byteCount")

        # Simple Field (status)
        write_buffer.write_unsigned_short(self.status, logical_name="status")

        # Simple Field (eventCount)
        write_buffer.write_unsigned_short(self.event_count, logical_name="eventCount")

        # Simple Field (messageCount)
        write_buffer.write_unsigned_short(
            self.message_count, logical_name="messageCount"
        )

        # Array Field (events)
        write_buffer.write_byte_array(self.events, logical_name="events")

        write_buffer.pop_context("ModbusPDUGetComEventLogResponse")

    def length_in_bytes(self) -> int:
        return int(math.ceil(float(self.get_length_in_bits() / 8.0)))

    def get_length_in_bits(self) -> int:
        length_in_bits: int = super().get_length_in_bits()
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
        if self.events != None:
            length_in_bits += 8 * len(self.events)

        return length_in_bits

    @staticmethod
    def static_parse_builder(read_buffer: ReadBuffer, response: bool):
        read_buffer.push_context("ModbusPDUGetComEventLogResponse")

        byte_count: int = read_implicit_field("byteCount", read_unsigned_short)

        self.status = read_simple_field("status", read_unsigned_int)

        self.event_count = read_simple_field("eventCount", read_unsigned_int)

        self.message_count = read_simple_field("messageCount", read_unsigned_int)

        self.events = read_buffer.read_byte_array("events", int(byte_count - int(6)))

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
        write_buffer_box_based: WriteBufferBoxBased = WriteBufferBoxBased(True, True)
        try:
            write_buffer_box_based.writeSerializable(self)
        except SerializationException as e:
            raise RuntimeException(e)

        return "\n" + str(write_buffer_box_based.get_box()) + "\n"


@dataclass
class ModbusPDUGetComEventLogResponseBuilder(ModbusPDUBuilder):
    status: int
    eventCount: int
    messageCount: int
    events: List[int]

    def __post_init__(self):
        pass

    def build(
        self,
    ) -> ModbusPDUGetComEventLogResponse:
        modbus_pdu_get_com_event_log_response: ModbusPDUGetComEventLogResponse = (
            ModbusPDUGetComEventLogResponse(
                self.status, self.event_count, self.message_count, self.events
            )
        )
        return modbus_pdu_get_com_event_log_response
