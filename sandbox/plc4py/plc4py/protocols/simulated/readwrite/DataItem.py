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
import math


class DataItem:
    @staticmethod
    def static_parse(read_buffer: ReadBuffer, data_type: str, number_of_values: int):
        if EvaluationHelper.equals(data_type, "_bool") and EvaluationHelper.equals(
            number_of_values, int(1)
        ):  # BOOL
            # Simple Field (value)
            value: bool = read_buffer.read_bit("")

            return PlcBOOL(value)
        if EvaluationHelper.equals(data_type, "_bool"):  # List
            # Array field (value)
            # Count array
            if numberOfValues > Integer.MAX_VALUE:
                raise ParseException(
                    "Array count of "
                    + (numberOfValues)
                    + " exceeds the maximum allowed count of "
                    + Integer.MAX_VALUE
                )

            item_count: int = int(numberOfValues)
            value: List[PlcValue] = []
            for cur_item in range(item_count):
                value.append(PlcBOOL(bool(read_buffer.read_bit(""))))

            return PlcList(value)
        if EvaluationHelper.equals(data_type, "_byte") and EvaluationHelper.equals(
            number_of_values, int(1)
        ):  # BYTE
            # Simple Field (value)
            value: int = read_buffer.read_unsigned_short(8, logical_name="")

            return PlcBYTE(value)
        if EvaluationHelper.equals(data_type, "_byte"):  # List
            # Array field (value)
            # Count array
            if numberOfValues > Integer.MAX_VALUE:
                raise ParseException(
                    "Array count of "
                    + (numberOfValues)
                    + " exceeds the maximum allowed count of "
                    + Integer.MAX_VALUE
                )

            item_count: int = int(numberOfValues)
            value: List[PlcValue] = []
            for cur_item in range(item_count):
                value.append(
                    PlcUINT(int(read_buffer.read_unsigned_short(8, logical_name="")))
                )

            return PlcList(value)
        if EvaluationHelper.equals(data_type, "_word") and EvaluationHelper.equals(
            number_of_values, int(1)
        ):  # WORD
            # Simple Field (value)
            value: int = read_buffer.read_unsigned_int(16, logical_name="")

            return PlcWORD(value)
        if EvaluationHelper.equals(data_type, "_word"):  # List
            # Array field (value)
            # Count array
            if numberOfValues > Integer.MAX_VALUE:
                raise ParseException(
                    "Array count of "
                    + (numberOfValues)
                    + " exceeds the maximum allowed count of "
                    + Integer.MAX_VALUE
                )

            item_count: int = int(numberOfValues)
            value: List[PlcValue] = []
            for cur_item in range(item_count):
                value.append(
                    PlcUDINT(int(read_buffer.read_unsigned_int(16, logical_name="")))
                )

            return PlcList(value)
        if EvaluationHelper.equals(data_type, "_dword") and EvaluationHelper.equals(
            number_of_values, int(1)
        ):  # DWORD
            # Simple Field (value)
            value: int = read_buffer.read_unsigned_long(32, logical_name="")

            return PlcDWORD(value)
        if EvaluationHelper.equals(data_type, "_dword"):  # List
            # Array field (value)
            # Count array
            if numberOfValues > Integer.MAX_VALUE:
                raise ParseException(
                    "Array count of "
                    + (numberOfValues)
                    + " exceeds the maximum allowed count of "
                    + Integer.MAX_VALUE
                )

            item_count: int = int(numberOfValues)
            value: List[PlcValue] = []
            for cur_item in range(item_count):
                value.append(
                    PlcULINT(int(read_buffer.read_unsigned_long(32, logical_name="")))
                )

            return PlcList(value)
        if EvaluationHelper.equals(data_type, "_lword") and EvaluationHelper.equals(
            number_of_values, int(1)
        ):  # LWORD
            # Simple Field (value)
            value: int = read_buffer.read_unsigned_big_integer(64, logical_name="")

            return PlcLWORD(value)
        if EvaluationHelper.equals(data_type, "_lword"):  # List
            # Array field (value)
            # Count array
            if numberOfValues > Integer.MAX_VALUE:
                raise ParseException(
                    "Array count of "
                    + (numberOfValues)
                    + " exceeds the maximum allowed count of "
                    + Integer.MAX_VALUE
                )

            item_count: int = int(numberOfValues)
            value: List[PlcValue] = []
            for cur_item in range(item_count):
                value.append(
                    PlcLINT(
                        int(read_buffer.read_unsigned_big_integer(64, logical_name=""))
                    )
                )

            return PlcList(value)
        if EvaluationHelper.equals(data_type, "_sint") and EvaluationHelper.equals(
            number_of_values, int(1)
        ):  # SINT
            # Simple Field (value)
            value: int = read_buffer.read_signed_byte(8, logical_name="")

            return PlcSINT(value)
        if EvaluationHelper.equals(data_type, "_sint"):  # List
            # Array field (value)
            # Count array
            if numberOfValues > Integer.MAX_VALUE:
                raise ParseException(
                    "Array count of "
                    + (numberOfValues)
                    + " exceeds the maximum allowed count of "
                    + Integer.MAX_VALUE
                )

            item_count: int = int(numberOfValues)
            value: List[PlcValue] = []
            for cur_item in range(item_count):
                value.append(
                    PlcSINT(int(read_buffer.read_signed_byte(8, logical_name="")))
                )

            return PlcList(value)
        if EvaluationHelper.equals(data_type, "_int") and EvaluationHelper.equals(
            number_of_values, int(1)
        ):  # INT
            # Simple Field (value)
            value: int = read_buffer.read_short(16, logical_name="")

            return PlcINT(value)
        if EvaluationHelper.equals(data_type, "_int"):  # List
            # Array field (value)
            # Count array
            if numberOfValues > Integer.MAX_VALUE:
                raise ParseException(
                    "Array count of "
                    + (numberOfValues)
                    + " exceeds the maximum allowed count of "
                    + Integer.MAX_VALUE
                )

            item_count: int = int(numberOfValues)
            value: List[PlcValue] = []
            for cur_item in range(item_count):
                value.append(PlcINT(int(read_buffer.read_short(16, logical_name=""))))

            return PlcList(value)
        if EvaluationHelper.equals(data_type, "_dint") and EvaluationHelper.equals(
            number_of_values, int(1)
        ):  # DINT
            # Simple Field (value)
            value: int = read_buffer.read_int(32, logical_name="")

            return PlcDINT(value)
        if EvaluationHelper.equals(data_type, "_dint"):  # List
            # Array field (value)
            # Count array
            if numberOfValues > Integer.MAX_VALUE:
                raise ParseException(
                    "Array count of "
                    + (numberOfValues)
                    + " exceeds the maximum allowed count of "
                    + Integer.MAX_VALUE
                )

            item_count: int = int(numberOfValues)
            value: List[PlcValue] = []
            for cur_item in range(item_count):
                value.append(PlcDINT(int(read_buffer.read_int(32, logical_name=""))))

            return PlcList(value)
        if EvaluationHelper.equals(data_type, "_lint") and EvaluationHelper.equals(
            number_of_values, int(1)
        ):  # LINT
            # Simple Field (value)
            value: int = read_buffer.read_long(64, logical_name="")

            return PlcLINT(value)
        if EvaluationHelper.equals(data_type, "_lint"):  # List
            # Array field (value)
            # Count array
            if numberOfValues > Integer.MAX_VALUE:
                raise ParseException(
                    "Array count of "
                    + (numberOfValues)
                    + " exceeds the maximum allowed count of "
                    + Integer.MAX_VALUE
                )

            item_count: int = int(numberOfValues)
            value: List[PlcValue] = []
            for cur_item in range(item_count):
                value.append(PlcLINT(int(read_buffer.read_long(64, logical_name=""))))

            return PlcList(value)
        if EvaluationHelper.equals(data_type, "_usint") and EvaluationHelper.equals(
            number_of_values, int(1)
        ):  # USINT
            # Simple Field (value)
            value: int = read_buffer.read_unsigned_short(8, logical_name="")

            return PlcUSINT(value)
        if EvaluationHelper.equals(data_type, "_usint"):  # List
            # Array field (value)
            # Count array
            if numberOfValues > Integer.MAX_VALUE:
                raise ParseException(
                    "Array count of "
                    + (numberOfValues)
                    + " exceeds the maximum allowed count of "
                    + Integer.MAX_VALUE
                )

            item_count: int = int(numberOfValues)
            value: List[PlcValue] = []
            for cur_item in range(item_count):
                value.append(
                    PlcUINT(int(read_buffer.read_unsigned_short(8, logical_name="")))
                )

            return PlcList(value)
        if EvaluationHelper.equals(data_type, "_uint") and EvaluationHelper.equals(
            number_of_values, int(1)
        ):  # UINT
            # Simple Field (value)
            value: int = read_buffer.read_unsigned_int(16, logical_name="")

            return PlcUINT(value)
        if EvaluationHelper.equals(data_type, "_uint"):  # List
            # Array field (value)
            # Count array
            if numberOfValues > Integer.MAX_VALUE:
                raise ParseException(
                    "Array count of "
                    + (numberOfValues)
                    + " exceeds the maximum allowed count of "
                    + Integer.MAX_VALUE
                )

            item_count: int = int(numberOfValues)
            value: List[PlcValue] = []
            for cur_item in range(item_count):
                value.append(
                    PlcUDINT(int(read_buffer.read_unsigned_int(16, logical_name="")))
                )

            return PlcList(value)
        if EvaluationHelper.equals(data_type, "_udint") and EvaluationHelper.equals(
            number_of_values, int(1)
        ):  # UDINT
            # Simple Field (value)
            value: int = read_buffer.read_unsigned_long(32, logical_name="")

            return PlcUDINT(value)
        if EvaluationHelper.equals(data_type, "_udint"):  # List
            # Array field (value)
            # Count array
            if numberOfValues > Integer.MAX_VALUE:
                raise ParseException(
                    "Array count of "
                    + (numberOfValues)
                    + " exceeds the maximum allowed count of "
                    + Integer.MAX_VALUE
                )

            item_count: int = int(numberOfValues)
            value: List[PlcValue] = []
            for cur_item in range(item_count):
                value.append(
                    PlcULINT(int(read_buffer.read_unsigned_long(32, logical_name="")))
                )

            return PlcList(value)
        if EvaluationHelper.equals(data_type, "_ulint") and EvaluationHelper.equals(
            number_of_values, int(1)
        ):  # ULINT
            # Simple Field (value)
            value: int = read_buffer.read_unsigned_big_integer(64, logical_name="")

            return PlcULINT(value)
        if EvaluationHelper.equals(data_type, "_ulint"):  # List
            # Array field (value)
            # Count array
            if numberOfValues > Integer.MAX_VALUE:
                raise ParseException(
                    "Array count of "
                    + (numberOfValues)
                    + " exceeds the maximum allowed count of "
                    + Integer.MAX_VALUE
                )

            item_count: int = int(numberOfValues)
            value: List[PlcValue] = []
            for cur_item in range(item_count):
                value.append(
                    PlcLINT(
                        int(read_buffer.read_unsigned_big_integer(64, logical_name=""))
                    )
                )

            return PlcList(value)
        if EvaluationHelper.equals(data_type, "_real") and EvaluationHelper.equals(
            number_of_values, int(1)
        ):  # REAL
            # Simple Field (value)
            value: float = read_buffer.read_float(32, logical_name="")

            return PlcREAL(value)
        if EvaluationHelper.equals(data_type, "_real"):  # List
            # Array field (value)
            # Count array
            if numberOfValues > Integer.MAX_VALUE:
                raise ParseException(
                    "Array count of "
                    + (numberOfValues)
                    + " exceeds the maximum allowed count of "
                    + Integer.MAX_VALUE
                )

            item_count: int = int(numberOfValues)
            value: List[PlcValue] = []
            for cur_item in range(item_count):
                value.append(
                    PlcREAL(float(read_buffer.read_float(32, logical_name="")))
                )

            return PlcList(value)
        if EvaluationHelper.equals(data_type, "_lreal") and EvaluationHelper.equals(
            number_of_values, int(1)
        ):  # LREAL
            # Simple Field (value)
            value: float = read_buffer.read_double(64, logical_name="")

            return PlcLREAL(value)
        if EvaluationHelper.equals(data_type, "_lreal"):  # List
            # Array field (value)
            # Count array
            if numberOfValues > Integer.MAX_VALUE:
                raise ParseException(
                    "Array count of "
                    + (numberOfValues)
                    + " exceeds the maximum allowed count of "
                    + Integer.MAX_VALUE
                )

            item_count: int = int(numberOfValues)
            value: List[PlcValue] = []
            for cur_item in range(item_count):
                value.append(
                    PlcLREAL(float(read_buffer.read_double(64, logical_name="")))
                )

            return PlcList(value)
        if EvaluationHelper.equals(data_type, "_char") and EvaluationHelper.equals(
            number_of_values, int(1)
        ):  # CHAR
            # Simple Field (value)
            value: str = read_buffer.read_string(8, logical_name="", encoding="")

            return PlcCHAR(value)
        if EvaluationHelper.equals(data_type, "_char"):  # List
            # Array field (value)
            # Count array
            if numberOfValues > Integer.MAX_VALUE:
                raise ParseException(
                    "Array count of "
                    + (numberOfValues)
                    + " exceeds the maximum allowed count of "
                    + Integer.MAX_VALUE
                )

            item_count: int = int(numberOfValues)
            value: List[PlcValue] = []
            for cur_item in range(item_count):
                value.append(
                    PlcSTRING(
                        str(read_buffer.read_string(8, logical_name="", encoding=""))
                    )
                )

            return PlcList(value)
        if EvaluationHelper.equals(data_type, "_wchar") and EvaluationHelper.equals(
            number_of_values, int(1)
        ):  # WCHAR
            # Simple Field (value)
            value: str = read_buffer.read_string(16, logical_name="", encoding="")

            return PlcWCHAR(value)
        if EvaluationHelper.equals(data_type, "_wchar"):  # List
            # Array field (value)
            # Count array
            if numberOfValues > Integer.MAX_VALUE:
                raise ParseException(
                    "Array count of "
                    + (numberOfValues)
                    + " exceeds the maximum allowed count of "
                    + Integer.MAX_VALUE
                )

            item_count: int = int(numberOfValues)
            value: List[PlcValue] = []
            for cur_item in range(item_count):
                value.append(
                    PlcSTRING(
                        str(read_buffer.read_string(16, logical_name="", encoding=""))
                    )
                )

            return PlcList(value)
        if EvaluationHelper.equals(data_type, "_string"):  # STRING
            # Simple Field (value)
            value: str = read_buffer.read_string(255, logical_name="", encoding="")

            return PlcSTRING(value)
        if EvaluationHelper.equals(data_type, "_wstring"):  # STRING
            # Simple Field (value)
            value: str = read_buffer.read_string(255, logical_name="", encoding="")

            return PlcSTRING(value)
        return None

    @staticmethod
    def static_serialize(
        writeBuffer: WriteBuffer, _value: PlcValue, dataType: str, numberOfValues: int
    ) -> None:
        static_serialize(
            writeBuffer, _value, dataType, numberOfValues, ByteOrder.BIG_ENDIAN
        )

    @staticmethod
    def static_serialize(
        writeBuffer: WriteBuffer,
        _value: PlcValue,
        dataType: str,
        numberOfValues: int,
        byteOrder: ByteOrder,
    ) -> None:
        if EvaluationHelper.equals(dataType, "BOOL") and EvaluationHelper.equals(
            numberOfValues, int(1)
        ):  # BOOL
            # Simple Field (value)
            value: bool = _value.getBool()
            writeBuffer.WriteBit("value", (value))
        if EvaluationHelper.equals(dataType, "BOOL"):  # List
            values: PlcList = _value

            for val in values.getList():
                value: bool = val.getBool()
                writeBuffer.WriteBit("value", (value))

        if EvaluationHelper.equals(dataType, "BYTE") and EvaluationHelper.equals(
            numberOfValues, int(1)
        ):  # BYTE
            # Simple Field (value)
            value: int = _value.getInt()
            writeBuffer.WriteUint8("value", 8, (value))
        if EvaluationHelper.equals(dataType, "BYTE"):  # List
            values: PlcList = _value

            for val in values.getList():
                value: int = val.getInt()
                writeBuffer.WriteUint8("value", 8, (value))

        if EvaluationHelper.equals(dataType, "WORD") and EvaluationHelper.equals(
            numberOfValues, int(1)
        ):  # WORD
            # Simple Field (value)
            value: int = _value.getInt()
            writeBuffer.WriteUint16("value", 16, (value))
        if EvaluationHelper.equals(dataType, "WORD"):  # List
            values: PlcList = _value

            for val in values.getList():
                value: int = val.getInt()
                writeBuffer.WriteUint16("value", 16, (value))

        if EvaluationHelper.equals(dataType, "DWORD") and EvaluationHelper.equals(
            numberOfValues, int(1)
        ):  # DWORD
            # Simple Field (value)
            value: int = _value.getInt()
            writeBuffer.WriteUint32("value", 32, (value))
        if EvaluationHelper.equals(dataType, "DWORD"):  # List
            values: PlcList = _value

            for val in values.getList():
                value: int = val.getInt()
                writeBuffer.WriteUint32("value", 32, (value))

        if EvaluationHelper.equals(dataType, "LWORD") and EvaluationHelper.equals(
            numberOfValues, int(1)
        ):  # LWORD
            # Simple Field (value)
            value: int = _value.getInt()
            writeBuffer.WriteUint64("value", 64, (value))
        if EvaluationHelper.equals(dataType, "LWORD"):  # List
            values: PlcList = _value

            for val in values.getList():
                value: int = val.getInt()
                writeBuffer.WriteUint64("value", 64, (value))

        if EvaluationHelper.equals(dataType, "SINT") and EvaluationHelper.equals(
            numberOfValues, int(1)
        ):  # SINT
            # Simple Field (value)
            value: int = _value.getInt()
            writeBuffer.WriteInt8("value", 8, (value))
        if EvaluationHelper.equals(dataType, "SINT"):  # List
            values: PlcList = _value

            for val in values.getList():
                value: int = val.getInt()
                writeBuffer.WriteInt8("value", 8, (value))

        if EvaluationHelper.equals(dataType, "INT") and EvaluationHelper.equals(
            numberOfValues, int(1)
        ):  # INT
            # Simple Field (value)
            value: int = _value.getInt()
            writeBuffer.WriteInt16("value", 16, (value))
        if EvaluationHelper.equals(dataType, "INT"):  # List
            values: PlcList = _value

            for val in values.getList():
                value: int = val.getInt()
                writeBuffer.WriteInt16("value", 16, (value))

        if EvaluationHelper.equals(dataType, "DINT") and EvaluationHelper.equals(
            numberOfValues, int(1)
        ):  # DINT
            # Simple Field (value)
            value: int = _value.getInt()
            writeBuffer.WriteInt32("value", 32, (value))
        if EvaluationHelper.equals(dataType, "DINT"):  # List
            values: PlcList = _value

            for val in values.getList():
                value: int = val.getInt()
                writeBuffer.WriteInt32("value", 32, (value))

        if EvaluationHelper.equals(dataType, "LINT") and EvaluationHelper.equals(
            numberOfValues, int(1)
        ):  # LINT
            # Simple Field (value)
            value: int = _value.getInt()
            writeBuffer.WriteInt64("value", 64, (value))
        if EvaluationHelper.equals(dataType, "LINT"):  # List
            values: PlcList = _value

            for val in values.getList():
                value: int = val.getInt()
                writeBuffer.WriteInt64("value", 64, (value))

        if EvaluationHelper.equals(dataType, "USINT") and EvaluationHelper.equals(
            numberOfValues, int(1)
        ):  # USINT
            # Simple Field (value)
            value: int = _value.getInt()
            writeBuffer.WriteUint8("value", 8, (value))
        if EvaluationHelper.equals(dataType, "USINT"):  # List
            values: PlcList = _value

            for val in values.getList():
                value: int = val.getInt()
                writeBuffer.WriteUint8("value", 8, (value))

        if EvaluationHelper.equals(dataType, "UINT") and EvaluationHelper.equals(
            numberOfValues, int(1)
        ):  # UINT
            # Simple Field (value)
            value: int = _value.getInt()
            writeBuffer.WriteUint16("value", 16, (value))
        if EvaluationHelper.equals(dataType, "UINT"):  # List
            values: PlcList = _value

            for val in values.getList():
                value: int = val.getInt()
                writeBuffer.WriteUint16("value", 16, (value))

        if EvaluationHelper.equals(dataType, "UDINT") and EvaluationHelper.equals(
            numberOfValues, int(1)
        ):  # UDINT
            # Simple Field (value)
            value: int = _value.getInt()
            writeBuffer.WriteUint32("value", 32, (value))
        if EvaluationHelper.equals(dataType, "UDINT"):  # List
            values: PlcList = _value

            for val in values.getList():
                value: int = val.getInt()
                writeBuffer.WriteUint32("value", 32, (value))

        if EvaluationHelper.equals(dataType, "ULINT") and EvaluationHelper.equals(
            numberOfValues, int(1)
        ):  # ULINT
            # Simple Field (value)
            value: int = _value.getInt()
            writeBuffer.WriteUint64("value", 64, (value))
        if EvaluationHelper.equals(dataType, "ULINT"):  # List
            values: PlcList = _value

            for val in values.getList():
                value: int = val.getInt()
                writeBuffer.WriteUint64("value", 64, (value))

        if EvaluationHelper.equals(dataType, "REAL") and EvaluationHelper.equals(
            numberOfValues, int(1)
        ):  # REAL
            # Simple Field (value)
            value: float = _value.getFloat()
            writeBuffer.WriteFloat32("value", 32, (value))
        if EvaluationHelper.equals(dataType, "REAL"):  # List
            values: PlcList = _value

            for val in values.getList():
                value: float = val.getFloat()
                writeBuffer.WriteFloat32("value", 32, (value))

        if EvaluationHelper.equals(dataType, "LREAL") and EvaluationHelper.equals(
            numberOfValues, int(1)
        ):  # LREAL
            # Simple Field (value)
            value: float = _value.getFloat()
            writeBuffer.WriteFloat64("value", 64, (value))
        if EvaluationHelper.equals(dataType, "LREAL"):  # List
            values: PlcList = _value

            for val in values.getList():
                value: float = val.getFloat()
                writeBuffer.WriteFloat64("value", 64, (value))

        if EvaluationHelper.equals(dataType, "CHAR") and EvaluationHelper.equals(
            numberOfValues, int(1)
        ):  # CHAR
            # Simple Field (value)
            value: str = _value.getStr()
            writeBuffer.WriteString("value", uint32(8), "UTF-8", (value))
        if EvaluationHelper.equals(dataType, "CHAR"):  # List
            values: PlcList = _value

            for val in values.getList():
                value: str = val.getStr()
                writeBuffer.WriteString("value", uint32(8), "UTF-8", (value))

        if EvaluationHelper.equals(dataType, "WCHAR") and EvaluationHelper.equals(
            numberOfValues, int(1)
        ):  # WCHAR
            # Simple Field (value)
            value: str = _value.getStr()
            writeBuffer.WriteString("value", uint32(16), "UTF-16", (value))
        if EvaluationHelper.equals(dataType, "WCHAR"):  # List
            values: PlcList = _value

            for val in values.getList():
                value: str = val.getStr()
                writeBuffer.WriteString("value", uint32(16), "UTF-16", (value))

        if EvaluationHelper.equals(dataType, "STRING"):  # STRING
            # Simple Field (value)
            value: str = _value.getStr()
            writeBuffer.WriteString("value", uint32(255), "UTF-8", (value))
        if EvaluationHelper.equals(dataType, "WSTRING"):  # STRING
            # Simple Field (value)
            value: str = _value.getStr()
            writeBuffer.WriteString("value", uint32(255), "UTF-16", (value))

    @staticmethod
    def get_length_in_bytes(
        _value: PlcValue, dataType: str, numberOfValues: int
    ) -> int:
        return int(
            math.ceil(float(getLengthInBits(_value, dataType, numberOfValues)) / 8.0)
        )

    @staticmethod
    def get_length_in_bits(_value: PlcValue, dataType: str, numberOfValues: int) -> int:
        sizeInBits: int = 0
        if EvaluationHelper.equals(dataType, "BOOL") and EvaluationHelper.equals(
            numberOfValues, int(1)
        ):  # BOOL
            # Simple Field (value)
            sizeInBits += 1
        if EvaluationHelper.equals(dataType, "BOOL"):  # List
            values: PlcList = _value
            sizeInBits += values.getList().size() * 1
        if EvaluationHelper.equals(dataType, "BYTE") and EvaluationHelper.equals(
            numberOfValues, int(1)
        ):  # BYTE
            # Simple Field (value)
            sizeInBits += 8
        if EvaluationHelper.equals(dataType, "BYTE"):  # List
            values: PlcList = _value
            sizeInBits += values.getList().size() * 8
        if EvaluationHelper.equals(dataType, "WORD") and EvaluationHelper.equals(
            numberOfValues, int(1)
        ):  # WORD
            # Simple Field (value)
            sizeInBits += 16
        if EvaluationHelper.equals(dataType, "WORD"):  # List
            values: PlcList = _value
            sizeInBits += values.getList().size() * 16
        if EvaluationHelper.equals(dataType, "DWORD") and EvaluationHelper.equals(
            numberOfValues, int(1)
        ):  # DWORD
            # Simple Field (value)
            sizeInBits += 32
        if EvaluationHelper.equals(dataType, "DWORD"):  # List
            values: PlcList = _value
            sizeInBits += values.getList().size() * 32
        if EvaluationHelper.equals(dataType, "LWORD") and EvaluationHelper.equals(
            numberOfValues, int(1)
        ):  # LWORD
            # Simple Field (value)
            sizeInBits += 64
        if EvaluationHelper.equals(dataType, "LWORD"):  # List
            values: PlcList = _value
            sizeInBits += values.getList().size() * 64
        if EvaluationHelper.equals(dataType, "SINT") and EvaluationHelper.equals(
            numberOfValues, int(1)
        ):  # SINT
            # Simple Field (value)
            sizeInBits += 8
        if EvaluationHelper.equals(dataType, "SINT"):  # List
            values: PlcList = _value
            sizeInBits += values.getList().size() * 8
        if EvaluationHelper.equals(dataType, "INT") and EvaluationHelper.equals(
            numberOfValues, int(1)
        ):  # INT
            # Simple Field (value)
            sizeInBits += 16
        if EvaluationHelper.equals(dataType, "INT"):  # List
            values: PlcList = _value
            sizeInBits += values.getList().size() * 16
        if EvaluationHelper.equals(dataType, "DINT") and EvaluationHelper.equals(
            numberOfValues, int(1)
        ):  # DINT
            # Simple Field (value)
            sizeInBits += 32
        if EvaluationHelper.equals(dataType, "DINT"):  # List
            values: PlcList = _value
            sizeInBits += values.getList().size() * 32
        if EvaluationHelper.equals(dataType, "LINT") and EvaluationHelper.equals(
            numberOfValues, int(1)
        ):  # LINT
            # Simple Field (value)
            sizeInBits += 64
        if EvaluationHelper.equals(dataType, "LINT"):  # List
            values: PlcList = _value
            sizeInBits += values.getList().size() * 64
        if EvaluationHelper.equals(dataType, "USINT") and EvaluationHelper.equals(
            numberOfValues, int(1)
        ):  # USINT
            # Simple Field (value)
            sizeInBits += 8
        if EvaluationHelper.equals(dataType, "USINT"):  # List
            values: PlcList = _value
            sizeInBits += values.getList().size() * 8
        if EvaluationHelper.equals(dataType, "UINT") and EvaluationHelper.equals(
            numberOfValues, int(1)
        ):  # UINT
            # Simple Field (value)
            sizeInBits += 16
        if EvaluationHelper.equals(dataType, "UINT"):  # List
            values: PlcList = _value
            sizeInBits += values.getList().size() * 16
        if EvaluationHelper.equals(dataType, "UDINT") and EvaluationHelper.equals(
            numberOfValues, int(1)
        ):  # UDINT
            # Simple Field (value)
            sizeInBits += 32
        if EvaluationHelper.equals(dataType, "UDINT"):  # List
            values: PlcList = _value
            sizeInBits += values.getList().size() * 32
        if EvaluationHelper.equals(dataType, "ULINT") and EvaluationHelper.equals(
            numberOfValues, int(1)
        ):  # ULINT
            # Simple Field (value)
            sizeInBits += 64
        if EvaluationHelper.equals(dataType, "ULINT"):  # List
            values: PlcList = _value
            sizeInBits += values.getList().size() * 64
        if EvaluationHelper.equals(dataType, "REAL") and EvaluationHelper.equals(
            numberOfValues, int(1)
        ):  # REAL
            # Simple Field (value)
            sizeInBits += 32
        if EvaluationHelper.equals(dataType, "REAL"):  # List
            values: PlcList = _value
            sizeInBits += values.getList().size() * 32
        if EvaluationHelper.equals(dataType, "LREAL") and EvaluationHelper.equals(
            numberOfValues, int(1)
        ):  # LREAL
            # Simple Field (value)
            sizeInBits += 64
        if EvaluationHelper.equals(dataType, "LREAL"):  # List
            values: PlcList = _value
            sizeInBits += values.getList().size() * 64
        if EvaluationHelper.equals(dataType, "CHAR") and EvaluationHelper.equals(
            numberOfValues, int(1)
        ):  # CHAR
            # Simple Field (value)
            sizeInBits += 8
        if EvaluationHelper.equals(dataType, "CHAR"):  # List
            values: PlcList = _value
            sizeInBits += values.getList().size() * 8
        if EvaluationHelper.equals(dataType, "WCHAR") and EvaluationHelper.equals(
            numberOfValues, int(1)
        ):  # WCHAR
            # Simple Field (value)
            sizeInBits += 16
        if EvaluationHelper.equals(dataType, "WCHAR"):  # List
            values: PlcList = _value
            sizeInBits += values.getList().size() * 16
        if EvaluationHelper.equals(dataType, "STRING"):  # STRING
            # Simple Field (value)
            sizeInBits += 255
        if EvaluationHelper.equals(dataType, "WSTRING"):  # STRING
            # Simple Field (value)
            sizeInBits += 255
        return sizeInBits
