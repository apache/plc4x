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
from typing import Awaitable, Dict

from plc4py.protocols.modbus.readwrite.DriverType import DriverType
from plc4py.protocols.modbus.readwrite.ModbusTcpADU import ModbusTcpADU
from plc4py.spi.generation.ReadBuffer import ReadBufferByteBased
from plc4py.spi.Plc4xBaseProtocol import Plc4xBaseProtocol
from plc4py.utils.GenericTypes import ByteOrder


@dataclass
class ModbusProtocol(Plc4xBaseProtocol):
    messages: Dict[int, Future] = field(default_factory=lambda: {})

    def packet_length_estimator(self, read_buffer: ReadBufferByteBased):
        """
        Estimate the length of the packet based on the provided `read_buffer`.

        Args:
            read_buffer (ReadBufferByteBased): The buffer containing the data to be read.

        Returns:
            bool: True if the estimated packet length plus the current position is less than the length of the buffer, False otherwise.
        """
        current_position = read_buffer.position
        logging.debug(
            f"Current position in packet_length_estimator: {current_position}"
        )
        read_buffer.position = 8 * 4
        packet_length: int = read_buffer.read_unsigned_short()
        logging.debug(f"Packet length in packet_length_estimator: {packet_length}")
        read_buffer.position = current_position
        logging.debug(
            f"Buffer length in packet_length_estimator: {len(read_buffer.bb)}"
        )
        return packet_length + current_position < len(read_buffer.bb)

    def data_received(self, data):
        """
        Unpack the adu and return the pdu.

        This function continuously processes incoming data until the complete
        ADU (ModbusTcpADU) is extracted. If the ADU transaction identifier is
        found in the messages dictionary, the corresponding PDU (data) is set
        as a result of the associated Future object. If the transaction
        identifier is not found, an error message is logged and the connection
        is closed.

        Args:
            data (bytes): The raw data received from the transport.
        """
        # Convert data to ReadBufferByteBased
        read_buffer = ReadBufferByteBased(
            bytearray(data), byte_order=ByteOrder.BIG_ENDIAN
        )

        # Continuously process incoming data until the complete ADU is
        # extracted.
        while self.packet_length_estimator(read_buffer):
            logging.debug(f"Processing {len(data)} bytes of data")
            # Parse the ADU from the ReadBufferByteBased.
            adu: ModbusTcpADU = ModbusTcpADU.static_parse_builder(
                read_buffer, DriverType.MODBUS_TCP, True
            ).build(True)
            # logging.debug(f"Received ADU: {adu}")
            # If the ADU transaction identifier is in the messages dictionary,
            # set the PDU (data) as a result of the associated Future object.
            if adu.transaction_identifier in self.messages:
                logging.debug(f"Found transaction ID {adu.transaction_identifier}")
                self.messages[adu.transaction_identifier].set_result(adu.pdu)
                self.messages.pop(adu.transaction_identifier)
            # If the ADU transaction identifier is not found, log an error
            # message and close the connection.
            else:
                logging.error("Unsolicited message returned")
                logging.debug(f"Transaction ID {adu.transaction_identifier} not found")
                self.close()

    def write_wait_for_response(self, data, transport, transaction_id, message_future):
        """
        Writes a message to the wire and records the transaction identifier to identify and route the response.

        Args:
            data (bytes): The data to be written to the wire.
            transport (asyncio.Transport): The transport to use for writing the data.
            transaction_id (int): The transaction identifier of the message.
            message_future (concurrent.futures.Future): The future to set the result of when the response is received.
        """
        # Record the message's transaction identifier and the future to set the result of.
        self.messages[transaction_id] = message_future
        # Write the data to the wire.
        transport.write(data)

    def close(self):
        """
        Clean up the message which didn't receive a response.

        This function iterates over the messages dictionary,
        sets the result of each message to None and removes it from the dictionary.
        After that, it sets the messages dictionary to None.
        """
        # Iterate over the messages dictionary
        for key, message in self.messages.items():
            # Log the removal of un-replied message
            logging.debug("Removing un-replied message with identifier " + str(key))
            # Set the result of the message to None
            message.set_result(None)
            # Remove the message from the dictionary
            self.messages.pop(key)
        # Set the messages dictionary to None
        self.messages = None
