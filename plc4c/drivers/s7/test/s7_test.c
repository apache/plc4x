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
#include <unity.h>

#include "plc4c/spi/read_buffer.h"
#include "tpkt_packet.h"
#include "stdio.h"

void parser_serializer_test_s7_read_write();

void s7_address_parser_test();

void internal_assert_arrays_equal(uint8_t* expected_array, uint8_t* actual_array, uint8_t num_bytes);

/*void internal_assert_arrays_equal(uint8_t* expected_array,
                                  plc4c_spi_write_buffer* write_buffer,
                                  uint8_t num_bytes) {
  for (int i = 0; i < num_bytes; i++) {
    uint8_t expected_value = *(expected_array + i);
    uint8_t actual_value = *(write_buffer->data + i);
    // Needed for debugging on remote machines: Output the entire arrays content.
    if(expected_value != actual_value) {
      printf("\n");
      for(int j = 0; j < num_bytes; j++) {
        bool different = *(expected_array + j) !=  *(write_buffer->data + j);
        if(different) {
            printf("\033[0;31m");
        }
        printf("E=%02X %s A=%02X | ", *(expected_array + j), ( different ? "!=" : "=="), *(write_buffer->data + j));
        if(different) {
            printf("\033[0m");
        }
      }
      printf("\n");
    }
    TEST_ASSERT_EQUAL_UINT8_MESSAGE(expected_value, actual_value, "Byte arrays differ");
  }
}*/

void internal_parse_serialize_test(uint8_t* payload,
                                   uint8_t payload_size) {
  // Create a new read_buffer instance
  plc4c_spi_read_buffer* read_buffer;
  plc4c_return_code return_code =
      plc4c_spi_read_buffer_create(payload, payload_size, &read_buffer);
  if (return_code != OK) {
    TEST_FAIL_MESSAGE("Error creating read buffer");
  }

  plc4c_s7_read_write_tpkt_packet* message = NULL;
  return_code = plc4c_s7_read_write_tpkt_packet_parse(plc4x_spi_context_background(), read_buffer, &message);
  if (return_code != OK) {
    TEST_FAIL_MESSAGE("Error parsing packet");
  }

  plc4c_spi_write_buffer* write_buffer;
  return_code = plc4c_spi_write_buffer_create(payload_size, &write_buffer);
  if (return_code != OK) {
    TEST_FAIL_MESSAGE("Error writing to buffer");
  }

  return_code =
      plc4c_s7_read_write_tpkt_packet_serialize(plc4x_spi_context_background(), write_buffer, message);

  if (return_code != OK) {
    TEST_FAIL_MESSAGE("Error serializing");
  }

  internal_assert_arrays_equal(payload, write_buffer->data, payload_size);

  printf("Success");
}

void parse_cotp_connection_request() {
  uint8_t payload[] = {0x03, 0x00, 0x00, 0x16, 0x11, 0xe0, 0x00, 0x00,
                       0x00, 0x0f, 0x00, 0xc2, 0x02, 0x01, 0x00, 0xc1,
                       0x02, 0x03, 0x11, 0xc0, 0x01, 0x0a};
  internal_parse_serialize_test(payload, sizeof(payload));
}

void parse_cotp_connection_response() {
  uint8_t payload[] = {0x03, 0x00, 0x00, 0x16, 0x11, 0xd0, 0x00, 0x0f,
                       0x00, 0x0b, 0x00, 0xc0, 0x01, 0x0a, 0xc1, 0x02,
                       0x03, 0x11, 0xc2, 0x02, 0x01, 0x00};
  internal_parse_serialize_test(payload, sizeof(payload));
}

void parse_s7_communication_setup_request() {
  uint8_t payload[] = {0x03, 0x00, 0x00, 0x19, 0x02, 0xf0, 0x81, 0x32, 0x01,
                       0x00, 0x00, 0x00, 0x00, 0x00, 0x08, 0x00, 0x00, 0xf0,
                       0x00, 0x00, 0x08, 0x00, 0x08, 0x03, 0xf0};
  internal_parse_serialize_test(payload, sizeof(payload));
}

void parse_s7_communication_setup_response() {
  uint8_t payload[] = {0x03, 0x00, 0x00, 0x1b, 0x02, 0xf0, 0x80, 0x32, 0x03,
                       0x00, 0x00, 0x00, 0x00, 0x00, 0x08, 0x00, 0x00, 0x00,
                       0x00, 0xf0, 0x00, 0x00, 0x03, 0x00, 0x03, 0x00, 0xf0};
  internal_parse_serialize_test(payload, sizeof(payload));
}

void parse_s7_read_plc_type_request() {
  uint8_t payload[] = {0x03, 0x00, 0x00, 0x21, 0x02, 0xf0, 0x82, 0x32, 0x07,
                       0x00, 0x00, 0x00, 0x01, 0x00, 0x08, 0x00, 0x08, 0x00,
                       0x01, 0x12, 0x04, 0x11, 0x44, 0x01, 0x00, 0xff, 0x09,
                       0x00, 0x04, 0x00, 0x11, 0x00, 0x00};
  internal_parse_serialize_test(payload, sizeof(payload));
}

void parse_s7_read_plc_type_response() {
  uint8_t payload[] = {
      0x03, /* protocolId */
      0x00,
      0x00, 0x7d, /* len */
        0x02, /* headerLength */
        0xf0, /* tpduCode */
          0x80, /* EOT & tpduRef */

        0x32, /* protocolId */
        0x07, /* messageType */
        0x00, 0x00,
        0x00, 0x01, /* tpduReference */
        0x00, 0x0c, /* parameterLength */
        0x00, 0x60, /* payloadLength */

        0x00, /* parameterType */
        0x01, /* numItems */
          0x12, /* itemType */
          0x08, /* itemLength */
          0x12, /* method */
          0x84, /* cpuFunctionType & cpuFunctionGroup */
          0x01, /* cpuSubfunction */
          0x01, /* sequenceNumber */
          0x00, /* dataUnitReferenceNumber */
          0x00, /* lastDataUnit */
          0x00, 0x00, /* errorCode */

        0xff, /* returnCode */
        0x09, /* transportSize */
        0x00, 0x5c, /* dataLength */
        0x00, 0x11, /* typeClass & sublistExtract & sublistList */
        0x00, 0x00, /* szlIndex */

          0x00, 0x1c, /* szlItemLength */
          0x00, 0x03, /* szlItemCount */

          0x00, 0x01, /* itemIndex */
          0x36, 0x45, 0x53, 0x37, 0x20, 0x32, 0x31, 0x32, 0x2d, 0x31, 0x42, 0x44, 0x33, 0x30, 0x2d, 0x30, 0x58, 0x42, 0x30, 0x20, /* mlfb */
          0x20, 0x20, /* moduleTypeId */
          0x00, 0x01, /* ausbg */
          0x20, 0x20, /* ausbe */

          0x00, 0x06, /* itemIndex */
          0x36, 0x45, 0x53, 0x37, 0x20, 0x32, 0x31, 0x32, 0x2d, 0x31, 0x42, 0x44, 0x33, 0x30, 0x2d, 0x30, 0x58, 0x42, 0x30, 0x20, /* mlfb */
          0x20, 0x20, /* moduleTypeId */
          0x00, 0x01, /* ausbg */
          0x20, 0x20, /* ausbe */

          0x00, 0x07, /* itemIndex */
          0x36, 0x45, 0x53, 0x37, 0x20, 0x32, 0x31, 0x32, 0x2d, 0x31, 0x42, 0x44, 0x33, 0x30, 0x2d, 0x30, 0x58, 0x42, 0x30, 0x20, /* mlfb */
          0x20, 0x20, /* moduleTypeId */
          0x56, 0x02, /* ausbg */
          0x00, 0x02}; /* ausbe */
  internal_parse_serialize_test(payload, sizeof(payload));
}

void parse_s7_read_request() {
  uint8_t payload[] = {
      0x03, 0x00, 0x00, 0x43, 0x02, 0xf0, 0x8b, 0x32, 0x01, 0x00, 0x00, 0x00,
      0x0b, 0x00, 0x32, 0x00, 0x00, 0x04, 0x04, 0x12, 0x0a, 0x10, 0x01, 0x00,
      0x01, 0x00, 0x00, 0x82, 0x00, 0x00, 0x00, 0x12, 0x0a, 0x10, 0x01, 0x00,
      0x01, 0x00, 0x00, 0x82, 0x00, 0x00, 0x00, 0x12, 0x0a, 0x10, 0x01, 0x00,
      0x01, 0x00, 0x00, 0x82, 0x00, 0x00, 0x00, 0x12, 0x0a, 0x10, 0x01, 0x00,
      0x01, 0x00, 0x00, 0x82, 0x00, 0x00, 0x00};
  internal_parse_serialize_test(payload, sizeof(payload));
}

void parse_s7_read_response() {
  uint8_t payload[] = {0x03, 0x00, 0x00, 0x2C, 0x02, 0xf0, 0x80, 0x32, 0x03,
                       0x00, 0x00, 0x00, 0x0b, 0x00, 0x02, 0x00, 0x17, 0x00,
                       0x00, 0x04, 0x04, 0xff, 0x03, 0x00, 0x01, 0x01, 0x00,
                       0xff, 0x03, 0x00, 0x01, 0x01, 0x00, 0xff, 0x03, 0x00,
                       0x01, 0x01, 0x00, 0xff, 0x03, 0x00, 0x01, 0x01};
  internal_parse_serialize_test(payload, sizeof(payload));
}

void parse_s7_read_error_response() {
  uint8_t payload[] = {0x03, 0x00, 0x00, 0x13, 0x02, 0xf0, 0x80,
                       0x32, 0x02, 0x00, 0x00, 0x00, 0x0a, 0x00,
                       0x00, 0x00, 0x00, 0x85, 0x00};
  internal_parse_serialize_test(payload, sizeof(payload));
}

void parse_s7_write_request() {
  uint8_t payload[] = {
      0x03, 0x00, 0x00, 0x5A, 0x02, 0xf0, 0x8e, 0x32, 0x01, 0x00, 0x00, 0x00,
      0x0e, 0x00, 0x32, 0x00, 0x17, 0x05, 0x04, 0x12, 0x0a, 0x10, 0x01, 0x00,
      0x01, 0x00, 0x00, 0x82, 0x00, 0x00, 0x00, 0x12, 0x0a, 0x10, 0x01, 0x00,
      0x01, 0x00, 0x00, 0x82, 0x00, 0x00, 0x01, 0x12, 0x0a, 0x10, 0x01, 0x00,
      0x01, 0x00, 0x00, 0x82, 0x00, 0x00, 0x02, 0x12, 0x0a, 0x10, 0x01, 0x00,
      0x01, 0x00, 0x00, 0x82, 0x00, 0x00, 0x03, 0xff, 0x03, 0x00, 0x01, 0x01,
      0x00, 0xff, 0x03, 0x00, 0x01, 0x01, 0x00, 0xff, 0x03, 0x00, 0x01, 0x01,
      0x00, 0xff, 0x03, 0x00, 0x01, 0x01};
  internal_parse_serialize_test(payload, sizeof(payload));
}

void parse_s7_write_response() {
  uint8_t payload[] = {0x03, 0x00, 0x00, 0x19, 0x02, 0xf0, 0x80, 0x32, 0x03,
                       0x00, 0x00, 0x00, 0x0e, 0x00, 0x02, 0x00, 0x04, 0x00,
                       0x00, 0x05, 0x04, 0xff, 0xff, 0xff, 0xff};
  internal_parse_serialize_test(payload, sizeof(payload));
}

void setUp(void) {}

void tearDown(void) {}

int main(void) {
  UNITY_BEGIN();

  RUN_TEST(parse_cotp_connection_request);
  RUN_TEST(parse_cotp_connection_response);
  RUN_TEST(parse_s7_communication_setup_request);
  RUN_TEST(parse_s7_communication_setup_response);
  RUN_TEST(parse_s7_read_plc_type_request);
  RUN_TEST(parse_s7_read_plc_type_response);
  RUN_TEST(parse_s7_read_request);
  RUN_TEST(parse_s7_read_response);
  RUN_TEST(parse_s7_read_error_response);
  RUN_TEST(parse_s7_write_request);
  RUN_TEST(parse_s7_write_response);

  // Run the address parser tests ...
  // TODO: Commented out as it seems to only fail while doing releases :-/
  //RUN_TEST(s7_address_parser_test);

  parser_serializer_test_s7_read_write();

  return UNITY_END();
}