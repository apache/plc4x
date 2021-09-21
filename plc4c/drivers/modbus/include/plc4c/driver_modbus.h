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
#ifndef PLC4C_DRIVER_MODBUS_H_
#define PLC4C_DRIVER_MODBUS_H_
#ifdef __cplusplus
extern "C" {
#endif

#include <plc4c/types.h>
#include <stdint.h>

#include "../../../../generated-sources/modbus/include/modbus_data_type.h"

struct plc4c_driver_modbus_config {
  uint16_t request_timeout;
  uint8_t unit_identifier;
  uint8_t communication_id_counter;
};
typedef struct plc4c_driver_modbus_config plc4c_driver_modbus_config;

plc4c_driver *plc4c_driver_modbus_create();

enum plc4c_driver_modbus_address_type {
  PLC4C_DRIVER_MODBUS_ADDRESS_TYPE_COIL = 0,
  PLC4C_DRIVER_MODBUS_ADDRESS_TYPE_DISCRETE_INPUT = 1,
  PLC4C_DRIVER_MODBUS_ADDRESS_TYPE_INPUT_REGISTER = 3,
  PLC4C_DRIVER_MODBUS_ADDRESS_TYPE_HOLDING_REGISTER = 4,
  PLC4C_DRIVER_MODBUS_ADDRESS_TYPE_EXTENDED_REGISTER = 6
};
typedef enum plc4c_driver_modbus_address_type plc4c_driver_modbus_address_type;

struct plc4c_driver_modbus_item {
  plc4c_driver_modbus_address_type type;
  uint16_t address;
  plc4c_modbus_read_write_modbus_data_type datatype;
  uint16_t num_elements;
};
typedef struct plc4c_driver_modbus_item plc4c_driver_modbus_item;

#ifdef __cplusplus
}
#endif
#endif  // PLC4C_DRIVER_MODBUS_H_