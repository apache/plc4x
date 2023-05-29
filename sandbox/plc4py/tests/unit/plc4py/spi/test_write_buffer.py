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

import pytest
from bitarray import bitarray

from plc4py.spi.generation.WriteBuffer import WriteBufferByteBased

from plc4py.utils.GenericTypes import ByteOrder


def test_write_buffer_get_bytes(mocker) -> None:
    wb: WriteBufferByteBased = WriteBufferByteBased(10, ByteOrder.LITTLE_ENDIAN)
    ba: memoryview = wb.get_bytes()
    assert 10 == ba.nbytes


def test_write_buffer_get_pos(mocker) -> None:
    wb: WriteBufferByteBased = WriteBufferByteBased(10, ByteOrder.LITTLE_ENDIAN)
    pos: int = wb.get_pos()
    assert 0 == pos


def test_write_buffer_set_bit_little_endian(mocker) -> None:
    wb: WriteBufferByteBased = WriteBufferByteBased(1, ByteOrder.LITTLE_ENDIAN)
    wb.write_bit(True)
    ba: memoryview = wb.get_bytes()
    assert b"\x01" == ba.tobytes()


def test_write_buffer_set_bit_big_endian(mocker) -> None:
    wb: WriteBufferByteBased = WriteBufferByteBased(1, ByteOrder.BIG_ENDIAN)
    wb.write_bit(True)
    ba: memoryview = wb.get_bytes()
    assert b"\x80" == ba.tobytes()


def test_write_buffer_set_bit_get_pos(mocker) -> None:
    wb: WriteBufferByteBased = WriteBufferByteBased(1, ByteOrder.LITTLE_ENDIAN)
    wb.write_bit(True)
    pos: int = wb.get_pos()
    assert 1 == pos


def test_write_buffer_set_bit_x2_little_endian(mocker) -> None:
    wb: WriteBufferByteBased = WriteBufferByteBased(1, ByteOrder.LITTLE_ENDIAN)
    wb.write_bit(True)
    wb.write_bit(True)
    ba: memoryview = wb.get_bytes()
    assert b"\x03" == ba.tobytes()


def test_write_buffer_set_bit_x2_big_endian(mocker) -> None:
    wb: WriteBufferByteBased = WriteBufferByteBased(1, ByteOrder.BIG_ENDIAN)
    wb.write_bit(True)
    wb.write_bit(True)
    ba: memoryview = wb.get_bytes()
    assert b"\xc0" == ba.tobytes()


def test_write_buffer_set_bit_get_pos_dual(mocker) -> None:
    wb: WriteBufferByteBased = WriteBufferByteBased(1, ByteOrder.LITTLE_ENDIAN)
    wb.write_bit(True)
    wb.write_bit(True)
    pos: int = wb.get_pos()
    assert pos == 2


def test_write_buffer_set_unsigned_byte_get_pos(mocker) -> None:
    wb: WriteBufferByteBased = WriteBufferByteBased(1, ByteOrder.LITTLE_ENDIAN)
    wb.write_unsigned_byte(0xFF, 8)
    pos: int = wb.get_pos()
    assert pos == 8


def test_write_buffer_set_unsigned_byte_little_endian(mocker) -> None:
    wb: WriteBufferByteBased = WriteBufferByteBased(1, ByteOrder.LITTLE_ENDIAN)
    wb.write_unsigned_byte(0x12, 8)
    ba: memoryview = wb.get_bytes()
    assert ba.obj == bitarray("01001000")


def test_write_buffer_set_unsigned_byte_big_endian(mocker) -> None:
    wb: WriteBufferByteBased = WriteBufferByteBased(1, ByteOrder.BIG_ENDIAN)
    wb.write_unsigned_byte(0x12, 8)
    ba: memoryview = wb.get_bytes()
    assert ba.obj == bitarray("00010010")


def test_write_buffer_set_unsigned_byte_little_endian_niblet(mocker) -> None:
    wb: WriteBufferByteBased = WriteBufferByteBased(1, ByteOrder.LITTLE_ENDIAN)
    wb.write_unsigned_byte(0x12, 4)
    ba: memoryview = wb.get_bytes()
    assert ba.obj == bitarray("01000000")


def test_write_buffer_set_unsigned_byte_big_endian_niblet(mocker) -> None:
    wb: WriteBufferByteBased = WriteBufferByteBased(1, ByteOrder.BIG_ENDIAN)
    wb.write_unsigned_byte(0x12, 4)
    ba: memoryview = wb.get_bytes()
    assert ba.obj == bitarray("00010000")


def test_write_buffer_write_unsigned_short_little_endian(mocker) -> None:
    wb: WriteBufferByteBased = WriteBufferByteBased(2, ByteOrder.LITTLE_ENDIAN)
    wb.write_unsigned_short(0x12, 16)
    ba: memoryview = wb.get_bytes()
    assert ba.obj == bitarray("01001000 00000000", endian="little")


def test_write_buffer_write_unsigned_short_big_endian(mocker) -> None:
    wb: WriteBufferByteBased = WriteBufferByteBased(2, ByteOrder.BIG_ENDIAN)
    wb.write_unsigned_short(0x12, 16)
    ba: memoryview = wb.get_bytes()
    assert ba.obj == bitarray("00010010 00000000", endian="big")


def test_write_buffer_write_unsigned_short_little_endian_dual(mocker) -> None:
    wb: WriteBufferByteBased = WriteBufferByteBased(2, ByteOrder.LITTLE_ENDIAN)
    wb.write_unsigned_short(0x12, 16)
    wb.write_unsigned_short(0x34, 16)
    ba: memoryview = wb.get_bytes()
    assert ba.obj == bitarray("01001000 00000000 00101100 00000000", endian="little")


def test_write_buffer_write_unsigned_short_big_endian_dual(mocker) -> None:
    wb: WriteBufferByteBased = WriteBufferByteBased(2, ByteOrder.BIG_ENDIAN)
    wb.write_unsigned_short(0x12, 16)
    wb.write_unsigned_short(0x34, 16)
    ba: memoryview = wb.get_bytes()
    assert ba.obj == bitarray("00010010 00000000 00110100 00000000", endian="big")


def test_write_buffer_write_unsigned_short_big_endian_full(mocker) -> None:
    wb: WriteBufferByteBased = WriteBufferByteBased(2, ByteOrder.BIG_ENDIAN)
    wb.write_unsigned_short(-1, 16)
    ba: memoryview = wb.get_bytes()
    assert ba.obj == bitarray("11111111 11111111", endian="big")


def test_write_buffer_write_unsigned_short_bit_big_endian_full(mocker) -> None:
    wb: WriteBufferByteBased = WriteBufferByteBased(2, ByteOrder.BIG_ENDIAN)
    wb.write_bit(True)
    wb.write_bit(False)
    wb.write_unsigned_short(-1, 16)
    ba: memoryview = wb.get_bytes()
    assert ba.obj == bitarray("10 11111111 11111111", endian="big")


def test_write_buffer_write_unsigned_short_ascii_encoding_little_endian(mocker) -> None:
    wb: WriteBufferByteBased = WriteBufferByteBased(2, ByteOrder.LITTLE_ENDIAN)
    wb.write_unsigned_short(1, 16, "ASCII Value of 1 - 0x31", encoding="ASCII")
    ba: memoryview = wb.get_bytes()
    assert ba.obj == bitarray("10001100", endian="little")


def test_write_buffer_write_unsigned_short_ascii_encoding_big_endian(mocker) -> None:
    wb: WriteBufferByteBased = WriteBufferByteBased(2, ByteOrder.BIG_ENDIAN)
    wb.write_unsigned_short(1, 16, "ASCII Value of 1 - 0x31", encoding="ASCII")
    ba: memoryview = wb.get_bytes()
    assert ba.obj == bitarray("00110001", endian="big")


def test_write_buffer_write_unsigned_int_little_endian(mocker) -> None:
    wb: WriteBufferByteBased = WriteBufferByteBased(2, ByteOrder.LITTLE_ENDIAN)
    # Testing incompatible format being sent to it.
    wb.write_unsigned_int(0x12, 32)
    ba: memoryview = wb.get_bytes()
    assert ba.obj == bitarray("01001000 00000000 00000000 00000000", endian="little")


def test_write_buffer_write_unsigned_int_big_endian(mocker) -> None:
    wb: WriteBufferByteBased = WriteBufferByteBased(2, ByteOrder.BIG_ENDIAN)
    # Testing incompatible format being sent to it.
    wb.write_unsigned_int(0x12, 32)
    ba: memoryview = wb.get_bytes()
    assert ba.obj == bitarray("00010010 00000000 00000000 00000000", endian="big")


def test_write_buffer_write_unsigned_long_little_endian(mocker) -> None:
    wb: WriteBufferByteBased = WriteBufferByteBased(2, ByteOrder.LITTLE_ENDIAN)
    wb.write_unsigned_long(0x12, 64)
    ba: memoryview = wb.get_bytes()
    assert ba.obj == bitarray(
        "01001000 00000000 00000000 00000000 00000000 00000000 00000000 00000000",
        endian="little",
    )


def test_write_buffer_write_unsigned_long_big_endian(mocker) -> None:
    wb: WriteBufferByteBased = WriteBufferByteBased(2, ByteOrder.BIG_ENDIAN)
    wb.write_unsigned_long(0x12, 64)
    ba: memoryview = wb.get_bytes()
    assert ba.obj == bitarray(
        "00010010 00000000 00000000 00000000 00000000 00000000 00000000 00000000",
        endian="big",
    )


def test_write_buffer_write_unsigned_long_ascii_encoding_little_endian(mocker) -> None:
    wb: WriteBufferByteBased = WriteBufferByteBased(8, ByteOrder.LITTLE_ENDIAN)
    wb.write_unsigned_long(
        11111111,
        64,
        "ASCII Value of 1111 1111 - 0x3131313131313131",
        encoding="ASCII",
    )
    ba: memoryview = wb.get_bytes()
    assert ba.obj == bitarray(
        "10001100 10001100 10001100 10001100 10001100 10001100 10001100 10001100",
        endian="little",
    )


def test_write_buffer_write_unsigned_long_ascii_encoding_big_endian(mocker) -> None:
    wb: WriteBufferByteBased = WriteBufferByteBased(8, ByteOrder.BIG_ENDIAN)
    wb.write_unsigned_long(
        11111111,
        64,
        "ASCII Value of 1111 1111 - 0x3131313131313131",
        encoding="ASCII",
    )
    ba: memoryview = wb.get_bytes()
    assert ba.obj == bitarray(
        "00110001 00110001 00110001 00110001 00110001 00110001 00110001 00110001",
        endian="big",
    )


def test_write_buffer_set_signed_byte(mocker) -> None:
    wb: WriteBufferByteBased = WriteBufferByteBased(1, ByteOrder.LITTLE_ENDIAN)
    wb.write_signed_byte(-1, 8)
    ba: memoryview = wb.get_bytes()
    assert ba.obj == bitarray("11111111", endian="little")


def test_write_buffer_set_signed_byte_three(mocker) -> None:
    wb: WriteBufferByteBased = WriteBufferByteBased(1, ByteOrder.LITTLE_ENDIAN)
    wb.write_signed_byte(3, 8)
    ba: memoryview = wb.get_bytes()
    assert ba.obj == bitarray("11000000", endian="little")


def test_write_buffer_set_float_little_endian(mocker) -> None:
    wb: WriteBufferByteBased = WriteBufferByteBased(1, ByteOrder.LITTLE_ENDIAN)
    wb.write_float(-1, 32)
    ba: memoryview = wb.get_bytes()
    assert ba.obj == bitarray("00000000000000000000000111111101", endian="little")


def test_write_buffer_set_float_big_endian(mocker) -> None:
    wb: WriteBufferByteBased = WriteBufferByteBased(1, ByteOrder.BIG_ENDIAN)
    wb.write_float(-1, 32)
    ba: memoryview = wb.get_bytes()
    assert ba.obj == bitarray("00000000000000001000000010111111", endian="big")
