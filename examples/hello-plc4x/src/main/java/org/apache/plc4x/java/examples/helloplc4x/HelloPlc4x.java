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

import java.util.Optional;
import java.util.concurrent.CompletableFuture;

public class HelloPlc4x {

    /**
     * Example code do demonstrate using PLC4X.
     *
     * @param args ignored.
     */
public static void main(String[] args) throws Exception {
    if (args.length < 2) {
        System.out.println("Usage: HelloPlc4x {connection-string} {address-string}+");
        System.out.println("Example: HelloPlc4x s7://10.10.64.20/1/1 %Q0.0:BOOL %Q0:BYTE");
        return;
    }

    // Establish a connection to the plc using the url provided as first argument
    try (PlcConnection plcConnection = new PlcDriverManager().getConnection(args[0])) {

        // Check if this connection support reading of data.
        Optional<PlcReadRequest.Builder> builderOptional = plcConnection.readRequestBuilder();
        if (!builderOptional.isPresent()) {
            System.err.println("This connection doesn't support reading.");
            return;
        }

        // Create a new read request:
        // - Give the single item requested the alias name "value"
        PlcReadRequest.Builder builder = builderOptional.get();
        for (int i = 1; i < args.length; i++) {
            builder.addItem("value-" + i, args[i]);
        }
        PlcReadRequest readRequest = builder.build();

        //////////////////////////////////////////////////////////
        // Read synchronously ...
        // NOTICE: the ".get()" immediately lets this thread pause until
        // the response is processed and available.
        System.out.println("Synchronous request ...");
        PlcReadResponse syncResponse = readRequest.execute().get();
        // Simply iterating over the field names returned in the response.
        printResponse(syncResponse);

        //////////////////////////////////////////////////////////
        // Read asynchronously ...
        // Register a callback executed as soon as a response arrives.
        System.out.println("Asynchronous request ...");
        CompletableFuture<? extends PlcReadResponse> asyncResponse = readRequest.execute();
        asyncResponse.whenComplete((readResponse, throwable) -> {
            if (readResponse != null) {
                printResponse(syncResponse);
            } else {
                System.err.println("An error occurred: " + throwable.getMessage());
                throwable.printStackTrace();
            }
        });
    }
}

private static void printResponse(PlcReadResponse response) {
    for (String fieldName : response.getFieldNames()) {
        if(response.getResponseCode(fieldName) == PlcResponseCode.OK) {
            int numValues = response.getNumberOfValues(fieldName);
            // If it's just one element, output just one single line.
            if(numValues == 1) {
                System.out.println("Value[" + fieldName + "]: " + response.getObject(fieldName));
            }
            // If it's more than one element, output each in a single row.
            else {
                System.out.println("Value[" + fieldName + "]:");
                for(int i = 0; i < numValues; i++) {
                    System.out.println(" - " + response.getObject(fieldName, i));
                }
            }
        }
        // Something went wrong, to output an error message instead.
        else {
            System.out.println("Error[" + fieldName + "]: " + response.getResponseCode(fieldName).name());
        }
    }
}

}
