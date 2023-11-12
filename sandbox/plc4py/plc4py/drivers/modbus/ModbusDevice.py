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
from asyncio import Transport
from dataclasses import dataclass, field
from typing import Dict, List

from plc4py.drivers.modbus.ModbusTag import (
    ModbusTagHoldingRegister,
    ModbusTagCoil,
    ModbusTagDiscreteInput,
    ModbusTagInputRegister,
)

from plc4py.api.exceptions.exceptions import PlcRuntimeException
from plc4py.drivers.modbus.ModbusConfiguration import ModbusConfiguration
from plc4py.protocols.modbus.readwrite.ModbusPDUReadCoilsRequest import (
    ModbusPDUReadCoilsRequest,
)
from plc4py.protocols.modbus.readwrite.ModbusPDUReadDiscreteInputsRequest import (
    ModbusPDUReadDiscreteInputsRequest,
)
from plc4py.protocols.modbus.readwrite.ModbusPDUReadInputRegistersRequest import (
    ModbusPDUReadInputRegistersRequest,
)
from plc4py.spi.generation.WriteBuffer import WriteBufferByteBased

from plc4py.api.messages.PlcRequest import PlcReadRequest
from plc4py.api.messages.PlcResponse import PlcReadResponse
from plc4py.api.value.PlcValue import PlcValue, PlcResponseCode
from plc4py.protocols.modbus.readwrite.ModbusPDUReadHoldingRegistersRequest import (
    ModbusPDUReadHoldingRegistersRequest,
)
from plc4py.protocols.modbus.readwrite.ModbusTcpADU import ModbusTcpADU
from plc4py.utils.GenericTypes import ByteOrder, AtomicInteger


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
        Reads one field from the Mock Device
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
            pdu = ModbusPDUReadCoilsRequest(tag.address, tag.quantity)
        elif isinstance(tag, ModbusTagDiscreteInput):
            pdu = ModbusPDUReadDiscreteInputsRequest(tag.address, tag.quantity)
        elif isinstance(tag, ModbusTagInputRegister):
            pdu = ModbusPDUReadInputRegistersRequest(tag.address, tag.quantity)
        elif isinstance(tag, ModbusTagHoldingRegister):
            pdu = ModbusPDUReadHoldingRegistersRequest(tag.address, tag.quantity)
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

        response = PlcReadResponse(PlcResponseCode.OK, [], {})
        return response
