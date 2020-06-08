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
#include "modbus_serial_adu.h"

// Parse function.
plc4c_return_code plc4c_modbus_read_write_modbus_serial_adu_parse(plc4c_spi_read_buffer* buf, bool response, plc4c_modbus_read_write_modbus_serial_adu** message) {
  uint16_t startPos = plc4c_spi_read_get_pos(buf);
  uint16_t curPos;

  // Pointer to the parsed datastructure.
  plc4c_modbus_read_write_modbus_serial_adu* msg = malloc(sizeof(plc4c_modbus_read_write_modbus_serial_adu));


  // Simple Field (transactionId)
  uint16_t transactionId = plc4c_spi_read_unsigned_int(buf, 16);
  msg->transaction_id = transactionId;

  // Reserved Field (Compartmentalized so the "reserved" variable can't leak)
  {
    uint16_t reserved = plc4c_spi_read_unsigned_int(buf, 16);
    if(reserved != (uint16_t) 0x0000) {
      printf("Expected constant value '%d' but got '%d' for reserved field.", 0x0000, reserved);
    }
  }

  // Simple Field (length)
  uint16_t length = plc4c_spi_read_unsigned_int(buf, 16);
  msg->length = length;

  // Simple Field (address)
  uint8_t address = plc4c_spi_read_unsigned_short(buf, 8);
  msg->address = address;

  // Simple Field (pdu)
  plc4c_modbus_read_write_modbus_pdu* pdu = NULL;
  plc4c_modbus_read_write_modbus_pdu_parse(buf, response, (void*) &pdu);
  msg->pdu = pdu;

  return OK;
}

plc4c_return_code plc4c_modbus_read_write_modbus_serial_adu_serialize(plc4c_spi_write_buffer* buf, plc4c_modbus_read_write_modbus_serial_adu* message) {
  return OK;
}
