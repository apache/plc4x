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

/**
 * CONNECTION CALLBACKS
 */

/**
 * OTHER FUNCTION DEFS FOR CONNECTION
 */

/**
 * CONNECTION FUNCTIONS
 */

/**
 * Function to terminate a connection to a PLC.
 *
 * @param connection
 * @param plc4c_promise
 */
plc4c_promise* plc4c_connection_disconnect(plc4c_connection *connection);

/**
 * Get the connection string from a given connection.
 */
char* plc4c_connection_get_connection_string(plc4c_connection *connection);

/**
 * Check if the current connection supports read operations.
 * @param connection reference to the connection
 * @return true if the connection supports reading, false otherwise
 */
bool plc4c_connection_supports_reading(plc4c_connection *connection);

/**
 * Initializes an empty read-request
 * @param connection connection that this read-request will be executed on.
 * @param num_items number of items we want to read.
 * @param addresses array of address strings.
 * @return pointer to new read-request.
 */
plc4c_read_request* plc4c_connection_create_read_request(plc4c_connection *connection, int num_items, char* addresses[]);

/**
 * Check if the current connection supports write operations.
 * @param connection reference to the connection
 * @return true if the connection supports writing, false otherwise
 */
bool plc4c_connection_supports_writing(plc4c_connection *connection);

/**
 * Initializes an empty write-request
 * @param connection connection that this write-request will be executed on.
 * @param num_items number of items we want to write.
 * @param addresses array of address strings.
 * @param values array of pointers to values.
 * @return pointer to new write_request
 */
plc4c_write_request* plc4c_connection_create_write_request(plc4c_connection *connection, int num_items, char* addresses[], void* values[]);

/**
 * Check if the current connection supports subscriptions.
 * @param connection reference to the connection
 * @return true if the connection supports subscriptions, false otherwise
 */
bool plc4c_connection_supports_subscriptions(plc4c_connection *connection);

#ifdef __cplusplus
}
#endif
#endif //PLC4C_CONNECTION_H_