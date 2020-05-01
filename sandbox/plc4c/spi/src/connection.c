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

bool plc4c_connection_is_connected(plc4c_connection *connection) {
    return true;
}

bool plc4c_connection_has_error(plc4c_connection *connection) {
    return false;
}

return_code plc4c_connection_disconnect(plc4c_connection *connection) {
    return OK;
}

void plc4c_connection_destroy(plc4c_connection *connection) {
    free(connection);
}

char* plc4c_connection_get_connection_string(plc4c_connection *connection) {
    return connection->connection_string;
}

bool plc4c_connection_supports_reading(plc4c_connection *connection) {
    return connection->supports_reading;
}

return_code plc4c_connection_create_read_request(plc4c_connection *connection, int num_items, char* addresses[], plc4c_read_request** read_request) {
    plc4c_read_request *new_read_request = malloc(sizeof(plc4c_read_request));
    new_read_request->connection = connection;
    plc4c_utils_list_create(&(new_read_request->items));
    for(int i = 0; i < num_items; i++) {
        plc4c_item *item = connection->driver->parse_address_function(addresses[i]);
        plc4c_utils_list_insert_tail_value(new_read_request->items, item);
    }
    *read_request = new_read_request;
    return OK;
}

bool plc4c_connection_supports_writing(plc4c_connection *connection) {
    return connection->supports_writing;
}

return_code plc4c_connection_create_write_request(plc4c_connection *connection, int num_items, char* addresses[], void* values[], plc4c_write_request** write_request) {
    plc4c_write_request* new_write_request = (plc4c_write_request*) malloc(sizeof(plc4c_write_request));
    new_write_request->num_items = num_items;
    new_write_request->items = malloc(num_items * sizeof(plc4c_write_item*));
    for(int i = 0; i < num_items; i++) {
        char* address = addresses[i];
        plc4c_item* addressItem = connection->driver->parse_address_function(address);

        plc4c_write_item* write_item = malloc(sizeof(plc4c_write_item));
        write_item->item = addressItem;
        write_item->value = values[i];

        new_write_request->items = write_item;
    }
    write_request = &new_write_request;

    return OK;
}
