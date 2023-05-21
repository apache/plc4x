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
from plc4py.protocols.modbus.readwrite.ModbusDataType import ModbusDataType
import math


class DataItem:
    @staticmethod
    def static_parse(
        read_buffer: ReadBuffer, data_type: ModbusDataType, number_of_values: c_uint16
    ):
        if EvaluationHelper.equals(
            data_type66, ModbusDataType.get_bool()
        ) and EvaluationHelper.equals(
            number_of_values66, c_uint16(1)
        ):  # BOOL
            # Reserved Field (Compartmentalized so the "reserved" variable can't leak)
            reserved: c_uint16 = read_buffer.readUnsignedInt("", 15)
            if reserved != c_uint16(0x0000):
                log.info(
                    "Expected constant value "
                    + str(0x0000)
                    + " but got "
                    + reserved
                    + " for reserved field."
                )

            # Simple Field (value)
            value: c_bool = read_buffer.readBit("")

            return PlcBOOL(value)
        if EvaluationHelper.equals(data_type66, ModbusDataType.get_bool()):  # List
            # Array field (value)
            # Count array
            if numberOfValues66 > Integer.MAX_VALUE:
                raise ParseException(
                    "Array count of "
                    + (numberOfValues66)
                    + " exceeds the maximum allowed count of "
                    + Integer.MAX_VALUE
                )

            item_count: int = int(numberOfValues66)
            value: List[PlcValue] = []
            for cur_item in range(item_count):
                value.append(PlcBOOL(c_bool(read_buffer.readBit(""))))

            return PlcList(value)
        if EvaluationHelper.equals(
            data_type66, ModbusDataType.get_byte()
        ) and EvaluationHelper.equals(
            number_of_values66, c_uint16(1)
        ):  # BYTE
            # Reserved Field (Compartmentalized so the "reserved" variable can't leak)
            reserved: c_uint8 = read_buffer.readUnsignedShort("", 8)
            if reserved != c_uint8(0x00):
                log.info(
                    "Expected constant value "
                    + str(0x00)
                    + " but got "
                    + reserved
                    + " for reserved field."
                )

            # Simple Field (value)
            value: c_uint8 = read_buffer.readUnsignedShort("", 8)

            return PlcBYTE(value)
        if EvaluationHelper.equals(data_type66, ModbusDataType.get_byte()):  # List
            # Array field (value)
            # Count array
            if numberOfValues66 * c_int32(8) > Integer.MAX_VALUE:
                raise ParseException(
                    "Array count of "
                    + (numberOfValues66 * c_int32(8))
                    + " exceeds the maximum allowed count of "
                    + Integer.MAX_VALUE
                )

            item_count: int = int(numberOfValues66 * c_int32(8))
            value: List[PlcValue] = []
            for cur_item in range(item_count):
                value.append(PlcBOOL(c_bool(read_buffer.readBit(""))))

            return PlcList(value)
        if EvaluationHelper.equals(data_type66, ModbusDataType.get_word()):  # WORD
            # Simple Field (value)
            value: c_uint16 = read_buffer.readUnsignedInt("", 16)

            return PlcWORD(value)
        if EvaluationHelper.equals(data_type66, ModbusDataType.get_dword()):  # DWORD
            # Simple Field (value)
            value: c_uint32 = read_buffer.readUnsignedLong("", 32)

            return PlcDWORD(value)
        if EvaluationHelper.equals(data_type66, ModbusDataType.get_lword()):  # LWORD
            # Simple Field (value)
            value: c_uint64 = read_buffer.readUnsignedBigInteger("", 64)

            return PlcLWORD(value)
        if EvaluationHelper.equals(
            data_type66, ModbusDataType.get_sint()
        ) and EvaluationHelper.equals(
            number_of_values66, c_uint16(1)
        ):  # SINT
            # Reserved Field (Compartmentalized so the "reserved" variable can't leak)
            reserved: c_uint8 = read_buffer.readUnsignedShort("", 8)
            if reserved != c_uint8(0x00):
                log.info(
                    "Expected constant value "
                    + str(0x00)
                    + " but got "
                    + reserved
                    + " for reserved field."
                )

            # Simple Field (value)
            value: c_int8 = read_buffer.readSignedByte("", 8)

            return PlcSINT(value)
        if EvaluationHelper.equals(data_type66, ModbusDataType.get_sint()):  # List
            # Array field (value)
            # Count array
            if numberOfValues66 > Integer.MAX_VALUE:
                raise ParseException(
                    "Array count of "
                    + (numberOfValues66)
                    + " exceeds the maximum allowed count of "
                    + Integer.MAX_VALUE
                )

            item_count: int = int(numberOfValues66)
            value: List[PlcValue] = []
            for cur_item in range(item_count):
                value.append(PlcSINT(c_int8(read_buffer.readSignedByte("", 8))))

            return PlcList(value)
        if EvaluationHelper.equals(
            data_type66, ModbusDataType.get_int()
        ) and EvaluationHelper.equals(
            number_of_values66, c_uint16(1)
        ):  # INT
            # Simple Field (value)
            value: c_int16 = read_buffer.readShort("", 16)

            return PlcINT(value)
        if EvaluationHelper.equals(data_type66, ModbusDataType.get_int()):  # List
            # Array field (value)
            # Count array
            if numberOfValues66 > Integer.MAX_VALUE:
                raise ParseException(
                    "Array count of "
                    + (numberOfValues66)
                    + " exceeds the maximum allowed count of "
                    + Integer.MAX_VALUE
                )

            item_count: int = int(numberOfValues66)
            value: List[PlcValue] = []
            for cur_item in range(item_count):
                value.append(PlcINT(c_int16(read_buffer.readShort("", 16))))

            return PlcList(value)
        if EvaluationHelper.equals(
            data_type66, ModbusDataType.get_dint()
        ) and EvaluationHelper.equals(
            number_of_values66, c_uint16(1)
        ):  # DINT
            # Simple Field (value)
            value: c_int32 = read_buffer.readInt("", 32)

            return PlcDINT(value)
        if EvaluationHelper.equals(data_type66, ModbusDataType.get_dint()):  # List
            # Array field (value)
            # Count array
            if numberOfValues66 > Integer.MAX_VALUE:
                raise ParseException(
                    "Array count of "
                    + (numberOfValues66)
                    + " exceeds the maximum allowed count of "
                    + Integer.MAX_VALUE
                )

            item_count: int = int(numberOfValues66)
            value: List[PlcValue] = []
            for cur_item in range(item_count):
                value.append(PlcDINT(c_int32(read_buffer.readInt("", 32))))

            return PlcList(value)
        if EvaluationHelper.equals(
            data_type66, ModbusDataType.get_lint()
        ) and EvaluationHelper.equals(
            number_of_values66, c_uint16(1)
        ):  # LINT
            # Simple Field (value)
            value: c_int64 = read_buffer.readLong("", 64)

            return PlcLINT(value)
        if EvaluationHelper.equals(data_type66, ModbusDataType.get_lint()):  # List
            # Array field (value)
            # Count array
            if numberOfValues66 > Integer.MAX_VALUE:
                raise ParseException(
                    "Array count of "
                    + (numberOfValues66)
                    + " exceeds the maximum allowed count of "
                    + Integer.MAX_VALUE
                )

            item_count: int = int(numberOfValues66)
            value: List[PlcValue] = []
            for cur_item in range(item_count):
                value.append(PlcLINT(c_int64(read_buffer.readLong("", 64))))

            return PlcList(value)
        if EvaluationHelper.equals(
            data_type66, ModbusDataType.get_usint()
        ) and EvaluationHelper.equals(
            number_of_values66, c_uint16(1)
        ):  # USINT
            # Reserved Field (Compartmentalized so the "reserved" variable can't leak)
            reserved: c_uint8 = read_buffer.readUnsignedShort("", 8)
            if reserved != c_uint8(0x00):
                log.info(
                    "Expected constant value "
                    + str(0x00)
                    + " but got "
                    + reserved
                    + " for reserved field."
                )

            # Simple Field (value)
            value: c_uint8 = read_buffer.readUnsignedShort("", 8)

            return PlcUSINT(value)
        if EvaluationHelper.equals(data_type66, ModbusDataType.get_usint()):  # List
            # Array field (value)
            # Count array
            if numberOfValues66 > Integer.MAX_VALUE:
                raise ParseException(
                    "Array count of "
                    + (numberOfValues66)
                    + " exceeds the maximum allowed count of "
                    + Integer.MAX_VALUE
                )

            item_count: int = int(numberOfValues66)
            value: List[PlcValue] = []
            for cur_item in range(item_count):
                value.append(PlcUINT(c_uint8(read_buffer.readUnsignedShort("", 8))))

            return PlcList(value)
        if EvaluationHelper.equals(
            data_type66, ModbusDataType.get_uint()
        ) and EvaluationHelper.equals(
            number_of_values66, c_uint16(1)
        ):  # UINT
            # Simple Field (value)
            value: c_uint16 = read_buffer.readUnsignedInt("", 16)

            return PlcUINT(value)
        if EvaluationHelper.equals(data_type66, ModbusDataType.get_uint()):  # List
            # Array field (value)
            # Count array
            if numberOfValues66 > Integer.MAX_VALUE:
                raise ParseException(
                    "Array count of "
                    + (numberOfValues66)
                    + " exceeds the maximum allowed count of "
                    + Integer.MAX_VALUE
                )

            item_count: int = int(numberOfValues66)
            value: List[PlcValue] = []
            for cur_item in range(item_count):
                value.append(PlcUDINT(c_uint16(read_buffer.readUnsignedInt("", 16))))

            return PlcList(value)
        if EvaluationHelper.equals(
            data_type66, ModbusDataType.get_udint()
        ) and EvaluationHelper.equals(
            number_of_values66, c_uint16(1)
        ):  # UDINT
            # Simple Field (value)
            value: c_uint32 = read_buffer.readUnsignedLong("", 32)

            return PlcUDINT(value)
        if EvaluationHelper.equals(data_type66, ModbusDataType.get_udint()):  # List
            # Array field (value)
            # Count array
            if numberOfValues66 > Integer.MAX_VALUE:
                raise ParseException(
                    "Array count of "
                    + (numberOfValues66)
                    + " exceeds the maximum allowed count of "
                    + Integer.MAX_VALUE
                )

            item_count: int = int(numberOfValues66)
            value: List[PlcValue] = []
            for cur_item in range(item_count):
                value.append(PlcULINT(c_uint32(read_buffer.readUnsignedLong("", 32))))

            return PlcList(value)
        if EvaluationHelper.equals(
            data_type66, ModbusDataType.get_ulint()
        ) and EvaluationHelper.equals(
            number_of_values66, c_uint16(1)
        ):  # ULINT
            # Simple Field (value)
            value: c_uint64 = read_buffer.readUnsignedBigInteger("", 64)

            return PlcULINT(value)
        if EvaluationHelper.equals(data_type66, ModbusDataType.get_ulint()):  # List
            # Array field (value)
            # Count array
            if numberOfValues66 > Integer.MAX_VALUE:
                raise ParseException(
                    "Array count of "
                    + (numberOfValues66)
                    + " exceeds the maximum allowed count of "
                    + Integer.MAX_VALUE
                )

            item_count: int = int(numberOfValues66)
            value: List[PlcValue] = []
            for cur_item in range(item_count):
                value.append(
                    PlcLINT(c_uint64(read_buffer.readUnsignedBigInteger("", 64)))
                )

            return PlcList(value)
        if EvaluationHelper.equals(
            data_type66, ModbusDataType.get_real()
        ) and EvaluationHelper.equals(
            number_of_values66, c_uint16(1)
        ):  # REAL
            # Simple Field (value)
            value: c_float = read_buffer.readFloat("", 32)

            return PlcREAL(value)
        if EvaluationHelper.equals(data_type66, ModbusDataType.get_real()):  # List
            # Array field (value)
            # Count array
            if numberOfValues66 > Integer.MAX_VALUE:
                raise ParseException(
                    "Array count of "
                    + (numberOfValues66)
                    + " exceeds the maximum allowed count of "
                    + Integer.MAX_VALUE
                )

            item_count: int = int(numberOfValues66)
            value: List[PlcValue] = []
            for cur_item in range(item_count):
                value.append(PlcREAL(c_float(read_buffer.readFloat("", 32))))

            return PlcList(value)
        if EvaluationHelper.equals(
            data_type66, ModbusDataType.get_lreal()
        ) and EvaluationHelper.equals(
            number_of_values66, c_uint16(1)
        ):  # LREAL
            # Simple Field (value)
            value: c_double = read_buffer.readDouble("", 64)

            return PlcLREAL(value)
        if EvaluationHelper.equals(data_type66, ModbusDataType.get_lreal()):  # List
            # Array field (value)
            # Count array
            if numberOfValues66 > Integer.MAX_VALUE:
                raise ParseException(
                    "Array count of "
                    + (numberOfValues66)
                    + " exceeds the maximum allowed count of "
                    + Integer.MAX_VALUE
                )

            item_count: int = int(numberOfValues66)
            value: List[PlcValue] = []
            for cur_item in range(item_count):
                value.append(PlcLREAL(c_double(read_buffer.readDouble("", 64))))

            return PlcList(value)
        if EvaluationHelper.equals(
            data_type66, ModbusDataType.get_char()
        ) and EvaluationHelper.equals(
            number_of_values66, c_uint16(1)
        ):  # CHAR
            # Simple Field (value)
            value: str = read_buffer.readString("", 8, "UTF-8")

            return PlcCHAR(value)
        if EvaluationHelper.equals(data_type66, ModbusDataType.get_char()):  # List
            # Array field (value)
            # Count array
            if numberOfValues66 > Integer.MAX_VALUE:
                raise ParseException(
                    "Array count of "
                    + (numberOfValues66)
                    + " exceeds the maximum allowed count of "
                    + Integer.MAX_VALUE
                )

            item_count: int = int(numberOfValues66)
            value: List[PlcValue] = []
            for cur_item in range(item_count):
                value.append(PlcSTRING(str(read_buffer.readString("", 8, "UTF-8"))))

            return PlcList(value)
        if EvaluationHelper.equals(
            data_type66, ModbusDataType.get_wchar()
        ) and EvaluationHelper.equals(
            number_of_values66, c_uint16(1)
        ):  # WCHAR
            # Simple Field (value)
            value: str = read_buffer.readString("", 16, "UTF-16")

            return PlcWCHAR(value)
        if EvaluationHelper.equals(data_type66, ModbusDataType.get_wchar()):  # List
            # Array field (value)
            # Count array
            if numberOfValues66 > Integer.MAX_VALUE:
                raise ParseException(
                    "Array count of "
                    + (numberOfValues66)
                    + " exceeds the maximum allowed count of "
                    + Integer.MAX_VALUE
                )

            item_count: int = int(numberOfValues66)
            value: List[PlcValue] = []
            for cur_item in range(item_count):
                value.append(PlcSTRING(str(read_buffer.readString("", 16, "UTF-16"))))

            return PlcList(value)
        return None

    @staticmethod
    def static_serialize(
        writeBuffer: WriteBuffer,
        _value: PlcValue,
        dataType: ModbusDataType,
        numberOfValues: c_uint16,
    ) -> None:
        static_serialize(
            writeBuffer, _value, dataType, numberOfValues, ByteOrder.BIG_ENDIAN
        )

    @staticmethod
    def static_serialize(
        writeBuffer: WriteBuffer,
        _value: PlcValue,
        dataType: ModbusDataType,
        numberOfValues: c_uint16,
        byteOrder: ByteOrder,
    ) -> None:
        if EvaluationHelper.equals(
            dataType66, ModbusDataType.get_bool()
        ) and EvaluationHelper.equals(
            numberOfValues66, c_uint16(1)
        ):  # BOOL
            # Reserved Field
            writeBuffer.WriteUint16("cuint160x0000", 15, c_uint16(0x0000))
            # Simple Field (value)
            value: c_bool = _value.getC_bool()
            writeBuffer.WriteBit("value", (value))
        if EvaluationHelper.equals(dataType66, ModbusDataType.get_bool()):  # List
            values: PlcList = _value

            for val in values.getList():
                value: c_bool = val.getC_bool()
                writeBuffer.WriteBit("value", (value))

        if EvaluationHelper.equals(
            dataType66, ModbusDataType.get_byte()
        ) and EvaluationHelper.equals(
            numberOfValues66, c_uint16(1)
        ):  # BYTE
            # Reserved Field
            writeBuffer.WriteUint8("cuint80x00", 8, c_uint8(0x00))
            # Simple Field (value)
            value: c_uint8 = _value.getC_uint8()
            writeBuffer.WriteUint8("value", 8, (value))
        if EvaluationHelper.equals(dataType66, ModbusDataType.get_byte()):  # List
            values: PlcList = _value

            for val in values.getList():
                value: c_bool = val.getC_bool()
                writeBuffer.WriteBit("value", (value))

        if EvaluationHelper.equals(dataType66, ModbusDataType.get_word()):  # WORD
            # Simple Field (value)
            value: c_uint16 = _value.getC_uint16()
            writeBuffer.WriteUint16("value", 16, (value))
        if EvaluationHelper.equals(dataType66, ModbusDataType.get_dword()):  # DWORD
            # Simple Field (value)
            value: c_uint32 = _value.getC_uint32()
            writeBuffer.WriteUint32("value", 32, (value))
        if EvaluationHelper.equals(dataType66, ModbusDataType.get_lword()):  # LWORD
            # Simple Field (value)
            value: c_uint64 = _value.getC_uint64()
            writeBuffer.WriteUint64("value", 64, (value))
        if EvaluationHelper.equals(
            dataType66, ModbusDataType.get_sint()
        ) and EvaluationHelper.equals(
            numberOfValues66, c_uint16(1)
        ):  # SINT
            # Reserved Field
            writeBuffer.WriteUint8("cuint80x00", 8, c_uint8(0x00))
            # Simple Field (value)
            value: c_int8 = _value.getC_int8()
            writeBuffer.WriteInt8("value", 8, (value))
        if EvaluationHelper.equals(dataType66, ModbusDataType.get_sint()):  # List
            values: PlcList = _value

            for val in values.getList():
                value: c_int8 = val.getC_int8()
                writeBuffer.WriteInt8("value", 8, (value))

        if EvaluationHelper.equals(
            dataType66, ModbusDataType.get_int()
        ) and EvaluationHelper.equals(
            numberOfValues66, c_uint16(1)
        ):  # INT
            # Simple Field (value)
            value: c_int16 = _value.getC_int16()
            writeBuffer.WriteInt16("value", 16, (value))
        if EvaluationHelper.equals(dataType66, ModbusDataType.get_int()):  # List
            values: PlcList = _value

            for val in values.getList():
                value: c_int16 = val.getC_int16()
                writeBuffer.WriteInt16("value", 16, (value))

        if EvaluationHelper.equals(
            dataType66, ModbusDataType.get_dint()
        ) and EvaluationHelper.equals(
            numberOfValues66, c_uint16(1)
        ):  # DINT
            # Simple Field (value)
            value: c_int32 = _value.getC_int32()
            writeBuffer.WriteInt32("value", 32, (value))
        if EvaluationHelper.equals(dataType66, ModbusDataType.get_dint()):  # List
            values: PlcList = _value

            for val in values.getList():
                value: c_int32 = val.getC_int32()
                writeBuffer.WriteInt32("value", 32, (value))

        if EvaluationHelper.equals(
            dataType66, ModbusDataType.get_lint()
        ) and EvaluationHelper.equals(
            numberOfValues66, c_uint16(1)
        ):  # LINT
            # Simple Field (value)
            value: c_int64 = _value.getC_int64()
            writeBuffer.WriteInt64("value", 64, (value))
        if EvaluationHelper.equals(dataType66, ModbusDataType.get_lint()):  # List
            values: PlcList = _value

            for val in values.getList():
                value: c_int64 = val.getC_int64()
                writeBuffer.WriteInt64("value", 64, (value))

        if EvaluationHelper.equals(
            dataType66, ModbusDataType.get_usint()
        ) and EvaluationHelper.equals(
            numberOfValues66, c_uint16(1)
        ):  # USINT
            # Reserved Field
            writeBuffer.WriteUint8("cuint80x00", 8, c_uint8(0x00))
            # Simple Field (value)
            value: c_uint8 = _value.getC_uint8()
            writeBuffer.WriteUint8("value", 8, (value))
        if EvaluationHelper.equals(dataType66, ModbusDataType.get_usint()):  # List
            values: PlcList = _value

            for val in values.getList():
                value: c_uint8 = val.getC_uint8()
                writeBuffer.WriteUint8("value", 8, (value))

        if EvaluationHelper.equals(
            dataType66, ModbusDataType.get_uint()
        ) and EvaluationHelper.equals(
            numberOfValues66, c_uint16(1)
        ):  # UINT
            # Simple Field (value)
            value: c_uint16 = _value.getC_uint16()
            writeBuffer.WriteUint16("value", 16, (value))
        if EvaluationHelper.equals(dataType66, ModbusDataType.get_uint()):  # List
            values: PlcList = _value

            for val in values.getList():
                value: c_uint16 = val.getC_uint16()
                writeBuffer.WriteUint16("value", 16, (value))

        if EvaluationHelper.equals(
            dataType66, ModbusDataType.get_udint()
        ) and EvaluationHelper.equals(
            numberOfValues66, c_uint16(1)
        ):  # UDINT
            # Simple Field (value)
            value: c_uint32 = _value.getC_uint32()
            writeBuffer.WriteUint32("value", 32, (value))
        if EvaluationHelper.equals(dataType66, ModbusDataType.get_udint()):  # List
            values: PlcList = _value

            for val in values.getList():
                value: c_uint32 = val.getC_uint32()
                writeBuffer.WriteUint32("value", 32, (value))

        if EvaluationHelper.equals(
            dataType66, ModbusDataType.get_ulint()
        ) and EvaluationHelper.equals(
            numberOfValues66, c_uint16(1)
        ):  # ULINT
            # Simple Field (value)
            value: c_uint64 = _value.getC_uint64()
            writeBuffer.WriteUint64("value", 64, (value))
        if EvaluationHelper.equals(dataType66, ModbusDataType.get_ulint()):  # List
            values: PlcList = _value

            for val in values.getList():
                value: c_uint64 = val.getC_uint64()
                writeBuffer.WriteUint64("value", 64, (value))

        if EvaluationHelper.equals(
            dataType66, ModbusDataType.get_real()
        ) and EvaluationHelper.equals(
            numberOfValues66, c_uint16(1)
        ):  # REAL
            # Simple Field (value)
            value: c_float = _value.getC_float()
            writeBuffer.WriteFloat32("value", 32, (value))
        if EvaluationHelper.equals(dataType66, ModbusDataType.get_real()):  # List
            values: PlcList = _value

            for val in values.getList():
                value: c_float = val.getC_float()
                writeBuffer.WriteFloat32("value", 32, (value))

        if EvaluationHelper.equals(
            dataType66, ModbusDataType.get_lreal()
        ) and EvaluationHelper.equals(
            numberOfValues66, c_uint16(1)
        ):  # LREAL
            # Simple Field (value)
            value: c_double = _value.getC_double()
            writeBuffer.WriteFloat64("value", 64, (value))
        if EvaluationHelper.equals(dataType66, ModbusDataType.get_lreal()):  # List
            values: PlcList = _value

            for val in values.getList():
                value: c_double = val.getC_double()
                writeBuffer.WriteFloat64("value", 64, (value))

        if EvaluationHelper.equals(
            dataType66, ModbusDataType.get_char()
        ) and EvaluationHelper.equals(
            numberOfValues66, c_uint16(1)
        ):  # CHAR
            # Simple Field (value)
            value: str = _value.getStr()
            writeBuffer.WriteString("value", uint32(8), "UTF-8", (value))
        if EvaluationHelper.equals(dataType66, ModbusDataType.get_char()):  # List
            values: PlcList = _value

            for val in values.getList():
                value: str = val.getStr()
                writeBuffer.WriteString("value", uint32(8), "UTF-8", (value))

        if EvaluationHelper.equals(
            dataType66, ModbusDataType.get_wchar()
        ) and EvaluationHelper.equals(
            numberOfValues66, c_uint16(1)
        ):  # WCHAR
            # Simple Field (value)
            value: str = _value.getStr()
            writeBuffer.WriteString("value", uint32(16), "UTF-16", (value))
        if EvaluationHelper.equals(dataType66, ModbusDataType.get_wchar()):  # List
            values: PlcList = _value

            for val in values.getList():
                value: str = val.getStr()
                writeBuffer.WriteString("value", uint32(16), "UTF-16", (value))

    @staticmethod
    def get_length_in_bytes(
        _value: PlcValue, dataType: ModbusDataType, numberOfValues: c_uint16
    ) -> int:
        return int(
            math.ceil(float(getLengthInBits(_value, dataType, numberOfValues)) / 8.0)
        )

    @staticmethod
    def get_length_in_bits(
        _value: PlcValue, dataType: ModbusDataType, numberOfValues: c_uint16
    ) -> int:
        sizeInBits: int = 0
        if EvaluationHelper.equals(
            dataType66, ModbusDataType.get_bool()
        ) and EvaluationHelper.equals(
            numberOfValues66, c_uint16(1)
        ):  # BOOL
            # Reserved Field
            sizeInBits += 15
            # Simple Field (value)
            sizeInBits += 1
        if EvaluationHelper.equals(dataType66, ModbusDataType.get_bool()):  # List
            values: PlcList = _value
            sizeInBits += values.getList().size() * 1
        if EvaluationHelper.equals(
            dataType66, ModbusDataType.get_byte()
        ) and EvaluationHelper.equals(
            numberOfValues66, c_uint16(1)
        ):  # BYTE
            # Reserved Field
            sizeInBits += 8
            # Simple Field (value)
            sizeInBits += 8
        if EvaluationHelper.equals(dataType66, ModbusDataType.get_byte()):  # List
            values: PlcList = _value
            sizeInBits += values.getList().size() * 1
        if EvaluationHelper.equals(dataType66, ModbusDataType.get_word()):  # WORD
            # Simple Field (value)
            sizeInBits += 16
        if EvaluationHelper.equals(dataType66, ModbusDataType.get_dword()):  # DWORD
            # Simple Field (value)
            sizeInBits += 32
        if EvaluationHelper.equals(dataType66, ModbusDataType.get_lword()):  # LWORD
            # Simple Field (value)
            sizeInBits += 64
        if EvaluationHelper.equals(
            dataType66, ModbusDataType.get_sint()
        ) and EvaluationHelper.equals(
            numberOfValues66, c_uint16(1)
        ):  # SINT
            # Reserved Field
            sizeInBits += 8
            # Simple Field (value)
            sizeInBits += 8
        if EvaluationHelper.equals(dataType66, ModbusDataType.get_sint()):  # List
            values: PlcList = _value
            sizeInBits += values.getList().size() * 8
        if EvaluationHelper.equals(
            dataType66, ModbusDataType.get_int()
        ) and EvaluationHelper.equals(
            numberOfValues66, c_uint16(1)
        ):  # INT
            # Simple Field (value)
            sizeInBits += 16
        if EvaluationHelper.equals(dataType66, ModbusDataType.get_int()):  # List
            values: PlcList = _value
            sizeInBits += values.getList().size() * 16
        if EvaluationHelper.equals(
            dataType66, ModbusDataType.get_dint()
        ) and EvaluationHelper.equals(
            numberOfValues66, c_uint16(1)
        ):  # DINT
            # Simple Field (value)
            sizeInBits += 32
        if EvaluationHelper.equals(dataType66, ModbusDataType.get_dint()):  # List
            values: PlcList = _value
            sizeInBits += values.getList().size() * 32
        if EvaluationHelper.equals(
            dataType66, ModbusDataType.get_lint()
        ) and EvaluationHelper.equals(
            numberOfValues66, c_uint16(1)
        ):  # LINT
            # Simple Field (value)
            sizeInBits += 64
        if EvaluationHelper.equals(dataType66, ModbusDataType.get_lint()):  # List
            values: PlcList = _value
            sizeInBits += values.getList().size() * 64
        if EvaluationHelper.equals(
            dataType66, ModbusDataType.get_usint()
        ) and EvaluationHelper.equals(
            numberOfValues66, c_uint16(1)
        ):  # USINT
            # Reserved Field
            sizeInBits += 8
            # Simple Field (value)
            sizeInBits += 8
        if EvaluationHelper.equals(dataType66, ModbusDataType.get_usint()):  # List
            values: PlcList = _value
            sizeInBits += values.getList().size() * 8
        if EvaluationHelper.equals(
            dataType66, ModbusDataType.get_uint()
        ) and EvaluationHelper.equals(
            numberOfValues66, c_uint16(1)
        ):  # UINT
            # Simple Field (value)
            sizeInBits += 16
        if EvaluationHelper.equals(dataType66, ModbusDataType.get_uint()):  # List
            values: PlcList = _value
            sizeInBits += values.getList().size() * 16
        if EvaluationHelper.equals(
            dataType66, ModbusDataType.get_udint()
        ) and EvaluationHelper.equals(
            numberOfValues66, c_uint16(1)
        ):  # UDINT
            # Simple Field (value)
            sizeInBits += 32
        if EvaluationHelper.equals(dataType66, ModbusDataType.get_udint()):  # List
            values: PlcList = _value
            sizeInBits += values.getList().size() * 32
        if EvaluationHelper.equals(
            dataType66, ModbusDataType.get_ulint()
        ) and EvaluationHelper.equals(
            numberOfValues66, c_uint16(1)
        ):  # ULINT
            # Simple Field (value)
            sizeInBits += 64
        if EvaluationHelper.equals(dataType66, ModbusDataType.get_ulint()):  # List
            values: PlcList = _value
            sizeInBits += values.getList().size() * 64
        if EvaluationHelper.equals(
            dataType66, ModbusDataType.get_real()
        ) and EvaluationHelper.equals(
            numberOfValues66, c_uint16(1)
        ):  # REAL
            # Simple Field (value)
            sizeInBits += 32
        if EvaluationHelper.equals(dataType66, ModbusDataType.get_real()):  # List
            values: PlcList = _value
            sizeInBits += values.getList().size() * 32
        if EvaluationHelper.equals(
            dataType66, ModbusDataType.get_lreal()
        ) and EvaluationHelper.equals(
            numberOfValues66, c_uint16(1)
        ):  # LREAL
            # Simple Field (value)
            sizeInBits += 64
        if EvaluationHelper.equals(dataType66, ModbusDataType.get_lreal()):  # List
            values: PlcList = _value
            sizeInBits += values.getList().size() * 64
        if EvaluationHelper.equals(
            dataType66, ModbusDataType.get_char()
        ) and EvaluationHelper.equals(
            numberOfValues66, c_uint16(1)
        ):  # CHAR
            # Simple Field (value)
            sizeInBits += 8
        if EvaluationHelper.equals(dataType66, ModbusDataType.get_char()):  # List
            values: PlcList = _value
            sizeInBits += values.getList().size() * 8
        if EvaluationHelper.equals(
            dataType66, ModbusDataType.get_wchar()
        ) and EvaluationHelper.equals(
            numberOfValues66, c_uint16(1)
        ):  # WCHAR
            # Simple Field (value)
            sizeInBits += 16
        if EvaluationHelper.equals(dataType66, ModbusDataType.get_wchar()):  # List
            values: PlcList = _value
            sizeInBits += values.getList().size() * 16
        return sizeInBits
