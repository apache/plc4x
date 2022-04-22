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
import asyncio
import logging
from dataclasses import dataclass, field
from typing import Awaitable, Type

import plc4py

from plc4py.api.PlcConnection import PlcConnection
from plc4py.api.exceptions.exceptions import PlcFieldParseException
from plc4py.api.messages.PlcField import PlcField
from plc4py.api.messages.PlcRequest import (
    ReadRequestBuilder,
    PlcReadRequest,
    PlcRequest,
)
from plc4py.api.messages.PlcResponse import PlcReadResponse, PlcResponse
from plc4py.api.value.PlcValue import PlcResponseCode, PlcValue
from plc4py.drivers.PlcConnectionLoader import PlcConnectionLoader
from plc4py.spi.messages.PlcReader import PlcReader
from plc4py.spi.messages.utils.ResponseItem import ResponseItem
from plc4py.spi.values.PlcBOOL import PlcBOOL
from plc4py.spi.values.PlcINT import PlcINT
from plc4py.drivers.mock.MockReadRequestBuilder import MockReadRequestBuilder


@dataclass
class MockPlcField(PlcField):
    datatype: str = "INT"


class MockPlcFieldHandler:
    @staticmethod
    def of(field: str) -> MockPlcField:
        try:
            datatype = field.split(":")[1]
            return MockPlcField(field, datatype)
        except IndexError:
            raise PlcFieldParseException


class MockDevice:
    def read(self, field) -> list[ResponseItem[PlcValue]]:
        logging.debug(f"Reading field {field} from Mock Device")
        plc_field = MockPlcFieldHandler.of(field)
        if plc_field.datatype == "BOOL":
            return [ResponseItem(PlcResponseCode.OK, PlcBOOL(True))]
        elif plc_field.datatype == "INT":
            return [ResponseItem(PlcResponseCode.OK, PlcINT(1))]
        else:
            raise PlcFieldParseException


@dataclass
class MockConnection(PlcConnection, PlcReader):
    _is_connected: bool = False
    device: MockDevice = field(default_factory=lambda: MockDevice())

    def connect(self):
        """
        Connect the Mock PLC connection
        :return:
        """
        self._is_connected = True
        self.device = MockDevice()

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
        return MockReadRequestBuilder()

    def execute(self, request: PlcRequest) -> Awaitable[PlcResponse]:
        """
        Executes a PlcRequest as long as it's already connected
        :param PlcRequest: Plc Request to execute
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
        if self.device is None:
            logging.error("No device is set in the mock connection!")
            return self._default_failed_request(PlcResponseCode.NOT_CONNECTED)

        loop = asyncio.get_running_loop()
        logging.debug("Sending read request to MockDevice")
        future = loop.create_future()

        async def _request(fut, req, device):
            try:
                response = PlcReadResponse(
                    PlcResponseCode.OK,
                    req.fields,
                    {field: device.read(field) for field in req.field_names},
                )
                fut.set_result(response)
            except PlcFieldParseException:
                fut.set_result(
                    PlcReadResponse(PlcResponseCode.INTERNAL_ERROR, req.fields, {})
                )

        loop.create_task(_request(future, request, self.device))
        return future


class MockConnectionLoader(PlcConnectionLoader):
    @staticmethod
    @plc4py.hookimpl
    def get_connection() -> Type[MockConnection]:
        return MockConnection

    @staticmethod
    @plc4py.hookimpl
    def key() -> str:
        return "mock"
