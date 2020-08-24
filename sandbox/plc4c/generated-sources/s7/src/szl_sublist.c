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

#include "szl_sublist.h"
#include <string.h>


// Create an empty NULL-struct
static const plc4c_s7_read_write_szl_sublist plc4c_s7_read_write_szl_sublist_null_const;

plc4c_s7_read_write_szl_sublist plc4c_s7_read_write_szl_sublist_null() {
  return plc4c_s7_read_write_szl_sublist_null_const;
}

plc4c_s7_read_write_szl_sublist plc4c_s7_read_write_szl_sublist_value_of(char* value_string) {
    if(strcmp(value_string, "MODULE_IDENTIFICATION") == 0) {
        return 0x11;
    }
    if(strcmp(value_string, "CPU_FEATURES") == 0) {
        return 0x12;
    }
    if(strcmp(value_string, "USER_MEMORY_AREA") == 0) {
        return 0x13;
    }
    if(strcmp(value_string, "SYSTEM_AREAS") == 0) {
        return 0x14;
    }
    if(strcmp(value_string, "BLOCK_TYPES") == 0) {
        return 0x15;
    }
    if(strcmp(value_string, "STATUS_MODULE_LEDS") == 0) {
        return 0x19;
    }
    if(strcmp(value_string, "COMPONENT_IDENTIFICATION") == 0) {
        return 0x1C;
    }
    if(strcmp(value_string, "INTERRUPT_STATUS") == 0) {
        return 0x22;
    }
    if(strcmp(value_string, "ASSIGNMENT_BETWEEN_PROCESS_IMAGE_PARTITIONS_AND_OBS") == 0) {
        return 0x25;
    }
    if(strcmp(value_string, "COMMUNICATION_STATUS_DATA") == 0) {
        return 0x32;
    }
    if(strcmp(value_string, "STATUS_SINGLE_MODULE_LED") == 0) {
        return 0x74;
    }
    if(strcmp(value_string, "DP_MASTER_SYSTEM_INFORMATION") == 0) {
        return 0x90;
    }
    if(strcmp(value_string, "MODULE_STATUS_INFORMATION") == 0) {
        return 0x91;
    }
    if(strcmp(value_string, "RACK_OR_STATION_STATUS_INFORMATION") == 0) {
        return 0x92;
    }
    if(strcmp(value_string, "RACK_OR_STATION_STATUS_INFORMATION_2") == 0) {
        return 0x94;
    }
    if(strcmp(value_string, "ADDITIONAL_DP_MASTER_SYSTEM_OR_PROFINET_IO_SYSTEM_INFORMATION") == 0) {
        return 0x95;
    }
    if(strcmp(value_string, "MODULE_STATUS_INFORMATION_PROFINET_IO_AND_PROFIBUS_DP") == 0) {
        return 0x96;
    }
    if(strcmp(value_string, "DIAGNOSTIC_BUFFER") == 0) {
        return 0xA0;
    }
    if(strcmp(value_string, "MODULE_DIAGNOSTIC_DATA") == 0) {
        return 0xB1;
    }
    return -1;
}

int plc4c_s7_read_write_szl_sublist_num_values() {
  return 19;
}

plc4c_s7_read_write_szl_sublist plc4c_s7_read_write_szl_sublist_value_for_index(int index) {
    switch(index) {
      case 0: {
        return 0x11;
      }
      case 1: {
        return 0x12;
      }
      case 2: {
        return 0x13;
      }
      case 3: {
        return 0x14;
      }
      case 4: {
        return 0x15;
      }
      case 5: {
        return 0x19;
      }
      case 6: {
        return 0x1C;
      }
      case 7: {
        return 0x22;
      }
      case 8: {
        return 0x25;
      }
      case 9: {
        return 0x32;
      }
      case 10: {
        return 0x74;
      }
      case 11: {
        return 0x90;
      }
      case 12: {
        return 0x91;
      }
      case 13: {
        return 0x92;
      }
      case 14: {
        return 0x94;
      }
      case 15: {
        return 0x95;
      }
      case 16: {
        return 0x96;
      }
      case 17: {
        return 0xA0;
      }
      case 18: {
        return 0xB1;
      }
      default: {
        return -1;
      }
    }
}
