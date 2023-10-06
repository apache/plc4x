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
from typing import Type, AsyncIterator, List, Dict
from pluggy import PluginManager  # type: ignore

from plc4py.api.PlcConnection import PlcConnection
from plc4py.api.PlcDriver import PlcDriver
from plc4py.spi.PlcDriverClassLoader import PlcDriverClassLoader
from plc4py.utils.ConnectionStringHandling import get_protocol_code


@dataclass
class PlcDriverManager:
    class_loader: PluginManager = field(default_factory=lambda: PluginManager("plc4py"))
    _driverMap: Dict[str, Type[PlcDriver]] = field(default_factory=lambda: {})

    def __post_init__(self):
        logging.info(
            f"Instantiating new PLC Driver Manager with class loader {self.class_loader}"
        )
        self.class_loader.add_hookspecs(PlcDriverClassLoader)
        logging.info("Registering available drivers...")

        import plc4py.drivers

        self.class_loader.register(plc4py.drivers)
        self.class_loader.load_setuptools_entrypoints("plc4py.drivers")
        self._driverMap = {
            key: loader
            for key, loader in zip(
                self.class_loader.hook.key(), self.class_loader.hook.get_driver()
            )
        }
        for driver in self._driverMap:
            logging.info(f"... {driver} .. OK")
        self.class_loader.check_pending()

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
        return await self._driverMap[protocol_code]().get_connection(url)

    def list_drivers(self) -> List[str]:
        """
        Returns the codes of the drivers which are currently registered at the PlcDriverManager
        :return: Set of driver codes for all drivers registered
        """
        return list(self._driverMap.keys())

    def get_driver(self, protocol_code: str) -> Type[PlcDriver]:
        """
        Returns suitable driver for protocol or throws an Exception.
        :param protocol_code: protocolCode protocol code identifying the driver
        :return: Driver instance for the given protocol
        """
        return self._driverMap[protocol_code]

    def get_driver_for_url(self, url: str) -> Type[PlcDriver]:
        """
        Returns the driver class that matches that identified within the connection string
        :param url: The plc connection string
        :return: the protocol code
        """
        protocol_code = get_protocol_code(url)
        return self._driverMap[protocol_code]
