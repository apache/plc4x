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

#include <plc4c/spi/evaluation_helper.h>

#include <stdbool.h>
#include <math.h>

bool plc4c_spi_evaluation_helper_equals(int a, int b) {
  return a == b;
}

double plc4c_spi_evaluation_helper_ceil(double a) {
  return ceil(a);
}

uint8_t plc4c_spi_evaluation_helper_count(plc4c_list* a) {
  return plc4c_utils_list_size(a);
}

uint8_t plc4c_spi_evaluation_helper_array_size_in_bytes(plc4c_list* a) {
  // TODO: This sort of can't work in C as we don't have the type information ...
  return 0;
}
