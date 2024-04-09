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

import asyncio
import logging
from dataclasses import dataclass, field
from typing import Awaitable, Type, List, Dict

from plc4py.spi.messages.PlcWriter import PlcWriter

from plc4py.drivers.mock import MockTag
from plc4py.drivers.mock.MockTag import MockTagBuilder
from plc4py.spi.messages.PlcRequest import DefaultReadRequestBuilder

import plc4py

from plc4py.api.PlcConnection import PlcConnection, PlcConnectionMetaData
from plc4py.api.PlcDriver import PlcDriver
from plc4py.api.authentication.PlcAuthentication import PlcAuthentication
from plc4py.api.exceptions.exceptions import PlcFieldParseException
from plc4py.api.messages.PlcField import PlcTag
from plc4py.api.messages.PlcRequest import (
    ReadRequestBuilder,
    PlcReadRequest,
    PlcRequest,
    PlcWriteRequest,
)
from plc4py.api.messages.PlcResponse import (
    PlcReadResponse,
    PlcResponse,
    PlcWriteResponse,
)
from plc4py.api.value.PlcValue import PlcResponseCode, PlcValue
from plc4py.drivers.PlcDriverLoader import PlcDriverLoader
from plc4py.spi.messages.PlcReader import PlcReader
from plc4py.spi.messages.utils.ResponseItem import ResponseItem
from plc4py.spi.values.PlcValues import PlcBOOL
from plc4py.spi.values.PlcValues import PlcINT


@dataclass
class MockDevice:
    def read(self, tag: MockTag) -> ResponseItem[PlcValue]:
        """
        Reads one field from the Mock Device
        """
        logging.debug(f"Reading field {str(tag)} from Mock Device")

        if tag.data_type == "BOOL":
            return ResponseItem(PlcResponseCode.OK, PlcBOOL(False))
        elif tag.data_type == "INT":
            return ResponseItem(PlcResponseCode.OK, PlcINT(0))
        else:
            raise PlcFieldParseException


@dataclass
class MockConnection(PlcConnection, PlcReader, PlcWriter, PlcConnectionMetaData):
    _is_connected: bool = False
    device: MockDevice = field(default_factory=lambda: MockDevice())

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
        connection._connect()
        return connection

    def is_connected(self) -> bool:
        """
        Return the current status of the Mock PLC Connection
        :return bool: True is connected
        """
        return self._is_connected

    def close(self) -> None:
        """
        Closes the connection to the remote PLC.
        :return:
        """
        self._is_connected = False

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

    async def _read(self, request: PlcReadRequest) -> PlcReadResponse:
        """
        Executes a PlcReadRequest
        """
        if self.device is None:
            logging.error("No device is set in the mock connection!")
            return self._default_failed_request(PlcResponseCode.NOT_CONNECTED)

        try:
            logging.debug("Sending read request to Mock Device")
            response = PlcReadResponse(
                PlcResponseCode.OK,
                {
                    tag_name: self.device.read(tag)
                    for tag_name, tag in request.tags.items()
                },
            )
            return response
        except Exception as e:
            # TODO:- This exception is very general and probably should be replaced
            return PlcReadResponse(PlcResponseCode.INTERNAL_ERROR, {})

    async def _write(self, request: PlcWriteRequest) -> PlcWriteResponse:
        """
        Executes a PlcReadRequest
        """
        if self.device is None:
            logging.error("No device is set in the mock connection!")
            return self._default_failed_request(PlcResponseCode.NOT_CONNECTED)

        try:
            logging.debug("Sending read request to MockDevice")
            response = PlcWriteResponse(
                PlcResponseCode.OK,
                {
                    tag_name: self.device.write(tag)
                    for tag_name, tag in request.tags.items()
                },
            )
            return response
        except Exception as e:
            # TODO:- This exception is very general and probably should be replaced
            return PlcWriteResponse(PlcResponseCode.INTERNAL_ERROR, request.tags)

    def is_read_supported(self) -> bool:
        """
        Indicates if the connection supports read requests.
        :return: True if connection supports reading, False otherwise
        """
        return True

    def is_write_supported(self) -> bool:
        """
        Indicates if the connection supports write requests.
        :return: True if connection supports writing, False otherwise
        """
        return False

    def is_subscribe_supported(self) -> bool:
        """
        Indicates if the connection supports subscription requests.
        :return: True if connection supports subscriptions, False otherwise
        """
        return False

    def is_browse_supported(self) -> bool:
        """
        Indicates if the connection supports browsing requests.
        :return: True if connection supports browsing, False otherwise
        """
        return False


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
    @plc4py.hookimpl
    def get_driver() -> Type[MockDriver]:
        return MockDriver

    @staticmethod
    @plc4py.hookimpl
    def key() -> str:
        return "mock"
