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
import asyncio
import logging
import math
from asyncio import Transport
from dataclasses import dataclass, field
from math import ceil
from typing import Dict, List

from bitarray import bitarray

from plc4py.api.exceptions.exceptions import PlcRuntimeException
from plc4py.api.messages.PlcRequest import PlcReadRequest, PlcWriteRequest
from plc4py.api.messages.PlcResponse import PlcReadResponse, PlcWriteResponse
from plc4py.api.value.PlcValue import PlcResponseCode, PlcValue
from plc4py.drivers.modbus.ModbusConfiguration import ModbusConfiguration
from plc4py.drivers.modbus.ModbusTag import (
    ModbusTagCoil,
    ModbusTagDiscreteInput,
    ModbusTagHoldingRegister,
    ModbusTagInputRegister,
)
from plc4py.protocols.modbus.readwrite.DataItem import DataItem
from plc4py.protocols.modbus.readwrite.ModbusPDUError import ModbusPDUError
from plc4py.protocols.modbus.readwrite.ModbusPDUReadCoilsRequest import (
    ModbusPDUReadCoilsRequest,
)
from plc4py.protocols.modbus.readwrite.ModbusPDUReadDiscreteInputsRequest import (
    ModbusPDUReadDiscreteInputsRequest,
)
from plc4py.protocols.modbus.readwrite.ModbusPDUReadHoldingRegistersRequest import (
    ModbusPDUReadHoldingRegistersRequest,
)
from plc4py.protocols.modbus.readwrite.ModbusPDUReadInputRegistersRequest import (
    ModbusPDUReadInputRegistersRequest,
)
from plc4py.protocols.modbus.readwrite.ModbusTcpADU import ModbusTcpADU
from plc4py.spi.generation.ReadBuffer import ReadBuffer, ReadBufferByteBased
from plc4py.spi.generation.WriteBuffer import WriteBufferByteBased
from plc4py.spi.messages.utils.ResponseItem import ResponseItem
from plc4py.spi.values.PlcValues import PlcList, PlcNull, PlcBOOL
from plc4py.utils.GenericTypes import AtomicInteger, ByteOrder

from plc4py.protocols.modbus.readwrite.ModbusPDUWriteMultipleCoilsRequest import (
    ModbusPDUWriteMultipleCoilsRequest,
)
from plc4py.protocols.modbus.readwrite.ModbusPDUWriteMultipleHoldingRegistersRequest import (
    ModbusPDUWriteMultipleHoldingRegistersRequestBuilder,
)

from plc4py.drivers.modbus.ModbusTag import ModbusTag
from plc4py.protocols.modbus.readwrite.ModbusDataType import ModbusDataType


@dataclass
class ModbusDevice:
    _configuration: ModbusConfiguration
    tags: Dict[str, PlcValue] = field(default_factory=lambda: {})
    _transaction_generator: AtomicInteger = field(
        default_factory=lambda: AtomicInteger()
    )

    async def read(
        self, request: PlcReadRequest, transport: Transport
    ) -> PlcReadResponse:
        """
        Reads one field from the Modbus Device
        """
        if len(request.tags) > 1:
            raise NotImplementedError(
                "The Modbus driver only supports reading single tags at once"
            )
        if len(request.tags) == 0:
            raise PlcRuntimeException("No tags have been specified to read")
        tag = request.tags[request.tag_names[0]]
        logging.debug(f"Reading tag {str(tag)} from Modbus Device")

        # Create future to be returned when a value is returned
        loop = asyncio.get_running_loop()
        message_future = loop.create_future()

        if isinstance(tag, ModbusTagCoil):
            if tag.data_type.value != ModbusDataType.BOOL.value:
                raise NotImplementedError(
                    f"Only BOOL data types can be used with the coil register area"
                )
            pdu = ModbusPDUReadCoilsRequest(tag.address, tag.quantity)
        elif isinstance(tag, ModbusTagDiscreteInput):
            if tag.data_type.value != ModbusDataType.BOOL.value:
                raise NotImplementedError(
                    f"Only BOOL data types can be used with the digital input register area"
                )
            pdu = ModbusPDUReadDiscreteInputsRequest(tag.address, tag.quantity)
        elif isinstance(tag, ModbusTagInputRegister):
            number_of_registers_per_item = tag.data_type.data_type_size / 2
            number_of_registers = ceil(tag.quantity * number_of_registers_per_item)
            pdu = ModbusPDUReadInputRegistersRequest(tag.address, number_of_registers)
        elif isinstance(tag, ModbusTagHoldingRegister):
            number_of_registers_per_item = tag.data_type.data_type_size / 2
            number_of_registers = ceil(tag.quantity * number_of_registers_per_item)
            pdu = ModbusPDUReadHoldingRegistersRequest(tag.address, number_of_registers)
        else:
            raise NotImplementedError(
                "Modbus tag type not implemented " + str(tag.__class__)
            )

        adu = ModbusTcpADU(
            False,
            self._transaction_generator.increment(),
            self._configuration.unit_identifier,
            pdu,
        )
        write_buffer = WriteBufferByteBased(adu.length_in_bytes(), ByteOrder.BIG_ENDIAN)
        adu.serialize(write_buffer)

        protocol = transport.protocol
        protocol.write_wait_for_response(
            write_buffer.get_bytes(),
            transport,
            adu.transaction_identifier,
            message_future,
        )

        await message_future
        result = message_future.result()

        if isinstance(result, ModbusPDUError):
            response_item = ResponseItem(
                PlcResponseCode.ACCESS_DENIED, PlcNull(result.exception_code)
            )

            response = PlcReadResponse(
                PlcResponseCode.OK, {request.tag_names[0]: response_item}
            )
            return response

        if isinstance(tag, ModbusTagCoil) or isinstance(tag, ModbusTagDiscreteInput):
            # As we need to do some funky stuff with the ordering of bits when reading bits
            # We aren't using the DataItem code for it

            # Normalize the array so we can just read the bits in one by one
            a = bitarray()
            a.frombytes(bytearray(result.value))
            a.bytereverse()
            read_buffer = ReadBufferByteBased(
                bytearray(a), self._configuration.byte_order
            )

            # If it's an array we need to wrap it in a PlcList
            quantity = request.tags[request.tag_names[0]].quantity
            if quantity == 1:
                returned_value = PlcBOOL(read_buffer.read_bit(""))
            else:
                returned_array = []
                for _ in range(quantity):
                    returned_array.append(PlcBOOL(read_buffer.read_bit("")))
                returned_value = PlcList(returned_array)
        else:
            read_buffer = ReadBufferByteBased(
                bytearray(result.value), self._configuration.byte_order
            )
            returned_value = DataItem.static_parse(
                read_buffer,
                request.tags[request.tag_names[0]].data_type,
                request.tags[request.tag_names[0]].quantity,
                True,
            )

        response_item = ResponseItem(PlcResponseCode.OK, returned_value)

        response = PlcReadResponse(
            PlcResponseCode.OK, {request.tag_names[0]: response_item}
        )
        return response

    async def write(
        self, request: PlcWriteRequest, transport: Transport
    ) -> PlcWriteResponse:
        """
        Writes one field from the Modbus Device
        """
        if len(request.tags) > 1:
            raise NotImplementedError(
                "The Modbus driver only supports writing single tags at once"
            )
        if len(request.tags) == 0:
            raise PlcRuntimeException("No tags have been specified to write")
        tag = request.tags[request.tag_names[0]]
        logging.debug(f"Writing tag {str(tag)} from Modbus Device")

        # Create future to be returned when a value is returned
        loop = asyncio.get_running_loop()
        message_future = loop.create_future()
        values = request.values[request.tag_names[0]]
        if isinstance(tag, ModbusTagCoil):
            values = self._serialize_data_items(tag, values)
            pdu = ModbusPDUWriteMultipleCoilsRequest(tag.address, tag.quantity, values)
        elif isinstance(tag, ModbusTagDiscreteInput):
            raise PlcRuntimeException(
                "Modbus doesn't support writing to discrete inputs"
            )
        elif isinstance(tag, ModbusTagInputRegister):
            raise PlcRuntimeException(
                "Modbus doesn't support writing to input registers"
            )
        elif isinstance(tag, ModbusTagHoldingRegister):
            values = self._serialize_data_items(tag, values)
            quantity = math.ceil(tag.quantity * (tag.data_type.data_type_size / 2))
            pdu = ModbusPDUWriteMultipleHoldingRegistersRequestBuilder(
                tag.address, quantity, values
            ).build()
        else:
            raise NotImplementedError(
                "Modbus tag type not implemented " + str(tag.__class__)
            )

        adu = ModbusTcpADU(
            False,
            self._transaction_generator.increment(),
            self._configuration.unit_identifier,
            pdu,
        )
        write_buffer = WriteBufferByteBased(adu.length_in_bytes(), ByteOrder.BIG_ENDIAN)
        adu.serialize(write_buffer)

        protocol = transport.protocol
        protocol.write_wait_for_response(
            write_buffer.get_bytes(),
            transport,
            adu.transaction_identifier,
            message_future,
        )

        await message_future
        result = message_future.result()
        if isinstance(result, ModbusPDUError):
            response_item = ResponseItem(PlcResponseCode.INVALID_ADDRESS, None)
        else:
            response_item = ResponseItem(PlcResponseCode.OK, None)
        write_response = PlcWriteResponse(
            PlcResponseCode.OK, {request.tag_names[0]: response_item}
        )
        return write_response

    def _serialize_data_items(self, tag: ModbusTag, values: PlcValue) -> List[int]:
        if (
            isinstance(tag, ModbusTagCoil)
            and (tag.data_type == ModbusDataType.BOOL)
            and not isinstance(values, PlcList)
        ):
            # Booleans are writen different depending on what registers they are being written to
            length = tag.quantity
            write_buffer = WriteBufferByteBased(length, self._configuration.byte_order)
            write_buffer.write_unsigned_short(int(0x0000), 7, "int0x0000")
            # Simple Field (value)
            value: bool = values.get_bool()
            write_buffer.write_bit((value), "value")
        elif (
            isinstance(tag, ModbusTagCoil)
            and (tag.data_type == ModbusDataType.BOOL)
            and isinstance(values, PlcList)
        ):
            length = math.ceil(tag.quantity / 8)
            write_buffer = WriteBufferByteBased(length, self._configuration.byte_order)
            for i in range(length):
                if i == length - 1:
                    sub_length = tag.quantity % 8
                else:
                    sub_length = 8
                sub_array = values.value[i * 8 : i * 8 + sub_length]
                sub_array.reverse()
                for j in range(8 - sub_length):
                    write_buffer.write_bit(False, "value")
                for item in sub_array:
                    # Simple Field (value)
                    value: bool = item.get_bool()
                    write_buffer.write_bit(value, "value")
        else:
            length = tag.quantity * tag.data_type.data_type_size
            if (length % 2) == 1:
                length += 1
            write_buffer = WriteBufferByteBased(length, self._configuration.byte_order)
            DataItem.static_serialize(
                write_buffer,
                values,
                tag.data_type,
                tag.quantity,
                True,
                self._configuration.byte_order,
            )
        return list(write_buffer.get_bytes().tobytes())
