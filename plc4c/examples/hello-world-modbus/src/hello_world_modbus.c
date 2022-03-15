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
#include <plc4c/driver_modbus.h>
#include <plc4c/plc4c.h>
#include <plc4c/transport_tcp.h>
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

enum plc4c_connection_state_t {
  CONNECTING,
  CONNECTED,
  READ_REQUEST_SENT,
  READ_RESPONSE_RECEIVED,
  DISCONNECTING,
  DISCONNECTED
};
typedef enum plc4c_connection_state_t plc4c_connection_state;

//#pragma clang diagnostic push
//#pragma ide diagnostic ignored "hicpp-multiway-paths-covered"

int main() {
  bool loop = true;
  plc4c_system *system = NULL;
  plc4c_connection *connection = NULL;
  plc4c_read_request *read_request = NULL;
  plc4c_read_request_execution *read_request_execution = NULL;

  // Create a new uninitialized plc4c_system
  printf("Creating new PLC4C System (Initializing inner data-structures) ... ");
  plc4c_return_code result = plc4c_system_create(&system);
  if (result != OK) {
    printf("FAILED\n");
    return -1;
  }
  printf("SUCCESS\n");

  // Manually register the "simulated" driver with the system.
  printf("Registering driver for the 'modbus' protocol ... ");
  plc4c_driver *modbus_driver = plc4c_driver_modbus_create();
  result = plc4c_system_add_driver(system, modbus_driver);
  if (result != OK) {
    printf("FAILED adding modbus driver\n");
    return -1;
  }
  printf("SUCCESS\n");

  plc4c_transport *tcp_transport = plc4c_transport_tcp_create();
  result = plc4c_system_add_transport(system, tcp_transport);
  if (result != OK) {
    printf("FAILED adding tcp transport\n");
    return -1;
  }
  printf("SUCCESS\n");

  // Initialize the plc4c_system
  // (loading of drivers, setting up other stuff ...)
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
  printf("Connecting to 'modbus-tcp://192.168.23.30' ... ");
  result = plc4c_system_connect(system, "modbus-tcp://192.168.23.30:502", &connection);
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
        printf("Preparing a read-request ... ");
        result =
            plc4c_connection_create_read_request(connection, &read_request);
        if (result != OK) {
          printf("FAILED\n");
          return -1;
        }
        printf("SUCCESS\n");

        result =
            plc4c_read_request_add_item(read_request, "BOOL", "holding-register:1:BOOL");
        if (result != OK) {
          printf("FAILED\n");
          return -1;
        }
        result =
            plc4c_read_request_add_item(read_request, "BYTE", "holding-register:2:BYTE");
        if (result != OK) {
          printf("FAILED\n");
          return -1;
        }
        result =
            plc4c_read_request_add_item(read_request, "WORD", "holding-register:3:WORD");
        if (result != OK) {
          printf("FAILED\n");
          return -1;
        }
        result =
            plc4c_read_request_add_item(read_request, "DWORD", "holding-register:4:DWORD");
        if (result != OK) {
          printf("FAILED\n");
          return -1;
        }
        result =
            plc4c_read_request_add_item(read_request, "LWORD", "holding-register:6:LWORD");
        if (result != OK) {
          printf("FAILED\n");
          return -1;
        }
        result =
            plc4c_read_request_add_item(read_request, "SINT", "holding-register:10:SINT");
        if (result != OK) {
          printf("FAILED\n");
          return -1;
        }
        result =
            plc4c_read_request_add_item(read_request, "USINT", "holding-register:11:USINT");
        if (result != OK) {
          printf("FAILED\n");
          return -1;
        }
        result =
            plc4c_read_request_add_item(read_request, "INT", "holding-register:12:INT");
        if (result != OK) {
          printf("FAILED\n");
          return -1;
        }
        result =
            plc4c_read_request_add_item(read_request, "UINT", "holding-register:13:UINT");
        if (result != OK) {
          printf("FAILED\n");
          return -1;
        }
        result =
            plc4c_read_request_add_item(read_request, "DINT", "holding-register:14:DINT");
        if (result != OK) {
          printf("FAILED\n");
          return -1;
        }
        result =
            plc4c_read_request_add_item(read_request, "UDINT", "holding-register:16:UDINT");
        if (result != OK) {
          printf("FAILED\n");
          return -1;
        }
        result =
            plc4c_read_request_add_item(read_request, "LINT", "holding-register:18:LINT");
        if (result != OK) {
          printf("FAILED\n");
          return -1;
        }
        result =
            plc4c_read_request_add_item(read_request, "ULINT", "holding-register:22:ULINT");
        if (result != OK) {
          printf("FAILED\n");
          return -1;
        }
        result =
            plc4c_read_request_add_item(read_request, "REAL", "holding-register:26:REAL");
        if (result != OK) {
          printf("FAILED\n");
          return -1;
        }
        result =
            plc4c_read_request_add_item(read_request, "REAL", "holding-register:28:LREAL");
        if (result != OK) {
          printf("FAILED\n");
          return -1;
        }
        // UTF-16 doesn't work well in a portable way
        printf("SUCCESS\n");

        // Execute the read-request.
        printf("Executing a read-request ... ");
        result =
            plc4c_read_request_execute(read_request, &read_request_execution);

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
            plc4c_utils_list_tail(read_response->items);
        while (cur_element != NULL) {
          plc4c_response_value_item *value_item = cur_element->value;

          printf("Value %s (%s): ", value_item->item->name,
                 plc4c_response_code_to_message(value_item->response_code));
          plc4c_data_printf(value_item->value);
          if (strcmp(value_item->item->name, "TIME") == 0) {
            printf(" (Duration in ms)");
          } else if (strcmp(value_item->item->name, "DATE") == 0) {
            printf(" (Days since 1990-01-01)");
          } else if ((strcmp(value_item->item->name, "TIME_OF_DAY") == 0) ||
                     (strcmp(value_item->item->name, "TOD") == 0)) {
            printf(" (Number of ms since midnight)");
          }
          printf("\n");

          cur_element = cur_element->next;
        }

        // Clean up.
        plc4c_read_response_destroy(read_response);
        plc4c_read_request_execution_destroy(read_request_execution);
        plc4c_read_request_destroy(read_request);

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
      default: {
      }
    }
  }

  // Make sure everything is cleaned up correctly.
  plc4c_system_shutdown(system);

  // Finally destroy the plc4c_system, freeing up all memory allocated by plc4c.
  plc4c_system_destroy(system);

  return 0;
}

//#pragma clang diagnostic pop
