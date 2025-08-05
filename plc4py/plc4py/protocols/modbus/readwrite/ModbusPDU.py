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
from plc4py.spi.generation.ReadBuffer import ReadBuffer
from plc4py.spi.generation.WriteBuffer import WriteBuffer
from plc4py.utils.ConnectionStringHandling import strtobool
import math


@dataclass
class ModbusPDU(ABC, PlcMessage):

    # Abstract accessors for discriminator values.
    @property
    def error_flag(self) -> bool:
        pass

    @property
    def function_flag(self) -> int:
        pass

    @property
    def response(self) -> bool:
        pass

    @abstractmethod
    def serialize_modbus_pdu_child(self, write_buffer: WriteBuffer) -> None:
        pass

    def serialize(self, write_buffer: WriteBuffer):
        write_buffer.push_context("ModbusPDU")

        # Discriminator Field (errorFlag) (Used as input to a switch field)
        write_buffer.write_bit(self.error_flag, logical_name="error_flag", bit_length=1)

        # Discriminator Field (functionFlag) (Used as input to a switch field)
        write_buffer.write_unsigned_byte(
            self.function_flag, logical_name="function_flag", bit_length=7
        )

        # Switch field (Serialize the sub-type)
        self.serialize_modbus_pdu_child(write_buffer)

        write_buffer.pop_context("ModbusPDU")

    def length_in_bytes(self) -> int:
        return int(math.ceil(float(self.length_in_bits() / 8.0)))

    def length_in_bits(self) -> int:
        length_in_bits: int = 0
        _value: ModbusPDU = self

        # Discriminator Field (errorFlag)
        length_in_bits += 1

        # Discriminator Field (functionFlag)
        length_in_bits += 7

        # Length of subtype elements will be added by sub-type...

        return length_in_bits

    @staticmethod
    def static_parse(read_buffer: ReadBuffer, **kwargs):

        if kwargs is None:
            raise PlcRuntimeException(
                "Wrong number of arguments, expected 1, but got None"
            )

        response: bool = False
        if isinstance(kwargs.get("response"), bool):
            response = bool(kwargs.get("response"))
        elif isinstance(kwargs.get("response"), str):
            response = bool(str(kwargs.get("response")))
        else:
            raise PlcRuntimeException(
                "Argument 0 expected to be of type bool or a string which is parseable but was "
                + kwargs.get("response").getClass().getName()
            )

        return ModbusPDU.static_parse_context(read_buffer, response)

    @staticmethod
    def static_parse_context(read_buffer: ReadBuffer, response: bool):
        read_buffer.push_context("ModbusPDU")

        if isinstance(response, str):
            response = bool(strtobool(response))

        error_flag: bool = read_buffer.read_bit(
            logical_name="error_flag", bit_length=1, response=response
        )

        function_flag: int = read_buffer.read_unsigned_byte(
            logical_name="function_flag", bit_length=7, response=response
        )

        # Switch Field (Depending on the discriminator values, passes the instantiation to a sub-type)
        builder: ModbusPDUBuilder = None
        from plc4py.protocols.modbus.readwrite.ModbusPDUError import ModbusPDUError

        if error_flag == bool(True):

            builder = ModbusPDUError.static_parse_builder(read_buffer, response)
        from plc4py.protocols.modbus.readwrite.ModbusPDUReadDiscreteInputsRequest import (
            ModbusPDUReadDiscreteInputsRequest,
        )

        if (
            error_flag == bool(False)
            and function_flag == int(0x02)
            and response == bool(False)
        ):

            builder = ModbusPDUReadDiscreteInputsRequest.static_parse_builder(
                read_buffer, response
            )
        from plc4py.protocols.modbus.readwrite.ModbusPDUReadDiscreteInputsResponse import (
            ModbusPDUReadDiscreteInputsResponse,
        )

        if (
            error_flag == bool(False)
            and function_flag == int(0x02)
            and response == bool(True)
        ):

            builder = ModbusPDUReadDiscreteInputsResponse.static_parse_builder(
                read_buffer, response
            )
        from plc4py.protocols.modbus.readwrite.ModbusPDUReadCoilsRequest import (
            ModbusPDUReadCoilsRequest,
        )

        if (
            error_flag == bool(False)
            and function_flag == int(0x01)
            and response == bool(False)
        ):

            builder = ModbusPDUReadCoilsRequest.static_parse_builder(
                read_buffer, response
            )
        from plc4py.protocols.modbus.readwrite.ModbusPDUReadCoilsResponse import (
            ModbusPDUReadCoilsResponse,
        )

        if (
            error_flag == bool(False)
            and function_flag == int(0x01)
            and response == bool(True)
        ):

            builder = ModbusPDUReadCoilsResponse.static_parse_builder(
                read_buffer, response
            )
        from plc4py.protocols.modbus.readwrite.ModbusPDUWriteSingleCoilRequest import (
            ModbusPDUWriteSingleCoilRequest,
        )

        if (
            error_flag == bool(False)
            and function_flag == int(0x05)
            and response == bool(False)
        ):

            builder = ModbusPDUWriteSingleCoilRequest.static_parse_builder(
                read_buffer, response
            )
        from plc4py.protocols.modbus.readwrite.ModbusPDUWriteSingleCoilResponse import (
            ModbusPDUWriteSingleCoilResponse,
        )

        if (
            error_flag == bool(False)
            and function_flag == int(0x05)
            and response == bool(True)
        ):

            builder = ModbusPDUWriteSingleCoilResponse.static_parse_builder(
                read_buffer, response
            )
        from plc4py.protocols.modbus.readwrite.ModbusPDUWriteMultipleCoilsRequest import (
            ModbusPDUWriteMultipleCoilsRequest,
        )

        if (
            error_flag == bool(False)
            and function_flag == int(0x0F)
            and response == bool(False)
        ):

            builder = ModbusPDUWriteMultipleCoilsRequest.static_parse_builder(
                read_buffer, response
            )
        from plc4py.protocols.modbus.readwrite.ModbusPDUWriteMultipleCoilsResponse import (
            ModbusPDUWriteMultipleCoilsResponse,
        )

        if (
            error_flag == bool(False)
            and function_flag == int(0x0F)
            and response == bool(True)
        ):

            builder = ModbusPDUWriteMultipleCoilsResponse.static_parse_builder(
                read_buffer, response
            )
        from plc4py.protocols.modbus.readwrite.ModbusPDUReadInputRegistersRequest import (
            ModbusPDUReadInputRegistersRequest,
        )

        if (
            error_flag == bool(False)
            and function_flag == int(0x04)
            and response == bool(False)
        ):

            builder = ModbusPDUReadInputRegistersRequest.static_parse_builder(
                read_buffer, response
            )
        from plc4py.protocols.modbus.readwrite.ModbusPDUReadInputRegistersResponse import (
            ModbusPDUReadInputRegistersResponse,
        )

        if (
            error_flag == bool(False)
            and function_flag == int(0x04)
            and response == bool(True)
        ):

            builder = ModbusPDUReadInputRegistersResponse.static_parse_builder(
                read_buffer, response
            )
        from plc4py.protocols.modbus.readwrite.ModbusPDUReadHoldingRegistersRequest import (
            ModbusPDUReadHoldingRegistersRequest,
        )

        if (
            error_flag == bool(False)
            and function_flag == int(0x03)
            and response == bool(False)
        ):

            builder = ModbusPDUReadHoldingRegistersRequest.static_parse_builder(
                read_buffer, response
            )
        from plc4py.protocols.modbus.readwrite.ModbusPDUReadHoldingRegistersResponse import (
            ModbusPDUReadHoldingRegistersResponse,
        )

        if (
            error_flag == bool(False)
            and function_flag == int(0x03)
            and response == bool(True)
        ):

            builder = ModbusPDUReadHoldingRegistersResponse.static_parse_builder(
                read_buffer, response
            )
        from plc4py.protocols.modbus.readwrite.ModbusPDUWriteSingleRegisterRequest import (
            ModbusPDUWriteSingleRegisterRequest,
        )

        if (
            error_flag == bool(False)
            and function_flag == int(0x06)
            and response == bool(False)
        ):

            builder = ModbusPDUWriteSingleRegisterRequest.static_parse_builder(
                read_buffer, response
            )
        from plc4py.protocols.modbus.readwrite.ModbusPDUWriteSingleRegisterResponse import (
            ModbusPDUWriteSingleRegisterResponse,
        )

        if (
            error_flag == bool(False)
            and function_flag == int(0x06)
            and response == bool(True)
        ):

            builder = ModbusPDUWriteSingleRegisterResponse.static_parse_builder(
                read_buffer, response
            )
        from plc4py.protocols.modbus.readwrite.ModbusPDUWriteMultipleHoldingRegistersRequest import (
            ModbusPDUWriteMultipleHoldingRegistersRequest,
        )

        if (
            error_flag == bool(False)
            and function_flag == int(0x10)
            and response == bool(False)
        ):

            builder = (
                ModbusPDUWriteMultipleHoldingRegistersRequest.static_parse_builder(
                    read_buffer, response
                )
            )
        from plc4py.protocols.modbus.readwrite.ModbusPDUWriteMultipleHoldingRegistersResponse import (
            ModbusPDUWriteMultipleHoldingRegistersResponse,
        )

        if (
            error_flag == bool(False)
            and function_flag == int(0x10)
            and response == bool(True)
        ):

            builder = (
                ModbusPDUWriteMultipleHoldingRegistersResponse.static_parse_builder(
                    read_buffer, response
                )
            )
        from plc4py.protocols.modbus.readwrite.ModbusPDUReadWriteMultipleHoldingRegistersRequest import (
            ModbusPDUReadWriteMultipleHoldingRegistersRequest,
        )

        if (
            error_flag == bool(False)
            and function_flag == int(0x17)
            and response == bool(False)
        ):

            builder = (
                ModbusPDUReadWriteMultipleHoldingRegistersRequest.static_parse_builder(
                    read_buffer, response
                )
            )
        from plc4py.protocols.modbus.readwrite.ModbusPDUReadWriteMultipleHoldingRegistersResponse import (
            ModbusPDUReadWriteMultipleHoldingRegistersResponse,
        )

        if (
            error_flag == bool(False)
            and function_flag == int(0x17)
            and response == bool(True)
        ):

            builder = (
                ModbusPDUReadWriteMultipleHoldingRegistersResponse.static_parse_builder(
                    read_buffer, response
                )
            )
        from plc4py.protocols.modbus.readwrite.ModbusPDUMaskWriteHoldingRegisterRequest import (
            ModbusPDUMaskWriteHoldingRegisterRequest,
        )

        if (
            error_flag == bool(False)
            and function_flag == int(0x16)
            and response == bool(False)
        ):

            builder = ModbusPDUMaskWriteHoldingRegisterRequest.static_parse_builder(
                read_buffer, response
            )
        from plc4py.protocols.modbus.readwrite.ModbusPDUMaskWriteHoldingRegisterResponse import (
            ModbusPDUMaskWriteHoldingRegisterResponse,
        )

        if (
            error_flag == bool(False)
            and function_flag == int(0x16)
            and response == bool(True)
        ):

            builder = ModbusPDUMaskWriteHoldingRegisterResponse.static_parse_builder(
                read_buffer, response
            )
        from plc4py.protocols.modbus.readwrite.ModbusPDUReadFifoQueueRequest import (
            ModbusPDUReadFifoQueueRequest,
        )

        if (
            error_flag == bool(False)
            and function_flag == int(0x18)
            and response == bool(False)
        ):

            builder = ModbusPDUReadFifoQueueRequest.static_parse_builder(
                read_buffer, response
            )
        from plc4py.protocols.modbus.readwrite.ModbusPDUReadFifoQueueResponse import (
            ModbusPDUReadFifoQueueResponse,
        )

        if (
            error_flag == bool(False)
            and function_flag == int(0x18)
            and response == bool(True)
        ):

            builder = ModbusPDUReadFifoQueueResponse.static_parse_builder(
                read_buffer, response
            )
        from plc4py.protocols.modbus.readwrite.ModbusPDUReadFileRecordRequest import (
            ModbusPDUReadFileRecordRequest,
        )

        if (
            error_flag == bool(False)
            and function_flag == int(0x14)
            and response == bool(False)
        ):

            builder = ModbusPDUReadFileRecordRequest.static_parse_builder(
                read_buffer, response
            )
        from plc4py.protocols.modbus.readwrite.ModbusPDUReadFileRecordResponse import (
            ModbusPDUReadFileRecordResponse,
        )

        if (
            error_flag == bool(False)
            and function_flag == int(0x14)
            and response == bool(True)
        ):

            builder = ModbusPDUReadFileRecordResponse.static_parse_builder(
                read_buffer, response
            )
        from plc4py.protocols.modbus.readwrite.ModbusPDUWriteFileRecordRequest import (
            ModbusPDUWriteFileRecordRequest,
        )

        if (
            error_flag == bool(False)
            and function_flag == int(0x15)
            and response == bool(False)
        ):

            builder = ModbusPDUWriteFileRecordRequest.static_parse_builder(
                read_buffer, response
            )
        from plc4py.protocols.modbus.readwrite.ModbusPDUWriteFileRecordResponse import (
            ModbusPDUWriteFileRecordResponse,
        )

        if (
            error_flag == bool(False)
            and function_flag == int(0x15)
            and response == bool(True)
        ):

            builder = ModbusPDUWriteFileRecordResponse.static_parse_builder(
                read_buffer, response
            )
        from plc4py.protocols.modbus.readwrite.ModbusPDUReadExceptionStatusRequest import (
            ModbusPDUReadExceptionStatusRequest,
        )

        if (
            error_flag == bool(False)
            and function_flag == int(0x07)
            and response == bool(False)
        ):

            builder = ModbusPDUReadExceptionStatusRequest.static_parse_builder(
                read_buffer, response
            )
        from plc4py.protocols.modbus.readwrite.ModbusPDUReadExceptionStatusResponse import (
            ModbusPDUReadExceptionStatusResponse,
        )

        if (
            error_flag == bool(False)
            and function_flag == int(0x07)
            and response == bool(True)
        ):

            builder = ModbusPDUReadExceptionStatusResponse.static_parse_builder(
                read_buffer, response
            )
        from plc4py.protocols.modbus.readwrite.ModbusPDUDiagnosticRequest import (
            ModbusPDUDiagnosticRequest,
        )

        if (
            error_flag == bool(False)
            and function_flag == int(0x08)
            and response == bool(False)
        ):

            builder = ModbusPDUDiagnosticRequest.static_parse_builder(
                read_buffer, response
            )
        from plc4py.protocols.modbus.readwrite.ModbusPDUDiagnosticResponse import (
            ModbusPDUDiagnosticResponse,
        )

        if (
            error_flag == bool(False)
            and function_flag == int(0x08)
            and response == bool(True)
        ):

            builder = ModbusPDUDiagnosticResponse.static_parse_builder(
                read_buffer, response
            )
        from plc4py.protocols.modbus.readwrite.ModbusPDUGetComEventCounterRequest import (
            ModbusPDUGetComEventCounterRequest,
        )

        if (
            error_flag == bool(False)
            and function_flag == int(0x0B)
            and response == bool(False)
        ):

            builder = ModbusPDUGetComEventCounterRequest.static_parse_builder(
                read_buffer, response
            )
        from plc4py.protocols.modbus.readwrite.ModbusPDUGetComEventCounterResponse import (
            ModbusPDUGetComEventCounterResponse,
        )

        if (
            error_flag == bool(False)
            and function_flag == int(0x0B)
            and response == bool(True)
        ):

            builder = ModbusPDUGetComEventCounterResponse.static_parse_builder(
                read_buffer, response
            )
        from plc4py.protocols.modbus.readwrite.ModbusPDUGetComEventLogRequest import (
            ModbusPDUGetComEventLogRequest,
        )

        if (
            error_flag == bool(False)
            and function_flag == int(0x0C)
            and response == bool(False)
        ):

            builder = ModbusPDUGetComEventLogRequest.static_parse_builder(
                read_buffer, response
            )
        from plc4py.protocols.modbus.readwrite.ModbusPDUGetComEventLogResponse import (
            ModbusPDUGetComEventLogResponse,
        )

        if (
            error_flag == bool(False)
            and function_flag == int(0x0C)
            and response == bool(True)
        ):

            builder = ModbusPDUGetComEventLogResponse.static_parse_builder(
                read_buffer, response
            )
        from plc4py.protocols.modbus.readwrite.ModbusPDUReportServerIdRequest import (
            ModbusPDUReportServerIdRequest,
        )

        if (
            error_flag == bool(False)
            and function_flag == int(0x11)
            and response == bool(False)
        ):

            builder = ModbusPDUReportServerIdRequest.static_parse_builder(
                read_buffer, response
            )
        from plc4py.protocols.modbus.readwrite.ModbusPDUReportServerIdResponse import (
            ModbusPDUReportServerIdResponse,
        )

        if (
            error_flag == bool(False)
            and function_flag == int(0x11)
            and response == bool(True)
        ):

            builder = ModbusPDUReportServerIdResponse.static_parse_builder(
                read_buffer, response
            )
        from plc4py.protocols.modbus.readwrite.ModbusPDUReadDeviceIdentificationRequest import (
            ModbusPDUReadDeviceIdentificationRequest,
        )

        if (
            error_flag == bool(False)
            and function_flag == int(0x2B)
            and response == bool(False)
        ):

            builder = ModbusPDUReadDeviceIdentificationRequest.static_parse_builder(
                read_buffer, response
            )
        from plc4py.protocols.modbus.readwrite.ModbusPDUReadDeviceIdentificationResponse import (
            ModbusPDUReadDeviceIdentificationResponse,
        )

        if (
            error_flag == bool(False)
            and function_flag == int(0x2B)
            and response == bool(True)
        ):

            builder = ModbusPDUReadDeviceIdentificationResponse.static_parse_builder(
                read_buffer, response
            )
        if builder is None:
            raise ParseException(
                "Unsupported case for discriminated type"
                + " parameters ["
                + "errorFlag="
                + str(error_flag)
                + " "
                + "functionFlag="
                + str(function_flag)
                + " "
                + "response="
                + str(response)
                + "]"
            )

        read_buffer.pop_context("ModbusPDU")
        # Create the instance
        _modbus_pdu: ModbusPDU = builder.build()
        return _modbus_pdu

    def equals(self, o: object) -> bool:
        if self == o:
            return True

        if not isinstance(o, ModbusPDU):
            return False

        that: ModbusPDU = ModbusPDU(o)
        return True

    def hash_code(self) -> int:
        return hash(self)

    def __str__(self) -> str:
        # TODO:- Implement a generic python object to probably json convertor or something.
        return ""


@dataclass
class ModbusPDUBuilder:
    def build(
        self,
    ) -> ModbusPDU:
        pass
