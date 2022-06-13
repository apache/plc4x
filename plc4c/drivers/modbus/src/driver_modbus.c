/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

#include "plc4c/driver_modbus.h"
#include "plc4c/driver_modbus_sm.h"
#include "plc4c/driver_modbus_encode_decode.h"

plc4c_return_code plc4c_driver_modbus_configure_function(plc4c_list* parameters,
                                                     void** configuration) {
  plc4c_driver_modbus_config* modbus_config = malloc(sizeof(plc4c_driver_modbus_config));
  if (modbus_config == NULL) {
    return NO_MEMORY;
  }

  // Initialize the parts that the user can influence.
  modbus_config->request_timeout = 5000;
  modbus_config->unit_identifier = 1;
  modbus_config->communication_id_counter = 1;

  // TODO: Apply the values from the parameters.

  *configuration = modbus_config;
  return OK;
}

plc4c_driver *plc4c_driver_modbus_create() {
  plc4c_driver *driver = (plc4c_driver *)malloc(sizeof(plc4c_driver));
  driver->protocol_code = "modbus-tcp";
  driver->protocol_name = "Modbus TCP";
  driver->default_transport_code = "tcp";
  driver->parse_address_function = &plc4c_driver_modbus_encode_address;
  driver->configure_function = &plc4c_driver_modbus_configure_function;
  driver->connect_function = &plc4c_driver_modbus_connect_function;
  driver->disconnect_function = &plc4c_driver_modbus_disconnect_function;
  driver->read_function = &plc4c_driver_modbus_read_function;
  driver->write_function = &plc4c_driver_modbus_write_function;
  driver->subscribe_function = NULL;
  driver->unsubscribe_function = NULL;
  driver->free_read_request_function = &plc4c_driver_modbus_free_read_request;
  driver->free_read_response_function = &plc4c_driver_modbus_free_read_response;
  driver->free_write_request_function = &plc4c_driver_modbus_free_write_request;
  driver->free_write_response_function = &plc4c_driver_modbus_free_write_response;
  driver->free_subscription_response_function = NULL;
  driver->free_unsubscription_response_function = NULL;
  return driver;
}
