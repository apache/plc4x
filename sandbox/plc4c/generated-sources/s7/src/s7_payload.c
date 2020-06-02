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

#include <plc4c/spi/read_buffer.h>
#include <plc4c/spi/write_buffer.h>
#include <plc4c/spi/evaluation_helper.h>

#include "s7_payload.h"

plc4c_return_code plc4c_s7_read_write_s7_payload_parse(plc4c_read_buffer buf, uint8_t messageType, plc4c_s7_read_write_s7_parameter parameter, plc4c_s7_read_write_s7_payload** message) {
  uint16_t start_pos = plc4c_spi_read_get_pos(buf);
  uint16_t cur_pos;

  plc4c_s7_read_write_s7_payload* msg = malloc(sizeof(plc4c_s7_read_write_s7_payload));

  // Switch Field (Depending on the discriminator values, passes the instantiation to a sub-type)
  if(plc4c_spi_evaluation_helper_equals(parameter.getDiscriminatorValues()[0], 0x04) && plc4c_spi_evaluation_helper_equals(messageType, 0x03)) {
    plc4c_s7_read_write_s7_payload_read_var_response_parse(buf, msg, messageType, parameter);
  } else 
  if(plc4c_spi_evaluation_helper_equals(parameter.getDiscriminatorValues()[0], 0x05) && plc4c_spi_evaluation_helper_equals(messageType, 0x01)) {
    plc4c_s7_read_write_s7_payload_write_var_request_parse(buf, msg, messageType, parameter);
  } else 
  if(plc4c_spi_evaluation_helper_equals(parameter.getDiscriminatorValues()[0], 0x05) && plc4c_spi_evaluation_helper_equals(messageType, 0x03)) {
    plc4c_s7_read_write_s7_payload_write_var_response_parse(buf, msg, messageType, parameter);
  } else 
  if(plc4c_spi_evaluation_helper_equals(parameter.getDiscriminatorValues()[0], 0x00) && plc4c_spi_evaluation_helper_equals(messageType, 0x07)) {
    plc4c_s7_read_write_s7_payload_user_data_parse(buf, msg, messageType, parameter);
  }

  return OK;
}

plc4c_return_code plc4c_s7_read_write_s7_payload_serialize(plc4c_write_buffer buf, plc4c_s7_read_write_s7_payload* message) {
  return OK;
}
