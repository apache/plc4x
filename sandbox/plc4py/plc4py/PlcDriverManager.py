#
# Licensed to the Apache Software Foundation (ASF) under one
# or more contributor license agreements.  See the NOTICE file
# distributed with this work for additional information
# regarding copyright ownership.  The ASF licenses this file
# to you under the Apache License, Version 2.0 (the
# "License"); you may not use this file except in compliance
# with the License.  You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing,
# software distributed under the License is distributed on an
# "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
# KIND, either express or implied.  See the License for the
# specific language governing permissions and limitations
# under the License.
#

import sys
from contextlib import contextmanager
from dataclasses import dataclass, field
from typing import Generator, Type

from pluggy import PluginManager  # type: ignore

from plc4py.api.PlcConnection import PlcConnection
from plc4py.drivers.modbus.ModbusConnection import ModbusConnectionLoader
from plc4py.spi.PlcDriverClassLoader import PlcDriverClassLoader


@dataclass
class PlcDriverManager:
    class_loader: PluginManager = field(default_factory=lambda: PluginManager("plc4py"))
    _driverMap: dict[str, Type[PlcConnection]] = field(default_factory=lambda: {})

    def __post_init__(self):
        self.class_loader.add_hookspecs(PlcDriverClassLoader)
        self.class_loader.register(ModbusConnectionLoader)
        self._driverMap = {key: loader for key, loader in zip(self.class_loader.hook.key(),
                                                              self.class_loader.hook.get_type())}

    @contextmanager
    def connection(self, url: str) -> Generator[PlcConnection, None, None]:
        """
        Context manager to handle connection.
        """
        conn = None
        try:
            conn = self.get_connection(url)
            yield conn
        finally:
            if conn is not None:
                conn.close()

    def get_connection(self, url: str) -> PlcConnection:
        driver_name = url.split(':', 1)[0]
        return self._driverMap[driver_name](url)

