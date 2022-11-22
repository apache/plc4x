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

add_executable(plc4c-examples-hello-world-s7
        src/hello_world_s7.c
        )

target_include_directories(
        plc4c-examples-hello-world-s7
        PUBLIC
        "../../api/include"
        "../../spi/include"
        "../../drivers/s7/include"
        "../../transports/tcp/include" )

add_executable(plc4c-examples-loopback-s7
		src/hello_world_s7_loopback.c
		)

target_include_directories(
        plc4c-examples-loopback-s7
        PUBLIC
        "../../api/include"
        "../../spi/include"
        "../../drivers/s7/include"
        "../../transports/tcp/include" )

IF (NOT WIN32)
    target_link_libraries(plc4c-examples-hello-world-s7
            plc4c-spi
            plc4c-driver-s7
            plc4c-transport-tcp
            m
    )
    target_link_libraries(plc4c-examples-loopback-s7
            plc4c-spi
            plc4c-driver-s7
            plc4c-transport-tcp
            m
    )
ELSE()
    target_link_libraries(plc4c-examples-hello-world-s7
            plc4c-spi
            plc4c-driver-s7
            plc4c-transport-tcp
    )
    target_link_libraries(plc4c-examples-loopback-s7
            plc4c-spi
            plc4c-driver-s7
            plc4c-transport-tcp
    )
ENDIF()

