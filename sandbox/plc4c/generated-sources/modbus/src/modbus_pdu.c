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
   .response = false, .function = 0x08, .error = false},
  {/* modbus_read_write_modbus_pdu_error */
   .response = -1, .function = -1, .error = true},
  {/* modbus_read_write_modbus_pdu_get_com_event_log_request */
   .response = false, .function = 0x0C, .error = false},
  {/* modbus_read_write_modbus_pdu_get_com_event_log_response */
   .response = true, .function = 0x0C, .error = false},
  {/* modbus_read_write_modbus_pdu_mask_write_holding_register_request */
   .response = false, .function = 0x16, .error = false},
  {/* modbus_read_write_modbus_pdu_mask_write_holding_register_response */
   .response = true, .function = 0x16, .error = false},
  {/* modbus_read_write_modbus_pdu_read_coils_request */
   .response = false, .function = 0x01, .error = false},
  {/* modbus_read_write_modbus_pdu_read_coils_response */
   .response = true, .function = 0x01, .error = false},
  {/* modbus_read_write_modbus_pdu_read_device_identification_request */
   .response = false, .function = 0x2B, .error = false},
  {/* modbus_read_write_modbus_pdu_read_device_identification_response */
   .response = true, .function = 0x2B, .error = false},
  {/* modbus_read_write_modbus_pdu_read_discrete_inputs_request */
   .response = false, .function = 0x02, .error = false},
  {/* modbus_read_write_modbus_pdu_read_discrete_inputs_response */
   .response = true, .function = 0x02, .error = false},
  {/* modbus_read_write_modbus_pdu_read_exception_status_request */
   .response = false, .function = 0x07, .error = false},
  {/* modbus_read_write_modbus_pdu_read_exception_status_response */
   .response = true, .function = 0x07, .error = false},
  {/* modbus_read_write_modbus_pdu_read_fifo_queue_request */
   .response = false, .function = 0x18, .error = false},
  {/* modbus_read_write_modbus_pdu_read_fifo_queue_response */
   .response = true, .function = 0x18, .error = false},
  {/* modbus_read_write_modbus_pdu_read_file_record_request */
   .response = false, .function = 0x14, .error = false},
  {/* modbus_read_write_modbus_pdu_read_file_record_response */
   .response = true, .function = 0x14, .error = false},
  {/* modbus_read_write_modbus_pdu_read_holding_registers_request */
   .response = false, .function = 0x03, .error = false},
  {/* modbus_read_write_modbus_pdu_read_holding_registers_response */
   .response = true, .function = 0x03, .error = false},
  {/* modbus_read_write_modbus_pdu_read_input_registers_request */
   .response = false, .function = 0x04, .error = false},
  {/* modbus_read_write_modbus_pdu_read_input_registers_response */
   .response = true, .function = 0x04, .error = false},
  {/* modbus_read_write_modbus_pdu_read_write_multiple_holding_registers_request */
   .response = false, .function = 0x17, .error = false},
  {/* modbus_read_write_modbus_pdu_read_write_multiple_holding_registers_response */
   .response = true, .function = 0x17, .error = false},
  {/* modbus_read_write_modbus_pdu_report_server_id_request */
   .response = false, .function = 0x11, .error = false},
  {/* modbus_read_write_modbus_pdu_report_server_id_response */
   .response = true, .function = 0x11, .error = false},
  {/* modbus_read_write_modbus_pdu_write_file_record_request */
   .response = false, .function = 0x15, .error = false},
  {/* modbus_read_write_modbus_pdu_write_file_record_response */
   .response = true, .function = 0x15, .error = false},
  {/* modbus_read_write_modbus_pdu_write_multiple_coils_request */
   .response = false, .function = 0x0F, .error = false},
  {/* modbus_read_write_modbus_pdu_write_multiple_coils_response */
   .response = true, .function = 0x0F, .error = false},
  {/* modbus_read_write_modbus_pdu_write_multiple_holding_registers_request */
   .response = false, .function = 0x10, .error = false},
  {/* modbus_read_write_modbus_pdu_write_multiple_holding_registers_response */
   .response = true, .function = 0x10, .error = false},
  {/* modbus_read_write_modbus_pdu_write_single_coil_request */
   .response = false, .function = 0x05, .error = false},
  {/* modbus_read_write_modbus_pdu_write_single_coil_response */
   .response = true, .function = 0x05, .error = false},
  {/* modbus_read_write_modbus_pdu_write_single_register_request */
   .response = false, .function = 0x06, .error = false},
  {/* modbus_read_write_modbus_pdu_write_single_register_response */
   .response = true, .function = 0x06, .error = false}
};

// Function returning the discriminator values for a given type constant.
plc4c_modbus_read_write_modbus_pdu_discriminator plc4c_modbus_read_write_modbus_pdu_get_discriminator(plc4c_modbus_read_write_modbus_pdu_type type) {
  return plc4c_modbus_read_write_modbus_pdu_discriminators[type];
}

// Parse function.
plc4c_return_code plc4c_modbus_read_write_modbus_pdu_parse(plc4c_spi_read_buffer* buf, bool response, plc4c_modbus_read_write_modbus_pdu** _message) {
  uint16_t startPos = plc4c_spi_read_get_pos(buf);
  uint16_t curPos;

  // Pointer to the parsed data structure.
  (*_message) = malloc(sizeof(plc4c_modbus_read_write_modbus_pdu));

  // Discriminator Field (error) (Used as input to a switch field)
  bool error = plc4c_spi_read_bit(buf);

  // Discriminator Field (function) (Used as input to a switch field)
  unsigned int function = plc4c_spi_read_unsigned_short(buf, 7);

  // Switch Field (Depending on the discriminator values, passes the instantiation to a sub-type)
  if(error == true) { /* ModbusPDUError */

  // Simple Field (exceptionCode)
  uint8_t exceptionCode = plc4c_spi_read_unsigned_short(buf, 8);
  (*_message)->modbus_pdu_error_exception_code = exceptionCode;
  } else 
  if((error == false) && (function == 0x02) && (response == false)) { /* ModbusPDUReadDiscreteInputsRequest */

  // Simple Field (startingAddress)
  uint16_t startingAddress = plc4c_spi_read_unsigned_int(buf, 16);
  (*_message)->modbus_pdu_read_discrete_inputs_request_starting_address = startingAddress;


  // Simple Field (quantity)
  uint16_t quantity = plc4c_spi_read_unsigned_int(buf, 16);
  (*_message)->modbus_pdu_read_discrete_inputs_request_quantity = quantity;
  } else 
  if((error == false) && (function == 0x02) && (response == true)) { /* ModbusPDUReadDiscreteInputsResponse */

  // Implicit Field (byteCount) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
  uint8_t byteCount = plc4c_spi_read_unsigned_short(buf, 8);


  // Array field (value)
  plc4c_list value;
  {
    // Count array
    uint8_t itemCount = byteCount;
    for(int curItem = 0; curItem < itemCount; curItem++) {
      
      int8_t value = plc4c_spi_read_byte(buf, 8);
      plc4c_utils_list_insert_head_value(&value, &value);
      plc4c_utils_list_insert_head_value(&value, &value);
    }
  }
  (*_message)->modbus_pdu_read_discrete_inputs_response_value = value;
  } else 
  if((error == false) && (function == 0x01) && (response == false)) { /* ModbusPDUReadCoilsRequest */

  // Simple Field (startingAddress)
  uint16_t startingAddress = plc4c_spi_read_unsigned_int(buf, 16);
  (*_message)->modbus_pdu_read_coils_request_starting_address = startingAddress;


  // Simple Field (quantity)
  uint16_t quantity = plc4c_spi_read_unsigned_int(buf, 16);
  (*_message)->modbus_pdu_read_coils_request_quantity = quantity;
  } else 
  if((error == false) && (function == 0x01) && (response == true)) { /* ModbusPDUReadCoilsResponse */

  // Implicit Field (byteCount) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
  uint8_t byteCount = plc4c_spi_read_unsigned_short(buf, 8);


  // Array field (value)
  plc4c_list value;
  {
    // Count array
    uint8_t itemCount = byteCount;
    for(int curItem = 0; curItem < itemCount; curItem++) {
      
      int8_t value = plc4c_spi_read_byte(buf, 8);
      plc4c_utils_list_insert_head_value(&value, &value);
      plc4c_utils_list_insert_head_value(&value, &value);
    }
  }
  (*_message)->modbus_pdu_read_coils_response_value = value;
  } else 
  if((error == false) && (function == 0x05) && (response == false)) { /* ModbusPDUWriteSingleCoilRequest */

  // Simple Field (address)
  uint16_t address = plc4c_spi_read_unsigned_int(buf, 16);
  (*_message)->modbus_pdu_write_single_coil_request_address = address;


  // Simple Field (value)
  uint16_t value = plc4c_spi_read_unsigned_int(buf, 16);
  (*_message)->modbus_pdu_write_single_coil_request_value = value;
  } else 
  if((error == false) && (function == 0x05) && (response == true)) { /* ModbusPDUWriteSingleCoilResponse */

  // Simple Field (address)
  uint16_t address = plc4c_spi_read_unsigned_int(buf, 16);
  (*_message)->modbus_pdu_write_single_coil_response_address = address;


  // Simple Field (value)
  uint16_t value = plc4c_spi_read_unsigned_int(buf, 16);
  (*_message)->modbus_pdu_write_single_coil_response_value = value;
  } else 
  if((error == false) && (function == 0x0F) && (response == false)) { /* ModbusPDUWriteMultipleCoilsRequest */

  // Simple Field (startingAddress)
  uint16_t startingAddress = plc4c_spi_read_unsigned_int(buf, 16);
  (*_message)->modbus_pdu_write_multiple_coils_request_starting_address = startingAddress;


  // Simple Field (quantity)
  uint16_t quantity = plc4c_spi_read_unsigned_int(buf, 16);
  (*_message)->modbus_pdu_write_multiple_coils_request_quantity = quantity;


  // Implicit Field (byteCount) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
  uint8_t byteCount = plc4c_spi_read_unsigned_short(buf, 8);


  // Array field (value)
  plc4c_list value;
  {
    // Count array
    uint8_t itemCount = byteCount;
    for(int curItem = 0; curItem < itemCount; curItem++) {
      
      int8_t value = plc4c_spi_read_byte(buf, 8);
      plc4c_utils_list_insert_head_value(&value, &value);
      plc4c_utils_list_insert_head_value(&value, &value);
    }
  }
  (*_message)->modbus_pdu_write_multiple_coils_request_value = value;
  } else 
  if((error == false) && (function == 0x0F) && (response == true)) { /* ModbusPDUWriteMultipleCoilsResponse */

  // Simple Field (startingAddress)
  uint16_t startingAddress = plc4c_spi_read_unsigned_int(buf, 16);
  (*_message)->modbus_pdu_write_multiple_coils_response_starting_address = startingAddress;


  // Simple Field (quantity)
  uint16_t quantity = plc4c_spi_read_unsigned_int(buf, 16);
  (*_message)->modbus_pdu_write_multiple_coils_response_quantity = quantity;
  } else 
  if((error == false) && (function == 0x04) && (response == false)) { /* ModbusPDUReadInputRegistersRequest */

  // Simple Field (startingAddress)
  uint16_t startingAddress = plc4c_spi_read_unsigned_int(buf, 16);
  (*_message)->modbus_pdu_read_input_registers_request_starting_address = startingAddress;


  // Simple Field (quantity)
  uint16_t quantity = plc4c_spi_read_unsigned_int(buf, 16);
  (*_message)->modbus_pdu_read_input_registers_request_quantity = quantity;
  } else 
  if((error == false) && (function == 0x04) && (response == true)) { /* ModbusPDUReadInputRegistersResponse */

  // Implicit Field (byteCount) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
  uint8_t byteCount = plc4c_spi_read_unsigned_short(buf, 8);


  // Array field (value)
  plc4c_list value;
  {
    // Count array
    uint8_t itemCount = byteCount;
    for(int curItem = 0; curItem < itemCount; curItem++) {
      
      int8_t value = plc4c_spi_read_byte(buf, 8);
      plc4c_utils_list_insert_head_value(&value, &value);
      plc4c_utils_list_insert_head_value(&value, &value);
    }
  }
  (*_message)->modbus_pdu_read_input_registers_response_value = value;
  } else 
  if((error == false) && (function == 0x03) && (response == false)) { /* ModbusPDUReadHoldingRegistersRequest */

  // Simple Field (startingAddress)
  uint16_t startingAddress = plc4c_spi_read_unsigned_int(buf, 16);
  (*_message)->modbus_pdu_read_holding_registers_request_starting_address = startingAddress;


  // Simple Field (quantity)
  uint16_t quantity = plc4c_spi_read_unsigned_int(buf, 16);
  (*_message)->modbus_pdu_read_holding_registers_request_quantity = quantity;
  } else 
  if((error == false) && (function == 0x03) && (response == true)) { /* ModbusPDUReadHoldingRegistersResponse */

  // Implicit Field (byteCount) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
  uint8_t byteCount = plc4c_spi_read_unsigned_short(buf, 8);


  // Array field (value)
  plc4c_list value;
  {
    // Count array
    uint8_t itemCount = byteCount;
    for(int curItem = 0; curItem < itemCount; curItem++) {
      
      int8_t value = plc4c_spi_read_byte(buf, 8);
      plc4c_utils_list_insert_head_value(&value, &value);
      plc4c_utils_list_insert_head_value(&value, &value);
    }
  }
  (*_message)->modbus_pdu_read_holding_registers_response_value = value;
  } else 
  if((error == false) && (function == 0x06) && (response == false)) { /* ModbusPDUWriteSingleRegisterRequest */

  // Simple Field (address)
  uint16_t address = plc4c_spi_read_unsigned_int(buf, 16);
  (*_message)->modbus_pdu_write_single_register_request_address = address;


  // Simple Field (value)
  uint16_t value = plc4c_spi_read_unsigned_int(buf, 16);
  (*_message)->modbus_pdu_write_single_register_request_value = value;
  } else 
  if((error == false) && (function == 0x06) && (response == true)) { /* ModbusPDUWriteSingleRegisterResponse */

  // Simple Field (address)
  uint16_t address = plc4c_spi_read_unsigned_int(buf, 16);
  (*_message)->modbus_pdu_write_single_register_response_address = address;


  // Simple Field (value)
  uint16_t value = plc4c_spi_read_unsigned_int(buf, 16);
  (*_message)->modbus_pdu_write_single_register_response_value = value;
  } else 
  if((error == false) && (function == 0x10) && (response == false)) { /* ModbusPDUWriteMultipleHoldingRegistersRequest */

  // Simple Field (startingAddress)
  uint16_t startingAddress = plc4c_spi_read_unsigned_int(buf, 16);
  (*_message)->modbus_pdu_write_multiple_holding_registers_request_starting_address = startingAddress;


  // Simple Field (quantity)
  uint16_t quantity = plc4c_spi_read_unsigned_int(buf, 16);
  (*_message)->modbus_pdu_write_multiple_holding_registers_request_quantity = quantity;


  // Implicit Field (byteCount) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
  uint8_t byteCount = plc4c_spi_read_unsigned_short(buf, 8);


  // Array field (value)
  plc4c_list value;
  {
    // Count array
    uint8_t itemCount = byteCount;
    for(int curItem = 0; curItem < itemCount; curItem++) {
      
      int8_t value = plc4c_spi_read_byte(buf, 8);
      plc4c_utils_list_insert_head_value(&value, &value);
      plc4c_utils_list_insert_head_value(&value, &value);
    }
  }
  (*_message)->modbus_pdu_write_multiple_holding_registers_request_value = value;
  } else 
  if((error == false) && (function == 0x10) && (response == true)) { /* ModbusPDUWriteMultipleHoldingRegistersResponse */

  // Simple Field (startingAddress)
  uint16_t startingAddress = plc4c_spi_read_unsigned_int(buf, 16);
  (*_message)->modbus_pdu_write_multiple_holding_registers_response_starting_address = startingAddress;


  // Simple Field (quantity)
  uint16_t quantity = plc4c_spi_read_unsigned_int(buf, 16);
  (*_message)->modbus_pdu_write_multiple_holding_registers_response_quantity = quantity;
  } else 
  if((error == false) && (function == 0x17) && (response == false)) { /* ModbusPDUReadWriteMultipleHoldingRegistersRequest */

  // Simple Field (readStartingAddress)
  uint16_t readStartingAddress = plc4c_spi_read_unsigned_int(buf, 16);
  (*_message)->modbus_pdu_read_write_multiple_holding_registers_request_read_starting_address = readStartingAddress;


  // Simple Field (readQuantity)
  uint16_t readQuantity = plc4c_spi_read_unsigned_int(buf, 16);
  (*_message)->modbus_pdu_read_write_multiple_holding_registers_request_read_quantity = readQuantity;


  // Simple Field (writeStartingAddress)
  uint16_t writeStartingAddress = plc4c_spi_read_unsigned_int(buf, 16);
  (*_message)->modbus_pdu_read_write_multiple_holding_registers_request_write_starting_address = writeStartingAddress;


  // Simple Field (writeQuantity)
  uint16_t writeQuantity = plc4c_spi_read_unsigned_int(buf, 16);
  (*_message)->modbus_pdu_read_write_multiple_holding_registers_request_write_quantity = writeQuantity;


  // Implicit Field (byteCount) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
  uint8_t byteCount = plc4c_spi_read_unsigned_short(buf, 8);


  // Array field (value)
  plc4c_list value;
  {
    // Count array
    uint8_t itemCount = byteCount;
    for(int curItem = 0; curItem < itemCount; curItem++) {
      
      int8_t value = plc4c_spi_read_byte(buf, 8);
      plc4c_utils_list_insert_head_value(&value, &value);
      plc4c_utils_list_insert_head_value(&value, &value);
    }
  }
  (*_message)->modbus_pdu_read_write_multiple_holding_registers_request_value = value;
  } else 
  if((error == false) && (function == 0x17) && (response == true)) { /* ModbusPDUReadWriteMultipleHoldingRegistersResponse */

  // Implicit Field (byteCount) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
  uint8_t byteCount = plc4c_spi_read_unsigned_short(buf, 8);


  // Array field (value)
  plc4c_list value;
  {
    // Count array
    uint8_t itemCount = byteCount;
    for(int curItem = 0; curItem < itemCount; curItem++) {
      
      int8_t value = plc4c_spi_read_byte(buf, 8);
      plc4c_utils_list_insert_head_value(&value, &value);
      plc4c_utils_list_insert_head_value(&value, &value);
    }
  }
  (*_message)->modbus_pdu_read_write_multiple_holding_registers_response_value = value;
  } else 
  if((error == false) && (function == 0x16) && (response == false)) { /* ModbusPDUMaskWriteHoldingRegisterRequest */

  // Simple Field (referenceAddress)
  uint16_t referenceAddress = plc4c_spi_read_unsigned_int(buf, 16);
  (*_message)->modbus_pdu_mask_write_holding_register_request_reference_address = referenceAddress;


  // Simple Field (andMask)
  uint16_t andMask = plc4c_spi_read_unsigned_int(buf, 16);
  (*_message)->modbus_pdu_mask_write_holding_register_request_and_mask = andMask;


  // Simple Field (orMask)
  uint16_t orMask = plc4c_spi_read_unsigned_int(buf, 16);
  (*_message)->modbus_pdu_mask_write_holding_register_request_or_mask = orMask;
  } else 
  if((error == false) && (function == 0x16) && (response == true)) { /* ModbusPDUMaskWriteHoldingRegisterResponse */

  // Simple Field (referenceAddress)
  uint16_t referenceAddress = plc4c_spi_read_unsigned_int(buf, 16);
  (*_message)->modbus_pdu_mask_write_holding_register_response_reference_address = referenceAddress;


  // Simple Field (andMask)
  uint16_t andMask = plc4c_spi_read_unsigned_int(buf, 16);
  (*_message)->modbus_pdu_mask_write_holding_register_response_and_mask = andMask;


  // Simple Field (orMask)
  uint16_t orMask = plc4c_spi_read_unsigned_int(buf, 16);
  (*_message)->modbus_pdu_mask_write_holding_register_response_or_mask = orMask;
  } else 
  if((error == false) && (function == 0x18) && (response == false)) { /* ModbusPDUReadFifoQueueRequest */

  // Simple Field (fifoPointerAddress)
  uint16_t fifoPointerAddress = plc4c_spi_read_unsigned_int(buf, 16);
  (*_message)->modbus_pdu_read_fifo_queue_request_fifo_pointer_address = fifoPointerAddress;
  } else 
  if((error == false) && (function == 0x18) && (response == true)) { /* ModbusPDUReadFifoQueueResponse */

  // Implicit Field (byteCount) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
  uint16_t byteCount = plc4c_spi_read_unsigned_int(buf, 16);


  // Implicit Field (fifoCount) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
  uint16_t fifoCount = plc4c_spi_read_unsigned_int(buf, 16);


  // Array field (fifoValue)
  plc4c_list fifoValue;
  {
    // Count array
    uint8_t itemCount = fifoCount;
    for(int curItem = 0; curItem < itemCount; curItem++) {
      
      uint16_t value = plc4c_spi_read_unsigned_int(buf, 16);
      plc4c_utils_list_insert_head_value(&fifoValue, &value);
      plc4c_utils_list_insert_head_value(&fifoValue, &value);
    }
  }
  (*_message)->modbus_pdu_read_fifo_queue_response_fifo_value = fifoValue;
  } else 
  if((error == false) && (function == 0x14) && (response == false)) { /* ModbusPDUReadFileRecordRequest */

  // Implicit Field (byteCount) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
  uint8_t byteCount = plc4c_spi_read_unsigned_short(buf, 8);


  // Array field (items)
  plc4c_list items;
  {
    // Length array
    uint8_t _itemsLength = byteCount;
    uint8_t itemsEndPos = plc4c_spi_read_get_pos(buf) + _itemsLength;
    while(plc4c_spi_read_get_pos(buf) < itemsEndPos) {
      plc4c_list* value = NULL;
      plc4c_modbus_read_write_modbus_pdu_read_file_record_request_item_parse(buf, (void*) &value);
      plc4c_utils_list_insert_head_value(&items, value);
    }
  }
  (*_message)->modbus_pdu_read_file_record_request_items = items;
  } else 
  if((error == false) && (function == 0x14) && (response == true)) { /* ModbusPDUReadFileRecordResponse */

  // Implicit Field (byteCount) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
  uint8_t byteCount = plc4c_spi_read_unsigned_short(buf, 8);


  // Array field (items)
  plc4c_list items;
  {
    // Length array
    uint8_t _itemsLength = byteCount;
    uint8_t itemsEndPos = plc4c_spi_read_get_pos(buf) + _itemsLength;
    while(plc4c_spi_read_get_pos(buf) < itemsEndPos) {
      plc4c_list* value = NULL;
      plc4c_modbus_read_write_modbus_pdu_read_file_record_response_item_parse(buf, (void*) &value);
      plc4c_utils_list_insert_head_value(&items, value);
    }
  }
  (*_message)->modbus_pdu_read_file_record_response_items = items;
  } else 
  if((error == false) && (function == 0x15) && (response == false)) { /* ModbusPDUWriteFileRecordRequest */

  // Implicit Field (byteCount) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
  uint8_t byteCount = plc4c_spi_read_unsigned_short(buf, 8);


  // Array field (items)
  plc4c_list items;
  {
    // Length array
    uint8_t _itemsLength = byteCount;
    uint8_t itemsEndPos = plc4c_spi_read_get_pos(buf) + _itemsLength;
    while(plc4c_spi_read_get_pos(buf) < itemsEndPos) {
      plc4c_list* value = NULL;
      plc4c_modbus_read_write_modbus_pdu_write_file_record_request_item_parse(buf, (void*) &value);
      plc4c_utils_list_insert_head_value(&items, value);
    }
  }
  (*_message)->modbus_pdu_write_file_record_request_items = items;
  } else 
  if((error == false) && (function == 0x15) && (response == true)) { /* ModbusPDUWriteFileRecordResponse */

  // Implicit Field (byteCount) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
  uint8_t byteCount = plc4c_spi_read_unsigned_short(buf, 8);


  // Array field (items)
  plc4c_list items;
  {
    // Length array
    uint8_t _itemsLength = byteCount;
    uint8_t itemsEndPos = plc4c_spi_read_get_pos(buf) + _itemsLength;
    while(plc4c_spi_read_get_pos(buf) < itemsEndPos) {
      plc4c_list* value = NULL;
      plc4c_modbus_read_write_modbus_pdu_write_file_record_response_item_parse(buf, (void*) &value);
      plc4c_utils_list_insert_head_value(&items, value);
    }
  }
  (*_message)->modbus_pdu_write_file_record_response_items = items;
  } else 
  if((error == false) && (function == 0x07) && (response == false)) { /* ModbusPDUReadExceptionStatusRequest */
  } else 
  if((error == false) && (function == 0x07) && (response == true)) { /* ModbusPDUReadExceptionStatusResponse */

  // Simple Field (value)
  uint8_t value = plc4c_spi_read_unsigned_short(buf, 8);
  (*_message)->modbus_pdu_read_exception_status_response_value = value;
  } else 
  if((error == false) && (function == 0x08) && (response == false)) { /* ModbusPDUDiagnosticRequest */

  // Simple Field (status)
  uint16_t status = plc4c_spi_read_unsigned_int(buf, 16);
  (*_message)->modbus_pdu_diagnostic_request_status = status;


  // Simple Field (eventCount)
  uint16_t eventCount = plc4c_spi_read_unsigned_int(buf, 16);
  (*_message)->modbus_pdu_diagnostic_request_event_count = eventCount;
  } else 
  if((error == false) && (function == 0x0C) && (response == false)) { /* ModbusPDUGetComEventLogRequest */
  } else 
  if((error == false) && (function == 0x0C) && (response == true)) { /* ModbusPDUGetComEventLogResponse */

  // Implicit Field (byteCount) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
  uint8_t byteCount = plc4c_spi_read_unsigned_short(buf, 8);


  // Simple Field (status)
  uint16_t status = plc4c_spi_read_unsigned_int(buf, 16);
  (*_message)->modbus_pdu_get_com_event_log_response_status = status;


  // Simple Field (eventCount)
  uint16_t eventCount = plc4c_spi_read_unsigned_int(buf, 16);
  (*_message)->modbus_pdu_get_com_event_log_response_event_count = eventCount;


  // Simple Field (messageCount)
  uint16_t messageCount = plc4c_spi_read_unsigned_int(buf, 16);
  (*_message)->modbus_pdu_get_com_event_log_response_message_count = messageCount;


  // Array field (events)
  plc4c_list events;
  {
    // Count array
    uint8_t itemCount = (byteCount) - (6);
    for(int curItem = 0; curItem < itemCount; curItem++) {
      
      int8_t value = plc4c_spi_read_byte(buf, 8);
      plc4c_utils_list_insert_head_value(&events, &value);
      plc4c_utils_list_insert_head_value(&events, &value);
    }
  }
  (*_message)->modbus_pdu_get_com_event_log_response_events = events;
  } else 
  if((error == false) && (function == 0x11) && (response == false)) { /* ModbusPDUReportServerIdRequest */
  } else 
  if((error == false) && (function == 0x11) && (response == true)) { /* ModbusPDUReportServerIdResponse */

  // Implicit Field (byteCount) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
  uint8_t byteCount = plc4c_spi_read_unsigned_short(buf, 8);


  // Array field (value)
  plc4c_list value;
  {
    // Count array
    uint8_t itemCount = byteCount;
    for(int curItem = 0; curItem < itemCount; curItem++) {
      
      int8_t value = plc4c_spi_read_byte(buf, 8);
      plc4c_utils_list_insert_head_value(&value, &value);
      plc4c_utils_list_insert_head_value(&value, &value);
    }
  }
  (*_message)->modbus_pdu_report_server_id_response_value = value;
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
