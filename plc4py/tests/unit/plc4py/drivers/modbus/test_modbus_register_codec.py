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
A Modbus register is 16 bits wide, so a value narrower than that is padded when it stands alone and
packed when it does not. Those rules used to be separate cases in the protocol description and now
live in ModbusRegisterCodec, which is why they are pinned down here directly.
"""

from plc4py.drivers.modbus import ModbusRegisterCodec
from plc4py.protocols.modbus.readwrite.ModbusDataType import ModbusDataType
from plc4py.spi.generation.ReadBuffer import ReadBufferByteBased
from plc4py.spi.generation.WriteBuffer import WriteBufferByteBased
from plc4py.spi.values.PlcValues import PlcBOOL, PlcList, PlcSINT, PlcUSINT
from plc4py.utils.GenericTypes import ByteOrder


def _parse(data, data_type, number_of_values, big_endian=True):
    read_buffer = ReadBufferByteBased(bytearray(data), ByteOrder.BIG_ENDIAN)
    return ModbusRegisterCodec.parse(
        read_buffer, data_type, number_of_values, big_endian, 1
    )


def _serialize(value, data_type, number_of_values, big_endian=True):
    length = ModbusRegisterCodec.length_in_bytes(data_type, number_of_values, 1)
    write_buffer = WriteBufferByteBased(length, ByteOrder.BIG_ENDIAN)
    ModbusRegisterCodec.serialize(
        write_buffer,
        value,
        data_type,
        number_of_values,
        big_endian,
        1,
        ByteOrder.BIG_ENDIAN,
    )
    return list(write_buffer.get_bytes().tobytes())


def test_pads_a_lone_byte_at_the_front_when_big_endian():
    assert _parse([0x00, 0x2A], ModbusDataType.USINT, 1).value == 0x2A
    assert _serialize(PlcUSINT(0x2A), ModbusDataType.USINT, 1) == [0x00, 0x2A]


def test_pads_a_lone_byte_at_the_end_when_little_endian():
    assert _parse([0x2A, 0x00], ModbusDataType.USINT, 1, big_endian=False).value == 0x2A
    assert _serialize(PlcUSINT(0x2A), ModbusDataType.USINT, 1, big_endian=False) == [
        0x2A,
        0x00,
    ]


def test_pads_a_lone_signed_byte_and_keeps_its_sign():
    assert _parse([0x00, 0xFF], ModbusDataType.SINT, 1).value == -1
    assert _parse([0xFF, 0x00], ModbusDataType.SINT, 1, big_endian=False).value == -1


def test_pads_a_lone_bool_to_a_whole_register():
    """Big endian leaves the bit last; little endian puts it after seven bits."""
    assert _parse([0x00, 0x01], ModbusDataType.BOOL, 1).value is True
    assert _parse([0x00, 0x00], ModbusDataType.BOOL, 1).value is False
    assert _parse([0x01, 0x00], ModbusDataType.BOOL, 1, big_endian=False).value is True


def test_leaves_a_lone_char_unpadded():
    """CHAR is not padded - it never has been, unlike the other single byte types."""
    assert ModbusRegisterCodec.length_in_bytes(ModbusDataType.CHAR, 1, 1) == 1


def test_never_pads_a_register_wide_value():
    assert _parse([0x12, 0x34], ModbusDataType.UINT, 1).value == 0x1234
    assert ModbusRegisterCodec.length_in_bytes(ModbusDataType.UINT, 1, 1) == 2


def test_packs_several_bytes_without_padding():
    value = _parse([0x01, 0x02, 0x03, 0x04], ModbusDataType.USINT, 4)

    assert len(value.value) == 4
    assert value.value[0].value == 1
    assert value.value[3].value == 4


def test_packs_bools_as_bits():
    value = _parse([0b10100000, 0x00], ModbusDataType.BOOL, 3)

    assert [item.value for item in value.value] == [True, False, True]


def test_pads_the_last_register_for_an_odd_count():
    three = PlcList([PlcUSINT(1), PlcUSINT(2), PlcUSINT(3)])

    assert ModbusRegisterCodec.length_in_bytes(ModbusDataType.USINT, 3, 1) == 4
    assert _serialize(three, ModbusDataType.USINT, 3) == [0x01, 0x02, 0x03, 0x00]


def test_adds_no_trailing_pad_for_an_even_count():
    two = PlcList([PlcUSINT(1), PlcUSINT(2)])

    assert ModbusRegisterCodec.length_in_bytes(ModbusDataType.USINT, 2, 1) == 2
    assert _serialize(two, ModbusDataType.USINT, 2) == [0x01, 0x02]


def test_pads_packed_bools_to_the_register_boundary():
    three_bools = PlcList([PlcBOOL(True), PlcBOOL(False), PlcBOOL(True)])

    assert ModbusRegisterCodec.length_in_bytes(ModbusDataType.BOOL, 3, 1) == 2
    assert _serialize(three_bools, ModbusDataType.BOOL, 3) == [0b10100000, 0x00]


def test_round_trips_a_padded_value():
    for big_endian in (True, False):
        written = _serialize(PlcSINT(-42), ModbusDataType.SINT, 1, big_endian)

        assert len(written) == 2, "a lone byte occupies a whole register"
        assert _parse(written, ModbusDataType.SINT, 1, big_endian).value == -42


def test_round_trips_a_packed_run():
    values = PlcList([PlcSINT(-1), PlcSINT(2), PlcSINT(-3)])
    written = _serialize(values, ModbusDataType.SINT, 3)

    read = _parse(written, ModbusDataType.SINT, 3)
    assert [item.value for item in read.value] == [-1, 2, -3]
