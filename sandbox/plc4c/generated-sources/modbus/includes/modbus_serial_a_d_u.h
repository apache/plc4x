/*
  Licensed to the Apache Software Foundation (ASF) under one
  or more contributor license agreements.  See the NOTICE file
  distributed with this work for additional information
  regarding copyright ownership.  The ASF licenses this file
  to you under the Apache License, Version 2.0 (the
  "License"); you may not use this file except in compliance
  with the License.  You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing,
  software distributed under the License is distributed on an
  "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
  KIND, either express or implied.  See the License for the
  specific language governing permissions and limitations
  under the License.
*/
#ifndef PLC4C_MODBUS_READ_WRITE_MODBUS_SERIAL_A_D_U_H_
#define PLC4C_MODBUS_READ_WRITE_MODBUS_SERIAL_A_D_U_H_
#ifdef __cplusplus
extern "C" {
#endif

#include <stdbool.h>
#include <stdint.h>
#include <plc4c/utils/list.h>
#include "modbus_p_d_u.h"

struct plc4c_modbus_read_write_modbus_serial_a_d_u {
  uint16_t transaction_id;
  uint16_t length;
  uint8_t address;
  plc4c_modbus_read_write_modbus_p_d_u pdu;
};
typedef struct plc4c_modbus_read_write_modbus_serial_a_d_u plc4c_modbus_read_write_modbus_serial_a_d_u;

#ifdef __cplusplus
}
#endif
#endif  // PLC4C_MODBUS_READ_WRITE_MODBUS_SERIAL_A_D_U_H_
