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

#include "s7_parameter_setup_communication.h"

plc4c_return_code plc4c_s7_read_write_s7_parameter_setup_communication_parse(plc4c_read_buffer buf, uint8_t messageType, plc4c_s7_read_write_s7_parameter_setup_communication** message) {
  uint16_t start_pos = plc4c_spi_read_get_pos(buf);
  uint16_t cur_pos;

  plc4c_s7_read_write_s7_parameter_setup_communication* msg = malloc(sizeof(plc4c_s7_read_write_s7_parameter_setup_communication));

  // Reserved Field (Compartmentalized so the "reserved" variable can't leak)
  {
    uint8_t reserved = plc4c_spi_read_unsigned_short(buf, 8);
    if(reserved != (uint8_t) 0x00) {
      LOGGER.info("Expected constant value " + 0x00 + " but got " + reserved + " for reserved field.");
    }
  }

  // Simple Field (maxAmqCaller)
  uint16_t maxAmqCaller = plc4c_spi_read_unsigned_int(buf, 16);
  msg.max_amq_caller = maxAmqCaller;

  // Simple Field (maxAmqCallee)
  uint16_t maxAmqCallee = plc4c_spi_read_unsigned_int(buf, 16);
  msg.max_amq_callee = maxAmqCallee;

  // Simple Field (pduLength)
  uint16_t pduLength = plc4c_spi_read_unsigned_int(buf, 16);
  msg.pdu_length = pduLength;

  return OK;
}

plc4c_return_code plc4c_s7_read_write_s7_parameter_setup_communication_serialize(plc4c_write_buffer buf, plc4c_s7_read_write_s7_parameter_setup_communication* message) {
  return OK;
}
