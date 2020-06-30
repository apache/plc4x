/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
#include "plc4c/spi/read_buffer.h"

#include <stdlib.h>
#include <unity.h>

void test_plc4c_spi_read_buffer_create_args(
    uint8_t* data, uint16_t length, plc4c_return_code expected_return_code) {
  printf("Running read_buffer create test with %d, expecting return code %d",
         length, expected_return_code);

  // Create a new read_buffer instance
  plc4c_spi_read_buffer* read_buffer;

  plc4c_return_code return_code =
      plc4c_spi_read_buffer_create(data, length, &read_buffer);

  TEST_ASSERT_EQUAL_INT(expected_return_code, return_code);
  if (expected_return_code != OK) {
    TEST_ASSERT_NULL(read_buffer);
  } else {
    TEST_ASSERT_NOT_NULL(read_buffer);
    free(read_buffer);
  }
  printf(" -> OK\n");
}

void test_plc4c_spi_read_buffer_create(void) {
  // Prepare input data
  uint8_t data[] = {1, 2, 3, 4, 5, 6, 7, 8};

  // Run test
  test_plc4c_spi_read_buffer_create_args(&data[0], 8, OK);
}

void test_plc4c_spi_read_buffer_get_total_bytes_args(
    plc4c_spi_read_buffer* read_buffer, uint16_t expected_length) {
  printf("Running read_buffer get_total_bytes test. Expecting %d length",
         expected_length);

  uint32_t length = plc4c_spi_read_get_total_bytes(read_buffer);

  TEST_ASSERT_EQUAL_INT(expected_length, length);
  printf(" -> OK\n");
}

void test_plc4c_spi_read_buffer_get_total_bytes(void) {
  // Prepare input data
  uint8_t data[] = {1, 2, 3, 4, 5, 6, 7, 8};
  plc4c_spi_read_buffer* read_buffer;
  plc4c_spi_read_buffer_create(data, 8, &read_buffer);

  // Run test
  test_plc4c_spi_read_buffer_get_total_bytes_args(read_buffer, 8);
}

void test_plc4c_spi_read_buffer_has_more_args(
    plc4c_spi_read_buffer* read_buffer, uint16_t num_bytes,
    bool expected_result) {
  printf(
      "Running read_buffer has_more test. Checking if %d bytes are available",
      num_bytes);

  bool result = plc4c_spi_read_has_more(read_buffer, num_bytes);

  TEST_ASSERT_EQUAL_INT(expected_result, result);
  printf(" -> OK\n");
}

void test_plc4c_spi_read_buffer_has_more(void) {
  // Prepare input data
  uint8_t data[] = {1, 2, 3, 4};
  plc4c_spi_read_buffer* read_buffer;
  plc4c_spi_read_buffer_create(data, 4, &read_buffer);

  // Run test
  test_plc4c_spi_read_buffer_has_more_args(read_buffer, 0, true);
  test_plc4c_spi_read_buffer_has_more_args(read_buffer, 1, true);
  test_plc4c_spi_read_buffer_has_more_args(read_buffer, 4, true);
  test_plc4c_spi_read_buffer_has_more_args(read_buffer, 14, true);
  test_plc4c_spi_read_buffer_has_more_args(read_buffer, 31, true);
  test_plc4c_spi_read_buffer_has_more_args(read_buffer, 32, true);

  // 4 bytes are 32 bits, so these should fail.
  test_plc4c_spi_read_buffer_has_more_args(read_buffer, 33, false);
  test_plc4c_spi_read_buffer_has_more_args(read_buffer, 50, false);
}

void test_plc4c_spi_read_buffer_get_bytes_args(
    plc4c_spi_read_buffer* read_buffer, uint16_t start_byte, uint16_t end_byte,
    plc4c_return_code expected_return_code, const uint8_t* expected_bytes,
    uint8_t expected_bytes_length) {
  printf(
      "Running read_buffer get_bytes test. Checking if reading from %d to %d "
      "bytes gives the correct response",
      start_byte, end_byte);

  uint8_t* read_bytes = NULL;
  plc4c_return_code result =
      plc4c_spi_read_get_bytes(read_buffer, start_byte, end_byte, &read_bytes);

  TEST_ASSERT_EQUAL_INT(expected_return_code, result);
  if (expected_return_code != OK) {
    TEST_ASSERT_NULL(read_bytes);
  } else {
    TEST_ASSERT_NOT_NULL(read_bytes);
    for (int i = 0; i < expected_bytes_length; i++) {
      uint8_t cur_byte = *read_buffer->data + start_byte + i;
      uint8_t expected_byte = *expected_bytes + i;
      TEST_ASSERT_EQUAL_INT(expected_byte, cur_byte);
    }
  }

  printf(" -> OK\n");
}

void test_plc4c_spi_read_buffer_get_bytes(void) {
  // Prepare input data
  uint8_t data[] = {1, 2, 3, 4, 5, 6, 7};
  plc4c_spi_read_buffer* read_buffer;
  plc4c_spi_read_buffer_create(data, 8, &read_buffer);

  // Run test
  uint8_t result_data1[] = {1};
  test_plc4c_spi_read_buffer_get_bytes_args(read_buffer, 0, 0, OK, result_data1,
                                            1);
  uint8_t result_data2[] = {2};
  test_plc4c_spi_read_buffer_get_bytes_args(read_buffer, 1, 1, OK, result_data2,
                                            1);
  uint8_t result_data3[] = {2, 3};
  test_plc4c_spi_read_buffer_get_bytes_args(read_buffer, 1, 2, OK, result_data3,
                                            2);
  uint8_t result_data4[] = {4, 5, 6};
  test_plc4c_spi_read_buffer_get_bytes_args(read_buffer, 3, 5, OK, result_data4,
                                            3);
  uint8_t result_data5[] = {4, 5, 6, 7};
  test_plc4c_spi_read_buffer_get_bytes_args(read_buffer, 3, 6, OK, result_data5,
                                            4);

  // These should fail for various reasons ...
  uint8_t result_data6[] = {};
  test_plc4c_spi_read_buffer_get_bytes_args(read_buffer, 0, 10, OUT_OF_RANGE,
                                            result_data6, 0);
  test_plc4c_spi_read_buffer_get_bytes_args(read_buffer, 6, 3, INVALID_ARGUMENT,
                                            result_data6, 0);
  test_plc4c_spi_read_buffer_get_bytes_args(NULL, 0, 0, NULL_VALUE,
                                            result_data6, 0);
}

void test_plc4c_spi_read_buffer_peek_byte_args(
    plc4c_spi_read_buffer* read_buffer, uint16_t peek_byte,
    plc4c_return_code expected_return_code, uint8_t expected_value) {
  printf(
      "Running read_buffer peek_byte test. Checking if peeking byte number %d "
      "gives the correct response",
      peek_byte);

  uint8_t peeked_byte = 0;
  plc4c_return_code result =
      plc4c_spi_read_peek_byte(read_buffer, peek_byte, &peeked_byte);

  TEST_ASSERT_EQUAL_INT(expected_return_code, result);
  TEST_ASSERT_EQUAL_INT(expected_value, peeked_byte);

  printf(" -> OK\n");
}

void test_plc4c_spi_read_buffer_peek_byte(void) {
  // Prepare input data
  uint8_t data[] = {1, 2, 3, 4, 5, 6, 7};
  plc4c_spi_read_buffer* read_buffer;
  plc4c_spi_read_buffer_create(data, 8, &read_buffer);

  // Run test
  test_plc4c_spi_read_buffer_peek_byte_args(read_buffer, 0, OK, 1);
  test_plc4c_spi_read_buffer_peek_byte_args(read_buffer, 4, OK, 5);
  test_plc4c_spi_read_buffer_peek_byte_args(read_buffer, 4, OK, 5);
  test_plc4c_spi_read_buffer_peek_byte_args(read_buffer, 6, OK, 7);
  // Bump the cur buffer position.
  read_buffer->curPosByte = 2;
  test_plc4c_spi_read_buffer_peek_byte_args(read_buffer, 2, OK, 5);

  // These should fail
  test_plc4c_spi_read_buffer_peek_byte_args(read_buffer, 8, OUT_OF_RANGE, 0);
}

void test_plc4c_spi_read_buffer_read_bit_args(
    plc4c_spi_read_buffer* read_buffer, bool expected_value) {
  printf(
      "Running read_buffer peek_byte test. Checking if reading a bit gives the "
      "correct response");

  bool value = false;
  plc4c_return_code result = plc4c_spi_read_bit(read_buffer, &value);

  TEST_ASSERT_EQUAL_INT(OK, result);
  TEST_ASSERT_EQUAL_INT(expected_value, value);

  printf(" -> OK\n");
}

void test_plc4c_spi_read_buffer_read_bit(void) {
  // Prepare input data
  uint8_t data[] = {1, 2, 3, 4, 5, 6, 7};
  plc4c_spi_read_buffer* read_buffer;
  plc4c_spi_read_buffer_create(data, 8, &read_buffer);

  // Run test
  test_plc4c_spi_read_buffer_read_bit_args(read_buffer, false);
  test_plc4c_spi_read_buffer_read_bit_args(read_buffer, false);
  test_plc4c_spi_read_buffer_read_bit_args(read_buffer, false);
  test_plc4c_spi_read_buffer_read_bit_args(read_buffer, false);
  test_plc4c_spi_read_buffer_read_bit_args(read_buffer, false);
  test_plc4c_spi_read_buffer_read_bit_args(read_buffer, false);
  test_plc4c_spi_read_buffer_read_bit_args(read_buffer, false);
  test_plc4c_spi_read_buffer_read_bit_args(read_buffer, true);
  test_plc4c_spi_read_buffer_read_bit_args(read_buffer, false);
  test_plc4c_spi_read_buffer_read_bit_args(read_buffer, false);
  test_plc4c_spi_read_buffer_read_bit_args(read_buffer, false);
  test_plc4c_spi_read_buffer_read_bit_args(read_buffer, false);
  test_plc4c_spi_read_buffer_read_bit_args(read_buffer, false);
  test_plc4c_spi_read_buffer_read_bit_args(read_buffer, false);
  test_plc4c_spi_read_buffer_read_bit_args(read_buffer, true);
  test_plc4c_spi_read_buffer_read_bit_args(read_buffer, false);
  test_plc4c_spi_read_buffer_read_bit_args(read_buffer, false);
  test_plc4c_spi_read_buffer_read_bit_args(read_buffer, false);
  test_plc4c_spi_read_buffer_read_bit_args(read_buffer, false);
  test_plc4c_spi_read_buffer_read_bit_args(read_buffer, false);
  test_plc4c_spi_read_buffer_read_bit_args(read_buffer, false);
  test_plc4c_spi_read_buffer_read_bit_args(read_buffer, false);
  test_plc4c_spi_read_buffer_read_bit_args(read_buffer, true);
  test_plc4c_spi_read_buffer_read_bit_args(read_buffer, true);
}

void test_plc4c_spi_read_buffer(void) {
  test_plc4c_spi_read_buffer_create();
  test_plc4c_spi_read_buffer_get_total_bytes();
  test_plc4c_spi_read_buffer_has_more();
  test_plc4c_spi_read_buffer_get_bytes();
  test_plc4c_spi_read_buffer_peek_byte();
  test_plc4c_spi_read_buffer_read_bit();
}