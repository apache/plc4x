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
import org.apache.plc4x.java.api.connection.PlcConnection;
import org.apache.plc4x.java.api.connection.PlcReader;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.messages.items.PlcReadRequestItem;
import org.apache.plc4x.java.api.messages.items.PlcReadResponseItem;
import org.apache.plc4x.java.api.model.PlcField;

import java.util.concurrent.CompletableFuture;

public class ManualPlc4XEtherNetIpTest {

    public static void main(String... args) {
        // Connection to our IoT Lab WAGO device.
        String connectionUrl = "eip://10.10.64.30:44818";
        try (PlcConnection plcConnection = new PlcDriverManager().getConnection(connectionUrl)) {
            System.out.println("PlcConnection " + plcConnection);

            // Get a reader instance.
            PlcReader reader = plcConnection.getReader().orElseThrow(() -> new RuntimeException("No Reader found"));

            // Parse the address.
            PlcField field = plcConnection.prepareField("#4#105#3");

            // Create a new read request.
            PlcReadRequest readRequest = new PlcReadRequest(new PlcReadRequestItem<>(Integer.class, field));

            // Execute the read operation.
            CompletableFuture<? extends PlcReadResponse> response = reader.read(readRequest);
            PlcReadResponse readResponse = response.get();

            // Output the response.
            PlcReadResponseItem responseItem = readResponse.getResponseItem()
                .orElseThrow(() -> new RuntimeException("No Item found"));
            System.out.println("ResponseItem " + responseItem);
            responseItem.getValues().stream().map(value -> "Value: " + value.toString()).forEach(System.out::println);
        } catch (Exception e) {
            e.printStackTrace();
            System.exit(1);
        }
        System.exit(0);
    }
}
