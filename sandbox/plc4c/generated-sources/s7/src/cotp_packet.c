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
#include "cotp_packet.h"

// Array of discriminator values that match the enum type constants.
// (The order is identical to the enum constants so we can use the
// enum constant to directly access a given types discriminator values)
const plc4c_s7_read_write_cotp_packet_discriminator plc4c_s7_read_write_cotp_packet_discriminators[] = {
  {/* s7_read_write_cotp_packet_connection_request */
   .tpduCode = 0xE0},
  {/* s7_read_write_cotp_packet_connection_response */
   .tpduCode = 0xD0},
  {/* s7_read_write_cotp_packet_data */
   .tpduCode = 0xF0},
  {/* s7_read_write_cotp_packet_disconnect_request */
   .tpduCode = 0x80},
  {/* s7_read_write_cotp_packet_disconnect_response */
   .tpduCode = 0xC0},
  {/* s7_read_write_cotp_packet_tpdu_error */
   .tpduCode = 0x70}
};

// Function returning the discriminator values for a given type constant.
plc4c_s7_read_write_cotp_packet_discriminator plc4c_s7_read_write_cotp_packet_get_discriminator(plc4c_s7_read_write_cotp_packet_type type) {
  return plc4c_s7_read_write_cotp_packet_discriminators[type];
}

// Parse function.
plc4c_return_code plc4c_s7_read_write_cotp_packet_parse(plc4c_spi_read_buffer* buf, uint16_t cotpLen, plc4c_s7_read_write_cotp_packet** _message) {
  uint16_t startPos = plc4c_spi_read_get_pos(buf);
  uint16_t curPos;

  // Pointer to the parsed data structure.
  (*_message) = malloc(sizeof(plc4c_s7_read_write_cotp_packet));
  if(*_message == NULL) {
    return NO_MEMORY;
  }


  // Implicit Field (headerLength) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
  uint8_t headerLength = plc4c_spi_read_unsigned_short(buf, 8);

  // Discriminator Field (tpduCode) (Used as input to a switch field)
  uint8_t tpduCode = plc4c_spi_read_unsigned_short(buf, 8);

  // Switch Field (Depending on the discriminator values, passes the instantiation to a sub-type)
  if(tpduCode == 0xF0) { /* COTPPacketData */
                    
    // Simple Field (eot)
    bool eot = plc4c_spi_read_bit(buf);
    (*_message)->cotp_packet_data_eot = eot;


                    
    // Simple Field (tpduRef)
    unsigned int tpduRef = plc4c_spi_read_unsigned_short(buf, 7);
    (*_message)->cotp_packet_data_tpdu_ref = tpduRef;

  } else 
  if(tpduCode == 0xE0) { /* COTPPacketConnectionRequest */
                    
    // Simple Field (destinationReference)
    uint16_t destinationReference = plc4c_spi_read_unsigned_int(buf, 16);
    (*_message)->cotp_packet_connection_request_destination_reference = destinationReference;


                    
    // Simple Field (sourceReference)
    uint16_t sourceReference = plc4c_spi_read_unsigned_int(buf, 16);
    (*_message)->cotp_packet_connection_request_source_reference = sourceReference;


                    
    // Enum field (protocolClass)
    plc4c_s7_read_write_cotp_protocol_class protocolClass = plc4c_spi_read_byte(buf, 8);
    (*_message)->cotp_packet_connection_request_protocol_class = protocolClass;

  } else 
  if(tpduCode == 0xD0) { /* COTPPacketConnectionResponse */
                    
    // Simple Field (destinationReference)
    uint16_t destinationReference = plc4c_spi_read_unsigned_int(buf, 16);
    (*_message)->cotp_packet_connection_response_destination_reference = destinationReference;


                    
    // Simple Field (sourceReference)
    uint16_t sourceReference = plc4c_spi_read_unsigned_int(buf, 16);
    (*_message)->cotp_packet_connection_response_source_reference = sourceReference;


                    
    // Enum field (protocolClass)
    plc4c_s7_read_write_cotp_protocol_class protocolClass = plc4c_spi_read_byte(buf, 8);
    (*_message)->cotp_packet_connection_response_protocol_class = protocolClass;

  } else 
  if(tpduCode == 0x80) { /* COTPPacketDisconnectRequest */
                    
    // Simple Field (destinationReference)
    uint16_t destinationReference = plc4c_spi_read_unsigned_int(buf, 16);
    (*_message)->cotp_packet_disconnect_request_destination_reference = destinationReference;


                    
    // Simple Field (sourceReference)
    uint16_t sourceReference = plc4c_spi_read_unsigned_int(buf, 16);
    (*_message)->cotp_packet_disconnect_request_source_reference = sourceReference;


                    
    // Enum field (protocolClass)
    plc4c_s7_read_write_cotp_protocol_class protocolClass = plc4c_spi_read_byte(buf, 8);
    (*_message)->cotp_packet_disconnect_request_protocol_class = protocolClass;

  } else 
  if(tpduCode == 0xC0) { /* COTPPacketDisconnectResponse */
                    
    // Simple Field (destinationReference)
    uint16_t destinationReference = plc4c_spi_read_unsigned_int(buf, 16);
    (*_message)->cotp_packet_disconnect_response_destination_reference = destinationReference;


                    
    // Simple Field (sourceReference)
    uint16_t sourceReference = plc4c_spi_read_unsigned_int(buf, 16);
    (*_message)->cotp_packet_disconnect_response_source_reference = sourceReference;

  } else 
  if(tpduCode == 0x70) { /* COTPPacketTpduError */
                    
    // Simple Field (destinationReference)
    uint16_t destinationReference = plc4c_spi_read_unsigned_int(buf, 16);
    (*_message)->cotp_packet_tpdu_error_destination_reference = destinationReference;


                    
    // Simple Field (rejectCause)
    uint8_t rejectCause = plc4c_spi_read_unsigned_short(buf, 8);
    (*_message)->cotp_packet_tpdu_error_reject_cause = rejectCause;

  }

  // Array field (parameters)
  curPos = plc4c_spi_read_get_pos(buf) - startPos;
  plc4c_list* parameters = malloc(sizeof(plc4c_list));
  if(parameters == NULL) {
    return NO_MEMORY;
  }
  {
    // Length array
    uint8_t _parametersLength = (((headerLength) + (1))) - (curPos);
    uint8_t parametersEndPos = plc4c_spi_read_get_pos(buf) + _parametersLength;
    while(plc4c_spi_read_get_pos(buf) < parametersEndPos) {
      plc4c_list* _value = NULL;
      plc4c_return_code _res = plc4c_s7_read_write_cotp_parameter_parse(buf, (((headerLength) + (1))) - (curPos), (void*) &_value);
      if(_res != OK) {
        return _res;
      }
      plc4c_utils_list_insert_head_value(parameters, _value);
      curPos = plc4c_spi_read_get_pos(buf) - startPos;
    }
  }
  (*_message)->parameters = parameters;

  // Optional Field (payload) (Can be skipped, if a given expression evaluates to false)
  curPos = plc4c_spi_read_get_pos(buf) - startPos;
  plc4c_s7_read_write_s7_message* payload = NULL;
  if((curPos) < (cotpLen)) {
    payload = malloc(sizeof(plc4c_s7_read_write_s7_message));
    if(payload == NULL) {
      return NO_MEMORY;
    }
    plc4c_return_code _res = plc4c_s7_read_write_s7_message_parse(buf, &payload);
    if(_res != OK) {
      return _res;
    }
    (*_message)->payload = payload;
  }

  return OK;
}

plc4c_return_code plc4c_s7_read_write_cotp_packet_serialize(plc4c_spi_write_buffer* buf, plc4c_s7_read_write_cotp_packet* _message) {

  // Discriminator Field (tpduCode)
  plc4c_spi_write_unsigned_short(buf, 8, plc4c_s7_read_write_cotp_packet_get_discriminator(_message->_type).tpduCode);

  // Array field (parameters)
  {
    uint8_t itemCount = plc4c_utils_list_size(_message->parameters);
    for(int curItem = 0; curItem < itemCount; curItem++) {
      plc4c_s7_read_write_cotp_parameter* _value = (plc4c_s7_read_write_cotp_parameter*) plc4c_utils_list_get_value(_message->parameters, curItem);
      plc4c_return_code _res = plc4c_s7_read_write_cotp_parameter_serialize(buf, (void*) &_value);
      if(_res != OK) {
        return _res;
      }
    }
  }

  // Optional Field (payload)
  if(_message->payload != NULL) {
    plc4c_s7_read_write_s7_message* _value = (plc4c_s7_read_write_s7_message*) _message->payload;
    plc4c_return_code _res = plc4c_s7_read_write_s7_message_serialize(buf, (void*) &_value);
    if(_res != OK) {
      return _res;
    }
  }

  return OK;
}
