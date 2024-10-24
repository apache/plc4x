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
from dataclasses import dataclass
from datetime import datetime
from typing import Any, Dict, List

from plc4py.api.value.PlcValue import PlcValue


class PlcINT(PlcValue[int]):
    pass


class PlcBYTE(PlcValue[int]):
    pass


class PlcCHAR(PlcValue[str]):
    pass


class PlcDATE(PlcValue[datetime]):
    pass


class PlcDATE_AND_TIME(PlcValue[datetime]):
    pass


class PlcDINT(PlcValue[int]):
    pass


class PlcDWORD(PlcValue[int]):
    pass


class PlcLDATE(PlcValue[int]):
    pass


class PlcLDATE_AND_TIME(PlcValue[int]):
    pass


class PlcLINT(PlcValue[int]):
    pass


class PlcList(PlcValue[List[Any]]):
    pass


class PlcLREAL(PlcValue[float]):
    pass


class PlcLTIME(PlcValue[int]):
    pass


class PlcLTIME_OF_DAY(PlcValue[int]):
    pass


class PlcLWORD(PlcValue[int]):
    pass


class PlcNull(PlcValue[None]):
    pass


class PlcRawByteArray(List[PlcValue[Any]]):
    pass


class PlcREAL(PlcValue[float]):
    pass


class PlcSINT(PlcValue[int]):
    pass


class PlcSTRING(PlcValue[str]):
    pass


class PlcStruct(PlcValue[Dict[str, PlcValue[str]]]):
    pass


class PlcTIME(PlcValue[int]):
    pass


class PlcTIME_OF_DAY(PlcValue[int]):
    pass


class PlcUBINT(PlcValue[int]):
    pass


class PlcUDINT(PlcValue[int]):
    pass


class PlcUINT(PlcValue[int]):
    pass


class PlcULINT(PlcValue[int]):
    pass


class PlcUSINT(PlcValue[int]):
    pass


class PlcWCHAR(PlcValue[str]):
    pass


class PlcWORD(PlcValue[int]):
    pass


class PlcWSTRING(PlcValue[str]):
    pass


class PlcBOOL(PlcValue[bool]):
    pass
