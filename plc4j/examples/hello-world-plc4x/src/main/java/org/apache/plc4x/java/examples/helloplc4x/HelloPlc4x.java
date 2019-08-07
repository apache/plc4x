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
package org.apache.plc4x.java.examples.helloplc4x;

import org.apache.plc4x.java.PlcDriverManager;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.concurrent.CompletableFuture;

public class HelloPlc4x {

    private static final Logger logger = LoggerFactory.getLogger(HelloPlc4x.class);

    /**
     * Example code do demonstrate using PLC4X.
     *
     * @param args ignored.
     */
    public static void main(String[] args) throws Exception {
        // Establish a connection to the plc using the url provided as first argument
        for (int j = 1; j <= 5; j++) {
            try (PlcConnection plcConnection = new PlcDriverManager().getConnection("s7://192.168.167.210/0/0")) {

                for (int i = 0; i <= 10; i++) {
                    try {
                        PlcReadResponse response = plcConnection.readRequestBuilder()
                            .addItem("field", "%DB400:DBW10:INT")
                            .build()
                            .execute()
                            .get();

                        System.out.println(response.getResponseCode("field"));
                    } catch (Exception e) {
                        System.out.println("error...");
                    }
                    Thread.sleep(100);
                }
            }
        }

        System.out.println("The loop is finished");
    }

}
