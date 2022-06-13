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

add_library(plc4cpp-api
        src/main/cpp/org/apache/plc4x/cpp/PlcDriverManager.cpp
        src/main/cpp/org/apache/plc4x/cpp/api/PlcConnection.cpp
        src/main/cpp/org/apache/plc4x/cpp/api/authentication/PlcAuthentication.cpp
        src/main/cpp/org/apache/plc4x/cpp/api/authentication/PlcUsernamePasswordAuthentication.cpp
        src/main/cpp/org/apache/plc4x/cpp/api/exceptions/PlcConnectionException.cpp
        src/main/cpp/org/apache/plc4x/cpp/api/exceptions/PlcException.cpp
        src/main/cpp/org/apache/plc4x/cpp/api/exceptions/PlcFieldRangeException.cpp
        src/main/cpp/org/apache/plc4x/cpp/api/exceptions/PlcIncompatibleDatatypeException.cpp
        src/main/cpp/org/apache/plc4x/cpp/api/exceptions/PlcInvalidFieldException.cpp
        src/main/cpp/org/apache/plc4x/cpp/api/exceptions/PlcIoException.cpp
        src/main/cpp/org/apache/plc4x/cpp/api/exceptions/PlcNotImplementedException.cpp
        src/main/cpp/org/apache/plc4x/cpp/api/exceptions/PlcProtocolException.cpp
        src/main/cpp/org/apache/plc4x/cpp/api/exceptions/PlcProtocolPayloadTooBigException.cpp
        src/main/cpp/org/apache/plc4x/cpp/api/exceptions/PlcRuntimeException.cpp
        src/main/cpp/org/apache/plc4x/cpp/api/exceptions/PlcTimeoutException.cpp
        src/main/cpp/org/apache/plc4x/cpp/api/exceptions/PlcUnsupportedDataTypeException.cpp
        src/main/cpp/org/apache/plc4x/cpp/api/exceptions/PlcUnsupportedOperationException.cpp
        src/main/cpp/org/apache/plc4x/cpp/api/messages/PlcFieldRequest.cpp
        src/main/cpp/org/apache/plc4x/cpp/api/messages/PlcFieldResponse.cpp
        src/main/cpp/org/apache/plc4x/cpp/api/messages/PlcMessage.cpp
        src/main/cpp/org/apache/plc4x/cpp/api/messages/PlcReadRequest.cpp
        src/main/cpp/org/apache/plc4x/cpp/api/messages/PlcReadResponse.cpp
        src/main/cpp/org/apache/plc4x/cpp/api/messages/PlcRequest.cpp
        src/main/cpp/org/apache/plc4x/cpp/api/messages/PlcRequestBuilder.cpp
        src/main/cpp/org/apache/plc4x/cpp/api/messages/PlcResponse.cpp
        src/main/cpp/org/apache/plc4x/cpp/api/messages/PlcSubscriptionEvent.cpp
        src/main/cpp/org/apache/plc4x/cpp/api/messages/PlcSubscriptionRequest.cpp
        src/main/cpp/org/apache/plc4x/cpp/api/messages/PlcSubscriptionResponse.cpp
        src/main/cpp/org/apache/plc4x/cpp/api/messages/PlcUnsubscriptionRequest.cpp
        src/main/cpp/org/apache/plc4x/cpp/api/messages/PlcUnsubscriptionResponse.cpp
        src/main/cpp/org/apache/plc4x/cpp/api/messages/PlcWriteRequest.cpp
        src/main/cpp/org/apache/plc4x/cpp/api/messages/PlcWriteResponse.cpp
        src/main/cpp/org/apache/plc4x/cpp/api/messages/REQUEST.cpp
        src/main/cpp/org/apache/plc4x/cpp/api/messages/RequestTemplate.cpp
        src/main/cpp/org/apache/plc4x/cpp/api/messages/RESPONSE.cpp
        src/main/cpp/org/apache/plc4x/cpp/api/messages/ResponseTemplate.cpp
        src/main/cpp/org/apache/plc4x/cpp/api/metadata/PlcConnectionMetadata.cpp
        src/main/cpp/org/apache/plc4x/cpp/api/model/Consumer.cpp
        src/main/cpp/org/apache/plc4x/cpp/api/model/PlcConsumerRegistration.cpp
        src/main/cpp/org/apache/plc4x/cpp/api/model/PlcField.cpp
        src/main/cpp/org/apache/plc4x/cpp/api/model/PlcSubscriptionHandle.cpp
        src/main/cpp/org/apache/plc4x/cpp/api/types/PlcClientDataType.cpp
        src/main/cpp/org/apache/plc4x/cpp/api/types/PlcResponseCode.cpp
        src/main/cpp/org/apache/plc4x/cpp/api/types/PlcSubscriptionType.cpp
        src/main/cpp/org/apache/plc4x/cpp/api/types/ValueTypeObject.cpp
        )

target_include_directories (plc4cpp-api PUBLIC ${Boost_INCLUDE_DIRS})
target_include_directories (plc4cpp-api PUBLIC ../utils/logger/src/main/cpp)

target_link_libraries (plc4cpp-api ${CMAKE_DL_LIBS})