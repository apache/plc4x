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
