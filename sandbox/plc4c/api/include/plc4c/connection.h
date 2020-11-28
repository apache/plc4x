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
#ifndef PLC4C_CONNECTION_H_
#define PLC4C_CONNECTION_H_
#ifdef __cplusplus
extern "C" {
#endif

#include <stdbool.h>

#include "types.h"
#include "utils/list.h"

/**
 * Initializes the connection members to their defaults.
 * @param connection plc4c_connection
 */
void plc4c_connection_initialize(plc4c_connection *connection);

/**
 * Destroy a previously created connection.
 * @param connection
 */
void plc4c_connection_destroy(plc4c_connection *connection);

/**
 * Returns the connection string for a given connection
 * @param connection plc4c_connection
 * @return connection string
 */
char *plc4c_connection_get_connection_string(plc4c_connection *connection);

/**
 * Set the connection string
 * @param connection plc4c_connection
 * @param connection_string
 */
void plc4c_connection_set_connection_string(plc4c_connection *connection,
                                            char *connection_string);

/**
 * Returns the protocol code for the given connection
 * @param connection plc4c_connection
 * @return protocol code
 */
char *plc4c_connection_get_protocol_code(plc4c_connection *connection);

/**
 * Set the protocol code
 * @param connection plc4c_connection
 * @param protocol_code
 */
void plc4c_connection_set_protocol_code(plc4c_connection *connection,
                                        char *protocol_code);

/**
 * Returns the transport code for the given connection
 * @param connection
 * @return transport code
 */
char *plc4c_connection_get_transport_code(plc4c_connection *connection);

/**
 * Set the transport code
 * @param connection plc4c_connection
 * @param transport_code
 */
void plc4c_connection_set_transport_code(plc4c_connection *connection,
                                         char *transport_code);

/**
 * Returns the transport for the given connection
 * @param connection plc4c_connection
 * @return transport plc4c_transport
 */
plc4c_transport *plc4c_connection_get_transport(plc4c_connection *connection);

/**
 * Set the transport
 * @param connection plc4c_connection
 * @param transport plc4c_transport
 */
void plc4c_connection_set_transport(plc4c_connection *connection,
                                    plc4c_transport *transport);

/**
 * Returns the transport connection information for a given connection
 * @param connection plc4c_connection
 * @return transport connection information
 */
char *plc4c_connection_get_transport_connect_information(
    plc4c_connection *connection);

/**
 * Set the transport connect information
 * @param connection plc4c_connection
 * @param transport_connect_information
 */
void plc4c_connection_set_transport_connect_information(
    plc4c_connection *connection, char *transport_connect_information);

/**
 * Returns the parameters for a given connection
 * @param connection plc4c_connection
 * @return parameters
 */
char *plc4c_connection_get_parameters(plc4c_connection *connection);

/**
 * Set the parameters
 * @param connection plc4c_connection
 * @param parameters
 */
void plc4c_connection_set_parameters(plc4c_connection *connection,
                                     char *parameters);

/**
 * Check if the connection is established successfully.
 *
 * @param connection the connection.
 * @return true if the connection is established.
 */
bool plc4c_connection_get_connected(plc4c_connection *connection);

/**
 * Sets the connected status of a given connection
 * @param connection plc4c_connection
 * @param connected bool is connected
 */
void plc4c_connection_set_connected(plc4c_connection *connection,
                                    bool connected);

/**
 * Returns the status of the disconnect flag for a given connection
 * @param connection plc4c_connection
 * @return if the flag is set
 */
bool plc4c_connection_get_disconnect(plc4c_connection *connection);

/**
 * Sets the status of the disconnect flag for a given connection
 * @param connection plc4c_connection
 * @param disconnect the new flag status
 */
void plc4c_connection_set_disconnect(plc4c_connection *connection,
                                     bool disconnect);

/**
 * Returns the system associated with a given connection
 * @param connection plc4c_connection
 * @return plc4c_system
 */
plc4c_system *plc4c_connection_get_system(plc4c_connection *connection);

/**
 * Sets a plc4c_system for a given connection
 * @param connection plc4c_connection
 * @param system plc4c_system
 */
void plc4c_connection_set_system(plc4c_connection *connection,
                                 plc4c_system *system);

/**
 * Returns the driver associated with a given connection
 * @param connection plc4c_connection
 * @return plc4c_driver
 */
plc4c_driver *plc4c_connection_get_driver(plc4c_connection *connection);

/**
 * Sets a plc4c_driver for a given connection
 * @param connection plc4c_connection
 * @param plc4c_driver driver
 */
void plc4c_connection_set_driver(plc4c_connection *connection,
                                 plc4c_driver *driver);

/**
 * Sets the configuration for a given connection.
 *
 * @param connection plc4c connection
 * @param configuration configuration
 */
void plc4c_connection_set_configuration(plc4c_connection *connection,
                                        void *configuration);
/**
 * Check if an error occurred while connecting.
 *
 * @param connection the connection.
 * @return true if an error occurred while connecting.
 */
bool plc4c_connection_has_error(plc4c_connection *connection);

/**
 * Function to terminate a connection to a PLC.
 *
 * @param connection
 * @param plc4c_return_code
 */
plc4c_return_code plc4c_connection_disconnect(plc4c_connection *connection);

/**
 * Check if the current connection supports read operations.
 *
 * @param connection reference to the connection
 * @return true if the connection supports reading, false otherwise
 */
bool plc4c_connection_get_supports_reading(plc4c_connection *connection);

/**
 * Initializes an empty read-request.
 *
 * @param connection connection that this read-request will be executed on.
 * @param read_request pointer to the read-request
 * @param plc4c_return_code
 */
plc4c_return_code plc4c_connection_create_read_request(
    plc4c_connection *connection, plc4c_read_request **read_request);

/**
 * Check if the current connection supports write operations.
 *
 * @param connection reference to the connection
 * @return true if the connection supports writing, false otherwise
 */
bool plc4c_connection_get_supports_writing(plc4c_connection *connection);

/**
 * Initializes an empty write-request.
 *
 * @param connection connection that this write-request will be executed on.
 * @param write_request pointer to the write-request
 * @return plc4c_return_code
 */
plc4c_return_code plc4c_connection_create_write_request(
    plc4c_connection *connection, plc4c_write_request **write_request);

/**
 * Check if the current connection supports subscriptions.
 *
 * @param connection reference to the connection
 * @return true if the connection supports subscriptions, false otherwise
 */
bool plc4c_connection_get_supports_subscriptions(plc4c_connection *connection);

/**
 * Initializes an empty subscription-request.
 *
 * @param connection connection that this subscription-request will be executed
 * on.
 * TODO: Somehow we need to provide the subscription type and depending on the
 * type some additional parameters (Cycle-time).
 * @param subscription_request pointer to the subscription-request
 * @return plc4c_return_code
 */
plc4c_return_code plc4c_connection_create_subscription_request(
    plc4c_connection *connection,
    plc4c_subscription_request **subscription_request);

/**
 * Initializes an empty unsubscription-request.
 *
 * @param connection connection that this unsubscription-request will be
 * executed on.
 * TODO: Somehow we need to provide the unsubscription type and depending on the
 * type some additional parameters (Cycle-time).
 * @param addresses list of address strings.
 * @param unsubscription_request pointer to the unsubscription-request
 * @return plc4c_return_code
 */
plc4c_return_code plc4c_connection_create_unsubscription_request(
    plc4c_connection *connection,
    plc4c_unsubscription_request **unsubscription_request);

/**
 * Returns the current number of running tasks for this connection
 * @param connection plc4c_connection
 * @return the count of running tasks
 */
int plc4c_connection_get_running_tasks_count(plc4c_connection *connection);

/**
 *
 * @param connection
 */
void plc4c_connection_running_tasks_clear(plc4c_connection *connection);

/**
 * Signals the plc4c_connection that a task has been added
 * @param connection plc4c_connection
 * @return the number of tasks after this operation
 */
int plc4c_connection_task_added(plc4c_connection *connection);

/**
 * Signals the plc4c_connection that a task has been completed and removed
 * @param connection plc4c_connection
 * @return the number of tasks after this operation
 */
int plc4c_connection_task_removed(plc4c_connection *connection);

#ifdef __cplusplus
}
#endif
#endif  // PLC4C_CONNECTION_H_