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
from asyncio import Protocol
from dataclasses import dataclass
from typing import Callable, Type

from plc4py.spi.transport.Plc4xBaseTransport import Plc4xBaseTransport
import plc4py
from plc4py.spi.transport.PlcTransportLoader import PlcTransportLoader

from plc4py.spi.transport.TCPTransport import TCPTransport


@dataclass
class MockTransport(Plc4xBaseTransport):
    """
    Wrapper for the Mock Transport
    """

    host: str
    port: int

    def write(self, data):
        logging.debug("Writing data to Mock Transport")

    def is_closing(self):
        return False

    @staticmethod
    async def create(
        protocol_factory: Callable[[], Protocol], host: str, port: int
    ) -> Plc4xBaseTransport:
        return MockTransport(None, None, host, port)


class MockTransportLoader(PlcTransportLoader):
    """
    Umas Driver Pluggy Hook Implementation, lets pluggy find the driver by name
    """

    @staticmethod
    @plc4py.spi.transport.hookimpl
    def get_transport() -> Type[MockTransport]:
        return MockTransport

    @staticmethod
    @plc4py.spi.transport.hookimpl
    def key() -> str:
        return "mock"
