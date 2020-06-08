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
#include "cotp_parameter.h"

// Array of discriminator values that match the enum type constants.
// (The order is identical to the enum constants so we can use the
// enum constant to directly access a given types discriminator values)
const plc4c_s7_read_write_cotp_parameter_discriminator plc4c_s7_read_write_cotp_parameter_discriminators[] = {
  {/* s7_read_write_cotp_parameter_called_tsap */
   .parameterType = 0xC2},
  {/* s7_read_write_cotp_parameter_calling_tsap */
   .parameterType = 0xC1},
  {/* s7_read_write_cotp_parameter_checksum */
   .parameterType = 0xC3},
  {/* s7_read_write_cotp_parameter_disconnect_additional_information */
   .parameterType = 0xE0},
  {/* s7_read_write_cotp_parameter_tpdu_size */
   .parameterType = 0xC0}
};

// Function returning the discriminator values for a given type constant.
plc4c_s7_read_write_cotp_parameter_discriminator plc4c_s7_read_write_cotp_parameter_get_discriminator(plc4c_s7_read_write_cotp_parameter_type type) {
  return plc4c_s7_read_write_cotp_parameter_discriminators[type];
}

// Parse function.
plc4c_return_code plc4c_s7_read_write_cotp_parameter_parse(plc4c_spi_read_buffer* buf, uint8_t rest, plc4c_s7_read_write_cotp_parameter** message) {
  uint16_t startPos = plc4c_spi_read_get_pos(buf);
  uint16_t curPos;

  // Pointer to the parsed datastructure.
  plc4c_s7_read_write_cotp_parameter* msg = malloc(sizeof(plc4c_s7_read_write_cotp_parameter));

  // Discriminator Field (parameterType) (Used as input to a switch field)
  uint8_t parameterType = plc4c_spi_read_unsigned_short(buf, 8);

  // Implicit Field (parameterLength) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
  uint8_t parameterLength = plc4c_spi_read_unsigned_short(buf, 8);

  // Switch Field (Depending on the discriminator values, passes the instantiation to a sub-type)
  if(parameterType == 0xC0) { /* COTPParameterTpduSize */
    plc4c_s7_read_write_cotp_tpdu_size* tpduSize = NULL;
    msg->cotp_parameter_tpdu_size_tpdu_size = tpduSize;
  } else 
  if(parameterType == 0xC1) { /* COTPParameterCallingTsap */
    uint16_t tsapId = -1;
    msg->cotp_parameter_calling_tsap_tsap_id = tsapId;
  } else 
  if(parameterType == 0xC2) { /* COTPParameterCalledTsap */
    uint16_t tsapId = -1;
    msg->cotp_parameter_called_tsap_tsap_id = tsapId;
  } else 
  if(parameterType == 0xC3) { /* COTPParameterChecksum */
    uint8_t crc = -1;
    msg->cotp_parameter_checksum_crc = crc;
  } else 
  if(parameterType == 0xE0) { /* COTPParameterDisconnectAdditionalInformation */
    plc4c_list data;
    msg->cotp_parameter_disconnect_additional_information_data = data;
  }

  return OK;
}

plc4c_return_code plc4c_s7_read_write_cotp_parameter_serialize(plc4c_spi_write_buffer* buf, plc4c_s7_read_write_cotp_parameter* message) {
  return OK;
}
