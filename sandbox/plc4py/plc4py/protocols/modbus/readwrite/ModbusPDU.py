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
from ctypes import c_bool
from ctypes import c_uint8
from plc4py.api.messages.PlcMessage import PlcMessage
from plc4py.spi.generation.WriteBuffer import WriteBuffer
import math


@dataclass
class ModbusPDU(ABC, PlcMessage):
    def __post_init__(self):
        super().__init__()

    # Abstract accessors for discriminator values.
    @property
    @abstractmethod
    def error_flag(self) -> c_bool:
        pass

    @property
    @abstractmethod
    def function_flag(self) -> c_uint8:
        pass

    @property
    @abstractmethod
    def response(self) -> c_bool:
        pass

    @abstractmethod
    def serialize_modbus_pdu_child(self, write_buffer: WriteBuffer) -> None:
        pass

    def serialize(self, write_buffer: WriteBuffer):
        write_buffer.push_context("ModbusPDU")

        # Discriminator Field (errorFlag) (Used as input to a switch field)
        write_buffer.write_boolean(self.error_flag(), logical_name="errorFlag")

        # Discriminator Field (functionFlag) (Used as input to a switch field)
        write_buffer.write_unsigned_byte(
            self.function_flag(), logical_name="functionFlag"
        )

        # Switch field (Serialize the sub-type)
        self.serialize_modbus_pdu_child(write_buffer)

        write_buffer.pop_context("ModbusPDU")

    def length_in_bytes(self) -> int:
        return int(math.ceil(float(self.get_length_in_bits() / 8.0)))

    def get_length_in_bits(self) -> int:
        length_in_bits: int = 0
        _value: ModbusPDU = self

        # Discriminator Field (errorFlag)
        length_in_bits += 1

        # Discriminator Field (functionFlag)
        length_in_bits += 7

        # Length of subtype elements will be added by sub-type...

        return length_in_bits

    def static_parse(read_buffer: ReadBuffer, args):
        if (args is None) or (args.length is not 1):
            raise PlcRuntimeException(
                "Wrong number of arguments, expected 1, but got " + args.length
            )

        response: c_bool = None
        if isinstance(args[0], c_bool):
            response = c_bool(args[0])
        elif isinstance(args[0], str):
            response = c_bool.valueOf(str(args[0]))
        else:
            raise PlcRuntimeException(
                "Argument 0 expected to be of type c_bool or a string which is parseable but was "
                + args[0].getClass().getName()
            )

        return staticParse(read_buffer, response)

    @staticmethod
    def static_parse_context(read_buffer: ReadBuffer, response: c_bool):
        read_buffer.pull_context("ModbusPDU")
        cur_pos: int = 0

        error_flag: c_bool = read_discriminator_field("errorFlag", read_boolean)

        function_flag: c_uint8 = read_discriminator_field(
            "functionFlag", read_unsigned_short
        )

        # Switch Field (Depending on the discriminator values, passes the instantiation to a sub-type)
        builder: ModbusPDUBuilder = None
        if EvaluationHelper.equals(errorFlag, c_bool(True)):
            builder = ModbusPDUError.staticParseBuilder(read_buffer, response)
        if (
            EvaluationHelper.equals(errorFlag, c_bool(False))
            and EvaluationHelper.equals(functionFlag, c_uint8(0x02))
            and EvaluationHelper.equals(response, c_bool(False))
        ):
            builder = ModbusPDUReadDiscreteInputsRequest.staticParseBuilder(
                read_buffer, response
            )
        if (
            EvaluationHelper.equals(errorFlag, c_bool(False))
            and EvaluationHelper.equals(functionFlag, c_uint8(0x02))
            and EvaluationHelper.equals(response, c_bool(True))
        ):
            builder = ModbusPDUReadDiscreteInputsResponse.staticParseBuilder(
                read_buffer, response
            )
        if (
            EvaluationHelper.equals(errorFlag, c_bool(False))
            and EvaluationHelper.equals(functionFlag, c_uint8(0x01))
            and EvaluationHelper.equals(response, c_bool(False))
        ):
            builder = ModbusPDUReadCoilsRequest.staticParseBuilder(
                read_buffer, response
            )
        if (
            EvaluationHelper.equals(errorFlag, c_bool(False))
            and EvaluationHelper.equals(functionFlag, c_uint8(0x01))
            and EvaluationHelper.equals(response, c_bool(True))
        ):
            builder = ModbusPDUReadCoilsResponse.staticParseBuilder(
                read_buffer, response
            )
        if (
            EvaluationHelper.equals(errorFlag, c_bool(False))
            and EvaluationHelper.equals(functionFlag, c_uint8(0x05))
            and EvaluationHelper.equals(response, c_bool(False))
        ):
            builder = ModbusPDUWriteSingleCoilRequest.staticParseBuilder(
                read_buffer, response
            )
        if (
            EvaluationHelper.equals(errorFlag, c_bool(False))
            and EvaluationHelper.equals(functionFlag, c_uint8(0x05))
            and EvaluationHelper.equals(response, c_bool(True))
        ):
            builder = ModbusPDUWriteSingleCoilResponse.staticParseBuilder(
                read_buffer, response
            )
        if (
            EvaluationHelper.equals(errorFlag, c_bool(False))
            and EvaluationHelper.equals(functionFlag, c_uint8(0x0F))
            and EvaluationHelper.equals(response, c_bool(False))
        ):
            builder = ModbusPDUWriteMultipleCoilsRequest.staticParseBuilder(
                read_buffer, response
            )
        if (
            EvaluationHelper.equals(errorFlag, c_bool(False))
            and EvaluationHelper.equals(functionFlag, c_uint8(0x0F))
            and EvaluationHelper.equals(response, c_bool(True))
        ):
            builder = ModbusPDUWriteMultipleCoilsResponse.staticParseBuilder(
                read_buffer, response
            )
        if (
            EvaluationHelper.equals(errorFlag, c_bool(False))
            and EvaluationHelper.equals(functionFlag, c_uint8(0x04))
            and EvaluationHelper.equals(response, c_bool(False))
        ):
            builder = ModbusPDUReadInputRegistersRequest.staticParseBuilder(
                read_buffer, response
            )
        if (
            EvaluationHelper.equals(errorFlag, c_bool(False))
            and EvaluationHelper.equals(functionFlag, c_uint8(0x04))
            and EvaluationHelper.equals(response, c_bool(True))
        ):
            builder = ModbusPDUReadInputRegistersResponse.staticParseBuilder(
                read_buffer, response
            )
        if (
            EvaluationHelper.equals(errorFlag, c_bool(False))
            and EvaluationHelper.equals(functionFlag, c_uint8(0x03))
            and EvaluationHelper.equals(response, c_bool(False))
        ):
            builder = ModbusPDUReadHoldingRegistersRequest.staticParseBuilder(
                read_buffer, response
            )
        if (
            EvaluationHelper.equals(errorFlag, c_bool(False))
            and EvaluationHelper.equals(functionFlag, c_uint8(0x03))
            and EvaluationHelper.equals(response, c_bool(True))
        ):
            builder = ModbusPDUReadHoldingRegistersResponse.staticParseBuilder(
                read_buffer, response
            )
        if (
            EvaluationHelper.equals(errorFlag, c_bool(False))
            and EvaluationHelper.equals(functionFlag, c_uint8(0x06))
            and EvaluationHelper.equals(response, c_bool(False))
        ):
            builder = ModbusPDUWriteSingleRegisterRequest.staticParseBuilder(
                read_buffer, response
            )
        if (
            EvaluationHelper.equals(errorFlag, c_bool(False))
            and EvaluationHelper.equals(functionFlag, c_uint8(0x06))
            and EvaluationHelper.equals(response, c_bool(True))
        ):
            builder = ModbusPDUWriteSingleRegisterResponse.staticParseBuilder(
                read_buffer, response
            )
        if (
            EvaluationHelper.equals(errorFlag, c_bool(False))
            and EvaluationHelper.equals(functionFlag, c_uint8(0x10))
            and EvaluationHelper.equals(response, c_bool(False))
        ):
            builder = ModbusPDUWriteMultipleHoldingRegistersRequest.staticParseBuilder(
                read_buffer, response
            )
        if (
            EvaluationHelper.equals(errorFlag, c_bool(False))
            and EvaluationHelper.equals(functionFlag, c_uint8(0x10))
            and EvaluationHelper.equals(response, c_bool(True))
        ):
            builder = ModbusPDUWriteMultipleHoldingRegistersResponse.staticParseBuilder(
                read_buffer, response
            )
        if (
            EvaluationHelper.equals(errorFlag, c_bool(False))
            and EvaluationHelper.equals(functionFlag, c_uint8(0x17))
            and EvaluationHelper.equals(response, c_bool(False))
        ):
            builder = (
                ModbusPDUReadWriteMultipleHoldingRegistersRequest.staticParseBuilder(
                    read_buffer, response
                )
            )
        if (
            EvaluationHelper.equals(errorFlag, c_bool(False))
            and EvaluationHelper.equals(functionFlag, c_uint8(0x17))
            and EvaluationHelper.equals(response, c_bool(True))
        ):
            builder = (
                ModbusPDUReadWriteMultipleHoldingRegistersResponse.staticParseBuilder(
                    read_buffer, response
                )
            )
        if (
            EvaluationHelper.equals(errorFlag, c_bool(False))
            and EvaluationHelper.equals(functionFlag, c_uint8(0x16))
            and EvaluationHelper.equals(response, c_bool(False))
        ):
            builder = ModbusPDUMaskWriteHoldingRegisterRequest.staticParseBuilder(
                read_buffer, response
            )
        if (
            EvaluationHelper.equals(errorFlag, c_bool(False))
            and EvaluationHelper.equals(functionFlag, c_uint8(0x16))
            and EvaluationHelper.equals(response, c_bool(True))
        ):
            builder = ModbusPDUMaskWriteHoldingRegisterResponse.staticParseBuilder(
                read_buffer, response
            )
        if (
            EvaluationHelper.equals(errorFlag, c_bool(False))
            and EvaluationHelper.equals(functionFlag, c_uint8(0x18))
            and EvaluationHelper.equals(response, c_bool(False))
        ):
            builder = ModbusPDUReadFifoQueueRequest.staticParseBuilder(
                read_buffer, response
            )
        if (
            EvaluationHelper.equals(errorFlag, c_bool(False))
            and EvaluationHelper.equals(functionFlag, c_uint8(0x18))
            and EvaluationHelper.equals(response, c_bool(True))
        ):
            builder = ModbusPDUReadFifoQueueResponse.staticParseBuilder(
                read_buffer, response
            )
        if (
            EvaluationHelper.equals(errorFlag, c_bool(False))
            and EvaluationHelper.equals(functionFlag, c_uint8(0x14))
            and EvaluationHelper.equals(response, c_bool(False))
        ):
            builder = ModbusPDUReadFileRecordRequest.staticParseBuilder(
                read_buffer, response
            )
        if (
            EvaluationHelper.equals(errorFlag, c_bool(False))
            and EvaluationHelper.equals(functionFlag, c_uint8(0x14))
            and EvaluationHelper.equals(response, c_bool(True))
        ):
            builder = ModbusPDUReadFileRecordResponse.staticParseBuilder(
                read_buffer, response
            )
        if (
            EvaluationHelper.equals(errorFlag, c_bool(False))
            and EvaluationHelper.equals(functionFlag, c_uint8(0x15))
            and EvaluationHelper.equals(response, c_bool(False))
        ):
            builder = ModbusPDUWriteFileRecordRequest.staticParseBuilder(
                read_buffer, response
            )
        if (
            EvaluationHelper.equals(errorFlag, c_bool(False))
            and EvaluationHelper.equals(functionFlag, c_uint8(0x15))
            and EvaluationHelper.equals(response, c_bool(True))
        ):
            builder = ModbusPDUWriteFileRecordResponse.staticParseBuilder(
                read_buffer, response
            )
        if (
            EvaluationHelper.equals(errorFlag, c_bool(False))
            and EvaluationHelper.equals(functionFlag, c_uint8(0x07))
            and EvaluationHelper.equals(response, c_bool(False))
        ):
            builder = ModbusPDUReadExceptionStatusRequest.staticParseBuilder(
                read_buffer, response
            )
        if (
            EvaluationHelper.equals(errorFlag, c_bool(False))
            and EvaluationHelper.equals(functionFlag, c_uint8(0x07))
            and EvaluationHelper.equals(response, c_bool(True))
        ):
            builder = ModbusPDUReadExceptionStatusResponse.staticParseBuilder(
                read_buffer, response
            )
        if (
            EvaluationHelper.equals(errorFlag, c_bool(False))
            and EvaluationHelper.equals(functionFlag, c_uint8(0x08))
            and EvaluationHelper.equals(response, c_bool(False))
        ):
            builder = ModbusPDUDiagnosticRequest.staticParseBuilder(
                read_buffer, response
            )
        if (
            EvaluationHelper.equals(errorFlag, c_bool(False))
            and EvaluationHelper.equals(functionFlag, c_uint8(0x08))
            and EvaluationHelper.equals(response, c_bool(True))
        ):
            builder = ModbusPDUDiagnosticResponse.staticParseBuilder(
                read_buffer, response
            )
        if (
            EvaluationHelper.equals(errorFlag, c_bool(False))
            and EvaluationHelper.equals(functionFlag, c_uint8(0x0B))
            and EvaluationHelper.equals(response, c_bool(False))
        ):
            builder = ModbusPDUGetComEventCounterRequest.staticParseBuilder(
                read_buffer, response
            )
        if (
            EvaluationHelper.equals(errorFlag, c_bool(False))
            and EvaluationHelper.equals(functionFlag, c_uint8(0x0B))
            and EvaluationHelper.equals(response, c_bool(True))
        ):
            builder = ModbusPDUGetComEventCounterResponse.staticParseBuilder(
                read_buffer, response
            )
        if (
            EvaluationHelper.equals(errorFlag, c_bool(False))
            and EvaluationHelper.equals(functionFlag, c_uint8(0x0C))
            and EvaluationHelper.equals(response, c_bool(False))
        ):
            builder = ModbusPDUGetComEventLogRequest.staticParseBuilder(
                read_buffer, response
            )
        if (
            EvaluationHelper.equals(errorFlag, c_bool(False))
            and EvaluationHelper.equals(functionFlag, c_uint8(0x0C))
            and EvaluationHelper.equals(response, c_bool(True))
        ):
            builder = ModbusPDUGetComEventLogResponse.staticParseBuilder(
                read_buffer, response
            )
        if (
            EvaluationHelper.equals(errorFlag, c_bool(False))
            and EvaluationHelper.equals(functionFlag, c_uint8(0x11))
            and EvaluationHelper.equals(response, c_bool(False))
        ):
            builder = ModbusPDUReportServerIdRequest.staticParseBuilder(
                read_buffer, response
            )
        if (
            EvaluationHelper.equals(errorFlag, c_bool(False))
            and EvaluationHelper.equals(functionFlag, c_uint8(0x11))
            and EvaluationHelper.equals(response, c_bool(True))
        ):
            builder = ModbusPDUReportServerIdResponse.staticParseBuilder(
                read_buffer, response
            )
        if (
            EvaluationHelper.equals(errorFlag, c_bool(False))
            and EvaluationHelper.equals(functionFlag, c_uint8(0x2B))
            and EvaluationHelper.equals(response, c_bool(False))
        ):
            builder = ModbusPDUReadDeviceIdentificationRequest.staticParseBuilder(
                read_buffer, response
            )
        if (
            EvaluationHelper.equals(errorFlag, c_bool(False))
            and EvaluationHelper.equals(functionFlag, c_uint8(0x2B))
            and EvaluationHelper.equals(response, c_bool(True))
        ):
            builder = ModbusPDUReadDeviceIdentificationResponse.staticParseBuilder(
                read_buffer, response
            )
        if builder is None:
            raise ParseException(
                "Unsupported case for discriminated type"
                + " parameters ["
                + "errorFlag="
                + errorFlag
                + " "
                + "functionFlag="
                + functionFlag
                + " "
                + "response="
                + response
                + "]"
            )

        read_buffer.close_context("ModbusPDU")
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
        write_buffer_box_based: WriteBufferBoxBased = WriteBufferBoxBased(True, True)
        try:
            write_buffer_box_based.writeSerializable(self)
        except SerializationException as e:
            raise RuntimeException(e)

        return "\n" + str(write_buffer_box_based.get_box()) + "\n"


class ModbusPDUBuilder:
    def build(
        self,
    ) -> ModbusPDU:
        pass
