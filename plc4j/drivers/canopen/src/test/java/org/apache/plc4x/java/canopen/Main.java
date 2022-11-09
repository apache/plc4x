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
package org.apache.plc4x.java.canopen;

import org.apache.plc4x.java.PlcDriverManager;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.messages.PlcWriteResponse;
import org.apache.plc4x.java.canopen.readwrite.CANOpenFrame;
import org.apache.plc4x.java.canopen.context.CANOpenDriverContext;
import org.apache.plc4x.java.canopen.listener.Callback;

import java.util.concurrent.CompletableFuture;

/**
 * Here we begin .. ;-)
 */
public class Main {

    public static void main(String[] args) throws Exception {
        PlcDriverManager driverManager = new PlcDriverManager();

        CANOpenDriverContext.CALLBACK.addCallback(new Callback() {
            @Override
            public void receive(CANOpenFrame frame) {
                //System.err.println("Received frame " + frame);
            }
        });

        PlcConnection connection = driverManager.getConnection("canopen:javacan://vcan0?nodeId=11");

        String value = "abcdef"; //UUID.randomUUID().toString();
        CompletableFuture<? extends PlcWriteResponse> response = connection.writeRequestBuilder()
            .addTagAddress("foo", "SDO:13:0x2000/0x0:VISIBLE_STRING", value)
            .build().execute();

        response.whenComplete((writeReply, writeError) -> {
            System.out.println("====================================");
            if (writeError != null) {
                System.out.println("Error ");
                writeError.printStackTrace();
            } else {
                System.out.println("Result " + writeReply.getResponseCode("foo") + " " + value);

                PlcReadRequest.Builder builder = connection.readRequestBuilder();
                builder.addTagAddress("foo", "SDO:13:0x2000/0x0:VISIBLE_STRING");
                CompletableFuture<? extends PlcReadResponse> future = builder.build().execute();
                future.whenComplete((readReply, readError) -> {
                    System.out.println("====================================");
                    if (readError != null) {
                        System.out.println("Error ");
                        readError.printStackTrace();
                    } else {
                        System.out.println("Result " + readReply.getString("foo"));
                    }
                });
            }
        });


//        while (true) {

//        }
    }

}
