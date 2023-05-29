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
from plc4py.spi.generation.ReadBuffer import ReadBuffer
from plc4py.spi.generation.WriteBuffer import WriteBuffer
import math


@dataclass
class Dummy(PlcMessage):
    dummy: int

    def __post_init__(self):
        super().__init__()

    def serialize(self, write_buffer: WriteBuffer):
        write_buffer.push_context("Dummy")

        # Simple Field (dummy)
        write_buffer.write_unsigned_short(self.dummy, logical_name="dummy")

        write_buffer.pop_context("Dummy")

    def length_in_bytes(self) -> int:
        return int(math.ceil(float(self.get_length_in_bits() / 8.0)))

    def get_length_in_bits(self) -> int:
        length_in_bits: int = 0
        _value: Dummy = self

        # Simple field (dummy)
        length_in_bits += 16

        return length_in_bits

    def static_parse(self, read_buffer: ReadBuffer, args):
        return self.static_parse_context(read_buffer)

    @staticmethod
    def static_parse_context(read_buffer: ReadBuffer):
        read_buffer.push_context("Dummy")

        self.dummy = read_simple_field(
            "dummy", read_unsigned_int, WithOption.WithByteOrder(get_bi_g__endian())
        )

        read_buffer.pop_context("Dummy")
        # Create the instance
        _dummy: Dummy = Dummy(dummy)
        return _dummy

    def equals(self, o: object) -> bool:
        if self == o:
            return True

        if not isinstance(o, Dummy):
            return False

        that: Dummy = Dummy(o)
        return (self.dummy == that.dummy) and True

    def hash_code(self) -> int:
        return hash(self)

    def __str__(self) -> str:
        write_buffer_box_based: WriteBufferBoxBased = WriteBufferBoxBased(True, True)
        try:
            write_buffer_box_based.writeSerializable(self)
        except SerializationException as e:
            raise RuntimeException(e)

        return "\n" + str(write_buffer_box_based.get_box()) + "\n"
