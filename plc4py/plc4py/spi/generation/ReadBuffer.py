#  Licensed to the Apache Software Foundation (ASF) under one
#  or more contributor license agreements.  See the NOTICE file
#  distributed with this work for additional information
#  regarding copyright ownership.  The ASF licenses this file
#  to you under the Apache License, Version 2.0 (the
#  "License"); you may not use this file except in compliance
#  with the License.  You may obtain a copy of the License at
#
#    https://www.apache.org/licenses/LICENSE-2.0
#
#  Unless required by applicable law or agreed to in writing,
#  software distributed under the License is distributed on an
#  "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
#  KIND, either express or implied.  See the License for the
#  specific language governing permissions and limitations
#  under the License.
import struct
import types
from abc import ABC, abstractmethod
from dataclasses import dataclass
from typing import Any, List, Union

import aenum
from bitarray import bitarray
from bitarray.util import ba2base, ba2int, zeros

from plc4py.api.exceptions.exceptions import SerializationException, ParseException
from plc4py.api.messages.PlcMessage import PlcMessage
from plc4py.utils.GenericTypes import ByteOrder, ByteOrderAware


class PositionAware:
    @abstractmethod
    def get_pos(self) -> int:
        raise NotImplementedError


@dataclass
class ReadBuffer(ByteOrderAware, PositionAware, ABC):
    byte_order: ByteOrder

    @abstractmethod
    def get_pos(self) -> int:
        raise NotImplementedError

    @abstractmethod
    def push_context(self, logical_name: str, **kwargs) -> None:
        raise NotImplementedError

    @abstractmethod
    def pop_context(self, logical_name: str, **kwargs) -> None:
        raise NotImplementedError

    @abstractmethod
    def read_bit(self, logical_name: str = "", **kwargs) -> bool:
        raise NotImplementedError

    def read_byte(self, logical_name: str = "", **kwargs) -> int:
        self.read_signed_byte(8, logical_name, **kwargs)

    @abstractmethod
    def read_byte_array(
        self, number_of_bytes: int, logical_name: str = "", **kwargs
    ) -> List[int]:
        raise NotImplementedError

    @abstractmethod
    def read_unsigned_byte(
        self, bit_length: int = 8, logical_name: str = "", **kwargs
    ) -> int:
        raise NotImplementedError

    @abstractmethod
    def read_unsigned_short(
        self, bit_length: int = 16, logical_name: str = "", **kwargs
    ) -> int:
        raise NotImplementedError

    @abstractmethod
    def read_unsigned_int(
        self, bit_length: int = 32, logical_name: str = "", **kwargs
    ) -> int:
        raise NotImplementedError

    @abstractmethod
    def read_unsigned_long(
        self, bit_length: int = 64, logical_name: str = "", **kwargs
    ) -> int:
        raise NotImplementedError

    @abstractmethod
    def read_signed_byte(
        self, bit_length: int = 8, logical_name: str = "", **kwargs
    ) -> int:
        raise NotImplementedError

    @abstractmethod
    def read_short(self, bit_length: int = 16, logical_name: str = "", **kwargs) -> int:
        raise NotImplementedError

    @abstractmethod
    def read_int(self, bit_length: int = 32, logical_name: str = "", **kwargs) -> int:
        raise NotImplementedError

    @abstractmethod
    def read_long(self, bit_length: int = 64, logical_name: str = "", **kwargs) -> int:
        raise NotImplementedError

    @abstractmethod
    def read_float(
        self, bit_length: int = 32, logical_name: str = "", **kwargs
    ) -> float:
        raise NotImplementedError

    @abstractmethod
    def read_double(
        self, bit_length: int = 64, logical_name: str = "", **kwargs
    ) -> float:
        raise NotImplementedError

    @abstractmethod
    def read_str(self, bit_length: int = -1, logical_name: str = "", **kwargs) -> str:
        raise NotImplementedError

    @abstractmethod
    def read_complex_array(self, logical_name: str = "", **kwargs) -> List[Any]:
        raise NotImplementedError

    @abstractmethod
    def read_complex(self, logical_name: str = "", read_function=None, **kwargs) -> Any:
        raise NotImplementedError

    @abstractmethod
    def read_enum(
        self,
        bit_length: int = -1,
        logical_name: str = "",
        read_function=None,
        **kwargs,
    ) -> Any:
        raise NotImplementedError

    @abstractmethod
    def read_array_field(
        self,
        logical_name: str = "",
        read_function=None,
        count: int = None,
        length: int = None,
        terminated=None,
        **kwargs,
    ) -> List[Any]:
        raise NotImplementedError


class ReadBufferByteBased(ReadBuffer):
    def __init__(
        self,
        bb: bytearray,
        byte_order: ByteOrder,
        bit_order: ByteOrder = ByteOrder.BIG_ENDIAN,
    ):
        if bit_order == ByteOrder.LITTLE_ENDIAN:
            bb = bitarray(buffer=bb, endian=ByteOrder.get_short_name(bit_order))
            bb.bytereverse()
        self.bb = bitarray(
            buffer=memoryview(bb), endian=ByteOrder.get_short_name(bit_order)
        )
        self.byte_order = byte_order
        self.position = 0

    def get_bytes(self) -> memoryview:
        return memoryview(self.bb)

    def get_pos(self) -> int:
        return self.position

    def push_context(self, logical_name: str, **kwargs) -> None:
        # byte buffer need no context handling
        pass

    def pop_context(self, logical_name: str, **kwargs) -> None:
        # Byte buffer doesn't need context handling
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
            result: int = ba2int(
                self.bb[self.position : self.position + bit_length],
                signed=False,
            )
            self.position += bit_length
            return result

    def read_unsigned_short(
        self, bit_length: int = 16, logical_name: str = "", **kwargs
    ) -> int:
        byte_order = kwargs.get("byte_order", self.byte_order)
        encoding = kwargs.get("encoding", "")
        if bit_length <= 0:
            raise SerializationException("unsigned short must contain at least 1 bit")
        elif bit_length > 16:
            raise SerializationException("unsigned short can only contain max 16 bits")
        else:
            if byte_order == ByteOrder.LITTLE_ENDIAN:
                endian_string = "<"
                padded = bitarray(
                    self.bb[self.position : self.position + bit_length]
                ) + (16 - bit_length) * bitarray("0")
            else:
                endian_string = ">"
                padded = (16 - bit_length) * bitarray("0") + bitarray(
                    self.bb[self.position : self.position + bit_length]
                )

            if encoding == "BCD":
                if bit_length % 4 != 0:
                    raise ParseException(
                        "'BCD' encoded fields must have a length that is a multiple of 4 bits long"
                    )
                result: int = 0
                for i in range(0, bit_length, 4):
                    digit: int = ba2int(padded[i : i + 4])
                    if digit > 9:
                        raise ParseException(
                            "'BCD' encoded value is not a correctly encoded BCD value"
                        )
                    multiplier = 10 ** ((int(bit_length - i) / 4) - 1)
                    result += multiplier * digit
                self.position += bit_length
                return result
            else:
                result: int = struct.unpack(endian_string + "H", padded)[0]
                self.position += bit_length
                return result

    def read_unsigned_int(
        self, bit_length: int = 32, logical_name: str = "", **kwargs
    ) -> int:
        byte_order = kwargs.get("byte_order", self.byte_order)
        encoding = kwargs.get("encoding", "")
        if bit_length <= 0:
            raise SerializationException("unsigned int must contain at least 1 bit")
        elif bit_length > 32:
            raise SerializationException("unsigned int can only contain max 32 bits")
        else:
            if byte_order == ByteOrder.LITTLE_ENDIAN:
                endian_string = "<"
                padded = bitarray(
                    self.bb[self.position : self.position + bit_length]
                ) + (32 - bit_length) * bitarray("0")
            else:
                endian_string = ">"
                padded = (32 - bit_length) * bitarray("0") + bitarray(
                    self.bb[self.position : self.position + bit_length]
                )
            if encoding == "BCD":
                if bit_length % 4 != 0:
                    raise ParseException(
                        "'BCD' encoded fields must have a length that is a multiple of 4 bits long"
                    )
                if byte_order == ByteOrder.LITTLE_ENDIAN:
                    padded = padded[8:16] + padded[:8] + padded[24:32] + padded[16:24]
                result: int = 0
                for i in range(0, bit_length, 4):
                    digit: int = ba2int(padded[i : i + 4])
                    if digit > 9:
                        raise ParseException(
                            "'BCD' encoded value is not a correctly encoded BCD value"
                        )
                    multiplier = 10 ** ((int(bit_length - i) / 4) - 1)
                    result += multiplier * digit
                self.position += bit_length
                return result
            else:
                result: int = struct.unpack(endian_string + "I", padded)[0]
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
            if byte_order == ByteOrder.LITTLE_ENDIAN:
                endian_string = "<"
                padded = bitarray(
                    self.bb[self.position : self.position + bit_length]
                ) + (64 - bit_length) * bitarray("0")
            else:
                endian_string = ">"
                padded = (64 - bit_length) * bitarray("0") + bitarray(
                    self.bb[self.position : self.position + bit_length]
                )
            result: int = struct.unpack(endian_string + "Q", padded)[0]
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
            result: int = ba2int(
                self.bb[self.position : self.position + bit_length],
                signed=False,
            )
            self.position += bit_length
            return result

    def read_short(self, bit_length: int = 16, logical_name: str = "", **kwargs) -> int:
        byte_order = kwargs.get("byte_order", self.byte_order)
        if bit_length <= 0:
            raise SerializationException("signed short must contain at least 1 bit")
        elif bit_length > 16:
            raise SerializationException("signed short can only contain max 16 bits")
        else:
            if byte_order == ByteOrder.LITTLE_ENDIAN:
                endian_string = "<"
                padded = bitarray(
                    self.bb[self.position : self.position + bit_length]
                ) + (16 - bit_length) * bitarray("0")
            else:
                endian_string = ">"
                padded = (16 - bit_length) * bitarray("0") + bitarray(
                    self.bb[self.position : self.position + bit_length]
                )
            result: int = struct.unpack(endian_string + "h", padded)[0]
            self.position += bit_length
            return result

    def read_int(self, bit_length: int = 32, logical_name: str = "", **kwargs) -> int:
        byte_order = kwargs.get("byte_order", self.byte_order)
        if bit_length <= 0:
            raise SerializationException("signed int must contain at least 1 bit")
        elif bit_length > 32:
            raise SerializationException("signed int can only contain max 32 bits")
        else:
            if byte_order == ByteOrder.LITTLE_ENDIAN:
                endian_string = "<"
                padded = bitarray(
                    self.bb[self.position : self.position + bit_length]
                ) + (32 - bit_length) * bitarray("0")
            else:
                endian_string = ">"
                padded = (32 - bit_length) * bitarray("0") + bitarray(
                    self.bb[self.position : self.position + bit_length]
                )
            if (
                byte_order == ByteOrder.BIG_ENDIAN_BYTE_SWAP
                or byte_order == ByteOrder.LITTLE_ENDIAN_BYTE_SWAP
            ):
                padded = padded[16:] + padded[0:16]
            result: int = struct.unpack(endian_string + "i", padded)[0]
            self.position += bit_length
            return result

    def read_long(self, bit_length: int = 64, logical_name: str = "", **kwargs) -> int:
        byte_order = kwargs.get("byte_order", self.byte_order)
        if bit_length <= 0:
            raise SerializationException("signed long must contain at least 1 bit")
        elif bit_length > 64:
            raise SerializationException("signed long can only contain max 64 bits")
        else:
            if byte_order == ByteOrder.LITTLE_ENDIAN:
                endian_string = "<"
                padded = bitarray(
                    self.bb[self.position : self.position + bit_length]
                ) + (64 - bit_length) * bitarray("0")
            else:
                endian_string = ">"
                padded = (64 - bit_length) * bitarray("0") + bitarray(
                    self.bb[self.position : self.position + bit_length]
                )
            if (
                byte_order == ByteOrder.BIG_ENDIAN_BYTE_SWAP
                or byte_order == ByteOrder.LITTLE_ENDIAN_BYTE_SWAP
            ):
                padded = padded[16:32] + padded[0:16] + padded[48:] + padded[32:48]
            result: int = struct.unpack(endian_string + "q", padded)[0]
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
            buffer = self.bb[self.position : self.position + bit_length]
            if (
                byte_order == ByteOrder.BIG_ENDIAN_BYTE_SWAP
                or byte_order == ByteOrder.LITTLE_ENDIAN_BYTE_SWAP
            ):
                buffer = buffer[16:] + buffer[0:16]
            result: float = struct.unpack(
                endianness + "f",
                buffer,
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
            buffer = self.bb[self.position : self.position + bit_length]
            if (
                byte_order == ByteOrder.BIG_ENDIAN_BYTE_SWAP
                or byte_order == ByteOrder.LITTLE_ENDIAN_BYTE_SWAP
            ):
                buffer = buffer[16:32] + buffer[0:16] + buffer[48:] + buffer[32:48]
            result: float = struct.unpack(
                endianness + "d",
                buffer,
            )[0]
            self.position += bit_length
            return result

    def read_complex(self, logical_name: str = "", read_function=None, **kwargs) -> Any:
        if isinstance(read_function, types.FunctionType):
            return read_function(logical_name=logical_name, read_buffer=self, **kwargs)

    def read_manual(self, logical_name: str = "", read_function=None, **kwargs) -> Any:
        if isinstance(read_function, types.FunctionType):
            return read_function()

    def read_enum(
        self,
        bit_length: int = -1,
        logical_name: str = "",
        read_function=None,
        **kwargs,
    ) -> Any:
        if isinstance(read_function, aenum._enum.EnumType):
            enum_return_value = read_function(
                ba2int(
                    self.bb[self.position : self.position + bit_length],
                    signed=False,
                )
            )
            return enum_return_value
        else:
            raise RuntimeError("read_enum called but read_function wasn't an enum")

    def read_complex_array(self, logical_name: str = "", **kwargs) -> List[Any]:
        raise NotImplementedError

    def read_str(self, bit_length: int = -1, logical_name: str = "", **kwargs) -> str:
        byte_order = kwargs.get("byte_order", self.byte_order)
        encoding = kwargs.get("encoding", "utf-8")
        result: bytes = struct.unpack(
            str(int(bit_length / 8)) + "s",
            self.bb[self.position : self.position + bit_length],
        )[0]
        self.position += bit_length
        result = result.decode(encoding=encoding, errors="strict")
        return result

    def read_array_field(
        self,
        logical_name: str = "",
        read_function=None,
        count: int = None,
        length: int = None,
        terminated=None,
        **kwargs,
    ) -> List[Any]:
        if count is not None:
            parsed_array = []
            for _ in range(count):
                parsed_array.append(
                    read_function(logical_name=logical_name, read_buffer=self, **kwargs)
                )
            return parsed_array
        elif length is not None:
            parsed_array = []
            end_pos = self.get_pos() + (length * 8)
            while self.get_pos() < end_pos:
                parsed_array.append(
                    read_function(logical_name=logical_name, read_buffer=self, **kwargs)
                )
            return parsed_array
        else:
            raise NotImplementedError
