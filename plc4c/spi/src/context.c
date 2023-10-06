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

#include "plc4c/spi/context.h"

plc4x_spi_context plc4x_spi_context_background() {
  plc4x_spi_context ctx = {};
  return ctx;
}

plc4x_spi_context plc4x_spi_context_create_array_context(plc4x_spi_context context, uint16_t numItems, uint16_t curItem) {
  // TODO: In general we would be taking the parent context and copying it over to the new, but currently we only have these settings, so there's no point in doing that.
  plc4x_spi_context ctx = {
      .numItems = numItems,
      .curItem = curItem
  };
  return ctx;
}

uint16_t plc4x_spi_context_get_num_items_from_context(plc4x_spi_context ctx) {
  return ctx.numItems;
}

uint16_t plc4x_spi_context_get_cur_item_from_context(plc4x_spi_context ctx) {
  return ctx.curItem;
}

bool plc4x_spi_context_get_last_item_from_context(plc4x_spi_context ctx) {
  return ctx.curItem == (ctx.numItems - 1);
}
