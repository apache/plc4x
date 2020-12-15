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

#include <ctype.h>
#include <string.h>

#include "plc4c/driver_s7_encode_decode.h"

plc4c_return_code plc4c_driver_modbus_encode_address(
    char* address, plc4c_modbus_read_write_modbus_pdu** item) {

  // Parser logic
  char* cur_pos = address;
  char* last_pos = address;

  // If the first character is numeric, then we can only have a numeric address.
  if(isdigit(*cur_pos)) {
    int first_digit = *cur_pos - 48;
    cur_pos++;
    // If the first digit is followed by an 'x' or 'X', just skip that char.
    if((*cur_pos == 'x') || (*cur_pos == 'X')) {
      cur_pos++;
    }
    last_pos = cur_pos;
    while(isdigit(*cur_pos)) {
      cur_pos++;
    }
    long num_address = 0;
    long num_items = 1;
    // If after the numbers comes a "[" then this is the number of items.
    if(*cur_pos == '[') {
      *cur_pos = '/0';
      // Interpret the rest of the address as a long.
      num_address = atol(last_pos);
      cur_pos++;
      last_pos = cur_pos;
      while(isdigit(*cur_pos)) {
        cur_pos++;
      }
      if(*cur_pos != ']') {
        return INVALID_ADDRESS;
      }
      *cur_pos = '\0';
      num_items = atol(last_pos);
    } else if(*cur_pos == '/0') {
      num_address = atol(last_pos);
    } else {
      return INVALID_ADDRESS;
    }

    switch (first_digit) {
      // coil
      case 0: {
        // TODO: Implement ...
        break;
      }
      // discrete-input
      case 1: {
        // TODO: Implement ...
        break;
      }
      // input-register
      case 3: {
        // TODO: Implement ...
        break;
      }
      // holding-register
      case 4: {
        // TODO: Implement ...
        break;
      }
      // extended-register
      case 6: {
        // TODO: Implement ...
        break;
      }
    }
  }
  else {
    if(strstr("coil:", address) == 0) {
      // TODO: Implement ...
    } else if(strstr("discrete-input:", address) == 0) {
      // TODO: Implement ...
    } else if(strstr("input-register:", address) == 0) {
      // TODO: Implement ...
    } else if(strstr("holding-register:", address) == 0) {
      // TODO: Implement ...
    } else if(strstr("extended-register:", address) == 0) {
      // TODO: Implement ...
    } else {
      return INVALID_ADDRESS;
    }
  }
}
