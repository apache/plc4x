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

#include <stdio.h>
#include <plc4c/spi/read_buffer.h>
#include <plc4c/spi/write_buffer.h>
#include <plc4c/spi/evaluation_helper.h>
#include "s7_parameter.h"

// Array of discriminator values that match the enum type constants.
// (The order is identical to the enum constants so we can use the
// enum constant to directly access a given types discriminator values)
const plc4c_s7_read_write_s7_parameter_discriminator plc4c_s7_read_write_s7_parameter_discriminators[] = {
  {/* s7_read_write_s7_parameter_read_var_request */
   .parameterType = 0x04, .messageType = 0x01},
  {/* s7_read_write_s7_parameter_read_var_response */
   .parameterType = 0x04, .messageType = 0x03},
  {/* s7_read_write_s7_parameter_setup_communication */
   .parameterType = 0xF0, .messageType = -1},
  {/* s7_read_write_s7_parameter_user_data */
   .parameterType = 0x00, .messageType = 0x07},
  {/* s7_read_write_s7_parameter_write_var_request */
   .parameterType = 0x05, .messageType = 0x01},
  {/* s7_read_write_s7_parameter_write_var_response */
   .parameterType = 0x05, .messageType = 0x03}
};

// Function returning the discriminator values for a given type constant.
plc4c_s7_read_write_s7_parameter_discriminator plc4c_s7_read_write_s7_parameter_get_discriminator(plc4c_s7_read_write_s7_parameter_type type) {
  return plc4c_s7_read_write_s7_parameter_discriminators[type];
}

// Parse function.
plc4c_return_code plc4c_s7_read_write_s7_parameter_parse(plc4c_spi_read_buffer* buf, uint8_t messageType, plc4c_s7_read_write_s7_parameter** _message) {
  uint16_t startPos = plc4c_spi_read_get_pos(buf);
  uint16_t curPos;

  // Pointer to the parsed data structure.
  (*_message) = malloc(sizeof(plc4c_s7_read_write_s7_parameter));
  if(*_message == NULL) {
    return NO_MEMORY;
  }

  // Discriminator Field (parameterType) (Used as input to a switch field)
  uint8_t parameterType = plc4c_spi_read_unsigned_short(buf, 8);

  // Switch Field (Depending on the discriminator values, passes the instantiation to a sub-type)
  if(parameterType == 0xF0) { /* S7ParameterSetupCommunication */

  // Reserved Field (Compartmentalized so the "reserved" variable can't leak)
  {
    uint8_t _reserved = plc4c_spi_read_unsigned_short(buf, 8);
    if(_reserved != 0x00) {
      printf("Expected constant value '%d' but got '%d' for reserved field.", 0x00, _reserved);
    }
  }


  // Simple Field (maxAmqCaller)
  uint16_t maxAmqCaller = plc4c_spi_read_unsigned_int(buf, 16);
  (*_message)->s7_parameter_setup_communication_max_amq_caller = maxAmqCaller;


  // Simple Field (maxAmqCallee)
  uint16_t maxAmqCallee = plc4c_spi_read_unsigned_int(buf, 16);
  (*_message)->s7_parameter_setup_communication_max_amq_callee = maxAmqCallee;


  // Simple Field (pduLength)
  uint16_t pduLength = plc4c_spi_read_unsigned_int(buf, 16);
  (*_message)->s7_parameter_setup_communication_pdu_length = pduLength;
  } else 
  if((parameterType == 0x04) && (messageType == 0x01)) { /* S7ParameterReadVarRequest */

  // Implicit Field (numItems) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
  uint8_t numItems = plc4c_spi_read_unsigned_short(buf, 8);


  // Array field (items)
  plc4c_list* items = malloc(sizeof(plc4c_list));
  if(items == NULL) {
    return NO_MEMORY;
  }
  {
    // Count array
    uint8_t itemCount = numItems;
    for(int curItem = 0; curItem < itemCount; curItem++) {
      bool lastItem = curItem == (itemCount - 1);
      plc4c_list* _value = NULL;
      plc4c_return_code _res = plc4c_s7_read_write_s7_var_request_parameter_item_parse(buf, (void*) &_value);
      if(_res != OK) {
        return _res;
      }
      plc4c_utils_list_insert_head_value(items, _value);
    }
  }
  (*_message)->s7_parameter_read_var_request_items = items;
  } else 
  if((parameterType == 0x04) && (messageType == 0x03)) { /* S7ParameterReadVarResponse */

  // Simple Field (numItems)
  uint8_t numItems = plc4c_spi_read_unsigned_short(buf, 8);
  (*_message)->s7_parameter_read_var_response_num_items = numItems;
  } else 
  if((parameterType == 0x05) && (messageType == 0x01)) { /* S7ParameterWriteVarRequest */

  // Implicit Field (numItems) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
  uint8_t numItems = plc4c_spi_read_unsigned_short(buf, 8);


  // Array field (items)
  plc4c_list* items = malloc(sizeof(plc4c_list));
  if(items == NULL) {
    return NO_MEMORY;
  }
  {
    // Count array
    uint8_t itemCount = numItems;
    for(int curItem = 0; curItem < itemCount; curItem++) {
      bool lastItem = curItem == (itemCount - 1);
      plc4c_list* _value = NULL;
      plc4c_return_code _res = plc4c_s7_read_write_s7_var_request_parameter_item_parse(buf, (void*) &_value);
      if(_res != OK) {
        return _res;
      }
      plc4c_utils_list_insert_head_value(items, _value);
    }
  }
  (*_message)->s7_parameter_write_var_request_items = items;
  } else 
  if((parameterType == 0x05) && (messageType == 0x03)) { /* S7ParameterWriteVarResponse */

  // Simple Field (numItems)
  uint8_t numItems = plc4c_spi_read_unsigned_short(buf, 8);
  (*_message)->s7_parameter_write_var_response_num_items = numItems;
  } else 
  if((parameterType == 0x00) && (messageType == 0x07)) { /* S7ParameterUserData */

  // Implicit Field (numItems) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
  uint8_t numItems = plc4c_spi_read_unsigned_short(buf, 8);


  // Array field (items)
  plc4c_list* items = malloc(sizeof(plc4c_list));
  if(items == NULL) {
    return NO_MEMORY;
  }
  {
    // Count array
    uint8_t itemCount = numItems;
    for(int curItem = 0; curItem < itemCount; curItem++) {
      bool lastItem = curItem == (itemCount - 1);
      plc4c_list* _value = NULL;
      plc4c_return_code _res = plc4c_s7_read_write_s7_parameter_user_data_item_parse(buf, (void*) &_value);
      if(_res != OK) {
        return _res;
      }
      plc4c_utils_list_insert_head_value(items, _value);
    }
  }
  (*_message)->s7_parameter_user_data_items = items;
  }

  return OK;
}

plc4c_return_code plc4c_s7_read_write_s7_parameter_serialize(plc4c_spi_write_buffer* buf, plc4c_s7_read_write_s7_parameter* message) {
  return OK;
}
