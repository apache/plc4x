/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package org.apache.plc4x.java.plc4x;

import org.apache.plc4x.java.PlcDriverManager;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.messages.PlcWriteRequest;
import org.apache.plc4x.java.api.messages.PlcWriteResponse;

import java.math.BigInteger;

public class Plc4xWriteClient {

    public static void main(String[] args) throws Exception {
        try (final PlcConnection connection = new PlcDriverManager().getConnection("plc4x://localhost?remote-connection-string=simulated%3A%2F%2Flocalhost")) {
            final PlcWriteRequest.Builder requestBuilder = connection.writeRequestBuilder();
            requestBuilder.addItem("test-BOOL", "STDOUT/foo:BOOL", true);
            requestBuilder.addItem("test-BYTE", "STDOUT/foo:BYTE", new boolean[] {true, true, false, true, false, true, false, true});
            requestBuilder.addItem("test-WORD", "STDOUT/foo:WORD", new boolean[] {true, true, false, true, false, true, false, true, false, false, false, false, true, true, true, true});
            requestBuilder.addItem("test-DWORD", "STDOUT/foo:DWORD", new boolean[] {true, true, false, true, false, true, false, true, false, false, false, false, true, true, true, true, true, false, true, false, true, false, true, false, true, false, true, false, true, false, true, false});
            requestBuilder.addItem("test-USINT", "STDOUT/foo:USINT", 12);
            requestBuilder.addItem("test-UINT", "STDOUT/foo:UINT", 12345);
            requestBuilder.addItem("test-UDINT", "STDOUT/foo:UDINT", 1234567890);
            requestBuilder.addItem("test-ULINT", "STDOUT/foo:ULINT", new BigInteger("12345678901234567"));
            requestBuilder.addItem("test-SINT", "STDOUT/foo:SINT", -12);
            requestBuilder.addItem("test-INT", "STDOUT/foo:INT", -12345);
            requestBuilder.addItem("test-DINT", "STDOUT/foo:DINT", -1234567890);
            requestBuilder.addItem("test-LINT", "STDOUT/foo:LINT", new BigInteger("-12345678901234567"));
            requestBuilder.addItem("test-REAL", "STDOUT/foo:REAL", 3.14159f);
            requestBuilder.addItem("test-LREAL", "STDOUT/foo:LREAL", 2.71828182845904523536028747135d);
            requestBuilder.addItem("test-CHAR", "STDOUT/foo:CHAR", "P");
            requestBuilder.addItem("test-WCHAR", "STDOUT/foo:WCHAR", "Ϡ");
            final PlcWriteRequest writeRequest = requestBuilder.build();
            final PlcWriteResponse writeResponse = writeRequest.execute().get();
            System.out.println(writeResponse);
        }
    }

}
