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
#include <plc4c/spi/types_private.h>
#include <stdlib.h>
#include <unity.h>

#include "plc4c/spi/system_private.h"

void test_system_plc4c_system_create_connection_args(
    char *connection_string, plc4c_return_code expected_return_code,
    char *expected_connection_string, char *expected_protocol_code,
    char *expected_transport_code, char *expected_transport_connect_information,
    char *expected_parameters) {

  printf("Running test with '%s'", connection_string);

  plc4c_connection *connection = NULL;
  plc4c_return_code result =
      plc4c_system_create_connection(connection_string, &connection);
  TEST_ASSERT_EQUAL(expected_return_code, result);
  if (expected_return_code != OK) {
    TEST_ASSERT_NULL(connection);
  } else {
    TEST_ASSERT_EQUAL_STRING(expected_connection_string,
                             connection->connection_string);
    TEST_ASSERT_EQUAL_STRING(expected_protocol_code, connection->protocol_code);
    TEST_ASSERT_EQUAL_STRING(expected_transport_code,
                             connection->transport_code);
    TEST_ASSERT_EQUAL_STRING(expected_transport_connect_information,
                             connection->transport_connect_information);
    TEST_ASSERT_EQUAL_STRING(expected_parameters, connection->parameters);
    free(connection);
  }
  printf(" -> OK\n");
}

void test_system_plc4c_system_create_connection(void) {
  test_system_plc4c_system_create_connection_args(
      "s7://1.2.3.4", OK, "s7://1.2.3.4", "s7", NULL, "1.2.3.4", NULL);
  test_system_plc4c_system_create_connection_args(
      "s7:tcp://1.2.3.4", OK, "s7:tcp://1.2.3.4", "s7", "tcp", "1.2.3.4", NULL);
  test_system_plc4c_system_create_connection_args("s7://1.2.3.4?params", OK,
                                                  "s7://1.2.3.4?params", "s7",
                                                  NULL, "1.2.3.4", "params");
  test_system_plc4c_system_create_connection_args(
      "s7:tcp://1.2.3.4?params", OK, "s7:tcp://1.2.3.4?params", "s7", "tcp",
      "1.2.3.4", "params");

  // A colon after the "://" shouldn't matter ...
  test_system_plc4c_system_create_connection_args(
      "s7://1.2.3.4:42", OK, "s7://1.2.3.4:42", "s7", NULL, "1.2.3.4:42", NULL);
  test_system_plc4c_system_create_connection_args(
      "s7://1.2.3.4?param=a:42", OK, "s7://1.2.3.4?param=a:42", "s7", NULL,
      "1.2.3.4", "param=a:42");

  // Well obviously the parser shouldn't be able to find anything here ...
  test_system_plc4c_system_create_connection_args(
      "hurz", INVALID_CONNECTION_STRING, NULL, NULL, NULL, NULL, NULL);
  // In these cases the parser expects a "//" after the second colon, which
  // isn't there ...
  test_system_plc4c_system_create_connection_args(
      "a:b:c://d", INVALID_CONNECTION_STRING, NULL, NULL, NULL, NULL, NULL);
  test_system_plc4c_system_create_connection_args(
      "a:b:d", INVALID_CONNECTION_STRING, NULL, NULL, NULL, NULL, NULL);

  // There should only be one question-mark ...
  test_system_plc4c_system_create_connection_args(
      "a://a?b?c", INVALID_CONNECTION_STRING, NULL, NULL, NULL, NULL, NULL);
}
