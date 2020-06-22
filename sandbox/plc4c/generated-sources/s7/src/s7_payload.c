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
#include "s7_payload.h"

// Array of discriminator values that match the enum type constants.
// (The order is identical to the enum constants so we can use the
// enum constant to directly access a given types discriminator values)
const plc4c_s7_read_write_s7_payload_discriminator plc4c_s7_read_write_s7_payload_discriminators[] = {
  {/* s7_read_write_s7_payload_read_var_response */
   .parameterParameterType = 0x04, .messageType = 0x03},
  {/* s7_read_write_s7_payload_user_data */
   .parameterParameterType = 0x00, .messageType = 0x07},
  {/* s7_read_write_s7_payload_write_var_request */
   .parameterParameterType = 0x05, .messageType = 0x01},
  {/* s7_read_write_s7_payload_write_var_response */
   .parameterParameterType = 0x05, .messageType = 0x03}
};

// Function returning the discriminator values for a given type constant.
plc4c_s7_read_write_s7_payload_discriminator plc4c_s7_read_write_s7_payload_get_discriminator(plc4c_s7_read_write_s7_payload_type type) {
  return plc4c_s7_read_write_s7_payload_discriminators[type];
}

// Parse function.
plc4c_return_code plc4c_s7_read_write_s7_payload_parse(plc4c_spi_read_buffer* buf, uint8_t messageType, plc4c_s7_read_write_s7_parameter* parameter, plc4c_s7_read_write_s7_payload** _message) {
  uint16_t startPos = plc4c_spi_read_get_pos(buf);
  uint16_t curPos;

  // Allocate enough memory to contain this data structure.
  (*_message) = malloc(sizeof(plc4c_s7_read_write_s7_payload));
  if(*_message == NULL) {
    return NO_MEMORY;
  }

  // Switch Field (Depending on the discriminator values, passes the instantiation to a sub-type)
  if((plc4c_s7_read_write_s7_parameter_get_discriminator(parameter->_type).parameterType == 0x04) && (messageType == 0x03)) { /* S7PayloadReadVarResponse */
                    
    // Array field (items)
    plc4c_list* items = malloc(sizeof(plc4c_list));
    if(items == NULL) {
      return NO_MEMORY;
    }
    {
      // Count array
      uint8_t itemCount = ((plc4c_s7_read_write_s7_parameter*) (parameter))->s7_parameter_read_var_response_num_items;
      for(int curItem = 0; curItem < itemCount; curItem++) {
        bool lastItem = curItem == (itemCount - 1);
                          plc4c_list* _value = NULL;
        plc4c_return_code _res = plc4c_s7_read_write_s7_var_payload_data_item_parse(buf, lastItem, (void*) &_value);
        if(_res != OK) {
          return _res;
        }
        plc4c_utils_list_insert_head_value(items, _value);
      }
    }
    (*_message)->s7_payload_read_var_response_items = items;

  } else 
  if((plc4c_s7_read_write_s7_parameter_get_discriminator(parameter->_type).parameterType == 0x05) && (messageType == 0x01)) { /* S7PayloadWriteVarRequest */
                    
    // Array field (items)
    plc4c_list* items = malloc(sizeof(plc4c_list));
    if(items == NULL) {
      return NO_MEMORY;
    }
    {
      // Count array
      uint8_t itemCount = plc4c_spi_evaluation_helper_count(((plc4c_s7_read_write_s7_parameter*) (parameter))->s7_parameter_write_var_request_items);
      for(int curItem = 0; curItem < itemCount; curItem++) {
        bool lastItem = curItem == (itemCount - 1);
                          plc4c_list* _value = NULL;
        plc4c_return_code _res = plc4c_s7_read_write_s7_var_payload_data_item_parse(buf, lastItem, (void*) &_value);
        if(_res != OK) {
          return _res;
        }
        plc4c_utils_list_insert_head_value(items, _value);
      }
    }
    (*_message)->s7_payload_write_var_request_items = items;

  } else 
  if((plc4c_s7_read_write_s7_parameter_get_discriminator(parameter->_type).parameterType == 0x05) && (messageType == 0x03)) { /* S7PayloadWriteVarResponse */
                    
    // Array field (items)
    plc4c_list* items = malloc(sizeof(plc4c_list));
    if(items == NULL) {
      return NO_MEMORY;
    }
    {
      // Count array
      uint8_t itemCount = ((plc4c_s7_read_write_s7_parameter*) (parameter))->s7_parameter_write_var_response_num_items;
      for(int curItem = 0; curItem < itemCount; curItem++) {
        bool lastItem = curItem == (itemCount - 1);
                          plc4c_list* _value = NULL;
        plc4c_return_code _res = plc4c_s7_read_write_s7_var_payload_status_item_parse(buf, (void*) &_value);
        if(_res != OK) {
          return _res;
        }
        plc4c_utils_list_insert_head_value(items, _value);
      }
    }
    (*_message)->s7_payload_write_var_response_items = items;

  } else 
  if((plc4c_s7_read_write_s7_parameter_get_discriminator(parameter->_type).parameterType == 0x00) && (messageType == 0x07)) { /* S7PayloadUserData */
                    
    // Array field (items)
    plc4c_list* items = malloc(sizeof(plc4c_list));
    if(items == NULL) {
      return NO_MEMORY;
    }
    {
      // Count array
      uint8_t itemCount = plc4c_spi_evaluation_helper_count(((plc4c_s7_read_write_s7_parameter*) (parameter))->s7_parameter_user_data_items);
      for(int curItem = 0; curItem < itemCount; curItem++) {
        bool lastItem = curItem == (itemCount - 1);
                          plc4c_list* _value = NULL;
        plc4c_return_code _res = plc4c_s7_read_write_s7_payload_user_data_item_parse(buf, ((plc4c_s7_read_write_s7_parameter_user_data_item*) (plc4c_utils_list_get(((plc4c_s7_read_write_s7_parameter*) (parameter))->s7_parameter_user_data_items, 0)))->s7_parameter_user_data_item_cpu_functions_cpu_function_type, (void*) &_value);
        if(_res != OK) {
          return _res;
        }
        plc4c_utils_list_insert_head_value(items, _value);
      }
    }
    (*_message)->s7_payload_user_data_items = items;

  }

  return OK;
}

plc4c_return_code plc4c_s7_read_write_s7_payload_serialize(plc4c_spi_write_buffer* buf, plc4c_s7_read_write_s7_payload* _message) {

  return OK;
}

uint8_t plc4c_s7_read_write_s7_payload_length_in_bytes(plc4c_s7_read_write_s7_payload* message) {
  return plc4c_s7_read_write_s7_payload_length_in_bits(message) / 8;
}

uint8_t plc4c_s7_read_write_s7_payload_length_in_bits(plc4c_s7_read_write_s7_payload* message) {
  return 0;
}

