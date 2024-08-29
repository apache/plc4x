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
import time

import pytest

from plc4py.PlcDriverManager import PlcDriverManager
from plc4py.api.value.PlcValue import PlcResponseCode
import logging

from plc4py.spi.values.PlcValues import PlcINT

logger = logging.getLogger("testing")


@pytest.mark.asyncio
async def manual_test_plc_driver_modbus_connect():
    """
    Test the connection to a Modbus PLC using PlcDriverManager.
    """
    # Initialize the PlcDriverManager
    driver_manager = PlcDriverManager()

    # Establish a connection to the Modbus PLC
    async with driver_manager.connection("modbus://1") as connection:
        # Check if the connection is successful
        assert connection.is_connected()

    # Ensure the connection is closed after exiting the context manager
    assert not connection.is_connected()


@pytest.mark.asyncio
@pytest.mark.xfail
async def test_plc_driver_modbus_read_coil():
    """
    Test reading data from a Modbus PLC.
    """
    log = logging.getLogger(__name__)

    # Initialize the PlcDriverManager
    driver_manager = PlcDriverManager()

    # Establish a connection to the Modbus PLC
    async with driver_manager.connection("modbus://192.168.174.128:502") as connection:
        with connection.read_request_builder() as builder:
            builder.add_item("Random Tag", "0x00001[9]")
            request = builder.build()
            response = await connection.execute(request)
            value = response.tags["Random Tag"].value


@pytest.mark.asyncio
@pytest.mark.xfail
async def test_plc_driver_modbus_read_coil_array():
    """
    Test reading data from a Modbus PLC.
    """
    log = logging.getLogger(__name__)

    # Initialize the PlcDriverManager
    driver_manager = PlcDriverManager()

    # Establish a connection to the Modbus PLC
    async with driver_manager.connection("modbus://192.168.174.128:502") as connection:
        with connection.read_request_builder() as builder:
            builder.add_item("Random Tag", "0x00001[10]")
            request = builder.build()
            response = await connection.execute(request)
            value = response.tags["Random Tag"].value


@pytest.mark.asyncio
async def test_plc_driver_modbus_read_contact_array():
    """
    Test reading data from a Modbus PLC.
    """
    log = logging.getLogger(__name__)

    # Initialize the PlcDriverManager
    driver_manager = PlcDriverManager()

    # Establish a connection to the Modbus PLC
    async with driver_manager.connection("modbus://192.168.174.128:502") as connection:
        with connection.read_request_builder() as builder:
            builder.add_item("Random Tag", "1x00001[10]")
            request = builder.build()
            response = await connection.execute(request)
            value = response.tags["Random Tag"].value
            pass


@pytest.mark.asyncio
async def test_plc_driver_modbus_read_input_register_array():
    """
    Test reading data from a Modbus PLC.
    """
    log = logging.getLogger(__name__)

    # Initialize the PlcDriverManager
    driver_manager = PlcDriverManager()

    # Establish a connection to the Modbus PLC
    async with driver_manager.connection("modbus://192.168.174.128:502") as connection:
        with connection.read_request_builder() as builder:
            builder.add_item("Random Tag", "3x00001")
            request = builder.build()
            response = await connection.execute(request)
            value = response.tags["Random Tag"].value
            pass


@pytest.mark.asyncio
async def test_plc_driver_modbus_read_holding_array():
    """
    Test reading data from a Modbus PLC.
    """
    log = logging.getLogger(__name__)

    # Initialize the PlcDriverManager
    driver_manager = PlcDriverManager()

    # Establish a connection to the Modbus PLC
    async with driver_manager.connection("modbus://192.168.174.128:502") as connection:
        with connection.read_request_builder() as builder:
            builder.add_item("Random Tag", "4x00001")
            request = builder.build()
            response = await connection.execute(request)
            value = response.tags["Random Tag"].value
            pass


@pytest.mark.asyncio
async def test_plc_driver_modbus_write_holding():
    """
    Test reading data from a Modbus PLC.
    """
    log = logging.getLogger(__name__)

    # Initialize the PlcDriverManager
    driver_manager = PlcDriverManager()

    # Establish a connection to the Modbus PLC
    async with driver_manager.connection("modbus://192.168.174.128:502") as connection:
        with connection.write_request_builder() as builder:
            builder.add_item("Random Tag", "4x00001", 2)
            request = builder.build()
            response = await connection.execute(request)
            value = response.tags["Random Tag"]
            pass