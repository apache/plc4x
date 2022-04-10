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

#include "cotp_protocol_class.h"
#include <string.h>

// Code generated by code-generation. DO NOT EDIT.


// Create an empty NULL-struct
static const plc4c_s7_read_write_cotp_protocol_class plc4c_s7_read_write_cotp_protocol_class_null_const;

plc4c_s7_read_write_cotp_protocol_class plc4c_s7_read_write_cotp_protocol_class_null() {
  return plc4c_s7_read_write_cotp_protocol_class_null_const;
}

// Parse function.
plc4c_return_code plc4c_s7_read_write_cotp_protocol_class_parse(plc4c_spi_read_buffer* readBuffer, plc4c_s7_read_write_cotp_protocol_class** _message) {
    plc4c_return_code _res = OK;

    // Allocate enough memory to contain this data structure.
    (*_message) = malloc(sizeof(plc4c_s7_read_write_cotp_protocol_class));
    if(*_message == NULL) {
        return NO_MEMORY;
    }

    _res = plc4c_spi_read_unsigned_byte(readBuffer, 8, (uint8_t*) *_message);

    return _res;
}

plc4c_return_code plc4c_s7_read_write_cotp_protocol_class_serialize(plc4c_spi_write_buffer* writeBuffer, plc4c_s7_read_write_cotp_protocol_class* _message) {
    plc4c_return_code _res = OK;

    _res = plc4c_spi_write_unsigned_byte(writeBuffer, 8, *_message);

    return _res;
}

plc4c_s7_read_write_cotp_protocol_class plc4c_s7_read_write_cotp_protocol_class_value_of(char* value_string) {
    if(strcmp(value_string, "CLASS_0") == 0) {
        return plc4c_s7_read_write_cotp_protocol_class_CLASS_0;
    }
    if(strcmp(value_string, "CLASS_1") == 0) {
        return plc4c_s7_read_write_cotp_protocol_class_CLASS_1;
    }
    if(strcmp(value_string, "CLASS_2") == 0) {
        return plc4c_s7_read_write_cotp_protocol_class_CLASS_2;
    }
    if(strcmp(value_string, "CLASS_3") == 0) {
        return plc4c_s7_read_write_cotp_protocol_class_CLASS_3;
    }
    if(strcmp(value_string, "CLASS_4") == 0) {
        return plc4c_s7_read_write_cotp_protocol_class_CLASS_4;
    }
    return -1;
}

int plc4c_s7_read_write_cotp_protocol_class_num_values() {
  return 5;
}

plc4c_s7_read_write_cotp_protocol_class plc4c_s7_read_write_cotp_protocol_class_value_for_index(int index) {
    switch(index) {
      case 0: {
        return plc4c_s7_read_write_cotp_protocol_class_CLASS_0;
      }
      case 1: {
        return plc4c_s7_read_write_cotp_protocol_class_CLASS_1;
      }
      case 2: {
        return plc4c_s7_read_write_cotp_protocol_class_CLASS_2;
      }
      case 3: {
        return plc4c_s7_read_write_cotp_protocol_class_CLASS_3;
      }
      case 4: {
        return plc4c_s7_read_write_cotp_protocol_class_CLASS_4;
      }
      default: {
        return -1;
      }
    }
}

uint16_t plc4c_s7_read_write_cotp_protocol_class_length_in_bytes(plc4c_s7_read_write_cotp_protocol_class* _message) {
    return plc4c_s7_read_write_cotp_protocol_class_length_in_bits(_message) / 8;
}

uint16_t plc4c_s7_read_write_cotp_protocol_class_length_in_bits(plc4c_s7_read_write_cotp_protocol_class* _message) {
    return 8;
}
