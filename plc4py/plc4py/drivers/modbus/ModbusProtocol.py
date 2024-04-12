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
import logging
from asyncio import Future
from dataclasses import dataclass, field
from typing import Dict, Awaitable

from plc4py.protocols.modbus.readwrite.DriverType import DriverType
from plc4py.spi.generation.ReadBuffer import ReadBufferByteBased

from plc4py.protocols.modbus.readwrite.ModbusTcpADU import ModbusTcpADU
from plc4py.spi.Plc4xBaseProtocol import Plc4xBaseProtocol
from plc4py.utils.GenericTypes import ByteOrder


@dataclass
class ModbusProtocol(Plc4xBaseProtocol):
    messages: Dict[int, Future] = field(default_factory=lambda: {})

    def packet_length_estimator(self, read_buffer: ReadBufferByteBased):
        current_position = read_buffer.position
        read_buffer.position = 8 * 4
        packet_length: int = read_buffer.read_unsigned_short()
        read_buffer.position = current_position
        return packet_length + current_position < len(read_buffer.bb)

    def data_received(self, data):
        """Unpack the adu and return the pdu"""
        read_buffer = ReadBufferByteBased(
            bytearray(data), byte_order=ByteOrder.BIG_ENDIAN
        )
        while self.packet_length_estimator(read_buffer):
            adu: ModbusTcpADU = ModbusTcpADU.static_parse_builder(
                read_buffer, DriverType.MODBUS_TCP, True
            ).build(True)
            if adu.transaction_identifier in self.messages:
                self.messages[adu.transaction_identifier].set_result(adu.pdu)
                self.messages.pop(adu.transaction_identifier)
            else:
                logging.error("Unsolicited message returned")
                self.close()

    def write_wait_for_response(self, data, transport, transaction_id, message_future):
        """Writes a message to the wire and records the identifier to identify and route the response"""
        self.messages[transaction_id] = message_future
        transport.write(data)

    def close(self):
        """Clean up the message which didn't receive a response"""
        for key, message in self.messages.items:
            logging.debug("Removing un-replied message with identifier " + str(key))
            message.set_result(None)
            self.messages.pop()
        self.messages = None
