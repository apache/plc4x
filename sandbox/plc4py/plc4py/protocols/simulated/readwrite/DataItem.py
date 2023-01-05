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
from ctypes import c_bool
from ctypes import c_double
from ctypes import c_float
from ctypes import c_int16
from ctypes import c_int32
from ctypes import c_int64
from ctypes import c_int8
from ctypes import c_uint16
from ctypes import c_uint32
from ctypes import c_uint64
from ctypes import c_uint8
from loguru import logging as log
import math


class DataItem:
    @staticmethod
    def static_parse(
        read_buffer: ReadBuffer, data_type: str, number_of_values: c_uint16
    ):
        if EvaluationHelper.equals(data_type, "_bool") and EvaluationHelper.equals(
            number_of_values, 1
        ):  # BOOL

            # Simple Field (value)
            value: c_bool = read_buffer.readBit("")

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
                value.append(PlcBOOL(c_bool(read_buffer.readBit(""))))

            return PlcList(value)
        if EvaluationHelper.equals(data_type, "_byte") and EvaluationHelper.equals(
            number_of_values, 1
        ):  # BYTE

            # Simple Field (value)
            value: c_uint8 = read_buffer.readUnsignedShort("", 8)

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
                value.append(PlcUINT(c_uint8(read_buffer.readUnsignedShort("", 8))))

            return PlcList(value)
        if EvaluationHelper.equals(data_type, "_word") and EvaluationHelper.equals(
            number_of_values, 1
        ):  # WORD

            # Simple Field (value)
            value: c_uint16 = read_buffer.readUnsignedInt("", 16)

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
                value.append(PlcUDINT(c_uint16(read_buffer.readUnsignedInt("", 16))))

            return PlcList(value)
        if EvaluationHelper.equals(data_type, "_dword") and EvaluationHelper.equals(
            number_of_values, 1
        ):  # DWORD

            # Simple Field (value)
            value: c_uint32 = read_buffer.readUnsignedLong("", 32)

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
                value.append(PlcULINT(c_uint32(read_buffer.readUnsignedLong("", 32))))

            return PlcList(value)
        if EvaluationHelper.equals(data_type, "_lword") and EvaluationHelper.equals(
            number_of_values, 1
        ):  # LWORD

            # Simple Field (value)
            value: c_uint64 = read_buffer.readUnsignedBigInteger("", 64)

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
                    PlcLINT(c_uint64(read_buffer.readUnsignedBigInteger("", 64)))
                )

            return PlcList(value)
        if EvaluationHelper.equals(data_type, "_sint") and EvaluationHelper.equals(
            number_of_values, 1
        ):  # SINT

            # Simple Field (value)
            value: c_int8 = read_buffer.readSignedByte("", 8)

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
                value.append(PlcSINT(c_int8(read_buffer.readSignedByte("", 8))))

            return PlcList(value)
        if EvaluationHelper.equals(data_type, "_int") and EvaluationHelper.equals(
            number_of_values, 1
        ):  # INT

            # Simple Field (value)
            value: c_int16 = read_buffer.readShort("", 16)

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
                value.append(PlcINT(c_int16(read_buffer.readShort("", 16))))

            return PlcList(value)
        if EvaluationHelper.equals(data_type, "_dint") and EvaluationHelper.equals(
            number_of_values, 1
        ):  # DINT

            # Simple Field (value)
            value: c_int32 = read_buffer.readInt("", 32)

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
                value.append(PlcDINT(c_int32(read_buffer.readInt("", 32))))

            return PlcList(value)
        if EvaluationHelper.equals(data_type, "_lint") and EvaluationHelper.equals(
            number_of_values, 1
        ):  # LINT

            # Simple Field (value)
            value: c_int64 = read_buffer.readLong("", 64)

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
                value.append(PlcLINT(c_int64(read_buffer.readLong("", 64))))

            return PlcList(value)
        if EvaluationHelper.equals(data_type, "_usint") and EvaluationHelper.equals(
            number_of_values, 1
        ):  # USINT

            # Simple Field (value)
            value: c_uint8 = read_buffer.readUnsignedShort("", 8)

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
                value.append(PlcUINT(c_uint8(read_buffer.readUnsignedShort("", 8))))

            return PlcList(value)
        if EvaluationHelper.equals(data_type, "_uint") and EvaluationHelper.equals(
            number_of_values, 1
        ):  # UINT

            # Simple Field (value)
            value: c_uint16 = read_buffer.readUnsignedInt("", 16)

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
                value.append(PlcUDINT(c_uint16(read_buffer.readUnsignedInt("", 16))))

            return PlcList(value)
        if EvaluationHelper.equals(data_type, "_udint") and EvaluationHelper.equals(
            number_of_values, 1
        ):  # UDINT

            # Simple Field (value)
            value: c_uint32 = read_buffer.readUnsignedLong("", 32)

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
                value.append(PlcULINT(c_uint32(read_buffer.readUnsignedLong("", 32))))

            return PlcList(value)
        if EvaluationHelper.equals(data_type, "_ulint") and EvaluationHelper.equals(
            number_of_values, 1
        ):  # ULINT

            # Simple Field (value)
            value: c_uint64 = read_buffer.readUnsignedBigInteger("", 64)

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
                    PlcLINT(c_uint64(read_buffer.readUnsignedBigInteger("", 64)))
                )

            return PlcList(value)
        if EvaluationHelper.equals(data_type, "_real") and EvaluationHelper.equals(
            number_of_values, 1
        ):  # REAL

            # Simple Field (value)
            value: c_float = read_buffer.readFloat("", 32)

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
                value.append(PlcREAL(c_float(read_buffer.readFloat("", 32))))

            return PlcList(value)
        if EvaluationHelper.equals(data_type, "_lreal") and EvaluationHelper.equals(
            number_of_values, 1
        ):  # LREAL

            # Simple Field (value)
            value: c_double = read_buffer.readDouble("", 64)

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
                value.append(PlcLREAL(c_double(read_buffer.readDouble("", 64))))

            return PlcList(value)
        if EvaluationHelper.equals(data_type, "_char") and EvaluationHelper.equals(
            number_of_values, 1
        ):  # CHAR

            # Simple Field (value)
            value: str = read_buffer.readString("", 8, "UTF-8")

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
                value.append(PlcSTRING(str(read_buffer.readString("", 8, "UTF-8"))))

            return PlcList(value)
        if EvaluationHelper.equals(data_type, "_wchar") and EvaluationHelper.equals(
            number_of_values, 1
        ):  # WCHAR

            # Simple Field (value)
            value: str = read_buffer.readString("", 16, "UTF-16")

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
                value.append(PlcSTRING(str(read_buffer.readString("", 16, "UTF-16"))))

            return PlcList(value)
        if EvaluationHelper.equals(data_type, "_string"):  # STRING

            # Simple Field (value)
            value: str = read_buffer.readString("", 255, "UTF-8")

            return PlcSTRING(value)
        if EvaluationHelper.equals(data_type, "_wstring"):  # STRING

            # Simple Field (value)
            value: str = read_buffer.readString("", 255, "UTF-16")

            return PlcSTRING(value)
        return None

    @staticmethod
    def static_serialize(
        writeBuffer: WriteBuffer,
        _value: PlcValue,
        dataType: str,
        numberOfValues: c_uint16,
    ) -> None:
        static_serialize(
            writeBuffer, _value, dataType, numberOfValues, ByteOrder.BIG_ENDIAN
        )

    @staticmethod
    def static_serialize(
        writeBuffer: WriteBuffer,
        _value: PlcValue,
        dataType: str,
        numberOfValues: c_uint16,
        byteOrder: ByteOrder,
    ) -> None:
        if EvaluationHelper.equals(dataType, "BOOL") and EvaluationHelper.equals(
            numberOfValues, 1
        ):  # BOOL

            # Simple Field (value)
            value: c_bool = _value.getC_bool()
            writeBuffer.writeBit("", bool((value)))
        if EvaluationHelper.equals(dataType, "BOOL"):  # List

            values: PlcList = _value

            for val in values.getList():
                value: c_bool = val.getC_bool()
                writeBuffer.writeBit("", bool((value)))

        if EvaluationHelper.equals(dataType, "BYTE") and EvaluationHelper.equals(
            numberOfValues, 1
        ):  # BYTE

            # Simple Field (value)
            value: c_uint8 = _value.getC_uint8()
            writeBuffer.writeUnsignedShort("", 8, (value).shortValue())
        if EvaluationHelper.equals(dataType, "BYTE"):  # List

            values: PlcList = _value

            for val in values.getList():
                value: c_uint8 = val.getC_uint8()
                writeBuffer.writeUnsignedShort("", 8, (value).shortValue())

        if EvaluationHelper.equals(dataType, "WORD") and EvaluationHelper.equals(
            numberOfValues, 1
        ):  # WORD

            # Simple Field (value)
            value: c_uint16 = _value.getC_uint16()
            writeBuffer.writeUnsignedInt("", 16, (value).intValue())
        if EvaluationHelper.equals(dataType, "WORD"):  # List

            values: PlcList = _value

            for val in values.getList():
                value: c_uint16 = val.getC_uint16()
                writeBuffer.writeUnsignedInt("", 16, (value).intValue())

        if EvaluationHelper.equals(dataType, "DWORD") and EvaluationHelper.equals(
            numberOfValues, 1
        ):  # DWORD

            # Simple Field (value)
            value: c_uint32 = _value.getC_uint32()
            writeBuffer.writeUnsignedLong("", 32, (value).longValue())
        if EvaluationHelper.equals(dataType, "DWORD"):  # List

            values: PlcList = _value

            for val in values.getList():
                value: c_uint32 = val.getC_uint32()
                writeBuffer.writeUnsignedLong("", 32, (value).longValue())

        if EvaluationHelper.equals(dataType, "LWORD") and EvaluationHelper.equals(
            numberOfValues, 1
        ):  # LWORD

            # Simple Field (value)
            value: c_uint64 = _value.getC_uint64()
            writeBuffer.writeUnsignedBigInteger("", 64, (value))
        if EvaluationHelper.equals(dataType, "LWORD"):  # List

            values: PlcList = _value

            for val in values.getList():
                value: c_uint64 = val.getC_uint64()
                writeBuffer.writeUnsignedBigInteger("", 64, (value))

        if EvaluationHelper.equals(dataType, "SINT") and EvaluationHelper.equals(
            numberOfValues, 1
        ):  # SINT

            # Simple Field (value)
            value: c_int8 = _value.getC_int8()
            writeBuffer.writeSignedByte("", 8, (value).byteValue())
        if EvaluationHelper.equals(dataType, "SINT"):  # List

            values: PlcList = _value

            for val in values.getList():
                value: c_int8 = val.getC_int8()
                writeBuffer.writeSignedByte("", 8, (value).byteValue())

        if EvaluationHelper.equals(dataType, "INT") and EvaluationHelper.equals(
            numberOfValues, 1
        ):  # INT

            # Simple Field (value)
            value: c_int16 = _value.getC_int16()
            writeBuffer.writeShort("", 16, (value).shortValue())
        if EvaluationHelper.equals(dataType, "INT"):  # List

            values: PlcList = _value

            for val in values.getList():
                value: c_int16 = val.getC_int16()
                writeBuffer.writeShort("", 16, (value).shortValue())

        if EvaluationHelper.equals(dataType, "DINT") and EvaluationHelper.equals(
            numberOfValues, 1
        ):  # DINT

            # Simple Field (value)
            value: c_int32 = _value.getC_int32()
            writeBuffer.writeInt("", 32, (value).intValue())
        if EvaluationHelper.equals(dataType, "DINT"):  # List

            values: PlcList = _value

            for val in values.getList():
                value: c_int32 = val.getC_int32()
                writeBuffer.writeInt("", 32, (value).intValue())

        if EvaluationHelper.equals(dataType, "LINT") and EvaluationHelper.equals(
            numberOfValues, 1
        ):  # LINT

            # Simple Field (value)
            value: c_int64 = _value.getC_int64()
            writeBuffer.writeLong("", 64, (value).longValue())
        if EvaluationHelper.equals(dataType, "LINT"):  # List

            values: PlcList = _value

            for val in values.getList():
                value: c_int64 = val.getC_int64()
                writeBuffer.writeLong("", 64, (value).longValue())

        if EvaluationHelper.equals(dataType, "USINT") and EvaluationHelper.equals(
            numberOfValues, 1
        ):  # USINT

            # Simple Field (value)
            value: c_uint8 = _value.getC_uint8()
            writeBuffer.writeUnsignedShort("", 8, (value).shortValue())
        if EvaluationHelper.equals(dataType, "USINT"):  # List

            values: PlcList = _value

            for val in values.getList():
                value: c_uint8 = val.getC_uint8()
                writeBuffer.writeUnsignedShort("", 8, (value).shortValue())

        if EvaluationHelper.equals(dataType, "UINT") and EvaluationHelper.equals(
            numberOfValues, 1
        ):  # UINT

            # Simple Field (value)
            value: c_uint16 = _value.getC_uint16()
            writeBuffer.writeUnsignedInt("", 16, (value).intValue())
        if EvaluationHelper.equals(dataType, "UINT"):  # List

            values: PlcList = _value

            for val in values.getList():
                value: c_uint16 = val.getC_uint16()
                writeBuffer.writeUnsignedInt("", 16, (value).intValue())

        if EvaluationHelper.equals(dataType, "UDINT") and EvaluationHelper.equals(
            numberOfValues, 1
        ):  # UDINT

            # Simple Field (value)
            value: c_uint32 = _value.getC_uint32()
            writeBuffer.writeUnsignedLong("", 32, (value).longValue())
        if EvaluationHelper.equals(dataType, "UDINT"):  # List

            values: PlcList = _value

            for val in values.getList():
                value: c_uint32 = val.getC_uint32()
                writeBuffer.writeUnsignedLong("", 32, (value).longValue())

        if EvaluationHelper.equals(dataType, "ULINT") and EvaluationHelper.equals(
            numberOfValues, 1
        ):  # ULINT

            # Simple Field (value)
            value: c_uint64 = _value.getC_uint64()
            writeBuffer.writeUnsignedBigInteger("", 64, (value))
        if EvaluationHelper.equals(dataType, "ULINT"):  # List

            values: PlcList = _value

            for val in values.getList():
                value: c_uint64 = val.getC_uint64()
                writeBuffer.writeUnsignedBigInteger("", 64, (value))

        if EvaluationHelper.equals(dataType, "REAL") and EvaluationHelper.equals(
            numberOfValues, 1
        ):  # REAL

            # Simple Field (value)
            value: c_float = _value.getC_float()
            writeBuffer.writeFloat("", 32, (value))
        if EvaluationHelper.equals(dataType, "REAL"):  # List

            values: PlcList = _value

            for val in values.getList():
                value: c_float = val.getC_float()
                writeBuffer.writeFloat("", 32, (value))

        if EvaluationHelper.equals(dataType, "LREAL") and EvaluationHelper.equals(
            numberOfValues, 1
        ):  # LREAL

            # Simple Field (value)
            value: c_double = _value.getC_double()
            writeBuffer.writeDouble("", 64, (value))
        if EvaluationHelper.equals(dataType, "LREAL"):  # List

            values: PlcList = _value

            for val in values.getList():
                value: c_double = val.getC_double()
                writeBuffer.writeDouble("", 64, (value))

        if EvaluationHelper.equals(dataType, "CHAR") and EvaluationHelper.equals(
            numberOfValues, 1
        ):  # CHAR

            # Simple Field (value)
            value: str = _value.getStr()
            writeBuffer.writeString("", 8, "UTF-8", (String)(value))
        if EvaluationHelper.equals(dataType, "CHAR"):  # List

            values: PlcList = _value

            for val in values.getList():
                value: str = val.getStr()
                writeBuffer.writeString("", 8, "UTF-8", (String)(value))

        if EvaluationHelper.equals(dataType, "WCHAR") and EvaluationHelper.equals(
            numberOfValues, 1
        ):  # WCHAR

            # Simple Field (value)
            value: str = _value.getStr()
            writeBuffer.writeString("", 16, "UTF-16", (String)(value))
        if EvaluationHelper.equals(dataType, "WCHAR"):  # List

            values: PlcList = _value

            for val in values.getList():
                value: str = val.getStr()
                writeBuffer.writeString("", 16, "UTF-16", (String)(value))

        if EvaluationHelper.equals(dataType, "STRING"):  # STRING

            # Simple Field (value)
            value: str = _value.getStr()
            writeBuffer.writeString("", 255, "UTF-8", (String)(value))
        if EvaluationHelper.equals(dataType, "WSTRING"):  # STRING

            # Simple Field (value)
            value: str = _value.getStr()
            writeBuffer.writeString("", 255, "UTF-16", (String)(value))

    @staticmethod
    def get_length_in_bytes(
        _value: PlcValue, dataType: str, numberOfValues: c_uint16
    ) -> int:
        return int(
            math.ceil(float(getLengthInBits(_value, dataType, numberOfValues)) / 8.0)
        )

    @staticmethod
    def get_length_in_bits(
        _value: PlcValue, dataType: str, numberOfValues: c_uint16
    ) -> int:
        sizeInBits: int = 0
        if EvaluationHelper.equals(dataType, "BOOL") and EvaluationHelper.equals(
            numberOfValues, 1
        ):  # BOOL
            # Simple Field (value)
            sizeInBits += 1
        if EvaluationHelper.equals(dataType, "BOOL"):  # List
            values: PlcList = _value
            sizeInBits += values.getList().size() * 1
        if EvaluationHelper.equals(dataType, "BYTE") and EvaluationHelper.equals(
            numberOfValues, 1
        ):  # BYTE
            # Simple Field (value)
            sizeInBits += 8
        if EvaluationHelper.equals(dataType, "BYTE"):  # List
            values: PlcList = _value
            sizeInBits += values.getList().size() * 8
        if EvaluationHelper.equals(dataType, "WORD") and EvaluationHelper.equals(
            numberOfValues, 1
        ):  # WORD
            # Simple Field (value)
            sizeInBits += 16
        if EvaluationHelper.equals(dataType, "WORD"):  # List
            values: PlcList = _value
            sizeInBits += values.getList().size() * 16
        if EvaluationHelper.equals(dataType, "DWORD") and EvaluationHelper.equals(
            numberOfValues, 1
        ):  # DWORD
            # Simple Field (value)
            sizeInBits += 32
        if EvaluationHelper.equals(dataType, "DWORD"):  # List
            values: PlcList = _value
            sizeInBits += values.getList().size() * 32
        if EvaluationHelper.equals(dataType, "LWORD") and EvaluationHelper.equals(
            numberOfValues, 1
        ):  # LWORD
            # Simple Field (value)
            sizeInBits += 64
        if EvaluationHelper.equals(dataType, "LWORD"):  # List
            values: PlcList = _value
            sizeInBits += values.getList().size() * 64
        if EvaluationHelper.equals(dataType, "SINT") and EvaluationHelper.equals(
            numberOfValues, 1
        ):  # SINT
            # Simple Field (value)
            sizeInBits += 8
        if EvaluationHelper.equals(dataType, "SINT"):  # List
            values: PlcList = _value
            sizeInBits += values.getList().size() * 8
        if EvaluationHelper.equals(dataType, "INT") and EvaluationHelper.equals(
            numberOfValues, 1
        ):  # INT
            # Simple Field (value)
            sizeInBits += 16
        if EvaluationHelper.equals(dataType, "INT"):  # List
            values: PlcList = _value
            sizeInBits += values.getList().size() * 16
        if EvaluationHelper.equals(dataType, "DINT") and EvaluationHelper.equals(
            numberOfValues, 1
        ):  # DINT
            # Simple Field (value)
            sizeInBits += 32
        if EvaluationHelper.equals(dataType, "DINT"):  # List
            values: PlcList = _value
            sizeInBits += values.getList().size() * 32
        if EvaluationHelper.equals(dataType, "LINT") and EvaluationHelper.equals(
            numberOfValues, 1
        ):  # LINT
            # Simple Field (value)
            sizeInBits += 64
        if EvaluationHelper.equals(dataType, "LINT"):  # List
            values: PlcList = _value
            sizeInBits += values.getList().size() * 64
        if EvaluationHelper.equals(dataType, "USINT") and EvaluationHelper.equals(
            numberOfValues, 1
        ):  # USINT
            # Simple Field (value)
            sizeInBits += 8
        if EvaluationHelper.equals(dataType, "USINT"):  # List
            values: PlcList = _value
            sizeInBits += values.getList().size() * 8
        if EvaluationHelper.equals(dataType, "UINT") and EvaluationHelper.equals(
            numberOfValues, 1
        ):  # UINT
            # Simple Field (value)
            sizeInBits += 16
        if EvaluationHelper.equals(dataType, "UINT"):  # List
            values: PlcList = _value
            sizeInBits += values.getList().size() * 16
        if EvaluationHelper.equals(dataType, "UDINT") and EvaluationHelper.equals(
            numberOfValues, 1
        ):  # UDINT
            # Simple Field (value)
            sizeInBits += 32
        if EvaluationHelper.equals(dataType, "UDINT"):  # List
            values: PlcList = _value
            sizeInBits += values.getList().size() * 32
        if EvaluationHelper.equals(dataType, "ULINT") and EvaluationHelper.equals(
            numberOfValues, 1
        ):  # ULINT
            # Simple Field (value)
            sizeInBits += 64
        if EvaluationHelper.equals(dataType, "ULINT"):  # List
            values: PlcList = _value
            sizeInBits += values.getList().size() * 64
        if EvaluationHelper.equals(dataType, "REAL") and EvaluationHelper.equals(
            numberOfValues, 1
        ):  # REAL
            # Simple Field (value)
            sizeInBits += 32
        if EvaluationHelper.equals(dataType, "REAL"):  # List
            values: PlcList = _value
            sizeInBits += values.getList().size() * 32
        if EvaluationHelper.equals(dataType, "LREAL") and EvaluationHelper.equals(
            numberOfValues, 1
        ):  # LREAL
            # Simple Field (value)
            sizeInBits += 64
        if EvaluationHelper.equals(dataType, "LREAL"):  # List
            values: PlcList = _value
            sizeInBits += values.getList().size() * 64
        if EvaluationHelper.equals(dataType, "CHAR") and EvaluationHelper.equals(
            numberOfValues, 1
        ):  # CHAR
            # Simple Field (value)
            sizeInBits += 8
        if EvaluationHelper.equals(dataType, "CHAR"):  # List
            values: PlcList = _value
            sizeInBits += values.getList().size() * 8
        if EvaluationHelper.equals(dataType, "WCHAR") and EvaluationHelper.equals(
            numberOfValues, 1
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
