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
package org.apache.plc4x.java.examples.helloplc4x.read;

import org.apache.plc4x.java.PlcDriverManager;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.concurrent.TimeUnit;

public class HelloPlc4xRead {

    private static final Logger logger = LoggerFactory.getLogger(HelloPlc4xRead.class);

    /**
     * Example code do demonstrate using PLC4X.
     *
     * @param args ignored.
     */
    public static void main(String[] args) throws Exception {
        CliOptions options = CliOptions.fromArgs(args);
        if (options == null) {
            CliOptions.printHelp();
            // Could not parse.
            System.exit(1);
        }

        // Establish a connection to the plc using the url provided as first argument
        try (PlcConnection plcConnection = new PlcDriverManager().getConnection(options.getConnectionString())) {

            // Check if this connection support reading of data.
            if (!plcConnection.getMetadata().canRead()) {
                logger.error("This connection doesn't support reading.");
                return;
            }

            // Create a new read request:
            // - Give the single item requested the alias name "value"
            PlcReadRequest.Builder builder = plcConnection.readRequestBuilder();
            for (int i = 0; i < options.getFieldAddress().length; i++) {
                builder.addFieldAddress("value-" + options.getFieldAddress()[i], options.getFieldAddress()[i]);
            }
            PlcReadRequest readRequest = builder.build();

            //////////////////////////////////////////////////////////
            // Read synchronously ...
            // NOTICE: the ".get()" immediately lets this thread pause until
            // the response is processed and available.
            logger.info("Synchronous request ...");
            PlcReadResponse syncResponse = readRequest.execute().get();
            // Simply iterating over the field names returned in the response.
            printResponse(syncResponse);

            /*PlcValue asPlcValue = syncResponse.getAsPlcValue();
            System.out.println(asPlcValue.toString());*/

            //////////////////////////////////////////////////////////
            // Read asynchronously ...
            // Register a callback executed as soon as a response arrives.
            /*logger.info("Asynchronous request ...");
            CompletionStage<? extends PlcReadResponse> asyncResponse = readRequest.execute();
            asyncResponse.whenComplete((readResponse, throwable) -> {
                if (readResponse != null) {
                    printResponse(readResponse);
                } else {
                    logger.error("An error occurred: " + throwable.getMessage(), throwable);
                }
            });*/

            // Give the async request a little time...
            TimeUnit.MILLISECONDS.sleep(1000);
            plcConnection.close();
            System.exit(0);
        }
    }

    private static void printResponse(PlcReadResponse response) {
        for (String fieldName : response.getFieldNames()) {
            if (response.getResponseCode(fieldName) == PlcResponseCode.OK) {
                int numValues = response.getNumberOfValues(fieldName);
                // If it's just one element, output just one single line.
                if (numValues == 1) {
                    logger.info("Value[{}]: {}", fieldName, response.getObject(fieldName));
                }
                // If it's more than one element, output each in a single row.
                else {
                    logger.info("Value[{}]:", fieldName);
                    for (int i = 0; i < numValues; i++) {
                        logger.info(" - {}", response.getObject(fieldName, i));
                    }
                }
            }
            // Something went wrong, to output an error message instead.
            else {
                logger.error("Error[{}]: {}", fieldName, response.getResponseCode(fieldName).name());
            }
        }
    }

}
