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
from abc import abstractmethod
from typing import Awaitable, Dict, Union

from plc4py.api.messages.PlcResponse import (
    PlcResponse,
    PlcTagResponse,
    PlcReadResponse,
    PlcWriteResponse,
    PlcBrowseResponse,
)
from plc4py.api.messages.PlcRequest import ReadRequestBuilder, PlcRequest
from plc4py.api.value.PlcValue import PlcResponseCode, PlcValue
from plc4py.spi.configuration.PlcConfiguration import PlcConfiguration
from plc4py.spi.messages.utils.ResponseItem import ResponseItem
from plc4py.utils.GenericTypes import GenericGenerator


class PlcConnection(GenericGenerator):
    def __init__(self, config: PlcConfiguration):
        self._configuration = config

    @abstractmethod
    def is_connected(self) -> bool:
        """
        Indicates if the connection is established to a remote PLC.
        :return: True if connection, False otherwise
        """
        pass

    @abstractmethod
    def close(self) -> None:
        """
        Closes the connection to the remote PLC.
        :return:
        """
        pass

    @abstractmethod
    def read_request_builder(self) -> ReadRequestBuilder:
        """
        :return: read request builder.
        """
        pass

    @abstractmethod
    def execute(self, request: PlcRequest) -> Awaitable[PlcResponse]:
        """
        Executes a PlcRequest as long as it's already connected
        :param request: Plc Request to execute
        :return: The response from the Plc/Device
        """
        pass

    def _default_failed_request(
        self, code: PlcResponseCode
    ) -> Union[PlcReadResponse, PlcWriteResponse, PlcBrowseResponse]:
        """
        Returns a default PlcResponse, mainly used in case of a failed request
        :param code: The response code to return
        :return: The PlcResponse
        """
        return PlcResponse(code)


class PlcConnectionMetaData:

    @abstractmethod
    def is_read_supported(self) -> bool:
        """
        Indicates if the connection supports read requests.
        :return: True if connection supports reading, False otherwise
        """
        pass

    @abstractmethod
    def is_write_supported(self) -> bool:
        """
        Indicates if the connection supports write requests.
        :return: True if connection supports writing, False otherwise
        """
        pass

    @abstractmethod
    def is_subscribe_supported(self) -> bool:
        """
        Indicates if the connection supports subscription requests.
        :return: True if connection supports subscriptions, False otherwise
        """
        pass

    @abstractmethod
    def is_browse_supported(self) -> bool:
        """
        Indicates if the connection supports browsing requests.
        :return: True if connection supports browsing, False otherwise
        """
        pass
