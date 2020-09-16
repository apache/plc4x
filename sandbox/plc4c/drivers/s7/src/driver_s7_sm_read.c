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

enum plc4c_driver_s7_read_states {
  PLC4C_DRIVER_S7_READ_INIT,
  PLC4C_DRIVER_S7_READ_FINISHED
};

plc4c_return_code plc4c_driver_s7_read_machine_function(
    plc4c_system_task* task) {
  plc4c_read_request_execution* read_request_execution = task->context;
  if (read_request_execution == NULL) {
    return INTERNAL_ERROR;
  }
  plc4c_read_request* read_request = read_request_execution->read_request;
  if (read_request == NULL) {
    return INTERNAL_ERROR;
  }
  plc4c_connection* connection = task->connection;
  if (connection == NULL) {
    return INTERNAL_ERROR;
  }

  switch (task->state_id) {
    case PLC4C_DRIVER_S7_READ_INIT: {
      plc4c_s7_read_write_tpkt_packet* s7_read_request_packet;
      plc4c_return_code return_code =
          plc4c_driver_s7_create_s7_read_request(
              read_request, &s7_read_request_packet);
      if (return_code != OK) {
        return return_code;
      }

      // Send the packet to the remote.
      return_code = plc4c_driver_s7_send_packet(connection, s7_read_request_packet);
      if (return_code != OK) {
        return return_code;
      }

      task->state_id = PLC4C_DRIVER_S7_READ_FINISHED;
      break;
    }
    case PLC4C_DRIVER_S7_READ_FINISHED: {
      // Read a response packet.
      plc4c_s7_read_write_tpkt_packet* s7_read_response_packet;
      plc4c_return_code return_code =
          plc4c_driver_s7_receive_packet(connection, &s7_read_response_packet);
      // If we haven't read enough to process a full message, just try again
      // next time.
      if (return_code == UNFINISHED) {
        return OK;
      } else if (return_code != OK) {
        return return_code;
      }

      // Check the response.
      plc4c_s7_read_write_s7_parameter* parameter = s7_read_response_packet->payload->payload->parameter;
      if(parameter->_type != plc4c_s7_read_write_s7_parameter_type_plc4c_s7_read_write_s7_parameter_read_var_response) {
        return INTERNAL_ERROR;
      }
      // Check if the number of items matches that of the request
      // (Otherwise we won't know how to interpret the items)
      if(parameter->s7_parameter_read_var_response_num_items != plc4c_utils_list_size(read_request->items)) {
        return INTERNAL_ERROR;
      }
      // Iterate over the request items and use the types to decode the
      // response items.
      plc4c_s7_read_write_s7_payload* payload = s7_read_response_packet->payload->payload->payload;
      plc4c_list_element* cur_request_item_element = plc4c_utils_list_tail(read_request->items);
      plc4c_list_element* cur_response_item_element = plc4c_utils_list_tail(payload->s7_payload_read_var_response_items);
      while((cur_request_item_element != NULL) && (cur_response_item_element != NULL)) {
        plc4c_item* cur_request_item = cur_request_item_element->value;
        plc4c_s7_read_write_s7_address* s7_address = cur_request_item->address;
        plc4c_s7_read_write_transport_size transport_size = s7_address->s7_address_any_transport_size;
        uint16_t num_elements = s7_address->s7_address_any_number_of_elements;

        plc4c_s7_read_write_s7_var_payload_data_item* cur_response_item = cur_response_item_element->value;
        

        cur_request_item_element = cur_request_item_element->next;
        cur_request_item_element = cur_response_item_element->next;
      }

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
