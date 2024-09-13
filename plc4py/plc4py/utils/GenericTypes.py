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
import threading
from dataclasses import dataclass
from enum import Enum, auto
from typing import Generator


class GenericGenerator(Generator):
    def __enter__(self):
        return self

    def send(self, _value):
        # TODO I need to figure out why this was added
        pass

    def throw(self):
        # TODO I need to figure out why this was added
        pass

    def __exit__(self, *args):
        # TODO I need to figure out why this was added
        pass


class ByteOrder(Enum):
    """
    Specifies the byte order for a message
    """

    LITTLE_ENDIAN = auto()
    BIG_ENDIAN = auto()
    LITTLE_ENDIAN_BYTE_SWAP = auto()
    BIG_ENDIAN_BYTE_SWAP = auto()

    def __new__(cls, value):
        obj = object.__new__(cls)
        obj._value_ = value
        return obj

    @staticmethod
    def get_short_name(order):
        if (
            order == ByteOrder.LITTLE_ENDIAN
            or order == ByteOrder.LITTLE_ENDIAN_BYTE_SWAP
        ):
            return "little"
        elif order == ByteOrder.BIG_ENDIAN or order == ByteOrder.BIG_ENDIAN_BYTE_SWAP:
            return "big"


@dataclass
class ByteOrderAware:
    byte_order: ByteOrder


class AtomicInteger:
    def __init__(self, seed=0):
        self._value = seed
        self._lock = threading.Lock()

    def increment(self, num=1):
        with self._lock:
            self._value += num
            return self._value

    def decrement(self, num=1):
        with self._lock:
            self._value -= num
            return self._value

    @property
    def value(self):
        return self._value
