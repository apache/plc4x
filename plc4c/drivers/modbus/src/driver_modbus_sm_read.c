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
#include "plc4c/driver_modbus_packets.h"
#include "modbus_adu.h"
#include "data_item.h"


/*
 * A Modbus register is 16 bits wide. data_item_parse encodes one value at its natural width and
 * knows nothing about that, so a value narrower than a register is padded here when it stands
 * alone and packed when several of them follow one another. Which half of the register a padded
 * value sits in depends on the byte order.
 */
static uint8_t plc4c_driver_modbus_padding_bits(
    plc4c_modbus_read_write_modbus_data_type data_type) {
  switch (data_type) {
    case plc4c_modbus_read_write_modbus_data_type_BOOL:
      return 15;
    case plc4c_modbus_read_write_modbus_data_type_BYTE:
    case plc4c_modbus_read_write_modbus_data_type_SINT:
    case plc4c_modbus_read_write_modbus_data_type_USINT:
      return 8;
    default:
      return 0;
  }
}

/* A lone BOOL sits between seven bits and eight when the byte order is little endian. */
static uint8_t plc4c_driver_modbus_leading_padding_bits(
    plc4c_modbus_read_write_modbus_data_type data_type) {
  return (data_type == plc4c_modbus_read_write_modbus_data_type_BOOL) ? 7 : 0;
}

static plc4c_return_code plc4c_driver_modbus_parse_registers(
    plc4c_spi_read_buffer* read_buffer,
    plc4c_modbus_read_write_modbus_data_type data_type,
    uint16_t num_elements, bool big_endian, uint16_t string_length,
    plc4c_data** data_item) {
  plc4c_return_code result;
  uint16_t _ignored = 0;
  if (num_elements == 1) {
    uint8_t padding = plc4c_driver_modbus_padding_bits(data_type);
    if (padding == 0) {
      return plc4c_modbus_read_write_data_item_parse(
          plc4x_spi_context_background(), read_buffer, data_type, string_length, data_item);
    }
    if (big_endian) {
      result = plc4c_spi_read_unsigned_short(read_buffer, padding, &_ignored);
      if (result != OK) {
        return result;
      }
      return plc4c_modbus_read_write_data_item_parse(
          plc4x_spi_context_background(), read_buffer, data_type, string_length, data_item);
    }
    uint8_t leading = plc4c_driver_modbus_leading_padding_bits(data_type);
    if (leading > 0) {
      result = plc4c_spi_read_unsigned_short(read_buffer, leading, &_ignored);
      if (result != OK) {
        return result;
      }
    }
    result = plc4c_modbus_read_write_data_item_parse(
        plc4x_spi_context_background(), read_buffer, data_type, string_length, data_item);
    if (result != OK) {
      return result;
    }
    return plc4c_spi_read_unsigned_short(read_buffer, padding - leading, &_ignored);
  }

  /* Several values are packed without padding between them. */
  plc4c_list* list;
  plc4c_utils_list_create(&list);
  for (uint16_t i = 0; i < num_elements; i++) {
    plc4c_data* element;
    result = plc4c_modbus_read_write_data_item_parse(
        plc4x_spi_context_background(), read_buffer, data_type, string_length, &element);
    if (result != OK) {
      return result;
    }
    plc4c_utils_list_insert_tail_value(list, element);
  }
  *data_item = plc4c_data_create_list_data(list);
  return OK;
}

enum plc4c_driver_modbus_read_states {
  PLC4C_DRIVER_MODBUS_READ_INIT,
  PLC4C_DRIVER_MODBUS_READ_SEND_ITEM_REQUEST,
  PLC4C_DRIVER_MODBUS_READ_HANDLE_ITEM_RESPONSE,
  PLC4C_DRIVER_MODBUS_READ_FINISHED
};

plc4c_return_code plc4c_driver_modbus_read_machine_function(
    plc4c_system_task* task) {
  plc4c_read_request_execution* read_request_execution = task->context;
  if (read_request_execution == NULL) {
    return INTERNAL_ERROR;
  }
  plc4c_read_request* read_request = read_request_execution->read_request;
  if (read_request == NULL) {
    return INTERNAL_ERROR;
  }
  plc4c_connection* connection = read_request->connection;
  if (connection == NULL) {
    return INTERNAL_ERROR;
  }
  plc4c_driver_modbus_config* modbus_config = connection->configuration;

  switch (task->state_id) {
    // First set the current item to the first item in the list (tail)
    case PLC4C_DRIVER_MODBUS_READ_INIT: {
      read_request_execution->cur_item = plc4c_utils_list_tail(read_request->items);

      // Create an empty read-response and attach that to the execution.
      plc4c_read_response* read_response = malloc(sizeof(plc4c_read_response));
      if(read_response == NULL) {
        return NO_MEMORY;
      }
      read_response->read_request = read_request;
      read_request_execution->read_response = read_response;
      plc4c_utils_list_create(&(read_response->items));
      read_request_execution->read_response = read_response;

      task->state_id = PLC4C_DRIVER_MODBUS_READ_SEND_ITEM_REQUEST;
      break;
    }

    case PLC4C_DRIVER_MODBUS_READ_SEND_ITEM_REQUEST: {
      plc4c_modbus_read_write_modbus_adu* modbus_read_request_packet;
      plc4c_return_code return_code =
          plc4c_driver_modbus_create_modbus_read_request(
              modbus_config,
              read_request_execution->cur_item->value,
              &modbus_read_request_packet);
      if (return_code != OK) {
        return return_code;
      }

      // Send the packet to the remote.
      return_code = plc4c_driver_modbus_send_packet(
          connection, modbus_read_request_packet);
      if (return_code != OK) {
        return return_code;
      }

      task->state_id = PLC4C_DRIVER_MODBUS_READ_HANDLE_ITEM_RESPONSE;
      break;
    }
    case PLC4C_DRIVER_MODBUS_READ_HANDLE_ITEM_RESPONSE: {
      // Read a response packet.
      plc4c_modbus_read_write_modbus_adu* modbus_read_response_packet;
      plc4c_return_code return_code = plc4c_driver_modbus_receive_packet(
          connection, &modbus_read_response_packet);
      // If we haven't read enough to process a full message, just try again
      // next time.
      if (return_code == UNFINISHED) {
        return OK;
      } else if (return_code != OK) {
        return return_code;
      }

      // Check if the response has the correct type.
      plc4c_item* read_request_item = read_request_execution->cur_item->value;
      plc4c_driver_modbus_item* modbus_item = read_request_item->address;
      plc4c_list* response_value;
      switch (modbus_item->type) {
        case PLC4C_DRIVER_MODBUS_ADDRESS_TYPE_COIL: {
          if(modbus_read_response_packet->modbus_tcp_adu_pdu->_type !=
              plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_read_coils_response) {
            return INTERNAL_ERROR;
          }
          response_value = modbus_read_response_packet->modbus_tcp_adu_pdu->modbus_pdu_read_coils_response_value;
          break;
        }
        case PLC4C_DRIVER_MODBUS_ADDRESS_TYPE_DISCRETE_INPUT: {
          if(modbus_read_response_packet->modbus_tcp_adu_pdu->_type !=
              plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_read_discrete_inputs_response) {
            return INTERNAL_ERROR;
          }
          response_value = modbus_read_response_packet->modbus_tcp_adu_pdu->modbus_pdu_read_discrete_inputs_response_value;
          break;
        }
        case PLC4C_DRIVER_MODBUS_ADDRESS_TYPE_INPUT_REGISTER: {
          if(modbus_read_response_packet->modbus_tcp_adu_pdu->_type !=
              plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_read_input_registers_response) {
            return INTERNAL_ERROR;
          }
          response_value = modbus_read_response_packet->modbus_tcp_adu_pdu->modbus_pdu_read_input_registers_response_value;
          break;
        }
        case PLC4C_DRIVER_MODBUS_ADDRESS_TYPE_HOLDING_REGISTER: {
          if(modbus_read_response_packet->modbus_tcp_adu_pdu->_type !=
              plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_read_holding_registers_response) {
            return INTERNAL_ERROR;
          }
          response_value = modbus_read_response_packet->modbus_tcp_adu_pdu->modbus_pdu_read_holding_registers_response_value;
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

      // Convert the list into an array.
      uint8_t* byte_array = plc4c_list_to_byte_array(response_value);
      if(byte_array == NULL) {
        return INTERNAL_ERROR;
      }

      // Create a new read-buffer for reading data from the uint8_t array.
      plc4c_spi_read_buffer* read_buffer;
      plc4c_return_code result = plc4c_spi_read_buffer_create(byte_array, plc4c_utils_list_size(response_value), &read_buffer);
      if(result != OK) {
        return result;
      }

      // Decode the items in the response ...
      plc4c_data* data_item;
      // The C tag has no notion of a string length yet, so 1 is passed - the value it has for every
      // data type that is not a string, which leaves the behaviour of all other types unchanged.
      // Padding and packing of sub-register values is handled by the helper above.
      result = plc4c_driver_modbus_parse_registers(read_buffer, modbus_item->datatype,
          modbus_item->num_elements, true, 1, &data_item);
      if(result != OK) {
        return result;
      }

      // Create a new response value-item
      plc4c_response_value_item* response_value_item = malloc(sizeof(plc4c_response_value_item));
      if(response_value_item == NULL) {
        return NO_MEMORY;
      }
      response_value_item->item = read_request_execution->cur_item->value;
      response_value_item->response_code = PLC4C_RESPONSE_CODE_OK;
      response_value_item->value = data_item;

      // Add the value-item to the list.
      plc4c_utils_list_insert_head_value(
          read_request_execution->read_response->items, response_value_item);

      // If there are more items to read, continue reading the next one.
      // Otherwise finish.
      if (read_request_execution->cur_item->next != NULL) {
        read_request_execution->cur_item = read_request_execution->cur_item->next;
        task->state_id = PLC4C_DRIVER_MODBUS_READ_SEND_ITEM_REQUEST;
      } else {
        task->state_id = PLC4C_DRIVER_MODBUS_READ_FINISHED;
      }
      break;
    }

    case PLC4C_DRIVER_MODBUS_READ_FINISHED: {

      // TODO: Return the results to the API ...
      task->completed = true;
      break;
    }
  }
  return OK;
}

plc4c_return_code plc4c_driver_modbus_read_function(
    plc4c_read_request_execution* read_request_execution,
    plc4c_system_task** task) {
  plc4c_system_task* new_task = malloc(sizeof(plc4c_system_task));
  if(new_task == NULL) {
    return NO_MEMORY;
  }
  new_task->state_id = PLC4C_DRIVER_MODBUS_READ_INIT;
  new_task->state_machine_function = &plc4c_driver_modbus_read_machine_function;
  new_task->completed = false;
  new_task->context = read_request_execution;
  new_task->connection = read_request_execution->read_request->connection;
  *task = new_task;
  return OK;
}

void plc4c_driver_modbus_free_read_request_item(
    plc4c_list_element* read_item_element) {
  plc4c_item* value_item =
      (plc4c_item*)read_item_element->value;
  plc4c_driver_modbus_item* modbus_item = (plc4c_driver_modbus_item*) value_item->address;
  free(modbus_item);
  value_item->address = NULL;
}

void plc4c_driver_modbus_free_read_request(plc4c_read_request* request) {
  // the request will be cleaned up elsewhere
  plc4c_utils_list_delete_elements(request->items,
                                   plc4c_driver_modbus_free_read_request_item);
}

void plc4c_driver_modbus_free_read_response_item(
    plc4c_list_element* read_item_element) {
  plc4c_response_value_item* value_item =
      (plc4c_response_value_item*)read_item_element->value;
  plc4c_data_destroy(value_item->value);
  value_item->value = NULL;
}

void plc4c_driver_modbus_free_read_response(plc4c_read_response* response) {
  // the request will be cleaned up elsewhere
  plc4c_utils_list_delete_elements(response->items,
                                   plc4c_driver_modbus_free_read_response_item);
}
