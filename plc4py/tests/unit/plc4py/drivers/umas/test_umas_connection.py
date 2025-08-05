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
import asyncio
import datetime
import logging
import time
from typing import AsyncGenerator

import pytest
import pytest_asyncio
from plc4py.api.PlcConnection import PlcConnection

from plc4py.api.value.PlcValue import PlcResponseCode
from plc4py.PlcDriverManager import PlcDriverManager
from plc4py.spi.values.PlcValues import (
    PlcBOOL,
    PlcINT,
    PlcREAL,
    PlcDINT,
    PlcSTRING,
    PlcTIME,
    PlcBYTE,
    PlcDATE,
    PlcDATE_AND_TIME,
    PlcList,
)


@pytest_asyncio.fixture
async def connection() -> AsyncGenerator[PlcConnection, None]:
    driver_manager = PlcDriverManager()
    async with driver_manager.connection("umas://127.0.0.1:502") as connection:
        yield connection


@pytest.mark.asyncio
@pytest.mark.xfail
async def test_plc_driver_umas_connect(connection):
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


@pytest.mark.asyncio
@pytest.mark.xfail
async def test_plc_driver_umas_read_boolean(connection):
    tag_name = "TESTING"
    await write_read(connection, tag_name, PlcBOOL(True))


@pytest.mark.asyncio
@pytest.mark.xfail
async def test_plc_driver_umas_read_boolean_with_data_type(connection):
    tag_name = "TESTING:BOOL"
    await write_read(connection, tag_name, PlcBOOL(True))


@pytest.mark.asyncio
@pytest.mark.xfail
async def test_plc_driver_umas_read_int(connection):
    tag_name = "TESTING_INT"
    await write_read(connection, tag_name, PlcINT(99))


@pytest.mark.asyncio
@pytest.mark.xfail
async def test_plc_driver_umas_read_int_with_data_type(connection):
    tag_name = "TESTING_INT:INT"
    await write_read(connection, tag_name, PlcINT(99))


@pytest.mark.asyncio
@pytest.mark.xfail
async def test_plc_driver_umas_read_dint(connection):
    tag_name = "TESTING_DINT"
    await write_read(connection, tag_name, PlcDINT(763539))


@pytest.mark.asyncio
@pytest.mark.xfail
async def test_plc_driver_umas_read_dint_with_data_type(connection):
    tag_name = "TESTING_DINT:DINT"
    await write_read(connection, tag_name, PlcDINT(763539))


@pytest.mark.asyncio
@pytest.mark.xfail
async def test_plc_driver_umas_read_ebool(connection):
    tag_name = "TESTING_EBOOL"
    await write_read(connection, tag_name, PlcBOOL(True))


@pytest.mark.asyncio
@pytest.mark.xfail
async def test_plc_driver_umas_read_ebool_with_data_type(connection):
    tag_name = "TESTING_EBOOL:BOOL"
    await write_read(connection, tag_name, PlcBOOL(True))


@pytest.mark.asyncio
@pytest.mark.xfail
async def test_plc_driver_umas_read_string(connection):
    tag_name = "TESTING_STRING"
    await write_read(connection, tag_name, PlcSTRING("Hello pyToddy!"))


@pytest.mark.asyncio
@pytest.mark.xfail
async def test_plc_driver_umas_read_string_with_data_type(connection):
    tag_name = "TESTING_STRING:STRING"
    await write_read(connection, tag_name, PlcSTRING("Hello pyToddy!"))


@pytest.mark.asyncio
@pytest.mark.xfail
async def test_plc_driver_umas_read_time(connection):
    tag_name = "TESTING_TIME"
    await write_read(connection, tag_name, PlcTIME(200000))


@pytest.mark.asyncio
@pytest.mark.xfail
async def test_plc_driver_umas_read_time_with_data_type(connection):
    tag_name = "TESTING_TIME:TIME"
    await write_read(connection, tag_name, PlcTIME(200000))


@pytest.mark.asyncio
@pytest.mark.xfail
async def test_plc_driver_umas_read_byte(connection):
    tag_name = "TESTING_BYTE"
    await write_read(connection, tag_name, PlcBYTE(253))


@pytest.mark.asyncio
@pytest.mark.xfail
async def test_plc_driver_umas_read_byte_with_data_type(connection):
    tag_name = "TESTING_BYTE:BYTE"
    await write_read(connection, tag_name, PlcBYTE(253))


@pytest.mark.asyncio
@pytest.mark.xfail
async def test_plc_driver_umas_read_date(connection):
    tag_name = "TESTING_DATE"
    await write_read(connection, tag_name, PlcDATE(datetime.datetime(2024, 10, 25)))


@pytest.mark.asyncio
@pytest.mark.xfail
async def test_plc_driver_umas_read_time_with_data_type(connection):
    tag_name = "TESTING_DATE:DATE"
    await write_read(connection, tag_name, PlcDATE(datetime.datetime(2025, 11, 22)))


@pytest.mark.asyncio
@pytest.mark.xfail
async def test_plc_driver_umas_read_dt(connection):
    tag_name = "TESTING_DT"
    await write_read(
        connection, tag_name, PlcDATE_AND_TIME(datetime.datetime(2000, 1, 10, 0, 40))
    )


@pytest.mark.asyncio
@pytest.mark.xfail
async def test_plc_driver_umas_read_dt_with_data_type(connection):
    tag_alias = "Random Tag"
    tag_name = "TESTING_DT:DATE_AND_TIME"
    await write_read(
        connection, tag_name, PlcDATE_AND_TIME(datetime.datetime(2002, 10, 31, 3, 38))
    )


@pytest.mark.asyncio
@pytest.mark.xfail
async def test_plc_driver_umas_read_array(connection):
    tag_name = "TESTING_BYTE_ARRAY:BYTE[10]"
    await write_read(
        connection,
        tag_name,
        PlcList(
            [
                PlcBYTE(1),
                PlcBYTE(2),
                PlcBYTE(3),
                PlcBYTE(4),
                PlcBYTE(5),
                PlcBYTE(6),
                PlcBYTE(7),
                PlcBYTE(8),
                PlcBYTE(9),
                PlcBYTE(10),
            ]
        ),
    )


@pytest.mark.asyncio
@pytest.mark.xfail
async def test_plc_driver_umas_browse(connection):
    with connection.browse_request_builder() as builder:
        builder.add_query("All Tags", "*")
        request = builder.build()

    future = connection.execute(request)
    response = await future

    pass
