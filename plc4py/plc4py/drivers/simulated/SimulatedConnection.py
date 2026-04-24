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
from dataclasses import dataclass, field
from typing import Type, Union

import plc4py
from plc4py.api.authentication.PlcAuthentication import PlcAuthentication
from plc4py.api.messages.PlcRequest import (
    PlcReadRequest,
    PlcRequest,
    PlcWriteRequest,
    ReadRequestBuilder,
    WriteRequestBuilder,
)
from plc4py.api.messages.PlcResponse import PlcResponse
from plc4py.api.PlcConnection import PlcConnection
from plc4py.api.PlcDriver import PlcDriver
from plc4py.api.value.PlcValue import PlcResponseCode
from plc4py.drivers.PlcDriverLoader import PlcDriverLoader
from plc4py.drivers.simulated.SimulatedConfiguration import SimulatedConfiguration
from plc4py.drivers.simulated.SimulatedDevice import SimulatedDevice
from plc4py.drivers.simulated.SimulatedTag import SimulatedTagBuilder
from plc4py.spi.messages.PlcReader import DefaultPlcReader
from plc4py.spi.messages.PlcRequest import (
    DefaultReadRequestBuilder,
    DefaultWriteRequestBuilder,
)
from plc4py.spi.messages.PlcWriter import DefaultPlcWriter
from plc4py.spi.transport.MockTransport import MockTransport
from plc4py.spi.transport.Plc4xBaseTransport import Plc4xBaseTransport


@dataclass
class SimulatedConnection(PlcConnection, DefaultPlcReader, DefaultPlcWriter):
    """
    In-process PLC connection backed by a :class:`SimulatedDevice`.

    The connection does not open a network socket; it dispatches read and
    write requests straight to the device. This makes it useful as a
    reference driver and in tests that should not depend on a real PLC.
    """

    _device: SimulatedDevice = field(default_factory=lambda: SimulatedDevice())
    _transport: Union[Plc4xBaseTransport, None] = None

    @staticmethod
    async def create(url: str) -> "SimulatedConnection":
        config = SimulatedConfiguration(url)
        connection = SimulatedConnection(
            _device=SimulatedDevice(name=config.device_name)
        )
        connection._transport = await MockTransport.create(None, None, None)
        return connection

    def read_request_builder(self) -> ReadRequestBuilder:
        return DefaultReadRequestBuilder(SimulatedTagBuilder)

    def write_request_builder(self) -> WriteRequestBuilder:
        return DefaultWriteRequestBuilder(SimulatedTagBuilder)

    async def execute(self, request: PlcRequest) -> PlcResponse:
        if not self.is_connected():
            return self._default_failed_request(PlcResponseCode.NOT_CONNECTED)
        if isinstance(request, PlcReadRequest):
            return await self._read(request)
        if isinstance(request, PlcWriteRequest):
            return await self._write(request)
        return self._default_failed_request(PlcResponseCode.NOT_CONNECTED)


class SimulatedDriver(PlcDriver):
    def __init__(self):
        super().__init__()
        self.protocol_code = "simulated"
        self.protocol_name = "Simulated"

    async def get_connection(
        self,
        url: str,
        authentication: PlcAuthentication = PlcAuthentication(),
    ) -> PlcConnection:
        return await SimulatedConnection.create(url)


class SimulatedDriverLoader(PlcDriverLoader):
    """
    Pluggy hook implementation that lets ``PlcDriverManager`` discover the
    simulated driver via the ``plc4py.drivers`` entry point.
    """

    @staticmethod
    @plc4py.drivers.hookimpl
    def get_driver() -> Type[SimulatedDriver]:
        return SimulatedDriver

    @staticmethod
    @plc4py.drivers.hookimpl
    def key() -> str:
        return "simulated"
