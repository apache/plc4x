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

bool loop = true;
plc4c_system *system = NULL;
plc4c_connection *connection = NULL;

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

void onDisconnectSuccess(plc4c_promise *promise) {
    // Terminate the execution loop.
    loop = false;
}

void onReadSuccess(plc4c_promise *promise) {
    // TODO: Do something with the result.

    plc4c_promise* disconnect_promise = plc4c_connection_disconnect(connection);
    plc4c_promise_set_success_callback(disconnect_promise, &onDisconnectSuccess);
}

void onLocalConnectionSuccess(plc4c_promise* promise) {
    if(plc4c_connection_supports_reading(connection)) {
        char* addresses[] = {"RANDOM/foo:INTEGER"};
        plc4c_read_request* read_request = plc4c_connection_create_read_request(connection, 1, addresses);
        plc4c_promise* read_promise = plc4c_read_request_execute(connection, read_request);
        // As the read_request is actually executed the next time the plc4c_system_loop
        // is executed, we can now register some callbacks.
        plc4c_promise_set_success_callback(read_promise, &onReadSuccess);
    }
}

void onLocalConnectionFailure(plc4c_promise* promise) {
    // TODO: Do something with the error.
}

int main() {
    // Create a new uninitialized plc4c_system
    return_code result = plc4c_system_create(&system);
    if (result != OK) {
        return -1;
    }

    // Manually register the "simulated" driver with the system.
    plc4c_driver *simulatedDriver = plc4c_driver_simulated_create();
    result = plc4c_system_add_driver(simulatedDriver);
    if (result != OK) {
        return -1;
    }

    // Initialize the plc4c_system (loading of drivers, setting up other stuff, ...)
    result = plc4c_system_init(system);
    if (result != OK) {
        return -1;
    }

    // Register the global callbacks.
    plc4c_system_set_on_connection(system, &onGlobalConnect);
    plc4c_system_set_on_disconnection(system, &onGlobalDisconnect);

    // Establish connections to remote devices
    // you may or may not care about the connection handle
    plc4c_promise* connect_promise = plc4c_system_connect(system, "s7://192.168.42.20", &connection);
    // Register some callbacks to be called as soon as the connection is established or fails.
    plc4c_promise_set_success_callback(connect_promise, &onLocalConnectionSuccess);
    plc4c_promise_set_failure_callback(connect_promise, &onLocalConnectionFailure);
    if (plc4c_promise_completed_unsuccessfully(connect_promise)) {
        return -1;
    }

    // Central program loop ...
    while (loop) {
        if (plc4c_system_loop(system) != OK) {
            break;
        }
    }

    // Make sure everything is cleaned up correctly.
    plc4c_system_shutdown(system);

    // Finally destroy the plc4c_system, freeing up all memory allocated by plc4c.
    plc4c_system_destroy(system);

    return 0;
}
