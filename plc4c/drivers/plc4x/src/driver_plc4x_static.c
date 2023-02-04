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
#include "plc4c/driver_plc4x_static.h"

/*
 *
 *   Static functions
 *
 */

uint8_t plc4c_spi_evaluation_helper_str_len(char* str) {
    return strlen(str);
}

char* plc4c_plc4x_read_write_parse_string(plc4c_spi_read_buffer* io, char* encoding) {
  if (strcmp(encoding, "UTF-8") == 0) {
    // Read the max length (which is not interesting for us.
    uint8_t length;
    plc4c_return_code res = plc4c_spi_read_unsigned_byte(io, 8, &length);
    if (res != OK) {
      return NULL;
    }
    char* result = malloc(sizeof(char) * (length + 1));
    if (result == NULL) {
      return NULL;
    }
    char* curPos = result;
    for(int i = 0; i < length; i++) {
      uint8_t val;
      plc4c_return_code res = plc4c_spi_read_unsigned_byte(io, 8, &val);
      if (res != OK) {
        return NULL;
      }
      *curPos = (char) val;
      curPos++;
    }
    *curPos = '\0';
    return result;
  } else if (strcmp(encoding, "UTF-16") == 0) {
  }
  return "";
}
