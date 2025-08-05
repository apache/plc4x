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
from plc4py.spi.generation.ReadBuffer import ReadBuffer
from plc4py.spi.generation.WriteBuffer import WriteBuffer
from plc4py.utils.GenericTypes import ByteOrder
from typing import Any
from typing import List
import math


@dataclass
class ModbusTcpADU:
    transaction_identifier: int
    unit_identifier: int
    pdu_array: List[int]
    PROTOCOL_IDENTIFIER: int = 0x0000

    def serialize(self, write_buffer: WriteBuffer):
        write_buffer.push_context("ModbusTcpADU")

        # Simple Field (transactionIdentifier)
        write_buffer.write_unsigned_short(
            self.transaction_identifier,
            bit_length=16,
            logical_name="transactionIdentifier",
        )

        # Const Field (protocolIdentifier)
        write_buffer.write_unsigned_short(
            self.PROTOCOL_IDENTIFIER, logical_name="protocolIdentifier"
        )

        # Implicit Field (length) (Used for parsing, but its value is not stored as it's implicitly given by the objects content)
        length: int = int(len(self.pdu_array)) + int(1)
        write_buffer.write_unsigned_short(length, logical_name="length")

        # Simple Field (unitIdentifier)
        write_buffer.write_unsigned_byte(
            self.unit_identifier, bit_length=8, logical_name="unitIdentifier"
        )

        # Array Field (pduArray)
        write_buffer.write_byte_array(self.pdu_array, logical_name="pdu_array")

        write_buffer.pop_context("ModbusTcpADU")

    def length_in_bytes(self) -> int:
        return int(math.ceil(float(self.length_in_bits() / 8.0)))

    def length_in_bits(self) -> int:
        length_in_bits: int = 0
        _value: ModbusTcpADU = self

        # Simple field (transactionIdentifier)
        length_in_bits += 16

        # Const Field (protocolIdentifier)
        length_in_bits += 16

        # Implicit Field (length)
        length_in_bits += 16

        # Simple field (unitIdentifier)
        length_in_bits += 8

        # Array field
        if self.pdu_array is not None:
            length_in_bits += 8 * len(self.pdu_array)

        return length_in_bits

    @staticmethod
    def static_parse(read_buffer: ReadBuffer, **kwargs):
        return ModbusTcpADU.static_parse_context(read_buffer)

    @staticmethod
    def static_parse_context(read_buffer: ReadBuffer):
        read_buffer.push_context("ModbusTcpADU")

        transaction_identifier: int = read_buffer.read_unsigned_short(
            logical_name="transaction_identifier",
            bit_length=16,
            byte_order=ByteOrder.BIG_ENDIAN,
        )

        PROTOCOL_IDENTIFIER: int = read_buffer.read_unsigned_short(
            logical_name="protocol_identifier", byte_order=ByteOrder.BIG_ENDIAN
        )

        length: int = read_buffer.read_unsigned_short(
            logical_name="length", byte_order=ByteOrder.BIG_ENDIAN
        )

        unit_identifier: int = read_buffer.read_unsigned_byte(
            logical_name="unit_identifier",
            bit_length=8,
            byte_order=ByteOrder.BIG_ENDIAN,
        )

        pdu_array: List[Any] = read_buffer.read_array_field(
            logical_name="pduArray",
            read_function=read_buffer.read_byte,
            count=length - int(1),
            byte_order=ByteOrder.BIG_ENDIAN,
        )

        read_buffer.pop_context("ModbusTcpADU")
        # Create the instance
        _modbus_tcp_adu: ModbusTcpADU = ModbusTcpADU(
            transaction_identifier, unit_identifier, pdu_array
        )
        return _modbus_tcp_adu

    def equals(self, o: object) -> bool:
        if self == o:
            return True

        if not isinstance(o, ModbusTcpADU):
            return False

        that: ModbusTcpADU = ModbusTcpADU(o)
        return (
            (self.transaction_identifier == that.transaction_identifier)
            and (self.unit_identifier == that.unit_identifier)
            and (self.pdu_array == that.pdu_array)
            and True
        )

    def hash_code(self) -> int:
        return hash(self)

    def __str__(self) -> str:
        # TODO:- Implement a generic python object to probably json convertor or something.
        return ""
