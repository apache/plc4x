/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements. See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership. The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License. You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

# Thrift Definition for PLC4X Interop Server

/**
 * The first thing to know about are types. The available types in Thrift are:
 *
 *  bool        Boolean, one byte
 *  i8 (byte)   Signed 8-bit integer
 *  i16         Signed 16-bit integer
 *  i32         Signed 32-bit integer
 *  i64         Signed 64-bit integer
 *  double      64-bit floating point value
 *  string      String
 *  binary      Blob (byte array)
 *  map<t1,t2>  Map from one type to another
 *  list<t1>    Ordered list of one type
 *  set<t1>     Set of unique elements of one type
 *
 * Did you also notice that Thrift supports C style comments?
 */

/**
 * Thrift files can namespace, package, or prefix their output in various
 * target languages.
 */

namespace java org.apache.plc4x.interop
namespace py org.apache.plc4x.interop
namespace csharp org.apache.plc4x.interop


enum RESPONSE_CODE {
  OK = 1,
  NOT_FOUND = 2,
  ACCESS_DENIED = 3,
  INVALID_ADDRESS = 4,
  INVALID_DATATYPE = 5,
  INTERNAL_ERROR = 6,
  RESPONSE_PENDING = 7
}

struct ConnectionHandle {
  1: i64 connectionId
}

exception PlcException {
  1: string url,
  2: string exceptionString
}

struct Request {
  1: map<string, string> fields
}

struct FieldResponse {
  1: RESPONSE_CODE responseCode,
  2: optional bool boolValue,
  3: optional i64 longValue,
  4: optional double doubleValue,
  5: optional string stringValue
}

struct Response {
  1: map<string, FieldResponse> fields
}

service InteropServer {

   ConnectionHandle connect(1: string connectionString) throws (1: PlcException connectionException),

   Response execute(1: ConnectionHandle handle, 2: Request request) throws (1: PlcException executionException),

   void close(1: ConnectionHandle handle)

}
