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
package org.apache.plc4x.java.s7.utils;

import org.apache.commons.io.FileExistsException;
import org.apache.commons.io.FileUtils;
import org.w3c.dom.Document;
import org.w3c.dom.NodeList;
import org.xml.sax.SAXException;

import javax.xml.parsers.DocumentBuilder;
import javax.xml.parsers.DocumentBuilderFactory;
import javax.xml.parsers.ParserConfigurationException;
import javax.xml.xpath.*;
import java.io.File;
import java.io.IOException;
import java.util.Arrays;

/**
 * Little helper class used to generate pcapng files for usage in Wireshark.
 * As I can't use the WireShark code to define the constant values for the s7 protocol as
 * this code is licensed under GPL, with this tool, I'll simply generate loads of pcapng files
 * each with slightly different values and hereby find out the constant values.
 */
public class PcapGenerator {

    private static final int CONNECT_TPDU_CODE = 99;
    private static final int CONNECT_TPDU_CLASS = 104;
    private static final int CONNECT_TPDU_PARAMETER_CODE = 105;

    private static final byte[] CONNECT_TEMPLATE = {
        // PCAP header
        // Global Header:
        //  magic number
        (byte) 0xD4, (byte) 0xC3, (byte) 0xB2, (byte) 0xA1,
        //  Version(major / minor)
        (byte) 0x02, (byte) 0x00, (byte) 0x04, (byte) 0x00,
        //  Timezone
        (byte) 0x00, (byte) 0x00, (byte) 0x00, (byte) 0x00,
        // 0
        (byte) 0x00, (byte) 0x00, (byte) 0x00, (byte) 0x00,
        //  Snapshot length
        (byte) 0xFF, (byte) 0xFF, (byte) 0x00, (byte) 0x00,
        //  Network
        (byte) 0x01, (byte) 0x00, (byte) 0x00, (byte) 0x00,
        // Packet header
        //  Timestamp (seconds)
        (byte) 0x8F, (byte) 0x46, (byte) 0x4E, (byte) 0x53,
        //  Timestamp (microseconds)
        (byte) 0x1B, (byte) 0x0C, (byte) 0x0C, (byte) 0x00,
        // Packet length (in file)
        (byte) 0x44, (byte) 0x00, (byte) 0x00, (byte) 0x00,
        // Packet length (real)
        (byte) 0x44, (byte) 0x00, (byte) 0x00, (byte) 0x00,
        ///////////////////////////////////////////////////
        // Packet content
        // Ethernet packet
        (byte) 0x00, (byte) 0x1B, (byte) 0x1B, (byte) 0x1A,
        (byte) 0xD6, (byte) 0xE0, (byte) 0xB8, (byte) 0x70,
        (byte) 0xF4, (byte) 0x6D, (byte) 0x2F, (byte) 0x58,
        (byte) 0x08, (byte) 0x00,
        // IP packet
        (byte) 0x45, (byte) 0x00, (byte) 0x00, (byte) 0x36,
        (byte) 0x6B, (byte) 0x6E, (byte) 0x40, (byte) 0x00,
        (byte) 0x40, (byte) 0x06, (byte) 0x44, (byte) 0xDF,
        (byte) 0x00, (byte) 0x00, (byte) 0x00, (byte) 0x00,
        (byte) 0x00, (byte) 0x00, (byte) 0x00, (byte) 0x00,
        // TCP packet
        (byte) 0xCC, (byte) 0xDE, (byte) 0x00, (byte) 0x66,
        (byte) 0x8B, (byte) 0x9F, (byte) 0x0B, (byte) 0x76,
        (byte) 0x00, (byte) 0x1A, (byte) 0x4D, (byte) 0x96,
        (byte) 0x50, (byte) 0x18, (byte) 0x72, (byte) 0x10,
        (byte) 0x1D, (byte) 0x43, (byte) 0x00, (byte) 0x00,
        // TPKT packet
        (byte) 0x03, (byte) 0x00, (byte) 0x00, (byte) 0x0E, // 68
        // 98
        (byte) 0x09,
        (byte) 0xE0 /* TPDU Code */,
        (byte) 0x00, (byte) 0x00,
        (byte) 0x00, (byte) 0x00,
        (byte) 0x00 /* TPDU Class */,
        (byte) 0xFF /* TPDU Parameter Code */,
        (byte) 0x01, (byte) 0x0a
    };

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
        (byte) 0x00, (byte) 0x00, (byte) 0x00, (byte) 0x00,
        (byte) 0x00, (byte) 0x00, (byte) 0x00, (byte) 0x00,
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
        //  Timezone
        (byte) 0x00, (byte) 0x00, (byte) 0x00, (byte) 0x00,
        (byte) 0x00, (byte) 0x0E, (byte) 0x00, (byte) 0x00,

        (byte) 0x04 /* Function Code [111] */, (byte) 0x01,
        (byte) 0x12 /* Specification Type [113] */, (byte) 0x0A,
        (byte) 0x10 /* Variable Addressing Mode [115] */,
        (byte) 0x02 /* Variable Type [116] */, (byte) 0x00, (byte) 0x40,
        (byte) 0x00, (byte) 0x01, (byte) 0x84 /* Memory Area [121] */, (byte) 0x00,
        (byte) 0x00, (byte) 0x00
    };

    private static final int READ_VARIABLE_RESPONSE_TRANSPORT_SIZE = 116;

    private static final byte[] READ_VARIABLE_RESPONSE_TEMPLATE = {
        // PCAP header
        // Global Header:
        //  magic number
        (byte) 0xD4, (byte) 0xC3, (byte) 0xB2, (byte) 0xA1,
        //  Version(major / minor)
        (byte) 0x02, (byte) 0x00, (byte) 0x04, (byte) 0x00,
        //  Timezone
        (byte) 0x00, (byte) 0x00, (byte) 0x00, (byte) 0x00,
        // 0
        (byte) 0x00, (byte) 0x00, (byte) 0x00, (byte) 0x00,
        //  Snapshot length
        (byte) 0xFF, (byte) 0xFF, (byte) 0x00, (byte) 0x00,
        //  Network
        (byte) 0x01, (byte) 0x00, (byte) 0x00, (byte) 0x00,
        // Packet header
        //  Timestamp (seconds)
        (byte) 0x00, (byte) 0x00, (byte) 0x00, (byte) 0x00,
        //  Timestamp (microseconds)
        (byte) 0x00, (byte) 0x00, (byte) 0x00, (byte) 0x00,
        // Packet length (in file)
        (byte) 0xAB, (byte) 0x00, (byte) 0x00, (byte) 0x00,
        // Packet length (real)
        (byte) 0xAB, (byte) 0x00, (byte) 0x00, (byte) 0x00,
        // Packet content
        (byte) 0x00, (byte) 0x1b, (byte) 0x1b, (byte) 0x1a,
        (byte) 0xd6, (byte) 0xe0, (byte) 0xb8, (byte) 0x70,
        (byte) 0xf4, (byte) 0x6d, (byte) 0x2f, (byte) 0x58,
        (byte) 0x08, (byte) 0x00, (byte) 0x45, (byte) 0x00,
        (byte) 0x00, (byte) 0x9d, (byte) 0x6b, (byte) 0x6e,
        (byte) 0x40, (byte) 0x00, (byte) 0x40, (byte) 0x06,
        (byte) 0x44, (byte) 0xdf, (byte) 0x00, (byte) 0x00,
        (byte) 0x00, (byte) 0x00, (byte) 0x00, (byte) 0x00,
        (byte) 0x00, (byte) 0x00, (byte) 0xcc, (byte) 0xde,
        (byte) 0x00, (byte) 0x66, (byte) 0x8b, (byte) 0x9f,
        (byte) 0x0b, (byte) 0x76, (byte) 0x00, (byte) 0x1a,
        (byte) 0x4d, (byte) 0x96, (byte) 0x50, (byte) 0x18,
        (byte) 0x72, (byte) 0x10, (byte) 0x1d, (byte) 0x43,
        (byte) 0x00, (byte) 0x00, (byte) 0x03, (byte) 0x00,
        (byte) 0x00, (byte) 0x75, (byte) 0x02, (byte) 0xf0,
        (byte) 0x80,
        // S7 Packet
        (byte) 0x32, (byte) 0x03, (byte) 0x00,
        (byte) 0x00, (byte) 0x08, (byte) 0x00, (byte) 0x00,
        (byte) 0x02, (byte) 0x00, (byte) 0x60, (byte) 0x00,
        (byte) 0x00,
        // Read Var Parameter
        (byte) 0x04, (byte) 0x01,
        // Data
        (byte) 0xff,
        (byte) 0x07, (byte) 0x00, (byte) 0x5c, (byte) 0x03,
        (byte) 0x00, (byte) 0x00, (byte) 0x00, (byte) 0x00,
        (byte) 0x00, (byte) 0x00, (byte) 0x00, (byte) 0x00,
        (byte) 0x00, (byte) 0x00, (byte) 0x00, (byte) 0x00,
        (byte) 0x00, (byte) 0x00, (byte) 0x00, (byte) 0x00,
        (byte) 0x00, (byte) 0x00, (byte) 0x00, (byte) 0x00,
        (byte) 0x00, (byte) 0x00, (byte) 0x00, (byte) 0x00,
        (byte) 0x00, (byte) 0x00, (byte) 0x00, (byte) 0x00,
        (byte) 0x00, (byte) 0x00, (byte) 0x00, (byte) 0x00,
        (byte) 0x00, (byte) 0x00, (byte) 0x00, (byte) 0x00,
        (byte) 0x00, (byte) 0x00, (byte) 0x00, (byte) 0x00,
        (byte) 0x00, (byte) 0x00, (byte) 0x00, (byte) 0x00,
        (byte) 0x00, (byte) 0x00, (byte) 0x00, (byte) 0x00,
        (byte) 0x00, (byte) 0x00, (byte) 0x00, (byte) 0x00,
        (byte) 0x00, (byte) 0x00, (byte) 0x00, (byte) 0x00,
        (byte) 0x00, (byte) 0x00, (byte) 0x00, (byte) 0x00,
        (byte) 0x00, (byte) 0x00, (byte) 0x00, (byte) 0x00,
        (byte) 0x00, (byte) 0x00, (byte) 0xaa, (byte) 0x00,
        (byte) 0x00, (byte) 0x00, (byte) 0x00, (byte) 0x00,
        (byte) 0x00, (byte) 0x00, (byte) 0x00, (byte) 0x00,
        (byte) 0x00, (byte) 0x00, (byte) 0x00, (byte) 0x00,
        (byte) 0x00, (byte) 0x00, (byte) 0x00, (byte) 0x00,
        (byte) 0x00, (byte) 0x00, (byte) 0x00, (byte) 0x00,
        (byte) 0x00, (byte) 0x00, (byte) 0x00
    };

    public static void main(String[] args) throws Exception {
        File outputDir = new File("target/out");
        System.out.println("Outputting to: " + outputDir.getAbsolutePath());

        //
        generateFiles(outputDir, "read-response/transport-size", READ_VARIABLE_RESPONSE_TRANSPORT_SIZE, READ_VARIABLE_RESPONSE_TEMPLATE,
            "//field[@name='s7comm.data.transportsize' and not(contains(@showname, 'Unknown'))]");

        // Connect variable generation
        // FIXME: Not all codes are output probably due to errors in the rest package structure
        generateFiles(outputDir, "connect/tpdu-code", CONNECT_TPDU_CODE, CONNECT_TEMPLATE,
            "//field[@name='cotp.type' and (substring(@value, 2, 1) = '0')]");
        generateFiles(outputDir, "connect/tpdu-class", CONNECT_TPDU_CLASS, CONNECT_TEMPLATE,
            "//field[@name='cotp.class' and (substring(@unmaskedvalue, 2, 1) = '0')]");
        generateFiles(outputDir, "connect/tpdu-parameter-code", CONNECT_TPDU_PARAMETER_CODE, CONNECT_TEMPLATE,
            "//field[@name='cotp.parameter_code' and not(contains(@showname, 'Unknown'))]");

        // Read variable generation
        generateFiles(outputDir, "read/message-type", READ_VARIABLE_MESSAGE_TYPE_BYTE, READ_VARIABLE_TEMPLATE,
            "//field[@name='s7comm.header.rosctr' and not(contains(@showname, 'Unknown'))]");
        generateFiles(outputDir, "read/function-code", READ_VARIABLE_FUNCTION_CODE_BYTE, READ_VARIABLE_TEMPLATE,
            "//field[@name='s7comm.param.func' and not(contains(@showname, 'Unknown'))]");
        generateFiles(outputDir, "read/specification-type", READ_VARIABLE_SPECIFICATION_TYPE_BYTE, READ_VARIABLE_TEMPLATE,
            "//field[@name='s7comm.param.item.varspec' and not(contains(../@showname, 'Unknown'))]");
        generateFiles(outputDir, "read/syntax-id", READ_VARIABLE_SYNTAX_ID_BYTE, READ_VARIABLE_TEMPLATE,
            "//field[@name='s7comm.param.item.syntaxid' and not(contains(@showname, 'Unknown'))]");
        generateFiles(outputDir, "read/transport-size", READ_VARIABLE_TRANSPORT_SIZE_BYTE, READ_VARIABLE_TEMPLATE,
            "//field[@name='s7comm.param.item.transp_size' and not(contains(@showname, 'Unknown'))]");
        generateFiles(outputDir, "read/memory-area", READ_VARIABLE_MEMORY_AREA, READ_VARIABLE_TEMPLATE,
            "//field[@name='s7comm.param.item.area' and not(contains(@showname, 'Unknown'))]");

        // Read Ack variable generation

        // Write variable generation

        // Write Ack variable generation
    }

    private static void generateFiles(File parentDir, String testName, int byteIndex, byte[] template, String matchXPath)
        throws IOException, ParserConfigurationException, XPathExpressionException, InterruptedException, SAXException {
        File testDir = new File(parentDir, testName);
        if (!testDir.exists() && !testDir.mkdirs()) {
            throw new FileExistsException("Error creating directory " + testDir.getAbsolutePath());
        }

        // Initialize the heavy weight XML stuff.
        DocumentBuilderFactory domFactory = DocumentBuilderFactory.newInstance();
        domFactory.setNamespaceAware(true);
        DocumentBuilder builder = domFactory.newDocumentBuilder();
        XPathFactory factory = XPathFactory.newInstance();
        XPath xpath = factory.newXPath();
        XPathExpression expr = xpath.compile(matchXPath);

        // Iterate over all possible value the current byte could have.
        for (int i = 0; i <= 255; i++) {

            // Generate a packet.
            byte[] copy = Arrays.copyOf(template, template.length);
            copy[byteIndex] = (byte) i;
            File output = new File(testDir, i + ".pcapng");
            FileUtils.writeByteArrayToFile(output, copy);

            // Use WireShark to decode the packet to an xml form.
            File decodedOutput = new File(testDir, i + ".xml");
            ProcessBuilder pb = new ProcessBuilder("tshark", "-T", "pdml", "-r", output.getAbsolutePath());
            pb.redirectOutput(ProcessBuilder.Redirect.to(decodedOutput));
            Process process = pb.start();
            process.waitFor();

            // Check if a given XPath is found -> A valid parameter was found.
            Document document = builder.parse(decodedOutput);
            NodeList list = (NodeList) expr.evaluate(document, XPathConstants.NODESET);

            // If a valid parameter is found, output that to the console.
            if (list.getLength() > 0) {
                String name = list.item(0).getAttributes().getNamedItem("showname").getNodeValue();
                String value = list.item(0).getAttributes().getNamedItem("value").getNodeValue();
                System.out.println("found option for " + testName + ": " + name + " = 0x" + value);
            }
        }
    }

}
