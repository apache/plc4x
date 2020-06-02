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

uint32_t plc4c_spi_read_get_pos(plc4c_spi_read_buffer* buf) {
  return 0;
}

uint32_t plc4c_spi_read_get_total_bytes(plc4c_spi_read_buffer* buf) {
  return 0;
}

bool plc4c_spi_read_has_more(plc4c_spi_read_buffer* buf, uint8_t num_bits) {
  return false;
}

uint8_t* plc4c_spi_read_get_bytes(plc4c_spi_read_buffer* buf, uint32_t start_pos_in_bytes, uint32_t end_pos_in_bytes) {
  return NULL;
}

uint8_t plc4c_spi_read_peek_byte(plc4c_spi_read_buffer* buf, uint32_t offset_in_bytes) {
  return 0;
}

bool plc4c_spi_read_bit(plc4c_spi_read_buffer* buf) {
  return false;
}

// Unsigned Integers ...

uint8_t plc4c_spi_read_unsigned_byte(plc4c_spi_read_buffer* buf, uint8_t num_bits) {
  return 0;
}

uint16_t plc4c_spi_read_unsigned_short(plc4c_spi_read_buffer* buf, uint8_t num_bits) {
  return 0;
}

uint32_t plc4c_spi_read_unsigned_int(plc4c_spi_read_buffer* buf, uint8_t num_bits) {
  return 0;
}

uint64_t plc4c_spi_read_unsigned_long(plc4c_spi_read_buffer* buf, uint8_t num_bits) {
  return 0;
}

// TODO: Not sure which type to use in this case ...
/*uint128_t plc4c_spi_read_unsigned_big_integer(plc4c_spi_read_buffer* buf, uint8_t num_bits) {
  return 0;
}*/

// Signed Integers ...

int8_t plc4c_spi_read_byte(plc4c_spi_read_buffer* buf, uint8_t num_bits) {
  return 0;
}

int16_t plc4c_spi_read_short(plc4c_spi_read_buffer* buf, uint8_t num_bits) {
  return 0;
}

int32_t plc4c_spi_read_int(plc4c_spi_read_buffer* buf, uint8_t num_bits) {
  return 0;
}

int64_t plc4c_spi_read_long(plc4c_spi_read_buffer* buf, uint8_t num_bits) {
  return 0;
}

// TODO: Not sure which type to use in this case ...
/*int128_t plc4c_spi_read_big_integer(plc4c_spi_read_buffer* buf, uint8_t num_bits);
 * return 0;
 * }*/

// Floating Point Numbers ...

float plc4c_spi_read_float(plc4c_spi_read_buffer* buf, uint8_t num_bits) {
  return 0;
}

double plc4c_spi_read_double(plc4c_spi_read_buffer* buf, uint8_t num_bits) {
  return 0;
}

// TODO: Not sure which type to use in this case ...
/*doubledouble plc4c_spi_read_big_decimal(plc4c_spi_read_buffer* buf, uint8_t num_bits);
 * return 0;
 * } */

char* plc4c_spi_read_string(plc4c_spi_read_buffer* buf, uint8_t num_bits, char* encoding) {
  return NULL;
}
