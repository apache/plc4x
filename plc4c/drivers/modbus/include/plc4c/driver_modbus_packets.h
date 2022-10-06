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
#ifndef PLC4C_DRIVER_MODBUS_PACKETS_H_
#define PLC4C_DRIVER_MODBUS_PACKETS_H_
#ifdef __cplusplus
extern "C" {
#endif

#include <plc4c/types.h>

#include "driver_modbus.h"
#include "modbus_adu.h"

plc4c_return_code plc4c_driver_modbus_send_packet(
    plc4c_connection* connection,
    plc4c_modbus_read_write_modbus_adu* packet);
plc4c_return_code plc4c_driver_modbus_receive_packet(
    plc4c_connection* connection,
    plc4c_modbus_read_write_modbus_adu** packet);

plc4c_return_code plc4c_driver_modbus_create_modbus_read_request(
    plc4c_driver_modbus_config* modbus_config,
    plc4c_item* read_request_item,
    plc4c_modbus_read_write_modbus_adu** modbus_read_request_packet);
plc4c_return_code plc4c_driver_modbus_create_modbus_write_request(
    plc4c_driver_modbus_item* write_request_item,
    plc4c_modbus_read_write_modbus_adu** modbus_read_request_packet);

#ifdef __cplusplus
}
#endif
#endif  // PLC4C_DRIVER_MODBUS_PACKETS_H_