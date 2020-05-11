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
#include <plc4c/driver_simulated.h>
#include <plc4c/plc4c.h>
#include <plc4c/transport_dummy.h>
#include <plc4c/utils/list.h>
#include <stdio.h>
#include <string.h>

#include "../../../spi/include/plc4c/spi/types_private.h"

int numOpenConnections = 0;

/**
 * Here we could implement something that keeps track of all open connections.
 * For example on embedded devices using the W5100 SPI Network device, this can
 * only handle 4 simultaneous connections.
 *
 * @param connection the connection that was just established
 */
void onGlobalConnect(plc4c_connection *cur_connection) {
  printf("Connected to %s",
         plc4c_connection_get_connection_string(cur_connection));
  numOpenConnections++;
}

void onGlobalDisconnect(plc4c_connection *cur_connection) {
  printf("Disconnected from %s",
         plc4c_connection_get_connection_string(cur_connection));
  numOpenConnections--;
}

void delete_address(plc4c_list_element *address_data_element) {
  // these are not malloc'd, no need to free
  address_data_element->value = NULL;
}

void delete_read_response_item(plc4c_list_element *response_read_item_element) {

}

void delete_write_response_item(
    plc4c_list_element *response_write_item_element) {}

enum plc4c_connection_state_t {
  CONNECTING,
  CONNECTED,
  READ_REQUEST_SENT,
  READ_RESPONSE_RECEIVED,
  WRITE_REQUEST_SENT,
  WRITE_RESPONSE_RECEIVED,
  DISCONNECTING,
  DISCONNECTED
};
typedef enum plc4c_connection_state_t plc4c_connection_state;

#pragma clang diagnostic push
#pragma ide diagnostic ignored "hicpp-multiway-paths-covered"

int main() {
  bool loop = true;
  plc4c_system *system = NULL;
  plc4c_connection *connection = NULL;
  plc4c_read_request *read_request = NULL;
  plc4c_write_request *write_request = NULL;
  plc4c_read_request_execution *read_request_execution = NULL;
  plc4c_write_request_execution *write_request_execution = NULL;

  // Create a new uninitialized plc4c_system
  printf("Creating new PLC4C System (Initializing inner data-structures) ... ");
  plc4c_return_code result = plc4c_system_create(&system);
  if (result != OK) {
    printf("FAILED\n");
    return -1;
  }
  printf("SUCCESS\n");

  // Manually register the "simulated" driver with the system.
  printf("Registering driver for the 'simulated' protocol ... ");
  plc4c_driver *simulated_driver = plc4c_driver_simulated_create();
  result = plc4c_system_add_driver(system, simulated_driver);
  if (result != OK) {
    printf("FAILED\n");
    return -1;
  }
  printf("SUCCESS\n");

  printf("Registering transport for the 'dummy' transport ... ");
  plc4c_transport *dummy_transport = plc4c_transport_dummy_create();
  result = plc4c_system_add_transport(system, dummy_transport);
  if (result != OK) {
    printf("FAILED\n");
    return -1;
  }
  printf("SUCCESS\n");

  // Initialize the plc4c_system (loading of drivers, setting up other stuff,
  // ...)
  printf(
      "Initializing the PLC4C system (Loading of drivers and transports) ... ");
  result = plc4c_system_init(system);
  if (result != OK) {
    printf("FAILED\n");
    return -1;
  }
  printf("SUCCESS\n");

  // Register the global callbacks.
  plc4c_system_set_on_connect_success_callback(system, &onGlobalConnect);
  plc4c_system_set_on_disconnect_success_callback(system, &onGlobalDisconnect);

  // Establish connections to remote devices
  // you may or may not care about the connection handle
  printf("Connecting to 'simulated://foo' ... ");
  result = plc4c_system_connect(system, "simulated://foo", &connection);
  if (result != OK) {
    printf("FAILED\n");
    return -1;
  }

  // Central program loop ...
  plc4c_connection_state state = CONNECTING;
  while (loop) {
    printf("* ");

    // Give plc4c a chance to do something.
    // This is where all I/O is being done.
    if (plc4c_system_loop(system) != OK) {
      printf("ERROR in the system loop\n");
      break;
    }

    // Depending on the current state, implement some logic.
    switch (state) {
      case CONNECTING: {
        // Check if the connection is established:
        if (plc4c_connection_get_connected(connection)) {
          printf("SUCCESS\n");
          state = CONNECTED;
        } else if (plc4c_connection_has_error(connection)) {
          printf("FAILED\n");
          return -1;
        }
        break;
      }
      case CONNECTED: {
        // Create a new read-request.
        printf("Preparing a read-request for 'RANDOM/foo:INTEGER' ... ");

        plc4c_list *address_list = NULL;
        plc4c_utils_list_create(&address_list);
        plc4c_utils_list_insert_head_value(address_list,
                                           (void *)"RANDOM/foo:INTEGER");
        result = plc4c_connection_create_read_request(connection, address_list,
                                                      &read_request);
        if (result != OK) {
          printf("FAILED\n");
          return -1;
        }
        printf("SUCCESS\n");

        // Execute the read-request.
        printf("Executing a read-request ... ");
        result =
            plc4c_read_request_execute(read_request, &read_request_execution);

        // As we only used these to create the request, they can now be released
        // again.
        plc4c_utils_list_delete_elements(address_list, &delete_address);
        free(address_list);

        if (result != OK) {
          printf("FAILED\n");
          return -1;
        } else {
          state = READ_REQUEST_SENT;
        }
        break;
      }
        // Wait until the read-request execution is finished.
      case READ_REQUEST_SENT: {
        if (plc4c_read_request_execution_check_finished_successfully(
                read_request_execution)) {
          printf("SUCCESS\n");
          state = READ_RESPONSE_RECEIVED;
        } else if (plc4c_read_request_execution_check_finished_with_error(
                       read_request_execution)) {
          printf("FAILED\n");
          return -1;
        }
        break;
      }
      case READ_RESPONSE_RECEIVED: {
        // Get the response for the given read-request.
        plc4c_read_response *read_response =
            plc4c_read_request_execution_get_response(read_request_execution);
        if (read_response == NULL) {
          printf("FAILED (No Response)\n");
          return -1;
        }

        // Iterate over all returned items.
        plc4c_list_element *cur_element =
            plc4c_utils_list_head(read_response->items);
        while (cur_element != NULL) {
          plc4c_response_value_item *value_item = cur_element->value;

          printf("Value %s (%s):", value_item->item->name,
                 plc4c_response_code_to_message(value_item->response_code));
          plc4c_data_printf(value_item->value);
          printf("\n");

          cur_element = cur_element->next;
        }

        // Clean up.
        plc4c_read_destroy_read_response(read_response);
        plc4c_read_request_execution_destroy(read_request_execution);
        plc4c_read_request_destroy(read_request);

        // Create a new write-request.
        printf("Preparing a write-request for 'STDOUT/foo:INTEGER' ... ");
        plc4c_list *address_list = NULL;
        plc4c_utils_list_create(&address_list);
        plc4c_utils_list_insert_head_value(address_list,
                                           (void *)"STDOUT/foo:STRING");
        plc4c_list *value_list = NULL;
        plc4c_utils_list_create(&value_list);
        char value[] = "bar";
        plc4c_utils_list_insert_head_value(
            value_list,
            plc4c_data_create_constant_string_data(strlen(value), value));
        result = plc4c_connection_create_write_request(
            connection, address_list, value_list, &write_request);

        // As we only used these to create the request, they can now be released
        // again.
        plc4c_utils_list_delete_elements(address_list, &delete_address);
        plc4c_utils_list_delete_elements(value_list, &delete_address);
        free(address_list);
        free(value_list);

        if (result != OK) {
          printf("FAILED\n");
          return -1;
        }
        printf("SUCCESS\n");

        // Execute the write-request.
        printf("Executing a write-request ... \n");
        result = plc4c_write_request_execute(write_request,
                                             &write_request_execution);
        if (result != OK) {
          printf("FAILED\n");
          return -1;
        } else {
          state = WRITE_REQUEST_SENT;
        }
        break;
      }
        // Wait until the write-request execution is finished.
      case WRITE_REQUEST_SENT: {
        if (plc4c_write_request_check_finished_successfully(
                write_request_execution)) {
          printf("SUCCESS\n");
          state = WRITE_RESPONSE_RECEIVED;
        } else if (plc4c_write_request_execution_check_completed_with_error(
                       write_request_execution)) {
          printf("FAILED\n");
          return -1;
        }
        break;
      }
      case WRITE_RESPONSE_RECEIVED: {
        plc4c_write_response *write_response =
            plc4c_write_request_execution_get_response(write_request_execution);

        // Iterate over the responses ...
        plc4c_list_element *cur_element =
            plc4c_utils_list_head(write_response->response_items);
        while (cur_element != NULL) {
          plc4c_response_item *response_item = cur_element->value;
          printf(" - Write Value %s (%s)\n", response_item->item->name,
                 plc4c_response_code_to_message(response_item->response_code));
          cur_element = cur_element->next;
        }

        // Clean up.
        plc4c_write_destroy_write_response(write_response);
        plc4c_write_request_execution_destroy(write_request_execution);

        // Disconnect.
        printf("Disconnecting ... ");
        result = plc4c_connection_disconnect(connection);
        if (result != OK) {
          printf("FAILED");
          return -1;
        }
        state = DISCONNECTING;
        break;
      }
        // Wait until the connection is disconnected
      case DISCONNECTING: {
        if (!plc4c_connection_get_connected(connection)) {
          printf("SUCCESS\n");
          // we could let the system shut this down,
          // or do it ourselves
          plc4c_system_remove_connection(system, connection);
          plc4c_connection_destroy(connection);
          state = DISCONNECTED;

          // Terminate the main program loop.
          loop = false;
        }
        break;
      }
      case DISCONNECTED: {
        // End the loop.
        loop = false;
        break;
      }
    }
  }

  // Make sure everything is cleaned up correctly.
  plc4c_system_shutdown(system);

  // Finally destroy the plc4c_system, freeing up all memory allocated by plc4c.
  plc4c_system_destroy(system);

  return 0;
}

#pragma clang diagnostic pop
