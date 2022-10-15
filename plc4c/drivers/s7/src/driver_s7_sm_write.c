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

#include <ctype.h>
#include <plc4c/plc4c.h>
#include <plc4c/spi/types_private.h>
#include <stdlib.h>
#include <string.h>
#include "plc4c/driver_s7.h"
#include "plc4c/driver_s7_encode_decode.h"
#include "plc4c/driver_s7_packets.h"
#include "cotp_protocol_class.h"
#include "tpkt_packet.h"

enum plc4c_driver_s7_write_states {
  PLC4C_DRIVER_S7_WRITE_INIT,
  PLC4C_DRIVER_S7_WRITE_FINISHED
};

// Forward declaration of helper function to stop PLC4C_DRIVER_S7_WRITE_FINISHED
// state become too big, TODO: move to some header or inline
plc4c_return_code plc4c_driver_s7_parse_write_response(
     plc4c_write_request_execution* execution, plc4c_s7_read_write_tpkt_packet* packet);


plc4c_return_code plc4c_driver_s7_sm_write_init(
    plc4c_connection* connection, plc4c_write_request_execution* execution) {
  
  plc4c_s7_read_write_tpkt_packet* packet;
  plc4c_return_code result;

  result = plc4c_driver_s7_create_s7_write_request(execution->write_request, &packet);
  if (result != OK) {
    return result;
  }

  // Send the packet to the remote. 
  result = plc4c_driver_s7_send_packet(connection, packet);
  if (result != OK) {
    return result;
  }
  plc4c_driver_s7_destroy_s7_write_request(packet);

  return OK;

}

plc4c_return_code plc4c_driver_s7_sm_write_finished(
    plc4c_connection* connection, plc4c_write_request_execution* execution) {
  
  plc4c_s7_read_write_tpkt_packet* packet;
  plc4c_return_code result;

  // Read a response packet. If we haven't read enough to process
  // a full message, just try again next time.
  result = plc4c_driver_s7_receive_packet(connection, &packet);
  if (result != OK) 
    return result;

  // Check the response packet s7 parameter is of correct type and 
  // number of parameter items match number of request items 
  if (packet->payload->payload->parameter->_type != 
      plc4c_s7_read_write_s7_parameter_type_plc4c_s7_read_write_s7_parameter_write_var_response)
    return INTERNAL_ERROR;

  if (packet->payload->payload->parameter->s7_parameter_read_var_response_num_items !=
      plc4c_utils_list_size(execution->write_request->items)) 
    return INTERNAL_ERROR;
  
  // Make a new response item and bind to the execution, also bind
  // the request to the response (useful of parsing)
  execution->write_response = malloc(sizeof(plc4c_write_response));
  if (execution->write_response == NULL) 
    return NO_MEMORY;
  execution->write_response->write_request = execution->write_request;
  
  // Set the write response status from the s7 payload response status
  result = plc4c_driver_s7_parse_write_response(execution, packet);
  
  if (result != OK)
    return result;

  plc4c_driver_s7_destroy_receive_packet(packet);
  return OK;
}

plc4c_return_code plc4c_driver_s7_write_machine_function(
    plc4c_system_task* task) {

  plc4c_write_request_execution* execution;
  plc4c_connection* connection;
  plc4c_return_code result;

  execution = task->context;
  connection = task->connection;

  if ((!execution) || (!execution->write_request) || (!connection))
    return INTERNAL_ERROR;

  switch (task->state_id) {

    case PLC4C_DRIVER_S7_WRITE_INIT: 
      result = plc4c_driver_s7_sm_write_init(connection, execution);
      if (result == OK)
        task->state_id = PLC4C_DRIVER_S7_WRITE_FINISHED;
      else
        return result;
      break;
    
    case PLC4C_DRIVER_S7_WRITE_FINISHED: 

      result = plc4c_driver_s7_sm_write_finished(connection, execution);
      if (result == OK)
        task->completed = true;
      else if (result == UNFINISHED) 
        return OK;
      else
        return result;
    
  }
  return OK;
}

plc4c_return_code plc4c_driver_s7_write_function(
    plc4c_write_request_execution* write_request_execution,
    plc4c_system_task** task) {

  plc4c_system_task* new_task = malloc(sizeof(plc4c_system_task));
  if(new_task == NULL)
    return NO_MEMORY;
  
  new_task->state_id = PLC4C_DRIVER_S7_WRITE_INIT;
  new_task->state_machine_function = &plc4c_driver_s7_write_machine_function;
  new_task->completed = false;
  new_task->context = write_request_execution;
  new_task->connection = plc4c_write_request_get_connection(write_request_execution->write_request);
  *task = new_task;
  return OK;
}


// TDOO: remove hacked include but think doing so needes a driver
// callback for request items due to malloc on: 
//  s7_var_request_parameter_item_address_address
// either way it dosnt belong in this file.
void plc4c_driver_s7_free_write_request_item(plc4c_list_element *element) {
  plc4c_request_value_item *item;
  item = element->value;

  plc4c_s7_read_write_s7_var_request_parameter_item *addr_item;
  addr_item = item->item->address;
  free(addr_item->s7_var_request_parameter_item_address_address );
  free(addr_item);

  free(item->item);
  // todo : require call expcitly or add here
  plc4c_data_destroy(item->value);
  free(item);
}

void plc4c_driver_s7_free_write_request(plc4c_write_request *request) {
  plc4c_utils_list_delete_elements(request->items,
      plc4c_driver_s7_free_write_request_item);
  free(request->items);
  // actual request free'd by caller
}

void plc4c_driver_s7_free_write_response_item(
    plc4c_list_element* write_item_element) {
  
  plc4c_response_item* response_item;
  response_item = write_item_element->value;
  // dont free response_item-item->item its managed by the request not response
  free(response_item);
}

void plc4c_driver_s7_free_write_response(plc4c_write_response* response) {
  
    plc4c_utils_list_delete_elements(response->response_items,
      plc4c_driver_s7_free_write_response_item);
    free(response->response_items);
}


plc4c_return_code plc4c_driver_s7_parse_write_response(
    plc4c_write_request_execution *execution, 
    plc4c_s7_read_write_tpkt_packet* packet) {

  // Locals first 4 are just used for walking two lists
  plc4c_list_element* request_elements;
  plc4c_request_value_item* request_item;
  plc4c_list_element* s7_payload_elements;
  plc4c_s7_read_write_s7_var_payload_status_item* s7_payload_item;
  plc4c_response_item* response_item;
  plc4c_return_code result;
  
  // Make the unfilled response a list
  plc4c_utils_list_create(&execution->write_response->response_items);

	// Iterate over the request items setting return codes as needed.
	request_elements = plc4c_utils_list_tail(execution->write_request->items);
	s7_payload_elements = plc4c_utils_list_tail(packet->payload->payload->payload->s7_payload_write_var_response_items);

	while ((request_elements != NULL) && (s7_payload_elements != NULL)) {
		
    request_item = request_elements->value;
		s7_payload_item = s7_payload_elements->value;

    // Make a new response item, bind the related request item
		response_item = malloc(sizeof(plc4c_response_item)); 
		if (response_item == NULL)
		  return NO_MEMORY;
    response_item->item = request_item->item;

    // TODO: better map transport error codes to response error codes
    if (s7_payload_item->return_code == plc4c_s7_read_write_data_transport_error_code_OK)
      response_item->response_code = PLC4C_RESPONSE_CODE_OK;
    else
      response_item->response_code = PLC4C_RESPONSE_CODE_INTERNAL_ERROR;

		plc4c_utils_list_insert_head_value(execution->write_response->response_items, response_item);

		request_elements = request_elements->next;
		s7_payload_elements = s7_payload_elements->next;
	}
  return OK;
}