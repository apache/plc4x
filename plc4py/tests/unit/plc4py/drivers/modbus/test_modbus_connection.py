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
from unittest import TestCase

import pytest

from plc4py.PlcDriverManager import PlcDriverManager
from plc4py.api.value.PlcValue import PlcResponseCode
import logging

from plc4py.spi.values.PlcValues import PlcINT, PlcREAL

logger = logging.getLogger("testing")
TEST_SERVER_IP = "192.168.190.132"


@pytest.mark.asyncio
@pytest.mark.xfail
async def manual_test_plc_driver_modbus_connect():
    """
    Test the connection to a Modbus PLC using PlcDriverManager.
    """
    # Initialize the PlcDriverManager
    driver_manager = PlcDriverManager()

    # Establish a connection to the Modbus PLC
    async with driver_manager.connection(f"modbus://{TEST_SERVER_IP}") as connection:
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
    async with driver_manager.connection(
        f"modbus://{TEST_SERVER_IP}:502"
    ) as connection:
        with connection.read_request_builder() as builder:
            builder.add_item("Random Tag", "0x00001")
            request = builder.build()
            response = await connection.execute(request)
            value = response.tags["Random Tag"].value
            assert value == True


@pytest.mark.asyncio
@pytest.mark.xfail
async def test_plc_driver_modbus_read_coil_non_bool():
    """
    Test reading data from a Modbus PLC.
    """
    log = logging.getLogger(__name__)

    # Initialize the PlcDriverManager
    driver_manager = PlcDriverManager()

    # Establish a connection to the Modbus PLC
    async with driver_manager.connection(
        f"modbus://{TEST_SERVER_IP}:502"
    ) as connection:
        with connection.read_request_builder() as builder:
            builder.add_item("Random Tag", "0x00001:REAL")
            request = builder.build()
            TestCase.assertRaises(
                await connection.execute(request), NotImplementedError
            )


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
    async with driver_manager.connection(
        f"modbus://{TEST_SERVER_IP}:502"
    ) as connection:
        with connection.read_request_builder() as builder:
            builder.add_item("Random Tag", "0x00001[2]")
            request = builder.build()
            response = await connection.execute(request)
            value = response.tags["Random Tag"].value
            assert value == [True, False]


@pytest.mark.asyncio
@pytest.mark.xfail
async def test_plc_driver_modbus_read_contacts():
    """
    Test reading data from a Modbus PLC.
    """
    log = logging.getLogger(__name__)

    # Initialize the PlcDriverManager
    driver_manager = PlcDriverManager()

    # Establish a connection to the Modbus PLC
    async with driver_manager.connection(
        f"modbus://{TEST_SERVER_IP}:502"
    ) as connection:
        with connection.read_request_builder() as builder:
            builder.add_item("Random Tag", "1x00001")
            request = builder.build()
            response = await connection.execute(request)
            value = response.tags["Random Tag"].value
            assert value == True


@pytest.mark.asyncio
@pytest.mark.xfail
async def test_plc_driver_modbus_read_contact_array():
    """
    Test reading data from a Modbus PLC.
    """
    log = logging.getLogger(__name__)

    # Initialize the PlcDriverManager
    driver_manager = PlcDriverManager()

    # Establish a connection to the Modbus PLC
    async with driver_manager.connection(
        f"modbus://{TEST_SERVER_IP}:502"
    ) as connection:
        with connection.read_request_builder() as builder:
            builder.add_item("Random Tag", "1x00001[2]")
            request = builder.build()
            response = await connection.execute(request)
            value = response.tags["Random Tag"].value
            assert value == [True, False]


@pytest.mark.asyncio
@pytest.mark.xfail
async def test_plc_driver_modbus_read_input_register():
    """
    Test reading data from a Modbus PLC.
    """
    log = logging.getLogger(__name__)

    # Initialize the PlcDriverManager
    driver_manager = PlcDriverManager()

    # Establish a connection to the Modbus PLC
    async with driver_manager.connection(
        f"modbus://{TEST_SERVER_IP}:502"
    ) as connection:
        with connection.read_request_builder() as builder:
            builder.add_item("Random Tag", "3x00001")
            request = builder.build()
            response = await connection.execute(request)
            value = response.tags["Random Tag"].value
            assert value == 333


@pytest.mark.asyncio
@pytest.mark.xfail
async def test_plc_driver_modbus_read_input_register_array():
    """
    Test reading data from a Modbus PLC.
    """
    log = logging.getLogger(__name__)

    # Initialize the PlcDriverManager
    driver_manager = PlcDriverManager()

    # Establish a connection to the Modbus PLC
    async with driver_manager.connection(
        f"modbus://{TEST_SERVER_IP}:502"
    ) as connection:
        with connection.read_request_builder() as builder:
            builder.add_item("Random Tag", "3x00001[2]")
            request = builder.build()
            response = await connection.execute(request)
            value = response.tags["Random Tag"].value
            assert value == [333, 0]


@pytest.mark.asyncio
@pytest.mark.xfail
async def test_plc_driver_modbus_read_holding():
    """
    Test reading data from a Modbus PLC.
    """
    log = logging.getLogger(__name__)

    # Initialize the PlcDriverManager
    driver_manager = PlcDriverManager()

    # Establish a connection to the Modbus PLC
    async with driver_manager.connection(
        f"modbus://{TEST_SERVER_IP}:502"
    ) as connection:
        with connection.read_request_builder() as builder:
            builder.add_item("Random Tag", "4x00001")
            request = builder.build()
            response = await connection.execute(request)
            value = response.tags["Random Tag"].value
            assert value == 874


@pytest.mark.asyncio
@pytest.mark.xfail
async def test_plc_driver_modbus_read_holding():
    """
    Test reading data from a Modbus PLC.
    """
    log = logging.getLogger(__name__)

    # Initialize the PlcDriverManager
    driver_manager = PlcDriverManager()

    # Establish a connection to the Modbus PLC
    async with driver_manager.connection(
        f"modbus://{TEST_SERVER_IP}:502"
    ) as connection:
        with connection.read_request_builder() as builder:
            builder.add_item("Random Tag", "4x00001[2]")
            request = builder.build()
            response = await connection.execute(request)
            value = response.tags["Random Tag"].value
            assert value == [874, 0]


@pytest.mark.asyncio
@pytest.mark.xfail
async def test_plc_driver_modbus_read_holding_real():
    """
    Test reading data from a Modbus PLC.
    """
    log = logging.getLogger(__name__)

    # Initialize the PlcDriverManager
    driver_manager = PlcDriverManager()

    # Establish a connection to the Modbus PLC
    async with driver_manager.connection(
        f"modbus://{TEST_SERVER_IP}:502?byte_order=BIG_ENDIAN_WORD_SWAP"
    ) as connection:
        with connection.read_request_builder() as builder:
            builder.add_item("Random Tag", "4x00011:REAL[2]")
            request = builder.build()
            response = await connection.execute(request)
            value = response.tags["Random Tag"].value
            assert value == [PlcREAL(value=874), PlcREAL(value=0.0)]


@pytest.mark.asyncio
@pytest.mark.xfail
async def test_plc_driver_modbus_read_holding_string_even():
    """
    Test reading data from a Modbus PLC.
    """
    log = logging.getLogger(__name__)

    # Initialize the PlcDriverManager
    driver_manager = PlcDriverManager()

    # Establish a connection to the Modbus PLC
    async with driver_manager.connection(
        f"modbus://{TEST_SERVER_IP}:502"
    ) as connection:
        with connection.read_request_builder() as builder:
            builder.add_item("Random Tag", "4x00041:CHAR[6]")
            request = builder.build()
            response = await connection.execute(request)
            value = response.tags["Random Tag"].value
            assert value == [b"F", b"A", b"F", b"B", b"C", b"B"]


@pytest.mark.asyncio
@pytest.mark.xfail
async def test_plc_driver_modbus_read_holding_string_odd():
    """
    Test reading data from a Modbus PLC.
    """
    log = logging.getLogger(__name__)

    # Initialize the PlcDriverManager
    driver_manager = PlcDriverManager()

    # Establish a connection to the Modbus PLC
    async with driver_manager.connection(
        f"modbus://{TEST_SERVER_IP}:502"
    ) as connection:
        with connection.read_request_builder() as builder:
            builder.add_item("Random Tag", "4x00041:CHAR[5]")
            request = builder.build()
            response = await connection.execute(request)
            value = response.tags["Random Tag"].value
            assert value == [b"F", b"A", b"F", b"B", b"C"]


@pytest.mark.asyncio
@pytest.mark.xfail
async def test_plc_driver_modbus_write_holding_int():
    """
    Test reading data from a Modbus PLC.
    """
    log = logging.getLogger(__name__)

    # Initialize the PlcDriverManager
    driver_manager = PlcDriverManager()

    # Establish a connection to the Modbus PLC
    async with driver_manager.connection(
        f"modbus://{TEST_SERVER_IP}:502"
    ) as connection:
        with connection.write_request_builder() as builder:
            builder.add_item("Random Tag", "4x00001", PlcINT(874))
            request = builder.build()
            response = await connection.execute(request)
            value = response.tags["Random Tag"]
            assert value.response_code == PlcResponseCode.OK


@pytest.mark.asyncio
@pytest.mark.xfail
async def test_plc_driver_modbus_write_holding_real():
    """
    Test reading data from a Modbus PLC.
    """
    log = logging.getLogger(__name__)

    # Initialize the PlcDriverManager
    driver_manager = PlcDriverManager()

    # Establish a connection to the Modbus PLC
    async with driver_manager.connection(
        f"modbus://{TEST_SERVER_IP}:502?byte_order=BIG_ENDIAN_WORD_SWAP"
    ) as connection:
        with connection.write_request_builder() as builder:
            builder.add_item("Random Tag", "4x00011:REAL", PlcREAL(874))
            request = builder.build()
            response = await connection.execute(request)
            value = response.tags["Random Tag"]
            assert value.response_code == PlcResponseCode.OK
