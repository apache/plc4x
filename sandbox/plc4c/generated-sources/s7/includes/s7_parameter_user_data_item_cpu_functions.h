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
#ifndef PLC4C_S7_READ_WRITE_S7_PARAMETER_USER_DATA_ITEM_CPU_FUNCTIONS_H_
#define PLC4C_S7_READ_WRITE_S7_PARAMETER_USER_DATA_ITEM_CPU_FUNCTIONS_H_
#ifdef __cplusplus
extern "C" {
#endif

#include <stdbool.h>
#include <stdint.h>
#include <plc4c/utils/list.h>

struct plc4c_s7_read_write_s7_parameter_user_data_item_cpu_functions {
  plc4c_s7_read_write_s7_parameter_user_data_item_type _type;
  uint8_t method;
  unsigned int cpu_function_type : 4;
  unsigned int cpu_function_group : 4;
  uint8_t cpu_subfunction;
  uint8_t sequence_number;
  uint8_t data_unit_reference_number;
  uint8_t last_data_unit;
  uint16_t error_code;
};
typedef struct plc4c_s7_read_write_s7_parameter_user_data_item_cpu_functions plc4c_s7_read_write_s7_parameter_user_data_item_cpu_functions;

#ifdef __cplusplus
}
#endif
#endif  // PLC4C_S7_READ_WRITE_S7_PARAMETER_USER_DATA_ITEM_CPU_FUNCTIONS_H_
