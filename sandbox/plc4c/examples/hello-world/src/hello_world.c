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
//#include <plc4c/driver_simulated.h>
#include <plc4c/driver_s7.h>
//#include <plc4c/driver_modbus.h>
#include <plc4c/plc4c.h>
//#include <plc4c/transport_dummy.h>
#include <plc4c/transport_tcp.h>
//#include <plc4c/transport_serial.h>
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
  /*printf("Registering driver for the 'simulated' protocol ... ");
  plc4c_driver *simulated_driver = plc4c_driver_simulated_create();
  result = plc4c_system_add_driver(system, simulated_driver);
  if (result != OK) {
    printf("FAILED adding simulated driver\n");
    return -1;
  }*/
  printf("Registering driver for the 's7' protocol ... ");
  plc4c_driver *s7_driver = plc4c_driver_s7_create();
  result = plc4c_system_add_driver(system, s7_driver);
  if (result != OK) {
    printf("FAILED adding s7 driver\n");
    return -1;
  }
  /*printf("Registering driver for the 'modbus' protocol ... ");
  plc4c_driver *modbus_driver = plc4c_driver_modbus_create();
  result = plc4c_system_add_driver(system, modbus_driver);
  if (result != OK) {
    printf("FAILED adding modbus driver\n");
    return -1;
  }*/
  printf("SUCCESS\n");

  /*printf("Registering transport for the 'dummy' transport ... ");
  plc4c_transport *dummy_transport = plc4c_transport_dummy_create();
  result = plc4c_system_add_transport(system, dummy_transport);
  if (result != OK) {
    printf("FAILED adding dummy transport\n");
    return -1;
  }*/
  plc4c_transport *tcp_transport = plc4c_transport_tcp_create();
  result = plc4c_system_add_transport(system, tcp_transport);
  if (result != OK) {
    printf("FAILED adding tcp transport\n");
    return -1;
  }
  /*plc4c_transport *serial_transport = plc4c_transport_serial_create();
  result = plc4c_system_add_transport(system, serial_transport);
  if (result != OK) {
    printf("FAILED adding serial transport\n");
    return -1;
  }*/
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
  printf("Connecting to 's7:tcp://192.168.23.30' ... ");
  result = plc4c_system_connect(system, "s7:tcp://192.168.23.30", &connection);
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

        printf("Adding an item for '%I0.0:BOOL' ... ");
        result =
            plc4c_read_request_add_item(read_request, "BOOL", "%DB4:0.0:BOOL");
        if (result != OK) {
          printf("FAILED\n");
          return -1;
        }
        result =
            plc4c_read_request_add_item(read_request, "BYTE", "%DB4:1:BYTE");
        if (result != OK) {
          printf("FAILED\n");
          return -1;
        }
        result =
            plc4c_read_request_add_item(read_request, "WORD", "%DB4:2:WORD");
        if (result != OK) {
          printf("FAILED\n");
          return -1;
        }
        result =
            plc4c_read_request_add_item(read_request, "DWORD", "%DB4:4:DWORD");
        if (result != OK) {
          printf("FAILED\n");
          return -1;
        }
        result =
            plc4c_read_request_add_item(read_request, "SINT", "%DB4:16:SINT");
        if (result != OK) {
          printf("FAILED\n");
          return -1;
        }
        result =
            plc4c_read_request_add_item(read_request, "USINT", "%DB4:17:USINT");
        if (result != OK) {
          printf("FAILED\n");
          return -1;
        }
        result =
            plc4c_read_request_add_item(read_request, "INT", "%DB4:18:INT");
        if (result != OK) {
          printf("FAILED\n");
          return -1;
        }
        result =
            plc4c_read_request_add_item(read_request, "UINT", "%DB4:20:UINT");
        if (result != OK) {
          printf("FAILED\n");
          return -1;
        }
        result =
            plc4c_read_request_add_item(read_request, "DINT", "%DB4:22:DINT");
        if (result != OK) {
          printf("FAILED\n");
          return -1;
        }
        result =
            plc4c_read_request_add_item(read_request, "UDINT", "%DB4:26:UDINT");
        if (result != OK) {
          printf("FAILED\n");
          return -1;
        }
        result =
            plc4c_read_request_add_item(read_request, "REAL", "%DB4:46:REAL");
        if (result != OK) {
          printf("FAILED\n");
          return -1;
        }
        // S7 reports "Not supported"
        /*result =
            plc4c_read_request_add_item(read_request, "TIME", "%DB4:58:TIME");
        if (result != OK) {
          printf("FAILED\n");
          return -1;
        }
        result =
            plc4c_read_request_add_item(read_request, "DATE", "%DB4:70:DATE");
        if (result != OK) {
          printf("FAILED\n");
          return -1;
        }
        result =
            plc4c_read_request_add_item(read_request, "TIME_OF_DAY", "%DB4:72:TIME_OF_DAY");
        if (result != OK) {
          printf("FAILED\n");
          return -1;
        }*/
        result =
            plc4c_read_request_add_item(read_request, "TOD", "%DB4:76:TOD");
        if (result != OK) {
          printf("FAILED\n");
          return -1;
        }
        /*result =
            plc4c_read_request_add_item(read_request, "CHAR", "%DB4:136:CHAR");
        if (result != OK) {
          printf("FAILED\n");
          return -1;
        }
        result =
            plc4c_read_request_add_item(read_request, "WCHAR", "%DB4:138:WCHAR");
        if (result != OK) {
          printf("FAILED\n");
          return -1;
        }
        result =
            plc4c_read_request_add_item(read_request, "STRING", "%DB4:140:STRING(10)");
        if (result != OK) {
          printf("FAILED\n");
          return -1;
        }
        result =
            plc4c_read_request_add_item(read_request, "WSTRING", "%DB4:396:WSTRING(10)");
        if (result != OK) {
          printf("FAILED\n");
          return -1;
        }*/
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

        // TODO: Comment out after implementing the write functionality.
        // Disconnect.
        printf("Disconnecting ... ");
        result = plc4c_connection_disconnect(connection);
        if (result != OK) {
          printf("FAILED");
          return -1;
        }
        state = DISCONNECTING;

        // Create a new write-request.
/*        printf("Preparing a write-request ... ");
        char value[] = "bar";
        result =
            plc4c_connection_create_write_request(connection, &write_request);
        if (result != OK) {
          printf("FAILED\n");
          return -1;
        }
        printf("SUCCESS\n");

        printf("Adding an item for 'STDOUT/foo:INTEGER' ... ");
        result = plc4c_write_request_add_item(
            write_request, "STDOUT/foo:STRING",
            plc4c_data_create_constant_string_data(strlen(value), value));
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
        }*/
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
