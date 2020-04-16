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

bool loop = true;

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
 */
void plc4c_init(void (*errorFuncPtr) (error_code)) {
}

/**
 * Function to clean up the PLC4C system (Free any still used resources, terminate live connections, ...)
 */
void plc4c_shutdown(void (*errorFuncPtr) (error_code)) {
}

/**
 * Function to initialize a connection to a PLC by parsing the given connection string
 * and returning a corresponding plc4c_connection data structure.
 *
 * @param connectionString PLC4X connection string
 * @param plc4c_connection success handler function that gets called as soon as the connection is established.
 * @param errorFuncPtr error handler function that takes an error_code as argument.
 * @return plc4c_connection data structure.
 */
void plc4c_connect(char* connectionString,
        void (*successFuncPtr) (plc4c_connection),
        void (*errorFuncPtr) (plc4c_connection, error_code)) {
}

/**
 * Function to terminate a connection to a PLC.
 *
 * @param connection pointer to a plc4c_connection data structure.
 * @param plc4c_connection success handler function that gets called as soon as the connection is terminated.
 * @param errorFuncPtr error handler function that gets called if anything goes wrong.
 */
void plc4c_disconnect(plc4c_connection* connection,
        void (*successFuncPtr) (plc4c_connection),
        void (*errorFuncPtr) (plc4c_connection, error_code)) {
}

/**
 * Function to give any drivers the chance to do their work.
 * In single-threaded environments we can't operate with event
 * handler loops as they would block the rest of the application.
 *
 * @param errorFuncPtr error handler function that gets called if anything goes wrong.
 */
void plc4c_loop(void (*errorFuncPtr) (error_code)) {
}

/**
 * Default error handler, that just translates a normal error_code to a message.
 *
 * @param errorCode error code.
 */
void onError(error_code errorCode) {
    printf(toErrorMessage(errorCode));
}

/**
 * Error handler that handles errors related to a specific connection.
 *
 * @param connection connection the error is related to.
 * @param errorCode error code.
 */
void onConnectionError(plc4c_connection connection, error_code errorCode) {
    printf(toErrorMessage(errorCode));
}

/**
 * Callback called as soon as the connection is established.
 *
 * @param connection connection object.
 */
void onConnect(plc4c_connection connection) {

}

/**
 * Callback called as soon as the connection is terminated.
 *
 * @param connection connection object.
 */
void onDisconnect(plc4c_connection connection) {

}

/**
 * Main application.
 *
 * @return return code (Usually 0)
 */
int main() {
    // Initialize the PLC4C system.
    plc4c_init(&onError);

    // Establish a connection to a remote device
    plc4c_connect("s7://192.168.42.20", &onConnect, &onConnectionError);

    // Central program loop ...
    while(loop) {
        plc4c_loop(&onError);
    }

    // Make sure everything is cleaned up correctly.
    plc4c_shutdown(&onError);

    return 0;
}