/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

#include <plc4c/spi/types_private.h>
#include <stdlib.h>

#include "modbus_adu.h"
#include "plc4c/driver_modbus.h"

#define MIN(a,b) (((a)<(b))?(a):(b))
#define MAX(a,b) (((a)>(b))?(a):(b))

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
        ((uint16_t)*(buffer_data + 4) << 8) | ((uint16_t)*(buffer_data + 5));
    packet_length += 6;
    return (int16_t) packet_length;
  }
  // In all other cases, we'll just have to wait for the next time.
  return 0;
}

plc4c_return_code plc4c_driver_modbus_send_packet(
    plc4c_connection* connection,
    plc4c_modbus_read_write_modbus_adu* packet) {
  // Get the size required to contain the serialized form of this packet.
  uint16_t packet_size =
      plc4c_modbus_read_write_modbus_adu_length_in_bytes(plc4x_spi_context_background(), packet);

  // Serialize this message to a byte-array.
  plc4c_spi_write_buffer* write_buffer;
  plc4c_return_code return_code =
      plc4c_spi_write_buffer_create(packet_size, &write_buffer);
  if (return_code != OK) {
    return return_code;
  }
  return_code = plc4c_modbus_read_write_modbus_adu_serialize(
      plc4x_spi_context_background(), write_buffer, packet);
  if (return_code != OK) {
    return return_code;
  }

  // Now send this to the recipient.
  return_code = connection->transport->send_message(
      connection->transport_configuration, write_buffer);
  if (return_code != OK) {
    return return_code;
  }

  return OK;
}

plc4c_return_code plc4c_driver_modbus_receive_packet(
    plc4c_connection* connection,
    plc4c_modbus_read_write_modbus_adu** packet) {
  // Check with the transport if there is a packet available.
  // If it is, get a read_buffer for reading it.
  plc4c_spi_read_buffer* read_buffer;
  plc4c_return_code return_code = connection->transport->select_message(
      connection->transport_configuration, 6,
      plc4c_driver_modbus_select_message_function, &read_buffer);
  // OK is only returned if a packet is available.
  if (return_code != OK) {
    return return_code;
  }

  // Parse the packet by consuming the read_buffer data.
  *packet = NULL;
  plc4c_modbus_read_write_driver_type driver_type =
      plc4c_modbus_read_write_driver_type_MODBUS_TCP;

  return_code =
      plc4c_modbus_read_write_modbus_adu_parse(
          plc4x_spi_context_background(), read_buffer,driver_type, true, packet);
  if (return_code != OK) {
    return return_code;
  }

  // In this case a packet was available and parsed.
  return OK;
}

plc4c_return_code plc4c_driver_modbus_create_modbus_read_request(
    plc4c_driver_modbus_config* modbus_config,
    plc4c_item* read_request_item,
    plc4c_modbus_read_write_modbus_adu** modbus_read_request_packet) {

  plc4c_modbus_read_write_modbus_pdu* pdu =
      malloc(sizeof(plc4c_modbus_read_write_modbus_pdu));
  if (pdu == NULL) {
    return NO_MEMORY;
  }

  plc4c_driver_modbus_item* modbus_item = read_request_item->address;
  uint16_t address = modbus_item->address - 1;
  uint16_t num_modbus_bytes = modbus_item->num_elements * plc4c_modbus_read_write_modbus_data_type_get_data_type_size(modbus_item->datatype);
  switch (modbus_item->type) {
    case PLC4C_DRIVER_MODBUS_ADDRESS_TYPE_COIL: {
      pdu->_type =
          plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_read_coils_request;
      pdu->modbus_pdu_read_coils_request_starting_address = address;
      pdu->modbus_pdu_read_coils_request_quantity = num_modbus_bytes;
      break;
    }
    case PLC4C_DRIVER_MODBUS_ADDRESS_TYPE_DISCRETE_INPUT: {
      pdu->_type = plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_read_discrete_inputs_request;
      pdu->modbus_pdu_read_discrete_inputs_request_starting_address = address;
      pdu->modbus_pdu_read_discrete_inputs_request_quantity = MAX(2, num_modbus_bytes) / 2;
      break;
    }
    case PLC4C_DRIVER_MODBUS_ADDRESS_TYPE_INPUT_REGISTER: {
      pdu->_type = plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_read_input_registers_request;
      pdu-> modbus_pdu_read_input_registers_request_starting_address = address;
      pdu-> modbus_pdu_read_input_registers_request_quantity = MAX(2, num_modbus_bytes) / 2;
      break;
    }
    case PLC4C_DRIVER_MODBUS_ADDRESS_TYPE_HOLDING_REGISTER: {
      pdu->_type = plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_read_holding_registers_request;
      pdu-> modbus_pdu_read_holding_registers_request_starting_address = address;
      pdu-> modbus_pdu_read_holding_registers_request_quantity = MAX(2, num_modbus_bytes) / 2;
      break;
    }
    case PLC4C_DRIVER_MODBUS_ADDRESS_TYPE_EXTENDED_REGISTER: {
      // TODO: Currently not supported.
      return INTERNAL_ERROR;
    }
    default: {
      return INVALID_ADDRESS;
    }
  }

  plc4c_modbus_read_write_modbus_adu* adu =
      malloc(sizeof(plc4c_modbus_read_write_modbus_adu));
  if (adu == NULL) {
    return NO_MEMORY;
  }
  adu->_type = plc4c_modbus_read_write_modbus_adu_type_plc4c_modbus_read_write_modbus_tcp_adu;
  adu->modbus_tcp_adu_transaction_identifier = modbus_config->communication_id_counter++;
  adu->modbus_tcp_adu_unit_identifier = modbus_config->unit_identifier;
  adu->modbus_tcp_adu_pdu = pdu;
  *modbus_read_request_packet = adu;

  return OK;
}

plc4c_return_code plc4c_driver_modbus_create_modbus_write_request(
    plc4c_write_request* write_request,
    plc4c_modbus_read_write_modbus_adu** modbus_read_request_packet) {
  // TODO: Implement this ...

  return OK;
}
