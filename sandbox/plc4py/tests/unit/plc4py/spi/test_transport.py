#
# Licensed to the Apache Software Foundation (ASF) under one
# or more contributor license agreements.  See the NOTICE file
# distributed with this work for additional information
# regarding copyright ownership.  The ASF licenses this file
# to you under the Apache License, Version 2.0 (the
# "License"); you may not use this file except in compliance
# with the License.  You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing,
# software distributed under the License is distributed on an
# "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
# KIND, either express or implied.  See the License for the
# specific language governing permissions and limitations
# under the License.
#
import asyncio
import threading
from asyncio import Protocol
import time
import socket
from unittest.mock import MagicMock

from plc4py.spi.transport.PLC4XBaseTransport import PLC4XBaseTransport
from plc4py.spi.transport.TCPTransport import TCPTransport


async def test_base_transport_is_reading(mocker) -> None:
    """
    Unit test for the Base PLC4X Transport, is reading.
    :param mocker:
    :return:
    """

    _transport = MagicMock()
    _protocol = MagicMock()
    transport = PLC4XBaseTransport(_protocol, _transport)

    connection_mock: MagicMock = mocker.patch.object(_transport, "is_reading()")
    connection_mock.return_value = True

    assert transport.is_reading()


async def test_base_transport_write(mocker) -> None:
    """
    Unit test for the Base PLC4X Transport, write
    :param mocker:
    :return:
    """

    _transport = MagicMock()
    _protocol = MagicMock()
    transport = PLC4XBaseTransport(_protocol, _transport)

    connection_mock: MagicMock = mocker.patch.object(_transport, "write()")
    connection_mock.return_value = None

    assert transport.write(b'This is a test') is None


async def create_socket_server(host, port, future):
    with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as s:
        s.bind((host, port))
        s.listen()
        conn, addr = s.accept()
        future.set_result()
        with conn:
            print(f"Connected by {addr}")
            while True:
                data = conn.recv(1024)
                if not data:
                    break
                conn.sendall(data)


async def test_tcp_transport(mocker) -> None:
    """
    Unit test for the TCP Transport, write
    :param mocker:
    :return:
    """
    HOST = "localhost"
    PORT = 8081

    loop = asyncio.get_running_loop()
    # Create a new Future object.
    future = loop.create_future()

    threading.Thread(target=create_socket_server(HOST, PORT, future)).start()
    await future

    def get_protocol() -> Protocol:
        protocol = MagicMock()
        protocol: MagicMock = mocker.patch.object(protocol, "data_received()")
        return protocol

    transport = TCPTransport(host=HOST, port=PORT, protocol_factory=get_protocol)
    await transport.connect()
    transport.write(b'PLC4X Test Packet')

    assert transport._protocol.call_args