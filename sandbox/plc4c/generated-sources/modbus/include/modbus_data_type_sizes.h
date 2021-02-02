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

#ifndef PLC4C_MODBUS_READ_WRITE_MODBUS_DATA_TYPE_SIZES_H_
#define PLC4C_MODBUS_READ_WRITE_MODBUS_DATA_TYPE_SIZES_H_

#include <stdbool.h>
#include <stdint.h>

#ifdef __cplusplus
extern "C" {
#endif

enum plc4c_modbus_read_write_modbus_data_type_sizes {
  plc4c_modbus_read_write_modbus_data_type_sizes_BOOL = 0,
  plc4c_modbus_read_write_modbus_data_type_sizes_BYTE = 1,
  plc4c_modbus_read_write_modbus_data_type_sizes_WORD = 2,
  plc4c_modbus_read_write_modbus_data_type_sizes_DWORD = 3,
  plc4c_modbus_read_write_modbus_data_type_sizes_LWORD = 4,
  plc4c_modbus_read_write_modbus_data_type_sizes_SINT = 5,
  plc4c_modbus_read_write_modbus_data_type_sizes_INT = 6,
  plc4c_modbus_read_write_modbus_data_type_sizes_DINT = 7,
  plc4c_modbus_read_write_modbus_data_type_sizes_LINT = 8,
  plc4c_modbus_read_write_modbus_data_type_sizes_USINT = 9,
  plc4c_modbus_read_write_modbus_data_type_sizes_UINT = 10,
  plc4c_modbus_read_write_modbus_data_type_sizes_UDINT = 11,
  plc4c_modbus_read_write_modbus_data_type_sizes_ULINT = 12,
  plc4c_modbus_read_write_modbus_data_type_sizes_REAL = 13,
  plc4c_modbus_read_write_modbus_data_type_sizes_LREAL = 14,
  plc4c_modbus_read_write_modbus_data_type_sizes_TIME = 15,
  plc4c_modbus_read_write_modbus_data_type_sizes_LTIME = 16,
  plc4c_modbus_read_write_modbus_data_type_sizes_DATE = 17,
  plc4c_modbus_read_write_modbus_data_type_sizes_LDATE = 18,
  plc4c_modbus_read_write_modbus_data_type_sizes_TIME_OF_DAY = 19,
  plc4c_modbus_read_write_modbus_data_type_sizes_LTIME_OF_DAY = 20,
  plc4c_modbus_read_write_modbus_data_type_sizes_DATE_AND_TIME = 21,
  plc4c_modbus_read_write_modbus_data_type_sizes_LDATE_AND_TIME = 22,
  plc4c_modbus_read_write_modbus_data_type_sizes_CHAR = 23,
  plc4c_modbus_read_write_modbus_data_type_sizes_WCHAR = 24,
  plc4c_modbus_read_write_modbus_data_type_sizes_STRING = 25,
  plc4c_modbus_read_write_modbus_data_type_sizes_WSTRING = 26
};
typedef enum plc4c_modbus_read_write_modbus_data_type_sizes plc4c_modbus_read_write_modbus_data_type_sizes;

// Get an empty NULL-struct
plc4c_modbus_read_write_modbus_data_type_sizes plc4c_modbus_read_write_modbus_data_type_sizes_null();

plc4c_modbus_read_write_modbus_data_type_sizes plc4c_modbus_read_write_modbus_data_type_sizes_value_of(char* value_string);

int plc4c_modbus_read_write_modbus_data_type_sizes_num_values();

plc4c_modbus_read_write_modbus_data_type_sizes plc4c_modbus_read_write_modbus_data_type_sizes_value_for_index(int index);

uint8_t plc4c_modbus_read_write_modbus_data_type_sizes_get_data_type_size(plc4c_modbus_read_write_modbus_data_type_sizes value);
plc4c_modbus_read_write_modbus_data_type_sizes plc4c_modbus_read_write_modbus_data_type_sizes_get_first_enum_for_field_data_type_size(uint8_t value);

#ifdef __cplusplus
}
#endif

#endif  // PLC4C_MODBUS_READ_WRITE_MODBUS_DATA_TYPE_SIZES_H_
