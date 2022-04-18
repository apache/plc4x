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

#include <plc4c/spi/types_private.h>
#include <plc4c/transport_serial.h>
#include <stdlib.h>

plc4c_return_code plc4c_transport_serial_open_function(void* config) {
  return OK;
}

plc4c_return_code plc4c_transport_serial_close_function(void* config) {
  return OK;
}

plc4c_return_code plc4c_transport_serial_send_message_function(
    void* transport_configuration, plc4c_spi_write_buffer* message) {
  return OK;
}

plc4c_return_code plc4c_transport_serial_select_message_function(
    void* transport_configuration, uint8_t min_size,
    accept_message_function accept_message, plc4c_spi_read_buffer** message) {
  return OK;
}

plc4c_transport *plc4c_transport_serial_create() {
  plc4c_transport *transport =
      (plc4c_transport *)malloc(sizeof(plc4c_transport));
  transport->transport_code = "serial";
  transport->transport_name = "Serial port transport";
  transport->open = &plc4c_transport_serial_open_function;
  transport->close = &plc4c_transport_serial_close_function;
  transport->send_message = &plc4c_transport_serial_send_message_function;
  transport->select_message = &plc4c_transport_serial_select_message_function;
  return transport;
}
