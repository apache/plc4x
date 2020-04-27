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

#include <stdlib.h>
#include <plc4c/connection.h>
#include <plc4c/spi/types_private.h>

plc4c_promise* plc4c_connection_disconnect(plc4c_connection *connection) {
    plc4c_promise* result = (plc4c_promise*) malloc(sizeof(plc4c_promise));
    result->returnCode = UNFINISHED;
    return result;
}

char* plc4c_connection_get_connection_string(plc4c_connection *connection) {
    return connection->connection_string;
}

bool plc4c_connection_supports_reading(plc4c_connection *connection) {
    return connection->supports_reading;
}

plc4c_read_request* plc4c_connection_create_read_request(plc4c_connection *connection, int num_items, char* addresses[]) {
    plc4c_read_request* read_request = (plc4c_read_request*) malloc(sizeof(plc4c_read_request));
    return read_request;
}

bool plc4c_connection_supports_writing(plc4c_connection *connection) {
    return connection->supports_writing;
}

plc4c_write_request* plc4c_connection_create_write_request(plc4c_connection *connection, int num_items, char* addresses[], void* values[]) {
    plc4c_write_request* write_request = (plc4c_write_request*) malloc(sizeof(plc4c_write_request));
    return write_request;
}
