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

#include <plc4c/types.h>

char *plc4c_return_code_to_message(plc4c_return_code return_code) {
  switch (return_code) {
    case UNFINISHED: {
      return "UNFINISHED";
    }
    case OK: {
      return "OK";
    }
    case NO_MEMORY: {
      return "OUT_OF_MEMORY";
    }
    case INVALID_CONNECTION_STRING: {
      return "INVALID CONNECTION STRING";
    }
    case NON_MATCHING_LISTS: {
      return "ITEM AND VALUE LISTS HAVE DIFFERENT SIZES";
    }
    case INVALID_LIST_SIZE: {
      return "INVALID LIST SIZE";
    }
    case NOT_REACHABLE: {
      return "DEVICE NOT REACHABLE";
    }
    case PERMISSION_DENIED: {
      return "PERMISSION DENIED";
    }

    case NO_DRIVER_AVAILABLE: {
      return "NO DRIVER FOUND";
    }
    case UNKNOWN_DRIVER: {
      return "UNKNOWN DRIVER";
    }

    case UNSPECIFIED_TRANSPORT: {
      return "TRANSPORT NOT SPECIFIED";
    }
    case NO_TRANSPORT_AVAILABLE: {
      return "NO TRANSPORT FOUND";
    }
    case UNKNOWN_TRANSPORT: {
      return "UNKNOWN TRANSPORT";
    }

    case UNKNOWN_ERROR: {
      return "UNKNOWN ERROR";
    }
    case INTERNAL_ERROR: {
      return "INTERNAL ERROR";
    }

    default: {
      return "UNKNOWN RETURN CODE";
    }
  }
}

char *plc4c_response_code_to_message(plc4c_response_code response_code) {
  switch (response_code) {
    case PLC4C_RESPONSE_CODE_OK: {
      return "OK";
    }
    case PLC4C_RESPONSE_CODE_NOT_FOUND: {
      return "NOT FOUND";
    }
    case PLC4C_RESPONSE_CODE_ACCESS_DENIED: {
      return "ACCESS DENIED";
    }
    case PLC4C_RESPONSE_CODE_INVALID_ADDRESS: {
      return "INVALID ADDRESS";
    }
    case PLC4C_RESPONSE_CODE_INVALID_DATATYPE: {
      return "INVALID_DATATYPE";
    }
    case PLC4C_RESPONSE_CODE_INTERNAL_ERROR: {
      return "INTERNAL ERROR";
    }
    case PLC4C_RESPONSE_CODE_RESPONSE_PENDING: {
      return "RESPONSE PENDING";
    }
    default: {
      return "UNKNOWN RESPONSE CODE";
    }
  }
}

char *plc4c_data_type_name(plc4c_data_type data_type) {
  switch (data_type) {
    // Boolean/Bit-String Types
    case PLC4C_BOOL:
      return "PLC4C_BOOL";
    case PLC4C_BYTE:
      return "PLC4C_BYTE";
    case PLC4C_WORD:
      return "PLC4C_WORD";
    case PLC4C_DWORD:
      return "PLC4C_DWORD";
    case PLC4C_LWORD:
      return "PLC4C_LWORD";
        // Integer Types (Unsigned/Signed)
        // 1-byte
    case PLC4C_USINT:
      return "PLC4C_USINT";
    case PLC4C_SINT:
      return "PLC4C_SINT";
        // 2-byte
    case PLC4C_UINT:
      return "PLC4C_UINT";
    case PLC4C_INT:
      return "PLC4C_INT";
        // 4-byte
    case PLC4C_UDINT:
      return "PLC4C_UDINT";
    case PLC4C_DINT:
      return "PLC4C_DINT";
        // 8-byte
    case PLC4C_ULINT:
      return "PLC4C_ULINT";
    case PLC4C_LINT:
      return "PLC4C_LINT";
        // Floating-Point Types
    case PLC4C_REAL:
      return "PLC4C_REAL";
    case PLC4C_LREAL:
      return "PLC4C_LREAL";
      // Time Types
    case PLC4C_TIME:
      return "PLC4C_TIME";
    case PLC4C_LTIME:
      return "PLC4C_LTIME";
        // Date Types
    case PLC4C_DATE:
      return "PLC4C_DATE";
    case PLC4C_LDATE:
      return "PLC4C_LDATE";
        // Time of day Types
    case PLC4C_TIME_OF_DAY:
      return "PLC4C_TIME_OF_DAY";
    case PLC4C_LTIME_OF_DAY:
      return "PLC4C_LTIME_OF_DAY";
        // Date and Time Types
    case PLC4C_DATE_AND_TIME:
      return "PLC4C_DATE_AND_TIME";
    case PLC4C_LDATE_AND_TIME:
      return "PLC4C_LDATE_AND_TIME";
        // Char and String Types
    case PLC4C_CHAR:
      return "PLC4C_CHAR";
    case PLC4C_WCHAR:
      return "PLC4C_WCHAR";
    case PLC4C_STRING:
      return "PLC4C_STRING";
    case PLC4C_WSTRING:
      return "PLC4C_WSTRING";

    case PLC4C_LIST:
      return "PLC4C_LIST";
    case PLC4C_STRUCT:
      return "PLC4C_STRUCT";

    default:
      return "UNKNOWN";
    }
  }
