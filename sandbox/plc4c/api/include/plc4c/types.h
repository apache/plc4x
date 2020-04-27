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
#ifndef PLC4C_TYPES_H_
#define PLC4C_TYPES_H_

#include <stdbool.h>

#ifdef __cplusplus
extern "C" {
#endif

/**
 *
 * PLC4C error codes
*/
typedef enum return_code {
    UNFINISHED,
    OK,
    UNKNOWN_ERROR,
    NO_MEMORY,
    INVALID_CONNECTION_STRING,
    NOT_REACHABLE,
    PERMISSION_DENIED,
    INTERNAL_ERROR
} return_code;

/**
 * Helper that translates from an return_code enum value to something a human can work with.
 *
 * @param err return code.
 * @return A human readable description.
 */
char *plc4c_return_code_to_message(return_code err);

/**
 * the plc4c system
 */
typedef struct plc4c_system_t plc4c_system;

/**
 * the plc4c_driver
 */
typedef struct plc4c_driver_t plc4c_driver;

/**
 * the plc4c_connection
 */
typedef struct plc4c_connection_t plc4c_connection;

/**
 * a plc4c read-request
 */
typedef struct plc4c_read_request_t plc4c_read_request;

/**
 * a plc4c write-request
 */
typedef struct plc4c_write_request_t plc4c_write_request;

/**
 * Return type for any form of async operation.
 */
typedef struct plc4c_promise_t plc4c_promise;

/**
 * Callback for any form of successful async operation.
 */
typedef void (*plc4c_success_callback)(plc4c_promise *promise);

/**
 * Callback for any form of failed async operation.
 */
typedef void (*plc4c_failure_callback)(plc4c_promise *promise);

/**
 * Function to register a success callback on a given plc4c_promise.
 * @param promise the promise to set the callback on
 * @param successCallback the callback
 */
void plc4c_promise_set_success_callback(plc4c_promise* promise, plc4c_success_callback successCallback);

/**
 * Function to register a failure callback on a given plc4c_promise.
 * @param promise the promise to set the callback on
 * @param successCallback the callback
 */
void plc4c_promise_set_failure_callback(plc4c_promise* promise, plc4c_failure_callback failureCallback);

/**
 * Check if a promise is completed
 * @param promise the promise
 * @return true if the promise is in a final state
 */
bool plc4c_promise_completed(plc4c_promise* promise);

/**
 * Check if a promise is completed successfully
 * @param promise the promise
 * @return true if the promise is the state OK
 */
bool plc4c_promise_completed_successfully(plc4c_promise* promise);

/**
 * Check if a promise is completed unsuccessfully
 * @param promise the promise
 * @return true if the promise is in a final state other than OK
 */
bool plc4c_promise_completed_unsuccessfully(plc4c_promise* promise);


#ifdef __cplusplus
}
#endif
#endif //PLC4C_TYPES_H_