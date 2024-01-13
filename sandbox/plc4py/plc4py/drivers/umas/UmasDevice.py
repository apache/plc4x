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

from bitarray import bitarray

from plc4py.protocols.umas.readwrite.DataItem import DataItem
from plc4py.protocols.umas.readwrite.UmasPDUError import UmasPDUError
from plc4py.spi.generation.ReadBuffer import ReadBuffer, ReadBufferByteBased

from plc4py.drivers.umas.UmasTag import (
    UmasTagHoldingRegister,
    UmasTagCoil,
    UmasTagDiscreteInput,
    UmasTagInputRegister,
)

from plc4py.api.exceptions.exceptions import PlcRuntimeException
from plc4py.drivers.umas.UmasConfiguration import UmasConfiguration
from plc4py.protocols.umas.readwrite.UmasPDUReadCoilsRequest import (
    UmasPDUReadCoilsRequest,
)
from plc4py.protocols.umas.readwrite.UmasPDUReadDiscreteInputsRequest import (
    UmasPDUReadDiscreteInputsRequest,
)
from plc4py.protocols.umas.readwrite.UmasPDUReadInputRegistersRequest import (
    UmasPDUReadInputRegistersRequest,
)
from plc4py.spi.generation.WriteBuffer import WriteBufferByteBased

from plc4py.api.messages.PlcRequest import PlcReadRequest
from plc4py.api.messages.PlcResponse import PlcReadResponse
from plc4py.api.value.PlcValue import PlcValue, PlcResponseCode
from plc4py.protocols.umas.readwrite.UmasPDUReadHoldingRegistersRequest import (
    UmasPDUReadHoldingRegistersRequest,
)
from plc4py.protocols.umas.readwrite.UmasTcpADU import UmasTcpADU
from plc4py.spi.messages.utils.ResponseItem import ResponseItem
from plc4py.spi.values.PlcValues import PlcList, PlcNull
from plc4py.utils.GenericTypes import ByteOrder, AtomicInteger


@dataclass
class UmasDevice:
    _configuration: UmasConfiguration
    tags: Dict[str, PlcValue] = field(default_factory=lambda: {})
    _transaction_generator: AtomicInteger = field(
        default_factory=lambda: AtomicInteger()
    )

    async def read(
        self, request: PlcReadRequest, transport: Transport
    ) -> PlcReadResponse:
        """
        Reads one field from the Umas Device
        """
        if len(request.tags) > 1:
            raise NotImplementedError(
                "The Umas driver only supports reading single tags at once"
            )
        if len(request.tags) == 0:
            raise PlcRuntimeException("No tags have been specified to read")
        tag = request.tags[request.tag_names[0]]
        logging.debug(f"Reading tag {str(tag)} from Umas Device")

        # Create future to be returned when a value is returned
        loop = asyncio.get_running_loop()
        message_future = loop.create_future()

        if isinstance(tag, UmasTagCoil):
            pdu = UmasPDUReadCoilsRequest(tag.address, tag.quantity)
        elif isinstance(tag, UmasTagDiscreteInput):
            pdu = UmasPDUReadDiscreteInputsRequest(tag.address, tag.quantity)
        elif isinstance(tag, UmasTagInputRegister):
            pdu = UmasPDUReadInputRegistersRequest(tag.address, tag.quantity)
        elif isinstance(tag, UmasTagHoldingRegister):
            pdu = UmasPDUReadHoldingRegistersRequest(tag.address, tag.quantity)
        else:
            raise NotImplementedError(
                "Umas tag type not implemented " + str(tag.__class__)
            )

        adu = UmasTcpADU(
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

        if isinstance(result, UmasPDUError):
            response_items = [
                ResponseItem(
                    PlcResponseCode.ACCESS_DENIED, PlcNull(result.exception_code)
                )
            ]

            response = PlcReadResponse(
                PlcResponseCode.OK, {request.tag_names[0]: response_items}
            )
            return response

        if isinstance(tag, UmasTagCoil) or isinstance(tag, UmasTagDiscreteInput):
            a = bitarray()
            a.frombytes(bytearray(result.value))
            a.bytereverse()
            read_buffer = ReadBufferByteBased(bytearray(a), ByteOrder.BIG_ENDIAN)
        else:
            read_buffer = ReadBufferByteBased(
                bytearray(result.value), ByteOrder.BIG_ENDIAN
            )
        returned_value = DataItem.static_parse(
            read_buffer,
            request.tags[request.tag_names[0]].data_type,
            request.tags[request.tag_names[0]].quantity,
        )

        response_items = [ResponseItem(PlcResponseCode.OK, returned_value)]

        response = PlcReadResponse(
            PlcResponseCode.OK, {request.tag_names[0]: response_items}
        )
        return response
