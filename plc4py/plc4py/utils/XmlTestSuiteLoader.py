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
import dataclasses
import importlib
import logging
import unittest
from abc import ABC, abstractmethod
from dataclasses import dataclass, field
from typing import Any, Iterator, List
from unittest import TestCase
from xml.etree import ElementTree
from xml.etree.ElementTree import XMLParser

from bitarray.util import ba2hex
from xsdata.formats.dataclass.models.generics import AnyElement
from xsdata.formats.dataclass.serializers import XmlSerializer
from xsdata.formats.dataclass.serializers.config import SerializerConfig
from xsdata.utils.text import camel_case, snake_case

from spi.generation.ReadBuffer import ReadBufferByteBased
from spi.generation.WriteBuffer import WriteBufferXmlBased, WriteBufferByteBased

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
        # Dummy code to generate test case
        modbus_adu = uninstantiated_class.static_parse_context(read_buffer, **kwargs)
        xml_buffer: WriteBufferXmlBased = WriteBufferXmlBased()
        modbus_adu.serialize(xml_buffer)
        # Could be used to generate the test case
        result = xml_buffer.to_xml_string()

        byte_buffer: WriteBufferByteBased = WriteBufferByteBased(size=len(self.test_case.raw), byte_order=self.test_suite.byte_order)
        modbus_adu.serialize(byte_buffer)

        comparision = self.test_case.raw == byte_buffer.get_bytes()
        if not comparision:
            print("---------------------Failed Assertion-----------------")
            # Ascii Box Again, might be time :(
            print(self.test_case.raw.hex())
            print(ba2hex(byte_buffer.bb))
            print(result)
            factory = SerializerConfig(xml_declaration=False, pretty_print=True)
            serializer = XmlSerializer(config=factory)
            ss = serializer.render(self.test_case.xml)
            print(ss)
            assert False

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
    def test_suite_tests(self) -> List[unittest.case]:
        # Get a list of test cases
        pass
