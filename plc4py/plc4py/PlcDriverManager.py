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
from contextlib import asynccontextmanager
from dataclasses import dataclass, field
from typing import AsyncIterator, Dict, List, Type

from pluggy import PluginManager  # type: ignore

from plc4py.api.PlcConnection import PlcConnection
from plc4py.api.PlcDriver import PlcDriver
from plc4py.spi.PlcDriverClassLoader import PlcDriverClassLoader
from plc4py.utils.ConnectionStringHandling import get_protocol_code
from plc4py.spi.PlcDriverClassLoader import PlcTransportClassLoader
from plc4py.spi.transport.Plc4xBaseTransport import Plc4xBaseTransport


@dataclass
class PlcDriverManager:
    class_loader: PluginManager = field(
        default_factory=lambda: PluginManager("plc4py.drivers")
    )
    transport_loader: PluginManager = field(
        default_factory=lambda: PluginManager("plc4py.transports")
    )
    _driver_map: Dict[str, Type[PlcDriver]] = field(default_factory=lambda: {})
    _transport_map: Dict[str, Type[Plc4xBaseTransport]] = field(
        default_factory=lambda: {}
    )

    def __post_init__(self):
        """
        Initialize the PLC Driver Manager.

        This function is called automatically when a new instance of the
        PlcDriverManager class is created.

        It registers the available drivers, by loading the entry points
        defined in the "plc4py.drivers" namespace.
        """
        # Log the class loader used
        logging.info(
            "Instantiating new PLC Driver Manager with class loader %s",
            self.class_loader,
        )

        # Add the PlcDriverClassLoader hookspecs to the class loader
        self.class_loader.add_hookspecs(PlcDriverClassLoader)

        # Log the registration of drivers
        logging.info("Registering available drivers...")

        # Register the plc4py.drivers package
        import plc4py.drivers

        self.class_loader.register(plc4py.drivers)

        # Load the setuptools entry points defined in the "plc4py.drivers" namespace
        self.class_loader.load_setuptools_entrypoints("plc4py.drivers")

        # Create a dictionary mapping the hook names to the PlcDriver instances
        self._driver_map = {
            key: loader
            for key, loader in zip(
                self.class_loader.hook.key(),
                self.class_loader.hook.get_driver(),
            )
        }

        # Log the successful registration of each driver
        for driver in self._driver_map:
            logging.info("... %s .. OK", driver)

        # Check for any pending plugins
        self.class_loader.check_pending()

        ######
        # Log the class loader used
        logging.info(
            "Instantiating new PLC Transport Manager with class loader %s",
            self.transport_loader,
        )

        # Add the PlcDriverClassLoader hookspecs to the class loader
        self.transport_loader.add_hookspecs(PlcTransportClassLoader)

        # Log the registration of drivers
        logging.info("Registering available transports...")

        # Register the plc4py.drivers package
        import plc4py.spi.transport

        self.transport_loader.register(plc4py.spi.transport)

        # Load the setuptools entry points defined in the "plc4py.drivers" namespace
        self.transport_loader.load_setuptools_entrypoints("plc4py.transports")

        # Create a dictionary mapping the hook names to the PlcDriver instances
        self._transport_map = {
            key: loader
            for key, loader in zip(
                self.transport_loader.hook.key(),
                self.transport_loader.hook.get_transport(),
            )
        }

        # Log the successful registration of each driver
        for transport in self._transport_map:
            logging.info("... %s .. OK", transport)

        # Check for any pending plugins
        self.transport_loader.check_pending()

    @asynccontextmanager
    async def connection(self, url: str) -> AsyncIterator[PlcConnection]:
        """
        Context manager to handle connection.

        :param url: plc connection string
        :return: plc connection generator
        """
        conn: PlcConnection = await self.get_connection(url)
        try:
            yield conn
        finally:
            conn.close()

    async def get_connection(self, url: str) -> PlcConnection:
        """
        Connects to a PLC using the given plc connection string using given authentication credentials.
        :param url: plc connection string.
        :return: plc connection
        """
        protocol_code = get_protocol_code(url)
        return await self._driver_map[protocol_code]().get_connection(url)

    def list_drivers(self) -> List[str]:
        """
        Returns the codes of the drivers which are currently registered at the PlcDriverManager
        :return: Set of driver codes for all drivers registered
        """
        return list(self._driver_map.keys())

    def get_driver(self, protocol_code: str) -> Type[PlcDriver]:
        """
        Returns suitable driver for protocol or throws an Exception.
        :param protocol_code: protocolCode protocol code identifying the driver
        :return: Driver instance for the given protocol
        """
        return self._driver_map[protocol_code]

    def get_driver_for_url(self, url: str) -> Type[PlcDriver]:
        """
        Returns the driver class that matches the one identified within the connection string.

        :param url: The plc connection string
        :return: The driver class that matches the protocol code identified within the connection string
        """
        # Extract the protocol code from the connection string
        protocol_code = get_protocol_code(url)

        # Retrieve the driver class using the protocol code
        return self._driver_map[protocol_code]
