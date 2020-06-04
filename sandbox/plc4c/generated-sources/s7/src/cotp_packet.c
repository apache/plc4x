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
#include "cotp_packet_data.h"
#include "cotp_packet_connection_request.h"
#include "cotp_packet_connection_response.h"
#include "cotp_packet_disconnect_request.h"
#include "cotp_packet_disconnect_response.h"
#include "cotp_packet_tpdu_error.h"
#include "cotp_packet.h"

// Array of discriminator values that match the enum type constants.
// (The order is identical to the enum constants so we can use the
// enum constant to directly access a given types discriminator values)
const plc4c_s7_read_write_cotp_packet_discriminator plc4c_s7_read_write_cotp_packet_discriminators[] = {
  {/* s7_read_write_cotp_packet_connection_request */
   .tpduCode = 0xE0},
  {/* s7_read_write_cotp_packet_connection_response */
   .tpduCode = 0xD0},
  {/* s7_read_write_cotp_packet_data */
   .tpduCode = 0xF0},
  {/* s7_read_write_cotp_packet_disconnect_request */
   .tpduCode = 0x80},
  {/* s7_read_write_cotp_packet_disconnect_response */
   .tpduCode = 0xC0},
  {/* s7_read_write_cotp_packet_tpdu_error */
   .tpduCode = 0x70}
};

// Function returning the discriminator values for a given type constant.
plc4c_s7_read_write_cotp_packet_discriminator plc4c_s7_read_write_cotp_packet_get_discriminator(plc4c_s7_read_write_cotp_packet_type type) {
  return plc4c_s7_read_write_cotp_packet_discriminators[type];
}

// Parse function.
plc4c_return_code plc4c_s7_read_write_cotp_packet_parse(plc4c_spi_read_buffer* buf, uint16_t cotpLen, plc4c_s7_read_write_cotp_packet** message) {
  uint16_t startPos = plc4c_spi_read_get_pos(buf);
  uint16_t curPos;

  // Pointer to the parsed datastructure.
  void* msg = NULL;
  // Factory function that allows filling the properties of this type
  void (*factory_ptr)()

  // Implicit Field (headerLength) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
  uint8_t headerLength = plc4c_spi_read_unsigned_short(buf, 8);

  // Discriminator Field (tpduCode) (Used as input to a switch field)
  uint8_t tpduCode = plc4c_spi_read_unsigned_short(buf, 8);

  // Switch Field (Depending on the discriminator values, passes the instantiation to a sub-type)
  if(tpduCode == 0xF0) {
    plc4c_s7_read_write_cotp_packet_data_parse(buf, cotpLen, &msg);
  } else 
  if(tpduCode == 0xE0) {
    plc4c_s7_read_write_cotp_packet_connection_request_parse(buf, cotpLen, &msg);
  } else 
  if(tpduCode == 0xD0) {
    plc4c_s7_read_write_cotp_packet_connection_response_parse(buf, cotpLen, &msg);
  } else 
  if(tpduCode == 0x80) {
    plc4c_s7_read_write_cotp_packet_disconnect_request_parse(buf, cotpLen, &msg);
  } else 
  if(tpduCode == 0xC0) {
    plc4c_s7_read_write_cotp_packet_disconnect_response_parse(buf, cotpLen, &msg);
  } else 
  if(tpduCode == 0x70) {
    plc4c_s7_read_write_cotp_packet_tpdu_error_parse(buf, cotpLen, &msg);
  }

  // Optional Field (payload) (Can be skipped, if a given expression evaluates to false)
  curPos = plc4c_spi_read_get_pos(buf) - startPos;
  plc4c_s7_read_write_s7_message* payload = NULL;
  if((curPos) < (cotpLen)) {
    plc4c_s7_read_write_s7_message_parse(buf, &payload);
  }

  return OK;
}

plc4c_return_code plc4c_s7_read_write_cotp_packet_serialize(plc4c_spi_write_buffer* buf, plc4c_s7_read_write_cotp_packet* message) {
  return OK;
}
