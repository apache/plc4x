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
class UmasMemoryBlock(ABC, PlcMessage):

    # Abstract accessors for discriminator values.
    @property
    def block_number(self) -> int:
        pass

    @property
    def offset(self) -> int:
        pass

    @abstractmethod
    def serialize_umas_memory_block_child(self, write_buffer: WriteBuffer) -> None:
        pass

    def serialize(self, write_buffer: WriteBuffer):
        write_buffer.push_context("UmasMemoryBlock")

        # Switch field (Serialize the sub-type)
        self.serialize_umas_memory_block_child(write_buffer)

        write_buffer.pop_context("UmasMemoryBlock")

    def length_in_bytes(self) -> int:
        return int(math.ceil(float(self.length_in_bits() / 8.0)))

    def length_in_bits(self) -> int:
        length_in_bits: int = 0
        _value: UmasMemoryBlock = self

        # Length of subtype elements will be added by sub-type...

        return length_in_bits

    @staticmethod
    def static_parse(read_buffer: ReadBuffer, **kwargs):

        if kwargs is None:
            raise PlcRuntimeException(
                "Wrong number of arguments, expected 2, but got None"
            )

        block_number: int = 0
        if isinstance(kwargs.get("block_number"), int):
            block_number = int(kwargs.get("block_number"))
        elif isinstance(kwargs.get("block_number"), str):
            block_number = int(str(kwargs.get("block_number")))
        else:
            raise PlcRuntimeException(
                "Argument 0 expected to be of type int or a string which is parseable but was "
                + kwargs.get("block_number").getClass().getName()
            )

        offset: int = 0
        if isinstance(kwargs.get("offset"), int):
            offset = int(kwargs.get("offset"))
        elif isinstance(kwargs.get("offset"), str):
            offset = int(str(kwargs.get("offset")))
        else:
            raise PlcRuntimeException(
                "Argument 1 expected to be of type int or a string which is parseable but was "
                + kwargs.get("offset").getClass().getName()
            )

        return UmasMemoryBlock.static_parse_context(read_buffer, block_number, offset)

    @staticmethod
    def static_parse_context(read_buffer: ReadBuffer, block_number: int, offset: int):
        read_buffer.push_context("UmasMemoryBlock")

        if isinstance(block_number, str):
            block_number = int(block_number)
        if isinstance(offset, str):
            offset = int(offset)

        # Switch Field (Depending on the discriminator values, passes the instantiation to a sub-type)
        builder: UmasMemoryBlockBuilder = None
        from plc4py.protocols.umas.readwrite.UmasMemoryBlockBasicInfo import (
            UmasMemoryBlockBasicInfo,
        )

        if block_number == int(0x30) and offset == int(0x00):

            builder = UmasMemoryBlockBasicInfo.static_parse_builder(
                read_buffer, block_number, offset
            )
        if builder is None:
            raise ParseException(
                "Unsupported case for discriminated type"
                + " parameters ["
                + "blockNumber="
                + str(block_number)
                + " "
                + "offset="
                + str(offset)
                + "]"
            )

        read_buffer.pop_context("UmasMemoryBlock")
        # Create the instance
        _umas_memory_block: UmasMemoryBlock = builder.build()
        return _umas_memory_block

    def equals(self, o: object) -> bool:
        if self == o:
            return True

        if not isinstance(o, UmasMemoryBlock):
            return False

        that: UmasMemoryBlock = UmasMemoryBlock(o)
        return True

    def hash_code(self) -> int:
        return hash(self)

    def __str__(self) -> str:
        # TODO:- Implement a generic python object to probably json convertor or something.
        return ""


@dataclass
class UmasMemoryBlockBuilder:
    def build(
        self,
    ) -> UmasMemoryBlock:
        pass
