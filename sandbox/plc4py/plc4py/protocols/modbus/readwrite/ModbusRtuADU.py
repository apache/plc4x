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

from dataclasses import dataclass

from ctypes import c_bool
from ctypes import c_uint16
from ctypes import c_uint8
from plc4py.api.messages.PlcMessage import PlcMessage
from plc4py.protocols.modbus.readwrite.DriverType import DriverType
from plc4py.protocols.modbus.readwrite.ModbusADU import ModbusADU
from plc4py.protocols.modbus.readwrite.ModbusADU import ModbusADUBuilder
from plc4py.protocols.modbus.readwrite.ModbusPDU import ModbusPDU
import math


@dataclass
class ModbusRtuADU(PlcMessage, ModbusADU):
    address: c_uint8
    pdu: ModbusPDU
    # Arguments.
    response: c_bool
    # Accessors for discriminator values.
    driver_type: DriverType = DriverType.MODBUS_RTU

    def __post_init__(self):
        super().__init__(self.response)

    def serialize_modbus_adu_child(self, write_buffer: WriteBuffer):
        position_aware: PositionAware = write_buffer
        start_pos: int = position_aware.get_pos()
        write_buffer.push_context("ModbusRtuADU")

        # Simple Field (address)
        write_simple_field(
            "address",
            self.address,
            write_unsigned_short(write_buffer, 8),
            WithOption.WithByteOrder(ByteOrder.BIG_ENDIAN),
        )

        # Simple Field (pdu)
        write_simple_field(
            "pdu",
            self.pdu,
            DataWriterComplexDefault(write_buffer),
            WithOption.WithByteOrder(ByteOrder.BIG_ENDIAN),
        )

        # Checksum Field (checksum) (Calculated)
        write_checksum_field(
            "crc",
            c_uint16(modbus.readwrite.utils.StaticHelper.rtuCrcCheck(address, pdu)),
            write_unsigned_int(write_buffer, 16),
        )

        write_buffer.pop_context("ModbusRtuADU")

    def length_in_bytes(self) -> int:
        return int(math.ceil(float(self.get_length_in_bits() / 8.0)))

    def get_length_in_bits(self) -> int:
        length_in_bits: int = super().get_length_in_bits()
        _value: ModbusRtuADU = self

        # Simple field (address)
        length_in_bits += 8

        # Simple field (pdu)
        length_in_bits += self.pdu.get_length_in_bits()

        # Checksum Field (checksum)
        length_in_bits += 16

        return length_in_bits

    @staticmethod
    def static_parse_builder(
        read_buffer: ReadBuffer, driver_type: DriverType, response: c_bool
    ):
        read_buffer.pull_context("ModbusRtuADU")
        position_aware: PositionAware = read_buffer
        start_pos: int = position_aware.get_pos()
        cur_pos: int = 0

        address: c_uint8 = read_simple_field(
            "address",
            read_unsigned_short(read_buffer, 8),
            WithOption.WithByteOrder(ByteOrder.BIG_ENDIAN),
        )

        pdu: ModbusPDU = read_simple_field(
            "pdu",
            DataReaderComplexDefault(
                ModbusPDU.static_parse(read_buffer, c_bool(response)), read_buffer
            ),
            WithOption.WithByteOrder(ByteOrder.BIG_ENDIAN),
        )

        crc: c_uint16 = read_checksum_field(
            "crc",
            read_unsigned_int(read_buffer, 16),
            (c_uint16)(modbus.readwrite.utils.StaticHelper.rtuCrcCheck(address, pdu)),
            WithOption.WithByteOrder(ByteOrder.BIG_ENDIAN),
        )

        read_buffer.close_context("ModbusRtuADU")
        # Create the instance
        return ModbusRtuADUBuilder(address, pdu, response)

    def equals(self, o: object) -> bool:
        if self == o:
            return True

        if not isinstance(o, ModbusRtuADU):
            return False

        that: ModbusRtuADU = ModbusRtuADU(o)
        return (
            (self.address == that.address)
            and (self.pdu == that.pdu)
            and super().equals(that)
            and True
        )

    def hash_code(self) -> int:
        return hash(self)

    def __str__(self) -> str:
        write_buffer_box_based: WriteBufferBoxBased = WriteBufferBoxBased(True, True)
        try:
            write_buffer_box_based.writeSerializable(self)
        except SerializationException as e:
            raise RuntimeException(e)

        return "\n" + str(write_buffer_box_based.get_box()) + "\n"


@dataclass
class ModbusRtuADUBuilder(ModbusADUBuilder):
    address: c_uint8
    pdu: ModbusPDU
    response: c_bool

    def __post_init__(self):
        pass

    def build(self, response: c_bool) -> ModbusRtuADU:
        modbus_rtu_adu: ModbusRtuADU = ModbusRtuADU(self.address, self.pdu, response)
        return modbus_rtu_adu
