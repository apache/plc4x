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

#include "cotp_tpdu_size.h"
#include <string.h>


// Create an empty NULL-struct
static const plc4c_s7_read_write_cotp_tpdu_size plc4c_s7_read_write_cotp_tpdu_size_null_const;

plc4c_s7_read_write_cotp_tpdu_size plc4c_s7_read_write_cotp_tpdu_size_null() {
  return plc4c_s7_read_write_cotp_tpdu_size_null_const;
}

plc4c_s7_read_write_cotp_tpdu_size plc4c_s7_read_write_cotp_tpdu_size_value_of(char* value_string) {
    if(strcmp(value_string, "SIZE_128") == 0) {
        return 0x07;
    }
    if(strcmp(value_string, "SIZE_256") == 0) {
        return 0x08;
    }
    if(strcmp(value_string, "SIZE_512") == 0) {
        return 0x09;
    }
    if(strcmp(value_string, "SIZE_1024") == 0) {
        return 0x0a;
    }
    if(strcmp(value_string, "SIZE_2048") == 0) {
        return 0x0b;
    }
    if(strcmp(value_string, "SIZE_4096") == 0) {
        return 0x0c;
    }
    if(strcmp(value_string, "SIZE_8192") == 0) {
        return 0x0d;
    }
    return -1;
}

int plc4c_s7_read_write_cotp_tpdu_size_num_values() {
  return 7;
}

plc4c_s7_read_write_cotp_tpdu_size plc4c_s7_read_write_cotp_tpdu_size_value_for_index(int index) {
    switch(index) {
      case 0: {
        return 0x07;
      }
      case 1: {
        return 0x08;
      }
      case 2: {
        return 0x09;
      }
      case 3: {
        return 0x0a;
      }
      case 4: {
        return 0x0b;
      }
      case 5: {
        return 0x0c;
      }
      case 6: {
        return 0x0d;
      }
      default: {
        return -1;
      }
    }
}

uint16_t plc4c_s7_read_write_cotp_tpdu_size_get_size_in_bytes(plc4c_s7_read_write_cotp_tpdu_size value) {
  switch(value) {
    case 7: { /* '0x07' */
      return 128;
    }
    case 8: { /* '0x08' */
      return 256;
    }
    case 9: { /* '0x09' */
      return 512;
    }
    case 10: { /* '0x0a' */
      return 1024;
    }
    case 11: { /* '0x0b' */
      return 2048;
    }
    case 12: { /* '0x0c' */
      return 4096;
    }
    case 13: { /* '0x0d' */
      return 8192;
    }
    default: {
      return 0;
    }
  }
}
