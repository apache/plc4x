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
package org.apache.plc4x.java.utils.capturereplay;

import org.apache.commons.cli.*;

import java.io.File;

public class CliOptions {

    private static final String OPTION_INPUT_FILE = "input-file";
    private static final String OPTION_OUTPUT_DEVICE = "output-device";
    private static final String OPTION_REPLAY_SPEED = "replay-speed";
    private static final String OPTION_LOOP = "loop";

    private static Options options;

    private final File inputFile;
    private final String outputDevice;
    private final float replaySpeed;
    private final boolean loop;

    public static CliOptions fromArgs(String[] args) {
        options = new Options();
        // Required arguments
        options.addOption(
            Option.builder()
                .type(String.class)
                .longOpt(OPTION_INPUT_FILE)
                .hasArg()
                .desc("Path to the PCAP(NG) file.")
                .required()
                .build());
        options.addOption(
            Option.builder()
                .type(String.class)
                .longOpt(OPTION_OUTPUT_DEVICE)
                .hasArgs()
                .desc("Name of the device that should output the packets")
                .required()
                .build());
        options.addOption(
            Option.builder()
                .type(Float.class)
                .longOpt(OPTION_REPLAY_SPEED)
                .hasArgs()
                .desc("Replay speed (1 = real time, 0 = as fast as possible, 0.5 = half speed, 2 = double speed")
                .build());
        options.addOption(
            Option.builder()
                .type(Boolean.class)
                .longOpt(OPTION_LOOP)
                .hasArgs()
                .desc("If set to TRUE it will start sending the packets as soon as it reaches the end")
                .build());

        CommandLineParser parser = new DefaultParser();
        CommandLine commandLine;
        try {
            commandLine = parser.parse(options, args);

            File inputFile = new File(commandLine.getOptionValue(OPTION_INPUT_FILE));
            String outputDevice = commandLine.getOptionValue(OPTION_OUTPUT_DEVICE);
            float replaySpeed = Float.parseFloat(commandLine.getOptionValue(OPTION_REPLAY_SPEED, "1"));
            boolean loop = Boolean.parseBoolean(commandLine.getOptionValue(OPTION_LOOP, "false"));

            return new CliOptions(inputFile, outputDevice, replaySpeed, loop);
        } catch (ParseException e) {
            System.err.println(e.getMessage());
            return null;
        }
    }

    public static void printHelp() {
        HelpFormatter formatter = new HelpFormatter();
        formatter.printHelp("CaptureReplay", options);
    }

    public CliOptions(File inputFile, String outputDevice, float replaySpeed, boolean loop) {
        this.inputFile = inputFile;
        this.outputDevice = outputDevice;
        this.replaySpeed = replaySpeed;
        this.loop = loop;
    }

    public File getInputFile() {
        return inputFile;
    }

    public String getOutputDevice() {
        return outputDevice;
    }

    public float getReplaySpeed() {
        return replaySpeed;
    }

    public boolean isLoop() {
        return loop;
    }

}
