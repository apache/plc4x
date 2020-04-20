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
#ifndef PLC4C_H_
#define PLC4C_H_
#ifdef __cplusplus
extern "C" {
#endif

#include <stdbool.h>


/**
 * Public API
 */

/**
 *
 * PLC4C error codes
*/
typedef enum error_code {
    OK,
    UNKNOWN_ERROR,
    NO_MEMORY,
    INVALID_CONNECTION_STRING,
    NOT_REACHABLE,
    PERMISSION_DENIED,
    INTERNAL_ERROR
} error_code;

/**
 * Helper that translates from an error_code enum value to something a human can work with.
 *
 * @param err error code.
 * @return A human readable error description.
 */
char *plc4c_error_code_to_error_message(error_code err);

/**
 * the plc4c system
 */
typedef struct plc4c_system_t plc4c_system;

/**
 * the plc4c_driver
 */
typedef struct plc4c_driver_t plc4c_driver;

/**
 * the plc4c_connection
 */
typedef struct plc4c_connection_t plc4c_connection;

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
error_code plc4c_init(plc4c_system *system);

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
 * Function to initialize a connection to a PLC by parsing the given connection string
 * and setting the passed connection system and also providing callbacks for success and
 * failure in case of problems.
 *
 * @param system
 * @param connectionString
 * @param connection
 * @param success_callback
 * @param error_callback
 * @return error_code INVALID_CONNECTION_STRING, NO_MEMORY
 */
error_code plc4c_system_connect_callback(plc4c_system *system,
                                         const char *connectionString,
                                         plc4c_connection **connection,
                                         plc4c_system_callback_on_connection success_callback,
                                         plc4c_system_callback_connection_error error_callback);

/**
 * Function to terminate a connection to a PLC.
 *
 * @param connection
 * @param plc4c_connection
 */
error_code plc4c_disconnect(plc4c_connection *connection);

/**
 * Function to terminate a connection to a PLC.
 *
 * @param connection
 * @param plc4c_connection
 */
error_code plc4c_disconnect_callback(plc4c_connection *connection,
                                     plc4c_system_callback_on_disconnection success_callback,
                                     plc4c_system_callback_disconnection_error error_callback);

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
#endif //PLC4C_H_