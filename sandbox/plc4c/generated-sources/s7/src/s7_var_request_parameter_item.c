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
#include "s7_var_request_parameter_item.h"

// Array of discriminator values that match the enum type constants.
// (The order is identical to the enum constants so we can use the
// enum constant to directly access a given types discriminator values)
const plc4c_s7_read_write_s7_var_request_parameter_item_discriminator plc4c_s7_read_write_s7_var_request_parameter_item_discriminators[] = {
  {/* s7_read_write_s7_var_request_parameter_item_address */
   .itemType = 0x12}
};

// Function returning the discriminator values for a given type constant.
plc4c_s7_read_write_s7_var_request_parameter_item_discriminator plc4c_s7_read_write_s7_var_request_parameter_item_get_discriminator(plc4c_s7_read_write_s7_var_request_parameter_item_type type) {
  return plc4c_s7_read_write_s7_var_request_parameter_item_discriminators[type];
}

// Parse function.
plc4c_return_code plc4c_s7_read_write_s7_var_request_parameter_item_parse(plc4c_spi_read_buffer* buf, plc4c_s7_read_write_s7_var_request_parameter_item** _message) {
  uint16_t startPos = plc4c_spi_read_get_pos(buf);
  uint16_t curPos;

  // Pointer to the parsed data structure.
  (*_message) = malloc(sizeof(plc4c_s7_read_write_s7_var_request_parameter_item));

  // Discriminator Field (itemType) (Used as input to a switch field)
  uint8_t itemType = plc4c_spi_read_unsigned_short(buf, 8);

  // Switch Field (Depending on the discriminator values, passes the instantiation to a sub-type)
  if(itemType == 0x12) { /* S7VarRequestParameterItemAddress */

  // Implicit Field (itemLength) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
  uint8_t itemLength = plc4c_spi_read_unsigned_short(buf, 8);


  // Simple Field (address)
  plc4c_s7_read_write_s7_address address;
  plc4c_s7_read_write_s7_address_parse(buf, (void*) &address);
  (*_message)->s7_var_request_parameter_item_address_address = address;
  }

  return OK;
}

plc4c_return_code plc4c_s7_read_write_s7_var_request_parameter_item_serialize(plc4c_spi_write_buffer* buf, plc4c_s7_read_write_s7_var_request_parameter_item* message) {
  return OK;
}
