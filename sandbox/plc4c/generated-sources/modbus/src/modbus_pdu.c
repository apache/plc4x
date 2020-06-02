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

#include <plc4c/spi/read_buffer.h>
#include <plc4c/spi/write_buffer.h>
#include <plc4c/spi/evaluation_helper.h>

#include "modbus_pdu.h"

plc4c_return_code plc4c_modbus_read_write_modbus_pdu_parse(plc4c_read_buffer buf, bool response, plc4c_modbus_read_write_modbus_pdu** message) {
  uint16_t start_pos = plc4c_spi_read_get_pos(buf);
  uint16_t cur_pos;

  plc4c_modbus_read_write_modbus_pdu* msg = malloc(sizeof(plc4c_modbus_read_write_modbus_pdu));

  // Implicit Field (error) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
  bool error = plc4c_spi_read_bit(buf);

  // Implicit Field (function) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
  unsigned int function = plc4c_spi_read_unsigned_short(buf, 7);

  // Switch Field (Depending on the discriminator values, passes the instantiation to a sub-type)
  if(plc4c_spi_evaluation_helper_equals(error, true)) {
    plc4c_modbus_read_write_modbus_pdu_error_parse(buf, msg, response);
  } else 
  if(plc4c_spi_evaluation_helper_equals(error, false) && plc4c_spi_evaluation_helper_equals(function, 0x02) && plc4c_spi_evaluation_helper_equals(response, false)) {
    plc4c_modbus_read_write_modbus_pdu_read_discrete_inputs_request_parse(buf, msg, response);
  } else 
  if(plc4c_spi_evaluation_helper_equals(error, false) && plc4c_spi_evaluation_helper_equals(function, 0x02) && plc4c_spi_evaluation_helper_equals(response, true)) {
    plc4c_modbus_read_write_modbus_pdu_read_discrete_inputs_response_parse(buf, msg, response);
  } else 
  if(plc4c_spi_evaluation_helper_equals(error, false) && plc4c_spi_evaluation_helper_equals(function, 0x01) && plc4c_spi_evaluation_helper_equals(response, false)) {
    plc4c_modbus_read_write_modbus_pdu_read_coils_request_parse(buf, msg, response);
  } else 
  if(plc4c_spi_evaluation_helper_equals(error, false) && plc4c_spi_evaluation_helper_equals(function, 0x01) && plc4c_spi_evaluation_helper_equals(response, true)) {
    plc4c_modbus_read_write_modbus_pdu_read_coils_response_parse(buf, msg, response);
  } else 
  if(plc4c_spi_evaluation_helper_equals(error, false) && plc4c_spi_evaluation_helper_equals(function, 0x05) && plc4c_spi_evaluation_helper_equals(response, false)) {
    plc4c_modbus_read_write_modbus_pdu_write_single_coil_request_parse(buf, msg, response);
  } else 
  if(plc4c_spi_evaluation_helper_equals(error, false) && plc4c_spi_evaluation_helper_equals(function, 0x05) && plc4c_spi_evaluation_helper_equals(response, true)) {
    plc4c_modbus_read_write_modbus_pdu_write_single_coil_response_parse(buf, msg, response);
  } else 
  if(plc4c_spi_evaluation_helper_equals(error, false) && plc4c_spi_evaluation_helper_equals(function, 0x0F) && plc4c_spi_evaluation_helper_equals(response, false)) {
    plc4c_modbus_read_write_modbus_pdu_write_multiple_coils_request_parse(buf, msg, response);
  } else 
  if(plc4c_spi_evaluation_helper_equals(error, false) && plc4c_spi_evaluation_helper_equals(function, 0x0F) && plc4c_spi_evaluation_helper_equals(response, true)) {
    plc4c_modbus_read_write_modbus_pdu_write_multiple_coils_response_parse(buf, msg, response);
  } else 
  if(plc4c_spi_evaluation_helper_equals(error, false) && plc4c_spi_evaluation_helper_equals(function, 0x04) && plc4c_spi_evaluation_helper_equals(response, false)) {
    plc4c_modbus_read_write_modbus_pdu_read_input_registers_request_parse(buf, msg, response);
  } else 
  if(plc4c_spi_evaluation_helper_equals(error, false) && plc4c_spi_evaluation_helper_equals(function, 0x04) && plc4c_spi_evaluation_helper_equals(response, true)) {
    plc4c_modbus_read_write_modbus_pdu_read_input_registers_response_parse(buf, msg, response);
  } else 
  if(plc4c_spi_evaluation_helper_equals(error, false) && plc4c_spi_evaluation_helper_equals(function, 0x03) && plc4c_spi_evaluation_helper_equals(response, false)) {
    plc4c_modbus_read_write_modbus_pdu_read_holding_registers_request_parse(buf, msg, response);
  } else 
  if(plc4c_spi_evaluation_helper_equals(error, false) && plc4c_spi_evaluation_helper_equals(function, 0x03) && plc4c_spi_evaluation_helper_equals(response, true)) {
    plc4c_modbus_read_write_modbus_pdu_read_holding_registers_response_parse(buf, msg, response);
  } else 
  if(plc4c_spi_evaluation_helper_equals(error, false) && plc4c_spi_evaluation_helper_equals(function, 0x06) && plc4c_spi_evaluation_helper_equals(response, false)) {
    plc4c_modbus_read_write_modbus_pdu_write_single_register_request_parse(buf, msg, response);
  } else 
  if(plc4c_spi_evaluation_helper_equals(error, false) && plc4c_spi_evaluation_helper_equals(function, 0x06) && plc4c_spi_evaluation_helper_equals(response, true)) {
    plc4c_modbus_read_write_modbus_pdu_write_single_register_response_parse(buf, msg, response);
  } else 
  if(plc4c_spi_evaluation_helper_equals(error, false) && plc4c_spi_evaluation_helper_equals(function, 0x10) && plc4c_spi_evaluation_helper_equals(response, false)) {
    plc4c_modbus_read_write_modbus_pdu_write_multiple_holding_registers_request_parse(buf, msg, response);
  } else 
  if(plc4c_spi_evaluation_helper_equals(error, false) && plc4c_spi_evaluation_helper_equals(function, 0x10) && plc4c_spi_evaluation_helper_equals(response, true)) {
    plc4c_modbus_read_write_modbus_pdu_write_multiple_holding_registers_response_parse(buf, msg, response);
  } else 
  if(plc4c_spi_evaluation_helper_equals(error, false) && plc4c_spi_evaluation_helper_equals(function, 0x17) && plc4c_spi_evaluation_helper_equals(response, false)) {
    plc4c_modbus_read_write_modbus_pdu_read_write_multiple_holding_registers_request_parse(buf, msg, response);
  } else 
  if(plc4c_spi_evaluation_helper_equals(error, false) && plc4c_spi_evaluation_helper_equals(function, 0x17) && plc4c_spi_evaluation_helper_equals(response, true)) {
    plc4c_modbus_read_write_modbus_pdu_read_write_multiple_holding_registers_response_parse(buf, msg, response);
  } else 
  if(plc4c_spi_evaluation_helper_equals(error, false) && plc4c_spi_evaluation_helper_equals(function, 0x16) && plc4c_spi_evaluation_helper_equals(response, false)) {
    plc4c_modbus_read_write_modbus_pdu_mask_write_holding_register_request_parse(buf, msg, response);
  } else 
  if(plc4c_spi_evaluation_helper_equals(error, false) && plc4c_spi_evaluation_helper_equals(function, 0x16) && plc4c_spi_evaluation_helper_equals(response, true)) {
    plc4c_modbus_read_write_modbus_pdu_mask_write_holding_register_response_parse(buf, msg, response);
  } else 
  if(plc4c_spi_evaluation_helper_equals(error, false) && plc4c_spi_evaluation_helper_equals(function, 0x18) && plc4c_spi_evaluation_helper_equals(response, false)) {
    plc4c_modbus_read_write_modbus_pdu_read_fifo_queue_request_parse(buf, msg, response);
  } else 
  if(plc4c_spi_evaluation_helper_equals(error, false) && plc4c_spi_evaluation_helper_equals(function, 0x18) && plc4c_spi_evaluation_helper_equals(response, true)) {
    plc4c_modbus_read_write_modbus_pdu_read_fifo_queue_response_parse(buf, msg, response);
  } else 
  if(plc4c_spi_evaluation_helper_equals(error, false) && plc4c_spi_evaluation_helper_equals(function, 0x14) && plc4c_spi_evaluation_helper_equals(response, false)) {
    plc4c_modbus_read_write_modbus_pdu_read_file_record_request_parse(buf, msg, response);
  } else 
  if(plc4c_spi_evaluation_helper_equals(error, false) && plc4c_spi_evaluation_helper_equals(function, 0x14) && plc4c_spi_evaluation_helper_equals(response, true)) {
    plc4c_modbus_read_write_modbus_pdu_read_file_record_response_parse(buf, msg, response);
  } else 
  if(plc4c_spi_evaluation_helper_equals(error, false) && plc4c_spi_evaluation_helper_equals(function, 0x15) && plc4c_spi_evaluation_helper_equals(response, false)) {
    plc4c_modbus_read_write_modbus_pdu_write_file_record_request_parse(buf, msg, response);
  } else 
  if(plc4c_spi_evaluation_helper_equals(error, false) && plc4c_spi_evaluation_helper_equals(function, 0x15) && plc4c_spi_evaluation_helper_equals(response, true)) {
    plc4c_modbus_read_write_modbus_pdu_write_file_record_response_parse(buf, msg, response);
  } else 
  if(plc4c_spi_evaluation_helper_equals(error, false) && plc4c_spi_evaluation_helper_equals(function, 0x07) && plc4c_spi_evaluation_helper_equals(response, false)) {
    plc4c_modbus_read_write_modbus_pdu_read_exception_status_request_parse(buf, msg, response);
  } else 
  if(plc4c_spi_evaluation_helper_equals(error, false) && plc4c_spi_evaluation_helper_equals(function, 0x07) && plc4c_spi_evaluation_helper_equals(response, true)) {
    plc4c_modbus_read_write_modbus_pdu_read_exception_status_response_parse(buf, msg, response);
  } else 
  if(plc4c_spi_evaluation_helper_equals(error, false) && plc4c_spi_evaluation_helper_equals(function, 0x08) && plc4c_spi_evaluation_helper_equals(response, false)) {
    plc4c_modbus_read_write_modbus_pdu_diagnostic_request_parse(buf, msg, response);
  } else 
  if(plc4c_spi_evaluation_helper_equals(error, false) && plc4c_spi_evaluation_helper_equals(function, 0x08) && plc4c_spi_evaluation_helper_equals(response, true)) {
    plc4c_modbus_read_write_modbus_pdu_diagnostic_response_parse(buf, msg, response);
  } else 
  if(plc4c_spi_evaluation_helper_equals(error, false) && plc4c_spi_evaluation_helper_equals(function, 0x0B) && plc4c_spi_evaluation_helper_equals(response, false)) {
    plc4c_modbus_read_write_modbus_pdu_get_com_event_counter_request_parse(buf, msg, response);
  } else 
  if(plc4c_spi_evaluation_helper_equals(error, false) && plc4c_spi_evaluation_helper_equals(function, 0x0B) && plc4c_spi_evaluation_helper_equals(response, true)) {
    plc4c_modbus_read_write_modbus_pdu_get_com_event_counter_response_parse(buf, msg, response);
  } else 
  if(plc4c_spi_evaluation_helper_equals(error, false) && plc4c_spi_evaluation_helper_equals(function, 0x0C) && plc4c_spi_evaluation_helper_equals(response, false)) {
    plc4c_modbus_read_write_modbus_pdu_get_com_event_log_request_parse(buf, msg, response);
  } else 
  if(plc4c_spi_evaluation_helper_equals(error, false) && plc4c_spi_evaluation_helper_equals(function, 0x0C) && plc4c_spi_evaluation_helper_equals(response, true)) {
    plc4c_modbus_read_write_modbus_pdu_get_com_event_log_response_parse(buf, msg, response);
  } else 
  if(plc4c_spi_evaluation_helper_equals(error, false) && plc4c_spi_evaluation_helper_equals(function, 0x11) && plc4c_spi_evaluation_helper_equals(response, false)) {
    plc4c_modbus_read_write_modbus_pdu_report_server_id_request_parse(buf, msg, response);
  } else 
  if(plc4c_spi_evaluation_helper_equals(error, false) && plc4c_spi_evaluation_helper_equals(function, 0x11) && plc4c_spi_evaluation_helper_equals(response, true)) {
    plc4c_modbus_read_write_modbus_pdu_report_server_id_response_parse(buf, msg, response);
  } else 
  if(plc4c_spi_evaluation_helper_equals(error, false) && plc4c_spi_evaluation_helper_equals(function, 0x2B) && plc4c_spi_evaluation_helper_equals(response, false)) {
    plc4c_modbus_read_write_modbus_pdu_read_device_identification_request_parse(buf, msg, response);
  } else 
  if(plc4c_spi_evaluation_helper_equals(error, false) && plc4c_spi_evaluation_helper_equals(function, 0x2B) && plc4c_spi_evaluation_helper_equals(response, true)) {
    plc4c_modbus_read_write_modbus_pdu_read_device_identification_response_parse(buf, msg, response);
  }

  return OK;
}

plc4c_return_code plc4c_modbus_read_write_modbus_pdu_serialize(plc4c_write_buffer buf, plc4c_modbus_read_write_modbus_pdu* message) {
  return OK;
}
