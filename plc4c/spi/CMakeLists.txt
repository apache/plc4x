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

add_library(plc4c-spi
        src/utils/list.c
        src/utils/queue.c
        src/connection.c
        src/data.c
        src/evaluation_helper.c
        src/read.c
        src/read_buffer.c
        src/subscribe.c
        src/system.c
        src/types.c
        src/write.c
        src/write_buffer.c
        )

target_include_directories(
        plc4c-spi
        PUBLIC
        "include"
        "../api/include")

IF (NOT WIN32)
    target_link_libraries(plc4c-spi
            m
            )
ELSE()
    target_link_libraries(plc4c-spi
            )
ENDIF()

if (BUILD_PHASE STREQUAL test-compile)
    add_executable(plc4c-spi-test
            test/system_test.c
            test/spi_test.c
            test/read_buffer_test.c
            test/write_buffer_test.c
            )
    target_link_libraries(plc4c-spi-test
            plc4c-spi
            unity

            )
    add_test(NAME plc4c-spi-test COMMAND plc4c-spi-test)
endif ()