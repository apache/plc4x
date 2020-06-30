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

#ifndef PLC4C_S7_READ_WRITE_MEMORY_AREA_H_
#define PLC4C_S7_READ_WRITE_MEMORY_AREA_H_
#ifdef __cplusplus
extern "C" {
#endif

#include <stdbool.h>

enum plc4c_s7_read_write_memory_area {
  plc4c_s7_read_write_memory_area_COUNTERS = 0x1C,
  plc4c_s7_read_write_memory_area_TIMERS = 0x1D,
  plc4c_s7_read_write_memory_area_DIRECT_PERIPHERAL_ACCESS = 0x80,
  plc4c_s7_read_write_memory_area_INPUTS = 0x81,
  plc4c_s7_read_write_memory_area_OUTPUTS = 0x82,
  plc4c_s7_read_write_memory_area_FLAGS_MARKERS = 0x83,
  plc4c_s7_read_write_memory_area_DATA_BLOCKS = 0x84,
  plc4c_s7_read_write_memory_area_INSTANCE_DATA_BLOCKS = 0x85,
  plc4c_s7_read_write_memory_area_LOCAL_DATA = 0x86
};
typedef enum plc4c_s7_read_write_memory_area plc4c_s7_read_write_memory_area;

// Create an empty NULL-struct
static const plc4c_s7_read_write_memory_area plc4c_s7_read_write_memory_area_null;


char * plc4c_s7_read_write_memory_area_get_short_name(plc4c_s7_read_write_memory_area value) {
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

#ifdef __cplusplus
}
#endif
#endif  // PLC4C_S7_READ_WRITE_MEMORY_AREA_H_
