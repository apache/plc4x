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
from plc4py.protocols.modbus.readwrite.DriverType import DriverType
from plc4py.protocols.modbus.readwrite.ModbusADU import ModbusADU
from plc4py.protocols.modbus.readwrite.ModbusADU import ModbusADUBuilder
from plc4py.protocols.modbus.readwrite.ModbusPDU import ModbusPDU
from plc4py.spi.generation.ReadBuffer import ReadBuffer
from plc4py.spi.generation.WriteBuffer import WriteBuffer
import math


@dataclass
class ModbusTcpADU(PlcMessage, ModbusADU):
    transaction_identifier: int
    unit_identifier: int
    pdu: ModbusPDU
    # Arguments.
    response: bool
    PROTOCOLIDENTIFIER: int = 0x0000
    # Accessors for discriminator values.
    driver_type: DriverType = DriverType.MODBUS_TCP

    def __post_init__(self):
        super().__init__(self.response)

    def serialize_modbus_adu_child(self, write_buffer: WriteBuffer):
        write_buffer.push_context("ModbusTcpADU")

        # Simple Field (transactionIdentifier)
        write_buffer.write_unsigned_short(
            self.transaction_identifier, logical_name="transactionIdentifier"
        )

        # Const Field (protocolIdentifier)
        write_buffer.write_unsigned_short(
            self.protocol_identifier.value, logical_name="protocolIdentifier"
        )

        # Implicit Field (length) (Used for parsing, but its value is not stored as it's implicitly given by the objects content)
        length: int = self.pdu.getlength_in_bytes(ctx) + int(1)
        write_buffer.write_unsigned_short(length, logical_name="length")

        # Simple Field (unitIdentifier)
        write_buffer.write_unsigned_byte(
            self.unit_identifier, logical_name="unitIdentifier"
        )

        # Simple Field (pdu)
        write_buffer.write_serializable(self.pdu, logical_name="pdu")

        write_buffer.pop_context("ModbusTcpADU")

    def length_in_bytes(self) -> int:
        return int(math.ceil(float(self.get_length_in_bits() / 8.0)))

    def get_length_in_bits(self) -> int:
        length_in_bits: int = super().get_length_in_bits()
        _value: ModbusTcpADU = self

        # Simple field (transactionIdentifier)
        length_in_bits += 16

        # Const Field (protocolIdentifier)
        length_in_bits += 16

        # Implicit Field (length)
        length_in_bits += 16

        # Simple field (unitIdentifier)
        length_in_bits += 8

        # Simple field (pdu)
        length_in_bits += self.pdu.get_length_in_bits()

        return length_in_bits

    @staticmethod
    def static_parse_builder(
        read_buffer: ReadBuffer, driver_type: DriverType, response: bool
    ):
        read_buffer.push_context("ModbusTcpADU")

        self.transaction_identifier = read_simple_field(
            "transactionIdentifier",
            read_unsigned_int,
            WithOption.WithByteOrder(get_bi_g__endian()),
        )

        self.protocol_identifier: int = read_const_field(
            "protocolIdentifier",
            read_unsigned_int,
            ModbusTcpADU.PROTOCOLIDENTIFIER,
            WithOption.WithByteOrder(get_bi_g__endian()),
        )

        length: int = read_implicit_field(
            "length", read_unsigned_int, WithOption.WithByteOrder(get_bi_g__endian())
        )

        self.unit_identifier = read_simple_field(
            "unitIdentifier",
            read_unsigned_short,
            WithOption.WithByteOrder(get_bi_g__endian()),
        )

        self.pdu = read_simple_field(
            "pdu",
            DataReaderComplexDefault(
                ModbusPDU.static_parse(read_buffer, bool(response)), read_buffer
            ),
            WithOption.WithByteOrder(get_bi_g__endian()),
        )

        read_buffer.pop_context("ModbusTcpADU")
        # Create the instance
        return ModbusTcpADUBuilder(
            transaction_identifier, unit_identifier, pdu, response
        )

    def equals(self, o: object) -> bool:
        if self == o:
            return True

        if not isinstance(o, ModbusTcpADU):
            return False

        that: ModbusTcpADU = ModbusTcpADU(o)
        return (
            (self.transaction_identifier == that.transaction_identifier)
            and (self.unit_identifier == that.unit_identifier)
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
class ModbusTcpADUBuilder(ModbusADUBuilder):
    transactionIdentifier: int
    unitIdentifier: int
    pdu: ModbusPDU
    response: bool

    def __post_init__(self):
        pass

    def build(self, response: bool) -> ModbusTcpADU:
        modbus_tcp_adu: ModbusTcpADU = ModbusTcpADU(
            self.transaction_identifier, self.unit_identifier, self.pdu, response
        )
        return modbus_tcp_adu
