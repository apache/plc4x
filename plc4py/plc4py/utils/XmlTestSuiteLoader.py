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
import unittest
from abc import ABC, abstractmethod
from dataclasses import dataclass, field
from typing import Any, Iterator
from xml.etree import ElementTree
from xml.etree.ElementTree import XMLParser


class ParserSerializerTestSuite:
    # TODO:- Do something with this
    pass


class TestCase:
    # TODO:- Do something with this
    pass


@dataclass
class TestCaseBuilder:
    element: ElementTree

    def build(self) -> TestCase:
        name = None
        description = None
        raw = None
        root_type = None
        parser_arguments = None
        for element in self.element:
            if element.tag == "name":
                name = element.text
            elif element.tag == "description":
                description = element.text
            elif element.tag == "raw":
                raw = element.text
            elif element.tag == "root-type":
                root_type = element.text
            elif element.tag == "parser-arguments":
                parser_arguments = element.text
            elif element.tag == "parser-raw":
                parser_raw = element.text

        return TestCase()

    #                 Element descriptionElement = testcaseXml.element(new QName("description"));
    #                 Element rawElement = testcaseXml.element(new QName("raw"));
    #                 Element rootTypeElement = testcaseXml.element(new QName("root-type"));
    #                 Element parserArgumentsElement = testcaseXml.element(new QName("parser-arguments"));
    #                 Element xmlElement = testcaseXml.element(new QName("xml"));
    #
    #                 String name = nameElement.getTextTrim();
    #                 String description = (descriptionElement != null) ? descriptionElement.getTextTrim() : null;
    #                 byte[] raw = Hex.decodeHex(rawElement.getTextTrim());
    #                 String rootType = rootTypeElement.getTextTrim();
    #
    #                 // Parse additional parser arguments.
    #                 List<String> parserArguments = new LinkedList<>();
    #                 if (parserArgumentsElement != null) {
    #                     for (Element element : parserArgumentsElement.elements()) {
    #                         parserArguments.add(element.getTextTrim());
    #                     }
    #                 }
    #                 Testcase testcase = new Testcase(testsuiteName, protocolName, outputFlavor, name, description, raw, rootType, parserArguments, xmlElement);
    #                 if (testcaseXml instanceof LocationAwareElement) {
    #                     // pass source location to test
    #                     testcase.setLocation(((LocationAwareElement) testcaseXml).getLocation());
    #                 }
    #                 testcases.add(testcase);
    #             }


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
    def run(
        self, test_suite: ParserSerializerTestSuite, test_case: TestCase
    ) -> Iterator[unittest.TestResult]:
        # Actually run the test suite
        pass
