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

#include <plc4c/connection.h>
#include <plc4c/spi/types_private.h>
#include <stdlib.h>
#include <string.h>

void plc4c_connection_initialize(plc4c_connection *new_connection) {
  new_connection->connected = false;
  new_connection->disconnect = false;
  new_connection->num_running_system_tasks = 0;

  new_connection->system = NULL;
  new_connection->driver = NULL;
  new_connection->transport = NULL;

  new_connection->supports_reading = NULL;
  new_connection->supports_writing = NULL;
  new_connection->supports_subscriptions = NULL;
}

void plc4c_connection_set_connection_string(plc4c_connection *connection,
                                            char *connection_string) {
  if (connection_string != NULL) {
    connection->connection_string =
        (char *)malloc((strlen(connection_string) + 1) * sizeof(char));
    strcpy(connection->connection_string, connection_string);
  } else {
    connection->connection_string = NULL;
  }
}

char *plc4c_connection_get_protocol_code(plc4c_connection *connection) {
  return connection->protocol_code;
}

void plc4c_connection_set_protocol_code(plc4c_connection *connection,
                                        char *protocol_code) {
  if (protocol_code != NULL) {
    connection->protocol_code =
        (char *)malloc((strlen(protocol_code) + 1) * sizeof(char));
    strcpy(connection->protocol_code, protocol_code);
  } else {
    connection->protocol_code = NULL;
  }
}

char *plc4c_connection_get_transport_code(plc4c_connection *connection) {
  return connection->transport_code;
}

void plc4c_connection_set_transport_code(plc4c_connection *connection,
                                         char *transport_code) {
  if (transport_code != NULL) {
    connection->transport_code =
        (char *)malloc((strlen(transport_code) + 1) * sizeof(char));
    strcpy(connection->transport_code, transport_code);
  } else {
    connection->transport_code = NULL;
  }
}

char *plc4c_connection_get_transport_connect_information(
    plc4c_connection *connection) {
  return connection->transport_connect_information;
}

void plc4c_connection_set_transport_connect_information(
    plc4c_connection *connection, char *transport_connect_information) {
  if (transport_connect_information != NULL) {
    connection->transport_connect_information = (char *)malloc(
        (strlen(transport_connect_information) + 1) * sizeof(char));
    strcpy(connection->transport_connect_information,
           transport_connect_information);
  } else {
    connection->transport_connect_information = NULL;
  }
}

char *plc4c_connection_get_parameters(plc4c_connection *connection) {
  return connection->parameters;
}

void plc4c_connection_set_parameters(plc4c_connection *connection,
                                     char *parameters) {
  if (parameters != NULL) {
    connection->parameters =
        (char *)malloc((strlen(parameters) + 1) * sizeof(char));
    strcpy(connection->parameters, parameters);
  } else {
    connection->parameters = NULL;
  }
}

bool plc4c_connection_get_connected(plc4c_connection *connection) {
  return connection->connected;
}

void plc4c_connection_set_connected(plc4c_connection *connection,
                                    bool connected) {
  connection->connected = connected;
}

bool plc4c_connection_get_disconnect(plc4c_connection *connection) {
  return connection->disconnect;
}

void plc4c_connection_set_disconnect(plc4c_connection *connection,
                                     bool disconnect) {
  connection->disconnect = disconnect;
}

plc4c_system *plc4c_connection_get_system(plc4c_connection *connection) {
  return connection->system;
}

void plc4c_connection_set_system(plc4c_connection *connection,
                                 plc4c_system *system) {
  connection->system = system;
}

bool plc4c_connection_has_error(plc4c_connection *connection) { return false; }

plc4c_return_code plc4c_connection_disconnect(plc4c_connection *connection) {
  plc4c_system_task *new_disconnection_task = NULL;
  plc4c_return_code result = connection->driver->disconnect_function(
      connection, &new_disconnection_task);
  if (result != OK) {
    return -1;
  }
  // Increment the number of running tasks for this connection.
  connection->num_running_system_tasks++;
  plc4c_utils_list_insert_tail_value(
      plc4c_system_get_task_list(plc4c_connection_get_system(connection)),
      new_disconnection_task);
  return OK;
}

void plc4c_connection_destroy(plc4c_connection *connection) {
  if (connection == NULL) {
    return;
  }
  if (connection->connection_string != NULL) {
    free(connection->connection_string);
  }
  if (connection->transport != NULL) {
    free(connection->transport);
  }
  if (connection->transport_code != NULL) {
    free(connection->transport_code);
  }
  if (connection->transport_connect_information != NULL) {
    free(connection->transport_connect_information);
  }
  if (connection->protocol_code != NULL) {
    free(connection->protocol_code);
  }
  if (connection->parameters != NULL) {
    free(connection->parameters);
  }
}

char *plc4c_connection_get_connection_string(plc4c_connection *connection) {
  return connection->connection_string;
}

bool plc4c_connection_get_supports_reading(plc4c_connection *connection) {
  return connection->supports_reading;
}

plc4c_return_code plc4c_connection_create_read_request(
    plc4c_connection *connection, plc4c_list *addresses,
    plc4c_read_request **read_request) {
  // NEED NULL ASSERTS

  // we need something to do
  if (plc4c_utils_list_size(addresses) == 0) {
    return INVALID_LIST_SIZE;
  }
  plc4c_read_request *new_read_request = malloc(sizeof(plc4c_read_request));
  new_read_request->connection = connection;
  plc4c_utils_list_create(&(new_read_request->items));
  plc4c_list_element *element = plc4c_utils_list_head(addresses);
  if (element != NULL) {
    do {
      plc4c_item *item =
          connection->driver->parse_address_function((char *)element->value);
      plc4c_utils_list_insert_tail_value(new_read_request->items, item);
      element = element->next;
    } while (element != NULL);
  }
  *read_request = new_read_request;
  return OK;
}

void plc4c_connection_destroy_read_response(
    plc4c_read_response *read_response) {
  read_response->read_request->connection->driver->free_read_response_function(
      read_response);
}

bool plc4c_connection_get_supports_writing(plc4c_connection *connection) {
  return connection->supports_writing;
}

plc4c_return_code plc4c_connection_create_write_request(
    plc4c_connection *connection, plc4c_list *addresses, plc4c_list *values,
    plc4c_write_request **write_request) {
  // NEED NULL ASSERTS

  // the address and value lists must match
  if (plc4c_utils_list_size(addresses) != plc4c_utils_list_size(values)) {
    return NON_MATCHING_LISTS;
  }

  // we need something to do
  if (plc4c_utils_list_size(addresses) == 0) {
    return INVALID_LIST_SIZE;
  }
  plc4c_write_request *new_write_request =
      (plc4c_write_request *)malloc(sizeof(plc4c_write_request));
  new_write_request->connection = connection;
  plc4c_utils_list_create(&(new_write_request->items));

  plc4c_list_element *address_element = plc4c_utils_list_head(addresses);
  plc4c_list_element *value_element = plc4c_utils_list_head(values);
  if (address_element != NULL && value_element != NULL) {
    do {
      char *address = (char *)address_element->value;
      // Parse an address string and get a driver-dependent data-structure
      // representing the address back.
      plc4c_item *address_item =
          connection->driver->parse_address_function(address);

      // Create a new value item, binding an address item to a value.
      plc4c_request_value_item *value_item =
          malloc(sizeof(plc4c_request_value_item));
      value_item->item = address_item;
      value_item->value = (plc4c_data *)value_element->value;

      // Add the new item ot the list of items.
      plc4c_utils_list_insert_tail_value(new_write_request->items, value_item);

      address_element = address_element->next;
      value_element = value_element->next;
    } while (address_element != NULL && value_element != NULL);
  }

  *write_request = new_write_request;
  return OK;
}

void plc4c_connection_destroy_write_response(
    plc4c_write_response *write_response) {
  write_response->write_request->connection->driver
      ->free_write_response_function(write_response);
}

int plc4c_connection_get_running_tasks_count(plc4c_connection *connection) {
  return connection->num_running_system_tasks;
}

int plc4c_connection_task_added(plc4c_connection *connection) {
  return ++connection->num_running_system_tasks;
}

int plc4c_connection_task_removed(plc4c_connection *connection) {
  return --connection->num_running_system_tasks;
}
