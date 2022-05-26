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

#include "stdbool.h"
#include "stdio.h"
#include "stdint.h"

void internal_assert_arrays_equal(uint8_t* expected_array,
                                  uint8_t* actual_array,
                                  uint8_t num_bytes) {
  for (int i = 0; i < num_bytes; i++) {
    uint8_t expected_value = *(expected_array + i);
    uint8_t actual_value = *(actual_array + i);
    // Needed for debugging on remote machines: Output the entire arrays content.
    if(expected_value != actual_value) {
      printf("\n");
      for(int j = 0; j < num_bytes; j++) {
        bool different = *(expected_array + j) !=  *(actual_array + j);
        if(different) {
          printf("\033[0;31m");
        }
        printf("E=%02X %s A=%02X | ", *(expected_array + j), ( different ? "!=" : "=="), *(actual_array + j));
        if(different) {
          printf("\033[0m");
        }
      }
      printf("\n");
    }
    //TEST_ASSERT_EQUAL_UINT8_MESSAGE(expected_value, actual_value, "Byte arrays differ");
  }
}
