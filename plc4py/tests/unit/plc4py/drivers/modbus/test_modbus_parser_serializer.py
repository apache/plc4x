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
import os

import pytest

from ...utils.ParserSerializerTestSuiteRunner import ParserSerializerTestsuiteRunner


def pytest_generate_tests(metafunc):
    path: str = os.path.join(
        os.path.dirname(__file__),
        "..",
        "..",
        "..",
        "..",
        "..",
        "..",
        "protocols",
        "modbus",
        "src",
        "test",
        "resources",
        "protocols",
        "modbus",
        "tcp",
        "ParserSerializerTestsuite.xml",
    )
    xml_loader = ParserSerializerTestsuiteRunner(path)

    test_suites = xml_loader.test_suite_tests

    metafunc.parametrize(
        "test_case", test_suites, ids=[test_case.name for test_case in test_suites]
    )


@pytest.mark.asyncio
async def test_parse_serializer_parse_xml_test(test_case) -> None:
    test_case.runTest()
