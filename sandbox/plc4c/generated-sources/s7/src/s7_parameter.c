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
   .messageType = 0x01, .parameterType = 0x04},
  {/* s7_read_write_s7_parameter_read_var_response */
   .messageType = 0x03, .parameterType = 0x04},
  {/* s7_read_write_s7_parameter_setup_communication */
   .messageType = -1, .parameterType = 0xF0},
  {/* s7_read_write_s7_parameter_user_data */
   .messageType = 0x07, .parameterType = 0x00},
  {/* s7_read_write_s7_parameter_write_var_request */
   .messageType = 0x01, .parameterType = 0x05},
  {/* s7_read_write_s7_parameter_write_var_response */
   .messageType = 0x03, .parameterType = 0x05}
};

// Function returning the discriminator values for a given type constant.
plc4c_s7_read_write_s7_parameter_discriminator plc4c_s7_read_write_s7_parameter_get_discriminator(plc4c_s7_read_write_s7_parameter_type type) {
  return plc4c_s7_read_write_s7_parameter_discriminators[type];
}

// Parse function.
plc4c_return_code plc4c_s7_read_write_s7_parameter_parse(plc4c_spi_read_buffer* buf, uint8_t messageType, plc4c_s7_read_write_s7_parameter** message) {
  uint16_t startPos = plc4c_spi_read_get_pos(buf);
  uint16_t curPos;

  // Pointer to the parsed datastructure.
  plc4c_s7_read_write_s7_parameter* msg = malloc(sizeof(plc4c_s7_read_write_s7_parameter));

  // Discriminator Field (parameterType) (Used as input to a switch field)
  uint8_t parameterType = plc4c_spi_read_unsigned_short(buf, 8);

  // Switch Field (Depending on the discriminator values, passes the instantiation to a sub-type)
  if(parameterType == 0xF0) { /* S7ParameterSetupCommunication */
    uint16_t maxAmqCaller = -1;
    msg->s7_parameter_setup_communication_max_amq_caller = maxAmqCaller;

    uint16_t maxAmqCallee = -1;
    msg->s7_parameter_setup_communication_max_amq_callee = maxAmqCallee;

    uint16_t pduLength = -1;
    msg->s7_parameter_setup_communication_pdu_length = pduLength;
  } else 
  if((parameterType == 0x04) && (messageType == 0x01)) { /* S7ParameterReadVarRequest */
    plc4c_list* items;
    msg->s7_parameter_read_var_request_items = items;
  } else 
  if((parameterType == 0x04) && (messageType == 0x03)) { /* S7ParameterReadVarResponse */
    uint8_t numItems = -1;
    msg->s7_parameter_read_var_response_num_items = numItems;
  } else 
  if((parameterType == 0x05) && (messageType == 0x01)) { /* S7ParameterWriteVarRequest */
    plc4c_list* items;
    msg->s7_parameter_write_var_request_items = items;
  } else 
  if((parameterType == 0x05) && (messageType == 0x03)) { /* S7ParameterWriteVarResponse */
    uint8_t numItems = -1;
    msg->s7_parameter_write_var_response_num_items = numItems;
  } else 
  if((parameterType == 0x00) && (messageType == 0x07)) { /* S7ParameterUserData */
    plc4c_list* items;
    msg->s7_parameter_user_data_items = items;
  }

  return OK;
}

plc4c_return_code plc4c_s7_read_write_s7_parameter_serialize(plc4c_spi_write_buffer* buf, plc4c_s7_read_write_s7_parameter* message) {
  return OK;
}
