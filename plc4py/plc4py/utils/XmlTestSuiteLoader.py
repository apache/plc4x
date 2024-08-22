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
import importlib
import logging
import unittest
from abc import ABC, abstractmethod
from dataclasses import dataclass, field
from typing import Any, Iterator, List
from unittest import TestCase
from xml.etree import ElementTree
from xml.etree.ElementTree import XMLParser

from xsdata.utils.text import camel_case, snake_case

from spi.generation.ReadBuffer import ReadBufferByteBased

logger = logging.getLogger(__name__)


class ParserSerializerTestSuite:
    # TODO:- Do something with this
    pass


class ParserSerializerTestCase(TestCase):

    def __init__(self, methodName="runTest"):
        super().__init__(methodName)
        self.name = None
        self.test_case = None
        self.test_suite = None

    def add_test_case(self, test_case):
        self.test_case = test_case
        self.name = test_case.name

    def add_test_suite(self, test_suite):
        self.test_suite = test_suite

    def runTest(self):
        logger.info("Running testcases %s", self.test_case.name)
        test_case_raw = self.test_case.raw
        if logger.isEnabledFor(logging.DEBUG):
            # TODO:- Ascii Box Writer please
            logger.debug(test_case_raw)
        read_buffer = ReadBufferByteBased(test_case_raw, self.test_suite.byte_order)
        try:
            protocol_name = self.test_suite.protocol_name
            output_flavor = self.test_suite.output_flavor.replace("-", "")
            root_type = self.test_case.root_type
            imported_module = importlib.import_module(
                f"protocols.{protocol_name}.{output_flavor}.{root_type}",
                package="plc4py",
            )
        except ModuleNotFoundError as e:
            logging.exception(e)
            raise e

        logger.debug("Parsing message")
        uninstantiated_class = getattr(imported_module, root_type)
        kwargs = {
            snake_case(element.qname): element.text
            for element in self.test_case.parser_arguments.local_element
        }
        modbus_adu = uninstantiated_class.static_parse_context(read_buffer, **kwargs)
        pass
        # Message parsedOutput = (Message) uninstantiated_class.parse(readBuffer, testcase.getParserArguments().toArray());
        #             LOGGER.trace("Validating and migrating");


#
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


@dataclass
class XmlTestSuiteLoader(unittest.TestLoader, ABC):
    # Parent class for test suites utilizing XML as representation mechanism.
    # <p>
    # It adds handling of resource loading and helps in navigating test framework to suite source files.
    #

    test_suite_document: str
    test_suite_document_xml: Iterator[tuple[str, Any]] = field(init=False)

    def __post_init__(self) -> None:
        # ElementTree.register_namespace('test', 'https://plc4x.apache.org/schemas/driver-testsuite.xsd')
        self.test_suite_document_xml = ElementTree.iterparse(
            self.test_suite_document, events=(["start"])
        )

    @abstractmethod
    def test_suite_tests(self) -> List[unittest.case]:
        # Get a list of test cases
        pass
