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

from abc import ABC
from abc import abstractmethod
from plc4py.api.exceptions.exceptions import ParseException
from plc4py.api.exceptions.exceptions import PlcRuntimeException
from plc4py.api.exceptions.exceptions import SerializationException
from plc4py.api.messages.PlcMessage import PlcMessage
from plc4py.protocols.modbus.readwrite.DriverType import DriverType
from plc4py.spi.generation.ReadBuffer import ReadBuffer
from plc4py.spi.generation.WriteBuffer import WriteBuffer
from plc4py.utils.ConnectionStringHandling import strtobool
import math


@dataclass
class ModbusADU(ABC, PlcMessage):
    # Arguments.
    response: bool

    # Abstract accessors for discriminator values.
    @property
    def driver_type(self) -> DriverType:
        pass

    @abstractmethod
    def serialize_modbus_adu_child(self, write_buffer: WriteBuffer) -> None:
        pass

    def serialize(self, write_buffer: WriteBuffer):
        write_buffer.push_context("ModbusADU")

        # Switch field (Serialize the sub-type)
        self.serialize_modbus_adu_child(write_buffer)

        write_buffer.pop_context("ModbusADU")

    def length_in_bytes(self) -> int:
        return int(math.ceil(float(self.length_in_bits() / 8.0)))

    def length_in_bits(self) -> int:
        length_in_bits: int = 0
        _value: ModbusADU = self

        # Length of subtype elements will be added by sub-type...

        return length_in_bits

    @staticmethod
    def static_parse(read_buffer: ReadBuffer, **kwargs):

        if kwargs is None:
            raise PlcRuntimeException(
                "Wrong number of arguments, expected 2, but got None"
            )

        driver_type: DriverType = 0
        if isinstance(kwargs.get("driver_type"), DriverType):
            driver_type = DriverType(kwargs.get("driver_type"))
        elif isinstance(kwargs.get("driver_type"), str):
            driver_type = DriverType(str(kwargs.get("driver_type")))
        else:
            raise PlcRuntimeException(
                "Argument 0 expected to be of type DriverType or a string which is parseable but was "
                + kwargs.get("driver_type").getClass().getName()
            )

        response: bool = False
        if isinstance(kwargs.get("response"), bool):
            response = bool(kwargs.get("response"))
        elif isinstance(kwargs.get("response"), str):
            response = bool(str(kwargs.get("response")))
        else:
            raise PlcRuntimeException(
                "Argument 1 expected to be of type bool or a string which is parseable but was "
                + kwargs.get("response").getClass().getName()
            )

        return ModbusADU.static_parse_context(read_buffer, driver_type, response)

    @staticmethod
    def static_parse_context(
        read_buffer: ReadBuffer, driver_type: DriverType, response: bool
    ):
        read_buffer.push_context("ModbusADU")

        if isinstance(driver_type, str):
            driver_type = DriverType[driver_type]
        if isinstance(response, str):
            response = bool(strtobool(response))

        # Switch Field (Depending on the discriminator values, passes the instantiation to a sub-type)
        builder: ModbusADUBuilder = None
        from plc4py.protocols.modbus.readwrite.ModbusTcpADU import ModbusTcpADU

        if driver_type == DriverType.MODBUS_TCP:

            builder = ModbusTcpADU.static_parse_builder(
                read_buffer, driver_type, response
            )
        from plc4py.protocols.modbus.readwrite.ModbusRtuADU import ModbusRtuADU

        if driver_type == DriverType.MODBUS_RTU:

            builder = ModbusRtuADU.static_parse_builder(
                read_buffer, driver_type, response
            )
        from plc4py.protocols.modbus.readwrite.ModbusAsciiADU import ModbusAsciiADU

        if driver_type == DriverType.MODBUS_ASCII:

            builder = ModbusAsciiADU.static_parse_builder(
                read_buffer, driver_type, response
            )
        if builder is None:
            raise ParseException(
                "Unsupported case for discriminated type"
                + " parameters ["
                + "driverType="
                + str(driver_type)
                + "]"
            )

        read_buffer.pop_context("ModbusADU")
        # Create the instance
        _modbus_adu: ModbusADU = builder.build(response)
        return _modbus_adu

    def equals(self, o: object) -> bool:
        if self == o:
            return True

        if not isinstance(o, ModbusADU):
            return False

        that: ModbusADU = ModbusADU(o)
        return True

    def hash_code(self) -> int:
        return hash(self)

    def __str__(self) -> str:
        # TODO:- Implement a generic python object to probably json convertor or something.
        return ""


@dataclass
class ModbusADUBuilder:
    def build(self, response: bool) -> ModbusADU:
        pass
