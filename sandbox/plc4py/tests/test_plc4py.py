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
from typing import cast

from plc4py import __version__
from plc4py.PlcDriverManager import PlcDriverManager
from plc4py.api.PlcConnection import PlcConnection
from plc4py.api.messages.PlcRequest import PlcFieldRequest
from plc4py.api.messages.PlcResponse import PlcReadResponse
from plc4py.api.value.PlcValue import PlcResponseCode
from plc4py.drivers.mock.MockConnection import MockConnection
from plc4py.drivers.modbus.ModbusConnection import ModbusConnection


def test_version():
    assert __version__ == "0.1.0"


async def test_plc_driver_manager_init():
    driver_manager = PlcDriverManager()
    async with driver_manager.connection("mock:tcp://127.0.0.1:502") as connection:
        assert isinstance(connection, PlcConnection)


async def manual_test_plc_driver_manager_init_modbus():
    driver_manager = PlcDriverManager()
    async with driver_manager.connection("nodbus:tcp://127.0.0.1:502") as connection:
        assert isinstance(connection, ModbusConnection)


async def test_plc_driver_manager_init_mock():
    driver_manager = PlcDriverManager()
    async with driver_manager.connection("mock:tcp://127.0.0.1:502") as connection:
        assert isinstance(connection, MockConnection)


async def test_plc_driver_manager_init_mock_read_request():
    driver_manager = PlcDriverManager()
    field = "1:BOOL"

    async with driver_manager.connection("mock:tcp://127.0.0.1:502") as connection:
        with connection.read_request_builder() as builder:
            builder.add_item(field)
            request: PlcFieldRequest = builder.build()
            response: PlcReadResponse = cast(
                PlcReadResponse, await connection.execute(request)
            )

    # verify that request has one field
    assert response.code == PlcResponseCode.OK
