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
package org.apache.plc4x.java.examples.hellonats;

import org.apache.commons.cli.*;

public class CliOptions {

    private static Options options;

    private final String natsServerConnectionString;
    private final String natsNodeName;
    private final String natsTopic;

    public static CliOptions fromArgs(String[] args) {
        options = new Options();
        // Required arguments
        options.addOption(
            Option.builder()
                .type(String.class)
                .longOpt("nats-server-connection-string")
                .hasArg()
                .desc("Nats Server Connection String")
                .required()
                .build());
        options.addOption(
            Option.builder()
                .type(String.class)
                .longOpt("nats-node-name")
                .hasArg()
                .desc("Name of this node.")
                .required()
                .build());
        options.addOption(
            Option.builder()
                .type(String.class)
                .longOpt("nats-topic")
                .hasArg()
                .desc("Name of the topic for receiving requests.")
                .required()
                .build());

        CommandLineParser parser = new DefaultParser();
        CommandLine commandLine;
        try {
            commandLine = parser.parse(options, args);

            String natsServerConnectionString = commandLine.getOptionValue("nats-server-connection-string");
            String natsNodeName = commandLine.getOptionValue("nats-node-name");
            String natsTopic = commandLine.getOptionValue("nats-topic");

            return new CliOptions(natsServerConnectionString, natsNodeName, natsTopic);
        } catch (ParseException e) {
            System.err.println(e.getMessage());
            return null;
        }
    }

    public static void printHelp() {
        HelpFormatter formatter = new HelpFormatter();
        formatter.printHelp("HelloPlc4x", options);
    }

    public CliOptions(String natsServerConnectionString, String natsNodeName, String natsTopic) {
        this.natsServerConnectionString = natsServerConnectionString;
        this.natsNodeName = natsNodeName;
        this.natsTopic = natsTopic;
    }

    public String getNatsServerConnectionString() {
        return natsServerConnectionString;
    }

    public String getNatsNodeName() {
        return natsNodeName;
    }

    public String getNatsTopic() {
        return natsTopic;
    }

}
