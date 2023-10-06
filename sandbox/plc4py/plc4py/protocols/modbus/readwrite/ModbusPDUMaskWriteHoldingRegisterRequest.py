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
from plc4py.protocols.modbus.readwrite.ModbusPDU import ModbusPDU
from plc4py.protocols.modbus.readwrite.ModbusPDU import ModbusPDUBuilder
from plc4py.spi.generation.ReadBuffer import ReadBuffer
from plc4py.spi.generation.WriteBuffer import WriteBuffer
import math


@dataclass
class ModbusPDUMaskWriteHoldingRegisterRequest(PlcMessage, ModbusPDU):
    reference_address: int
    and_mask: int
    or_mask: int
    # Accessors for discriminator values.
    error_flag: bool = False
    function_flag: int = 0x16
    response: bool = False

    def __post_init__(self):
        super().__init__()

    def serialize_modbus_pdu_child(self, write_buffer: WriteBuffer):
        write_buffer.push_context("ModbusPDUMaskWriteHoldingRegisterRequest")

        # Simple Field (referenceAddress)
        write_buffer.write_unsigned_short(
            self.reference_address, logical_name="referenceAddress"
        )

        # Simple Field (andMask)
        write_buffer.write_unsigned_short(self.and_mask, logical_name="andMask")

        # Simple Field (orMask)
        write_buffer.write_unsigned_short(self.or_mask, logical_name="orMask")

        write_buffer.pop_context("ModbusPDUMaskWriteHoldingRegisterRequest")

    def length_in_bytes(self) -> int:
        return int(math.ceil(float(self.get_length_in_bits() / 8.0)))

    def get_length_in_bits(self) -> int:
        length_in_bits: int = super().get_length_in_bits()
        _value: ModbusPDUMaskWriteHoldingRegisterRequest = self

        # Simple field (referenceAddress)
        length_in_bits += 16

        # Simple field (andMask)
        length_in_bits += 16

        # Simple field (orMask)
        length_in_bits += 16

        return length_in_bits

    @staticmethod
    def static_parse_builder(read_buffer: ReadBuffer, response: bool):
        read_buffer.push_context("ModbusPDUMaskWriteHoldingRegisterRequest")

        self.reference_address = read_simple_field(
            "referenceAddress", read_unsigned_int
        )

        self.and_mask = read_simple_field("andMask", read_unsigned_int)

        self.or_mask = read_simple_field("orMask", read_unsigned_int)

        read_buffer.pop_context("ModbusPDUMaskWriteHoldingRegisterRequest")
        # Create the instance
        return ModbusPDUMaskWriteHoldingRegisterRequestBuilder(
            reference_address, and_mask, or_mask
        )

    def equals(self, o: object) -> bool:
        if self == o:
            return True

        if not isinstance(o, ModbusPDUMaskWriteHoldingRegisterRequest):
            return False

        that: ModbusPDUMaskWriteHoldingRegisterRequest = (
            ModbusPDUMaskWriteHoldingRegisterRequest(o)
        )
        return (
            (self.reference_address == that.reference_address)
            and (self.and_mask == that.and_mask)
            and (self.or_mask == that.or_mask)
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
class ModbusPDUMaskWriteHoldingRegisterRequestBuilder(ModbusPDUBuilder):
    referenceAddress: int
    andMask: int
    orMask: int

    def __post_init__(self):
        pass

    def build(
        self,
    ) -> ModbusPDUMaskWriteHoldingRegisterRequest:
        modbus_pdu_mask_write_holding_register_request: ModbusPDUMaskWriteHoldingRegisterRequest = ModbusPDUMaskWriteHoldingRegisterRequest(
            self.reference_address, self.and_mask, self.or_mask
        )
        return modbus_pdu_mask_write_holding_register_request
