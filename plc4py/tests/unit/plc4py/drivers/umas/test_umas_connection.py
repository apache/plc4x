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
from plc4py.spi.values.PlcValues import PlcBOOL, PlcINT, PlcREAL


@pytest_asyncio.fixture
async def connection() -> AsyncGenerator[PlcConnection, None]:
    driver_manager = PlcDriverManager()
    async with driver_manager.connection("umas://192.168.190.152:502") as connection:
        yield connection


@pytest.mark.asyncio
@pytest.mark.xfail
async def test_plc_driver_umas_connect(connection):
    assert connection.is_connected


@pytest.mark.asyncio
@pytest.mark.xfail
async def test_plc_driver_umas_read_boolean(connection):
    tag_alias = "Random Tag"
    tag_name = "TESTING"
    with connection.read_request_builder() as builder:
        builder.add_item(tag_alias, tag_name)
        request = builder.build()
    future = connection.execute(request)
    response = await future
    value = response.tags[tag_alias].value
    response_code = response.tags[tag_alias].response_code
    assert value == True
    assert response_code == PlcResponseCode.OK


@pytest.mark.asyncio
@pytest.mark.xfail
async def test_plc_driver_umas_read_boolean_with_data_type(connection):
    tag_alias = "Random Tag"
    tag_name = "TESTING:BOOL"
    with connection.read_request_builder() as builder:
        builder.add_item(tag_alias, tag_name)
        request = builder.build()
    future = connection.execute(request)
    response = await future
    value = response.tags[tag_alias].value
    response_code = response.tags[tag_alias].response_code
    assert value == True
    assert response_code == PlcResponseCode.OK


@pytest.mark.asyncio
@pytest.mark.xfail
async def test_plc_driver_umas_read_int(connection):
    tag_alias = "Random Tag"
    tag_name = "TESTING_INT"
    with connection.read_request_builder() as builder:
        builder.add_item(tag_alias, tag_name)
        request = builder.build()
    future = connection.execute(request)
    response = await future
    value = response.tags[tag_alias].value
    response_code = response.tags[tag_alias].response_code
    assert value == 99
    assert response_code == PlcResponseCode.OK


@pytest.mark.asyncio
@pytest.mark.xfail
async def test_plc_driver_umas_read_int_with_data_type(connection):
    tag_alias = "Random Tag"
    tag_name = "TESTING_INT:INT"
    with connection.read_request_builder() as builder:
        builder.add_item(tag_alias, tag_name)
        request = builder.build()
    future = connection.execute(request)
    response = await future
    value = response.tags[tag_alias].value
    response_code = response.tags[tag_alias].response_code
    assert value == 99
    assert response_code == PlcResponseCode.OK


@pytest.mark.asyncio
@pytest.mark.xfail
async def test_plc_driver_umas_read_dint(connection):
    tag_alias = "Random Tag"
    tag_name = "TESTING_DINT"
    with connection.read_request_builder() as builder:
        builder.add_item(tag_alias, tag_name)
        request = builder.build()
    future = connection.execute(request)
    response = await future
    value = response.tags[tag_alias].value
    response_code = response.tags[tag_alias].response_code
    assert value == 763539
    assert response_code == PlcResponseCode.OK


@pytest.mark.asyncio
@pytest.mark.xfail
async def test_plc_driver_umas_read_dint_with_data_type(connection):
    tag_alias = "Random Tag"
    tag_name = "TESTING_DINT:DINT"
    with connection.read_request_builder() as builder:
        builder.add_item(tag_alias, tag_name)
        request = builder.build()
    future = connection.execute(request)
    response = await future
    value = response.tags[tag_alias].value
    response_code = response.tags[tag_alias].response_code
    assert value == 763539
    assert response_code == PlcResponseCode.OK


@pytest.mark.asyncio
@pytest.mark.xfail
async def test_plc_driver_umas_read_ebool(connection):
    tag_alias = "Random Tag"
    tag_name = "TESTING_EBOOL"
    with connection.read_request_builder() as builder:
        builder.add_item(tag_alias, tag_name)
        request = builder.build()
    future = connection.execute(request)
    response = await future
    value = response.tags[tag_alias].value
    response_code = response.tags[tag_alias].response_code
    assert value == True
    assert response_code == PlcResponseCode.OK


@pytest.mark.asyncio
@pytest.mark.xfail
async def test_plc_driver_umas_read_ebool_with_data_type(connection):
    tag_alias = "Random Tag"
    tag_name = "TESTING_EBOOL:BOOL"
    with connection.read_request_builder() as builder:
        builder.add_item(tag_alias, tag_name)
        request = builder.build()
    future = connection.execute(request)
    response = await future
    value = response.tags[tag_alias].value
    response_code = response.tags[tag_alias].response_code
    assert value == True
    assert response_code == PlcResponseCode.OK


@pytest.mark.asyncio
@pytest.mark.xfail
async def test_plc_driver_umas_read_string(connection):
    tag_alias = "Random Tag"
    tag_name = "TESTING_STRING"
    with connection.read_request_builder() as builder:
        builder.add_item(tag_alias, tag_name)
        request = builder.build()
    future = connection.execute(request)
    response = await future
    value = response.tags[tag_alias].value
    response_code = response.tags[tag_alias].response_code
    assert value == "Hello World"
    assert response_code == PlcResponseCode.OK


@pytest.mark.asyncio
@pytest.mark.xfail
async def test_plc_driver_umas_read_string_with_data_type(connection):
    tag_alias = "Random Tag"
    tag_name = "TESTING_STRING:STRING"
    with connection.read_request_builder() as builder:
        builder.add_item(tag_alias, tag_name)
        request = builder.build()
    future = connection.execute(request)
    response = await future
    value = response.tags[tag_alias].value
    response_code = response.tags[tag_alias].response_code
    assert value == "Hello World"
    assert response_code == PlcResponseCode.OK


@pytest.mark.asyncio
@pytest.mark.xfail
async def test_plc_driver_umas_read_time(connection):
    tag_alias = "Random Tag"
    tag_name = "TESTING_TIME"
    with connection.read_request_builder() as builder:
        builder.add_item(tag_alias, tag_name)
        request = builder.build()
    future = connection.execute(request)
    response = await future
    value = response.tags[tag_alias].value
    response_code = response.tags[tag_alias].response_code
    assert value == 200000
    assert response_code == PlcResponseCode.OK


@pytest.mark.asyncio
@pytest.mark.xfail
async def test_plc_driver_umas_read_time_with_data_type(connection):
    tag_alias = "Random Tag"
    tag_name = "TESTING_TIME:TIME"
    with connection.read_request_builder() as builder:
        builder.add_item(tag_alias, tag_name)
        request = builder.build()
    future = connection.execute(request)
    response = await future
    value = response.tags[tag_alias].value
    response_code = response.tags[tag_alias].response_code
    assert value == 200000
    assert response_code == PlcResponseCode.OK


@pytest.mark.asyncio
@pytest.mark.xfail
async def test_plc_driver_umas_read_byte(connection):
    tag_alias = "Random Tag"
    tag_name = "TESTING_BYTE"
    with connection.read_request_builder() as builder:
        builder.add_item(tag_alias, tag_name)
        request = builder.build()
    future = connection.execute(request)
    response = await future
    value = response.tags[tag_alias].value
    response_code = response.tags[tag_alias].response_code
    assert value == 253
    assert response_code == PlcResponseCode.OK


@pytest.mark.asyncio
@pytest.mark.xfail
async def test_plc_driver_umas_read_byte_with_data_type(connection):
    tag_alias = "Random Tag"
    tag_name = "TESTING_BYTE:BYTE"
    with connection.read_request_builder() as builder:
        builder.add_item(tag_alias, tag_name)
        request = builder.build()
    future = connection.execute(request)
    response = await future
    value = response.tags[tag_alias].value
    response_code = response.tags[tag_alias].response_code
    assert value == 253
    assert response_code == PlcResponseCode.OK


@pytest.mark.asyncio
@pytest.mark.xfail
async def test_plc_driver_umas_read_date(connection):
    tag_alias = "Random Tag"
    tag_name = "TESTING_DATE"
    with connection.read_request_builder() as builder:
        builder.add_item(tag_alias, tag_name)
        request = builder.build()
    future = connection.execute(request)
    response = await future
    value = response.tags[tag_alias].value
    response_code = response.tags[tag_alias].response_code
    assert value == datetime.datetime(2024, 10, 25)
    assert response_code == PlcResponseCode.OK


@pytest.mark.asyncio
@pytest.mark.xfail
async def test_plc_driver_umas_read_time_with_data_type(connection):
    tag_alias = "Random Tag"
    tag_name = "TESTING_DATE:DATE"
    with connection.read_request_builder() as builder:
        builder.add_item(tag_alias, tag_name)
        request = builder.build()
    future = connection.execute(request)
    response = await future
    value = response.tags[tag_alias].value
    response_code = response.tags[tag_alias].response_code
    assert value == datetime.datetime(2024, 10, 25)
    assert response_code == PlcResponseCode.OK


@pytest.mark.asyncio
@pytest.mark.xfail
async def test_plc_driver_umas_read_dt(connection):
    tag_alias = "Random Tag"
    tag_name = "TESTING_DT"
    with connection.read_request_builder() as builder:
        builder.add_item(tag_alias, tag_name)
        request = builder.build()
    future = connection.execute(request)
    response = await future
    value = response.tags[tag_alias].value
    response_code = response.tags[tag_alias].response_code
    assert value == datetime.datetime(2000, 1, 10, 0, 40)
    assert response_code == PlcResponseCode.OK


@pytest.mark.asyncio
@pytest.mark.xfail
async def test_plc_driver_umas_read_dt_with_data_type(connection):
    tag_alias = "Random Tag"
    tag_name = "TESTING_DT:DATE_AND_TIME"
    with connection.read_request_builder() as builder:
        builder.add_item(tag_alias, tag_name)
        request = builder.build()
    future = connection.execute(request)
    response = await future
    value = response.tags[tag_alias].value
    response_code = response.tags[tag_alias].response_code
    assert value == datetime.datetime(2000, 1, 10, 0, 40)
    assert response_code == PlcResponseCode.OK


@pytest.mark.asyncio
@pytest.mark.xfail
async def test_plc_driver_umas_browse(connection):
    with connection.browse_request_builder() as builder:
        builder.add_query("All Tags", "*")
        request = builder.build()

    future = connection.execute(request)
    response = await future

    pass
