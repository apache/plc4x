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
#ifndef PLC4C_DRIVER_S7_H_
#define PLC4C_DRIVER_S7_H_
#ifdef __cplusplus
extern "C" {
#endif

#include <plc4c/types.h>
#include <stdint.h>
#include <time.h>

#include "../../../../generated-sources/s7/include/cotp_tpdu_size.h"
#include "../../../../spi/include/plc4c/spi/read_buffer.h"

enum plc4c_driver_s7_controller_type {
  PLC4C_DRIVER_S7_CONTROLLER_TYPE_ANY = 0,
  PLC4C_DRIVER_S7_CONTROLLER_TYPE_S7_300 = 1,
  PLC4C_DRIVER_S7_CONTROLLER_TYPE_S7_400 = 2,
  PLC4C_DRIVER_S7_CONTROLLER_TYPE_S7_1200 = 3,
  PLC4C_DRIVER_S7_CONTROLLER_TYPE_S7_1500 = 4,
  PLC4C_DRIVER_S7_CONTROLLER_TYPE_LOGO = 5
};
typedef enum plc4c_driver_s7_controller_type plc4c_driver_s7_controller_type;

enum plc4c_driver_s7_device_group {
  PLC4C_DRIVER_S7_DEVICE_GROUP_PG_OR_PC = 1,
  PLC4C_DRIVER_S7_DEVICE_GROUP_OS = 2,
  PLC4C_DRIVER_S7_DEVICE_GROUP_OTHERS = 3
};
typedef enum plc4c_driver_s7_device_group plc4c_driver_s7_device_group;

struct plc4c_driver_s7_config {
  uint8_t local_rack;
  uint8_t local_slot;
  uint8_t remote_rack;
  uint8_t remote_slot;
  uint16_t calling_tsap_id;
  uint16_t called_tsap_id;
  plc4c_s7_read_write_cotp_tpdu_size cotp_tpdu_size;
  uint16_t pdu_size;
  uint8_t max_amq_caller;
  uint8_t max_amq_callee;
  plc4c_driver_s7_controller_type controller_type;

  uint16_t pdu_id;
};
typedef struct plc4c_driver_s7_config plc4c_driver_s7_config;

plc4c_driver *plc4c_driver_s7_create();

#ifdef __cplusplus
}
#endif
#endif  // PLC4C_DRIVER_S7_H_