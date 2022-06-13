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

# Collect all the generated files
aux_source_directory(../../target/generated-sources/drivers/proxy PROXY_GENERATED_SOURCES)

message(STATUS "Found generated sources: ${PROXY_GENERATED_SOURCES}")

add_library(plc4cpp-driver-proxy SHARED
        src/main/cpp/org/apache/plc4x/cpp/proxy/dllexports.cpp
        src/main/cpp/org/apache/plc4x/cpp/proxy/ProxyDriver.cpp
		src/main/cpp/org/apache/plc4x/cpp/proxy/connection/ProxyConnection.cpp
        ${PROXY_GENERATED_SOURCES}
            )

target_include_directories (plc4cpp-driver-proxy PUBLIC ${Boost_INCLUDE_DIRS})
target_include_directories (plc4cpp-driver-proxy PUBLIC ${PLC4CPP_ROOT_DIR}/target/lib/thrift/include/)
target_include_directories (plc4cpp-driver-proxy PUBLIC ../../api/src/main/cpp)
target_include_directories (plc4cpp-driver-proxy PUBLIC ../../protocols/driver-bases/base/src/main/cpp)
target_include_directories (plc4cpp-driver-proxy PUBLIC ../../protocols/s7/src/main/cpp)
target_include_directories (plc4cpp-driver-proxy PUBLIC ../../utils/logger/src/main/cpp)
target_include_directories (plc4cpp-driver-proxy PUBLIC ../../utils/systemconfig/src/main/cpp)

target_link_libraries (plc4cpp-driver-proxy
                        ${Boost_LIBRARIES}
                        ${PLC4CPP_ROOT_DIR}/target/lib/thrift/lib/libparse.a ${PLC4CPP_ROOT_DIR}/target/lib/thrift/lib/libthrift.a ${PLC4CPP_ROOT_DIR}/target/lib/thrift/lib/libthriftz.a ${PLC4CPP_ROOT_DIR}/target/lib/thrift/lib/libtutorialgencpp.a
                        plc4cpp-utils-logger 
                        plc4cpp-api 
                        plc4cpp-protocols-driver-base-base
                        plc4cpp-utils-systemconfig
                      )
