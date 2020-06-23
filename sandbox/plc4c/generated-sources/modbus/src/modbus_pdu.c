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

  // Allocate enough memory to contain this data structure.
  (*_message) = malloc(sizeof(plc4c_modbus_read_write_modbus_pdu));
  if(*_message == NULL) {
    return NO_MEMORY;
  }

  // Discriminator Field (error) (Used as input to a switch field)
  bool error = plc4c_spi_read_bit(buf);

  // Discriminator Field (function) (Used as input to a switch field)
  unsigned int function = plc4c_spi_read_unsigned_short(buf, 7);

  // Switch Field (Depending on the discriminator values, passes the instantiation to a sub-type)
  if(error == true) { /* ModbusPDUError */
    (*_message)->_type = plc4c_modbus_read_write_modbus_pdu_type_modbus_read_write_modbus_pdu_error;
                    
    // Simple Field (exceptionCode)
    uint8_t exceptionCode = plc4c_spi_read_unsigned_short(buf, 8);
    (*_message)->modbus_pdu_error_exception_code = exceptionCode;

  } else 
  if((error == false) && (function == 0x02) && (response == false)) { /* ModbusPDUReadDiscreteInputsRequest */
    (*_message)->_type = plc4c_modbus_read_write_modbus_pdu_type_modbus_read_write_modbus_pdu_read_discrete_inputs_request;
                    
    // Simple Field (startingAddress)
    uint16_t startingAddress = plc4c_spi_read_unsigned_int(buf, 16);
    (*_message)->modbus_pdu_read_discrete_inputs_request_starting_address = startingAddress;


                    
    // Simple Field (quantity)
    uint16_t quantity = plc4c_spi_read_unsigned_int(buf, 16);
    (*_message)->modbus_pdu_read_discrete_inputs_request_quantity = quantity;

  } else 
  if((error == false) && (function == 0x02) && (response == true)) { /* ModbusPDUReadDiscreteInputsResponse */
    (*_message)->_type = plc4c_modbus_read_write_modbus_pdu_type_modbus_read_write_modbus_pdu_read_discrete_inputs_response;
                    
    // Implicit Field (byteCount) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
    uint8_t byteCount = plc4c_spi_read_unsigned_short(buf, 8);


                    
    // Array field (value)
    plc4c_list* value = malloc(sizeof(plc4c_list));
    if(value == NULL) {
      return NO_MEMORY;
    }
    {
      // Count array
      uint8_t itemCount = byteCount;
      for(int curItem = 0; curItem < itemCount; curItem++) {
        
                  
        int8_t _value = plc4c_spi_read_byte(buf, 8);
        plc4c_utils_list_insert_head_value(value, &_value);
      }
    }
    (*_message)->modbus_pdu_read_discrete_inputs_response_value = value;

  } else 
  if((error == false) && (function == 0x01) && (response == false)) { /* ModbusPDUReadCoilsRequest */
    (*_message)->_type = plc4c_modbus_read_write_modbus_pdu_type_modbus_read_write_modbus_pdu_read_coils_request;
                    
    // Simple Field (startingAddress)
    uint16_t startingAddress = plc4c_spi_read_unsigned_int(buf, 16);
    (*_message)->modbus_pdu_read_coils_request_starting_address = startingAddress;


                    
    // Simple Field (quantity)
    uint16_t quantity = plc4c_spi_read_unsigned_int(buf, 16);
    (*_message)->modbus_pdu_read_coils_request_quantity = quantity;

  } else 
  if((error == false) && (function == 0x01) && (response == true)) { /* ModbusPDUReadCoilsResponse */
    (*_message)->_type = plc4c_modbus_read_write_modbus_pdu_type_modbus_read_write_modbus_pdu_read_coils_response;
                    
    // Implicit Field (byteCount) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
    uint8_t byteCount = plc4c_spi_read_unsigned_short(buf, 8);


                    
    // Array field (value)
    plc4c_list* value = malloc(sizeof(plc4c_list));
    if(value == NULL) {
      return NO_MEMORY;
    }
    {
      // Count array
      uint8_t itemCount = byteCount;
      for(int curItem = 0; curItem < itemCount; curItem++) {
        
                  
        int8_t _value = plc4c_spi_read_byte(buf, 8);
        plc4c_utils_list_insert_head_value(value, &_value);
      }
    }
    (*_message)->modbus_pdu_read_coils_response_value = value;

  } else 
  if((error == false) && (function == 0x05) && (response == false)) { /* ModbusPDUWriteSingleCoilRequest */
    (*_message)->_type = plc4c_modbus_read_write_modbus_pdu_type_modbus_read_write_modbus_pdu_write_single_coil_request;
                    
    // Simple Field (address)
    uint16_t address = plc4c_spi_read_unsigned_int(buf, 16);
    (*_message)->modbus_pdu_write_single_coil_request_address = address;


                    
    // Simple Field (value)
    uint16_t value = plc4c_spi_read_unsigned_int(buf, 16);
    (*_message)->modbus_pdu_write_single_coil_request_value = value;

  } else 
  if((error == false) && (function == 0x05) && (response == true)) { /* ModbusPDUWriteSingleCoilResponse */
    (*_message)->_type = plc4c_modbus_read_write_modbus_pdu_type_modbus_read_write_modbus_pdu_write_single_coil_response;
                    
    // Simple Field (address)
    uint16_t address = plc4c_spi_read_unsigned_int(buf, 16);
    (*_message)->modbus_pdu_write_single_coil_response_address = address;


                    
    // Simple Field (value)
    uint16_t value = plc4c_spi_read_unsigned_int(buf, 16);
    (*_message)->modbus_pdu_write_single_coil_response_value = value;

  } else 
  if((error == false) && (function == 0x0F) && (response == false)) { /* ModbusPDUWriteMultipleCoilsRequest */
    (*_message)->_type = plc4c_modbus_read_write_modbus_pdu_type_modbus_read_write_modbus_pdu_write_multiple_coils_request;
                    
    // Simple Field (startingAddress)
    uint16_t startingAddress = plc4c_spi_read_unsigned_int(buf, 16);
    (*_message)->modbus_pdu_write_multiple_coils_request_starting_address = startingAddress;


                    
    // Simple Field (quantity)
    uint16_t quantity = plc4c_spi_read_unsigned_int(buf, 16);
    (*_message)->modbus_pdu_write_multiple_coils_request_quantity = quantity;


                    
    // Implicit Field (byteCount) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
    uint8_t byteCount = plc4c_spi_read_unsigned_short(buf, 8);


                    
    // Array field (value)
    plc4c_list* value = malloc(sizeof(plc4c_list));
    if(value == NULL) {
      return NO_MEMORY;
    }
    {
      // Count array
      uint8_t itemCount = byteCount;
      for(int curItem = 0; curItem < itemCount; curItem++) {
        
                  
        int8_t _value = plc4c_spi_read_byte(buf, 8);
        plc4c_utils_list_insert_head_value(value, &_value);
      }
    }
    (*_message)->modbus_pdu_write_multiple_coils_request_value = value;

  } else 
  if((error == false) && (function == 0x0F) && (response == true)) { /* ModbusPDUWriteMultipleCoilsResponse */
    (*_message)->_type = plc4c_modbus_read_write_modbus_pdu_type_modbus_read_write_modbus_pdu_write_multiple_coils_response;
                    
    // Simple Field (startingAddress)
    uint16_t startingAddress = plc4c_spi_read_unsigned_int(buf, 16);
    (*_message)->modbus_pdu_write_multiple_coils_response_starting_address = startingAddress;


                    
    // Simple Field (quantity)
    uint16_t quantity = plc4c_spi_read_unsigned_int(buf, 16);
    (*_message)->modbus_pdu_write_multiple_coils_response_quantity = quantity;

  } else 
  if((error == false) && (function == 0x04) && (response == false)) { /* ModbusPDUReadInputRegistersRequest */
    (*_message)->_type = plc4c_modbus_read_write_modbus_pdu_type_modbus_read_write_modbus_pdu_read_input_registers_request;
                    
    // Simple Field (startingAddress)
    uint16_t startingAddress = plc4c_spi_read_unsigned_int(buf, 16);
    (*_message)->modbus_pdu_read_input_registers_request_starting_address = startingAddress;


                    
    // Simple Field (quantity)
    uint16_t quantity = plc4c_spi_read_unsigned_int(buf, 16);
    (*_message)->modbus_pdu_read_input_registers_request_quantity = quantity;

  } else 
  if((error == false) && (function == 0x04) && (response == true)) { /* ModbusPDUReadInputRegistersResponse */
    (*_message)->_type = plc4c_modbus_read_write_modbus_pdu_type_modbus_read_write_modbus_pdu_read_input_registers_response;
                    
    // Implicit Field (byteCount) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
    uint8_t byteCount = plc4c_spi_read_unsigned_short(buf, 8);


                    
    // Array field (value)
    plc4c_list* value = malloc(sizeof(plc4c_list));
    if(value == NULL) {
      return NO_MEMORY;
    }
    {
      // Count array
      uint8_t itemCount = byteCount;
      for(int curItem = 0; curItem < itemCount; curItem++) {
        
                  
        int8_t _value = plc4c_spi_read_byte(buf, 8);
        plc4c_utils_list_insert_head_value(value, &_value);
      }
    }
    (*_message)->modbus_pdu_read_input_registers_response_value = value;

  } else 
  if((error == false) && (function == 0x03) && (response == false)) { /* ModbusPDUReadHoldingRegistersRequest */
    (*_message)->_type = plc4c_modbus_read_write_modbus_pdu_type_modbus_read_write_modbus_pdu_read_holding_registers_request;
                    
    // Simple Field (startingAddress)
    uint16_t startingAddress = plc4c_spi_read_unsigned_int(buf, 16);
    (*_message)->modbus_pdu_read_holding_registers_request_starting_address = startingAddress;


                    
    // Simple Field (quantity)
    uint16_t quantity = plc4c_spi_read_unsigned_int(buf, 16);
    (*_message)->modbus_pdu_read_holding_registers_request_quantity = quantity;

  } else 
  if((error == false) && (function == 0x03) && (response == true)) { /* ModbusPDUReadHoldingRegistersResponse */
    (*_message)->_type = plc4c_modbus_read_write_modbus_pdu_type_modbus_read_write_modbus_pdu_read_holding_registers_response;
                    
    // Implicit Field (byteCount) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
    uint8_t byteCount = plc4c_spi_read_unsigned_short(buf, 8);


                    
    // Array field (value)
    plc4c_list* value = malloc(sizeof(plc4c_list));
    if(value == NULL) {
      return NO_MEMORY;
    }
    {
      // Count array
      uint8_t itemCount = byteCount;
      for(int curItem = 0; curItem < itemCount; curItem++) {
        
                  
        int8_t _value = plc4c_spi_read_byte(buf, 8);
        plc4c_utils_list_insert_head_value(value, &_value);
      }
    }
    (*_message)->modbus_pdu_read_holding_registers_response_value = value;

  } else 
  if((error == false) && (function == 0x06) && (response == false)) { /* ModbusPDUWriteSingleRegisterRequest */
    (*_message)->_type = plc4c_modbus_read_write_modbus_pdu_type_modbus_read_write_modbus_pdu_write_single_register_request;
                    
    // Simple Field (address)
    uint16_t address = plc4c_spi_read_unsigned_int(buf, 16);
    (*_message)->modbus_pdu_write_single_register_request_address = address;


                    
    // Simple Field (value)
    uint16_t value = plc4c_spi_read_unsigned_int(buf, 16);
    (*_message)->modbus_pdu_write_single_register_request_value = value;

  } else 
  if((error == false) && (function == 0x06) && (response == true)) { /* ModbusPDUWriteSingleRegisterResponse */
    (*_message)->_type = plc4c_modbus_read_write_modbus_pdu_type_modbus_read_write_modbus_pdu_write_single_register_response;
                    
    // Simple Field (address)
    uint16_t address = plc4c_spi_read_unsigned_int(buf, 16);
    (*_message)->modbus_pdu_write_single_register_response_address = address;


                    
    // Simple Field (value)
    uint16_t value = plc4c_spi_read_unsigned_int(buf, 16);
    (*_message)->modbus_pdu_write_single_register_response_value = value;

  } else 
  if((error == false) && (function == 0x10) && (response == false)) { /* ModbusPDUWriteMultipleHoldingRegistersRequest */
    (*_message)->_type = plc4c_modbus_read_write_modbus_pdu_type_modbus_read_write_modbus_pdu_write_multiple_holding_registers_request;
                    
    // Simple Field (startingAddress)
    uint16_t startingAddress = plc4c_spi_read_unsigned_int(buf, 16);
    (*_message)->modbus_pdu_write_multiple_holding_registers_request_starting_address = startingAddress;


                    
    // Simple Field (quantity)
    uint16_t quantity = plc4c_spi_read_unsigned_int(buf, 16);
    (*_message)->modbus_pdu_write_multiple_holding_registers_request_quantity = quantity;


                    
    // Implicit Field (byteCount) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
    uint8_t byteCount = plc4c_spi_read_unsigned_short(buf, 8);


                    
    // Array field (value)
    plc4c_list* value = malloc(sizeof(plc4c_list));
    if(value == NULL) {
      return NO_MEMORY;
    }
    {
      // Count array
      uint8_t itemCount = byteCount;
      for(int curItem = 0; curItem < itemCount; curItem++) {
        
                  
        int8_t _value = plc4c_spi_read_byte(buf, 8);
        plc4c_utils_list_insert_head_value(value, &_value);
      }
    }
    (*_message)->modbus_pdu_write_multiple_holding_registers_request_value = value;

  } else 
  if((error == false) && (function == 0x10) && (response == true)) { /* ModbusPDUWriteMultipleHoldingRegistersResponse */
    (*_message)->_type = plc4c_modbus_read_write_modbus_pdu_type_modbus_read_write_modbus_pdu_write_multiple_holding_registers_response;
                    
    // Simple Field (startingAddress)
    uint16_t startingAddress = plc4c_spi_read_unsigned_int(buf, 16);
    (*_message)->modbus_pdu_write_multiple_holding_registers_response_starting_address = startingAddress;


                    
    // Simple Field (quantity)
    uint16_t quantity = plc4c_spi_read_unsigned_int(buf, 16);
    (*_message)->modbus_pdu_write_multiple_holding_registers_response_quantity = quantity;

  } else 
  if((error == false) && (function == 0x17) && (response == false)) { /* ModbusPDUReadWriteMultipleHoldingRegistersRequest */
    (*_message)->_type = plc4c_modbus_read_write_modbus_pdu_type_modbus_read_write_modbus_pdu_read_write_multiple_holding_registers_request;
                    
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
    plc4c_list* value = malloc(sizeof(plc4c_list));
    if(value == NULL) {
      return NO_MEMORY;
    }
    {
      // Count array
      uint8_t itemCount = byteCount;
      for(int curItem = 0; curItem < itemCount; curItem++) {
        
                  
        int8_t _value = plc4c_spi_read_byte(buf, 8);
        plc4c_utils_list_insert_head_value(value, &_value);
      }
    }
    (*_message)->modbus_pdu_read_write_multiple_holding_registers_request_value = value;

  } else 
  if((error == false) && (function == 0x17) && (response == true)) { /* ModbusPDUReadWriteMultipleHoldingRegistersResponse */
    (*_message)->_type = plc4c_modbus_read_write_modbus_pdu_type_modbus_read_write_modbus_pdu_read_write_multiple_holding_registers_response;
                    
    // Implicit Field (byteCount) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
    uint8_t byteCount = plc4c_spi_read_unsigned_short(buf, 8);


                    
    // Array field (value)
    plc4c_list* value = malloc(sizeof(plc4c_list));
    if(value == NULL) {
      return NO_MEMORY;
    }
    {
      // Count array
      uint8_t itemCount = byteCount;
      for(int curItem = 0; curItem < itemCount; curItem++) {
        
                  
        int8_t _value = plc4c_spi_read_byte(buf, 8);
        plc4c_utils_list_insert_head_value(value, &_value);
      }
    }
    (*_message)->modbus_pdu_read_write_multiple_holding_registers_response_value = value;

  } else 
  if((error == false) && (function == 0x16) && (response == false)) { /* ModbusPDUMaskWriteHoldingRegisterRequest */
    (*_message)->_type = plc4c_modbus_read_write_modbus_pdu_type_modbus_read_write_modbus_pdu_mask_write_holding_register_request;
                    
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
    (*_message)->_type = plc4c_modbus_read_write_modbus_pdu_type_modbus_read_write_modbus_pdu_mask_write_holding_register_response;
                    
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
    (*_message)->_type = plc4c_modbus_read_write_modbus_pdu_type_modbus_read_write_modbus_pdu_read_fifo_queue_request;
                    
    // Simple Field (fifoPointerAddress)
    uint16_t fifoPointerAddress = plc4c_spi_read_unsigned_int(buf, 16);
    (*_message)->modbus_pdu_read_fifo_queue_request_fifo_pointer_address = fifoPointerAddress;

  } else 
  if((error == false) && (function == 0x18) && (response == true)) { /* ModbusPDUReadFifoQueueResponse */
    (*_message)->_type = plc4c_modbus_read_write_modbus_pdu_type_modbus_read_write_modbus_pdu_read_fifo_queue_response;
                    
    // Implicit Field (byteCount) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
    uint16_t byteCount = plc4c_spi_read_unsigned_int(buf, 16);


                    
    // Implicit Field (fifoCount) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
    uint16_t fifoCount = plc4c_spi_read_unsigned_int(buf, 16);


                    
    // Array field (fifoValue)
    plc4c_list* fifoValue = malloc(sizeof(plc4c_list));
    if(fifoValue == NULL) {
      return NO_MEMORY;
    }
    {
      // Count array
      uint8_t itemCount = fifoCount;
      for(int curItem = 0; curItem < itemCount; curItem++) {
        
                  
        uint16_t _value = plc4c_spi_read_unsigned_int(buf, 16);
        plc4c_utils_list_insert_head_value(fifoValue, &_value);
      }
    }
    (*_message)->modbus_pdu_read_fifo_queue_response_fifo_value = fifoValue;

  } else 
  if((error == false) && (function == 0x14) && (response == false)) { /* ModbusPDUReadFileRecordRequest */
    (*_message)->_type = plc4c_modbus_read_write_modbus_pdu_type_modbus_read_write_modbus_pdu_read_file_record_request;
                    
    // Implicit Field (byteCount) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
    uint8_t byteCount = plc4c_spi_read_unsigned_short(buf, 8);


                    
    // Array field (items)
    plc4c_list* items = malloc(sizeof(plc4c_list));
    if(items == NULL) {
      return NO_MEMORY;
    }
    {
      // Length array
      uint8_t _itemsLength = byteCount;
      uint8_t itemsEndPos = plc4c_spi_read_get_pos(buf) + _itemsLength;
      while(plc4c_spi_read_get_pos(buf) < itemsEndPos) {
        plc4c_list* _value = NULL;
        plc4c_return_code _res = plc4c_modbus_read_write_modbus_pdu_read_file_record_request_item_parse(buf, (void*) &_value);
        if(_res != OK) {
          return _res;
        }
        plc4c_utils_list_insert_head_value(items, _value);
      }
    }
    (*_message)->modbus_pdu_read_file_record_request_items = items;

  } else 
  if((error == false) && (function == 0x14) && (response == true)) { /* ModbusPDUReadFileRecordResponse */
    (*_message)->_type = plc4c_modbus_read_write_modbus_pdu_type_modbus_read_write_modbus_pdu_read_file_record_response;
                    
    // Implicit Field (byteCount) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
    uint8_t byteCount = plc4c_spi_read_unsigned_short(buf, 8);


                    
    // Array field (items)
    plc4c_list* items = malloc(sizeof(plc4c_list));
    if(items == NULL) {
      return NO_MEMORY;
    }
    {
      // Length array
      uint8_t _itemsLength = byteCount;
      uint8_t itemsEndPos = plc4c_spi_read_get_pos(buf) + _itemsLength;
      while(plc4c_spi_read_get_pos(buf) < itemsEndPos) {
        plc4c_list* _value = NULL;
        plc4c_return_code _res = plc4c_modbus_read_write_modbus_pdu_read_file_record_response_item_parse(buf, (void*) &_value);
        if(_res != OK) {
          return _res;
        }
        plc4c_utils_list_insert_head_value(items, _value);
      }
    }
    (*_message)->modbus_pdu_read_file_record_response_items = items;

  } else 
  if((error == false) && (function == 0x15) && (response == false)) { /* ModbusPDUWriteFileRecordRequest */
    (*_message)->_type = plc4c_modbus_read_write_modbus_pdu_type_modbus_read_write_modbus_pdu_write_file_record_request;
                    
    // Implicit Field (byteCount) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
    uint8_t byteCount = plc4c_spi_read_unsigned_short(buf, 8);


                    
    // Array field (items)
    plc4c_list* items = malloc(sizeof(plc4c_list));
    if(items == NULL) {
      return NO_MEMORY;
    }
    {
      // Length array
      uint8_t _itemsLength = byteCount;
      uint8_t itemsEndPos = plc4c_spi_read_get_pos(buf) + _itemsLength;
      while(plc4c_spi_read_get_pos(buf) < itemsEndPos) {
        plc4c_list* _value = NULL;
        plc4c_return_code _res = plc4c_modbus_read_write_modbus_pdu_write_file_record_request_item_parse(buf, (void*) &_value);
        if(_res != OK) {
          return _res;
        }
        plc4c_utils_list_insert_head_value(items, _value);
      }
    }
    (*_message)->modbus_pdu_write_file_record_request_items = items;

  } else 
  if((error == false) && (function == 0x15) && (response == true)) { /* ModbusPDUWriteFileRecordResponse */
    (*_message)->_type = plc4c_modbus_read_write_modbus_pdu_type_modbus_read_write_modbus_pdu_write_file_record_response;
                    
    // Implicit Field (byteCount) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
    uint8_t byteCount = plc4c_spi_read_unsigned_short(buf, 8);


                    
    // Array field (items)
    plc4c_list* items = malloc(sizeof(plc4c_list));
    if(items == NULL) {
      return NO_MEMORY;
    }
    {
      // Length array
      uint8_t _itemsLength = byteCount;
      uint8_t itemsEndPos = plc4c_spi_read_get_pos(buf) + _itemsLength;
      while(plc4c_spi_read_get_pos(buf) < itemsEndPos) {
        plc4c_list* _value = NULL;
        plc4c_return_code _res = plc4c_modbus_read_write_modbus_pdu_write_file_record_response_item_parse(buf, (void*) &_value);
        if(_res != OK) {
          return _res;
        }
        plc4c_utils_list_insert_head_value(items, _value);
      }
    }
    (*_message)->modbus_pdu_write_file_record_response_items = items;

  } else 
  if((error == false) && (function == 0x07) && (response == false)) { /* ModbusPDUReadExceptionStatusRequest */
    (*_message)->_type = plc4c_modbus_read_write_modbus_pdu_type_modbus_read_write_modbus_pdu_read_exception_status_request;
  } else 
  if((error == false) && (function == 0x07) && (response == true)) { /* ModbusPDUReadExceptionStatusResponse */
    (*_message)->_type = plc4c_modbus_read_write_modbus_pdu_type_modbus_read_write_modbus_pdu_read_exception_status_response;
                    
    // Simple Field (value)
    uint8_t value = plc4c_spi_read_unsigned_short(buf, 8);
    (*_message)->modbus_pdu_read_exception_status_response_value = value;

  } else 
  if((error == false) && (function == 0x08) && (response == false)) { /* ModbusPDUDiagnosticRequest */
    (*_message)->_type = plc4c_modbus_read_write_modbus_pdu_type_modbus_read_write_modbus_pdu_diagnostic_request;
                    
    // Simple Field (status)
    uint16_t status = plc4c_spi_read_unsigned_int(buf, 16);
    (*_message)->modbus_pdu_diagnostic_request_status = status;


                    
    // Simple Field (eventCount)
    uint16_t eventCount = plc4c_spi_read_unsigned_int(buf, 16);
    (*_message)->modbus_pdu_diagnostic_request_event_count = eventCount;

  } else 
  if((error == false) && (function == 0x0C) && (response == false)) { /* ModbusPDUGetComEventLogRequest */
    (*_message)->_type = plc4c_modbus_read_write_modbus_pdu_type_modbus_read_write_modbus_pdu_get_com_event_log_request;
  } else 
  if((error == false) && (function == 0x0C) && (response == true)) { /* ModbusPDUGetComEventLogResponse */
    (*_message)->_type = plc4c_modbus_read_write_modbus_pdu_type_modbus_read_write_modbus_pdu_get_com_event_log_response;
                    
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
    plc4c_list* events = malloc(sizeof(plc4c_list));
    if(events == NULL) {
      return NO_MEMORY;
    }
    {
      // Count array
      uint8_t itemCount = (byteCount) - (6);
      for(int curItem = 0; curItem < itemCount; curItem++) {
        
                  
        int8_t _value = plc4c_spi_read_byte(buf, 8);
        plc4c_utils_list_insert_head_value(events, &_value);
      }
    }
    (*_message)->modbus_pdu_get_com_event_log_response_events = events;

  } else 
  if((error == false) && (function == 0x11) && (response == false)) { /* ModbusPDUReportServerIdRequest */
    (*_message)->_type = plc4c_modbus_read_write_modbus_pdu_type_modbus_read_write_modbus_pdu_report_server_id_request;
  } else 
  if((error == false) && (function == 0x11) && (response == true)) { /* ModbusPDUReportServerIdResponse */
    (*_message)->_type = plc4c_modbus_read_write_modbus_pdu_type_modbus_read_write_modbus_pdu_report_server_id_response;
                    
    // Implicit Field (byteCount) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
    uint8_t byteCount = plc4c_spi_read_unsigned_short(buf, 8);


                    
    // Array field (value)
    plc4c_list* value = malloc(sizeof(plc4c_list));
    if(value == NULL) {
      return NO_MEMORY;
    }
    {
      // Count array
      uint8_t itemCount = byteCount;
      for(int curItem = 0; curItem < itemCount; curItem++) {
        
                  
        int8_t _value = plc4c_spi_read_byte(buf, 8);
        plc4c_utils_list_insert_head_value(value, &_value);
      }
    }
    (*_message)->modbus_pdu_report_server_id_response_value = value;

  } else 
  if((error == false) && (function == 0x2B) && (response == false)) { /* ModbusPDUReadDeviceIdentificationRequest */
    (*_message)->_type = plc4c_modbus_read_write_modbus_pdu_type_modbus_read_write_modbus_pdu_read_device_identification_request;
  } else 
  if((error == false) && (function == 0x2B) && (response == true)) { /* ModbusPDUReadDeviceIdentificationResponse */
    (*_message)->_type = plc4c_modbus_read_write_modbus_pdu_type_modbus_read_write_modbus_pdu_read_device_identification_response;
  }

  return OK;
}

plc4c_return_code plc4c_modbus_read_write_modbus_pdu_serialize(plc4c_spi_write_buffer* buf, plc4c_modbus_read_write_modbus_pdu* _message) {

  // Discriminator Field (error)
  plc4c_spi_write_bit(buf, plc4c_modbus_read_write_modbus_pdu_get_discriminator(_message->_type).error);

  // Discriminator Field (function)
  plc4c_spi_write_unsigned_short(buf, 7, plc4c_modbus_read_write_modbus_pdu_get_discriminator(_message->_type).function);

  return OK;
}

uint8_t plc4c_modbus_read_write_modbus_pdu_length_in_bytes(plc4c_modbus_read_write_modbus_pdu* _message) {
  return plc4c_modbus_read_write_modbus_pdu_length_in_bits(_message) / 8;
}

uint8_t plc4c_modbus_read_write_modbus_pdu_length_in_bits(plc4c_modbus_read_write_modbus_pdu* _message) {
  uint8_t lengthInBits = 0;

  // Discriminator Field (error)
  lengthInBits += 1;

  // Discriminator Field (function)
  lengthInBits += 7;

  // Depending of the current type, add the length of sub-type elements ...
  switch(_message->_type) {
    case plc4c_modbus_read_write_modbus_pdu_type_modbus_read_write_modbus_pdu_error: {

      // Simple field (exceptionCode)
      lengthInBits += 8;
      break;
    }
    case plc4c_modbus_read_write_modbus_pdu_type_modbus_read_write_modbus_pdu_read_discrete_inputs_request: {

      // Simple field (startingAddress)
      lengthInBits += 16;


      // Simple field (quantity)
      lengthInBits += 16;
      break;
    }
    case plc4c_modbus_read_write_modbus_pdu_type_modbus_read_write_modbus_pdu_read_discrete_inputs_response: {

      // Implicit Field (byteCount)
      lengthInBits += 8;


      // Array field
      lengthInBits += 8 * plc4c_utils_list_size(_message->modbus_pdu_read_discrete_inputs_response_value);
      break;
    }
    case plc4c_modbus_read_write_modbus_pdu_type_modbus_read_write_modbus_pdu_read_coils_request: {

      // Simple field (startingAddress)
      lengthInBits += 16;


      // Simple field (quantity)
      lengthInBits += 16;
      break;
    }
    case plc4c_modbus_read_write_modbus_pdu_type_modbus_read_write_modbus_pdu_read_coils_response: {

      // Implicit Field (byteCount)
      lengthInBits += 8;


      // Array field
      lengthInBits += 8 * plc4c_utils_list_size(_message->modbus_pdu_read_coils_response_value);
      break;
    }
    case plc4c_modbus_read_write_modbus_pdu_type_modbus_read_write_modbus_pdu_write_single_coil_request: {

      // Simple field (address)
      lengthInBits += 16;


      // Simple field (value)
      lengthInBits += 16;
      break;
    }
    case plc4c_modbus_read_write_modbus_pdu_type_modbus_read_write_modbus_pdu_write_single_coil_response: {

      // Simple field (address)
      lengthInBits += 16;


      // Simple field (value)
      lengthInBits += 16;
      break;
    }
    case plc4c_modbus_read_write_modbus_pdu_type_modbus_read_write_modbus_pdu_write_multiple_coils_request: {

      // Simple field (startingAddress)
      lengthInBits += 16;


      // Simple field (quantity)
      lengthInBits += 16;


      // Implicit Field (byteCount)
      lengthInBits += 8;


      // Array field
      lengthInBits += 8 * plc4c_utils_list_size(_message->modbus_pdu_write_multiple_coils_request_value);
      break;
    }
    case plc4c_modbus_read_write_modbus_pdu_type_modbus_read_write_modbus_pdu_write_multiple_coils_response: {

      // Simple field (startingAddress)
      lengthInBits += 16;


      // Simple field (quantity)
      lengthInBits += 16;
      break;
    }
    case plc4c_modbus_read_write_modbus_pdu_type_modbus_read_write_modbus_pdu_read_input_registers_request: {

      // Simple field (startingAddress)
      lengthInBits += 16;


      // Simple field (quantity)
      lengthInBits += 16;
      break;
    }
    case plc4c_modbus_read_write_modbus_pdu_type_modbus_read_write_modbus_pdu_read_input_registers_response: {

      // Implicit Field (byteCount)
      lengthInBits += 8;


      // Array field
      lengthInBits += 8 * plc4c_utils_list_size(_message->modbus_pdu_read_input_registers_response_value);
      break;
    }
    case plc4c_modbus_read_write_modbus_pdu_type_modbus_read_write_modbus_pdu_read_holding_registers_request: {

      // Simple field (startingAddress)
      lengthInBits += 16;


      // Simple field (quantity)
      lengthInBits += 16;
      break;
    }
    case plc4c_modbus_read_write_modbus_pdu_type_modbus_read_write_modbus_pdu_read_holding_registers_response: {

      // Implicit Field (byteCount)
      lengthInBits += 8;


      // Array field
      lengthInBits += 8 * plc4c_utils_list_size(_message->modbus_pdu_read_holding_registers_response_value);
      break;
    }
    case plc4c_modbus_read_write_modbus_pdu_type_modbus_read_write_modbus_pdu_write_single_register_request: {

      // Simple field (address)
      lengthInBits += 16;


      // Simple field (value)
      lengthInBits += 16;
      break;
    }
    case plc4c_modbus_read_write_modbus_pdu_type_modbus_read_write_modbus_pdu_write_single_register_response: {

      // Simple field (address)
      lengthInBits += 16;


      // Simple field (value)
      lengthInBits += 16;
      break;
    }
    case plc4c_modbus_read_write_modbus_pdu_type_modbus_read_write_modbus_pdu_write_multiple_holding_registers_request: {

      // Simple field (startingAddress)
      lengthInBits += 16;


      // Simple field (quantity)
      lengthInBits += 16;


      // Implicit Field (byteCount)
      lengthInBits += 8;


      // Array field
      lengthInBits += 8 * plc4c_utils_list_size(_message->modbus_pdu_write_multiple_holding_registers_request_value);
      break;
    }
    case plc4c_modbus_read_write_modbus_pdu_type_modbus_read_write_modbus_pdu_write_multiple_holding_registers_response: {

      // Simple field (startingAddress)
      lengthInBits += 16;


      // Simple field (quantity)
      lengthInBits += 16;
      break;
    }
    case plc4c_modbus_read_write_modbus_pdu_type_modbus_read_write_modbus_pdu_read_write_multiple_holding_registers_request: {

      // Simple field (readStartingAddress)
      lengthInBits += 16;


      // Simple field (readQuantity)
      lengthInBits += 16;


      // Simple field (writeStartingAddress)
      lengthInBits += 16;


      // Simple field (writeQuantity)
      lengthInBits += 16;


      // Implicit Field (byteCount)
      lengthInBits += 8;


      // Array field
      lengthInBits += 8 * plc4c_utils_list_size(_message->modbus_pdu_read_write_multiple_holding_registers_request_value);
      break;
    }
    case plc4c_modbus_read_write_modbus_pdu_type_modbus_read_write_modbus_pdu_read_write_multiple_holding_registers_response: {

      // Implicit Field (byteCount)
      lengthInBits += 8;


      // Array field
      lengthInBits += 8 * plc4c_utils_list_size(_message->modbus_pdu_read_write_multiple_holding_registers_response_value);
      break;
    }
    case plc4c_modbus_read_write_modbus_pdu_type_modbus_read_write_modbus_pdu_mask_write_holding_register_request: {

      // Simple field (referenceAddress)
      lengthInBits += 16;


      // Simple field (andMask)
      lengthInBits += 16;


      // Simple field (orMask)
      lengthInBits += 16;
      break;
    }
    case plc4c_modbus_read_write_modbus_pdu_type_modbus_read_write_modbus_pdu_mask_write_holding_register_response: {

      // Simple field (referenceAddress)
      lengthInBits += 16;


      // Simple field (andMask)
      lengthInBits += 16;


      // Simple field (orMask)
      lengthInBits += 16;
      break;
    }
    case plc4c_modbus_read_write_modbus_pdu_type_modbus_read_write_modbus_pdu_read_fifo_queue_request: {

      // Simple field (fifoPointerAddress)
      lengthInBits += 16;
      break;
    }
    case plc4c_modbus_read_write_modbus_pdu_type_modbus_read_write_modbus_pdu_read_fifo_queue_response: {

      // Implicit Field (byteCount)
      lengthInBits += 16;


      // Implicit Field (fifoCount)
      lengthInBits += 16;


      // Array field
      lengthInBits += 16 * plc4c_utils_list_size(_message->modbus_pdu_read_fifo_queue_response_fifo_value);
      break;
    }
    case plc4c_modbus_read_write_modbus_pdu_type_modbus_read_write_modbus_pdu_read_file_record_request: {

      // Implicit Field (byteCount)
      lengthInBits += 8;


      // Array field
      if(_message->modbus_pdu_read_file_record_request_items != NULL) {
        plc4c_list_element* curElement = _message->modbus_pdu_read_file_record_request_items->head;
        while (curElement != NULL) {
          lengthInBits += plc4c_modbus_read_write_modbus_pdu_read_file_record_request_item_length_in_bits((plc4c_modbus_read_write_modbus_pdu_read_file_record_request_item*) curElement->value);
          curElement = curElement->next;
        }
      }
      break;
    }
    case plc4c_modbus_read_write_modbus_pdu_type_modbus_read_write_modbus_pdu_read_file_record_response: {

      // Implicit Field (byteCount)
      lengthInBits += 8;


      // Array field
      if(_message->modbus_pdu_read_file_record_response_items != NULL) {
        plc4c_list_element* curElement = _message->modbus_pdu_read_file_record_response_items->head;
        while (curElement != NULL) {
          lengthInBits += plc4c_modbus_read_write_modbus_pdu_read_file_record_response_item_length_in_bits((plc4c_modbus_read_write_modbus_pdu_read_file_record_response_item*) curElement->value);
          curElement = curElement->next;
        }
      }
      break;
    }
    case plc4c_modbus_read_write_modbus_pdu_type_modbus_read_write_modbus_pdu_write_file_record_request: {

      // Implicit Field (byteCount)
      lengthInBits += 8;


      // Array field
      if(_message->modbus_pdu_write_file_record_request_items != NULL) {
        plc4c_list_element* curElement = _message->modbus_pdu_write_file_record_request_items->head;
        while (curElement != NULL) {
          lengthInBits += plc4c_modbus_read_write_modbus_pdu_write_file_record_request_item_length_in_bits((plc4c_modbus_read_write_modbus_pdu_write_file_record_request_item*) curElement->value);
          curElement = curElement->next;
        }
      }
      break;
    }
    case plc4c_modbus_read_write_modbus_pdu_type_modbus_read_write_modbus_pdu_write_file_record_response: {

      // Implicit Field (byteCount)
      lengthInBits += 8;


      // Array field
      if(_message->modbus_pdu_write_file_record_response_items != NULL) {
        plc4c_list_element* curElement = _message->modbus_pdu_write_file_record_response_items->head;
        while (curElement != NULL) {
          lengthInBits += plc4c_modbus_read_write_modbus_pdu_write_file_record_response_item_length_in_bits((plc4c_modbus_read_write_modbus_pdu_write_file_record_response_item*) curElement->value);
          curElement = curElement->next;
        }
      }
      break;
    }
    case plc4c_modbus_read_write_modbus_pdu_type_modbus_read_write_modbus_pdu_read_exception_status_request: {
      break;
    }
    case plc4c_modbus_read_write_modbus_pdu_type_modbus_read_write_modbus_pdu_read_exception_status_response: {

      // Simple field (value)
      lengthInBits += 8;
      break;
    }
    case plc4c_modbus_read_write_modbus_pdu_type_modbus_read_write_modbus_pdu_diagnostic_request: {

      // Simple field (status)
      lengthInBits += 16;


      // Simple field (eventCount)
      lengthInBits += 16;
      break;
    }
    case plc4c_modbus_read_write_modbus_pdu_type_modbus_read_write_modbus_pdu_get_com_event_log_request: {
      break;
    }
    case plc4c_modbus_read_write_modbus_pdu_type_modbus_read_write_modbus_pdu_get_com_event_log_response: {

      // Implicit Field (byteCount)
      lengthInBits += 8;


      // Simple field (status)
      lengthInBits += 16;


      // Simple field (eventCount)
      lengthInBits += 16;


      // Simple field (messageCount)
      lengthInBits += 16;


      // Array field
      lengthInBits += 8 * plc4c_utils_list_size(_message->modbus_pdu_get_com_event_log_response_events);
      break;
    }
    case plc4c_modbus_read_write_modbus_pdu_type_modbus_read_write_modbus_pdu_report_server_id_request: {
      break;
    }
    case plc4c_modbus_read_write_modbus_pdu_type_modbus_read_write_modbus_pdu_report_server_id_response: {

      // Implicit Field (byteCount)
      lengthInBits += 8;


      // Array field
      lengthInBits += 8 * plc4c_utils_list_size(_message->modbus_pdu_report_server_id_response_value);
      break;
    }
    case plc4c_modbus_read_write_modbus_pdu_type_modbus_read_write_modbus_pdu_read_device_identification_request: {
      break;
    }
    case plc4c_modbus_read_write_modbus_pdu_type_modbus_read_write_modbus_pdu_read_device_identification_response: {
      break;
    }
  }

  return lengthInBits;
}

