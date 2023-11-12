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
from plc4py.protocols.modbus.readwrite.ModbusDataType import ModbusDataType
from plc4py.spi.generation.ReadBuffer import ReadBuffer
from plc4py.spi.generation.WriteBuffer import WriteBuffer
from typing import List
import math


class DataItem:
    @staticmethod
    def static_parse(
        read_buffer: ReadBuffer, data_type: ModbusDataType, number_of_values: int
    ):
        if data_type == ModbusDataType.BOOL and number_of_values == int(1):  # BOOL
            # Reserved Field (Compartmentalized so the "reserved" variable can't leak)
            reserved: int = read_buffer.read_unsigned_int(15, logical_name="")
            if reserved != int(0x0000):
                log.info(
                    "Expected constant value "
                    + str(0x0000)
                    + " but got "
                    + reserved
                    + " for reserved field."
                )

            # Simple Field (value)
            value: bool = read_buffer.read_bit("")

            return PlcBOOL(value)
        if data_type == ModbusDataType.BOOL:  # List
            # Array field (value)
            # Count array
            item_count: int = int(number_of_values)
            value: List[PlcValue] = []
            for cur_item in range(item_count):
                value.append(PlcBOOL(bool(read_buffer.read_bit(""))))

            return PlcList(value)
        if data_type == ModbusDataType.BYTE and number_of_values == int(1):  # BYTE
            # Reserved Field (Compartmentalized so the "reserved" variable can't leak)
            reserved: int = read_buffer.read_unsigned_short(8, logical_name="")
            if reserved != int(0x00):
                log.info(
                    "Expected constant value "
                    + str(0x00)
                    + " but got "
                    + reserved
                    + " for reserved field."
                )

            # Simple Field (value)
            value: int = read_buffer.read_unsigned_short(8, logical_name="")

            return PlcBYTE(value)
        if data_type == ModbusDataType.BYTE:  # List
            # Array field (value)
            # Count array
            item_count: int = int(number_of_values * int(8))
            value: List[PlcValue] = []
            for cur_item in range(item_count):
                value.append(PlcBOOL(bool(read_buffer.read_bit(""))))

            return PlcList(value)
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
            value: int = read_buffer.read_unsigned_big_integer(64, logical_name="")

            return PlcLWORD(value)
        if data_type == ModbusDataType.SINT and number_of_values == int(1):  # SINT
            # Reserved Field (Compartmentalized so the "reserved" variable can't leak)
            reserved: int = read_buffer.read_unsigned_short(8, logical_name="")
            if reserved != int(0x00):
                log.info(
                    "Expected constant value "
                    + str(0x00)
                    + " but got "
                    + reserved
                    + " for reserved field."
                )

            # Simple Field (value)
            value: int = read_buffer.read_signed_byte(8, logical_name="")

            return PlcSINT(value)
        if data_type == ModbusDataType.SINT:  # List
            # Array field (value)
            # Count array
            item_count: int = int(number_of_values)
            value: List[PlcValue] = []
            for cur_item in range(item_count):
                value.append(
                    PlcSINT(int(read_buffer.read_signed_byte(8, logical_name="")))
                )

            return PlcList(value)
        if data_type == ModbusDataType.INT and number_of_values == int(1):  # INT
            # Simple Field (value)
            value: int = read_buffer.read_short(16, logical_name="")

            return PlcINT(value)
        if data_type == ModbusDataType.INT:  # List
            # Array field (value)
            # Count array
            item_count: int = int(number_of_values)
            value: List[PlcValue] = []
            for cur_item in range(item_count):
                value.append(PlcINT(int(read_buffer.read_short(16, logical_name=""))))

            return PlcList(value)
        if data_type == ModbusDataType.DINT and number_of_values == int(1):  # DINT
            # Simple Field (value)
            value: int = read_buffer.read_int(32, logical_name="")

            return PlcDINT(value)
        if data_type == ModbusDataType.DINT:  # List
            # Array field (value)
            # Count array
            item_count: int = int(number_of_values)
            value: List[PlcValue] = []
            for cur_item in range(item_count):
                value.append(PlcDINT(int(read_buffer.read_int(32, logical_name=""))))

            return PlcList(value)
        if data_type == ModbusDataType.LINT and number_of_values == int(1):  # LINT
            # Simple Field (value)
            value: int = read_buffer.read_long(64, logical_name="")

            return PlcLINT(value)
        if data_type == ModbusDataType.LINT:  # List
            # Array field (value)
            # Count array
            item_count: int = int(number_of_values)
            value: List[PlcValue] = []
            for cur_item in range(item_count):
                value.append(PlcLINT(int(read_buffer.read_long(64, logical_name=""))))

            return PlcList(value)
        if data_type == ModbusDataType.USINT and number_of_values == int(1):  # USINT
            # Reserved Field (Compartmentalized so the "reserved" variable can't leak)
            reserved: int = read_buffer.read_unsigned_short(8, logical_name="")
            if reserved != int(0x00):
                log.info(
                    "Expected constant value "
                    + str(0x00)
                    + " but got "
                    + reserved
                    + " for reserved field."
                )

            # Simple Field (value)
            value: int = read_buffer.read_unsigned_short(8, logical_name="")

            return PlcUSINT(value)
        if data_type == ModbusDataType.USINT:  # List
            # Array field (value)
            # Count array
            item_count: int = int(number_of_values)
            value: List[PlcValue] = []
            for cur_item in range(item_count):
                value.append(
                    PlcUINT(int(read_buffer.read_unsigned_short(8, logical_name="")))
                )

            return PlcList(value)
        if data_type == ModbusDataType.UINT and number_of_values == int(1):  # UINT
            # Simple Field (value)
            value: int = read_buffer.read_unsigned_int(16, logical_name="")

            return PlcUINT(value)
        if data_type == ModbusDataType.UINT:  # List
            # Array field (value)
            # Count array
            item_count: int = int(number_of_values)
            value: List[PlcValue] = []
            for cur_item in range(item_count):
                value.append(
                    PlcUDINT(int(read_buffer.read_unsigned_int(16, logical_name="")))
                )

            return PlcList(value)
        if data_type == ModbusDataType.UDINT and number_of_values == int(1):  # UDINT
            # Simple Field (value)
            value: int = read_buffer.read_unsigned_long(32, logical_name="")

            return PlcUDINT(value)
        if data_type == ModbusDataType.UDINT:  # List
            # Array field (value)
            # Count array
            item_count: int = int(number_of_values)
            value: List[PlcValue] = []
            for cur_item in range(item_count):
                value.append(
                    PlcULINT(int(read_buffer.read_unsigned_long(32, logical_name="")))
                )

            return PlcList(value)
        if data_type == ModbusDataType.ULINT and number_of_values == int(1):  # ULINT
            # Simple Field (value)
            value: int = read_buffer.read_unsigned_big_integer(64, logical_name="")

            return PlcULINT(value)
        if data_type == ModbusDataType.ULINT:  # List
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
        if data_type == ModbusDataType.REAL and number_of_values == int(1):  # REAL
            # Simple Field (value)
            value: float = read_buffer.read_float(32, logical_name="")

            return PlcREAL(value)
        if data_type == ModbusDataType.REAL:  # List
            # Array field (value)
            # Count array
            item_count: int = int(number_of_values)
            value: List[PlcValue] = []
            for cur_item in range(item_count):
                value.append(
                    PlcREAL(float(read_buffer.read_float(32, logical_name="")))
                )

            return PlcList(value)
        if data_type == ModbusDataType.LREAL and number_of_values == int(1):  # LREAL
            # Simple Field (value)
            value: float = read_buffer.read_double(64, logical_name="")

            return PlcLREAL(value)
        if data_type == ModbusDataType.LREAL:  # List
            # Array field (value)
            # Count array
            item_count: int = int(number_of_values)
            value: List[PlcValue] = []
            for cur_item in range(item_count):
                value.append(
                    PlcLREAL(float(read_buffer.read_double(64, logical_name="")))
                )

            return PlcList(value)
        if data_type == ModbusDataType.CHAR and number_of_values == int(1):  # CHAR
            # Simple Field (value)
            value: str = read_buffer.read_string(8, logical_name="", encoding="")

            return PlcCHAR(value)
        if data_type == ModbusDataType.CHAR:  # List
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
        if data_type == ModbusDataType.WCHAR and number_of_values == int(1):  # WCHAR
            # Simple Field (value)
            value: str = read_buffer.read_string(16, logical_name="", encoding="")

            return PlcWCHAR(value)
        if data_type == ModbusDataType.WCHAR:  # List
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
        return None

    @staticmethod
    def static_serialize(
        write_buffer: WriteBuffer,
        _value: PlcValue,
        data_type: ModbusDataType,
        number_of_values: int,
        byte_order: ByteOrder,
    ) -> None:
        if data_type == ModbusDataType.BOOL and number_of_values == int(1):  # BOOL
            # Reserved Field
            write_buffer.write_unsigned_short(int(0x0000), 15, "int0x0000")
            # Simple Field (value)
            value: bool = _value.get_bool()
            write_buffer.write_bit((value), "value")
        if data_type == ModbusDataType.BOOL:  # List
            values: PlcList = _value

            for val in values.getList():
                value: bool = val.get_bool()
                write_buffer.write_bit((value), "value")

        if data_type == ModbusDataType.BYTE and number_of_values == int(1):  # BYTE
            # Reserved Field
            write_buffer.write_byte(int(0x00), 8, "int0x00")
            # Simple Field (value)
            value: int = _value.get_int()
            write_buffer.write_byte((value), 8, "value")
        if data_type == ModbusDataType.BYTE:  # List
            values: PlcList = _value

            for val in values.getList():
                value: bool = val.get_bool()
                write_buffer.write_bit((value), "value")

        if data_type == ModbusDataType.WORD:  # WORD
            # Simple Field (value)
            value: int = _value.get_int()
            write_buffer.write_unsigned_short((value), 16, "value")
        if data_type == ModbusDataType.DWORD:  # DWORD
            # Simple Field (value)
            value: int = _value.get_int()
            write_buffer.write_unsigned_int((value), 32, "value")
        if data_type == ModbusDataType.LWORD:  # LWORD
            # Simple Field (value)
            value: int = _value.get_int()
            write_buffer.write_unsigned_long((value), 64, "value")
        if data_type == ModbusDataType.SINT and number_of_values == int(1):  # SINT
            # Reserved Field
            write_buffer.write_byte(int(0x00), 8, "int0x00")
            # Simple Field (value)
            value: int = _value.get_int()
            write_buffer.write_signed_byte((value), 8, "value")
        if data_type == ModbusDataType.SINT:  # List
            values: PlcList = _value

            for val in values.getList():
                value: int = val.get_int()
                write_buffer.write_signed_byte((value), 8, "value")

        if data_type == ModbusDataType.INT and number_of_values == int(1):  # INT
            # Simple Field (value)
            value: int = _value.get_int()
            write_buffer.write_short((value), 16, "value")
        if data_type == ModbusDataType.INT:  # List
            values: PlcList = _value

            for val in values.getList():
                value: int = val.get_int()
                write_buffer.write_short((value), 16, "value")

        if data_type == ModbusDataType.DINT and number_of_values == int(1):  # DINT
            # Simple Field (value)
            value: int = _value.get_int()
            write_buffer.write_int((value), 32, "value")
        if data_type == ModbusDataType.DINT:  # List
            values: PlcList = _value

            for val in values.getList():
                value: int = val.get_int()
                write_buffer.write_int((value), 32, "value")

        if data_type == ModbusDataType.LINT and number_of_values == int(1):  # LINT
            # Simple Field (value)
            value: int = _value.get_int()
            write_buffer.write_long((value), 64, "value")
        if data_type == ModbusDataType.LINT:  # List
            values: PlcList = _value

            for val in values.getList():
                value: int = val.get_int()
                write_buffer.write_long((value), 64, "value")

        if data_type == ModbusDataType.USINT and number_of_values == int(1):  # USINT
            # Reserved Field
            write_buffer.write_byte(int(0x00), 8, "int0x00")
            # Simple Field (value)
            value: int = _value.get_int()
            write_buffer.write_byte((value), 8, "value")
        if data_type == ModbusDataType.USINT:  # List
            values: PlcList = _value

            for val in values.getList():
                value: int = val.get_int()
                write_buffer.write_byte((value), 8, "value")

        if data_type == ModbusDataType.UINT and number_of_values == int(1):  # UINT
            # Simple Field (value)
            value: int = _value.get_int()
            write_buffer.write_unsigned_short((value), 16, "value")
        if data_type == ModbusDataType.UINT:  # List
            values: PlcList = _value

            for val in values.getList():
                value: int = val.get_int()
                write_buffer.write_unsigned_short((value), 16, "value")

        if data_type == ModbusDataType.UDINT and number_of_values == int(1):  # UDINT
            # Simple Field (value)
            value: int = _value.get_int()
            write_buffer.write_unsigned_int((value), 32, "value")
        if data_type == ModbusDataType.UDINT:  # List
            values: PlcList = _value

            for val in values.getList():
                value: int = val.get_int()
                write_buffer.write_unsigned_int((value), 32, "value")

        if data_type == ModbusDataType.ULINT and number_of_values == int(1):  # ULINT
            # Simple Field (value)
            value: int = _value.get_int()
            write_buffer.write_unsigned_long((value), 64, "value")
        if data_type == ModbusDataType.ULINT:  # List
            values: PlcList = _value

            for val in values.getList():
                value: int = val.get_int()
                write_buffer.write_unsigned_long((value), 64, "value")

        if data_type == ModbusDataType.REAL and number_of_values == int(1):  # REAL
            # Simple Field (value)
            value: float = _value.get_float()
            write_buffer.write_float((value), 32, "value")
        if data_type == ModbusDataType.REAL:  # List
            values: PlcList = _value

            for val in values.getList():
                value: float = val.get_float()
                write_buffer.write_float((value), 32, "value")

        if data_type == ModbusDataType.LREAL and number_of_values == int(1):  # LREAL
            # Simple Field (value)
            value: float = _value.get_float()
            write_buffer.write_double((value), 64, "value")
        if data_type == ModbusDataType.LREAL:  # List
            values: PlcList = _value

            for val in values.getList():
                value: float = val.get_float()
                write_buffer.write_double((value), 64, "value")

        if data_type == ModbusDataType.CHAR and number_of_values == int(1):  # CHAR
            # Simple Field (value)
            value: str = _value.get_str()
            write_buffer.write_str((value), uint32(8), "UTF-8", "value")
        if data_type == ModbusDataType.CHAR:  # List
            values: PlcList = _value

            for val in values.getList():
                value: str = val.get_str()
                write_buffer.write_str((value), uint32(8), "UTF-8", "value")

        if data_type == ModbusDataType.WCHAR and number_of_values == int(1):  # WCHAR
            # Simple Field (value)
            value: str = _value.get_str()
            write_buffer.write_str((value), uint32(16), "UTF-16", "value")
        if data_type == ModbusDataType.WCHAR:  # List
            values: PlcList = _value

            for val in values.getList():
                value: str = val.get_str()
                write_buffer.write_str((value), uint32(16), "UTF-16", "value")

    @staticmethod
    def get_length_in_bytes(
        _value: PlcValue, data_type: ModbusDataType, number_of_values: int
    ) -> int:
        return int(
            math.ceil(
                float(get_length_in_bits(_value, data_type, number_of_values)) / 8.0
            )
        )

    @staticmethod
    def get_length_in_bits(
        _value: PlcValue, data_type: ModbusDataType, number_of_values: int
    ) -> int:
        size_in_bits: int = 0
        if data_type == ModbusDataType.BOOL and number_of_values == int(1):  # BOOL
            # Reserved Field
            size_in_bits += 15
            # Simple Field (value)
            size_in_bits += 1
        if data_type == ModbusDataType.BOOL:  # List
            values: PlcList = _value
            size_in_bits += values.get_list().size() * 1
        if data_type == ModbusDataType.BYTE and number_of_values == int(1):  # BYTE
            # Reserved Field
            size_in_bits += 8
            # Simple Field (value)
            size_in_bits += 8
        if data_type == ModbusDataType.BYTE:  # List
            values: PlcList = _value
            size_in_bits += values.get_list().size() * 1
        if data_type == ModbusDataType.WORD:  # WORD
            # Simple Field (value)
            size_in_bits += 16
        if data_type == ModbusDataType.DWORD:  # DWORD
            # Simple Field (value)
            size_in_bits += 32
        if data_type == ModbusDataType.LWORD:  # LWORD
            # Simple Field (value)
            size_in_bits += 64
        if data_type == ModbusDataType.SINT and number_of_values == int(1):  # SINT
            # Reserved Field
            size_in_bits += 8
            # Simple Field (value)
            size_in_bits += 8
        if data_type == ModbusDataType.SINT:  # List
            values: PlcList = _value
            size_in_bits += values.get_list().size() * 8
        if data_type == ModbusDataType.INT and number_of_values == int(1):  # INT
            # Simple Field (value)
            size_in_bits += 16
        if data_type == ModbusDataType.INT:  # List
            values: PlcList = _value
            size_in_bits += values.get_list().size() * 16
        if data_type == ModbusDataType.DINT and number_of_values == int(1):  # DINT
            # Simple Field (value)
            size_in_bits += 32
        if data_type == ModbusDataType.DINT:  # List
            values: PlcList = _value
            size_in_bits += values.get_list().size() * 32
        if data_type == ModbusDataType.LINT and number_of_values == int(1):  # LINT
            # Simple Field (value)
            size_in_bits += 64
        if data_type == ModbusDataType.LINT:  # List
            values: PlcList = _value
            size_in_bits += values.get_list().size() * 64
        if data_type == ModbusDataType.USINT and number_of_values == int(1):  # USINT
            # Reserved Field
            size_in_bits += 8
            # Simple Field (value)
            size_in_bits += 8
        if data_type == ModbusDataType.USINT:  # List
            values: PlcList = _value
            size_in_bits += values.get_list().size() * 8
        if data_type == ModbusDataType.UINT and number_of_values == int(1):  # UINT
            # Simple Field (value)
            size_in_bits += 16
        if data_type == ModbusDataType.UINT:  # List
            values: PlcList = _value
            size_in_bits += values.get_list().size() * 16
        if data_type == ModbusDataType.UDINT and number_of_values == int(1):  # UDINT
            # Simple Field (value)
            size_in_bits += 32
        if data_type == ModbusDataType.UDINT:  # List
            values: PlcList = _value
            size_in_bits += values.get_list().size() * 32
        if data_type == ModbusDataType.ULINT and number_of_values == int(1):  # ULINT
            # Simple Field (value)
            size_in_bits += 64
        if data_type == ModbusDataType.ULINT:  # List
            values: PlcList = _value
            size_in_bits += values.get_list().size() * 64
        if data_type == ModbusDataType.REAL and number_of_values == int(1):  # REAL
            # Simple Field (value)
            size_in_bits += 32
        if data_type == ModbusDataType.REAL:  # List
            values: PlcList = _value
            size_in_bits += values.get_list().size() * 32
        if data_type == ModbusDataType.LREAL and number_of_values == int(1):  # LREAL
            # Simple Field (value)
            size_in_bits += 64
        if data_type == ModbusDataType.LREAL:  # List
            values: PlcList = _value
            size_in_bits += values.get_list().size() * 64
        if data_type == ModbusDataType.CHAR and number_of_values == int(1):  # CHAR
            # Simple Field (value)
            size_in_bits += 8
        if data_type == ModbusDataType.CHAR:  # List
            values: PlcList = _value
            size_in_bits += values.get_list().size() * 8
        if data_type == ModbusDataType.WCHAR and number_of_values == int(1):  # WCHAR
            # Simple Field (value)
            size_in_bits += 16
        if data_type == ModbusDataType.WCHAR:  # List
            values: PlcList = _value
            size_in_bits += values.get_list().size() * 16
        return size_in_bits
