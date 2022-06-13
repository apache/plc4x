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
package org.apache.plc4x.java.examples.integration.edgent;

import org.apache.commons.cli.*;

public class CliOptions {

    private static Options options;

    private final String connectionString;
    private final String fieldAddress;
    private final int pollingInterval;

    public static CliOptions fromArgs(String[] args) {
        options = new Options();
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
                .longOpt("field-address")
                .hasArg()
                .desc("Field Address.")
                .required()
                .build());
        options.addOption(
            Option.builder()
                .type(Integer.class)
                .longOpt("polling-interval")
                .hasArg()
                .desc("Polling Interval (milliseconds).")
                .required()
                .build());

        CommandLineParser parser = new DefaultParser();
        CommandLine commandLine;
        try {
            commandLine = parser.parse(options, args);

            String connectionString = commandLine.getOptionValue("connection-string");
            String fieldAddress = commandLine.getOptionValue("field-address");
            int pollingInterval = Integer.parseInt(commandLine.getOptionValue("polling-interval"));

            return new CliOptions(connectionString, fieldAddress, pollingInterval);
        } catch (ParseException e) {
            System.err.println(e.getMessage());
            return null;
        }
    }

    public static void printHelp() {
        HelpFormatter formatter = new HelpFormatter();
        formatter.printHelp("PlcLogger", options);
    }

    public CliOptions(String connectionString, String fieldAddress, int pollingInterval) {
        this.connectionString = connectionString;
        this.fieldAddress = fieldAddress;
        this.pollingInterval = pollingInterval;
    }

    public String getConnectionString() {
        return connectionString;
    }

    public String getFieldAddress() {
        return fieldAddress;
    }

    public int getPollingInterval() {
        return pollingInterval;
    }

}
