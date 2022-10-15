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
from abc import abstractmethod
from dataclasses import dataclass

from plc4py.api.PlcConnection import PlcConnection
from plc4py.api.authentication.PlcAuthentication import PlcAuthentication
from plc4py.api.exceptions.exceptions import PlcNotImplementedException
from plc4py.api.messages.PlcDiscovery import PlcDiscoveryRequestBuilder


@dataclass
class PlcDriverMetaData:
    """
    Information about driver capabilities
    """

    """Indicates that the driver supports discovery"""
    can_discover: bool = False


class PlcDriver:
    """
    General interface defining the minimal methods required for adding a new type of driver to the PLC4PY system.

    <b>Note that each driver has to add a setuptools entrypoint as plc4x.driver in order to be loaded by pluggy</b>
    """

    def __init__(self):
        self.protocol_code: str
        self.protocol_name: str

    @property
    def metadata(self):
        return PlcDriverMetaData()

    @abstractmethod
    async def get_connection(
        self, url: str, authentication: PlcAuthentication = PlcAuthentication()
    ) -> PlcConnection:
        """
        Connects to a PLC using the given plc connection string.
        :param url: plc connection string
        :param authentication: authentication credentials.
        :return PlcConnection: PLC Connection object
        """
        pass

    def discovery_request_builder(self) -> PlcDiscoveryRequestBuilder:
        """
        Discovery Request Builder aids in generating a discovery request for this protocol
        :return builder: Discovery request builder
        """
        raise PlcNotImplementedException(f"Not implemented for {self.protocol_name}")
