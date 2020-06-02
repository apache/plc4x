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
#ifndef PLC4C_MODBUS_READ_WRITE_MODBUS_PDU_H_
#define PLC4C_MODBUS_READ_WRITE_MODBUS_PDU_H_
#ifdef __cplusplus
extern "C" {
#endif

#include <stdbool.h>
#include <stdint.h>
#include <plc4c/utils/list.h>

// Enum assigning each sub-type an individual id.
enum plc4c_modbus_read_write_modbus_pdu_type {
  plc4c_modbus_read_write_modbus_pdu_type_modbus_read_write_modbus_pdu_error = 0,
  plc4c_modbus_read_write_modbus_pdu_type_modbus_read_write_modbus_pdu_read_discrete_inputs_request = 1,
  plc4c_modbus_read_write_modbus_pdu_type_modbus_read_write_modbus_pdu_read_discrete_inputs_response = 2,
  plc4c_modbus_read_write_modbus_pdu_type_modbus_read_write_modbus_pdu_read_coils_request = 3,
  plc4c_modbus_read_write_modbus_pdu_type_modbus_read_write_modbus_pdu_read_coils_response = 4,
  plc4c_modbus_read_write_modbus_pdu_type_modbus_read_write_modbus_pdu_write_single_coil_request = 5,
  plc4c_modbus_read_write_modbus_pdu_type_modbus_read_write_modbus_pdu_write_single_coil_response = 6,
  plc4c_modbus_read_write_modbus_pdu_type_modbus_read_write_modbus_pdu_write_multiple_coils_request = 7,
  plc4c_modbus_read_write_modbus_pdu_type_modbus_read_write_modbus_pdu_write_multiple_coils_response = 8,
  plc4c_modbus_read_write_modbus_pdu_type_modbus_read_write_modbus_pdu_read_input_registers_request = 9,
  plc4c_modbus_read_write_modbus_pdu_type_modbus_read_write_modbus_pdu_read_input_registers_response = 10,
  plc4c_modbus_read_write_modbus_pdu_type_modbus_read_write_modbus_pdu_read_holding_registers_request = 11,
  plc4c_modbus_read_write_modbus_pdu_type_modbus_read_write_modbus_pdu_read_holding_registers_response = 12,
  plc4c_modbus_read_write_modbus_pdu_type_modbus_read_write_modbus_pdu_write_single_register_request = 13,
  plc4c_modbus_read_write_modbus_pdu_type_modbus_read_write_modbus_pdu_write_single_register_response = 14,
  plc4c_modbus_read_write_modbus_pdu_type_modbus_read_write_modbus_pdu_write_multiple_holding_registers_request = 15,
  plc4c_modbus_read_write_modbus_pdu_type_modbus_read_write_modbus_pdu_write_multiple_holding_registers_response = 16,
  plc4c_modbus_read_write_modbus_pdu_type_modbus_read_write_modbus_pdu_read_write_multiple_holding_registers_request = 17,
  plc4c_modbus_read_write_modbus_pdu_type_modbus_read_write_modbus_pdu_read_write_multiple_holding_registers_response = 18,
  plc4c_modbus_read_write_modbus_pdu_type_modbus_read_write_modbus_pdu_mask_write_holding_register_request = 19,
  plc4c_modbus_read_write_modbus_pdu_type_modbus_read_write_modbus_pdu_mask_write_holding_register_response = 20,
  plc4c_modbus_read_write_modbus_pdu_type_modbus_read_write_modbus_pdu_read_fifo_queue_request = 21,
  plc4c_modbus_read_write_modbus_pdu_type_modbus_read_write_modbus_pdu_read_fifo_queue_response = 22,
  plc4c_modbus_read_write_modbus_pdu_type_modbus_read_write_modbus_pdu_read_file_record_request = 23,
  plc4c_modbus_read_write_modbus_pdu_type_modbus_read_write_modbus_pdu_read_file_record_response = 24,
  plc4c_modbus_read_write_modbus_pdu_type_modbus_read_write_modbus_pdu_write_file_record_request = 25,
  plc4c_modbus_read_write_modbus_pdu_type_modbus_read_write_modbus_pdu_write_file_record_response = 26,
  plc4c_modbus_read_write_modbus_pdu_type_modbus_read_write_modbus_pdu_read_exception_status_request = 27,
  plc4c_modbus_read_write_modbus_pdu_type_modbus_read_write_modbus_pdu_read_exception_status_response = 28,
  plc4c_modbus_read_write_modbus_pdu_type_modbus_read_write_modbus_pdu_diagnostic_request = 29,
  plc4c_modbus_read_write_modbus_pdu_type_modbus_read_write_modbus_pdu_get_com_event_log_request = 30,
  plc4c_modbus_read_write_modbus_pdu_type_modbus_read_write_modbus_pdu_get_com_event_log_response = 31,
  plc4c_modbus_read_write_modbus_pdu_type_modbus_read_write_modbus_pdu_report_server_id_request = 32,
  plc4c_modbus_read_write_modbus_pdu_type_modbus_read_write_modbus_pdu_report_server_id_response = 33,
  plc4c_modbus_read_write_modbus_pdu_type_modbus_read_write_modbus_pdu_read_device_identification_request = 34,
  plc4c_modbus_read_write_modbus_pdu_type_modbus_read_write_modbus_pdu_read_device_identification_response = 35};
typedef enum plc4c_modbus_read_write_modbus_pdu_type plc4c_modbus_read_write_modbus_pdu_type;

struct plc4c_modbus_read_write_modbus_pdu {
  plc4c_modbus_read_write_modbus_pdu_type _type;
};
typedef struct plc4c_modbus_read_write_modbus_pdu plc4c_modbus_read_write_modbus_pdu;

plc4c_return_code plc4c_modbus_read_write_modbus_pdu_parse(plc4c_spi_read_buffer* buf, bool response, plc4c_modbus_read_write_modbus_pdu** message);

plc4c_return_code plc4c_modbus_read_write_modbus_pdu_serialize(plc4c_spi_write_buffer* buf, plc4c_modbus_read_write_modbus_pdu* message);

#ifdef __cplusplus
}
#endif
#endif  // PLC4C_MODBUS_READ_WRITE_MODBUS_PDU_H_
