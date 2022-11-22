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
package org.apache.plc4x.java.examples.integration.iotdb;

import org.apache.commons.cli.*;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class CliOptions {

    private static final Logger LOGGER = LoggerFactory.getLogger(CliOptions.class);

    private static Options options;

    private final String connectionString;
    private final String tagAddress;
    private final int pollingInterval;
    private final String iotdbIpPort;
    private final String user;
    private final String password;
    private final String storageGroup;
    private final String device;
    private final String datatype;
    private final boolean useJDBC;

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
                .longOpt("tag-address")
                .hasArg()
                .desc("Tag Address.")
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
        options.addOption(
            Option.builder()
                .type(Integer.class)
                .longOpt("iotdb-address")
                .hasArg()
                .desc("The address and port of IoTDB server. format: ip:port")
                .required()
                .build());
        options.addOption(
            Option.builder()
                .type(Integer.class)
                .longOpt("iotdb-user-name")
                .hasArg()
                .desc("The connection user that has privilege to write data into IoTDB")
                .required()
                .build());
        options.addOption(
            Option.builder()
                .type(Integer.class)
                .longOpt("iotdb-user-password")
                .hasArg()
                .desc("The connection user password that has privilege to write data into IoTDB")
                .required()
                .build());
        options.addOption(
            Option.builder()
                .type(Integer.class)
                .longOpt("iotdb-sg")
                .hasArg()
                .desc("The Storage group name, e.g., testapp")
                .required()
                .build());
        options.addOption(
            Option.builder()
                .type(Integer.class)
                .longOpt("iotdb-device")
                .hasArg()
                .desc("The device name, e.g., mitsubishi.D58501")
                .required()
                .build());
        options.addOption(
            Option.builder()
                .type(Integer.class)
                .longOpt("iotdb-datatype")
                .hasArg()
                .desc("The data type of the tag")
                .required()
                .build());
        options.addOption(
            Option.builder()
                .type(Boolean.class)
                .longOpt("use-jdbc")
                .hasArg()
                .desc("Whether use JDBC API or not")
                .build());

        CommandLineParser parser = new DefaultParser();
        CommandLine commandLine;
        try {
            commandLine = parser.parse(options, args);

            String connectionString = commandLine.getOptionValue("connection-string");
            String tagAddress = commandLine.getOptionValue("tag-address");
            int pollingInterval = Integer.parseInt(commandLine.getOptionValue("polling-interval"));
            String iotdbIpPort = commandLine.getOptionValue("iotdb-address");
            String user = commandLine.getOptionValue("iotdb-user-name");
            String password = commandLine.getOptionValue("iotdb-user-password");
            String storageGroup = commandLine.getOptionValue("iotdb-sg");
            String device = commandLine.getOptionValue("iotdb-device");
            String datatype = commandLine.getOptionValue("iotdb-datatype");
            boolean useJDBC = Boolean.valueOf(commandLine.getOptionValue("use-jdbc", "false"));

            return new CliOptions(connectionString, tagAddress, pollingInterval, iotdbIpPort, user, password, storageGroup, device, datatype, useJDBC);
        } catch (ParseException e) {
            LOGGER.error(e.getMessage());
            return null;
        }
    }

    public static void printHelp() {
        HelpFormatter formatter = new HelpFormatter();
        formatter.printHelp("PlcLogger", options);
    }

    public CliOptions(String connectionString, String tagAddress, int pollingInterval, String iotdbIpPort, String user, String password, String storageGroup, String device, String datatype, boolean useJDBC) {
        this.connectionString = connectionString;
        this.tagAddress = tagAddress;
        this.pollingInterval = pollingInterval;
        this.iotdbIpPort = iotdbIpPort;
        this.user = user;
        this.password = password;
        this.storageGroup = storageGroup;
        this.device = device;
        this.datatype = datatype;
        this.useJDBC = useJDBC;
    }

    public String getConnectionString() {
        return connectionString;
    }

    public String getTagAddress() {
        return tagAddress;
    }

    public int getPollingInterval() {
        return pollingInterval;
    }

    public String getIotdbIpPort() {
        return iotdbIpPort;
    }

    public String getUser() {
        return user;
    }

    public String getPassword() {
        return password;
    }

    public String getStorageGroup() {
        return storageGroup;
    }

    public String getDevice() {
        return device;
    }

    public String getDatatype() {
        return datatype;
    }

    public boolean isUseJDBC() {
        return useJDBC;
    }

}
