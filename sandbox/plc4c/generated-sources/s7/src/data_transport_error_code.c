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

#include "data_transport_error_code.h"
#include <string.h>


// Create an empty NULL-struct
static const plc4c_s7_read_write_data_transport_error_code plc4c_s7_read_write_data_transport_error_code_null_const;

plc4c_s7_read_write_data_transport_error_code plc4c_s7_read_write_data_transport_error_code_null() {
  return plc4c_s7_read_write_data_transport_error_code_null_const;
}

plc4c_s7_read_write_data_transport_error_code plc4c_s7_read_write_data_transport_error_code_value_of(char* value_string) {
    if(strcmp(value_string, "RESERVED") == 0) {
        return 0x00;
    }
    if(strcmp(value_string, "ACCESS_DENIED") == 0) {
        return 0x03;
    }
    if(strcmp(value_string, "INVALID_ADDRESS") == 0) {
        return 0x05;
    }
    if(strcmp(value_string, "DATA_TYPE_NOT_SUPPORTED") == 0) {
        return 0x06;
    }
    if(strcmp(value_string, "NOT_FOUND") == 0) {
        return 0x0A;
    }
    if(strcmp(value_string, "OK") == 0) {
        return 0xFF;
    }
    return -1;
}

int plc4c_s7_read_write_data_transport_error_code_num_values() {
  return 6;
}

plc4c_s7_read_write_data_transport_error_code plc4c_s7_read_write_data_transport_error_code_value_for_index(int index) {
    switch(index) {
      case 0: {
        return 0x00;
      }
      case 1: {
        return 0x03;
      }
      case 2: {
        return 0x05;
      }
      case 3: {
        return 0x06;
      }
      case 4: {
        return 0x0A;
      }
      case 5: {
        return 0xFF;
      }
      default: {
        return -1;
      }
    }
}
