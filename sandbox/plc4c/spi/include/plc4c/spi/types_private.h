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
#ifndef PLC4C_SPI_TYPES_PRIVATE_H_
#define PLC4C_SPI_TYPES_PRIVATE_H_

#include <plc4c/types.h>
#include <plc4c/system.h>
#include <plc4c/utils/list.h>
#include <plc4c/utils/queue.h>

typedef struct plc4c_item_t plc4c_item;
typedef struct plc4c_driver_list_item_t plc4c_driver_list_item;
typedef struct plc4c_transport_list_item_t plc4c_transport_list_item;
typedef struct plc4c_connection_list_item_t plc4c_connection_list_item;
typedef struct plc4c_write_item_t plc4c_write_item;

typedef plc4c_item *(*plc4c_connection_parse_address_item)(const char *address_string);

typedef struct plc4c_system_task_t plc4c_system_task;
typedef plc4c_item (*plc4c_system_task_state_machine_function)(plc4c_system_task* task);
typedef return_code (*plc4c_connection_connect_function)(plc4c_system *system, plc4c_connection *connection, plc4c_system_task** task);
typedef return_code (*plc4c_connection_read_function)(plc4c_system *system, plc4c_connection *connection, plc4c_read_request *read_request, plc4c_system_task** task);
typedef return_code (*plc4c_connection_write_function)(plc4c_system *system, plc4c_connection *connection, plc4c_write_request *write_request, plc4c_system_task** task);

struct plc4c_system_t {
    /* drivers */
    plc4c_list *driver_list;

    /* transports */
    plc4c_list *transport_list;

    /* connections */
    plc4c_list *connection_list;

    /* tasks */
    plc4c_list *task_list;

    /* callbacks */
    plc4c_system_on_driver_load_success_callback on_driver_load_success_callback;
    plc4c_system_on_driver_load_failure_callback on_driver_load_failure_callback;
    plc4c_system_on_connect_success_callback on_connect_success_callback;
    plc4c_system_on_connect_failure_callback on_connect_failure_callback;
    plc4c_system_on_disconnect_success_callback on_disconnect_success_callback;
    plc4c_system_on_disconnect_failure_callback on_disconnect_failure_callback;
    plc4c_system_on_loop_failure_callback on_loop_failure_callback;
};

struct plc4c_item_t {
};

struct plc4c_driver_t {
    char *protocol_code;
    char *protocol_name;
    char *default_transport_code;
    plc4c_connection_parse_address_item parse_address_function;
    plc4c_connection_connect_function connect_function;
    plc4c_connection_read_function read_function;
    plc4c_connection_write_function write_function;
};

struct plc4c_driver_list_item_t {
    plc4c_driver *driver;
    plc4c_driver_list_item *next;
};


struct plc4c_transport_t {
    char *transport_code;

    // TODO: add the send and receive function references here ...
};

struct plc4c_transport_list_item_t {
    plc4c_transport *transport;
    plc4c_transport_list_item *next;
};


struct plc4c_connection_t {
    char *connection_string;
    char *protocol_code;
    char *transport_code;
    char *transport_connect_information;
    char *parameters;

    plc4c_driver* driver;
    plc4c_transport* transport;
    bool supports_reading;
    bool supports_writing;
    bool supports_subscriptions;
};

struct plc4c_connection_list_item_t {
    plc4c_connection *connection;
    plc4c_connection_list_item *prev;
    plc4c_connection_list_item *next;
};


struct plc4c_read_request_t {
    plc4c_connection *connection;
    int num_items;
    plc4c_item items[];
};

struct plc4c_write_item_t {
    plc4c_item *item;
    void *value;
};

struct plc4c_write_request_t {
    plc4c_connection *connection;
    int num_items;
    plc4c_write_item *items;
};

struct plc4c_system_task_t {
    int state_id;
    plc4c_system_task_state_machine_function state_machine_function;
    void* context;
    bool completed;
};

#endif //PLC4C_SPI_TYPES_PRIVATE_H_
