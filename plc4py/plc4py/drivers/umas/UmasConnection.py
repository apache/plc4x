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
from typing import Awaitable, Type

import plc4py
from plc4py.api.authentication.PlcAuthentication import PlcAuthentication
from plc4py.api.exceptions.exceptions import PlcConnectionException
from plc4py.api.messages.PlcRequest import (
    BrowseRequestBuilder,
    PlcBrowseRequest,
    PlcReadRequest,
    PlcRequest,
    ReadRequestBuilder,
    WriteRequestBuilder,
    PlcWriteRequest,
)
from plc4py.api.messages.PlcResponse import (
    PlcResponse,
)
from plc4py.api.PlcConnection import PlcConnection, PlcConnectionMetaData
from plc4py.api.PlcDriver import PlcDriver
from plc4py.api.value.PlcValue import PlcResponseCode
from plc4py.drivers.PlcDriverLoader import PlcDriverLoader
from plc4py.drivers.umas.UmasConfiguration import UmasConfiguration
from plc4py.drivers.umas.UmasDevice import UmasDevice
from plc4py.drivers.umas.UmasProtocol import UmasProtocol
from plc4py.drivers.umas.UmasTag import UmasTagBuilder
from plc4py.spi.messages.PlcBrowser import DefaultPlcBrowser
from plc4py.spi.messages.PlcReader import DefaultPlcReader
from plc4py.spi.messages.PlcRequest import (
    DefaultBrowseRequestBuilder,
    DefaultReadRequestBuilder,
    DefaultWriteRequestBuilder,
)
from plc4py.spi.transport.Plc4xBaseTransport import Plc4xBaseTransport
from plc4py.spi.transport.TCPTransport import TCPTransport
from plc4py.spi.messages.PlcWriter import DefaultPlcWriter


class UmasConnection(
    PlcConnection, DefaultPlcReader, DefaultPlcWriter, DefaultPlcBrowser
):
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
                f"Time out while opening UMAS connection to device - {config.host}",
                e,
            )
        return connection

    def read_request_builder(self) -> ReadRequestBuilder:
        """
        :return: read request builder.
        """
        return DefaultReadRequestBuilder(UmasTagBuilder)

    def write_request_builder(self) -> WriteRequestBuilder:
        """
        :return: write request builder.
        """
        return DefaultWriteRequestBuilder(UmasTagBuilder)

    def browse_request_builder(self) -> BrowseRequestBuilder:
        """
        :return: browse request builder.
        """
        return DefaultBrowseRequestBuilder()

    async def execute(self, request: PlcRequest) -> PlcResponse:
        """
        Executes a PlcRequest as long as it's already connected
        :param request: Plc Request to execute
        :return: The response from the Plc/Device
        """
        if not self.is_connected():
            return self._default_failed_request(PlcResponseCode.NOT_CONNECTED)

        if isinstance(request, PlcReadRequest):
            return await self._read(request)

        if isinstance(request, PlcWriteRequest):
            return await self._write(request)

        if isinstance(request, PlcBrowseRequest):
            return await self._browse(request)

        return self._default_failed_request(PlcResponseCode.NOT_CONNECTED)


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
    @plc4py.drivers.hookimpl
    def get_driver() -> Type[UmasDriver]:
        return UmasDriver

    @staticmethod
    @plc4py.drivers.hookimpl
    def key() -> str:
        return "umas"
