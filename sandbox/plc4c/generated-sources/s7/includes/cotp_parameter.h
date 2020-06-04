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
#ifndef PLC4C_S7_READ_WRITE_COTP_PARAMETER_H_
#define PLC4C_S7_READ_WRITE_COTP_PARAMETER_H_
#ifdef __cplusplus
extern "C" {
#endif

#include <stdbool.h>
#include <stdint.h>
#include <plc4c/utils/list.h>

// Structure used to contain the discriminator values for discriminated types using this as a parent
struct plc4c_s7_read_write_cotp_parameter_discriminator {
  uint8_t parameterType;
};
typedef struct plc4c_s7_read_write_cotp_parameter_discriminator plc4c_s7_read_write_cotp_parameter_discriminator;

// Enum assigning each sub-type an individual id.
enum plc4c_s7_read_write_cotp_parameter_type {
  plc4c_s7_read_write_cotp_parameter_type_s7_read_write_cotp_parameter_tpdu_size = 0,
  plc4c_s7_read_write_cotp_parameter_type_s7_read_write_cotp_parameter_calling_tsap = 1,
  plc4c_s7_read_write_cotp_parameter_type_s7_read_write_cotp_parameter_called_tsap = 2,
  plc4c_s7_read_write_cotp_parameter_type_s7_read_write_cotp_parameter_checksum = 3,
  plc4c_s7_read_write_cotp_parameter_type_s7_read_write_cotp_parameter_disconnect_additional_information = 4};
typedef enum plc4c_s7_read_write_cotp_parameter_type plc4c_s7_read_write_cotp_parameter_type;

// Function to get the discriminator values for a given type.
plc4c_s7_read_write_cotp_parameter_discriminator plc4c_s7_read_write_cotp_parameter_get_discriminator(plc4c_s7_read_write_cotp_parameter_type type);

struct plc4c_s7_read_write_cotp_parameter {
  plc4c_s7_read_write_cotp_parameter_type _type;
};
typedef struct plc4c_s7_read_write_cotp_parameter plc4c_s7_read_write_cotp_parameter;

plc4c_return_code plc4c_s7_read_write_cotp_parameter_parse(plc4c_spi_read_buffer* buf, uint8_t rest, plc4c_s7_read_write_cotp_parameter** message);

plc4c_return_code plc4c_s7_read_write_cotp_parameter_serialize(plc4c_spi_write_buffer* buf, plc4c_s7_read_write_cotp_parameter* message);

#ifdef __cplusplus
}
#endif
#endif  // PLC4C_S7_READ_WRITE_COTP_PARAMETER_H_
