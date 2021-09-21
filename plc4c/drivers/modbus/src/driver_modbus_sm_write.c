/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

#include <plc4c/plc4c.h>
#include <plc4c/spi/types_private.h>
#include <stdlib.h>
#include <string.h>
#include "plc4c/driver_modbus_packets.h"
#include "modbus_tcp_adu.h"

enum plc4c_driver_modbus_write_states {
  PLC4C_DRIVER_MODBUS_WRITE_INIT,
  PLC4C_DRIVER_MODBUS_WRITE_FINISHED
};

plc4c_return_code plc4c_driver_modbus_write_machine_function(
    plc4c_system_task* task) {
  plc4c_write_request_execution* write_request_execution = task->context;
  if (write_request_execution == NULL) {
    return INTERNAL_ERROR;
  }
  plc4c_write_request* write_request = write_request_execution->write_request;
  if (write_request == NULL) {
    return INTERNAL_ERROR;
  }
  plc4c_connection* connection = task->context;
  if (connection == NULL) {
    return INTERNAL_ERROR;
  }

  switch (task->state_id) {
    case PLC4C_DRIVER_MODBUS_WRITE_INIT: {
      plc4c_modbus_read_write_modbus_tcp_adu* modbus_write_request_packet;
      plc4c_return_code return_code =
          plc4c_driver_modbus_create_modbus_write_request(
              write_request, &modbus_write_request_packet);
      if (return_code != OK) {
        return return_code;
      }

      // Send the packet to the remote.
      return_code = plc4c_driver_modbus_send_packet(
          connection, modbus_write_request_packet);
      if (return_code != OK) {
        return return_code;
      }

      task->state_id = PLC4C_DRIVER_MODBUS_WRITE_FINISHED;
      break;
    }
    case PLC4C_DRIVER_MODBUS_WRITE_FINISHED: {
      // Read a response packet.
      plc4c_modbus_read_write_modbus_tcp_adu* modbus_write_response_packet;
      plc4c_return_code return_code =
          plc4c_driver_modbus_receive_packet(
              connection, &modbus_write_response_packet);
      // If we haven't read enough to process a full message, just try again
      // next time.
      if (return_code == UNFINISHED) {
        return OK;
      } else if (return_code != OK) {
        return return_code;
      }

      // TODO: Check the response ...
      // TODO: Decode the return codes in the response ...
      // TODO: Return the results to the API ...

      task->completed = true;
      break;
    }
  }
  return OK;
}

plc4c_return_code plc4c_driver_modbus_write_function(
    plc4c_write_request_execution* write_request_execution,
    plc4c_system_task** task) {
  plc4c_system_task* new_task = malloc(sizeof(plc4c_system_task));
  new_task->state_id = PLC4C_DRIVER_MODBUS_WRITE_INIT;
  new_task->state_machine_function = &plc4c_driver_modbus_write_machine_function;
  new_task->completed = false;
  new_task->context = write_request_execution;
  new_task->connection = write_request_execution->system_task->connection;
  *task = new_task;
  return OK;
}

void plc4c_driver_modbus_free_write_response_item(
    plc4c_list_element* write_item_element) {
  plc4c_response_value_item* value_item =
      (plc4c_response_value_item*)write_item_element->value;
  // do not delete the plc4c_item
  // we also, in THIS case don't delete the random value which isn't really
  // a pointer
  // free(value_item->value);
  value_item->value = NULL;
}

void plc4c_driver_modbus_free_write_response(plc4c_write_response* response) {
  // the request will be cleaned up elsewhere
  plc4c_utils_list_delete_elements(response->response_items,
                                   &plc4c_driver_modbus_free_write_response_item);
}

