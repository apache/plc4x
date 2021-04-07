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

#include <plc4c/spi/system_private.h>
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

void test_plc4c_spi_read_get_total_bytes_args(
    plc4c_spi_read_buffer* read_buffer, uint16_t expected_length) {
  printf("Running read_buffer get_total_bytes test. Expecting %d length",
         expected_length);

  uint32_t length = plc4c_spi_read_get_total_bytes(read_buffer);

  TEST_ASSERT_EQUAL_INT(expected_length, length);
  printf(" -> OK\n");
}

void test_plc4c_spi_read_get_total_bytes(void) {
  // Prepare input data
  uint8_t data[] = {1, 2, 3, 4, 5, 6, 7, 8};
  plc4c_spi_read_buffer* read_buffer;
  plc4c_spi_read_buffer_create(data, 8, &read_buffer);

  // Run test
  test_plc4c_spi_read_get_total_bytes_args(read_buffer, 8);
}

void test_plc4c_spi_read_has_more_args(plc4c_spi_read_buffer* read_buffer,
                                       uint16_t num_bytes,
                                       bool expected_result) {
  printf(
      "Running read_buffer has_more test. Checking if %d bytes are available",
      num_bytes);

  bool result = plc4c_spi_read_has_more(read_buffer, num_bytes);

  TEST_ASSERT_EQUAL_INT(expected_result, result);
  printf(" -> OK\n");
}

void test_plc4c_spi_read_has_more(void) {
  // Prepare input data
  uint8_t data[] = {1, 2, 3, 4};
  plc4c_spi_read_buffer* read_buffer;
  plc4c_spi_read_buffer_create(data, 4, &read_buffer);

  // Run test
  test_plc4c_spi_read_has_more_args(read_buffer, 0, true);
  test_plc4c_spi_read_has_more_args(read_buffer, 1, true);
  test_plc4c_spi_read_has_more_args(read_buffer, 4, true);
  test_plc4c_spi_read_has_more_args(read_buffer, 14, true);
  test_plc4c_spi_read_has_more_args(read_buffer, 31, true);
  test_plc4c_spi_read_has_more_args(read_buffer, 32, true);

  // 4 bytes are 32 bits, so these should fail.
  test_plc4c_spi_read_has_more_args(read_buffer, 33, false);
  test_plc4c_spi_read_has_more_args(read_buffer, 50, false);
}

void test_plc4c_spi_read_get_bytes_args(plc4c_spi_read_buffer* read_buffer,
                                        uint16_t start_byte, uint16_t end_byte,
                                        plc4c_return_code expected_return_code,
                                        const uint8_t* expected_bytes,
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

void test_plc4c_spi_read_get_bytes(void) {
  // Prepare input data
  uint8_t data[] = {1, 2, 3, 4, 5, 6, 7, 8};
  plc4c_spi_read_buffer* read_buffer;
  plc4c_spi_read_buffer_create(data, 8, &read_buffer);

  // Run test
  uint8_t result_data1[] = {1};
  test_plc4c_spi_read_get_bytes_args(read_buffer, 0, 0, OK, result_data1, 1);
  uint8_t result_data2[] = {2};
  test_plc4c_spi_read_get_bytes_args(read_buffer, 1, 1, OK, result_data2, 1);
  uint8_t result_data3[] = {2, 3};
  test_plc4c_spi_read_get_bytes_args(read_buffer, 1, 2, OK, result_data3, 2);
  uint8_t result_data4[] = {4, 5, 6};
  test_plc4c_spi_read_get_bytes_args(read_buffer, 3, 5, OK, result_data4, 3);
  uint8_t result_data5[] = {4, 5, 6, 7};
  test_plc4c_spi_read_get_bytes_args(read_buffer, 3, 6, OK, result_data5, 4);

  // These should fail for various reasons ...
  uint8_t result_data6[] = {};
  test_plc4c_spi_read_get_bytes_args(read_buffer, 0, 10, OUT_OF_RANGE,
                                     result_data6, 0);
  test_plc4c_spi_read_get_bytes_args(read_buffer, 6, 3, INVALID_ARGUMENT,
                                     result_data6, 0);
  test_plc4c_spi_read_get_bytes_args(NULL, 0, 0, NULL_VALUE, result_data6, 0);
}

void test_plc4c_spi_read_peek_byte_args(plc4c_spi_read_buffer* read_buffer,
                                        uint16_t peek_byte,
                                        plc4c_return_code expected_return_code,
                                        uint8_t expected_value) {
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

void test_plc4c_spi_read_peek_byte(void) {
  // Prepare input data
  uint8_t data[] = {1, 2, 3, 4, 5, 6, 7, 8};
  plc4c_spi_read_buffer* read_buffer;
  plc4c_spi_read_buffer_create(data, 8, &read_buffer);

  // Run test
  test_plc4c_spi_read_peek_byte_args(read_buffer, 0, OK, 1);
  test_plc4c_spi_read_peek_byte_args(read_buffer, 4, OK, 5);
  test_plc4c_spi_read_peek_byte_args(read_buffer, 4, OK, 5);
  test_plc4c_spi_read_peek_byte_args(read_buffer, 6, OK, 7);
  // Bump the cur buffer position.
  read_buffer->curPosByte = 2;
  test_plc4c_spi_read_peek_byte_args(read_buffer, 2, OK, 5);

  // These should fail
  test_plc4c_spi_read_peek_byte_args(read_buffer, 8, OUT_OF_RANGE, 0);
}

void test_plc4c_spi_read_read_bit_args(plc4c_spi_read_buffer* read_buffer,
                                       bool expected_value) {
  printf(
      "Running read_buffer peek_byte test. Checking if reading a bit gives the "
      "correct response");

  bool value = false;
  plc4c_return_code result = plc4c_spi_read_bit(read_buffer, &value);

  TEST_ASSERT_EQUAL_INT(OK, result);
  TEST_ASSERT_EQUAL_INT(expected_value, value);

  printf(" -> OK\n");
}

void test_plc4c_spi_read_read_bit(void) {
  // Prepare input data
  uint8_t data[] = {1, 2, 3, 4, 5, 6, 7, 8};
  plc4c_spi_read_buffer* read_buffer;
  plc4c_spi_read_buffer_create(data, 8, &read_buffer);

  // Run test
  test_plc4c_spi_read_read_bit_args(read_buffer, false);
  test_plc4c_spi_read_read_bit_args(read_buffer, false);
  test_plc4c_spi_read_read_bit_args(read_buffer, false);
  test_plc4c_spi_read_read_bit_args(read_buffer, false);
  test_plc4c_spi_read_read_bit_args(read_buffer, false);
  test_plc4c_spi_read_read_bit_args(read_buffer, false);
  test_plc4c_spi_read_read_bit_args(read_buffer, false);
  test_plc4c_spi_read_read_bit_args(read_buffer, true);
  test_plc4c_spi_read_read_bit_args(read_buffer, false);
  test_plc4c_spi_read_read_bit_args(read_buffer, false);
  test_plc4c_spi_read_read_bit_args(read_buffer, false);
  test_plc4c_spi_read_read_bit_args(read_buffer, false);
  test_plc4c_spi_read_read_bit_args(read_buffer, false);
  test_plc4c_spi_read_read_bit_args(read_buffer, false);
  test_plc4c_spi_read_read_bit_args(read_buffer, true);
  test_plc4c_spi_read_read_bit_args(read_buffer, false);
  test_plc4c_spi_read_read_bit_args(read_buffer, false);
  test_plc4c_spi_read_read_bit_args(read_buffer, false);
  test_plc4c_spi_read_read_bit_args(read_buffer, false);
  test_plc4c_spi_read_read_bit_args(read_buffer, false);
  test_plc4c_spi_read_read_bit_args(read_buffer, false);
  test_plc4c_spi_read_read_bit_args(read_buffer, false);
  test_plc4c_spi_read_read_bit_args(read_buffer, true);
  test_plc4c_spi_read_read_bit_args(read_buffer, true);
}

void test_plc4c_spi_read_unsigned_byte_args(char* message,
    plc4c_spi_read_buffer* read_buffer, uint8_t num_bits,
    plc4c_return_code expected_return_code, uint8_t expected_value) {
  printf("Running read_buffer read_unsigned_byte test: %s", message);

  uint8_t value = 0;
  plc4c_return_code result =
      plc4c_spi_read_unsigned_byte(read_buffer, num_bits, &value);

  TEST_ASSERT_EQUAL_INT(expected_return_code, result);
  TEST_ASSERT_EQUAL_INT(expected_value, value);

  printf(" -> OK\n");
}

void test_plc4c_spi_read_unsigned_byte(void) {
  // Prepare input data
  uint8_t data[] = {1, 2, 3, 4, 5, 6, 7, 8};
  plc4c_spi_read_buffer* read_buffer;
  plc4c_spi_read_buffer_create(data, 8, &read_buffer);
  // Run test
  // Read all the full bytes
  test_plc4c_spi_read_unsigned_byte_args("Simple full unsigned byte 1", read_buffer, 8, OK, 1);
  test_plc4c_spi_read_unsigned_byte_args("Simple full unsigned byte 2", read_buffer, 8, OK, 2);
  test_plc4c_spi_read_unsigned_byte_args("Simple full unsigned byte 3", read_buffer, 8, OK, 3);
  test_plc4c_spi_read_unsigned_byte_args("Simple full unsigned byte 4", read_buffer, 8, OK, 4);
  test_plc4c_spi_read_unsigned_byte_args("Simple full unsigned byte 5", read_buffer, 8, OK, 5);
  test_plc4c_spi_read_unsigned_byte_args("Simple full unsigned byte 6", read_buffer, 8, OK, 6);
  test_plc4c_spi_read_unsigned_byte_args("Simple full unsigned byte 7", read_buffer, 8, OK, 7);
  test_plc4c_spi_read_unsigned_byte_args("Simple full unsigned byte 8", read_buffer, 8, OK, 8);
  // Read a 9th byte (buffer only has 8) (results in error)
  test_plc4c_spi_read_unsigned_byte_args("Exceed read-buffer size", read_buffer, 8, OUT_OF_RANGE, 0);
  plc4c_spi_read_buffer_destroy(read_buffer);

  uint8_t data2[] = {186, 117};
  plc4c_spi_read_buffer* read_buffer2;
  plc4c_spi_read_buffer_create(data2, 2, &read_buffer2);
  // Read part of a byte (fits in one byte)
  test_plc4c_spi_read_unsigned_byte_args("Simple 4 bits of unsigned byte", read_buffer2, 4, OK, 11);
  // Read part of a byte (finishes one byte)
  test_plc4c_spi_read_unsigned_byte_args("Simple 4 bits of unsigned byte, finishing rest of first byte", read_buffer2, 4, OK, 10);
  test_plc4c_spi_read_unsigned_byte_args("Simple 4 bits of unsigned byte", read_buffer2, 4, OK, 7);
  test_plc4c_spi_read_unsigned_byte_args("Simple 4 bits of unsigned byte", read_buffer2, 4, OK, 5);
  // Read part of a byte (spans two bytes)
  read_buffer2->curPosByte = 0;
  read_buffer2->curPosBit = 5;
  test_plc4c_spi_read_unsigned_byte_args("Simple 6 bits of unsigned byte starting at bit 5 (flowing over to next byte)", read_buffer2, 6, OK, 19);
  test_plc4c_spi_read_unsigned_byte_args("Simple 4 bits of unsigned byte starting at bit 3", read_buffer2, 4, OK, 10);
  // Read more than a byte (results in error)
  read_buffer2->curPosByte = 0;
  read_buffer2->curPosBit = 5;
  test_plc4c_spi_read_unsigned_byte_args("Exceed read-buffer size (Part 2)", read_buffer2, 10, OUT_OF_RANGE, 0);
}

void test_plc4c_spi_read_unsigned_short_args(char* message,
    plc4c_spi_read_buffer* read_buffer, uint8_t num_bits,
    plc4c_return_code expected_return_code, uint16_t expected_value) {
  printf("Running read_buffer read_unsigned_short test: %s", message);

  uint16_t value = 0;
  plc4c_return_code result =
      plc4c_spi_read_unsigned_short(read_buffer, num_bits, &value);

  TEST_ASSERT_EQUAL_INT(expected_return_code, result);
  TEST_ASSERT_EQUAL_INT(expected_value, value);

  printf(" -> OK\n");
}

void test_plc4c_spi_read_unsigned_short(void) {
  // Prepare input data
  uint8_t data[] = {1, 2, 3, 4, 5, 6, 7, 8};
  plc4c_spi_read_buffer* read_buffer;
  plc4c_spi_read_buffer_create(data, 8, &read_buffer);
  // Run test
  // Read all the full short
  test_plc4c_spi_read_unsigned_short_args("Simple full short 1", read_buffer, 16, OK, 258);
  test_plc4c_spi_read_unsigned_short_args("Simple full short 2", read_buffer, 16, OK, 772);
  test_plc4c_spi_read_unsigned_short_args("Simple full short 3", read_buffer, 16, OK, 1286);
  test_plc4c_spi_read_unsigned_short_args("Simple full short 4", read_buffer, 16, OK, 1800);

  // Read a short that spans across 3 bytes
  read_buffer->curPosByte = 0;
  read_buffer->curPosBit = 3;
  test_plc4c_spi_read_unsigned_short_args("Full short starting at bit 3", read_buffer, 14, OK, 516);

  // Finish a partial short
  read_buffer->curPosByte = 0;
  read_buffer->curPosBit = 1;
  test_plc4c_spi_read_unsigned_short_args("Short starting at bit 1 but beading full last byte", read_buffer, 15, OK, 258);

  // Read only fractions of a short
  read_buffer->curPosByte = 0;
  read_buffer->curPosBit = 0;
  test_plc4c_spi_read_unsigned_short_args("Short only reading 2 bits", read_buffer, 2, OK, 0);
  test_plc4c_spi_read_unsigned_short_args("Short only reading 4 bits", read_buffer, 4, OK, 0);
  test_plc4c_spi_read_unsigned_short_args("Short only reading 6 bits", read_buffer, 6, OK, 16);
  test_plc4c_spi_read_unsigned_short_args("Short only reading 8 bits", read_buffer, 8, OK, 32);
  test_plc4c_spi_read_unsigned_short_args("Short only reading 10 bits", read_buffer, 10, OK, 193);
  test_plc4c_spi_read_unsigned_short_args("Short only reading 12 bits", read_buffer, 12, OK, 20);
  test_plc4c_spi_read_unsigned_short_args("Short only reading 14 bits", read_buffer, 14, OK, 1543);

  // We only have 8 bytes, so with this we would exceed the range.
  test_plc4c_spi_read_unsigned_short_args("Try to read mode bytes than the buffer has available", read_buffer, 16, OUT_OF_RANGE, 0);
  test_plc4c_spi_read_unsigned_short_args("Try to read too many bytes for a short", read_buffer, 18, OUT_OF_RANGE, 0);
}

void test_plc4c_spi_read_unsigned_int_args(char* message,
    plc4c_spi_read_buffer* read_buffer, uint8_t num_bits,
    plc4c_return_code expected_return_code, uint32_t expected_value) {
  printf("Running read_buffer read_unsigned_int test: %s", message);

  uint32_t value = 0;
  plc4c_return_code result =
      plc4c_spi_read_unsigned_int(read_buffer, num_bits, &value);

  TEST_ASSERT_EQUAL_INT(expected_return_code, result);
  TEST_ASSERT_EQUAL_INT(expected_value, value);

  printf(" -> OK\n");
}

void test_plc4c_spi_read_unsigned_int(void) {
  // Prepare input data
  uint8_t data[] = {1, 2, 3, 4, 5, 6, 7, 8};
  plc4c_spi_read_buffer* read_buffer;
  plc4c_spi_read_buffer_create(data, 8, &read_buffer);
  // Run test
  // Read all the full short
  test_plc4c_spi_read_unsigned_int_args("Simple full int 1", read_buffer, 32, OK, 16909060);
  test_plc4c_spi_read_unsigned_int_args("Simple full int 1", read_buffer, 32, OK, 84281096);

  // Read a full int starting somewhere in between.
  read_buffer->curPosByte = 0;
  read_buffer->curPosBit = 3;
  test_plc4c_spi_read_unsigned_int_args("Full int starting at bit 3", read_buffer, 32, OK, 135272480);

  // Read an int starting somewhere in between but ending at a full byte.
  read_buffer->curPosByte = 0;
  read_buffer->curPosBit = 3;
  test_plc4c_spi_read_unsigned_int_args("Full int starting at bit 3", read_buffer, 29, OK, 16909060);

  // Read a full int starting somewhere in between.
  read_buffer->curPosByte = 0;
  read_buffer->curPosBit = 3;
  test_plc4c_spi_read_unsigned_int_args("Int with only 5 bit", read_buffer, 5, OK, 1);
  test_plc4c_spi_read_unsigned_int_args("Int with only 7 bit", read_buffer, 7, OK, 1);
  test_plc4c_spi_read_unsigned_int_args("Int with only 9 bit", read_buffer, 9, OK, 3);
  test_plc4c_spi_read_unsigned_int_args("Int with only 15 bit", read_buffer, 15, OK, 514);
  test_plc4c_spi_read_unsigned_int_args("Int with only 16 bit", read_buffer, 16, OK, 33539);
  read_buffer->curPosByte = 0;
  read_buffer->curPosBit = 3;
  test_plc4c_spi_read_unsigned_int_args("Int with only 17 bit", read_buffer, 17, OK, 4128);
  test_plc4c_spi_read_unsigned_int_args("Int with only 26 bit", read_buffer, 26, OK, 12648769);
  read_buffer->curPosByte = 0;
  read_buffer->curPosBit = 3;
  test_plc4c_spi_read_unsigned_int_args("Int with only 30 bit", read_buffer, 30, OK, 33818120);
  read_buffer->curPosByte = 0;
  read_buffer->curPosBit = 3;
  test_plc4c_spi_read_unsigned_int_args("Int with only 31 bit", read_buffer, 31, OK, 67636240);
}

void test_plc4c_spi_read_unsigned_long_args(char* message,
    plc4c_spi_read_buffer* read_buffer, uint8_t num_bits,
    plc4c_return_code expected_return_code, uint64_t expected_value) {
  printf("Running read_buffer read_unsigned_long test: %s", message);

  uint64_t value = 0;
  plc4c_return_code result =
      plc4c_spi_read_unsigned_long(read_buffer, num_bits, &value);

  TEST_ASSERT_EQUAL_INT(expected_return_code, result);
  TEST_ASSERT_EQUAL_INT(expected_value, value);

  printf(" -> OK\n");
}

void test_plc4c_spi_read_unsigned_long(void) {
  // Prepare input data
  uint8_t data[] = {1, 2, 3, 4, 5, 6, 7, 8, 9};
  plc4c_spi_read_buffer* read_buffer;
  plc4c_spi_read_buffer_create(data, 9, &read_buffer);
  // Run test
  // Read all the full long
  test_plc4c_spi_read_unsigned_long_args("Simple full long", read_buffer, 64, OK, 72623859790382856);

  // Read a full long starting somewhere in between.
  read_buffer->curPosByte = 0;
  read_buffer->curPosBit = 3;
  test_plc4c_spi_read_unsigned_long_args("Full long starting at bit 3", read_buffer, 64, OK, 580990878323062848);

  // Read a long starting somewhere in between but ending at a full byte.
  read_buffer->curPosByte = 0;
  read_buffer->curPosBit = 3;
  test_plc4c_spi_read_unsigned_long_args("Full long starting at bit 3", read_buffer, 61, OK, 72623859790382856);
}

void test_plc4c_spi_read_signed_byte_args(char* message,
                                            plc4c_spi_read_buffer* read_buffer, uint8_t num_bits,
                                            plc4c_return_code expected_return_code, int8_t expected_value) {
  printf("Running read_buffer read_signed_byte test: %s", message);

  int8_t value = -1;
  plc4c_return_code result =
      plc4c_spi_read_signed_byte(read_buffer, num_bits, &value);

  TEST_ASSERT_EQUAL_INT(expected_return_code, result);
  TEST_ASSERT_EQUAL_INT(expected_value, value);

  printf(" -> OK\n");
}

void test_plc4c_spi_read_signed_byte(void) {
  // Prepare input data
  uint8_t data[] = {1, 2, 3, 4, 5, 6, 7, 8};
  plc4c_spi_read_buffer* read_buffer;
  plc4c_spi_read_buffer_create(data, 8, &read_buffer);
  // Run test
  // Read all the full bytes
  test_plc4c_spi_read_signed_byte_args("Simple full signed byte 1", read_buffer, 8, OK, 1);
  test_plc4c_spi_read_signed_byte_args("Simple full signed byte 2", read_buffer, 8, OK, 2);
  test_plc4c_spi_read_signed_byte_args("Simple full signed byte 3", read_buffer, 8, OK, 3);
  test_plc4c_spi_read_signed_byte_args("Simple full signed byte 4", read_buffer, 8, OK, 4);
  test_plc4c_spi_read_signed_byte_args("Simple full signed byte 5", read_buffer, 8, OK, 5);
  test_plc4c_spi_read_signed_byte_args("Simple full signed byte 6", read_buffer, 8, OK, 6);
  test_plc4c_spi_read_signed_byte_args("Simple full signed byte 7", read_buffer, 8, OK, 7);
  // Read a 9th byte (buffer only has 8) (results in error)
//  test_plc4c_spi_read_signed_byte_args("Exceed read-buffer size", read_buffer, 8, OUT_OF_RANGE, -1);
  plc4c_spi_read_buffer_destroy(read_buffer);

  uint8_t data2[] = {186, 117};
  plc4c_spi_read_buffer* read_buffer2;
  plc4c_spi_read_buffer_create(data2, 2, &read_buffer2);
  // Read part of a byte (fits in one byte)
  test_plc4c_spi_read_signed_byte_args("Simple 4 bits of signed byte", read_buffer2, 4, OK, -5);
  // Read part of a byte (finishes one byte)
  test_plc4c_spi_read_signed_byte_args("Simple 4 bits of signed byte, finishing rest of first byte", read_buffer2, 4, OK, -6);
  test_plc4c_spi_read_signed_byte_args("Simple 4 bits of signed byte", read_buffer2, 4, OK, 7);
  test_plc4c_spi_read_signed_byte_args("Simple 4 bits of signed byte", read_buffer2, 4, OK, 5);
  // Read part of a byte (spans two bytes)
  read_buffer2->curPosByte = 0;
  read_buffer2->curPosBit = 5;
  test_plc4c_spi_read_signed_byte_args("Simple 6 bits of signed byte starting at bit 5 (flowing over to next byte)", read_buffer2, 6, OK, 19);
  test_plc4c_spi_read_signed_byte_args("Simple 4 bits of signed byte starting at bit 3", read_buffer2, 4, OK, -6);
  // Read more than a byte (results in error)
  read_buffer2->curPosByte = 0;
  read_buffer2->curPosBit = 5;
  test_plc4c_spi_read_signed_byte_args("Exceed read-buffer size (Part 2)", read_buffer2, 10, OUT_OF_RANGE, -1);
}

void test_plc4c_spi_read_signed_short_args(char* message,
                                          plc4c_spi_read_buffer* read_buffer, uint8_t num_bits,
                                          plc4c_return_code expected_return_code, int16_t expected_value) {
  printf("Running read_buffer read_signed_byte test: %s", message);

  int16_t value = -1;
  plc4c_return_code result =
      plc4c_spi_read_signed_short(read_buffer, num_bits, &value);

  TEST_ASSERT_EQUAL_INT(expected_return_code, result);
  TEST_ASSERT_EQUAL_INT(expected_value, value);

  printf(" -> OK\n");
}

void test_plc4c_spi_read_signed_short(void) {
  // Prepare input data
  uint8_t data[] = {255, 214, 3};
  plc4c_spi_read_buffer* read_buffer;
  plc4c_spi_read_buffer_create(data, 3, &read_buffer);
  // Run test
  // Read all the full bytes
  test_plc4c_spi_read_signed_short_args("Simple full signed short 1", read_buffer, 16, OK, -42);

  // Read the only part of a short (having to fill up 1s)
  read_buffer->curPosByte = 0;
  read_buffer->curPosBit = 4;
  test_plc4c_spi_read_signed_short_args("Simple 12 bit signed short", read_buffer, 12, OK, -42);

  // Read an even shorter part of a short (having to fill up even more 1s)
  read_buffer->curPosByte = 1;
  read_buffer->curPosBit = 1;
  test_plc4c_spi_read_signed_short_args("Simple 7 bit signed short", read_buffer, 7, OK, -42);

  // Read an even shorter part of a short (This time however the value should
  // be positive and hence the higher level byte should be filled with 0s)
  read_buffer->curPosByte = 1;
  read_buffer->curPosBit = 2;
  test_plc4c_spi_read_signed_short_args("Simple 6 bit signed short", read_buffer, 6, OK, 22);
}

void test_plc4c_spi_read_signed_int_args(char* message,
                                           plc4c_spi_read_buffer* read_buffer, uint8_t num_bits,
                                           plc4c_return_code expected_return_code, int32_t expected_value) {
  printf("Running read_buffer read_signed_byte test: %s", message);

  int32_t value = -1;
  plc4c_return_code result =
      plc4c_spi_read_signed_int(read_buffer, num_bits, &value);

  TEST_ASSERT_EQUAL_INT(expected_return_code, result);
  TEST_ASSERT_EQUAL_INT(expected_value, value);

  printf(" -> OK\n");
}

void test_plc4c_spi_read_signed_int(void) {
  // Prepare input data
  uint8_t data[] = {255, 255, 255, 214, 3};
  plc4c_spi_read_buffer* read_buffer;
  plc4c_spi_read_buffer_create(data, 5, &read_buffer);
  // Run test
  // Read all the full bytes
  test_plc4c_spi_read_signed_int_args("Simple full signed int", read_buffer, 32, OK, -42);
  // Read the only part of a int (having to fill up 1s)
  read_buffer->curPosByte = 0;
  read_buffer->curPosBit = 4;
  test_plc4c_spi_read_signed_int_args("Simple 28 bit signed int", read_buffer, 28, OK, -42);
  // Read the only part of a int (having to fill up 1s)
  read_buffer->curPosByte = 1;
  read_buffer->curPosBit = 4;
  test_plc4c_spi_read_signed_int_args("Simple 20 bit signed int", read_buffer, 20, OK, -42);
  // Read the only part of a int (having to fill up 1s)
  read_buffer->curPosByte = 2;
  read_buffer->curPosBit = 4;
  test_plc4c_spi_read_signed_int_args("Simple 12 bit signed int", read_buffer, 12, OK, -42);

  // Read an even shorter part of a int (having to fill up even more 1s)
  read_buffer->curPosByte = 3;
  read_buffer->curPosBit = 1;
  test_plc4c_spi_read_signed_int_args("Simple 7 bit signed int", read_buffer, 7, OK, -42);

  // Read an even shorter part of a int (This time however the value should
  // be positive and hence the higher level byte should be filled with 0s)
  read_buffer->curPosByte = 3;
  read_buffer->curPosBit = 2;
  test_plc4c_spi_read_signed_int_args("Simple 6 bit signed int", read_buffer, 6, OK, 22);
}

void test_plc4c_spi_read_signed_long_args(char* message,
                                         plc4c_spi_read_buffer* read_buffer, uint8_t num_bits,
                                         plc4c_return_code expected_return_code, int64_t expected_value) {
  printf("Running read_buffer read_signed_byte test: %s", message);

  int64_t value = -1;
  plc4c_return_code result =
      plc4c_spi_read_signed_long(read_buffer, num_bits, &value);

  TEST_ASSERT_EQUAL_INT(expected_return_code, result);
  TEST_ASSERT_EQUAL_INT(expected_value, value);

  printf(" -> OK\n");
}

void test_plc4c_spi_read_signed_long(void) {
  // Prepare input data
  uint8_t data[] = {255, 255, 255, 255, 255, 255, 255, 214, 3};
  plc4c_spi_read_buffer* read_buffer;
  plc4c_spi_read_buffer_create(data, 9, &read_buffer);
  // Run test
  // Read all the full bytes
  test_plc4c_spi_read_signed_long_args("Simple full signed long", read_buffer, 64, OK, -42);
  // Read the only part of a long (having to fill up 1s)
  read_buffer->curPosByte = 0;
  read_buffer->curPosBit = 4;
  test_plc4c_spi_read_signed_long_args("Simple 60 bit signed long", read_buffer, 60, OK, -42);
  // Read the only part of a long (having to fill up 1s)
  read_buffer->curPosByte = 1;
  read_buffer->curPosBit = 4;
  test_plc4c_spi_read_signed_long_args("Simple 52 bit signed long", read_buffer, 52, OK, -42);
  // Read the only part of a long (having to fill up 1s)
  read_buffer->curPosByte = 2;
  read_buffer->curPosBit = 4;
  test_plc4c_spi_read_signed_long_args("Simple 44 bit signed long", read_buffer, 44, OK, -42);
  // Read the only part of a long (having to fill up 1s)
  read_buffer->curPosByte = 3;
  read_buffer->curPosBit = 4;
  test_plc4c_spi_read_signed_long_args("Simple 36 bit signed long", read_buffer, 36, OK, -42);
  // Read the only part of a long (having to fill up 1s)
  read_buffer->curPosByte = 4;
  read_buffer->curPosBit = 4;
  test_plc4c_spi_read_signed_long_args("Simple 28 bit signed long", read_buffer, 28, OK, -42);
  // Read the only part of a long (having to fill up 1s)
  read_buffer->curPosByte = 5;
  read_buffer->curPosBit = 4;
  test_plc4c_spi_read_signed_long_args("Simple 20 bit signed long", read_buffer, 20, OK, -42);
  // Read the only part of a long (having to fill up 1s)
  read_buffer->curPosByte = 6;
  read_buffer->curPosBit = 4;
  test_plc4c_spi_read_signed_long_args("Simple 12 bit signed long", read_buffer, 12, OK, -42);

  // Read an even shorter part of a long (having to fill up even more 1s)
  read_buffer->curPosByte = 7;
  read_buffer->curPosBit = 1;
  test_plc4c_spi_read_signed_long_args("Simple 7 bit signed long", read_buffer, 7, OK, -42);

  // Read an even shorter part of a long (This time however the value should
  // be positive and hence the higher level bytes should be filled with 0s)
  read_buffer->curPosByte = 7;
  read_buffer->curPosBit = 2;
  test_plc4c_spi_read_signed_long_args("Simple 6 bit signed long", read_buffer, 6, OK, 22);
}

void test_plc4c_spi_read_float_args(char* message,
                                    plc4c_spi_read_buffer* read_buffer, uint8_t num_bits,
                                    plc4c_return_code expected_return_code, float expected_value) {
    printf("Running read_buffer read_float test: %s", message);

    float value = -1;
    plc4c_return_code result =
            plc4c_spi_read_float(read_buffer, num_bits, &value);

    TEST_ASSERT_EQUAL_INT(expected_return_code, result);
    TEST_ASSERT_EQUAL_INT(expected_value, value);

    printf(" -> OK\n");
}

void test_plc4c_spi_read_float(void) {
    // Prepare input data
    uint8_t data[] = {0x40, 0x49, 0x0f, 0xdb};
    plc4c_spi_read_buffer* read_buffer;
    plc4c_spi_read_buffer_create(data, 4, &read_buffer);
    // Run test
    // Read all the full bytes
    test_plc4c_spi_read_float_args("Simple 16 bit float 1", read_buffer, 16, OK, -24);
    test_plc4c_spi_read_float_args("Simple 16 bit float 2", read_buffer, 16, OK, 990);
    // Read the only part of a long (having to fill up 1s)
    read_buffer->curPosByte = 0;
    read_buffer->curPosBit = 0;
    test_plc4c_spi_read_float_args("Simple 32 bit float", read_buffer, 32, OK, 3.14159274);
    // Read the only part of a long (having to fill up 1s)
}

void test_plc4c_spi_read_double_args(char* message,
                                    plc4c_spi_read_buffer* read_buffer, uint8_t num_bits,
                                    plc4c_return_code expected_return_code, double expected_value) {
    printf("Running read_buffer read_double test: %s", message);

    double value = -1;
    plc4c_return_code result =
            plc4c_spi_read_double(read_buffer, num_bits, &value);

    TEST_ASSERT_EQUAL_INT(expected_return_code, result);
    TEST_ASSERT_EQUAL_INT(expected_value, value);

    printf(" -> OK\n");
}

void test_plc4c_spi_read_double(void) {
    // Prepare input data
    uint8_t data[] = {0x40, 0x09, 0x21, 0xfb, 0x54, 0x44, 0x2d, 0x18};
    plc4c_spi_read_buffer* read_buffer;
    plc4c_spi_read_buffer_create(data, 8, &read_buffer);
    test_plc4c_spi_read_double_args("Simple 64 bit float", read_buffer, 64, OK, 3.1415926535897931);
    // Read the only part of a long (having to fill up 1s)
}

void test_plc4c_spi_read_string_args(char* message,
                                     plc4c_spi_read_buffer* read_buffer, uint8_t num_bits, char* encoding,
                                     plc4c_return_code expected_return_code, char* expected_value) {
    printf("Running read_buffer read_string test: %s", message);

    char* value = NULL;
    plc4c_return_code result =
            plc4c_spi_read_string(read_buffer, num_bits, encoding, &value);

    TEST_ASSERT_EQUAL_INT(expected_return_code, result);
    TEST_ASSERT_EQUAL_STRING(expected_value, value);

    printf(" -> OK\n");
}

void test_plc4c_spi_read_string(void) {
    // Prepare input data
    uint8_t data[] = {0x48, 0x75, 0x72, 0x7a};
    plc4c_spi_read_buffer* read_buffer;
    plc4c_spi_read_buffer_create(data, 4, &read_buffer);
    test_plc4c_spi_read_string_args("Simple 32 bit string (4 chars)", read_buffer, 32, "UTF-8", OK, "Hurz");
    // Read the only part of a long (having to fill up 1s)
}

void test_plc4c_spi_read_buffer(void) {
  printf("Test being run in %s\n\n", (plc4c_is_bigendian() ? "Big Endian" : "Little Endian"));

  test_plc4c_spi_read_buffer_create();
  test_plc4c_spi_read_get_total_bytes();
  test_plc4c_spi_read_has_more();
  test_plc4c_spi_read_get_bytes();
  test_plc4c_spi_read_peek_byte();

  test_plc4c_spi_read_read_bit();

  test_plc4c_spi_read_unsigned_byte();
  test_plc4c_spi_read_unsigned_short();
  test_plc4c_spi_read_unsigned_int();
  test_plc4c_spi_read_unsigned_long();

  test_plc4c_spi_read_signed_byte();
  test_plc4c_spi_read_signed_short();
  test_plc4c_spi_read_signed_int();
  test_plc4c_spi_read_signed_long();

  test_plc4c_spi_read_float();
  test_plc4c_spi_read_double();

  // TODO: Commented out ... needs fixing on Windows
  //test_plc4c_spi_read_string();
}