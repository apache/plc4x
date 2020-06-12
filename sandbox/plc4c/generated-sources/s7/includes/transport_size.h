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

#ifndef PLC4C_S7_READ_WRITE_TRANSPORT_SIZE_H_
#define PLC4C_S7_READ_WRITE_TRANSPORT_SIZE_H_
#ifdef __cplusplus
extern "C" {
#endif

#include <stdbool.h>
#include "transport_size.h"
#include "data_transport_size.h"

enum plc4c_s7_read_write_transport_size {
  plc4c_s7_read_write_transport_size_BOOL = 0x01,
  plc4c_s7_read_write_transport_size_BYTE = 0x02,
  plc4c_s7_read_write_transport_size_WORD = 0x04,
  plc4c_s7_read_write_transport_size_DWORD = 0x06,
  plc4c_s7_read_write_transport_size_LWORD = 0x00,
  plc4c_s7_read_write_transport_size_INT = 0x05,
  plc4c_s7_read_write_transport_size_UINT = 0x05,
  plc4c_s7_read_write_transport_size_SINT = 0x02,
  plc4c_s7_read_write_transport_size_USINT = 0x02,
  plc4c_s7_read_write_transport_size_DINT = 0x07,
  plc4c_s7_read_write_transport_size_UDINT = 0x07,
  plc4c_s7_read_write_transport_size_LINT = 0x00,
  plc4c_s7_read_write_transport_size_ULINT = 0x00,
  plc4c_s7_read_write_transport_size_REAL = 0x08,
  plc4c_s7_read_write_transport_size_LREAL = 0x00,
  plc4c_s7_read_write_transport_size_CHAR = 0x03,
  plc4c_s7_read_write_transport_size_WCHAR = 0x13,
  plc4c_s7_read_write_transport_size_STRING = 0x03,
  plc4c_s7_read_write_transport_size_WSTRING = 0x00,
  plc4c_s7_read_write_transport_size_TIME = 0x0B,
  plc4c_s7_read_write_transport_size_LTIME = 0x00,
  plc4c_s7_read_write_transport_size_DATE = 0x02,
  plc4c_s7_read_write_transport_size_TIME_OF_DAY = 0x02,
  plc4c_s7_read_write_transport_size_DATE_AND_TIME = 0x02
};
typedef enum plc4c_s7_read_write_transport_size plc4c_s7_read_write_transport_size;


bool plc4c_s7_read_write_transport_size_get_supported__s7_300(plc4c_s7_read_write_transport_size value) {
  switch(value) {
    case 0: { /* '0x00' */
      return false;
    }
    case 1: { /* '0x01' */
      return true;
    }
    case 2: { /* '0x02' */
      return true;
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
      return true;
    }
    case 7: { /* '0x07' */
      return true;
    }
    case 8: { /* '0x08' */
      return true;
    }
    case 11: { /* '0x0B' */
      return true;
    }
    case 19: { /* '0x13' */
      return false;
    }
    default: {
      return 0;
    }
  }
}

bool plc4c_s7_read_write_transport_size_get_supported__logo(plc4c_s7_read_write_transport_size value) {
  switch(value) {
    case 0: { /* '0x00' */
      return false;
    }
    case 1: { /* '0x01' */
      return true;
    }
    case 2: { /* '0x02' */
      return true;
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
      return true;
    }
    case 7: { /* '0x07' */
      return true;
    }
    case 8: { /* '0x08' */
      return true;
    }
    case 11: { /* '0x0B' */
      return true;
    }
    case 19: { /* '0x13' */
      return true;
    }
    default: {
      return 0;
    }
  }
}

uint8_t plc4c_s7_read_write_transport_size_get_size_in_bytes(plc4c_s7_read_write_transport_size value) {
  switch(value) {
    case 0: { /* '0x00' */
      return 8;
    }
    case 1: { /* '0x01' */
      return 1;
    }
    case 2: { /* '0x02' */
      return 1;
    }
    case 3: { /* '0x03' */
      return 1;
    }
    case 4: { /* '0x04' */
      return 2;
    }
    case 5: { /* '0x05' */
      return 2;
    }
    case 6: { /* '0x06' */
      return 4;
    }
    case 7: { /* '0x07' */
      return 4;
    }
    case 8: { /* '0x08' */
      return 4;
    }
    case 11: { /* '0x0B' */
      return 4;
    }
    case 19: { /* '0x13' */
      return 2;
    }
    default: {
      return 0;
    }
  }
}

bool plc4c_s7_read_write_transport_size_get_supported__s7_400(plc4c_s7_read_write_transport_size value) {
  switch(value) {
    case 0: { /* '0x00' */
      return false;
    }
    case 1: { /* '0x01' */
      return true;
    }
    case 2: { /* '0x02' */
      return true;
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
      return true;
    }
    case 7: { /* '0x07' */
      return true;
    }
    case 8: { /* '0x08' */
      return true;
    }
    case 11: { /* '0x0B' */
      return true;
    }
    case 19: { /* '0x13' */
      return false;
    }
    default: {
      return 0;
    }
  }
}

bool plc4c_s7_read_write_transport_size_get_supported__s7_1200(plc4c_s7_read_write_transport_size value) {
  switch(value) {
    case 0: { /* '0x00' */
      return false;
    }
    case 1: { /* '0x01' */
      return true;
    }
    case 2: { /* '0x02' */
      return true;
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
      return true;
    }
    case 7: { /* '0x07' */
      return true;
    }
    case 8: { /* '0x08' */
      return true;
    }
    case 11: { /* '0x0B' */
      return true;
    }
    case 19: { /* '0x13' */
      return true;
    }
    default: {
      return 0;
    }
  }
}

uint8_t plc4c_s7_read_write_transport_size_get_size_code(plc4c_s7_read_write_transport_size value) {
  switch(value) {
    case 0: { /* '0x00' */
      return 'X';
    }
    case 1: { /* '0x01' */
      return 'X';
    }
    case 2: { /* '0x02' */
      return 'B';
    }
    case 3: { /* '0x03' */
      return 'B';
    }
    case 4: { /* '0x04' */
      return 'W';
    }
    case 5: { /* '0x05' */
      return 'W';
    }
    case 6: { /* '0x06' */
      return 'D';
    }
    case 7: { /* '0x07' */
      return 'D';
    }
    case 8: { /* '0x08' */
      return 'D';
    }
    case 11: { /* '0x0B' */
      return 'X';
    }
    case 19: { /* '0x13' */
      return 'X';
    }
    default: {
      return 0;
    }
  }
}

bool plc4c_s7_read_write_transport_size_get_supported__s7_1500(plc4c_s7_read_write_transport_size value) {
  switch(value) {
    case 0: { /* '0x00' */
      return true;
    }
    case 1: { /* '0x01' */
      return true;
    }
    case 2: { /* '0x02' */
      return true;
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
      return true;
    }
    case 7: { /* '0x07' */
      return true;
    }
    case 8: { /* '0x08' */
      return true;
    }
    case 11: { /* '0x0B' */
      return true;
    }
    case 19: { /* '0x13' */
      return true;
    }
    default: {
      return 0;
    }
  }
}

plc4c_s7_read_write_data_transport_size plc4c_s7_read_write_transport_size_get_data_transport_size(plc4c_s7_read_write_transport_size value) {
  switch(value) {
    case 0: { /* '0x00' */
      return -1;
    }
    case 1: { /* '0x01' */
      return plc4c_s7_read_write_data_transport_size_BIT;
    }
    case 2: { /* '0x02' */
      return plc4c_s7_read_write_data_transport_size_BYTE_WORD_DWORD;
    }
    case 3: { /* '0x03' */
      return plc4c_s7_read_write_data_transport_size_BYTE_WORD_DWORD;
    }
    case 4: { /* '0x04' */
      return plc4c_s7_read_write_data_transport_size_BYTE_WORD_DWORD;
    }
    case 5: { /* '0x05' */
      return plc4c_s7_read_write_data_transport_size_BYTE_WORD_DWORD;
    }
    case 6: { /* '0x06' */
      return plc4c_s7_read_write_data_transport_size_BYTE_WORD_DWORD;
    }
    case 7: { /* '0x07' */
      return plc4c_s7_read_write_data_transport_size_BYTE_WORD_DWORD;
    }
    case 8: { /* '0x08' */
      return plc4c_s7_read_write_data_transport_size_BYTE_WORD_DWORD;
    }
    case 11: { /* '0x0B' */
      return -1;
    }
    case 19: { /* '0x13' */
      return -1;
    }
    default: {
      return 0;
    }
  }
}

plc4c_s7_read_write_transport_size plc4c_s7_read_write_transport_size_get_base_type(plc4c_s7_read_write_transport_size value) {
  switch(value) {
    case 0: { /* '0x00' */
      return -1;
    }
    case 1: { /* '0x01' */
      return -1;
    }
    case 2: { /* '0x02' */
      return -1;
    }
    case 3: { /* '0x03' */
      return -1;
    }
    case 4: { /* '0x04' */
      return -1;
    }
    case 5: { /* '0x05' */
      return -1;
    }
    case 6: { /* '0x06' */
      return plc4c_s7_read_write_transport_size_WORD;
    }
    case 7: { /* '0x07' */
      return plc4c_s7_read_write_transport_size_INT;
    }
    case 8: { /* '0x08' */
      return -1;
    }
    case 11: { /* '0x0B' */
      return -1;
    }
    case 19: { /* '0x13' */
      return -1;
    }
    default: {
      return 0;
    }
  }
}

uint8_t plc4c_s7_read_write_transport_size_get_data_protocol_id(plc4c_s7_read_write_transport_size value) {
  switch(value) {
    case 0: { /* '0x00' */
      return 14;
    }
    case 1: { /* '0x01' */
      return 01;
    }
    case 2: { /* '0x02' */
      return 11;
    }
    case 3: { /* '0x03' */
      return 41;
    }
    case 4: { /* '0x04' */
      return 12;
    }
    case 5: { /* '0x05' */
      return 23;
    }
    case 6: { /* '0x06' */
      return 13;
    }
    case 7: { /* '0x07' */
      return 25;
    }
    case 8: { /* '0x08' */
      return 31;
    }
    case 11: { /* '0x0B' */
      return 51;
    }
    case 19: { /* '0x13' */
      return 42;
    }
    default: {
      return 0;
    }
  }
}

#ifdef __cplusplus
}
#endif
#endif  // PLC4C_S7_READ_WRITE_TRANSPORT_SIZE_H_
