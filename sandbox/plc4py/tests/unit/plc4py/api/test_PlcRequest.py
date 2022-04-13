#
# Licensed to the Apache Software Foundation (ASF) under one
# or more contributor license agreements.  See the NOTICE file
# distributed with this work for additional information
# regarding copyright ownership.  The ASF licenses this file
# to you under the Apache License, Version 2.0 (the
# "License"); you may not use this file except in compliance
# with the License.  You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing,
# software distributed under the License is distributed on an
# "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
# KIND, either express or implied.  See the License for the
# specific language governing permissions and limitations
# under the License.
#

from plc4py.api.PlcConnection import PlcConnection
from plc4py.api.messages.PlcRequest import (
    PlcFieldRequest,
)
from tests.unit.plc4py.api.test.MockPlcConection import MockPlcConnection


def test_read_request_builder_empty_request(mocker) -> None:
    connection: PlcConnection = MockPlcConnection()

    # the connection function is supposed to support context manager
    # so using it in a with statement should result in close being called on the connection
    with connection.read_request_builder() as builder:
        request: PlcFieldRequest = builder.build()
    assert len(request.field_names) == 0


def test_read_request_builder_non_empty_request(mocker) -> None:
    connection: PlcConnection = MockPlcConnection()

    # the connection function is supposed to support context manager
    # so using it in a with statement should result in close being called on the connection
    with connection.read_request_builder() as builder:
        builder.add_item("1:BOOL")
        request: PlcFieldRequest = builder.build()

    # verify that request has one field
    assert request.field_names == ["1:BOOL"]
    assert len(request.field_names) == 1
