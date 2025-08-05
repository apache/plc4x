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
from typing import AnyStr, Pattern

from plc4py.api.exceptions.exceptions import PlcFieldParseException
from plc4py.api.messages.PlcField import PlcTag
from plc4py.protocols.umas.readwrite.UmasDataType import UmasDataType
from plc4py.spi.messages.PlcRequest import TagBuilder


class UmasTag(PlcTag):
    _ADDRESS_PATTERN: str = (
        "^(?P<tag>[%a-zA-Z_.0-9]+):?(?P<dataType>[A-Z_]*)(\[(?P<elementNb>[0-9]*)\])?"
    )

    _ADDRESS_COMPILED: Pattern[AnyStr] = re.compile(_ADDRESS_PATTERN)
    _DEFAULT_DATA_TYPE: UmasDataType = UmasDataType.INT

    def __init__(self, tag_name: str, quantity: int, data_type: UmasDataType):
        self.tag_name: str = tag_name
        self.quantity: int = quantity
        self.data_type: str = data_type

    @classmethod
    def matches(cls, address_string: str):
        return cls._ADDRESS_COMPILED.match(address_string) is not None

    @classmethod
    def _matcher(cls, address_string):
        match = cls._ADDRESS_COMPILED.match(address_string)
        if match is not None:
            return match

    @classmethod
    def create(cls, address_string):
        matcher = cls._matcher(address_string)
        tag_name: str = matcher.group("tag")
        quantity: int = (
            int(matcher.group("elementNb"))
            if "elementNb" in matcher.groupdict()
            and matcher.group("elementNb") is not None
            and len(matcher.group("elementNb")) is not 0
            else 1
        )
        data_type = (
            matcher.group("dataType")
            if "dataType" in matcher.groupdict()
            and matcher.group("dataType") is not None
            and len(matcher.group("dataType")) is not 0
            else None
        )
        return cls(tag_name, quantity, data_type)


class UmasTagBuilder(TagBuilder):
    @staticmethod
    def create(address_string: str) -> UmasTag:
        if UmasTag.matches(address_string):
            return UmasTag.create(address_string)
        else:
            raise PlcFieldParseException("Unable to parse address: " + address_string)
