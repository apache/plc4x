/*
* Licensed to the Apache Software Foundation (ASF) under one
* or more contributor license agreements.  See the NOTICE file
* distributed with this work for additional information
* regarding copyright ownership.  The ASF licenses this file
* to you under the Apache License, Version 2.0 (the
* "License"); you may not use this file except in compliance
* with the License.  You may obtain a copy of the License at
*
*   https://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing,
* software distributed under the License is distributed on an
* "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
* KIND, either express or implied.  See the License for the
* specific language governing permissions and limitations
* under the License.
*/

#ifndef PLC4C_CONTEXT_H
#define PLC4C_CONTEXT_H

#include <stdint.h>
#include <stdbool.h>

struct plc4x_spi_context {
  uint16_t numItems;
  uint16_t curItem;
};
typedef struct plc4x_spi_context plc4x_spi_context;

plc4x_spi_context plc4x_spi_context_background();

plc4x_spi_context plc4x_spi_context_create_array_context(plc4x_spi_context context, uint16_t numItems, uint16_t curItem);

uint16_t plc4x_spi_context_get_num_items_from_context(plc4x_spi_context ctx);

uint16_t plc4x_spi_context_get_cur_item_from_context(plc4x_spi_context ctx);

bool plc4x_spi_context_get_last_item_from_context(plc4x_spi_context ctx);

#endif  // PLC4C_CONTEXT_H
