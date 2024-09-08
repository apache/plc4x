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
from plc4py.protocols.modbus.readwrite.ModbusDataType import ModbusDataType
from plc4py.spi.messages.PlcRequest import TagBuilder


class ModbusTag(PlcTag):
    _ADDRESS_PATTERN: str = (
        r"(?P<address>\d+)(:(?P<datatype>[a-zA-Z_]+))?(\[(?P<quantity>\d+)\])?"
    )
    _FIXED_DIGIT_MODBUS_PATTERN: str = (
        r"(?P<address>\d{4,5})?(:(?P<datatype>[a-zA-Z_]+))?(\[(?P<quantity>\d+)\])?"
    )
    _PROTOCOL_ADDRESS_OFFSET: int = 1
    _REGISTER_MAX_ADDRESS: int = 65535

    _ADDRESS_SHORTER_PATTERN: Pattern[AnyStr] = re.compile(_FIXED_DIGIT_MODBUS_PATTERN)
    _ADDRESS_SHORT_PATTERN: Pattern[AnyStr] = re.compile(_FIXED_DIGIT_MODBUS_PATTERN)
    _DEFAULT_DATA_TYPE: ModbusDataType = ModbusDataType.INT

    _QUANTITY_MAX: int = 120

    def __init__(self, address: int, quantity: int, data_type: ModbusDataType):
        self.address: int = address
        self.quantity: int = quantity
        self.data_type: ModbusDataType = data_type

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
        address: int = (
            int(matcher.group("address")) - ModbusTag._PROTOCOL_ADDRESS_OFFSET
        )
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
            ModbusDataType[matcher.group("datatype")]
            if "datatype" in matcher.groupdict()
            and matcher.group("datatype") is not None
            else cls._DEFAULT_DATA_TYPE
        )
        return cls(address, quantity, data_type)


class ModbusTagCoil(ModbusTag):
    _ADDRESS_PREFIX: str = "0x"
    _ADDRESS_PATTERN: Pattern[AnyStr] = re.compile("coil:" + ModbusTag._ADDRESS_PATTERN)
    _ADDRESS_SHORTER_PATTERN: Pattern[AnyStr] = re.compile(
        "0" + ModbusTag._FIXED_DIGIT_MODBUS_PATTERN
    )
    _ADDRESS_SHORT_PATTERN: Pattern[AnyStr] = re.compile(
        "0x" + ModbusTag._FIXED_DIGIT_MODBUS_PATTERN
    )
    _DEFAULT_DATA_TYPE: ModbusDataType = ModbusDataType.BOOL
    _QUANTITY_MAX: int = 2000


class ModbusTagDiscreteInput(ModbusTag):
    _ADDRESS_PREFIX: str = "1x"
    _ADDRESS_PATTERN: Pattern[AnyStr] = re.compile(
        "discrete-input:" + ModbusTag._ADDRESS_PATTERN
    )
    _ADDRESS_SHORTER_PATTERN: Pattern[AnyStr] = re.compile(
        "1" + ModbusTag._FIXED_DIGIT_MODBUS_PATTERN
    )
    _ADDRESS_SHORT_PATTERN: Pattern[AnyStr] = re.compile(
        "1x" + ModbusTag._FIXED_DIGIT_MODBUS_PATTERN
    )
    _DEFAULT_DATA_TYPE: ModbusDataType = ModbusDataType.BOOL
    _QUANTITY_MAX: int = 2000


class ModbusTagInputRegister(ModbusTag):
    _ADDRESS_PREFIX: str = "3x"
    _ADDRESS_PATTERN: Pattern[AnyStr] = re.compile(
        "input-register:" + ModbusTag._ADDRESS_PATTERN
    )
    _ADDRESS_SHORTER_PATTERN: Pattern[AnyStr] = re.compile(
        "3" + ModbusTag._FIXED_DIGIT_MODBUS_PATTERN
    )
    _ADDRESS_SHORT_PATTERN: Pattern[AnyStr] = re.compile(
        "3x" + ModbusTag._FIXED_DIGIT_MODBUS_PATTERN
    )


class ModbusTagHoldingRegister(ModbusTag):
    _ADDRESS_PREFIX: str = "4x"
    _ADDRESS_PATTERN: Pattern[AnyStr] = re.compile(
        "holding-register:" + ModbusTag._ADDRESS_PATTERN
    )
    _ADDRESS_SHORTER_PATTERN: Pattern[AnyStr] = re.compile(
        "4" + ModbusTag._FIXED_DIGIT_MODBUS_PATTERN
    )
    _ADDRESS_SHORT_PATTERN: Pattern[AnyStr] = re.compile(
        "4x" + ModbusTag._FIXED_DIGIT_MODBUS_PATTERN
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
