/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package org.apache.plc4x.java.plc4x;

import org.apache.plc4x.java.DefaultPlcDriverManager;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;

public class Plc4xReadClient {

    public static void main(String[] args) throws Exception {
        try (final PlcConnection connection = new DefaultPlcDriverManager().getConnection("plc4x://localhost?remote-connection-string=simulated%3A%2F%2Flocalhost")) {
            final PlcReadRequest.Builder requestBuilder = connection.readRequestBuilder();
            requestBuilder.addTagAddress("test-BOOL", "RANDOM/foo:BOOL");
            requestBuilder.addTagAddress("test-BYTE", "RANDOM/foo:BYTE");
            requestBuilder.addTagAddress("test-WORD", "RANDOM/foo:WORD");
            requestBuilder.addTagAddress("test-DWORD", "RANDOM/foo:DWORD");
            requestBuilder.addTagAddress("test-USINT", "RANDOM/foo:USINT");
            requestBuilder.addTagAddress("test-UINT", "RANDOM/foo:UINT");
            requestBuilder.addTagAddress("test-UDINT", "RANDOM/foo:UDINT");
            requestBuilder.addTagAddress("test-ULINT", "RANDOM/foo:ULINT");
            requestBuilder.addTagAddress("test-SINT", "RANDOM/foo:SINT");
            requestBuilder.addTagAddress("test-INT", "RANDOM/foo:INT");
            requestBuilder.addTagAddress("test-DINT", "RANDOM/foo:DINT");
            requestBuilder.addTagAddress("test-LINT", "RANDOM/foo:LINT");
            requestBuilder.addTagAddress("test-REAL", "RANDOM/foo:REAL");
            requestBuilder.addTagAddress("test-LREAL", "RANDOM/foo:LREAL");
            requestBuilder.addTagAddress("test-CHAR", "RANDOM/foo:CHAR");
            requestBuilder.addTagAddress("test-WCHAR", "RANDOM/foo:WCHAR");
            final PlcReadRequest readRequest = requestBuilder.build();
            final PlcReadResponse readResponse = readRequest.execute().get();
            System.out.println(readResponse);
        }
    }

}
