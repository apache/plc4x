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
#ifndef PLC4C_READ_H_
#define PLC4C_READ_H_
#ifdef __cplusplus
extern "C" {
#endif

#include "plc4c/types.h"

/**
 * Destroys a given read-request.
 * @param read_request the read-request.
 */
void plc4c_read_request_destroy(plc4c_read_request *read_request);

/**
 * Function for adding a new item to a given request.
 * @param read_request the read-request.
 * @param address address string
 * @return return code
 */
plc4c_return_code plc4c_read_request_add_item(plc4c_read_request *read_request,
                                              char *field_name,
                                              char *address);

/**
 * Returns the plc4c_connection for a give read request
 * @param read_request plc4c_read_request
 * @return plc4c_connection
 */
plc4c_connection *plc4c_read_request_get_connection(
    plc4c_read_request *read_request);

/**
 * Sets the plc4c_connection for a given read request
 * @param read_request plc4c_read_request
 * @param connection plc4c_connection
 */
void plc4c_read_request_set_connection(plc4c_read_request *read_request,
                                       plc4c_connection *connection);

/**
 * Actually executes the read-request.
 * @param connection connection this read-request will be executed on.
 * @param read_request the read-request object.
 * @param read_request_execution pointer to a data-structure handling one
 * execution of the read-request.
 * @return plc4c_return_code
 */
plc4c_return_code plc4c_read_request_execute(
    plc4c_read_request *read_request,
    plc4c_read_request_execution **read_request_execution);

/**
 * Destroys a given read-request execution.
 *
 * @param read_request_execution the read-request execution.
 */
void plc4c_read_request_execution_destroy(
    plc4c_read_request_execution *read_request_execution);

/**
 * Check if the read-request is completed successfully.
 *
 * @param read_request_execution the read-request execution.
 * @return true if the read-request is completed successfully.
 */
bool plc4c_read_request_execution_check_finished_successfully(
    plc4c_read_request_execution *read_request_execution);

/**
 * Check if the read-request is completed unsuccessfully.
 *
 * @param read_request_execution the read-request execution.
 * @return true if the read-request is completed with an error.
 */
bool plc4c_read_request_execution_check_finished_with_error(
    plc4c_read_request_execution *read_request_execution);

/**
 * Retrieve the read-response from a given read-request execution.
 *
 * @param read_request_execution the read-request execution.
 * @return the read-response.
 */
plc4c_read_response *plc4c_read_request_execution_get_response(
    plc4c_read_request_execution *read_request_execution);

/**
 * Destroys a given read_response
 * @param read_response the read_response
 */
void plc4c_read_response_destroy(plc4c_read_response *read_response);

#ifdef __cplusplus
}
#endif
#endif  // PLC4C_READ_H_