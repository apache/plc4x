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
import logging
import random
import string
from asyncio import Transport
from dataclasses import dataclass, field
from typing import Callable, Dict, Optional

from plc4py.api.exceptions.exceptions import PlcRuntimeException
from plc4py.api.messages.PlcRequest import PlcReadRequest, PlcWriteRequest
from plc4py.api.messages.PlcResponse import PlcReadResponse, PlcWriteResponse
from plc4py.api.value.PlcValue import PlcResponseCode, PlcValue
from plc4py.drivers.simulated.SimulatedTag import SimulatedTag, SimulatedTagType
from plc4py.spi.messages.utils.ResponseItem import ResponseItem
from plc4py.spi.values.PlcValues import (
    PlcBOOL,
    PlcBYTE,
    PlcCHAR,
    PlcDINT,
    PlcDWORD,
    PlcINT,
    PlcLINT,
    PlcLREAL,
    PlcLWORD,
    PlcList,
    PlcREAL,
    PlcSINT,
    PlcSTRING,
    PlcUDINT,
    PlcUINT,
    PlcULINT,
    PlcUSINT,
    PlcWCHAR,
    PlcWORD,
)


def _random_scalar(data_type: str) -> Optional[PlcValue]:
    if data_type == "BOOL":
        return PlcBOOL(bool(random.getrandbits(1)))
    if data_type == "BYTE":
        return PlcBYTE(random.randint(0, 0xFF))
    if data_type == "WORD":
        return PlcWORD(random.randint(0, 0xFFFF))
    if data_type == "DWORD":
        return PlcDWORD(random.randint(0, 0xFFFFFFFF))
    if data_type == "LWORD":
        return PlcLWORD(random.randint(0, 0xFFFFFFFFFFFFFFFF))
    if data_type == "SINT":
        return PlcSINT(random.randint(-0x80, 0x7F))
    if data_type == "INT":
        return PlcINT(random.randint(-0x8000, 0x7FFF))
    if data_type == "DINT":
        return PlcDINT(random.randint(-0x80000000, 0x7FFFFFFF))
    if data_type == "LINT":
        return PlcLINT(
            random.randint(-0x8000000000000000, 0x7FFFFFFFFFFFFFFF)
        )
    if data_type == "USINT":
        return PlcUSINT(random.randint(0, 0xFF))
    if data_type == "UINT":
        return PlcUINT(random.randint(0, 0xFFFF))
    if data_type == "UDINT":
        return PlcUDINT(random.randint(0, 0xFFFFFFFF))
    if data_type == "ULINT":
        return PlcULINT(random.randint(0, 0xFFFFFFFFFFFFFFFF))
    if data_type == "REAL":
        return PlcREAL(random.uniform(-1.0e6, 1.0e6))
    if data_type == "LREAL":
        return PlcLREAL(random.uniform(-1.0e12, 1.0e12))
    if data_type == "CHAR":
        return PlcCHAR(random.choice(string.ascii_letters))
    if data_type == "WCHAR":
        return PlcWCHAR(random.choice(string.ascii_letters))
    if data_type in ("STRING", "WSTRING"):
        length = random.randint(1, 16)
        return PlcSTRING(
            "".join(random.choices(string.ascii_letters, k=length))
        )
    return None


def _random_value(data_type: str, quantity: int) -> Optional[PlcValue]:
    if quantity <= 1:
        return _random_scalar(data_type)
    values = [_random_scalar(data_type) for _ in range(quantity)]
    if any(v is None for v in values):
        return None
    return PlcList(values)


@dataclass
class SimulatedDevice:
    """
    Virtual PLC that keeps STATE tags in-memory, serves RANDOM tags with
    freshly generated values, and writes STDOUT tags to the log.
    """

    name: str = "simulated"
    _state: Dict[str, PlcValue] = field(default_factory=dict)

    async def read(
        self, request: PlcReadRequest, transport: Transport
    ) -> PlcReadResponse:
        if len(request.tags) == 0:
            raise PlcRuntimeException("No tags have been specified to read")
        response_items: Dict[str, ResponseItem] = {}
        for tag_name, tag in request.tags.items():
            if not isinstance(tag, SimulatedTag):
                response_items[tag_name] = ResponseItem(
                    PlcResponseCode.INVALID_ADDRESS, None
                )
                continue
            response_items[tag_name] = self._read_tag(tag)
        return PlcReadResponse(PlcResponseCode.OK, response_items)

    async def write(
        self, request: PlcWriteRequest, transport: Transport
    ) -> PlcWriteResponse:
        if len(request.tags) == 0:
            raise PlcRuntimeException("No tags have been specified to write")
        response_items: Dict[str, ResponseItem] = {}
        for tag_name, tag in request.tags.items():
            value = request.values.get(tag_name)
            if not isinstance(tag, SimulatedTag):
                response_items[tag_name] = ResponseItem(
                    PlcResponseCode.INVALID_ADDRESS, None
                )
                continue
            response_items[tag_name] = self._write_tag(tag, value)
        return PlcWriteResponse(PlcResponseCode.OK, response_items)

    def _read_tag(self, tag: SimulatedTag) -> ResponseItem:
        if tag.tag_type == SimulatedTagType.RANDOM:
            value = _random_value(tag.data_type, tag.quantity)
            if value is None:
                return ResponseItem(PlcResponseCode.INVALID_DATATYPE, None)
            return ResponseItem(PlcResponseCode.OK, value)
        if tag.tag_type == SimulatedTagType.STATE:
            value = self._state.get(tag.name)
            if value is None:
                return ResponseItem(PlcResponseCode.NOT_FOUND, None)
            return ResponseItem(PlcResponseCode.OK, value)
        # STDOUT tags cannot be read
        return ResponseItem(PlcResponseCode.ACCESS_DENIED, None)

    def _write_tag(
        self, tag: SimulatedTag, value: Optional[PlcValue]
    ) -> ResponseItem:
        if value is None:
            return ResponseItem(PlcResponseCode.INVALID_DATA, None)
        if tag.tag_type == SimulatedTagType.STATE:
            self._state[tag.name] = value
            return ResponseItem(PlcResponseCode.OK, None)
        if tag.tag_type == SimulatedTagType.STDOUT:
            logging.info(
                "TEST PLC STDOUT [%s/%s]: %s", self.name, tag.name, value.get_raw()
            )
            return ResponseItem(PlcResponseCode.OK, None)
        # RANDOM tags cannot be written
        return ResponseItem(PlcResponseCode.ACCESS_DENIED, None)
