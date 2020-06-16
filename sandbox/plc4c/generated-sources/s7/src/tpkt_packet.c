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
#include "tpkt_packet.h"

// Parse function.
plc4c_return_code plc4c_s7_read_write_tpkt_packet_parse(plc4c_spi_read_buffer* buf, plc4c_s7_read_write_tpkt_packet** _message) {
  uint16_t startPos = plc4c_spi_read_get_pos(buf);
  uint16_t curPos;

  // Pointer to the parsed data structure.
  (*_message) = malloc(sizeof(plc4c_s7_read_write_tpkt_packet));


  // Const Field (protocolId)
  uint8_t protocolId = plc4c_spi_read_unsigned_short(buf, 8);
  if(protocolId != S7_READ_WRITE_TPKT_PACKET_PROTOCOL_ID) {
    return PARSE_ERROR;
    // throw new ParseException("Expected constant value " + S7_READ_WRITE_TPKT_PACKET_PROTOCOL_ID + " but got " + protocolId);
  }

  // Reserved Field (Compartmentalized so the "reserved" variable can't leak)
  {
    uint8_t _reserved = plc4c_spi_read_unsigned_short(buf, 8);
    if(_reserved != 0x00) {
      printf("Expected constant value '%d' but got '%d' for reserved field.", 0x00, _reserved);
    }
  }

  // Implicit Field (len) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
  uint16_t len = plc4c_spi_read_unsigned_int(buf, 16);

  // Simple Field (payload)
  plc4c_s7_read_write_cotp_packet payload;
  plc4c_s7_read_write_cotp_packet_parse(buf, (len) - (4), (void*) &payload);
  (*_message)->payload = payload;

  return OK;
}

plc4c_return_code plc4c_s7_read_write_tpkt_packet_serialize(plc4c_spi_write_buffer* buf, plc4c_s7_read_write_tpkt_packet* message) {
  return OK;
}
