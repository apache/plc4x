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
package org.apache.plc4x.java.examples.helloplc4x.write;

import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.PlcDriverManager;
import org.apache.plc4x.java.api.messages.*;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class HelloPlc4xWrite {

    private static final Logger LOGGER = LoggerFactory.getLogger(HelloPlc4xWrite.class);

    private final CliOptions options;

    public HelloPlc4xWrite(CliOptions options) {
        this.options = options;
    }

    public void run() throws Exception {
        // Establish a connection to the plc.
        try (PlcConnection plcConnection = PlcDriverManager.getDefault().getConnectionManager().getConnection(options.getConnectionString())) {

            // Check if this connection support subscriptions.
            if (!plcConnection.getMetadata().canWrite()) {
                LOGGER.error("This connection doesn't support writing.");
                return;
            }

            if(options.getTagValues().length != options.getTagAddress().length) {
                LOGGER.error("The number of values doesn't match the number of addresses.");
                return;
            }

            // Create a new read request:
            // - Give the single item requested the alias name "value"
            final PlcWriteRequest.Builder builder = plcConnection.writeRequestBuilder();
            for (int i = 0; i < options.getTagAddress().length; i++) {
                //If an array value is passed instead of a single value then convert to a String array
                if ((options.getTagValues()[i].charAt(0) == '[') && (options.getTagValues()[i].charAt(options.getTagValues()[i].length() - 1) == ']')) {
                    String[] values = options.getTagValues()[i].substring(1,options.getTagValues()[i].length() - 1).split(",");
                    builder.addTagAddress("value-" + i, options.getTagAddress()[i], values);
                } else {
                    builder.addTagAddress("value-" + i, options.getTagAddress()[i], options.getTagValues()[i]);
                }
            }
            PlcWriteRequest writeRequest = builder.build();

            // Execute the write request.
            final PlcWriteResponse writeResponse = writeRequest.execute().get();

            // Attach handlers for the incoming data.
            for (String tagName : writeResponse.getTagNames()) {
                LOGGER.info(String.format("Return code for %s was %s",
                    tagName, writeResponse.getResponseCode(tagName)));
            }
        }
    }

    /**
     * Example code do demonstrate using PLC4X Subcription API.
     */
    public static void main(String[] args) throws Exception {
        CliOptions options = CliOptions.fromArgs(args);
        if (options == null) {
            CliOptions.printHelp();
            // Could not parse.
            System.exit(1);
        }

        HelloPlc4xWrite subscriptionApplication = new HelloPlc4xWrite(options);

        subscriptionApplication.run();

        System.exit(0);
    }

}
