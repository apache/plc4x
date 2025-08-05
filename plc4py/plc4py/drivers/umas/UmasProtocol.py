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
from typing import Awaitable, Dict, Tuple

from plc4py.protocols.umas.readwrite.ModbusPDU import ModbusPDU
from plc4py.protocols.umas.readwrite.ModbusTcpADU import ModbusTcpADU
from plc4py.protocols.umas.readwrite.UmasPDU import UmasPDU, UmasPDUBuilder
from plc4py.protocols.umas.readwrite.UmasPDUItem import UmasPDUItem
from plc4py.spi.generation.ReadBuffer import ReadBufferByteBased
from plc4py.spi.generation.WriteBuffer import WriteBufferByteBased
from plc4py.spi.Plc4xBaseProtocol import Plc4xBaseProtocol
from plc4py.utils.GenericTypes import AtomicInteger, ByteOrder


@dataclass
class UmasProtocol(Plc4xBaseProtocol):
    unit_identifier: int = 1
    messages: Dict[int, Tuple[int, Future]] = field(default_factory=lambda: {})
    _transaction_generator: AtomicInteger = field(
        default_factory=lambda: AtomicInteger()
    )

    def data_received(self, data):
        """Unpack the adu and return the pdu"""
        read_buffer = ReadBufferByteBased(
            bytearray(data), byte_order=ByteOrder.BIG_ENDIAN
        )
        adu: ModbusTcpADU = ModbusTcpADU.static_parse(read_buffer)
        if adu.transaction_identifier in self.messages:
            read_buffer = ReadBufferByteBased(
                bytearray(adu.pdu_array), byte_order=ByteOrder.LITTLE_ENDIAN
            )
            pdu: ModbusPDU = ModbusPDU.static_parse(
                read_buffer,
                umas_request_function_key=self.messages[adu.transaction_identifier][0],
                byte_length=len(adu.pdu_array) - 1,
            )
            if isinstance(pdu, UmasPDU):
                self.messages[adu.transaction_identifier][1].set_result(pdu.item)
            else:
                logging.error("Modbus Error Message Received")
                self.close()
        else:
            logging.error("Unsolicited message returned")
            self.close()

    def write_wait_for_response(
        self, umas_pdu_item: UmasPDUItem, transport, message_future
    ):
        pdu = UmasPDUBuilder(umas_pdu_item).build(umas_pdu_item.umas_function_key, 0)

        write_buffer = WriteBufferByteBased(
            pdu.length_in_bytes(), byte_order=ByteOrder.LITTLE_ENDIAN
        )
        pdu.serialize(write_buffer)

        adu = ModbusTcpADU(
            self._transaction_generator.increment(),
            self.unit_identifier,
            write_buffer.get_bytes().tolist(),
        )

        adu_write_buffer = WriteBufferByteBased(
            adu.length_in_bytes(), ByteOrder.BIG_ENDIAN
        )
        adu.serialize(adu_write_buffer)

        """Writes a message to the wire and records the identifier to identify and route the response"""
        self.messages[adu.transaction_identifier] = (
            pdu.umas_request_function_key,
            message_future,
        )
        transport.write(adu_write_buffer.get_bytes())

    def close(self):
        """Clean up the message which didn't receive a response"""
        for key, message in self.messages.items:
            logging.debug("Removing un-replied message with identifier " + str(key))
            message.set_result(None)
        self.messages = None
