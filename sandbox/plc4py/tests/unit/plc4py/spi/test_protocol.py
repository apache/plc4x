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

import pytest


from plc4py.spi.Plc4xBaseProtocol import Plc4xBaseProtocol
from plc4py.spi.transport.TCPTransport import TCPTransport
from tests.unit.plc4py.spi.tcp.server import Server

HOST = "localhost"
PORT = 9998


@pytest.fixture(scope="session")
def tcp_server():
    tcp_server = Server(HOST, PORT)
    with tcp_server:
        tcp_server = threading.Thread(target=tcp_server.listen_for_traffic)
        tcp_server.daemon = True
        tcp_server.start()
        yield tcp_server


async def test_tcp_protocol(mocker, tcp_server) -> None:
    """
    Unit test for a not implemented protocol
    :param mocker:
    :return:
    """
    message = b"PLC4X Test Packet"
    loop = asyncio.get_running_loop()
    future = loop.create_future()

    def get_protocol(future) -> asyncio.Protocol:
        protocol = Plc4xBaseProtocol(future)
        return protocol

    transport = await TCPTransport.create(
        protocol_factory=lambda: get_protocol(future), host=HOST, port=PORT
    )
    transport.write(message)
    result = await future

    assert result == message
