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
#include <plc4c/spi/evaluation_helper.h>
#include "s7_var_payload_data_item.h"


// Parse function.
plc4c_return_code plc4c_s7_read_write_s7_var_payload_data_item_parse(plc4c_spi_read_buffer* io, bool lastItem, plc4c_s7_read_write_s7_var_payload_data_item** _message) {
  uint16_t startPos = plc4c_spi_read_get_pos(io);
  uint16_t curPos;
  plc4c_return_code _res = OK;

  // Allocate enough memory to contain this data structure.
  (*_message) = malloc(sizeof(plc4c_s7_read_write_s7_var_payload_data_item));
  if(*_message == NULL) {
    return NO_MEMORY;
  }

  // Enum field (returnCode)
  plc4c_s7_read_write_data_transport_error_code returnCode = plc4c_s7_read_write_data_transport_error_code_null();
  _res = plc4c_spi_read_unsigned_byte(io, 8, (uint8_t*) &returnCode);
  if(_res != OK) {
    return _res;
  }
  (*_message)->return_code = returnCode;

  // Enum field (transportSize)
  plc4c_s7_read_write_data_transport_size transportSize = plc4c_s7_read_write_data_transport_size_null();
  _res = plc4c_spi_read_unsigned_byte(io, 8, (uint8_t*) &transportSize);
  if(_res != OK) {
    return _res;
  }
  (*_message)->transport_size = transportSize;

  // Implicit Field (dataLength) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
  uint16_t dataLength = 0;
  _res = plc4c_spi_read_unsigned_short(io, 16, (uint16_t*) &dataLength);
  if(_res != OK) {
    return _res;
  }

  // Array field (data)
  plc4c_list* data = NULL;
  plc4c_utils_list_create(&data);
  if(data == NULL) {
    return NO_MEMORY;
  }
  {
    // Count array
    uint16_t itemCount = ((plc4c_s7_read_write_data_transport_size_get_size_in_bits(transportSize)) ? plc4c_spi_evaluation_helper_ceil((dataLength) / (8.0)) : dataLength);
    for(int curItem = 0; curItem < itemCount; curItem++) {
      
      int8_t* _value = malloc(sizeof(int8_t));
      _res = plc4c_spi_read_signed_byte(io, 8, (int8_t*) _value);
      if(_res != OK) {
        return _res;
      }
      plc4c_utils_list_insert_head_value(data, _value);
    }
  }
  (*_message)->data = data;

  // Padding Field (padding)
  {
    int _timesPadding = (int) ((plc4c_spi_read_has_more(io, 8)) && (((lastItem) ? 0 : (plc4c_spi_evaluation_helper_count(data)) % (2))));
    while (_timesPadding-- > 0) {
      // Just read the padding data and ignore it
      uint8_t _paddingValue = 0;
      _res = plc4c_spi_read_unsigned_byte(io, 8, (uint8_t*) &_paddingValue);
      if(_res != OK) {
        return _res;
      }
    }
  }

  return OK;
}

plc4c_return_code plc4c_s7_read_write_s7_var_payload_data_item_serialize(plc4c_spi_write_buffer* io, plc4c_s7_read_write_s7_var_payload_data_item* _message, bool lastItem) {
  plc4c_return_code _res = OK;

  // Enum field (returnCode)
  _res = plc4c_spi_write_unsigned_byte(io, 8, _message->return_code);
  if(_res != OK) {
    return _res;
  }

  // Enum field (transportSize)
  _res = plc4c_spi_write_unsigned_byte(io, 8, _message->transport_size);
  if(_res != OK) {
    return _res;
  }

  // Implicit Field (dataLength) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
  _res = plc4c_spi_write_unsigned_short(io, 16, (plc4c_spi_evaluation_helper_count(_message->data)) * ((((((_message->transport_size) == (plc4c_s7_read_write_data_transport_size_BIT))) ? 1 : (((plc4c_s7_read_write_data_transport_size_get_size_in_bits(_message->transport_size)) ? 8 : 1))))));
  if(_res != OK) {
    return _res;
  }

  // Array field (data)
  {
    uint8_t itemCount = plc4c_utils_list_size(_message->data);
    for(int curItem = 0; curItem < itemCount; curItem++) {

      int8_t* _value = (int8_t*) plc4c_utils_list_get_value(_message->data, curItem);
      plc4c_spi_write_signed_byte(io, 8, *_value);
    }
  }

  // Padding Field (padding)
  {
    int _timesPadding = (int) (((lastItem) ? 0 : (plc4c_spi_evaluation_helper_count(_message->data)) % (2)));
    while (_timesPadding-- > 0) {
      // Just output the default padding data
      _res = plc4c_spi_write_unsigned_byte(io, 8, 0);
      if(_res != OK) {
        return _res;
      }
    }
  }

  return OK;
}

uint16_t plc4c_s7_read_write_s7_var_payload_data_item_length_in_bytes(plc4c_s7_read_write_s7_var_payload_data_item* _message) {
  return plc4c_s7_read_write_s7_var_payload_data_item_length_in_bits(_message) / 8;
}

uint16_t plc4c_s7_read_write_s7_var_payload_data_item_length_in_bits(plc4c_s7_read_write_s7_var_payload_data_item* _message) {
  uint16_t lengthInBits = 0;

  // Enum Field (returnCode)
  lengthInBits += 8;

  // Enum Field (transportSize)
  lengthInBits += 8;

  // Implicit Field (dataLength)
  lengthInBits += 16;

  // Array field
  lengthInBits += 8 * plc4c_utils_list_size(_message->data);

  // Padding Field (padding)
 int _needsPadding = (int) (((false) ? 0 : (plc4c_spi_evaluation_helper_count(_message->data)) % (2)));
 while(_needsPadding-- > 0) {
    lengthInBits += 8;
  }

  return lengthInBits;
}

