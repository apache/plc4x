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
from asyncio import Transport
from dataclasses import dataclass, field
from typing import Awaitable, Type, List, Dict, Union

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

from plc4py.spi.transport.Plc4xBaseTransport import Plc4xBaseTransport


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
                response_items[tag_name] = ResponseItem(PlcResponseCode.OK, PlcBOOL(False))
            elif tag.data_type == "INT":
                response_items[tag_name] = ResponseItem(PlcResponseCode.OK, PlcINT(0))
            else:
                raise PlcFieldParseException
        return PlcReadResponse(PlcResponseCode.OK, response_items)


@dataclass
class MockConnection(PlcConnection, PlcReader, PlcWriter, PlcConnectionMetaData):
    _is_connected: bool = False
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

    def _check_connection(self) -> bool:
        """
        Checks if a ModbusDevice is set.

        :return: True if no device is set, False otherwise.
        """
        """
        A ModbusDevice is only set if the device was successfully connected during the constructor.
        If no device is set, it's not possible to execute any read or write requests.

        This method is used to prevent calling methods on the ModbusConnection which are not possible
        if no device is set.
        """
        return self._device is None

    async def _read(self, request: PlcReadRequest) -> PlcReadResponse:
        """
        Executes a PlcReadRequest

        This method sends a read request to the connected modbus device and waits for a response.
        The response is then returned as a PlcReadResponse.

        If no device is set, an error is logged and a PlcResponseCode.NOT_CONNECTED is returned.
        If an error occurs during the execution of the read request, a PlcResponseCode.INTERNAL_ERROR is
        returned.

        :param request: PlcReadRequest to execute
        :return: PlcReadResponse
        """
        if self._check_connection():
            logging.error("No device is set in the Mock connection!")
            return self._default_failed_request(PlcResponseCode.NOT_CONNECTED)

        # TODO: Insert Optimizer base on data from a browse request
        try:
            logging.debug("Sending read request to Mock Device")
            response = await asyncio.wait_for(
                self._device.read(request, self._transport), 10
            )
            return response
        except Exception:
            # TODO:- This exception is very general and probably should be replaced
            return PlcReadResponse(PlcResponseCode.INTERNAL_ERROR, {})

    async def _write(self, request: PlcWriteRequest) -> PlcWriteResponse:
        """
        Executes a PlcWriteRequest

        This method sends a write request to the connected Modbus device and waits for a response.
        The response is then returned as a PlcWriteResponse.

        If no device is set, an error is logged and a PlcWriteResponse with the
        PlcResponseCode.NOT_CONNECTED code is returned.
        If an error occurs during the execution of the write request, a
        PlcWriteResponse with the PlcResponseCode.INTERNAL_ERROR code is returned.

        :param request: PlcWriteRequest to execute
        :return: PlcWriteResponse
        """
        if self._check_connection():
            # If no device is set, log an error and return a response with the NOT_CONNECTED code
            logging.error("No device is set in the Mock connection!")
            return self._default_failed_request(PlcResponseCode.NOT_CONNECTED)

        try:
            # Send the write request to the device and wait for a response
            logging.debug("Sending write request to Mock Device")
            response = await asyncio.wait_for(
                self._device.write(request, self._transport), 5
            )
            # Return the response
            return response
        except Exception:
            # If an error occurs during the execution of the write request, return a response with
            # the INTERNAL_ERROR code. This exception is very general and probably should be replaced.
            # TODO:- This exception is very general and probably should be replaced
            return PlcWriteResponse(PlcResponseCode.INTERNAL_ERROR, {})

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
