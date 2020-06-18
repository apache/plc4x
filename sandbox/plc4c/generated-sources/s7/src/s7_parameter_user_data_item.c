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
#include "s7_parameter_user_data_item.h"

// Array of discriminator values that match the enum type constants.
// (The order is identical to the enum constants so we can use the
// enum constant to directly access a given types discriminator values)
const plc4c_s7_read_write_s7_parameter_user_data_item_discriminator plc4c_s7_read_write_s7_parameter_user_data_item_discriminators[] = {
  {/* s7_read_write_s7_parameter_user_data_item_cpu_functions */
   .itemType = 0x12}
};

// Function returning the discriminator values for a given type constant.
plc4c_s7_read_write_s7_parameter_user_data_item_discriminator plc4c_s7_read_write_s7_parameter_user_data_item_get_discriminator(plc4c_s7_read_write_s7_parameter_user_data_item_type type) {
  return plc4c_s7_read_write_s7_parameter_user_data_item_discriminators[type];
}

// Parse function.
plc4c_return_code plc4c_s7_read_write_s7_parameter_user_data_item_parse(plc4c_spi_read_buffer* buf, plc4c_s7_read_write_s7_parameter_user_data_item** _message) {
  uint16_t startPos = plc4c_spi_read_get_pos(buf);
  uint16_t curPos;

  // Pointer to the parsed data structure.
  (*_message) = malloc(sizeof(plc4c_s7_read_write_s7_parameter_user_data_item));
  if(*_message == NULL) {
    return NO_MEMORY;
  }

  // Discriminator Field (itemType) (Used as input to a switch field)
  uint8_t itemType = plc4c_spi_read_unsigned_short(buf, 8);

  // Switch Field (Depending on the discriminator values, passes the instantiation to a sub-type)
  if(itemType == 0x12) { /* S7ParameterUserDataItemCPUFunctions */
                    
    // Implicit Field (itemLength) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
    uint8_t itemLength = plc4c_spi_read_unsigned_short(buf, 8);


                    
    // Simple Field (method)
    uint8_t method = plc4c_spi_read_unsigned_short(buf, 8);
    (*_message)->s7_parameter_user_data_item_cpu_functions_method = method;


                    
    // Simple Field (cpuFunctionType)
    unsigned int cpuFunctionType = plc4c_spi_read_unsigned_byte(buf, 4);
    (*_message)->s7_parameter_user_data_item_cpu_functions_cpu_function_type = cpuFunctionType;


                    
    // Simple Field (cpuFunctionGroup)
    unsigned int cpuFunctionGroup = plc4c_spi_read_unsigned_byte(buf, 4);
    (*_message)->s7_parameter_user_data_item_cpu_functions_cpu_function_group = cpuFunctionGroup;


                    
    // Simple Field (cpuSubfunction)
    uint8_t cpuSubfunction = plc4c_spi_read_unsigned_short(buf, 8);
    (*_message)->s7_parameter_user_data_item_cpu_functions_cpu_subfunction = cpuSubfunction;


                    
    // Simple Field (sequenceNumber)
    uint8_t sequenceNumber = plc4c_spi_read_unsigned_short(buf, 8);
    (*_message)->s7_parameter_user_data_item_cpu_functions_sequence_number = sequenceNumber;


                    
    // Optional Field (dataUnitReferenceNumber) (Can be skipped, if a given expression evaluates to false)
    uint8_t* dataUnitReferenceNumber = NULL;
    if((cpuFunctionType) == (8)) {
      dataUnitReferenceNumber = malloc(sizeof(uint8_t));
      if(dataUnitReferenceNumber == NULL) {
        return NO_MEMORY;
      }
      *dataUnitReferenceNumber = plc4c_spi_read_unsigned_short(buf, 8);
      (*_message)->s7_parameter_user_data_item_cpu_functions_data_unit_reference_number = dataUnitReferenceNumber;
    }


                    
    // Optional Field (lastDataUnit) (Can be skipped, if a given expression evaluates to false)
    uint8_t* lastDataUnit = NULL;
    if((cpuFunctionType) == (8)) {
      lastDataUnit = malloc(sizeof(uint8_t));
      if(lastDataUnit == NULL) {
        return NO_MEMORY;
      }
      *lastDataUnit = plc4c_spi_read_unsigned_short(buf, 8);
      (*_message)->s7_parameter_user_data_item_cpu_functions_last_data_unit = lastDataUnit;
    }


                    
    // Optional Field (errorCode) (Can be skipped, if a given expression evaluates to false)
    uint16_t* errorCode = NULL;
    if((cpuFunctionType) == (8)) {
      errorCode = malloc(sizeof(uint16_t));
      if(errorCode == NULL) {
        return NO_MEMORY;
      }
      *errorCode = plc4c_spi_read_unsigned_int(buf, 16);
      (*_message)->s7_parameter_user_data_item_cpu_functions_error_code = errorCode;
    }

  }

  return OK;
}

plc4c_return_code plc4c_s7_read_write_s7_parameter_user_data_item_serialize(plc4c_spi_write_buffer* buf, plc4c_s7_read_write_s7_parameter_user_data_item* _message) {

  // Discriminator Field (itemType)
  plc4c_spi_write_unsigned_short(buf, 8, plc4c_s7_read_write_s7_parameter_user_data_item_get_discriminator(_message->_type).itemType);

  return OK;
}
