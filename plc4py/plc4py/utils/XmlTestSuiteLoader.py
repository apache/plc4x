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
from unittest import TestCase
from xml.etree import ElementTree
from xml.etree.ElementTree import XMLParser


class ParserSerializerTestSuite:
    # TODO:- Do something with this
    pass


class ParserSerializerTestCase(TestCase):

    def __init__(self, methodName='runTest'):
        super().__init__(methodName)
        self.name = None
        self.test_case = None

    def add_test_case(self, test_case):
        self.test_case = test_case
        self.name = test_case.name

    def runTest(self):
        pass


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
