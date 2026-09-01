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
How many bytes a write puts on the wire, taken from the driver rather than from the codec.

The driver used to size its buffer as ``quantity * data_type_size`` rounded up to an even number,
which disagreed with what it then wrote: three booleans got a six byte buffer and three bits of
content. The size now comes from the same code that decides the layout, so these two cannot drift
apart - and since the resulting byte counts changed for the sub-register types, they are pinned
down here.
"""

from plc4py.drivers.modbus.ModbusConfiguration import ModbusConfiguration
from plc4py.drivers.modbus.ModbusDevice import ModbusDevice
from plc4py.drivers.modbus.ModbusTag import ModbusTagHoldingRegister
from plc4py.protocols.modbus.readwrite.ModbusDataType import ModbusDataType
from plc4py.spi.values.PlcValues import PlcBOOL, PlcList, PlcUSINT, PlcUINT, PlcUDINT


def _device():
    return ModbusDevice(ModbusConfiguration("modbus://127.0.0.1:502"))


def _write(data_type, quantity, values):
    tag = ModbusTagHoldingRegister(1, quantity, data_type)
    return _device()._serialize_data_items(tag, values)


def test_three_bools_occupy_a_single_register():
    """Three bits round up to one register, not to three."""
    values = PlcList([PlcBOOL(True), PlcBOOL(False), PlcBOOL(True)])

    written = _write(ModbusDataType.BOOL, 3, values)

    assert len(written) == 2
    assert written == [0b10100000, 0x00]


def test_three_bytes_occupy_two_registers():
    """Two bytes to a register, so three of them need two registers with a pad byte."""
    values = PlcList([PlcUSINT(1), PlcUSINT(2), PlcUSINT(3)])

    written = _write(ModbusDataType.USINT, 3, values)

    assert len(written) == 4
    assert written == [0x01, 0x02, 0x03, 0x00]


def test_two_bytes_fill_one_register_exactly():
    values = PlcList([PlcUSINT(1), PlcUSINT(2)])

    written = _write(ModbusDataType.USINT, 2, values)

    assert len(written) == 2
    assert written == [0x01, 0x02]


def test_a_lone_byte_still_occupies_a_whole_register():
    written = _write(ModbusDataType.USINT, 1, PlcUSINT(0x2A))

    assert len(written) == 2
    assert written == [0x00, 0x2A]


def test_register_wide_values_are_unchanged():
    """The types that already filled whole registers must write exactly as before."""
    values = PlcList([PlcUINT(0x1234), PlcUINT(0x5678)])

    written = _write(ModbusDataType.UINT, 2, values)

    assert len(written) == 4
    assert written == [0x12, 0x34, 0x56, 0x78]


def test_wider_than_a_register_values_are_unchanged():
    written = _write(ModbusDataType.UDINT, 2, PlcList([PlcUDINT(1), PlcUDINT(2)]))

    assert len(written) == 8
    assert written == [0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x02]
