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
from typing import AsyncGenerator
from unittest import TestCase

import pytest
import pytest_asyncio
from plc4py.api.PlcConnection import PlcConnection

from plc4py.PlcDriverManager import PlcDriverManager
from plc4py.api.value.PlcValue import PlcResponseCode
import logging

from plc4py.spi.values.PlcValues import PlcINT, PlcREAL, PlcList, PlcBOOL, PlcCHAR

logger = logging.getLogger("testing")


@pytest_asyncio.fixture
async def connection() -> AsyncGenerator[PlcConnection, None]:
    driver_manager = PlcDriverManager()
    async with driver_manager.connection("modbus://127.0.0.1:502") as connection:
        yield connection


@pytest.mark.asyncio
@pytest.mark.xfail
async def test_plc_driver_modbus_connect(connection):
    assert connection.is_connected


async def write_read(connection, tag_name, tag_value):
    tag_alias = "Some Random Alias"
    with connection.write_request_builder() as builder:
        builder.add_item(tag_alias, tag_name, tag_value)
        write_request = builder.build()
    with connection.read_request_builder() as builder:
        builder.add_item(tag_alias, tag_name)
        read_request = builder.build()
    future = connection.execute(write_request)
    response = await future
    assert response.response_code == PlcResponseCode.OK

    future = connection.execute(read_request)
    response = await future
    assert response.response_code == PlcResponseCode.OK

    value = response.tags[tag_alias].value
    response_code = response.tags[tag_alias].response_code
    assert value == tag_value
    assert response_code == PlcResponseCode.OK


async def simple_read(connection, tag_name, tag_value):
    tag_alias = "Some Random Alias"
    with connection.read_request_builder() as builder:
        builder.add_item(tag_alias, tag_name)
        read_request = builder.build()

    future = connection.execute(read_request)
    response = await future
    assert response.response_code == PlcResponseCode.OK

    value = response.tags[tag_alias].value
    response_code = response.tags[tag_alias].response_code
    assert value == tag_value
    assert response_code == PlcResponseCode.OK


@pytest.mark.asyncio
@pytest.mark.xfail
async def test_plc_driver_modbus_read_boolean(connection):
    tag_name = "0x00001"
    await write_read(connection, tag_name, PlcBOOL(True))


@pytest.mark.asyncio
@pytest.mark.xfail
async def test_plc_driver_modbus_read_int(connection):
    tag_name = "4x00001"
    await write_read(connection, tag_name, PlcINT(83))


@pytest.mark.asyncio
@pytest.mark.xfail
async def test_plc_driver_modbus_read_real(connection):
    tag_name = "4x00001"
    await write_read(connection, tag_name, PlcINT(83))


@pytest.mark.asyncio
@pytest.mark.xfail
async def test_plc_driver_modbus_read_boolean_array(connection):
    tag_name = "0x00001[2]"
    await write_read(connection, tag_name, PlcList([PlcBOOL(False), PlcBOOL(True)]))


@pytest.mark.asyncio
@pytest.mark.xfail
async def test_plc_driver_modbus_read_int(connection):
    tag_name = "4x00001"
    await write_read(connection, tag_name, PlcINT(83))


@pytest.mark.asyncio
@pytest.mark.xfail
async def test_plc_driver_modbus_read_bool_array(connection):
    tag_name = "0x00001[2]"
    await write_read(connection, tag_name, PlcList([PlcBOOL(True), PlcBOOL(False)]))


@pytest.mark.asyncio
@pytest.mark.xfail
async def test_plc_driver_modbus_read_contacts(connection):
    tag_name = "1x00001"
    await simple_read(connection, tag_name, PlcBOOL(False))


@pytest.mark.asyncio
@pytest.mark.xfail
async def test_plc_driver_modbus_read_bool_array_discrete_input(connection):
    tag_name = "1x00001[2]"
    await simple_read(connection, tag_name, PlcList([PlcBOOL(False), PlcBOOL(False)]))


@pytest.mark.asyncio
@pytest.mark.xfail
async def test_plc_driver_modbus_read_input_register(connection):
    tag_name = "3x00001"
    await simple_read(connection, tag_name, PlcINT(0))


@pytest.mark.asyncio
@pytest.mark.xfail
async def test_plc_driver_modbus_read_input_register_array(connection):
    tag_name = "3x00001[2]"
    await simple_read(connection, tag_name, PlcList([PlcINT(0), PlcINT(0)]))


@pytest.mark.asyncio
@pytest.mark.xfail
async def test_plc_driver_modbus_read_holding_register(connection):
    tag_name = "4x00001"
    await write_read(connection, tag_name, PlcINT(334))


@pytest.mark.asyncio
@pytest.mark.xfail
async def test_plc_driver_modbus_read_holdiong_register_array(connection):
    tag_name = "4x00001[2]"
    await write_read(connection, tag_name, PlcList([PlcINT(334), PlcINT(0)]))


@pytest.mark.asyncio
@pytest.mark.xfail
async def test_plc_driver_modbus_read_holding_register_real(connection):
    tag_name = "4x00011:REAL[2]"
    await write_read(
        connection, tag_name, PlcList([PlcREAL(value=874), PlcREAL(value=0.0)])
    )


@pytest.mark.asyncio
@pytest.mark.xfail
async def test_plc_driver_modbus_read_holding_char(connection):
    tag_name = "4x00041:CHAR"
    await write_read(
        connection,
        tag_name,
        PlcCHAR("F"),
    )


@pytest.mark.asyncio
@pytest.mark.xfail
async def test_plc_driver_modbus_read_holding_string_even(connection):
    tag_name = "4x00041:CHAR[6]"
    await write_read(
        connection,
        tag_name,
        PlcList(
            [
                PlcCHAR(value="F"),
                PlcCHAR(value="A"),
                PlcCHAR(value="F"),
                PlcCHAR(value="B"),
                PlcCHAR(value="C"),
                PlcCHAR(value="B"),
            ]
        ),
    )


@pytest.mark.asyncio
@pytest.mark.xfail
async def test_plc_driver_modbus_read_holding_string_odd(connection):
    tag_name = "4x00041:CHAR[5]"
    await write_read(
        connection,
        tag_name,
        PlcList(
            [
                PlcCHAR(value=b"F"),
                PlcCHAR(value=b"A"),
                PlcCHAR(value=b"F"),
                PlcCHAR(value=b"B"),
                PlcCHAR(value=b"C"),
            ]
        ),
    )
