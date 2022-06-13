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

add_library(plc4cpp-driver-s7 SHARED
        src/main/cpp/org/apache/plc4x/cpp/s7/dllexports.cpp
        src/main/cpp/org/apache/plc4x/cpp/s7/S7PlcDriver.cpp
		src/main/cpp/org/apache/plc4x/cpp/s7/connection/S7PlcConnection.cpp
            )

target_include_directories (plc4cpp-driver-s7 PUBLIC ${Boost_INCLUDE_DIRS})
target_include_directories (plc4cpp-driver-s7 PUBLIC ../../api/src/main/cpp)
target_include_directories (plc4cpp-driver-s7 PUBLIC ../../protocols/driver-bases/base/src/main/cpp)
target_include_directories (plc4cpp-driver-s7 PUBLIC ../../protocols/s7/src/main/cpp)
target_include_directories (plc4cpp-driver-s7 PUBLIC ../../utils/logger/src/main/cpp)
target_include_directories (plc4cpp-driver-s7 PUBLIC ../../utils/systemconfig/src/main/cpp)

target_link_libraries (plc4cpp-driver-s7
                        ${Boost_LIBRARIES}
                        plc4cpp-utils-logger 
                        plc4cpp-api 
                        plc4cpp-protocols-driver-base-base
                        plc4cpp-utils-systemconfig
                      )

if(MINGW)
    target_link_libraries (plc4cpp-driver-s7 ws2_32)
endif()
