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
package org.apache.plc4x.java.profinet.utils;

import org.apache.commons.io.FileUtils;

import java.io.File;
import java.io.IOException;
import java.util.Arrays;

/**
 * Little helper class used to generate pcapng files for usage in Wireshark.
 * As I can't use the WireShark code to define the constant values for the profinet protocol as
 * this code is licensed under GPL, with this tool, I'll simply generate loads of pcapng files
 * each with slightly different values and hereby find out the constant values.
 */
public class PcapGenerator {

    private static final int READ_VARIABLE_MESSAGE_TYPE_BYTE = 102;
    private static final int READ_VARIABLE_FUNCTION_CODE_BYTE = 111;
    private static final int READ_VARIABLE_SPECIFICATION_TYPE_BYTE = 113;
    private static final int READ_VARIABLE_SYNTAX_ID_BYTE = 115;
    private static final int READ_VARIABLE_TRANSPORT_SIZE_BYTE = 116;
    private static final int READ_VARIABLE_MEMORY_AREA = 121;

    private static final byte[] READ_VARIABLE_TEMPLATE = {
        // PCAP header
        (byte) 0xD4, (byte) 0xC3, (byte) 0xB2, (byte) 0xA1,
        (byte) 0x02, (byte) 0x00, (byte) 0x04, (byte) 0x00,
        (byte) 0x00, (byte) 0x00, (byte) 0x00, (byte) 0x00,
        (byte) 0x00, (byte) 0x00, (byte) 0x00, (byte) 0x00,
        (byte) 0xFF, (byte) 0xFF, (byte) 0x00, (byte) 0x00,
        (byte) 0x01, (byte) 0x00, (byte) 0x00, (byte) 0x00,
        (byte) 0x8F, (byte) 0x46, (byte) 0x4E, (byte) 0x53,
        (byte) 0x1B, (byte) 0x0C, (byte) 0x0C, (byte) 0x00,
        (byte) 0x55, (byte) 0x00, (byte) 0x00, (byte) 0x00,
        (byte) 0x55, (byte) 0x00, (byte) 0x00, (byte) 0x00,
        // Packet content
        // Ethernet packet
        (byte) 0x00, (byte) 0x1B, (byte) 0x1B, (byte) 0x1A,
        (byte) 0xD6, (byte) 0xE0, (byte) 0xB8, (byte) 0x70,
        (byte) 0xF4, (byte) 0x6D, (byte) 0x2F, (byte) 0x58,
        (byte) 0x08, (byte) 0x00,
        // IP packet
        (byte) 0x45, (byte) 0x00, (byte) 0x00, (byte) 0x47,
        (byte) 0x6B, (byte) 0x6E, (byte) 0x40, (byte) 0x00,
        (byte) 0x40, (byte) 0x06, (byte) 0x44, (byte) 0xDF,
        (byte) 0x86, (byte) 0xF9, (byte) 0x3E, (byte) 0xCE,
        (byte) 0x86, (byte) 0xF9, (byte) 0x3D, (byte) 0xA3,
        // TCP packet
        (byte) 0xCC, (byte) 0xDE, (byte) 0x00, (byte) 0x66,
        (byte) 0x8B, (byte) 0x9F, (byte) 0x0B, (byte) 0x76,
        (byte) 0x00, (byte) 0x1A, (byte) 0x4D, (byte) 0x96,
        (byte) 0x50, (byte) 0x18, (byte) 0x72, (byte) 0x10,
        (byte) 0x1D, (byte) 0x43, (byte) 0x00, (byte) 0x00,
        // TPKT packet
        (byte) 0x03, (byte) 0x00, (byte) 0x00, (byte) 0x1F,
        //
        (byte) 0x02, (byte) 0xF0, (byte) 0x80,
        // S7 packet
        // Job header
        (byte) 0x32, (byte) 0x01 /* Message Type [102] */,
        (byte) 0x00, (byte) 0x00, (byte) 0x00, (byte) 0x00,
        (byte) 0x00, (byte) 0x0E, (byte) 0x00, (byte) 0x00,

        (byte) 0x04 /* Function Code [111] */, (byte) 0x01,
        (byte) 0x12 /* Specification Type [113] */, (byte) 0x0A,
        (byte) 0x10 /* Variable Addressing Mode [115] */,
        (byte) 0x02 /* Variable Type [116] */, (byte) 0x00, (byte) 0x40,
        (byte) 0x00, (byte) 0x01, (byte) 0x84 /* Memory Area [121] */, (byte) 0x00,
        (byte) 0x00, (byte) 0x00
    };

    public static void main(String[] args) throws IOException {
        File outputDir = new File("target/out");
        System.out.println("Outputting to: " + outputDir.getAbsolutePath());

        // Read variable generation
        generateFiles(outputDir, "read/message-type", READ_VARIABLE_MESSAGE_TYPE_BYTE, READ_VARIABLE_TEMPLATE);
        generateFiles(outputDir, "read/function-code", READ_VARIABLE_FUNCTION_CODE_BYTE, READ_VARIABLE_TEMPLATE);
        generateFiles(outputDir, "read/specification-type", READ_VARIABLE_SPECIFICATION_TYPE_BYTE, READ_VARIABLE_TEMPLATE);
        generateFiles(outputDir, "read/syntax-id", READ_VARIABLE_SYNTAX_ID_BYTE, READ_VARIABLE_TEMPLATE);
        generateFiles(outputDir, "read/transport-size", READ_VARIABLE_TRANSPORT_SIZE_BYTE, READ_VARIABLE_TEMPLATE);
        generateFiles(outputDir, "read/memory-area", READ_VARIABLE_MEMORY_AREA, READ_VARIABLE_TEMPLATE);

        // Read Ack variable generation

        // Write variable generation

        // Write Ack variable generation
    }

    private static void generateFiles(File parentDir, String testName, int byteIndex, byte[] template) throws IOException {
        File testDir = new File(parentDir, testName);
        if(!testDir.mkdirs()) {
            throw new RuntimeException("Error creating directory " + testDir.getAbsolutePath());
        }
        for(int i = 0; i <= 255; i++) {
            byte[] copy = Arrays.copyOf(template, template.length);
            copy[byteIndex] = (byte) i;
            File output = new File(testDir, i + ".pcapng");
            FileUtils.writeByteArrayToFile(output, copy);
        }
    }

}
