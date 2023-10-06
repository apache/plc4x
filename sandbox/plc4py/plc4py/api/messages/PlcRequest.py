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
from dataclasses import dataclass, field
from typing import Union, List

from plc4py.api.messages.PlcField import PlcField
from plc4py.api.messages.PlcMessage import PlcMessage
from plc4py.utils.GenericTypes import GenericGenerator


class PlcRequest(PlcMessage):
    """
    Base type for all messages sent from the plc4x system to a connected plc.
    """


@dataclass
class PlcFieldRequest(PlcRequest):
    fields: List[PlcField] = field(default_factory=lambda: [])

    @property
    def field_names(self):
        return [field.name for field in self.fields]


@dataclass
class PlcReadRequest(PlcFieldRequest):
    """
    Base type for all messages sent from the plc4x system to a connected plc.
    """


class ReadRequestBuilder(GenericGenerator):
    @abstractmethod
    def build(self) -> PlcReadRequest:
        pass

    @abstractmethod
    def add_item(self, field_query: Union[str, PlcField]) -> None:
        pass
