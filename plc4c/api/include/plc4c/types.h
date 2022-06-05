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
#ifndef PLC4C_TYPES_H_
#define PLC4C_TYPES_H_

#include <stdbool.h>

#ifdef __cplusplus
extern "C" {
#endif

/**
 * PLC4C error codes
 */
typedef enum plc4c_return_code {
  UNFINISHED,
  OK,
  NO_MEMORY,
  INVALID_CONNECTION_STRING,
  INVALID_ADDRESS,
  NON_MATCHING_LISTS,
  INVALID_LIST_SIZE,
  NULL_VALUE,
  OUT_OF_RANGE,
  INVALID_ARGUMENT,
  NOT_REACHABLE,
  PERMISSION_DENIED,

  NO_DRIVER_AVAILABLE,
  UNKNOWN_DRIVER,

  UNSPECIFIED_TRANSPORT,
  NO_TRANSPORT_AVAILABLE,
  UNKNOWN_TRANSPORT,

  ALREADY_CONNECTED,

  PARSE_ERROR,

  UNKNOWN_ERROR,
  CONNECTION_ERROR,
  INTERNAL_ERROR,
  NOT_IMPLEMENTED
} plc4c_return_code;

typedef enum plc4c_response_code {
  PLC4C_RESPONSE_CODE_OK,
  PLC4C_RESPONSE_CODE_NOT_FOUND,
  PLC4C_RESPONSE_CODE_ACCESS_DENIED,
  PLC4C_RESPONSE_CODE_INVALID_ADDRESS,
  PLC4C_RESPONSE_CODE_INVALID_DATATYPE,
  PLC4C_RESPONSE_CODE_INTERNAL_ERROR,
  PLC4C_RESPONSE_CODE_RESPONSE_PENDING
} plc4c_response_code;

/**
 * PLC4C data types
 */
typedef enum plc4c_data_type {
  // Boolean/Bit-String Types
  PLC4C_BOOL,
  PLC4C_BYTE,
  PLC4C_WORD,
  PLC4C_DWORD,
  PLC4C_LWORD,
  // Integer Types (Unsigned/Signed)
  // 1-byte
  PLC4C_USINT,
  PLC4C_SINT,
  // 2-byte
  PLC4C_UINT,
  PLC4C_INT,
  // 4-byte
  PLC4C_UDINT,
  PLC4C_DINT,
  // 8-byte
  PLC4C_ULINT,
  PLC4C_LINT,
  // Floating-Point Types
  PLC4C_REAL,
  PLC4C_LREAL,
  // Time Types
  PLC4C_TIME,
  PLC4C_LTIME,
  // Date Types
  PLC4C_DATE,
  PLC4C_LDATE,
  // Time of day Types
  PLC4C_TIME_OF_DAY,
  PLC4C_LTIME_OF_DAY,
  // Date and Time Types
  PLC4C_DATE_AND_TIME,
  PLC4C_LDATE_AND_TIME,
  // Char and String Types
  PLC4C_CHAR,
  PLC4C_WCHAR,
  PLC4C_STRING,
  PLC4C_WSTRING,

  PLC4C_LIST,
  PLC4C_STRUCT
} plc4c_data_type;

/**
 * Helper that translates from a return_code enum value to something a human can
 * work with.
 *
 * @param err return code.
 * @return A human readable description.
 */
char *plc4c_return_code_to_message(plc4c_return_code err);

/**
 * Helper that translates from a plc4c_response_code enum value to something a
 * human can work with.
 *
 * @param response_code plc4c_response_code.
 * @return A human readable description.
 */
char *plc4c_response_code_to_message(plc4c_response_code response_code);

/**
 * Helper function translates from a plc4c_data_type enum value to something a
 * human can work with.
 * @param data_type plc4c_data_type
 * @return string representation
 */
char *plc4c_data_type_name(plc4c_data_type data_type);

/**
 * The plc4c system.
 */
typedef struct plc4c_system_t plc4c_system;

/**
 * The plc4c_driver.
 */
typedef struct plc4c_driver_t plc4c_driver;

/**
 * The plc4c_transport.
 */
typedef struct plc4c_transport_t plc4c_transport;

/**
 * The plc4c_connection.
 */
typedef struct plc4c_connection_t plc4c_connection;

/**
 * A plc4c read-request.
 */
typedef struct plc4c_read_request_t plc4c_read_request;

/**
 * A plc4c read-request-execution.
 */
typedef struct plc4c_read_request_execution_t plc4c_read_request_execution;

/**
 * A plc4c read-request response.
 */
typedef struct plc4c_read_response_t plc4c_read_response;

/**
 * A plc4c write-request.
 */
typedef struct plc4c_write_request_t plc4c_write_request;

/**
 * A plc4c write-request-execution.
 */
typedef struct plc4c_write_request_execution_t plc4c_write_request_execution;

/**
 * A plc4c write-request response.
 */
typedef struct plc4c_write_response_t plc4c_write_response;

/**
 * A plc4c subscription-request.
 */
typedef struct plc4c_subscription_request_t plc4c_subscription_request;

/**
 * A plc4c subscription-request-execution.
 */
typedef struct plc4c_subscription_request_execution_t
    plc4c_subscription_request_execution;

/**
 * A plc4c subscription-response.
 */
typedef struct plc4c_subscription_response_t plc4c_subscription_response;

/**
 * A plc4s subscription-event.
 */
typedef struct plc4c_subscription_event_t plc4c_subscription_event;

/**
 * A plc4c unsubscription-request.
 */
typedef struct plc4c_unsubscription_request_t plc4c_unsubscription_request;

/**
 * A plc4c unsubscription-request-execution.
 */
typedef struct plc4c_unsubscription_request_execution_t
    plc4c_unsubscription_request_execution;

/**
 * A plc4c unsubscription-response.
 */
typedef struct plc4c_unsubscription_response_t plc4c_unsubscription_response;

typedef struct plc4c_data_t plc4c_data;

#ifdef __cplusplus
}
#endif
#endif  // PLC4C_TYPES_H_