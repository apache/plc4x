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
from abc import abstractmethod

from plc4py.api.messages.PlcRequest import PlcReadRequest
from plc4py.api.messages.PlcResponse import PlcReadResponse
from plc4py.api.value.PlcValue import PlcResponseCode


class PlcReader:
    """
    Interface implemented by all PlcConnections that are able to read from remote resources.
    """

    @abstractmethod
    async def _read(self, request: PlcReadRequest) -> PlcReadResponse:
        """
        Executes a PlcReadRequest

        This method sends a read request to the connected device and waits for a response.
        The response is then returned as a PlcReadResponse.

        If no device is set, an error is logged and a PlcResponseCode.NOT_CONNECTED is returned.
        If an error occurs during the execution of the read request, a PlcResponseCode.INTERNAL_ERROR is
        returned.

        :param request: PlcReadRequest to execute
        :return: PlcReadResponse
        """
        pass

    @abstractmethod
    def is_read_supported(self) -> bool:
        """
        Indicates if the connection supports read requests.
        :return: True if connection supports reading, False otherwise
        """
        pass


class DefaultPlcReader(PlcReader):
    """
    Interface implemented by all PlcConnections that are able to read from remote resources.
    """

    def __init__(self):
        self._transport = None
        self._device = None

    async def _read(self, request: PlcReadRequest) -> PlcReadResponse:
        """
        Executes a PlcReadRequest

        This method sends a read request to the connected device and waits for a response.
        The response is then returned as a PlcReadResponse.

        If no device is set, an error is logged and a PlcResponseCode.NOT_CONNECTED is returned.
        If an error occurs during the execution of the read request, a PlcResponseCode.INTERNAL_ERROR is
        returned.

        :param request: PlcReadRequest to execute
        :return: PlcReadResponse
        """

        # TODO: Insert Optimizer base on data from a browse request
        try:
            logging.debug("Sending read request to Device")
            response = await asyncio.wait_for(
                self._device.read(request, self._transport), 10
            )
            return response
        except Exception as e:
            # TODO:- This exception is very general and probably should be replaced
            raise e

    def is_read_supported(self) -> bool:
        """
        Indicates if the connection supports read requests.
        :return: True if connection supports reading, False otherwise
        """
        return True
