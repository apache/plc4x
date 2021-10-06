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

#ifndef PLC4C_DRIVER_S7_STATIC_H
#define PLC4C_DRIVER_S7_STATIC_H

#include <plc4c/spi/read_buffer.h>
#include <plc4c/spi/write_buffer.h>
#include <time.h>

/*
 *
 *   Static functions
 *
 */

uint16_t plc4c_s7_read_write_event_item_length(plc4c_spi_read_buffer* io, uint16_t valueLength);

uint16_t plc4c_s7_read_write_right_shift3(plc4c_spi_read_buffer* io);

plc4c_return_code plc4c_s7_read_write_left_shift3(plc4c_spi_write_buffer* io, uint16_t value);

uint8_t plc4c_s7_read_write_bcd_to_int(plc4c_spi_read_buffer* io);

plc4c_return_code plc4c_s7_read_write_byte_to_bcd(plc4c_spi_write_buffer* writeBuffer, uint8_t value);

uint16_t plc4c_s7_read_write_s7msec_to_int(plc4c_spi_read_buffer* io);

plc4c_return_code plc4c_s7_read_write_int_to_s7msec(plc4c_spi_write_buffer* writeBuffer, uint16_t value);

char* plc4c_s7_read_write_parse_s7_string(plc4c_spi_read_buffer* io, int32_t stringLength, char* encoding);

char* plc4c_s7_read_write_parse_s7_char(plc4c_spi_read_buffer* io, char* encoding);

time_t plc4c_s7_read_write_parse_tia_time(plc4c_spi_read_buffer* io);

time_t plc4c_s7_read_write_parse_s5_time(plc4c_spi_read_buffer* io);

time_t plc4c_s7_read_write_parse_tia_l_time(plc4c_spi_read_buffer* io);

time_t plc4c_s7_read_write_parse_tia_date(plc4c_spi_read_buffer* io);

time_t plc4c_s7_read_write_parse_tia_time_of_day(plc4c_spi_read_buffer* io);

time_t plc4c_s7_read_write_parse_tia_date_time(plc4c_spi_read_buffer* io);

#endif  // PLC4C_DRIVER_S7_STATIC_H
