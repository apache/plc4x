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

#include "plc4c/driver_plc4x.h"

#include "plc4c/driver_plc4x_encode_decode.h"
#include "plc4c/driver_plc4x_sm.h"

plc4c_return_code plc4c_driver_plc4x_configure_function(plc4c_list* parameters,
                                                     void** configuration) {
  plc4c_driver_plc4x_config* plc4x_config = malloc(sizeof(plc4c_driver_plc4x_config));
  if (plc4x_config == NULL) {
    return NO_MEMORY;
  }

  // Initialize the parts that the user can influence.
  plc4x_config->remote_connection_string = NULL;
  plc4x_config->request_timeout = 5000;

  // TODO: Apply the values from the parameters.

  *configuration = plc4x_config;
  return OK;
}

plc4c_driver* plc4c_driver_plc4x_create() {
  plc4c_driver* driver = (plc4c_driver*)malloc(sizeof(plc4c_driver));
  driver->protocol_code = "plc4c";
  driver->protocol_name = "PLC4X (Proxy-Protocol)";
  driver->default_transport_code = "tcp";
  driver->parse_address_function = &plc4c_driver_plc4x_encode_address;
  driver->configure_function = &plc4c_driver_plc4x_configure_function;
  driver->connect_function = &plc4c_driver_plc4x_connect_function;
  driver->disconnect_function = &plc4c_driver_plc4x_disconnect_function;
  driver->read_function = &plc4c_driver_plc4x_read_function;
  driver->write_function = &plc4c_driver_plc4x_write_function;
  driver->subscribe_function = NULL;
  driver->unsubscribe_function = NULL;
  driver->free_read_request_function = &plc4c_driver_plc4x_free_read_request;
  driver->free_write_request_function = &plc4c_driver_plc4x_free_write_request;
  driver->free_read_response_function = &plc4c_driver_plc4x_free_read_response;
  driver->free_write_response_function = &plc4c_driver_plc4x_free_write_response;
  driver->free_subscription_response_function = NULL;
  driver->free_unsubscription_response_function = NULL;
  return driver;
}
