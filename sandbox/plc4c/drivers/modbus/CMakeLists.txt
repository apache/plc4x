#[[
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
]]

include_directories("include" "../../api/include" "../../spi/include"
    "../../generated-sources/modbus/include")

# Add the generated sources
file(GLOB generatedSources "../../generated-sources/modbus/src/*.c")

add_library(plc4c-driver-modbus
        src/driver_modbus.c
        src/driver_modbus_encode_decode.c
        src/driver_modbus_packets.c
        src/driver_modbus_sm_connect.c
        src/driver_modbus_sm_disconnect.c
        src/driver_modbus_sm_read.c
        src/driver_modbus_sm_write.c
        ${generatedSources})

target_link_libraries(plc4c-driver-modbus plc4c-spi ${CMAKE_DL_LIBS})