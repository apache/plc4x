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
package org.apache.plc4x.java.examples.cloud.azure;

import org.apache.commons.cli.*;

public class CliOptions {

    private static Options options;

    private final String plc4xConnectionString;
    private final String plc4xTagAddress;
    private final String iotHubConnectionString;

    public static CliOptions fromArgs(String[] args) {
        options = new Options();
        // Required arguments
        options.addOption(
            Option.builder()
                .type(String.class)
                .longOpt("plc4x-connection-string")
                .hasArg()
                .desc("Connection String")
                .required()
                .build());
        options.addOption(
            Option.builder()
                .type(String.class)
                .longOpt("plc4x-tag-address")
                .hasArg()
                .desc("Tag Address.")
                .required()
                .build());
        options.addOption(
            Option.builder()
                .type(String.class)
                .longOpt("iot-hub-connection-string")
                .hasArg()
                .desc("IoT Hub Connection String.")
                .required()
                .build());

        CommandLineParser parser = new DefaultParser();
        CommandLine commandLine;
        try {
            commandLine = parser.parse(options, args);

            String plc4xConnectionString = commandLine.getOptionValue("plc4x-connection-string");
            String plc4xTagAddress = commandLine.getOptionValue("plc4x-tag-address");
            String iotHubConnectionString = commandLine.getOptionValue("iot-hub-connection-string");

            return new CliOptions(plc4xConnectionString, plc4xTagAddress, iotHubConnectionString);
        } catch (ParseException e) {
            System.err.println(e.getMessage());
            return null;
        }
    }

    public static void printHelp() {
        HelpFormatter formatter = new HelpFormatter();
        formatter.printHelp("S7PlcToAzureIoTHubSample", options);
    }

    public CliOptions(String plc4xConnectionString, String plc4xTagAddress, String iotHubConnectionString) {
        this.plc4xConnectionString = plc4xConnectionString;
        this.plc4xTagAddress = plc4xTagAddress;
        this.iotHubConnectionString = iotHubConnectionString;
    }

    public String getPlc4xConnectionString() {
        return plc4xConnectionString;
    }

    public String getPlc4xTagAddress() {
        return plc4xTagAddress;
    }

    protected String getIotHubConnectionString() {
        return iotHubConnectionString;
    }

}
