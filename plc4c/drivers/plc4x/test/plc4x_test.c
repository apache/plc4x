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
#include <string.h>
#include <unity.h>

#include "plc4c/driver_plc4x_packets.h"

void parser_serializer_test_plc4x_read_write();

void setUp(void) {}

void tearDown(void) {}

// Every code the proxy can put on the wire, and what a plc4c caller sees for it.
void test_plc4x_response_code_mapping(void) {
  struct {
    plc4c_plc4x_read_write_plc4x_response_code wire;
    plc4c_response_code expected;
  } cases[] = {
      {plc4c_plc4x_read_write_plc4x_response_code_OK, PLC4C_RESPONSE_CODE_OK},
      {plc4c_plc4x_read_write_plc4x_response_code_NOT_FOUND, PLC4C_RESPONSE_CODE_NOT_FOUND},
      {plc4c_plc4x_read_write_plc4x_response_code_ACCESS_DENIED, PLC4C_RESPONSE_CODE_ACCESS_DENIED},
      {plc4c_plc4x_read_write_plc4x_response_code_INVALID_ADDRESS, PLC4C_RESPONSE_CODE_INVALID_ADDRESS},
      {plc4c_plc4x_read_write_plc4x_response_code_INVALID_DATATYPE, PLC4C_RESPONSE_CODE_INVALID_DATATYPE},
      {plc4c_plc4x_read_write_plc4x_response_code_INVALID_DATA, PLC4C_RESPONSE_CODE_INVALID_DATA},
      {plc4c_plc4x_read_write_plc4x_response_code_INTERNAL_ERROR, PLC4C_RESPONSE_CODE_INTERNAL_ERROR},
      {plc4c_plc4x_read_write_plc4x_response_code_REMOTE_BUSY, PLC4C_RESPONSE_CODE_REMOTE_BUSY},
      {plc4c_plc4x_read_write_plc4x_response_code_REMOTE_ERROR, PLC4C_RESPONSE_CODE_REMOTE_ERROR},
      {plc4c_plc4x_read_write_plc4x_response_code_UNSUPPORTED, PLC4C_RESPONSE_CODE_UNSUPPORTED},
      {plc4c_plc4x_read_write_plc4x_response_code_RESPONSE_PENDING, PLC4C_RESPONSE_CODE_RESPONSE_PENDING},
      {plc4c_plc4x_read_write_plc4x_response_code_NOT_READY, PLC4C_RESPONSE_CODE_NOT_READY},
      {plc4c_plc4x_read_write_plc4x_response_code_OUT_OF_RANGE, PLC4C_RESPONSE_CODE_OUT_OF_RANGE},
      {plc4c_plc4x_read_write_plc4x_response_code_REQUEST_TIMEOUT, PLC4C_RESPONSE_CODE_REQUEST_TIMEOUT},
  };
  int num_cases = sizeof(cases) / sizeof(cases[0]);

  // A code added to the wire enum has to be added here, so that its mapping is a decision.
  TEST_ASSERT_EQUAL_INT(plc4c_plc4x_read_write_plc4x_response_code_num_values(), num_cases);
  for (int i = 0; i < num_cases; i++) {
    TEST_ASSERT_EQUAL_INT_MESSAGE(cases[i].expected,
        plc4c_driver_plc4x_map_response_code(cases[i].wire),
        plc4c_response_code_to_message(cases[i].expected));
    TEST_ASSERT_NOT_EQUAL_INT(0, strcmp("UNKNOWN RESPONSE CODE",
        plc4c_response_code_to_message(cases[i].expected)));
  }
  // A value the protocol doesn't define can't be reported as anything in particular.
  TEST_ASSERT_EQUAL_INT(PLC4C_RESPONSE_CODE_INTERNAL_ERROR,
      plc4c_driver_plc4x_map_response_code(
          (plc4c_plc4x_read_write_plc4x_response_code) 0xFF));
}

int main(void) {
  UNITY_BEGIN();

  // Execute the auto-generated serializer/parser tests.
  parser_serializer_test_plc4x_read_write();
  RUN_TEST(test_plc4x_response_code_mapping);

  return UNITY_END();
}