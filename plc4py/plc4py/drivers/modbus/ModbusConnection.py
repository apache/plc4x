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

    """
    Initialization of ModbusConnection object

    :param config: Modbus configuration object
    :param transport: Plc4xBaseTransport object used for the TCP connection
    """
    def __init__(self, config: ModbusConfiguration, transport: Plc4xBaseTransport):
        """
        Initializes a ModbusConnection object

        :param config: Modbus configuration object
        :param transport: Plc4xBaseTransport object used for the TCP connection
        """
        super().__init__(config)
        self._configuration: ModbusConfiguration = config
        """
        ModbusDevice object that encapsulates the PLC device configuration
        and provides methods for device related operations
        """
        self._device: ModbusDevice = ModbusDevice(self._configuration)
        """
        Plc4xBaseTransport object used for the TCP connection
        """
        self._transport: Plc4xBaseTransport = transport

    @staticmethod
    async def create(url: str) -> "ModbusConnection":
        """
        Static Factory to return an instance of a ModbusConnection.
        It creates the TCP connection to the Modbus device before returning.

        :param url: PLC4X connection string of the Modbus TCP connection
        :return ModbusConnection instance using the configuration from the url provided
        """
        # Extract configuration from the given url
        config = ModbusConfiguration(url)
        # Create a future to manage the creation of the TCP connection
        loop = asyncio.get_running_loop()
        connection_future = loop.create_future()
        # Create the TCP connection to the Modbus device.
        # The creation process is asynchronous and takes up to 10 seconds.
        # If the creation fails, the Future is cancelled.
        try:
            transport = await asyncio.wait_for(
                TCPTransport.create(
                    protocol_factory=lambda: ModbusProtocol(connection_future),
                    host=config.host,
                    port=config.port,
                ),
                10,
            )
        except asyncio.TimeoutError:
            # If the creation times out, cancel the Future
            connection_future.cancel()
            raise ValueError(
                f"Failed to connect to {url} in 10 seconds. "
                "Please check the PLC device is running and the connection string is correct."
            )
        # If the creation succeeds, return a ModbusConnection object
        # using the configuration from the url provided and the transport
        return ModbusConnection(config, transport)

    def is_connected(self) -> bool:
        """
        Indicates if the connection is established to a remote PLC.

        :return: True if connection, False otherwise
        """
        """
        The function checks if the connection to the remote PLC is established.
        The connection is considered established if the transport object is not None
        and the transport is not closing.

        :return bool: True if connection, False otherwise
        """
        if self._transport is not None:
            # The transport is not None if a connection is established
            return not self._transport.is_closing()
        else:
            # The transport is None if no connection is established
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
        elif isinstance(request, PlcWriteRequest):
            return await self._write(request)

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
            # If no device is set, log an error and return a response with the NOT_CONNECTED code
            logging.error("No device is set in the modbus connection!")
            return self._default_failed_request(PlcResponseCode.NOT_CONNECTED)

        try:
            # Send the read request to the device and wait for a response
            logging.debug("Sending read request to Modbus Device")
            response = await asyncio.wait_for(
                self._device.read(request, self._transport), 5
            )
            # Return the response
            return response
        except Exception:
            # If an error occurs during the execution of the read request, return a response with
            # the INTERNAL_ERROR code
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
            logging.error("No device is set in the modbus connection!")
            return self._default_failed_request(PlcResponseCode.NOT_CONNECTED)

        try:
            # Send the write request to the device and wait for a response
            logging.debug("Sending write request to Modbus Device")
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


class ModbusDriver(PlcDriver):
    def __init__(self):
        super().__init__()
        self.protocol_code = "modbus-tcp"
        self.protocol_name = "Modbus TCP"

    """
    Creates a new PlcConnection instance using the given URL.

    This method is an asynchronous entry point to connect to a PLC using the Modbus TCP protocol.
    It takes a URL (e.g. "tcp://<ip address>:<port>") as input and an optional PlcAuthentication instance.
    The authentication instance is used to authenticate the connection. If it is not provided, the connection will
    be attempted without authentication.

    :param url: URL of the PLC to connect to (e.g. "tcp://<ip address>:<port>")
    :param authentication: Optional PlcAuthentication instance used to authenticate the connection.
    :return: An instance of PlcConnection that is connected to the PLC
    """
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
        """
        Returns the ModbusDriver class used to instantiate Modbus plc drivers.

        This method is an entry point for pluggy to load the Modbus driver.

        :return Type[ModbusDriver]: ModbusDriver class used to instantiate Modbus plc drivers
        """
        return ModbusDriver

    @staticmethod
    @plc4py.hookimpl
    def key() -> str:
        """
        Returns the unique key for the Modbus driver.

        The key is used by PlcDriverManager to identify the driver.

        :return str: Unique key for the Modbus driver
        """
        return "modbus"
