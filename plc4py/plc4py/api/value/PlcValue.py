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
from dataclasses import dataclass
from enum import Enum, auto
from typing import Generic, List, TypeVar

T = TypeVar("T")


@dataclass
class PlcValue(Generic[T], ABC):
    value: T

    def get_bool(self) -> bool:
        return bool(self.value)

    def get_float(self) -> float:
        return float(self.value)

    def get_str(self) -> str:
        return str(self.value)

    def get_int(self) -> int:
        return int(self.value)

    def get_list(self) -> List["PlcValue"]:
        return self.value

    def get_raw(self):
        return self.value

    def __len__(self):
        if isinstance(self.value, list):
            return len(self.value)
        return 1

    def __eq__(self, other):
        """Compare PlcValue with another or a native data type"""
        if isinstance(other, PlcValue):
            return (self.value == other.value) and (
                self.__class__.__name__ == other.__class__.__name__
            )
        else:
            return other == self.value


class PlcResponseCode(Enum):
    OK = auto()
    NOT_FOUND = auto()
    ACCESS_DENIED = auto()
    INVALID_ADDRESS = auto()
    INVALID_DATATYPE = auto()
    INVALID_DATA = auto()
    INTERNAL_ERROR = auto()
    REMOTE_BUSY = auto()
    REMOTE_ERROR = auto()
    UNSUPPORTED = auto()
    RESPONSE_PENDING = auto()
    NOT_CONNECTED = auto()
