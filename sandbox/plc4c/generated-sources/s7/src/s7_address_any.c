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

#include "s7_address_any.h"

plc4c_return_code plc4c_s7_read_write_s7_address_any_parse(plc4c_read_buffer buf, plc4c_s7_read_write_s7_address_any** message) {
  uint16_t start_pos = plc4c_spi_read_get_pos(buf);
  uint16_t cur_pos;

  plc4c_s7_read_write_s7_address_any* msg = malloc(sizeof(plc4c_s7_read_write_s7_address_any));

  // Enum field (transportSize)
  plc4c_s7_read_write_transport_size transportSize = plc4c_s7_read_write_transport_size.valueOf(plc4c_spi_read_byte(buf, 8));

  // Simple Field (numberOfElements)
  uint16_t numberOfElements = plc4c_spi_read_unsigned_int(buf, 16);
  msg.number_of_elements = numberOfElements;

  // Simple Field (dbNumber)
  uint16_t dbNumber = plc4c_spi_read_unsigned_int(buf, 16);
  msg.db_number = dbNumber;

  // Enum field (area)
  plc4c_s7_read_write_memory_area area = plc4c_s7_read_write_memory_area.valueOf(plc4c_spi_read_byte(buf, 8));

  // Reserved Field (Compartmentalized so the "reserved" variable can't leak)
  {
    unsigned int reserved = plc4c_spi_read_unsigned_short(buf, 5);
    if(reserved != (unsigned int) 0x00) {
      LOGGER.info("Expected constant value " + 0x00 + " but got " + reserved + " for reserved field.");
    }
  }

  // Simple Field (byteAddress)
  uint16_t byteAddress = plc4c_spi_read_unsigned_int(buf, 16);
  msg.byte_address = byteAddress;

  // Simple Field (bitAddress)
  unsigned int bitAddress = plc4c_spi_read_unsigned_byte(buf, 3);
  msg.bit_address = bitAddress;

  return OK;
}

plc4c_return_code plc4c_s7_read_write_s7_address_any_serialize(plc4c_write_buffer buf, plc4c_s7_read_write_s7_address_any* message) {
  return OK;
}
