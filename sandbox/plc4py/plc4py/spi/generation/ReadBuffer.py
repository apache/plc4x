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
import struct
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
from bitarray.util import zeros, ba2int, ba2base

from plc4py.api.exceptions.exceptions import SerializationException
from plc4py.api.messages.PlcMessage import PlcMessage
from plc4py.utils.GenericTypes import ByteOrder, ByteOrderAware


class PositionAware:
    def get_pos(self) -> int:
        raise NotImplementedError


@dataclass
class ReadBuffer(ByteOrderAware, PositionAware):
    byte_order: ByteOrder

    def get_pos(self) -> int:
        raise NotImplementedError

    def push_context(self, logical_name: str, **kwargs) -> None:
        raise NotImplementedError

    def pop_context(self, logical_name: str, **kwargs) -> None:
        raise NotImplementedError

    def read_bit(self, logical_name: str = "", **kwargs) -> bool:
        raise NotImplementedError

    def read_byte(self, logical_name: str = "", **kwargs) -> int:
        self.read_signed_byte(8, logical_name, **kwargs)

    def read_byte_array(
        self, number_of_bytes: int, logical_name: str = "", **kwargs
    ) -> List[int]:
        raise NotImplementedError

    def read_unsigned_byte(
        self, bit_length: int = 8, logical_name: str = "", **kwargs
    ) -> int:
        raise NotImplementedError

    def read_unsigned_short(
        self, bit_length: int = 16, logical_name: str = "", **kwargs
    ) -> int:
        raise NotImplementedError

    def read_unsigned_int(
        self, bit_length: int = 32, logical_name: str = "", **kwargs
    ) -> int:
        raise NotImplementedError

    def read_unsigned_long(
        self, bit_length: int = 64, logical_name: str = "", **kwargs
    ) -> int:
        raise NotImplementedError

    def read_signed_byte(
        self, bit_length: int = 8, logical_name: str = "", **kwargs
    ) -> int:
        raise NotImplementedError

    def read_short(self, bit_length: int = 16, logical_name: str = "", **kwargs) -> int:
        raise NotImplementedError

    def read_int(self, bit_length: int = 32, logical_name: str = "", **kwargs) -> int:
        raise NotImplementedError

    def read_long(self, bit_length: int = 64, logical_name: str = "", **kwargs) -> int:
        raise NotImplementedError

    def read_float(
        self, bit_length: int = 32, logical_name: str = "", **kwargs
    ) -> float:
        raise NotImplementedError

    def read_double(
        self, bit_length: int = 64, logical_name: str = "", **kwargs
    ) -> float:
        raise NotImplementedError

    def read_str(self, bit_length: int = -1, logical_name: str = "", **kwargs) -> str:
        raise NotImplementedError

    def read_virtual(self, logical_name: str = "", **kwargs) -> str:
        raise NotImplementedError

    def read_complex_array(self, logical_name: str = "", **kwargs) -> List[PlcMessage]:
        raise NotImplementedError


class ReadBufferByteBased(ReadBuffer):
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

    def __init__(self, bb: bytearray, byte_order: ByteOrder):
        if byte_order == ByteOrder.LITTLE_ENDIAN:
            bb = bitarray(buffer=bb)
            bb.bytereverse()
        self.bb = bitarray(
            buffer=memoryview(bb), endian=ByteOrder.get_short_name(byte_order)
        )
        self.byte_order = byte_order

    def get_bytes(self) -> memoryview:
        return memoryview(self.bb)

    def get_pos(self) -> int:
        return self.position

    def push_context(self, logical_name: str, **kwargs) -> None:
        # byte buffer need no context handling
        pass

    def read_bit(self, logical_name: str = "", **kwargs) -> bool:
        result: bool = bool(self.bb[self.position])
        self.position += 1
        return result

    def read_byte(self, logical_name: str = "", **kwargs) -> int:
        value = self.read_signed_byte(8, logical_name, **kwargs)
        return value

    def read_byte_array(
        self, number_of_bytes: int, logical_name: str = "", **kwargs
    ) -> List[int]:
        result: List[int] = [] * number_of_bytes
        for index in range(0, number_of_bytes):
            result[index] = self.read_signed_byte(8, logical_name, **kwargs)
        return result

    def read_unsigned_byte(
        self, bit_length: int = 8, logical_name: str = "", **kwargs
    ) -> int:
        byte_order = kwargs.get("byte_order", self.byte_order)
        if bit_length <= 0:
            raise SerializationException("unsigned byte must contain at least 1 bit")
        elif bit_length > 8:
            raise SerializationException("unsigned byte can only contain max 8 bits")
        else:
            result: int = ba2int(self.bb[self.position : bit_length], signed=False)
            self.position += bit_length
            return result

    def read_unsigned_short(
        self, bit_length: int = 16, logical_name: str = "", **kwargs
    ) -> int:
        byte_order = kwargs.get("byte_order", self.byte_order)
        if bit_length <= 0:
            raise SerializationException("unsigned short must contain at least 1 bit")
        elif bit_length > 16:
            raise SerializationException("unsigned short can only contain max 16 bits")
        else:
            result: int = ba2int(self.bb[self.position : bit_length], signed=False)
            self.position += bit_length
            return result

    def read_unsigned_int(
        self, bit_length: int = 32, logical_name: str = "", **kwargs
    ) -> int:
        byte_order = kwargs.get("byte_order", self.byte_order)
        if bit_length <= 0:
            raise SerializationException("unsigned int must contain at least 1 bit")
        elif bit_length > 32:
            raise SerializationException("unsigned int can only contain max 32 bits")
        else:
            result: int = ba2int(self.bb[self.position : bit_length], signed=False)
            self.position += bit_length
            return result

    def read_unsigned_long(
        self, bit_length: int = 64, logical_name: str = "", **kwargs
    ) -> int:
        byte_order = kwargs.get("byte_order", self.byte_order)
        if bit_length <= 0:
            raise SerializationException("unsigned long must contain at least 1 bit")
        elif bit_length > 64:
            raise SerializationException("unsigned long can only contain max 64 bits")
        else:
            result: int = ba2int(self.bb[self.position : bit_length], signed=False)
            self.position += bit_length
            return result

    def read_signed_byte(
        self, bit_length: int = 8, logical_name: str = "", **kwargs
    ) -> int:
        byte_order = kwargs.get("byte_order", self.byte_order)
        if bit_length <= 0:
            raise SerializationException("signed byte must contain at least 1 bit")
        elif bit_length > 8:
            raise SerializationException("signed byte can only contain max 8 bits")
        else:
            result: int = ba2int(self.bb[self.position : bit_length], signed=False)
            self.position += bit_length
            return result

    def read_short(self, bit_length: int = 16, logical_name: str = "", **kwargs) -> int:
        byte_order = kwargs.get("byte_order", self.byte_order)
        if bit_length <= 0:
            raise SerializationException("signed short must contain at least 1 bit")
        elif bit_length > 16:
            raise SerializationException("signed short can only contain max 16 bits")
        else:
            result: int = ba2int(self.bb[self.position : bit_length], signed=False)
            self.position += bit_length
            return result

    def read_int(self, bit_length: int = 32, logical_name: str = "", **kwargs) -> int:
        byte_order = kwargs.get("byte_order", self.byte_order)
        if bit_length <= 0:
            raise SerializationException("signed int must contain at least 1 bit")
        elif bit_length > 32:
            raise SerializationException("signed int can only contain max 32 bits")
        else:
            result: int = ba2int(self.bb[self.position : bit_length], signed=False)
            self.position += bit_length
            return result

    def read_long(self, bit_length: int = 64, logical_name: str = "", **kwargs) -> int:
        byte_order = kwargs.get("byte_order", self.byte_order)
        if bit_length <= 0:
            raise SerializationException("signed long must contain at least 1 bit")
        elif bit_length > 64:
            raise SerializationException("signed long can only contain max 64 bits")
        else:
            result: int = ba2int(self.bb[self.position : bit_length], signed=False)
            self.position += bit_length
            return result

    def read_float(
        self, bit_length: int = 32, logical_name: str = "", **kwargs
    ) -> float:
        byte_order = kwargs.get("byte_order", self.byte_order)
        if bit_length != 32:
            raise SerializationException("Float can only contain 32 bits")
        else:
            endianness: str = ">"
            if byte_order == ByteOrder.LITTLE_ENDIAN:
                endianness = "<"
            result: float = struct.unpack(
                endianness + "f", self.bb[self.position : bit_length]
            )[0]
            self.position += bit_length
            return result

    def read_double(
        self, bit_length: int = 64, logical_name: str = "", **kwargs
    ) -> float:
        byte_order = kwargs.get("byte_order", self.byte_order)
        if bit_length != 64:
            raise SerializationException("Double can only contain 64 bits")
        else:
            endianness: str = ">"
            if byte_order == ByteOrder.LITTLE_ENDIAN:
                endianness = "<"
            result: float = struct.unpack(
                endianness + "d", self.bb[self.position : bit_length]
            )[0]
            self.position += bit_length
            return result
