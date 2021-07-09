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

package org.apache.plc4x.java.examples.cloud.aws;

import org.apache.commons.cli.*;

import java.util.UUID;

/**
 * Command line options for the MQTT example.
 */
public class CliOptions {

    private static Options options;

    //Required Argumnets
    private final String endpoint;
    private final String rootca;
    private final String cert;
    private final String key;
    private final String userName;
    private final String password;
    private final String connectionString;
    private final String[] fieldAddress;

    //Optional Argumnets
    private final String clientId;
    private final String topic;
    private final String message;
    private final int messagesToPublish;
    private final String protocolName;
    private final String[] authParams;



    public static CliOptions fromArgs(String[] args) {
        options = new Options();
        // Required arguments
        options.addOption(
            Option.builder()
                .type(String.class)
                .longOpt("endpoint")
                .hasArg()
                .desc("AWS IoT service endpoint hostname")
                .required()
                .build());
        options.addOption(
            Option.builder()
                .type(String.class)
                .longOpt("rootca")
                .hasArg()
                .desc("Path to the root certificate")
                .required()
                .build());
        options.addOption(
            Option.builder()
                .type(String.class)
                .longOpt("cert")
                .hasArg()
                .desc("Path to the IoT thing certificate")
                .required()
                .build());
        options.addOption(
            Option.builder()
                .type(String.class)
                .longOpt("key")
                .hasArg()
                .desc("Path to the IoT thing private key")
                .required()
                .build());
        options.addOption(
            Option.builder()
                .type(String.class)
                .longOpt("username")
                .hasArg()
                .desc("Username to use as part of the connection/authentication process")
                .required()
                .build());
        options.addOption(
            Option.builder()
                .type(String.class)
                .longOpt("password")
                .hasArg()
                .desc("Password to use as part of the connection/authentication process")
                .required()
                .build());
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
                .hasArgs()
                .desc("Field Addresses (Space separated).")
                .required()
                .build());

        // Optional arguments.
        options.addOption(
            Option.builder()
                .type(String.class)
                .longOpt("clientId")
                .hasArg()
                .desc("Client ID to use when connecting (optional)")
                .build());
        options.addOption(
            Option.builder()
                .type(String.class)
                .longOpt("topic")
                .hasArg()
                .desc("Topic to subscribe/publish to (optional)")
                .build());
        options.addOption(
            Option.builder()
                .type(String.class)
                .longOpt("message")
                .hasArg()
                .desc("Message to publish (optional)")
                .build());
        options.addOption(
            Option.builder()
                .type(Number.class)
                .longOpt("count")
                .hasArg()
                .desc("Number of messages to publish (optional)")
                .build());
        options.addOption(
            Option.builder()
                .type(String.class)
                .longOpt("protocol")
                .hasArg()
                .desc("Communication protocol to use; defaults to mqtt (optional)")
                .build());
        options.addOption(
            Option.builder()
                .type(String.class)
                .longOpt("auth_params")
                .hasArg()
                .desc("Comma delimited list of auth parameters. For websockets these will be set as headers. For raw mqtt these will be appended to user_name. (optional)")
                .build());

        CommandLineParser parser = new DefaultParser();
        CommandLine commandLine;
        try {
            commandLine = parser.parse(options, args);

            String endpoint = commandLine.getOptionValue("endpoint");
            String rootca = commandLine.getOptionValue("rootca");
            String cert = commandLine.getOptionValue("cert");
            String key = commandLine.getOptionValue("key");
            String userName = commandLine.getOptionValue("username");
            String password = commandLine.getOptionValue("password");
            String connectionString = commandLine.getOptionValue("connection-string");
            String[] fieldAddress = commandLine.getOptionValues("field-addresses");


            String clientId = "test-" + UUID.randomUUID().toString();
            if (commandLine.hasOption("clientId")) {
                clientId = commandLine.getOptionValue("clientId");
            }

            String topic = "test/topic";
            if (commandLine.hasOption("topic")) {
                topic = commandLine.getOptionValue("topic");
            }

            String message = "Hello Plc4x World!";
            if (commandLine.hasOption("message")) {
                message = commandLine.getOptionValue("message");
            }

            int messagesToPublish = 10;
            if (commandLine.hasOption("count")) {
                messagesToPublish = ((Number) commandLine.getParsedOptionValue("count")).intValue();
            }

            String protocolName = "mqtt";
            if (commandLine.hasOption("protocol")) {
                protocolName = commandLine.getOptionValue("protocol");
            }

            String[] authParams = null;
            if (commandLine.hasOption("auth_params")) {
                authParams = commandLine.getOptionValues("auth_params");
            }
            return new CliOptions(endpoint,rootca,cert,key,userName,password,connectionString,fieldAddress,clientId,topic,message,messagesToPublish,protocolName,authParams);

        } catch (ParseException e) {
            System.err.println(e.getMessage());
            return null;
        }
    }

    public static void printHelp() {
        HelpFormatter formatter = new HelpFormatter();
        formatter.printHelp("AWSIoTCoreSample", options);
    }

    public CliOptions(String endpoint, String rootca, String cert, String key, String userName, String password, String connectionString, String[] fieldAddress, String clientId, String topic, String message, int messagesToPublish, String protocolName, String[] authParams) {
        this.endpoint = endpoint;
        this.rootca = rootca;
        this.cert = cert;
        this.key = key;
        this.userName = userName;
        this.password = password;
        this.connectionString = connectionString;
        this.fieldAddress = fieldAddress;
        this.clientId = clientId;
        this.topic = topic;
        this.message = message;
        this.messagesToPublish = messagesToPublish;
        this.protocolName = protocolName;
        this.authParams = authParams;
    }


    public String getEndpoint() {
        return endpoint;
    }

    public String getRootca() {
        return rootca;
    }

    public String getCert() {
        return cert;
    }

    public String getKey() {
        return key;
    }

    public String getUserName() { return userName; }

    public String getPassword() {
        return password;
    }

    public String getClientId() {
        return clientId;
    }

    public String getTopic() {
        return topic;
    }

    public String getMessage() {
        return message;
    }

    public int getMessagesToPublish() {
        return messagesToPublish;
    }

    public String getProtocolName() {
        return protocolName;
    }

    public String[] getAuthParams() {
        return authParams;
    }

    public String getConnectionString() { return connectionString;    }

    public String[] getFieldAddress() { return fieldAddress;    }


}