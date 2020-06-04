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
#include <plc4c/spi/read_buffer.h>
#include <plc4c/spi/write_buffer.h>
#include <plc4c/spi/evaluation_helper.h>
#include "modbus_pdu_error.h"
#include "modbus_pdu_read_discrete_inputs_request.h"
#include "modbus_pdu_read_discrete_inputs_response.h"
#include "modbus_pdu_read_coils_request.h"
#include "modbus_pdu_read_coils_response.h"
#include "modbus_pdu_write_single_coil_request.h"
#include "modbus_pdu_write_single_coil_response.h"
#include "modbus_pdu_write_multiple_coils_request.h"
#include "modbus_pdu_write_multiple_coils_response.h"
#include "modbus_pdu_read_input_registers_request.h"
#include "modbus_pdu_read_input_registers_response.h"
#include "modbus_pdu_read_holding_registers_request.h"
#include "modbus_pdu_read_holding_registers_response.h"
#include "modbus_pdu_write_single_register_request.h"
#include "modbus_pdu_write_single_register_response.h"
#include "modbus_pdu_write_multiple_holding_registers_request.h"
#include "modbus_pdu_write_multiple_holding_registers_response.h"
#include "modbus_pdu_read_write_multiple_holding_registers_request.h"
#include "modbus_pdu_read_write_multiple_holding_registers_response.h"
#include "modbus_pdu_mask_write_holding_register_request.h"
#include "modbus_pdu_mask_write_holding_register_response.h"
#include "modbus_pdu_read_fifo_queue_request.h"
#include "modbus_pdu_read_fifo_queue_response.h"
#include "modbus_pdu_read_file_record_request.h"
#include "modbus_pdu_read_file_record_response.h"
#include "modbus_pdu_write_file_record_request.h"
#include "modbus_pdu_write_file_record_response.h"
#include "modbus_pdu_read_exception_status_request.h"
#include "modbus_pdu_read_exception_status_response.h"
#include "modbus_pdu_diagnostic_request.h"
#include "modbus_pdu_get_com_event_log_request.h"
#include "modbus_pdu_get_com_event_log_response.h"
#include "modbus_pdu_report_server_id_request.h"
#include "modbus_pdu_report_server_id_response.h"
#include "modbus_pdu_read_device_identification_request.h"
#include "modbus_pdu_read_device_identification_response.h"
#include "modbus_pdu.h"

// Array of discriminator values that match the enum type constants.
// (The order is identical to the enum constants so we can use the
// enum constant to directly access a given types discriminator values)
const plc4c_modbus_read_write_modbus_pdu_discriminator plc4c_modbus_read_write_modbus_pdu_discriminators[] = {
  {/* modbus_read_write_modbus_pdu_diagnostic_request */
   .error = false, .function = 0x08, .response = false},
  {/* modbus_read_write_modbus_pdu_error */
   .error = true, .function = -1, .response = -1},
  {/* modbus_read_write_modbus_pdu_get_com_event_log_request */
   .error = false, .function = 0x0C, .response = false},
  {/* modbus_read_write_modbus_pdu_get_com_event_log_response */
   .error = false, .function = 0x0C, .response = true},
  {/* modbus_read_write_modbus_pdu_mask_write_holding_register_request */
   .error = false, .function = 0x16, .response = false},
  {/* modbus_read_write_modbus_pdu_mask_write_holding_register_response */
   .error = false, .function = 0x16, .response = true},
  {/* modbus_read_write_modbus_pdu_read_coils_request */
   .error = false, .function = 0x01, .response = false},
  {/* modbus_read_write_modbus_pdu_read_coils_response */
   .error = false, .function = 0x01, .response = true},
  {/* modbus_read_write_modbus_pdu_read_device_identification_request */
   .error = false, .function = 0x2B, .response = false},
  {/* modbus_read_write_modbus_pdu_read_device_identification_response */
   .error = false, .function = 0x2B, .response = true},
  {/* modbus_read_write_modbus_pdu_read_discrete_inputs_request */
   .error = false, .function = 0x02, .response = false},
  {/* modbus_read_write_modbus_pdu_read_discrete_inputs_response */
   .error = false, .function = 0x02, .response = true},
  {/* modbus_read_write_modbus_pdu_read_exception_status_request */
   .error = false, .function = 0x07, .response = false},
  {/* modbus_read_write_modbus_pdu_read_exception_status_response */
   .error = false, .function = 0x07, .response = true},
  {/* modbus_read_write_modbus_pdu_read_fifo_queue_request */
   .error = false, .function = 0x18, .response = false},
  {/* modbus_read_write_modbus_pdu_read_fifo_queue_response */
   .error = false, .function = 0x18, .response = true},
  {/* modbus_read_write_modbus_pdu_read_file_record_request */
   .error = false, .function = 0x14, .response = false},
  {/* modbus_read_write_modbus_pdu_read_file_record_response */
   .error = false, .function = 0x14, .response = true},
  {/* modbus_read_write_modbus_pdu_read_holding_registers_request */
   .error = false, .function = 0x03, .response = false},
  {/* modbus_read_write_modbus_pdu_read_holding_registers_response */
   .error = false, .function = 0x03, .response = true},
  {/* modbus_read_write_modbus_pdu_read_input_registers_request */
   .error = false, .function = 0x04, .response = false},
  {/* modbus_read_write_modbus_pdu_read_input_registers_response */
   .error = false, .function = 0x04, .response = true},
  {/* modbus_read_write_modbus_pdu_read_write_multiple_holding_registers_request */
   .error = false, .function = 0x17, .response = false},
  {/* modbus_read_write_modbus_pdu_read_write_multiple_holding_registers_response */
   .error = false, .function = 0x17, .response = true},
  {/* modbus_read_write_modbus_pdu_report_server_id_request */
   .error = false, .function = 0x11, .response = false},
  {/* modbus_read_write_modbus_pdu_report_server_id_response */
   .error = false, .function = 0x11, .response = true},
  {/* modbus_read_write_modbus_pdu_write_file_record_request */
   .error = false, .function = 0x15, .response = false},
  {/* modbus_read_write_modbus_pdu_write_file_record_response */
   .error = false, .function = 0x15, .response = true},
  {/* modbus_read_write_modbus_pdu_write_multiple_coils_request */
   .error = false, .function = 0x0F, .response = false},
  {/* modbus_read_write_modbus_pdu_write_multiple_coils_response */
   .error = false, .function = 0x0F, .response = true},
  {/* modbus_read_write_modbus_pdu_write_multiple_holding_registers_request */
   .error = false, .function = 0x10, .response = false},
  {/* modbus_read_write_modbus_pdu_write_multiple_holding_registers_response */
   .error = false, .function = 0x10, .response = true},
  {/* modbus_read_write_modbus_pdu_write_single_coil_request */
   .error = false, .function = 0x05, .response = false},
  {/* modbus_read_write_modbus_pdu_write_single_coil_response */
   .error = false, .function = 0x05, .response = true},
  {/* modbus_read_write_modbus_pdu_write_single_register_request */
   .error = false, .function = 0x06, .response = false},
  {/* modbus_read_write_modbus_pdu_write_single_register_response */
   .error = false, .function = 0x06, .response = true}
};

// Function returning the discriminator values for a given type constant.
plc4c_modbus_read_write_modbus_pdu_discriminator plc4c_modbus_read_write_modbus_pdu_get_discriminator(plc4c_modbus_read_write_modbus_pdu_type type) {
  return plc4c_modbus_read_write_modbus_pdu_discriminators[type];
}

// Parse function.
plc4c_return_code plc4c_modbus_read_write_modbus_pdu_parse(plc4c_spi_read_buffer* buf, bool response, plc4c_modbus_read_write_modbus_pdu** message) {
  uint16_t startPos = plc4c_spi_read_get_pos(buf);
  uint16_t curPos;

  // Pointer to the parsed datastructure.
  void* msg = NULL;
  // Factory function that allows filling the properties of this type
  void (*factory_ptr)()

  // Implicit Field (error) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
  bool error = plc4c_spi_read_bit(buf);

  // Implicit Field (function) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
  unsigned int function = plc4c_spi_read_unsigned_short(buf, 7);

  // Switch Field (Depending on the discriminator values, passes the instantiation to a sub-type)
  if(error == true) {
    plc4c_modbus_read_write_modbus_pdu_error_parse(buf, response, &msg);
  } else 
  if((error == false) && (function == 0x02) && (response == false)) {
    plc4c_modbus_read_write_modbus_pdu_read_discrete_inputs_request_parse(buf, response, &msg);
  } else 
  if((error == false) && (function == 0x02) && (response == true)) {
    plc4c_modbus_read_write_modbus_pdu_read_discrete_inputs_response_parse(buf, response, &msg);
  } else 
  if((error == false) && (function == 0x01) && (response == false)) {
    plc4c_modbus_read_write_modbus_pdu_read_coils_request_parse(buf, response, &msg);
  } else 
  if((error == false) && (function == 0x01) && (response == true)) {
    plc4c_modbus_read_write_modbus_pdu_read_coils_response_parse(buf, response, &msg);
  } else 
  if((error == false) && (function == 0x05) && (response == false)) {
    plc4c_modbus_read_write_modbus_pdu_write_single_coil_request_parse(buf, response, &msg);
  } else 
  if((error == false) && (function == 0x05) && (response == true)) {
    plc4c_modbus_read_write_modbus_pdu_write_single_coil_response_parse(buf, response, &msg);
  } else 
  if((error == false) && (function == 0x0F) && (response == false)) {
    plc4c_modbus_read_write_modbus_pdu_write_multiple_coils_request_parse(buf, response, &msg);
  } else 
  if((error == false) && (function == 0x0F) && (response == true)) {
    plc4c_modbus_read_write_modbus_pdu_write_multiple_coils_response_parse(buf, response, &msg);
  } else 
  if((error == false) && (function == 0x04) && (response == false)) {
    plc4c_modbus_read_write_modbus_pdu_read_input_registers_request_parse(buf, response, &msg);
  } else 
  if((error == false) && (function == 0x04) && (response == true)) {
    plc4c_modbus_read_write_modbus_pdu_read_input_registers_response_parse(buf, response, &msg);
  } else 
  if((error == false) && (function == 0x03) && (response == false)) {
    plc4c_modbus_read_write_modbus_pdu_read_holding_registers_request_parse(buf, response, &msg);
  } else 
  if((error == false) && (function == 0x03) && (response == true)) {
    plc4c_modbus_read_write_modbus_pdu_read_holding_registers_response_parse(buf, response, &msg);
  } else 
  if((error == false) && (function == 0x06) && (response == false)) {
    plc4c_modbus_read_write_modbus_pdu_write_single_register_request_parse(buf, response, &msg);
  } else 
  if((error == false) && (function == 0x06) && (response == true)) {
    plc4c_modbus_read_write_modbus_pdu_write_single_register_response_parse(buf, response, &msg);
  } else 
  if((error == false) && (function == 0x10) && (response == false)) {
    plc4c_modbus_read_write_modbus_pdu_write_multiple_holding_registers_request_parse(buf, response, &msg);
  } else 
  if((error == false) && (function == 0x10) && (response == true)) {
    plc4c_modbus_read_write_modbus_pdu_write_multiple_holding_registers_response_parse(buf, response, &msg);
  } else 
  if((error == false) && (function == 0x17) && (response == false)) {
    plc4c_modbus_read_write_modbus_pdu_read_write_multiple_holding_registers_request_parse(buf, response, &msg);
  } else 
  if((error == false) && (function == 0x17) && (response == true)) {
    plc4c_modbus_read_write_modbus_pdu_read_write_multiple_holding_registers_response_parse(buf, response, &msg);
  } else 
  if((error == false) && (function == 0x16) && (response == false)) {
    plc4c_modbus_read_write_modbus_pdu_mask_write_holding_register_request_parse(buf, response, &msg);
  } else 
  if((error == false) && (function == 0x16) && (response == true)) {
    plc4c_modbus_read_write_modbus_pdu_mask_write_holding_register_response_parse(buf, response, &msg);
  } else 
  if((error == false) && (function == 0x18) && (response == false)) {
    plc4c_modbus_read_write_modbus_pdu_read_fifo_queue_request_parse(buf, response, &msg);
  } else 
  if((error == false) && (function == 0x18) && (response == true)) {
    plc4c_modbus_read_write_modbus_pdu_read_fifo_queue_response_parse(buf, response, &msg);
  } else 
  if((error == false) && (function == 0x14) && (response == false)) {
    plc4c_modbus_read_write_modbus_pdu_read_file_record_request_parse(buf, response, &msg);
  } else 
  if((error == false) && (function == 0x14) && (response == true)) {
    plc4c_modbus_read_write_modbus_pdu_read_file_record_response_parse(buf, response, &msg);
  } else 
  if((error == false) && (function == 0x15) && (response == false)) {
    plc4c_modbus_read_write_modbus_pdu_write_file_record_request_parse(buf, response, &msg);
  } else 
  if((error == false) && (function == 0x15) && (response == true)) {
    plc4c_modbus_read_write_modbus_pdu_write_file_record_response_parse(buf, response, &msg);
  } else 
  if((error == false) && (function == 0x07) && (response == false)) {
    plc4c_modbus_read_write_modbus_pdu_read_exception_status_request_parse(buf, response, &msg);
  } else 
  if((error == false) && (function == 0x07) && (response == true)) {
    plc4c_modbus_read_write_modbus_pdu_read_exception_status_response_parse(buf, response, &msg);
  } else 
  if((error == false) && (function == 0x08) && (response == false)) {
    plc4c_modbus_read_write_modbus_pdu_diagnostic_request_parse(buf, response, &msg);
  } else 
  if((error == false) && (function == 0x0C) && (response == false)) {
    plc4c_modbus_read_write_modbus_pdu_get_com_event_log_request_parse(buf, response, &msg);
  } else 
  if((error == false) && (function == 0x0C) && (response == true)) {
    plc4c_modbus_read_write_modbus_pdu_get_com_event_log_response_parse(buf, response, &msg);
  } else 
  if((error == false) && (function == 0x11) && (response == false)) {
    plc4c_modbus_read_write_modbus_pdu_report_server_id_request_parse(buf, response, &msg);
  } else 
  if((error == false) && (function == 0x11) && (response == true)) {
    plc4c_modbus_read_write_modbus_pdu_report_server_id_response_parse(buf, response, &msg);
  } else 
  if((error == false) && (function == 0x2B) && (response == false)) {
    plc4c_modbus_read_write_modbus_pdu_read_device_identification_request_parse(buf, response, &msg);
  } else 
  if((error == false) && (function == 0x2B) && (response == true)) {
    plc4c_modbus_read_write_modbus_pdu_read_device_identification_response_parse(buf, response, &msg);
  }

  return OK;
}

plc4c_return_code plc4c_modbus_read_write_modbus_pdu_serialize(plc4c_spi_write_buffer* buf, plc4c_modbus_read_write_modbus_pdu* message) {
  return OK;
}
