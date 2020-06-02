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

#include "s7_var_request_parameter_item_address.h"

plc4c_return_code plc4c_s7_read_write_s7_var_request_parameter_item_address_parse(plc4c_read_buffer buf, plc4c_s7_read_write_s7_var_request_parameter_item_address** message) {
  uint16_t start_pos = plc4c_spi_read_get_pos(buf);
  uint16_t cur_pos;

  plc4c_s7_read_write_s7_var_request_parameter_item_address* msg = malloc(sizeof(plc4c_s7_read_write_s7_var_request_parameter_item_address));

  // Implicit Field (itemLength) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
  uint8_t itemLength = plc4c_spi_read_unsigned_short(buf, 8);

  // Simple Field (address)
  plc4c_s7_read_write_s7_address address = plc4c_s7_read_write_s7_address_parse(buf);
  msg.address = address;

  return OK;
}

plc4c_return_code plc4c_s7_read_write_s7_var_request_parameter_item_address_serialize(plc4c_write_buffer buf, plc4c_s7_read_write_s7_var_request_parameter_item_address* message) {
  return OK;
}
