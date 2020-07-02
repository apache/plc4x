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

// As we're doing some operations where byte-order is important, we need this
// little helper to find out if we're on a big- or little-endian machine.
const int i = 1;
#define plc4c_is_bigendian() (((char)&i) == 0)

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
    {255, 0, 0, 0, 0, 0, 0, 0}};

uint8_t plc4c_spi_read_unsigned_byte_internal(uint8_t data, uint8_t num_bits,
                                              uint8_t from_bit) {
  return (data & bit_matrix[num_bits - 1][from_bit]) >>
         (((unsigned int)8) - (from_bit + num_bits));
}

uint8_t plc4c_spi_read_unsigned_byte_get_byte_internal(
    plc4c_spi_read_buffer* buf, uint8_t offset) {
  return *(buf->data + (buf->curPosByte + offset));
}

plc4c_return_code plc4c_spi_read_unsigned_bits_internal(
    plc4c_spi_read_buffer* buf, uint8_t num_bits, uint8_t* value) {
  if (buf == NULL) {
    return NULL_VALUE;
  }

  // If the bit-offset is currently 0 and we're reading
  // a full byte, go this shortcut.
  if ((buf->curPosBit == 0) && (num_bits % 8 == 0)) {
    if (buf->curPosByte >= (buf->length - 1)) {
      return OUT_OF_RANGE;
    }

    // Find how many full bytes we'll be reading.
    uint8_t num_bytes = num_bits / 8;
    // If this is little endian, go to the end of the range.
    if (!plc4c_is_bigendian()) {
      value = value + (num_bytes - 1);
    }
    // Read each of these.
    for (int i = 0; i < num_bytes; i++) {
      *value = plc4c_spi_read_unsigned_byte_get_byte_internal(buf, 0);
      // Move the read-pointer to the next byte.
      buf->curPosByte++;
      // Move the write-pointer to the next byte.
      value = plc4c_is_bigendian() ? value + 1 : value - 1;
    }
    return OK;
  }

  // in this case the current byte alone is enough to service this request.
  else if ((((unsigned int)8) - buf->curPosBit) >= num_bits) {
    if (buf->curPosByte > (buf->length - 1)) {
      return OUT_OF_RANGE;
    }
    *value = plc4c_spi_read_unsigned_byte_internal(
        plc4c_spi_read_unsigned_byte_get_byte_internal(buf, 0), num_bits,
        buf->curPosBit);
    if (buf->curPosBit + num_bits == 8) {
      buf->curPosByte++;
      buf->curPosBit = 0;
    } else {
      buf->curPosBit += num_bits;
    }
    return OK;
  }

  // In this case we also need more than one following byte.
  else {
    uint8_t* original_value = value;
    // Find out how many bytes we need to read
    uint8_t num_bytes = ((buf->curPosBit + num_bits) / 8) -
                        (((buf->curPosBit + num_bits) % 8 == 0) ? 1 : 0);
    if ((buf->curPosByte + num_bytes) > buf->length) {
      return OUT_OF_RANGE;
    }
    // If this is little endian, go to the end of the range.
    if (!plc4c_is_bigendian()) {
      value = value + (num_bytes - 1);
    }
    // Find out how many of the bits will be read from the first byte
    // It's actually just the rest of the byte as we checked if it all fits
    // in one byte in the else-block before.
    uint8_t num_bits_first_byte = 8 - buf->curPosBit;
    // Having read the bits from the first byte, see how many bits will
    // have to be read in the last byte.
    // (This will also be the number of bits that everything needs to be
    // shifted left)
    uint8_t num_bits_last_byte = (num_bits - num_bits_first_byte) % 8;
    // All other bytes in-between will just read all 8 bits.

    // Read the bits of the first byte
    uint8_t part_first_byte = plc4c_spi_read_unsigned_byte_internal(
        plc4c_spi_read_unsigned_byte_get_byte_internal(buf, 0),
        num_bits_first_byte, buf->curPosBit);
    buf->curPosByte++;
    // Shift the result left by the amount of bits in the last byte
    *value = part_first_byte << num_bits_last_byte;

    // For each of the following bytes.
    for (int i = 0; i < num_bytes; i++) {
      // Get the next full byte.
      uint8_t cur_byte = plc4c_spi_read_unsigned_byte_get_byte_internal(buf, 0);
      // Get the rest of the bits that belong to the previous byte.
      uint8_t rest_previous_byte = plc4c_spi_read_unsigned_byte_internal(
          cur_byte, num_bits_last_byte, 0);
      // Add that to the end of the last one
      *value = *value | rest_previous_byte;
      // If this is not the last byte ... continue
      if (i < (num_bits / 8)) {
        // Set value to the next address.
        value = plc4c_is_bigendian() ? value + 1 : value - 1;
        // Initialize the first part of the following byte with the rest of the
        // next byte.
        *value = plc4c_spi_read_unsigned_byte_internal(
                     cur_byte, (8 - num_bits_last_byte), num_bits_last_byte)
                 << num_bits_last_byte;
        buf->curPosByte++;
      }
    }

    // Update the buffer position
    buf->curPosBit = num_bits_last_byte;
    return OK;
  }
}

plc4c_return_code plc4c_spi_read_buffer_create(uint8_t* data, uint16_t length,
                                               plc4c_spi_read_buffer** buffer) {
  *buffer = malloc(sizeof(plc4c_spi_read_buffer));
  if (*buffer == NULL) {
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

plc4c_return_code plc4c_spi_read_get_bytes(plc4c_spi_read_buffer* buf,
                                           uint16_t start_pos_in_bytes,
                                           uint16_t end_pos_in_bytes,
                                           uint8_t** dest) {
  if (buf == NULL) {
    return NULL_VALUE;
  }
  if (dest == NULL) {
    return NULL_VALUE;
  }
  // Check if the arguments for start and stop position are correct.
  if (end_pos_in_bytes < start_pos_in_bytes) {
    return INVALID_ARGUMENT;
  }
  if (end_pos_in_bytes > buf->length) {
    return OUT_OF_RANGE;
  }
  uint16_t num_bytes = end_pos_in_bytes - start_pos_in_bytes;

  *dest = malloc(sizeof(uint8_t) * num_bytes);
  if (*dest == NULL) {
    return NO_MEMORY;
  }

  // Copy the requested bytes to the output.
  memcpy(*dest, buf->data, num_bytes);
  return OK;
}

plc4c_return_code plc4c_spi_read_peek_byte(plc4c_spi_read_buffer* buf,
                                           uint16_t offset_in_bytes,
                                           uint8_t* value) {
  if (buf == NULL) {
    return NULL_VALUE;
  }
  if (buf->curPosByte + offset_in_bytes > buf->length) {
    return OUT_OF_RANGE;
  }
  *value = *(buf->data + (buf->curPosByte + offset_in_bytes));
  return OK;
}

plc4c_return_code plc4c_spi_read_bit(plc4c_spi_read_buffer* buf, bool* value) {
  uint8_t cur_byte = (*buf->data) + buf->curPosByte;
  // We have to invert the position as bit 0 will be the first
  // (most significant bit).
  unsigned int bit_pos = ((unsigned int)7) - buf->curPosBit;
  // Get the bit's value.
  *value = ((cur_byte >> bit_pos) & 1) != 0;
  // If this was the last bit in this byte, move on to the next one.
  if (buf->curPosBit == 7) {
    buf->curPosByte++;
    buf->curPosBit = 0;
  } else {
    buf->curPosBit++;
  }
  return OK;
}

// Unsigned Integers ...

plc4c_return_code plc4c_spi_read_unsigned_byte(plc4c_spi_read_buffer* buf,
                                               uint8_t num_bits,
                                               uint8_t* value) {
  // If more than 8 bits are requested, return an error.
  if (num_bits > 8) {
    return OUT_OF_RANGE;
  }
  // Get the bits.
  return plc4c_spi_read_unsigned_bits_internal(buf, num_bits, value);
}

plc4c_return_code plc4c_spi_read_unsigned_short(plc4c_spi_read_buffer* buf,
                                                uint8_t num_bits,
                                                uint16_t* value) {
  // If more than 16 bits are requested, return an error.
  if (num_bits > 16) {
    return OUT_OF_RANGE;
  }
  // Get the bits.
  plc4c_return_code res =
      plc4c_spi_read_unsigned_bits_internal(buf, num_bits, value);
  // Shift the bits to the right position.
  if (res == OK) {
    if (num_bits <= 8) {
      *value = *value >> 8;
    }
  }
  return res;
}

plc4c_return_code plc4c_spi_read_unsigned_int(plc4c_spi_read_buffer* buf,
                                              uint8_t num_bits,
                                              uint32_t* value) {
  // If more than 32 bits are requested, return an error.
  if (num_bits > 32) {
    return OUT_OF_RANGE;
  }
  // Get the bits.
  plc4c_return_code res =
      plc4c_spi_read_unsigned_bits_internal(buf, num_bits, value);
  // Shift the bits to the right position.
  if (res == OK) {
    if (num_bits <= 8) {
      *value = *value >> 24;
    } else if (num_bits <= 16) {
      *value = *value >> 16;
    } else if (num_bits <= 24) {
      *value = *value >> 8;
    }
  }
  return res;
}

plc4c_return_code plc4c_spi_read_unsigned_long(plc4c_spi_read_buffer* buf,
                                               uint8_t num_bits,
                                               uint64_t* value) {
  // If more than 64 bits are requested, return an error.
  if (num_bits > 64) {
    return OUT_OF_RANGE;
  }
  // Get the bits.
  plc4c_return_code res =
      plc4c_spi_read_unsigned_bits_internal(buf, num_bits, value);
  // Shift the bits to the right position.
  if (res == OK) {
    if (num_bits <= 8) {
      *value = *value >> 56;
    } else if (num_bits <= 16) {
      *value = *value >> 48;
    } else if (num_bits <= 24) {
      *value = *value >> 40;
    } else if (num_bits <= 32) {
      *value = *value >> 32;
    } else if (num_bits <= 40) {
      *value = *value >> 24;
    } else if (num_bits <= 48) {
      *value = *value >> 16;
    } else if (num_bits <= 56) {
      *value = *value >> 8;
    }
  }
  return res;
}

// TODO: Not sure which type to use in this case ...
/*uint128_t plc4c_spi_read_unsigned_big_integer(plc4c_spi_read_buffer* buf,
uint8_t num_bits) { return OK;
}*/

// Signed Integers ...

plc4c_return_code plc4c_spi_read_signed_byte(plc4c_spi_read_buffer* buf,
                                             uint8_t num_bits, int8_t* value) {
  return plc4c_spi_read_unsigned_byte(buf, num_bits, (uint8_t*)value);
}

plc4c_return_code plc4c_spi_read_signed_short(plc4c_spi_read_buffer* buf,
                                              uint8_t num_bits,
                                              int16_t* value) {
  return plc4c_spi_read_unsigned_byte(buf, num_bits, (uint16_t*)value);
}

plc4c_return_code plc4c_spi_read_signed_int(plc4c_spi_read_buffer* buf,
                                            uint8_t num_bits, int32_t* value) {
  return plc4c_spi_read_unsigned_byte(buf, num_bits, (uint32_t*)value);
}

plc4c_return_code plc4c_spi_read_signed_long(plc4c_spi_read_buffer* buf,
                                             uint8_t num_bits, int64_t* value) {
  return plc4c_spi_read_unsigned_byte(buf, num_bits, (uint64_t*)value);
}

// TODO: Not sure which type to use in this case ...
/*int128_t plc4c_spi_read_signed_big_integer(plc4c_spi_read_buffer* buf, uint8_t
 * num_bits); return OK;
 * }*/

// Floating Point Numbers ...

plc4c_return_code plc4c_spi_read_float(plc4c_spi_read_buffer* buf,
                                       uint8_t num_bits, float* value) {
  return OK;
}

plc4c_return_code plc4c_spi_read_double(plc4c_spi_read_buffer* buf,
                                        uint8_t num_bits, double* value) {
  return OK;
}

// TODO: Not sure which type to use in this case ...
/*doubledouble plc4c_spi_read_big_decimal(plc4c_spi_read_buffer* buf, uint8_t
 * num_bits); return 0;
 * } */

plc4c_return_code plc4c_spi_read_string(plc4c_spi_read_buffer* buf,
                                        uint8_t num_bits, char* encoding,
                                        char** value) {
  return OK;
}
