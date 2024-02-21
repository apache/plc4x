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
from plc4py.api.messages.PlcResponse import (
    PlcResponse,
    PlcReadResponse,
    PlcWriteResponse,
    PlcTagResponse,
)
from plc4py.api.messages.PlcRequest import (
    ReadRequestBuilder,
    PlcRequest,
    PlcReadRequest,
    PlcWriteRequest,
)
from plc4py.api.value.PlcValue import PlcResponseCode
from plc4py.drivers.PlcDriverLoader import PlcDriverLoader
from plc4py.drivers.modbus.ModbusConfiguration import ModbusConfiguration
from plc4py.drivers.modbus.ModbusDevice import ModbusDevice
from plc4py.drivers.modbus.ModbusProtocol import ModbusProtocol
from plc4py.drivers.modbus.ModbusTag import ModbusTagBuilder
from plc4py.spi.messages.PlcReader import PlcReader
from plc4py.spi.messages.PlcRequest import DefaultReadRequestBuilder
from plc4py.spi.messages.PlcWriter import PlcWriter
from plc4py.spi.transport.Plc4xBaseTransport import Plc4xBaseTransport
from plc4py.spi.transport.TCPTransport import TCPTransport


class ModbusConnection(PlcConnection, PlcReader, PlcWriter, PlcConnectionMetaData):
    """
    Modbus TCP PLC connection implementation
    """

    def __init__(self, config: ModbusConfiguration, transport: Plc4xBaseTransport):
        super().__init__(config)
        self._configuration: ModbusConfiguration
        self._device: ModbusDevice = ModbusDevice(self._configuration)
        self._transport: Plc4xBaseTransport = transport

    @staticmethod
    async def create(url: str):
        """
        Static Factory to return an instance of a ModbusConnection.
        It creates the TCP connection to the Modbus device before returning.

        :param url: PLC4X connection string of the Modbus TCP connection
        :return ModbusConnection instance using the configuration from the url provided
        """
        config = ModbusConfiguration(url)
        loop = asyncio.get_running_loop()
        connection_future = loop.create_future()
        transport = await asyncio.wait_for(
            TCPTransport.create(
                protocol_factory=lambda: ModbusProtocol(connection_future),
                host=config.host,
                port=config.port,
            ),
            10,
        )
        return ModbusConnection(config, transport)

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
        return DefaultReadRequestBuilder(ModbusTagBuilder)

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

        return self._default_failed_request(PlcResponseCode.NOT_CONNECTED)

    def _read(self, request: PlcReadRequest) -> Awaitable[PlcReadResponse]:
        """
        Executes a PlcReadRequest
        """
        if self._device is None:
            logging.error("No device is set in the modbus connection!")
            return self._default_failed_read_request(PlcResponseCode.NOT_CONNECTED)

        async def _request(req, device) -> PlcReadResponse:
            try:
                response = await asyncio.wait_for(device.read(req, self._transport), 5)
                return response
            except Exception as e:
                # TODO:- This exception is very general and probably should be replaced
                return PlcReadResponse(PlcResponseCode.INTERNAL_ERROR, {})

        logging.debug("Sending read request to ModbusDevice")
        future = asyncio.ensure_future(_request(request, self._device))
        return future

    def _write(self, request: PlcWriteRequest) -> Awaitable[PlcTagResponse]:
        """
        Executes a PlcWriteRequest
        """
        if self._device is None:
            logging.error("No device is set in the modbus connection!")
            return self._default_failed_write_request(PlcResponseCode.NOT_CONNECTED)

        async def _request(req, device) -> PlcWriteResponse:
            try:
                response = await asyncio.wait_for(device.write(req, self._transport), 5)
                return response
            except Exception as e:
                # TODO:- This exception is very general and probably should be replaced
                return PlcWriteResponse(PlcResponseCode.INTERNAL_ERROR, {})

        logging.debug("Sending write request to ModbusDevice")
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
        return True

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


class ModbusDriver(PlcDriver):
    def __init__(self):
        super().__init__()
        self.protocol_code = "modbus-tcp"
        self.protocol_name = "Modbus TCP"

    async def get_connection(
        self, url: str, authentication: PlcAuthentication = PlcAuthentication()
    ) -> PlcConnection:
        """
        Connects to a PLC using the given plc connection string.
        :param url: plc connection string
        :param authentication: authentication credentials.
        :return PlcConnection: PLC Connection object
        """
        return await ModbusConnection.create(url)


class ModbusDriverLoader(PlcDriverLoader):
    """
    Modbus Driver Pluggy Hook Implmentation, lets pluggy find the driver by name
    """

    @staticmethod
    @plc4py.hookimpl
    def get_driver() -> Type[ModbusDriver]:
        return ModbusDriver

    @staticmethod
    @plc4py.hookimpl
    def key() -> str:
        return "modbus"
