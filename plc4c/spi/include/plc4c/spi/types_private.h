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

#include <plc4c/data.h>
#include <plc4c/system.h>
#include <plc4c/types.h>
#include <plc4c/utils/list.h>
#include <plc4c/utils/queue.h>

#include "read_buffer.h"
#include "write_buffer.h"

typedef struct plc4c_item_t plc4c_item;
typedef struct plc4c_driver_list_item_t plc4c_driver_list_item;
typedef struct plc4c_transport_list_item_t plc4c_transport_list_item;
typedef struct plc4c_connection_list_item_t plc4c_connection_list_item;
typedef struct plc4c_request_value_item_t plc4c_request_value_item;
typedef struct plc4c_response_value_item_t plc4c_response_value_item;
typedef struct plc4c_response_item_t plc4c_response_item;
typedef struct plc4c_response_subscription_item_t plc4c_response_subscription_item;
typedef struct plc4c_response_unsubscription_item_t plc4c_response_unsubscription_item;

typedef plc4c_return_code (*plc4c_connection_configure_function)(
    plc4c_list* parameters, void** configuration);

typedef plc4c_return_code (*plc4c_connection_parse_address_item)(
    char *address_string, void** encoded_address);

typedef plc4c_return_code (*plc4c_connection_encode_value_item)(
    plc4c_item *item, void *value, void **encoded_value);

typedef struct plc4c_system_task_t plc4c_system_task;

typedef plc4c_return_code (*plc4c_system_task_state_machine_function)(
    plc4c_system_task *task);

typedef plc4c_return_code (*plc4c_connection_connect_function)(
    plc4c_connection *connection, plc4c_system_task **task);

typedef plc4c_return_code (*plc4c_connection_disconnect_function)(
    plc4c_connection *connection, plc4c_system_task **task);

typedef plc4c_return_code (*plc4c_connection_read_function)(
    plc4c_read_request_execution *read_request_execution,
    plc4c_system_task **task);

typedef plc4c_return_code (*plc4c_connection_write_function)(
    plc4c_write_request_execution *write_request_execution,
    plc4c_system_task **task);

typedef plc4c_return_code (*plc4c_connection_subscribe_function)(
    plc4c_subscription_request_execution *subscription_request_execution,
    plc4c_system_task **task);

typedef plc4c_return_code (*plc4c_connection_unsubscribe_function)(
    plc4c_unsubscription_request_execution *unsubscription_request_execution,
    plc4c_system_task **task);

typedef void (*plc4c_connect_free_read_request_function)(
    plc4c_read_request *request);

typedef void (*plc4c_connect_free_write_request_function)(
    plc4c_write_request *request);
    
typedef void (*plc4c_connect_free_read_response_function)(
    plc4c_read_response *response);

typedef void (*plc4c_connect_free_write_response_function)(
    plc4c_write_response *response);

typedef void (*plc4c_connect_free_subscription_response_function)(
    plc4c_subscription_response *response);

typedef void (*plc4c_connect_free_unsubscription_response_function)(
    plc4c_unsubscription_response *response);

typedef plc4c_return_code (*plc4c_transport_configure_function)(
    char* transport_connect_information, plc4c_list* parameters, void** configuration);

// TODO: Implement the argument.
typedef plc4c_return_code (*plc4c_transport_open_function)(void* config);

// TODO: Implement the argument.
typedef plc4c_return_code (*plc4c_transport_close_function)(void* config);

typedef plc4c_return_code (*plc4c_transport_send_message_function)(
    void* transport_configuration, plc4c_spi_write_buffer* message);

// Helper function that tells the transport what to do with the current input
// on positive return value: a read-buffer with given number of bytes is
// returned. RESPONSE_ACCEPT_INCOMPLETE or 0: Currently not enough data is
// or this content is not applicable for the current task. Don't consume
// anything and potentially re-check next time. If however the value is negative
// the buffer seems to contain corrupt data, clean up by skipping the number of
// bytes you get by making the negative value a positive and not returning any
// read-buffer.
typedef int16_t (*accept_message_function)(
    uint8_t* data, uint16_t length);

typedef plc4c_return_code (*plc4c_transport_select_message_function)(
    void* transport_configuration, uint8_t min_size, accept_message_function accept_message, plc4c_spi_read_buffer** message);

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
  char *name;
  void *address;
};

struct plc4c_driver_t {
  char *protocol_code;
  char *protocol_name;
  char *default_transport_code;
  plc4c_connection_configure_function configure_function;
  plc4c_connection_parse_address_item parse_address_function;
  plc4c_connection_connect_function connect_function;
  plc4c_connection_disconnect_function disconnect_function;
  plc4c_connection_read_function read_function;
  plc4c_connection_write_function write_function;
  plc4c_connection_subscribe_function subscribe_function;
  plc4c_connection_unsubscribe_function unsubscribe_function;
  plc4c_connect_free_read_request_function free_read_request_function;
  plc4c_connect_free_write_request_function free_write_request_function;
  plc4c_connect_free_read_response_function free_read_response_function;
  plc4c_connect_free_write_response_function free_write_response_function;
  plc4c_connect_free_subscription_response_function free_subscription_response_function;
  plc4c_connect_free_unsubscription_response_function free_unsubscription_response_function;
};

struct plc4c_driver_list_item_t {
  plc4c_driver *driver;
  plc4c_driver_list_item *next;
};

struct plc4c_transport_message_t {
  int length;
  uint8_t data[];
};

struct plc4c_transport_t {
  char *transport_code;
  char* transport_name;

  plc4c_transport_configure_function configure;
  plc4c_transport_open_function open;
  plc4c_transport_close_function close;
  plc4c_transport_send_message_function send_message;
  // Function that uses a function passed in to see if a given system-task
  // will be able to consume the current content of the transports input
  // buffer.
  plc4c_transport_select_message_function select_message;
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
  void *transport_configuration;
  char *parameters;
  void *configuration;

  bool connected;
  // Internal flag indicating the connection should be disconnected
  bool disconnect;
  // Number of system_tasks currently still active in the system.
  int num_running_system_tasks;

  plc4c_system *system;
  plc4c_driver *driver;
  plc4c_transport *transport;

  bool supports_reading;
  bool supports_writing;
  bool supports_subscriptions;
};

struct plc4c_connection_list_item_t {
  plc4c_connection *connection;
  plc4c_connection_list_item *prev;
  plc4c_connection_list_item *next;
};

struct plc4c_data_t {
  plc4c_data_type data_type;
  size_t size;
  union {
    bool bool_value;
    uint8_t byte_value;
    uint16_t word_value;
    uint32_t dword_value;
    uint64_t lword_value;

    uint8_t usint_value;
    int8_t sint_value;
    uint16_t uint_value;
    int16_t int_value;
    uint32_t udint_value;
    int32_t dint_value;
    uint64_t ulint_value;
    int64_t lint_value;

    float real_value;
    double lreal_value;

    uint16_t date_value;
    uint32_t time_value;
    uint64_t ltime_value;
    uint32_t time_of_day_value;
    uint32_t date_and_time_value;

    char char_value;
    wchar_t wchar_value;
    char *string_value;
    wchar_t *wstring_value;

    plc4c_list* list_value;
  } data;

  plc4c_data_custom_destroy custom_destroy;
  plc4c_data_custom_printf custom_printf;
};

struct plc4c_request_value_item_t {
  plc4c_item *item;
  plc4c_data *value;
};

struct plc4c_response_value_item_t {
  plc4c_item *item;
  plc4c_response_code response_code;
  plc4c_data *value;
};

struct plc4c_response_subscription_item_t {
  plc4c_item *item;
  plc4c_response_code response_code;
  // This is highly coupled to the protocol used.
  void *subscription_handle;
};

struct plc4c_request_unsubscription_item_t {
  plc4c_item *item;
  // This is highly coupled to the protocol used.
  void *subscription_handle;
};

struct plc4c_response_item_t {
  plc4c_item *item;
  plc4c_response_code response_code;
};

struct plc4c_read_request_t {
  plc4c_connection *connection;
  plc4c_list *items;
};

struct plc4c_write_request_t {
  plc4c_connection *connection;
  plc4c_list *items;
};

struct plc4c_subscription_request_t {
  plc4c_connection *connection;
  plc4c_list *items;
};

struct plc4c_unsubscription_request_t {
  plc4c_connection *connection;
  plc4c_list *items;
};

struct plc4c_read_request_execution_t {
  plc4c_read_request *read_request;
  plc4c_read_response *read_response;
  plc4c_list_element* cur_item;
  plc4c_system_task *system_task;
};

struct plc4c_write_request_execution_t {
  plc4c_write_request *write_request;
  plc4c_write_response *write_response;
  plc4c_system_task *system_task;
};

struct plc4c_subscription_request_execution_t {
  plc4c_subscription_request *subscription_request;
  plc4c_subscription_response *subscription_response;
  plc4c_system_task *system_task;
};

struct plc4c_unsubscription_request_execution_t {
  plc4c_unsubscription_request *unsubscription_request;
  plc4c_unsubscription_response *unsubscription_response;
  plc4c_system_task *system_task;
};

struct plc4c_read_response_t {
  plc4c_read_request *read_request;
  plc4c_list *items;
};

struct plc4c_write_response_t {
  plc4c_write_request *write_request;
  plc4c_list *response_items;
};

struct plc4c_subscription_response_t {
  plc4c_subscription_request *subscription_request;
  plc4c_list *response_items;
};

struct plc4c_subscription_event_t {
  plc4c_subscription_request *subscription_request;
  plc4c_list *event_items;
};

struct plc4c_unsubscription_response_t {
  plc4c_unsubscription_request *unsubscription_request;
  plc4c_list *response_items;
};

struct plc4c_system_task_t {
  int state_id;
  plc4c_system_task_state_machine_function state_machine_function;
  bool completed;

  void *context;
  // Reference to the connection that owns this task
  plc4c_connection *connection;
};

#endif  // PLC4C_SPI_TYPES_PRIVATE_H_
