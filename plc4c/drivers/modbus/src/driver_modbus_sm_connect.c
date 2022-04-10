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

#include <plc4c/driver_modbus.h>
#include <plc4c/plc4c.h>
#include <plc4c/spi/types_private.h>

enum plc4c_driver_modbus_connect_states {
  PLC4C_DRIVER_MODBUS_CONNECT_INIT,
  PLC4C_DRIVER_MODBUS_CONNECT_TRANSPORT_CONNECT,
  PLC4C_DRIVER_MODBUS_CONNECT_FINISHED
};

plc4c_return_code plc4c_driver_modbus_connect_machine_function(
    plc4c_system_task *task) {
  plc4c_connection *connection = task->context;
  if (connection == NULL) {
    return INTERNAL_ERROR;
  }
  if (plc4c_connection_get_connected(connection)) {
    return ALREADY_CONNECTED;
  }

  plc4c_driver_modbus_config* configuration = connection->configuration;

  switch (task->state_id) {
    case PLC4C_DRIVER_MODBUS_CONNECT_INIT: {
      configuration->communication_id_counter = 1;

      task->state_id = PLC4C_DRIVER_MODBUS_CONNECT_TRANSPORT_CONNECT;
      break;
    }
    case PLC4C_DRIVER_MODBUS_CONNECT_TRANSPORT_CONNECT: {
      plc4c_return_code return_code = connection->transport->open(connection->transport_configuration);
      if(return_code != OK) {
        return return_code;
      }

      task->state_id = PLC4C_DRIVER_MODBUS_CONNECT_FINISHED;
      break;
    }
    case PLC4C_DRIVER_MODBUS_CONNECT_FINISHED: {
      plc4c_connection_set_connected(connection, true);
      task->completed = true;
      break;
    }
  }
  return OK;
}

plc4c_return_code plc4c_driver_modbus_connect_function(
    plc4c_connection *connection, plc4c_system_task **task) {

  plc4c_system_task *new_task = malloc(sizeof(plc4c_system_task));
  // There's nothing to do here, so no need for a state-machine.
  new_task->state_id = PLC4C_DRIVER_MODBUS_CONNECT_INIT;
  new_task->state_machine_function =
      &plc4c_driver_modbus_connect_machine_function;
  new_task->completed = false;
  new_task->context = connection;
  new_task->connection = connection;
  *task = new_task;
  return OK;
}
