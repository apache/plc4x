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
package org.apache.plc4x.test.generator

import groovy.xml.MarkupBuilder
import groovyjarjarpicocli.CommandLine
import org.apache.commons.lang3.SystemUtils
import org.apache.plc4x.java.spi.generation.ByteOrder
import org.apache.plc4x.java.spi.generation.ReadBufferByteBased
import org.apache.plc4x.java.spi.generation.WriteBufferXmlBased
import org.pcap4j.core.PcapHandle
import org.pcap4j.core.Pcaps
import org.pcap4j.packet.EthernetPacket
import org.pcap4j.packet.Packet
import org.pcap4j.packet.TcpPacket
import org.pcap4j.packet.UdpPacket

import java.util.function.Consumer

@CommandLine.Command(name = "pstg", version = "0.1-Alpha", mixinStandardHelpOptions = true)
class ParserSerializerTestsuiteGenerator implements Runnable {

    @CommandLine.Option(names = ["-d", "--debug"], description = "show debug information")
    boolean debug

    @CommandLine.Option(names = ["-l", "--little-endian"], description = "render a little endian attribute")
    boolean littleEndian

    @CommandLine.Option(names = ["-t", "--test-suite-name"], description = "render a little endian attribute", defaultValue = "TODO: name me")
    String testSuiteName

    @CommandLine.Option(names = ["-p", "--protocol-name"], description = "output flavor of the driver", defaultValue = "TODO: name me")
    String protocolName

    @CommandLine.Option(names = ["-o", "--output-flavor"], description = "output flavor of the driver", defaultValue = "read-write")
    String flavor

    @CommandLine.Parameters(paramLabel = "root-message-type-class", description = "fqcn of the root message type")
    String rootMessageTypeClass

    @CommandLine.Parameters(paramLabel = "inputpcap", description = "pcap to consume")
    String pcapFile

    @CommandLine.Parameters(paramLabel = "outputtestfile", description = "test suite file to write")
    String xmlTestSuiteFile

    static Consumer<Integer> exitFunc = System::exit

    static void main(String... args) {
        if (SystemUtils.IS_OS_MAC) {
            // On my Intel Mac I found the libs in: "/usr/local/Cellar/libpcap/1.10.1/lib"
            // On my M1 Mac I found the libs in: "/opt/homebrew/Cellar/libpcap/1.10.1/lib"
            if (new File("/usr/local/Cellar/libpcap/1.10.1/lib").exists()) {
                System.getProperties().setProperty("jna.library.path", "/usr/local/Cellar/libpcap/1.10.1/lib");
            } else if (new File("/opt/homebrew/opt/libpcap/lib").exists()) {
                System.getProperties().setProperty("jna.library.path", "/opt/homebrew/opt/libpcap/lib");
            }
        }

        int exitCode = new CommandLine(new ParserSerializerTestsuiteGenerator()).execute(args)
        exitFunc.accept(exitCode)
    }

    @Override
    void run() {
        generateOutput()
    }

    void generateOutput() {
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
        generateXmlTestSuite xmlTestSuiteFile, testSuiteName, littleEndian, rootMessageTypeClass, testMap

        infoOutput "Done"
    }

    def generateXmlTestSuite(String xmlTestSuiteFile, String testSuiteName, boolean littleEndian, String rootMessageTypeClass, Map<String, byte[]> testMap) {
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

      https://www.apache.org/licenses/LICENSE-2.0

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
            protocolName protocolName
            outputFlavor flavor
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
                                def xmlWriter = new WriteBufferXmlBased()
                                message.serialize(xmlWriter)
                                def messageString = xmlWriter.getXmlString()
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

    List<byte[]> readPayloads(String pcapFile) {
        def values = []
        def pcapHandle = Pcaps.openOffline(pcapFile, PcapHandle.TimestampPrecision.NANO)
        Packet packet
        try {
            while ((packet = pcapHandle.nextPacket) != null) {
                def udpPacket = packet.get(UdpPacket.class)
                if (udpPacket != null) {
                    values << udpPacket.payload.rawData
                } else {
                    def tcpPacket = packet.get(TcpPacket.class)
                    if (tcpPacket != null) {
                        values << tcpPacket.payload.rawData
                    } else {
                        def ethernetPacket = packet.get(EthernetPacket.class)
                        if (ethernetPacket != null) {
                            values << ethernetPacket.rawData
                        } else {
                            values << new byte[]{0, 0, 0, 0}
                        }
                    }
                }
            }
        } catch (Exception e) {
            e.printStackTrace()
        } finally {
            pcapHandle.close()
        }
        values
    }

    def checkForRequirements() {
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

    def errorOutput(String format, def ... values) {
        System.err.printf("$format\n", values)
    }

    def infoOutput(def format, def ... values) {
        System.out.printf("$format\n", values)
    }

    def debugOutput(def format, def ... values) {
        if (!debug) return
        System.out.printf("$format\n", values)
    }

}
