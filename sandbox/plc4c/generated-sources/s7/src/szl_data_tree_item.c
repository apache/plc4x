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
#include "szl_data_tree_item.h"

// Parse function.
plc4c_return_code plc4c_s7_read_write_szl_data_tree_item_parse(plc4c_spi_read_buffer* buf, plc4c_s7_read_write_szl_data_tree_item** _message) {
  uint16_t startPos = plc4c_spi_read_get_pos(buf);
  uint16_t curPos;

  // Pointer to the parsed data structure.
  (*_message) = malloc(sizeof(plc4c_s7_read_write_szl_data_tree_item));


  // Simple Field (itemIndex)
  uint16_t itemIndex = plc4c_spi_read_unsigned_int(buf, 16);
  (*_message)->item_index = itemIndex;

  // Array field (mlfb)
  plc4c_list* mlfb = malloc(sizeof(plc4c_list));
  if(mlfb == NULL) {
    return NO_MEMORY;
  }
  {
    // Count array
    uint8_t itemCount = 20;
    for(int curItem = 0; curItem < itemCount; curItem++) {
      
      int8_t _value = plc4c_spi_read_byte(buf, 8);
      plc4c_utils_list_insert_head_value(mlfb, &_value);
    }
  }
  (*_message)->mlfb = mlfb;

  // Simple Field (moduleTypeId)
  uint16_t moduleTypeId = plc4c_spi_read_unsigned_int(buf, 16);
  (*_message)->module_type_id = moduleTypeId;

  // Simple Field (ausbg)
  uint16_t ausbg = plc4c_spi_read_unsigned_int(buf, 16);
  (*_message)->ausbg = ausbg;

  // Simple Field (ausbe)
  uint16_t ausbe = plc4c_spi_read_unsigned_int(buf, 16);
  (*_message)->ausbe = ausbe;

  return OK;
}

plc4c_return_code plc4c_s7_read_write_szl_data_tree_item_serialize(plc4c_spi_write_buffer* buf, plc4c_s7_read_write_szl_data_tree_item* message) {
  return OK;
}
