#  Licensed to the Apache Software Foundation (ASF) under one
#  or more contributor license agreements.  See the NOTICE file
#  distributed with this work for additional information
#  regarding copyright ownership.  The ASF licenses this file
#  to you under the Apache License, Version 2.0 (the
#  "License"); you may not use this file except in compliance
#  with the License.  You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
#  Unless required by applicable law or agreed to in writing,
#  software distributed under the License is distributed on an
#  "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
#  KIND, either express or implied.  See the License for the
#  specific language governing permissions and limitations
#  under the License.
import logging
import unittest
from dataclasses import dataclass, field
from typing import Iterator, List
from driver_testsuite import DriverTestsuite
from xsdata.formats.dataclass.parsers import XmlParser

from api.exceptions.exceptions import ParseException
from utils.XmlTestSuiteLoader import (
    ParserSerializerTestSuite,
    TestCase,
    TestCaseBuilder,
    XmlTestSuiteLoader,
)

logger = logging.getLogger(__name__)


@dataclass
class ParserSerializerTestsuiteRunner(XmlTestSuiteLoader):
    test_suite_document: str
    auto_migrate: bool = False  # if set to true if will automigrate and on the next run test should be green
    ignored_test_cases: List[str] = field(default_factory=list)

    @property
    def test_suite_tests(self) -> Iterator[unittest.case]:
        parser = XmlParser()
        test_suite = parser.parse(self.test_suite_document, DriverTestsuite)

        test_suite = self._parse_test_suite()
        dynamic_tests: List[unittest.case] = []
        for test_case in test_suite:
            if not test_case.__name__ in self.ignored_test_cases:
                dynamic_tests.append(test_case)

        return iter(dynamic_tests)

    def _parse_test_suite(self):
        root = self.test_suite_document_xml.__next__()
        byte_order = root[1].attrib.get("byteOrder", "LITTLE_ENDIAN")
        test_suite_name = None
        protocol_name = None
        output_flavor = None
        driver_name = None
        found_header = False
        test_cases: List[TestCase] = []

        for name, element in self.test_suite_document_xml:
            if element.tag == "name":
                test_suite_name = element.text
            elif element.tag == "protocolName":
                protocol_name = element.text
            elif element.tag == "outputFlavor":
                output_flavor = element.text
            elif element.tag == "driver-name":
                driver_name = element.text
            elif element.tag == "testcase":
                test_cases.append(TestCaseBuilder(element).build())
            found_header = (
                test_suite_name is not None
                and protocol_name is not None
                and output_flavor is not None
                and driver_name is not None
            )

        if not found_header:
            raise ParseException()
        pass

    def run(
        self, test_suite: ParserSerializerTestSuite, test_case: TestCase
    ) -> Iterator[unittest.TestResult]:
        return self.test_suite_document_xml is not None


#
#
#             List<Element> testcasesXml = testsuiteXml.elements(new QName("testcase"));
#             List<Testcase> testcases = new ArrayList<>(testcasesXml.size());
#             for (Element testcaseXml : testcasesXml) {
#
#             LOGGER.info(String.format("Found %d testcases.", testcases.size()));
#             return new ParserSerializerTestsuite(testsuiteName, testcases, byteOrder, options);
#         } catch (DocumentException e) {
#             throw new ParserSerializerTestsuiteException("Error parsing testsuite xml", e);
#         } catch (DecoderException e) {
#             throw new ParserSerializerTestsuiteException("Error parsing testcase raw data", e);
#         }
#     }
#
#     private void run(ParserSerializerTestsuite testSuite, Testcase testcase) throws ParserSerializerTestsuiteException {
#         LOGGER.info("Running testcase {}", testcase);
#         byte[] testcaseRaw = testcase.getRaw();
#         if (LOGGER.isDebugEnabled())
#             LOGGER.debug("\n{}", AsciiBoxWriter.DEFAULT.boxString("Testcase raw input", org.apache.plc4x.java.spi.utils.hex.Hex.dump(testcaseRaw), 0));
#         ReadBufferByteBased readBuffer = new ReadBufferByteBased(testcaseRaw, testSuite.getByteOrder());
#
#         try {
#             MessageInput<?> messageInput = MessageResolver.getMessageIOStaticLinked(
#                 testSuite.getOptions(),
#                 testcase.getRootType()
#             );
#
#             ////////////////////////////////////////////////////////////////////////////////////////////////////////////
#             // Parse the raw bytes into a message
#             ////////////////////////////////////////////////////////////////////////////////////////////////////////////
#
#             LOGGER.trace("Parsing message");
#             Message parsedOutput = (Message) messageInput.parse(readBuffer, testcase.getParserArguments().toArray());
#             LOGGER.trace("Validating and migrating");
#
#             ////////////////////////////////////////////////////////////////////////////////////////////////////////////
#             // Compare the parsed message with the reference XML
#             ////////////////////////////////////////////////////////////////////////////////////////////////////////////
#
#             // In this case no reference xml has been provided
#             // (This is usually during development)
#             if (testcase.getXml().elements().size() == 0) {
#                 WriteBufferXmlBased writeBufferXmlBased = new WriteBufferXmlBased();
#                 parsedOutput.serialize(writeBufferXmlBased);
#                 String xmlString = writeBufferXmlBased.getXmlString();
#                 throw new ParserSerializerTestsuiteException("Missing reference xml element. Parsed: \n" + xmlString);
#             }
#             // If more than one root element is provided, the testcase is corrupt.
#             else if (testcase.getXml().elements().size() > 1) {
#                 throw new ParserSerializerTestsuiteException("Too many element roots in testcase");
#             }
#             boolean migrated = MessageValidatorAndMigrator.validateOutboundMessageAndMigrate(
#                 testcase.getName(),
#                 messageInput,
#                 testcase.getXml().elements().get(0),
#                 testcase.getParserArguments(),
#                 testcaseRaw,
#                 testSuite.getByteOrder(),
#                 autoMigrate,
#                 suiteUri
#             );
#             if (migrated) {
#                 LOGGER.warn("Migrated testcase {}", testcase);
#             }
#             LOGGER.debug("Parsed message {}", parsedOutput);
#             LOGGER.info("Parsing passed for testcase {}", testcase);
#
#             ////////////////////////////////////////////////////////////////////////////////////////////////////////////
#             // Serialize the parsed message to a byte array
#             ////////////////////////////////////////////////////////////////////////////////////////////////////////////
#
#             LOGGER.trace("Writing message back again");
#             WriteBufferByteBased writeBuffer = new WriteBufferByteBased(parsedOutput.getLengthInBytes(), testSuite.getByteOrder());
#             parsedOutput.serialize(writeBuffer);
#             LOGGER.info("Serializing passed for testcase {}", testcase);
#             byte[] data = writeBuffer.getBytes();
#             if (testcaseRaw.length != data.length) {
#                 LOGGER.info("Expected a byte array with a length of {} but got one with {}", testcaseRaw.length, data.length);
#             }
#
#             ////////////////////////////////////////////////////////////////////////////////////////////////////////////
#             // Compare the serialized bytes to the initial raw array
#             ////////////////////////////////////////////////////////////////////////////////////////////////////////////
#
#             if (!Arrays.equals(testcaseRaw, data)) {
#                 // This goes to std out on purpose to preserve coloring
#                 System.out.println(HexDiff.diffHex(testcaseRaw, data));
#                 throw new ParserSerializerTestsuiteException("Differences were found after serializing.\nExpected: " +
#                     Hex.encodeHexString(testcaseRaw) + "\nBut Got:  " + Hex.encodeHexString(data) + "");
#             }
#         } catch (SerializationException | ParseException e) {
#             throw new ParserSerializerTestsuiteException("Unable to parse message", e);
#         }
#     }
#
# }
