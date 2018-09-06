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
import org.apache.plc4x.java.api.connection.PlcConnection;
import org.apache.plc4x.java.api.connection.PlcReader;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.Optional;
import java.util.concurrent.CompletableFuture;

public class HelloPlc4x {

    private static final Logger logger = LoggerFactory.getLogger(HelloPlc4x.class);

    /**
     * Example code do demonstrate using PLC4X.
     *
     * @param args ignored.
     */
    @SuppressWarnings("unchecked")
    public static void main(String[] args) {
        if (args.length < 2) {
            System.out.println("Usage: HelloPlc4x {connection-string} {address-string}+");
            System.out.println("Example: HelloPlc4x s7://10.10.64.30/1/1 %I0.0:BOOLEAN %DB1.DBX38:BYTE");
            return;
        }

        // Establish a connection to the plc using the url provided as first argument
        try (PlcConnection plcConnection = new PlcDriverManager().getConnection(args[0])) {

            Optional<PlcReader> reader = plcConnection.getReader();

            // Check if this connection support reading of data.
            if (reader.isPresent()) {
                PlcReader plcReader = reader.get();

                // Create a new read request:
                // - Give the single item requested the alias name "value"
                PlcReadRequest.Builder builder = plcReader.readRequestBuilder();
                for (int i = 1; i < args.length; i++) {
                    builder.addItem("value-" + i, args[i]);
                }
                PlcReadRequest plcReadRequest = builder.build();

                //////////////////////////////////////////////////////////
                // Read synchronously ...
                // NOTICE: the ".get()" immediately lets this thread pause till
                // the response is processed and available.
                PlcReadResponse<?> syncResponse = plcReader.read(plcReadRequest).get();
                // Simply iterating over the field names returned in the response.
                for (String fieldName : syncResponse.getFieldNames()) {
                    System.out.println("Value[" + fieldName + "]: " + syncResponse.getObject(fieldName));
                }

                //////////////////////////////////////////////////////////
                // Read asynchronously ...
                // Register a callback executed as soon as a response arives.
                CompletableFuture<PlcReadResponse<?>> asyncResponse = plcReader.read(plcReadRequest);
                asyncResponse.whenComplete((readResponse, throwable) -> {
                    if (readResponse != null) {
                        // Directly asking for fields by name.
                        for (int i = 1; i < args.length; i++) {
                            System.out.println("Value[value-" + i + "]: " + syncResponse.getObject("value-" + i));
                        }
                    } else {
                        logger.error("An error occurred", throwable);
                    }
                });
            } else {
                logger.info("This connection doesn't support reading.");
            }
        }
        // Catch any exception or the application won't be able to finish if something goes wrong.
        catch (Exception e) {
            logger.error("An error occurred", e);
        }
    }

}
