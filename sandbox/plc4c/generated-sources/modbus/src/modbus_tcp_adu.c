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
#include "modbus_tcp_adu.h"

// Parse function.
plc4c_return_code plc4c_modbus_read_write_modbus_tcp_adu_parse(plc4c_spi_read_buffer* buf, bool response, plc4c_modbus_read_write_modbus_tcp_adu** _message) {
  uint16_t startPos = plc4c_spi_read_get_pos(buf);
  uint16_t curPos;
  plc4c_return_code _res = OK;

  // Allocate enough memory to contain this data structure.
  (*_message) = malloc(sizeof(plc4c_modbus_read_write_modbus_tcp_adu));
  if(*_message == NULL) {
    return NO_MEMORY;
  }

  // Simple Field (transactionIdentifier)
  uint16_t transactionIdentifier = 0;
  _res = plc4c_spi_read_unsigned_int(buf, 16, (uint32_t*) &transactionIdentifier);
  if(_res != OK) {
    return _res;
  }
  (*_message)->transaction_identifier = transactionIdentifier;

  // Const Field (protocolIdentifier)
  uint16_t protocolIdentifier = 0;
  _res = plc4c_spi_read_unsigned_int(buf, 16, (uint32_t*) &protocolIdentifier);
  if(_res != OK) {
    return _res;
  }
  if(protocolIdentifier != MODBUS_READ_WRITE_MODBUS_TCP_ADU_PROTOCOL_IDENTIFIER) {
    return PARSE_ERROR;
    // throw new ParseException("Expected constant value " + MODBUS_READ_WRITE_MODBUS_TCP_ADU_PROTOCOL_IDENTIFIER + " but got " + protocolIdentifier);
  }

  // Implicit Field (length) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
  uint16_t length = 0;
  _res = plc4c_spi_read_unsigned_int(buf, 16, (uint32_t*) &length);
  if(_res != OK) {
    return _res;
  }

  // Simple Field (unitIdentifier)
  uint8_t unitIdentifier = 0;
  _res = plc4c_spi_read_unsigned_short(buf, 8, (uint16_t*) &unitIdentifier);
  if(_res != OK) {
    return _res;
  }
  (*_message)->unit_identifier = unitIdentifier;

  // Simple Field (pdu)
  plc4c_modbus_read_write_modbus_pdu* pdu;
  _res = plc4c_modbus_read_write_modbus_pdu_parse(buf, response, (void*) &pdu);
  if(_res != OK) {
    return _res;
  }
  (*_message)->pdu = pdu;

  return OK;
}

plc4c_return_code plc4c_modbus_read_write_modbus_tcp_adu_serialize(plc4c_spi_write_buffer* buf, plc4c_modbus_read_write_modbus_tcp_adu* _message) {
  plc4c_return_code _res = OK;

  // Simple Field (transactionIdentifier)
  _res = plc4c_spi_write_unsigned_int(buf, 16, _message->transaction_identifier);
  if(_res != OK) {
    return _res;
  }

  // Const Field (protocolIdentifier)
  plc4c_spi_write_unsigned_int(buf, 16, MODBUS_READ_WRITE_MODBUS_TCP_ADU_PROTOCOL_IDENTIFIER);

  // Implicit Field (length) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
  _res = plc4c_spi_write_unsigned_int(buf, 16, (plc4c_modbus_read_write_modbus_pdu_length_in_bytes(_message->pdu)) + (1));
  if(_res != OK) {
    return _res;
  }

  // Simple Field (unitIdentifier)
  _res = plc4c_spi_write_unsigned_short(buf, 8, _message->unit_identifier);
  if(_res != OK) {
    return _res;
  }

  // Simple Field (pdu)
  _res = plc4c_modbus_read_write_modbus_pdu_serialize(buf, _message->pdu);
  if(_res != OK) {
    return _res;
  }

  return OK;
}

uint8_t plc4c_modbus_read_write_modbus_tcp_adu_length_in_bytes(plc4c_modbus_read_write_modbus_tcp_adu* _message) {
  return plc4c_modbus_read_write_modbus_tcp_adu_length_in_bits(_message) / 8;
}

uint8_t plc4c_modbus_read_write_modbus_tcp_adu_length_in_bits(plc4c_modbus_read_write_modbus_tcp_adu* _message) {
  uint8_t lengthInBits = 0;

  // Simple field (transactionIdentifier)
  lengthInBits += 16;

  // Const Field (protocolIdentifier)
  lengthInBits += 16;

  // Implicit Field (length)
  lengthInBits += 16;

  // Simple field (unitIdentifier)
  lengthInBits += 8;

  // Simple field (pdu)
  lengthInBits += plc4c_modbus_read_write_modbus_pdu_length_in_bits(_message->pdu);

  return lengthInBits;
}

