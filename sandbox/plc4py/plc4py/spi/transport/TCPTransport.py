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
from asyncio import Protocol
from dataclasses import dataclass
from typing import Callable

from plc4py.spi.transport.Plc4xBaseTransport import Plc4xBaseTransport


@dataclass
class TCPTransport(Plc4xBaseTransport):
    """
    Wrapper for the TCP Transport
    """

    host: str
    port: int

    def write(self, data):
        self._transport.write(data)

    @staticmethod
    async def create(
        protocol_factory: Callable[[], Protocol], host: str, port: int
    ) -> Plc4xBaseTransport:
        loop = asyncio.get_running_loop()
        coro = loop.create_connection(protocol_factory, host, port)
        _transport, _protocol = await coro
        return TCPTransport(_transport, _protocol, host, port)
