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

from plc4py.api.exceptions.exceptions import PlcRuntimeException
from plc4py.api.exceptions.exceptions import SerializationException
from plc4py.api.messages.PlcMessage import PlcMessage
from plc4py.protocols.modbus import StaticHelper
from plc4py.protocols.modbus.readwrite.DriverType import DriverType
from plc4py.protocols.modbus.readwrite.ModbusADU import ModbusADU
from plc4py.protocols.modbus.readwrite.ModbusPDU import ModbusPDU
from plc4py.spi.generation.ReadBuffer import ReadBuffer
from plc4py.spi.generation.WriteBuffer import WriteBuffer
from plc4py.utils.ConnectionStringHandling import strtobool
from plc4py.utils.GenericTypes import ByteOrder
from typing import ClassVar
import math


@dataclass
class ModbusRtuADU(ModbusADU):
    address: int
    pdu: ModbusPDU
    # Arguments.
    response: bool
    # Accessors for discriminator values.
    driver_type: ClassVar[DriverType] = DriverType.MODBUS_RTU

    def serialize_modbus_adu_child(self, write_buffer: WriteBuffer):
        write_buffer.push_context("ModbusRtuADU")

        # Simple Field (address)
        write_buffer.write_unsigned_byte(
            self.address, bit_length=8, logical_name="address"
        )

        # Simple Field (pdu)
        write_buffer.write_serializable(self.pdu, logical_name="pdu")

        # Checksum Field (checksum) (Calculated)
        write_buffer.write_unsigned_short(
            int(StaticHelper.rtu_crc_check(address, pdu)), logical_name="crc"
        )

        write_buffer.pop_context("ModbusRtuADU")

    def length_in_bytes(self) -> int:
        return int(math.ceil(float(self.length_in_bits() / 8.0)))

    def length_in_bits(self) -> int:
        length_in_bits: int = super().length_in_bits()
        _value: ModbusRtuADU = self

        # Simple field (address)
        length_in_bits += 8

        # Simple field (pdu)
        length_in_bits += self.pdu.length_in_bits()

        # Checksum Field (checksum)
        length_in_bits += 16

        return length_in_bits

    @staticmethod
    def static_parse_builder(
        read_buffer: ReadBuffer, driver_type: DriverType, response: bool
    ):
        read_buffer.push_context("ModbusRtuADU")

        if isinstance(driver_type, str):
            driver_type = DriverType[driver_type]
        if isinstance(response, str):
            response = bool(strtobool(response))

        address: int = read_buffer.read_unsigned_byte(
            logical_name="address",
            bit_length=8,
            byte_order=ByteOrder.BIG_ENDIAN,
            driver_type=driver_type,
            response=response,
        )

        pdu: ModbusPDU = read_buffer.read_complex(
            read_function=ModbusPDU.static_parse,
            logical_name="pdu",
            byte_order=ByteOrder.BIG_ENDIAN,
            driver_type=driver_type,
            response=response,
        )

        crc: int = read_buffer.read_unsigned_short(
            logical_name="crc",
            byte_order=ByteOrder.BIG_ENDIAN,
            driver_type=driver_type,
            response=response,
        )

        read_buffer.pop_context("ModbusRtuADU")
        # Create the instance
        return ModbusRtuADUBuilder(address, pdu)

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
        # TODO:- Implement a generic python object to probably json convertor or something.
        return ""


@dataclass
class ModbusRtuADUBuilder:
    address: int
    pdu: ModbusPDU

    def build(
        self,
        response: bool,
    ) -> ModbusRtuADU:
        modbus_rtu_adu: ModbusRtuADU = ModbusRtuADU(response, self.address, self.pdu)
        return modbus_rtu_adu
