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

#include "s7_parameter_write_var_request.h"

plc4c_return_code plc4c_s7_read_write_s7_parameter_write_var_request_parse(plc4c_spi_read_buffer* buf, uint8_t messageType, plc4c_s7_read_write_s7_parameter_write_var_request** message) {
  uint16_t startPos = plc4c_spi_read_get_pos(buf);
  uint16_t curPos;

  plc4c_s7_read_write_s7_parameter_write_var_request* msg = malloc(sizeof(plc4c_s7_read_write_s7_parameter_write_var_request));

  // Implicit Field (numItems) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
  uint8_t numItems = plc4c_spi_read_unsigned_short(buf, 8);

  return OK;
}

plc4c_return_code plc4c_s7_read_write_s7_parameter_write_var_request_serialize(plc4c_spi_write_buffer* buf, plc4c_s7_read_write_s7_parameter_write_var_request* message) {
  return OK;
}
