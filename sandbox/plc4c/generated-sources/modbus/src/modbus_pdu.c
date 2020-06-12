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

  // Pointer to the parsed data structure.
  plc4c_modbus_read_write_modbus_pdu* msg = malloc(sizeof(plc4c_modbus_read_write_modbus_pdu));
  // Implicit Field (error) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
  bool error = plc4c_spi_read_bit(buf);

  // Implicit Field (function) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
  unsigned int function = plc4c_spi_read_unsigned_short(buf, 7);

  // Switch Field (Depending on the discriminator values, passes the instantiation to a sub-type)
  if(error == true) { /* ModbusPDUError */
    uint8_t exceptionCode = -1;
    msg->modbus_pdu_error_exception_code = exceptionCode;
  } else 
  if((error == false) && (function == 0x02) && (response == false)) { /* ModbusPDUReadDiscreteInputsRequest */
    uint16_t startingAddress = -1;
    msg->modbus_pdu_read_discrete_inputs_request_starting_address = startingAddress;

    uint16_t quantity = -1;
    msg->modbus_pdu_read_discrete_inputs_request_quantity = quantity;
  } else 
  if((error == false) && (function == 0x02) && (response == true)) { /* ModbusPDUReadDiscreteInputsResponse */
    plc4c_list value;
    msg->modbus_pdu_read_discrete_inputs_response_value = value;
  } else 
  if((error == false) && (function == 0x01) && (response == false)) { /* ModbusPDUReadCoilsRequest */
    uint16_t startingAddress = -1;
    msg->modbus_pdu_read_coils_request_starting_address = startingAddress;

    uint16_t quantity = -1;
    msg->modbus_pdu_read_coils_request_quantity = quantity;
  } else 
  if((error == false) && (function == 0x01) && (response == true)) { /* ModbusPDUReadCoilsResponse */
    plc4c_list value;
    msg->modbus_pdu_read_coils_response_value = value;
  } else 
  if((error == false) && (function == 0x05) && (response == false)) { /* ModbusPDUWriteSingleCoilRequest */
    uint16_t address = -1;
    msg->modbus_pdu_write_single_coil_request_address = address;

    uint16_t value = -1;
    msg->modbus_pdu_write_single_coil_request_value = value;
  } else 
  if((error == false) && (function == 0x05) && (response == true)) { /* ModbusPDUWriteSingleCoilResponse */
    uint16_t address = -1;
    msg->modbus_pdu_write_single_coil_response_address = address;

    uint16_t value = -1;
    msg->modbus_pdu_write_single_coil_response_value = value;
  } else 
  if((error == false) && (function == 0x0F) && (response == false)) { /* ModbusPDUWriteMultipleCoilsRequest */
    uint16_t startingAddress = -1;
    msg->modbus_pdu_write_multiple_coils_request_starting_address = startingAddress;

    uint16_t quantity = -1;
    msg->modbus_pdu_write_multiple_coils_request_quantity = quantity;

    plc4c_list value;
    msg->modbus_pdu_write_multiple_coils_request_value = value;
  } else 
  if((error == false) && (function == 0x0F) && (response == true)) { /* ModbusPDUWriteMultipleCoilsResponse */
    uint16_t startingAddress = -1;
    msg->modbus_pdu_write_multiple_coils_response_starting_address = startingAddress;

    uint16_t quantity = -1;
    msg->modbus_pdu_write_multiple_coils_response_quantity = quantity;
  } else 
  if((error == false) && (function == 0x04) && (response == false)) { /* ModbusPDUReadInputRegistersRequest */
    uint16_t startingAddress = -1;
    msg->modbus_pdu_read_input_registers_request_starting_address = startingAddress;

    uint16_t quantity = -1;
    msg->modbus_pdu_read_input_registers_request_quantity = quantity;
  } else 
  if((error == false) && (function == 0x04) && (response == true)) { /* ModbusPDUReadInputRegistersResponse */
    plc4c_list value;
    msg->modbus_pdu_read_input_registers_response_value = value;
  } else 
  if((error == false) && (function == 0x03) && (response == false)) { /* ModbusPDUReadHoldingRegistersRequest */
    uint16_t startingAddress = -1;
    msg->modbus_pdu_read_holding_registers_request_starting_address = startingAddress;

    uint16_t quantity = -1;
    msg->modbus_pdu_read_holding_registers_request_quantity = quantity;
  } else 
  if((error == false) && (function == 0x03) && (response == true)) { /* ModbusPDUReadHoldingRegistersResponse */
    plc4c_list value;
    msg->modbus_pdu_read_holding_registers_response_value = value;
  } else 
  if((error == false) && (function == 0x06) && (response == false)) { /* ModbusPDUWriteSingleRegisterRequest */
    uint16_t address = -1;
    msg->modbus_pdu_write_single_register_request_address = address;

    uint16_t value = -1;
    msg->modbus_pdu_write_single_register_request_value = value;
  } else 
  if((error == false) && (function == 0x06) && (response == true)) { /* ModbusPDUWriteSingleRegisterResponse */
    uint16_t address = -1;
    msg->modbus_pdu_write_single_register_response_address = address;

    uint16_t value = -1;
    msg->modbus_pdu_write_single_register_response_value = value;
  } else 
  if((error == false) && (function == 0x10) && (response == false)) { /* ModbusPDUWriteMultipleHoldingRegistersRequest */
    uint16_t startingAddress = -1;
    msg->modbus_pdu_write_multiple_holding_registers_request_starting_address = startingAddress;

    uint16_t quantity = -1;
    msg->modbus_pdu_write_multiple_holding_registers_request_quantity = quantity;

    plc4c_list value;
    msg->modbus_pdu_write_multiple_holding_registers_request_value = value;
  } else 
  if((error == false) && (function == 0x10) && (response == true)) { /* ModbusPDUWriteMultipleHoldingRegistersResponse */
    uint16_t startingAddress = -1;
    msg->modbus_pdu_write_multiple_holding_registers_response_starting_address = startingAddress;

    uint16_t quantity = -1;
    msg->modbus_pdu_write_multiple_holding_registers_response_quantity = quantity;
  } else 
  if((error == false) && (function == 0x17) && (response == false)) { /* ModbusPDUReadWriteMultipleHoldingRegistersRequest */
    uint16_t readStartingAddress = -1;
    msg->modbus_pdu_read_write_multiple_holding_registers_request_read_starting_address = readStartingAddress;

    uint16_t readQuantity = -1;
    msg->modbus_pdu_read_write_multiple_holding_registers_request_read_quantity = readQuantity;

    uint16_t writeStartingAddress = -1;
    msg->modbus_pdu_read_write_multiple_holding_registers_request_write_starting_address = writeStartingAddress;

    uint16_t writeQuantity = -1;
    msg->modbus_pdu_read_write_multiple_holding_registers_request_write_quantity = writeQuantity;

    plc4c_list value;
    msg->modbus_pdu_read_write_multiple_holding_registers_request_value = value;
  } else 
  if((error == false) && (function == 0x17) && (response == true)) { /* ModbusPDUReadWriteMultipleHoldingRegistersResponse */
    plc4c_list value;
    msg->modbus_pdu_read_write_multiple_holding_registers_response_value = value;
  } else 
  if((error == false) && (function == 0x16) && (response == false)) { /* ModbusPDUMaskWriteHoldingRegisterRequest */
    uint16_t referenceAddress = -1;
    msg->modbus_pdu_mask_write_holding_register_request_reference_address = referenceAddress;

    uint16_t andMask = -1;
    msg->modbus_pdu_mask_write_holding_register_request_and_mask = andMask;

    uint16_t orMask = -1;
    msg->modbus_pdu_mask_write_holding_register_request_or_mask = orMask;
  } else 
  if((error == false) && (function == 0x16) && (response == true)) { /* ModbusPDUMaskWriteHoldingRegisterResponse */
    uint16_t referenceAddress = -1;
    msg->modbus_pdu_mask_write_holding_register_response_reference_address = referenceAddress;

    uint16_t andMask = -1;
    msg->modbus_pdu_mask_write_holding_register_response_and_mask = andMask;

    uint16_t orMask = -1;
    msg->modbus_pdu_mask_write_holding_register_response_or_mask = orMask;
  } else 
  if((error == false) && (function == 0x18) && (response == false)) { /* ModbusPDUReadFifoQueueRequest */
    uint16_t fifoPointerAddress = -1;
    msg->modbus_pdu_read_fifo_queue_request_fifo_pointer_address = fifoPointerAddress;
  } else 
  if((error == false) && (function == 0x18) && (response == true)) { /* ModbusPDUReadFifoQueueResponse */
    plc4c_list fifoValue;
    msg->modbus_pdu_read_fifo_queue_response_fifo_value = fifoValue;
  } else 
  if((error == false) && (function == 0x14) && (response == false)) { /* ModbusPDUReadFileRecordRequest */
    plc4c_list* items;
    msg->modbus_pdu_read_file_record_request_items = items;
  } else 
  if((error == false) && (function == 0x14) && (response == true)) { /* ModbusPDUReadFileRecordResponse */
    plc4c_list* items;
    msg->modbus_pdu_read_file_record_response_items = items;
  } else 
  if((error == false) && (function == 0x15) && (response == false)) { /* ModbusPDUWriteFileRecordRequest */
    plc4c_list* items;
    msg->modbus_pdu_write_file_record_request_items = items;
  } else 
  if((error == false) && (function == 0x15) && (response == true)) { /* ModbusPDUWriteFileRecordResponse */
    plc4c_list* items;
    msg->modbus_pdu_write_file_record_response_items = items;
  } else 
  if((error == false) && (function == 0x07) && (response == false)) { /* ModbusPDUReadExceptionStatusRequest */
  } else 
  if((error == false) && (function == 0x07) && (response == true)) { /* ModbusPDUReadExceptionStatusResponse */
    uint8_t value = -1;
    msg->modbus_pdu_read_exception_status_response_value = value;
  } else 
  if((error == false) && (function == 0x08) && (response == false)) { /* ModbusPDUDiagnosticRequest */
    uint16_t status = -1;
    msg->modbus_pdu_diagnostic_request_status = status;

    uint16_t eventCount = -1;
    msg->modbus_pdu_diagnostic_request_event_count = eventCount;
  } else 
  if((error == false) && (function == 0x0C) && (response == false)) { /* ModbusPDUGetComEventLogRequest */
  } else 
  if((error == false) && (function == 0x0C) && (response == true)) { /* ModbusPDUGetComEventLogResponse */
    uint16_t status = -1;
    msg->modbus_pdu_get_com_event_log_response_status = status;

    uint16_t eventCount = -1;
    msg->modbus_pdu_get_com_event_log_response_event_count = eventCount;

    uint16_t messageCount = -1;
    msg->modbus_pdu_get_com_event_log_response_message_count = messageCount;

    plc4c_list events;
    msg->modbus_pdu_get_com_event_log_response_events = events;
  } else 
  if((error == false) && (function == 0x11) && (response == false)) { /* ModbusPDUReportServerIdRequest */
  } else 
  if((error == false) && (function == 0x11) && (response == true)) { /* ModbusPDUReportServerIdResponse */
    plc4c_list value;
    msg->modbus_pdu_report_server_id_response_value = value;
  } else 
  if((error == false) && (function == 0x2B) && (response == false)) { /* ModbusPDUReadDeviceIdentificationRequest */
  } else 
  if((error == false) && (function == 0x2B) && (response == true)) { /* ModbusPDUReadDeviceIdentificationResponse */
  }


  return OK;
}

plc4c_return_code plc4c_modbus_read_write_modbus_pdu_serialize(plc4c_spi_write_buffer* buf, plc4c_modbus_read_write_modbus_pdu* message) {
  return OK;
}
