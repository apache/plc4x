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
import importlib
import logging
import unittest
from abc import ABC, abstractmethod
from dataclasses import dataclass, field
from typing import Any, Iterator, List
from unittest import TestCase
from xml.etree import ElementTree

from bitarray.util import ba2hex
from xsdata.formats.dataclass.serializers import XmlSerializer
from xsdata.formats.dataclass.serializers.config import SerializerConfig
from xsdata.utils.text import snake_case

from plc4py.api.exceptions.exceptions import SerializationException
from plc4py.spi.generation.ReadBuffer import ReadBufferByteBased
from plc4py.spi.generation.WriteBuffer import WriteBufferXmlBased, WriteBufferByteBased

logger = logging.getLogger(__name__)


class ParserSerializerTestCase(TestCase):

    def __init__(self, methodName="runTest"):
        super().__init__(methodName)
        self.name = None
        self.test_case = None
        self.test_suite = None
        self.migrate = False

    def add_test_case(self, test_case):
        self.test_case = test_case
        self.name = test_case.name

    def add_test_suite(self, test_suite):
        self.test_suite = test_suite

    def add_migrate(self, migrate: bool):
        self.migrate = migrate

    def runTest(self):
        logger.info("Running testcases %s", self.test_case.name)
        test_case_raw = self.test_case.raw
        logger.debug(test_case_raw.hex())

        """ 
            Parse the raw bytes and generate to plc4py class 
        """
        read_buffer = ReadBufferByteBased(test_case_raw, self.test_suite.byte_order)
        try:
            protocol_name = self.test_suite.protocol_name
            output_flavor = self.test_suite.output_flavor.replace("-", "")
            root_type = self.test_case.root_type
            imported_module = importlib.import_module(
                f"plc4py.protocols.{protocol_name}.{output_flavor}.{root_type}",
                package="plc4py",
            )
        except ModuleNotFoundError as e:
            logging.exception(
                "Failed to find plc4py class for %s %s",
                self.test_suite.protocol_name,
                e,
            )
            raise e

        logger.debug("Parsing message")
        uninstantiated_class = getattr(imported_module, root_type)
        kwargs = {
            snake_case(element.qname): element.text
            for element in self.test_case.parser_arguments.local_element
        }
        modbus_adu = uninstantiated_class.static_parse_context(read_buffer, **kwargs)

        """ 
            Serialize the Model back to the raw bytes for comparison 
        """
        byte_buffer: WriteBufferByteBased = WriteBufferByteBased(
            size=len(self.test_case.raw), byte_order=self.test_suite.byte_order
        )
        modbus_adu.serialize(byte_buffer)

        comparison = self.test_case.raw == byte_buffer.get_bytes()
        if not comparison:

            """
            Serialize generate class to xml
            """
            exception: SerializationException = SerializationException()
            exception.add_note(
                "---------------------Failed Assertion------------------"
            )
            xml_buffer: WriteBufferXmlBased = WriteBufferXmlBased()
            modbus_adu.serialize(xml_buffer)
            result = xml_buffer.to_xml_string()

            exception.add_note("")
            exception.add_note(
                "---------------------Initial Raw Bytes-----------------"
            )
            exception.add_note(self.test_case.raw.hex())
            exception.add_note("")
            exception.add_note(
                "---------------------Serialized Raw Bytes from Model---"
            )
            exception.add_note(ba2hex(byte_buffer.bb))
            exception.add_note("")
            exception.add_note(
                "---------------------Initial XML Case------------------"
            )
            """
                Serialize the model back to xml for user comparison
            """
            factory = SerializerConfig(xml_declaration=False, pretty_print=True)
            serializer = XmlSerializer(config=factory)
            ss = serializer.render(self.test_case.xml)
            exception.add_note(ss)
            exception.add_note("")

            exception.add_note("---------------------Serialized XML case--------------")
            exception.add_note(result)
            raise exception
