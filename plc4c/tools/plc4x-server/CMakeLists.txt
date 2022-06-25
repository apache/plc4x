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

SET(generatedSources ${PLC4C_ROOT_DIR}/generated-sources/plc4x/src/plc4x_constants.c)
SET(generatedHeaders ${PLC4C_ROOT_DIR}/generated-sources/plc4x/include/plc4x_constants.h)

# Add the sources
add_executable(plc4c-tools-plc4x-server
        src/plc4x_server.c
        ${generatedHeaders}
        ${generatedSources}
        )

target_include_directories(plc4c-tools-plc4x-server
        PRIVATE
        ${PLC4C_ROOT_DIR}/api/include
        ${PLC4C_ROOT_DIR}/spi/include
        ${PLC4C_ROOT_DIR}/drivers/plc4x/include
        ${PLC4C_ROOT_DIR}/generated-sources/plc4x/include)

if(WIN32)
    target_link_libraries(plc4c-tools-plc4x-server
            plc4c-driver-plc4x
            plc4c-transport-tcp
            plc4c-spi
            wsock32
            ws2_32
            ${CMAKE_DL_LIBS}
    )
else()
    target_link_libraries(plc4c-tools-plc4x-server
            plc4c-driver-plc4x
            plc4c-transport-tcp
            plc4c-spi
            ${CMAKE_DL_LIBS}
    )
endif()
