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

import pytest

from plc4py import __version__
from plc4py.PlcDriverManager import PlcDriverManager
from plc4py.api.PlcConnection import PlcConnection
from plc4py.api.messages.PlcRequest import PlcTagRequest
from plc4py.api.messages.PlcResponse import PlcReadResponse
from plc4py.api.value.PlcValue import PlcResponseCode
from plc4py.drivers.mock.MockConnection import MockConnection
from plc4py.drivers.modbus.ModbusConnection import ModbusConnection


def test_version():
    """
    This test verifies the version of the package.

    The version is a string that is expected to match the version specified in the
    `__version__` constant at the top of the file. If the versions do not match,
    the test will fail.

    :return: None
    """
    assert __version__ == "0.1.0"


@pytest.mark.asyncio
async def test_plc_driver_manager_init():
    """
    This test verifies that the PlcDriverManager class is able to create a connection
    successfully.

    :return: None
    """
    driver_manager = PlcDriverManager()
    async with driver_manager.connection("mock:tcp://127.0.0.1:502") as connection:
        assert isinstance(connection, PlcConnection)


@pytest.mark.asyncio
async def manual_test_plc_driver_manager_init_modbus():
    """
    This test verifies that the PlcDriverManager class is able to create a connection
    successfully with the Modbus protocol.

    The test creates a PlcDriverManager object and then uses the connection method
    to get a connection to a Modbus PLC. The connection should be an instance of
    the ModbusConnection class.

    """
    driver_manager = PlcDriverManager()
    async with driver_manager.connection("modbus:tcp://127.0.0.1:502") as connection:
        assert isinstance(connection, ModbusConnection)


@pytest.mark.asyncio
async def test_plc_driver_manager_init_mock():
    """
    This test verifies that the PlcDriverManager class is able to create a connection
    successfully with the Mock driver.

    The test creates a PlcDriverManager object and then uses the connection method
    to get a connection to a Mock PLC. The connection should be an instance of
    the MockConnection class.

    """
    driver_manager = PlcDriverManager()
    async with driver_manager.connection("mock:tcp://127.0.0.1:502") as connection:
        assert isinstance(connection, MockConnection)


@pytest.mark.asyncio
async def test_plc_driver_manager_init_mock_read_request():
    """
    This test verifies that the PlcDriverManager class is able to create a
    connection to a mock PLC and then execute a read request successfully.

    The test creates a PlcDriverManager object and then uses the connection method
    to get a connection to a Mock PLC. It then uses the read_request_builder
    to construct a read request with one item. The request is then executed
    and the result should be a PlcReadResponse with a response code of OK.

    """
    driver_manager = PlcDriverManager()
    tag = "1:BOOL"

    async with driver_manager.connection("mock:tcp://127.0.0.1:502") as connection:
        # Create a read request builder
        with connection.read_request_builder() as builder:
            # Add the tag to the request
            builder.add_item("Random Tag", tag)
            # Build the request
            request: PlcTagRequest = builder.build()
            # Execute the request
            response = await connection.execute(request)

    # Verify that the request has one field
    assert response.response_code == PlcResponseCode.OK
