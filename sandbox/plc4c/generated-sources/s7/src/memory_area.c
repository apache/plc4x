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

#include "memory_area.h"
#include <string.h>


// Create an empty NULL-struct
static const plc4c_s7_read_write_memory_area plc4c_s7_read_write_memory_area_null_const;

plc4c_s7_read_write_memory_area plc4c_s7_read_write_memory_area_null() {
  return plc4c_s7_read_write_memory_area_null_const;
}

plc4c_s7_read_write_memory_area plc4c_s7_read_write_memory_area_value_of(char* value_string) {
    if(strcmp(value_string, "COUNTERS") == 0) {
        return 0x1C;
    }
    if(strcmp(value_string, "TIMERS") == 0) {
        return 0x1D;
    }
    if(strcmp(value_string, "DIRECT_PERIPHERAL_ACCESS") == 0) {
        return 0x80;
    }
    if(strcmp(value_string, "INPUTS") == 0) {
        return 0x81;
    }
    if(strcmp(value_string, "OUTPUTS") == 0) {
        return 0x82;
    }
    if(strcmp(value_string, "FLAGS_MARKERS") == 0) {
        return 0x83;
    }
    if(strcmp(value_string, "DATA_BLOCKS") == 0) {
        return 0x84;
    }
    if(strcmp(value_string, "INSTANCE_DATA_BLOCKS") == 0) {
        return 0x85;
    }
    if(strcmp(value_string, "LOCAL_DATA") == 0) {
        return 0x86;
    }
    return -1;
}

int plc4c_s7_read_write_memory_area_num_values() {
  return 9;
}

plc4c_s7_read_write_memory_area plc4c_s7_read_write_memory_area_value_for_index(int index) {
    switch(index) {
      case 0: {
        return 0x1C;
      }
      case 1: {
        return 0x1D;
      }
      case 2: {
        return 0x80;
      }
      case 3: {
        return 0x81;
      }
      case 4: {
        return 0x82;
      }
      case 5: {
        return 0x83;
      }
      case 6: {
        return 0x84;
      }
      case 7: {
        return 0x85;
      }
      case 8: {
        return 0x86;
      }
      default: {
        return -1;
      }
    }
}

char* plc4c_s7_read_write_memory_area_get_short_name(plc4c_s7_read_write_memory_area value) {
  switch(value) {
    case 28: { /* '0x1C' */
      return "C";
    }
    case 29: { /* '0x1D' */
      return "T";
    }
    case 128: { /* '0x80' */
      return "D";
    }
    case 129: { /* '0x81' */
      return "I";
    }
    case 130: { /* '0x82' */
      return "Q";
    }
    case 131: { /* '0x83' */
      return "M";
    }
    case 132: { /* '0x84' */
      return "DB";
    }
    case 133: { /* '0x85' */
      return "DBI";
    }
    case 134: { /* '0x86' */
      return "LD";
    }
    default: {
      return 0;
    }
  }
}
