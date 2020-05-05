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
    return connection->connected;
}

bool plc4c_connection_has_error(plc4c_connection *connection) {
    return false;
}

return_code plc4c_connection_disconnect(plc4c_connection *connection) {
    plc4c_system_task *new_disconnection_task = NULL;
    return_code result = connection->driver->disconnect_function(connection, &new_disconnection_task);
    if (result != OK) {
        return -1;
    }
    plc4c_utils_list_insert_tail_value(connection->system->task_list, new_disconnection_task);
    return OK;
}

void plc4c_connection_destroy(plc4c_connection *connection) {
    free(connection);
}

char *plc4c_connection_get_connection_string(plc4c_connection *connection) {
    return connection->connection_string;
}

bool plc4c_connection_supports_reading(plc4c_connection *connection) {
    return connection->supports_reading;
}

return_code plc4c_connection_create_read_request(plc4c_connection *connection,
                                                 plc4c_list *addresses,
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
            plc4c_item *item = connection->driver->parse_address_function((char *) element->value);
            plc4c_utils_list_insert_tail_value(new_read_request->items, item);
            element = element->next;
        } while (element != NULL);
    }
    *read_request = new_read_request;
    return OK;
}

void plc4c_connection_read_response_destroy(plc4c_read_response *read_response) {
    read_response->read_request->connection->driver->free_read_response_function(read_response);
}

bool plc4c_connection_supports_writing(plc4c_connection *connection) {
    return connection->supports_writing;
}

return_code plc4c_connection_create_write_request(plc4c_connection *connection,
                                                  plc4c_list *addresses,
                                                  plc4c_list *values,
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
    plc4c_write_request *new_write_request = (plc4c_write_request *) malloc(sizeof(plc4c_write_request));
    new_write_request->connection = connection;
    plc4c_utils_list_create(&(new_write_request->items));

    plc4c_list_element *address_element = plc4c_utils_list_head(addresses);
    plc4c_list_element *value_element = plc4c_utils_list_head(values);
    if (address_element != NULL && value_element != NULL) {
        do {
            char *address = (char *) address_element->value;
            // Parse an address string and get a driver-dependent data-structure representing the address back.
            plc4c_item *address_item = connection->driver->parse_address_function(address);

            // Create a new value item, binding an address item to a value.
            plc4c_value_item *value_item = malloc(sizeof(plc4c_value_item));
            value_item->item = address_item;
            connection->driver->encode_value_function(address_item, value_element->value, &value_item->value);

            // Add the new item ot the list of items.
            plc4c_utils_list_insert_tail_value(new_write_request->items, value_item);

            address_element = address_element->next;
            value_element = value_element->next;
        } while (address_element != NULL && value_element != NULL);
    }

    *write_request = new_write_request;
    return OK;
}

void plc4c_connection_write_response_destroy(plc4c_write_response *write_response) {
    write_response->write_request->connection->driver->free_write_response_function(write_response);
}
