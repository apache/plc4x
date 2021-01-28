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

#include "transport_size.h"
#include <string.h>

#include "data_transport_size.h"
#include "transport_size.h"

// Create an empty NULL-struct
static const plc4c_s7_read_write_transport_size plc4c_s7_read_write_transport_size_null_const;

plc4c_s7_read_write_transport_size plc4c_s7_read_write_transport_size_null() {
  return plc4c_s7_read_write_transport_size_null_const;
}

plc4c_s7_read_write_transport_size plc4c_s7_read_write_transport_size_value_of(char* value_string) {
    if(strcmp(value_string, "LWORD") == 0) {
        return plc4c_s7_read_write_transport_size_LWORD;
    }
    if(strcmp(value_string, "BOOL") == 0) {
        return plc4c_s7_read_write_transport_size_BOOL;
    }
    if(strcmp(value_string, "BYTE") == 0) {
        return plc4c_s7_read_write_transport_size_BYTE;
    }
    if(strcmp(value_string, "CHAR") == 0) {
        return plc4c_s7_read_write_transport_size_CHAR;
    }
    if(strcmp(value_string, "WORD") == 0) {
        return plc4c_s7_read_write_transport_size_WORD;
    }
    if(strcmp(value_string, "INT") == 0) {
        return plc4c_s7_read_write_transport_size_INT;
    }
    if(strcmp(value_string, "DWORD") == 0) {
        return plc4c_s7_read_write_transport_size_DWORD;
    }
    if(strcmp(value_string, "DINT") == 0) {
        return plc4c_s7_read_write_transport_size_DINT;
    }
    if(strcmp(value_string, "REAL") == 0) {
        return plc4c_s7_read_write_transport_size_REAL;
    }
    if(strcmp(value_string, "DATE") == 0) {
        return plc4c_s7_read_write_transport_size_DATE;
    }
    if(strcmp(value_string, "TIME") == 0) {
        return plc4c_s7_read_write_transport_size_TIME;
    }
    if(strcmp(value_string, "DATE_AND_TIME") == 0) {
        return plc4c_s7_read_write_transport_size_DATE_AND_TIME;
    }
    if(strcmp(value_string, "WCHAR") == 0) {
        return plc4c_s7_read_write_transport_size_WCHAR;
    }
    if(strcmp(value_string, "LREAL") == 0) {
        return plc4c_s7_read_write_transport_size_LREAL;
    }
    return -1;
}

int plc4c_s7_read_write_transport_size_num_values() {
  return 26;
}

plc4c_s7_read_write_transport_size plc4c_s7_read_write_transport_size_value_for_index(int index) {
    switch(index) {
      case 0: {
        return plc4c_s7_read_write_transport_size_LWORD;
      }
      case 1: {
        return plc4c_s7_read_write_transport_size_BOOL;
      }
      case 2: {
        return plc4c_s7_read_write_transport_size_BYTE;
      }
      case 3: {
        return plc4c_s7_read_write_transport_size_CHAR;
      }
      case 4: {
        return plc4c_s7_read_write_transport_size_WORD;
      }
      case 5: {
        return plc4c_s7_read_write_transport_size_INT;
      }
      case 6: {
        return plc4c_s7_read_write_transport_size_DWORD;
      }
      case 7: {
        return plc4c_s7_read_write_transport_size_DINT;
      }
      case 8: {
        return plc4c_s7_read_write_transport_size_REAL;
      }
      case 9: {
        return plc4c_s7_read_write_transport_size_DATE;
      }
      case 10: {
        return plc4c_s7_read_write_transport_size_TIME;
      }
      case 11: {
        return plc4c_s7_read_write_transport_size_DATE_AND_TIME;
      }
      case 12: {
        return plc4c_s7_read_write_transport_size_WCHAR;
      }
      case 13: {
        return plc4c_s7_read_write_transport_size_LREAL;
      }
      default: {
        return -1;
      }
    }
}

bool plc4c_s7_read_write_transport_size_get_supported__s7_300(plc4c_s7_read_write_transport_size value) {
  switch(value) {
    case plc4c_s7_read_write_transport_size_LWORD: { /* '0x00' */
      return false;
    }
    case plc4c_s7_read_write_transport_size_BOOL: { /* '0x01' */
      return true;
    }
    case plc4c_s7_read_write_transport_size_BYTE: { /* '0x02' */
      return true;
    }
    case plc4c_s7_read_write_transport_size_CHAR: { /* '0x03' */
      return true;
    }
    case plc4c_s7_read_write_transport_size_WORD: { /* '0x04' */
      return true;
    }
    case plc4c_s7_read_write_transport_size_INT: { /* '0x05' */
      return true;
    }
    case plc4c_s7_read_write_transport_size_DWORD: { /* '0x06' */
      return true;
    }
    case plc4c_s7_read_write_transport_size_DINT: { /* '0x07' */
      return true;
    }
    case plc4c_s7_read_write_transport_size_REAL: { /* '0x08' */
      return true;
    }
    case plc4c_s7_read_write_transport_size_DATE: { /* '0x09' */
      return true;
    }
    case plc4c_s7_read_write_transport_size_TIME: { /* '0x0B' */
      return true;
    }
    case plc4c_s7_read_write_transport_size_DATE_AND_TIME: { /* '0x0F' */
      return true;
    }
    case plc4c_s7_read_write_transport_size_WCHAR: { /* '0x13' */
      return false;
    }
    case plc4c_s7_read_write_transport_size_LREAL: { /* '0x30' */
      return false;
    }
    default: {
      return 0;
    }
  }
}

bool plc4c_s7_read_write_transport_size_get_supported__logo(plc4c_s7_read_write_transport_size value) {
  switch(value) {
    case plc4c_s7_read_write_transport_size_LWORD: { /* '0x00' */
      return false;
    }
    case plc4c_s7_read_write_transport_size_BOOL: { /* '0x01' */
      return true;
    }
    case plc4c_s7_read_write_transport_size_BYTE: { /* '0x02' */
      return true;
    }
    case plc4c_s7_read_write_transport_size_CHAR: { /* '0x03' */
      return true;
    }
    case plc4c_s7_read_write_transport_size_WORD: { /* '0x04' */
      return true;
    }
    case plc4c_s7_read_write_transport_size_INT: { /* '0x05' */
      return true;
    }
    case plc4c_s7_read_write_transport_size_DWORD: { /* '0x06' */
      return true;
    }
    case plc4c_s7_read_write_transport_size_DINT: { /* '0x07' */
      return true;
    }
    case plc4c_s7_read_write_transport_size_REAL: { /* '0x08' */
      return true;
    }
    case plc4c_s7_read_write_transport_size_DATE: { /* '0x09' */
      return true;
    }
    case plc4c_s7_read_write_transport_size_TIME: { /* '0x0B' */
      return true;
    }
    case plc4c_s7_read_write_transport_size_DATE_AND_TIME: { /* '0x0F' */
      return false;
    }
    case plc4c_s7_read_write_transport_size_WCHAR: { /* '0x13' */
      return true;
    }
    case plc4c_s7_read_write_transport_size_LREAL: { /* '0x30' */
      return false;
    }
    default: {
      return 0;
    }
  }
}

uint8_t plc4c_s7_read_write_transport_size_get_size_in_bytes(plc4c_s7_read_write_transport_size value) {
  switch(value) {
    case plc4c_s7_read_write_transport_size_LWORD: { /* '0x00' */
      return 8;
    }
    case plc4c_s7_read_write_transport_size_BOOL: { /* '0x01' */
      return 1;
    }
    case plc4c_s7_read_write_transport_size_BYTE: { /* '0x02' */
      return 1;
    }
    case plc4c_s7_read_write_transport_size_CHAR: { /* '0x03' */
      return 1;
    }
    case plc4c_s7_read_write_transport_size_WORD: { /* '0x04' */
      return 2;
    }
    case plc4c_s7_read_write_transport_size_INT: { /* '0x05' */
      return 2;
    }
    case plc4c_s7_read_write_transport_size_DWORD: { /* '0x06' */
      return 4;
    }
    case plc4c_s7_read_write_transport_size_DINT: { /* '0x07' */
      return 4;
    }
    case plc4c_s7_read_write_transport_size_REAL: { /* '0x08' */
      return 4;
    }
    case plc4c_s7_read_write_transport_size_DATE: { /* '0x09' */
      return 2;
    }
    case plc4c_s7_read_write_transport_size_TIME: { /* '0x0B' */
      return 4;
    }
    case plc4c_s7_read_write_transport_size_DATE_AND_TIME: { /* '0x0F' */
      return 12;
    }
    case plc4c_s7_read_write_transport_size_WCHAR: { /* '0x13' */
      return 2;
    }
    case plc4c_s7_read_write_transport_size_LREAL: { /* '0x30' */
      return 8;
    }
    default: {
      return 0;
    }
  }
}

bool plc4c_s7_read_write_transport_size_get_supported__s7_400(plc4c_s7_read_write_transport_size value) {
  switch(value) {
    case plc4c_s7_read_write_transport_size_LWORD: { /* '0x00' */
      return false;
    }
    case plc4c_s7_read_write_transport_size_BOOL: { /* '0x01' */
      return true;
    }
    case plc4c_s7_read_write_transport_size_BYTE: { /* '0x02' */
      return true;
    }
    case plc4c_s7_read_write_transport_size_CHAR: { /* '0x03' */
      return true;
    }
    case plc4c_s7_read_write_transport_size_WORD: { /* '0x04' */
      return true;
    }
    case plc4c_s7_read_write_transport_size_INT: { /* '0x05' */
      return true;
    }
    case plc4c_s7_read_write_transport_size_DWORD: { /* '0x06' */
      return true;
    }
    case plc4c_s7_read_write_transport_size_DINT: { /* '0x07' */
      return true;
    }
    case plc4c_s7_read_write_transport_size_REAL: { /* '0x08' */
      return true;
    }
    case plc4c_s7_read_write_transport_size_DATE: { /* '0x09' */
      return true;
    }
    case plc4c_s7_read_write_transport_size_TIME: { /* '0x0B' */
      return true;
    }
    case plc4c_s7_read_write_transport_size_DATE_AND_TIME: { /* '0x0F' */
      return true;
    }
    case plc4c_s7_read_write_transport_size_WCHAR: { /* '0x13' */
      return false;
    }
    case plc4c_s7_read_write_transport_size_LREAL: { /* '0x30' */
      return false;
    }
    default: {
      return 0;
    }
  }
}

bool plc4c_s7_read_write_transport_size_get_supported__s7_1200(plc4c_s7_read_write_transport_size value) {
  switch(value) {
    case plc4c_s7_read_write_transport_size_LWORD: { /* '0x00' */
      return false;
    }
    case plc4c_s7_read_write_transport_size_BOOL: { /* '0x01' */
      return true;
    }
    case plc4c_s7_read_write_transport_size_BYTE: { /* '0x02' */
      return true;
    }
    case plc4c_s7_read_write_transport_size_CHAR: { /* '0x03' */
      return true;
    }
    case plc4c_s7_read_write_transport_size_WORD: { /* '0x04' */
      return true;
    }
    case plc4c_s7_read_write_transport_size_INT: { /* '0x05' */
      return true;
    }
    case plc4c_s7_read_write_transport_size_DWORD: { /* '0x06' */
      return true;
    }
    case plc4c_s7_read_write_transport_size_DINT: { /* '0x07' */
      return true;
    }
    case plc4c_s7_read_write_transport_size_REAL: { /* '0x08' */
      return true;
    }
    case plc4c_s7_read_write_transport_size_DATE: { /* '0x09' */
      return true;
    }
    case plc4c_s7_read_write_transport_size_TIME: { /* '0x0B' */
      return true;
    }
    case plc4c_s7_read_write_transport_size_DATE_AND_TIME: { /* '0x0F' */
      return false;
    }
    case plc4c_s7_read_write_transport_size_WCHAR: { /* '0x13' */
      return true;
    }
    case plc4c_s7_read_write_transport_size_LREAL: { /* '0x30' */
      return true;
    }
    default: {
      return 0;
    }
  }
}

uint8_t plc4c_s7_read_write_transport_size_get_size_code(plc4c_s7_read_write_transport_size value) {
  switch(value) {
    case plc4c_s7_read_write_transport_size_LWORD: { /* '0x00' */
      return 'X';
    }
    case plc4c_s7_read_write_transport_size_BOOL: { /* '0x01' */
      return 'X';
    }
    case plc4c_s7_read_write_transport_size_BYTE: { /* '0x02' */
      return 'B';
    }
    case plc4c_s7_read_write_transport_size_CHAR: { /* '0x03' */
      return 'B';
    }
    case plc4c_s7_read_write_transport_size_WORD: { /* '0x04' */
      return 'W';
    }
    case plc4c_s7_read_write_transport_size_INT: { /* '0x05' */
      return 'W';
    }
    case plc4c_s7_read_write_transport_size_DWORD: { /* '0x06' */
      return 'D';
    }
    case plc4c_s7_read_write_transport_size_DINT: { /* '0x07' */
      return 'D';
    }
    case plc4c_s7_read_write_transport_size_REAL: { /* '0x08' */
      return 'D';
    }
    case plc4c_s7_read_write_transport_size_DATE: { /* '0x09' */
      return 'X';
    }
    case plc4c_s7_read_write_transport_size_TIME: { /* '0x0B' */
      return 'X';
    }
    case plc4c_s7_read_write_transport_size_DATE_AND_TIME: { /* '0x0F' */
      return 'X';
    }
    case plc4c_s7_read_write_transport_size_WCHAR: { /* '0x13' */
      return 'X';
    }
    case plc4c_s7_read_write_transport_size_LREAL: { /* '0x30' */
      return 'X';
    }
    default: {
      return 0;
    }
  }
}

bool plc4c_s7_read_write_transport_size_get_supported__s7_1500(plc4c_s7_read_write_transport_size value) {
  switch(value) {
    case plc4c_s7_read_write_transport_size_LWORD: { /* '0x00' */
      return true;
    }
    case plc4c_s7_read_write_transport_size_BOOL: { /* '0x01' */
      return true;
    }
    case plc4c_s7_read_write_transport_size_BYTE: { /* '0x02' */
      return true;
    }
    case plc4c_s7_read_write_transport_size_CHAR: { /* '0x03' */
      return true;
    }
    case plc4c_s7_read_write_transport_size_WORD: { /* '0x04' */
      return true;
    }
    case plc4c_s7_read_write_transport_size_INT: { /* '0x05' */
      return true;
    }
    case plc4c_s7_read_write_transport_size_DWORD: { /* '0x06' */
      return true;
    }
    case plc4c_s7_read_write_transport_size_DINT: { /* '0x07' */
      return true;
    }
    case plc4c_s7_read_write_transport_size_REAL: { /* '0x08' */
      return true;
    }
    case plc4c_s7_read_write_transport_size_DATE: { /* '0x09' */
      return true;
    }
    case plc4c_s7_read_write_transport_size_TIME: { /* '0x0B' */
      return true;
    }
    case plc4c_s7_read_write_transport_size_DATE_AND_TIME: { /* '0x0F' */
      return true;
    }
    case plc4c_s7_read_write_transport_size_WCHAR: { /* '0x13' */
      return true;
    }
    case plc4c_s7_read_write_transport_size_LREAL: { /* '0x30' */
      return true;
    }
    default: {
      return 0;
    }
  }
}

plc4c_s7_read_write_data_transport_size plc4c_s7_read_write_transport_size_get_data_transport_size(plc4c_s7_read_write_transport_size value) {
  switch(value) {
    case plc4c_s7_read_write_transport_size_LWORD: { /* '0x00' */
      return -1;
    }
    case plc4c_s7_read_write_transport_size_BOOL: { /* '0x01' */
      return plc4c_s7_read_write_data_transport_size_BIT;
    }
    case plc4c_s7_read_write_transport_size_BYTE: { /* '0x02' */
      return plc4c_s7_read_write_data_transport_size_BYTE_WORD_DWORD;
    }
    case plc4c_s7_read_write_transport_size_CHAR: { /* '0x03' */
      return plc4c_s7_read_write_data_transport_size_BYTE_WORD_DWORD;
    }
    case plc4c_s7_read_write_transport_size_WORD: { /* '0x04' */
      return plc4c_s7_read_write_data_transport_size_BYTE_WORD_DWORD;
    }
    case plc4c_s7_read_write_transport_size_INT: { /* '0x05' */
      return plc4c_s7_read_write_data_transport_size_INTEGER;
    }
    case plc4c_s7_read_write_transport_size_DWORD: { /* '0x06' */
      return plc4c_s7_read_write_data_transport_size_BYTE_WORD_DWORD;
    }
    case plc4c_s7_read_write_transport_size_DINT: { /* '0x07' */
      return plc4c_s7_read_write_data_transport_size_INTEGER;
    }
    case plc4c_s7_read_write_transport_size_REAL: { /* '0x08' */
      return plc4c_s7_read_write_data_transport_size_BYTE_WORD_DWORD;
    }
    case plc4c_s7_read_write_transport_size_DATE: { /* '0x09' */
      return plc4c_s7_read_write_data_transport_size_BYTE_WORD_DWORD;
    }
    case plc4c_s7_read_write_transport_size_TIME: { /* '0x0B' */
      return -1;
    }
    case plc4c_s7_read_write_transport_size_DATE_AND_TIME: { /* '0x0F' */
      return -1;
    }
    case plc4c_s7_read_write_transport_size_WCHAR: { /* '0x13' */
      return -1;
    }
    case plc4c_s7_read_write_transport_size_LREAL: { /* '0x30' */
      return -1;
    }
    default: {
      return 0;
    }
  }
}

plc4c_s7_read_write_transport_size plc4c_s7_read_write_transport_size_get_base_type(plc4c_s7_read_write_transport_size value) {
  switch(value) {
    case plc4c_s7_read_write_transport_size_LWORD: { /* '0x00' */
      return -1;
    }
    case plc4c_s7_read_write_transport_size_BOOL: { /* '0x01' */
      return -1;
    }
    case plc4c_s7_read_write_transport_size_BYTE: { /* '0x02' */
      return -1;
    }
    case plc4c_s7_read_write_transport_size_CHAR: { /* '0x03' */
      return -1;
    }
    case plc4c_s7_read_write_transport_size_WORD: { /* '0x04' */
      return -1;
    }
    case plc4c_s7_read_write_transport_size_INT: { /* '0x05' */
      return -1;
    }
    case plc4c_s7_read_write_transport_size_DWORD: { /* '0x06' */
      return plc4c_s7_read_write_transport_size_WORD;
    }
    case plc4c_s7_read_write_transport_size_DINT: { /* '0x07' */
      return plc4c_s7_read_write_transport_size_INT;
    }
    case plc4c_s7_read_write_transport_size_REAL: { /* '0x08' */
      return -1;
    }
    case plc4c_s7_read_write_transport_size_DATE: { /* '0x09' */
      return -1;
    }
    case plc4c_s7_read_write_transport_size_TIME: { /* '0x0B' */
      return -1;
    }
    case plc4c_s7_read_write_transport_size_DATE_AND_TIME: { /* '0x0F' */
      return -1;
    }
    case plc4c_s7_read_write_transport_size_WCHAR: { /* '0x13' */
      return -1;
    }
    case plc4c_s7_read_write_transport_size_LREAL: { /* '0x30' */
      return plc4c_s7_read_write_transport_size_REAL;
    }
    default: {
      return 0;
    }
  }
}

char* plc4c_s7_read_write_transport_size_get_data_protocol_id(plc4c_s7_read_write_transport_size value) {
  switch(value) {
    case plc4c_s7_read_write_transport_size_LWORD: { /* '0x00' */
      return "IEC61131_LWORD";
    }
    case plc4c_s7_read_write_transport_size_BOOL: { /* '0x01' */
      return "IEC61131_BOOL";
    }
    case plc4c_s7_read_write_transport_size_BYTE: { /* '0x02' */
      return "IEC61131_BYTE";
    }
    case plc4c_s7_read_write_transport_size_CHAR: { /* '0x03' */
      return "IEC61131_CHAR";
    }
    case plc4c_s7_read_write_transport_size_WORD: { /* '0x04' */
      return "IEC61131_WORD";
    }
    case plc4c_s7_read_write_transport_size_INT: { /* '0x05' */
      return "IEC61131_INT";
    }
    case plc4c_s7_read_write_transport_size_DWORD: { /* '0x06' */
      return "IEC61131_DWORD";
    }
    case plc4c_s7_read_write_transport_size_DINT: { /* '0x07' */
      return "IEC61131_DINT";
    }
    case plc4c_s7_read_write_transport_size_REAL: { /* '0x08' */
      return "IEC61131_REAL";
    }
    case plc4c_s7_read_write_transport_size_DATE: { /* '0x09' */
      return "IEC61131_DATE";
    }
    case plc4c_s7_read_write_transport_size_TIME: { /* '0x0B' */
      return "IEC61131_TIME";
    }
    case plc4c_s7_read_write_transport_size_DATE_AND_TIME: { /* '0x0F' */
      return "IEC61131_DATE_AND_TIME";
    }
    case plc4c_s7_read_write_transport_size_WCHAR: { /* '0x13' */
      return "IEC61131_WCHAR";
    }
    case plc4c_s7_read_write_transport_size_LREAL: { /* '0x30' */
      return "IEC61131_LREAL";
    }
    default: {
      return 0;
    }
  }
}
