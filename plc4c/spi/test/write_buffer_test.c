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
#include "plc4c/spi/write_buffer.h"

#include <stdlib.h>
#include <unity.h>

void internal_write_buffer_assert_arrays_equal(const uint8_t* expected_array, plc4c_spi_write_buffer* write_buffer, uint8_t num_bytes) {
  for(int i = 0; i < num_bytes; i++) {
    uint8_t expected_value = *(expected_array + i);
    uint8_t actual_value = *(write_buffer->data + i);
    TEST_ASSERT_EQUAL_UINT8_MESSAGE(
        expected_value, actual_value, "Byte arrays differ");
  }
}

void test_plc4c_spi_write_buffer_create_args(
    uint16_t length, plc4c_return_code expected_return_code) {
  printf("Running write_buffer create test with %d, expecting return code %d",
         length, expected_return_code);

  // Create a new write_buffer instance
  plc4c_spi_write_buffer* write_buffer;

  plc4c_return_code return_code =
      plc4c_spi_write_buffer_create(length, &write_buffer);

  TEST_ASSERT_EQUAL_INT(expected_return_code, return_code);
  if (expected_return_code != OK) {
    TEST_ASSERT_NULL(write_buffer);
  } else {
    TEST_ASSERT_NOT_NULL(write_buffer);

    // Check all bytes are 0-initialized
    uint8_t* data = plc4c_spi_write_get_data(write_buffer);
    for(int i = 0; i < length; i++) {
      if(*data != 0) {
        TEST_FAIL_MESSAGE("Data array not 0-initialized");
      }
      data++;
    }

    free(write_buffer);
  }

  printf(" -> OK\n");
}

void test_plc4c_spi_write_buffer_create(void) {
  // Run test
  test_plc4c_spi_write_buffer_create_args(8, OK);
}

void test_plc4c_spi_write_write_bit_args(plc4c_spi_write_buffer* write_buffer,
                                       bool value) {
  printf(
      "Running write_buffer write_bit test. Checking if writing a bit gives the "
      "correct response");

  plc4c_return_code result = plc4c_spi_write_bit(write_buffer, value);

  TEST_ASSERT_EQUAL_INT(OK, result);

  printf(" -> OK\n");
}

void test_plc4c_spi_write_write_bit(void) {
  plc4c_spi_write_buffer* write_buffer;
  plc4c_spi_write_buffer_create(5, &write_buffer);

  // Run test
  test_plc4c_spi_write_write_bit_args(write_buffer, false);
  TEST_ASSERT_EQUAL_UINT8(0, write_buffer->curPosByte);
  TEST_ASSERT_EQUAL_UINT8(1, write_buffer->curPosBit);
  test_plc4c_spi_write_write_bit_args(write_buffer, false);
  TEST_ASSERT_EQUAL_UINT8(0, write_buffer->curPosByte);
  TEST_ASSERT_EQUAL_UINT8(2, write_buffer->curPosBit);
  test_plc4c_spi_write_write_bit_args(write_buffer, false);
  TEST_ASSERT_EQUAL_UINT8(0, write_buffer->curPosByte);
  TEST_ASSERT_EQUAL_UINT8(3, write_buffer->curPosBit);
  test_plc4c_spi_write_write_bit_args(write_buffer, false);
  TEST_ASSERT_EQUAL_UINT8(0, write_buffer->curPosByte);
  TEST_ASSERT_EQUAL_UINT8(4, write_buffer->curPosBit);
  test_plc4c_spi_write_write_bit_args(write_buffer, false);
  TEST_ASSERT_EQUAL_UINT8(0, write_buffer->curPosByte);
  TEST_ASSERT_EQUAL_UINT8(5, write_buffer->curPosBit);
  test_plc4c_spi_write_write_bit_args(write_buffer, false);
  TEST_ASSERT_EQUAL_UINT8(0, write_buffer->curPosByte);
  TEST_ASSERT_EQUAL_UINT8(6, write_buffer->curPosBit);
  test_plc4c_spi_write_write_bit_args(write_buffer, false);
  TEST_ASSERT_EQUAL_UINT8(0, write_buffer->curPosByte);
  TEST_ASSERT_EQUAL_UINT8(7, write_buffer->curPosBit);
  test_plc4c_spi_write_write_bit_args(write_buffer, true);
  TEST_ASSERT_EQUAL_UINT8(1, write_buffer->curPosByte);
  TEST_ASSERT_EQUAL_UINT8(0, write_buffer->curPosBit);

  test_plc4c_spi_write_write_bit_args(write_buffer, false);
  test_plc4c_spi_write_write_bit_args(write_buffer, false);
  test_plc4c_spi_write_write_bit_args(write_buffer, false);
  test_plc4c_spi_write_write_bit_args(write_buffer, false);
  test_plc4c_spi_write_write_bit_args(write_buffer, false);
  test_plc4c_spi_write_write_bit_args(write_buffer, false);
  test_plc4c_spi_write_write_bit_args(write_buffer, true);
  test_plc4c_spi_write_write_bit_args(write_buffer, false);

  test_plc4c_spi_write_write_bit_args(write_buffer, false);
  test_plc4c_spi_write_write_bit_args(write_buffer, false);
  test_plc4c_spi_write_write_bit_args(write_buffer, false);
  test_plc4c_spi_write_write_bit_args(write_buffer, false);
  test_plc4c_spi_write_write_bit_args(write_buffer, false);
  test_plc4c_spi_write_write_bit_args(write_buffer, false);
  test_plc4c_spi_write_write_bit_args(write_buffer, true);
  test_plc4c_spi_write_write_bit_args(write_buffer, true);

  test_plc4c_spi_write_write_bit_args(write_buffer, true);
  test_plc4c_spi_write_write_bit_args(write_buffer, true);
  test_plc4c_spi_write_write_bit_args(write_buffer, true);
  test_plc4c_spi_write_write_bit_args(write_buffer, true);
  test_plc4c_spi_write_write_bit_args(write_buffer, true);
  test_plc4c_spi_write_write_bit_args(write_buffer, true);
  test_plc4c_spi_write_write_bit_args(write_buffer, true);
  test_plc4c_spi_write_write_bit_args(write_buffer, true);

  test_plc4c_spi_write_write_bit_args(write_buffer, true);
  test_plc4c_spi_write_write_bit_args(write_buffer, true);
  test_plc4c_spi_write_write_bit_args(write_buffer, true);
  test_plc4c_spi_write_write_bit_args(write_buffer, true);
  test_plc4c_spi_write_write_bit_args(write_buffer, false);
  test_plc4c_spi_write_write_bit_args(write_buffer, false);
  test_plc4c_spi_write_write_bit_args(write_buffer, false);
  test_plc4c_spi_write_write_bit_args(write_buffer, false);

  // Compare the content of the write buffer with the expected one.
  uint8_t expectedData[] = {1, 2, 3, 255, 240};
  TEST_ASSERT_EQUAL_UINT8(5, write_buffer->curPosByte);
  TEST_ASSERT_EQUAL_UINT8(0, write_buffer->curPosBit);
  internal_write_buffer_assert_arrays_equal((uint8_t*) &expectedData, write_buffer, 5);
}

void test_plc4c_spi_write_unsigned_byte_args(char* message,
                                            plc4c_spi_write_buffer* write_buffer, uint8_t num_bits,
                                            plc4c_return_code expected_return_code, uint8_t value) {
  printf("Running write_buffer write_unsigned_byte test: %s", message);

  plc4c_return_code result =
      plc4c_spi_write_unsigned_byte(write_buffer, num_bits, value);

  TEST_ASSERT_EQUAL_INT(expected_return_code, result);

  printf(" -> OK\n");
}

void test_plc4c_spi_write_unsigned_byte(void) {
  plc4c_spi_write_buffer* write_buffer;
  plc4c_spi_write_buffer_create(8, &write_buffer);
  // Run test
  // Write full bytes
  test_plc4c_spi_write_unsigned_byte_args("Simple full unsigned byte 1", write_buffer, 8, OK, 1);
  test_plc4c_spi_write_unsigned_byte_args("Simple full unsigned byte 2", write_buffer, 8, OK, 2);
  test_plc4c_spi_write_unsigned_byte_args("Simple full unsigned byte 3", write_buffer, 8, OK, 3);
  test_plc4c_spi_write_unsigned_byte_args("Simple full unsigned byte 4", write_buffer, 8, OK, 4);
  test_plc4c_spi_write_unsigned_byte_args("Simple full unsigned byte 5", write_buffer, 8, OK, 5);
  test_plc4c_spi_write_unsigned_byte_args("Simple full unsigned byte 6", write_buffer, 8, OK, 6);
  test_plc4c_spi_write_unsigned_byte_args("Simple full unsigned byte 7", write_buffer, 8, OK, 7);
  test_plc4c_spi_write_unsigned_byte_args("Simple full unsigned byte 8", write_buffer, 8, OK, 8);

  // Check if the array contains the expected data.
  uint8_t expected_data[] = {1, 2, 3, 4, 5, 6, 7, 8};
  TEST_ASSERT_EQUAL_UINT8(8, write_buffer->curPosByte);
  TEST_ASSERT_EQUAL_UINT8(0, write_buffer->curPosBit);
  internal_write_buffer_assert_arrays_equal((uint8_t*) &expected_data, write_buffer, 8);

  // Write a 9th byte (buffer only has 8) (results in error)
  test_plc4c_spi_write_unsigned_byte_args("Exceed write-buffer size", write_buffer, 8, OUT_OF_RANGE, 0);
  plc4c_spi_write_buffer_destroy(write_buffer);

  plc4c_spi_write_buffer_create(2, &write_buffer);
  // Write part of a byte (fits in one byte)
  test_plc4c_spi_write_unsigned_byte_args("Simple 4 bits of unsigned byte", write_buffer, 4, OK, 11);
  TEST_ASSERT_EQUAL_UINT8(0, write_buffer->curPosByte);
  TEST_ASSERT_EQUAL_UINT8(4, write_buffer->curPosBit);
  // Write part of a byte (finishes one byte)
  test_plc4c_spi_write_unsigned_byte_args("Simple 4 bits of unsigned byte, finishing rest of first byte", write_buffer, 4, OK, 10);
  TEST_ASSERT_EQUAL_UINT8(1, write_buffer->curPosByte);
  TEST_ASSERT_EQUAL_UINT8(0, write_buffer->curPosBit);
  test_plc4c_spi_write_unsigned_byte_args("Simple 4 bits of unsigned byte", write_buffer, 4, OK, 7);
  TEST_ASSERT_EQUAL_UINT8(1, write_buffer->curPosByte);
  TEST_ASSERT_EQUAL_UINT8(4, write_buffer->curPosBit);
  test_plc4c_spi_write_unsigned_byte_args("Simple 4 bits of unsigned byte", write_buffer, 4, OK, 5);
  TEST_ASSERT_EQUAL_UINT8(2, write_buffer->curPosByte);
  TEST_ASSERT_EQUAL_UINT8(0, write_buffer->curPosBit);
  uint8_t expected_data_2[] = {186, 117};
  internal_write_buffer_assert_arrays_equal((uint8_t*) &expected_data_2, write_buffer, 2);
  plc4c_spi_write_buffer_destroy(write_buffer);

  // Write part of a byte (spans two bytes)
  plc4c_spi_write_buffer_create(2, &write_buffer);
  write_buffer->curPosByte = 0;
  write_buffer->curPosBit = 5;
  test_plc4c_spi_write_unsigned_byte_args("Simple 7 bits of unsigned byte starting at bit 5 (flowing over to next byte)", write_buffer, 7, OK, 19);
  TEST_ASSERT_EQUAL_UINT8(1, write_buffer->curPosByte);
  TEST_ASSERT_EQUAL_UINT8(4, write_buffer->curPosBit);
  test_plc4c_spi_write_unsigned_byte_args("Simple 4 bits of unsigned byte starting at bit 3", write_buffer, 4, OK, 10);
  TEST_ASSERT_EQUAL_UINT8(2, write_buffer->curPosByte);
  TEST_ASSERT_EQUAL_UINT8(0, write_buffer->curPosBit);
  uint8_t expected_data_3[] = {1, 58};
  internal_write_buffer_assert_arrays_equal((uint8_t*) &expected_data_3, write_buffer, 2);
  plc4c_spi_write_buffer_destroy(write_buffer);

  // Write more than a byte (results in error)
  plc4c_spi_write_buffer_create(2, &write_buffer);
  write_buffer->curPosByte = 0;
  write_buffer->curPosBit = 5;
  test_plc4c_spi_write_unsigned_byte_args("Exceed write-buffer size (Part 2)", write_buffer, 10, OUT_OF_RANGE, 0);
  plc4c_spi_write_buffer_destroy(write_buffer);
}

void test_plc4c_spi_write_unsigned_short_args(char* message,
                                             plc4c_spi_write_buffer* write_buffer, uint8_t num_bits,
                                             plc4c_return_code expected_return_code,
                                              uint16_t value) {
  printf("Running write_buffer write_unsigned_short test: %s", message);

  plc4c_return_code result =
      plc4c_spi_write_unsigned_short(write_buffer, num_bits, value);

  TEST_ASSERT_EQUAL_INT(expected_return_code, result);

  printf(" -> OK\n");
}

void test_plc4c_spi_write_unsigned_short(void) {
  // Prepare input data
  plc4c_spi_write_buffer* write_buffer;
  plc4c_spi_write_buffer_create(8, &write_buffer);
  // Run test
  // Write all the full short
  test_plc4c_spi_write_unsigned_short_args("Simple full short 1", write_buffer, 16, OK, 258);
  TEST_ASSERT_EQUAL_UINT8(2, write_buffer->curPosByte);
  TEST_ASSERT_EQUAL_UINT8(0, write_buffer->curPosBit);
  test_plc4c_spi_write_unsigned_short_args("Simple full short 2", write_buffer, 16, OK, 772);
  TEST_ASSERT_EQUAL_UINT8(4, write_buffer->curPosByte);
  TEST_ASSERT_EQUAL_UINT8(0, write_buffer->curPosBit);
  test_plc4c_spi_write_unsigned_short_args("Simple full short 3", write_buffer, 16, OK, 1286);
  TEST_ASSERT_EQUAL_UINT8(6, write_buffer->curPosByte);
  TEST_ASSERT_EQUAL_UINT8(0, write_buffer->curPosBit);
  test_plc4c_spi_write_unsigned_short_args("Simple full short 4", write_buffer, 16, OK, 1800);
  TEST_ASSERT_EQUAL_UINT8(8, write_buffer->curPosByte);
  TEST_ASSERT_EQUAL_UINT8(0, write_buffer->curPosBit);
  uint8_t expected_data[] = {1, 2, 3, 4, 5, 6, 7, 8};
  internal_write_buffer_assert_arrays_equal((uint8_t*) &expected_data, write_buffer, 8);
  plc4c_spi_write_buffer_destroy(write_buffer);

  // Write a short that spans across 3 bytes
  plc4c_spi_write_buffer_create(3, &write_buffer);
  write_buffer->curPosByte = 0;
  write_buffer->curPosBit = 3;
  test_plc4c_spi_write_unsigned_short_args("Full short starting at bit 3", write_buffer, 14, OK, 517);
  TEST_ASSERT_EQUAL_UINT8(2, write_buffer->curPosByte);
  TEST_ASSERT_EQUAL_UINT8(1, write_buffer->curPosBit);
  uint8_t expected_data_2[] = {1, 2, 128};
  internal_write_buffer_assert_arrays_equal((uint8_t*) &expected_data_2, write_buffer, 3);
  plc4c_spi_write_buffer_destroy(write_buffer);

  // Finish a partial short
  plc4c_spi_write_buffer_create(2, &write_buffer);
  write_buffer->curPosByte = 0;
  write_buffer->curPosBit = 1;
  test_plc4c_spi_write_unsigned_short_args("Short starting at bit 1 but writing full last byte", write_buffer, 15, OK, 258);
  TEST_ASSERT_EQUAL_UINT8(2, write_buffer->curPosByte);
  TEST_ASSERT_EQUAL_UINT8(0, write_buffer->curPosBit);
  uint8_t expected_data_3[] = {1, 2};
  internal_write_buffer_assert_arrays_equal((uint8_t*) &expected_data_3, write_buffer, 2);
  plc4c_spi_write_buffer_destroy(write_buffer);

  // Write only fractions of a short
  plc4c_spi_write_buffer_create(8, &write_buffer);
  write_buffer->curPosByte = 0;
  write_buffer->curPosBit = 0;
  test_plc4c_spi_write_unsigned_short_args("Short only writing 2 byte", write_buffer, 2, OK, 0);
  TEST_ASSERT_EQUAL_UINT8(0, write_buffer->curPosByte);
  TEST_ASSERT_EQUAL_UINT8(2, write_buffer->curPosBit);
  test_plc4c_spi_write_unsigned_short_args("Short only writing 4 byte", write_buffer, 4, OK, 0);
  TEST_ASSERT_EQUAL_UINT8(0, write_buffer->curPosByte);
  TEST_ASSERT_EQUAL_UINT8(6, write_buffer->curPosBit);
  test_plc4c_spi_write_unsigned_short_args("Short only writing 6 byte", write_buffer, 6, OK, 16);
  TEST_ASSERT_EQUAL_UINT8(1, write_buffer->curPosByte);
  TEST_ASSERT_EQUAL_UINT8(4, write_buffer->curPosBit);
  test_plc4c_spi_write_unsigned_short_args("Short only writing 8 byte", write_buffer, 8, OK, 32);
  TEST_ASSERT_EQUAL_UINT8(2, write_buffer->curPosByte);
  TEST_ASSERT_EQUAL_UINT8(4, write_buffer->curPosBit);
  test_plc4c_spi_write_unsigned_short_args("Short only writing 10 byte", write_buffer, 10, OK, 193);
  TEST_ASSERT_EQUAL_UINT8(3, write_buffer->curPosByte);
  TEST_ASSERT_EQUAL_UINT8(6, write_buffer->curPosBit);
  test_plc4c_spi_write_unsigned_short_args("Short only writing 12 byte", write_buffer, 12, OK, 20);
  TEST_ASSERT_EQUAL_UINT8(5, write_buffer->curPosByte);
  TEST_ASSERT_EQUAL_UINT8(2, write_buffer->curPosBit);
  test_plc4c_spi_write_unsigned_short_args("Short only writing 14 byte", write_buffer, 14, OK, 1543);
  TEST_ASSERT_EQUAL_UINT8(7, write_buffer->curPosByte);
  TEST_ASSERT_EQUAL_UINT8(0, write_buffer->curPosBit);
  uint8_t expected_data_4[] = {1, 2, 3, 4, 5, 6};
  internal_write_buffer_assert_arrays_equal((uint8_t*) &expected_data_4, write_buffer, 6);

  // We only have 8 bytes, so with this we would exceed the range.
  test_plc4c_spi_write_unsigned_short_args("Try to write more bytes than the buffer has available", write_buffer, 16, OUT_OF_RANGE, 0);

  plc4c_spi_write_buffer_create(8, &write_buffer);
  write_buffer->curPosByte = 0;
  write_buffer->curPosBit = 0;
  test_plc4c_spi_write_unsigned_short_args("Try to write too many bytes for a short", write_buffer, 18, OUT_OF_RANGE, 0);
}

void test_plc4c_spi_write_unsigned_int_args(char* message,
                                           plc4c_spi_write_buffer* write_buffer, uint8_t num_bits,
                                           plc4c_return_code expected_return_code, uint32_t value) {
  printf("Running write_buffer write_unsigned_int test: %s", message);

  plc4c_return_code result =
      plc4c_spi_write_unsigned_int(write_buffer, num_bits, value);

  TEST_ASSERT_EQUAL_INT(expected_return_code, result);

  printf(" -> OK\n");
}

void test_plc4c_spi_write_unsigned_int(void) {
  // Write all the full short
  plc4c_spi_write_buffer* write_buffer;
  plc4c_spi_write_buffer_create(8, &write_buffer);
  test_plc4c_spi_write_unsigned_int_args("Simple full int 1", write_buffer, 32, OK, 16909060);
  TEST_ASSERT_EQUAL_UINT8(4, write_buffer->curPosByte);
  TEST_ASSERT_EQUAL_UINT8(0, write_buffer->curPosBit);
  test_plc4c_spi_write_unsigned_int_args("Simple full int 1", write_buffer, 32, OK, 84281096);
  TEST_ASSERT_EQUAL_UINT8(8, write_buffer->curPosByte);
  TEST_ASSERT_EQUAL_UINT8(0, write_buffer->curPosBit);
  uint8_t expected_data[] = {1, 2, 3, 4, 5, 6, 7, 8};
  internal_write_buffer_assert_arrays_equal((uint8_t*) &expected_data, write_buffer, 8);
  plc4c_spi_write_buffer_destroy(write_buffer);

  // Write a full int starting somewhere in between.
  plc4c_spi_write_buffer_create(8, &write_buffer);
  write_buffer->curPosByte = 0;
  write_buffer->curPosBit = 3;
  test_plc4c_spi_write_unsigned_int_args("Full int starting at bit 3", write_buffer, 32, OK, 135272480);
  TEST_ASSERT_EQUAL_UINT8(4, write_buffer->curPosByte);
  TEST_ASSERT_EQUAL_UINT8(3, write_buffer->curPosBit);
  uint8_t expected_data_2[] = {1, 2, 3, 4, 0};
  internal_write_buffer_assert_arrays_equal((uint8_t*) &expected_data_2, write_buffer, 5);
  plc4c_spi_write_buffer_destroy(write_buffer);

  // Write an int starting somewhere in between but ending at a full byte.
  plc4c_spi_write_buffer_create(8, &write_buffer);
  write_buffer->curPosByte = 0;
  write_buffer->curPosBit = 3;
  test_plc4c_spi_write_unsigned_int_args("Full int starting at bit 3", write_buffer, 29, OK, 16909060);
  TEST_ASSERT_EQUAL_UINT8(4, write_buffer->curPosByte);
  TEST_ASSERT_EQUAL_UINT8(0, write_buffer->curPosBit);
  internal_write_buffer_assert_arrays_equal((uint8_t*) &expected_data, write_buffer, 4);
  plc4c_spi_write_buffer_destroy(write_buffer);

  // Write a full int starting somewhere in between.
  plc4c_spi_write_buffer_create(8, &write_buffer);
  write_buffer->curPosByte = 0;
  write_buffer->curPosBit = 3;
  test_plc4c_spi_write_unsigned_int_args("Int with only 5 bit", write_buffer, 5, OK, 1);
  TEST_ASSERT_EQUAL_UINT8(1, write_buffer->curPosByte);
  TEST_ASSERT_EQUAL_UINT8(0, write_buffer->curPosBit);
  test_plc4c_spi_write_unsigned_int_args("Int with only 7 bit", write_buffer, 7, OK, 1);
  TEST_ASSERT_EQUAL_UINT8(1, write_buffer->curPosByte);
  TEST_ASSERT_EQUAL_UINT8(7, write_buffer->curPosBit);
  test_plc4c_spi_write_unsigned_int_args("Int with only 9 bit", write_buffer, 9, OK, 3);
  TEST_ASSERT_EQUAL_UINT8(3, write_buffer->curPosByte);
  TEST_ASSERT_EQUAL_UINT8(0, write_buffer->curPosBit);
  test_plc4c_spi_write_unsigned_int_args("Int with only 15 bit", write_buffer, 15, OK, 514);
  TEST_ASSERT_EQUAL_UINT8(4, write_buffer->curPosByte);
  TEST_ASSERT_EQUAL_UINT8(7, write_buffer->curPosBit);
  test_plc4c_spi_write_unsigned_int_args("Int with only 16 bit", write_buffer, 16, OK, 33539);
  TEST_ASSERT_EQUAL_UINT8(6, write_buffer->curPosByte);
  TEST_ASSERT_EQUAL_UINT8(7, write_buffer->curPosBit);
  uint8_t expected_data_3[] = {1, 2, 3, 4, 5, 6, 6};
  internal_write_buffer_assert_arrays_equal((uint8_t*) &expected_data_3, write_buffer, 7);
  plc4c_spi_write_buffer_destroy(write_buffer);

  plc4c_spi_write_buffer_create(8, &write_buffer);
  write_buffer->curPosByte = 0;
  write_buffer->curPosBit = 3;
  test_plc4c_spi_write_unsigned_int_args("Int with only 17 bit", write_buffer, 17, OK, 4128);
  TEST_ASSERT_EQUAL_UINT8(2, write_buffer->curPosByte);
  TEST_ASSERT_EQUAL_UINT8(4, write_buffer->curPosBit);
  test_plc4c_spi_write_unsigned_int_args("Int with only 26 bit", write_buffer, 26, OK, 12648769);
  TEST_ASSERT_EQUAL_UINT8(5, write_buffer->curPosByte);
  TEST_ASSERT_EQUAL_UINT8(6, write_buffer->curPosBit);
  uint8_t expected_data_4[] = {1, 2, 3, 4, 5, 4};
  internal_write_buffer_assert_arrays_equal((uint8_t*) &expected_data_4, write_buffer, 5);
  plc4c_spi_write_buffer_destroy(write_buffer);

  plc4c_spi_write_buffer_create(8, &write_buffer);
  write_buffer->curPosByte = 0;
  write_buffer->curPosBit = 3;
  test_plc4c_spi_write_unsigned_int_args("Int with only 30 bit", write_buffer, 30, OK, 33818120);
  TEST_ASSERT_EQUAL_UINT8(4, write_buffer->curPosByte);
  TEST_ASSERT_EQUAL_UINT8(1, write_buffer->curPosBit);
  uint8_t expected_data_5[] = {1, 2, 3, 4, 0};
  internal_write_buffer_assert_arrays_equal((uint8_t*) &expected_data_5, write_buffer, 4);
  plc4c_spi_write_buffer_destroy(write_buffer);

  plc4c_spi_write_buffer_create(8, &write_buffer);
  write_buffer->curPosByte = 0;
  write_buffer->curPosBit = 3;
  test_plc4c_spi_write_unsigned_int_args("Int with only 31 bit", write_buffer, 31, OK, 67636240);
  TEST_ASSERT_EQUAL_UINT8(4, write_buffer->curPosByte);
  TEST_ASSERT_EQUAL_UINT8(2, write_buffer->curPosBit);
  internal_write_buffer_assert_arrays_equal((uint8_t*) &expected_data_5, write_buffer, 4);
  plc4c_spi_write_buffer_destroy(write_buffer);
}

void test_plc4c_spi_write_unsigned_long_args(char* message,
                                            plc4c_spi_write_buffer* write_buffer, uint8_t num_bits,
                                            plc4c_return_code expected_return_code, uint64_t value) {
  printf("Running write_buffer write_unsigned_long test: %s", message);

  plc4c_return_code result =
      plc4c_spi_write_unsigned_long(write_buffer, num_bits, value);

  TEST_ASSERT_EQUAL_INT(expected_return_code, result);

  printf(" -> OK\n");
}

void test_plc4c_spi_write_unsigned_long(void) {
  plc4c_spi_write_buffer* write_buffer;
  plc4c_spi_write_buffer_create(8, &write_buffer);
  // Write all the full long
  test_plc4c_spi_write_unsigned_long_args("Simple full long", write_buffer, 64, OK, 72623859790382856);
  TEST_ASSERT_EQUAL_UINT8(8, write_buffer->curPosByte);
  TEST_ASSERT_EQUAL_UINT8(0, write_buffer->curPosBit);
  uint8_t expected_data[] = {1, 2, 3, 4, 5, 6, 7, 8};
  internal_write_buffer_assert_arrays_equal((uint8_t*) &expected_data, write_buffer, 8);
  plc4c_spi_write_buffer_destroy(write_buffer);

  // Write a full long starting somewhere in between.
  plc4c_spi_write_buffer_create(9, &write_buffer);
  write_buffer->curPosByte = 0;
  write_buffer->curPosBit = 3;
  test_plc4c_spi_write_unsigned_long_args("Full long starting at bit 3", write_buffer, 64, OK, 580990878323062848);
  TEST_ASSERT_EQUAL_UINT8(8, write_buffer->curPosByte);
  TEST_ASSERT_EQUAL_UINT8(3, write_buffer->curPosBit);
  uint8_t expected_data_2[] = {1, 2, 3, 4, 5, 6, 7, 8, 0};
  internal_write_buffer_assert_arrays_equal((uint8_t*) &expected_data_2, write_buffer, 9);
  plc4c_spi_write_buffer_destroy(write_buffer);

  // Write a long starting somewhere in between but ending at a full byte.
  plc4c_spi_write_buffer_create(9, &write_buffer);
  write_buffer->curPosByte = 0;
  write_buffer->curPosBit = 3;
  test_plc4c_spi_write_unsigned_long_args("Full long starting at bit 3", write_buffer, 61, OK, 72623859790382856);
  TEST_ASSERT_EQUAL_UINT8(8, write_buffer->curPosByte);
  TEST_ASSERT_EQUAL_UINT8(0, write_buffer->curPosBit);
  internal_write_buffer_assert_arrays_equal((uint8_t*) &expected_data, write_buffer, 8);
  plc4c_spi_write_buffer_destroy(write_buffer);
}

void test_plc4c_spi_write_signed_byte_args(char* message,
                                          plc4c_spi_write_buffer* write_buffer, uint8_t num_bits,
                                          plc4c_return_code expected_return_code, int8_t value) {
  printf("Running write_buffer write_signed_byte test: %s", message);

  plc4c_return_code result =
      plc4c_spi_write_signed_byte(write_buffer, num_bits, value);

  TEST_ASSERT_EQUAL_INT(expected_return_code, result);

  printf(" -> OK\n");
}

void test_plc4c_spi_write_signed_byte(void) {
  plc4c_spi_write_buffer* write_buffer;
  plc4c_spi_write_buffer_create(8, &write_buffer);
  // Write all the full bytes
  test_plc4c_spi_write_signed_byte_args("Simple full signed byte 1", write_buffer, 8, OK, 1);
  TEST_ASSERT_EQUAL_UINT8(1, write_buffer->curPosByte);
  TEST_ASSERT_EQUAL_UINT8(0, write_buffer->curPosBit);
  test_plc4c_spi_write_signed_byte_args("Simple full signed byte 2", write_buffer, 8, OK, 2);
  TEST_ASSERT_EQUAL_UINT8(2, write_buffer->curPosByte);
  TEST_ASSERT_EQUAL_UINT8(0, write_buffer->curPosBit);
  test_plc4c_spi_write_signed_byte_args("Simple full signed byte 3", write_buffer, 8, OK, 3);
  TEST_ASSERT_EQUAL_UINT8(3, write_buffer->curPosByte);
  TEST_ASSERT_EQUAL_UINT8(0, write_buffer->curPosBit);
  test_plc4c_spi_write_signed_byte_args("Simple full signed byte 4", write_buffer, 8, OK, 4);
  TEST_ASSERT_EQUAL_UINT8(4, write_buffer->curPosByte);
  TEST_ASSERT_EQUAL_UINT8(0, write_buffer->curPosBit);
  test_plc4c_spi_write_signed_byte_args("Simple full signed byte 5", write_buffer, 8, OK, 5);
  TEST_ASSERT_EQUAL_UINT8(5, write_buffer->curPosByte);
  TEST_ASSERT_EQUAL_UINT8(0, write_buffer->curPosBit);
  test_plc4c_spi_write_signed_byte_args("Simple full signed byte 6", write_buffer, 8, OK, 6);
  TEST_ASSERT_EQUAL_UINT8(6, write_buffer->curPosByte);
  TEST_ASSERT_EQUAL_UINT8(0, write_buffer->curPosBit);
  test_plc4c_spi_write_signed_byte_args("Simple full signed byte 7", write_buffer, 8, OK, 7);
  TEST_ASSERT_EQUAL_UINT8(7, write_buffer->curPosByte);
  TEST_ASSERT_EQUAL_UINT8(0, write_buffer->curPosBit);
  test_plc4c_spi_write_signed_byte_args("Simple full signed byte 8", write_buffer, 8, OK, 8);
  TEST_ASSERT_EQUAL_UINT8(8, write_buffer->curPosByte);
  TEST_ASSERT_EQUAL_UINT8(0, write_buffer->curPosBit);
  uint8_t expected_data[] = {1, 2, 3, 4, 5, 6, 7, 8};
  internal_write_buffer_assert_arrays_equal((uint8_t*) &expected_data, write_buffer, 8);

  // Write a 9th byte (buffer only has 8) (results in error)
  test_plc4c_spi_write_signed_byte_args("Exceed write-buffer size", write_buffer, 8, OUT_OF_RANGE, -1);
  plc4c_spi_write_buffer_destroy(write_buffer);

  plc4c_spi_write_buffer_create(8, &write_buffer);
  // write part of a byte (fits in one byte)
  test_plc4c_spi_write_signed_byte_args("Simple 4 bits of signed byte", write_buffer, 4, OK, -5);
  TEST_ASSERT_EQUAL_UINT8(0, write_buffer->curPosByte);
  TEST_ASSERT_EQUAL_UINT8(4, write_buffer->curPosBit);
  // write part of a byte (finishes one byte)
  test_plc4c_spi_write_signed_byte_args("Simple 4 bits of signed byte, finishing rest of first byte", write_buffer, 4, OK, -6);
  TEST_ASSERT_EQUAL_UINT8(1, write_buffer->curPosByte);
  TEST_ASSERT_EQUAL_UINT8(0, write_buffer->curPosBit);
  test_plc4c_spi_write_signed_byte_args("Simple 4 bits of signed byte", write_buffer, 4, OK, 7);
  TEST_ASSERT_EQUAL_UINT8(1, write_buffer->curPosByte);
  TEST_ASSERT_EQUAL_UINT8(4, write_buffer->curPosBit);
  test_plc4c_spi_write_signed_byte_args("Simple 4 bits of signed byte", write_buffer, 4, OK, 5);
  TEST_ASSERT_EQUAL_UINT8(2, write_buffer->curPosByte);
  TEST_ASSERT_EQUAL_UINT8(0, write_buffer->curPosBit);
  uint8_t expected_data_2[] = {186, 117};
  internal_write_buffer_assert_arrays_equal((uint8_t*) &expected_data_2, write_buffer, 2);
  plc4c_spi_write_buffer_destroy(write_buffer);

  // write part of a byte (spans two bytes)
  plc4c_spi_write_buffer_create(8, &write_buffer);
  write_buffer->curPosByte = 0;
  write_buffer->curPosBit = 5;
  test_plc4c_spi_write_signed_byte_args("Simple 6 bits of signed byte starting at bit 5 (flowing over to next byte)", write_buffer, 6, OK, 19);
  TEST_ASSERT_EQUAL_UINT8(1, write_buffer->curPosByte);
  TEST_ASSERT_EQUAL_UINT8(3, write_buffer->curPosBit);
  test_plc4c_spi_write_signed_byte_args("Simple 4 bits of signed byte starting at bit 3", write_buffer, 4, OK, -6);
  TEST_ASSERT_EQUAL_UINT8(1, write_buffer->curPosByte);
  TEST_ASSERT_EQUAL_UINT8(7, write_buffer->curPosBit);
  uint8_t expected_data_3[] = {2, 116};
  internal_write_buffer_assert_arrays_equal((uint8_t*) &expected_data_3, write_buffer, 2);
  plc4c_spi_write_buffer_destroy(write_buffer);

  // write more than a byte (results in error)
  plc4c_spi_write_buffer_create(8, &write_buffer);
  test_plc4c_spi_write_signed_byte_args("Exceed write-buffer size (Part 2)", write_buffer, 10, OUT_OF_RANGE, -1);
  plc4c_spi_write_buffer_destroy(write_buffer);
}

void test_plc4c_spi_write_signed_short_args(char* message,
                                           plc4c_spi_write_buffer* write_buffer, uint8_t num_bits,
                                           plc4c_return_code expected_return_code, int16_t value) {
  printf("Running write_buffer write_signed_byte test: %s", message);

  plc4c_return_code result =
      plc4c_spi_write_signed_short(write_buffer, num_bits, value);

  TEST_ASSERT_EQUAL_INT(expected_return_code, result);

  printf(" -> OK\n");
}

void test_plc4c_spi_write_signed_short(void) {
  plc4c_spi_write_buffer* write_buffer;
  plc4c_spi_write_buffer_create(2, &write_buffer);
  // write all the full bytes
  test_plc4c_spi_write_signed_short_args("Simple full signed short 1", write_buffer, 16, OK, -42);
  TEST_ASSERT_EQUAL_UINT8(2, write_buffer->curPosByte);
  TEST_ASSERT_EQUAL_UINT8(0, write_buffer->curPosBit);
  uint8_t expected_data[] = {255, 214};
  internal_write_buffer_assert_arrays_equal((uint8_t*) &expected_data, write_buffer, 2);
  plc4c_spi_write_buffer_destroy(write_buffer);

  // write the only part of a short (having to fill up 1s)
  plc4c_spi_write_buffer_create(2, &write_buffer);
  write_buffer->curPosByte = 0;
  write_buffer->curPosBit = 4;
  test_plc4c_spi_write_signed_short_args("Simple 12 bit signed short", write_buffer, 12, OK, -42);
  TEST_ASSERT_EQUAL_UINT8(2, write_buffer->curPosByte);
  TEST_ASSERT_EQUAL_UINT8(0, write_buffer->curPosBit);
  uint8_t expected_data_2[] = {15, 214};
  internal_write_buffer_assert_arrays_equal((uint8_t*) &expected_data_2, write_buffer, 2);
  plc4c_spi_write_buffer_destroy(write_buffer);

  // write an even shorter part of a short (having to fill up even more 1s)
  plc4c_spi_write_buffer_create(2, &write_buffer);
  write_buffer->curPosByte = 1;
  write_buffer->curPosBit = 1;
  test_plc4c_spi_write_signed_short_args("Simple 7 bit signed short", write_buffer, 7, OK, -42);
  TEST_ASSERT_EQUAL_UINT8(2, write_buffer->curPosByte);
  TEST_ASSERT_EQUAL_UINT8(0, write_buffer->curPosBit);
  uint8_t expected_data_3[] = {0, 86};
  internal_write_buffer_assert_arrays_equal((uint8_t*) &expected_data_3, write_buffer, 2);
  plc4c_spi_write_buffer_destroy(write_buffer);

  // write an even shorter part of a short (This time however the value should
  // be positive and hence the higher level byte should be filled with 0s)
  plc4c_spi_write_buffer_create(2, &write_buffer);
  write_buffer->curPosByte = 1;
  write_buffer->curPosBit = 2;
  test_plc4c_spi_write_signed_short_args("Simple 6 bit signed short", write_buffer, 6, OK, 22);
  TEST_ASSERT_EQUAL_UINT8(2, write_buffer->curPosByte);
  TEST_ASSERT_EQUAL_UINT8(0, write_buffer->curPosBit);
  uint8_t expected_data_4[] = {0, 22};
  internal_write_buffer_assert_arrays_equal((uint8_t*) &expected_data_4, write_buffer, 2);
  plc4c_spi_write_buffer_destroy(write_buffer);
}

void test_plc4c_spi_write_signed_int_args(char* message,
                                         plc4c_spi_write_buffer* write_buffer, uint8_t num_bits,
                                         plc4c_return_code expected_return_code, int32_t value) {
  printf("Running write_buffer write_signed_byte test: %s", message);

  plc4c_return_code result =
      plc4c_spi_write_signed_int(write_buffer, num_bits, value);

  TEST_ASSERT_EQUAL_INT(expected_return_code, result);

  printf(" -> OK\n");
}

void test_plc4c_spi_write_signed_int(void) {
  plc4c_spi_write_buffer* write_buffer;
  plc4c_spi_write_buffer_create(5, &write_buffer);
  // write all the full bytes
  test_plc4c_spi_write_signed_int_args("Simple full signed int", write_buffer, 32, OK, -42);
  TEST_ASSERT_EQUAL_UINT8(4, write_buffer->curPosByte);
  TEST_ASSERT_EQUAL_UINT8(0, write_buffer->curPosBit);
  uint8_t expected_data[] = {255, 255, 255, 214};
  internal_write_buffer_assert_arrays_equal((uint8_t*) &expected_data, write_buffer, 4);
  plc4c_spi_write_buffer_destroy(write_buffer);

  // write the only part of a int (having to fill up 1s)
  plc4c_spi_write_buffer_create(5, &write_buffer);
  write_buffer->curPosByte = 0;
  write_buffer->curPosBit = 4;
  test_plc4c_spi_write_signed_int_args("Simple 28 bit signed int", write_buffer, 28, OK, -42);
  TEST_ASSERT_EQUAL_UINT8(4, write_buffer->curPosByte);
  TEST_ASSERT_EQUAL_UINT8(0, write_buffer->curPosBit);
  uint8_t expected_data_2[] = {15, 255, 255, 214};
  internal_write_buffer_assert_arrays_equal((uint8_t*) &expected_data_2, write_buffer, 4);
  plc4c_spi_write_buffer_destroy(write_buffer);

  // write the only part of a int (having to fill up 1s)
  plc4c_spi_write_buffer_create(5, &write_buffer);
  write_buffer->curPosByte = 1;
  write_buffer->curPosBit = 4;
  test_plc4c_spi_write_signed_int_args("Simple 20 bit signed int", write_buffer, 20, OK, -42);
  TEST_ASSERT_EQUAL_UINT8(4, write_buffer->curPosByte);
  TEST_ASSERT_EQUAL_UINT8(0, write_buffer->curPosBit);
  uint8_t expected_data_3[] = {0, 15, 255, 214};
  internal_write_buffer_assert_arrays_equal((uint8_t*) &expected_data_3, write_buffer, 4);
  plc4c_spi_write_buffer_destroy(write_buffer);

  // write the only part of a int (having to fill up 1s)
  plc4c_spi_write_buffer_create(5, &write_buffer);
  write_buffer->curPosByte = 2;
  write_buffer->curPosBit = 4;
  test_plc4c_spi_write_signed_int_args("Simple 12 bit signed int", write_buffer, 12, OK, -42);
  TEST_ASSERT_EQUAL_UINT8(4, write_buffer->curPosByte);
  TEST_ASSERT_EQUAL_UINT8(0, write_buffer->curPosBit);
  uint8_t expected_data_4[] = {0, 0, 15, 214};
  internal_write_buffer_assert_arrays_equal((uint8_t*) &expected_data_4, write_buffer, 4);
  plc4c_spi_write_buffer_destroy(write_buffer);

  // write an even shorter part of a int (having to fill up even more 1s)
  plc4c_spi_write_buffer_create(5, &write_buffer);
  write_buffer->curPosByte = 3;
  write_buffer->curPosBit = 1;
  test_plc4c_spi_write_signed_int_args("Simple 7 bit signed int", write_buffer, 7, OK, -42);
  TEST_ASSERT_EQUAL_UINT8(4, write_buffer->curPosByte);
  TEST_ASSERT_EQUAL_UINT8(0, write_buffer->curPosBit);
  uint8_t expected_data_5[] = {0, 0, 0, 86};
  internal_write_buffer_assert_arrays_equal((uint8_t*) &expected_data_5, write_buffer, 4);
  plc4c_spi_write_buffer_destroy(write_buffer);

  // write an even shorter part of a int (This time however the value should
  // be positive and hence the higher level byte should be filled with 0s)
  plc4c_spi_write_buffer_create(5, &write_buffer);
  write_buffer->curPosByte = 3;
  write_buffer->curPosBit = 2;
  test_plc4c_spi_write_signed_int_args("Simple 6 bit signed int", write_buffer, 6, OK, 22);
  TEST_ASSERT_EQUAL_UINT8(4, write_buffer->curPosByte);
  TEST_ASSERT_EQUAL_UINT8(0, write_buffer->curPosBit);
  uint8_t expected_data_6[] = {0, 0, 0, 22};
  internal_write_buffer_assert_arrays_equal((uint8_t*) &expected_data_6, write_buffer, 4);
  plc4c_spi_write_buffer_destroy(write_buffer);
}

void test_plc4c_spi_write_signed_long_args(char* message,
                                          plc4c_spi_write_buffer* write_buffer, uint8_t num_bits,
                                          plc4c_return_code expected_return_code, int64_t value) {
  printf("Running write_buffer write_signed_byte test: %s", message);

  plc4c_return_code result =
      plc4c_spi_write_signed_long(write_buffer, num_bits, value);

  TEST_ASSERT_EQUAL_INT(expected_return_code, result);

  printf(" -> OK\n");
}

void test_plc4c_spi_write_signed_long(void) {
  plc4c_spi_write_buffer* write_buffer;
  // write all the full bytes
  plc4c_spi_write_buffer_create(8, &write_buffer);
  test_plc4c_spi_write_signed_long_args("Simple full signed long", write_buffer, 64, OK, -42);
  TEST_ASSERT_EQUAL_UINT8(8, write_buffer->curPosByte);
  TEST_ASSERT_EQUAL_UINT8(0, write_buffer->curPosBit);
  uint8_t expected_data[] = {255, 255, 255, 255, 255, 255, 255, 214};
  internal_write_buffer_assert_arrays_equal((uint8_t*) &expected_data, write_buffer, 8);
  plc4c_spi_write_buffer_destroy(write_buffer);

  // write the only part of a long (having to fill up 1s)
  plc4c_spi_write_buffer_create(8, &write_buffer);
  write_buffer->curPosByte = 0;
  write_buffer->curPosBit = 4;
  test_plc4c_spi_write_signed_long_args("Simple 60 bit signed long", write_buffer, 60, OK, -42);
  uint8_t expected_data_2[] = {15, 255, 255, 255, 255, 255, 255, 214};
  internal_write_buffer_assert_arrays_equal((uint8_t*) &expected_data_2, write_buffer, 8);
  plc4c_spi_write_buffer_destroy(write_buffer);

  // write the only part of a long (having to fill up 1s)
  plc4c_spi_write_buffer_create(8, &write_buffer);
  write_buffer->curPosByte = 1;
  write_buffer->curPosBit = 4;
  test_plc4c_spi_write_signed_long_args("Simple 52 bit signed long", write_buffer, 52, OK, -42);
  uint8_t expected_data_3[] = {0, 15, 255, 255, 255, 255, 255, 214};
  internal_write_buffer_assert_arrays_equal((uint8_t*) &expected_data_3, write_buffer, 8);
  plc4c_spi_write_buffer_destroy(write_buffer);

  // write the only part of a long (having to fill up 1s)
  plc4c_spi_write_buffer_create(8, &write_buffer);
  write_buffer->curPosByte = 2;
  write_buffer->curPosBit = 4;
  test_plc4c_spi_write_signed_long_args("Simple 44 bit signed long", write_buffer, 44, OK, -42);
  uint8_t expected_data_4[] = {0, 0, 15, 255, 255, 255, 255, 214};
  internal_write_buffer_assert_arrays_equal((uint8_t*) &expected_data_4, write_buffer, 8);
  plc4c_spi_write_buffer_destroy(write_buffer);

  // write the only part of a long (having to fill up 1s)
  plc4c_spi_write_buffer_create(8, &write_buffer);
  write_buffer->curPosByte = 3;
  write_buffer->curPosBit = 4;
  test_plc4c_spi_write_signed_long_args("Simple 36 bit signed long", write_buffer, 36, OK, -42);
  uint8_t expected_data_5[] = {0, 0, 0, 15, 255, 255, 255, 214};
  internal_write_buffer_assert_arrays_equal((uint8_t*) &expected_data_5, write_buffer, 8);
  plc4c_spi_write_buffer_destroy(write_buffer);

  // write the only part of a long (having to fill up 1s)
  plc4c_spi_write_buffer_create(8, &write_buffer);
  write_buffer->curPosByte = 4;
  write_buffer->curPosBit = 4;
  test_plc4c_spi_write_signed_long_args("Simple 28 bit signed long", write_buffer, 28, OK, -42);
  uint8_t expected_data_6[] = {0, 0, 0, 0, 15, 255, 255, 214};
  internal_write_buffer_assert_arrays_equal((uint8_t*) &expected_data_6, write_buffer, 8);
  plc4c_spi_write_buffer_destroy(write_buffer);

  // write the only part of a long (having to fill up 1s)
  plc4c_spi_write_buffer_create(8, &write_buffer);
  write_buffer->curPosByte = 5;
  write_buffer->curPosBit = 4;
  test_plc4c_spi_write_signed_long_args("Simple 20 bit signed long", write_buffer, 20, OK, -42);
  uint8_t expected_data_7[] = {0, 0, 0, 0, 0, 15, 255, 214};
  internal_write_buffer_assert_arrays_equal((uint8_t*) &expected_data_7, write_buffer, 8);
  plc4c_spi_write_buffer_destroy(write_buffer);

  // write the only part of a long (having to fill up 1s)
  plc4c_spi_write_buffer_create(8, &write_buffer);
  write_buffer->curPosByte = 6;
  write_buffer->curPosBit = 4;
  test_plc4c_spi_write_signed_long_args("Simple 12 bit signed long", write_buffer, 12, OK, -42);
  uint8_t expected_data_8[] = {0, 0, 0, 0, 0, 0, 15, 214};
  internal_write_buffer_assert_arrays_equal((uint8_t*) &expected_data_8, write_buffer, 8);
  plc4c_spi_write_buffer_destroy(write_buffer);

  // write an even shorter part of a long (having to fill up even more 1s)
  plc4c_spi_write_buffer_create(8, &write_buffer);
  write_buffer->curPosByte = 7;
  write_buffer->curPosBit = 1;
  test_plc4c_spi_write_signed_long_args("Simple 7 bit signed long", write_buffer, 7, OK, -42);
  uint8_t expected_data_9[] = {0, 0, 0, 0, 0, 0, 0, 86};
  internal_write_buffer_assert_arrays_equal((uint8_t*) &expected_data_9, write_buffer, 8);
  plc4c_spi_write_buffer_destroy(write_buffer);

  // write an even shorter part of a long (This time however the value should
  // be positive and hence the higher level bytes should be filled with 0s)
  plc4c_spi_write_buffer_create(8, &write_buffer);
  write_buffer->curPosByte = 7;
  write_buffer->curPosBit = 2;
  test_plc4c_spi_write_signed_long_args("Simple 6 bit signed long", write_buffer, 6, OK, 22);
  uint8_t expected_data_10[] = {0, 0, 0, 0, 0, 0, 0, 22};
  internal_write_buffer_assert_arrays_equal((uint8_t*) &expected_data_10, write_buffer, 8);
  plc4c_spi_write_buffer_destroy(write_buffer);
}

void test_plc4c_spi_write_float_args(char* message,
                                    plc4c_spi_write_buffer* write_buffer, uint8_t num_bits,
                                    plc4c_return_code expected_return_code, float value) {
  printf("Running write_buffer write_float test: %s", message);

  plc4c_return_code result =
      plc4c_spi_write_float(write_buffer, num_bits, value);

  TEST_ASSERT_EQUAL_INT(expected_return_code, result);

  printf(" -> OK\n");
}

void test_plc4c_spi_write_float(void) {
  // Prepare input data
  plc4c_spi_write_buffer* write_buffer;
  plc4c_spi_write_buffer_create(4, &write_buffer);
  test_plc4c_spi_write_float_args("Simple 32 bit float", write_buffer, 32, OK, (float) 3.14159274);
  uint8_t expected_data[] = {0x40, 0x49, 0x0f, 0xdb};
  internal_write_buffer_assert_arrays_equal((uint8_t*) &expected_data, write_buffer, 4);
  plc4c_spi_write_buffer_destroy(write_buffer);
}

void test_plc4c_spi_write_double_args(char* message,
                                     plc4c_spi_write_buffer* write_buffer, uint8_t num_bits,
                                     plc4c_return_code expected_return_code, double value) {
  printf("Running write_buffer write_double test: %s", message);

  plc4c_return_code result =
      plc4c_spi_write_double(write_buffer, num_bits, value);

  TEST_ASSERT_EQUAL_INT(expected_return_code, result);

  printf(" -> OK\n");
}

void test_plc4c_spi_write_double(void) {
  // Prepare input data
  plc4c_spi_write_buffer* write_buffer;
  plc4c_spi_write_buffer_create(8, &write_buffer);
  test_plc4c_spi_write_double_args("Simple 64 bit float", write_buffer, 64, OK, 3.1415926535897931);
  uint8_t expected_data[] = {0x40, 0x09, 0x21, 0xfb, 0x54, 0x44, 0x2d, 0x18};
  internal_write_buffer_assert_arrays_equal((uint8_t*) &expected_data, write_buffer, 8);
  plc4c_spi_write_buffer_destroy(write_buffer);
}

void test_plc4c_spi_write_string_args(char* message,
                                     plc4c_spi_write_buffer* write_buffer, uint8_t num_bits, char* encoding,
                                     plc4c_return_code expected_return_code, char* value) {
  printf("Running write_buffer write_string test: %s", message);

  plc4c_return_code result =
      plc4c_spi_write_string(write_buffer, num_bits, encoding, value);

  TEST_ASSERT_EQUAL_INT(expected_return_code, result);

  printf(" -> OK\n");
}

void test_plc4c_spi_write_string(void) {
  // Prepare input data
  plc4c_spi_write_buffer* write_buffer;
  plc4c_spi_write_buffer_create(32, &write_buffer);
  test_plc4c_spi_write_string_args("Simple 32 bit string (4 chars)", write_buffer, 32, "UTF-8", OK, "Hurz");
  uint8_t expected_data[] = {0x48, 0x75, 0x72, 0x7a};
  internal_write_buffer_assert_arrays_equal((uint8_t*) &expected_data, write_buffer, 4);
  plc4c_spi_write_buffer_destroy(write_buffer);
}

void test_plc4c_spi_write_buffer(void) {
  test_plc4c_spi_write_buffer_create();

  test_plc4c_spi_write_write_bit();

  test_plc4c_spi_write_unsigned_byte();
  test_plc4c_spi_write_unsigned_short();
  test_plc4c_spi_write_unsigned_int();
  test_plc4c_spi_write_unsigned_long();

  test_plc4c_spi_write_signed_byte();
  test_plc4c_spi_write_signed_short();
  test_plc4c_spi_write_signed_int();
  test_plc4c_spi_write_signed_long();

  test_plc4c_spi_write_float();
  test_plc4c_spi_write_double();

  test_plc4c_spi_write_string();
}