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
from plc4py.protocols.umas import StaticHelper
from plc4py.protocols.umas.readwrite.DriverType import DriverType
from plc4py.protocols.umas.readwrite.UmasADU import UmasADU
from plc4py.protocols.umas.readwrite.UmasADU import UmasADUBuilder
from plc4py.protocols.umas.readwrite.UmasPDU import UmasPDU
from plc4py.spi.generation.ReadBuffer import ReadBuffer
from plc4py.spi.generation.WriteBuffer import WriteBuffer
from plc4py.utils.GenericTypes import ByteOrder
import math


@dataclass
class UmasAsciiADU(UmasADU):
    address: int
    pdu: UmasPDU
    # Arguments.
    response: bool
    # Accessors for discriminator values.
    driver_type: DriverType = DriverType.Umas_ASCII

    def serialize_umas_adu_child(self, write_buffer: WriteBuffer):
        write_buffer.push_context("UmasAsciiADU")

        # Simple Field (address)
        write_buffer.write_unsigned_byte(self.address, logical_name="address")

        # Simple Field (pdu)
        write_buffer.write_serializable(self.pdu, logical_name="pdu")

        # Checksum Field (checksum) (Calculated)
        write_buffer.write_unsigned_byte(
            int(StaticHelper.ascii_lrc_check(address, pdu)), logical_name="crc"
        )

        write_buffer.pop_context("UmasAsciiADU")

    def length_in_bytes(self) -> int:
        return int(math.ceil(float(self.length_in_bits() / 8.0)))

    def length_in_bits(self) -> int:
        length_in_bits: int = super().length_in_bits()
        _value: UmasAsciiADU = self

        # Simple field (address)
        length_in_bits += 8

        # Simple field (pdu)
        length_in_bits += self.pdu.length_in_bits()

        # Checksum Field (checksum)
        length_in_bits += 8

        return length_in_bits

    @staticmethod
    def static_parse_builder(
        read_buffer: ReadBuffer, driver_type: DriverType, response: bool
    ):
        read_buffer.push_context("UmasAsciiADU")

        address: int = read_buffer.read_unsigned_byte(
            logical_name="address",
            bit_length=8,
            byte_order=ByteOrder.BIG_ENDIAN,
            driver_type=driver_type,
            response=response,
        )

        pdu: UmasPDU = read_buffer.read_complex(
            read_function=UmasPDU.static_parse,
            logical_name="pdu",
            byte_order=ByteOrder.BIG_ENDIAN,
            driver_type=driver_type,
            response=response,
        )

        crc: int = read_buffer.read_unsigned_byte(
            logical_name="crc",
            byte_order=ByteOrder.BIG_ENDIAN,
            driver_type=driver_type,
            response=response,
        )

        read_buffer.pop_context("UmasAsciiADU")
        # Create the instance
        return UmasAsciiADUBuilder(address, pdu)

    def equals(self, o: object) -> bool:
        if self == o:
            return True

        if not isinstance(o, UmasAsciiADU):
            return False

        that: UmasAsciiADU = UmasAsciiADU(o)
        return (
            (self.address == that.address)
            and (self.pdu == that.pdu)
            and super().equals(that)
            and True
        )

    def hash_code(self) -> int:
        return hash(self)

    def __str__(self) -> str:
        pass
        # write_buffer_box_based: WriteBufferBoxBased = WriteBufferBoxBased(True, True)
        # try:
        #    write_buffer_box_based.writeSerializable(self)
        # except SerializationException as e:
        #    raise PlcRuntimeException(e)

        # return "\n" + str(write_buffer_box_based.get_box()) + "\n"


@dataclass
class UmasAsciiADUBuilder(UmasADUBuilder):
    address: int
    pdu: UmasPDU

    def build(
        self,
        response: bool,
    ) -> UmasAsciiADU:
        umas_ascii_adu: UmasAsciiADU = UmasAsciiADU(response, self.address, self.pdu)
        return umas_ascii_adu
