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
#ifndef PLC4C_WRITE_H_
#define PLC4C_WRITE_H_
#ifdef __cplusplus
extern "C" {
#endif

#include "plc4c/types.h"

/**
 * Destroys a given write-request.
 * @param write_request the write-request.
 */
void plc4c_write_request_destroy(plc4c_write_request *write_request);

/**
 * Function for adding a new item to a given request.
 * @param write_request the write-request.
 * @param address address string
 * @param value value of the given resource
 * @return return code
 */
plc4c_return_code plc4c_write_request_add_item(
    plc4c_write_request *write_request, char *address, plc4c_data *value);

/**
 * Returns the plc4c_connection for a give write request
 * @param write_request plc4c_write_request
 * @return plc4c_connection
 */
plc4c_connection *plc4c_write_request_get_connection(
    plc4c_write_request *write_request);

/**
 * Sets the plc4c_connection for a given write request
 * @param write_request plc4c_write_request
 * @param connection plc4c_connection
 */
void plc4c_write_request_set_connection(plc4c_write_request *write_request,
                                        plc4c_connection *connection);

/**
 * Actually executes the write-request.
 *
 * @param write_request the write-request object.
 * @param write_request_execution pointer to a data-structure handling one
 * execution of the write-request.
 * @return plc4c_return_code
 */
plc4c_return_code plc4c_write_request_execute(
    plc4c_write_request *write_request,
    plc4c_write_request_execution **write_request_execution);

/**
 * Destroys a given write-request execution.
 *
 * @param write_request_execution the write-request execution.
 */
void plc4c_write_request_execution_destroy(
    plc4c_write_request_execution *write_request_execution);

/**
 * Check if the write-request is completed successfully.
 *
 * @param write_request_execution the write-request execution.
 * @return true if the write-request is completed successfully.
 */
bool plc4c_write_request_check_finished_successfully(
    plc4c_write_request_execution *write_request_execution);

/**
 * Check if the write-request is completed unsuccessfully.
 *
 * @param write_request_execution the write-request execution.
 * @return true if the write-request is completed with an error.
 */
bool plc4c_write_request_execution_check_completed_with_error(
    plc4c_write_request_execution *write_request_execution);

/**
 * Retrieve the write-response from a given write-request execution.
 *
 * @param write_request_execution the write-request execution.
 * @return the write-response.
 */
plc4c_write_response *plc4c_write_request_execution_get_response(
    plc4c_write_request_execution *write_request_execution);

/**
 * Destroys a given write_response
 * @param write_response the write_response
 */
void plc4c_write_response_destroy(plc4c_write_response *write_response);

plc4c_return_code plc4x_write_execution_status(plc4c_write_request_execution *execution);

#ifdef __cplusplus
}
#endif
#endif  // PLC4C_WRITE_H_
