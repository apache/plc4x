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
from asyncio import Transport
from dataclasses import dataclass, field
from typing import Dict

from plc4py.api.messages.PlcRequest import PlcReadRequest
from plc4py.api.messages.PlcResponse import PlcReadResponse
from plc4py.api.value.PlcValue import PlcValue
from plc4py.drivers.umas.UmasConfiguration import UmasConfiguration
from plc4py.utils.GenericTypes import AtomicInteger


@dataclass
class UmasDevice:
    _configuration: UmasConfiguration
    tags: Dict[str, PlcValue] = field(default_factory=lambda: {})
    _transaction_generator: AtomicInteger = field(
        default_factory=lambda: AtomicInteger()
    )

    async def connect(self):
        pass

    async def read(
        self, request: PlcReadRequest, transport: Transport
    ) -> PlcReadResponse:
        """
        Reads one field from the Umas Device
        """
        pass
