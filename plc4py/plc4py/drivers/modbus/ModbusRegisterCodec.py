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
"""
Reads and writes values as they are laid out in Modbus registers.

``DataItem`` encodes one value at its natural width and knows nothing about registers. What it does
not cover is that a register is 16 bits wide: a lone value narrower than that occupies a whole
register and is padded, while several of them are packed. Which half of the register a padded value
sits in depends on the byte order. That is Modbus knowledge, so it lives here rather than in the
protocol description.
"""

from plc4py.api.value.PlcValue import PlcValue
from plc4py.protocols.modbus.readwrite.DataItem import DataItem
from plc4py.protocols.modbus.readwrite.ModbusDataType import ModbusDataType
from plc4py.spi.values.PlcValues import PlcList

# Width of one value of each type, in bits. Strings are sized from their declared length instead.
_WIDTH_BITS = {
    ModbusDataType.BOOL: 1,
    ModbusDataType.BYTE: 8,
    ModbusDataType.SINT: 8,
    ModbusDataType.USINT: 8,
    ModbusDataType.CHAR: 8,
    ModbusDataType.WORD: 16,
    ModbusDataType.INT: 16,
    ModbusDataType.UINT: 16,
    ModbusDataType.WCHAR: 16,
    ModbusDataType.DWORD: 32,
    ModbusDataType.DINT: 32,
    ModbusDataType.UDINT: 32,
    ModbusDataType.REAL: 32,
    ModbusDataType.LWORD: 64,
    ModbusDataType.LINT: 64,
    ModbusDataType.ULINT: 64,
    ModbusDataType.LREAL: 64,
}

# What a lone value needs on top of its own width to fill a register. CHAR is deliberately left
# unpadded, as it always has been.
_PADDING_BITS = {
    ModbusDataType.BOOL: 15,
    ModbusDataType.BYTE: 8,
    ModbusDataType.SINT: 8,
    ModbusDataType.USINT: 8,
}


def width_bits(data_type: ModbusDataType, string_length: int) -> int:
    if data_type == ModbusDataType.STRING:
        return string_length * 8
    if data_type == ModbusDataType.WSTRING:
        return string_length * 16
    return _WIDTH_BITS.get(data_type, 0)


def padding_bits(data_type: ModbusDataType) -> int:
    return _PADDING_BITS.get(data_type, 0)


def _leading_padding_bits_little_endian(data_type: ModbusDataType) -> int:
    """A lone BOOL sits between seven bits and eight; a byte simply comes first."""
    return 7 if data_type == ModbusDataType.BOOL else 0


def trailing_padding_bits(data_type: ModbusDataType, number_of_values: int) -> int:
    """The bits needed to round a packed run of values up to a whole register."""
    width = width_bits(data_type, 1)
    if width >= 16:
        return 0
    remainder = (width * number_of_values) % 16
    return 0 if remainder == 0 else 16 - remainder


def length_in_bytes(
    data_type: ModbusDataType, number_of_values: int, string_length: int
) -> int:
    if number_of_values == 1:
        return (width_bits(data_type, string_length) + padding_bits(data_type) + 7) // 8
    bits = width_bits(data_type, string_length) * number_of_values
    bits += trailing_padding_bits(data_type, number_of_values)
    return (bits + 7) // 8


def parse(
    read_buffer,
    data_type: ModbusDataType,
    number_of_values: int,
    big_endian: bool,
    string_length: int,
) -> PlcValue:
    """Reads number_of_values values, returning the value itself for one and a list for several."""
    if number_of_values == 1:
        padding = padding_bits(data_type)
        if padding == 0:
            return DataItem.static_parse(read_buffer, data_type, string_length)
        if big_endian:
            read_buffer.read_unsigned_short(padding, logical_name="padding")
            return DataItem.static_parse(read_buffer, data_type, string_length)
        leading = _leading_padding_bits_little_endian(data_type)
        if leading > 0:
            read_buffer.read_unsigned_short(leading, logical_name="padding")
        value = DataItem.static_parse(read_buffer, data_type, string_length)
        read_buffer.read_unsigned_short(padding - leading, logical_name="padding")
        return value

    # Several values are packed without padding between them; a trailing pad is left unread.
    return PlcList(
        [
            DataItem.static_parse(read_buffer, data_type, string_length)
            for _ in range(number_of_values)
        ]
    )


def serialize(
    write_buffer,
    value: PlcValue,
    data_type: ModbusDataType,
    number_of_values: int,
    big_endian: bool,
    string_length: int,
    byte_order=None,
) -> None:
    """Writes a value, or every element of a list, with the layout parse expects."""
    if number_of_values == 1:
        padding = padding_bits(data_type)
        if padding == 0:
            DataItem.static_serialize(
                write_buffer, value, data_type, string_length, byte_order
            )
            return
        if big_endian:
            write_buffer.write_unsigned_short(0, padding, logical_name="padding")
            DataItem.static_serialize(
                write_buffer, value, data_type, string_length, byte_order
            )
            return
        leading = _leading_padding_bits_little_endian(data_type)
        if leading > 0:
            write_buffer.write_unsigned_short(0, leading, logical_name="padding")
        DataItem.static_serialize(
            write_buffer, value, data_type, string_length, byte_order
        )
        write_buffer.write_unsigned_short(0, padding - leading, logical_name="padding")
        return

    for index in range(number_of_values):
        element = value.get_list()[index] if isinstance(value, PlcList) else value
        DataItem.static_serialize(
            write_buffer, element, data_type, string_length, byte_order
        )
    trailing = trailing_padding_bits(data_type, number_of_values)
    if trailing > 0:
        write_buffer.write_unsigned_short(0, trailing, logical_name="padding")
