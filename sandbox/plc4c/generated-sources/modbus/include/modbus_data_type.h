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

#ifndef PLC4C_MODBUS_READ_WRITE_MODBUS_DATA_TYPE_H_
#define PLC4C_MODBUS_READ_WRITE_MODBUS_DATA_TYPE_H_

#include <stdbool.h>
#include <stdint.h>

#ifdef __cplusplus
extern "C" {
#endif

enum plc4c_modbus_read_write_modbus_data_type {
  plc4c_modbus_read_write_modbus_data_type_NULL = 00,
  plc4c_modbus_read_write_modbus_data_type_BOOL = 01,
  plc4c_modbus_read_write_modbus_data_type_BYTE = 10,
  plc4c_modbus_read_write_modbus_data_type_WORD = 11,
  plc4c_modbus_read_write_modbus_data_type_DWORD = 12,
  plc4c_modbus_read_write_modbus_data_type_LWORD = 13,
  plc4c_modbus_read_write_modbus_data_type_SINT = 20,
  plc4c_modbus_read_write_modbus_data_type_INT = 21,
  plc4c_modbus_read_write_modbus_data_type_DINT = 22,
  plc4c_modbus_read_write_modbus_data_type_LINT = 23,
  plc4c_modbus_read_write_modbus_data_type_USINT = 24,
  plc4c_modbus_read_write_modbus_data_type_UINT = 25,
  plc4c_modbus_read_write_modbus_data_type_UDINT = 26,
  plc4c_modbus_read_write_modbus_data_type_ULINT = 27,
  plc4c_modbus_read_write_modbus_data_type_REAL = 30,
  plc4c_modbus_read_write_modbus_data_type_LREAL = 31,
  plc4c_modbus_read_write_modbus_data_type_TIME = 40,
  plc4c_modbus_read_write_modbus_data_type_LTIME = 41,
  plc4c_modbus_read_write_modbus_data_type_DATE = 50,
  plc4c_modbus_read_write_modbus_data_type_LDATE = 51,
  plc4c_modbus_read_write_modbus_data_type_TIME_OF_DAY = 60,
  plc4c_modbus_read_write_modbus_data_type_LTIME_OF_DAY = 61,
  plc4c_modbus_read_write_modbus_data_type_DATE_AND_TIME = 70,
  plc4c_modbus_read_write_modbus_data_type_LDATE_AND_TIME = 71,
  plc4c_modbus_read_write_modbus_data_type_CHAR = 80,
  plc4c_modbus_read_write_modbus_data_type_WCHAR = 81,
  plc4c_modbus_read_write_modbus_data_type_STRING = 82,
  plc4c_modbus_read_write_modbus_data_type_WSTRING = 83
};
typedef enum plc4c_modbus_read_write_modbus_data_type plc4c_modbus_read_write_modbus_data_type;

// Get an empty NULL-struct
plc4c_modbus_read_write_modbus_data_type plc4c_modbus_read_write_modbus_data_type_null();

plc4c_modbus_read_write_modbus_data_type plc4c_modbus_read_write_modbus_data_type_value_of(char* value_string);

int plc4c_modbus_read_write_modbus_data_type_num_values();

plc4c_modbus_read_write_modbus_data_type plc4c_modbus_read_write_modbus_data_type_value_for_index(int index);

uint8_t plc4c_modbus_read_write_modbus_data_type_get_data_type_size(plc4c_modbus_read_write_modbus_data_type value);

#ifdef __cplusplus
}
#endif

#endif  // PLC4C_MODBUS_READ_WRITE_MODBUS_DATA_TYPE_H_
