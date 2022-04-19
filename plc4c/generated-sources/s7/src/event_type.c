/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

#include "event_type.h"
#include <string.h>

// Code generated by code-generation. DO NOT EDIT.


// Create an empty NULL-struct
static const plc4c_s7_read_write_event_type plc4c_s7_read_write_event_type_null_const;

plc4c_s7_read_write_event_type plc4c_s7_read_write_event_type_null() {
  return plc4c_s7_read_write_event_type_null_const;
}

// Parse function.
plc4c_return_code plc4c_s7_read_write_event_type_parse(plc4c_spi_read_buffer* readBuffer, plc4c_s7_read_write_event_type* _message) {
    plc4c_return_code _res = OK;

    _res = plc4c_spi_read_unsigned_byte(readBuffer, 8, (uint8_t*) _message);

    return _res;
}

plc4c_return_code plc4c_s7_read_write_event_type_serialize(plc4c_spi_write_buffer* writeBuffer, plc4c_s7_read_write_event_type* _message) {
    plc4c_return_code _res = OK;

    _res = plc4c_spi_write_unsigned_byte(writeBuffer, 8, *_message);

    return _res;
}

plc4c_s7_read_write_event_type plc4c_s7_read_write_event_type_value_of(char* value_string) {
    if(strcmp(value_string, "MODE") == 0) {
        return plc4c_s7_read_write_event_type_MODE;
    }
    if(strcmp(value_string, "SYS") == 0) {
        return plc4c_s7_read_write_event_type_SYS;
    }
    if(strcmp(value_string, "USR") == 0) {
        return plc4c_s7_read_write_event_type_USR;
    }
    if(strcmp(value_string, "ALM") == 0) {
        return plc4c_s7_read_write_event_type_ALM;
    }
    return -1;
}

int plc4c_s7_read_write_event_type_num_values() {
  return 4;
}

plc4c_s7_read_write_event_type plc4c_s7_read_write_event_type_value_for_index(int index) {
    switch(index) {
      case 0: {
        return plc4c_s7_read_write_event_type_MODE;
      }
      case 1: {
        return plc4c_s7_read_write_event_type_SYS;
      }
      case 2: {
        return plc4c_s7_read_write_event_type_USR;
      }
      case 3: {
        return plc4c_s7_read_write_event_type_ALM;
      }
      default: {
        return -1;
      }
    }
}

uint16_t plc4c_s7_read_write_event_type_length_in_bytes(plc4c_s7_read_write_event_type* _message) {
    return plc4c_s7_read_write_event_type_length_in_bits(_message) / 8;
}

uint16_t plc4c_s7_read_write_event_type_length_in_bits(plc4c_s7_read_write_event_type* _message) {
    return 8;
}
