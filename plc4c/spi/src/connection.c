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

plc4c_transport *plc4c_connection_get_transport(plc4c_connection *connection) {
  return connection->transport;
}

void plc4c_connection_set_transport(plc4c_connection *connection,
                                    plc4c_transport *transport) {
  connection->transport = transport;
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

plc4c_driver *plc4c_connection_get_driver(plc4c_connection *connection) {
  return connection->driver;
}

void plc4c_connection_set_driver(plc4c_connection *connection,
                                 plc4c_driver *driver) {
  connection->driver = driver;
}

void plc4c_connection_set_configuration(plc4c_connection *connection,
                                        void *configuration) {
  connection->configuration = configuration;
}

bool plc4c_connection_has_error(plc4c_connection *connection) { 
  return false; 
}

plc4c_return_code plc4c_connection_disconnect(plc4c_connection *connection) {
  
  plc4c_system_task *new_disconnection_task = NULL;
  plc4c_return_code result;
  
  result = connection->driver->disconnect_function(connection, &new_disconnection_task);
  if (result != OK) {
    return -1;
  }
  // Increment the number of running tasks for this connection.
  connection->num_running_system_tasks++;
  plc4c_system *system = plc4c_connection_get_system(connection);
  plc4c_utils_list_insert_tail_value(system->task_list, new_disconnection_task);

  return OK;
}


void plc4c_connection_destroy(plc4c_connection *connection) {
  
  if (connection == NULL) {
    return;
  }
  if (connection->connection_string != NULL) {
    free(connection->connection_string);
    connection->connection_string = NULL;
  }
  if (connection->configuration != NULL) {
    free(connection->configuration);
    connection->configuration = NULL;
  }
  if (connection->transport_configuration != NULL) {
    free(connection->transport_configuration);
    connection->transport_configuration = NULL;
  }
  if (connection->transport_code != NULL) {
    free(connection->transport_code);
    connection->transport_code = NULL;
  }
  if (connection->transport_connect_information != NULL) {
    free(connection->transport_connect_information);
    connection->transport_connect_information = NULL;
  }
  if (connection->protocol_code != NULL) {
    free(connection->protocol_code);
    connection->protocol_code = NULL;
  }
  if (connection->parameters != NULL) {
    free(connection->parameters);
    connection->parameters = NULL;
  }
  // TODO: verify free, seems too obvious to omit next line
  free(connection);
  connection = NULL;
}

char *plc4c_connection_get_connection_string(plc4c_connection *connection) {
  return connection->connection_string;
}

bool plc4c_connection_get_supports_reading(plc4c_connection *connection) {
  return connection->supports_reading;
}

plc4c_return_code plc4c_connection_create_read_request(
    plc4c_connection *connection, plc4c_read_request **read_request) {
  // NEED NULL ASSERTS

  plc4c_read_request *new_read_request = malloc(sizeof(plc4c_read_request));
  new_read_request->connection = connection;
  plc4c_utils_list_create(&(new_read_request->items));

  *read_request = new_read_request;
  return OK;
}

bool plc4c_connection_get_supports_writing(plc4c_connection *connection) {
  return connection->supports_writing;
}

plc4c_return_code plc4c_connection_create_write_request(
    plc4c_connection *connection, plc4c_write_request **write_request) {
  // NEED NULL ASSERTS

  plc4c_write_request *new_write_request =
      (plc4c_write_request *)malloc(sizeof(plc4c_write_request));
  new_write_request->connection = connection;
  plc4c_utils_list_create(&(new_write_request->items));

  *write_request = new_write_request;
  return OK;
}

bool plc4c_connection_get_supports_subscriptions(plc4c_connection *connection) {
  return connection->supports_subscriptions;
}

plc4c_return_code plc4c_connection_create_subscription_request(
    plc4c_connection *connection,
    plc4c_subscription_request **subscription_request) {
  // NEED NULL ASSERTS

  plc4c_subscription_request *new_subscription_request =
      malloc(sizeof(plc4c_subscription_request));
  new_subscription_request->connection = connection;
  plc4c_utils_list_create(&(new_subscription_request->items));

  *subscription_request = new_subscription_request;
  return OK;
}

plc4c_return_code plc4c_connection_create_unsubscription_request(
    plc4c_connection *connection,
    plc4c_unsubscription_request **unsubscription_request) {
  // NEED NULL ASSERTS

  plc4c_unsubscription_request *new_unsubscription_request =
      malloc(sizeof(plc4c_unsubscription_request));
  new_unsubscription_request->connection = connection;
  plc4c_utils_list_create(&(new_unsubscription_request->items));

  *unsubscription_request = new_unsubscription_request;
  return OK;
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
