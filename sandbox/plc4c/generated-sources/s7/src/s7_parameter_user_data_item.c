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
plc4c_return_code plc4c_s7_read_write_s7_parameter_user_data_item_parse(plc4c_spi_read_buffer* buf, plc4c_s7_read_write_s7_parameter_user_data_item** message) {
  uint16_t startPos = plc4c_spi_read_get_pos(buf);
  uint16_t curPos;

  // Pointer to the parsed data structure.
  plc4c_s7_read_write_s7_parameter_user_data_item* msg = malloc(sizeof(plc4c_s7_read_write_s7_parameter_user_data_item));

  // Discriminator Field (itemType) (Used as input to a switch field)
  uint8_t itemType = plc4c_spi_read_unsigned_short(buf, 8);

  // Switch Field (Depending on the discriminator values, passes the instantiation to a sub-type)
  if(itemType == 0x12) { /* S7ParameterUserDataItemCPUFunctions */
    uint8_t method = -1;
    msg->s7_parameter_user_data_item_cpu_functions_method = method;

    unsigned int cpuFunctionType = -1;
    msg->s7_parameter_user_data_item_cpu_functions_cpu_function_type = cpuFunctionType;

    unsigned int cpuFunctionGroup = -1;
    msg->s7_parameter_user_data_item_cpu_functions_cpu_function_group = cpuFunctionGroup;

    uint8_t cpuSubfunction = -1;
    msg->s7_parameter_user_data_item_cpu_functions_cpu_subfunction = cpuSubfunction;

    uint8_t sequenceNumber = -1;
    msg->s7_parameter_user_data_item_cpu_functions_sequence_number = sequenceNumber;

    uint8_t dataUnitReferenceNumber = -1;
    msg->s7_parameter_user_data_item_cpu_functions_data_unit_reference_number = dataUnitReferenceNumber;

    uint8_t lastDataUnit = -1;
    msg->s7_parameter_user_data_item_cpu_functions_last_data_unit = lastDataUnit;

    uint16_t errorCode = -1;
    msg->s7_parameter_user_data_item_cpu_functions_error_code = errorCode;
  }

  return OK;
}

plc4c_return_code plc4c_s7_read_write_s7_parameter_user_data_item_serialize(plc4c_spi_write_buffer* buf, plc4c_s7_read_write_s7_parameter_user_data_item* message) {
  return OK;
}
