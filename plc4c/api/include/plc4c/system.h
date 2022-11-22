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
#ifndef PLC4C_SYSTEM_H_
#define PLC4C_SYSTEM_H_
#ifdef __cplusplus
extern "C" {
#endif

#include <plc4c/utils/list.h>
#include <stdbool.h>

#include "types.h"

/**
 * SYSTEM CALLBACKS
 */

/**
 * Function pointer for a callback called when a driver is loaded.
 * Set in plc4c_system @see plc4c_system_set_on_driver_load_success_callback()
 *
 * @param driver
 */
typedef void (*plc4c_system_on_driver_load_success_callback)(
    plc4c_driver *driver);

/**
 * Function pointer for a callback called when loading a driver fails.
 * Set in plc4c_system @see plc4c_system_set_on_driver_load_failure_callback
 * NOTE: driver_name could be a pointer to the configuration for the driver
 * instead....
 *
 * @param driver_name
 * @param error_code
 */
typedef void (*plc4c_system_on_driver_load_failure_callback)(
    char *driver_name, plc4c_return_code error);

/**
 * Function pointer for a callback called when is successfully made
 * Set in plc4c_system @see plc4c_system_set_on_connect_success_callback().
 *
 * @param connection
 */
typedef void (*plc4c_system_on_connect_success_callback)(
    plc4c_connection *connection);

/**
 * Function pointer for a callback called when connecting fails.
 * Set in plc4c_system @see plc4c_system_set_on_connect_failure_callback.
 *
 * @param connection_string
 * @param error_code
 */
typedef void (*plc4c_system_on_connect_failure_callback)(
    char *connection_string, plc4c_return_code error);

/**
 * Function pointer for a callback called when is successfully made
 * Set in plc4c_system @see plc4c_system_set_on_disconnect_success_callback().
 *
 * @param connection
 */
typedef void (*plc4c_system_on_disconnect_success_callback)(
    plc4c_connection *connection);

/**
 * Function pointer for a callback called when connecting fails.
 * Set in plc4c_system @see plc4c_system_set_on_disconnect_failure_callback.
 *
 * @param connection
 * @param error_code
 */
typedef void (*plc4c_system_on_disconnect_failure_callback)(
    plc4c_connection *connection, plc4c_return_code error);

/**
 * Function pointer for a callback called when a driver returns an error.
 * Set in plc4c_system @see plc4c_system_set_on_loop_failure_callback.
 *
 * @param driver
 * @param connection
 * @param error_code
 */
typedef void (*plc4c_system_on_loop_failure_callback)(
    plc4c_driver *driver, plc4c_connection *connection,
    plc4c_return_code error);

/**
 * OTHER FUNCTION DEFS FOR SYSTEM
 */

/**
 * SYSTEM FUNCTIONS
 */

/**
 * Function to create a plc4c_system.
 *
 * @param system
 * @return NO_MEMORY if failed to create system
 */
plc4c_return_code plc4c_system_create(plc4c_system **system);

/**
 * Function to initialize the PLC4C system (Initialize the driver manager and
 * the list of enabled drivers)
 *
 * @param system
 * @return plc4c_return_code
 */
plc4c_return_code plc4c_system_init(plc4c_system *system);

/**
 * Function to destroy a plc4c_system.
 * This will also destroy all connections associated with the system.
 *
 * @param system
 */
void plc4c_system_destroy(plc4c_system *system);

/**
 * Returns the plc4c_list for tasks for a given system
 * TODO: Possilby this should be an SPI function
 * @param system plc4c_system
 * @return plc4c_list
 */
plc4c_list *plc4c_system_get_task_list(plc4c_system *system);

/**
 * Sets the plc4c_list for tasks for a given system
 * TODO: Possilby this should be an SPI function
 * @param system plc4c_system
 * @param task_list plc4c_list
 */
void plc4c_system_set_task_list(plc4c_system *system, plc4c_list *task_list);

void plc4c_system_remove_connection(plc4c_system *system,
                                    plc4c_connection *connection);

/**
 * Function to set the on_driver_loaded callback for the plc4c system.
 *
 * @param system
 * @param callback plc4c_system_callback_on_driver
 */
void plc4c_system_set_on_driver_load_success_callback(
    plc4c_system *system,
    plc4c_system_on_driver_load_success_callback callback);

/**
 * Function to set the error callback for loading drivers for the plc4c system.
 *
 * @param system
 * @param callback plc4c_system_callback_driver_load_error
 */
void plc4c_system_set_on_driver_load_failure_callback(
    plc4c_system *system,
    plc4c_system_on_driver_load_failure_callback callback);

/**
 * Function to set the on_connection callback for the plc4c system.
 *
 * @param system
 * @param callback plc4c_system_callback_on_connection
 */
void plc4c_system_set_on_connect_success_callback(
    plc4c_system *system, plc4c_system_on_connect_success_callback callback);

/**
 * Function to set the error callback for making connections for the plc4c
 * system.
 *
 * @param system
 * @param callback plc4c_system_on_connect_failure_callback
 */
void plc4c_system_set_on_connect_failure_callback(
    plc4c_system *system, plc4c_system_on_connect_failure_callback callback);

/**
 * Function to set the on_disconnection callback for the plc4c system.
 *
 * @param system
 * @param callback plc4c_system_callback_on_disconnection
 */
void plc4c_system_set_on_disconnect_success_callback(
    plc4c_system *system, plc4c_system_on_disconnect_success_callback callback);

/**
 * Function to set the error callback for shutting down connections for the
 * plc4c system.
 *
 * @param system
 * @param callback
 */
void plc4c_system_set_on_disconnect_failure_callback(
    plc4c_system *system, plc4c_system_on_disconnect_failure_callback callback);

/**
 * Function to set the error callback loops.
 *
 * @param system
 * @param callback plc4c_system_callback_loop_error
 */
void plc4c_system_set_on_loop_failure_callback(
    plc4c_system *system, plc4c_system_on_loop_failure_callback callback);

/**
 * Function to manually add a driver to the system.
 *
 * @param system the system the driver should be added to.
 * @param driver instance of the driver
 * @return plc4c_return_code
 */
plc4c_return_code plc4c_system_add_driver(plc4c_system *system,
                                          plc4c_driver *driver);

/**
 * Function to manually add a transport to the system.
 *
 * @param system the system the transport should be added to.
 * @param transport instance of the transport
 * @return plc4c_return_code
 */
plc4c_return_code plc4c_system_add_transport(plc4c_system *system,
                                             plc4c_transport *transport);

/**
 * Function to clean up the PLC4C system (Free any still used resources,
 * terminate live connections, ...)
 *
 * @param system
 */
void plc4c_system_shutdown(plc4c_system *system);

/**
 * Function to initialize a connection to a PLC by parsing the given connection
 * string and setting the passed connection system
 *
 * @param system
 * @param connection_string
 * @param ponter to where the connection object will be created.
 * @return plc4c_return_code INVALID_CONNECTION_STRING, NO_MEMORY
 */
plc4c_return_code plc4c_system_connect(plc4c_system *system,
                                       char *connectionString,
                                       plc4c_connection **connection);

/**
 * Function to give any drivers the chance to do their work.
 * In single-threaded environments we can't operate with event
 * handler loops as they would block the rest of the application.
 *
 * @param system the system instance
 * @return plc4c_return_code
 */
plc4c_return_code plc4c_system_loop(plc4c_system *system);

#ifdef __cplusplus
}
#endif
#endif  // PLC4C_SYSTEM_H_