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

cmake_minimum_required(VERSION 3.7)
set(CMAKE_CXX_STANDARD 11)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_CXX_EXTENSIONS OFF)

add_library(plc4cpp-utils-logger src/main/cpp/org/apache/plc4x/cpp/utils/logger/BLogger.cpp
        src/main/cpp/org/apache/plc4x/cpp/utils/logger/DatDmp.cpp
		src/main/cpp/org/apache/plc4x/cpp/utils/logger/DbgTrace.cpp
		src/main/cpp/org/apache/plc4x/cpp/utils/logger/ErrorCategory.cpp
		src/main/cpp/org/apache/plc4x/cpp/utils/logger/ErrorInfoException.cpp
		src/main/cpp/org/apache/plc4x/cpp/utils/logger/ExLog.cpp
        )

#[[
    Import the headers for boost and the plc4cpp-utils module
]]
target_include_directories (plc4cpp-utils-logger PUBLIC target/dependencies/include)
