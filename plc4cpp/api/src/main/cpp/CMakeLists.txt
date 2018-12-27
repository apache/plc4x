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

add_library(api org/apache/plc4x/cpp/api/PlcConnection.cpp
        org/apache/plc4x/cpp/api/authentication/PlcAuthentication.cpp
        org/apache/plc4x/cpp/api/authentication/PlcUsernamePasswordAuthentication.cpp
        org/apache/plc4x/cpp/api/exceptions/PlcConnectionException.cpp
        org/apache/plc4x/cpp/api/exceptions/PlcException.cpp
        org/apache/plc4x/cpp/api/exceptions/PlcFieldRangeException.cpp
        org/apache/plc4x/cpp/api/exceptions/PlcIncompatibleDatatypeException.cpp
        org/apache/plc4x/cpp/api/exceptions/PlcInvalidFieldException.cpp
        org/apache/plc4x/cpp/api/exceptions/PlcIoException.cpp
        org/apache/plc4x/cpp/api/exceptions/PlcNotImplementedException.cpp
        org/apache/plc4x/cpp/api/exceptions/PlcProtocolException.cpp
        org/apache/plc4x/cpp/api/exceptions/PlcProtocolPayloadTooBigException.cpp
        org/apache/plc4x/cpp/api/exceptions/PlcRuntimeException.cpp
        org/apache/plc4x/cpp/api/exceptions/PlcTimeoutException.cpp
        org/apache/plc4x/cpp/api/exceptions/PlcUnsupportedDataTypeException.cpp
        org/apache/plc4x/cpp/api/exceptions/PlcUnsupportedOperationException.cpp
        org/apache/plc4x/cpp/api/messages/PlcFieldRequest.cpp
        org/apache/plc4x/cpp/api/messages/PlcFieldResponse.cpp
        org/apache/plc4x/cpp/api/messages/PlcMessage.cpp
        org/apache/plc4x/cpp/api/messages/PlcReadRequest.cpp
        org/apache/plc4x/cpp/api/messages/PlcReadRequestBuilder.cpp
        org/apache/plc4x/cpp/api/messages/PlcReadResponse.cpp
        org/apache/plc4x/cpp/api/messages/PlcRequest.cpp
        org/apache/plc4x/cpp/api/messages/PlcRequestBuilder.cpp
        org/apache/plc4x/cpp/api/messages/PlcResponse.cpp
        org/apache/plc4x/cpp/api/messages/PlcSubscriptionEvent.cpp
        org/apache/plc4x/cpp/api/messages/PlcSubscriptionRequest.cpp
        org/apache/plc4x/cpp/api/messages/PlcSubscriptionResponse.cpp
        org/apache/plc4x/cpp/api/messages/PlcUnsubscriptionRequest.cpp
        org/apache/plc4x/cpp/api/messages/PlcUnsubscriptionResponse.cpp
        org/apache/plc4x/cpp/api/messages/PlcWriteRequest.cpp
        org/apache/plc4x/cpp/api/messages/PlcWriteResponse.cpp
        org/apache/plc4x/cpp/api/messages/REQUEST.cpp
        org/apache/plc4x/cpp/api/messages/RequestTemplate.cpp
        org/apache/plc4x/cpp/api/messages/RESPONSE.cpp
        org/apache/plc4x/cpp/api/messages/ResponseTemplate.cpp
        org/apache/plc4x/cpp/api/metadata/PlcConnectionMetadata.cpp
        org/apache/plc4x/cpp/api/model/Consumer.cpp
        org/apache/plc4x/cpp/api/model/PlcConsumerRegistration.cpp
        org/apache/plc4x/cpp/api/model/PlcField.cpp
        org/apache/plc4x/cpp/api/model/PlcSubscriptionHandle.cpp
        org/apache/plc4x/cpp/api/types/PlcClientDataType.cpp
        org/apache/plc4x/cpp/api/types/PlcResponseCode.cpp
        org/apache/plc4x/cpp/api/types/PlcSubscriptionType.cpp
        org/apache/plc4x/cpp/api/types/ValueTypeObject.cpp
        )

target_include_directories (api PUBLIC ../../../../libs/libs/boost/include)

