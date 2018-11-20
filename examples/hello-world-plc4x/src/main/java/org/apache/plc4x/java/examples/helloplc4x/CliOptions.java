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

import org.apache.commons.cli.*;

public class CliOptions {

    private final String connectionString;
    private final String[] fieldAddress;

    /** Construct an MqttExampleOptions class from command line flags. */
    public static CliOptions fromFlags(String[] args) {
        org.apache.commons.cli.Options options = new org.apache.commons.cli.Options();
        // Required arguments
        options.addOption(
            Option.builder()
                .type(String.class)
                .longOpt("connection-string")
                .hasArg()
                .desc("Connection String")
                .required()
                .build());
        options.addOption(
            Option.builder()
                .type(String.class)
                .longOpt("field-addresses")
                .hasArg()
                .desc("Field Addresses (Space separated).")
                .required()
                .build());

        CommandLineParser parser = new DefaultParser();
        CommandLine commandLine;
        try {
            commandLine = parser.parse(options, args);

            String connectionString = commandLine.getOptionValue("connection-string");
            String[] fieldAddress = commandLine.getOptionValues("field-addresses");

            return new CliOptions(connectionString, fieldAddress);
        } catch (ParseException e) {
            System.err.println(e.getMessage());
            return null;
        }
    }

    public CliOptions(String connectionString, String[] fieldAddress) {
        this.connectionString = connectionString;
        this.fieldAddress = fieldAddress;
    }

    public String getConnectionString() {
        return connectionString;
    }

    public String[] getFieldAddress() {
        return fieldAddress;
    }

}
