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

from plc4py.api.messages.PlcMessage import PlcMessage
from plc4py.protocols.modbus import StaticHelper
from plc4py.protocols.modbus.readwrite.DriverType import DriverType
from plc4py.protocols.modbus.readwrite.ModbusADU import ModbusADU
from plc4py.protocols.modbus.readwrite.ModbusADU import ModbusADUBuilder
from plc4py.protocols.modbus.readwrite.ModbusPDU import ModbusPDU
from plc4py.spi.generation.ReadBuffer import ReadBuffer
from plc4py.spi.generation.WriteBuffer import WriteBuffer
import math


@dataclass
class ModbusAsciiADU(PlcMessage, ModbusADU):
    address: int
    pdu: ModbusPDU
    # Arguments.
    response: bool
    # Accessors for discriminator values.
    driver_type: DriverType = DriverType.MODBUS_ASCII

    def __post_init__(self):
        super().__init__(self.response)

    def serialize_modbus_adu_child(self, write_buffer: WriteBuffer):
        write_buffer.push_context("ModbusAsciiADU")

        # Simple Field (address)
        write_buffer.write_unsigned_byte(self.address, logical_name="address")

        # Simple Field (pdu)
        write_buffer.write_serializable(self.pdu, logical_name="pdu")

        # Checksum Field (checksum) (Calculated)
        write_buffer.write_unsigned_byte(
            int(StaticHelper.ascii_lrc_check(self.address, self.pdu)),
            logical_name="crc",
        )

        write_buffer.pop_context("ModbusAsciiADU")

    def length_in_bytes(self) -> int:
        return int(math.ceil(float(self.get_length_in_bits() / 8.0)))

    def get_length_in_bits(self) -> int:
        length_in_bits: int = super().get_length_in_bits()
        _value: ModbusAsciiADU = self

        # Simple field (address)
        length_in_bits += 8

        # Simple field (pdu)
        length_in_bits += self.pdu.get_length_in_bits()

        # Checksum Field (checksum)
        length_in_bits += 8

        return length_in_bits

    @staticmethod
    def static_parse_builder(
        read_buffer: ReadBuffer, driver_type: DriverType, response: bool
    ):
        read_buffer.push_context("ModbusAsciiADU")

        self.address = read_simple_field(
            "address", read_unsigned_short, WithOption.WithByteOrder(get_bi_g__endian())
        )

        self.pdu = read_simple_field(
            "pdu",
            DataReaderComplexDefault(
                ModbusPDU.static_parse(read_buffer, bool(response)), read_buffer
            ),
            WithOption.WithByteOrder(get_bi_g__endian()),
        )

        crc: int = read_checksum_field(
            "crc",
            read_unsigned_short,
            (int)(ascii_lrc_check(self.address, self.pdu)),
            WithOption.WithByteOrder(get_bi_g__endian()),
        )

        read_buffer.pop_context("ModbusAsciiADU")
        # Create the instance
        return ModbusAsciiADUBuilder(address, pdu, response)

    def equals(self, o: object) -> bool:
        if self == o:
            return True

        if not isinstance(o, ModbusAsciiADU):
            return False

        that: ModbusAsciiADU = ModbusAsciiADU(o)
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
class ModbusAsciiADUBuilder(ModbusADUBuilder):
    address: int
    pdu: ModbusPDU
    response: bool

    def __post_init__(self):
        pass

    def build(self, response: bool) -> ModbusAsciiADU:
        modbus_ascii_adu: ModbusAsciiADU = ModbusAsciiADU(
            self.address, self.pdu, response
        )
        return modbus_ascii_adu
