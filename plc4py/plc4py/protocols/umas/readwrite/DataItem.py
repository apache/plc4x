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
from plc4py.protocols.umas import StaticHelper
from plc4py.protocols.umas.readwrite.UmasDataType import UmasDataType
from plc4py.spi.generation.ReadBuffer import ReadBuffer
from plc4py.spi.generation.WriteBuffer import WriteBuffer
from plc4py.spi.values.PlcValues import PlcBOOL
from plc4py.spi.values.PlcValues import PlcBYTE
from plc4py.spi.values.PlcValues import PlcDATE
from plc4py.spi.values.PlcValues import PlcDATE_AND_TIME
from plc4py.spi.values.PlcValues import PlcDINT
from plc4py.spi.values.PlcValues import PlcDWORD
from plc4py.spi.values.PlcValues import PlcINT
from plc4py.spi.values.PlcValues import PlcList
from plc4py.spi.values.PlcValues import PlcREAL
from plc4py.spi.values.PlcValues import PlcSTRING
from plc4py.spi.values.PlcValues import PlcTIME
from plc4py.spi.values.PlcValues import PlcTIME_OF_DAY
from plc4py.spi.values.PlcValues import PlcUDINT
from plc4py.spi.values.PlcValues import PlcUINT
from plc4py.spi.values.PlcValues import PlcWORD
from plc4py.utils.GenericTypes import ByteOrder
from typing import List
from typing import cast
import datetime
import logging
import math


class DataItem:
    @staticmethod
    def static_parse(
        read_buffer: ReadBuffer, data_type: UmasDataType, number_of_values: int
    ):
        if data_type == UmasDataType.BOOL and number_of_values == int(1):  # BOOL

            # Reserved Field (Compartmentalized so the "reserved" variable can't leak)
            reserved: int = read_buffer.read_unsigned_short(7, logical_name="")
            if reserved != int(0x0000):
                logging.warning(
                    "Expected constant value "
                    + str(0x0000)
                    + " but got "
                    + str(reserved)
                    + " for reserved field."
                )

            # Simple Field (value)
            value: bool = read_buffer.read_bit("")

            return PlcBOOL(value)
        if data_type == UmasDataType.EBOOL and number_of_values == int(1):  # BOOL

            # Reserved Field (Compartmentalized so the "reserved" variable can't leak)
            reserved: int = read_buffer.read_unsigned_short(7, logical_name="")
            if reserved != int(0x0000):
                logging.warning(
                    "Expected constant value "
                    + str(0x0000)
                    + " but got "
                    + str(reserved)
                    + " for reserved field."
                )

            # Simple Field (value)
            value: bool = read_buffer.read_bit("")

            return PlcBOOL(value)
        if data_type == UmasDataType.BYTE and number_of_values == int(1):  # BYTE

            # Simple Field (value)
            value: int = read_buffer.read_byte("")

            return PlcBYTE(value)
        if data_type == UmasDataType.BYTE:  # List
            # Array field (value)
            # Count array
            item_count: int = int(number_of_values)
            value: List[PlcValue] = []
            for _ in range(item_count):
                value.append(read_buffer.read_byte(""))

            return PlcList(value)
        if data_type == UmasDataType.WORD:  # WORD

            # Simple Field (value)
            value: int = read_buffer.read_unsigned_int(16, logical_name="")

            return PlcWORD(value)
        if data_type == UmasDataType.DWORD:  # DWORD

            # Simple Field (value)
            value: int = read_buffer.read_unsigned_long(32, logical_name="")

            return PlcDWORD(value)
        if data_type == UmasDataType.INT and number_of_values == int(1):  # INT

            # Simple Field (value)
            value: int = read_buffer.read_short(16, logical_name="")

            return PlcINT(value)
        if data_type == UmasDataType.INT:  # List
            # Array field (value)
            # Count array
            item_count: int = int(number_of_values)
            value: List[PlcValue] = []
            for _ in range(item_count):
                value.append(read_buffer.read_short(16, logical_name=""))

            return PlcList(value)
        if data_type == UmasDataType.DINT and number_of_values == int(1):  # DINT

            # Simple Field (value)
            value: int = read_buffer.read_int(32, logical_name="")

            return PlcDINT(value)
        if data_type == UmasDataType.DINT:  # List
            # Array field (value)
            # Count array
            item_count: int = int(number_of_values)
            value: List[PlcValue] = []
            for _ in range(item_count):
                value.append(read_buffer.read_int(32, logical_name=""))

            return PlcList(value)
        if data_type == UmasDataType.UINT and number_of_values == int(1):  # UINT

            # Simple Field (value)
            value: int = read_buffer.read_unsigned_int(16, logical_name="")

            return PlcUINT(value)
        if data_type == UmasDataType.UINT:  # List
            # Array field (value)
            # Count array
            item_count: int = int(number_of_values)
            value: List[PlcValue] = []
            for _ in range(item_count):
                value.append(read_buffer.read_unsigned_int(16, logical_name=""))

            return PlcList(value)
        if data_type == UmasDataType.UDINT and number_of_values == int(1):  # UDINT

            # Simple Field (value)
            value: int = read_buffer.read_unsigned_long(32, logical_name="")

            return PlcUDINT(value)
        if data_type == UmasDataType.UDINT:  # List
            # Array field (value)
            # Count array
            item_count: int = int(number_of_values)
            value: List[PlcValue] = []
            for _ in range(item_count):
                value.append(read_buffer.read_unsigned_long(32, logical_name=""))

            return PlcList(value)
        if data_type == UmasDataType.REAL and number_of_values == int(1):  # REAL

            # Simple Field (value)
            value: float = read_buffer.read_float(32, logical_name="")

            return PlcREAL(value)
        if data_type == UmasDataType.REAL:  # List
            # Array field (value)
            # Count array
            item_count: int = int(number_of_values)
            value: List[PlcValue] = []
            for _ in range(item_count):
                value.append(read_buffer.read_float(32, logical_name=""))

            return PlcList(value)
        if data_type == UmasDataType.STRING and number_of_values == int(1):  # STRING
            # Manual Field (value)
            value: str = (str)(
                StaticHelper.parse_terminated_string_bytes(
                    read_buffer, number_of_values
                )
            )

            return PlcSTRING(value)
        if data_type == UmasDataType.STRING:  # List
            # Array field (value)
            # Count array
            item_count: int = int(number_of_values)
            value: List[PlcValue] = []
            for _ in range(item_count):
                value.append(read_buffer.read_float(32, logical_name=""))

            return PlcList(value)
        if data_type == UmasDataType.TIME and number_of_values == int(1):  # TIME

            # Simple Field (value)
            value: int = read_buffer.read_unsigned_long(32, logical_name="")

            return PlcTIME(value)
        if data_type == UmasDataType.TIME:  # List
            # Array field (value)
            # Count array
            item_count: int = int(number_of_values)
            value: List[PlcValue] = []
            for _ in range(item_count):
                value.append(read_buffer.read_unsigned_long(32, logical_name=""))

            return PlcList(value)
        if data_type == UmasDataType.DATE and number_of_values == int(1):  # DATE

            # Simple Field (day)
            day: int = read_buffer.read_unsigned_short(
                8, logical_name="", encoding="BCD"
            )

            # Simple Field (month)
            month: int = read_buffer.read_unsigned_short(
                8, logical_name="", encoding="BCD"
            )

            # Simple Field (year)
            year: int = read_buffer.read_unsigned_int(
                16, logical_name="", encoding="BCD"
            )

            value: datetime = datetime.datetime(int(year), int(month), int(day))
            return PlcDATE(value)
        if data_type == UmasDataType.TOD and number_of_values == int(1):  # TIME_OF_DAY

            # Simple Field (value)
            value: int = read_buffer.read_unsigned_long(32, logical_name="")

            return PlcTIME_OF_DAY(value)
        if data_type == UmasDataType.TOD:  # List
            # Array field (value)
            # Count array
            item_count: int = int(number_of_values)
            value: List[PlcValue] = []
            for _ in range(item_count):
                value.append(read_buffer.read_unsigned_long(32, logical_name=""))

            return PlcList(value)
        if data_type == UmasDataType.DATE_AND_TIME and number_of_values == int(
            1
        ):  # DATE_AND_TIME

            # Simple Field (unused)
            unused: int = read_buffer.read_unsigned_short(8, logical_name="")

            # Simple Field (seconds)
            seconds: int = read_buffer.read_unsigned_short(
                8, logical_name="", encoding="BCD"
            )

            # Simple Field (minutes)
            minutes: int = read_buffer.read_unsigned_short(
                8, logical_name="", encoding="BCD"
            )

            # Simple Field (hour)
            hour: int = read_buffer.read_unsigned_short(
                8, logical_name="", encoding="BCD"
            )

            # Simple Field (day)
            day: int = read_buffer.read_unsigned_short(
                8, logical_name="", encoding="BCD"
            )

            # Simple Field (month)
            month: int = read_buffer.read_unsigned_short(
                8, logical_name="", encoding="BCD"
            )

            # Simple Field (year)
            year: int = read_buffer.read_unsigned_int(
                16, logical_name="", encoding="BCD"
            )

            value: datetime = datetime.datetime(
                int(year), int(month), int(day), int(hour), int(minutes), int(seconds)
            )
            return PlcDATE_AND_TIME(value)
        return None

    @staticmethod
    def static_serialize(
        write_buffer: WriteBuffer,
        _value: PlcValue,
        data_type: UmasDataType,
        number_of_values: int,
        byte_order: ByteOrder,
    ) -> None:
        if data_type == UmasDataType.BOOL and number_of_values == int(1):  # BOOL
            # Reserved Field
            write_buffer.write_byte(int(0x0000), 7, "int0x0000")
            # Simple Field (value)
            value: bool = _value.get_bool()
            write_buffer.write_bit((value), "value")

        elif data_type == UmasDataType.EBOOL and number_of_values == int(1):  # BOOL
            # Reserved Field
            write_buffer.write_byte(int(0x0000), 7, "int0x0000")
            # Simple Field (value)
            value: bool = _value.get_bool()
            write_buffer.write_bit((value), "value")

        elif data_type == UmasDataType.BYTE and number_of_values == int(1):  # BYTE
            # Simple Field (value)
            value: int = _value.get_int()
            write_buffer.write_byte((value), 8, "value")

        elif data_type == UmasDataType.BYTE:  # List
            values: PlcList = cast(PlcList, _value)
            for val in values.get_list():
                value: List[int] = val.get_raw()
                write_buffer.write_byte(value, 8, "")

        elif data_type == UmasDataType.WORD:  # WORD
            # Simple Field (value)
            value: int = _value.get_int()
            write_buffer.write_unsigned_short((value), 16, "value")

        elif data_type == UmasDataType.DWORD:  # DWORD
            # Simple Field (value)
            value: int = _value.get_int()
            write_buffer.write_unsigned_int((value), 32, "value")

        elif data_type == UmasDataType.INT and number_of_values == int(1):  # INT
            # Simple Field (value)
            value: int = _value.get_int()
            write_buffer.write_short((value), 16, "value")

        elif data_type == UmasDataType.INT:  # List
            values: PlcList = cast(PlcList, _value)
            for val in values.get_list():
                value: int = val.get_int()
                write_buffer.write_short((value), 16, "value")

        elif data_type == UmasDataType.DINT and number_of_values == int(1):  # DINT
            # Simple Field (value)
            value: int = _value.get_int()
            write_buffer.write_int((value), 32, "value")

        elif data_type == UmasDataType.DINT:  # List
            values: PlcList = cast(PlcList, _value)
            for val in values.get_list():
                value: int = val.get_int()
                write_buffer.write_int((value), 32, "value")

        elif data_type == UmasDataType.UINT and number_of_values == int(1):  # UINT
            # Simple Field (value)
            value: int = _value.get_int()
            write_buffer.write_unsigned_short((value), 16, "value")

        elif data_type == UmasDataType.UINT:  # List
            values: PlcList = cast(PlcList, _value)
            for val in values.get_list():
                value: int = val.get_int()
                write_buffer.write_unsigned_short((value), 16, "value")

        elif data_type == UmasDataType.UDINT and number_of_values == int(1):  # UDINT
            # Simple Field (value)
            value: int = _value.get_int()
            write_buffer.write_unsigned_int((value), 32, "value")

        elif data_type == UmasDataType.UDINT:  # List
            values: PlcList = cast(PlcList, _value)
            for val in values.get_list():
                value: int = val.get_int()
                write_buffer.write_unsigned_int((value), 32, "value")

        elif data_type == UmasDataType.REAL and number_of_values == int(1):  # REAL
            # Simple Field (value)
            value: float = _value.get_float()
            write_buffer.write_float((value), 32, "value")

        elif data_type == UmasDataType.REAL:  # List
            values: PlcList = cast(PlcList, _value)
            for val in values.get_list():
                value: float = val.get_float()
                write_buffer.write_float((value), 32, "value")

        elif data_type == UmasDataType.STRING and number_of_values == int(1):  # STRING
            # Manual Field (value)
            value: PlcValue = _value
            StaticHelper.serialize_terminated_string(
                write_buffer, value, number_of_values
            )
        elif data_type == UmasDataType.STRING:  # List
            values: PlcList = cast(PlcList, _value)
            for val in values.get_list():
                value: float = val.get_float()
                write_buffer.write_float((value), 32, "value")

        elif data_type == UmasDataType.TIME and number_of_values == int(1):  # TIME
            # Simple Field (value)
            value: int = _value.get_int()
            write_buffer.write_unsigned_int((value), 32, "value")

        elif data_type == UmasDataType.TIME:  # List
            values: PlcList = cast(PlcList, _value)
            for val in values.get_list():
                value: int = val.get_int()
                write_buffer.write_unsigned_int((value), 32, "value")

        elif data_type == UmasDataType.DATE and number_of_values == int(1):  # DATE
            # Simple Field (day)
            day: int = 0
            write_buffer.write_byte((day), 8, "day")

            # Simple Field (month)
            month: int = 0
            write_buffer.write_byte((month), 8, "month")

            # Simple Field (year)
            year: int = 0
            write_buffer.write_unsigned_short((year), 16, "year")

        elif data_type == UmasDataType.TOD and number_of_values == int(
            1
        ):  # TIME_OF_DAY
            # Simple Field (value)
            value: int = _value.get_int()
            write_buffer.write_unsigned_int((value), 32, "value")

        elif data_type == UmasDataType.TOD:  # List
            values: PlcList = cast(PlcList, _value)
            for val in values.get_list():
                value: int = val.get_int()
                write_buffer.write_unsigned_int((value), 32, "value")

        elif data_type == UmasDataType.DATE_AND_TIME and number_of_values == int(
            1
        ):  # DATE_AND_TIME
            # Simple Field (unused)
            unused: int = 0
            write_buffer.write_byte((unused), 8, "unused")

            # Simple Field (seconds)
            seconds: int = 0
            write_buffer.write_byte((seconds), 8, "seconds")

            # Simple Field (minutes)
            minutes: int = 0
            write_buffer.write_byte((minutes), 8, "minutes")

            # Simple Field (hour)
            hour: int = 0
            write_buffer.write_byte((hour), 8, "hour")

            # Simple Field (day)
            day: int = 0
            write_buffer.write_byte((day), 8, "day")

            # Simple Field (month)
            month: int = 0
            write_buffer.write_byte((month), 8, "month")

            # Simple Field (year)
            year: int = 0
            write_buffer.write_unsigned_short((year), 16, "year")

    @staticmethod
    def get_length_in_bytes(
        _value: PlcValue, data_type: UmasDataType, number_of_values: int
    ) -> int:
        return int(
            math.ceil(
                float(DataItem.get_length_in_bits(_value, data_type, number_of_values))
                / 8.0
            )
        )

    @staticmethod
    def get_length_in_bits(
        _value: PlcValue, data_type: UmasDataType, number_of_values: int
    ) -> int:
        size_in_bits: int = 0
        if data_type == UmasDataType.BOOL and number_of_values == int(1):  # BOOL
            # Reserved Field
            size_in_bits += 7
            # Simple Field (value)
            size_in_bits += 1
        elif data_type == UmasDataType.EBOOL and number_of_values == int(1):  # BOOL
            # Reserved Field
            size_in_bits += 7
            # Simple Field (value)
            size_in_bits += 1
        elif data_type == UmasDataType.BYTE and number_of_values == int(1):  # BYTE
            # Simple Field (value)
            size_in_bits += 8
        elif data_type == UmasDataType.BYTE:  # List
            values: PlcList = cast(PlcList, _value)
            size_in_bits += len(values.get_list()) * 8
        elif data_type == UmasDataType.WORD:  # WORD
            # Simple Field (value)
            size_in_bits += 16
        elif data_type == UmasDataType.DWORD:  # DWORD
            # Simple Field (value)
            size_in_bits += 32
        elif data_type == UmasDataType.INT and number_of_values == int(1):  # INT
            # Simple Field (value)
            size_in_bits += 16
        elif data_type == UmasDataType.INT:  # List
            values: PlcList = cast(PlcList, _value)
            size_in_bits += len(values.get_list()) * 16
        elif data_type == UmasDataType.DINT and number_of_values == int(1):  # DINT
            # Simple Field (value)
            size_in_bits += 32
        elif data_type == UmasDataType.DINT:  # List
            values: PlcList = cast(PlcList, _value)
            size_in_bits += len(values.get_list()) * 32
        elif data_type == UmasDataType.UINT and number_of_values == int(1):  # UINT
            # Simple Field (value)
            size_in_bits += 16
        elif data_type == UmasDataType.UINT:  # List
            values: PlcList = cast(PlcList, _value)
            size_in_bits += len(values.get_list()) * 16
        elif data_type == UmasDataType.UDINT and number_of_values == int(1):  # UDINT
            # Simple Field (value)
            size_in_bits += 32
        elif data_type == UmasDataType.UDINT:  # List
            values: PlcList = cast(PlcList, _value)
            size_in_bits += len(values.get_list()) * 32
        elif data_type == UmasDataType.REAL and number_of_values == int(1):  # REAL
            # Simple Field (value)
            size_in_bits += 32
        elif data_type == UmasDataType.REAL:  # List
            values: PlcList = cast(PlcList, _value)
            size_in_bits += len(values.get_list()) * 32
        elif data_type == UmasDataType.STRING and number_of_values == int(1):  # STRING
            # Manual Field (value)
            size_in_bits += number_of_values * int(8)
        elif data_type == UmasDataType.STRING:  # List
            values: PlcList = cast(PlcList, _value)
            size_in_bits += len(values.get_list()) * 32
        elif data_type == UmasDataType.TIME and number_of_values == int(1):  # TIME
            # Simple Field (value)
            size_in_bits += 32
        elif data_type == UmasDataType.TIME:  # List
            values: PlcList = cast(PlcList, _value)
            size_in_bits += len(values.get_list()) * 32
        elif data_type == UmasDataType.DATE and number_of_values == int(1):  # DATE
            # Simple Field (day)
            size_in_bits += 8
            # Simple Field (month)
            size_in_bits += 8
            # Simple Field (year)
            size_in_bits += 16
        elif data_type == UmasDataType.TOD and number_of_values == int(
            1
        ):  # TIME_OF_DAY
            # Simple Field (value)
            size_in_bits += 32
        elif data_type == UmasDataType.TOD:  # List
            values: PlcList = cast(PlcList, _value)
            size_in_bits += len(values.get_list()) * 32
        elif data_type == UmasDataType.DATE_AND_TIME and number_of_values == int(
            1
        ):  # DATE_AND_TIME
            # Simple Field (unused)
            size_in_bits += 8
            # Simple Field (seconds)
            size_in_bits += 8
            # Simple Field (minutes)
            size_in_bits += 8
            # Simple Field (hour)
            size_in_bits += 8
            # Simple Field (day)
            size_in_bits += 8
            # Simple Field (month)
            size_in_bits += 8
            # Simple Field (year)
            size_in_bits += 16

        return size_in_bits
