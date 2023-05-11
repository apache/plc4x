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
from ctypes import c_bool, c_byte

import pytest
from bitarray import bitarray

from plc4py.spi.generation.WriteBuffer import WriteBufferByteBased

from plc4py.utils.GenericTypes import ByteOrder


def test_write_buffer_get_bytes(mocker) -> None:
    wb: WriteBufferByteBased = WriteBufferByteBased(10, ByteOrder.LITTLE_ENDIAN)
    ba: memoryview = wb.get_bytes()
    assert(10 == ba.nbytes)


def test_write_buffer_get_pos(mocker) -> None:
    wb: WriteBufferByteBased = WriteBufferByteBased(10, ByteOrder.LITTLE_ENDIAN)
    pos: int = wb.get_pos()
    assert(0 == pos)


def test_write_buffer_set_bit_little_endian(mocker) -> None:
    wb: WriteBufferByteBased = WriteBufferByteBased(1, ByteOrder.LITTLE_ENDIAN)
    wb.write_bit(c_bool(True), "Test String")
    ba: memoryview = wb.get_bytes()
    assert(b'\x01' == ba.tobytes())


def test_write_buffer_set_bit_big_endian(mocker) -> None:
    wb: WriteBufferByteBased = WriteBufferByteBased(1, ByteOrder.BIG_ENDIAN)
    wb.write_bit(c_bool(True), "Test String")
    ba: memoryview = wb.get_bytes()
    assert(b'\x80' == ba.tobytes())


def test_write_buffer_set_bit_get_pos(mocker) -> None:
    wb: WriteBufferByteBased = WriteBufferByteBased(1, ByteOrder.LITTLE_ENDIAN)
    wb.write_bit(c_bool(True), "Test String")
    pos: int = wb.get_pos()
    assert(1 == pos)


def test_write_buffer_set_bit_x2_little_endian(mocker) -> None:
    wb: WriteBufferByteBased = WriteBufferByteBased(1, ByteOrder.LITTLE_ENDIAN)
    wb.write_bit(c_bool(True), "Test String 1")
    wb.write_bit(c_bool(True), "Test String 2")
    ba: memoryview = wb.get_bytes()
    assert (b'\x03' == ba.tobytes())


def test_write_buffer_set_bit_x2_big_endian(mocker) -> None:
    wb: WriteBufferByteBased = WriteBufferByteBased(1, ByteOrder.BIG_ENDIAN)
    wb.write_bit(c_bool(True), "Test String 1")
    wb.write_bit(c_bool(True), "Test String 2")
    ba: memoryview = wb.get_bytes()
    assert (b'\xc0' == ba.tobytes())


def test_write_buffer_set_bit_get_pos_dual(mocker) -> None:
    wb: WriteBufferByteBased = WriteBufferByteBased(1, ByteOrder.LITTLE_ENDIAN)
    wb.write_bit(c_bool(True), "Test String 1")
    wb.write_bit(c_bool(True), "Test String 2")
    pos: int = wb.get_pos()
    assert(pos == 2)


def test_write_buffer_set_unsigned_byte_get_pos(mocker) -> None:
    wb: WriteBufferByteBased = WriteBufferByteBased(1, ByteOrder.LITTLE_ENDIAN)
    wb.write_unsigned_byte(c_byte(0xFF), 8, "Test String 1")
    pos: int = wb.get_pos()
    assert(pos == 8)


def test_write_buffer_set_unsigned_byte_little_endian(mocker) -> None:
    wb: WriteBufferByteBased = WriteBufferByteBased(1, ByteOrder.LITTLE_ENDIAN)
    wb.write_unsigned_byte(c_byte(0x12), 8, "Test String 1")
    ba: memoryview = wb.get_bytes()
    assert (ba.obj == bitarray("01001000"))


def test_write_buffer_set_unsigned_byte_big_endian(mocker) -> None:
    wb: WriteBufferByteBased = WriteBufferByteBased(1, ByteOrder.BIG_ENDIAN)
    wb.write_unsigned_byte(c_byte(0x12), 8, "Test String 1")
    ba: memoryview = wb.get_bytes()
    assert (ba.obj == bitarray("00010010"))


def test_write_buffer_set_unsigned_byte_little_endian(mocker) -> None:
    wb: WriteBufferByteBased = WriteBufferByteBased(1, ByteOrder.LITTLE_ENDIAN)
    wb.write_unsigned_byte(c_byte(0x12), 4, "Test String 1")
    ba: memoryview = wb.get_bytes()
    assert (ba.obj == bitarray("01000000"))


def test_write_buffer_set_unsigned_byte_big_endian(mocker) -> None:
    wb: WriteBufferByteBased = WriteBufferByteBased(1, ByteOrder.BIG_ENDIAN)
    wb.write_unsigned_byte(c_byte(0x12), 4, "Test String 1")
    ba: memoryview = wb.get_bytes()
    assert (ba.obj == bitarray("00010000"))