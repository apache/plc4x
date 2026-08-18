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

from plc4py.api.value.PlcValue import PlcValue
from plc4py.protocols.modbus.readwrite.ModbusDataType import ModbusDataType
from plc4py.spi.generation.ReadBuffer import ReadBuffer
from plc4py.spi.generation.WriteBuffer import WriteBuffer
from plc4py.spi.values.PlcValues import PlcBOOL
from plc4py.spi.values.PlcValues import PlcBYTE
from plc4py.spi.values.PlcValues import PlcCHAR
from plc4py.spi.values.PlcValues import PlcDINT
from plc4py.spi.values.PlcValues import PlcDWORD
from plc4py.spi.values.PlcValues import PlcINT
from plc4py.spi.values.PlcValues import PlcLINT
from plc4py.spi.values.PlcValues import PlcLREAL
from plc4py.spi.values.PlcValues import PlcLWORD
from plc4py.spi.values.PlcValues import PlcREAL
from plc4py.spi.values.PlcValues import PlcSINT
from plc4py.spi.values.PlcValues import PlcSTRING
from plc4py.spi.values.PlcValues import PlcUDINT
from plc4py.spi.values.PlcValues import PlcUINT
from plc4py.spi.values.PlcValues import PlcULINT
from plc4py.spi.values.PlcValues import PlcUSINT
from plc4py.spi.values.PlcValues import PlcWCHAR
from plc4py.spi.values.PlcValues import PlcWORD
from plc4py.spi.values.PlcValues import PlcWSTRING
from plc4py.utils.GenericTypes import ByteOrder
import math


class DataItem:
    @staticmethod
    def static_parse(
        read_buffer: ReadBuffer, data_type: ModbusDataType, string_length: int
    ):
        if data_type == ModbusDataType.BOOL:  # BOOL

            # Simple Field (value)
            value: bool = read_buffer.read_bit("")

            return PlcBOOL(value)
        if data_type == ModbusDataType.BYTE:  # BYTE

            # Simple Field (value)
            value: int = read_buffer.read_unsigned_short(8, logical_name="")

            return PlcBYTE(value)
        if data_type == ModbusDataType.WORD:  # WORD

            # Simple Field (value)
            value: int = read_buffer.read_unsigned_int(16, logical_name="")

            return PlcWORD(value)
        if data_type == ModbusDataType.DWORD:  # DWORD

            # Simple Field (value)
            value: int = read_buffer.read_unsigned_long(32, logical_name="")

            return PlcDWORD(value)
        if data_type == ModbusDataType.LWORD:  # LWORD

            # Simple Field (value)
            value: int = read_buffer.read_unsigned_long(64, logical_name="")

            return PlcLWORD(value)
        if data_type == ModbusDataType.SINT:  # SINT

            # Simple Field (value)
            value: int = read_buffer.read_signed_byte(8, logical_name="")

            return PlcSINT(value)
        if data_type == ModbusDataType.INT:  # INT

            # Simple Field (value)
            value: int = read_buffer.read_short(16, logical_name="")

            return PlcINT(value)
        if data_type == ModbusDataType.DINT:  # DINT

            # Simple Field (value)
            value: int = read_buffer.read_int(32, logical_name="")

            return PlcDINT(value)
        if data_type == ModbusDataType.LINT:  # LINT

            # Simple Field (value)
            value: int = read_buffer.read_long(64, logical_name="")

            return PlcLINT(value)
        if data_type == ModbusDataType.USINT:  # USINT

            # Simple Field (value)
            value: int = read_buffer.read_unsigned_short(8, logical_name="")

            return PlcUSINT(value)
        if data_type == ModbusDataType.UINT:  # UINT

            # Simple Field (value)
            value: int = read_buffer.read_unsigned_int(16, logical_name="")

            return PlcUINT(value)
        if data_type == ModbusDataType.UDINT:  # UDINT

            # Simple Field (value)
            value: int = read_buffer.read_unsigned_long(32, logical_name="")

            return PlcUDINT(value)
        if data_type == ModbusDataType.ULINT:  # ULINT

            # Simple Field (value)
            value: int = read_buffer.read_unsigned_long(64, logical_name="")

            return PlcULINT(value)
        if data_type == ModbusDataType.REAL:  # REAL

            # Simple Field (value)
            value: float = read_buffer.read_float(32, logical_name="")

            return PlcREAL(value)
        if data_type == ModbusDataType.LREAL:  # LREAL

            # Simple Field (value)
            value: float = read_buffer.read_double(64, logical_name="")

            return PlcLREAL(value)
        if data_type == ModbusDataType.CHAR:  # CHAR

            # Simple Field (value)
            value: str = read_buffer.read_str(8, logical_name="")

            return PlcCHAR(value)
        if data_type == ModbusDataType.WCHAR:  # WCHAR

            # Simple Field (value)
            value: str = read_buffer.read_str(16, logical_name="")

            return PlcWCHAR(value)
        if data_type == ModbusDataType.STRING:  # STRING

            # Simple Field (value)
            value: str = read_buffer.read_str(-1, logical_name="")

            return PlcSTRING(value)
        if data_type == ModbusDataType.WSTRING:  # WSTRING

            # Simple Field (value)
            value: str = read_buffer.read_str(-1, logical_name="")

            return PlcWSTRING(value)
        return None

    @staticmethod
    def static_serialize(
        write_buffer: WriteBuffer,
        _value: PlcValue,
        data_type: ModbusDataType,
        string_length: int,
        byte_order: ByteOrder,
    ) -> None:
        if data_type == ModbusDataType.BOOL:  # BOOL
            # Simple Field (value)
            value: bool = _value.get_bool()
            write_buffer.write_bit((value), "value")

        elif data_type == ModbusDataType.BYTE:  # BYTE
            # Simple Field (value)
            value: int = _value.get_int()
            write_buffer.write_byte((value), 8, "value")

        elif data_type == ModbusDataType.WORD:  # WORD
            # Simple Field (value)
            value: int = _value.get_int()
            write_buffer.write_unsigned_short((value), 16, "value")

        elif data_type == ModbusDataType.DWORD:  # DWORD
            # Simple Field (value)
            value: int = _value.get_int()
            write_buffer.write_unsigned_int((value), 32, "value")

        elif data_type == ModbusDataType.LWORD:  # LWORD
            # Simple Field (value)
            value: int = _value.get_int()
            write_buffer.write_unsigned_long((value), 64, "value")

        elif data_type == ModbusDataType.SINT:  # SINT
            # Simple Field (value)
            value: int = _value.get_int()
            write_buffer.write_signed_byte((value), 8, "value")

        elif data_type == ModbusDataType.INT:  # INT
            # Simple Field (value)
            value: int = _value.get_int()
            write_buffer.write_short((value), 16, "value")

        elif data_type == ModbusDataType.DINT:  # DINT
            # Simple Field (value)
            value: int = _value.get_int()
            write_buffer.write_int((value), 32, "value")

        elif data_type == ModbusDataType.LINT:  # LINT
            # Simple Field (value)
            value: int = _value.get_int()
            write_buffer.write_long((value), 64, "value")

        elif data_type == ModbusDataType.USINT:  # USINT
            # Simple Field (value)
            value: int = _value.get_int()
            write_buffer.write_byte((value), 8, "value")

        elif data_type == ModbusDataType.UINT:  # UINT
            # Simple Field (value)
            value: int = _value.get_int()
            write_buffer.write_unsigned_short((value), 16, "value")

        elif data_type == ModbusDataType.UDINT:  # UDINT
            # Simple Field (value)
            value: int = _value.get_int()
            write_buffer.write_unsigned_int((value), 32, "value")

        elif data_type == ModbusDataType.ULINT:  # ULINT
            # Simple Field (value)
            value: int = _value.get_int()
            write_buffer.write_unsigned_long((value), 64, "value")

        elif data_type == ModbusDataType.REAL:  # REAL
            # Simple Field (value)
            value: float = _value.get_float()
            write_buffer.write_float((value), 32, "value")

        elif data_type == ModbusDataType.LREAL:  # LREAL
            # Simple Field (value)
            value: float = _value.get_float()
            write_buffer.write_double((value), 64, "value")

        elif data_type == ModbusDataType.CHAR:  # CHAR
            # Simple Field (value)
            value: str = _value.get_str()
            write_buffer.write_str((value), 8, "value", "UTF-8")

        elif data_type == ModbusDataType.WCHAR:  # WCHAR
            # Simple Field (value)
            value: str = _value.get_str()
            write_buffer.write_str((value), 16, "value", "UTF-8")

        elif data_type == ModbusDataType.STRING:  # STRING
            # Simple Field (value)
            value: str = _value.get_str()
            write_buffer.write_str((value), string_length * (8), "value", "UTF-8")

        elif data_type == ModbusDataType.WSTRING:  # WSTRING
            # Simple Field (value)
            value: str = _value.get_str()
            write_buffer.write_str((value), string_length * (16), "value", "UTF-8")

    @staticmethod
    def get_length_in_bytes(
        _value: PlcValue, data_type: ModbusDataType, string_length: int
    ) -> int:
        return int(
            math.ceil(
                float(DataItem.get_length_in_bits(_value, data_type, string_length))
                / 8.0
            )
        )

    @staticmethod
    def get_length_in_bits(
        _value: PlcValue, data_type: ModbusDataType, string_length: int
    ) -> int:
        size_in_bits: int = 0
        if data_type == ModbusDataType.BOOL:  # BOOL
            # Simple Field (value)
            size_in_bits += 1
        elif data_type == ModbusDataType.BYTE:  # BYTE
            # Simple Field (value)
            size_in_bits += 8
        elif data_type == ModbusDataType.WORD:  # WORD
            # Simple Field (value)
            size_in_bits += 16
        elif data_type == ModbusDataType.DWORD:  # DWORD
            # Simple Field (value)
            size_in_bits += 32
        elif data_type == ModbusDataType.LWORD:  # LWORD
            # Simple Field (value)
            size_in_bits += 64
        elif data_type == ModbusDataType.SINT:  # SINT
            # Simple Field (value)
            size_in_bits += 8
        elif data_type == ModbusDataType.INT:  # INT
            # Simple Field (value)
            size_in_bits += 16
        elif data_type == ModbusDataType.DINT:  # DINT
            # Simple Field (value)
            size_in_bits += 32
        elif data_type == ModbusDataType.LINT:  # LINT
            # Simple Field (value)
            size_in_bits += 64
        elif data_type == ModbusDataType.USINT:  # USINT
            # Simple Field (value)
            size_in_bits += 8
        elif data_type == ModbusDataType.UINT:  # UINT
            # Simple Field (value)
            size_in_bits += 16
        elif data_type == ModbusDataType.UDINT:  # UDINT
            # Simple Field (value)
            size_in_bits += 32
        elif data_type == ModbusDataType.ULINT:  # ULINT
            # Simple Field (value)
            size_in_bits += 64
        elif data_type == ModbusDataType.REAL:  # REAL
            # Simple Field (value)
            size_in_bits += 32
        elif data_type == ModbusDataType.LREAL:  # LREAL
            # Simple Field (value)
            size_in_bits += 64
        elif data_type == ModbusDataType.CHAR:  # CHAR
            # Simple Field (value)
            size_in_bits += 8
        elif data_type == ModbusDataType.WCHAR:  # WCHAR
            # Simple Field (value)
            size_in_bits += 16
        elif data_type == ModbusDataType.STRING:  # STRING
            # Simple Field (value)
            size_in_bits += -1
        elif data_type == ModbusDataType.WSTRING:  # WSTRING
            # Simple Field (value)
            size_in_bits += -1

        return size_in_bits
