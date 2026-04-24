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
import pytest

from plc4py.PlcDriverManager import PlcDriverManager
from plc4py.api.PlcConnection import PlcConnection
from plc4py.api.value.PlcValue import PlcResponseCode
from plc4py.drivers.simulated.SimulatedConnection import SimulatedConnection
from plc4py.spi.values.PlcValues import PlcBOOL, PlcDINT, PlcINT, PlcREAL, PlcSTRING


@pytest.mark.asyncio
async def test_driver_manager_returns_simulated_connection():
    driver_manager = PlcDriverManager()
    async with driver_manager.connection("simulated://plc-1") as connection:
        assert isinstance(connection, PlcConnection)
        assert isinstance(connection, SimulatedConnection)


@pytest.mark.asyncio
async def test_state_tag_round_trip():
    """Writing a STATE tag and then reading it back returns the written value."""
    driver_manager = PlcDriverManager()
    async with driver_manager.connection("simulated://plc-1") as connection:
        with connection.write_request_builder() as builder:
            builder.add_item("counter", "STATE/counter:INT", PlcINT(42))
            write_request = builder.build()
        write_response = await connection.execute(write_request)
        assert write_response.response_code == PlcResponseCode.OK
        assert write_response.tags["counter"].response_code == PlcResponseCode.OK

        with connection.read_request_builder() as builder:
            builder.add_item("counter", "STATE/counter:INT")
            read_request = builder.build()
        read_response = await connection.execute(read_request)
        assert read_response.response_code == PlcResponseCode.OK
        assert read_response.tags["counter"].response_code == PlcResponseCode.OK
        assert read_response.tags["counter"].value == PlcINT(42)


@pytest.mark.asyncio
async def test_state_tag_read_before_write_is_not_found():
    driver_manager = PlcDriverManager()
    async with driver_manager.connection("simulated://plc-1") as connection:
        with connection.read_request_builder() as builder:
            builder.add_item("unknown", "STATE/unknown:INT")
            request = builder.build()
        response = await connection.execute(request)
        assert response.response_code == PlcResponseCode.OK
        assert response.tags["unknown"].response_code == PlcResponseCode.NOT_FOUND


@pytest.mark.asyncio
async def test_random_tag_returns_value():
    driver_manager = PlcDriverManager()
    async with driver_manager.connection("simulated://plc-1") as connection:
        with connection.read_request_builder() as builder:
            builder.add_item("temperature", "RANDOM/temperature:REAL")
            request = builder.build()
        response = await connection.execute(request)
        assert response.response_code == PlcResponseCode.OK
        assert response.tags["temperature"].response_code == PlcResponseCode.OK
        assert isinstance(response.tags["temperature"].value, PlcREAL)


@pytest.mark.asyncio
async def test_random_tag_cannot_be_written():
    driver_manager = PlcDriverManager()
    async with driver_manager.connection("simulated://plc-1") as connection:
        with connection.write_request_builder() as builder:
            builder.add_item("rnd", "RANDOM/rnd:INT", PlcINT(1))
            request = builder.build()
        response = await connection.execute(request)
        assert response.response_code == PlcResponseCode.OK
        assert (
            response.tags["rnd"].response_code == PlcResponseCode.ACCESS_DENIED
        )


@pytest.mark.asyncio
async def test_stdout_tag_write_ok(caplog):
    driver_manager = PlcDriverManager()
    async with driver_manager.connection("simulated://plc-1") as connection:
        with connection.write_request_builder() as builder:
            builder.add_item(
                "msg", "STDOUT/msg:STRING", PlcSTRING("hello")
            )
            request = builder.build()
        response = await connection.execute(request)
        assert response.response_code == PlcResponseCode.OK
        assert response.tags["msg"].response_code == PlcResponseCode.OK


@pytest.mark.asyncio
async def test_stdout_tag_cannot_be_read():
    driver_manager = PlcDriverManager()
    async with driver_manager.connection("simulated://plc-1") as connection:
        with connection.read_request_builder() as builder:
            builder.add_item("msg", "STDOUT/msg:STRING")
            request = builder.build()
        response = await connection.execute(request)
        assert response.response_code == PlcResponseCode.OK
        assert (
            response.tags["msg"].response_code == PlcResponseCode.ACCESS_DENIED
        )


@pytest.mark.asyncio
async def test_state_tag_array_round_trip():
    driver_manager = PlcDriverManager()
    async with driver_manager.connection("simulated://plc-1") as connection:
        from plc4py.spi.values.PlcValues import PlcList

        values = PlcList([PlcDINT(i) for i in range(5)])
        with connection.write_request_builder() as builder:
            builder.add_item("samples", "STATE/samples:DINT[5]", values)
            write_request = builder.build()
        write_response = await connection.execute(write_request)
        assert write_response.tags["samples"].response_code == PlcResponseCode.OK

        with connection.read_request_builder() as builder:
            builder.add_item("samples", "STATE/samples:DINT[5]")
            read_request = builder.build()
        read_response = await connection.execute(read_request)
        assert read_response.tags["samples"].response_code == PlcResponseCode.OK
        stored = read_response.tags["samples"].value
        assert isinstance(stored, PlcList)
        assert [v.get_int() for v in stored.get_list()] == [0, 1, 2, 3, 4]


@pytest.mark.asyncio
async def test_multi_tag_read():
    """Reads with multiple tags are dispatched to each tag independently."""
    driver_manager = PlcDriverManager()
    async with driver_manager.connection("simulated://plc-1") as connection:
        with connection.write_request_builder() as builder:
            builder.add_item("flag", "STATE/flag:BOOL", PlcBOOL(True))
            await connection.execute(builder.build())

        with connection.read_request_builder() as builder:
            builder.add_item("flag", "STATE/flag:BOOL")
            builder.add_item("rnd", "RANDOM/rnd:INT")
            request = builder.build()
        response = await connection.execute(request)

        assert response.tags["flag"].response_code == PlcResponseCode.OK
        assert response.tags["flag"].value == PlcBOOL(True)
        assert response.tags["rnd"].response_code == PlcResponseCode.OK
