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
class UmasVariableBlock(ABC, PlcMessage):
    # Abstract accessors for discriminator values.
    @property
    def record_format(self) -> int:
        pass

    @abstractmethod
    def serialize_umas_variable_block_child(self, write_buffer: WriteBuffer) -> None:
        pass

    def serialize(self, write_buffer: WriteBuffer):
        write_buffer.push_context("UmasVariableBlock")

        # Switch field (Serialize the sub-type)
        self.serialize_umas_variable_block_child(write_buffer)

        write_buffer.pop_context("UmasVariableBlock")

    def length_in_bytes(self) -> int:
        return int(math.ceil(float(self.length_in_bits() / 8.0)))

    def length_in_bits(self) -> int:
        length_in_bits: int = 0
        _value: UmasVariableBlock = self

        # Length of subtype elements will be added by sub-type...

        return length_in_bits

    @staticmethod
    def static_parse(read_buffer: ReadBuffer, **kwargs):
        if kwargs is None:
            raise PlcRuntimeException(
                "Wrong number of arguments, expected 1, but got None"
            )

        record_format: int = 0
        if isinstance(kwargs.get("record_format"), int):
            record_format = int(kwargs.get("record_format"))
        elif isinstance(kwargs.get("record_format"), str):
            record_format = int(str(kwargs.get("record_format")))
        else:
            raise PlcRuntimeException(
                "Argument 0 expected to be of type int or a string which is parseable but was "
                + kwargs.get("record_format").getClass().getName()
            )

        return UmasVariableBlock.static_parse_context(read_buffer, record_format)

    @staticmethod
    def static_parse_context(read_buffer: ReadBuffer, record_format: int):
        read_buffer.push_context("UmasVariableBlock")

        # Switch Field (Depending on the discriminator values, passes the instantiation to a sub-type)
        builder: UmasVariableBlockBuilder = None
        from plc4py.protocols.umas.readwrite.UmasPDUReadUnlocatedVariableNamesResponse import (
            UmasPDUReadUnlocatedVariableNamesResponse,
        )

        if record_format == int(0xDD02):
            builder = UmasPDUReadUnlocatedVariableNamesResponse.static_parse_builder(
                read_buffer, record_format
            )
        from plc4py.protocols.umas.readwrite.UmasPDUReadDatatypeNamesResponse import (
            UmasPDUReadDatatypeNamesResponse,
        )

        if record_format == int(0xDD03):
            builder = UmasPDUReadDatatypeNamesResponse.static_parse_builder(
                read_buffer, record_format
            )
        if builder is None:
            raise ParseException(
                "Unsupported case for discriminated type"
                + " parameters ["
                + "recordFormat="
                + str(record_format)
                + "]"
            )

        read_buffer.pop_context("UmasVariableBlock")
        # Create the instance
        _umas_variable_block: UmasVariableBlock = builder.build()
        return _umas_variable_block

    def equals(self, o: object) -> bool:
        if self == o:
            return True

        if not isinstance(o, UmasVariableBlock):
            return False

        that: UmasVariableBlock = UmasVariableBlock(o)
        return True

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
class UmasVariableBlockBuilder:
    def build(
        self,
    ) -> UmasVariableBlock:
        pass
