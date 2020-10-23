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

#include "modbus_data_type_sizes.h"
#include <string.h>


// Create an empty NULL-struct
static const plc4c_modbus_read_write_modbus_data_type_sizes plc4c_modbus_read_write_modbus_data_type_sizes_null_const;

plc4c_modbus_read_write_modbus_data_type_sizes plc4c_modbus_read_write_modbus_data_type_sizes_null() {
  return plc4c_modbus_read_write_modbus_data_type_sizes_null_const;
}

plc4c_modbus_read_write_modbus_data_type_sizes plc4c_modbus_read_write_modbus_data_type_sizes_value_of(char* value_string) {
    if(strcmp(value_string, "BOOL") == 0) {
        return plc4c_modbus_read_write_modbus_data_type_sizes_BOOL;
    }
    if(strcmp(value_string, "BYTE") == 0) {
        return plc4c_modbus_read_write_modbus_data_type_sizes_BYTE;
    }
    if(strcmp(value_string, "CHAR") == 0) {
        return plc4c_modbus_read_write_modbus_data_type_sizes_CHAR;
    }
    if(strcmp(value_string, "DATE") == 0) {
        return plc4c_modbus_read_write_modbus_data_type_sizes_DATE;
    }
    if(strcmp(value_string, "DATE_AND_TIME") == 0) {
        return plc4c_modbus_read_write_modbus_data_type_sizes_DATE_AND_TIME;
    }
    if(strcmp(value_string, "DINT") == 0) {
        return plc4c_modbus_read_write_modbus_data_type_sizes_DINT;
    }
    if(strcmp(value_string, "DWORD") == 0) {
        return plc4c_modbus_read_write_modbus_data_type_sizes_DWORD;
    }
    if(strcmp(value_string, "INT") == 0) {
        return plc4c_modbus_read_write_modbus_data_type_sizes_INT;
    }
    if(strcmp(value_string, "LDATE") == 0) {
        return plc4c_modbus_read_write_modbus_data_type_sizes_LDATE;
    }
    if(strcmp(value_string, "LDATE_AND_TIME") == 0) {
        return plc4c_modbus_read_write_modbus_data_type_sizes_LDATE_AND_TIME;
    }
    if(strcmp(value_string, "LINT") == 0) {
        return plc4c_modbus_read_write_modbus_data_type_sizes_LINT;
    }
    if(strcmp(value_string, "LREAL") == 0) {
        return plc4c_modbus_read_write_modbus_data_type_sizes_LREAL;
    }
    if(strcmp(value_string, "LTIME") == 0) {
        return plc4c_modbus_read_write_modbus_data_type_sizes_LTIME;
    }
    if(strcmp(value_string, "LTIME_OF_DAY") == 0) {
        return plc4c_modbus_read_write_modbus_data_type_sizes_LTIME_OF_DAY;
    }
    if(strcmp(value_string, "LWORD") == 0) {
        return plc4c_modbus_read_write_modbus_data_type_sizes_LWORD;
    }
    if(strcmp(value_string, "REAL") == 0) {
        return plc4c_modbus_read_write_modbus_data_type_sizes_REAL;
    }
    if(strcmp(value_string, "SINT") == 0) {
        return plc4c_modbus_read_write_modbus_data_type_sizes_SINT;
    }
    if(strcmp(value_string, "STRING") == 0) {
        return plc4c_modbus_read_write_modbus_data_type_sizes_STRING;
    }
    if(strcmp(value_string, "TIME") == 0) {
        return plc4c_modbus_read_write_modbus_data_type_sizes_TIME;
    }
    if(strcmp(value_string, "TIME_OF_DAY") == 0) {
        return plc4c_modbus_read_write_modbus_data_type_sizes_TIME_OF_DAY;
    }
    if(strcmp(value_string, "UDINT") == 0) {
        return plc4c_modbus_read_write_modbus_data_type_sizes_UDINT;
    }
    if(strcmp(value_string, "UINT") == 0) {
        return plc4c_modbus_read_write_modbus_data_type_sizes_UINT;
    }
    if(strcmp(value_string, "ULINT") == 0) {
        return plc4c_modbus_read_write_modbus_data_type_sizes_ULINT;
    }
    if(strcmp(value_string, "USINT") == 0) {
        return plc4c_modbus_read_write_modbus_data_type_sizes_USINT;
    }
    if(strcmp(value_string, "WCHAR") == 0) {
        return plc4c_modbus_read_write_modbus_data_type_sizes_WCHAR;
    }
    if(strcmp(value_string, "WORD") == 0) {
        return plc4c_modbus_read_write_modbus_data_type_sizes_WORD;
    }
    if(strcmp(value_string, "WSTRING") == 0) {
        return plc4c_modbus_read_write_modbus_data_type_sizes_WSTRING;
    }
    return -1;
}

int plc4c_modbus_read_write_modbus_data_type_sizes_num_values() {
  return 27;
}

plc4c_modbus_read_write_modbus_data_type_sizes plc4c_modbus_read_write_modbus_data_type_sizes_value_for_index(int index) {
    switch(index) {
      case 0: {
        return plc4c_modbus_read_write_modbus_data_type_sizes_BOOL;
      }
      case 1: {
        return plc4c_modbus_read_write_modbus_data_type_sizes_BYTE;
      }
      case 2: {
        return plc4c_modbus_read_write_modbus_data_type_sizes_CHAR;
      }
      case 3: {
        return plc4c_modbus_read_write_modbus_data_type_sizes_DATE;
      }
      case 4: {
        return plc4c_modbus_read_write_modbus_data_type_sizes_DATE_AND_TIME;
      }
      case 5: {
        return plc4c_modbus_read_write_modbus_data_type_sizes_DINT;
      }
      case 6: {
        return plc4c_modbus_read_write_modbus_data_type_sizes_DWORD;
      }
      case 7: {
        return plc4c_modbus_read_write_modbus_data_type_sizes_INT;
      }
      case 8: {
        return plc4c_modbus_read_write_modbus_data_type_sizes_LDATE;
      }
      case 9: {
        return plc4c_modbus_read_write_modbus_data_type_sizes_LDATE_AND_TIME;
      }
      case 10: {
        return plc4c_modbus_read_write_modbus_data_type_sizes_LINT;
      }
      case 11: {
        return plc4c_modbus_read_write_modbus_data_type_sizes_LREAL;
      }
      case 12: {
        return plc4c_modbus_read_write_modbus_data_type_sizes_LTIME;
      }
      case 13: {
        return plc4c_modbus_read_write_modbus_data_type_sizes_LTIME_OF_DAY;
      }
      case 14: {
        return plc4c_modbus_read_write_modbus_data_type_sizes_LWORD;
      }
      case 15: {
        return plc4c_modbus_read_write_modbus_data_type_sizes_REAL;
      }
      case 16: {
        return plc4c_modbus_read_write_modbus_data_type_sizes_SINT;
      }
      case 17: {
        return plc4c_modbus_read_write_modbus_data_type_sizes_STRING;
      }
      case 18: {
        return plc4c_modbus_read_write_modbus_data_type_sizes_TIME;
      }
      case 19: {
        return plc4c_modbus_read_write_modbus_data_type_sizes_TIME_OF_DAY;
      }
      case 20: {
        return plc4c_modbus_read_write_modbus_data_type_sizes_UDINT;
      }
      case 21: {
        return plc4c_modbus_read_write_modbus_data_type_sizes_UINT;
      }
      case 22: {
        return plc4c_modbus_read_write_modbus_data_type_sizes_ULINT;
      }
      case 23: {
        return plc4c_modbus_read_write_modbus_data_type_sizes_USINT;
      }
      case 24: {
        return plc4c_modbus_read_write_modbus_data_type_sizes_WCHAR;
      }
      case 25: {
        return plc4c_modbus_read_write_modbus_data_type_sizes_WORD;
      }
      case 26: {
        return plc4c_modbus_read_write_modbus_data_type_sizes_WSTRING;
      }
      default: {
        return -1;
      }
    }
}

uint8_t plc4c_modbus_read_write_modbus_data_type_sizes_get_data_type_size(plc4c_modbus_read_write_modbus_data_type_sizes value) {
  switch(value) {
    case plc4c_modbus_read_write_modbus_data_type_sizes_BOOL: { /* 'IEC61131_BOOL' */
      return 1;
    }
    case plc4c_modbus_read_write_modbus_data_type_sizes_BYTE: { /* 'IEC61131_BYTE' */
      return 1;
    }
    case plc4c_modbus_read_write_modbus_data_type_sizes_CHAR: { /* 'IEC61131_CHAR' */
      return 1;
    }
    case plc4c_modbus_read_write_modbus_data_type_sizes_DATE: { /* 'IEC61131_DATE' */
      return 8;
    }
    case plc4c_modbus_read_write_modbus_data_type_sizes_DATE_AND_TIME: { /* 'IEC61131_DATE_AND_TIME' */
      return 8;
    }
    case plc4c_modbus_read_write_modbus_data_type_sizes_DINT: { /* 'IEC61131_DINT' */
      return 4;
    }
    case plc4c_modbus_read_write_modbus_data_type_sizes_DWORD: { /* 'IEC61131_DWORD' */
      return 4;
    }
    case plc4c_modbus_read_write_modbus_data_type_sizes_INT: { /* 'IEC61131_INT' */
      return 2;
    }
    case plc4c_modbus_read_write_modbus_data_type_sizes_LDATE: { /* 'IEC61131_LDATE' */
      return 8;
    }
    case plc4c_modbus_read_write_modbus_data_type_sizes_LDATE_AND_TIME: { /* 'IEC61131_LDATE_AND_TIME' */
      return 8;
    }
    case plc4c_modbus_read_write_modbus_data_type_sizes_LINT: { /* 'IEC61131_LINT' */
      return 8;
    }
    case plc4c_modbus_read_write_modbus_data_type_sizes_LREAL: { /* 'IEC61131_LREAL' */
      return 8;
    }
    case plc4c_modbus_read_write_modbus_data_type_sizes_LTIME: { /* 'IEC61131_LTIME' */
      return 8;
    }
    case plc4c_modbus_read_write_modbus_data_type_sizes_LTIME_OF_DAY: { /* 'IEC61131_LTIME_OF_DAY' */
      return 8;
    }
    case plc4c_modbus_read_write_modbus_data_type_sizes_LWORD: { /* 'IEC61131_LWORD' */
      return 8;
    }
    case plc4c_modbus_read_write_modbus_data_type_sizes_REAL: { /* 'IEC61131_REAL' */
      return 4;
    }
    case plc4c_modbus_read_write_modbus_data_type_sizes_SINT: { /* 'IEC61131_SINT' */
      return 1;
    }
    case plc4c_modbus_read_write_modbus_data_type_sizes_STRING: { /* 'IEC61131_STRING' */
      return 1;
    }
    case plc4c_modbus_read_write_modbus_data_type_sizes_TIME: { /* 'IEC61131_TIME' */
      return 8;
    }
    case plc4c_modbus_read_write_modbus_data_type_sizes_TIME_OF_DAY: { /* 'IEC61131_TIME_OF_DAY' */
      return 8;
    }
    case plc4c_modbus_read_write_modbus_data_type_sizes_UDINT: { /* 'IEC61131_UDINT' */
      return 4;
    }
    case plc4c_modbus_read_write_modbus_data_type_sizes_UINT: { /* 'IEC61131_UINT' */
      return 2;
    }
    case plc4c_modbus_read_write_modbus_data_type_sizes_ULINT: { /* 'IEC61131_ULINT' */
      return 8;
    }
    case plc4c_modbus_read_write_modbus_data_type_sizes_USINT: { /* 'IEC61131_USINT' */
      return 1;
    }
    case plc4c_modbus_read_write_modbus_data_type_sizes_WCHAR: { /* 'IEC61131_WCHAR' */
      return 2;
    }
    case plc4c_modbus_read_write_modbus_data_type_sizes_WORD: { /* 'IEC61131_WORD' */
      return 2;
    }
    case plc4c_modbus_read_write_modbus_data_type_sizes_WSTRING: { /* 'IEC61131_WSTRING' */
      return 2;
    }
    default: {
      return 0;
    }
  }
}
