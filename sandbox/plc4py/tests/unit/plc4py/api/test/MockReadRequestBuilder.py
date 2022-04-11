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
from asyncio import Future
from dataclasses import dataclass, field
from typing import Union, TypeVar, cast

from plc4py.api.messages.PlcMessage import PlcMessage
from plc4py.api.messages.PlcRequest import (
    ReadRequestBuilder,
    PlcField,
    PlcReadRequest,
)
from plc4py.api.messages.PlcResponse import PlcReadResponse


class MockPlcReadResponse(PlcReadResponse):
    def get_request(self) -> PlcMessage:
        return PlcMessage()


T = TypeVar("T", bound=PlcReadResponse)


class MockPlcReadRequest(PlcReadRequest[T]):
    def __init__(self, fields: list[PlcField] = []):
        super().__init__(fields)

    async def _execute(self, future: Future[T]):
        future.set_result(cast(T, MockPlcReadResponse()))

    def execute(self) -> Future[T]:
        loop = asyncio.get_running_loop()
        future = loop.create_future()
        loop.create_task(self._execute(future))
        return future


@dataclass
class MockReadRequestBuilder(ReadRequestBuilder):
    items: list[PlcField] = field(default_factory=lambda: [])

    def build(self) -> PlcReadRequest:
        return MockPlcReadRequest(self.items)

    def add_item(self, field_query: Union[str, PlcField]) -> None:
        field_temp: PlcField = (
            PlcField(field_query) if isinstance(field_query, str) else field_query
        )
        self.items.append(field_temp)
