#
# Licensed to the Apache Software Foundation (ASF) under one
# or more contributor license agreements.  See the NOTICE file
# distributed with this work for additional information
# regarding copyright ownership.  The ASF licenses this file
# to you under the Apache License, Version 2.0 (the
# "License"); you may not use this file except in compliance
# with the License.  You may obtain a copy of the License at
#
#     https://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing,
# software distributed under the License is distributed on an
# "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
# KIND, either express or implied.  See the License for the
# specific language governing permissions and limitations
# under the License.
#

import logging
from asyncio import Transport
from dataclasses import dataclass, field
from typing import Type, Union

import plc4py
from plc4py.api.authentication.PlcAuthentication import PlcAuthentication
from plc4py.api.exceptions.exceptions import PlcFieldParseException
from plc4py.api.messages.PlcRequest import (
    PlcReadRequest,
    PlcRequest,
    ReadRequestBuilder,
)
from plc4py.api.messages.PlcResponse import (
    PlcReadResponse,
    PlcResponse,
)
from plc4py.api.PlcConnection import PlcConnection, PlcConnectionMetaData
from plc4py.api.PlcDriver import PlcDriver
from plc4py.api.value.PlcValue import PlcResponseCode
from plc4py.drivers.mock.MockTag import MockTagBuilder
from plc4py.drivers.PlcDriverLoader import PlcDriverLoader
from plc4py.spi.messages.PlcReader import DefaultPlcReader
from plc4py.spi.messages.PlcRequest import DefaultReadRequestBuilder
from plc4py.spi.messages.PlcWriter import DefaultPlcWriter
from plc4py.spi.messages.utils.ResponseItem import ResponseItem
from plc4py.spi.transport.MockTransport import MockTransport
from plc4py.spi.transport.Plc4xBaseTransport import Plc4xBaseTransport
from plc4py.spi.values.PlcValues import PlcBOOL, PlcINT


@dataclass
class MockDevice:
    async def read(
        self, request: PlcReadRequest, transport: Transport
    ) -> PlcReadResponse:
        """
        Reads one field from the Mock Device
        """
        response_items = {}
        for tag_name, tag in request.tags.items():
            logging.debug(f"Reading field {str(tag)} from Mock Device")

            if tag.data_type == "BOOL":
                response_items[tag_name] = ResponseItem(
                    PlcResponseCode.OK, PlcBOOL(False)
                )
            elif tag.data_type == "INT":
                response_items[tag_name] = ResponseItem(PlcResponseCode.OK, PlcINT(0))
            else:
                raise PlcFieldParseException
        return PlcReadResponse(PlcResponseCode.OK, response_items)


@dataclass
class MockConnection(PlcConnection, DefaultPlcReader, DefaultPlcWriter):
    _device: MockDevice = field(default_factory=lambda: MockDevice())
    _transport: Union[Plc4xBaseTransport, None] = None

    def _connect(self):
        """
        Connect the Mock PLC connection
        :return:
        """
        self._is_connected = True

    @staticmethod
    async def create(url):
        # config = PlcConfiguration(url)
        connection = MockConnection()
        connection._transport = await MockTransport.create(None, None, None)
        connection._connect()
        return connection

    def read_request_builder(self) -> ReadRequestBuilder:
        """
        :return: read request builder.
        """
        return DefaultReadRequestBuilder(MockTagBuilder)

    async def execute(self, request: PlcRequest) -> PlcResponse:
        """
        Executes a PlcRequest as long as it's already connected
        :param PlcRequest: Plc Request to execute
        :return: The response from the Plc/Device
        """
        if not self.is_connected():
            return self._default_failed_request(PlcResponseCode.NOT_CONNECTED)

        if isinstance(request, PlcReadRequest):
            return await self._read(request)

        return self._default_failed_request(PlcResponseCode.NOT_CONNECTED)


class MockDriver(PlcDriver):
    def __init__(self):
        self.protocol_code = "mock"
        self.protocol_name = "Mock"

    async def get_connection(
        self, url: str, authentication: PlcAuthentication = PlcAuthentication()
    ) -> PlcConnection:
        """
        Connects to a PLC using the given plc connection string.
        :param url: plc connection string
        :param authentication: authentication credentials.
        :return PlcConnection: PLC Connection object
        """
        return await MockConnection.create(url)


class MockDriverLoader(PlcDriverLoader):
    """
    Mock Connection Loader, after adding this to the setuptools entry point
    pluggy should be able to find this and import it.
    """

    @staticmethod
    @plc4py.drivers.hookimpl
    def get_driver() -> Type[MockDriver]:
        return MockDriver

    @staticmethod
    @plc4py.drivers.hookimpl
    def key() -> str:
        return "mock"
