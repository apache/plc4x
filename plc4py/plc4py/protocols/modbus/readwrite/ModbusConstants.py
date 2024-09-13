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

from plc4py.api.exceptions.exceptions import PlcRuntimeException
from plc4py.api.exceptions.exceptions import SerializationException
from plc4py.api.messages.PlcMessage import PlcMessage
from plc4py.spi.generation.ReadBuffer import ReadBuffer
from plc4py.spi.generation.WriteBuffer import WriteBuffer
import math


@dataclass
class ModbusConstants:
    MODBUS_TCP_DEFAULT_PORT: int = int(502)

    def serialize(self, write_buffer: WriteBuffer):
        write_buffer.push_context("ModbusConstants")

        # Const Field (modbusTcpDefaultPort)
        write_buffer.write_unsigned_short(
            self.MODBUS_TCP_DEFAULT_PORT, logical_name="modbusTcpDefaultPort"
        )

        write_buffer.pop_context("ModbusConstants")

    def length_in_bytes(self) -> int:
        return int(math.ceil(float(self.length_in_bits() / 8.0)))

    def length_in_bits(self) -> int:
        length_in_bits: int = 0
        _value: ModbusConstants = self

        # Const Field (modbusTcpDefaultPort)
        length_in_bits += 16

        return length_in_bits

    @staticmethod
    def static_parse(read_buffer: ReadBuffer, **kwargs):
        return ModbusConstants.static_parse_context(read_buffer)

    @staticmethod
    def static_parse_context(read_buffer: ReadBuffer):
        read_buffer.push_context("ModbusConstants")

        MODBUS_TCP_DEFAULT_PORT: int = read_buffer.read_unsigned_short(
            logical_name="modbus_tcp_default_port"
        )

        read_buffer.pop_context("ModbusConstants")
        # Create the instance
        _modbus_constants: ModbusConstants = ModbusConstants()
        return _modbus_constants

    def equals(self, o: object) -> bool:
        if self == o:
            return True

        if not isinstance(o, ModbusConstants):
            return False

        that: ModbusConstants = ModbusConstants(o)
        return True

    def hash_code(self) -> int:
        return hash(self)

    def __str__(self) -> str:
        # TODO:- Implement a generic python object to probably json convertor or something.
        return ""
