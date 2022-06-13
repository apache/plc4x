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

cmake_minimum_required(VERSION 3.7)
set(CMAKE_CXX_STANDARD 11)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_CXX_EXTENSIONS OFF)
set(CMAKE_POSITION_INDEPENDENT_CODE ON)

# Save the root directory
set(PLC4CPP_ROOT_DIR ${CMAKE_SOURCE_DIR})
# Define some user-defined profiles
set(WITH_PROXIES OFF CACHE BOOL "Enable building Thrift based proxy driver")

# Tell the tooling where to find our Boost installation
set(BOOST_INCLUDEDIR "${PLC4CPP_ROOT_DIR}/target/lib/boost/include")
set(BOOST_LIBRARYDIR "${PLC4CPP_ROOT_DIR}/target/lib/boost/lib")

# Disable usage of Boost-CMake
set(Boost_NO_BOOST_CMAKE ON)
# Force CMake to use the boost version provided in BOOST_ROOT
set(Boost_NO_SYSTEM_PATHS ON)
set(Boost_USE_STATIC_LIBS ON)
set(Boost_USE_MULTITHREADED ON)
set(Boost_ARCHITECTURE "-x64")

# Ensure the listed modules are found and environment variables for boost are set.
find_package(Boost COMPONENTS thread date_time system regex log)

# Trace all CMAKE Variables
get_cmake_property(_variableNames VARIABLES)
list (SORT _variableNames)
foreach (_variableName ${_variableNames})
    message(STATUS "${_variableName}=${${_variableName}}")
endforeach()

#[[
    Build all the modules of PLC4Cpp
]]
add_subdirectory(api)
add_subdirectory(protocols)
add_subdirectory(drivers)
add_subdirectory(utils)
add_subdirectory(examples)