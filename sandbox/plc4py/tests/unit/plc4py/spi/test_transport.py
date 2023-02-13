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

import threading
import asyncio

from asyncio import Transport

from typing import cast
from unittest.mock import MagicMock

import pytest

from plc4py.spi.transport.Plc4xBaseTransport import Plc4xBaseTransport
from plc4py.spi.transport.TCPTransport import TCPTransport
from tests.unit.plc4py.spi.tcp.server import Server

HOST = "localhost"
PORT = 9999


@pytest.fixture(scope="session")
def tcp_server():
    tcp_server = Server(HOST, PORT)
    with tcp_server:
        tcp_server = threading.Thread(target=tcp_server.listen_for_traffic)
        tcp_server.daemon = True
        tcp_server.start()
        yield tcp_server


async def test_base_transport_is_reading(mocker) -> None:
    """
    Unit test for the Base PLC4X Transport, is reading.
    :param mocker:
    :return:
    """

    _transport: MagicMock = mocker.patch.object(Transport, "is_reading")
    _transport.is_reading.return_value = True

    _protocol = MagicMock()
    transport = Plc4xBaseTransport(_transport, _protocol)

    assert transport.is_reading()


async def test_tcp_transport(mocker, tcp_server) -> None:
    """
    Unit test for the TCP Transport, write
    :param mocker:
    :return:
    """
    message = b"PLC4X Test Packet"
    loop = asyncio.get_running_loop()
    future = loop.create_future()

    def get_protocol(future) -> asyncio.Protocol:
        protocol: MagicMock = mocker.patch.object(
            asyncio.Protocol, attribute="data_received"
        )
        protocol.attach_mock(protocol, attribute="data_received")
        protocol.data_received.side_effect = future.set_result
        return protocol

    transport = await TCPTransport.create(
        protocol_factory=lambda: get_protocol(future),
        host=HOST,
        port=PORT,
    )
    transport.write(message)
    await future

    cast(MagicMock, transport._protocol).assert_called_with(message)
