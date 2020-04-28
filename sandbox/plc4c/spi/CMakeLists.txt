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

include_directories("include" "../api/include")

file(GLOB sources "src/*.c")

add_library(plc4c-spi SHARED ${sources})

if(BUILD_PHASE STREQUAL test-compile)
    file(GLOB testSources "test/*.c")
    add_executable(plc4c-spi-test ${testSources} test/system_test.c)
    target_link_libraries (plc4c-spi-test
            plc4c-spi
            unity
            )
    add_test(NAME plc4c-spi-test COMMAND plc4c-spi-test)
endif()