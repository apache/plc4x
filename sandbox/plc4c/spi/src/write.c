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
#include <plc4c/write.h>
#include <stdlib.h>

plc4c_connection *plc4c_write_request_get_connection(
    plc4c_write_request *write_request) {
  return write_request->connection;
}

void plc4c_write_request_set_connection(plc4c_write_request *write_request,
                                        plc4c_connection *connection) {
  write_request->connection = connection;
}

plc4c_return_code plc4c_write_request_add_item(
    plc4c_write_request *write_request, char *address, plc4c_data *value) {
  // Parse an address string and get a driver-dependent data-structure
  // representing the address back.
  plc4c_request_value_item *value_item;
  plc4c_item *address_item;
  plc4c_return_code result;

  address_item = malloc(sizeof(plc4c_item));
  if(address_item == NULL) 
    return NO_MEMORY;
  
  result = write_request->connection->driver->parse_address_function(
          address, &address_item->address);
  
  if(result != OK) 
    return result;
  
  // Create a new value item, binding an address item to a value.
  value_item = malloc(sizeof(plc4c_request_value_item));
  value_item->item = address_item;
  value_item->value = value;

  plc4c_utils_list_insert_head_value(write_request->items, value_item);
  return OK;
}

plc4c_return_code plc4c_write_request_execute(
    plc4c_write_request *write_request,
    plc4c_write_request_execution **write_request_execution) {

  plc4c_write_request_execution *new_write_request_execution;
  plc4c_connection *connection;
  plc4c_driver *driver;
  plc4c_system *system;
  plc4c_list *system_tasks;

  // Extract some plc4c pointers for clarity
  connection = plc4c_write_request_get_connection(write_request);
  system = plc4c_connection_get_system(connection);
  driver = plc4c_connection_get_driver(connection);
  system_tasks = plc4c_system_get_task_list(system);

  // Inject the default write context into the system task.
  new_write_request_execution = malloc(sizeof(plc4c_write_request_execution));
  new_write_request_execution->write_request = write_request;
  new_write_request_execution->write_response = NULL;
  new_write_request_execution->system_task = NULL;
  
  driver->write_function(new_write_request_execution, &(new_write_request_execution->system_task));

  // Increment the number of running tasks for this connection.
  plc4c_connection_task_added(connection);

  // Add the new task to the task-list.
  plc4c_utils_list_insert_tail_value(system_tasks, new_write_request_execution->system_task);

  *write_request_execution = new_write_request_execution;
  return OK;
}

bool plc4c_write_request_check_finished_successfully(
    plc4c_write_request_execution *write_request_execution) {

  if (write_request_execution == NULL) {
    return true;
  }
  if (write_request_execution->system_task == NULL) {
    return true;
  }
  return write_request_execution->system_task->completed;
}

bool plc4c_write_request_execution_check_completed_with_error(
    plc4c_write_request_execution *write_request_execution) {
  return false;
}

plc4c_write_response *plc4c_write_request_execution_get_response(
    plc4c_write_request_execution *write_request_execution) {
  if (write_request_execution == NULL) {
    return NULL;
  }
  return write_request_execution->write_response;
}




void plc4c_write_request_destroy(plc4c_write_request *write_request) {
  write_request->connection->driver->free_write_request_function(write_request);
  free(write_request);
}

void plc4c_write_request_execution_destroy(
    plc4c_write_request_execution *write_request_execution) {
  free(write_request_execution);
}

void plc4c_write_response_destroy(
    plc4c_write_response *write_response) {
  write_response->write_request->connection->driver
      ->free_write_response_function(write_response);
  free(write_response);
}

