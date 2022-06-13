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

add_library(plc4cpp-protocols-driver-base-base
        src/main/cpp/org/apache/plc4x/cpp/base/connection/AbstractPlcConnection.cpp
        src/main/cpp/org/apache/plc4x/cpp/base/connection/BoostConnection.cpp
        src/main/cpp/org/apache/plc4x/cpp/base/connection/ChannelFactory.cpp
        src/main/cpp/org/apache/plc4x/cpp/base/connection/DefaultPlcFieldHandler.cpp
        src/main/cpp/org/apache/plc4x/cpp/base/connection/PlcFieldHandler.cpp
        src/main/cpp/org/apache/plc4x/cpp/base/events/ConnectedEvent.cpp
        src/main/cpp/org/apache/plc4x/cpp/base/events/ConnectEvent.cpp
        src/main/cpp/org/apache/plc4x/cpp/base/messages/items/BaseDefaultFieldItem.cpp
        src/main/cpp/org/apache/plc4x/cpp/base/messages/items/DefaultBigDecimalFieldItem.cpp
        src/main/cpp/org/apache/plc4x/cpp/base/messages/items/DefaultBigIntegerFieldItem.cpp
        src/main/cpp/org/apache/plc4x/cpp/base/messages/items/DefaultBooleanFieldItem.cpp
        src/main/cpp/org/apache/plc4x/cpp/base/messages/items/DefaultByteArrayFieldItem.cpp
        src/main/cpp/org/apache/plc4x/cpp/base/messages/items/DefaultByteFieldItem.cpp
        src/main/cpp/org/apache/plc4x/cpp/base/messages/items/DefaultDoubleFieldItem.cpp
        src/main/cpp/org/apache/plc4x/cpp/base/messages/items/DefaultFloatFieldItem.cpp
        src/main/cpp/org/apache/plc4x/cpp/base/messages/items/DefaultIntegerFieldItem.cpp
        src/main/cpp/org/apache/plc4x/cpp/base/messages/items/DefaultLocalDateFieldItem.cpp
        src/main/cpp/org/apache/plc4x/cpp/base/messages/items/DefaultLocalDateTimeFieldItem.cpp
        src/main/cpp/org/apache/plc4x/cpp/base/messages/items/DefaultLocalTimeFieldItem.cpp
        src/main/cpp/org/apache/plc4x/cpp/base/messages/items/DefaultLongFieldItem.cpp
        src/main/cpp/org/apache/plc4x/cpp/base/messages/items/DefaultShortFieldItem.cpp
        src/main/cpp/org/apache/plc4x/cpp/base/messages/items/DefaultStringFieldItem.cpp
        src/main/cpp/org/apache/plc4x/cpp/base/messages/DefaultPlcReadRequest.cpp
        src/main/cpp/org/apache/plc4x/cpp/base/messages/DefaultPlcWriteRequest.cpp
        src/main/cpp/org/apache/plc4x/cpp/base/messages/InternalPlcFieldRequest.cpp
        src/main/cpp/org/apache/plc4x/cpp/base/messages/InternalPlcFieldResponse.cpp
        src/main/cpp/org/apache/plc4x/cpp/base/messages/InternalPlcProprietaryRequest.cpp
        src/main/cpp/org/apache/plc4x/cpp/base/messages/InternalPlcProprietaryResponse.cpp
        src/main/cpp/org/apache/plc4x/cpp/base/messages/InternalPlcReadRequest.cpp
        src/main/cpp/org/apache/plc4x/cpp/base/messages/InternalPlcRequest.cpp
        src/main/cpp/org/apache/plc4x/cpp/base/messages/InternalPlcResponse.cpp
        src/main/cpp/org/apache/plc4x/cpp/base/messages/InternalPlcWriteRequest.cpp
        src/main/cpp/org/apache/plc4x/cpp/base/messages/PlcProprietaryRequest.cpp
        src/main/cpp/org/apache/plc4x/cpp/base/messages/PlcProprietaryResponse.cpp
        src/main/cpp/org/apache/plc4x/cpp/base/messages/PlcReader.cpp
        src/main/cpp/org/apache/plc4x/cpp/base/messages/PlcWriter.cpp
        )

target_include_directories (plc4cpp-protocols-driver-base-base PUBLIC ${Boost_INCLUDE_DIRS})
target_include_directories (plc4cpp-protocols-driver-base-base PUBLIC ../../../api/src/main/cpp)
target_include_directories (plc4cpp-protocols-driver-base-base PUBLIC ../../../protocols/driver-bases/base/src/main/cpp)
target_include_directories (plc4cpp-protocols-driver-base-base PUBLIC ../../../utils/logger/src/main/cpp)

target_link_libraries (plc4cpp-protocols-driver-base-base
                        ${Boost_LIBRARIES}
                        plc4cpp-utils-logger
                        plc4cpp-api
                      )
