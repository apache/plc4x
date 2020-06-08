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
#ifndef PLC4C_S7_READ_WRITE_S7_MESSAGE_H_
#define PLC4C_S7_READ_WRITE_S7_MESSAGE_H_
#ifdef __cplusplus
extern "C" {
#endif

#include <stdbool.h>
#include <stdint.h>
#include <plc4c/utils/list.h>
#include "s7_parameter.h"
#include "s7_payload.h"

// Structure used to contain the discriminator values for discriminated types using this as a parent
struct plc4c_s7_read_write_s7_message_discriminator {
  uint8_t messageType;
};
typedef struct plc4c_s7_read_write_s7_message_discriminator plc4c_s7_read_write_s7_message_discriminator;

// Enum assigning each sub-type an individual id.
enum plc4c_s7_read_write_s7_message_type {
  plc4c_s7_read_write_s7_message_type_s7_read_write_s7_message_request = 0,
  plc4c_s7_read_write_s7_message_type_s7_read_write_s7_message_response = 1,
  plc4c_s7_read_write_s7_message_type_s7_read_write_s7_message_response_data = 2,
  plc4c_s7_read_write_s7_message_type_s7_read_write_s7_message_user_data = 3};
typedef enum plc4c_s7_read_write_s7_message_type plc4c_s7_read_write_s7_message_type;

// Function to get the discriminator values for a given type.
plc4c_s7_read_write_s7_message_discriminator plc4c_s7_read_write_s7_message_get_discriminator(plc4c_s7_read_write_s7_message_type type);

// Constant values.
const uint8_t S7_READ_WRITE_S7_MESSAGE_PROTOCOL_ID = 0x32;

struct plc4c_s7_read_write_s7_message {
  /* This is an abstract type so this property saves the type of this typed union */
  plc4c_s7_read_write_s7_message_type _type;
  /* Properties */
  uint8_t protocol_id;
  uint16_t tpdu_reference;
  union {
    struct { /* S7MessageRequest */
    };
    struct { /* S7MessageResponse */
      uint8_t s7_message_response_error_class;
      uint8_t s7_message_response_error_code;
    };
    struct { /* S7MessageResponseData */
      uint8_t s7_message_response_data_error_class;
      uint8_t s7_message_response_data_error_code;
    };
    struct { /* S7MessageUserData */
    };
  };
  plc4c_s7_read_write_s7_parameter* parameter;
  plc4c_s7_read_write_s7_payload* payload;
};
typedef struct plc4c_s7_read_write_s7_message plc4c_s7_read_write_s7_message;

plc4c_return_code plc4c_s7_read_write_s7_message_parse(plc4c_spi_read_buffer* buf, plc4c_s7_read_write_s7_message** message);

plc4c_return_code plc4c_s7_read_write_s7_message_serialize(plc4c_spi_write_buffer* buf, plc4c_s7_read_write_s7_message* message);

#ifdef __cplusplus
}
#endif
#endif  // PLC4C_S7_READ_WRITE_S7_MESSAGE_H_
