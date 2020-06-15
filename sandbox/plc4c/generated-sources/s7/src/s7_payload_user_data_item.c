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
#include "s7_payload_user_data_item.h"

// Array of discriminator values that match the enum type constants.
// (The order is identical to the enum constants so we can use the
// enum constant to directly access a given types discriminator values)
const plc4c_s7_read_write_s7_payload_user_data_item_discriminator plc4c_s7_read_write_s7_payload_user_data_item_discriminators[] = {
  {/* s7_read_write_s7_payload_user_data_item_cpu_function_read_szl_request */
   .cpuFunctionType = 0x04},
  {/* s7_read_write_s7_payload_user_data_item_cpu_function_read_szl_response */
   .cpuFunctionType = 0x08}
};

// Function returning the discriminator values for a given type constant.
plc4c_s7_read_write_s7_payload_user_data_item_discriminator plc4c_s7_read_write_s7_payload_user_data_item_get_discriminator(plc4c_s7_read_write_s7_payload_user_data_item_type type) {
  return plc4c_s7_read_write_s7_payload_user_data_item_discriminators[type];
}

// Parse function.
plc4c_return_code plc4c_s7_read_write_s7_payload_user_data_item_parse(plc4c_spi_read_buffer* buf, unsigned int cpuFunctionType, plc4c_s7_read_write_s7_payload_user_data_item** message) {
  uint16_t startPos = plc4c_spi_read_get_pos(buf);
  uint16_t curPos;

  // Pointer to the parsed data structure.
  plc4c_s7_read_write_s7_payload_user_data_item* msg = malloc(sizeof(plc4c_s7_read_write_s7_payload_user_data_item));


  // Enum field (returnCode)
  plc4c_s7_read_write_data_transport_error_code returnCode = plc4c_spi_read_byte(buf, 8);
  msg->return_code = returnCode;

  // Enum field (transportSize)
  plc4c_s7_read_write_data_transport_size transportSize = plc4c_spi_read_byte(buf, 8);
  msg->transport_size = transportSize;

  // Implicit Field (dataLength) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
  uint16_t dataLength = plc4c_spi_read_unsigned_int(buf, 16);

  // Simple Field (szlId)
  plc4c_s7_read_write_szl_id* szlId = NULL;
  plc4c_s7_read_write_szl_id_parse(buf, (void*) &szlId);
  msg->szl_id = szlId;

  // Simple Field (szlIndex)
  uint16_t szlIndex = plc4c_spi_read_unsigned_int(buf, 16);
  msg->szl_index = szlIndex;

  // Switch Field (Depending on the discriminator values, passes the instantiation to a sub-type)
  if(cpuFunctionType == 0x04) { /* S7PayloadUserDataItemCpuFunctionReadSzlRequest */
  } else 
  if(cpuFunctionType == 0x08) { /* S7PayloadUserDataItemCpuFunctionReadSzlResponse */
    plc4c_list* items;
    msg->s7_payload_user_data_item_cpu_function_read_szl_response_items = items;
  }

  return OK;
}

plc4c_return_code plc4c_s7_read_write_s7_payload_user_data_item_serialize(plc4c_spi_write_buffer* buf, plc4c_s7_read_write_s7_payload_user_data_item* message) {
  return OK;
}
