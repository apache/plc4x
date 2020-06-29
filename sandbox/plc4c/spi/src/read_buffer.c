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
#include <string.h>

plc4c_return_code plc4c_spi_read_buffer_create(uint8_t* data, uint16_t length, plc4c_spi_read_buffer** buffer) {
  *buffer = malloc(sizeof(plc4c_spi_read_buffer));
  if(*buffer == NULL) {
    return NO_MEMORY;
  }

  (*buffer)->data = data;
  (*buffer)->length = length;
  (*buffer)->curPosByte = 0;
  (*buffer)->curPosBit = 0;

  return OK;
}

void plc4c_spi_read_buffer_destroy(plc4c_spi_read_buffer* buffer) {
  free(buffer);
}

uint32_t plc4c_spi_read_get_pos(plc4c_spi_read_buffer* buf) {
  return buf->curPosByte;
}

uint32_t plc4c_spi_read_get_total_bytes(plc4c_spi_read_buffer* buf) {
  return buf->length;
}

bool plc4c_spi_read_has_more(plc4c_spi_read_buffer* buf, uint16_t num_bits) {
  return (((buf->length - buf->curPosByte) * 8) - buf->curPosBit) >= num_bits;
}

plc4c_return_code plc4c_spi_read_get_bytes(plc4c_spi_read_buffer* buf, uint16_t start_pos_in_bytes, uint16_t end_pos_in_bytes, uint8_t** dest) {
  if(buf == NULL) {
    return NULL_VALUE;
  }
  if(dest == NULL) {
    return NULL_VALUE;
  }
  // Check if the arguments for start and stop position are correct.
  if(end_pos_in_bytes < start_pos_in_bytes) {
    return INVALID_ARGUMENT;
  }
  if(end_pos_in_bytes > buf->length) {
    return OUT_OF_RANGE;
  }
  uint16_t num_bytes = end_pos_in_bytes - start_pos_in_bytes;

  *dest = malloc(sizeof(uint8_t) * num_bytes);
  if(*dest == NULL) {
    return NO_MEMORY;
  }

  // Copy the requested bytes to the output.
  memcpy(*dest, buf->data, num_bytes);
  return OK;
}

uint8_t plc4c_spi_read_peek_byte(plc4c_spi_read_buffer* buf, uint16_t offset_in_bytes) {
  if(buf == NULL) {
    return 0;
  }
  if(buf->curPosByte + offset_in_bytes > buf->length) {
    return 0;
  }
  return (*buf->data) + (buf->curPosByte + offset_in_bytes);
}

bool plc4c_spi_read_bit(plc4c_spi_read_buffer* buf) {
  uint8_t cur_byte = (*buf->data) + buf->curPosByte;
  // We have to invert the position as bit 0 will be the first
  // (most significant bit).
  unsigned int bit_pos = ((unsigned int) 7) - buf->curPosBit;
  // Get the bit's value.
  bool value = ((cur_byte >> bit_pos) & 1) != 0;
  // If this was the last bit in this byte, move on to the next one.
  if(buf->curPosBit == 7) {
    buf->curPosByte++;
    buf->curPosBit = 0;
  } else {
    buf->curPosBit++;
  }
  return value;
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
