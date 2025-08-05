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

from plc4py.api.messages.PlcRequest import PlcWriteRequest
from plc4py.api.messages.PlcResponse import PlcWriteResponse
from plc4py.api.value.PlcValue import PlcResponseCode


class PlcWriter:
    """
    Interface implemented by all PlcConnections that are able to write from remote resources.
    """

    @abstractmethod
    async def _write(self, request: PlcWriteRequest) -> PlcWriteResponse:
        """
        Executes a PlcWriteRequest

        This method sends a write request to the connected Modbus device and waits for a response.
        The response is then returned as a PlcWriteResponse.

        If no device is set, an error is logged and a PlcWriteResponse with the
        PlcResponseCode.NOT_CONNECTED code is returned.
        If an error occurs during the execution of the write request, a
        PlcWriteResponse with the PlcResponseCode.INTERNAL_ERROR code is returned.

        :param request: PlcWriteRequest to execute
        :return: PlcWriteResponse
        """
        pass

    @abstractmethod
    def is_write_supported(self) -> bool:
        """
        Indicates if the connection supports write requests.
        :return: True if connection supports writing, False otherwise
        """
        pass


class DefaultPlcWriter(PlcWriter):
    """
    Interface implemented by all PlcConnections that are able to write from remote resources.
    """

    def __init__(self):
        self._transport = None
        self._device = None

    async def _write(self, request: PlcWriteRequest) -> PlcWriteResponse:
        """
        Executes a PlcWriteRequest

        This method sends a write request to the connected Modbus device and waits for a response.
        The response is then returned as a PlcWriteResponse.

        If no device is set, an error is logged and a PlcWriteResponse with the
        PlcResponseCode.NOT_CONNECTED code is returned.
        If an error occurs during the execution of the write request, a
        PlcWriteResponse with the PlcResponseCode.INTERNAL_ERROR code is returned.

        :param request: PlcWriteRequest to execute
        :return: PlcWriteResponse
        """

        try:
            # Send the write request to the device and wait for a response
            logging.debug("Sending write request to Device")
            response = await asyncio.wait_for(
                self._device.write(request, self._transport), 15
            )
            # Return the response
            return response
        except Exception as e:
            # If an error occurs during the execution of the write request, return a response with
            # Still haven't found a nice way to return an error
            # the INTERNAL_ERROR code. This exception is very general and probably should be replaced.
            # TODO:- This exception is very general and probably should be replaced
            raise e

    def is_write_supported(self) -> bool:
        """
        Indicates if the connection supports write requests.
        :return: True if connection supports writing, False otherwise
        """
        return True
