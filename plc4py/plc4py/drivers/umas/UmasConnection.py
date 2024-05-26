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
from typing import Type, Awaitable

import plc4py
from plc4py.api.PlcConnection import PlcConnection, PlcConnectionMetaData
from plc4py.api.PlcDriver import PlcDriver
from plc4py.api.authentication.PlcAuthentication import PlcAuthentication
from plc4py.api.exceptions.exceptions import PlcConnectionException
from plc4py.api.messages.PlcResponse import (
    PlcResponse,
    PlcReadResponse,
    PlcBrowseResponse,
    PlcTagResponse,
    PlcWriteResponse,
)
from plc4py.api.messages.PlcRequest import (
    ReadRequestBuilder,
    PlcRequest,
    PlcReadRequest,
    BrowseRequestBuilder,
    PlcBrowseRequest,
    PlcWriteRequest,
)
from plc4py.api.value.PlcValue import PlcResponseCode
from plc4py.drivers.PlcDriverLoader import PlcDriverLoader
from plc4py.drivers.umas.UmasConfiguration import UmasConfiguration
from plc4py.drivers.umas.UmasDevice import UmasDevice
from plc4py.drivers.umas.UmasProtocol import UmasProtocol
from plc4py.drivers.umas.UmasTag import UmasTagBuilder
from plc4py.spi.messages.PlcReader import PlcReader
from plc4py.spi.messages.PlcRequest import (
    DefaultReadRequestBuilder,
    DefaultBrowseRequestBuilder,
)
from plc4py.spi.messages.PlcWriter import PlcWriter
from plc4py.spi.transport.Plc4xBaseTransport import Plc4xBaseTransport
from plc4py.spi.transport.TCPTransport import TCPTransport


class UmasConnection(PlcConnection, PlcReader, PlcWriter, PlcConnectionMetaData):
    """
    Umas TCP PLC connection implementation
    """

    DEFAULT_TCP_CONNECTION_TIMEOUT = 10
    DEFAULT_UMAS_CONNECTION_TIMEOUT = 60  # Long timeout to account for te 15s delay when first initiating a connection. Possibly just the simulator though.
    log = logging.getLogger(__name__)

    def __init__(self, config: UmasConfiguration, transport: Plc4xBaseTransport):
        super().__init__(config)
        self._configuration: UmasConfiguration
        self._device: UmasDevice = UmasDevice(self._configuration)
        self._transport: Plc4xBaseTransport = transport

    @staticmethod
    async def create(url: str):
        """
        Static Factory to return an instance of a UmasConnection.
        It creates the TCP connection to the Umas device before returning.

        :param url: PLC4X connection string of the Umas TCP connection
        :return UmasConnection: instance using the configuration from the url provided
        """
        config = UmasConfiguration(url)
        loop = asyncio.get_running_loop()
        connection_future = loop.create_future()
        try:
            logging.debug(f"Opening TCP connection for UMAS device {config.host}")
            transport = await asyncio.wait_for(
                TCPTransport.create(
                    protocol_factory=lambda: UmasProtocol(
                        connection_future, False, config.unit_identifier
                    ),
                    host=config.host,
                    port=config.port,
                ),
                UmasConnection.DEFAULT_TCP_CONNECTION_TIMEOUT,
            )
            connection = UmasConnection(config, transport)
        except TimeoutError as e:
            raise PlcConnectionException(
                f"Timeout while opening TCP connection to umas device - {config.host}",
                e,
            )

        try:
            await asyncio.wait_for(
                connection._device.connect(transport),
                timeout=UmasConnection.DEFAULT_UMAS_CONNECTION_TIMEOUT,
            )
        except TimeoutError as e:
            raise PlcConnectionException(
                f"Time out while opening UMAS connection to device - {config.host}", e
            )
        return connection

    def is_connected(self) -> bool:
        """
        Indicates if the connection is established to a remote PLC.
        :return: True if connection, False otherwise
        """
        if self._transport is not None:
            return not self._transport.is_closing()
        else:
            return False

    def close(self) -> None:
        """
        Closes the connection to the remote PLC.
        :return:
        """
        if self._transport is not None:
            self._transport.close()

    def read_request_builder(self) -> ReadRequestBuilder:
        """
        :return: read request builder.
        """
        return DefaultReadRequestBuilder(UmasTagBuilder)

    def browse_request_builder(self) -> BrowseRequestBuilder:
        """
        :return: browse request builder.
        """
        return DefaultBrowseRequestBuilder()

    def execute(self, request: PlcRequest) -> Awaitable[PlcResponse]:
        """
        Executes a PlcRequest as long as it's already connected
        :param request: Plc Request to execute
        :return: The response from the Plc/Device
        """
        if not self.is_connected():
            return self._default_failed_request(PlcResponseCode.NOT_CONNECTED)

        if isinstance(request, PlcReadRequest):
            return self._read(request)

        if isinstance(request, PlcBrowseRequest):
            return self._browse(request)

        return self._default_failed_request(PlcResponseCode.NOT_CONNECTED)

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
            logging.error("No device is set in the Umas connection!")
            return self._default_failed_request(PlcResponseCode.NOT_CONNECTED)

        # TODO: Insert Optimizer base on data from a browse request
        try:
            logging.debug("Sending read request to Umas Device")
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
            logging.error("No device is set in the Umas connection!")
            return self._default_failed_request(PlcResponseCode.NOT_CONNECTED)

        try:
            # Send the write request to the device and wait for a response
            logging.debug("Sending write request to Umas Device")
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

    def _browse(self, request: PlcBrowseRequest) -> Awaitable[PlcBrowseResponse]:
        """
        Executes a PlcBrowseRequest
        """
        if self._device is None:
            logging.error("No device is set in the Umas connection!")
            return self._default_failed_request(PlcResponseCode.NOT_CONNECTED)

        async def _request(req, device) -> PlcBrowseResponse:
            try:
                response = await asyncio.wait_for(
                    device.browse(req, self._transport), 10
                )
                return response
            except Exception:
                # TODO:- This exception is very general and probably should be replaced
                self.log.exception(
                    "Caught an exception while executing a browse request"
                )
                return PlcBrowseResponse(PlcResponseCode.INTERNAL_ERROR, {})

        logging.debug("Sending browse request to UmasDevice")
        future = asyncio.ensure_future(_request(request, self._device))
        return future

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
        return True


class UmasDriver(PlcDriver):
    def __init__(self):
        super().__init__()
        self.protocol_code = "umas"
        self.protocol_name = "umas"

    async def get_connection(
        self, url: str, authentication: PlcAuthentication = PlcAuthentication()
    ) -> PlcConnection:
        """
        Connects to a PLC using the given plc connection string.
        :param url: plc connection string
        :param authentication: authentication credentials.
        :return PlcConnection: PLC Connection object
        """
        return await UmasConnection.create(url)


class UmasDriverLoader(PlcDriverLoader):
    """
    Umas Driver Pluggy Hook Implementation, lets pluggy find the driver by name
    """

    @staticmethod
    @plc4py.hookimpl
    def get_driver() -> Type[UmasDriver]:
        return UmasDriver

    @staticmethod
    @plc4py.hookimpl
    def key() -> str:
        return "umas"
