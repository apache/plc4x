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

from __future__ import annotations

from asyncio import Transport, BaseTransport, BaseProtocol
from dataclasses import dataclass


@dataclass
class Plc4xBaseTransport(Transport):
    """
    Wrapper for the Python Transport. Instead of using the generic asyncio.create_connection.
    Using a custom transport allows for it to be modified if required.
    """

    _transport: BaseTransport
    _protocol: BaseProtocol

    @property
    def protocol(self) -> BaseProtocol:
        return self._protocol

    @protocol.setter
    def protocol(self, protocol: BaseProtocol):
        self._protocol = protocol

    def is_closing(self):
        if self._transport is not None:
            return self._transport.is_closing()

    # def get_write_buffer_limits(self):
    #    return self._write_transport.get_write_buffer_limits()

    def close(self):
        if self._transport is not None:
            self._transport.close()
