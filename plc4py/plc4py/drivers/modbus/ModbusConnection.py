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
from typing import Type

import plc4py
from plc4py.api.authentication.PlcAuthentication import PlcAuthentication
from plc4py.api.messages.PlcRequest import (
    PlcReadRequest,
    PlcRequest,
    PlcWriteRequest,
    ReadRequestBuilder,
    WriteRequestBuilder,
)
from plc4py.api.messages.PlcResponse import (
    PlcResponse,
)
from plc4py.api.PlcConnection import PlcConnection, PlcConnectionMetaData
from plc4py.api.PlcDriver import PlcDriver
from plc4py.api.value.PlcValue import PlcResponseCode
from plc4py.drivers.modbus.ModbusConfiguration import ModbusConfiguration
from plc4py.drivers.modbus.ModbusDevice import ModbusDevice
from plc4py.drivers.modbus.ModbusProtocol import ModbusProtocol
from plc4py.drivers.modbus.ModbusTag import ModbusTagBuilder
from plc4py.drivers.PlcDriverLoader import PlcDriverLoader
from plc4py.spi.messages.PlcReader import DefaultPlcReader
from plc4py.spi.messages.PlcRequest import (
    DefaultReadRequestBuilder,
    DefaultWriteRequestBuilder,
)
from plc4py.spi.messages.PlcWriter import DefaultPlcWriter
from plc4py.spi.transport.Plc4xBaseTransport import Plc4xBaseTransport
from plc4py.spi.transport.TCPTransport import TCPTransport


class ModbusConnection(PlcConnection, DefaultPlcReader, DefaultPlcWriter):
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
                15,
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

    def read_request_builder(self) -> ReadRequestBuilder:
        """
        :return: read request builder.
        """
        return DefaultReadRequestBuilder(ModbusTagBuilder)

    def write_request_builder(self) -> WriteRequestBuilder:
        """
        :return: write request builder.
        """
        return DefaultWriteRequestBuilder(ModbusTagBuilder)


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
    @plc4py.drivers.hookimpl
    def get_driver() -> Type[ModbusDriver]:
        """
        Returns the ModbusDriver class used to instantiate Modbus plc drivers.

        This method is an entry point for pluggy to load the Modbus driver.

        :return Type[ModbusDriver]: ModbusDriver class used to instantiate Modbus plc drivers
        """
        return ModbusDriver

    @staticmethod
    @plc4py.drivers.hookimpl
    def key() -> str:
        """
        Returns the unique key for the Modbus driver.

        The key is used by PlcDriverManager to identify the driver.

        :return str: Unique key for the Modbus driver
        """
        return "modbus"
