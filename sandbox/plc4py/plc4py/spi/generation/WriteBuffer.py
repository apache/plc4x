#  Licensed to the Apache Software Foundation (ASF) under one
#  or more contributor license agreements.  See the NOTICE file
#  distributed with this work for additional information
#  regarding copyright ownership.  The ASF licenses this file
#  to you under the Apache License, Version 2.0 (the
#  "License"); you may not use this file except in compliance
#  with the License.  You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
#  Unless required by applicable law or agreed to in writing,
#  software distributed under the License is distributed on an
#  "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
#  KIND, either express or implied.  See the License for the
#  specific language governing permissions and limitations
#  under the License.
from ctypes import (
    c_byte,
    c_ubyte,
    c_uint16,
    c_uint32,
    c_uint64,
    c_int16,
    c_int32,
    c_int64,
    c_float,
    c_double,
    c_int8,
    c_uint8,
)
from dataclasses import dataclass
from typing import List, Union

from bitarray import bitarray
from bitarray.util import zeros

from plc4py.api.exceptions.exceptions import SerializationException
from plc4py.api.messages.PlcMessage import PlcMessage
from plc4py.utils.GenericTypes import ByteOrder, ByteOrderAware


class PositionAware:
    def get_pos(self) -> int:
        raise NotImplementedError


@dataclass
class WriteBuffer(ByteOrderAware, PositionAware):
    byte_order: ByteOrder

    def get_pos(self) -> int:
        raise NotImplementedError

    def push_context(self, logical_name: str, **kwargs) -> None:
        raise NotImplementedError

    def pop_context(self, logical_name: str, **kwargs) -> None:
        raise NotImplementedError

    def write_bit(self, value: bool, logical_name: str = "", **kwargs) -> None:
        raise NotImplementedError

    def write_byte(self, value: int, logical_name: str = "", **kwargs) -> None:
        self.write_signed_byte(value, 8, logical_name, **kwargs)

    def write_byte_array(
        self, value: List[int], logical_name: str = "", **kwargs
    ) -> None:
        raise NotImplementedError

    def write_unsigned_byte(
        self, value: int, bit_length: int = 8, logical_name: str = "", **kwargs
    ) -> None:
        raise NotImplementedError

    def write_unsigned_short(
        self, value: int, bit_length: int = 16, logical_name: str = "", **kwargs
    ) -> None:
        raise NotImplementedError

    def write_unsigned_int(
        self, value: int, bit_length: int = 32, logical_name: str = "", **kwargs
    ) -> None:
        raise NotImplementedError

    def write_unsigned_long(
        self, value: int, bit_length: int = 64, logical_name: str = "", **kwargs
    ) -> None:
        raise NotImplementedError

    def write_signed_byte(
        self, value: int, bit_length: int = 8, logical_name: str = "", **kwargs
    ) -> None:
        raise NotImplementedError

    def write_short(
        self, value: int, bit_length: int = 16, logical_name: str = "", **kwargs
    ) -> None:
        raise NotImplementedError

    def write_int(
        self, value: int, bit_length: int = 32, logical_name: str = "", **kwargs
    ) -> None:
        raise NotImplementedError

    def write_long(
        self, value: int, bit_length: int = 64, logical_name: str = "", **kwargs
    ) -> None:
        raise NotImplementedError

    def write_float(
        self, value: float, bit_length: int = 32, logical_name: str = "", **kwargs
    ) -> None:
        raise NotImplementedError

    def write_double(
        self, value: float, bit_length: int = 64, logical_name: str = "", **kwargs
    ) -> None:
        raise NotImplementedError

    def write_str(
        self, value: str, bit_length: int = -1, logical_name: str = "", **kwargs
    ) -> None:
        raise NotImplementedError

    def write_virtual(self, value: str, logical_name: str = "", **kwargs) -> None:
        raise NotImplementedError

    def write_complex_array(
        self, value: List[PlcMessage], logical_name: str = "", **kwargs
    ) -> None:
        raise NotImplementedError

    #
    # This method can be used to influence serializing (e.g. intercept whole types and render them in a simplified form)
    #
    # @param value the value to be serialized
    # @throws SerializationException if something goes wrong
    #
    def write_serializable(self, value, logical_name="") -> None:
        value.serialize(self)


class WriteBufferByteBased(WriteBuffer):
    byte_order: ByteOrder
    position: int = 0
    bb: bitarray

    NUMERIC_UNION = Union[
        c_ubyte,
        c_byte,
        c_uint8,
        c_uint16,
        c_uint32,
        c_uint64,
        c_int8,
        c_int16,
        c_int32,
        c_int64,
        c_float,
        c_double,
    ]

    def __init__(self, size: int, byte_order: ByteOrder):
        self.bb = zeros(size * 8, endian=ByteOrder.get_short_name(byte_order))
        self.byte_order = byte_order

    def get_bytes(self) -> memoryview:
        return memoryview(self.bb)

    def get_pos(self) -> int:
        return self.position

    def push_context(self, logical_name: str, **kwargs) -> None:
        # byte buffer need no context handling
        pass

    def write_bit(self, value: bool, logical_name: str = "", **kwargs) -> None:
        self.bb[self.position] = value
        self.position += 1

    def write_byte(self, value: int, logical_name: str = "", **kwargs) -> None:
        self.write_signed_byte(value, 8, logical_name, **kwargs)

    def write_byte_array(
        self, value: List[int], logical_name: str = "", **kwargs
    ) -> None:
        for a_byte in value:
            self.write_signed_byte(a_byte, 8, logical_name, **kwargs)

    def write_unsigned_byte(
        self, value: int, bit_length: int = 8, logical_name: str = "", **kwargs
    ) -> None:
        if bit_length <= 0:
            raise SerializationException("unsigned byte must contain at least 1 bit")
        elif bit_length > 8:
            raise SerializationException("unsigned byte can only contain max 8 bits")
        else:
            self._handle_numeric_encoding(c_ubyte(value), bit_length, **kwargs)

    def write_unsigned_short(
        self, value: int, bit_length: int = 16, logical_name: str = "", **kwargs
    ) -> None:
        if bit_length <= 0:
            raise SerializationException("unsigned short must contain at least 1 bit")
        elif bit_length > 16:
            raise SerializationException("unsigned short can only contain max 16 bits")
        else:
            self._handle_numeric_encoding(c_uint16(value), bit_length, **kwargs)

    def write_unsigned_int(
        self, value: int, bit_length: int = 32, logical_name: str = "", **kwargs
    ) -> None:
        if bit_length <= 0:
            raise SerializationException("unsigned int must contain at least 1 bit")
        elif bit_length > 32:
            raise SerializationException("unsigned int can only contain max 32 bits")
        else:
            self._handle_numeric_encoding(c_uint32(value), bit_length, **kwargs)

    def write_unsigned_long(
        self, value: int, bit_length: int = 64, logical_name: str = "", **kwargs
    ) -> None:
        if bit_length <= 0:
            raise SerializationException("unsigned long must contain at least 1 bit")
        elif bit_length > 64:
            raise SerializationException("unsigned long can only contain max 16 bits")
        else:
            self._handle_numeric_encoding(c_uint64(value), bit_length, **kwargs)

    def write_signed_byte(
        self, value: int, bit_length: int = 8, logical_name: str = "", **kwargs
    ) -> None:
        if bit_length <= 0:
            raise SerializationException("Signed byte must contain at least 1 bit")
        elif bit_length > 8:
            raise SerializationException("Signed byte can only contain max 8 bits")
        self._handle_numeric_encoding(c_int8(value), bit_length, **kwargs)

    def write_short(
        self, value: int, bit_length: int = 16, logical_name: str = "", **kwargs
    ) -> None:
        if bit_length <= 0:
            raise SerializationException("Signed short must contain at least 1 bit")
        elif bit_length > 16:
            raise SerializationException("Signed short can only contain max 16 bits")
        self._handle_numeric_encoding(c_int16(value), bit_length, **kwargs)

    def write_int(
        self, value: int, bit_length: int = 32, logical_name: str = "", **kwargs
    ) -> None:
        if bit_length <= 0:
            raise SerializationException("Signed int must contain at least 1 bit")
        elif bit_length > 32:
            raise SerializationException("Signed int can only contain max 32 bits")
        self._handle_numeric_encoding(c_int32(value), bit_length, **kwargs)

    def write_long(
        self, value: int, bit_length: int = 64, logical_name: str = "", **kwargs
    ) -> None:
        if bit_length <= 0:
            raise SerializationException("Signed long must contain at least 1 bit")
        elif bit_length > 64:
            raise SerializationException("Signed long can only contain max 64 bits")
        self._handle_numeric_encoding(c_int64(value), bit_length, **kwargs)

    def write_float(
        self, value: float, bit_length: int = 32, logical_name: str = "", **kwargs
    ) -> None:
        if bit_length <= 0:
            raise SerializationException("Float must contain at least 1 bit")
        elif bit_length > 32:
            raise SerializationException("Float can only contain max 32 bits")
        self._handle_numeric_encoding(c_float(value), bit_length, **kwargs)

    def write_double(
        self, value: float, bit_length: int = 64, logical_name: str = "", **kwargs
    ) -> None:
        if bit_length <= 0:
            raise SerializationException("Double must contain at least 1 bit")
        elif bit_length > 64:
            raise SerializationException("Double can only contain max 64 bits")
        self._handle_numeric_encoding(c_double(value), bit_length, **kwargs)

    def write_complex_array(
        self, value: List[PlcMessage], logical_name: str = "", **kwargs
    ) -> None:
        for item in value:
            self.push_context(logical_name, **kwargs)
            self.write_serializable(item, logical_name="")
            self.pop_context(logical_name, **kwargs)

    def _handle_numeric_encoding(self, value: NUMERIC_UNION, bit_length: int, **kwargs):
        byte_order = kwargs.get("byte_order", self.byte_order)
        value_encoding: str = kwargs.get("encoding", "default")
        if value_encoding == "ASCII":
            if bit_length % 8 != 0:
                raise SerializationException(
                    "'ASCII' encoded fields must have a length that is a multiple of 8 bits long"
                )
            char_len: int = int(bit_length / 8)
            max_value: int = int(10**char_len - 1)
            if value.value > max_value:
                raise SerializationException(
                    "Provided value of "
                    + str(value)
                    + " exceeds the max value of "
                    + str(max_value)
                )
            string_value: str = "{}".format(value.value)
            src = bitarray(endian=ByteOrder.get_short_name(byte_order))
            src.frombytes(bytearray(string_value, value_encoding))
            self.bb[self.position : bit_length] = src[:bit_length]
            self.position += bit_length
        elif value_encoding == "default":
            src = bitarray(endian=ByteOrder.get_short_name(byte_order))
            src.frombytes(value)
            self.bb[self.position : bit_length] = src[:bit_length]
            self.position += bit_length
