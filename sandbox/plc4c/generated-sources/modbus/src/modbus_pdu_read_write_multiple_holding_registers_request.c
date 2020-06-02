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

#include "modbus_pdu_read_write_multiple_holding_registers_request.h"

plc4c_return_code plc4c_modbus_read_write_modbus_pdu_read_write_multiple_holding_registers_request_parse(plc4c_read_buffer buf, bool response, plc4c_modbus_read_write_modbus_pdu_read_write_multiple_holding_registers_request** message) {
  uint16_t start_pos = plc4c_spi_read_get_pos(buf);
  uint16_t cur_pos;

  plc4c_modbus_read_write_modbus_pdu_read_write_multiple_holding_registers_request* msg = malloc(sizeof(plc4c_modbus_read_write_modbus_pdu_read_write_multiple_holding_registers_request));

  // Simple Field (readStartingAddress)
  uint16_t readStartingAddress = plc4c_spi_read_unsigned_int(buf, 16);
  msg.read_starting_address = readStartingAddress;

  // Simple Field (readQuantity)
  uint16_t readQuantity = plc4c_spi_read_unsigned_int(buf, 16);
  msg.read_quantity = readQuantity;

  // Simple Field (writeStartingAddress)
  uint16_t writeStartingAddress = plc4c_spi_read_unsigned_int(buf, 16);
  msg.write_starting_address = writeStartingAddress;

  // Simple Field (writeQuantity)
  uint16_t writeQuantity = plc4c_spi_read_unsigned_int(buf, 16);
  msg.write_quantity = writeQuantity;

  // Implicit Field (byteCount) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
  uint8_t byteCount = plc4c_spi_read_unsigned_short(buf, 8);

  return OK;
}

plc4c_return_code plc4c_modbus_read_write_modbus_pdu_read_write_multiple_holding_registers_request_serialize(plc4c_write_buffer buf, plc4c_modbus_read_write_modbus_pdu_read_write_multiple_holding_registers_request* message) {
  return OK;
}
