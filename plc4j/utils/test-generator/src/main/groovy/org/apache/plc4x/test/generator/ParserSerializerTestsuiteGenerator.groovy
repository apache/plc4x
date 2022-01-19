/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package org.apache.plc4x.test.generator

import groovy.cli.picocli.CliBuilder
import groovy.xml.MarkupBuilder
import org.apache.plc4x.java.spi.generation.ByteOrder
import org.apache.plc4x.java.spi.generation.ReadBufferByteBased
import org.pcap4j.core.PcapHandle
import org.pcap4j.core.Pcaps
import org.pcap4j.packet.Packet
import org.pcap4j.packet.TcpPacket
import org.pcap4j.packet.UdpPacket

class ParserSerializerTestsuiteGenerator {

    static boolean debug = false

    static void main(String... args) {
        def options = parseAndGetOptions(args)
        def pcapFile = options.arguments()[0]
        String testSuiteName = options.t
        boolean littleEnding = options.l
        String rootMessageTypeClass = options.c
        String xmlTestSuiteFile = options.arguments()[1]

        generateOutput(pcapFile, xmlTestSuiteFile, testSuiteName, rootMessageTypeClass, littleEnding)
    }

    static void generateOutput(String pcapFile, String xmlTestSuiteFile, String testSuiteName, String rootMessageTypeClass, boolean littleEnding) {
        checkForRequirements()
        // 1. Extract headers out of the existing pcap file
        debugOutput "Extracting headers"
        def infoFields = "tshark -r $pcapFile -T fields -e _ws.col.No. -e _ws.col.Info".execute().text.split("\n")
        debugOutput "Found $infoFields.length messages"

        debugOutput "Extracting payload"
        def values = readPayloads(pcapFile)
        assert infoFields.length == values.size()
        Map<String, byte[]> testMap = [infoFields, values].transpose().collectEntries { [it[0], it[1]] }

        debugOutput "Generating xml to $xmlTestSuiteFile"
        generateXmlTestSuite xmlTestSuiteFile, testSuiteName, littleEnding, rootMessageTypeClass, testMap

        infoOutput "Done"
    }

    static def generateXmlTestSuite(String xmlTestSuiteFile, String testSuiteName, boolean littleEndian, String rootMessageTypeClass, Map<String, byte[]> testMap) {
        def writer = new FileWriter(xmlTestSuiteFile)
        def xml = new MarkupBuilder(writer)

        def byteOrder = "BIG_ENDIAN"
        if (littleEndian)
            byteOrder = "LITTLE_ENDIAN"
        xml.mkp.xmlDeclaration(version: "1.0", encoding: "utf-8")
        xml.mkp.comment '''
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
'''
        xml.mkp.yield "\n"
        xml."test:testsuite"(["xmlns:test": "https://plc4x.apache.org/schemas/parser-serializer-testsuite.xsd", "byteOrder": byteOrder]) {
            name(testSuiteName)
            mkp.yield "\n"
            // TODO: get name
            protocolName "bacnet"
            // TODO: get flavor
            outputFlavor "read-write"
            mkp.yield "\n"

            testMap.each { testEntry ->
                testcase {
                    name testEntry.key
                    raw testEntry.value.encodeHex()
                    "root-type" Class.forName(rootMessageTypeClass).simpleName
                    delegate.xml {
                        mkp.yield "\n"
                        mkp.yieldUnescaped(() -> {
                            try {
                                def clazz = Class.forName(rootMessageTypeClass)
                                def message = clazz."staticParse"(new ReadBufferByteBased(testEntry.value, littleEndian ? ByteOrder.LITTLE_ENDIAN : ByteOrder.BIG_ENDIAN))
                                def messageString = message.toString()
                                // TODO ugyWorkaround
                                messageString = messageString.split("\n").collect { "      $it" }.join("\n") + "\n    "
                                return messageString
                            } catch (e) {
                                return e.toString()
                            }
                        }())
                    }
                }
                mkp.yield "\n"
            }
        }
    }

    static List<byte[]> readPayloads(String pcapFile) {
        def values = []
        def pcapHandle = Pcaps.openOffline(pcapFile, PcapHandle.TimestampPrecision.NANO)
        Packet packet
        while ((packet = pcapHandle.nextPacket) != null) {
            def udpPacket = packet.get(UdpPacket.class)
            if (udpPacket != null) {
                values << udpPacket.payload.rawData
            } else {
                def tcpPacket = packet.get(TcpPacket.class)
                if (tcpPacket != null) {
                    values << tcpPacket.payload.rawData
                } else {
                    values << new byte[]{0, 0, 0, 0}
                }
            }
        }
        pcapHandle.close()
        values
    }

    static def parseAndGetOptions(String... args) {
        def cli = new CliBuilder(usage: 'ptsg [options] <inputpcap> <outputtestfile>')
        cli.d(longOpt: 'debug', 'Debug output')
        cli.l(longOpt: 'littleEndian', 'Little Endian')
        cli.t(longOpt: 'testSuiteName', args: 1, argName: 'testSuiteName', defaultValue: 'TODO: name me', 'test suite name')
        cli.c(longOpt: 'rootMessageTypeClass', args: 1, argName: 'rootMessageTypeClass', required: true, 'fqcn of the root message type')
        def options = cli.parse(args)
        assert options
        if (options.d) debug = true

        assert options.arguments().size() == 2
        options
    }

    static def checkForRequirements() {
        // check for sh
        try {
            debugOutput "sh -c 'echo we have a shell'".execute().text
        } catch (IOException e) {
            errorOutput "$ParserSerializerTestsuiteGenerator.class needs a shell to work properly"
            throw e
        }
        // check for thark
        try {
            debugOutput "tshark --version".execute().text
        } catch (IOException e) {
            errorOutput "$ParserSerializerTestsuiteGenerator.class needs tshark to work properly"
            throw e
        }
    }

    static def errorOutput(String format, def ... values) {
        System.err.printf("$format\n", values)
    }

    static def infoOutput(def format, def ... values) {
        System.out.printf("$format\n", values)
    }

    static def debugOutput(def format, def ... values) {
        if (!debug) return
        System.out.printf("$format\n", values)
    }
}
