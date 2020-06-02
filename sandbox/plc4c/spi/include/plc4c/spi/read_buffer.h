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
#ifndef PLC4C_READ_BUFFER_H_
#define PLC4C_READ_BUFFER_H_

#include <stdint.h>
#include <plc4c/system.h>
#include <plc4c/types.h>

struct plc4c_spi_read_buffer {

};
typedef struct plc4c_spi_read_buffer plc4c_spi_read_buffer;

uint32_t plc4c_spi_read_get_pos(plc4c_spi_read_buffer* buf);

uint32_t plc4c_spi_read_get_total_bytes(plc4c_spi_read_buffer* buf);

bool plc4c_spi_read_has_more(plc4c_spi_read_buffer* buf, uint8_t num_bits);

uint8_t* plc4c_spi_read_get_bytes(plc4c_spi_read_buffer* buf, uint32_t start_pos_in_bytes, uint32_t end_pos_in_bytes);

uint8_t plc4c_spi_read_peek_byte(plc4c_spi_read_buffer* buf, uint32_t offset_in_bytes);

bool plc4c_spi_read_bit(plc4c_spi_read_buffer* buf);

// Unsigned Integers ...

uint8_t plc4c_spi_read_unsigned_byte(plc4c_spi_read_buffer* buf, uint8_t num_bits);

uint16_t plc4c_spi_read_unsigned_short(plc4c_spi_read_buffer* buf, uint8_t num_bits);

uint32_t plc4c_spi_read_unsigned_int(plc4c_spi_read_buffer* buf, uint8_t num_bits);

uint64_t plc4c_spi_read_unsigned_long(plc4c_spi_read_buffer* buf, uint8_t num_bits);

// TODO: Not sure which type to use in this case ...
//uint128_t plc4c_spi_read_unsigned_big_integer(plc4c_spi_read_buffer* buf, uint8_t num_bits);

// Signed Integers ...

int8_t plc4c_spi_read_byte(plc4c_spi_read_buffer* buf, uint8_t num_bits);

int16_t plc4c_spi_read_short(plc4c_spi_read_buffer* buf, uint8_t num_bits);

int32_t plc4c_spi_read_int(plc4c_spi_read_buffer* buf, uint8_t num_bits);

int64_t plc4c_spi_read_long(plc4c_spi_read_buffer* buf, uint8_t num_bits);

// TODO: Not sure which type to use in this case ...
//int128_t plc4c_spi_read_big_integer(plc4c_spi_read_buffer* buf, uint8_t num_bits);

// Floating Point Numbers ...

float plc4c_spi_read_float(plc4c_spi_read_buffer* buf, uint8_t num_bits);

double plc4c_spi_read_double(plc4c_spi_read_buffer* buf, uint8_t num_bits);

// TODO: Not sure which type to use in this case ...
//doubledouble plc4c_spi_read_big_decimal(plc4c_spi_read_buffer* buf, uint8_t num_bits);

char* plc4c_spi_read_string(plc4c_spi_read_buffer* buf, uint8_t num_bits, char* encoding);

#endif  // PLC4C_READ_BUFFER_H_