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
#include <plc4c/spi/evaluation_helper.h>
#include "modbus_pdu.h"

// Array of discriminator values that match the enum type constants.
// (The order is identical to the enum constants so we can use the
// enum constant to directly access a given types discriminator values)
const plc4c_modbus_read_write_modbus_pdu_discriminator plc4c_modbus_read_write_modbus_pdu_discriminators[] = {
  {/* plc4c_modbus_read_write_modbus_pdu_error */
    .errorFlag = true,
       .functionFlag = -1,
       .response = -1   },
  {/* plc4c_modbus_read_write_modbus_pdu_read_discrete_inputs_request */
    .errorFlag = false,
       .functionFlag = 0x02,
       .response = false   },
  {/* plc4c_modbus_read_write_modbus_pdu_read_discrete_inputs_response */
    .errorFlag = false,
       .functionFlag = 0x02,
       .response = true   },
  {/* plc4c_modbus_read_write_modbus_pdu_read_coils_request */
    .errorFlag = false,
       .functionFlag = 0x01,
       .response = false   },
  {/* plc4c_modbus_read_write_modbus_pdu_read_coils_response */
    .errorFlag = false,
       .functionFlag = 0x01,
       .response = true   },
  {/* plc4c_modbus_read_write_modbus_pdu_write_single_coil_request */
    .errorFlag = false,
       .functionFlag = 0x05,
       .response = false   },
  {/* plc4c_modbus_read_write_modbus_pdu_write_single_coil_response */
    .errorFlag = false,
       .functionFlag = 0x05,
       .response = true   },
  {/* plc4c_modbus_read_write_modbus_pdu_write_multiple_coils_request */
    .errorFlag = false,
       .functionFlag = 0x0F,
       .response = false   },
  {/* plc4c_modbus_read_write_modbus_pdu_write_multiple_coils_response */
    .errorFlag = false,
       .functionFlag = 0x0F,
       .response = true   },
  {/* plc4c_modbus_read_write_modbus_pdu_read_input_registers_request */
    .errorFlag = false,
       .functionFlag = 0x04,
       .response = false   },
  {/* plc4c_modbus_read_write_modbus_pdu_read_input_registers_response */
    .errorFlag = false,
       .functionFlag = 0x04,
       .response = true   },
  {/* plc4c_modbus_read_write_modbus_pdu_read_holding_registers_request */
    .errorFlag = false,
       .functionFlag = 0x03,
       .response = false   },
  {/* plc4c_modbus_read_write_modbus_pdu_read_holding_registers_response */
    .errorFlag = false,
       .functionFlag = 0x03,
       .response = true   },
  {/* plc4c_modbus_read_write_modbus_pdu_write_single_register_request */
    .errorFlag = false,
       .functionFlag = 0x06,
       .response = false   },
  {/* plc4c_modbus_read_write_modbus_pdu_write_single_register_response */
    .errorFlag = false,
       .functionFlag = 0x06,
       .response = true   },
  {/* plc4c_modbus_read_write_modbus_pdu_write_multiple_holding_registers_request */
    .errorFlag = false,
       .functionFlag = 0x10,
       .response = false   },
  {/* plc4c_modbus_read_write_modbus_pdu_write_multiple_holding_registers_response */
    .errorFlag = false,
       .functionFlag = 0x10,
       .response = true   },
  {/* plc4c_modbus_read_write_modbus_pdu_read_write_multiple_holding_registers_request */
    .errorFlag = false,
       .functionFlag = 0x17,
       .response = false   },
  {/* plc4c_modbus_read_write_modbus_pdu_read_write_multiple_holding_registers_response */
    .errorFlag = false,
       .functionFlag = 0x17,
       .response = true   },
  {/* plc4c_modbus_read_write_modbus_pdu_mask_write_holding_register_request */
    .errorFlag = false,
       .functionFlag = 0x16,
       .response = false   },
  {/* plc4c_modbus_read_write_modbus_pdu_mask_write_holding_register_response */
    .errorFlag = false,
       .functionFlag = 0x16,
       .response = true   },
  {/* plc4c_modbus_read_write_modbus_pdu_read_fifo_queue_request */
    .errorFlag = false,
       .functionFlag = 0x18,
       .response = false   },
  {/* plc4c_modbus_read_write_modbus_pdu_read_fifo_queue_response */
    .errorFlag = false,
       .functionFlag = 0x18,
       .response = true   },
  {/* plc4c_modbus_read_write_modbus_pdu_read_file_record_request */
    .errorFlag = false,
       .functionFlag = 0x14,
       .response = false   },
  {/* plc4c_modbus_read_write_modbus_pdu_read_file_record_response */
    .errorFlag = false,
       .functionFlag = 0x14,
       .response = true   },
  {/* plc4c_modbus_read_write_modbus_pdu_write_file_record_request */
    .errorFlag = false,
       .functionFlag = 0x15,
       .response = false   },
  {/* plc4c_modbus_read_write_modbus_pdu_write_file_record_response */
    .errorFlag = false,
       .functionFlag = 0x15,
       .response = true   },
  {/* plc4c_modbus_read_write_modbus_pdu_read_exception_status_request */
    .errorFlag = false,
       .functionFlag = 0x07,
       .response = false   },
  {/* plc4c_modbus_read_write_modbus_pdu_read_exception_status_response */
    .errorFlag = false,
       .functionFlag = 0x07,
       .response = true   },
  {/* plc4c_modbus_read_write_modbus_pdu_diagnostic_request */
    .errorFlag = false,
       .functionFlag = 0x08,
       .response = false   },
  {/* plc4c_modbus_read_write_modbus_pdu_diagnostic_response */
    .errorFlag = false,
       .functionFlag = 0x08,
       .response = true   },
  {/* plc4c_modbus_read_write_modbus_pdu_get_com_event_counter_request */
    .errorFlag = false,
       .functionFlag = 0x0B,
       .response = false   },
  {/* plc4c_modbus_read_write_modbus_pdu_get_com_event_counter_response */
    .errorFlag = false,
       .functionFlag = 0x0B,
       .response = true   },
  {/* plc4c_modbus_read_write_modbus_pdu_get_com_event_log_request */
    .errorFlag = false,
       .functionFlag = 0x0C,
       .response = false   },
  {/* plc4c_modbus_read_write_modbus_pdu_get_com_event_log_response */
    .errorFlag = false,
       .functionFlag = 0x0C,
       .response = true   },
  {/* plc4c_modbus_read_write_modbus_pdu_report_server_id_request */
    .errorFlag = false,
       .functionFlag = 0x11,
       .response = false   },
  {/* plc4c_modbus_read_write_modbus_pdu_report_server_id_response */
    .errorFlag = false,
       .functionFlag = 0x11,
       .response = true   },
  {/* plc4c_modbus_read_write_modbus_pdu_read_device_identification_request */
    .errorFlag = false,
       .functionFlag = 0x2B,
       .response = false   },
  {/* plc4c_modbus_read_write_modbus_pdu_read_device_identification_response */
    .errorFlag = false,
       .functionFlag = 0x2B,
       .response = true   }
};

// Function returning the discriminator values for a given type constant.
plc4c_modbus_read_write_modbus_pdu_discriminator plc4c_modbus_read_write_modbus_pdu_get_discriminator(plc4c_modbus_read_write_modbus_pdu_type type) {
  return plc4c_modbus_read_write_modbus_pdu_discriminators[type];
}

// Create an empty NULL-struct
static const plc4c_modbus_read_write_modbus_pdu plc4c_modbus_read_write_modbus_pdu_null_const;

plc4c_modbus_read_write_modbus_pdu plc4c_modbus_read_write_modbus_pdu_null() {
  return plc4c_modbus_read_write_modbus_pdu_null_const;
}


// Parse function.
plc4c_return_code plc4c_modbus_read_write_modbus_pdu_parse(plc4c_spi_read_buffer* io, bool response, plc4c_modbus_read_write_modbus_pdu** _message) {
  uint16_t startPos = plc4c_spi_read_get_pos(io);
  uint16_t curPos;
  plc4c_return_code _res = OK;

  // Allocate enough memory to contain this data structure.
  (*_message) = malloc(sizeof(plc4c_modbus_read_write_modbus_pdu));
  if(*_message == NULL) {
    return NO_MEMORY;
  }
        // Discriminator Field (errorFlag)
  // Discriminator Field (errorFlag) (Used as input to a switch field)
  bool errorFlag = false;
  _res = plc4c_spi_read_bit(io, (bool*) &errorFlag);
  if(_res != OK) {
    return _res;
  }
        // Discriminator Field (functionFlag)
  // Discriminator Field (functionFlag) (Used as input to a switch field)
  uint8_t functionFlag = 0;
  _res = plc4c_spi_read_unsigned_byte(io, 7, (uint8_t*) &functionFlag);
  if(_res != OK) {
    return _res;
  }
  // Switch Field (Depending on the discriminator values, passes the instantiation to a sub-type)
  if(errorFlag == true) { /* ModbusPDUError */
    (*_message)->_type = plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_error;
                    
    // Enum field (exceptionCode)
    plc4c_modbus_read_write_modbus_error_code exceptionCode = plc4c_modbus_read_write_modbus_error_code_null();
    _res = plc4c_spi_read_unsigned_byte(io, 8, (uint8_t*) &exceptionCode);
    if(_res != OK) {
      return _res;
    }
    (*_message)->modbus_pdu_error_exception_code = exceptionCode;

  } else 
  if((errorFlag == false) && (functionFlag == 0x02) && (response == false)) { /* ModbusPDUReadDiscreteInputsRequest */
    (*_message)->_type = plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_read_discrete_inputs_request;
                    
    // Simple Field (startingAddress)
    uint16_t startingAddress = 0;
    _res = plc4c_spi_read_unsigned_short(io, 16, (uint16_t*) &startingAddress);
    if(_res != OK) {
      return _res;
    }
    (*_message)->modbus_pdu_read_discrete_inputs_request_starting_address = startingAddress;


                    
    // Simple Field (quantity)
    uint16_t quantity = 0;
    _res = plc4c_spi_read_unsigned_short(io, 16, (uint16_t*) &quantity);
    if(_res != OK) {
      return _res;
    }
    (*_message)->modbus_pdu_read_discrete_inputs_request_quantity = quantity;

  } else 
  if((errorFlag == false) && (functionFlag == 0x02) && (response == true)) { /* ModbusPDUReadDiscreteInputsResponse */
    (*_message)->_type = plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_read_discrete_inputs_response;
                    
    // Implicit Field (byteCount) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
    uint8_t byteCount = 0;
    _res = plc4c_spi_read_unsigned_byte(io, 8, (uint8_t*) &byteCount);
    if(_res != OK) {
      return _res;
    }


                    
    // Array field (value)
    plc4c_list* value = NULL;
    plc4c_utils_list_create(&value);
    if(value == NULL) {
      return NO_MEMORY;
    }
    {
      // Count array
      uint8_t itemCount = byteCount;
      for(int curItem = 0; curItem < itemCount; curItem++) {
        
        int8_t* _value = malloc(sizeof(int8_t));
        _res = plc4c_spi_read_signed_byte(io, 8, (int8_t*) _value);
        if(_res != OK) {
          return _res;
        }
        plc4c_utils_list_insert_head_value(value, _value);
      }
    }
    (*_message)->modbus_pdu_read_discrete_inputs_response_value = value;

  } else 
  if((errorFlag == false) && (functionFlag == 0x01) && (response == false)) { /* ModbusPDUReadCoilsRequest */
    (*_message)->_type = plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_read_coils_request;
                    
    // Simple Field (startingAddress)
    uint16_t startingAddress = 0;
    _res = plc4c_spi_read_unsigned_short(io, 16, (uint16_t*) &startingAddress);
    if(_res != OK) {
      return _res;
    }
    (*_message)->modbus_pdu_read_coils_request_starting_address = startingAddress;


                    
    // Simple Field (quantity)
    uint16_t quantity = 0;
    _res = plc4c_spi_read_unsigned_short(io, 16, (uint16_t*) &quantity);
    if(_res != OK) {
      return _res;
    }
    (*_message)->modbus_pdu_read_coils_request_quantity = quantity;

  } else 
  if((errorFlag == false) && (functionFlag == 0x01) && (response == true)) { /* ModbusPDUReadCoilsResponse */
    (*_message)->_type = plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_read_coils_response;
                    
    // Implicit Field (byteCount) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
    uint8_t byteCount = 0;
    _res = plc4c_spi_read_unsigned_byte(io, 8, (uint8_t*) &byteCount);
    if(_res != OK) {
      return _res;
    }


                    
    // Array field (value)
    plc4c_list* value = NULL;
    plc4c_utils_list_create(&value);
    if(value == NULL) {
      return NO_MEMORY;
    }
    {
      // Count array
      uint8_t itemCount = byteCount;
      for(int curItem = 0; curItem < itemCount; curItem++) {
        
        int8_t* _value = malloc(sizeof(int8_t));
        _res = plc4c_spi_read_signed_byte(io, 8, (int8_t*) _value);
        if(_res != OK) {
          return _res;
        }
        plc4c_utils_list_insert_head_value(value, _value);
      }
    }
    (*_message)->modbus_pdu_read_coils_response_value = value;

  } else 
  if((errorFlag == false) && (functionFlag == 0x05) && (response == false)) { /* ModbusPDUWriteSingleCoilRequest */
    (*_message)->_type = plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_write_single_coil_request;
                    
    // Simple Field (address)
    uint16_t address = 0;
    _res = plc4c_spi_read_unsigned_short(io, 16, (uint16_t*) &address);
    if(_res != OK) {
      return _res;
    }
    (*_message)->modbus_pdu_write_single_coil_request_address = address;


                    
    // Simple Field (value)
    uint16_t value = 0;
    _res = plc4c_spi_read_unsigned_short(io, 16, (uint16_t*) &value);
    if(_res != OK) {
      return _res;
    }
    (*_message)->modbus_pdu_write_single_coil_request_value = value;

  } else 
  if((errorFlag == false) && (functionFlag == 0x05) && (response == true)) { /* ModbusPDUWriteSingleCoilResponse */
    (*_message)->_type = plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_write_single_coil_response;
                    
    // Simple Field (address)
    uint16_t address = 0;
    _res = plc4c_spi_read_unsigned_short(io, 16, (uint16_t*) &address);
    if(_res != OK) {
      return _res;
    }
    (*_message)->modbus_pdu_write_single_coil_response_address = address;


                    
    // Simple Field (value)
    uint16_t value = 0;
    _res = plc4c_spi_read_unsigned_short(io, 16, (uint16_t*) &value);
    if(_res != OK) {
      return _res;
    }
    (*_message)->modbus_pdu_write_single_coil_response_value = value;

  } else 
  if((errorFlag == false) && (functionFlag == 0x0F) && (response == false)) { /* ModbusPDUWriteMultipleCoilsRequest */
    (*_message)->_type = plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_write_multiple_coils_request;
                    
    // Simple Field (startingAddress)
    uint16_t startingAddress = 0;
    _res = plc4c_spi_read_unsigned_short(io, 16, (uint16_t*) &startingAddress);
    if(_res != OK) {
      return _res;
    }
    (*_message)->modbus_pdu_write_multiple_coils_request_starting_address = startingAddress;


                    
    // Simple Field (quantity)
    uint16_t quantity = 0;
    _res = plc4c_spi_read_unsigned_short(io, 16, (uint16_t*) &quantity);
    if(_res != OK) {
      return _res;
    }
    (*_message)->modbus_pdu_write_multiple_coils_request_quantity = quantity;


                    
    // Implicit Field (byteCount) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
    uint8_t byteCount = 0;
    _res = plc4c_spi_read_unsigned_byte(io, 8, (uint8_t*) &byteCount);
    if(_res != OK) {
      return _res;
    }


                    
    // Array field (value)
    plc4c_list* value = NULL;
    plc4c_utils_list_create(&value);
    if(value == NULL) {
      return NO_MEMORY;
    }
    {
      // Count array
      uint8_t itemCount = byteCount;
      for(int curItem = 0; curItem < itemCount; curItem++) {
        
        int8_t* _value = malloc(sizeof(int8_t));
        _res = plc4c_spi_read_signed_byte(io, 8, (int8_t*) _value);
        if(_res != OK) {
          return _res;
        }
        plc4c_utils_list_insert_head_value(value, _value);
      }
    }
    (*_message)->modbus_pdu_write_multiple_coils_request_value = value;

  } else 
  if((errorFlag == false) && (functionFlag == 0x0F) && (response == true)) { /* ModbusPDUWriteMultipleCoilsResponse */
    (*_message)->_type = plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_write_multiple_coils_response;
                    
    // Simple Field (startingAddress)
    uint16_t startingAddress = 0;
    _res = plc4c_spi_read_unsigned_short(io, 16, (uint16_t*) &startingAddress);
    if(_res != OK) {
      return _res;
    }
    (*_message)->modbus_pdu_write_multiple_coils_response_starting_address = startingAddress;


                    
    // Simple Field (quantity)
    uint16_t quantity = 0;
    _res = plc4c_spi_read_unsigned_short(io, 16, (uint16_t*) &quantity);
    if(_res != OK) {
      return _res;
    }
    (*_message)->modbus_pdu_write_multiple_coils_response_quantity = quantity;

  } else 
  if((errorFlag == false) && (functionFlag == 0x04) && (response == false)) { /* ModbusPDUReadInputRegistersRequest */
    (*_message)->_type = plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_read_input_registers_request;
                    
    // Simple Field (startingAddress)
    uint16_t startingAddress = 0;
    _res = plc4c_spi_read_unsigned_short(io, 16, (uint16_t*) &startingAddress);
    if(_res != OK) {
      return _res;
    }
    (*_message)->modbus_pdu_read_input_registers_request_starting_address = startingAddress;


                    
    // Simple Field (quantity)
    uint16_t quantity = 0;
    _res = plc4c_spi_read_unsigned_short(io, 16, (uint16_t*) &quantity);
    if(_res != OK) {
      return _res;
    }
    (*_message)->modbus_pdu_read_input_registers_request_quantity = quantity;

  } else 
  if((errorFlag == false) && (functionFlag == 0x04) && (response == true)) { /* ModbusPDUReadInputRegistersResponse */
    (*_message)->_type = plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_read_input_registers_response;
                    
    // Implicit Field (byteCount) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
    uint8_t byteCount = 0;
    _res = plc4c_spi_read_unsigned_byte(io, 8, (uint8_t*) &byteCount);
    if(_res != OK) {
      return _res;
    }


                    
    // Array field (value)
    plc4c_list* value = NULL;
    plc4c_utils_list_create(&value);
    if(value == NULL) {
      return NO_MEMORY;
    }
    {
      // Count array
      uint8_t itemCount = byteCount;
      for(int curItem = 0; curItem < itemCount; curItem++) {
        
        int8_t* _value = malloc(sizeof(int8_t));
        _res = plc4c_spi_read_signed_byte(io, 8, (int8_t*) _value);
        if(_res != OK) {
          return _res;
        }
        plc4c_utils_list_insert_head_value(value, _value);
      }
    }
    (*_message)->modbus_pdu_read_input_registers_response_value = value;

  } else 
  if((errorFlag == false) && (functionFlag == 0x03) && (response == false)) { /* ModbusPDUReadHoldingRegistersRequest */
    (*_message)->_type = plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_read_holding_registers_request;
                    
    // Simple Field (startingAddress)
    uint16_t startingAddress = 0;
    _res = plc4c_spi_read_unsigned_short(io, 16, (uint16_t*) &startingAddress);
    if(_res != OK) {
      return _res;
    }
    (*_message)->modbus_pdu_read_holding_registers_request_starting_address = startingAddress;


                    
    // Simple Field (quantity)
    uint16_t quantity = 0;
    _res = plc4c_spi_read_unsigned_short(io, 16, (uint16_t*) &quantity);
    if(_res != OK) {
      return _res;
    }
    (*_message)->modbus_pdu_read_holding_registers_request_quantity = quantity;

  } else 
  if((errorFlag == false) && (functionFlag == 0x03) && (response == true)) { /* ModbusPDUReadHoldingRegistersResponse */
    (*_message)->_type = plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_read_holding_registers_response;
                    
    // Implicit Field (byteCount) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
    uint8_t byteCount = 0;
    _res = plc4c_spi_read_unsigned_byte(io, 8, (uint8_t*) &byteCount);
    if(_res != OK) {
      return _res;
    }


                    
    // Array field (value)
    plc4c_list* value = NULL;
    plc4c_utils_list_create(&value);
    if(value == NULL) {
      return NO_MEMORY;
    }
    {
      // Count array
      uint8_t itemCount = byteCount;
      for(int curItem = 0; curItem < itemCount; curItem++) {
        
        int8_t* _value = malloc(sizeof(int8_t));
        _res = plc4c_spi_read_signed_byte(io, 8, (int8_t*) _value);
        if(_res != OK) {
          return _res;
        }
        plc4c_utils_list_insert_head_value(value, _value);
      }
    }
    (*_message)->modbus_pdu_read_holding_registers_response_value = value;

  } else 
  if((errorFlag == false) && (functionFlag == 0x06) && (response == false)) { /* ModbusPDUWriteSingleRegisterRequest */
    (*_message)->_type = plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_write_single_register_request;
                    
    // Simple Field (address)
    uint16_t address = 0;
    _res = plc4c_spi_read_unsigned_short(io, 16, (uint16_t*) &address);
    if(_res != OK) {
      return _res;
    }
    (*_message)->modbus_pdu_write_single_register_request_address = address;


                    
    // Simple Field (value)
    uint16_t value = 0;
    _res = plc4c_spi_read_unsigned_short(io, 16, (uint16_t*) &value);
    if(_res != OK) {
      return _res;
    }
    (*_message)->modbus_pdu_write_single_register_request_value = value;

  } else 
  if((errorFlag == false) && (functionFlag == 0x06) && (response == true)) { /* ModbusPDUWriteSingleRegisterResponse */
    (*_message)->_type = plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_write_single_register_response;
                    
    // Simple Field (address)
    uint16_t address = 0;
    _res = plc4c_spi_read_unsigned_short(io, 16, (uint16_t*) &address);
    if(_res != OK) {
      return _res;
    }
    (*_message)->modbus_pdu_write_single_register_response_address = address;


                    
    // Simple Field (value)
    uint16_t value = 0;
    _res = plc4c_spi_read_unsigned_short(io, 16, (uint16_t*) &value);
    if(_res != OK) {
      return _res;
    }
    (*_message)->modbus_pdu_write_single_register_response_value = value;

  } else 
  if((errorFlag == false) && (functionFlag == 0x10) && (response == false)) { /* ModbusPDUWriteMultipleHoldingRegistersRequest */
    (*_message)->_type = plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_write_multiple_holding_registers_request;
                    
    // Simple Field (startingAddress)
    uint16_t startingAddress = 0;
    _res = plc4c_spi_read_unsigned_short(io, 16, (uint16_t*) &startingAddress);
    if(_res != OK) {
      return _res;
    }
    (*_message)->modbus_pdu_write_multiple_holding_registers_request_starting_address = startingAddress;


                    
    // Simple Field (quantity)
    uint16_t quantity = 0;
    _res = plc4c_spi_read_unsigned_short(io, 16, (uint16_t*) &quantity);
    if(_res != OK) {
      return _res;
    }
    (*_message)->modbus_pdu_write_multiple_holding_registers_request_quantity = quantity;


                    
    // Implicit Field (byteCount) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
    uint8_t byteCount = 0;
    _res = plc4c_spi_read_unsigned_byte(io, 8, (uint8_t*) &byteCount);
    if(_res != OK) {
      return _res;
    }


                    
    // Array field (value)
    plc4c_list* value = NULL;
    plc4c_utils_list_create(&value);
    if(value == NULL) {
      return NO_MEMORY;
    }
    {
      // Count array
      uint8_t itemCount = byteCount;
      for(int curItem = 0; curItem < itemCount; curItem++) {
        
        int8_t* _value = malloc(sizeof(int8_t));
        _res = plc4c_spi_read_signed_byte(io, 8, (int8_t*) _value);
        if(_res != OK) {
          return _res;
        }
        plc4c_utils_list_insert_head_value(value, _value);
      }
    }
    (*_message)->modbus_pdu_write_multiple_holding_registers_request_value = value;

  } else 
  if((errorFlag == false) && (functionFlag == 0x10) && (response == true)) { /* ModbusPDUWriteMultipleHoldingRegistersResponse */
    (*_message)->_type = plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_write_multiple_holding_registers_response;
                    
    // Simple Field (startingAddress)
    uint16_t startingAddress = 0;
    _res = plc4c_spi_read_unsigned_short(io, 16, (uint16_t*) &startingAddress);
    if(_res != OK) {
      return _res;
    }
    (*_message)->modbus_pdu_write_multiple_holding_registers_response_starting_address = startingAddress;


                    
    // Simple Field (quantity)
    uint16_t quantity = 0;
    _res = plc4c_spi_read_unsigned_short(io, 16, (uint16_t*) &quantity);
    if(_res != OK) {
      return _res;
    }
    (*_message)->modbus_pdu_write_multiple_holding_registers_response_quantity = quantity;

  } else 
  if((errorFlag == false) && (functionFlag == 0x17) && (response == false)) { /* ModbusPDUReadWriteMultipleHoldingRegistersRequest */
    (*_message)->_type = plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_read_write_multiple_holding_registers_request;
                    
    // Simple Field (readStartingAddress)
    uint16_t readStartingAddress = 0;
    _res = plc4c_spi_read_unsigned_short(io, 16, (uint16_t*) &readStartingAddress);
    if(_res != OK) {
      return _res;
    }
    (*_message)->modbus_pdu_read_write_multiple_holding_registers_request_read_starting_address = readStartingAddress;


                    
    // Simple Field (readQuantity)
    uint16_t readQuantity = 0;
    _res = plc4c_spi_read_unsigned_short(io, 16, (uint16_t*) &readQuantity);
    if(_res != OK) {
      return _res;
    }
    (*_message)->modbus_pdu_read_write_multiple_holding_registers_request_read_quantity = readQuantity;


                    
    // Simple Field (writeStartingAddress)
    uint16_t writeStartingAddress = 0;
    _res = plc4c_spi_read_unsigned_short(io, 16, (uint16_t*) &writeStartingAddress);
    if(_res != OK) {
      return _res;
    }
    (*_message)->modbus_pdu_read_write_multiple_holding_registers_request_write_starting_address = writeStartingAddress;


                    
    // Simple Field (writeQuantity)
    uint16_t writeQuantity = 0;
    _res = plc4c_spi_read_unsigned_short(io, 16, (uint16_t*) &writeQuantity);
    if(_res != OK) {
      return _res;
    }
    (*_message)->modbus_pdu_read_write_multiple_holding_registers_request_write_quantity = writeQuantity;


                    
    // Implicit Field (byteCount) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
    uint8_t byteCount = 0;
    _res = plc4c_spi_read_unsigned_byte(io, 8, (uint8_t*) &byteCount);
    if(_res != OK) {
      return _res;
    }


                    
    // Array field (value)
    plc4c_list* value = NULL;
    plc4c_utils_list_create(&value);
    if(value == NULL) {
      return NO_MEMORY;
    }
    {
      // Count array
      uint8_t itemCount = byteCount;
      for(int curItem = 0; curItem < itemCount; curItem++) {
        
        int8_t* _value = malloc(sizeof(int8_t));
        _res = plc4c_spi_read_signed_byte(io, 8, (int8_t*) _value);
        if(_res != OK) {
          return _res;
        }
        plc4c_utils_list_insert_head_value(value, _value);
      }
    }
    (*_message)->modbus_pdu_read_write_multiple_holding_registers_request_value = value;

  } else 
  if((errorFlag == false) && (functionFlag == 0x17) && (response == true)) { /* ModbusPDUReadWriteMultipleHoldingRegistersResponse */
    (*_message)->_type = plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_read_write_multiple_holding_registers_response;
                    
    // Implicit Field (byteCount) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
    uint8_t byteCount = 0;
    _res = plc4c_spi_read_unsigned_byte(io, 8, (uint8_t*) &byteCount);
    if(_res != OK) {
      return _res;
    }


                    
    // Array field (value)
    plc4c_list* value = NULL;
    plc4c_utils_list_create(&value);
    if(value == NULL) {
      return NO_MEMORY;
    }
    {
      // Count array
      uint8_t itemCount = byteCount;
      for(int curItem = 0; curItem < itemCount; curItem++) {
        
        int8_t* _value = malloc(sizeof(int8_t));
        _res = plc4c_spi_read_signed_byte(io, 8, (int8_t*) _value);
        if(_res != OK) {
          return _res;
        }
        plc4c_utils_list_insert_head_value(value, _value);
      }
    }
    (*_message)->modbus_pdu_read_write_multiple_holding_registers_response_value = value;

  } else 
  if((errorFlag == false) && (functionFlag == 0x16) && (response == false)) { /* ModbusPDUMaskWriteHoldingRegisterRequest */
    (*_message)->_type = plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_mask_write_holding_register_request;
                    
    // Simple Field (referenceAddress)
    uint16_t referenceAddress = 0;
    _res = plc4c_spi_read_unsigned_short(io, 16, (uint16_t*) &referenceAddress);
    if(_res != OK) {
      return _res;
    }
    (*_message)->modbus_pdu_mask_write_holding_register_request_reference_address = referenceAddress;


                    
    // Simple Field (andMask)
    uint16_t andMask = 0;
    _res = plc4c_spi_read_unsigned_short(io, 16, (uint16_t*) &andMask);
    if(_res != OK) {
      return _res;
    }
    (*_message)->modbus_pdu_mask_write_holding_register_request_and_mask = andMask;


                    
    // Simple Field (orMask)
    uint16_t orMask = 0;
    _res = plc4c_spi_read_unsigned_short(io, 16, (uint16_t*) &orMask);
    if(_res != OK) {
      return _res;
    }
    (*_message)->modbus_pdu_mask_write_holding_register_request_or_mask = orMask;

  } else 
  if((errorFlag == false) && (functionFlag == 0x16) && (response == true)) { /* ModbusPDUMaskWriteHoldingRegisterResponse */
    (*_message)->_type = plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_mask_write_holding_register_response;
                    
    // Simple Field (referenceAddress)
    uint16_t referenceAddress = 0;
    _res = plc4c_spi_read_unsigned_short(io, 16, (uint16_t*) &referenceAddress);
    if(_res != OK) {
      return _res;
    }
    (*_message)->modbus_pdu_mask_write_holding_register_response_reference_address = referenceAddress;


                    
    // Simple Field (andMask)
    uint16_t andMask = 0;
    _res = plc4c_spi_read_unsigned_short(io, 16, (uint16_t*) &andMask);
    if(_res != OK) {
      return _res;
    }
    (*_message)->modbus_pdu_mask_write_holding_register_response_and_mask = andMask;


                    
    // Simple Field (orMask)
    uint16_t orMask = 0;
    _res = plc4c_spi_read_unsigned_short(io, 16, (uint16_t*) &orMask);
    if(_res != OK) {
      return _res;
    }
    (*_message)->modbus_pdu_mask_write_holding_register_response_or_mask = orMask;

  } else 
  if((errorFlag == false) && (functionFlag == 0x18) && (response == false)) { /* ModbusPDUReadFifoQueueRequest */
    (*_message)->_type = plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_read_fifo_queue_request;
                    
    // Simple Field (fifoPointerAddress)
    uint16_t fifoPointerAddress = 0;
    _res = plc4c_spi_read_unsigned_short(io, 16, (uint16_t*) &fifoPointerAddress);
    if(_res != OK) {
      return _res;
    }
    (*_message)->modbus_pdu_read_fifo_queue_request_fifo_pointer_address = fifoPointerAddress;

  } else 
  if((errorFlag == false) && (functionFlag == 0x18) && (response == true)) { /* ModbusPDUReadFifoQueueResponse */
    (*_message)->_type = plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_read_fifo_queue_response;
                    
    // Implicit Field (byteCount) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
    uint16_t byteCount = 0;
    _res = plc4c_spi_read_unsigned_short(io, 16, (uint16_t*) &byteCount);
    if(_res != OK) {
      return _res;
    }


                    
    // Implicit Field (fifoCount) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
    uint16_t fifoCount = 0;
    _res = plc4c_spi_read_unsigned_short(io, 16, (uint16_t*) &fifoCount);
    if(_res != OK) {
      return _res;
    }


                    
    // Array field (fifoValue)
    plc4c_list* fifoValue = NULL;
    plc4c_utils_list_create(&fifoValue);
    if(fifoValue == NULL) {
      return NO_MEMORY;
    }
    {
      // Count array
      uint8_t itemCount = fifoCount;
      for(int curItem = 0; curItem < itemCount; curItem++) {
        
        uint16_t* _value = malloc(sizeof(uint16_t));
        _res = plc4c_spi_read_unsigned_short(io, 16, (uint16_t*) _value);
        if(_res != OK) {
          return _res;
        }
        plc4c_utils_list_insert_head_value(fifoValue, _value);
      }
    }
    (*_message)->modbus_pdu_read_fifo_queue_response_fifo_value = fifoValue;

  } else 
  if((errorFlag == false) && (functionFlag == 0x14) && (response == false)) { /* ModbusPDUReadFileRecordRequest */
    (*_message)->_type = plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_read_file_record_request;
                    
    // Implicit Field (byteCount) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
    uint8_t byteCount = 0;
    _res = plc4c_spi_read_unsigned_byte(io, 8, (uint8_t*) &byteCount);
    if(_res != OK) {
      return _res;
    }


                    
    // Array field (items)
    plc4c_list* items = NULL;
    plc4c_utils_list_create(&items);
    if(items == NULL) {
      return NO_MEMORY;
    }
    {
      // Length array
      uint8_t _itemsLength = byteCount;
      uint8_t itemsEndPos = plc4c_spi_read_get_pos(io) + _itemsLength;
      while(plc4c_spi_read_get_pos(io) < itemsEndPos) {
        plc4c_modbus_read_write_modbus_pdu_read_file_record_request_item* _value = NULL;
        _res = plc4c_modbus_read_write_modbus_pdu_read_file_record_request_item_parse(io, (void*) &_value);
        if(_res != OK) {
          return _res;
        }
        plc4c_utils_list_insert_head_value(items, _value);
      }
    }
    (*_message)->modbus_pdu_read_file_record_request_items = items;

  } else 
  if((errorFlag == false) && (functionFlag == 0x14) && (response == true)) { /* ModbusPDUReadFileRecordResponse */
    (*_message)->_type = plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_read_file_record_response;
                    
    // Implicit Field (byteCount) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
    uint8_t byteCount = 0;
    _res = plc4c_spi_read_unsigned_byte(io, 8, (uint8_t*) &byteCount);
    if(_res != OK) {
      return _res;
    }


                    
    // Array field (items)
    plc4c_list* items = NULL;
    plc4c_utils_list_create(&items);
    if(items == NULL) {
      return NO_MEMORY;
    }
    {
      // Length array
      uint8_t _itemsLength = byteCount;
      uint8_t itemsEndPos = plc4c_spi_read_get_pos(io) + _itemsLength;
      while(plc4c_spi_read_get_pos(io) < itemsEndPos) {
        plc4c_modbus_read_write_modbus_pdu_read_file_record_response_item* _value = NULL;
        _res = plc4c_modbus_read_write_modbus_pdu_read_file_record_response_item_parse(io, (void*) &_value);
        if(_res != OK) {
          return _res;
        }
        plc4c_utils_list_insert_head_value(items, _value);
      }
    }
    (*_message)->modbus_pdu_read_file_record_response_items = items;

  } else 
  if((errorFlag == false) && (functionFlag == 0x15) && (response == false)) { /* ModbusPDUWriteFileRecordRequest */
    (*_message)->_type = plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_write_file_record_request;
                    
    // Implicit Field (byteCount) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
    uint8_t byteCount = 0;
    _res = plc4c_spi_read_unsigned_byte(io, 8, (uint8_t*) &byteCount);
    if(_res != OK) {
      return _res;
    }


                    
    // Array field (items)
    plc4c_list* items = NULL;
    plc4c_utils_list_create(&items);
    if(items == NULL) {
      return NO_MEMORY;
    }
    {
      // Length array
      uint8_t _itemsLength = byteCount;
      uint8_t itemsEndPos = plc4c_spi_read_get_pos(io) + _itemsLength;
      while(plc4c_spi_read_get_pos(io) < itemsEndPos) {
        plc4c_modbus_read_write_modbus_pdu_write_file_record_request_item* _value = NULL;
        _res = plc4c_modbus_read_write_modbus_pdu_write_file_record_request_item_parse(io, (void*) &_value);
        if(_res != OK) {
          return _res;
        }
        plc4c_utils_list_insert_head_value(items, _value);
      }
    }
    (*_message)->modbus_pdu_write_file_record_request_items = items;

  } else 
  if((errorFlag == false) && (functionFlag == 0x15) && (response == true)) { /* ModbusPDUWriteFileRecordResponse */
    (*_message)->_type = plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_write_file_record_response;
                    
    // Implicit Field (byteCount) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
    uint8_t byteCount = 0;
    _res = plc4c_spi_read_unsigned_byte(io, 8, (uint8_t*) &byteCount);
    if(_res != OK) {
      return _res;
    }


                    
    // Array field (items)
    plc4c_list* items = NULL;
    plc4c_utils_list_create(&items);
    if(items == NULL) {
      return NO_MEMORY;
    }
    {
      // Length array
      uint8_t _itemsLength = byteCount;
      uint8_t itemsEndPos = plc4c_spi_read_get_pos(io) + _itemsLength;
      while(plc4c_spi_read_get_pos(io) < itemsEndPos) {
        plc4c_modbus_read_write_modbus_pdu_write_file_record_response_item* _value = NULL;
        _res = plc4c_modbus_read_write_modbus_pdu_write_file_record_response_item_parse(io, (void*) &_value);
        if(_res != OK) {
          return _res;
        }
        plc4c_utils_list_insert_head_value(items, _value);
      }
    }
    (*_message)->modbus_pdu_write_file_record_response_items = items;

  } else 
  if((errorFlag == false) && (functionFlag == 0x07) && (response == false)) { /* ModbusPDUReadExceptionStatusRequest */
    (*_message)->_type = plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_read_exception_status_request;
  } else 
  if((errorFlag == false) && (functionFlag == 0x07) && (response == true)) { /* ModbusPDUReadExceptionStatusResponse */
    (*_message)->_type = plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_read_exception_status_response;
                    
    // Simple Field (value)
    uint8_t value = 0;
    _res = plc4c_spi_read_unsigned_byte(io, 8, (uint8_t*) &value);
    if(_res != OK) {
      return _res;
    }
    (*_message)->modbus_pdu_read_exception_status_response_value = value;

  } else 
  if((errorFlag == false) && (functionFlag == 0x08) && (response == false)) { /* ModbusPDUDiagnosticRequest */
    (*_message)->_type = plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_diagnostic_request;
                    
    // Simple Field (subFunction)
    uint16_t subFunction = 0;
    _res = plc4c_spi_read_unsigned_short(io, 16, (uint16_t*) &subFunction);
    if(_res != OK) {
      return _res;
    }
    (*_message)->modbus_pdu_diagnostic_request_sub_function = subFunction;


                    
    // Simple Field (data)
    uint16_t data = 0;
    _res = plc4c_spi_read_unsigned_short(io, 16, (uint16_t*) &data);
    if(_res != OK) {
      return _res;
    }
    (*_message)->modbus_pdu_diagnostic_request_data = data;

  } else 
  if((errorFlag == false) && (functionFlag == 0x08) && (response == true)) { /* ModbusPDUDiagnosticResponse */
    (*_message)->_type = plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_diagnostic_response;
                    
    // Simple Field (subFunction)
    uint16_t subFunction = 0;
    _res = plc4c_spi_read_unsigned_short(io, 16, (uint16_t*) &subFunction);
    if(_res != OK) {
      return _res;
    }
    (*_message)->modbus_pdu_diagnostic_response_sub_function = subFunction;


                    
    // Simple Field (data)
    uint16_t data = 0;
    _res = plc4c_spi_read_unsigned_short(io, 16, (uint16_t*) &data);
    if(_res != OK) {
      return _res;
    }
    (*_message)->modbus_pdu_diagnostic_response_data = data;

  } else 
  if((errorFlag == false) && (functionFlag == 0x0B) && (response == false)) { /* ModbusPDUGetComEventCounterRequest */
    (*_message)->_type = plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_get_com_event_counter_request;
  } else 
  if((errorFlag == false) && (functionFlag == 0x0B) && (response == true)) { /* ModbusPDUGetComEventCounterResponse */
    (*_message)->_type = plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_get_com_event_counter_response;
                    
    // Simple Field (status)
    uint16_t status = 0;
    _res = plc4c_spi_read_unsigned_short(io, 16, (uint16_t*) &status);
    if(_res != OK) {
      return _res;
    }
    (*_message)->modbus_pdu_get_com_event_counter_response_status = status;


                    
    // Simple Field (eventCount)
    uint16_t eventCount = 0;
    _res = plc4c_spi_read_unsigned_short(io, 16, (uint16_t*) &eventCount);
    if(_res != OK) {
      return _res;
    }
    (*_message)->modbus_pdu_get_com_event_counter_response_event_count = eventCount;

  } else 
  if((errorFlag == false) && (functionFlag == 0x0C) && (response == false)) { /* ModbusPDUGetComEventLogRequest */
    (*_message)->_type = plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_get_com_event_log_request;
  } else 
  if((errorFlag == false) && (functionFlag == 0x0C) && (response == true)) { /* ModbusPDUGetComEventLogResponse */
    (*_message)->_type = plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_get_com_event_log_response;
                    
    // Implicit Field (byteCount) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
    uint8_t byteCount = 0;
    _res = plc4c_spi_read_unsigned_byte(io, 8, (uint8_t*) &byteCount);
    if(_res != OK) {
      return _res;
    }


                    
    // Simple Field (status)
    uint16_t status = 0;
    _res = plc4c_spi_read_unsigned_short(io, 16, (uint16_t*) &status);
    if(_res != OK) {
      return _res;
    }
    (*_message)->modbus_pdu_get_com_event_log_response_status = status;


                    
    // Simple Field (eventCount)
    uint16_t eventCount = 0;
    _res = plc4c_spi_read_unsigned_short(io, 16, (uint16_t*) &eventCount);
    if(_res != OK) {
      return _res;
    }
    (*_message)->modbus_pdu_get_com_event_log_response_event_count = eventCount;


                    
    // Simple Field (messageCount)
    uint16_t messageCount = 0;
    _res = plc4c_spi_read_unsigned_short(io, 16, (uint16_t*) &messageCount);
    if(_res != OK) {
      return _res;
    }
    (*_message)->modbus_pdu_get_com_event_log_response_message_count = messageCount;


                    
    // Array field (events)
    plc4c_list* events = NULL;
    plc4c_utils_list_create(&events);
    if(events == NULL) {
      return NO_MEMORY;
    }
    {
      // Count array
      uint8_t itemCount = (byteCount) - (6);
      for(int curItem = 0; curItem < itemCount; curItem++) {
        
        int8_t* _value = malloc(sizeof(int8_t));
        _res = plc4c_spi_read_signed_byte(io, 8, (int8_t*) _value);
        if(_res != OK) {
          return _res;
        }
        plc4c_utils_list_insert_head_value(events, _value);
      }
    }
    (*_message)->modbus_pdu_get_com_event_log_response_events = events;

  } else 
  if((errorFlag == false) && (functionFlag == 0x11) && (response == false)) { /* ModbusPDUReportServerIdRequest */
    (*_message)->_type = plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_report_server_id_request;
  } else 
  if((errorFlag == false) && (functionFlag == 0x11) && (response == true)) { /* ModbusPDUReportServerIdResponse */
    (*_message)->_type = plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_report_server_id_response;
                    
    // Implicit Field (byteCount) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
    uint8_t byteCount = 0;
    _res = plc4c_spi_read_unsigned_byte(io, 8, (uint8_t*) &byteCount);
    if(_res != OK) {
      return _res;
    }


                    
    // Array field (value)
    plc4c_list* value = NULL;
    plc4c_utils_list_create(&value);
    if(value == NULL) {
      return NO_MEMORY;
    }
    {
      // Count array
      uint8_t itemCount = byteCount;
      for(int curItem = 0; curItem < itemCount; curItem++) {
        
        int8_t* _value = malloc(sizeof(int8_t));
        _res = plc4c_spi_read_signed_byte(io, 8, (int8_t*) _value);
        if(_res != OK) {
          return _res;
        }
        plc4c_utils_list_insert_head_value(value, _value);
      }
    }
    (*_message)->modbus_pdu_report_server_id_response_value = value;

  } else 
  if((errorFlag == false) && (functionFlag == 0x2B) && (response == false)) { /* ModbusPDUReadDeviceIdentificationRequest */
    (*_message)->_type = plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_read_device_identification_request;
  } else 
  if((errorFlag == false) && (functionFlag == 0x2B) && (response == true)) { /* ModbusPDUReadDeviceIdentificationResponse */
    (*_message)->_type = plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_read_device_identification_response;
  }

  return OK;
}

plc4c_return_code plc4c_modbus_read_write_modbus_pdu_serialize(plc4c_spi_write_buffer* io, plc4c_modbus_read_write_modbus_pdu* _message) {
  plc4c_return_code _res = OK;
  // Enumerated Discriminator Field (errorFlag)
            // org.apache.plc4x.plugins.codegenerator.language.mspec.model.fields.DefaultDiscriminatorField@4888d1ea
            // org.apache.plc4x.plugins.codegenerator.types.references.DefaultBooleanTypeReference@217b0952
            // errorFlag
            // plc4c_modbus_read_write_modbus_pdu
  plc4c_spi_write_bit(io, plc4c_modbus_read_write_modbus_pdu_get_discriminator(_message->_type).errorFlag);
  // Enumerated Discriminator Field (functionFlag)
            // org.apache.plc4x.plugins.codegenerator.language.mspec.model.fields.DefaultDiscriminatorField@38ba8b45
            // org.apache.plc4x.plugins.codegenerator.types.references.DefaultIntegerTypeReference@41f23499
            // functionFlag
            // plc4c_modbus_read_write_modbus_pdu
  plc4c_spi_write_unsigned_byte(io, 7, plc4c_modbus_read_write_modbus_pdu_get_discriminator(_message->_type).functionFlag);

  // Switch Field (Depending of the current type, serialize the sub-type elements)
  switch(_message->_type) {
    case plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_error: {

      // Enum field (exceptionCode)
      _res = plc4c_spi_write_unsigned_byte(io, 8, _message->modbus_pdu_error_exception_code);
      if(_res != OK) {
        return _res;
      }

      break;
    }
    case plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_read_discrete_inputs_request: {

      // Simple Field (startingAddress)
      _res = plc4c_spi_write_unsigned_short(io, 16, _message->modbus_pdu_read_discrete_inputs_request_starting_address);
      if(_res != OK) {
        return _res;
      }

      // Simple Field (quantity)
      _res = plc4c_spi_write_unsigned_short(io, 16, _message->modbus_pdu_read_discrete_inputs_request_quantity);
      if(_res != OK) {
        return _res;
      }

      break;
    }
    case plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_read_discrete_inputs_response: {

      // Implicit Field (byteCount) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
      _res = plc4c_spi_write_unsigned_byte(io, 8, plc4c_spi_evaluation_helper_count(_message->modbus_pdu_read_discrete_inputs_response_value));
      if(_res != OK) {
        return _res;
      }

      // Array field (value)
      {
        uint8_t itemCount = plc4c_utils_list_size(_message->modbus_pdu_read_discrete_inputs_response_value);
        for(int curItem = 0; curItem < itemCount; curItem++) {

          int8_t* _value = (int8_t*) plc4c_utils_list_get_value(_message->modbus_pdu_read_discrete_inputs_response_value, curItem);
          plc4c_spi_write_signed_byte(io, 8, *_value);
        }
      }

      break;
    }
    case plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_read_coils_request: {

      // Simple Field (startingAddress)
      _res = plc4c_spi_write_unsigned_short(io, 16, _message->modbus_pdu_read_coils_request_starting_address);
      if(_res != OK) {
        return _res;
      }

      // Simple Field (quantity)
      _res = plc4c_spi_write_unsigned_short(io, 16, _message->modbus_pdu_read_coils_request_quantity);
      if(_res != OK) {
        return _res;
      }

      break;
    }
    case plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_read_coils_response: {

      // Implicit Field (byteCount) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
      _res = plc4c_spi_write_unsigned_byte(io, 8, plc4c_spi_evaluation_helper_count(_message->modbus_pdu_read_coils_response_value));
      if(_res != OK) {
        return _res;
      }

      // Array field (value)
      {
        uint8_t itemCount = plc4c_utils_list_size(_message->modbus_pdu_read_coils_response_value);
        for(int curItem = 0; curItem < itemCount; curItem++) {

          int8_t* _value = (int8_t*) plc4c_utils_list_get_value(_message->modbus_pdu_read_coils_response_value, curItem);
          plc4c_spi_write_signed_byte(io, 8, *_value);
        }
      }

      break;
    }
    case plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_write_single_coil_request: {

      // Simple Field (address)
      _res = plc4c_spi_write_unsigned_short(io, 16, _message->modbus_pdu_write_single_coil_request_address);
      if(_res != OK) {
        return _res;
      }

      // Simple Field (value)
      _res = plc4c_spi_write_unsigned_short(io, 16, _message->modbus_pdu_write_single_coil_request_value);
      if(_res != OK) {
        return _res;
      }

      break;
    }
    case plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_write_single_coil_response: {

      // Simple Field (address)
      _res = plc4c_spi_write_unsigned_short(io, 16, _message->modbus_pdu_write_single_coil_response_address);
      if(_res != OK) {
        return _res;
      }

      // Simple Field (value)
      _res = plc4c_spi_write_unsigned_short(io, 16, _message->modbus_pdu_write_single_coil_response_value);
      if(_res != OK) {
        return _res;
      }

      break;
    }
    case plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_write_multiple_coils_request: {

      // Simple Field (startingAddress)
      _res = plc4c_spi_write_unsigned_short(io, 16, _message->modbus_pdu_write_multiple_coils_request_starting_address);
      if(_res != OK) {
        return _res;
      }

      // Simple Field (quantity)
      _res = plc4c_spi_write_unsigned_short(io, 16, _message->modbus_pdu_write_multiple_coils_request_quantity);
      if(_res != OK) {
        return _res;
      }

      // Implicit Field (byteCount) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
      _res = plc4c_spi_write_unsigned_byte(io, 8, plc4c_spi_evaluation_helper_count(_message->modbus_pdu_write_multiple_coils_request_value));
      if(_res != OK) {
        return _res;
      }

      // Array field (value)
      {
        uint8_t itemCount = plc4c_utils_list_size(_message->modbus_pdu_write_multiple_coils_request_value);
        for(int curItem = 0; curItem < itemCount; curItem++) {

          int8_t* _value = (int8_t*) plc4c_utils_list_get_value(_message->modbus_pdu_write_multiple_coils_request_value, curItem);
          plc4c_spi_write_signed_byte(io, 8, *_value);
        }
      }

      break;
    }
    case plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_write_multiple_coils_response: {

      // Simple Field (startingAddress)
      _res = plc4c_spi_write_unsigned_short(io, 16, _message->modbus_pdu_write_multiple_coils_response_starting_address);
      if(_res != OK) {
        return _res;
      }

      // Simple Field (quantity)
      _res = plc4c_spi_write_unsigned_short(io, 16, _message->modbus_pdu_write_multiple_coils_response_quantity);
      if(_res != OK) {
        return _res;
      }

      break;
    }
    case plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_read_input_registers_request: {

      // Simple Field (startingAddress)
      _res = plc4c_spi_write_unsigned_short(io, 16, _message->modbus_pdu_read_input_registers_request_starting_address);
      if(_res != OK) {
        return _res;
      }

      // Simple Field (quantity)
      _res = plc4c_spi_write_unsigned_short(io, 16, _message->modbus_pdu_read_input_registers_request_quantity);
      if(_res != OK) {
        return _res;
      }

      break;
    }
    case plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_read_input_registers_response: {

      // Implicit Field (byteCount) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
      _res = plc4c_spi_write_unsigned_byte(io, 8, plc4c_spi_evaluation_helper_count(_message->modbus_pdu_read_input_registers_response_value));
      if(_res != OK) {
        return _res;
      }

      // Array field (value)
      {
        uint8_t itemCount = plc4c_utils_list_size(_message->modbus_pdu_read_input_registers_response_value);
        for(int curItem = 0; curItem < itemCount; curItem++) {

          int8_t* _value = (int8_t*) plc4c_utils_list_get_value(_message->modbus_pdu_read_input_registers_response_value, curItem);
          plc4c_spi_write_signed_byte(io, 8, *_value);
        }
      }

      break;
    }
    case plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_read_holding_registers_request: {

      // Simple Field (startingAddress)
      _res = plc4c_spi_write_unsigned_short(io, 16, _message->modbus_pdu_read_holding_registers_request_starting_address);
      if(_res != OK) {
        return _res;
      }

      // Simple Field (quantity)
      _res = plc4c_spi_write_unsigned_short(io, 16, _message->modbus_pdu_read_holding_registers_request_quantity);
      if(_res != OK) {
        return _res;
      }

      break;
    }
    case plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_read_holding_registers_response: {

      // Implicit Field (byteCount) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
      _res = plc4c_spi_write_unsigned_byte(io, 8, plc4c_spi_evaluation_helper_count(_message->modbus_pdu_read_holding_registers_response_value));
      if(_res != OK) {
        return _res;
      }

      // Array field (value)
      {
        uint8_t itemCount = plc4c_utils_list_size(_message->modbus_pdu_read_holding_registers_response_value);
        for(int curItem = 0; curItem < itemCount; curItem++) {

          int8_t* _value = (int8_t*) plc4c_utils_list_get_value(_message->modbus_pdu_read_holding_registers_response_value, curItem);
          plc4c_spi_write_signed_byte(io, 8, *_value);
        }
      }

      break;
    }
    case plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_write_single_register_request: {

      // Simple Field (address)
      _res = plc4c_spi_write_unsigned_short(io, 16, _message->modbus_pdu_write_single_register_request_address);
      if(_res != OK) {
        return _res;
      }

      // Simple Field (value)
      _res = plc4c_spi_write_unsigned_short(io, 16, _message->modbus_pdu_write_single_register_request_value);
      if(_res != OK) {
        return _res;
      }

      break;
    }
    case plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_write_single_register_response: {

      // Simple Field (address)
      _res = plc4c_spi_write_unsigned_short(io, 16, _message->modbus_pdu_write_single_register_response_address);
      if(_res != OK) {
        return _res;
      }

      // Simple Field (value)
      _res = plc4c_spi_write_unsigned_short(io, 16, _message->modbus_pdu_write_single_register_response_value);
      if(_res != OK) {
        return _res;
      }

      break;
    }
    case plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_write_multiple_holding_registers_request: {

      // Simple Field (startingAddress)
      _res = plc4c_spi_write_unsigned_short(io, 16, _message->modbus_pdu_write_multiple_holding_registers_request_starting_address);
      if(_res != OK) {
        return _res;
      }

      // Simple Field (quantity)
      _res = plc4c_spi_write_unsigned_short(io, 16, _message->modbus_pdu_write_multiple_holding_registers_request_quantity);
      if(_res != OK) {
        return _res;
      }

      // Implicit Field (byteCount) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
      _res = plc4c_spi_write_unsigned_byte(io, 8, plc4c_spi_evaluation_helper_count(_message->modbus_pdu_write_multiple_holding_registers_request_value));
      if(_res != OK) {
        return _res;
      }

      // Array field (value)
      {
        uint8_t itemCount = plc4c_utils_list_size(_message->modbus_pdu_write_multiple_holding_registers_request_value);
        for(int curItem = 0; curItem < itemCount; curItem++) {

          int8_t* _value = (int8_t*) plc4c_utils_list_get_value(_message->modbus_pdu_write_multiple_holding_registers_request_value, curItem);
          plc4c_spi_write_signed_byte(io, 8, *_value);
        }
      }

      break;
    }
    case plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_write_multiple_holding_registers_response: {

      // Simple Field (startingAddress)
      _res = plc4c_spi_write_unsigned_short(io, 16, _message->modbus_pdu_write_multiple_holding_registers_response_starting_address);
      if(_res != OK) {
        return _res;
      }

      // Simple Field (quantity)
      _res = plc4c_spi_write_unsigned_short(io, 16, _message->modbus_pdu_write_multiple_holding_registers_response_quantity);
      if(_res != OK) {
        return _res;
      }

      break;
    }
    case plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_read_write_multiple_holding_registers_request: {

      // Simple Field (readStartingAddress)
      _res = plc4c_spi_write_unsigned_short(io, 16, _message->modbus_pdu_read_write_multiple_holding_registers_request_read_starting_address);
      if(_res != OK) {
        return _res;
      }

      // Simple Field (readQuantity)
      _res = plc4c_spi_write_unsigned_short(io, 16, _message->modbus_pdu_read_write_multiple_holding_registers_request_read_quantity);
      if(_res != OK) {
        return _res;
      }

      // Simple Field (writeStartingAddress)
      _res = plc4c_spi_write_unsigned_short(io, 16, _message->modbus_pdu_read_write_multiple_holding_registers_request_write_starting_address);
      if(_res != OK) {
        return _res;
      }

      // Simple Field (writeQuantity)
      _res = plc4c_spi_write_unsigned_short(io, 16, _message->modbus_pdu_read_write_multiple_holding_registers_request_write_quantity);
      if(_res != OK) {
        return _res;
      }

      // Implicit Field (byteCount) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
      _res = plc4c_spi_write_unsigned_byte(io, 8, plc4c_spi_evaluation_helper_count(_message->modbus_pdu_read_write_multiple_holding_registers_request_value));
      if(_res != OK) {
        return _res;
      }

      // Array field (value)
      {
        uint8_t itemCount = plc4c_utils_list_size(_message->modbus_pdu_read_write_multiple_holding_registers_request_value);
        for(int curItem = 0; curItem < itemCount; curItem++) {

          int8_t* _value = (int8_t*) plc4c_utils_list_get_value(_message->modbus_pdu_read_write_multiple_holding_registers_request_value, curItem);
          plc4c_spi_write_signed_byte(io, 8, *_value);
        }
      }

      break;
    }
    case plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_read_write_multiple_holding_registers_response: {

      // Implicit Field (byteCount) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
      _res = plc4c_spi_write_unsigned_byte(io, 8, plc4c_spi_evaluation_helper_count(_message->modbus_pdu_read_write_multiple_holding_registers_response_value));
      if(_res != OK) {
        return _res;
      }

      // Array field (value)
      {
        uint8_t itemCount = plc4c_utils_list_size(_message->modbus_pdu_read_write_multiple_holding_registers_response_value);
        for(int curItem = 0; curItem < itemCount; curItem++) {

          int8_t* _value = (int8_t*) plc4c_utils_list_get_value(_message->modbus_pdu_read_write_multiple_holding_registers_response_value, curItem);
          plc4c_spi_write_signed_byte(io, 8, *_value);
        }
      }

      break;
    }
    case plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_mask_write_holding_register_request: {

      // Simple Field (referenceAddress)
      _res = plc4c_spi_write_unsigned_short(io, 16, _message->modbus_pdu_mask_write_holding_register_request_reference_address);
      if(_res != OK) {
        return _res;
      }

      // Simple Field (andMask)
      _res = plc4c_spi_write_unsigned_short(io, 16, _message->modbus_pdu_mask_write_holding_register_request_and_mask);
      if(_res != OK) {
        return _res;
      }

      // Simple Field (orMask)
      _res = plc4c_spi_write_unsigned_short(io, 16, _message->modbus_pdu_mask_write_holding_register_request_or_mask);
      if(_res != OK) {
        return _res;
      }

      break;
    }
    case plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_mask_write_holding_register_response: {

      // Simple Field (referenceAddress)
      _res = plc4c_spi_write_unsigned_short(io, 16, _message->modbus_pdu_mask_write_holding_register_response_reference_address);
      if(_res != OK) {
        return _res;
      }

      // Simple Field (andMask)
      _res = plc4c_spi_write_unsigned_short(io, 16, _message->modbus_pdu_mask_write_holding_register_response_and_mask);
      if(_res != OK) {
        return _res;
      }

      // Simple Field (orMask)
      _res = plc4c_spi_write_unsigned_short(io, 16, _message->modbus_pdu_mask_write_holding_register_response_or_mask);
      if(_res != OK) {
        return _res;
      }

      break;
    }
    case plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_read_fifo_queue_request: {

      // Simple Field (fifoPointerAddress)
      _res = plc4c_spi_write_unsigned_short(io, 16, _message->modbus_pdu_read_fifo_queue_request_fifo_pointer_address);
      if(_res != OK) {
        return _res;
      }

      break;
    }
    case plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_read_fifo_queue_response: {

      // Implicit Field (byteCount) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
      _res = plc4c_spi_write_unsigned_short(io, 16, (((plc4c_spi_evaluation_helper_count(_message->modbus_pdu_read_fifo_queue_response_fifo_value)) * (2))) + (2));
      if(_res != OK) {
        return _res;
      }

      // Implicit Field (fifoCount) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
      _res = plc4c_spi_write_unsigned_short(io, 16, (((plc4c_spi_evaluation_helper_count(_message->modbus_pdu_read_fifo_queue_response_fifo_value)) * (2))) / (2));
      if(_res != OK) {
        return _res;
      }

      // Array field (fifoValue)
      {
        uint8_t itemCount = plc4c_utils_list_size(_message->modbus_pdu_read_fifo_queue_response_fifo_value);
        for(int curItem = 0; curItem < itemCount; curItem++) {

          uint16_t* _value = (uint16_t*) plc4c_utils_list_get_value(_message->modbus_pdu_read_fifo_queue_response_fifo_value, curItem);
          plc4c_spi_write_unsigned_short(io, 16, *_value);
        }
      }

      break;
    }
    case plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_read_file_record_request: {

      // Implicit Field (byteCount) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
      _res = plc4c_spi_write_unsigned_byte(io, 8, plc4c_spi_evaluation_helper_array_size_in_bytes(_message->modbus_pdu_read_file_record_request_items));
      if(_res != OK) {
        return _res;
      }

      // Array field (items)
      {
        uint8_t itemCount = plc4c_utils_list_size(_message->modbus_pdu_read_file_record_request_items);
        for(int curItem = 0; curItem < itemCount; curItem++) {
          bool lastItem = curItem == (itemCount - 1);
          plc4c_modbus_read_write_modbus_pdu_read_file_record_request_item* _value = (plc4c_modbus_read_write_modbus_pdu_read_file_record_request_item*) plc4c_utils_list_get_value(_message->modbus_pdu_read_file_record_request_items, curItem);
          _res = plc4c_modbus_read_write_modbus_pdu_read_file_record_request_item_serialize(io, (void*) _value);
          if(_res != OK) {
            return _res;
          }
        }
      }

      break;
    }
    case plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_read_file_record_response: {

      // Implicit Field (byteCount) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
      _res = plc4c_spi_write_unsigned_byte(io, 8, plc4c_spi_evaluation_helper_array_size_in_bytes(_message->modbus_pdu_read_file_record_response_items));
      if(_res != OK) {
        return _res;
      }

      // Array field (items)
      {
        uint8_t itemCount = plc4c_utils_list_size(_message->modbus_pdu_read_file_record_response_items);
        for(int curItem = 0; curItem < itemCount; curItem++) {
          bool lastItem = curItem == (itemCount - 1);
          plc4c_modbus_read_write_modbus_pdu_read_file_record_response_item* _value = (plc4c_modbus_read_write_modbus_pdu_read_file_record_response_item*) plc4c_utils_list_get_value(_message->modbus_pdu_read_file_record_response_items, curItem);
          _res = plc4c_modbus_read_write_modbus_pdu_read_file_record_response_item_serialize(io, (void*) _value);
          if(_res != OK) {
            return _res;
          }
        }
      }

      break;
    }
    case plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_write_file_record_request: {

      // Implicit Field (byteCount) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
      _res = plc4c_spi_write_unsigned_byte(io, 8, plc4c_spi_evaluation_helper_array_size_in_bytes(_message->modbus_pdu_write_file_record_request_items));
      if(_res != OK) {
        return _res;
      }

      // Array field (items)
      {
        uint8_t itemCount = plc4c_utils_list_size(_message->modbus_pdu_write_file_record_request_items);
        for(int curItem = 0; curItem < itemCount; curItem++) {
          bool lastItem = curItem == (itemCount - 1);
          plc4c_modbus_read_write_modbus_pdu_write_file_record_request_item* _value = (plc4c_modbus_read_write_modbus_pdu_write_file_record_request_item*) plc4c_utils_list_get_value(_message->modbus_pdu_write_file_record_request_items, curItem);
          _res = plc4c_modbus_read_write_modbus_pdu_write_file_record_request_item_serialize(io, (void*) _value);
          if(_res != OK) {
            return _res;
          }
        }
      }

      break;
    }
    case plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_write_file_record_response: {

      // Implicit Field (byteCount) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
      _res = plc4c_spi_write_unsigned_byte(io, 8, plc4c_spi_evaluation_helper_array_size_in_bytes(_message->modbus_pdu_write_file_record_response_items));
      if(_res != OK) {
        return _res;
      }

      // Array field (items)
      {
        uint8_t itemCount = plc4c_utils_list_size(_message->modbus_pdu_write_file_record_response_items);
        for(int curItem = 0; curItem < itemCount; curItem++) {
          bool lastItem = curItem == (itemCount - 1);
          plc4c_modbus_read_write_modbus_pdu_write_file_record_response_item* _value = (plc4c_modbus_read_write_modbus_pdu_write_file_record_response_item*) plc4c_utils_list_get_value(_message->modbus_pdu_write_file_record_response_items, curItem);
          _res = plc4c_modbus_read_write_modbus_pdu_write_file_record_response_item_serialize(io, (void*) _value);
          if(_res != OK) {
            return _res;
          }
        }
      }

      break;
    }
    case plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_read_exception_status_request: {

      break;
    }
    case plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_read_exception_status_response: {

      // Simple Field (value)
      _res = plc4c_spi_write_unsigned_byte(io, 8, _message->modbus_pdu_read_exception_status_response_value);
      if(_res != OK) {
        return _res;
      }

      break;
    }
    case plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_diagnostic_request: {

      // Simple Field (subFunction)
      _res = plc4c_spi_write_unsigned_short(io, 16, _message->modbus_pdu_diagnostic_request_sub_function);
      if(_res != OK) {
        return _res;
      }

      // Simple Field (data)
      _res = plc4c_spi_write_unsigned_short(io, 16, _message->modbus_pdu_diagnostic_request_data);
      if(_res != OK) {
        return _res;
      }

      break;
    }
    case plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_diagnostic_response: {

      // Simple Field (subFunction)
      _res = plc4c_spi_write_unsigned_short(io, 16, _message->modbus_pdu_diagnostic_response_sub_function);
      if(_res != OK) {
        return _res;
      }

      // Simple Field (data)
      _res = plc4c_spi_write_unsigned_short(io, 16, _message->modbus_pdu_diagnostic_response_data);
      if(_res != OK) {
        return _res;
      }

      break;
    }
    case plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_get_com_event_counter_request: {

      break;
    }
    case plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_get_com_event_counter_response: {

      // Simple Field (status)
      _res = plc4c_spi_write_unsigned_short(io, 16, _message->modbus_pdu_get_com_event_counter_response_status);
      if(_res != OK) {
        return _res;
      }

      // Simple Field (eventCount)
      _res = plc4c_spi_write_unsigned_short(io, 16, _message->modbus_pdu_get_com_event_counter_response_event_count);
      if(_res != OK) {
        return _res;
      }

      break;
    }
    case plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_get_com_event_log_request: {

      break;
    }
    case plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_get_com_event_log_response: {

      // Implicit Field (byteCount) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
      _res = plc4c_spi_write_unsigned_byte(io, 8, (plc4c_spi_evaluation_helper_count(_message->modbus_pdu_get_com_event_log_response_events)) + (6));
      if(_res != OK) {
        return _res;
      }

      // Simple Field (status)
      _res = plc4c_spi_write_unsigned_short(io, 16, _message->modbus_pdu_get_com_event_log_response_status);
      if(_res != OK) {
        return _res;
      }

      // Simple Field (eventCount)
      _res = plc4c_spi_write_unsigned_short(io, 16, _message->modbus_pdu_get_com_event_log_response_event_count);
      if(_res != OK) {
        return _res;
      }

      // Simple Field (messageCount)
      _res = plc4c_spi_write_unsigned_short(io, 16, _message->modbus_pdu_get_com_event_log_response_message_count);
      if(_res != OK) {
        return _res;
      }

      // Array field (events)
      {
        uint8_t itemCount = plc4c_utils_list_size(_message->modbus_pdu_get_com_event_log_response_events);
        for(int curItem = 0; curItem < itemCount; curItem++) {

          int8_t* _value = (int8_t*) plc4c_utils_list_get_value(_message->modbus_pdu_get_com_event_log_response_events, curItem);
          plc4c_spi_write_signed_byte(io, 8, *_value);
        }
      }

      break;
    }
    case plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_report_server_id_request: {

      break;
    }
    case plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_report_server_id_response: {

      // Implicit Field (byteCount) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
      _res = plc4c_spi_write_unsigned_byte(io, 8, plc4c_spi_evaluation_helper_count(_message->modbus_pdu_report_server_id_response_value));
      if(_res != OK) {
        return _res;
      }

      // Array field (value)
      {
        uint8_t itemCount = plc4c_utils_list_size(_message->modbus_pdu_report_server_id_response_value);
        for(int curItem = 0; curItem < itemCount; curItem++) {

          int8_t* _value = (int8_t*) plc4c_utils_list_get_value(_message->modbus_pdu_report_server_id_response_value, curItem);
          plc4c_spi_write_signed_byte(io, 8, *_value);
        }
      }

      break;
    }
    case plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_read_device_identification_request: {

      break;
    }
    case plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_read_device_identification_response: {

      break;
    }
  }

  return OK;
}

uint16_t plc4c_modbus_read_write_modbus_pdu_length_in_bytes(plc4c_modbus_read_write_modbus_pdu* _message) {
  return plc4c_modbus_read_write_modbus_pdu_length_in_bits(_message) / 8;
}

uint16_t plc4c_modbus_read_write_modbus_pdu_length_in_bits(plc4c_modbus_read_write_modbus_pdu* _message) {
  uint16_t lengthInBits = 0;

        // Discriminator Field (errorFlag)
                lengthInBits += 1;

        // Discriminator Field (functionFlag)
                lengthInBits += 7;

  // Depending of the current type, add the length of sub-type elements ...
  switch(_message->_type) {
    case plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_error: {

      // Enum Field (exceptionCode)
      lengthInBits += 8;

      break;
    }
    case plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_read_discrete_inputs_request: {

      // Simple field (startingAddress)
      lengthInBits += 16;


      // Simple field (quantity)
      lengthInBits += 16;

      break;
    }
    case plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_read_discrete_inputs_response: {

      // Implicit Field (byteCount)
      lengthInBits += 8;


      // Array field
      lengthInBits += 8 * plc4c_utils_list_size(_message->modbus_pdu_read_discrete_inputs_response_value);

      break;
    }
    case plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_read_coils_request: {

      // Simple field (startingAddress)
      lengthInBits += 16;


      // Simple field (quantity)
      lengthInBits += 16;

      break;
    }
    case plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_read_coils_response: {

      // Implicit Field (byteCount)
      lengthInBits += 8;


      // Array field
      lengthInBits += 8 * plc4c_utils_list_size(_message->modbus_pdu_read_coils_response_value);

      break;
    }
    case plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_write_single_coil_request: {

      // Simple field (address)
      lengthInBits += 16;


      // Simple field (value)
      lengthInBits += 16;

      break;
    }
    case plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_write_single_coil_response: {

      // Simple field (address)
      lengthInBits += 16;


      // Simple field (value)
      lengthInBits += 16;

      break;
    }
    case plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_write_multiple_coils_request: {

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
    case plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_write_multiple_coils_response: {

      // Simple field (startingAddress)
      lengthInBits += 16;


      // Simple field (quantity)
      lengthInBits += 16;

      break;
    }
    case plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_read_input_registers_request: {

      // Simple field (startingAddress)
      lengthInBits += 16;


      // Simple field (quantity)
      lengthInBits += 16;

      break;
    }
    case plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_read_input_registers_response: {

      // Implicit Field (byteCount)
      lengthInBits += 8;


      // Array field
      lengthInBits += 8 * plc4c_utils_list_size(_message->modbus_pdu_read_input_registers_response_value);

      break;
    }
    case plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_read_holding_registers_request: {

      // Simple field (startingAddress)
      lengthInBits += 16;


      // Simple field (quantity)
      lengthInBits += 16;

      break;
    }
    case plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_read_holding_registers_response: {

      // Implicit Field (byteCount)
      lengthInBits += 8;


      // Array field
      lengthInBits += 8 * plc4c_utils_list_size(_message->modbus_pdu_read_holding_registers_response_value);

      break;
    }
    case plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_write_single_register_request: {

      // Simple field (address)
      lengthInBits += 16;


      // Simple field (value)
      lengthInBits += 16;

      break;
    }
    case plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_write_single_register_response: {

      // Simple field (address)
      lengthInBits += 16;


      // Simple field (value)
      lengthInBits += 16;

      break;
    }
    case plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_write_multiple_holding_registers_request: {

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
    case plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_write_multiple_holding_registers_response: {

      // Simple field (startingAddress)
      lengthInBits += 16;


      // Simple field (quantity)
      lengthInBits += 16;

      break;
    }
    case plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_read_write_multiple_holding_registers_request: {

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
    case plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_read_write_multiple_holding_registers_response: {

      // Implicit Field (byteCount)
      lengthInBits += 8;


      // Array field
      lengthInBits += 8 * plc4c_utils_list_size(_message->modbus_pdu_read_write_multiple_holding_registers_response_value);

      break;
    }
    case plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_mask_write_holding_register_request: {

      // Simple field (referenceAddress)
      lengthInBits += 16;


      // Simple field (andMask)
      lengthInBits += 16;


      // Simple field (orMask)
      lengthInBits += 16;

      break;
    }
    case plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_mask_write_holding_register_response: {

      // Simple field (referenceAddress)
      lengthInBits += 16;


      // Simple field (andMask)
      lengthInBits += 16;


      // Simple field (orMask)
      lengthInBits += 16;

      break;
    }
    case plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_read_fifo_queue_request: {

      // Simple field (fifoPointerAddress)
      lengthInBits += 16;

      break;
    }
    case plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_read_fifo_queue_response: {

      // Implicit Field (byteCount)
      lengthInBits += 16;


      // Implicit Field (fifoCount)
      lengthInBits += 16;


      // Array field
      lengthInBits += 16 * plc4c_utils_list_size(_message->modbus_pdu_read_fifo_queue_response_fifo_value);

      break;
    }
    case plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_read_file_record_request: {

      // Implicit Field (byteCount)
      lengthInBits += 8;


      // Array field
      if(_message->modbus_pdu_read_file_record_request_items != NULL) {
        plc4c_list_element* curElement = _message->modbus_pdu_read_file_record_request_items->tail;
        while (curElement != NULL) {
          lengthInBits += plc4c_modbus_read_write_modbus_pdu_read_file_record_request_item_length_in_bits((plc4c_modbus_read_write_modbus_pdu_read_file_record_request_item*) curElement->value);
          curElement = curElement->next;
        }
      }

      break;
    }
    case plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_read_file_record_response: {

      // Implicit Field (byteCount)
      lengthInBits += 8;


      // Array field
      if(_message->modbus_pdu_read_file_record_response_items != NULL) {
        plc4c_list_element* curElement = _message->modbus_pdu_read_file_record_response_items->tail;
        while (curElement != NULL) {
          lengthInBits += plc4c_modbus_read_write_modbus_pdu_read_file_record_response_item_length_in_bits((plc4c_modbus_read_write_modbus_pdu_read_file_record_response_item*) curElement->value);
          curElement = curElement->next;
        }
      }

      break;
    }
    case plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_write_file_record_request: {

      // Implicit Field (byteCount)
      lengthInBits += 8;


      // Array field
      if(_message->modbus_pdu_write_file_record_request_items != NULL) {
        plc4c_list_element* curElement = _message->modbus_pdu_write_file_record_request_items->tail;
        while (curElement != NULL) {
          lengthInBits += plc4c_modbus_read_write_modbus_pdu_write_file_record_request_item_length_in_bits((plc4c_modbus_read_write_modbus_pdu_write_file_record_request_item*) curElement->value);
          curElement = curElement->next;
        }
      }

      break;
    }
    case plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_write_file_record_response: {

      // Implicit Field (byteCount)
      lengthInBits += 8;


      // Array field
      if(_message->modbus_pdu_write_file_record_response_items != NULL) {
        plc4c_list_element* curElement = _message->modbus_pdu_write_file_record_response_items->tail;
        while (curElement != NULL) {
          lengthInBits += plc4c_modbus_read_write_modbus_pdu_write_file_record_response_item_length_in_bits((plc4c_modbus_read_write_modbus_pdu_write_file_record_response_item*) curElement->value);
          curElement = curElement->next;
        }
      }

      break;
    }
    case plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_read_exception_status_request: {

      break;
    }
    case plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_read_exception_status_response: {

      // Simple field (value)
      lengthInBits += 8;

      break;
    }
    case plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_diagnostic_request: {

      // Simple field (subFunction)
      lengthInBits += 16;


      // Simple field (data)
      lengthInBits += 16;

      break;
    }
    case plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_diagnostic_response: {

      // Simple field (subFunction)
      lengthInBits += 16;


      // Simple field (data)
      lengthInBits += 16;

      break;
    }
    case plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_get_com_event_counter_request: {

      break;
    }
    case plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_get_com_event_counter_response: {

      // Simple field (status)
      lengthInBits += 16;


      // Simple field (eventCount)
      lengthInBits += 16;

      break;
    }
    case plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_get_com_event_log_request: {

      break;
    }
    case plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_get_com_event_log_response: {

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
    case plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_report_server_id_request: {

      break;
    }
    case plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_report_server_id_response: {

      // Implicit Field (byteCount)
      lengthInBits += 8;


      // Array field
      lengthInBits += 8 * plc4c_utils_list_size(_message->modbus_pdu_report_server_id_response_value);

      break;
    }
    case plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_read_device_identification_request: {

      break;
    }
    case plc4c_modbus_read_write_modbus_pdu_type_plc4c_modbus_read_write_modbus_pdu_read_device_identification_response: {

      break;
    }
  }

  return lengthInBits;
}

