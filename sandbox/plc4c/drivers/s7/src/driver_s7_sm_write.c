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
plc4c_return_code plc4c_driver_s7_parse_write_responce(plc4c_write_request* request, 
    plc4c_write_response* response, plc4c_s7_read_write_tpkt_packet* packet);

plc4c_return_code plc4c_driver_s7_write_machine_function(
    plc4c_system_task* task) {

  plc4c_write_request_execution* write_request_execution;
  plc4c_write_request* write_request;
  plc4c_connection* connection;
  plc4c_s7_read_write_tpkt_packet* write_packet;
  plc4c_return_code return_code;

  write_request_execution = task->context;
  if (write_request_execution == NULL) {
    return INTERNAL_ERROR;
  }
  write_request = write_request_execution->write_request;
  if (write_request == NULL) {
    return INTERNAL_ERROR;
  }
  connection = task->connection;
  if (connection == NULL) {
    return INTERNAL_ERROR;
  }

  switch (task->state_id) {
    case PLC4C_DRIVER_S7_WRITE_INIT: {
      return_code = plc4c_driver_s7_create_s7_write_request(write_request, &write_packet);
      if (return_code != OK) {
        return return_code;
      }

      // Send the packet to the remote.
      return_code = plc4c_driver_s7_send_packet(connection, write_packet);
      if (return_code != OK) {
        return return_code;
      }

      task->state_id = PLC4C_DRIVER_S7_WRITE_FINISHED;
      break;
    }
    case PLC4C_DRIVER_S7_WRITE_FINISHED: {

      plc4c_s7_read_write_s7_message* s7_packet;
      plc4c_write_response* write_response;

      // Read a response packet.
      return_code = plc4c_driver_s7_receive_packet(connection, &write_packet);
      // If we haven't read enough to process a full message, just try again
      // next time.
      if (return_code == UNFINISHED) {
        return OK;
      } else if (return_code != OK) {
        return return_code;
      }

      // Check the response
      s7_packet = write_packet->payload->payload;
      if (s7_packet->parameter->_type != plc4c_s7_read_write_s7_parameter_type_plc4c_s7_read_write_s7_parameter_write_var_response) {
        return INTERNAL_ERROR;
      }
      // Check if the number of items matches that of the request
      // (Otherwise we won't know how to interpret the items)
      if (s7_packet->parameter->s7_parameter_read_var_response_num_items != plc4c_utils_list_size(write_request->items)) {
        return INTERNAL_ERROR;
      }

      write_response = malloc(sizeof(plc4c_write_response));
      if (write_response == NULL) {
        return NO_MEMORY;
      }
      write_response->write_request = write_request;
      write_request_execution->write_response = write_response;
      plc4c_utils_list_create(&write_response->response_items);

      return_code = plc4c_driver_s7_parse_write_responce(write_request, write_response, write_packet);
      if (return_code != OK)
        return return_code;

      // TODO: Return the results to the API ...
      task->completed = true;
      break;
    }
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

void plc4c_driver_s7_free_write_response_item(
    plc4c_list_element* write_item_element) {
  
  plc4c_response_item* value_item;
  value_item = (plc4c_response_item*)write_item_element->value;
  /*
      // do not delete the plc4c_item
      // we also, in THIS case don't delete the random value which
      // isn't really a pointer. 

    TODO: what does above comment mean? Possibly written as we where
    casting to plc4c_response_value_item insted of plc4c_response_item
    in which case random explosion would probably occur on NULL'ing and 
    freeing. Cast correctly I think comment is no more valid...
  */
  //free(value_item->item);
  //value_item->item = NULL;
}

void plc4c_driver_s7_free_write_response(plc4c_write_response* response) {
  // the request will be cleaned up elsewhere
  plc4c_utils_list_delete_elements(response->response_items,
                                   &plc4c_driver_s7_free_write_response_item);
}

plc4c_return_code plc4c_driver_s7_parse_write_responce(
                                plc4c_write_request* request, 
                                plc4c_write_response* response,
                                plc4c_s7_read_write_tpkt_packet* packet) {
  
  // Locals
  plc4c_s7_read_write_s7_message* s7_packet;
  plc4c_list_element* request_list_element;
  plc4c_list_element* response_list_element;
  
  plc4c_request_value_item* request_item;
  plc4c_s7_read_write_s7_var_payload_status_item* s7_payload_status;
  plc4c_response_item* response_item;
  plc4c_return_code result;
  
	// Iterate over the request items and use the types to decode the
	// response items. TODO: Decode the return codes in the response ...
  s7_packet = packet->payload->payload;
	request_list_element = plc4c_utils_list_tail(request->items);
	response_list_element = plc4c_utils_list_tail(s7_packet->payload->s7_payload_write_var_response_items);

	while ((request_list_element != NULL) && (response_list_element != NULL)) {
		
    request_item = request_list_element->value;
		s7_payload_status = response_list_element->value;

		// Create a new response value-item
		response_item = malloc(sizeof(plc4c_response_item));
		if (response_item == NULL)
		  return NO_MEMORY;
		
		response_item->item = request_item->item;
    if (s7_payload_status->return_code == plc4c_s7_read_write_data_transport_error_code_OK)
      response_item->response_code = PLC4C_RESPONSE_CODE_OK;
    else
      // TODO: how to map plc4c_s7_read_write_data_transport_error_code to 
      // plc4c_responce_code, same issue in driver_s7_sm_read.c
      response_item->response_code = PLC4C_RESPONSE_CODE_INTERNAL_ERROR;


		// Add the value-item to the list.
		plc4c_utils_list_insert_head_value(response->response_items, response_item);

		request_list_element = request_list_element->next;
		response_list_element = response_list_element->next;
	}
  return OK;
}