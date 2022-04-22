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
from unittest.mock import MagicMock

from plc4py.PlcDriverManager import PlcDriverManager
from plc4py.drivers.mock.MockConnection import MockConnection


def test_connection_context_manager_impl_close_called(mocker) -> None:
    manager: PlcDriverManager = PlcDriverManager()

    # getup a plain return value for get_connection
    connection_mock: MagicMock = mocker.patch.object(manager, "get_connection")
    connection_mock.return_value = MockConnection()

    # the connection function is supposed to support context manager
    # so using it in a with statement should result in close being called on the connection
    with manager.connection("foo://bar") as conn:
        close_mock: MagicMock = mocker.patch.object(conn, "close")
        print(conn.__doc__)

    # verify that close was called by the context manager
    close_mock.assert_called()
