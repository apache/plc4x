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
        "(?P<address>\d+)(:(?P<datatype>[a-zA-Z_]+))?(\[(?P<quantity>\d+)\])?"
    )
    _FIXED_DIGIT_Umas_PATTERN: str = (
        "(?P<address>\d{4,5})?(:(?P<datatype>[a-zA-Z_]+))?(\[(?P<quantity>\d+)\])?"
    )
    _PROTOCOL_ADDRESS_OFFSET: int = 1
    _REGISTER_MAX_ADDRESS: int = 65535

    _ADDRESS_SHORTER_PATTERN: Pattern[AnyStr] = re.compile(_FIXED_DIGIT_Umas_PATTERN)
    _ADDRESS_SHORT_PATTERN: Pattern[AnyStr] = re.compile(_FIXED_DIGIT_Umas_PATTERN)
    _DEFAULT_DATA_TYPE: UmasDataType = UmasDataType.INT

    _QUANTITY_MAX: int = 120

    def __init__(self, address: int, quantity: int, data_type: UmasDataType):
        self.address: int = address
        self.quantity: int = quantity
        self.data_type: UmasDataType = data_type

    @classmethod
    def matches(cls, address_string: str):
        return (
            cls._ADDRESS_PATTERN.match(address_string) is not None
            or cls._ADDRESS_SHORTER_PATTERN.match(address_string) is not None
            or cls._ADDRESS_SHORT_PATTERN.match(address_string) is not None
        )

    @classmethod
    def _matcher(cls, address_string):
        match = cls._ADDRESS_PATTERN.match(address_string)
        if match is not None:
            return match
        match = cls._ADDRESS_SHORT_PATTERN.match(address_string)
        if match is not None:
            return match
        match = cls._ADDRESS_SHORTER_PATTERN.match(address_string)
        if match is not None:
            return match

    @classmethod
    def create(cls, address_string):
        matcher = cls._matcher(address_string)
        address: int = int(matcher.group("address")) - UmasTag._PROTOCOL_ADDRESS_OFFSET
        if address > cls._REGISTER_MAX_ADDRESS:
            raise PlcFieldParseException(
                "Address must be less than or equal to "
                + str(cls._REGISTER_MAX_ADDRESS)
                + ". Was "
                + str(address + cls._PROTOCOL_ADDRESS_OFFSET)
            )

        quantity: int = (
            int(matcher.group("quantity"))
            if "quantity" in matcher.groupdict()
            and matcher.group("quantity") is not None
            else 1
        )
        if (address + quantity) > cls._REGISTER_MAX_ADDRESS:
            raise PlcFieldParseException(
                "Last requested address is out of range, should be between "
                + str(cls._PROTOCOL_ADDRESS_OFFSET)
                + " and "
                + str(cls._REGISTER_MAX_ADDRESS)
                + ". Was "
                + str(address + cls._PROTOCOL_ADDRESS_OFFSET + (quantity - 1))
            )

        if quantity > cls._QUANTITY_MAX:
            raise PlcFieldParseException(
                "quantity may not be larger than 2000. Was " + str(quantity)
            )

        data_type = (
            UmasDataType(matcher.group("datatype"))
            if "datatype" in matcher.groupdict()
            and matcher.group("datatype") is not None
            else cls._DEFAULT_DATA_TYPE
        )
        return cls(address, quantity, data_type)


class UmasTagCoil(UmasTag):
    _ADDRESS_PREFIX: str = "0x"
    _ADDRESS_PATTERN: Pattern[AnyStr] = re.compile("coil:" + UmasTag._ADDRESS_PATTERN)
    _ADDRESS_SHORTER_PATTERN: Pattern[AnyStr] = re.compile(
        "0" + UmasTag._FIXED_DIGIT_Umas_PATTERN
    )
    _ADDRESS_SHORT_PATTERN: Pattern[AnyStr] = re.compile(
        "0x" + UmasTag._FIXED_DIGIT_Umas_PATTERN
    )
    _DEFAULT_DATA_TYPE: UmasDataType = UmasDataType.BOOL
    _QUANTITY_MAX: int = 2000


class UmasTagDiscreteInput(UmasTag):
    _ADDRESS_PREFIX: str = "1x"
    _ADDRESS_PATTERN: Pattern[AnyStr] = re.compile(
        "discrete-input:" + UmasTag._ADDRESS_PATTERN
    )
    _ADDRESS_SHORTER_PATTERN: Pattern[AnyStr] = re.compile(
        "1" + UmasTag._FIXED_DIGIT_Umas_PATTERN
    )
    _ADDRESS_SHORT_PATTERN: Pattern[AnyStr] = re.compile(
        "1x" + UmasTag._FIXED_DIGIT_Umas_PATTERN
    )
    _DEFAULT_DATA_TYPE: UmasDataType = UmasDataType.BOOL
    _QUANTITY_MAX: int = 2000


class UmasTagInputRegister(UmasTag):
    _ADDRESS_PREFIX: str = "3x"
    _ADDRESS_PATTERN: Pattern[AnyStr] = re.compile(
        "input-register:" + UmasTag._ADDRESS_PATTERN
    )
    _ADDRESS_SHORTER_PATTERN: Pattern[AnyStr] = re.compile(
        "3" + UmasTag._FIXED_DIGIT_Umas_PATTERN
    )
    _ADDRESS_SHORT_PATTERN: Pattern[AnyStr] = re.compile(
        "3x" + UmasTag._FIXED_DIGIT_Umas_PATTERN
    )


class UmasTagHoldingRegister(UmasTag):
    _ADDRESS_PREFIX: str = "4x"
    _ADDRESS_PATTERN: Pattern[AnyStr] = re.compile(
        "holding-register:" + UmasTag._ADDRESS_PATTERN
    )
    _ADDRESS_SHORTER_PATTERN: Pattern[AnyStr] = re.compile(
        "4" + UmasTag._FIXED_DIGIT_Umas_PATTERN
    )
    _ADDRESS_SHORT_PATTERN: Pattern[AnyStr] = re.compile(
        "4x" + UmasTag._FIXED_DIGIT_Umas_PATTERN
    )


class UmasTagBuilder(TagBuilder):
    @staticmethod
    def create(address_string: str) -> UmasTag:
        if UmasTagCoil.matches(address_string):
            return UmasTagCoil.create(address_string)
        elif UmasTagDiscreteInput.matches(address_string):
            return UmasTagDiscreteInput.create(address_string)
        elif UmasTagInputRegister.matches(address_string):
            return UmasTagInputRegister.create(address_string)
        elif UmasTagHoldingRegister.matches(address_string):
            return UmasTagHoldingRegister.create(address_string)
        else:
            raise PlcFieldParseException("Unable to parse address: " + address_string)
