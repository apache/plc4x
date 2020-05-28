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
#ifndef PLC4C_S7_READ_WRITE_S7_PARAMETER_H_
#define PLC4C_S7_READ_WRITE_S7_PARAMETER_H_
#ifdef __cplusplus
extern "C" {
#endif

#include <stdbool.h>
#include <stdint.h>
#include <plc4c/utils/list.h>

struct plc4c_s7_read_write_s7_parameter {
  plc4c_s7_read_write_s7_parameter_type _type;
};
typedef struct plc4c_s7_read_write_s7_parameter plc4c_s7_read_write_s7_parameter;

// Enum assigning each sub-type an individual id.
enum plc4c_s7_read_write_s7_parameter_type {
  plc4c_s7_read_write_s7_parameter_type_s7_read_write_s7_parameter_setup_communication = 0;
  plc4c_s7_read_write_s7_parameter_type_s7_read_write_s7_parameter_read_var_request = 1;
  plc4c_s7_read_write_s7_parameter_type_s7_read_write_s7_parameter_read_var_response = 2;
  plc4c_s7_read_write_s7_parameter_type_s7_read_write_s7_parameter_write_var_request = 3;
  plc4c_s7_read_write_s7_parameter_type_s7_read_write_s7_parameter_write_var_response = 4;
  plc4c_s7_read_write_s7_parameter_type_s7_read_write_s7_parameter_user_data = 5;
}
typedef enum plc4c_s7_read_write_s7_parameter_types plc4c_s7_read_write_s7_parameter_types;

#ifdef __cplusplus
}
#endif
#endif  // PLC4C_S7_READ_WRITE_S7_PARAMETER_H_
