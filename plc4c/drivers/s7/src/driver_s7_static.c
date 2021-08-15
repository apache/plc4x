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

#include <plc4c/spi/read_buffer.h>
#include <plc4c/spi/write_buffer.h>
#include <string.h>

/*
 *
 *   Static functions
 *
 */

uint16_t plc4c_s7_read_write_event_item_length(plc4c_spi_read_buffer* io, uint16_t valueLength) {
  return 0;
}

uint16_t plc4c_s7_read_write_right_shift3(plc4c_spi_read_buffer* io) {
  return 0;
}

void plc4c_s7_read_write_left_shift3(plc4c_spi_write_buffer* io, uint16_t value) {
}

uint8_t plc4c_s7_read_write_bcd_to_int(plc4c_spi_read_buffer* io) {
  return 0;
}

uint16_t plc4c_s7_read_write_s7msec_to_int(plc4c_spi_read_buffer* io) {
  return 0;
}

char* plc4c_s7_read_write_parse_s7_string(plc4c_spi_read_buffer* io,
                                          int32_t stringLength,
                                          char* encoding) {
  if (strcmp(encoding, "UTF-8") == 0) {
    // Read the max length (which is not interesting for us.
    uint8_t maxLen;
    plc4c_return_code res = plc4c_spi_read_unsigned_byte(io, 8, &maxLen);
    if (res != OK) {
      return NULL;
    }
    // Read the effective length of the string.
    uint8_t effectiveStringLength;
    res = plc4c_spi_read_unsigned_byte(io, 8, &effectiveStringLength);
    if (res != OK) {
      return NULL;
    }
    char* result = malloc(sizeof(char) * (effectiveStringLength + 1));
    if (result == NULL) {
      return NULL;
    }
    char* curPos = result;
    for(int i = 0; i < effectiveStringLength; i++) {
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

char* plc4c_s7_read_write_parse_s7_char(plc4c_spi_read_buffer* io,
                                        char* encoding) {
  if (strcmp(encoding, "UTF-8") == 0) {
    char* result = malloc(sizeof(char) * 2);
    if (result == NULL) {
      return NULL;
    }
    uint8_t val;
    plc4c_return_code res = plc4c_spi_read_unsigned_byte(io, 8, &val);
    if (res != OK) {
      return NULL;
    }
    *result = (char) val;
    *(result+1) = '\0';
    return result;
  } else if (strcmp(encoding, "UTF-16") == 0) {
  }
  return "";
}

time_t plc4c_s7_read_write_parse_tia_time(plc4c_spi_read_buffer* io) {
  // TODO: Implement ...
  return 0;
}

time_t plc4c_s7_read_write_parse_s5_time(plc4c_spi_read_buffer* io) {
  // TODO: Implement ...
  return 0;
}

time_t plc4c_s7_read_write_parse_tia_l_time(plc4c_spi_read_buffer* io) {
  // TODO: Implement ...
  return 0;
}

time_t plc4c_s7_read_write_parse_tia_date(plc4c_spi_read_buffer* io) {
  // TODO: Implement ...
  return 0;
}

time_t plc4c_s7_read_write_parse_tia_time_of_day(plc4c_spi_read_buffer* io) {
  // TODO: Implement ...
  return 0;
}

time_t plc4c_s7_read_write_parse_tia_date_time(plc4c_spi_read_buffer* io) {
  // TODO: Implement ...
  return 0;
}
