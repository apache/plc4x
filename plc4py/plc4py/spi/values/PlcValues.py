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
from typing import List, Any, Dict

from plc4py.api.value.PlcValue import PlcValue


@dataclass
class PlcINT(PlcValue[int]):
    pass


@dataclass
class PlcBYTE(PlcValue[int]):
    pass


@dataclass
class PlcCHAR(PlcValue[str]):
    pass


@dataclass
class PlcDATE(PlcValue[int]):
    pass


@dataclass
class PlcDATE_AND_TIME(PlcValue[int]):
    pass


@dataclass
class PlcDINT(PlcValue[int]):
    pass


@dataclass
class PlcDWORD(PlcValue[int]):
    pass


@dataclass
class PlcLDATE(PlcValue[int]):
    pass


@dataclass
class PlcLDATE_AND_TIME(PlcValue[int]):
    pass


@dataclass
class PlcLINT(PlcValue[int]):
    pass


@dataclass
class PlcList(PlcValue[List[Any]]):
    pass


@dataclass
class PlcLREAL(PlcValue[float]):
    pass


@dataclass
class PlcLTIME(PlcValue[int]):
    pass


@dataclass
class PlcLTIME_OF_DAY(PlcValue[int]):
    pass


@dataclass
class PlcLWORD(PlcValue[int]):
    pass


@dataclass
class PlcNull(PlcValue[None]):
    pass


@dataclass
class PlcRawByteArray(List[PlcValue[Any]]):
    pass


@dataclass
class PlcREAL(PlcValue[float]):
    pass


@dataclass
class PlcSINT(PlcValue[int]):
    pass


@dataclass
class PlcSTRING(PlcValue[str]):
    pass


@dataclass
class PlcStruct(PlcValue[Dict[str, PlcValue[str]]]):
    pass


@dataclass
class PlcTIME(PlcValue[int]):
    pass


@dataclass
class PlcTIME_OF_DAY(PlcValue[int]):
    pass


@dataclass
class PlcUBINT(PlcValue[int]):
    pass


@dataclass
class PlcUDINT(PlcValue[int]):
    pass


@dataclass
class PlcUINT(PlcValue[int]):
    pass


@dataclass
class PlcULINT(PlcValue[int]):
    pass


@dataclass
class PlcUSINT(PlcValue[int]):
    pass


@dataclass
class PlcWCHAR(PlcValue[str]):
    pass


@dataclass
class PlcWORD(PlcValue[int]):
    pass


@dataclass
class PlcWSTRING(PlcValue[str]):
    pass


@dataclass
class PlcBOOL(PlcValue[bool]):
    pass
