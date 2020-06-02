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

#include "s7_message.h"

plc4c_return_code plc4c_s7_read_write_s7_message_parse(plc4c_read_buffer buf, plc4c_s7_read_write_s7_message** message) {
  uint16_t start_pos = plc4c_spi_read_get_pos(buf);
  uint16_t cur_pos;

  plc4c_s7_read_write_s7_message* msg = malloc(sizeof(plc4c_s7_read_write_s7_message));

  // Const Field (protocolId)
  uint8_t protocolId = plc4c_spi_read_unsigned_short(buf, 8);
  if(protocolId != S7Message.PROTOCOLID) {
    throw new ParseException("Expected constant value " + S7Message.PROTOCOLID + " but got " + protocolId);
  }

  // Discriminator Field (messageType) (Used as input to a switch field)
  uint8_t messageType = plc4c_spi_read_unsigned_short(buf, 8);

  // Reserved Field (Compartmentalized so the "reserved" variable can't leak)
  {
    uint16_t reserved = plc4c_spi_read_unsigned_int(buf, 16);
    if(reserved != (uint16_t) 0x0000) {
      LOGGER.info("Expected constant value " + 0x0000 + " but got " + reserved + " for reserved field.");
    }
  }

  // Simple Field (tpduReference)
  uint16_t tpduReference = plc4c_spi_read_unsigned_int(buf, 16);
  msg.tpdu_reference = tpduReference;

  // Implicit Field (parameterLength) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
  uint16_t parameterLength = plc4c_spi_read_unsigned_int(buf, 16);

  // Implicit Field (payloadLength) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
  uint16_t payloadLength = plc4c_spi_read_unsigned_int(buf, 16);

  // Switch Field (Depending on the discriminator values, passes the instantiation to a sub-type)
  if(plc4c_spi_evaluation_helper_equals(messageType, 0x01)) {
    plc4c_s7_read_write_s7_message_request_parse(buf, msg);
  } else 
  if(plc4c_spi_evaluation_helper_equals(messageType, 0x02)) {
    plc4c_s7_read_write_s7_message_response_parse(buf, msg);
  } else 
  if(plc4c_spi_evaluation_helper_equals(messageType, 0x03)) {
    plc4c_s7_read_write_s7_message_response_data_parse(buf, msg);
  } else 
  if(plc4c_spi_evaluation_helper_equals(messageType, 0x07)) {
    plc4c_s7_read_write_s7_message_user_data_parse(buf, msg);
  }

  // Optional Field (parameter) (Can be skipped, if a given expression evaluates to false)
  plc4c_s7_read_write_s7_parameter parameter = NULL;
  if((parameterLength) > (0)) {
    parameter = plc4c_s7_read_write_s7_parameter_parse(buf, messageType);
  }

  // Optional Field (payload) (Can be skipped, if a given expression evaluates to false)
  plc4c_s7_read_write_s7_payload payload = NULL;
  if((payloadLength) > (0)) {
    payload = plc4c_s7_read_write_s7_payload_parse(buf, messageType, parameter);
  }

  return OK;
}

plc4c_return_code plc4c_s7_read_write_s7_message_serialize(plc4c_write_buffer buf, plc4c_s7_read_write_s7_message* message) {
  return OK;
}
