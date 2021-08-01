/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

#include "plc4c/driver_s7.h"

#include <string.h>

#include "plc4c/driver_s7_encode_decode.h"
#include "plc4c/driver_s7_sm.h"

plc4c_return_code plc4c_driver_s7_configure_function(plc4c_list* parameters,
                                                     void** configuration) {
  plc4c_driver_s7_config* s7_config = malloc(sizeof(plc4c_driver_s7_config));
  if (s7_config == NULL) {
    return NO_MEMORY;
  }

  // Initialize the parts that the user can influence.
  s7_config->local_rack = 1;
  s7_config->local_slot = 1;
  s7_config->remote_rack = 0;
  s7_config->remote_slot = 0;
  s7_config->pdu_size = 1024;
  s7_config->max_amq_caller = 8;
  s7_config->max_amq_callee = 8;
  s7_config->controller_type = PLC4C_DRIVER_S7_CONTROLLER_TYPE_ANY;

  // TODO: Apply the values from the parameters.

  *configuration = s7_config;
  return OK;
}

plc4c_driver* plc4c_driver_s7_create() {
  plc4c_driver* driver = (plc4c_driver*)malloc(sizeof(plc4c_driver));
  driver->protocol_code = "s7";
  driver->protocol_name = "Siemens S7 (Basic)";
  driver->default_transport_code = "tcp";
  driver->parse_address_function = &plc4c_driver_s7_encode_address;
  driver->configure_function = &plc4c_driver_s7_configure_function;
  driver->connect_function = &plc4c_driver_s7_connect_function;
  driver->disconnect_function = &plc4c_driver_s7_disconnect_function;
  driver->read_function = &plc4c_driver_s7_read_function;
  driver->write_function = &plc4c_driver_s7_write_function;
  driver->subscribe_function = NULL;
  driver->unsubscribe_function = NULL;
  driver->free_read_request_function = &plc4c_driver_s7_free_read_request;
  driver->free_write_request_function = &plc4c_driver_s7_free_write_request;
  driver->free_read_response_function = &plc4c_driver_s7_free_read_response;
  driver->free_write_response_function = &plc4c_driver_s7_free_write_response;
  driver->free_subscription_response_function = NULL;
  driver->free_unsubscription_response_function = NULL;
  return driver;
}

/*
 *
 *   Static functions
 *
 */

char* plc4c_s7_read_write_parse_s7_string(plc4c_spi_read_buffer* io,
                                          int32_t stringLength,
                                          char* encoding) {
  if (strcmp(encoding, "UTF-8") == 0) {
    // Read the max length (which is not interesting for us.
    uint8_t maxLen;
    plc4c_return_code res = plc4c_spi_read_unsigned_byte(io, 8, &maxLen);
    if (res != OK) {
      return NULL;
    }
    // Read the effective length of the string.
    uint8_t effectiveStringLength;
    res = plc4c_spi_read_unsigned_byte(io, 8, &effectiveStringLength);
    if (res != OK) {
      return NULL;
    }
    char* result = malloc(sizeof(char) * (effectiveStringLength + 1));
    if (result == NULL) {
      return NULL;
    }
    char* curPos = result;
    for(int i = 0; i < effectiveStringLength; i++) {
      uint8_t val;
      plc4c_return_code res = plc4c_spi_read_unsigned_byte(io, 8, &val);
      if (res != OK) {
        return NULL;
      }
      *curPos = (char) val;
      curPos++;
    }
    *curPos = '\0';
    return result;
  } else if (strcmp(encoding, "UTF-16") == 0) {
  }
  return "";
}

char* plc4c_s7_read_write_parse_s7_char(plc4c_spi_read_buffer* io,
                                        char* encoding) {
  if (strcmp(encoding, "UTF-8") == 0) {
    char* result = malloc(sizeof(char) * 2);
    if (result == NULL) {
      return NULL;
    }
    uint8_t val;
    plc4c_return_code res = plc4c_spi_read_unsigned_byte(io, 8, &val);
    if (res != OK) {
      return NULL;
    }
    *result = (char) val;
    *(result+1) = '\0';
    return result;
  } else if (strcmp(encoding, "UTF-16") == 0) {
  }
  return "";
}

time_t plc4c_s7_read_write_parse_tia_time(plc4c_spi_read_buffer* io) {
  // TODO: Implement ...
  return 0;
}

time_t plc4c_s7_read_write_parse_s5_time(plc4c_spi_read_buffer* io) {
  // TODO: Implement ...
  return 0;
}

time_t plc4c_s7_read_write_parse_tia_l_time(plc4c_spi_read_buffer* io) {
  // TODO: Implement ...
  return 0;
}

time_t plc4c_s7_read_write_parse_tia_date(plc4c_spi_read_buffer* io) {
  // TODO: Implement ...
  return 0;
}

time_t plc4c_s7_read_write_parse_tia_time_of_day(plc4c_spi_read_buffer* io) {
  // TODO: Implement ...
  return 0;
}

time_t plc4c_s7_read_write_parse_tia_date_time(plc4c_spi_read_buffer* io) {
  // TODO: Implement ...
  return 0;
}
