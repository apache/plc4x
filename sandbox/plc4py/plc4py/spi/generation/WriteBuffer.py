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

from ctypes import *
from dataclasses import dataclass
from typing import List

from plc4py.spi.values.common import Serializable
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

    def write_bit(self, value: c_bool, logical_name: str = "", **kwargs) -> None:
        raise NotImplementedError

    def write_byte(self, value: c_byte, logical_name: str = "", **kwargs) -> None:
        self.write_signed_byte(value, logical_name, **kwargs)

    def write_byte_array(self, value: List[c_byte], logical_name: str = "", **kwargs) -> None:
        raise NotImplementedError

    def write_unsigned_byte(self, value: c_ubyte, logical_name: str = "", **kwargs) -> None:
        raise NotImplementedError

    def write_unsigned_short(self, value: c_uint16, logical_name: str = "", **kwargs) -> None:
        raise NotImplementedError

    def write_unsigned_int(self, value: c_uint32, logical_name: str = "", **kwargs) -> None:
        raise NotImplementedError

    def write_unsigned_long(self, value: c_uint64, logical_name: str = "", **kwargs) -> None:
        raise NotImplementedError

    def write_signed_byte(self, value: c_byte, logical_name: str = "", **kwargs) -> None:
        raise NotImplementedError

    def write_short(self, value: c_int16, logical_name: str = "", **kwargs) -> None:
        raise NotImplementedError

    def write_int(self, value: c_int32, logical_name: str = "", **kwargs) -> None:
        raise NotImplementedError

    def write_long(self, value: c_int64, logical_name: str = "", **kwargs) -> None:
        raise NotImplementedError

    def write_float(self, value: c_float, logical_name: str = "", **kwargs) -> None:
        raise NotImplementedError

    def write_double(self, value: c_double, logical_name: str = "", **kwargs) -> None:
        raise NotImplementedError

    def write_str(self, value: str, logical_name: str = "", bit_length: int = -1, **kwargs) -> None:
        raise NotImplementedError

    def write_virtual(self, value: str, logical_name: str = "", **kwargs) -> None:
        raise NotImplementedError

    #
    # This method can be used to influence serializing (e.g. intercept whole types and render them in a simplified form)
    #
    # @param value the value to be serialized
    # @throws SerializationException if something goes wrong
    #
    def write_serializable(self, value: Serializable) -> None:
        value.serialize(self)
