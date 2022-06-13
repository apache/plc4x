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

#include <plc4c/spi/system_private.h>
#include <plc4c/spi/write_buffer.h>
#include <string.h>

// This matrix contains constants for reading X bits starting with bit Y.
static const uint8_t write_bit_matrix[8][8] = {
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

void plc4c_spi_write_unsigned_byte_internal(uint8_t* output_data,
                                            uint8_t num_bits, uint8_t from_bit,
                                            uint8_t value) {
  if (num_bits + from_bit > 8) {
    return;
  }
  uint8_t shifted_value = value << (((unsigned int)8) - (from_bit + num_bits));
  *output_data =
      *output_data | (shifted_value & write_bit_matrix[num_bits - 1][from_bit]);
}

uint8_t plc4c_spi_write_get_byte_internal(plc4c_spi_write_buffer* buf,
                                          uint8_t offset) {
  return *(buf->data + (buf->curPosByte + offset));
}

void plc4c_spi_write_put_byte_internal(plc4c_spi_write_buffer* buf,
                                       uint8_t offset, uint8_t value) {
  *(buf->data + (buf->curPosByte + offset)) = value;
}

plc4c_return_code plc4c_spi_write_unsigned_bits_internal(
    plc4c_spi_write_buffer* buf, uint8_t num_bits, uint8_t* value) {
  if (buf == NULL) {
    return NULL_VALUE;
  }
  if (buf->curPosByte + ((buf->curPosBit + num_bits) / 8) > (buf->length)) {
    return OUT_OF_RANGE;
  }

  // If the bit-offset is currently 0 and we're writing
  // a full byte, go this shortcut.
  if ((buf->curPosBit == 0) && (num_bits % 8 == 0)) {
    // Find how many full bytes we'll be writing.
    uint8_t num_bytes = num_bits / 8;
    // If this is little endian, go to the end of the range.
    if (!plc4c_is_bigendian()) {
      value = value + (num_bytes - 1);
    }
    // Write each of these.
    for (int i = 0; i < num_bytes; i++) {
      plc4c_spi_write_put_byte_internal(buf, 0, *value);
      // Move the write-pointer to the next byte.
      buf->curPosByte++;
      value = plc4c_is_bigendian() ? value + 1 : value - 1;
    }
    return OK;
  }

  // in this case the current byte alone is enough to service this request.
  else if ((((unsigned int)8) - buf->curPosBit) >= num_bits) {
    plc4c_spi_write_unsigned_byte_internal((buf->data + buf->curPosByte),
                                           num_bits, buf->curPosBit, *value);
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
    // Find out how many bytes we need to write to the output.
    uint8_t num_write_bytes = ((buf->curPosBit + num_bits) / 8) + 1;
    // Do a quick range check for the input
    // (for the output, the calling function is responsible)
    if ((buf->curPosByte + num_write_bytes -
         ((((buf->curPosBit + num_bits) % 8) == 0) ? 1 : 0)) > buf->length) {
      return OUT_OF_RANGE;
    }

    // Find out how many byte we need to read from the input.
    uint8_t num_read_bytes = (num_bits / 8) + ((num_bits % 8 != 0) ? 1 : 0);
    // The first read-byte is the only one that can be incomplete.
    uint8_t num_read_bits_first_byte = num_bits % 8;
    if(num_read_bits_first_byte == 0) {
      num_read_bits_first_byte = 8;
    }
    // All others will obviously have all 8 bits read.

    // Find out how many of the bits will be written to the first byte
    // It's actually just the rest of the byte as we checked if it all fits
    // in one byte in the else-block before.
    uint8_t num_write_bits_first_byte = 8 - buf->curPosBit;
    // Having written the bits from the first byte, see how many bits will
    // have to be written in the last byte (If we are finishing at a byte
    // border, the last byte will have 0 bits written.
    uint8_t num_write_bits_last_byte =
        (num_bits - num_write_bits_first_byte) % 8;
    // All in-between will obviously have all 8 bits written.

    // If this is little endian, go to the end of the range of the read input,
    // as in this case we have to read the value from the back to the front.
    if (!plc4c_is_bigendian()) {
      value = value + (num_read_bytes - 1);
    }

    // For the first byte the end will be as much as fits in the
    // current byte, after this this will be updated to the inverse
    // of the bits written to the last byte.
    uint8_t num_bits_end_of_write_byte = num_write_bits_first_byte;

    // If the number of read-bits of the first read-byte aren't enough
    // to fill the rest of the first write-byte, we need to start reading
    // all of this first and then continue filling up with the next full
    // read-byte.
    if (num_read_bits_first_byte < num_write_bits_first_byte) {
      // Calculate the number of bits we can't fill with the first read-byte.
      uint8_t remaining_bits =
          num_write_bits_first_byte - num_read_bits_first_byte;
      // Write the content of the first read byte to the output.
      plc4c_spi_write_unsigned_byte_internal(buf->data + buf->curPosByte,
                                             num_read_bits_first_byte,
                                             buf->curPosBit,
                                             *value);

      // Move to the next byte in the input.
      value = plc4c_is_bigendian() ? value + 1 : value - 1;

      // Update the read-pointer.
      buf->curPosBit += num_read_bits_first_byte;

      // Decrease the number of write bits for the first byte.
      num_bits_end_of_write_byte = remaining_bits;

      // Decrement the number of bits to write in total.
      num_bits -= num_read_bits_first_byte;
    }

    // For each of the following bytes.
    while (num_bits > 0) {
      // Output the first part of the current read-byte as last part of the
      // current output byte.
      uint8_t fragment = *value >> num_write_bits_last_byte;
      plc4c_spi_write_unsigned_byte_internal(buf->data + buf->curPosByte,
                                             num_bits_end_of_write_byte,
                                             buf->curPosBit, fragment);
      // Move the write buffer to the next byte and reset the bit position.
      buf->curPosByte++;
      buf->curPosBit = 0;

      // Decrement the number of remaining bits.
      num_bits -= num_bits_end_of_write_byte;

      // From now on the end of the current write bit will be the rest
      // of what doesn't fit into the last byte.
      num_bits_end_of_write_byte = 8 - num_write_bits_last_byte;

      // Only if there are remaining bits, continue writing.
      if (num_bits > 0) {
        fragment =
            (*value & write_bit_matrix[num_write_bits_last_byte - 1][8 - num_write_bits_last_byte]);
        plc4c_spi_write_unsigned_byte_internal(buf->data + buf->curPosByte,
                                               num_write_bits_last_byte,
                                               buf->curPosBit, fragment);

        // Move the write buffer to the next bit position.
        buf->curPosBit = num_write_bits_last_byte;

        // Decrement the number of remaining bits.
        num_bits -= num_write_bits_last_byte;

        if (num_bits > 0) {
          // Move to the next byte in the input.
          value = plc4c_is_bigendian() ? value + 1 : value - 1;
        }
      }
    }
    return OK;
  }
}

plc4c_return_code plc4c_spi_write_buffer_create(
    uint16_t length, plc4c_spi_write_buffer** buffer) {
  *buffer = malloc(sizeof(plc4c_spi_write_buffer));
  if (*buffer == NULL) {
    return NO_MEMORY;
  }

  (*buffer)->data = calloc(length, sizeof(uint8_t));
  if ((*buffer)->data == NULL) {
    return NO_MEMORY;
  }
  (*buffer)->length = length;
  (*buffer)->curPosByte = 0;
  (*buffer)->curPosBit = 0;

  return OK;
}

void plc4c_spi_write_buffer_destroy(plc4c_spi_write_buffer* buffer) {
  free(buffer);
}

uint8_t* plc4c_spi_write_get_data(plc4c_spi_write_buffer* buf) {
  return buf->data;
}

uint32_t plc4c_spi_write_get_pos(plc4c_spi_write_buffer* buf) {
  return buf->curPosByte;
}

plc4c_return_code plc4c_spi_write_get_bytes(plc4c_spi_write_buffer* buf,
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

  // Get a 0-initialized buffer.
  *dest = calloc(num_bytes, sizeof(uint8_t));
  if (*dest == NULL) {
    return NO_MEMORY;
  }

  // Copy the requested bytes to the output.
  memcpy(*dest, buf->data, num_bytes);
  return OK;
}

plc4c_return_code plc4c_spi_write_bit(plc4c_spi_write_buffer* buf, bool value) {
  // Only if the value is "true" will the content of the
  // buffer look any different.
  if (value) {
    uint8_t cur_byte = *(buf->data + buf->curPosByte);
    // We have to invert the position as bit 0 will be the first
    // (most significant bit).
    unsigned int bit_pos = ((unsigned int)7) - buf->curPosBit;
    uint8_t bitValue = ((uint8_t)1) << bit_pos;
    *(buf->data + buf->curPosByte) = cur_byte | bitValue;
  }
  // If this was the last bit in this byte, move on to the next one.
  if (buf->curPosBit == (unsigned int)7) {
    buf->curPosByte++;
    buf->curPosBit = 0;
  } else {
    buf->curPosBit++;
  }
  return OK;
}

plc4c_return_code plc4c_spi_write_char(plc4c_spi_write_buffer* buf, char value) {
    return plc4c_spi_write_signed_int(buf, 8, (int8_t) value);
}

// Unsigned Integers ...

plc4c_return_code plc4c_spi_write_unsigned_byte(plc4c_spi_write_buffer* buf,
                                                uint8_t num_bits,
                                                uint8_t value) {
  // If more than 8 bits are requested, return an error.
  if (num_bits > 8) {
    return OUT_OF_RANGE;
  }
  // Write the bits.
  return plc4c_spi_write_unsigned_bits_internal(buf, num_bits, &value);
}

plc4c_return_code plc4c_spi_write_unsigned_short(plc4c_spi_write_buffer* buf,
                                                 uint8_t num_bits,
                                                 uint16_t value) {
  // If more than 16 bits are requested, return an error.
  if (num_bits > 16) {
    return OUT_OF_RANGE;
  }
  // Write the bits.
  return plc4c_spi_write_unsigned_bits_internal(buf, num_bits, &value);
}

plc4c_return_code plc4c_spi_write_unsigned_int(plc4c_spi_write_buffer* buf,
                                               uint8_t num_bits,
                                               uint32_t value) {
  // If more than 32 bits are requested, return an error.
  if (num_bits > 32) {
    return OUT_OF_RANGE;
  }
  // Write the bits.
  return plc4c_spi_write_unsigned_bits_internal(buf, num_bits, &value);
}

plc4c_return_code plc4c_spi_write_unsigned_long(plc4c_spi_write_buffer* buf,
                                                uint8_t num_bits,
                                                uint64_t value) {
  // If more than 64 bits are requested, return an error.
  if (num_bits > 64) {
    return OUT_OF_RANGE;
  }
  // Write the bits.
  return plc4c_spi_write_unsigned_bits_internal(buf, num_bits, &value);
}

// TODO: Not sure which type to use in this case ...
/*plc4c_return_code
 * plc4c_spi_write_unsigned_big_integer(plc4c_spi_write_buffer* buf, uint8_t
 * num_bits, uint128_t value) {
 * } */

// Signed Integers ...

plc4c_return_code plc4c_spi_write_signed_byte(plc4c_spi_write_buffer* buf,
                                              uint8_t num_bits, int8_t value) {
  return plc4c_spi_write_unsigned_byte(buf, num_bits, (uint8_t) value);
}

plc4c_return_code plc4c_spi_write_signed_short(plc4c_spi_write_buffer* buf,
                                               uint8_t num_bits,
                                               int16_t value) {
  return plc4c_spi_write_unsigned_short(buf, num_bits, (uint16_t) value);
}

plc4c_return_code plc4c_spi_write_signed_int(plc4c_spi_write_buffer* buf,
                                             uint8_t num_bits, int32_t value) {
  return plc4c_spi_write_unsigned_int(buf, num_bits, (uint32_t) value);
}

plc4c_return_code plc4c_spi_write_signed_long(plc4c_spi_write_buffer* buf,
                                              uint8_t num_bits, int64_t value) {
  return plc4c_spi_write_unsigned_long(buf, num_bits, (uint64_t) value);
}

// TODO: Not sure which type to use in this case ...
/*plc4c_return_code plc4c_spi_write_signed_big_integer(plc4c_spi_write_buffer*
 * buf, uint8_t num_bits, int128_t) {
 * } */

// Floating Point Numbers ...

plc4c_return_code plc4c_spi_write_float(plc4c_spi_write_buffer* buf,
                                        uint8_t num_bits, float value) {
  // Half precision floats (16 bit) are currently not implemented.
  if(num_bits != 32) {
    return NOT_IMPLEMENTED;
  }
  // Use this little helper to convert the 32 bit
  // float into a 32 bit unsigned int.
  union {
    float f;
    uint32_t u;
  } helper;
  helper.f = value;
  return plc4c_spi_write_unsigned_int(buf, num_bits, helper.u);
}

plc4c_return_code plc4c_spi_write_double(plc4c_spi_write_buffer* buf,
                                         uint8_t num_bits, double value) {
  if(num_bits != 64) {
    return NOT_IMPLEMENTED;
  }
  // Use this little helper to convert the 64 bit
  // float into a 64 bit unsigned int.
  union {
    double d;
    uint64_t u;
  } helper;
  helper.d = value;
  return plc4c_spi_write_unsigned_long(buf, num_bits, (uint64_t) helper.u);
}

// TODO: Not sure which type to use in this case ...
/*void plc4c_spi_write_big_decimal(plc4c_spi_write_buffer* buf, uint8_t
 * num_bits, doubledouble value) {
 * } */

plc4c_return_code plc4c_spi_write_string(plc4c_spi_write_buffer* buf,
                                         uint8_t num_bits, char* encoding,
                                         char* value) {
  // Right now we only support utf-8 and utf-16.
  if((strcmp(encoding,"UTF-8") != 0) && (strcmp(encoding,"UTF-16") != 0)) {
    return INVALID_ARGUMENT;
  }
  // Simply output the bytes to the buffer.
  for(int i = 0; (i < (num_bits / 8)); i++) {
    plc4c_spi_write_unsigned_byte(buf, 8, (uint8_t*) *value);
    value++;
  }
  return OK;
}
