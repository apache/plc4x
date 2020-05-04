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
#include <stdio.h>
#include "../../../../../api/src/main/include/plc4c.h"

void onConnection2Success(plc4c_connection *connection) {

}

void onConnection2Error(const char *connection_string, error_code error) {

}

int main() {
  bool loop = true;
  plc4c_system *system = NULL;
  plc4c_connection *connection1 = NULL;
  plc4c_connection *connection2 = NULL;

  error_code error = plc4c_system_create(&system);
  if (error != OK) {
    return -1;
  }

  error = plc4c_init(system);
  if (error != OK) {
    return -1;
  }

  // Establish connections to remote devices
  // you may or may not care about the connection handle
  error = plc4c_system_connect(system, "s7://192.168.42.20", &connection1);
  if (error != OK) {
    return -1;
  }

  error = plc4c_system_connect_callback(system, "s7://192.168.42.22", &connection2,
          &onConnection2Success, &onConnection2Error);
  if (error != OK) {
    return -1;
  }

  // Central program loop ...
  while (loop) {
    if (plc4c_system_loop(system) != OK) {
      break;
    }
  }

  // Make sure everything is cleaned up correctly.
  plc4c_system_shutdown(system);
  plc4c_system_destroy(system);

  return 0;
}
