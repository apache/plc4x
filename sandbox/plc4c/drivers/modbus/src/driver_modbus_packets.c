/*
  Licensed to the Apache Software Foundation (ASF) under one
  or more contributor license agreements.  See the NOTICE file
  distributed with this work for additional information
  regarding copyright ownership.  The ASF licenses this file
  to you under the Apache License, Version 2.0 (the
  "License"); you may not use this file except in compliance
  with the License.  You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing,
  software distributed under the License is distributed on an
  "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
  KIND, either express or implied.  See the License for the
  specific language governing permissions and limitations
  under the License.
*/

#include <ctype.h>
#include <plc4c/spi/types_private.h>
#include <stdlib.h>
#include <string.h>

#include "modbus_tcp_adu.h"
#include "plc4c/driver_modbus.h"
#include "plc4c/driver_s7_encode_decode.h"

/**
 * Function used by the driver to tell the transport if there is a full
 * packet the driver can operate on available.
 *
 * @param buffer_data pointer to the buffer
 * @param buffer_length length of the buffer
 * @return positive integer = length of the packet, 0 = not enough,
 * try again later, negative integer = remove the given amount of bytes
 * from the buffer.
 */
int16_t plc4c_driver_modbus_select_message_function(uint8_t* buffer_data,
                                                uint16_t buffer_length) {
  // The length information is located in bytes 5 and 6
  if (buffer_length >= 6) {
    uint16_t packet_length =
        ((uint16_t) (buffer_data + 4) << 8) |
        ((uint16_t) (buffer_data + 5));
    packet_length += 6;
    if (buffer_length >= packet_length) {
      return packet_length;
    }
  }
  // In all other cases, we'll just have to wait for the next time.
  return 0;
}

plc4c_return_code plc4c_driver_modbus_send_packet(plc4c_connection* connection,
                              plc4c_modbus_read_write_modbus_tcp_adu* packet) {
  // Get the size required to contain the serialized form of this packet.
  uint16_t packet_size =
      plc4c_modbus_read_write_modbus_tcp_adu_length_in_bytes(packet);

  // Serialize this message to a byte-array.
  plc4c_spi_write_buffer* write_buffer;
  plc4c_return_code return_code =
      plc4c_spi_write_buffer_create(packet_size, &write_buffer);
  if (return_code != OK) {
    return return_code;
  }
  plc4c_modbus_read_write_modbus_tcp_adu_serialize(write_buffer, packet);

  // Now send this to the recipient.
  return_code = connection->transport->send_message(
      connection->transport_configuration, write_buffer);
  if (return_code != OK) {
    return return_code;
  }

  return OK;
}

plc4c_return_code plc4c_driver_modbus_receive_packet(plc4c_connection* connection,
                                 plc4c_modbus_read_write_modbus_tcp_adu** packet) {
  // Check with the transport if there is a packet available.
  // If it is, get a read_buffer for reading it.
  plc4c_spi_read_buffer* read_buffer;
  plc4c_return_code return_code = connection->transport->select_message(
      connection->transport_configuration,
      6, plc4c_driver_modbus_select_message_function,
      &read_buffer);
  // OK is only returned if a packet is available.
  if (return_code != OK) {
    return return_code;
  }

  // Parse the packet by consuming the read_buffer data.
  *packet = NULL;
  return_code = plc4c_modbus_read_write_modbus_tcp_adu_parse(read_buffer, true, packet);
  if (return_code != OK) {
    return return_code;
  }

  // In this case a packet was available and parsed.
  return OK;
}

plc4c_return_code plc4c_driver_modbus_create_modbus_read_request(
    plc4c_read_request* read_request,
    plc4c_list** modbus_read_request_packets) {
  // Initialize the packet list.
  plc4c_utils_list_create(modbus_read_request_packets);

  // For every item in the request, create a separate packet.
  plc4c_list_element* item = read_request->items->tail;
  while (item != NULL) {
    // Get the item address from the API request.
    char* itemAddress = item->value;

    // Create a packet from the current item.
    plc4c_modbus_read_write_modbus_pdu* packet;
    plc4c_return_code result = plc4c_driver_modbus_encode_address(
        itemAddress, &packet);
    if(result != OK) {
      return result;
    }

    // Add the packet to the list of packets.
    plc4c_utils_list_insert_head_value(*modbus_read_request_packets, packet);

    // Proceed with the next item.
    item = item->next;
  }
  return OK;
}

plc4c_return_code plc4c_driver_modbus_create_modbus_write_request(
    plc4c_write_request* write_request,
    plc4c_modbus_read_write_modbus_tcp_adu** modbus_read_request_packet) {

  // TODO: Implement this ...

  return OK;
}

