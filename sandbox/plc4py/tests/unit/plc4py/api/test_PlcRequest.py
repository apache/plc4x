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

from plc4py.api.PlcConnection import PlcConnection
from plc4py.api.messages.PlcField import PlcField
from plc4py.api.messages.PlcRequest import (
    PlcFieldRequest,
)
from plc4py.api.messages.PlcResponse import PlcReadResponse
from plc4py.api.value.PlcValue import PlcResponseCode
from plc4py.spi.messages.utils.ResponseItem import ResponseItem
from plc4py.spi.values.PlcBOOL import PlcBOOL
from plc4py.spi.values.PlcINT import PlcINT
from plc4py.drivers.mock.MockConnection import MockConnection


def test_read_request_builder_empty_request(mocker) -> None:
    """
    Create an empty Plc Request and confirm that it gets created without any Fields
    :param mocker:
    :return:
    """
    connection: PlcConnection = MockConnection()

    # the connection function is supposed to support context manager
    # so using it in a with statement should result in close being called on the connection
    with connection.read_request_builder() as builder:
        request: PlcFieldRequest = builder.build()
    assert len(request.field_names) == 0


def test_read_request_builder_non_empty_request(mocker) -> None:
    """
    Add a field to the reuqest and confirm that field gets added
    :param mocker:
    :return:
    """
    connection: PlcConnection = MockConnection()

    # the connection function is supposed to support context manager
    # so using it in a with statement should result in close being called on the connection
    with connection.read_request_builder() as builder:
        builder.add_item("1:BOOL")
        request: PlcFieldRequest = builder.build()

    # verify that request has one field
    assert request.field_names == ["1:BOOL"]
    assert len(request.field_names) == 1


@pytest.mark.asyncio
async def test_read_request_builder_non_empty_request_not_connected(mocker) -> None:
    """
    Create a request with a field and then confirm an empty response gets returned with a NOT_CONNECTED code
    :param mocker:
    :return:
    """
    connection: PlcConnection = MockConnection()

    # the connection function is supposed to support context manager
    # so using it in a with statement should result in close being called on the connection
    with connection.read_request_builder() as builder:
        builder.add_item("1:BOOL")
        request: PlcFieldRequest = builder.build()
        response = await connection.execute(request)

    # verify that request has one field
    assert response.code == PlcResponseCode.NOT_CONNECTED


@pytest.mark.asyncio
async def test_read_request_builder_non_empty_request_connected_bool(mocker) -> None:
    """
    Create a request with a field and then confirm a non empty response gets returned with a OK code
    :param mocker:
    :return:
    """
    connection: PlcConnection = await MockConnection.create("mock://localhost")
    field = "1:BOOL"

    # the connection function is supposed to support context manager
    # so using it in a with statement should result in close being called on the connection
    with connection.read_request_builder() as builder:
        builder.add_item(field)
        request: PlcFieldRequest = builder.build()
        response: PlcReadResponse = cast(
            PlcReadResponse, await connection.execute(request)
        )

    # verify that request has one field
    assert response.code == PlcResponseCode.OK

    value = response.values[field][0].value
    assert not value.get_bool()


@pytest.mark.asyncio
async def test_read_request_builder_non_empty_request_connected_int(mocker) -> None:
    """
    Create a request with a field and then confirm a non empty response gets returned with a OK code
    :param mocker:
    :return:
    """
    connection: PlcConnection = await MockConnection.create("mock://localhost")
    field = "1:INT"

    # the connection function is supposed to support context manager
    # so using it in a with statement should result in close being called on the connection
    with connection.read_request_builder() as builder:
        builder.add_item(field)
        request: PlcFieldRequest = builder.build()
        response: PlcReadResponse = cast(
            PlcReadResponse, await connection.execute(request)
        )

    # verify that request has one field
    assert response.code == PlcResponseCode.OK

    value = response.values[field][0].value
    assert value.get_int() == 0


def test_read_response_boolean_response(mocker) -> None:
    """
    Create a Plc Response with a boolean field, confirm that a boolean gets returned
    :param mocker:
    :return:
    """
    response = PlcReadResponse(
        PlcResponseCode.OK,
        [PlcField("1:BOOL")],
        {"1:BOOL": [ResponseItem(PlcResponseCode.OK, PlcBOOL(True))]},
    )
    assert response.get_boolean("1:BOOL")
    assert isinstance(response.get_plc_value("1:BOOL"), PlcBOOL)


def test_read_response_int_response(mocker) -> None:
    """
    Create a Plc Response with an int field, confirm that an int gets returned
    :param mocker:
    :return:
    """
    response = PlcReadResponse(
        PlcResponseCode.OK,
        [PlcField("1:INT")],
        {"1:INT": [ResponseItem(PlcResponseCode.OK, PlcINT(10))]},
    )
    assert response.get_int("1:INT") == 10
    assert isinstance(response.get_plc_value("1:INT"), PlcINT)
