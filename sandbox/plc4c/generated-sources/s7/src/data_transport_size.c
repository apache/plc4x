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

#include "data_transport_size.h"
#include <string.h>


// Create an empty NULL-struct
static const plc4c_s7_read_write_data_transport_size plc4c_s7_read_write_data_transport_size_null_const;

plc4c_s7_read_write_data_transport_size plc4c_s7_read_write_data_transport_size_null() {
  return plc4c_s7_read_write_data_transport_size_null_const;
}

plc4c_s7_read_write_data_transport_size plc4c_s7_read_write_data_transport_size_value_of(char* value_string) {
    if(strcmp(value_string, "NULL") == 0) {
        return 0x00;
    }
    if(strcmp(value_string, "BIT") == 0) {
        return 0x03;
    }
    if(strcmp(value_string, "BYTE_WORD_DWORD") == 0) {
        return 0x04;
    }
    if(strcmp(value_string, "INTEGER") == 0) {
        return 0x05;
    }
    if(strcmp(value_string, "DINTEGER") == 0) {
        return 0x06;
    }
    if(strcmp(value_string, "REAL") == 0) {
        return 0x07;
    }
    if(strcmp(value_string, "OCTET_STRING") == 0) {
        return 0x09;
    }
    return -1;
}

int plc4c_s7_read_write_data_transport_size_num_values() {
  return 7;
}

plc4c_s7_read_write_data_transport_size plc4c_s7_read_write_data_transport_size_value_for_index(int index) {
    switch(index) {
      case 0: {
        return 0x00;
      }
      case 1: {
        return 0x03;
      }
      case 2: {
        return 0x04;
      }
      case 3: {
        return 0x05;
      }
      case 4: {
        return 0x06;
      }
      case 5: {
        return 0x07;
      }
      case 6: {
        return 0x09;
      }
      default: {
        return -1;
      }
    }
}

bool plc4c_s7_read_write_data_transport_size_get_size_in_bits(plc4c_s7_read_write_data_transport_size value) {
  switch(value) {
    case 0: { /* '0x00' */
      return false;
    }
    case 3: { /* '0x03' */
      return true;
    }
    case 4: { /* '0x04' */
      return true;
    }
    case 5: { /* '0x05' */
      return true;
    }
    case 6: { /* '0x06' */
      return false;
    }
    case 7: { /* '0x07' */
      return false;
    }
    case 9: { /* '0x09' */
      return false;
    }
    default: {
      return 0;
    }
  }
}
