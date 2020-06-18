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
#include "s7_address.h"

// Array of discriminator values that match the enum type constants.
// (The order is identical to the enum constants so we can use the
// enum constant to directly access a given types discriminator values)
const plc4c_s7_read_write_s7_address_discriminator plc4c_s7_read_write_s7_address_discriminators[] = {
  {/* s7_read_write_s7_address_any */
   .addressType = 0x10}
};

// Function returning the discriminator values for a given type constant.
plc4c_s7_read_write_s7_address_discriminator plc4c_s7_read_write_s7_address_get_discriminator(plc4c_s7_read_write_s7_address_type type) {
  return plc4c_s7_read_write_s7_address_discriminators[type];
}

// Parse function.
plc4c_return_code plc4c_s7_read_write_s7_address_parse(plc4c_spi_read_buffer* buf, plc4c_s7_read_write_s7_address** _message) {
  uint16_t startPos = plc4c_spi_read_get_pos(buf);
  uint16_t curPos;

  // Pointer to the parsed data structure.
  (*_message) = malloc(sizeof(plc4c_s7_read_write_s7_address));
  if(*_message == NULL) {
    return NO_MEMORY;
  }

  // Discriminator Field (addressType) (Used as input to a switch field)
  uint8_t addressType = plc4c_spi_read_unsigned_short(buf, 8);

  // Switch Field (Depending on the discriminator values, passes the instantiation to a sub-type)
  if(addressType == 0x10) { /* S7AddressAny */
                    
    // Enum field (transportSize)
    plc4c_s7_read_write_transport_size transportSize = plc4c_spi_read_byte(buf, 8);
    (*_message)->s7_address_any_transport_size = transportSize;


                    
    // Simple Field (numberOfElements)
    uint16_t numberOfElements = plc4c_spi_read_unsigned_int(buf, 16);
    (*_message)->s7_address_any_number_of_elements = numberOfElements;


                    
    // Simple Field (dbNumber)
    uint16_t dbNumber = plc4c_spi_read_unsigned_int(buf, 16);
    (*_message)->s7_address_any_db_number = dbNumber;


                    
    // Enum field (area)
    plc4c_s7_read_write_memory_area area = plc4c_spi_read_byte(buf, 8);
    (*_message)->s7_address_any_area = area;


                    
    // Reserved Field (Compartmentalized so the "reserved" variable can't leak)
    {
      unsigned int _reserved = plc4c_spi_read_unsigned_short(buf, 5);
      if(_reserved != 0x00) {
        printf("Expected constant value '%d' but got '%d' for reserved field.", 0x00, _reserved);
      }
    }


                    
    // Simple Field (byteAddress)
    uint16_t byteAddress = plc4c_spi_read_unsigned_int(buf, 16);
    (*_message)->s7_address_any_byte_address = byteAddress;


                    
    // Simple Field (bitAddress)
    unsigned int bitAddress = plc4c_spi_read_unsigned_byte(buf, 3);
    (*_message)->s7_address_any_bit_address = bitAddress;

  }

  return OK;
}

plc4c_return_code plc4c_s7_read_write_s7_address_serialize(plc4c_spi_write_buffer* buf, plc4c_s7_read_write_s7_address* message) {
  return OK;
}
