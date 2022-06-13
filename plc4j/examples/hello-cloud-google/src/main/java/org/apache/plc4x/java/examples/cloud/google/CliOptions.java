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
package org.apache.plc4x.java.examples.cloud.google;

import org.apache.commons.cli.*;

/** Command line options for the MQTT example. */
public class CliOptions {

    private static Options options;

    private final String projectId;
    private final String registryId;
    private final String deviceId;
    private final String privateKeyFile;
    private final String algorithm;
    private final String cloudRegion;
    private final int tokenExpMins;
    private final String mqttBridgeHostname;
    private final short mqttBridgePort;
    private final String messageType;

    public static CliOptions fromArgs(String[] args) {
        options = new Options();
        // Required arguments
        options.addOption(
            Option.builder()
                .type(String.class)
                .longOpt("project-id")
                .hasArg()
                .desc("GCP cloud project name.")
                .required()
                .build());
        options.addOption(
            Option.builder()
                .type(String.class)
                .longOpt("registry-id")
                .hasArg()
                .desc("Cloud IoT Core registry id.")
                .required()
                .build());
        options.addOption(
            Option.builder()
                .type(String.class)
                .longOpt("device-id")
                .hasArg()
                .desc("Cloud IoT Core device id.")
                .required()
                .build());
        options.addOption(
            Option.builder()
                .type(String.class)
                .longOpt("private-key-file")
                .hasArg()
                .desc("Path to private key file.")
                .required()
                .build());
        options.addOption(
            Option.builder()
                .type(String.class)
                .longOpt("algorithm")
                .hasArg()
                .desc("Encryption algorithm to use to generate the JWT. Either 'RS256' or 'ES256'.")
                .required()
                .build());

        // Optional arguments.
        options.addOption(
            Option.builder()
                .type(String.class)
                .longOpt("cloud-region")
                .hasArg()
                .desc("GCP cloud region.")
                .build());
        options.addOption(
            Option.builder()
                .type(String.class)
                .longOpt("mqtt-bridge-hostname")
                .hasArg()
                .desc("MQTT bridge hostname.")
                .build());
        options.addOption(
            Option.builder()
                .type(Number.class)
                .longOpt("token-exp-minutes")
                .hasArg()
                .desc("Minutes to JWT token refresh (token expiration time).")
                .build());
        options.addOption(
            Option.builder()
                .type(Number.class)
                .longOpt("mqtt-bridge-port")
                .hasArg()
                .desc("MQTT bridge port.")
                .build());
        options.addOption(
            Option.builder()
                .type(String.class)
                .longOpt("message-type")
                .hasArg()
                .desc("Indicates whether the message is a telemetry event or a device state message")
                .build());

        CommandLineParser parser = new DefaultParser();
        CommandLine commandLine;
        try {
            commandLine = parser.parse(options, args);

            String projectId = commandLine.getOptionValue("project-id");
            String registryId = commandLine.getOptionValue("registry-id");
            String deviceId = commandLine.getOptionValue("device-id");
            String privateKeyFile = commandLine.getOptionValue("private-key-file");
            String algorithm = commandLine.getOptionValue("algorithm");
            String cloudRegion = "europe-west1";
            if (commandLine.hasOption("cloud-region")) {
                cloudRegion = commandLine.getOptionValue("cloud-region");
            }
            int tokenExpMins = 20;
            if (commandLine.hasOption("token-exp-minutes")) {
                tokenExpMins =
                    ((Number) commandLine.getParsedOptionValue("token-exp-minutes")).intValue();
            }
            String mqttBridgeHostname = "mqtt.googleapis.com";
            if (commandLine.hasOption("mqtt-bridge-hostname")) {
                mqttBridgeHostname = commandLine.getOptionValue("mqtt-bridge-hostname");
            }
            short mqttBridgePort = 8883;
            if (commandLine.hasOption("mqtt-bridge-port")) {
                mqttBridgePort =
                    ((Number) commandLine.getParsedOptionValue("mqtt-bridge-port")).shortValue();
            }
            String messageType = "event";
            if (commandLine.hasOption("message-type")) {
                messageType = commandLine.getOptionValue("message-type");
            }

            return new CliOptions(projectId, registryId, deviceId, privateKeyFile, algorithm, cloudRegion, tokenExpMins,
                mqttBridgeHostname, mqttBridgePort, messageType);
        } catch (ParseException e) {
            System.err.println(e.getMessage());
            return null;
        }
    }

    public static void printHelp() {
        HelpFormatter formatter = new HelpFormatter();
        formatter.printHelp("S7PlcToGoogleIoTCoreSample", options);
    }

    public CliOptions(String projectId, String registryId, String deviceId, String privateKeyFile, String algorithm,
                      String cloudRegion, int tokenExpMins, String mqttBridgeHostname, short mqttBridgePort,
                      String messageType) {
        this.projectId = projectId;
        this.registryId = registryId;
        this.deviceId = deviceId;
        this.privateKeyFile = privateKeyFile;
        this.algorithm = algorithm;
        this.cloudRegion = cloudRegion;
        this.tokenExpMins = tokenExpMins;
        this.mqttBridgeHostname = mqttBridgeHostname;
        this.mqttBridgePort = mqttBridgePort;
        this.messageType = messageType;
    }

    public String getProjectId() {
        return projectId;
    }

    public String getRegistryId() {
        return registryId;
    }

    public String getDeviceId() {
        return deviceId;
    }

    public String getPrivateKeyFile() {
        return privateKeyFile;
    }

    public String getAlgorithm() {
        return algorithm;
    }

    public String getCloudRegion() {
        return cloudRegion;
    }

    public int getTokenExpMins() {
        return tokenExpMins;
    }

    public String getMqttBridgeHostname() {
        return mqttBridgeHostname;
    }

    public short getMqttBridgePort() {
        return mqttBridgePort;
    }

    public String getMessageType() {
        return messageType;
    }

}