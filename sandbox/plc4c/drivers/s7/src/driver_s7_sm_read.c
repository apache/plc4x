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
plc4c_return_code plc4c_driver_s7_parse_read_responce(
    plc4c_read_request_execution *execution, plc4c_s7_read_write_tpkt_packet* packet);

plc4c_return_code plc4c_driver_s7_sm_read_init(
    plc4c_connection* connection, plc4c_read_request_execution* execution) {
  
  plc4c_s7_read_write_tpkt_packet* packet;
  plc4c_return_code result;

  result = plc4c_driver_s7_create_s7_read_request(execution->read_request, &packet);
  if (result != OK) 
    return result;

  // Send the packet to the remote.
  result = plc4c_driver_s7_send_packet(connection, packet);
  if (result != OK) 
    return result;
  
  plc4c_driver_s7_destroy_s7_read_request(packet);
  return OK;
}

plc4c_return_code plc4c_driver_s7_sm_read_finished(
    plc4c_connection* connection, plc4c_read_request_execution* execution) {

  plc4c_s7_read_write_tpkt_packet* packet;
  plc4c_return_code result;

  // Read a response packet.If we haven't read enough to 
  // process a full message, just try again next time.
  result = plc4c_driver_s7_receive_packet(connection, &packet);
  if (result != OK) 
    return result;
  
  // Check the response. the number of items matches that of 
  //the request
  if (packet->payload->payload->parameter->_type != 
      plc4c_s7_read_write_s7_parameter_type_plc4c_s7_read_write_s7_parameter_read_var_response) 
    return INTERNAL_ERROR;
  
  if (packet->payload->payload->parameter->s7_parameter_read_var_response_num_items != 
      plc4c_utils_list_size(execution->read_request->items)) 
    return INTERNAL_ERROR;
  
  execution->read_response = malloc(sizeof(plc4c_read_response));
  if (execution->read_response == NULL) 
    return NO_MEMORY;
  execution->read_response->read_request = execution->read_request;

  result = plc4c_driver_s7_parse_read_responce(execution, packet);

  if (result != OK)
    return result;

  plc4c_driver_s7_destroy_receive_packet(packet);
  return OK;

}


plc4c_return_code plc4c_driver_s7_read_machine_function(
    plc4c_system_task* task) {

  plc4c_read_request_execution* execution;
  plc4c_connection* connection;
  plc4c_return_code result;

  execution = task->context;
  connection = task->connection;

  if ((!execution) || (!execution->read_request) || (!connection)) 
    return INTERNAL_ERROR;

  switch (task->state_id) {

    case PLC4C_DRIVER_S7_READ_INIT: 
      result = plc4c_driver_s7_sm_read_init(connection, execution);
      if (result != OK) 
        return result;
      task->state_id = PLC4C_DRIVER_S7_READ_FINISHED;
      break;
    
    case PLC4C_DRIVER_S7_READ_FINISHED: 
      result = plc4c_driver_s7_sm_read_finished(connection,execution);
      if (result == OK)
        task->completed = true;
      else if (result == UNFINISHED) 
        return OK;
      else
        return result;
      break;
    
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

void plc4c_driver_s7_free_read_request_item(plc4c_list_element *element){
  plc4c_item *item;
  item = element->value;
  plc4c_s7_read_write_s7_var_request_parameter_item *addr_item;
  addr_item = item->address;
  free(addr_item->s7_var_request_parameter_item_address_address );
  free(addr_item);
  free(item);
}


void plc4c_driver_s7_free_read_request(plc4c_read_request *request) {
  plc4c_utils_list_delete_elements(request->items,
      &plc4c_driver_s7_free_read_request_item);
  free(request->items);
  request->items = NULL;
  request->connection = NULL;
  // actual request free'd by caller
}

void plc4c_driver_s7_free_read_response_item(
    plc4c_list_element* element) {

  plc4c_response_value_item* value_item;
  value_item = element->value;
  plc4c_data_destroy(value_item->value);
  // dont free value_item->item its managed by the request not responce

  free(value_item);
}

void plc4c_driver_s7_free_read_response(plc4c_read_response* response) {
  // the request will be cleaned up elsewhere
  plc4c_utils_list_delete_elements(response->items,
                                   &plc4c_driver_s7_free_read_response_item);
  free(response->items);
}

plc4c_return_code plc4c_driver_s7_parse_read_responce( 
                                plc4c_read_request_execution* execution,
                                plc4c_s7_read_write_tpkt_packet* packet) {

  // TODO this function is too big and complex: modularise

  // Locals :(
  plc4c_s7_read_write_s7_message* s7_packet;
  plc4c_list_element* request_elements;
  plc4c_list_element* payload_elements;
  plc4c_item* request_item;
  plc4c_s7_read_write_s7_var_request_parameter_item* request_address;
  plc4c_s7_read_write_transport_size transport_size;
  plc4c_s7_read_write_s7_var_payload_data_item* payload_item;
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

  // Make a new list for holding the responce value items
  plc4c_utils_list_create(&execution->read_response->items);

  // Iterate over the request items and use the types to decode the
  // response items.
  request_elements = plc4c_utils_list_tail(execution->read_request->items);
  payload_elements = plc4c_utils_list_tail(packet->payload->payload->payload->s7_payload_read_var_response_items);
  
  while ((request_elements != NULL) && (payload_elements != NULL)) {
    
    request_item = request_elements->value;
    payload_item = payload_elements->value;

    // Get the protocol id for the current item from the corresponding
    // request item. Also get the number of elements, if it's an array.
    request_address = request_item->address;
    transport_size = request_address->s7_var_request_parameter_item_address_address->s7_address_any_transport_size;
    num_elements = request_address->s7_var_request_parameter_item_address_address->s7_address_any_number_of_elements;
    data_protocol_id = plc4c_s7_read_write_transport_size_get_data_protocol_id(transport_size);
    
    if (transport_size == plc4c_s7_read_write_transport_size_STRING) {
      // TODO: This needs to be changed to read arrays of strings.
      string_length = num_elements;
      num_elements = 1;
    } else {
      string_length = 0;
    }

    // Convert the linked list with uint8_t elements into an array of uint8_t.
    byte_array = plc4c_list_to_byte_array(payload_item->data);
    if (byte_array == NULL) 
      return INTERNAL_ERROR;

    // Create a new read-buffer for reading data from the uint8_t array.
    list_size = plc4c_utils_list_size(payload_item->data);
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

    response_value_item->item = request_item;

    // TODO: use other than just INTERNAL_ERROR on transport failure
    if (payload_item->return_code == plc4c_s7_read_write_data_transport_error_code_OK)
      response_value_item->response_code = PLC4C_RESPONSE_CODE_OK;
    else 
      response_value_item->response_code = PLC4C_RESPONSE_CODE_INTERNAL_ERROR;
    
    response_value_item->value = data_item;

    // Add the value-item to the list.
    plc4c_utils_list_insert_head_value(execution->read_response->items, response_value_item);
    free(read_buffer);
    free(byte_array);
    request_elements = request_elements->next;
    payload_elements = payload_elements->next;
  }
  return OK;
}
