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
from plc4py.spi.generation.WriteBuffer import WriteBufferByteBased

from plc4py.protocols.modbus.readwrite.ModbusPDUReadDiscreteInputsRequest import (
    ModbusPDUReadDiscreteInputsRequestBuilder,
)
from plc4py.utils.GenericTypes import ByteOrder


async def test_modbus_discreate_inputs_request():
    request = ModbusPDUReadDiscreteInputsRequestBuilder(0, 10).build()
    assert request is not None


async def test_modbus_discreate_inputs_request():
    request = ModbusPDUReadDiscreteInputsRequestBuilder(0, 10).build()
    size = request.length_in_bytes()
    write_buffer = WriteBufferByteBased(size, ByteOrder.LITTLE_ENDIAN)
    serialize = request.serialize_modbus_pdu_child(write_buffer)

    assert request is not None