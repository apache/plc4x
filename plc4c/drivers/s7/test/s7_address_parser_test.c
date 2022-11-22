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
#include <memory_area.h>
#include <s7_var_request_parameter_item.h>
#include <transport_size.h>
#include <unity.h>

#include "plc4c/spi/read_buffer.h"

plc4c_return_code plc4c_driver_s7_encode_address(char* address, plc4c_s7_read_write_s7_var_request_parameter_item** item);

void internal_parse_addresses_test(
    char* address, plc4c_return_code return_code,
    plc4c_s7_read_write_transport_size transport_size,
    plc4c_s7_read_write_memory_area memory_area, uint16_t db_number,
    uint16_t byte_address, uint8_t bit_address, uint16_t number_of_elements) {

  // Call the parser ...
  plc4c_s7_read_write_s7_var_request_parameter_item* item;
  plc4c_return_code act_return_code =
      plc4c_driver_s7_encode_address(address, &item);

  // Check the result ...
  if(act_return_code != OK) {
    TEST_ASSERT_EQUAL_INT(return_code, act_return_code);
  } else {
    TEST_ASSERT_EQUAL_INT(item->_type, plc4c_s7_read_write_s7_var_request_parameter_item_type_plc4c_s7_read_write_s7_var_request_parameter_item_address);
    TEST_ASSERT_EQUAL_INT(return_code, act_return_code);
    TEST_ASSERT_EQUAL_INT(transport_size, item->s7_var_request_parameter_item_address_address->s7_address_any_transport_size);
    TEST_ASSERT_EQUAL_INT(memory_area, item->s7_var_request_parameter_item_address_address->s7_address_any_area);
    TEST_ASSERT_EQUAL_INT(db_number, item->s7_var_request_parameter_item_address_address->s7_address_any_db_number);
    TEST_ASSERT_EQUAL_INT(byte_address, item->s7_var_request_parameter_item_address_address->s7_address_any_byte_address);
    TEST_ASSERT_EQUAL_INT(bit_address, item->s7_var_request_parameter_item_address_address->s7_address_any_bit_address);
    TEST_ASSERT_EQUAL_INT(number_of_elements, item->s7_var_request_parameter_item_address_address->s7_address_any_number_of_elements);
  }
}

void parse_addresses_test() {
  internal_parse_addresses_test(
      "%I23.2:BOOL",OK,
      plc4c_s7_read_write_transport_size_BOOL,
      plc4c_s7_read_write_memory_area_INPUTS,
      0, 23, 2, 1);
  internal_parse_addresses_test(
      "%I23.2:BOOL[3]",OK,
      plc4c_s7_read_write_transport_size_BOOL,
      plc4c_s7_read_write_memory_area_INPUTS,
      0, 23, 2, 3);
  internal_parse_addresses_test(
      "%Q23.2:BOOL[1]",OK,
      plc4c_s7_read_write_transport_size_BOOL,
      plc4c_s7_read_write_memory_area_OUTPUTS,
      0, 23, 2, 1);
  internal_parse_addresses_test(
      "%I23.0:BYTE",OK,
      plc4c_s7_read_write_transport_size_BYTE,
      plc4c_s7_read_write_memory_area_INPUTS,
      0, 23, 0, 1);
  internal_parse_addresses_test(
      "%I23:BYTE",OK,
      plc4c_s7_read_write_transport_size_BYTE,
      plc4c_s7_read_write_memory_area_INPUTS,
      0, 23, 0, 1);
  internal_parse_addresses_test(
      "%DB23.DBX3.4:BOOL",OK,
      plc4c_s7_read_write_transport_size_BOOL,
      plc4c_s7_read_write_memory_area_DATA_BLOCKS,
      23, 3, 4, 1);
  internal_parse_addresses_test(
      "%DB23.DBX3.4:BOOL[4]",OK,
      plc4c_s7_read_write_transport_size_BOOL,
      plc4c_s7_read_write_memory_area_DATA_BLOCKS,
      23, 3, 4, 4);
  internal_parse_addresses_test(
      "%DB23.DBW3.0:INT",OK,
      plc4c_s7_read_write_transport_size_INT,
      plc4c_s7_read_write_memory_area_DATA_BLOCKS,
      23, 3, 0, 1);
  internal_parse_addresses_test(
      "%DB23.DBW3:INT",OK,
      plc4c_s7_read_write_transport_size_INT,
      plc4c_s7_read_write_memory_area_DATA_BLOCKS,
      23, 3, 0, 1);
  // TODO: add more tests ...
  internal_parse_addresses_test(
      "10-01-00-01-00-2D-84-00-00-08",OK,
      plc4c_s7_read_write_transport_size_BOOL,
      plc4c_s7_read_write_memory_area_DATA_BLOCKS,
      45, 1, 0, 1);
}

void s7_address_parser_test() {
  parse_addresses_test();
}