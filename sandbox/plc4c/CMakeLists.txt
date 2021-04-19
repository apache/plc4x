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

cmake_minimum_required(VERSION 3.16)

# Set the name of this project
project(PLC4C)

# Set the C language level to C11
set(CMAKE_C_STANDARD 11)

# Save the root directory
set(PLC4C_ROOT_DIR ${CMAKE_SOURCE_DIR})

# Our link to the maven lifecycle
set(BUILD_PHASE test-compile CACHE STRING "Phase of the Maven build we are executing cmake")

# Access the Unity version the maven build is providing us with.
set(UNITY_VERSION 2.5.0 CACHE STRING "Version of the used Unity test framework")

# Depending on the phase of the build we are currently running, initialize
# The test system.
if (BUILD_PHASE STREQUAL compile)
    # Nothing really to do here ... just need it to check the known values.
elseif (BUILD_PHASE STREQUAL test-compile)
    # Initialize the test subsystem as well as the Unity library sources
    include(CTest)

    # Make the Unity sources available to the build as "Unity" library
    #add_subdirectory(target/dependency/Unity-${UNITY_VERSION})

    #add_library(Unity STATIC
    #  target/dependency/Unity-${UNITY_VERSION}/src/unity.c
    #)
    #target_include_directories(Unity PUBLIC
    #  target/dependency/Unity-${UNITY_VERSION}/src
    #)
else ()
    # Output an error
    message(FATAL_ERROR "Given BUILD_PHASE unknown. Known values are 'compile' and 'test-compile'")
endif ()

set(DEBUG_OUTPUT ON CACHE BOOL "Enable outputting of debug information")
if (DEBUG_OUTPUT)
    # Trace all CMAKE Variables
    get_cmake_property(_variableNames VARIABLES)
    list(SORT _variableNames)
    foreach (_variableName ${_variableNames})
        message(STATUS "${_variableName}=${${_variableName}}")
    endforeach ()
endif ()

#[[
    Build all the modules of PLC4C
]]
# Core stuff
add_subdirectory(api)
add_subdirectory(spi)
add_subdirectory(transports)
add_subdirectory(drivers)
#add_subdirectory(integrations)
add_subdirectory(examples)
