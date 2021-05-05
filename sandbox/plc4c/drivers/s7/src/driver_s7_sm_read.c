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
#include "plc4c/driver_s7.h"
#include "plc4c/driver_s7_encode_decode.h"
#include "plc4c/driver_s7_packets.h"
#include "cotp_protocol_class.h"
#include "tpkt_packet.h"
#include "data_item.h"

enum plc4c_driver_s7_read_states {
  PLC4C_DRIVER_S7_READ_INIT,
  PLC4C_DRIVER_S7_READ_FINISHED
};

// Forward declaration of helper function to stop PLC4C_DRIVER_S7_READ_FINISHED
// state become too big, TODO: move to some header or inline
plc4c_return_code plc4c_driver_s7_parse_read_responce(plc4c_read_request* request,
    plc4c_read_response* response, plc4c_s7_read_write_tpkt_packet* packet);

plc4c_return_code plc4c_driver_s7_read_machine_function(
    plc4c_system_task* task) {

  plc4c_read_request_execution* read_request_execution;
  plc4c_read_request* read_request;
  plc4c_connection* connection;
  plc4c_s7_read_write_tpkt_packet* s7_read_packet;
  plc4c_return_code return_code;

  read_request_execution = task->context;
  if (read_request_execution == NULL) {
    return INTERNAL_ERROR;
  }
  read_request = read_request_execution->read_request;
  if (read_request == NULL) {
    return INTERNAL_ERROR;
  }
  connection = task->connection;
  if (connection == NULL) {
    return INTERNAL_ERROR;
  }

  switch (task->state_id) {
    case PLC4C_DRIVER_S7_READ_INIT: {
      return_code = plc4c_driver_s7_create_s7_read_request(read_request, &s7_read_packet);
      if (return_code != OK) {
        return return_code;
      }

      // Send the packet to the remote.
      return_code = plc4c_driver_s7_send_packet(connection, s7_read_packet);
      if (return_code != OK) {
        return return_code;
      }

      task->state_id = PLC4C_DRIVER_S7_READ_FINISHED;
      break;
    }
    case PLC4C_DRIVER_S7_READ_FINISHED: {
      
      plc4c_s7_read_write_s7_parameter* parameter;
      plc4c_read_response* read_response;

      // Read a response packet.
      return_code = plc4c_driver_s7_receive_packet(connection, &s7_read_packet);
      // If we haven't read enough to process a full message, just try again
      // next time.
      if (return_code == UNFINISHED) {
        return OK;
      } else if (return_code != OK) {
        return return_code;
      }

      // Check the response.
      parameter = s7_read_packet->payload->payload->parameter;
      if (parameter->_type != plc4c_s7_read_write_s7_parameter_type_plc4c_s7_read_write_s7_parameter_read_var_response) {
        return INTERNAL_ERROR;
      }
      // Check if the number of items matches that of the request
      // (Otherwise we won't know how to interpret the items)
      if (parameter->s7_parameter_read_var_response_num_items != plc4c_utils_list_size(read_request->items)) {
        return INTERNAL_ERROR;
      }

      read_response = malloc(sizeof(plc4c_read_response));
      if (read_response == NULL) {
        return NO_MEMORY;
      }
      read_response->read_request = read_request;
      read_request_execution->read_response = read_response;
      plc4c_utils_list_create(&read_response->items);

      return_code = plc4c_driver_s7_parse_read_responce(read_request, 
          read_response, s7_read_packet);
      if (return_code != OK)
        return return_code;

      // TODO: Return the results to the API ...
      task->completed = true;
      break;
    }
  }
  return OK;
}

plc4c_return_code plc4c_driver_s7_read_function(
    plc4c_read_request_execution* read_request_execution,
    plc4c_system_task** task) {
  plc4c_system_task* new_task = malloc(sizeof(plc4c_system_task));
  if(new_task == NULL) {
    return NO_MEMORY;
  }
  new_task->state_id = PLC4C_DRIVER_S7_READ_INIT;
  new_task->state_machine_function = &plc4c_driver_s7_read_machine_function;
  new_task->completed = false;
  new_task->context = read_request_execution;
  new_task->connection = read_request_execution->read_request->connection;
  *task = new_task;
  return OK;
}

void plc4c_driver_s7_free_read_response_item(
    plc4c_list_element* read_item_element) {
  plc4c_response_value_item* value_item =
      (plc4c_response_value_item*)read_item_element->value;
  plc4c_data_destroy(value_item->value);
  value_item->value = NULL;
}

void plc4c_driver_s7_free_read_response(plc4c_read_response* response) {
  // the request will be cleaned up elsewhere
  plc4c_utils_list_delete_elements(response->items,
                                   &plc4c_driver_s7_free_read_response_item);
}

plc4c_return_code plc4c_driver_s7_parse_read_responce( 
                                plc4c_read_request* request, 
                                plc4c_read_response* response,
                                plc4c_s7_read_write_tpkt_packet* packet) {

  // Locals
  plc4c_s7_read_write_s7_message* s7_packet;
  plc4c_list_element* request_list_element;
  plc4c_list_element* response_list_element;
  plc4c_item* cur_request_item;
  plc4c_s7_read_write_s7_var_request_parameter_item* s7_address;
  plc4c_s7_read_write_transport_size transport_size;

  plc4c_s7_read_write_s7_var_payload_data_item* s7_payload_data_item;
  plc4c_spi_read_buffer* read_buffer;
  plc4c_data* data_item;
  plc4c_response_value_item* response_value_item;

  char* data_protocol_id;
  uint16_t num_elements;
  int32_t string_length;
  uint8_t* byte_array;
  size_t list_size;
  size_t idx;
  enum plc4c_return_code result;

  // Iterate over the request items and use the types to decode the
  // response items.
  s7_packet = packet->payload->payload;
  request_list_element = plc4c_utils_list_tail(request->items);
  response_list_element = plc4c_utils_list_tail(s7_packet->payload->s7_payload_read_var_response_items);
  
  while ((request_list_element != NULL) && (response_list_element != NULL)) {
    
    cur_request_item = request_list_element->value;

    // Get the protocol id for the current item from the corresponding
    // request item. Also get the number of elements, if it's an array.
    s7_address = cur_request_item->address;
    transport_size = s7_address->s7_var_request_parameter_item_address_address->s7_address_any_transport_size;
    data_protocol_id = plc4c_s7_read_write_transport_size_get_data_protocol_id(transport_size);
    num_elements = s7_address->s7_var_request_parameter_item_address_address->s7_address_any_number_of_elements;
    string_length = 0;
    if (transport_size == plc4c_s7_read_write_transport_size_STRING) {
      // TODO: This needs to be changed to read arrays of strings.
      string_length = num_elements;
      num_elements = 1;
    }

    // Convert the linked list with uint8_t elements into an array of uint8_t.
    s7_payload_data_item = response_list_element->value;
    byte_array = plc4c_list_to_byte_array(s7_payload_data_item->data);
    if (byte_array == NULL) 
      return INTERNAL_ERROR;

    // Create a new read-buffer for reading data from the uint8_t array.
    list_size = plc4c_utils_list_size(s7_payload_data_item->data);
    result = plc4c_spi_read_buffer_create(byte_array, list_size, &read_buffer);
    if (result != OK) 
      return result;
  
    // TODO: check if elements > 1 is always a list
    if (num_elements > 1) {
      plc4c_list *all_list;
      plc4c_data* all_data_item;
      plc4c_utils_list_create(&all_list);
      all_data_item = plc4c_data_create_list_data(*all_list);
      for (idx = 0; idx < num_elements ; idx++) {
        plc4c_s7_read_write_data_item_parse(read_buffer, data_protocol_id, string_length, &data_item);
        plc4c_utils_list_insert_head_value(&all_data_item->data.list_value, (void*)data_item);
      }
      data_item = all_data_item;
    } else {
      plc4c_s7_read_write_data_item_parse(read_buffer, data_protocol_id, string_length, &data_item);
    }

    // Create a new response value-item
    response_value_item = malloc(sizeof(plc4c_response_value_item));
    if (response_value_item == NULL) 
      return NO_MEMORY;

    response_value_item->item = cur_request_item;

    // TODO: use other than just INTERNAL_ERROR on transport failure
    if (s7_payload_data_item->return_code == plc4c_s7_read_write_data_transport_error_code_OK)
      response_value_item->response_code = PLC4C_RESPONSE_CODE_OK;
    else 
      response_value_item->response_code = PLC4C_RESPONSE_CODE_INTERNAL_ERROR;
    
    response_value_item->value = data_item;

    // Add the value-item to the list.
    plc4c_utils_list_insert_head_value(response->items, response_value_item);

    request_list_element = request_list_element->next;
    response_list_element = response_list_element->next;
  }
  return OK;
}