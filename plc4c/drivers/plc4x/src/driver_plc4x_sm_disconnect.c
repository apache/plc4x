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

#include <plc4c/plc4c.h>
#include <plc4c/spi/types_private.h>
#include <stdlib.h>
#include "plc4c/driver_plc4x.h"

enum plc4c_driver_plc4x_disconnect_states {
  PLC4C_DRIVER_PLC4X_DISCONNECT_INIT,
  PLC4C_DRIVER_PLC4X_DISCONNECT_WAIT_TASKS_FINISHED,
  PLC4C_DRIVER_PLC4X_DISCONNECT_FINISHED
};

plc4c_return_code plc4c_driver_plc4x_disconnect_machine_function(
    plc4c_system_task* task) {
  plc4c_connection* connection = task->context;
  if (connection == NULL) {
    return INTERNAL_ERROR;
  }

  switch (task->state_id) {
    case PLC4C_DRIVER_PLC4X_DISCONNECT_INIT: {
      plc4c_connection_set_disconnect(connection, true);
      task->state_id = PLC4C_DRIVER_PLC4X_DISCONNECT_WAIT_TASKS_FINISHED;
      break;
    }
    case PLC4C_DRIVER_PLC4X_DISCONNECT_WAIT_TASKS_FINISHED: {
      // The disconnect system-task also counts.
      if (plc4c_connection_get_running_tasks_count(connection) == 1) {
        plc4c_connection_set_connected(connection, false);
        task->completed = true;
        task->state_id = PLC4C_DRIVER_PLC4X_DISCONNECT_FINISHED;
      }
      break;
    }
    case PLC4C_DRIVER_PLC4X_DISCONNECT_FINISHED: {
      // Do nothing
      break;
    }
    default: {
      return INTERNAL_ERROR;
    }
  }
  return OK;
}

plc4c_return_code plc4c_driver_plc4x_disconnect_function(
    plc4c_connection* connection, plc4c_system_task** task) {
  plc4c_system_task* new_task = malloc(sizeof(plc4c_system_task));
  new_task->state_id = PLC4C_DRIVER_PLC4X_DISCONNECT_INIT;
  new_task->state_machine_function =
      &plc4c_driver_plc4x_disconnect_machine_function;
  new_task->completed = false;
  new_task->context = connection;
  new_task->connection = connection;
  *task = new_task;
  return OK;
}
