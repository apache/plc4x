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
#include<stdio.h>
#include<stdbool.h>

/**
 Some error types that can typically occur.
 */
enum error_code {
    OK,
    NOT_REACHABLE,
    PERMISSION_DENIED,
    INTERNAL_ERROR
};
typedef enum error_code error_code;

/**
 * Helper that translates from an error_code enum value to something a human can work with.
 *
 * @param err error code.
 * @return A human readable error description.
 */
char* toErrorMessage(error_code err) {
    if(err) {
        switch (err) {
            case OK:
                return "OK";
            case NOT_REACHABLE:
                return "requested destination not reachable/not available";
            case PERMISSION_DENIED:
                return "permission denied";
            case INTERNAL_ERROR:
                return "An internal server error occured";
            default:
                return "WTF?!?!";
        }
    }
    return "WTF?!?!";
}

/**
 * Data structure containing all data related to an individual connection.
 */
struct plc4c_connection {
    bool connected;
    error_code error;
};
typedef struct plc4c_connection plc4c_connection;

/**
 * Function to initialize the PLC4C system (Initialize the driver manager and the list of enabled drivers)
 *
 * @return an error code
 */
error_code plc4c_init() {
    return OK;
}

/**
 * Function to clean up the PLC4C system (Free any still used resources, terminate live connections, ...)
 *
 * @return an error code
 */
error_code plc4c_shutdown() {
    return OK;
}

/**
 * Function to initialize a connection to a PLC by parsing the given connection string
 * and filling in the plc4c_connection data structure.
 *
 * @param connectionString PLC4X connection string
 * @param connection (return) pointer to a plc4c_connection data structure.
 * @return an error code
 */
error_code plc4c_connect(char* connectionString, plc4c_connection** connection) {
    return OK;
}

/**
 * Function to terminate a connection to a PLC.
 *
 * @param connection pointer to a plc4c_connection data structure.
 * @return an error code
 */
error_code plc4c_disconnect(plc4c_connection* connection) {
    return OK;
}

/**
 * Function to give any drivers the chance to do their work.
 * In single-threaded environments we can't operate with event
 * handler loops as they would block the rest of the application.
 *
 * @return an error code (But only if there's something bad going
 * happening with PLC4C and not just something going wrong in a connection)
 */
error_code plc4c_loop() {
    return OK;
}


/**
 * Main application.
 *
 * @return return code (Usually 0)
 */
int main() {
    // Initialize the PLC4C system.
    error_code res = plc4c_init();
    if(res != OK) {
        printf("ERROR: An error initializing PLC4C driver manager");
    }

    // Establish a connection to a remote device
    plc4c_connection* connection = NULL;
    res = plc4c_connect("s7://192.168.42.20", &connection);
    // Check if all was ok
    if(res != OK) {
        printf("ERROR: An error occurred while connecting");
        goto cleanup;
    }

    bool loop = true;
    while(loop) {
        // Give the plc4c internals the chance to do their work.
        // NOTE: If for example one of two connections has a problem, this will not return
        // an error code as we would have no way of knowing whihc connection was having issues.
        res = plc4c_loop();
        if(res != OK) {
            printf("ERROR: An error occurred in the loop");
            goto cleanup;
        }

        // As long as we're not connected and not in an error state, abort this loop iteration.
        if(!connection->connected) {
            // If the connection is in an error state, abort the loop.
            if(connection->error) {
                printf("An connection error occurred");
                goto cleanup;
            }
            continue;
        }

        // TODO: Do stuff ...
        loop = false;
    }

    // If we're connected, gracefully shut down the connection.
    if(connection->connected) {
        res = plc4c_disconnect(connection);
        if(res != OK) {
            printf("ERROR: An error while disconnecting");
            goto cleanup;
        }
        // Give the connection some time to say goodbye.
        while (1) {
            res = plc4c_loop();
            if(res != OK) {
                printf("ERROR: An error occurred in the loop");
                goto cleanup;
            }
            if(!connection->connected) {
                printf("Disconnected");
                goto cleanup;
            }
        }
    }

    cleanup:
    // Make sure everything is cleaned up correctly.
    plc4c_shutdown();

    return 0;
}