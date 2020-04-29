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

#include <stdlib.h>
#include <plc4c/spi/types_private.h>
#include <plc4c/driver_simulated.h>

// State definitions
enum read_states {
    READ_INIT
};

enum write_states {
    WRITE_INIT
};

struct plc4c_driver_simulated_item_t {
    // Sort of senseless, but for keeping track of the fact that it's an plc4c_item.
    plc4c_item super;

    // Actual properties goes here ...
};
typedef struct plc4c_driver_simulated_item_t plc4c_driver_simulated_item;

plc4c_item *plc4c_driver_simulated_parse_address(const char *address_string) {
    plc4c_driver_simulated_item *item = (plc4c_driver_simulated_item *) malloc(sizeof(plc4c_driver_simulated_item));
    // TODO: Actually to the parsing of the address ...
    return (plc4c_item *) item;
}

return_code plc4c_driver_simulated_connect_function(plc4c_system *system, plc4c_connection *connection,
                                                    plc4c_system_task **task) {
    plc4c_system_task *new_task = malloc(sizeof(plc4c_system_task));
    new_task->context = connection;
    // There's nothing to do here, so no need for a state-machine.
    new_task->state_id = -1;
    new_task->state_machine_function = NULL;
    // We're setting this to true as there's nothing to do.
    new_task->completed = true;
    *task = new_task;
    return OK;
}

return_code plc4c_driver_simulated_read_function(plc4c_system *system, plc4c_connection *connection,
                                                 plc4c_read_request *read_request, plc4c_system_task **task) {
    plc4c_system_task *new_task = malloc(sizeof(plc4c_system_task));
    new_task->state_id = READ_INIT;
    new_task->state_machine_function = NULL;
    new_task->context = read_request;
    *task = new_task;
    return OK;
}

return_code plc4c_driver_simulated_write_function(plc4c_system *system, plc4c_connection *connection,
                                                  plc4c_write_request *write_request, plc4c_system_task **task) {
    plc4c_system_task *new_task = malloc(sizeof(plc4c_system_task));
    new_task->state_id = WRITE_INIT;
    new_task->state_machine_function = NULL;
    new_task->context = write_request;
    *task = new_task;
    return OK;
}

plc4c_driver *plc4c_driver_simulated_create() {
    plc4c_driver *driver = (plc4c_driver *) malloc(sizeof(plc4c_driver));
    driver->protocol_code = "simulated";
    driver->protocol_name = "Simulated PLC4X Datasource";
    driver->default_transport_code = "dummy";
    driver->parse_address_function = &plc4c_driver_simulated_parse_address;
    driver->connect_function = &plc4c_driver_simulated_connect_function;
    driver->read_function = &plc4c_driver_simulated_read_function;
    driver->write_function = &plc4c_driver_simulated_write_function;
    return driver;
}

