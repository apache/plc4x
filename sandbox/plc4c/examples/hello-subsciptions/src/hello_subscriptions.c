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

enum plc4c_connection_state_t {
  CONNECTING,
  CONNECTED,
  SUBSCRIPTION_REQUEST_SENT,
  SUBSCRIPTION_RESPONSE_RECEIVED,
  READING_EVENTS,
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
  plc4c_subscription_request *subscription_request = NULL;
  plc4c_subscription_request_execution *subscription_request_execution = NULL;
  int num_events = 0;
  void *subscription_handle = NULL;

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
        // Create a new subscription-request.
        printf("Preparing a subscription-request ... ");
        result = plc4c_connection_create_subscription_request(
            connection, &subscription_request);
        if (result != OK) {
          printf("FAILED\n");
          return -1;
        }
        printf("SUCCESS\n");

        printf("Adding cyclic item for 'RANDOM/foo:INTEGER' ... ");
        result = plc4c_subscription_request_add_cyclic_item(
            subscription_request, "RANDOM/foo:INTEGER", 500);
        if (result != OK) {
          printf("FAILED\n");
          return -1;
        }
        printf("SUCCESS\n");

        // Execute the subscription-request.
        printf("Executing a subscription-request ... ");
        result = plc4c_subscription_request_execute(
            subscription_request, &subscription_request_execution);
        if (result != OK) {
          printf("FAILED\n");
          return -1;
        } else {
          state = SUBSCRIPTION_REQUEST_SENT;
        }
        break;
      }
        // Wait until the read-request execution is finished.
      case SUBSCRIPTION_REQUEST_SENT: {
        if (plc4c_subscription_request_execution_check_finished_successfully(
                subscription_request_execution)) {
          printf("SUCCESS\n");
          state = SUBSCRIPTION_RESPONSE_RECEIVED;
        } else if (
            plc4c_subscription_request_execution_check_finished_with_error(
                subscription_request_execution)) {
          printf("FAILED\n");
          return -1;
        }
        break;
      }
      case SUBSCRIPTION_RESPONSE_RECEIVED: {
        // Get the response for the given subscription-request.
        plc4c_subscription_response *subscription_response =
            plc4c_subscription_request_execution_get_response(
                subscription_request_execution);
        if (subscription_response == NULL) {
          printf("FAILED (No Response)\n");
          return -1;
        }

        // Save the subsciption handle.
        plc4c_list_element *cur_element =
            plc4c_utils_list_head(subscription_response->response_items);
        plc4c_response_subscription_item *subscription_item_item =
            cur_element->value;
        printf("Value %s (%s):", subscription_item_item->item->name,
               plc4c_response_code_to_message(
                   subscription_item_item->response_code));
        // Check if the response was ok.
        if (subscription_item_item->response_code == OK) {
          // Save the subscription handles...
          subscription_handle = subscription_item_item->subscription_handle;
        } else {
          printf("FAILED (No Response)\n");
          return -1;
        }

        // Clean up the subscription..
        plc4c_subscription_response_destroy(subscription_response);
        plc4c_subscription_request_execution_destroy(
            subscription_request_execution);
        plc4c_subscription_request_destroy(subscription_request);

        state = READING_EVENTS;
        break;
      }
        // Wait for 10 incoming events.
      case READING_EVENTS: {
        // Check if an event is available ...
        printf("Checking Events ... ");
        if (plc4c_subscription_check_data_available(subscription_handle)) {
          printf("New events available.\n");

          printf("Getting Events ... ");
          plc4c_list *events = NULL;
          result = plc4c_subscription_get_subscription_events(
              subscription_handle, &events);
          if (result != OK) {
            printf("FAILED\n");
            return -1;
          }
          printf("SUCCESS\n");

          plc4c_list_element *cur_element = plc4c_utils_list_head(events);
          while (cur_element != NULL) {
            plc4c_response_subscription_item *subscription_event_item =
                cur_element->value;

            // TODO: Do something with the event ...
            printf("Got Event %s", subscription_event_item->item->name);

            // Increment the number of processed events.
            num_events++;

            cur_element = cur_element->next;
          }

          // If at least 10 events have been processed, disconnect.
          if (num_events > 10) {
            // Disconnect.
            printf("Disconnecting ... ");
            result = plc4c_connection_disconnect(connection);
            if (result != OK) {
              printf("FAILED");
              return -1;
            }
            state = DISCONNECTING;
          }
        } else {
          printf("No events.\n");
        }
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
