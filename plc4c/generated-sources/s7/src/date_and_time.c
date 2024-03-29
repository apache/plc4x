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

#include <stdio.h>
#include <plc4c/spi/context.h>
#include <plc4c/spi/evaluation_helper.h>
#include <plc4c/driver_s7_static.h>

#include "date_and_time.h"

// Code generated by code-generation. DO NOT EDIT.


// Parse function.
plc4c_return_code plc4c_s7_read_write_date_and_time_parse(plc4x_spi_context ctx, plc4c_spi_read_buffer* readBuffer, plc4c_s7_read_write_date_and_time** _message) {
  uint16_t startPos = plc4c_spi_read_get_pos(readBuffer);
  plc4c_return_code _res = OK;

  // Allocate enough memory to contain this data structure.
  (*_message) = malloc(sizeof(plc4c_s7_read_write_date_and_time));
  if(*_message == NULL) {
    return NO_MEMORY;
  }

  // Simple Field (year)
  uint8_t year = 0;
  _res = plc4c_spi_read_unsigned_byte(readBuffer, 8, (uint8_t*) &year);
  if(_res != OK) {
    return _res;
  }
  (*_message)->year = year;

  // Simple Field (month)
  uint8_t month = 0;
  _res = plc4c_spi_read_unsigned_byte(readBuffer, 8, (uint8_t*) &month);
  if(_res != OK) {
    return _res;
  }
  (*_message)->month = month;

  // Simple Field (day)
  uint8_t day = 0;
  _res = plc4c_spi_read_unsigned_byte(readBuffer, 8, (uint8_t*) &day);
  if(_res != OK) {
    return _res;
  }
  (*_message)->day = day;

  // Simple Field (hour)
  uint8_t hour = 0;
  _res = plc4c_spi_read_unsigned_byte(readBuffer, 8, (uint8_t*) &hour);
  if(_res != OK) {
    return _res;
  }
  (*_message)->hour = hour;

  // Simple Field (minutes)
  uint8_t minutes = 0;
  _res = plc4c_spi_read_unsigned_byte(readBuffer, 8, (uint8_t*) &minutes);
  if(_res != OK) {
    return _res;
  }
  (*_message)->minutes = minutes;

  // Simple Field (seconds)
  uint8_t seconds = 0;
  _res = plc4c_spi_read_unsigned_byte(readBuffer, 8, (uint8_t*) &seconds);
  if(_res != OK) {
    return _res;
  }
  (*_message)->seconds = seconds;

  // Simple Field (msec)
  uint16_t msec = 0;
  _res = plc4c_spi_read_unsigned_short(readBuffer, 12, (uint16_t*) &msec);
  if(_res != OK) {
    return _res;
  }
  (*_message)->msec = msec;

  // Simple Field (dow)
  uint8_t dow = 0;
  _res = plc4c_spi_read_unsigned_byte(readBuffer, 4, (uint8_t*) &dow);
  if(_res != OK) {
    return _res;
  }
  (*_message)->dow = dow;

  return OK;
}

plc4c_return_code plc4c_s7_read_write_date_and_time_serialize(plc4x_spi_context ctx, plc4c_spi_write_buffer* writeBuffer, plc4c_s7_read_write_date_and_time* _message) {
  plc4c_return_code _res = OK;

  // Simple Field (year)
  _res = plc4c_spi_write_unsigned_byte(writeBuffer, 8, _message->year);
  if(_res != OK) {
    return _res;
  }

  // Simple Field (month)
  _res = plc4c_spi_write_unsigned_byte(writeBuffer, 8, _message->month);
  if(_res != OK) {
    return _res;
  }

  // Simple Field (day)
  _res = plc4c_spi_write_unsigned_byte(writeBuffer, 8, _message->day);
  if(_res != OK) {
    return _res;
  }

  // Simple Field (hour)
  _res = plc4c_spi_write_unsigned_byte(writeBuffer, 8, _message->hour);
  if(_res != OK) {
    return _res;
  }

  // Simple Field (minutes)
  _res = plc4c_spi_write_unsigned_byte(writeBuffer, 8, _message->minutes);
  if(_res != OK) {
    return _res;
  }

  // Simple Field (seconds)
  _res = plc4c_spi_write_unsigned_byte(writeBuffer, 8, _message->seconds);
  if(_res != OK) {
    return _res;
  }

  // Simple Field (msec)
  _res = plc4c_spi_write_unsigned_short(writeBuffer, 12, _message->msec);
  if(_res != OK) {
    return _res;
  }

  // Simple Field (dow)
  _res = plc4c_spi_write_unsigned_byte(writeBuffer, 4, _message->dow);
  if(_res != OK) {
    return _res;
  }

  return OK;
}

uint16_t plc4c_s7_read_write_date_and_time_length_in_bytes(plc4x_spi_context ctx, plc4c_s7_read_write_date_and_time* _message) {
  return plc4c_s7_read_write_date_and_time_length_in_bits(ctx, _message) / 8;
}

uint16_t plc4c_s7_read_write_date_and_time_length_in_bits(plc4x_spi_context ctx, plc4c_s7_read_write_date_and_time* _message) {
  uint16_t lengthInBits = 0;

  // Simple field (year)
  lengthInBits += 8;

  // Simple field (month)
  lengthInBits += 8;

  // Simple field (day)
  lengthInBits += 8;

  // Simple field (hour)
  lengthInBits += 8;

  // Simple field (minutes)
  lengthInBits += 8;

  // Simple field (seconds)
  lengthInBits += 8;

  // Simple field (msec)
  lengthInBits += 12;

  // Simple field (dow)
  lengthInBits += 4;

  return lengthInBits;
}

