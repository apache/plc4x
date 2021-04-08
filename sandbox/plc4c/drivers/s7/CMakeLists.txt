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
    "../../generated-sources/s7/include")

# Add the generated sources
file(GLOB generatedSources "../../generated-sources/s7/src/*.c")

add_library(plc4c-driver-s7
        src/driver_s7.c
        src/driver_s7_encode_decode.c
        src/driver_s7_packets.c
        src/driver_s7_sm_connect.c
        src/driver_s7_sm_disconnect.c
        src/driver_s7_sm_read.c
        src/driver_s7_sm_write.c
        ${generatedSources})

target_link_libraries(plc4c-driver-s7
        plc4c-spi)

if (BUILD_PHASE STREQUAL test-compile)
    add_executable(plc4c-driver-s7-test
            test/s7_address_parser_test.c
            test/s7_test.c
            )

    IF (NOT WIN32)
        target_link_libraries(plc4c-driver-s7-test
                plc4c-driver-s7
                unity
                m
                )
    ELSE()
        target_link_libraries(plc4c-driver-s7-test
                plc4c-driver-s7
                unity
                )
    ENDIF()
    add_test(NAME plc4c-driver-s7-test COMMAND plc4c-driver-s7-test)
endif ()