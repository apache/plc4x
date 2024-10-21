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
import logging
import time

import pytest

from plc4py.api.value.PlcValue import PlcResponseCode
from plc4py.PlcDriverManager import PlcDriverManager
from plc4py.spi.values.PlcValues import PlcBOOL, PlcINT, PlcREAL


@pytest.mark.asyncio
@pytest.mark.xfail
async def manual_test_plc_driver_umas_connect():
    driver_manager = PlcDriverManager()
    async with driver_manager.connection("umas://127.0.0.1:5555") as connection:
        assert connection.is_connected()
    assert not connection.is_connected()


@pytest.mark.asyncio
@pytest.mark.xfail
async def test_plc_driver_umas_read():
    log = logging.getLogger(__name__)

    driver_manager = PlcDriverManager()
    async with driver_manager.connection("umas://192.168.190.152:502") as connection:
        with connection.read_request_builder() as builder:
            #builder.add_item(f"Random Tag {1}", "TESTING_10:BOOL")
            builder.add_item(f"Random Tag {2}", "TESTING_REAL:REAL")

            request = builder.build()

        future = connection.execute(request)
        response = await future
        value = response.tags["Random Tag 1"].value
        assert value == 0.0


@pytest.mark.asyncio
@pytest.mark.xfail
async def test_plc_driver_umas_write():
    log = logging.getLogger(__name__)

    driver_manager = PlcDriverManager()
    async with driver_manager.connection("umas://192.168.190.152:502") as connection:
        with connection.write_request_builder() as builder:
            # builder.add_item(f"Random Tag {1}", "TESTING_10:BOOL", PlcBOOL(True))
            # builder.add_item(f"Random Tag {2}", "TESTING_INT:INT", PlcINT(10))
            # builder.add_item(f"Random Tag {3}", "TESTING_EBOOL:BOOL", PlcBOOL(True))
            builder.add_item(f"Random Tag {4}", "TESTING_REAL:REAL", PlcREAL(3.18))
            request = builder.build()

        future = connection.execute(request)
        response = await future
        value = response.tags["Random Tag 1"].response_code
        assert value == PlcResponseCode.OK


@pytest.mark.asyncio
@pytest.mark.xfail
async def test_plc_driver_umas_browse():
    driver_manager = PlcDriverManager()
    async with driver_manager.connection("umas://192.168.190.174:502") as connection:
        with connection.browse_request_builder() as builder:
            builder.add_query("All Tags", "*")
            request = builder.build()

        future = connection.execute(request)
        response = await future

        pass
