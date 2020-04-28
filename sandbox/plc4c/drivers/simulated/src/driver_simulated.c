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

#include <stdlib.h>
#include <plc4c/spi/types_private.h>
#include <plc4c/driver_simulated.h>

struct plc4c_driver_simulated_item_t {
    // Sort of senseless, but for keeping track of the fact that it's an plc4c_item.
    plc4c_item super;

    // Actual properties goes here ...
};
typedef struct plc4c_driver_simulated_item_t plc4c_driver_simulated_item ;

plc4c_item *plc4c_driver_simulated_parse_address(const char* address_string) {
    plc4c_driver_simulated_item* item = (plc4c_driver_simulated_item*) malloc(sizeof(plc4c_driver_simulated_item));
    return (plc4c_item*) item;
}

plc4c_driver *plc4c_driver_simulated_create() {
    plc4c_driver* driver = (plc4c_driver*) malloc(sizeof(plc4c_driver));
    driver->protocol_code = "simulated";
    driver->protocol_name = "Simulated PLC4X Datasource";
    driver->default_transport_code = "dummy";
    driver->parse_address_function = &plc4c_driver_simulated_parse_address;
    return driver;
}

