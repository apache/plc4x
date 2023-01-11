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

from asyncio import Protocol, Future
from dataclasses import dataclass
from typing import Union


@dataclass
class Plc4xBaseProtocol(Protocol):
    handler: Future
    connected: bool = False

    def data_received(self, data) -> None:
        self.handler.set_result(data)

    def connection_made(self):
        self.connected = True

    def connection_lost(self, exc: Union[Exception, None]) -> None:
        self.connected = False
        raise ConnectionError
