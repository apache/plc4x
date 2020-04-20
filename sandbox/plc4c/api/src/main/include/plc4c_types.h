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
#ifdef __cplusplus
extern "C" {
#endif

/**
 *
 * PLC4C error codes
*/
typedef enum error_code {
    OK,
    UNKNOWN_ERROR,
    NO_MEMORY,
    INVALID_CONNECTION_STRING,
    NOT_REACHABLE,
    PERMISSION_DENIED,
    INTERNAL_ERROR
} error_code;

/**
 * Helper that translates from an error_code enum value to something a human can work with.
 *
 * @param err error code.
 * @return A human readable error description.
 */
char *plc4c_error_code_to_error_message(error_code err);

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

#ifdef __cplusplus
}
#endif
#endif //PLC4C_TYPES_H_