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
#include <unity.h>

#include "plc4c/spi/read_buffer.h"
#include "tpkt_packet.h"

void internal_assert_arrays_equal(const uint8_t* expected_array, plc4c_spi_write_buffer* write_buffer, uint8_t num_bytes) {
  for(int i = 0; i < num_bytes; i++) {
    uint8_t expected_value = *(expected_array + i);
    uint8_t actual_value = *(write_buffer->data + i);
    TEST_ASSERT_EQUAL_UINT8_MESSAGE(
        expected_value, actual_value, "Byte arrays differ");
  }
}

void parse_cotp_connection_request() {
  uint8_t cotp_connection_request_bytes[] = {
    0x03,0x00,0x00,0x16,0x11,0xe0,0x00,0x00,0x00,0x0f,0x00,
    0xc2,0x02,0x01,0x00,0xc1,0x02,0x03,0x11,0xc0,0x01,0x0a};

  // Create a new read_buffer instance
  plc4c_spi_read_buffer* read_buffer;
  plc4c_return_code return_code =
      plc4c_spi_read_buffer_create(cotp_connection_request_bytes,
                                   sizeof(cotp_connection_request_bytes),
                                   &read_buffer);
  if(return_code != OK) {
    TEST_FAIL_MESSAGE("Error");
  }

  plc4c_s7_read_write_tpkt_packet* message = NULL;
  return_code = plc4c_s7_read_write_tpkt_packet_parse(read_buffer,
                                            &message);
  if(return_code != OK) {
    TEST_FAIL_MESSAGE("Error");
  }

  plc4c_spi_write_buffer* write_buffer;
  return_code =
      plc4c_spi_write_buffer_create(sizeof(cotp_connection_request_bytes),
                                    &write_buffer);
  if(return_code != OK) {
    TEST_FAIL_MESSAGE("Error");
  }

  return_code = plc4c_s7_read_write_tpkt_packet_serialize(write_buffer, message);

  if(return_code != OK) {
    TEST_FAIL_MESSAGE("Error");
  }

  internal_assert_arrays_equal(cotp_connection_request_bytes, write_buffer,
                               sizeof(cotp_connection_request_bytes));
  printf("Success");

}

void setUp(void) {}

void tearDown(void) {}

int main(void) {
  UNITY_BEGIN();

  RUN_TEST(parse_cotp_connection_request);

  return UNITY_END();
}