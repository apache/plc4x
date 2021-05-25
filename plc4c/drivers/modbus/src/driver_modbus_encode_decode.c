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
#include <plc4c/driver_modbus.h>
#include <string.h>

#include "plc4c/driver_s7_encode_decode.h"

plc4c_return_code plc4c_driver_modbus_encode_address(char* address,
                                                     void** item) {
  plc4c_driver_modbus_item* modbus_item =
      malloc(sizeof(plc4c_driver_modbus_item));
  if (modbus_item == NULL) {
    return NO_MEMORY;
  }
  // The overall default is 1
  modbus_item->num_elements = 1;

  // Parser logic
  char* cur_pos = address;
  char* last_pos = address;

  // If the first character is numeric, then we can only have a numeric address.
  if (isdigit(*cur_pos)) {
    int first_digit = *cur_pos - 48;
    cur_pos++;
    // If the first digit is followed by an 'x' or 'X', just skip that char.
    if ((*cur_pos == 'x') || (*cur_pos == 'X')) {
      cur_pos++;
    }
    last_pos = cur_pos;
    // In case of a numeric address, the first digit defines the type of address
    switch (first_digit) {
      // coil
      case 0: {
        modbus_item->type = PLC4C_DRIVER_MODBUS_ADDRESS_TYPE_COIL;
        modbus_item->datatype = plc4c_modbus_read_write_modbus_data_type_BOOL;
        break;
      }
      // discrete-input
      case 1: {
        modbus_item->type = PLC4C_DRIVER_MODBUS_ADDRESS_TYPE_DISCRETE_INPUT;
        modbus_item->datatype = plc4c_modbus_read_write_modbus_data_type_UINT;
        break;
      }
      // input-register
      case 3: {
        modbus_item->type = PLC4C_DRIVER_MODBUS_ADDRESS_TYPE_INPUT_REGISTER;
        modbus_item->datatype = plc4c_modbus_read_write_modbus_data_type_UINT;
        break;
      }
      // holding-register
      case 4: {
        modbus_item->type = PLC4C_DRIVER_MODBUS_ADDRESS_TYPE_HOLDING_REGISTER;
        modbus_item->datatype = plc4c_modbus_read_write_modbus_data_type_UINT;
        break;
      }
      // extended-register
      case 6: {
        modbus_item->type = PLC4C_DRIVER_MODBUS_ADDRESS_TYPE_EXTENDED_REGISTER;
        modbus_item->datatype = plc4c_modbus_read_write_modbus_data_type_UINT;
        break;
      }
      default: {
        return INVALID_ADDRESS;
      }
    }
  }
  // If the first character isn't a digit, it must be a name of the field-type.
  else {
    if (((char*)strstr(address, "coil:")) != NULL) {
      modbus_item->type = PLC4C_DRIVER_MODBUS_ADDRESS_TYPE_COIL;
      modbus_item->datatype = plc4c_modbus_read_write_modbus_data_type_BOOL;
      cur_pos += 5;
    } else if (((char*)strstr(address, "discrete-register:")) != NULL) {
      modbus_item->type = PLC4C_DRIVER_MODBUS_ADDRESS_TYPE_DISCRETE_INPUT;
      modbus_item->datatype = plc4c_modbus_read_write_modbus_data_type_UINT;
      cur_pos += 18;
    } else if (((char*)strstr(address, "input-register:")) != NULL) {
      modbus_item->type = PLC4C_DRIVER_MODBUS_ADDRESS_TYPE_INPUT_REGISTER;
      modbus_item->datatype = plc4c_modbus_read_write_modbus_data_type_UINT;
      cur_pos += 15;
    } else if (((char*)strstr(address, "holding-register:")) != NULL) {
      modbus_item->type = PLC4C_DRIVER_MODBUS_ADDRESS_TYPE_HOLDING_REGISTER;
      modbus_item->datatype = plc4c_modbus_read_write_modbus_data_type_UINT;
      cur_pos += 17;
    } else if (((char*)strstr(address, "extended-register:")) != NULL) {
      modbus_item->type = PLC4C_DRIVER_MODBUS_ADDRESS_TYPE_EXTENDED_REGISTER;
      modbus_item->datatype = plc4c_modbus_read_write_modbus_data_type_UINT;
      cur_pos += 18;
    } else {
      return INVALID_ADDRESS;
    }
  }

  // Now consume all of the digits.
  last_pos = cur_pos;
  while (isdigit(*cur_pos)) {
    cur_pos++;
  }
  if (last_pos == cur_pos) {
    return INVALID_ADDRESS;
  }
  // Parse the current segment as number
  int len = cur_pos - last_pos;
  char* address_str = malloc(sizeof(char) * (len + 1));
  strncpy(address_str, last_pos, len);
  *(address_str + len) = '\0';
  modbus_item->address = (uint16_t)atol(address_str);

  // If a datatype is provided, parse that now
  if (*cur_pos == ':') {
    cur_pos++;

    // Inspect the substring, if this matches and of the supported datatypes
    last_pos = cur_pos;
    while ((*cur_pos != '\0') && (*cur_pos != '[')) {
      cur_pos++;
    }
    len = cur_pos - last_pos;
    char* datatype_str = malloc(sizeof(char) * (len + 1));
    strncpy(datatype_str, last_pos, len);
    *(datatype_str + len) = '\0';
    modbus_item->datatype =
        plc4c_modbus_read_write_modbus_data_type_value_of(datatype_str);

    // If a number of elements is provided, parse that now.
    if (*cur_pos == '[') {
      cur_pos++;
      last_pos = cur_pos;
      while (isdigit(*cur_pos)) {
        cur_pos++;
      }
      if (*cur_pos != ']') {
        return INVALID_ADDRESS;
      }
      *cur_pos = '\0';
      modbus_item->num_elements = atol(last_pos);
    }
  }

  // Pass back the result
  *item = modbus_item;
  return OK;
}
