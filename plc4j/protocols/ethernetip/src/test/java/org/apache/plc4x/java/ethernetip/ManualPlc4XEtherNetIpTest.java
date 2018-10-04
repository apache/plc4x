/*
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
 */
package org.apache.plc4x.java.ethernetip;

import org.apache.plc4x.java.PlcDriverManager;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;

import java.util.concurrent.CompletableFuture;

public class ManualPlc4XEtherNetIpTest {

    public static void main(String... args) {
        // Connection to our IoT Lab WAGO device.
        String connectionUrl = "eip://10.10.64.30:44818";
        try (PlcConnection plcConnection = new PlcDriverManager().getConnection(connectionUrl)) {
            System.out.println("PlcConnection " + plcConnection);

            PlcReadRequest readRequest = plcConnection.readRequestBuilder().orElseThrow(() -> new RuntimeException("Reading not supported"))
                .addItem("field", "#4#105#3").build();

            // Execute the read operation.
            CompletableFuture<? extends PlcReadResponse> response = readRequest.execute();
            PlcReadResponse readResponse = response.get();

            // Output the response.
            for (String fieldName : readResponse.getFieldNames()) {
                readResponse.getAllLongs(fieldName).stream().map(
                    value -> "Value: " + value.toString()).forEach(System.out::println);
            }
        } catch (Exception e) {
            e.printStackTrace();
            System.exit(1);
        }
        System.exit(0);
    }

}
