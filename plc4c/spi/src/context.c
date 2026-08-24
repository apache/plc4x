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

#include <ctype.h>
#include <stdio.h>
#include <stdlib.h>

plc4x_spi_context plc4x_spi_context_background() {
  plc4x_spi_context ctx = {};
  return ctx;
}

plc4x_spi_context plc4x_spi_context_create_array_context(plc4x_spi_context context, uint16_t numItems, uint16_t curItem) {
  // The element sits as deep as the array does, so the depth comes across with
  // it - an array of a self-nesting type would otherwise start counting from
  // nothing at every element and escape the bound entirely.
  plc4x_spi_context ctx = {
      .numItems = numItems,
      .curItem = curItem,
      .depth = context.depth
  };
  return ctx;
}

uint16_t plc4x_spi_context_resolve_max_depth(const char* configured) {
  if (configured == NULL) {
    return PLC4C_DEFAULT_MAX_NESTING_DEPTH;
  }
  while (*configured != '\0' && isspace((unsigned char) *configured)) {
    configured++;
  }
  if (*configured == '\0') {
    return PLC4C_DEFAULT_MAX_NESTING_DEPTH;
  }
  char* end = NULL;
  long depth = strtol(configured, &end, 10);
  while (end != NULL && *end != '\0' && isspace((unsigned char) *end)) {
    end++;
  }
  if (end == configured || (end != NULL && *end != '\0')) {
    fprintf(stderr, "%s is not a number (%s), keeping the maximum nesting depth of %d\n",
            PLC4C_MAX_NESTING_DEPTH_ENV, configured, PLC4C_DEFAULT_MAX_NESTING_DEPTH);
    return PLC4C_DEFAULT_MAX_NESTING_DEPTH;
  }
  if (depth < 1 || depth > UINT16_MAX) {
    fprintf(stderr, "%s must be between 1 and %d but was %ld, keeping the maximum nesting depth of %d\n",
            PLC4C_MAX_NESTING_DEPTH_ENV, UINT16_MAX, depth, PLC4C_DEFAULT_MAX_NESTING_DEPTH);
    return PLC4C_DEFAULT_MAX_NESTING_DEPTH;
  }
  return (uint16_t) depth;
}

uint16_t plc4x_spi_context_get_max_depth() {
  // Resolved on first use and kept, since the environment cannot change under a
  // running process. Two threads arriving together compute the same answer.
  static uint16_t maxDepth = 0;
  if (maxDepth == 0) {
    maxDepth = plc4x_spi_context_resolve_max_depth(getenv(PLC4C_MAX_NESTING_DEPTH_ENV));
  }
  return maxDepth;
}

plc4c_return_code plc4x_spi_context_enter_type(plc4x_spi_context* context) {
  if (context->depth >= plc4x_spi_context_get_max_depth()) {
    return PARSE_ERROR;
  }
  context->depth++;
  return OK;
}

uint16_t plc4x_spi_context_get_depth(plc4x_spi_context ctx) {
  return ctx.depth;
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
