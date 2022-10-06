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
#ifndef PLC4C_DRIVER_PLC4X_H_
#define PLC4C_DRIVER_PLC4X_H_
#ifdef __cplusplus
extern "C" {
#endif

#include <plc4c/types.h>
#include <stdint.h>
#include <time.h>

#include "../../../../spi/include/plc4c/spi/read_buffer.h"

struct plc4c_driver_plc4x_config {
  char* remote_connection_string;
  uint32_t request_timeout;

  uint16_t connection_id;
  uint16_t request_id;
};
typedef struct plc4c_driver_plc4x_config plc4c_driver_plc4x_config;

plc4c_driver *plc4c_driver_plc4x_create();

#ifdef __cplusplus
}
#endif
#endif  // PLC4C_DRIVER_PLC4X_H_