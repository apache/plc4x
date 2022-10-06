/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
#ifndef PLC4C_DRIVER_PLC4X_SM_H_
#define PLC4C_DRIVER_PLC4X_SM_H_
#ifdef __cplusplus
extern "C" {
#endif

#include <plc4c/types.h>
#include <plc4c/spi/types_private.h>

plc4c_return_code plc4c_driver_plc4x_connect_function(
    plc4c_connection* connection, plc4c_system_task** task);

plc4c_return_code plc4c_driver_plc4x_disconnect_function(
    plc4c_connection* connection, plc4c_system_task** task);

plc4c_return_code plc4c_driver_plc4x_read_function(
    plc4c_read_request_execution* read_request_execution,
    plc4c_system_task** task);
void plc4c_driver_plc4x_free_read_request(plc4c_read_request* request);
void plc4c_driver_plc4x_free_read_response(plc4c_read_response* response);

plc4c_return_code plc4c_driver_plc4x_write_function(
    plc4c_write_request_execution* write_request_execution,
    plc4c_system_task** task);
void plc4c_driver_plc4x_free_write_request(plc4c_write_request* request);
void plc4c_driver_plc4x_free_write_response(plc4c_write_response* response);

#ifdef __cplusplus
}
#endif
#endif  // PLC4C_DRIVER_PLC4X_SM_H_