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

import static org.xmlunit.matchers.CompareMatcher.isIdenticalTo
import static spock.util.matcher.HamcrestSupport.*
import spock.lang.Specification


import java.nio.file.Files

class ParserSerializerTestsuiteGeneratorSpec extends Specification {
    def "Test main with an example pcap"() {
        given:
        def testSuitePath = Files.createTempFile("parser-serializer-testsuite", ".xml")
        URL pcap = ParserSerializerTestsuiteGeneratorSpec.getResource("/bacnet-stack-services.cap");
        File pcapFile = new File(pcap.toURI());
        ParserSerializerTestsuiteGenerator.exitFunc = (it)-> println("exiting with $it")

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
