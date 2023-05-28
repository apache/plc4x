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
from plc4py.api.exceptions.exceptions import PlcRuntimeException
from plc4py.api.messages.PlcMessage import PlcMessage
from plc4py.protocols.modbus.readwrite.DriverType import DriverType
from plc4py.spi.generation.ReadBuffer import ReadBuffer
from plc4py.spi.generation.WriteBuffer import WriteBuffer
import math


@dataclass
class ModbusADU(ABC, PlcMessage):
    # Arguments.
    response: bool

    def __post_init__(self):
        super().__init__()

    # Abstract accessors for discriminator values.
    @property
    @abstractmethod
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
        return int(math.ceil(float(self.get_length_in_bits() / 8.0)))

    def get_length_in_bits(self) -> int:
        length_in_bits: int = 0
        _value: ModbusADU = self

        # Length of subtype elements will be added by sub-type...

        return length_in_bits

    def static_parse(self, read_buffer: ReadBuffer, args):
        if args is None:
            raise PlcRuntimeException(
                "Wrong number of arguments, expected 2, but got None"
            )
        elif args.length != 2:
            raise PlcRuntimeException(
                "Wrong number of arguments, expected 2, but got " + str(len(args))
            )

        driverType: DriverType = 0
        if isinstance(args[0], DriverType):
            driverType = DriverType(args[0])
        elif isinstance(args[0], str):
            driverType = DriverType(str(args[0]))
        else:
            raise PlcRuntimeException(
                "Argument 0 expected to be of type DriverType or a string which is parseable but was "
                + args[0].getClass().getName()
            )

        response: bool = False
        if isinstance(args[1], bool):
            response = bool(args[1])
        elif isinstance(args[1], str):
            response = bool(str(args[1]))
        else:
            raise PlcRuntimeException(
                "Argument 1 expected to be of type bool or a string which is parseable but was "
                + args[1].getClass().getName()
            )

        return self.static_parse_context(read_buffer, driverType, response)

    @staticmethod
    def static_parse_context(
        read_buffer: ReadBuffer, driver_type: DriverType, response: bool
    ):
        read_buffer.push_context("ModbusADU")

        # Switch Field (Depending on the discriminator values, passes the instantiation to a sub-type)
        builder: ModbusADUBuilder = None
        if EvaluationHelper.equals(driverType, DriverType.get_modbu_s__tcp()):
            builder = ModbusTcpADU.staticParseBuilder(read_buffer, driverType, response)
        if EvaluationHelper.equals(driverType, DriverType.get_modbu_s__rtu()):
            builder = ModbusRtuADU.staticParseBuilder(read_buffer, driverType, response)
        if EvaluationHelper.equals(driverType, DriverType.get_modbu_s__ascii()):
            builder = ModbusAsciiADU.staticParseBuilder(
                read_buffer, driverType, response
            )
        if builder is None:
            raise ParseException(
                "Unsupported case for discriminated type"
                + " parameters ["
                + "driverType="
                + driverType
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
        write_buffer_box_based: WriteBufferBoxBased = WriteBufferBoxBased(True, True)
        try:
            write_buffer_box_based.writeSerializable(self)
        except SerializationException as e:
            raise RuntimeException(e)

        return "\n" + str(write_buffer_box_based.get_box()) + "\n"


class ModbusADUBuilder:
    def build(self, response: bool) -> ModbusADU:
        pass
