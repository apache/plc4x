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

#ifndef PLC4C_DRIVER_PLC4X_STATIC_H
#define PLC4C_DRIVER_PLC4X_STATIC_H
#ifdef __cplusplus
extern "C" {
#endif

#include <plc4c/spi/read_buffer.h>
#include <plc4c/spi/write_buffer.h>
#include <time.h>
#include <stdint.h>

/*
 *
 *   Static functions
 *
 */

uint8_t plc4c_spi_evaluation_helper_str_len(char* str);

char* plc4c_plc4x_read_write_parse_string(plc4c_spi_read_buffer* io, char* encoding);

#ifdef __cplusplus
}
#endif
#endif  // PLC4C_DRIVER_PLC4X_STATIC_H
