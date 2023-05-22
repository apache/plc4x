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
from enum import Enum, auto
from typing import Generator


# TODO: Figure out what the parameters are and if we need this
class GenericGenerator(Generator):
    def __enter__(self):
        return self

    def send(self, _value, blah):
        pass

    def throw(self):
        pass

    def __exit__(self, *args):
        pass


class ByteOrder(Enum):
    """
    Specifies the byte order for a message
    """

    LITTLE_ENDIAN = auto()
    BIG_ENDIAN = auto()

    def __new__(cls, value):
        obj = object.__new__(cls)
        obj._value_ = value
        return obj

    @staticmethod
    def get_short_name(order):
        if order == ByteOrder.LITTLE_ENDIAN:
            return "little"
        elif order == ByteOrder.BIG_ENDIAN:
            return "big"


@dataclass
class ByteOrderAware:
    byte_order: ByteOrder
