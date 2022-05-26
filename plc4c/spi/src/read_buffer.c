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
#include <plc4c/spi/system_private.h>
#include <string.h>
#include <math.h>

// This matrix contains constants for reading X bits starting with bit Y.
static const uint8_t read_bit_matrix[8][8] = {
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
  return (data & read_bit_matrix[num_bits - 1][from_bit]) >>
         (((unsigned int) 8) - (from_bit + num_bits));
}

uint8_t plc4c_spi_read_unsigned_byte_get_byte_internal(
    plc4c_spi_read_buffer* buf, uint8_t offset) {
  return *(buf->data + (buf->curPosByte + offset));
}

plc4c_return_code plc4c_spi_read_unsigned_bits_internal(
    plc4c_spi_read_buffer* buf, uint8_t num_bits, void* opaque_value) {

  // Cast void input to uint8_t, so we can walk the bytes, without
  // casting prior to calling this fcn. Size of values is really
  // given in the num_bits.
  uint8_t* value = (void*)opaque_value;

  if (buf == NULL) {
    return NULL_VALUE;
  }
  // Check if there are enough bytes in total left.
  if (!plc4c_spi_read_has_more(buf, num_bits)) {
    return OUT_OF_RANGE;
  }

  // If the bit-offset is currently 0, and we're reading
  // a full byte, go this shortcut.
  if ((buf->curPosBit == 0) && (num_bits % 8 == 0)) {
    if (buf->curPosByte > (buf->length - 1)) {
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
    // TODO: For debugging ... (Just that I can see the values in the debugger)
    uint8_t* original_value = value;

    // Find out how many bytes we need to read.
    uint8_t num_bytes_to_read = ((buf->curPosBit + num_bits) / 8) + 1;
    // Find out how many byte we need to write in the output.
    uint8_t num_bytes_to_write = (num_bits / 8) + ((num_bits % 8 != 0) ? 1 : 0);
    // Do a quick range check for the input
    // (for the output, the calling function is responsible)
    if ((buf->curPosByte + num_bytes_to_read - ((((buf->curPosBit + num_bits) % 8) == 0) ? 1 : 0)) > buf->length) {
      return OUT_OF_RANGE;
    }

    // If this is little endian, go to the end of the range as in this
    // case we have to fill the result from the back to the front.
    if (!plc4c_is_bigendian()) {
      value = value + (num_bytes_to_write - 1);
    }

    // Find out how many of the bits will be read from the first byte
    // It's actually just the rest of the byte as we checked if it all fits
    // in one byte in the else-block before.
    uint8_t num_bits_first_byte = 8 - buf->curPosBit;
    // Having read the bits from the first byte, see how many bits will
    // have to be read in the last byte (If we are finishing at a byte border,
    // the last byte will have 0 bits read).
    uint8_t num_bits_last_byte = (num_bits - num_bits_first_byte) % 8;
    // All in-between will obviously have all 8 bits read.

    // Read the bits of the first byte
    uint8_t cur_byte = plc4c_spi_read_unsigned_byte_get_byte_internal(buf, 0);

    // In the case that the number of bits read from the first and last
    // byte are more than 8, we got to put that excess data into its own
    // output byte.
    if(num_bits_first_byte + num_bits_last_byte > 8) {
      uint8_t excess_bits = num_bits_first_byte + num_bits_last_byte - 8;
      *value = plc4c_spi_read_unsigned_byte_internal(
          cur_byte, excess_bits, buf->curPosBit);
      // Move on to the next output byte
      value = plc4c_is_bigendian() ? value + 1 : value - 1;
      // Update the read-pointer.
      buf->curPosBit += excess_bits;
      // Change the number of bits read, as we already read some.
      num_bits_first_byte = num_bits_first_byte - excess_bits;
    }

    uint8_t high_level_part = plc4c_spi_read_unsigned_byte_internal(
        cur_byte, num_bits_first_byte, buf->curPosBit);
    // For each of the following bytes.
    for (int i = 1; i < num_bytes_to_read; i++) {
      // Shift the high level part by the amount of bits in the last byte
      *value = high_level_part << num_bits_last_byte;

      // We're done with this input byte, move on to the next one.
      buf->curPosByte++;

      // Get the next full byte.
      cur_byte = plc4c_spi_read_unsigned_byte_get_byte_internal(buf, 0);

      // Add the remaining bits of the current output byte.
      if (num_bits_last_byte != 0) {
        // Get the rest of the bits that belong to the previous byte.
        uint8_t low_level_part = plc4c_spi_read_unsigned_byte_internal(
            cur_byte, num_bits_last_byte, 0);
        // Add that to the end of the current output byte.
        *value = *value | low_level_part;

        // Here, we're finished reading the current output byte, so
        // move the output pointer to the next byte.
        value = plc4c_is_bigendian() ? value + 1 : value - 1;

        // The remaining parts of this byte will become the highest level
        // part of the next byte.
        high_level_part = plc4c_spi_read_unsigned_byte_internal(
            cur_byte, 8 - num_bits_last_byte, num_bits_last_byte);
      }
      // If this value happens to end at the end of a byte, there are no
      // remaining bits, so we can simply pass the current byte along.
      else {
        // In this case the last byte would contain 0 bits,
        // so we just abort here
        if(i == (num_bytes_to_read - 1)) {
          break;
        }
        // Here, we're finished reading the current output byte, so
        // move the output pointer to the next byte.
        value = plc4c_is_bigendian() ? value + 1 : value - 1;

        // Effectively this complete byte will become the next output byte.
        high_level_part = cur_byte;
      }
    }

    // Update the buffer position
    buf->curPosBit = num_bits_last_byte;
    return OK;
  }
}

bool plc4c_spi_read_buffer_is_negative_internal(uint8_t num_bits, int8_t value) {
  int8_t tmp_value = value >> num_bits;
  return (tmp_value & 1) != 0;
}

bool plc4c_spi_fill_sign_internal(uint8_t num_bits, int8_t* value) {
  // Find out how many bytes the value has.
  uint8_t num_bytes_total = (num_bits / 8) + ((num_bits % 8 != 0) ? 1 : 0);

  // If this is big endian, go to the highest level byte.
  if (!plc4c_is_bigendian()) {
    value = value + (num_bytes_total - 1);
  }

  if(plc4c_spi_read_buffer_is_negative_internal((num_bits - 1) % 8, *value)) {
    // Set all bits above {num_bits} to 1
    int8_t tmp_value = *value;
    if(num_bits % 8 != 0) {
      tmp_value = tmp_value | read_bit_matrix[7 - (num_bits % 8)][0];
    }
    *value = tmp_value;
    return true;
  }
  return false;
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
  uint8_t cur_byte = *(buf->data + buf->curPosByte);
  // We have to invert the position as bit 0 will be the first
  // (most significant bit).
  unsigned int bit_pos = ((unsigned int)7) - buf->curPosBit;
  // Get the bit's value.
  *value = ((cur_byte >> bit_pos) & 1) != 0;
  // If this was the last bit in this byte, move on to the next one.
  if (buf->curPosBit == (unsigned int) 7) {
    buf->curPosByte++;
    buf->curPosBit = 0;
  } else {
    buf->curPosBit++;
  }
  return OK;
}

plc4c_return_code plc4c_spi_read_char(plc4c_spi_read_buffer* buf, char* value) {
    return plc4c_spi_read_signed_int(buf, 8, (int8_t*) value);
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
  if ((res == OK) && plc4c_is_bigendian()) {
    if (num_bits <= 8) {
      *value >>= 8;
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
  if ((res == OK) && plc4c_is_bigendian()) {
    if (num_bits <= 8) {
      *value >>= 24;
    } else if (num_bits <= 16) {
      *value >>= 16;
    } else if (num_bits <= 24) {
      *value >>= 8;
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
  plc4c_return_code res = plc4c_spi_read_unsigned_bits_internal(buf, num_bits, value);
  // Shift the bits to the right position.
  if ((res == OK) && plc4c_is_bigendian()) {
    if (num_bits <= 8) {
      *value >>= 56;
    } else if (num_bits <= 16) {
      *value >>= 48;
    } else if (num_bits <= 24) {
      *value >>= 40;
    } else if (num_bits <= 32) {
      *value >>= 32;
    } else if (num_bits <= 40) {
      *value >>= 24;
    } else if (num_bits <= 48) {
      *value >>= 16;
    } else if (num_bits <= 56) {
      *value >>= 8;
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

  plc4c_return_code res = plc4c_spi_read_unsigned_byte(buf, num_bits, (uint8_t*) value);
  if(res == OK) {
    plc4c_spi_fill_sign_internal(num_bits, value);
  }
  return res;
}

plc4c_return_code plc4c_spi_read_signed_short(plc4c_spi_read_buffer* buf,
                                              uint8_t num_bits,
                                              int16_t* value) {

  plc4c_return_code res = plc4c_spi_read_unsigned_short(buf, num_bits, (uint16_t*) value);
  if(res == OK) {
    if(plc4c_spi_fill_sign_internal(num_bits, (int8_t*) value)) {
      // Potentially fill all higher level bytes with 255
      if(num_bits <= 8) {
        *value |= 0xFF00; // use hex, its less like a magic number and typos easier to see
      }
    } else {
      // Potentially fill all higher level bytes with 0
      if(num_bits <= 8) {
        *value &= 0xFF;
      }
    }
  }
  return res;
}

plc4c_return_code plc4c_spi_read_signed_int(plc4c_spi_read_buffer* buf,
                                            uint8_t num_bits, int32_t* value) {

  plc4c_return_code res = plc4c_spi_read_unsigned_int(buf, num_bits, (uint32_t*) value);
  if(res == OK) {
    if(plc4c_spi_fill_sign_internal(num_bits, (int8_t*) value)) {
      // Potentially fill all higher level bytes with 255
      if(num_bits <= 8) {
        *value |= 0xFFFFFF00; 
      } else if(num_bits <= 16) {
        *value |= 0xFFFF0000; 
      } else if(num_bits <= 24) {
        *value |= 0xFF000000;
      }
    } else {
      // Potentially fill all higher level bytes with 0
      if(num_bits <= 8) {
        *value &= 0xFF;
      } else if(num_bits <= 16) {
        *value &= 0xFFFF ;
      } else if(num_bits <= 24) {
        *value &= 0xFFFFFF;
      }
    }
  }
  return res;
}

plc4c_return_code plc4c_spi_read_signed_long(plc4c_spi_read_buffer* buf,
                                             uint8_t num_bits, int64_t* value) {

  plc4c_return_code res = plc4c_spi_read_unsigned_long(buf, num_bits, (uint64_t*) value);
  if(res == OK) {
    if(plc4c_spi_fill_sign_internal(num_bits, (int8_t*) value)) {
      // Potentially fill all higher level bytes with 255
      if(num_bits <= 8) {
        *value |= 0xFFFFFFFFFFFFFF00;
      } else if(num_bits <= 16) {
        *value |= 0xFFFFFFFFFFFF0000;
      } else if(num_bits <= 24) {
        *value |= 0xFFFFFFFFFF000000;
      } else if(num_bits <= 32) {
        *value |= 0xFFFFFFFF00000000;
      } else if(num_bits <= 40) {
        *value |= 0xFFFFFF0000000000;
      } else if(num_bits <= 48) {
        *value |= 0xFFFF000000000000;
      } else if(num_bits <= 56) {
        *value |= 0xFF00000000000000;
      }
    } else {
      // Potentially fill all higher level bytes with 0
      if(num_bits <= 8) {
        *value &=  0xFF;
      } else if(num_bits <= 16) {
        *value &= 0xFFFF;
      } else if(num_bits <= 24) {
        *value &= 0xFFFFFF;
      } else if(num_bits <= 32) {
        *value &= 0xFFFFFFFF;
      } else if(num_bits <= 40) {
        *value &= 0xFFFFFFFFFF;
      } else if(num_bits <= 48) {
        *value &= 0xFFFFFFFFFFFF;
      } else if(num_bits <= 56) {
        *value &= 0xFFFFFFFFFFFFFF;
      }
    }
  }
  return res;
}

// TODO: Not sure which type to use in this case ...
/*int128_t plc4c_spi_read_signed_big_integer(plc4c_spi_read_buffer* buf, uint8_t
 * num_bits); return OK;
 * }*/

// Floating Point Numbers ...

plc4c_return_code plc4c_spi_read_float(plc4c_spi_read_buffer* buf,
                                       uint8_t num_bits, float* value) {
  if(num_bits == 16) {
      // https://en.wikipedia.org/wiki/Half-precision_floating-point_format
      bool sign = false;
      plc4c_return_code res = plc4c_spi_read_bit(buf, &sign);
      if(res != OK) {
        return res;
      }
      uint8_t exponent = 0;
      res = plc4c_spi_read_unsigned_byte(buf, 5, &exponent);
      if(res != OK) {
        return res;
      }
      uint16_t fraction = 0;
      res = plc4c_spi_read_unsigned_short(buf, 10, &fraction);
      if(res != OK) {
        return res;
      }

      if((exponent >= 1) && (exponent <= 30)) {
        *value = (sign ? (float) 1 : (float) -1) * ((float) (2 ^ (exponent - 15))) * ((float) (1 + (fraction / 10)));
      } else if(exponent == 0) {
        if (fraction == 0) {
          *value = 0.0f;
        } else {
          *value = (sign ? (float) 1 : (float) -1) * ((float) (2 ^ (-14))) * ((float) (fraction / 10));
        }
      } else if(exponent == 31) {
        if (fraction == 0) {
          *value = sign ? INFINITY : -INFINITY;
        } else {
          *value = NAN;
        }
      } else {
        return INVALID_ARGUMENT;
      }
  } else if(num_bits == 32) {
      plc4c_return_code res = plc4c_spi_read_unsigned_int(buf, 32, (uint32_t*) value);
      if(res != OK) {
          return res;
      }
  } else {
      return INVALID_ARGUMENT;
  }
  return OK;
}

plc4c_return_code plc4c_spi_read_double(plc4c_spi_read_buffer* buf,
                                        uint8_t num_bits, double* value) {
  if(num_bits == 64) {
    plc4c_return_code res = plc4c_spi_read_unsigned_long(buf, 64, (uint64_t*) value);
    if(res != OK) {
      return res;
    }
  } else {
    return INVALID_ARGUMENT;
  }
  return OK;
}

// TODO: Not sure which type to use in this case ...
/*doubledouble plc4c_spi_read_big_decimal(plc4c_spi_read_buffer* buf, uint8_t
 * num_bits); return 0;
 * } */

plc4c_return_code plc4c_spi_read_string(plc4c_spi_read_buffer* buf,
                                        uint8_t num_bits, char* encoding,
                                        char** value) {

  // Right now we only support utf-8.
  if(strcmp(encoding,"UTF-8") != 0) {
      return INVALID_ARGUMENT;
  }

  // Allocate enough chars to contain the string and termination character.
  char* str = malloc(sizeof(char) * ((num_bits / 8) + 1));
  char* cur_str = str;
  // Read all the bytes one by one.
  for(int i = 0; (i < (num_bits / 8)) && plc4c_spi_read_has_more(buf, 8); i++) {
    plc4c_spi_read_unsigned_byte(buf, 8, (uint8_t*) cur_str);
    cur_str++;
  }
  // Terminate the string.
  *cur_str = '\0';
  *value = str;
  return OK;
}
