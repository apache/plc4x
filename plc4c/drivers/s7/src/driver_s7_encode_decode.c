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

#include <cotp_protocol_class.h>
#include <ctype.h>
#include <plc4c/driver_s7.h>
#include <stdlib.h>
#include <string.h>
#include <tpkt_packet.h>
#include "plc4c/driver_s7_encode_decode.h"

uint16_t plc4c_driver_s7_encode_tsap_id(
    plc4c_driver_s7_device_group device_group, uint8_t rack, uint8_t slot) {
  return (device_group << 8) |
         ((uint16_t) ((uint16_t) rack & (uint16_t) 0x000F) << (uint16_t) 4) |
         ((uint16_t) slot & (uint16_t) 0x000F);
}

uint16_t plc4c_driver_s7_get_nearest_matching_tpdu_size(uint16_t pdu_size) {
  for (int i = 0; i < plc4c_s7_read_write_cotp_tpdu_size_num_values(); i++) {
    uint16_t cur_value =
        plc4c_s7_read_write_cotp_tpdu_size_get_size_in_bytes(
            plc4c_s7_read_write_cotp_tpdu_size_value_for_index(i));
    if (cur_value >= pdu_size) {
      return plc4c_s7_read_write_cotp_tpdu_size_value_for_index(i);
    }
  }
  return 0;
}

plc4c_driver_s7_controller_type decode_controller_type(char* article_number) {
  char* prefix = "6ES7 ";
  // If this article-number doesn't start with this prefix, we can't decode it.
  if (strncmp(prefix, article_number, (size_t) strlen(prefix)) != 0) {
    return PLC4C_DRIVER_S7_CONTROLLER_TYPE_ANY;
  }
  char model = *(article_number + 5);
  switch (model) {
    case '2':
      return PLC4C_DRIVER_S7_CONTROLLER_TYPE_S7_1200;
    case '5':
      return PLC4C_DRIVER_S7_CONTROLLER_TYPE_S7_1500;
    case '3':
      return PLC4C_DRIVER_S7_CONTROLLER_TYPE_S7_300;
    case '4':
      return PLC4C_DRIVER_S7_CONTROLLER_TYPE_S7_400;
    default:
      return PLC4C_DRIVER_S7_CONTROLLER_TYPE_ANY;
  }
}

int8_t decode_hex_char(char c) {
  if (('0' <= c) && (c <= '9')) {
    return (int8_t) (c - 48);
  }
  if (('A' <= c) && (c <= 'D')) {
    return (int8_t) (c - 55);
  }
  return -1;
}

plc4c_return_code decode_byte(const char* from_ptr, const char* to_ptr, uint8_t* value) {
  if (to_ptr - from_ptr != 2) {
    return INTERNAL_ERROR;
  }

  int8_t first_char = decode_hex_char(*from_ptr);
  int8_t second_char = decode_hex_char(*(from_ptr + 1));
  if ((first_char == -1) || (second_char == -1)) {
    return INTERNAL_ERROR;
  }
  *value = ((uint8_t) ((uint8_t) first_char << (uint8_t) 4)) | (uint8_t) second_char;
  return OK;
}


plc4c_return_code plc4c_driver_s7_encode_address(char* address, void** item) {

  plc4c_s7_read_write_s7_var_request_parameter_item* s7_item;
  
  s7_item = malloc(sizeof(plc4c_s7_read_write_s7_var_request_parameter_item));
  s7_item->_type =
      plc4c_s7_read_write_s7_var_request_parameter_item_type_plc4c_s7_read_write_s7_var_request_parameter_item_address;

  // Java Regexp:
  // ADDRESS_PATTERN =
  // ^%(?<memoryArea>.)(?<transferSizeCode>[XBWD]?)(?<byteOffset>\d{1,7})(.(?<bitOffset>[0-7]))?:(?<dataType>[a-zA-Z_]+)(\[(?<numElements>\d+)])?
  // DATA_BLOCK_ADDRESS_PATTERN =
  // ^%DB(?<blockNumber>\d{1,5}).DB(?<transferSizeCode>[XBWD]?)(?<byteOffset>\d{1,7})(.(?<bitOffset>[0-7]))?:(?<dataType>[a-zA-Z_]+)(\[(?<numElements>\d+)])?
  // DATA_BLOCK_SHORT_PATTERN =
  // ^%DB(?<blockNumber>\d{1,5}).(?<byteOffset>\d{1,7})(.(?<bitOffset>[0-7]))?:(?<dataType>[a-zA-Z_]+)(\[(?<numElements>\d+)])?
  // DATA_BLOCK_STRING_ADDRESS_PATTERN =
  // ^%DB(?<blockNumber>\d{1,5}).DB(?<transferSizeCode>[XBWD]?)(?<byteOffset>\d{1,7}):STRING\((?<stringLength>\d{1,3})\)(\[(?<numElements>\d+)])?
  // DATA_BLOCK_STRING_SHORT_PATTERN =
  // ^%DB(?<blockNumber>\d{1,5}):(?<byteOffset>\d{1,7}):STRING\((?<stringLength>\d{1,3})\)(\[(?<numElements>\d+)])?
  // PLC_PROXY_ADDRESS_PATTERN =
  // [0-9A-F]{2}-[0-9A-F]{2}-[0-9A-F]{2}-[0-9A-F]{2}-[0-9A-F]{2}-[0-9A-F]{2}-[0-9A-F]{2}-[0-9A-F]{2}-[0-9A-F]{2}-[0-9A-F]{2}

  // Parser logic
  char* cur_pos = address;
  char* last_pos = address;
  char* string_encoding = NULL;
  // - Does it start with "%"?
  if (*cur_pos == '%') {
    cur_pos++;

    char* memory_area = NULL;
    char* block_number = NULL;
    char* transfer_size_code = NULL;
    char* byte_offset = NULL;
    char* bit_offset = NULL;
    char* data_type = NULL;
    char* string_length = NULL;
    char* num_elements = NULL;

    ////////////////////////////////////////////////////////////////////////////
    // First extract the different parts of the address
    ////////////////////////////////////////////////////////////////////////////

    if (!isalpha(*cur_pos)) {
      free(s7_item);
      return INVALID_ADDRESS;
    }

    last_pos = cur_pos;
    while (isalpha(*cur_pos)) {
      cur_pos++;
    }
    uint16_t len = cur_pos - last_pos;
    memory_area = malloc(sizeof(char) * (len + 1));
    strncpy(memory_area, last_pos, len);
    *(memory_area + len) = '\0';

    // If it's a DB-block, get the block_number
    if (strcmp(memory_area, "DB") == 0) {
      last_pos = cur_pos;
      while (isdigit(*cur_pos)) {
        cur_pos++;
      }
      len = cur_pos - last_pos;
      block_number = malloc(sizeof(char) * (len + 1));
      strncpy(block_number, last_pos, len);
      *(block_number + len) = '\0';

      // Skip the "."
      cur_pos++;
    }

    // If the next is not a digit it might be DB, DB{transferSizeCode}
    // or {transferSizeCode}
    if (!isdigit(*cur_pos)) {
      last_pos = cur_pos;
      while (!isdigit(*cur_pos)) {
        cur_pos++;
      }
      len = cur_pos - last_pos;
      // If it's at least 2 digits long, it's DB{transferSizeCode} or
      // "DB". So we get rid of the "DB" prefix, as this has no value for us.
      if (len >= 2) {
        last_pos += 2;
      }
      // If it's 1 or 3 long it contains a "transferSizeCode", which is just
      // one char long.
      if ((len == 1) || (len == 3)) {
        transfer_size_code = malloc(sizeof(char) * 2);
        *transfer_size_code = *last_pos;
        *(transfer_size_code + 1) = '\0';
      }
    }

    // Next comes the byte_offset
    last_pos = cur_pos;
    while (isdigit(*cur_pos)) {
      cur_pos++;
    }
    len = cur_pos - last_pos;
    byte_offset = malloc(sizeof(char) * (len + 1));
    strncpy(byte_offset, last_pos, len);
    *(byte_offset + len) = '\0';

    // Parse the bit_offset
    if (*cur_pos == '.') {
      cur_pos++;
      // Next comes the byte_offset
      last_pos = cur_pos;
      while (isdigit(*cur_pos)) {
        cur_pos++;
      }
      len = cur_pos - last_pos;
      bit_offset = malloc(sizeof(char) * (len + 1));
      strncpy(bit_offset, last_pos, len);
      *(bit_offset + len) = '\0';
    }

    // Skip the ":" char.
    cur_pos++;

    // Next comes the data_type
    last_pos = cur_pos;
    while (isalpha(*cur_pos)) {
      cur_pos++;
    }
    len = cur_pos - last_pos;
    data_type = malloc(sizeof(char) * (len + 1));
    strncpy(data_type, last_pos, len);
    *(data_type + len) = '\0';

    if ((*cur_pos == '(') && (strcmp(data_type, "STRING") == 0)) {
      // Next comes the string_length
      last_pos = ++cur_pos;
      while (isdigit(*cur_pos)) {
        cur_pos++;
      }
      len = cur_pos - last_pos;
      string_length = malloc(sizeof(char) * (len + 1));
      strncpy(string_length, last_pos, len);
      *(string_length + len) = '\0';

      // Skip the ")"
      cur_pos++;
    }

    if (*cur_pos == '[') {
      // Next comes the num_elements
      cur_pos++;
      last_pos = cur_pos;
      while (isdigit(*cur_pos)) {
        cur_pos++;
      }
      len = cur_pos - last_pos;
      num_elements = malloc(sizeof(char) * (len + 1));
      strncpy(num_elements, last_pos, len);
      *(num_elements + len) = '\0';
    }

    if (*cur_pos == '|') {
      // Next comes the num_elements
      cur_pos++;
      last_pos = cur_pos;
      while (*cur_pos != '\0') {
        cur_pos++;
      }
      len = cur_pos - last_pos;
      string_encoding = malloc(sizeof(char) * (len + 1));
      strncpy(string_encoding, last_pos, len);
      *(string_encoding + len) = '\0';
    }

    ////////////////////////////////////////////////////////////////////////////
    // Now parse the contents.
    ////////////////////////////////////////////////////////////////////////////

    plc4c_s7_read_write_s7_address* any_address = malloc(sizeof(plc4c_s7_read_write_s7_address));
    if(any_address == NULL) {
      free(memory_area);
      free(block_number);
      free(transfer_size_code);
      free(byte_offset);
      free(bit_offset);
      free(data_type);
      free(string_length);
      free(num_elements);
      free(string_encoding);
      free(s7_item);
      return NO_MEMORY;
    }
    any_address->_type = plc4c_s7_read_write_s7_address_type_plc4c_s7_read_write_s7_address_any;

    any_address->s7_address_any_area = plc4c_s7_read_write_memory_area_null();
    for(int i = 0; i < plc4c_s7_read_write_memory_area_num_values(); i++) {
      plc4c_s7_read_write_memory_area ma = plc4c_s7_read_write_memory_area_value_for_index(i);
      if(strcmp(plc4c_s7_read_write_memory_area_get_short_name(ma), memory_area) == 0) {
        any_address->s7_address_any_area = ma;
        break;
      }
    }
    if (any_address->s7_address_any_area == plc4c_s7_read_write_memory_area_null()) {
      free(memory_area);
      free(block_number);
      free(transfer_size_code);
      free(byte_offset);
      free(bit_offset);
      free(data_type);
      free(string_length);
      free(num_elements);
      free(string_encoding);
      free(any_address);
      free(s7_item);
      return INVALID_ADDRESS;
    }
    free(memory_area);

    if (block_number != NULL) {
      any_address->s7_address_any_db_number = strtol(block_number, 0, 10);
    } else {
      any_address->s7_address_any_db_number = 0;
    }
    free(block_number);

    any_address->s7_address_any_byte_address = strtol(byte_offset, 0, 10);
    free(byte_offset);

    if (bit_offset != NULL) {
      any_address->s7_address_any_bit_address = strtol(bit_offset, 0, 10);
    } else {
      any_address->s7_address_any_bit_address = 0;
    }
    free(bit_offset);

    any_address->s7_address_any_transport_size =
        plc4c_s7_read_write_transport_size_value_of(data_type);
    free(data_type);

    if (num_elements != NULL) {
      any_address->s7_address_any_number_of_elements =
          strtol(num_elements, 0, 10);
    } else {
      any_address->s7_address_any_number_of_elements = 1;
    }
    free(num_elements);

    // TODO: THis should be moved to "driver_s7_packets.c->plc4c_return_code plc4c_driver_s7_create_s7_read_request"
    if (any_address->s7_address_any_transport_size ==
        plc4c_s7_read_write_transport_size_TIME ||
        any_address->s7_address_any_transport_size ==
            plc4c_s7_read_write_transport_size_LINT ||
        any_address->s7_address_any_transport_size ==
            plc4c_s7_read_write_transport_size_ULINT ||
        any_address->s7_address_any_transport_size ==
            plc4c_s7_read_write_transport_size_LWORD ||
        any_address->s7_address_any_transport_size ==
            plc4c_s7_read_write_transport_size_LREAL ||
        any_address->s7_address_any_transport_size ==
            plc4c_s7_read_write_transport_size_REAL ||
        any_address->s7_address_any_transport_size ==
            plc4c_s7_read_write_transport_size_LTIME ||
        any_address->s7_address_any_transport_size ==
            plc4c_s7_read_write_transport_size_DATE ||
        any_address->s7_address_any_transport_size ==
            plc4c_s7_read_write_transport_size_TIME_OF_DAY ||
        any_address->s7_address_any_transport_size ==
            plc4c_s7_read_write_transport_size_DATE_AND_TIME) {
      any_address->s7_address_any_transport_size = plc4c_s7_read_write_transport_size_BYTE;
        any_address->s7_address_any_number_of_elements =
          plc4c_s7_read_write_transport_size_length_in_bytes(&(any_address->s7_address_any_transport_size)) *
            any_address->s7_address_any_number_of_elements;
    } else if (any_address->s7_address_any_transport_size ==
         plc4c_s7_read_write_transport_size_STRING) {
      any_address->s7_address_any_transport_size = plc4c_s7_read_write_transport_size_BYTE;
      if (string_length != NULL) {
        any_address->s7_address_any_number_of_elements =
            (strtol(string_length, 0, 10) +2) *
                any_address->s7_address_any_number_of_elements;
      } else {
        any_address->s7_address_any_number_of_elements =
            256 * any_address->s7_address_any_number_of_elements;
      }
    } else if (any_address->s7_address_any_transport_size ==
        plc4c_s7_read_write_transport_size_WSTRING) {
      any_address->s7_address_any_transport_size = plc4c_s7_read_write_transport_size_BYTE;
      if (string_length != NULL) {
        any_address->s7_address_any_number_of_elements =
            (strtol(string_length, 0, 10) +2) * 2 *
            any_address->s7_address_any_number_of_elements;
      } else {
        any_address->s7_address_any_number_of_elements =
            512 * any_address->s7_address_any_number_of_elements;
      }
    } else if (any_address->s7_address_any_transport_size ==
               plc4c_s7_read_write_transport_size_TOD) {
      any_address->s7_address_any_transport_size = plc4c_s7_read_write_transport_size_TIME_OF_DAY;
    }
    free(string_length);

    // Check the optional transport size code.
    if(transfer_size_code != NULL) {
      if(plc4c_s7_read_write_transport_size_get_short_name(any_address->s7_address_any_transport_size) != *transfer_size_code) {
        free(transfer_size_code);
        free(any_address);
        free(s7_item);
        return INVALID_ADDRESS;
      }
    }
    free(transfer_size_code);

    s7_item->s7_var_request_parameter_item_address_address = any_address;
  }
  // - Else -> PLC_PROXY_ADDRESS_PATTERN
  else {
    //   - parse the sequence of 2 digit Hex numbers into an array of 10 bytes
    uint8_t* raw_data = malloc(sizeof(uint8_t) * 10);
    if (raw_data == NULL) {
      free(s7_item);
      return NO_MEMORY;
    }
    cur_pos += 2;
    for (int i = 0; i < 10; i++) {
      plc4c_return_code return_code =
          decode_byte(last_pos, cur_pos, raw_data + i);
      if (return_code != OK) {
        free(raw_data);
        free(s7_item);
        return return_code;
      }
      if (i < 9) {
        if (*cur_pos != '-') {
          free(raw_data);
          free(s7_item);
          return INVALID_ADDRESS;
        }
        // Move to the next segment.
        last_pos += 3;
        cur_pos += 3;
      }
    }
    //   - create a plc4c_spi_read_buffer from the 10 byte array
    plc4c_spi_read_buffer* read_buffer;
    plc4c_return_code return_code =
        plc4c_spi_read_buffer_create(raw_data, 10, &read_buffer);
    if (return_code != OK) {
      free(read_buffer);
      free(raw_data);
      free(s7_item);
      return return_code;
    }
    //   - plc4c_s7_read_write_s7_var_request_parameter_item_parse function to
    //   parse the byte array
    //   - directly add the resulting struct to the request
    plc4c_s7_read_write_s7_address_parse(
        plc4x_spi_context_background(), read_buffer, &s7_item->s7_var_request_parameter_item_address_address);

    free(read_buffer);
    free(raw_data);
  }
  plc4c_s7_read_write_s7_var_request_parameter_item_field* s7_item_field;

  s7_item_field = malloc(sizeof(plc4c_s7_read_write_s7_var_request_parameter_item_field));
  s7_item_field->parameter_item = s7_item;
  s7_item_field->s7_address_any_encoding_of_string = string_encoding;
  *item = s7_item_field;

  return OK;
}
