/*
  Licensed to the Apache Software Foundation (ASF) under one
  or more contributor license agreements.  See the NOTICE file
  distributed with this work for additional information
  regarding copyright ownership.  The ASF licenses this file
  to you under the Apache License, Version 2.0 (the
  "License"); you may not use this file except in compliance
  with the License.  You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing,
  software distributed under the License is distributed on an
  "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
  KIND, either express or implied.  See the License for the
  specific language governing permissions and limitations
  under the License.
*/

#include <stdio.h>
#include <time.h>
#include <plc4c/data.h>
#include <plc4c/spi/evaluation_helper.h>
#include <plc4c/driver_s7.h>
#include "data_item.h"

// Parse function.
plc4c_return_code plc4c_s7_read_write_data_item_parse(plc4c_spi_read_buffer* io, uint8_t dataProtocolId, int32_t stringLength, plc4c_data** data_item) {
    uint16_t startPos = plc4c_spi_read_get_pos(io);
    uint16_t curPos;
    plc4c_return_code _res = OK;

        if(dataProtocolId == 01) { /* BOOL */

                // Reserved Field (Compartmentalized so the "reserved" variable can't leak)
                {
                    unsigned int _reserved = 0;
                    _res = plc4c_spi_read_unsigned_byte(io, 7, (uint8_t*) &_reserved);
                    if(_res != OK) {
                        return _res;
                    }
                    if(_reserved != 0x00) {
                      printf("Expected constant value '%d' but got '%d' for reserved field.", 0x00, _reserved);
                    }
                }

                // Simple Field (value)
                bool value = false;
                _res = plc4c_spi_read_bit(io, (bool*) &value);
                if(_res != OK) {
                    return _res;
                }

                *data_item = plc4c_data_create_bool_data(value);

        } else 
        if(dataProtocolId == 11) { /* BOOL */

                    // Array field (value)
        } else 
        if(dataProtocolId == 12) { /* BOOL */

                    // Array field (value)
        } else 
        if(dataProtocolId == 13) { /* BOOL */

                    // Array field (value)
        } else 
        if(dataProtocolId == 14) { /* BOOL */

                    // Array field (value)
        } else 
        if(dataProtocolId == 21) { /* SINT */

                // Simple Field (value)
                int8_t value = 0;
                _res = plc4c_spi_read_signed_byte(io, 8, (int8_t*) &value);
                if(_res != OK) {
                    return _res;
                }

                *data_item = plc4c_data_create_int8_t_data(value);

        } else 
        if(dataProtocolId == 22) { /* USINT */

                // Simple Field (value)
                uint8_t value = 0;
                _res = plc4c_spi_read_unsigned_byte(io, 8, (uint8_t*) &value);
                if(_res != OK) {
                    return _res;
                }

                *data_item = plc4c_data_create_uint8_t_data(value);

        } else 
        if(dataProtocolId == 23) { /* INT */

                // Simple Field (value)
                int16_t value = 0;
                _res = plc4c_spi_read_signed_short(io, 16, (int16_t*) &value);
                if(_res != OK) {
                    return _res;
                }

                *data_item = plc4c_data_create_int16_t_data(value);

        } else 
        if(dataProtocolId == 24) { /* UINT */

                // Simple Field (value)
                uint16_t value = 0;
                _res = plc4c_spi_read_unsigned_short(io, 16, (uint16_t*) &value);
                if(_res != OK) {
                    return _res;
                }

                *data_item = plc4c_data_create_uint16_t_data(value);

        } else 
        if(dataProtocolId == 25) { /* DINT */

                // Simple Field (value)
                int32_t value = 0;
                _res = plc4c_spi_read_signed_int(io, 32, (int32_t*) &value);
                if(_res != OK) {
                    return _res;
                }

                *data_item = plc4c_data_create_int32_t_data(value);

        } else 
        if(dataProtocolId == 26) { /* UDINT */

                // Simple Field (value)
                uint32_t value = 0;
                _res = plc4c_spi_read_unsigned_int(io, 32, (uint32_t*) &value);
                if(_res != OK) {
                    return _res;
                }

                *data_item = plc4c_data_create_uint32_t_data(value);

        } else 
        if(dataProtocolId == 27) { /* LINT */

                // Simple Field (value)
                int64_t value = 0;
                _res = plc4c_spi_read_signed_long(io, 64, (int64_t*) &value);
                if(_res != OK) {
                    return _res;
                }

                *data_item = plc4c_data_create_int64_t_data(value);

        } else 
        if(dataProtocolId == 28) { /* ULINT */

                // Simple Field (value)
                uint64_t value = 0;
                _res = plc4c_spi_read_unsigned_long(io, 64, (uint64_t*) &value);
                if(_res != OK) {
                    return _res;
                }

                *data_item = plc4c_data_create_uint64_t_data(value);

        } else 
        if(dataProtocolId == 31) { /* REAL */

                // Simple Field (value)
                float value = 0.0;
                _res = plc4c_spi_read_float(io, 32, (float*) &value);
                if(_res != OK) {
                    return _res;
                }

                *data_item = plc4c_data_create_float_data(value);

        } else 
        if(dataProtocolId == 32) { /* LREAL */

                // Simple Field (value)
                double value = 0.0;
                _res = plc4c_spi_read_double(io, 64, (double*) &value);
                if(_res != OK) {
                    return _res;
                }

                *data_item = plc4c_data_create_double_data(value);

        } else 
        if(dataProtocolId == 41) { /* String */
        } else 
        if(dataProtocolId == 42) { /* String */
        } else 
        if(dataProtocolId == 43) { /* String */

                    // Manual Field (value)
                    char* value = (char*) (plc4c_s7_read_write_parse_s7_string(io, stringLength, "UTF-8"));
        } else 
        if(dataProtocolId == 44) { /* String */

                    // Manual Field (value)
                    char* value = (char*) (plc4c_s7_read_write_parse_s7_string(io, stringLength, "UTF-16"));
        } else 
        if(dataProtocolId == 51) { /* Time */

                    // Manual Field (value)
                    time_t value = (time_t) (plc4c_s7_read_write_parse_tia_time(io));
        } else 
        if(dataProtocolId == 52) { /* Time */

                    // Manual Field (value)
                    time_t value = (time_t) (plc4c_s7_read_write_parse_s5_time(io));
        } else 
        if(dataProtocolId == 53) { /* Time */

                    // Manual Field (value)
                    time_t value = (time_t) (plc4c_s7_read_write_parse_tia_l_time(io));
        } else 
        if(dataProtocolId == 54) { /* Date */

                    // Manual Field (value)
                    time_t value = (time_t) (plc4c_s7_read_write_parse_tia_date(io));
        } else 
        if(dataProtocolId == 55) { /* Time */

                    // Manual Field (value)
                    time_t value = (time_t) (plc4c_s7_read_write_parse_tia_time_of_day(io));
        } else 
        if(dataProtocolId == 56) { /* DateTime */

                    // Manual Field (value)
                    time_t value = (time_t) (plc4c_s7_read_write_parse_tia_date_time(io));
        }

  return OK;
}

plc4c_return_code plc4c_s7_read_write_data_item_serialize(plc4c_spi_write_buffer* io, plc4c_data** data_item) {
  plc4c_return_code _res = OK;

  return OK;
}

uint16_t plc4c_s7_read_write_data_item_length_in_bytes(plc4c_data* data_item) {
  return plc4c_s7_read_write_data_item_length_in_bits(data_item) / 8;
}

uint16_t plc4c_s7_read_write_data_item_length_in_bits(plc4c_data* data_item) {
  uint16_t lengthInBits = 0;

  return lengthInBits;
}

