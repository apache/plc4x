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
from collections import OrderedDict
from dataclasses import dataclass, field
from typing import Dict, List, Union

from plc4py.api.messages.PlcField import PlcTag
from plc4py.api.messages.PlcMessage import PlcMessage
from plc4py.api.value.PlcValue import PlcValue
from plc4py.utils.GenericTypes import GenericGenerator


class PlcRequest(PlcMessage):
    """
    Base type for all messages sent from the plc4x system to a connected plc.
    """


@dataclass
class PlcTagRequest(PlcRequest):
    tags: Dict[str, PlcTag] = field(default_factory=lambda: OrderedDict())

    @property
    def tag_names(self):
        return [tag_name for tag_name in self.tags.keys()]


@dataclass
class PlcQueryRequest(PlcRequest):
    queries: Dict[str, str] = field(default_factory=lambda: OrderedDict())

    @property
    def query_names(self):
        return [quary_name for quary_name in self.queries.keys()]


@dataclass
class PlcReadRequest(PlcTagRequest):
    """
    Base type for all read messages sent from the plc4x system to a connected plc.
    """


@dataclass
class PlcWriteRequest(PlcTagRequest):
    """
    Base type for all write messages sent from the plc4x system to a connected plc.
    """

    values: Dict[str, PlcValue] = field(default_factory=lambda: OrderedDict())


@dataclass
class PlcBrowseRequest(PlcQueryRequest):
    """
    Base type for all messages sent from the plc4x system to a connected plc.
    """


class ReadRequestBuilder(GenericGenerator):
    @abstractmethod
    def build(self) -> PlcReadRequest:
        pass

    @abstractmethod
    def add_item(self, tag_name: str, address_string: str) -> None:
        pass


class WriteRequestBuilder(GenericGenerator):
    @abstractmethod
    def build(self) -> PlcWriteRequest:
        pass

    @abstractmethod
    def add_item(self, tag_name: str, address_string: str, value: PlcValue) -> None:
        pass


class BrowseRequestBuilder(GenericGenerator):
    @abstractmethod
    def build(self) -> PlcBrowseRequest:
        pass

    @abstractmethod
    def add_query(self, query_name: str, query: str) -> None:
        pass
