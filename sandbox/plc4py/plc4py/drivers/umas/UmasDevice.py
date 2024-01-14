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
from asyncio import Transport
from dataclasses import dataclass, field
from typing import Dict

from plc4py.spi.generation.WriteBuffer import WriteBufferByteBased

from plc4py.api.messages.PlcRequest import PlcReadRequest
from plc4py.api.messages.PlcResponse import PlcReadResponse
from plc4py.api.value.PlcValue import PlcValue
from plc4py.drivers.umas.UmasConfiguration import UmasConfiguration
from plc4py.protocols.umas.readwrite.ModbusTcpADU import ModbusTcpADU
from plc4py.protocols.umas.readwrite.UmasPDU import UmasPDU
from plc4py.protocols.umas.readwrite.UmasPDURequest import (
    UmasPDURequest,
    UmasPDURequestBuilder,
)
from plc4py.utils.GenericTypes import AtomicInteger, ByteOrder


@dataclass
class UmasDevice:
    _configuration: UmasConfiguration
    tags: Dict[str, PlcValue] = field(default_factory=lambda: {})
    _transaction_generator: AtomicInteger = field(
        default_factory=lambda: AtomicInteger()
    )

    async def connect(self, transport: Transport):
        # Create future to be returned when a value is returned
        loop = asyncio.get_running_loop()
        message_future = loop.create_future()

        pdu = UmasPDU(False, UmasPDURequestBuilder().build(0))

        adu = ModbusTcpADU(
            self._transaction_generator.increment(),
            self._configuration.unit_identifier,
            pdu,
            False,
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
        pass

    async def read(
        self, request: PlcReadRequest, transport: Transport
    ) -> PlcReadResponse:
        """
        Reads one field from the Umas Device
        """
        pass
