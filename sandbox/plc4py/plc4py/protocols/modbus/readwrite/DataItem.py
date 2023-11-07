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
from plc4py.protocols.modbus.readwrite.ModbusDataType import ModbusDataType
import math


class DataItem:
    @staticmethod
    def static_parse(
        read_buffer: ReadBuffer, data_type: ModbusDataType, number_of_values: int
    ):
        if EvaluationHelper.equals(
            data_type, ModbusDataType.BOOL
        ) and EvaluationHelper.equals(
            number_of_values, int(1)
        ):  # BOOL
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
        if EvaluationHelper.equals(data_type, ModbusDataType.BOOL):  # List
            # Array field (value)
            # Count array
            if number_of_values > Integer.MAX_VALUE:
                raise ParseException(
                    "Array count of "
                    + (number_of_values)
                    + " exceeds the maximum allowed count of "
                    + Integer.MAX_VALUE
                )

            item_count: int = int(number_of_values)
            value: List[PlcValue] = []
            for cur_item in range(item_count):
                value.append(PlcBOOL(bool(read_buffer.read_bit(""))))

            return PlcList(value)
        if EvaluationHelper.equals(
            data_type, ModbusDataType.BYTE
        ) and EvaluationHelper.equals(
            number_of_values, int(1)
        ):  # BYTE
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
        if EvaluationHelper.equals(data_type, ModbusDataType.BYTE):  # List
            # Array field (value)
            # Count array
            if number_of_values * int(8) > Integer.MAX_VALUE:
                raise ParseException(
                    "Array count of "
                    + (number_of_values * int(8))
                    + " exceeds the maximum allowed count of "
                    + Integer.MAX_VALUE
                )

            item_count: int = int(number_of_values * int(8))
            value: List[PlcValue] = []
            for cur_item in range(item_count):
                value.append(PlcBOOL(bool(read_buffer.read_bit(""))))

            return PlcList(value)
        if EvaluationHelper.equals(data_type, ModbusDataType.WORD):  # WORD
            # Simple Field (value)
            value: int = read_buffer.read_unsigned_int(16, logical_name="")

            return PlcWORD(value)
        if EvaluationHelper.equals(data_type, ModbusDataType.DWORD):  # DWORD
            # Simple Field (value)
            value: int = read_buffer.read_unsigned_long(32, logical_name="")

            return PlcDWORD(value)
        if EvaluationHelper.equals(data_type, ModbusDataType.LWORD):  # LWORD
            # Simple Field (value)
            value: int = read_buffer.read_unsigned_big_integer(64, logical_name="")

            return PlcLWORD(value)
        if EvaluationHelper.equals(
            data_type, ModbusDataType.SINT
        ) and EvaluationHelper.equals(
            number_of_values, int(1)
        ):  # SINT
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
        if EvaluationHelper.equals(data_type, ModbusDataType.SINT):  # List
            # Array field (value)
            # Count array
            if number_of_values > Integer.MAX_VALUE:
                raise ParseException(
                    "Array count of "
                    + (number_of_values)
                    + " exceeds the maximum allowed count of "
                    + Integer.MAX_VALUE
                )

            item_count: int = int(number_of_values)
            value: List[PlcValue] = []
            for cur_item in range(item_count):
                value.append(
                    PlcSINT(int(read_buffer.read_signed_byte(8, logical_name="")))
                )

            return PlcList(value)
        if EvaluationHelper.equals(
            data_type, ModbusDataType.INT
        ) and EvaluationHelper.equals(
            number_of_values, int(1)
        ):  # INT
            # Simple Field (value)
            value: int = read_buffer.read_short(16, logical_name="")

            return PlcINT(value)
        if EvaluationHelper.equals(data_type, ModbusDataType.INT):  # List
            # Array field (value)
            # Count array
            if number_of_values > Integer.MAX_VALUE:
                raise ParseException(
                    "Array count of "
                    + (number_of_values)
                    + " exceeds the maximum allowed count of "
                    + Integer.MAX_VALUE
                )

            item_count: int = int(number_of_values)
            value: List[PlcValue] = []
            for cur_item in range(item_count):
                value.append(PlcINT(int(read_buffer.read_short(16, logical_name=""))))

            return PlcList(value)
        if EvaluationHelper.equals(
            data_type, ModbusDataType.DINT
        ) and EvaluationHelper.equals(
            number_of_values, int(1)
        ):  # DINT
            # Simple Field (value)
            value: int = read_buffer.read_int(32, logical_name="")

            return PlcDINT(value)
        if EvaluationHelper.equals(data_type, ModbusDataType.DINT):  # List
            # Array field (value)
            # Count array
            if number_of_values > Integer.MAX_VALUE:
                raise ParseException(
                    "Array count of "
                    + (number_of_values)
                    + " exceeds the maximum allowed count of "
                    + Integer.MAX_VALUE
                )

            item_count: int = int(number_of_values)
            value: List[PlcValue] = []
            for cur_item in range(item_count):
                value.append(PlcDINT(int(read_buffer.read_int(32, logical_name=""))))

            return PlcList(value)
        if EvaluationHelper.equals(
            data_type, ModbusDataType.LINT
        ) and EvaluationHelper.equals(
            number_of_values, int(1)
        ):  # LINT
            # Simple Field (value)
            value: int = read_buffer.read_long(64, logical_name="")

            return PlcLINT(value)
        if EvaluationHelper.equals(data_type, ModbusDataType.LINT):  # List
            # Array field (value)
            # Count array
            if number_of_values > Integer.MAX_VALUE:
                raise ParseException(
                    "Array count of "
                    + (number_of_values)
                    + " exceeds the maximum allowed count of "
                    + Integer.MAX_VALUE
                )

            item_count: int = int(number_of_values)
            value: List[PlcValue] = []
            for cur_item in range(item_count):
                value.append(PlcLINT(int(read_buffer.read_long(64, logical_name=""))))

            return PlcList(value)
        if EvaluationHelper.equals(
            data_type, ModbusDataType.USINT
        ) and EvaluationHelper.equals(
            number_of_values, int(1)
        ):  # USINT
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
        if EvaluationHelper.equals(data_type, ModbusDataType.USINT):  # List
            # Array field (value)
            # Count array
            if number_of_values > Integer.MAX_VALUE:
                raise ParseException(
                    "Array count of "
                    + (number_of_values)
                    + " exceeds the maximum allowed count of "
                    + Integer.MAX_VALUE
                )

            item_count: int = int(number_of_values)
            value: List[PlcValue] = []
            for cur_item in range(item_count):
                value.append(
                    PlcUINT(int(read_buffer.read_unsigned_short(8, logical_name="")))
                )

            return PlcList(value)
        if EvaluationHelper.equals(
            data_type, ModbusDataType.UINT
        ) and EvaluationHelper.equals(
            number_of_values, int(1)
        ):  # UINT
            # Simple Field (value)
            value: int = read_buffer.read_unsigned_int(16, logical_name="")

            return PlcUINT(value)
        if EvaluationHelper.equals(data_type, ModbusDataType.UINT):  # List
            # Array field (value)
            # Count array
            if number_of_values > Integer.MAX_VALUE:
                raise ParseException(
                    "Array count of "
                    + (number_of_values)
                    + " exceeds the maximum allowed count of "
                    + Integer.MAX_VALUE
                )

            item_count: int = int(number_of_values)
            value: List[PlcValue] = []
            for cur_item in range(item_count):
                value.append(
                    PlcUDINT(int(read_buffer.read_unsigned_int(16, logical_name="")))
                )

            return PlcList(value)
        if EvaluationHelper.equals(
            data_type, ModbusDataType.UDINT
        ) and EvaluationHelper.equals(
            number_of_values, int(1)
        ):  # UDINT
            # Simple Field (value)
            value: int = read_buffer.read_unsigned_long(32, logical_name="")

            return PlcUDINT(value)
        if EvaluationHelper.equals(data_type, ModbusDataType.UDINT):  # List
            # Array field (value)
            # Count array
            if number_of_values > Integer.MAX_VALUE:
                raise ParseException(
                    "Array count of "
                    + (number_of_values)
                    + " exceeds the maximum allowed count of "
                    + Integer.MAX_VALUE
                )

            item_count: int = int(number_of_values)
            value: List[PlcValue] = []
            for cur_item in range(item_count):
                value.append(
                    PlcULINT(int(read_buffer.read_unsigned_long(32, logical_name="")))
                )

            return PlcList(value)
        if EvaluationHelper.equals(
            data_type, ModbusDataType.ULINT
        ) and EvaluationHelper.equals(
            number_of_values, int(1)
        ):  # ULINT
            # Simple Field (value)
            value: int = read_buffer.read_unsigned_big_integer(64, logical_name="")

            return PlcULINT(value)
        if EvaluationHelper.equals(data_type, ModbusDataType.ULINT):  # List
            # Array field (value)
            # Count array
            if number_of_values > Integer.MAX_VALUE:
                raise ParseException(
                    "Array count of "
                    + (number_of_values)
                    + " exceeds the maximum allowed count of "
                    + Integer.MAX_VALUE
                )

            item_count: int = int(number_of_values)
            value: List[PlcValue] = []
            for cur_item in range(item_count):
                value.append(
                    PlcLINT(
                        int(read_buffer.read_unsigned_big_integer(64, logical_name=""))
                    )
                )

            return PlcList(value)
        if EvaluationHelper.equals(
            data_type, ModbusDataType.REAL
        ) and EvaluationHelper.equals(
            number_of_values, int(1)
        ):  # REAL
            # Simple Field (value)
            value: float = read_buffer.read_float(32, logical_name="")

            return PlcREAL(value)
        if EvaluationHelper.equals(data_type, ModbusDataType.REAL):  # List
            # Array field (value)
            # Count array
            if number_of_values > Integer.MAX_VALUE:
                raise ParseException(
                    "Array count of "
                    + (number_of_values)
                    + " exceeds the maximum allowed count of "
                    + Integer.MAX_VALUE
                )

            item_count: int = int(number_of_values)
            value: List[PlcValue] = []
            for cur_item in range(item_count):
                value.append(
                    PlcREAL(float(read_buffer.read_float(32, logical_name="")))
                )

            return PlcList(value)
        if EvaluationHelper.equals(
            data_type, ModbusDataType.LREAL
        ) and EvaluationHelper.equals(
            number_of_values, int(1)
        ):  # LREAL
            # Simple Field (value)
            value: float = read_buffer.read_double(64, logical_name="")

            return PlcLREAL(value)
        if EvaluationHelper.equals(data_type, ModbusDataType.LREAL):  # List
            # Array field (value)
            # Count array
            if number_of_values > Integer.MAX_VALUE:
                raise ParseException(
                    "Array count of "
                    + (number_of_values)
                    + " exceeds the maximum allowed count of "
                    + Integer.MAX_VALUE
                )

            item_count: int = int(number_of_values)
            value: List[PlcValue] = []
            for cur_item in range(item_count):
                value.append(
                    PlcLREAL(float(read_buffer.read_double(64, logical_name="")))
                )

            return PlcList(value)
        if EvaluationHelper.equals(
            data_type, ModbusDataType.CHAR
        ) and EvaluationHelper.equals(
            number_of_values, int(1)
        ):  # CHAR
            # Simple Field (value)
            value: str = read_buffer.read_string(8, logical_name="", encoding="")

            return PlcCHAR(value)
        if EvaluationHelper.equals(data_type, ModbusDataType.CHAR):  # List
            # Array field (value)
            # Count array
            if number_of_values > Integer.MAX_VALUE:
                raise ParseException(
                    "Array count of "
                    + (number_of_values)
                    + " exceeds the maximum allowed count of "
                    + Integer.MAX_VALUE
                )

            item_count: int = int(number_of_values)
            value: List[PlcValue] = []
            for cur_item in range(item_count):
                value.append(
                    PlcSTRING(
                        str(read_buffer.read_string(8, logical_name="", encoding=""))
                    )
                )

            return PlcList(value)
        if EvaluationHelper.equals(
            data_type, ModbusDataType.WCHAR
        ) and EvaluationHelper.equals(
            number_of_values, int(1)
        ):  # WCHAR
            # Simple Field (value)
            value: str = read_buffer.read_string(16, logical_name="", encoding="")

            return PlcWCHAR(value)
        if EvaluationHelper.equals(data_type, ModbusDataType.WCHAR):  # List
            # Array field (value)
            # Count array
            if number_of_values > Integer.MAX_VALUE:
                raise ParseException(
                    "Array count of "
                    + (number_of_values)
                    + " exceeds the maximum allowed count of "
                    + Integer.MAX_VALUE
                )

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
        writeBuffer: WriteBuffer,
        _value: PlcValue,
        dataType: ModbusDataType,
        numberOfValues: int,
    ) -> None:
        static_serialize(
            writeBuffer, _value, dataType, numberOfValues, ByteOrder.BIG_ENDIAN
        )

    @staticmethod
    def static_serialize(
        writeBuffer: WriteBuffer,
        _value: PlcValue,
        dataType: ModbusDataType,
        numberOfValues: int,
        byteOrder: ByteOrder,
    ) -> None:
        if EvaluationHelper.equals(
            data_type, ModbusDataType.BOOL
        ) and EvaluationHelper.equals(
            number_of_values, int(1)
        ):  # BOOL
            # Reserved Field
            writeBuffer.WriteUint16("int0x0000", 15, int(0x0000))
            # Simple Field (value)
            value: bool = _value.getBool()
            writeBuffer.WriteBit("value", (value))
        if EvaluationHelper.equals(data_type, ModbusDataType.BOOL):  # List
            values: PlcList = _value

            for val in values.getList():
                value: bool = val.getBool()
                writeBuffer.WriteBit("value", (value))

        if EvaluationHelper.equals(
            data_type, ModbusDataType.BYTE
        ) and EvaluationHelper.equals(
            number_of_values, int(1)
        ):  # BYTE
            # Reserved Field
            writeBuffer.WriteUint8("int0x00", 8, int(0x00))
            # Simple Field (value)
            value: int = _value.getInt()
            writeBuffer.WriteUint8("value", 8, (value))
        if EvaluationHelper.equals(data_type, ModbusDataType.BYTE):  # List
            values: PlcList = _value

            for val in values.getList():
                value: bool = val.getBool()
                writeBuffer.WriteBit("value", (value))

        if EvaluationHelper.equals(data_type, ModbusDataType.WORD):  # WORD
            # Simple Field (value)
            value: int = _value.getInt()
            writeBuffer.WriteUint16("value", 16, (value))
        if EvaluationHelper.equals(data_type, ModbusDataType.DWORD):  # DWORD
            # Simple Field (value)
            value: int = _value.getInt()
            writeBuffer.WriteUint32("value", 32, (value))
        if EvaluationHelper.equals(data_type, ModbusDataType.LWORD):  # LWORD
            # Simple Field (value)
            value: int = _value.getInt()
            writeBuffer.WriteUint64("value", 64, (value))
        if EvaluationHelper.equals(
            data_type, ModbusDataType.SINT
        ) and EvaluationHelper.equals(
            number_of_values, int(1)
        ):  # SINT
            # Reserved Field
            writeBuffer.WriteUint8("int0x00", 8, int(0x00))
            # Simple Field (value)
            value: int = _value.getInt()
            writeBuffer.WriteInt8("value", 8, (value))
        if EvaluationHelper.equals(data_type, ModbusDataType.SINT):  # List
            values: PlcList = _value

            for val in values.getList():
                value: int = val.getInt()
                writeBuffer.WriteInt8("value", 8, (value))

        if EvaluationHelper.equals(
            data_type, ModbusDataType.INT
        ) and EvaluationHelper.equals(
            number_of_values, int(1)
        ):  # INT
            # Simple Field (value)
            value: int = _value.getInt()
            writeBuffer.WriteInt16("value", 16, (value))
        if EvaluationHelper.equals(data_type, ModbusDataType.INT):  # List
            values: PlcList = _value

            for val in values.getList():
                value: int = val.getInt()
                writeBuffer.WriteInt16("value", 16, (value))

        if EvaluationHelper.equals(
            data_type, ModbusDataType.DINT
        ) and EvaluationHelper.equals(
            number_of_values, int(1)
        ):  # DINT
            # Simple Field (value)
            value: int = _value.getInt()
            writeBuffer.WriteInt32("value", 32, (value))
        if EvaluationHelper.equals(data_type, ModbusDataType.DINT):  # List
            values: PlcList = _value

            for val in values.getList():
                value: int = val.getInt()
                writeBuffer.WriteInt32("value", 32, (value))

        if EvaluationHelper.equals(
            data_type, ModbusDataType.LINT
        ) and EvaluationHelper.equals(
            number_of_values, int(1)
        ):  # LINT
            # Simple Field (value)
            value: int = _value.getInt()
            writeBuffer.WriteInt64("value", 64, (value))
        if EvaluationHelper.equals(data_type, ModbusDataType.LINT):  # List
            values: PlcList = _value

            for val in values.getList():
                value: int = val.getInt()
                writeBuffer.WriteInt64("value", 64, (value))

        if EvaluationHelper.equals(
            data_type, ModbusDataType.USINT
        ) and EvaluationHelper.equals(
            number_of_values, int(1)
        ):  # USINT
            # Reserved Field
            writeBuffer.WriteUint8("int0x00", 8, int(0x00))
            # Simple Field (value)
            value: int = _value.getInt()
            writeBuffer.WriteUint8("value", 8, (value))
        if EvaluationHelper.equals(data_type, ModbusDataType.USINT):  # List
            values: PlcList = _value

            for val in values.getList():
                value: int = val.getInt()
                writeBuffer.WriteUint8("value", 8, (value))

        if EvaluationHelper.equals(
            data_type, ModbusDataType.UINT
        ) and EvaluationHelper.equals(
            number_of_values, int(1)
        ):  # UINT
            # Simple Field (value)
            value: int = _value.getInt()
            writeBuffer.WriteUint16("value", 16, (value))
        if EvaluationHelper.equals(data_type, ModbusDataType.UINT):  # List
            values: PlcList = _value

            for val in values.getList():
                value: int = val.getInt()
                writeBuffer.WriteUint16("value", 16, (value))

        if EvaluationHelper.equals(
            data_type, ModbusDataType.UDINT
        ) and EvaluationHelper.equals(
            number_of_values, int(1)
        ):  # UDINT
            # Simple Field (value)
            value: int = _value.getInt()
            writeBuffer.WriteUint32("value", 32, (value))
        if EvaluationHelper.equals(data_type, ModbusDataType.UDINT):  # List
            values: PlcList = _value

            for val in values.getList():
                value: int = val.getInt()
                writeBuffer.WriteUint32("value", 32, (value))

        if EvaluationHelper.equals(
            data_type, ModbusDataType.ULINT
        ) and EvaluationHelper.equals(
            number_of_values, int(1)
        ):  # ULINT
            # Simple Field (value)
            value: int = _value.getInt()
            writeBuffer.WriteUint64("value", 64, (value))
        if EvaluationHelper.equals(data_type, ModbusDataType.ULINT):  # List
            values: PlcList = _value

            for val in values.getList():
                value: int = val.getInt()
                writeBuffer.WriteUint64("value", 64, (value))

        if EvaluationHelper.equals(
            data_type, ModbusDataType.REAL
        ) and EvaluationHelper.equals(
            number_of_values, int(1)
        ):  # REAL
            # Simple Field (value)
            value: float = _value.getFloat()
            writeBuffer.WriteFloat32("value", 32, (value))
        if EvaluationHelper.equals(data_type, ModbusDataType.REAL):  # List
            values: PlcList = _value

            for val in values.getList():
                value: float = val.getFloat()
                writeBuffer.WriteFloat32("value", 32, (value))

        if EvaluationHelper.equals(
            data_type, ModbusDataType.LREAL
        ) and EvaluationHelper.equals(
            number_of_values, int(1)
        ):  # LREAL
            # Simple Field (value)
            value: float = _value.getFloat()
            writeBuffer.WriteFloat64("value", 64, (value))
        if EvaluationHelper.equals(data_type, ModbusDataType.LREAL):  # List
            values: PlcList = _value

            for val in values.getList():
                value: float = val.getFloat()
                writeBuffer.WriteFloat64("value", 64, (value))

        if EvaluationHelper.equals(
            data_type, ModbusDataType.CHAR
        ) and EvaluationHelper.equals(
            number_of_values, int(1)
        ):  # CHAR
            # Simple Field (value)
            value: str = _value.getStr()
            writeBuffer.WriteString("value", uint32(8), "UTF-8", (value))
        if EvaluationHelper.equals(data_type, ModbusDataType.CHAR):  # List
            values: PlcList = _value

            for val in values.getList():
                value: str = val.getStr()
                writeBuffer.WriteString("value", uint32(8), "UTF-8", (value))

        if EvaluationHelper.equals(
            data_type, ModbusDataType.WCHAR
        ) and EvaluationHelper.equals(
            number_of_values, int(1)
        ):  # WCHAR
            # Simple Field (value)
            value: str = _value.getStr()
            writeBuffer.WriteString("value", uint32(16), "UTF-16", (value))
        if EvaluationHelper.equals(data_type, ModbusDataType.WCHAR):  # List
            values: PlcList = _value

            for val in values.getList():
                value: str = val.getStr()
                writeBuffer.WriteString("value", uint32(16), "UTF-16", (value))

    @staticmethod
    def get_length_in_bytes(
        _value: PlcValue, dataType: ModbusDataType, numberOfValues: int
    ) -> int:
        return int(
            math.ceil(float(getLengthInBits(_value, dataType, numberOfValues)) / 8.0)
        )

    @staticmethod
    def get_length_in_bits(
        _value: PlcValue, dataType: ModbusDataType, numberOfValues: int
    ) -> int:
        sizeInBits: int = 0
        if EvaluationHelper.equals(
            data_type, ModbusDataType.BOOL
        ) and EvaluationHelper.equals(
            number_of_values, int(1)
        ):  # BOOL
            # Reserved Field
            sizeInBits += 15
            # Simple Field (value)
            sizeInBits += 1
        if EvaluationHelper.equals(data_type, ModbusDataType.BOOL):  # List
            values: PlcList = _value
            sizeInBits += values.getList().size() * 1
        if EvaluationHelper.equals(
            data_type, ModbusDataType.BYTE
        ) and EvaluationHelper.equals(
            number_of_values, int(1)
        ):  # BYTE
            # Reserved Field
            sizeInBits += 8
            # Simple Field (value)
            sizeInBits += 8
        if EvaluationHelper.equals(data_type, ModbusDataType.BYTE):  # List
            values: PlcList = _value
            sizeInBits += values.getList().size() * 1
        if EvaluationHelper.equals(data_type, ModbusDataType.WORD):  # WORD
            # Simple Field (value)
            sizeInBits += 16
        if EvaluationHelper.equals(data_type, ModbusDataType.DWORD):  # DWORD
            # Simple Field (value)
            sizeInBits += 32
        if EvaluationHelper.equals(data_type, ModbusDataType.LWORD):  # LWORD
            # Simple Field (value)
            sizeInBits += 64
        if EvaluationHelper.equals(
            data_type, ModbusDataType.SINT
        ) and EvaluationHelper.equals(
            number_of_values, int(1)
        ):  # SINT
            # Reserved Field
            sizeInBits += 8
            # Simple Field (value)
            sizeInBits += 8
        if EvaluationHelper.equals(data_type, ModbusDataType.SINT):  # List
            values: PlcList = _value
            sizeInBits += values.getList().size() * 8
        if EvaluationHelper.equals(
            data_type, ModbusDataType.INT
        ) and EvaluationHelper.equals(
            number_of_values, int(1)
        ):  # INT
            # Simple Field (value)
            sizeInBits += 16
        if EvaluationHelper.equals(data_type, ModbusDataType.INT):  # List
            values: PlcList = _value
            sizeInBits += values.getList().size() * 16
        if EvaluationHelper.equals(
            data_type, ModbusDataType.DINT
        ) and EvaluationHelper.equals(
            number_of_values, int(1)
        ):  # DINT
            # Simple Field (value)
            sizeInBits += 32
        if EvaluationHelper.equals(data_type, ModbusDataType.DINT):  # List
            values: PlcList = _value
            sizeInBits += values.getList().size() * 32
        if EvaluationHelper.equals(
            data_type, ModbusDataType.LINT
        ) and EvaluationHelper.equals(
            number_of_values, int(1)
        ):  # LINT
            # Simple Field (value)
            sizeInBits += 64
        if EvaluationHelper.equals(data_type, ModbusDataType.LINT):  # List
            values: PlcList = _value
            sizeInBits += values.getList().size() * 64
        if EvaluationHelper.equals(
            data_type, ModbusDataType.USINT
        ) and EvaluationHelper.equals(
            number_of_values, int(1)
        ):  # USINT
            # Reserved Field
            sizeInBits += 8
            # Simple Field (value)
            sizeInBits += 8
        if EvaluationHelper.equals(data_type, ModbusDataType.USINT):  # List
            values: PlcList = _value
            sizeInBits += values.getList().size() * 8
        if EvaluationHelper.equals(
            data_type, ModbusDataType.UINT
        ) and EvaluationHelper.equals(
            number_of_values, int(1)
        ):  # UINT
            # Simple Field (value)
            sizeInBits += 16
        if EvaluationHelper.equals(data_type, ModbusDataType.UINT):  # List
            values: PlcList = _value
            sizeInBits += values.getList().size() * 16
        if EvaluationHelper.equals(
            data_type, ModbusDataType.UDINT
        ) and EvaluationHelper.equals(
            number_of_values, int(1)
        ):  # UDINT
            # Simple Field (value)
            sizeInBits += 32
        if EvaluationHelper.equals(data_type, ModbusDataType.UDINT):  # List
            values: PlcList = _value
            sizeInBits += values.getList().size() * 32
        if EvaluationHelper.equals(
            data_type, ModbusDataType.ULINT
        ) and EvaluationHelper.equals(
            number_of_values, int(1)
        ):  # ULINT
            # Simple Field (value)
            sizeInBits += 64
        if EvaluationHelper.equals(data_type, ModbusDataType.ULINT):  # List
            values: PlcList = _value
            sizeInBits += values.getList().size() * 64
        if EvaluationHelper.equals(
            data_type, ModbusDataType.REAL
        ) and EvaluationHelper.equals(
            number_of_values, int(1)
        ):  # REAL
            # Simple Field (value)
            sizeInBits += 32
        if EvaluationHelper.equals(data_type, ModbusDataType.REAL):  # List
            values: PlcList = _value
            sizeInBits += values.getList().size() * 32
        if EvaluationHelper.equals(
            data_type, ModbusDataType.LREAL
        ) and EvaluationHelper.equals(
            number_of_values, int(1)
        ):  # LREAL
            # Simple Field (value)
            sizeInBits += 64
        if EvaluationHelper.equals(data_type, ModbusDataType.LREAL):  # List
            values: PlcList = _value
            sizeInBits += values.getList().size() * 64
        if EvaluationHelper.equals(
            data_type, ModbusDataType.CHAR
        ) and EvaluationHelper.equals(
            number_of_values, int(1)
        ):  # CHAR
            # Simple Field (value)
            sizeInBits += 8
        if EvaluationHelper.equals(data_type, ModbusDataType.CHAR):  # List
            values: PlcList = _value
            sizeInBits += values.getList().size() * 8
        if EvaluationHelper.equals(
            data_type, ModbusDataType.WCHAR
        ) and EvaluationHelper.equals(
            number_of_values, int(1)
        ):  # WCHAR
            # Simple Field (value)
            sizeInBits += 16
        if EvaluationHelper.equals(data_type, ModbusDataType.WCHAR):  # List
            values: PlcList = _value
            sizeInBits += values.getList().size() * 16
        return sizeInBits
