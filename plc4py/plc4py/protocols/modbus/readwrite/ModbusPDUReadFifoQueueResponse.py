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
class ModbusPDUReadFifoQueueResponse(ModbusPDU):
    fifo_value: List[int]
    # Accessors for discriminator values.
    error_flag: ClassVar[bool] = False
    function_flag: ClassVar[int] = 0x18
    response: ClassVar[bool] = True

    def serialize_modbus_pdu_child(self, write_buffer: WriteBuffer):
        write_buffer.push_context("ModbusPDUReadFifoQueueResponse")

        # Implicit Field (byte_count) (Used for parsing, but its value is not stored as it's implicitly given by the objects content)
        byte_count: int = (int(len(self.fifo_value)) * int(2)) + int(2)
        write_buffer.write_unsigned_short(byte_count, logical_name="byte_count")

        # Implicit Field (fifo_count) (Used for parsing, but its value is not stored as it's implicitly given by the objects content)
        fifo_count: int = (int(len(self.fifo_value)) * int(2)) / int(2)
        write_buffer.write_unsigned_short(fifo_count, logical_name="fifo_count")

        # Array Field (fifoValue)
        write_buffer.write_simple_array(
            self.fifo_value,
            write_buffer.write_unsigned_short,
            logical_name="fifo_value",
        )

        write_buffer.pop_context("ModbusPDUReadFifoQueueResponse")

    def length_in_bytes(self) -> int:
        return int(math.ceil(float(self.length_in_bits() / 8.0)))

    def length_in_bits(self) -> int:
        length_in_bits: int = super().length_in_bits()
        _value: ModbusPDUReadFifoQueueResponse = self

        # Implicit Field (byteCount)
        length_in_bits += 16

        # Implicit Field (fifoCount)
        length_in_bits += 16

        # Array field
        if self.fifo_value is not None:
            length_in_bits += 16 * len(self.fifo_value)

        return length_in_bits

    @staticmethod
    def static_parse_builder(read_buffer: ReadBuffer, response: bool):
        read_buffer.push_context("ModbusPDUReadFifoQueueResponse")

        if isinstance(response, str):
            response = bool(strtobool(response))

        byte_count: int = read_buffer.read_unsigned_short(
            logical_name="byte_count", response=response
        )

        fifo_count: int = read_buffer.read_unsigned_short(
            logical_name="fifo_count", response=response
        )

        fifo_value: List[Any] = read_buffer.read_array_field(
            logical_name="fifoValue",
            read_function=read_buffer.read_unsigned_short,
            count=fifo_count,
            response=response,
        )

        read_buffer.pop_context("ModbusPDUReadFifoQueueResponse")
        # Create the instance
        return ModbusPDUReadFifoQueueResponseBuilder(fifo_value)

    def equals(self, o: object) -> bool:
        if self == o:
            return True

        if not isinstance(o, ModbusPDUReadFifoQueueResponse):
            return False

        that: ModbusPDUReadFifoQueueResponse = ModbusPDUReadFifoQueueResponse(o)
        return (self.fifo_value == that.fifo_value) and super().equals(that) and True

    def hash_code(self) -> int:
        return hash(self)

    def __str__(self) -> str:
        # TODO:- Implement a generic python object to probably json convertor or something.
        return ""


@dataclass
class ModbusPDUReadFifoQueueResponseBuilder:
    fifo_value: List[int]

    def build(
        self,
    ) -> ModbusPDUReadFifoQueueResponse:
        modbus_pduread_fifo_queue_response: ModbusPDUReadFifoQueueResponse = (
            ModbusPDUReadFifoQueueResponse(self.fifo_value)
        )
        return modbus_pduread_fifo_queue_response
