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

struct plc4c_system_t {
  /* drivers */

  /* connections */

  /* callbacks */
};

struct plc4c_item_t {
};
typedef struct plc4c_item_t plc4c_item;

typedef plc4c_item* (*plc4c_connection_parse_address_item)(const char *address_string);

struct plc4c_driver_t {
    char* protocol_code;
    char* protocol_name;
    plc4c_connection_parse_address_item parse_address_function;
};

struct plc4c_connection_t {
    plc4c_driver driver;
    char* connection_string;
    bool supports_reading;
    bool supports_writing;
    bool supports_subscriptions;
};

struct plc4c_promise_t {
    return_code returnCode;
    plc4c_success_callback successCallback;
    plc4c_failure_callback failureCallback;
};

struct plc4c_read_request_t {
    plc4c_connection* connection;
    int num_items;
    plc4c_item items[];
};

struct plc_write_item_t {
    plc4c_item item;
    void* value;
};
typedef struct plc_write_item_t plc_write_item;

struct plc4c_write_request_t {
    plc4c_connection* connection;
    int num_items;
    plc_write_item items[];
};

#endif //PLC4C_SPI_TYPES_PRIVATE_H_
