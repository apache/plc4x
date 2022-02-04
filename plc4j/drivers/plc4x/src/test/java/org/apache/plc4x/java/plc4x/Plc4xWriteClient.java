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

public class Plc4xWriteClient {

    public static void main(String[] args) throws Exception {
        try (final PlcConnection connection = new PlcDriverManager().getConnection("plc4x://localhost?remote-connection-string=simulated%3A%2F%2Flocalhost")) {
            final PlcWriteRequest writeRequest = connection.writeRequestBuilder().addItem("lalala", "STDOUT/foo:INT", (short) 23).build();
            final PlcWriteResponse writeResponse = writeRequest.execute().get();
            System.out.println(writeResponse);
        }
    }

}
