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
import re
from dataclasses import dataclass
from enum import Enum
from typing import ClassVar, Pattern

from plc4py.api.exceptions.exceptions import PlcFieldParseException
from plc4py.api.messages.PlcField import PlcTag
from plc4py.spi.messages.PlcRequest import TagBuilder


class SimulatedTagType(Enum):
    RANDOM = "RANDOM"
    STATE = "STATE"
    STDOUT = "STDOUT"


@dataclass
class SimulatedTag(PlcTag):
    _ADDRESS_PATTERN: ClassVar[Pattern] = re.compile(
        r"^(?P<type>\w+)/(?P<name>[a-zA-Z0-9_.]+):(?P<dataType>[a-zA-Z0-9]+)(\[(?P<numElements>\d+)\])?$"
    )

    tag_type: SimulatedTagType
    name: str
    data_type: str
    quantity: int = 1

    @classmethod
    def matches(cls, address_string: str) -> bool:
        return cls._ADDRESS_PATTERN.match(address_string) is not None

    @classmethod
    def create(cls, address_string: str) -> "SimulatedTag":
        matcher = cls._ADDRESS_PATTERN.match(address_string)
        if matcher is None:
            raise PlcFieldParseException(
                f"Unable to parse address: {address_string}"
            )
        try:
            tag_type = SimulatedTagType[matcher.group("type").upper()]
        except KeyError:
            raise PlcFieldParseException(
                f"Invalid simulated tag type: {matcher.group('type')}"
            )
        name = matcher.group("name")
        data_type = matcher.group("dataType").upper()
        num_elements = matcher.group("numElements")
        quantity = int(num_elements) if num_elements is not None else 1
        return cls(
            address=address_string,
            tag_type=tag_type,
            name=name,
            data_type=data_type,
            quantity=quantity,
        )


class SimulatedTagBuilder(TagBuilder):
    @staticmethod
    def create(address_string: str) -> SimulatedTag:
        return SimulatedTag.create(address_string)
