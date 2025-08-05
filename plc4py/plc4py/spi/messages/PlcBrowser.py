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

from plc4py.api.messages.PlcRequest import PlcBrowseRequest
from plc4py.api.messages.PlcResponse import PlcBrowseResponse
from plc4py.api.value.PlcValue import PlcResponseCode


class PlcBrowser:
    """
    Interface implemented by all PlcConnections that are able to read from remote resources.
    """

    @abstractmethod
    async def _browse(self, request: PlcBrowseRequest) -> PlcBrowseResponse:
        """
        Executes a PlcBrowseRequest
        """
        pass

    @abstractmethod
    def is_browse_supported(self) -> bool:
        """
        Indicates if the connection supports browsing requests.
        :return: True if connection supports browsing, False otherwise
        """
        pass


class DefaultPlcBrowser(PlcBrowser):
    """
    Interface implemented by all PlcConnections that are able to read from remote resources.
    """

    def __init__(self):
        self._transport = None
        self._device = None

    async def _browse(self, request: PlcBrowseRequest) -> PlcBrowseResponse:
        """
        Executes a PlcBrowseRequest
        """

        try:
            # Send the browse request to the device and wait for a response
            logging.debug("Sending browse request to Device")
            response = await asyncio.wait_for(
                self._device.browse(request, self._transport), 5
            )
            # Return the response
            return response
        except Exception as e:
            # TODO:- This exception is very general and probably should be replaced
            raise e

    def is_browse_supported(self) -> bool:
        """
        Indicates if the connection supports browsing requests.
        :return: True if connection supports browsing, False otherwise
        """
        return True
