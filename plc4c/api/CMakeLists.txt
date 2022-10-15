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

#[[
file(GLOB sources "src/*.c")

add_library(plc4c-api SHARED ${sources})

target_include_directories(plc4c-api PUBLIC "include")

if(BUILD_PHASE STREQUAL test-compile)
    file(GLOB testSources "test/*.c")
    add_executable(plc4c-api-test ${testSources})
    target_link_libraries (plc4c-api-test
                           plc4c-api
                           unity
    )
    add_test(NAME plc4c-api-test COMMAND plc4c-api-test)
endif()
]]