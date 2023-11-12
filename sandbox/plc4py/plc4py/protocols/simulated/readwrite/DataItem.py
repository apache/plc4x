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

from abc import staticmethod
from loguru import logging as log
from plc4py.api.value.PlcValue import PlcValue
from plc4py.spi.generation.ReadBuffer import ReadBuffer
from plc4py.spi.generation.WriteBuffer import WriteBuffer
from typing import List
import math


class DataItem:
    @staticmethod
    def static_parse(read_buffer: ReadBuffer, data_type: str, number_of_values: int):
        if data_type == "_bool" and number_of_values == int(1):  # BOOL
            # Simple Field (value)
            value: bool = read_buffer.read_bit("")

            return PlcBOOL(value)
        if data_type == "_bool":  # List
            # Array field (value)
            # Count array
            item_count: int = int(number_of_values)
            value: List[PlcValue] = []
            for cur_item in range(item_count):
                value.append(PlcBOOL(bool(read_buffer.read_bit(""))))

            return PlcList(value)
        if data_type == "_byte" and number_of_values == int(1):  # BYTE
            # Simple Field (value)
            value: int = read_buffer.read_unsigned_short(8, logical_name="")

            return PlcBYTE(value)
        if data_type == "_byte":  # List
            # Array field (value)
            # Count array
            item_count: int = int(number_of_values)
            value: List[PlcValue] = []
            for cur_item in range(item_count):
                value.append(
                    PlcUINT(int(read_buffer.read_unsigned_short(8, logical_name="")))
                )

            return PlcList(value)
        if data_type == "_word" and number_of_values == int(1):  # WORD
            # Simple Field (value)
            value: int = read_buffer.read_unsigned_int(16, logical_name="")

            return PlcWORD(value)
        if data_type == "_word":  # List
            # Array field (value)
            # Count array
            item_count: int = int(number_of_values)
            value: List[PlcValue] = []
            for cur_item in range(item_count):
                value.append(
                    PlcUDINT(int(read_buffer.read_unsigned_int(16, logical_name="")))
                )

            return PlcList(value)
        if data_type == "_dword" and number_of_values == int(1):  # DWORD
            # Simple Field (value)
            value: int = read_buffer.read_unsigned_long(32, logical_name="")

            return PlcDWORD(value)
        if data_type == "_dword":  # List
            # Array field (value)
            # Count array
            item_count: int = int(number_of_values)
            value: List[PlcValue] = []
            for cur_item in range(item_count):
                value.append(
                    PlcULINT(int(read_buffer.read_unsigned_long(32, logical_name="")))
                )

            return PlcList(value)
        if data_type == "_lword" and number_of_values == int(1):  # LWORD
            # Simple Field (value)
            value: int = read_buffer.read_unsigned_big_integer(64, logical_name="")

            return PlcLWORD(value)
        if data_type == "_lword":  # List
            # Array field (value)
            # Count array
            item_count: int = int(number_of_values)
            value: List[PlcValue] = []
            for cur_item in range(item_count):
                value.append(
                    PlcLINT(
                        int(read_buffer.read_unsigned_big_integer(64, logical_name=""))
                    )
                )

            return PlcList(value)
        if data_type == "_sint" and number_of_values == int(1):  # SINT
            # Simple Field (value)
            value: int = read_buffer.read_signed_byte(8, logical_name="")

            return PlcSINT(value)
        if data_type == "_sint":  # List
            # Array field (value)
            # Count array
            item_count: int = int(number_of_values)
            value: List[PlcValue] = []
            for cur_item in range(item_count):
                value.append(
                    PlcSINT(int(read_buffer.read_signed_byte(8, logical_name="")))
                )

            return PlcList(value)
        if data_type == "_int" and number_of_values == int(1):  # INT
            # Simple Field (value)
            value: int = read_buffer.read_short(16, logical_name="")

            return PlcINT(value)
        if data_type == "_int":  # List
            # Array field (value)
            # Count array
            item_count: int = int(number_of_values)
            value: List[PlcValue] = []
            for cur_item in range(item_count):
                value.append(PlcINT(int(read_buffer.read_short(16, logical_name=""))))

            return PlcList(value)
        if data_type == "_dint" and number_of_values == int(1):  # DINT
            # Simple Field (value)
            value: int = read_buffer.read_int(32, logical_name="")

            return PlcDINT(value)
        if data_type == "_dint":  # List
            # Array field (value)
            # Count array
            item_count: int = int(number_of_values)
            value: List[PlcValue] = []
            for cur_item in range(item_count):
                value.append(PlcDINT(int(read_buffer.read_int(32, logical_name=""))))

            return PlcList(value)
        if data_type == "_lint" and number_of_values == int(1):  # LINT
            # Simple Field (value)
            value: int = read_buffer.read_long(64, logical_name="")

            return PlcLINT(value)
        if data_type == "_lint":  # List
            # Array field (value)
            # Count array
            item_count: int = int(number_of_values)
            value: List[PlcValue] = []
            for cur_item in range(item_count):
                value.append(PlcLINT(int(read_buffer.read_long(64, logical_name=""))))

            return PlcList(value)
        if data_type == "_usint" and number_of_values == int(1):  # USINT
            # Simple Field (value)
            value: int = read_buffer.read_unsigned_short(8, logical_name="")

            return PlcUSINT(value)
        if data_type == "_usint":  # List
            # Array field (value)
            # Count array
            item_count: int = int(number_of_values)
            value: List[PlcValue] = []
            for cur_item in range(item_count):
                value.append(
                    PlcUINT(int(read_buffer.read_unsigned_short(8, logical_name="")))
                )

            return PlcList(value)
        if data_type == "_uint" and number_of_values == int(1):  # UINT
            # Simple Field (value)
            value: int = read_buffer.read_unsigned_int(16, logical_name="")

            return PlcUINT(value)
        if data_type == "_uint":  # List
            # Array field (value)
            # Count array
            item_count: int = int(number_of_values)
            value: List[PlcValue] = []
            for cur_item in range(item_count):
                value.append(
                    PlcUDINT(int(read_buffer.read_unsigned_int(16, logical_name="")))
                )

            return PlcList(value)
        if data_type == "_udint" and number_of_values == int(1):  # UDINT
            # Simple Field (value)
            value: int = read_buffer.read_unsigned_long(32, logical_name="")

            return PlcUDINT(value)
        if data_type == "_udint":  # List
            # Array field (value)
            # Count array
            item_count: int = int(number_of_values)
            value: List[PlcValue] = []
            for cur_item in range(item_count):
                value.append(
                    PlcULINT(int(read_buffer.read_unsigned_long(32, logical_name="")))
                )

            return PlcList(value)
        if data_type == "_ulint" and number_of_values == int(1):  # ULINT
            # Simple Field (value)
            value: int = read_buffer.read_unsigned_big_integer(64, logical_name="")

            return PlcULINT(value)
        if data_type == "_ulint":  # List
            # Array field (value)
            # Count array
            item_count: int = int(number_of_values)
            value: List[PlcValue] = []
            for cur_item in range(item_count):
                value.append(
                    PlcLINT(
                        int(read_buffer.read_unsigned_big_integer(64, logical_name=""))
                    )
                )

            return PlcList(value)
        if data_type == "_real" and number_of_values == int(1):  # REAL
            # Simple Field (value)
            value: float = read_buffer.read_float(32, logical_name="")

            return PlcREAL(value)
        if data_type == "_real":  # List
            # Array field (value)
            # Count array
            item_count: int = int(number_of_values)
            value: List[PlcValue] = []
            for cur_item in range(item_count):
                value.append(
                    PlcREAL(float(read_buffer.read_float(32, logical_name="")))
                )

            return PlcList(value)
        if data_type == "_lreal" and number_of_values == int(1):  # LREAL
            # Simple Field (value)
            value: float = read_buffer.read_double(64, logical_name="")

            return PlcLREAL(value)
        if data_type == "_lreal":  # List
            # Array field (value)
            # Count array
            item_count: int = int(number_of_values)
            value: List[PlcValue] = []
            for cur_item in range(item_count):
                value.append(
                    PlcLREAL(float(read_buffer.read_double(64, logical_name="")))
                )

            return PlcList(value)
        if data_type == "_char" and number_of_values == int(1):  # CHAR
            # Simple Field (value)
            value: str = read_buffer.read_string(8, logical_name="", encoding="")

            return PlcCHAR(value)
        if data_type == "_char":  # List
            # Array field (value)
            # Count array
            item_count: int = int(number_of_values)
            value: List[PlcValue] = []
            for cur_item in range(item_count):
                value.append(
                    PlcSTRING(
                        str(read_buffer.read_string(8, logical_name="", encoding=""))
                    )
                )

            return PlcList(value)
        if data_type == "_wchar" and number_of_values == int(1):  # WCHAR
            # Simple Field (value)
            value: str = read_buffer.read_string(16, logical_name="", encoding="")

            return PlcWCHAR(value)
        if data_type == "_wchar":  # List
            # Array field (value)
            # Count array
            item_count: int = int(number_of_values)
            value: List[PlcValue] = []
            for cur_item in range(item_count):
                value.append(
                    PlcSTRING(
                        str(read_buffer.read_string(16, logical_name="", encoding=""))
                    )
                )

            return PlcList(value)
        if data_type == "_string":  # STRING
            # Simple Field (value)
            value: str = read_buffer.read_string(255, logical_name="", encoding="")

            return PlcSTRING(value)
        if data_type == "_wstring":  # STRING
            # Simple Field (value)
            value: str = read_buffer.read_string(255, logical_name="", encoding="")

            return PlcSTRING(value)
        return None

    @staticmethod
    def static_serialize(
        write_buffer: WriteBuffer,
        _value: PlcValue,
        data_type: str,
        number_of_values: int,
        byte_order: ByteOrder,
    ) -> None:
        if data_type == "BOOL" and number_of_values == int(1):  # BOOL
            # Simple Field (value)
            value: bool = _value.get_bool()
            write_buffer.write_bit((value), "value")
        if data_type == "BOOL":  # List
            values: PlcList = _value

            for val in values.getList():
                value: bool = val.get_bool()
                write_buffer.write_bit((value), "value")

        if data_type == "BYTE" and number_of_values == int(1):  # BYTE
            # Simple Field (value)
            value: int = _value.get_int()
            write_buffer.write_byte((value), 8, "value")
        if data_type == "BYTE":  # List
            values: PlcList = _value

            for val in values.getList():
                value: int = val.get_int()
                write_buffer.write_byte((value), 8, "value")

        if data_type == "WORD" and number_of_values == int(1):  # WORD
            # Simple Field (value)
            value: int = _value.get_int()
            write_buffer.write_unsigned_short((value), 16, "value")
        if data_type == "WORD":  # List
            values: PlcList = _value

            for val in values.getList():
                value: int = val.get_int()
                write_buffer.write_unsigned_short((value), 16, "value")

        if data_type == "DWORD" and number_of_values == int(1):  # DWORD
            # Simple Field (value)
            value: int = _value.get_int()
            write_buffer.write_unsigned_int((value), 32, "value")
        if data_type == "DWORD":  # List
            values: PlcList = _value

            for val in values.getList():
                value: int = val.get_int()
                write_buffer.write_unsigned_int((value), 32, "value")

        if data_type == "LWORD" and number_of_values == int(1):  # LWORD
            # Simple Field (value)
            value: int = _value.get_int()
            write_buffer.write_unsigned_long((value), 64, "value")
        if data_type == "LWORD":  # List
            values: PlcList = _value

            for val in values.getList():
                value: int = val.get_int()
                write_buffer.write_unsigned_long((value), 64, "value")

        if data_type == "SINT" and number_of_values == int(1):  # SINT
            # Simple Field (value)
            value: int = _value.get_int()
            write_buffer.write_signed_byte((value), 8, "value")
        if data_type == "SINT":  # List
            values: PlcList = _value

            for val in values.getList():
                value: int = val.get_int()
                write_buffer.write_signed_byte((value), 8, "value")

        if data_type == "INT" and number_of_values == int(1):  # INT
            # Simple Field (value)
            value: int = _value.get_int()
            write_buffer.write_short((value), 16, "value")
        if data_type == "INT":  # List
            values: PlcList = _value

            for val in values.getList():
                value: int = val.get_int()
                write_buffer.write_short((value), 16, "value")

        if data_type == "DINT" and number_of_values == int(1):  # DINT
            # Simple Field (value)
            value: int = _value.get_int()
            write_buffer.write_int((value), 32, "value")
        if data_type == "DINT":  # List
            values: PlcList = _value

            for val in values.getList():
                value: int = val.get_int()
                write_buffer.write_int((value), 32, "value")

        if data_type == "LINT" and number_of_values == int(1):  # LINT
            # Simple Field (value)
            value: int = _value.get_int()
            write_buffer.write_long((value), 64, "value")
        if data_type == "LINT":  # List
            values: PlcList = _value

            for val in values.getList():
                value: int = val.get_int()
                write_buffer.write_long((value), 64, "value")

        if data_type == "USINT" and number_of_values == int(1):  # USINT
            # Simple Field (value)
            value: int = _value.get_int()
            write_buffer.write_byte((value), 8, "value")
        if data_type == "USINT":  # List
            values: PlcList = _value

            for val in values.getList():
                value: int = val.get_int()
                write_buffer.write_byte((value), 8, "value")

        if data_type == "UINT" and number_of_values == int(1):  # UINT
            # Simple Field (value)
            value: int = _value.get_int()
            write_buffer.write_unsigned_short((value), 16, "value")
        if data_type == "UINT":  # List
            values: PlcList = _value

            for val in values.getList():
                value: int = val.get_int()
                write_buffer.write_unsigned_short((value), 16, "value")

        if data_type == "UDINT" and number_of_values == int(1):  # UDINT
            # Simple Field (value)
            value: int = _value.get_int()
            write_buffer.write_unsigned_int((value), 32, "value")
        if data_type == "UDINT":  # List
            values: PlcList = _value

            for val in values.getList():
                value: int = val.get_int()
                write_buffer.write_unsigned_int((value), 32, "value")

        if data_type == "ULINT" and number_of_values == int(1):  # ULINT
            # Simple Field (value)
            value: int = _value.get_int()
            write_buffer.write_unsigned_long((value), 64, "value")
        if data_type == "ULINT":  # List
            values: PlcList = _value

            for val in values.getList():
                value: int = val.get_int()
                write_buffer.write_unsigned_long((value), 64, "value")

        if data_type == "REAL" and number_of_values == int(1):  # REAL
            # Simple Field (value)
            value: float = _value.get_float()
            write_buffer.write_float((value), 32, "value")
        if data_type == "REAL":  # List
            values: PlcList = _value

            for val in values.getList():
                value: float = val.get_float()
                write_buffer.write_float((value), 32, "value")

        if data_type == "LREAL" and number_of_values == int(1):  # LREAL
            # Simple Field (value)
            value: float = _value.get_float()
            write_buffer.write_double((value), 64, "value")
        if data_type == "LREAL":  # List
            values: PlcList = _value

            for val in values.getList():
                value: float = val.get_float()
                write_buffer.write_double((value), 64, "value")

        if data_type == "CHAR" and number_of_values == int(1):  # CHAR
            # Simple Field (value)
            value: str = _value.get_str()
            write_buffer.write_str((value), uint32(8), "UTF-8", "value")
        if data_type == "CHAR":  # List
            values: PlcList = _value

            for val in values.getList():
                value: str = val.get_str()
                write_buffer.write_str((value), uint32(8), "UTF-8", "value")

        if data_type == "WCHAR" and number_of_values == int(1):  # WCHAR
            # Simple Field (value)
            value: str = _value.get_str()
            write_buffer.write_str((value), uint32(16), "UTF-16", "value")
        if data_type == "WCHAR":  # List
            values: PlcList = _value

            for val in values.getList():
                value: str = val.get_str()
                write_buffer.write_str((value), uint32(16), "UTF-16", "value")

        if data_type == "STRING":  # STRING
            # Simple Field (value)
            value: str = _value.get_str()
            write_buffer.write_str((value), uint32(255), "UTF-8", "value")
        if data_type == "WSTRING":  # STRING
            # Simple Field (value)
            value: str = _value.get_str()
            write_buffer.write_str((value), uint32(255), "UTF-16", "value")

    @staticmethod
    def get_length_in_bytes(
        _value: PlcValue, data_type: str, number_of_values: int
    ) -> int:
        return int(
            math.ceil(
                float(get_length_in_bits(_value, data_type, number_of_values)) / 8.0
            )
        )

    @staticmethod
    def get_length_in_bits(
        _value: PlcValue, data_type: str, number_of_values: int
    ) -> int:
        size_in_bits: int = 0
        if data_type == "BOOL" and number_of_values == int(1):  # BOOL
            # Simple Field (value)
            size_in_bits += 1
        if data_type == "BOOL":  # List
            values: PlcList = _value
            size_in_bits += values.get_list().size() * 1
        if data_type == "BYTE" and number_of_values == int(1):  # BYTE
            # Simple Field (value)
            size_in_bits += 8
        if data_type == "BYTE":  # List
            values: PlcList = _value
            size_in_bits += values.get_list().size() * 8
        if data_type == "WORD" and number_of_values == int(1):  # WORD
            # Simple Field (value)
            size_in_bits += 16
        if data_type == "WORD":  # List
            values: PlcList = _value
            size_in_bits += values.get_list().size() * 16
        if data_type == "DWORD" and number_of_values == int(1):  # DWORD
            # Simple Field (value)
            size_in_bits += 32
        if data_type == "DWORD":  # List
            values: PlcList = _value
            size_in_bits += values.get_list().size() * 32
        if data_type == "LWORD" and number_of_values == int(1):  # LWORD
            # Simple Field (value)
            size_in_bits += 64
        if data_type == "LWORD":  # List
            values: PlcList = _value
            size_in_bits += values.get_list().size() * 64
        if data_type == "SINT" and number_of_values == int(1):  # SINT
            # Simple Field (value)
            size_in_bits += 8
        if data_type == "SINT":  # List
            values: PlcList = _value
            size_in_bits += values.get_list().size() * 8
        if data_type == "INT" and number_of_values == int(1):  # INT
            # Simple Field (value)
            size_in_bits += 16
        if data_type == "INT":  # List
            values: PlcList = _value
            size_in_bits += values.get_list().size() * 16
        if data_type == "DINT" and number_of_values == int(1):  # DINT
            # Simple Field (value)
            size_in_bits += 32
        if data_type == "DINT":  # List
            values: PlcList = _value
            size_in_bits += values.get_list().size() * 32
        if data_type == "LINT" and number_of_values == int(1):  # LINT
            # Simple Field (value)
            size_in_bits += 64
        if data_type == "LINT":  # List
            values: PlcList = _value
            size_in_bits += values.get_list().size() * 64
        if data_type == "USINT" and number_of_values == int(1):  # USINT
            # Simple Field (value)
            size_in_bits += 8
        if data_type == "USINT":  # List
            values: PlcList = _value
            size_in_bits += values.get_list().size() * 8
        if data_type == "UINT" and number_of_values == int(1):  # UINT
            # Simple Field (value)
            size_in_bits += 16
        if data_type == "UINT":  # List
            values: PlcList = _value
            size_in_bits += values.get_list().size() * 16
        if data_type == "UDINT" and number_of_values == int(1):  # UDINT
            # Simple Field (value)
            size_in_bits += 32
        if data_type == "UDINT":  # List
            values: PlcList = _value
            size_in_bits += values.get_list().size() * 32
        if data_type == "ULINT" and number_of_values == int(1):  # ULINT
            # Simple Field (value)
            size_in_bits += 64
        if data_type == "ULINT":  # List
            values: PlcList = _value
            size_in_bits += values.get_list().size() * 64
        if data_type == "REAL" and number_of_values == int(1):  # REAL
            # Simple Field (value)
            size_in_bits += 32
        if data_type == "REAL":  # List
            values: PlcList = _value
            size_in_bits += values.get_list().size() * 32
        if data_type == "LREAL" and number_of_values == int(1):  # LREAL
            # Simple Field (value)
            size_in_bits += 64
        if data_type == "LREAL":  # List
            values: PlcList = _value
            size_in_bits += values.get_list().size() * 64
        if data_type == "CHAR" and number_of_values == int(1):  # CHAR
            # Simple Field (value)
            size_in_bits += 8
        if data_type == "CHAR":  # List
            values: PlcList = _value
            size_in_bits += values.get_list().size() * 8
        if data_type == "WCHAR" and number_of_values == int(1):  # WCHAR
            # Simple Field (value)
            size_in_bits += 16
        if data_type == "WCHAR":  # List
            values: PlcList = _value
            size_in_bits += values.get_list().size() * 16
        if data_type == "STRING":  # STRING
            # Simple Field (value)
            size_in_bits += 255
        if data_type == "WSTRING":  # STRING
            # Simple Field (value)
            size_in_bits += 255
        return size_in_bits
