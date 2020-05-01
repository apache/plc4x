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
#include <stdio.h>
#include <plc4c/plc4c.h>
#include <plc4c/driver_simulated.h>
#include <plc4c/transport_dummy.h>


int numOpenConnections = 0;

/**
 * Here we could implement something that keeps track of all open connections.
 * For example on embedded devices using the W5100 SPI Network device, this can
 * only handle 4 simultaneous connections.
 *
 * @param connection the connection that was just established
 */
void onGlobalConnect(plc4c_connection *cur_connection) {
    printf("Connected to %s", plc4c_connection_get_connection_string(cur_connection));
    numOpenConnections++;
}

void onGlobalDisconnect(plc4c_connection *cur_connection) {
    printf("Disconnected from %s", plc4c_connection_get_connection_string(cur_connection));
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
typedef enum plc4c_connection_state_t plc4c_connection_state ;

#pragma clang diagnostic push
#pragma ide diagnostic ignored "hicpp-multiway-paths-covered"
int main() {
    bool loop = true;
    plc4c_system *system = NULL;
    plc4c_connection *connection = NULL;
    plc4c_read_request *read_request = NULL;
    plc4c_read_request_execution *read_request_execution = NULL;

    // Create a new uninitialized plc4c_system
    return_code result = plc4c_system_create(&system);
    if (result != OK) {
        return -1;
    }

    // Manually register the "simulated" driver with the system.
    plc4c_driver *simulated_driver = plc4c_driver_simulated_create();
    result = plc4c_system_add_driver(system, simulated_driver);
    if (result != OK) {
        return -1;
    }

    plc4c_transport *dummy_transport = plc4c_transport_dummy_create();
    result = plc4c_system_add_transport(system, dummy_transport);
    if (result != OK) {
        return -1;
    }

    // Initialize the plc4c_system (loading of drivers, setting up other stuff, ...)
    result = plc4c_system_init(system);
    if (result != OK) {
        return -1;
    }

    // Register the global callbacks.
    plc4c_system_set_on_connect_success_callback(system, &onGlobalConnect);
    plc4c_system_set_on_disconnect_success_callback(system, &onGlobalDisconnect);

    // Establish connections to remote devices
    // you may or may not care about the connection handle
    result = plc4c_system_connect(system, "simulated://foo", &connection);
    if (result != OK) {
        return -1;
    }

    // Central program loop ...
    plc4c_connection_state state = CONNECTING;
    while (loop) {
        // Give plc4c a chance to do something.
        // This is where all I/O is being done.
        if (plc4c_system_loop(system) != OK) {
            break;
        }

        // Depending on the current state, implement some logic.
        switch (state) {
            case CONNECTING: {
                // Check if the connection is established:
                if (plc4c_connection_is_connected(connection)) {
                    state = CONNECTED;
                } else if (plc4c_connection_has_error(connection)) {
                    return -1;
                }
                break;
            }
            case CONNECTED: {
                // Create a new read-request.
                char *addresses[] = {"RANDOM/foo:INTEGER"};
                result = plc4c_connection_create_read_request(connection, 1, addresses, &read_request);
                if(result != OK) {
                    return -1;
                }

                // Execute the read-request.
                result = plc4c_read_request_execute(read_request, &read_request_execution);
                if(result != OK) {
                    return -1;
                } else {
                    state = READ_REQUEST_SENT;
                }
                break;
            }
            // Wait until the read-request execution is finished.
            case READ_REQUEST_SENT: {
                if(plc4c_read_request_finished_successfully(read_request_execution)) {
                    state = READ_RESPONSE_RECEIVED;
                } else if(plc4c_read_request_has_error(read_request_execution)) {
                    return -1;
                }
                break;
            }
            case READ_RESPONSE_RECEIVED: {
                // Get the response for the given read-request.
                plc4c_read_response *response = plc4c_read_request_get_response(read_request_execution);

                // TODO: Do something sensible ...

                // Clean up.
                plc4c_read_request_execution_destroy(read_request_execution);
                plc4c_read_request_destroy(read_request);

                // Disconnect.
                result = plc4c_connection_disconnect(connection);
                if(result != OK) {
                    return -1;
                }
                state = DISCONNECTING;
                break;
            }
            // Wait until the connection is disconnected
            case DISCONNECTING: {
                if(!plc4c_connection_is_connected(connection)) {
                    plc4c_connection_destroy(connection);
                    state = DISCONNECTED;
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
