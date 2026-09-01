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

#include <unity.h>

// Every parse function takes the context by value and passes its own copy on to
// the types it contains, so entering a type on that copy is all the depth
// tracking there is - the count unwinds by itself when the function returns.
// These tests stand in for that call chain.

static plc4c_return_code descend(plc4x_spi_context ctx, uint16_t levels) {
  for (uint16_t i = 0; i < levels; i++) {
    plc4c_return_code result = plc4x_spi_context_enter_type(&ctx);
    if (result != OK) {
      return result;
    }
  }
  return OK;
}

void test_plc4c_spi_context_allows_nesting_up_to_the_bound(void) {
  plc4x_spi_context ctx = plc4x_spi_context_background();
  TEST_ASSERT_EQUAL_INT(OK, descend(ctx, plc4x_spi_context_get_max_depth()));
}

void test_plc4c_spi_context_rejects_nesting_past_the_bound(void) {
  plc4x_spi_context ctx = plc4x_spi_context_background();
  TEST_ASSERT_EQUAL_INT(PARSE_ERROR,
                        descend(ctx, plc4x_spi_context_get_max_depth() + 1));
}

void test_plc4c_spi_context_depth_unwinds_with_the_caller(void) {
  plc4x_spi_context ctx = plc4x_spi_context_background();
  // A sibling type parsed after a deep one starts from the caller's depth again,
  // because the deep one only ever raised its own copy.
  TEST_ASSERT_EQUAL_INT(OK, descend(ctx, plc4x_spi_context_get_max_depth()));
  TEST_ASSERT_EQUAL_INT(0, plc4x_spi_context_get_depth(ctx));
  TEST_ASSERT_EQUAL_INT(OK, descend(ctx, plc4x_spi_context_get_max_depth()));
}

void test_plc4c_spi_context_array_context_keeps_the_depth(void) {
  plc4x_spi_context ctx = plc4x_spi_context_background();
  TEST_ASSERT_EQUAL_INT(OK, plc4x_spi_context_enter_type(&ctx));
  TEST_ASSERT_EQUAL_INT(OK, plc4x_spi_context_enter_type(&ctx));

  // An array element gets a context carrying its position. It must carry the
  // depth too, or an array of a self-nesting type escapes the bound entirely.
  plc4x_spi_context item = plc4x_spi_context_create_array_context(ctx, 4, 1);
  TEST_ASSERT_EQUAL_INT(2, plc4x_spi_context_get_depth(item));
  TEST_ASSERT_EQUAL_INT(4, plc4x_spi_context_get_num_items_from_context(item));
  TEST_ASSERT_EQUAL_INT(1, plc4x_spi_context_get_cur_item_from_context(item));
}

void test_plc4c_spi_context_resolve_max_depth(void) {
  TEST_ASSERT_EQUAL_INT(PLC4C_DEFAULT_MAX_NESTING_DEPTH,
                        plc4x_spi_context_resolve_max_depth(NULL));
  TEST_ASSERT_EQUAL_INT(PLC4C_DEFAULT_MAX_NESTING_DEPTH,
                        plc4x_spi_context_resolve_max_depth(""));
  TEST_ASSERT_EQUAL_INT(PLC4C_DEFAULT_MAX_NESTING_DEPTH,
                        plc4x_spi_context_resolve_max_depth("   "));
  TEST_ASSERT_EQUAL_INT(PLC4C_DEFAULT_MAX_NESTING_DEPTH,
                        plc4x_spi_context_resolve_max_depth("plenty"));
  TEST_ASSERT_EQUAL_INT(PLC4C_DEFAULT_MAX_NESTING_DEPTH,
                        plc4x_spi_context_resolve_max_depth("0"));
  TEST_ASSERT_EQUAL_INT(PLC4C_DEFAULT_MAX_NESTING_DEPTH,
                        plc4x_spi_context_resolve_max_depth("-1"));
  TEST_ASSERT_EQUAL_INT(64, plc4x_spi_context_resolve_max_depth("64"));
  TEST_ASSERT_EQUAL_INT(4096, plc4x_spi_context_resolve_max_depth(" 4096 "));
}

void test_plc4c_spi_context(void) {
  test_plc4c_spi_context_allows_nesting_up_to_the_bound();
  test_plc4c_spi_context_rejects_nesting_past_the_bound();
  test_plc4c_spi_context_depth_unwinds_with_the_caller();
  test_plc4c_spi_context_array_context_keeps_the_depth();
  test_plc4c_spi_context_resolve_max_depth();
}
