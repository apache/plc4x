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
from re import Pattern
from typing import AnyStr

from plc4py.api.exceptions.exceptions import PlcFieldParseException
from plc4py.api.messages.PlcField import PlcTag
from plc4py.protocols.modbus.readwrite.ModbusDataType import ModbusDataType
from plc4py.spi.messages.PlcRequest import TagBuilder


@dataclass
class ModbusTag(PlcTag):
    address: int
    quantity: int
    data_type: ModbusDataType

    ADDRESS_PATTERN: str = (
        "(?P<address>\d+)(:(?P<datatype>[a-zA-Z_]+))?([(?P<quantity>\d+)])?"
    )
    FIXED_DIGIT_MODBUS_PATTERN: str = (
        "(?P<address>\d{4,5})?(:(?P<datatype>[a-zA-Z_]+))?([(?P<quantity>\d+)])?"
    )
    PROTOCOL_ADDRESS_OFFSET: int = 1
    REGISTER_MAX_ADDRESS: int = 65535

    ADDRESS_SHORTER_PATTERN: Pattern[AnyStr] = re.compile(FIXED_DIGIT_MODBUS_PATTERN)
    ADDRESS_SHORT_PATTERN: Pattern[AnyStr] = re.compile(FIXED_DIGIT_MODBUS_PATTERN)
    DEFAULT_DATA_TYPE: ModbusDataType = ModbusDataType.INT

    QUANTITY_MAX: int = 120

    @classmethod
    def matches(cls, address_string: str):
        return (
            cls.ADDRESS_PATTERN.match(address_string) is not None
            or cls.ADDRESS_SHORTER_PATTERN.match(address_string) is not None
            or cls.ADDRESS_SHORT_PATTERN.match(address_string) is not None
        )

    @classmethod
    def _matcher(cls, address_string):
        match = cls.ADDRESS_PATTERN.match(address_string)
        if match is not None:
            return match
        match = cls.ADDRESS_SHORTER_PATTERN.match(address_string)
        if match is not None:
            return match
        match = cls.ADDRESS_SHORT_PATTERN.match(address_string)
        if match is not None:
            return match

    @classmethod
    def create(cls, address_string):
        matcher = cls._matcher(address_string)
        address: int = int(matcher.group("address")) - ModbusTag.PROTOCOL_ADDRESS_OFFSET
        if address > cls.REGISTER_MAX_ADDRESS:
            raise PlcFieldParseException(
                "Address must be less than or equal to "
                + str(cls.REGISTER_MAX_ADDRESS)
                + ". Was "
                + str(address + cls.PROTOCOL_ADDRESS_OFFSET)
            )

        quantity: int = (
            int(matcher.group("quantity"))
            if "quantity" in matcher.groupdict()
            and matcher.group("datatype") is not None
            else 1
        )
        if (address + quantity) > cls.REGISTER_MAX_ADDRESS:
            raise PlcFieldParseException(
                "Last requested address is out of range, should be between "
                + str(cls.PROTOCOL_ADDRESS_OFFSET)
                + " and "
                + str(cls.REGISTER_MAX_ADDRESS)
                + ". Was "
                + str(address + cls.PROTOCOL_ADDRESS_OFFSET + (quantity - 1))
            )

        if quantity > cls.QUANTITY_MAX:
            raise PlcFieldParseException(
                "quantity may not be larger than 2000. Was " + str(quantity)
            )

        data_type = (
            ModbusDataType(matcher.group("datatype"))
            if "datatype" in matcher.groupdict()
            and matcher.group("datatype") is not None
            else ModbusDataType.INT
        )
        return cls(address, quantity, data_type)


class ModbusTagCoil(ModbusTag):
    ADDRESS_PREFIX: str = "0x"
    ADDRESS_PATTERN: Pattern[AnyStr] = re.compile("coil:" + ModbusTag.ADDRESS_PATTERN)
    ADDRESS_SHORTER_PATTERN: Pattern[AnyStr] = re.compile(
        "0" + ModbusTag.FIXED_DIGIT_MODBUS_PATTERN
    )
    ADDRESS_SHORT_PATTERN: Pattern[AnyStr] = re.compile(
        "0x" + ModbusTag.FIXED_DIGIT_MODBUS_PATTERN
    )
    DEFAULT_DATA_TYPE: ModbusDataType = ModbusDataType.BOOL


class ModbusTagDiscreteInput(ModbusTag):
    ADDRESS_PREFIX: str = "1x"
    ADDRESS_PATTERN: Pattern[AnyStr] = re.compile(
        "discrete-input:" + ModbusTag.ADDRESS_PATTERN
    )
    ADDRESS_SHORTER_PATTERN: Pattern[AnyStr] = re.compile(
        "1" + ModbusTag.FIXED_DIGIT_MODBUS_PATTERN
    )
    ADDRESS_SHORT_PATTERN: Pattern[AnyStr] = re.compile(
        "1x" + ModbusTag.FIXED_DIGIT_MODBUS_PATTERN
    )
    DEFAULT_DATA_TYPE: ModbusDataType = ModbusDataType.BOOL


class ModbusTagInputRegister(ModbusTag):
    ADDRESS_PREFIX: str = "3x"
    ADDRESS_PATTERN: Pattern[AnyStr] = re.compile(
        "input-register:" + ModbusTag.ADDRESS_PATTERN
    )
    ADDRESS_SHORTER_PATTERN: Pattern[AnyStr] = re.compile(
        "3" + ModbusTag.FIXED_DIGIT_MODBUS_PATTERN
    )
    ADDRESS_SHORT_PATTERN: Pattern[AnyStr] = re.compile(
        "3x" + ModbusTag.FIXED_DIGIT_MODBUS_PATTERN
    )


class ModbusTagHoldingRegister(ModbusTag):
    ADDRESS_PREFIX: str = "4x"
    ADDRESS_PATTERN: Pattern[AnyStr] = re.compile(
        "holding-register:" + ModbusTag.ADDRESS_PATTERN
    )
    ADDRESS_SHORTER_PATTERN: Pattern[AnyStr] = re.compile(
        "4" + ModbusTag.FIXED_DIGIT_MODBUS_PATTERN
    )
    ADDRESS_SHORT_PATTERN: Pattern[AnyStr] = re.compile(
        "4x" + ModbusTag.FIXED_DIGIT_MODBUS_PATTERN
    )


class ModbusTagBuilder(TagBuilder):
    @staticmethod
    def create(address_string: str) -> ModbusTag:
        if ModbusTagCoil.matches(address_string):
            return ModbusTagCoil.create(address_string)
        elif ModbusTagDiscreteInput.matches(address_string):
            return ModbusTagDiscreteInput.create(address_string)
        elif ModbusTagInputRegister.matches(address_string):
            return ModbusTagInputRegister.create(address_string)
        elif ModbusTagHoldingRegister.matches(address_string):
            return ModbusTagHoldingRegister.create(address_string)
        else:
            raise PlcFieldParseException("Unable to parse address: " + address_string)
