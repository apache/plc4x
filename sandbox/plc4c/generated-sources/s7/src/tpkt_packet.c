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

#include "tpkt_packet.h"

plc4c_return_code plc4c_s7_read_write_tpkt_packet_parse(plc4c_read_buffer buf, plc4c_s7_read_write_tpkt_packet** message) {
  uint16_t start_pos = plc4c_spi_read_get_pos(buf);
  uint16_t cur_pos;

  plc4c_s7_read_write_tpkt_packet* msg = malloc(sizeof(plc4c_s7_read_write_tpkt_packet));

  // Const Field (protocolId)
  uint8_t protocolId = plc4c_spi_read_unsigned_short(buf, 8);
  if(protocolId != TPKTPacket.PROTOCOLID) {
    throw new ParseException("Expected constant value " + TPKTPacket.PROTOCOLID + " but got " + protocolId);
  }

  // Reserved Field (Compartmentalized so the "reserved" variable can't leak)
  {
    uint8_t reserved = plc4c_spi_read_unsigned_short(buf, 8);
    if(reserved != (uint8_t) 0x00) {
      LOGGER.info("Expected constant value " + 0x00 + " but got " + reserved + " for reserved field.");
    }
  }

  // Implicit Field (len) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
  uint16_t len = plc4c_spi_read_unsigned_int(buf, 16);

  // Simple Field (payload)
  plc4c_s7_read_write_cotp_packet payload = plc4c_s7_read_write_cotp_packet_parse(buf, (len) - (4));
  msg.payload = payload;

  return OK;
}

plc4c_return_code plc4c_s7_read_write_tpkt_packet_serialize(plc4c_write_buffer buf, plc4c_s7_read_write_tpkt_packet* message) {
  return OK;
}
