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
#include "plc4c/driver_plc4x.h"
#include "plc4c/driver_plc4x_sm.h"
#include "plc4c/driver_plc4x_packets.h"
#include "plc4x_message.h"


enum plc4c_driver_plc4x_connect_states {
  PLC4C_DRIVER_PLC4X_CONNECT_INIT,
  PLC4C_DRIVER_PLC4X_CONNECT_TRANSPORT_CONNECT,
  PLC4C_DRIVER_PLC4X_CONNECT_SEND_CONNECT_REQUEST,
  PLC4C_DRIVER_PLC4X_CONNECT_RECEIVE_CONNECT_RESPONSE,
  PLC4C_DRIVER_PLC4X_CONNECT_FINISHED
};

/**
 * State machine function for establishing a connection to a remote S7 device.
 * @param task the current system task
 * @return return code of the current state machine step execution
 */
plc4c_return_code plc4c_driver_plc4x_connect_machine_function(plc4c_system_task* task) {

  plc4c_connection* connection;
  plc4c_driver_plc4x_config* configuration;
  plc4c_return_code return_code;
  plc4c_plc4x_read_write_plc4x_message* packet;

  connection = task->connection;
  if (connection == NULL) {
    return INTERNAL_ERROR;
  }
  // If we were already connected, return an error
  if (plc4c_connection_get_connected(connection)) {
    return ALREADY_CONNECTED;
  }
  configuration = connection->configuration;

  // Initialize the pdu id (The first messages are hard-coded)
//  configuration->pdu_id = 1;

  switch (task->state_id) {
    // Initialize some internal data-structures.
    case PLC4C_DRIVER_PLC4X_CONNECT_INIT: {
      // Calculate some internal settings from the values provided in the
      // configuration.

      task->state_id = PLC4C_DRIVER_PLC4X_CONNECT_TRANSPORT_CONNECT;
      break;
    }
    case PLC4C_DRIVER_PLC4X_CONNECT_TRANSPORT_CONNECT: {
      return_code = connection->transport->open(connection->transport_configuration);
      if(return_code != OK) {
        return return_code;
      }

      task->state_id = PLC4C_DRIVER_PLC4X_CONNECT_SEND_CONNECT_REQUEST;
      break;
    }
    // Send a connection request.
    case PLC4C_DRIVER_PLC4X_CONNECT_SEND_CONNECT_REQUEST: {
      // Get a connection response for the settings in the config.
      return_code = plc4c_driver_plc4x_create_connection_request(configuration, &packet);
      if (return_code != OK) {
        return return_code;
      }

      // Send the packet to the remote.
      return_code = plc4c_driver_plc4x_send_packet(connection, packet);
      if (return_code != OK) {
        return return_code;
      }
      plc4c_driver_plc4x_destroy_connection_request(packet);
      task->state_id = PLC4C_DRIVER_PLC4X_CONNECT_RECEIVE_CONNECT_RESPONSE;
      break;
    }
    // Receive a connection response.
    case PLC4C_DRIVER_PLC4X_CONNECT_RECEIVE_CONNECT_RESPONSE: {
      // Read a response packet.
      return_code = plc4c_driver_plc4x_receive_packet(connection, &packet);
      // If we haven't read enough to process a full message, just try again
      // next time.
      if (return_code == UNFINISHED) {
        return OK;
      } else if (return_code != OK) {
        return return_code;
      }

      plc4c_driver_plc4x_destroy_packet(packet);
      // If we got the expected response, continue with the next higher level
      // of connection.
      task->state_id = PLC4C_DRIVER_PLC4X_CONNECT_FINISHED;
      break;
    }
    // Clean up some internal data-structures.
    case PLC4C_DRIVER_PLC4X_CONNECT_FINISHED: {
      plc4c_connection_set_connected(connection, true);
      task->completed = true;
      break;
    }
    // If an unexpected state id was received, this is not really something
    // we can recover from.
    default: {
      task->completed = true;
      return INTERNAL_ERROR;
    }
  }
  return OK;
}

plc4c_return_code plc4c_driver_plc4x_connect_function(plc4c_connection* connection,
                                                   plc4c_system_task** task) {
                                                     
  plc4c_system_task* new_task = malloc(sizeof(plc4c_system_task));
  // There's nothing to do here, so no need for a state-machine.
  new_task->state_id = PLC4C_DRIVER_PLC4X_CONNECT_INIT;
  new_task->state_machine_function = &plc4c_driver_plc4x_connect_machine_function;
  new_task->completed = false;
  new_task->context = connection;
  new_task->connection = connection;
  *task = new_task;
  return OK;
}
