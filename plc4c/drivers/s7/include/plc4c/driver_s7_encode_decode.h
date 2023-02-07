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
#ifndef PLC4C_DRIVER_S7_ENCODE_DECODE_H_
#define PLC4C_DRIVER_S7_ENCODE_DECODE_H_
#ifdef __cplusplus
extern "C" {
#endif

#include <stdlib.h>

#include "plc4c/driver_s7.h"
#include "s7_var_request_parameter_item.h"

struct plc4c_s7_read_write_s7_var_request_parameter_item_field {
  plc4c_s7_read_write_s7_var_request_parameter_item* parameter_item;
  char* s7_address_any_encoding_of_string;
};
typedef struct plc4c_s7_read_write_s7_var_request_parameter_item_field plc4c_s7_read_write_s7_var_request_parameter_item_field;

uint16_t plc4c_driver_s7_encode_tsap_id(
    plc4c_driver_s7_device_group device_group, uint8_t rack, uint8_t slot);

uint16_t plc4c_driver_s7_get_nearest_matching_tpdu_size(uint16_t pdu_size);

plc4c_driver_s7_controller_type decode_controller_type(char* article_number);

plc4c_return_code decode_byte(const char* from_ptr, const char* to_ptr, uint8_t* value);

plc4c_return_code plc4c_driver_s7_encode_address(
    char* address, void** item);

#ifdef __cplusplus
}
#endif
#endif  // PLC4C_DRIVER_S7_ENCODE_DECODE_H_