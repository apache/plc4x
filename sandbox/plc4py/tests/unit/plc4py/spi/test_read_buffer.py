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

from plc4py.spi.generation.ReadBuffer import ReadBufferByteBased
from plc4py.spi.generation.WriteBuffer import WriteBufferByteBased

from plc4py.utils.GenericTypes import ByteOrder


def test_read_buffer_get_bytes(mocker) -> None:
    bb: bytearray = bytearray(b"\x01")
    rb: ReadBufferByteBased = ReadBufferByteBased(bb, ByteOrder.LITTLE_ENDIAN)
    ba: memoryview = rb.get_bytes()
    assert 1 == ba.nbytes


def test_read_buffer_get_pos(mocker) -> None:
    bb: bytearray = bytearray(b"\x01")
    rb: ReadBufferByteBased = ReadBufferByteBased(bb, ByteOrder.LITTLE_ENDIAN)
    pos: int = rb.get_pos()
    assert 0 == pos


def test_read_buffer_read_bit_little_endian(mocker) -> None:
    bb: bytearray = bytearray(b"\x01")
    rb: ReadBufferByteBased = ReadBufferByteBased(bb, ByteOrder.LITTLE_ENDIAN)
    assert False == rb.read_bit("Bit")


def test_read_buffer_read_bit_big_endian(mocker) -> None:
    bb: bytearray = bytearray(b"\x01")
    rb: ReadBufferByteBased = ReadBufferByteBased(bb, ByteOrder.BIG_ENDIAN)
    assert False == rb.read_bit("Bit")


def test_read_buffer_read_bit_get_pos(mocker) -> None:
    bb: bytearray = bytearray(b"\x01")
    rb: ReadBufferByteBased = ReadBufferByteBased(bb, ByteOrder.BIG_ENDIAN)
    rb.read_bit("Bit")
    pos: int = rb.get_pos()
    assert 1 == pos


def test_read_buffer_read_bit_x2_get_pos(mocker) -> None:
    bb: bytearray = bytearray(b"\x01")
    rb: ReadBufferByteBased = ReadBufferByteBased(bb, ByteOrder.BIG_ENDIAN)
    rb.read_bit("Bit")
    rb.read_bit("Bit")
    pos: int = rb.get_pos()
    assert 2 == pos


def test_read_buffer_read_bit_x2_little_endian(mocker) -> None:
    bb: bytearray = bytearray(b"\x03")
    rb: ReadBufferByteBased = ReadBufferByteBased(bb, ByteOrder.LITTLE_ENDIAN)
    assert rb.read_bit("Bit") == False
    assert rb.read_bit("Bit") == False


def test_read_buffer_read_bit_x2_big_endian(mocker) -> None:
    bb: bytearray = bytearray(b"\x03")
    rb: ReadBufferByteBased = ReadBufferByteBased(bb, ByteOrder.BIG_ENDIAN)
    assert rb.read_bit("Bit") == False
    assert rb.read_bit("Bit") == False


def test_read_buffer_read_unsigned_byte_get_pos(mocker) -> None:
    bb: bytearray = bytearray(b"\x03")
    rb: ReadBufferByteBased = ReadBufferByteBased(bb, ByteOrder.BIG_ENDIAN)
    assert rb.read_unsigned_byte(logical_name="Byte") == 3
    pos: int = rb.get_pos()
    assert pos == 8


def test_read_buffer_read_unsigned_byte_little_endian(mocker) -> None:
    bb: bytearray = bytearray(b"\x03")
    rb: ReadBufferByteBased = ReadBufferByteBased(bb, ByteOrder.LITTLE_ENDIAN)
    assert rb.read_unsigned_byte(logical_name="Byte") == 192


def test_read_buffer_read_unsigned_byte_big_endian(mocker) -> None:
    bb: bytearray = bytearray(b"\x12")
    rb: ReadBufferByteBased = ReadBufferByteBased(bb, ByteOrder.BIG_ENDIAN)
    assert rb.read_unsigned_byte(logical_name="Byte") == 18


def test_read_buffer_read_unsigned_byte_little_endian_niblet(mocker) -> None:
    bb: bytearray = bytearray(b"\x84")
    rb: ReadBufferByteBased = ReadBufferByteBased(bb, ByteOrder.LITTLE_ENDIAN)
    assert rb.read_unsigned_byte(bit_length=4, logical_name="Byte") == 1


def test_read_buffer_read_unsigned_byte_big_endian_niblet(mocker) -> None:
    bb: bytearray = bytearray(b"\x84")
    rb: ReadBufferByteBased = ReadBufferByteBased(bb, ByteOrder.BIG_ENDIAN)
    assert rb.read_unsigned_byte(bit_length=4, logical_name="Byte") == 8


def test_read_buffer_read_unsigned_short_little_endian(mocker) -> None:
    bb: bytearray = bytearray(b"\x03")
    rb: ReadBufferByteBased = ReadBufferByteBased(bb, ByteOrder.LITTLE_ENDIAN)
    assert rb.read_unsigned_short(logical_name="Byte") == 192


def test_read_buffer_read_unsigned_short_big_endian(mocker) -> None:
    bb: bytearray = bytearray(b"\x03")
    rb: ReadBufferByteBased = ReadBufferByteBased(bb, ByteOrder.BIG_ENDIAN)
    assert rb.read_unsigned_short(logical_name="Byte") == 3


def test_read_buffer_read_unsigned_int_little_endian(mocker) -> None:
    bb: bitarray = bitarray("0000 0000 0000 0000 0000 0000 0000 0011")
    rb: ReadBufferByteBased = ReadBufferByteBased(bb, ByteOrder.LITTLE_ENDIAN)
    assert rb.read_unsigned_int(logical_name="Byte") == 3221225472


def test_read_buffer_read_unsigned_int_big_endian(mocker) -> None:
    bb: bitarray = bitarray("0000 0000 0000 0000 0000 0000 0000 0011")
    rb: ReadBufferByteBased = ReadBufferByteBased(bb, ByteOrder.BIG_ENDIAN)
    assert rb.read_unsigned_int(logical_name="Byte") == 3


def test_read_buffer_read_float_little_endian(mocker) -> None:
    bb: bitarray = bitarray("0000 0000 0000 0000 0000 0000 0000 0011")
    rb: ReadBufferByteBased = ReadBufferByteBased(bb, ByteOrder.LITTLE_ENDIAN)
    assert rb.read_float(logical_name="Byte") == -2.0


def test_read_buffer_read_float_big_endian(mocker) -> None:
    bb: bitarray = bitarray("1100 0000 0000 0000 0000 0000 0000 0000")
    rb: ReadBufferByteBased = ReadBufferByteBased(bb, ByteOrder.BIG_ENDIAN)
    assert rb.read_float(logical_name="Byte") == -2.0
