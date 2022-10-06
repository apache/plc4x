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
#ifndef PLC4C_DRIVER_PLC4X_PACKETS_H_
#define PLC4C_DRIVER_PLC4X_PACKETS_H_
#ifdef __cplusplus
extern "C" {
#endif

#include <plc4c/types.h>

#include "plc4x_message.h"
#include "plc4c/driver_plc4x.h"

plc4c_return_code plc4c_driver_plc4x_send_packet(
    plc4c_connection* connection, plc4c_plc4x_read_write_plc4x_message* packet);

plc4c_return_code plc4c_driver_plc4x_receive_packet(
    plc4c_connection* connection, plc4c_plc4x_read_write_plc4x_message** packet);
void plc4c_driver_plc4x_destroy_packet(
    plc4c_plc4x_read_write_plc4x_message* packet);

plc4c_return_code plc4c_driver_plc4x_create_connection_request(
    plc4c_driver_plc4x_config* configuration,
    plc4c_plc4x_read_write_plc4x_message** connect_request_packet);
void plc4c_driver_plc4x_destroy_connection_request(
    plc4c_plc4x_read_write_plc4x_message* packet);

plc4c_return_code plc4c_driver_plc4x_create_plc4x_read_request(
    plc4c_read_request* read_request,
    plc4c_plc4x_read_write_plc4x_message** plc4x_read_request_packet);
void plc4c_driver_plc4x_destroy_plc4x_read_request(
    plc4c_plc4x_read_write_plc4x_message* packet);

plc4c_return_code plc4c_driver_plc4x_create_plc4x_write_request(
    plc4c_write_request* write_request,
    plc4c_plc4x_read_write_plc4x_message** plc4x_write_request_packet);
void plc4c_driver_plc4x_destroy_plc4x_write_request(
    plc4c_plc4x_read_write_plc4x_message* packet);

#ifdef __cplusplus
}
#endif
#endif  // PLC4C_DRIVER_PLC4X_PACKETS_H_