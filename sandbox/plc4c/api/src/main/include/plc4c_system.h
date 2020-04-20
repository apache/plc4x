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
#ifndef PLC4C_SYSTEM_H_
#define PLC4C_SYSTEM_H_
#ifdef __cplusplus
extern "C" {
#endif

#include <stdbool.h>
#include "plc4c_types.h"

/**
 * SYSTEM CALLBACKS
 */

/**
 * Function pointer for a callback called when a driver is loaded.
 * Set in plc4c_system @see plc4c_system_set_on_driver_loaded()
 * @param driver
 */
typedef void (*plc4c_system_callback_on_driver_loaded)(plc4c_driver *driver);

/**
 * Function pointer for a callback called when loading a driver fails.
 * Set in plc4c_system @see plc4c_system_set_on_driver_load_error
 * NOTE: driver_name could be a pointer to the configuration for the driver instead....
 * @param driver_name
 * @param error_code
 */
typedef void (*plc4c_system_callback_driver_load_error)(const char *driver_name, error_code error);

/**
 * Function pointer for a callback called when is successfully made
 * Set in plc4c_system @see plc4c_system_set_on_connection()
 * @param connection
 */
typedef void (*plc4c_system_callback_on_connection)(plc4c_connection *connection);

/**
 * Function pointer for a callback called when connecting fails.
 * Set in plc4c_system @see plc4c_system_set_on_connection_error
 * @param connection_string
 * @param error_code
 */
typedef void (*plc4c_system_callback_connection_error)(const char *connection_string, error_code error);

/**
 * Function pointer for a callback called when is successfully made
 * Set in plc4c_system @see plc4c_system_set_on_connection()
 * @param connection
 */
typedef void (*plc4c_system_callback_on_disconnection)(plc4c_connection *connection);

/**
 * Function pointer for a callback called when connecting fails.
 * Set in plc4c_system @see plc4c_system_set_on_connection_error
 * @param connection
 * @param error_code
 */
typedef void (*plc4c_system_callback_disconnection_error)(plc4c_connection *connection, error_code error);

/**
 * Function pointer for a callback called when a driver returns an error
 * @param driver
 * @param connection
 * @param error_code
 */
typedef void(*plc4c_system_callback_loop_error)
        (plc4c_driver *driver, plc4c_connection *connection, error_code error);

/**
 * OTHER FUNCTION DEFS FOR SYSTEM
 */


/**
 * SYSTEM FUNCTIONS
 */

/**
 * Function to create a plc4c_system
 * @param system
 * @return NO_MEMORY if failed to create system
 */
error_code plc4c_system_create(plc4c_system **system);

/**
 * Function to destroy a plc4c_system
 * This will also destroy all connections associated with the system
 * @param system
 */
void plc4c_system_destroy(plc4c_system *system);

/**
 * Function to set the on_driver_loaded callback for the plc4c system
 * @param system
 * @param callback plc4c_system_callback_on_driver
 */
void plc4c_system_set_on_driver_loaded(plc4c_system *system,
                                       plc4c_system_callback_on_driver_loaded callback);

/**
 * Function to set the error callback for loading drivers for the plc4c system
 * @param system
 * @param callback plc4c_system_callback_driver_load_error
 */
void plc4c_system_set_on_driver_load_error(plc4c_system *system,
                                           plc4c_system_callback_driver_load_error callback);

/**
 * Function to set the on_connection callback for the plc4c system
 * @param system
 * @param callback plc4c_system_callback_on_connection
 */
void plc4c_system_set_on_connection(plc4c_system *system,
                                    plc4c_system_callback_on_connection callback);

/**
 * Function to set the error callback for making connections for the plc4c system
 * @param system
 * @param callback plc4c_system_callback_connection_error
 */
void plc4c_system_set_on_connection_error(plc4c_system *system,
                                          plc4c_system_callback_connection_error callback);

/**
 * Function to set the on_disconnection callback for the plc4c system
 * @param system
 * @param callback plc4c_system_callback_on_disconnection
 */
void plc4c_system_set_on_disconnection(plc4c_system *system,
                                       plc4c_system_callback_on_disconnection callback);

/**
 * Function to set the error callback for shutting down connections for the plc4c system
 * @param system
 * @param callback
 */
void plc4c_system_set_on_disconnection_error(plc4c_system *system,
                                             plc4c_system_callback_disconnection_error callback);

/**
 * Function to set the error callback loops
 * @param system
 * @param callback plc4c_system_callback_loop_error
 */
void plc4c_system_set_on_loop_error(plc4c_system *system,
                                    plc4c_system_callback_loop_error callback);

/**
 * Function to initialize the PLC4C system (Initialize the driver manager and the list of enabled drivers)
 * @param system
 * @return error_code
 */
error_code plc4c_system_init(plc4c_system *system);

/**
 * Function to clean up the PLC4C system (Free any still used resources, terminate live connections, ...)
 * @param system
 */
void plc4c_system_shutdown(plc4c_system *system);

/**
 * Function to initialize a connection to a PLC by parsing the given connection string
 * and setting the passed connection system
 *
 * @param system
 * @param connectionString
 * @return error_code INVALID_CONNECTION_STRING, NO_MEMORY
 */
error_code plc4c_system_connect(plc4c_system *system,
                                const char *connectionString,
                                plc4c_connection **connection);

/**
 * Function to give any drivers the chance to do their work.
 * In single-threaded environments we can't operate with event
 * handler loops as they would block the rest of the application.
 *
 * @return error_code
 */
error_code plc4c_system_loop();

#ifdef __cplusplus
}
#endif
#endif //PLC4C_SYSTEM_H_