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
from abc import ABC, abstractmethod
from dataclasses import dataclass, field
from typing import List, Iterator, Any, Tuple
from xml.etree import ElementTree

from xsdata.formats.dataclass.parsers import XmlParser

from .XmlTestSuiteLoader import (
    ParserSerializerTestCase,
)
from .generated.parser_serializer_testsuite import Testsuite

logger = logging.getLogger(__name__)


@dataclass
class XmlTestSuiteLoader(unittest.TestLoader, ABC):
    # Parent class for test suites utilizing XML as representation mechanism.
    # <p>
    # It adds handling of resource loading and helps in navigating test framework to suite source files.
    #

    test_suite_document: str
    test_suite_document_xml: Iterator[Tuple[str, Any]] = field(init=False)

    def __post_init__(self) -> None:
        # ElementTree.register_namespace('test', 'https://plc4x.apache.org/schemas/driver-testsuite.xsd')
        self.test_suite_document_xml = ElementTree.iterparse(
            self.test_suite_document, events=(["start"])
        )

    @abstractmethod
    def test_suite_tests(self) -> List[unittest.TestCase]:
        # Get a list of test cases
        pass


@dataclass
class ParserSerializerTestsuiteRunner(XmlTestSuiteLoader):
    test_suite_document: str
    auto_migrate: bool = (
        False  # if set to true if will automigrate and on the next run test should be green
    )
    ignored_test_cases: List[str] = field(default_factory=list)

    @property
    def test_suite_tests(self) -> List[unittest.TestCase]:
        parser = XmlParser()
        test_suite_xml = parser.parse(self.test_suite_document, Testsuite)

        dynamic_tests: List[unittest.TestCase] = []
        for test_case in test_suite_xml.testcase:
            if not test_case.name in self.ignored_test_cases:
                test_suite = ParserSerializerTestCase()
                test_suite.add_test_case(test_case)
                test_suite.add_test_suite(test_suite_xml)
                test_suite.add_migrate(self.auto_migrate)
                dynamic_tests.append(test_suite)
        logger.info("Found %s testcases.", len(test_suite_xml.testcase))
        return dynamic_tests
