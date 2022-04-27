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

from typing import Type, Awaitable

import plc4py
from plc4py.api.PlcConnection import PlcConnection
from plc4py.api.PlcDriver import PlcDriver
from plc4py.api.authentication.PlcAuthentication import PlcAuthentication
from plc4py.api.messages.PlcResponse import PlcResponse
from plc4py.api.messages.PlcRequest import ReadRequestBuilder
from plc4py.drivers.PlcDriverLoader import PlcDriverLoader


class ModbusConnection(PlcConnection):
    """A hook implementation namespace."""

    def __init__(self, url: str):
        super().__init__(url)

    def connect(self):
        """
        Establishes the connection to the remote PLC.
        """
        pass

    def is_connected(self) -> bool:
        """
        Indicates if the connection is established to a remote PLC.
        :return: True if connection, False otherwise
        """
        pass

    def close(self) -> None:
        """
        Closes the connection to the remote PLC.
        :return:
        """
        pass

    def read_request_builder(self) -> ReadRequestBuilder:
        """
        :return: read request builder.
        """
        pass

    def execute(self, PlcRequest) -> Awaitable[PlcResponse]:
        """
        Executes a PlcRequest as long as it's already connected
        :param PlcRequest: Plc Request to execute
        :return: The response from the Plc/Device
        """
        pass


class ModbusDriver(PlcDriver):
    def __init__(self):
        self.protocol_code = "modbus-tcp"
        self.protocol_name = "Modbus TCP"

    def get_connection(
        self, url: str, authentication: PlcAuthentication = PlcAuthentication()
    ) -> PlcConnection:
        """
        Connects to a PLC using the given plc connection string.
        :param url: plc connection string
        :param authentication: authentication credentials.
        :return PlcConnection: PLC Connection object
        """
        return ModbusConnection(url)


class ModbusDriverLoader(PlcDriverLoader):
    @staticmethod
    @plc4py.hookimpl
    def get_driver() -> Type[ModbusDriver]:
        return ModbusDriver

    @staticmethod
    @plc4py.hookimpl
    def key() -> str:
        return "modbus"
