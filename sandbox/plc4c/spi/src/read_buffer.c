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

uint8_t bit_mask[9] = {0, 1, 3, 7, 15, 31, 63, 127, 255};
// This matrix contains constants for reading X bits starting with bit Y.
static const uint8_t bit_matrix[8][8] = {
    // Reading 1 bit
    {128, 64, 32, 16, 8, 4, 2, 1},
    // Reading 2 bits
    {192, 96, 48, 24, 12, 6, 3, 0},
    // Reading 3 bits
    {224, 112, 56, 28, 14, 7, 0, 0},
    // Reading 4 bits
    {240, 120, 60, 30, 15, 0, 0, 0},
    // Reading 5 bits
    {248, 124, 62, 31, 0, 0, 0, 0},
    // Reading 6 bits
    {252, 126, 63, 0, 0, 0, 0, 0},
    // Reading 7 bits
    {254, 127, 0, 0, 0, 0, 0, 0},
    // Reading 8 bits
    {255, 0, 0, 0, 0, 0, 0, 0}
};
uint8_t plc4c_spi_read_unsigned_byte_internal(uint8_t data, uint8_t num_bits, uint8_t from_bit) {
  return (data & bit_matrix[num_bits - 1][from_bit]) >> (((unsigned int) 8) - (from_bit + num_bits));
}
uint8_t plc4c_spi_read_unsigned_byte_get_byte_internal(plc4c_spi_read_buffer* buf, uint8_t offset) {
  uint8_t value = 0;
  plc4c_spi_read_peek_byte(buf, offset, &value);
  return value;
}

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

plc4c_return_code plc4c_spi_read_peek_byte(plc4c_spi_read_buffer* buf, uint16_t offset_in_bytes, uint8_t* value) {
  if(buf == NULL) {
    return NULL_VALUE;
  }
  if(buf->curPosByte + offset_in_bytes > buf->length) {
    return OUT_OF_RANGE;
  }
  *value = *(buf->data + (buf->curPosByte + offset_in_bytes));
  return OK;
}

plc4c_return_code plc4c_spi_read_bit(plc4c_spi_read_buffer* buf, bool* value) {
  uint8_t cur_byte = (*buf->data) + buf->curPosByte;
  // We have to invert the position as bit 0 will be the first
  // (most significant bit).
  unsigned int bit_pos = ((unsigned int) 7) - buf->curPosBit;
  // Get the bit's value.
  *value = ((cur_byte >> bit_pos) & 1) != 0;
  // If this was the last bit in this byte, move on to the next one.
  if(buf->curPosBit == 7) {
    buf->curPosByte++;
    buf->curPosBit = 0;
  } else {
    buf->curPosBit++;
  }
  return OK;
}

// Unsigned Integers ...

plc4c_return_code plc4c_spi_read_unsigned_byte(plc4c_spi_read_buffer* buf, uint8_t num_bits, uint8_t* value) {
  // If more than 8 bits are requested, return an error.
  if(num_bits > 8) {
    return OUT_OF_RANGE;
  }
  if(buf->curPosByte >= (buf->length - 1)) {
    return OUT_OF_RANGE;
  }
  // If the bit-offset is currently 0 and we're reading
  // a full byte, go this shortcut.
  if((buf->curPosBit == 0) && (num_bits == 8)){
    *value = plc4c_spi_read_unsigned_byte_get_byte_internal(buf, 0);
    buf->curPosByte++;
    return OK;
  }
  // in this case the current byte alone is enough to service this request.
  else if((8 - buf->curPosBit) >= num_bits) {
    *value = plc4c_spi_read_unsigned_byte_internal(plc4c_spi_read_unsigned_byte_get_byte_internal(buf, 0), num_bits, buf->curPosBit);
    if(buf->curPosBit + num_bits == 8) {
      buf->curPosByte++;
      buf->curPosBit = 0;
    } else {
      buf->curPosBit += num_bits;
    }
    return OK;
  }
  // In this case we also need the following byte.
  else {
    // Read the remaining bits of the first byte
    uint8_t cur_byte = plc4c_spi_read_unsigned_byte_internal(plc4c_spi_read_unsigned_byte_get_byte_internal(buf, 0), 8 - buf->curPosBit, buf->curPosBit);
    // Shift the result left by the amount of bits in the second byte
    uint8_t num_bits_remaining = num_bits - (8 - buf->curPosBit);
    cur_byte = cur_byte << num_bits_remaining;
    // Read the remaining bits from the second byte
    uint8_t next_byte = plc4c_spi_read_unsigned_byte_get_byte_internal(buf, 1);
    next_byte = plc4c_spi_read_unsigned_byte_internal(next_byte, num_bits_remaining, 0);;
    // Bring the two parts together to assemble the resulting value
    *value = cur_byte | next_byte;
    // Update the buffer position
    buf->curPosByte = buf->curPosByte + ((buf->curPosBit + num_bits) / 8);
    buf->curPosBit = (buf->curPosBit + num_bits) % 8;
    return OK;
  }
}

plc4c_return_code plc4c_spi_read_unsigned_short(plc4c_spi_read_buffer* buf, uint8_t num_bits, uint16_t* value) {
  // If the bit-offset is currently 0, then we simply read a byte ...
  if(buf->curPosBit == 0) {
    uint16_t cur_short = (*buf->data) + buf->curPosByte;
    buf->curPosByte += 2;
    return cur_short;
  } else {
    uint8_t cur_short = (*buf->data) + buf->curPosByte;
    cur_short = cur_short << buf->curPosBit;
    uint8_t next_byte = (*buf->data) + (buf->curPosByte + 1);
    next_byte = next_byte >> buf->curPosBit;
    uint8_t virtual_byte = cur_short | next_byte;
    buf->curPosByte++;
    return virtual_byte;
  }
}

plc4c_return_code plc4c_spi_read_unsigned_int(plc4c_spi_read_buffer* buf, uint8_t num_bits, uint32_t* value) {
  return OK;
}

plc4c_return_code plc4c_spi_read_unsigned_long(plc4c_spi_read_buffer* buf, uint8_t num_bits, uint64_t* value) {
  return OK;
}

// TODO: Not sure which type to use in this case ...
/*uint128_t plc4c_spi_read_unsigned_big_integer(plc4c_spi_read_buffer* buf, uint8_t num_bits) {
  return OK;
}*/

// Signed Integers ...

plc4c_return_code plc4c_spi_read_byte(plc4c_spi_read_buffer* buf, uint8_t num_bits, int8_t* value) {
  return OK;
}

plc4c_return_code plc4c_spi_read_short(plc4c_spi_read_buffer* buf, uint8_t num_bits, int16_t* value) {
  return OK;
}

plc4c_return_code plc4c_spi_read_int(plc4c_spi_read_buffer* buf, uint8_t num_bits, int32_t* value) {
  return OK;
}

plc4c_return_code plc4c_spi_read_long(plc4c_spi_read_buffer* buf, uint8_t num_bits, int64_t* value) {
  return OK;
}

// TODO: Not sure which type to use in this case ...
/*int128_t plc4c_spi_read_big_integer(plc4c_spi_read_buffer* buf, uint8_t num_bits);
 * return OK;
 * }*/

// Floating Point Numbers ...

plc4c_return_code plc4c_spi_read_float(plc4c_spi_read_buffer* buf, uint8_t num_bits, float* value) {
  return OK;
}

plc4c_return_code plc4c_spi_read_double(plc4c_spi_read_buffer* buf, uint8_t num_bits, double* value) {
  return OK;
}

// TODO: Not sure which type to use in this case ...
/*doubledouble plc4c_spi_read_big_decimal(plc4c_spi_read_buffer* buf, uint8_t num_bits);
 * return 0;
 * } */

plc4c_return_code plc4c_spi_read_string(plc4c_spi_read_buffer* buf, uint8_t num_bits, char* encoding, char** value) {
  return OK;
}
