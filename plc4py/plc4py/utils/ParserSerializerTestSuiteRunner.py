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
from utils.generated.driver_testsuite import DriverTestsuite
from xsdata.formats.dataclass.parsers import XmlParser

from api.exceptions.exceptions import ParseException
from utils.XmlTestSuiteLoader import (
    ParserSerializerTestSuite,
    XmlTestSuiteLoader,
    ParserSerializerTestCase,
)
from utils.generated.parser_serializer_testsuite import Testsuite

logger = logging.getLogger(__name__)


@dataclass
class ParserSerializerTestsuiteRunner(XmlTestSuiteLoader):
    test_suite_document: str
    auto_migrate: bool = (
        False  # if set to true if will automigrate and on the next run test should be green
    )
    ignored_test_cases: List[str] = field(default_factory=list)

    @property
    def test_suite_tests(self) -> List[unittest.case]:
        parser = XmlParser()
        test_suite_xml = parser.parse(self.test_suite_document, Testsuite)

        dynamic_tests: List[unittest.case] = []
        for test_case in test_suite_xml.testcase:
            if not test_case.name in self.ignored_test_cases:
                test_suite = ParserSerializerTestCase()
                test_suite.add_test_case(test_case)
                test_suite.add_test_suite(test_suite_xml)
                dynamic_tests.append(test_suite)
        logger.info("Found %s testcases.", len(test_suite_xml.testcase))
        return dynamic_tests
