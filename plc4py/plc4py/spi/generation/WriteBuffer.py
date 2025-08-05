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
from abc import ABCMeta
from ctypes import (
    c_byte,
    c_double,
    c_float,
    c_int8,
    c_int16,
    c_int32,
    c_int64,
    c_ubyte,
    c_uint8,
    c_uint16,
    c_uint32,
    c_uint64,
)
from dataclasses import dataclass
from typing import List, Union

from bitarray import bitarray
from bitarray.util import zeros
from typing_extensions import override
import xml.etree.ElementTree as ET

from plc4py.api.exceptions.exceptions import SerializationException
from plc4py.api.messages.PlcMessage import PlcMessage
from plc4py.utils.GenericTypes import ByteOrder, ByteOrderAware
from xsdata.utils.text import camel_case


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

    def write_byte(
        self, value: int, bit_length: int = 8, logical_name: str = "", **kwargs
    ) -> None:
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
        self,
        value: str,
        bit_length: int = -1,
        logical_name: str = "",
        encoding: str = "UTF-8",
        **kwargs,
    ) -> None:
        raise NotImplementedError

    def write_complex_array(
        self, value: List[PlcMessage], logical_name: str = "", **kwargs
    ) -> None:
        raise NotImplementedError

    def write_manual(self, logical_name: str = "", write_function=None, **kwargs):
        if isinstance(write_function, types.FunctionType):
            return write_function()

    #
    # This method can be used to influence serializing (e.g. intercept whole types and render them in a simplified form)
    #
    # @param value the value to be serialized
    # @throws SerializationException if something goes wrong
    #
    def write_serializable(self, value, logical_name="") -> None:
        value.serialize(self)


class WriteBufferByteBased(WriteBuffer, metaclass=ABCMeta):
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
        # This refers to the bit alignment, which we always use big bit endianess
        self.bb = zeros(size * 8, endian=ByteOrder.get_short_name(ByteOrder.BIG_ENDIAN))
        self.byte_order = byte_order
        self.position: int = 0

    def get_bytes(self) -> memoryview:
        return memoryview(self.bb)

    def get_pos(self) -> int:
        return self.position

    def push_context(self, logical_name: str, **kwargs) -> None:
        # byte buffer need no context handling
        pass

    def pop_context(self, logical_name: str, **kwargs) -> None:
        # Byte Based Buffer doesn't need a context.
        pass

    def write_bit(self, value: bool, logical_name: str = "", **kwargs) -> None:
        self.bb[self.position] = value
        self.position += 1

    def write_byte(
        self, value: int, bit_length: int = 8, logical_name: str = "", **kwargs
    ) -> None:
        self.write_unsigned_byte(value, bit_length, logical_name, **kwargs)

    def write_byte_array(
        self, value: List[int], logical_name: str = "", **kwargs
    ) -> None:
        for a_byte in value:
            self.write_unsigned_byte(a_byte, 8, logical_name, **kwargs)

    def write_unsigned_byte(
        self, value: int, bit_length: int = 8, logical_name: str = "", **kwargs
    ) -> None:
        if bit_length <= 0:
            raise SerializationException("unsigned byte must contain at least 1 bit")
        elif bit_length > 8:
            raise SerializationException("unsigned byte can only contain max 8 bits")
        else:
            self._handle_numeric_encoding(
                int(value), bit_length, numeric_format="B", **kwargs
            )

    def write_unsigned_short(
        self,
        value: int,
        bit_length: int = 16,
        logical_name: str = "",
        **kwargs,
    ) -> None:
        if bit_length <= 0:
            raise SerializationException("unsigned short must contain at least 1 bit")
        elif bit_length > 16:
            raise SerializationException("unsigned short can only contain max 16 bits")
        else:
            self._handle_numeric_encoding(
                int(value), bit_length, numeric_format="H", **kwargs
            )

    def write_unsigned_int(
        self,
        value: int,
        bit_length: int = 32,
        logical_name: str = "",
        **kwargs,
    ) -> None:
        if bit_length <= 0:
            raise SerializationException("unsigned int must contain at least 1 bit")
        elif bit_length > 32:
            raise SerializationException("unsigned int can only contain max 32 bits")
        else:
            self._handle_numeric_encoding(
                int(value), bit_length, numeric_format="I", **kwargs
            )

    def write_unsigned_long(
        self,
        value: int,
        bit_length: int = 64,
        logical_name: str = "",
        **kwargs,
    ) -> None:
        if bit_length <= 0:
            raise SerializationException("unsigned long must contain at least 1 bit")
        elif bit_length > 64:
            raise SerializationException("unsigned long can only contain max 16 bits")
        else:
            self._handle_numeric_encoding(
                int(value), bit_length, numeric_format="Q", **kwargs
            )

    def write_signed_byte(
        self, value: int, bit_length: int = 8, logical_name: str = "", **kwargs
    ) -> None:
        if bit_length <= 0:
            raise SerializationException("Signed byte must contain at least 1 bit")
        elif bit_length > 8:
            raise SerializationException("Signed byte can only contain max 8 bits")
        self._handle_numeric_encoding(
            int(value), bit_length, numeric_format="b", **kwargs
        )

    def write_short(
        self,
        value: int,
        bit_length: int = 16,
        logical_name: str = "",
        **kwargs,
    ) -> None:
        if bit_length <= 0:
            raise SerializationException("Signed short must contain at least 1 bit")
        elif bit_length > 16:
            raise SerializationException("Signed short can only contain max 16 bits")
        self._handle_numeric_encoding(
            int(value), bit_length, numeric_format="h", **kwargs
        )

    def write_int(
        self,
        value: int,
        bit_length: int = 32,
        logical_name: str = "",
        **kwargs,
    ) -> None:
        if bit_length <= 0:
            raise SerializationException("Signed int must contain at least 1 bit")
        elif bit_length > 32:
            raise SerializationException("Signed int can only contain max 32 bits")
        self._handle_numeric_encoding(
            int(value), bit_length, numeric_format="i", **kwargs
        )

    def write_long(
        self,
        value: int,
        bit_length: int = 64,
        logical_name: str = "",
        **kwargs,
    ) -> None:
        if bit_length <= 0:
            raise SerializationException("Signed long must contain at least 1 bit")
        elif bit_length > 64:
            raise SerializationException("Signed long can only contain max 64 bits")
        self._handle_numeric_encoding(
            int(value), bit_length, numeric_format="q", **kwargs
        )

    def write_float(
        self,
        value: float,
        bit_length: int = 32,
        logical_name: str = "",
        **kwargs,
    ) -> None:
        if bit_length <= 0:
            raise SerializationException("Float must contain at least 1 bit")
        elif bit_length > 32:
            raise SerializationException("Float can only contain max 32 bits")
        self._handle_numeric_encoding(value, bit_length, numeric_format="f", **kwargs)

    def write_double(
        self,
        value: float,
        bit_length: int = 64,
        logical_name: str = "",
        **kwargs,
    ) -> None:
        if bit_length <= 0:
            raise SerializationException("Double must contain at least 1 bit")
        elif bit_length > 64:
            raise SerializationException("Double can only contain max 64 bits")
        self._handle_numeric_encoding(value, bit_length, numeric_format="d", **kwargs)

    def write_str(
        self,
        value: str,
        bit_length: int = -1,
        logical_name: str = "",
        encoding: str = "UTF-8",
        **kwargs,
    ) -> None:
        bit_order = kwargs.get("bit_order", ByteOrder.BIG_ENDIAN)
        src = bitarray(endian=ByteOrder.get_short_name(bit_order))
        src.frombytes(value.encode(encoding))
        self.bb[self.position : self.position + bit_length] = src[:bit_length]
        self.position += bit_length

    def write_complex_array(
        self, value: List[PlcMessage], logical_name: str = "", **kwargs
    ) -> None:
        for item in value:
            self.push_context(logical_name, **kwargs)
            self.write_serializable(item, logical_name="")
            self.pop_context(logical_name, **kwargs)

    def _handle_numeric_encoding(self, value, bit_length: int, **kwargs):
        bit_order = kwargs.get("bit_order", ByteOrder.BIG_ENDIAN)
        byte_order = kwargs.get("byte_order", self.byte_order)
        numeric_format = kwargs.get("numeric_format", "I")
        value_encoding: str = kwargs.get("encoding", "default")
        if value_encoding == "ASCII":
            if bit_length % 8 != 0:
                raise SerializationException(
                    "'ASCII' encoded fields must have a length that is a multiple of 8 bits long"
                )
            string_value: str = "{}".format(value.value)
            src = bitarray(endian=ByteOrder.get_short_name(bit_order))
            src.frombytes(bytearray(string_value, value_encoding))
            self.bb[self.position : self.position + bit_length] = src[:bit_length]
            self.position += bit_length
        elif value_encoding == "default":
            src = bitarray(endian=ByteOrder.get_short_name(bit_order))
            endianness: str = ">"
            if byte_order == ByteOrder.LITTLE_ENDIAN:
                endianness = "<"
            if not isinstance(value, int):
                pass

            result: bytes = struct.pack(
                endianness + numeric_format,
                value,
            )

            src.frombytes(result)
            if (
                byte_order == ByteOrder.BIG_ENDIAN_BYTE_SWAP
                or byte_order == ByteOrder.LITTLE_ENDIAN_BYTE_SWAP
            ):
                if 32 < bit_length <= 64:
                    src = src[16:32] + src[0:16] + src[48:] + src[32:48]
                elif 16 < bit_length <= 32:
                    src = src[16:] + src[0:16]
            if bit_length < 8:
                self.bb[self.position : self.position + bit_length] = src[-bit_length:]
            else:
                self.bb[self.position : self.position + bit_length] = src[:bit_length]
            self.position += bit_length


class WriteBufferXmlBased(WriteBuffer, metaclass=ABCMeta):
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

    def __init__(self, byte_order: ByteOrder = ByteOrder.BIG_ENDIAN):
        # This refers to the bit alignment, which we always use big bit endianess
        self.byte_order = byte_order
        self.position = 0
        self.stack = []

    def get_bytes(self) -> memoryview:
        return memoryview(bytearray([]))

    def get_pos(self) -> int:
        return self.position

    def push_context(self, logical_name: str, **kwargs) -> None:
        new_element: ET.Element = ET.Element(logical_name)
        if len(self.stack) > 0:
            self.stack[-1].append(new_element)
        else:
            self.stack.append(new_element)
        self.stack.append(new_element)

    def pop_context(self, logical_name: str, **kwargs) -> None:
        self.stack.pop()

    def write_bit(self, value: bool, logical_name: str = "", **kwargs) -> None:
        data_type: str = "bit"
        data: str = str(value)
        if "bit_length" in kwargs:
            kwargs["bit_length"] = str(kwargs["bit_length"])
        else:
            kwargs["bit_length"] = str(1)
        self._create_and_append(camel_case(logical_name), data_type, data, **kwargs)

    def write_byte(self, value: int, logical_name: str = "", **kwargs) -> None:
        data_type: str = "byte"
        data: str = str(value)
        if "bit_length" in kwargs:
            kwargs["bit_length"] = str(kwargs["bit_length"])
        else:
            kwargs["bit_length"] = str(8)
        self._create_and_append(camel_case(logical_name), data_type, data, **kwargs)

    def write_byte_array(
        self, value: List[int], logical_name: str = "", **kwargs
    ) -> None:
        data_type: str = "byte"
        data: str = str(value)
        if "bit_length" in kwargs:
            kwargs["bit_length"] = str(kwargs["bit_length"])
        else:
            kwargs["bit_length"] = str(len(value) * 8)
        self._create_and_append(camel_case(logical_name), data_type, data, **kwargs)

    def write_unsigned_byte(
        self, value: int, bit_length: int = 8, logical_name: str = "", **kwargs
    ) -> None:
        data_type: str = "uint"
        data: str = str(value)
        bit_length: str = str(bit_length)
        self._create_and_append(
            camel_case(logical_name), data_type, data, bit_length, **kwargs
        )

    def write_unsigned_short(
        self,
        value: int,
        bit_length: int = 16,
        logical_name: str = "",
        **kwargs,
    ) -> None:
        data_type: str = "uint"
        data: str = str(value)
        bit_length: str = str(bit_length)
        self._create_and_append(
            camel_case(logical_name), data_type, data, bit_length, **kwargs
        )

    def write_unsigned_int(
        self,
        value: int,
        bit_length: int = 32,
        logical_name: str = "",
        **kwargs,
    ) -> None:
        data_type: str = "udint"
        data: str = str(value)
        bit_length: str = str(bit_length)
        self._create_and_append(
            camel_case(logical_name), data_type, data, bit_length, **kwargs
        )

    def write_unsigned_long(
        self,
        value: int,
        bit_length: int = 64,
        logical_name: str = "",
        **kwargs,
    ) -> None:
        data_type: str = "ulint"
        data: str = str(value)
        bit_length: str = str(bit_length)
        self._create_and_append(
            camel_case(logical_name), data_type, data, bit_length, **kwargs
        )

    def write_signed_byte(
        self, value: int, bit_length: int = 8, logical_name: str = "", **kwargs
    ) -> None:
        data_type: str = "byte"
        data: str = str(value)
        bit_length: str = str(bit_length)
        self._create_and_append(
            camel_case(logical_name), data_type, data, bit_length, **kwargs
        )

    def write_short(
        self,
        value: int,
        bit_length: int = 16,
        logical_name: str = "",
        **kwargs,
    ) -> None:
        data_type: str = "int"
        data: str = str(value)
        bit_length: str = str(bit_length)
        self._create_and_append(
            camel_case(logical_name), data_type, data, bit_length, **kwargs
        )

    def write_int(
        self,
        value: int,
        bit_length: int = 32,
        logical_name: str = "",
        **kwargs,
    ) -> None:
        data_type: str = "dint"
        data: str = str(value)
        bit_length: str = str(bit_length)
        self._create_and_append(
            camel_case(logical_name), data_type, data, bit_length, **kwargs
        )

    def write_long(
        self,
        value: int,
        bit_length: int = 64,
        logical_name: str = "",
        **kwargs,
    ) -> None:
        data_type: str = "lint"
        data: str = str(value)
        bit_length: str = str(bit_length)
        self._create_and_append(
            camel_case(logical_name), data_type, data, bit_length, **kwargs
        )

    def write_float(
        self,
        value: float,
        bit_length: int = 32,
        logical_name: str = "",
        **kwargs,
    ) -> None:
        data_type: str = "real"
        data: str = str(value)
        bit_length: str = str(bit_length)
        self._create_and_append(
            camel_case(logical_name), data_type, data, bit_length, **kwargs
        )

    def write_double(
        self,
        value: float,
        bit_length: int = 64,
        logical_name: str = "",
        **kwargs,
    ) -> None:
        data_type: str = "lreal"
        data: str = str(value)
        bit_length: str = str(bit_length)
        self._create_and_append(
            camel_case(logical_name), data_type, data, bit_length, **kwargs
        )

    def write_complex_array(
        self, value: List[PlcMessage], logical_name: str = "", **kwargs
    ) -> None:
        for _ in value:
            new_element: ET.Element = ET.Element(logical_name)
            self.stack[-1].append(new_element)
            self.stack.append(new_element)

    def _create_and_append(
        self, logical_name: str, data_type: str, data: str, bit_length: str, **kwargs
    ) -> None:
        new_element: ET.Element = ET.Element(logical_name)
        new_element.set("dataType", data_type)
        new_element.set("bitlength", bit_length)
        new_element.text = data
        if len(self.stack) > 0:
            self.stack[-1].append(new_element)
        else:
            self.stack.append(new_element)

    def to_xml_string(self) -> str:
        ET.indent(self.stack[0], space="\t", level=0)
        return ET.tostring(self.stack[0]).decode("utf-8")
