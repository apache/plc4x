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
#include "s7_var_payload_data_item.h"

// Parse function.
plc4c_return_code plc4c_s7_read_write_s7_var_payload_data_item_parse(plc4c_spi_read_buffer* buf, bool lastItem, plc4c_s7_read_write_s7_var_payload_data_item** message) {
  uint16_t startPos = plc4c_spi_read_get_pos(buf);
  uint16_t curPos;

  // Pointer to the parsed data structure.
  plc4c_s7_read_write_s7_var_payload_data_item* msg = malloc(sizeof(plc4c_s7_read_write_s7_var_payload_data_item));


  // Enum field (returnCode)
  plc4c_s7_read_write_data_transport_error_code returnCode = plc4c_spi_read_byte(buf, 8);
  msg->return_code = returnCode;

  // Enum field (transportSize)
  plc4c_s7_read_write_data_transport_size transportSize = plc4c_spi_read_byte(buf, 8);
  msg->transport_size = transportSize;

  // Simple Field (dataLength)
  uint16_t dataLength = plc4c_spi_read_unsigned_int(buf, 16);
  msg->data_length = dataLength;

  // Padding Field (pad)
  bool _padNeedsPadding = (bool) ((plc4c_spi_read_has_more(buf, 8)) && ((!(lastItem)) && (((((plc4c_utils_list_size(&data)) % (2))) == (1)))));
  if(_padNeedsPadding) {
    // Just read the padding data and ignore it
    plc4c_spi_read_unsigned_short(buf, 8);
  }

  return OK;
}

plc4c_return_code plc4c_s7_read_write_s7_var_payload_data_item_serialize(plc4c_spi_write_buffer* buf, plc4c_s7_read_write_s7_var_payload_data_item* message) {
  return OK;
}
