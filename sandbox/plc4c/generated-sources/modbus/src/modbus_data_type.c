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

#include "modbus_data_type.h"
#include <string.h>


// Create an empty NULL-struct
static const plc4c_modbus_read_write_modbus_data_type plc4c_modbus_read_write_modbus_data_type_null_const;

plc4c_modbus_read_write_modbus_data_type plc4c_modbus_read_write_modbus_data_type_null() {
  return plc4c_modbus_read_write_modbus_data_type_null_const;
}

plc4c_modbus_read_write_modbus_data_type plc4c_modbus_read_write_modbus_data_type_value_of(char* value_string) {
    if(strcmp(value_string, "NULL") == 0) {
        return plc4c_modbus_read_write_modbus_data_type_NULL;
    }
    if(strcmp(value_string, "BOOL") == 0) {
        return plc4c_modbus_read_write_modbus_data_type_BOOL;
    }
    if(strcmp(value_string, "BYTE") == 0) {
        return plc4c_modbus_read_write_modbus_data_type_BYTE;
    }
    if(strcmp(value_string, "WORD") == 0) {
        return plc4c_modbus_read_write_modbus_data_type_WORD;
    }
    if(strcmp(value_string, "DWORD") == 0) {
        return plc4c_modbus_read_write_modbus_data_type_DWORD;
    }
    if(strcmp(value_string, "LWORD") == 0) {
        return plc4c_modbus_read_write_modbus_data_type_LWORD;
    }
    if(strcmp(value_string, "SINT") == 0) {
        return plc4c_modbus_read_write_modbus_data_type_SINT;
    }
    if(strcmp(value_string, "INT") == 0) {
        return plc4c_modbus_read_write_modbus_data_type_INT;
    }
    if(strcmp(value_string, "DINT") == 0) {
        return plc4c_modbus_read_write_modbus_data_type_DINT;
    }
    if(strcmp(value_string, "LINT") == 0) {
        return plc4c_modbus_read_write_modbus_data_type_LINT;
    }
    if(strcmp(value_string, "USINT") == 0) {
        return plc4c_modbus_read_write_modbus_data_type_USINT;
    }
    if(strcmp(value_string, "UINT") == 0) {
        return plc4c_modbus_read_write_modbus_data_type_UINT;
    }
    if(strcmp(value_string, "UDINT") == 0) {
        return plc4c_modbus_read_write_modbus_data_type_UDINT;
    }
    if(strcmp(value_string, "ULINT") == 0) {
        return plc4c_modbus_read_write_modbus_data_type_ULINT;
    }
    if(strcmp(value_string, "REAL") == 0) {
        return plc4c_modbus_read_write_modbus_data_type_REAL;
    }
    if(strcmp(value_string, "LREAL") == 0) {
        return plc4c_modbus_read_write_modbus_data_type_LREAL;
    }
    if(strcmp(value_string, "TIME") == 0) {
        return plc4c_modbus_read_write_modbus_data_type_TIME;
    }
    if(strcmp(value_string, "LTIME") == 0) {
        return plc4c_modbus_read_write_modbus_data_type_LTIME;
    }
    if(strcmp(value_string, "DATE") == 0) {
        return plc4c_modbus_read_write_modbus_data_type_DATE;
    }
    if(strcmp(value_string, "LDATE") == 0) {
        return plc4c_modbus_read_write_modbus_data_type_LDATE;
    }
    if(strcmp(value_string, "TIME_OF_DAY") == 0) {
        return plc4c_modbus_read_write_modbus_data_type_TIME_OF_DAY;
    }
    if(strcmp(value_string, "LTIME_OF_DAY") == 0) {
        return plc4c_modbus_read_write_modbus_data_type_LTIME_OF_DAY;
    }
    if(strcmp(value_string, "DATE_AND_TIME") == 0) {
        return plc4c_modbus_read_write_modbus_data_type_DATE_AND_TIME;
    }
    if(strcmp(value_string, "LDATE_AND_TIME") == 0) {
        return plc4c_modbus_read_write_modbus_data_type_LDATE_AND_TIME;
    }
    if(strcmp(value_string, "CHAR") == 0) {
        return plc4c_modbus_read_write_modbus_data_type_CHAR;
    }
    if(strcmp(value_string, "WCHAR") == 0) {
        return plc4c_modbus_read_write_modbus_data_type_WCHAR;
    }
    if(strcmp(value_string, "STRING") == 0) {
        return plc4c_modbus_read_write_modbus_data_type_STRING;
    }
    if(strcmp(value_string, "WSTRING") == 0) {
        return plc4c_modbus_read_write_modbus_data_type_WSTRING;
    }
    return -1;
}

int plc4c_modbus_read_write_modbus_data_type_num_values() {
  return 28;
}

plc4c_modbus_read_write_modbus_data_type plc4c_modbus_read_write_modbus_data_type_value_for_index(int index) {
    switch(index) {
      case 0: {
        return plc4c_modbus_read_write_modbus_data_type_NULL;
      }
      case 1: {
        return plc4c_modbus_read_write_modbus_data_type_BOOL;
      }
      case 2: {
        return plc4c_modbus_read_write_modbus_data_type_BYTE;
      }
      case 3: {
        return plc4c_modbus_read_write_modbus_data_type_WORD;
      }
      case 4: {
        return plc4c_modbus_read_write_modbus_data_type_DWORD;
      }
      case 5: {
        return plc4c_modbus_read_write_modbus_data_type_LWORD;
      }
      case 6: {
        return plc4c_modbus_read_write_modbus_data_type_SINT;
      }
      case 7: {
        return plc4c_modbus_read_write_modbus_data_type_INT;
      }
      case 8: {
        return plc4c_modbus_read_write_modbus_data_type_DINT;
      }
      case 9: {
        return plc4c_modbus_read_write_modbus_data_type_LINT;
      }
      case 10: {
        return plc4c_modbus_read_write_modbus_data_type_USINT;
      }
      case 11: {
        return plc4c_modbus_read_write_modbus_data_type_UINT;
      }
      case 12: {
        return plc4c_modbus_read_write_modbus_data_type_UDINT;
      }
      case 13: {
        return plc4c_modbus_read_write_modbus_data_type_ULINT;
      }
      case 14: {
        return plc4c_modbus_read_write_modbus_data_type_REAL;
      }
      case 15: {
        return plc4c_modbus_read_write_modbus_data_type_LREAL;
      }
      case 16: {
        return plc4c_modbus_read_write_modbus_data_type_TIME;
      }
      case 17: {
        return plc4c_modbus_read_write_modbus_data_type_LTIME;
      }
      case 18: {
        return plc4c_modbus_read_write_modbus_data_type_DATE;
      }
      case 19: {
        return plc4c_modbus_read_write_modbus_data_type_LDATE;
      }
      case 20: {
        return plc4c_modbus_read_write_modbus_data_type_TIME_OF_DAY;
      }
      case 21: {
        return plc4c_modbus_read_write_modbus_data_type_LTIME_OF_DAY;
      }
      case 22: {
        return plc4c_modbus_read_write_modbus_data_type_DATE_AND_TIME;
      }
      case 23: {
        return plc4c_modbus_read_write_modbus_data_type_LDATE_AND_TIME;
      }
      case 24: {
        return plc4c_modbus_read_write_modbus_data_type_CHAR;
      }
      case 25: {
        return plc4c_modbus_read_write_modbus_data_type_WCHAR;
      }
      case 26: {
        return plc4c_modbus_read_write_modbus_data_type_STRING;
      }
      case 27: {
        return plc4c_modbus_read_write_modbus_data_type_WSTRING;
      }
      default: {
        return -1;
      }
    }
}

uint8_t plc4c_modbus_read_write_modbus_data_type_get_data_type_size(plc4c_modbus_read_write_modbus_data_type value) {
  switch(value) {
    case plc4c_modbus_read_write_modbus_data_type_NULL: { /* '00' */
      return 0;
    }
    case plc4c_modbus_read_write_modbus_data_type_BOOL: { /* '01' */
      return 1;
    }
    case plc4c_modbus_read_write_modbus_data_type_BYTE: { /* '10' */
      return 1;
    }
    case plc4c_modbus_read_write_modbus_data_type_WORD: { /* '11' */
      return 2;
    }
    case plc4c_modbus_read_write_modbus_data_type_DWORD: { /* '12' */
      return 4;
    }
    case plc4c_modbus_read_write_modbus_data_type_LWORD: { /* '13' */
      return 8;
    }
    case plc4c_modbus_read_write_modbus_data_type_SINT: { /* '20' */
      return 1;
    }
    case plc4c_modbus_read_write_modbus_data_type_INT: { /* '21' */
      return 2;
    }
    case plc4c_modbus_read_write_modbus_data_type_DINT: { /* '22' */
      return 4;
    }
    case plc4c_modbus_read_write_modbus_data_type_LINT: { /* '23' */
      return 8;
    }
    case plc4c_modbus_read_write_modbus_data_type_USINT: { /* '24' */
      return 1;
    }
    case plc4c_modbus_read_write_modbus_data_type_UINT: { /* '25' */
      return 2;
    }
    case plc4c_modbus_read_write_modbus_data_type_UDINT: { /* '26' */
      return 4;
    }
    case plc4c_modbus_read_write_modbus_data_type_ULINT: { /* '27' */
      return 8;
    }
    case plc4c_modbus_read_write_modbus_data_type_REAL: { /* '30' */
      return 4;
    }
    case plc4c_modbus_read_write_modbus_data_type_LREAL: { /* '31' */
      return 8;
    }
    case plc4c_modbus_read_write_modbus_data_type_TIME: { /* '40' */
      return 8;
    }
    case plc4c_modbus_read_write_modbus_data_type_LTIME: { /* '41' */
      return 8;
    }
    case plc4c_modbus_read_write_modbus_data_type_DATE: { /* '50' */
      return 8;
    }
    case plc4c_modbus_read_write_modbus_data_type_LDATE: { /* '51' */
      return 8;
    }
    case plc4c_modbus_read_write_modbus_data_type_TIME_OF_DAY: { /* '60' */
      return 8;
    }
    case plc4c_modbus_read_write_modbus_data_type_LTIME_OF_DAY: { /* '61' */
      return 8;
    }
    case plc4c_modbus_read_write_modbus_data_type_DATE_AND_TIME: { /* '70' */
      return 8;
    }
    case plc4c_modbus_read_write_modbus_data_type_LDATE_AND_TIME: { /* '71' */
      return 8;
    }
    case plc4c_modbus_read_write_modbus_data_type_CHAR: { /* '80' */
      return 1;
    }
    case plc4c_modbus_read_write_modbus_data_type_WCHAR: { /* '81' */
      return 2;
    }
    case plc4c_modbus_read_write_modbus_data_type_STRING: { /* '82' */
      return 1;
    }
    case plc4c_modbus_read_write_modbus_data_type_WSTRING: { /* '83' */
      return 2;
    }
    default: {
      return 0;
    }
  }
}
