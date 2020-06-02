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

#include "cotp_packet.h"

plc4c_return_code plc4c_s7_read_write_cotp_packet_parse(plc4c_read_buffer buf, uint16_t cotpLen, plc4c_s7_read_write_cotp_packet** message) {
  uint16_t start_pos = plc4c_spi_read_get_pos(buf);
  uint16_t cur_pos;

  plc4c_s7_read_write_cotp_packet* msg = malloc(sizeof(plc4c_s7_read_write_cotp_packet));

  // Implicit Field (headerLength) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
  uint8_t headerLength = plc4c_spi_read_unsigned_short(buf, 8);

  // Discriminator Field (tpduCode) (Used as input to a switch field)
  uint8_t tpduCode = plc4c_spi_read_unsigned_short(buf, 8);

  // Switch Field (Depending on the discriminator values, passes the instantiation to a sub-type)
  if(plc4c_spi_evaluation_helper_equals(tpduCode, 0xF0)) {
    plc4c_s7_read_write_cotp_packet_data_parse(buf, msg, cotpLen);
  } else 
  if(plc4c_spi_evaluation_helper_equals(tpduCode, 0xE0)) {
    plc4c_s7_read_write_cotp_packet_connection_request_parse(buf, msg, cotpLen);
  } else 
  if(plc4c_spi_evaluation_helper_equals(tpduCode, 0xD0)) {
    plc4c_s7_read_write_cotp_packet_connection_response_parse(buf, msg, cotpLen);
  } else 
  if(plc4c_spi_evaluation_helper_equals(tpduCode, 0x80)) {
    plc4c_s7_read_write_cotp_packet_disconnect_request_parse(buf, msg, cotpLen);
  } else 
  if(plc4c_spi_evaluation_helper_equals(tpduCode, 0xC0)) {
    plc4c_s7_read_write_cotp_packet_disconnect_response_parse(buf, msg, cotpLen);
  } else 
  if(plc4c_spi_evaluation_helper_equals(tpduCode, 0x70)) {
    plc4c_s7_read_write_cotp_packet_tpdu_error_parse(buf, msg, cotpLen);
  }

  // Optional Field (payload) (Can be skipped, if a given expression evaluates to false)
  curPos = io.getPos() - startPos;
  plc4c_s7_read_write_s7_message payload = NULL;
  if((curPos) < (cotpLen)) {
    payload = plc4c_s7_read_write_s7_message_parse(buf);
  }

  return OK;
}

plc4c_return_code plc4c_s7_read_write_cotp_packet_serialize(plc4c_write_buffer buf, plc4c_s7_read_write_cotp_packet* message) {
  return OK;
}
