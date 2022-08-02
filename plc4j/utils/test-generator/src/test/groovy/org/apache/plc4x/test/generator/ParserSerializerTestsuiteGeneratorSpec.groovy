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

import org.opentest4j.TestAbortedException
import spock.lang.IgnoreIf

import static org.xmlunit.matchers.CompareMatcher.isIdenticalTo
import static spock.util.matcher.HamcrestSupport.*
import spock.lang.Specification


import java.nio.file.Files

// TODO: find out how to access surefire groups...
@IgnoreIf({ env["ENABLE_ALL_TESTS"] == null || env["ENABLE_ALL_TESTS"] == "false" })
class ParserSerializerTestsuiteGeneratorSpec extends Specification {
    def "Test main with an example pcap"() {
        given:
        try {
            "sh -c 'echo 1'".execute()
        } catch (e) {
            throw new TestAbortedException("no sh", e)
        }
        try {
            "tshark --version".execute()
        } catch (e) {
            throw new TestAbortedException("no tshark", e)
        }
        // On macs the libpcap version installed is usually 1.9.x
        // This causes errors. We therefore need to skip this test on such devices.
        try {
            def sout = new StringBuilder(), serr = new StringBuilder()
            def proc = "tcpdump --version".execute()
            proc.consumeProcessOutput(sout, serr)
            proc.waitForOrKill(1000)
            def output = serr + sout
            if (!output.contains("libpcap")) {
                throw new TestAbortedException("no libpcap")
            }
            def libpcapVersion = output.substring(output.indexOf("libpcap version ") + 16)
            libpcapVersion = libpcapVersion.substring(0, libpcapVersion.indexOf("\n"))
            // Check the libpcapVersion is at least 1.10.0
            if (!libpcapVersion.startsWith("1.10")) {
                throw new TestAbortedException("minimum libpcap version 1.10.0 expected")
            }
        } catch (e) {
            throw new TestAbortedException("no tcpdump", e)
        }
        if (!new File('/bin/sh').canExecute()) throw new TestAbortedException("No bin sh")
        def testSuitePath = Files.createTempFile("parser-serializer-testsuite", ".xml")
        URL pcap = ParserSerializerTestsuiteGeneratorSpec.getResource("/bacnet-stack-services.cap");
        File pcapFile = new File(pcap.toURI());
        ParserSerializerTestsuiteGenerator.exitFunc = (it) -> println("exiting with $it")

        when:
        ParserSerializerTestsuiteGenerator.main("-d", "-t TODO: name me", "-l", DummyMessageRootType.class.name, pcapFile.path, testSuitePath.toString())

        then:
        assert Files.exists(testSuitePath)
        expect:
        def expected = ParserSerializerTestsuiteGeneratorSpec.getResource("/ParserSerializerTestSuite.xml").text
        def actual = testSuitePath.toFile().text
        that actual, isIdenticalTo(expected).ignoreComments().ignoreWhitespace()
    }
}
