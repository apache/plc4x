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
from abc import abstractmethod, ABC
from dataclasses import dataclass, field
from typing import Any, Iterator
from xml.etree import ElementTree


class ParserSerializerTestSuite:
    # TODO:- Do something with this
    pass


class TestCase:
    # TODO:- Do something with this
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
        self.test_suite_document_xml = ElementTree.iterparse(self.test_suite_document)

    @abstractmethod
    def run(self, test_suite: ParserSerializerTestSuite, test_case: TestCase) -> Iterator[unittest.TestResult]:
        # Actually run the test suite
        pass
