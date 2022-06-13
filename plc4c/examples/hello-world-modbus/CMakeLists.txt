#[[
  Licensed to the Apache Software Foundation (ASF) under one
  or more contributor license agreements.  See the NOTICE file
  distributed with this work for additional information
  regarding copyright ownership.  The ASF licenses this file
  to you under the Apache License, Version 2.0 (the
  "License"); you may not use this file except in compliance
  with the License.  You may obtain a copy of the License at

      https://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing,
  software distributed under the License is distributed on an
  "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
  KIND, either express or implied.  See the License for the
  specific language governing permissions and limitations
  under the License.
]]

include_directories("../../api/include"
        "../../spi/include"
        "../../drivers/modbus/include"
        "../../transports/tcp/include")

add_executable(plc4c-examples-hello-world-modbus
        src/hello_world_modbus.c
        )

IF (NOT WIN32)
    target_link_libraries(plc4c-examples-hello-world-modbus
            plc4c-spi
            plc4c-driver-modbus
            plc4c-transport-tcp
            m
    )
ELSE()
    target_link_libraries(plc4c-examples-hello-world-modbus
            plc4c-spi
            plc4c-driver-modbus
            plc4c-transport-tcp
            unity
    )
ENDIF()
