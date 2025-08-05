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
from abc import ABC
from dataclasses import dataclass, field

from typing import Generic, TypeVar, Union, Dict, List

from plc4py.api.messages.PlcField import PlcTag
from plc4py.api.messages.PlcResponse import PlcResponseCode
from plc4py.api.value.PlcValue import PlcValue

T = TypeVar("T", bound=Union[PlcValue, None])


@dataclass
class ResponseItem(Generic[T], ABC):
    response_code: PlcResponseCode
    value: T


@dataclass
class ArrayInfo:
    size: int
    lower_bound: int
    upper_bound: int


R = TypeVar("R", bound=Union[PlcTag, None])


@dataclass
class PlcBrowseItem(Generic[R], ABC):
    tag: R
    name: str
    readable: bool
    writeable: bool
    subscribable: bool
    publishable: bool
    array: bool
    array_info: List[ArrayInfo] = field(default_factory=lambda: [])
    children: Dict[str, "PlcBrowseItem"] = field(default_factory=lambda: {})
    options: Dict[str, PlcValue] = field(default_factory=lambda: {})
