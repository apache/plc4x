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
#include "plc4c/driver_plc4x.h"
#include "plc4c/driver_plc4x_packets.h"
#include "plc4x_message.h"

enum plc4c_driver_plc4x_read_states {
  PLC4C_DRIVER_PLC4X_READ_INIT,
  PLC4C_DRIVER_PLC4X_READ_FINISHED
};

// Forward declaration of helper function to stop PLC4C_DRIVER_PLC4X_READ_FINISHED
// state become too big, TODO: move to some header or inline
plc4c_return_code plc4c_driver_plc4x_parse_read_response(
    plc4c_read_request_execution *execution, plc4c_plc4x_read_write_plc4x_message* packet);

plc4c_return_code plc4c_driver_plc4x_sm_read_init(
    plc4c_connection* connection, plc4c_read_request_execution* execution) {
  
  plc4c_plc4x_read_write_plc4x_message* packet;
  plc4c_return_code result;

  result = plc4c_driver_plc4x_create_plc4x_read_request(execution->read_request, &packet);
  if (result != OK) {
    return result;
  }

  // Send the packet to the remote.
  result = plc4c_driver_plc4x_send_packet(connection, packet);
  if (result != OK) {
    return result;
  }
  
  plc4c_driver_plc4x_destroy_plc4x_read_request(packet);
  return OK;
}

plc4c_return_code plc4c_driver_plc4x_sm_read_finished(
    plc4c_connection* connection, plc4c_read_request_execution* execution) {

  plc4c_plc4x_read_write_plc4x_message* packet;
  plc4c_return_code result;

  // Read a response packet.If we haven't read enough to 
  // process a full message, just try again next time.
  result = plc4c_driver_plc4x_receive_packet(connection, &packet);
  if (result != OK)
    return result;
  
  // Check the response. the number of items matches that of 
  //the request
  if (packet->_type !=
      plc4c_plc4x_read_write_plc4x_message_type_plc4c_plc4x_read_write_plc4x_read_response)
    return INTERNAL_ERROR;

  // TODO: Check the connection id and the request id both match
  //if (packet->plc4x_read_response_connection_id == connection.)
  
  if (plc4c_utils_list_size(packet->plc4x_read_request_fields) !=
      plc4c_utils_list_size(execution->read_request->items))
    return INTERNAL_ERROR;
  
  execution->read_response = malloc(sizeof(plc4c_read_response));
  if (execution->read_response == NULL)
    return NO_MEMORY;
  execution->read_response->read_request = execution->read_request;

  result = plc4c_driver_plc4x_parse_read_response(execution, packet);

  if (result != OK)
    return result;

  plc4c_driver_plc4x_destroy_packet(packet);
  return OK;

}


plc4c_return_code plc4c_driver_plc4x_read_machine_function(
    plc4c_system_task* task) {

  plc4c_read_request_execution* execution;
  plc4c_connection* connection;
  plc4c_return_code result;

  execution = task->context;
  connection = task->connection;

  if ((!execution) || (!execution->read_request) || (!connection)) 
    return INTERNAL_ERROR;

  switch (task->state_id) {

    case PLC4C_DRIVER_PLC4X_READ_INIT:
      result = plc4c_driver_plc4x_sm_read_init(connection, execution);
      if (result != OK) 
        return result;
      task->state_id = PLC4C_DRIVER_PLC4X_READ_FINISHED;
      break;
    
    case PLC4C_DRIVER_PLC4X_READ_FINISHED:
      result = plc4c_driver_plc4x_sm_read_finished(connection,execution);
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

plc4c_return_code plc4c_driver_plc4x_read_function(
    plc4c_read_request_execution* read_request_execution,
    plc4c_system_task** task) {
  plc4c_system_task* new_task = malloc(sizeof(plc4c_system_task));
  if(new_task == NULL) {
    return NO_MEMORY;
  }
  new_task->state_id = PLC4C_DRIVER_PLC4X_READ_INIT;
  new_task->state_machine_function = &plc4c_driver_plc4x_read_machine_function;
  new_task->completed = false;
  new_task->context = read_request_execution;
  new_task->connection = read_request_execution->read_request->connection;
  *task = new_task;
  return OK;
}

void plc4c_driver_plc4x_free_read_request_item(plc4c_list_element *element){
  plc4c_item *item;
  item = element->value;
  plc4c_plc4x_read_write_plc4x_field_request *addr_item;
  addr_item = item->address;
  free(addr_item->field->name);
  free(addr_item->field->field_query);
  free(addr_item->field);
  free(item);
}


void plc4c_driver_plc4x_free_read_request(plc4c_read_request *request) {
  plc4c_utils_list_delete_elements(request->items,
      &plc4c_driver_plc4x_free_read_request_item);
  free(request->items);
  request->items = NULL;
  request->connection = NULL;
  // actual request freed by caller
}

void plc4c_driver_plc4x_free_read_response_item(
    plc4c_list_element* element) {

  plc4c_response_value_item* value_item;
  value_item = element->value;
  plc4c_data_destroy(value_item->value);
  // don't free value_item->item its managed by the request not response

  free(value_item);
}

void plc4c_driver_plc4x_free_read_response(plc4c_read_response* response) {
  // the request will be cleaned up elsewhere
  plc4c_utils_list_delete_elements(response->items,
                                   &plc4c_driver_plc4x_free_read_response_item);
  free(response->items);
}

plc4c_return_code plc4c_driver_plc4x_parse_read_response(
                                plc4c_read_request_execution* execution,
                                plc4c_plc4x_read_write_plc4x_message* packet) {

  // Locals :(
  plc4c_list_element* request_elements;
  plc4c_list_element* payload_elements;
  plc4c_item* request_item;
  plc4c_plc4x_read_write_plc4x_field_value_response * payload_item;
  plc4c_response_value_item* response_value_item;

  // Make a new list for holding the response value items
  plc4c_utils_list_create(&execution->read_response->items);

  // Iterate over the request items and use the types to decode the
  // response items.
  request_elements = plc4c_utils_list_tail(execution->read_request->items);
  payload_elements = plc4c_utils_list_tail(packet->plc4x_read_response_fields);
  
  while ((request_elements != NULL) && (payload_elements != NULL)) {
    request_item = request_elements->value;
    payload_item = payload_elements->value;

    // Create a new response value-item
    response_value_item = malloc(sizeof(plc4c_response_value_item));
    if (response_value_item == NULL) 
      return NO_MEMORY;

    response_value_item->item = request_item;

    if (payload_item->response_code == plc4c_plc4x_read_write_plc4x_response_code_OK) {
      response_value_item->response_code = PLC4C_RESPONSE_CODE_OK;
    } else {
      response_value_item->response_code = PLC4C_RESPONSE_CODE_INTERNAL_ERROR;
    }
    
    response_value_item->value = payload_item->value;

    // Add the value-item to the list.
    plc4c_utils_list_insert_head_value(execution->read_response->items, response_value_item);
    request_elements = request_elements->next;
    payload_elements = payload_elements->next;
  }

  return OK;
}
