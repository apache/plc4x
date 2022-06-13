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
#ifndef PLC4C_WRITE_BUFFER_H_
#define PLC4C_WRITE_BUFFER_H_

#include <stdint.h>
#include <plc4c/system.h>
#include <plc4c/types.h>

struct plc4c_spi_write_buffer {
  // Pointer to the data itself
  uint8_t* data;
  // Total size of the data array.
  uint16_t length;
  // Current full byte position
  uint16_t curPosByte;
  // Current bit-position inside the current byte
  unsigned int curPosBit : 4;
};
typedef struct plc4c_spi_write_buffer plc4c_spi_write_buffer;

plc4c_return_code plc4c_spi_write_buffer_create(uint16_t length, plc4c_spi_write_buffer** buffer);

void plc4c_spi_write_buffer_destroy(plc4c_spi_write_buffer* buffer);

uint8_t* plc4c_spi_write_get_data(plc4c_spi_write_buffer* buf);

uint32_t plc4c_spi_write_get_pos(plc4c_spi_write_buffer* buf);

plc4c_return_code plc4c_spi_write_get_bytes(plc4c_spi_write_buffer* buf, uint16_t start_pos_in_bytes, uint16_t end_pos_in_bytes, uint8_t** dest);

plc4c_return_code plc4c_spi_write_bit(plc4c_spi_write_buffer* buf, bool value);

plc4c_return_code plc4c_spi_write_char(plc4c_spi_write_buffer* buf, char value);

// Unsigned Integers ...

plc4c_return_code plc4c_spi_write_unsigned_byte(plc4c_spi_write_buffer* buf, uint8_t num_bits, uint8_t value);

plc4c_return_code plc4c_spi_write_unsigned_short(plc4c_spi_write_buffer* buf, uint8_t num_bits, uint16_t value);

plc4c_return_code plc4c_spi_write_unsigned_int(plc4c_spi_write_buffer* buf, uint8_t num_bits, uint32_t value);

plc4c_return_code plc4c_spi_write_unsigned_long(plc4c_spi_write_buffer* buf, uint8_t num_bits, uint64_t value);

// TODO: Not sure which type to use in this case ...
//void plc4c_spi_write_unsigned_big_integer(plc4c_spi_write_buffer* buf, uint8_t num_bits, uint128_t value);

// Signed Integers ...

plc4c_return_code plc4c_spi_write_signed_byte(plc4c_spi_write_buffer* buf, uint8_t num_bits, int8_t value);

plc4c_return_code plc4c_spi_write_signed_short(plc4c_spi_write_buffer* buf, uint8_t num_bits, int16_t value);

plc4c_return_code plc4c_spi_write_signed_int(plc4c_spi_write_buffer* buf, uint8_t num_bits, int32_t value);

plc4c_return_code plc4c_spi_write_signed_long(plc4c_spi_write_buffer* buf, uint8_t num_bits, int64_t value);

// TODO: Not sure which type to use in this case ...
//plc4c_return_code plc4c_spi_write_signed_big_integer(plc4c_spi_write_buffer* buf, uint8_t num_bits, int128_t);

// Floating Point Numbers ...

plc4c_return_code plc4c_spi_write_float(plc4c_spi_write_buffer* buf, uint8_t num_bits, float value);

plc4c_return_code plc4c_spi_write_double(plc4c_spi_write_buffer* buf, uint8_t num_bits, double value);

// TODO: Not sure which type to use in this case ...
//plc4c_return_code plc4c_spi_write_big_decimal(plc4c_spi_write_buffer* buf, uint8_t num_bits, doubledouble value);

plc4c_return_code plc4c_spi_write_string(plc4c_spi_write_buffer* buf, uint8_t num_bits, char* encoding, char* value);

#endif  // PLC4C_WRITE_BUFFER_H_