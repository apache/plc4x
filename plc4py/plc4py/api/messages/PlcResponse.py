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
from dataclasses import dataclass
from typing import Dict, List, cast

from plc4py.api.messages.PlcField import PlcTag
from plc4py.api.messages.PlcMessage import PlcMessage
from plc4py.api.value.PlcValue import PlcResponseCode, PlcValue
from plc4py.spi.messages.utils.ResponseItem import ResponseItem, PlcBrowseItem


@dataclass
class PlcResponse(PlcMessage):
    """
    Base type for all response messages sent as response for a prior request
    from a plc to the plc4x system.
    """

    response_code: PlcResponseCode


@dataclass
class PlcTagResponse(PlcResponse):
    tags: Dict[str, ResponseItem[PlcValue]]

    @property
    def tag_names(self):
        return [tag_name for tag_name in self.tags.keys()]


@dataclass
class PlcReadResponse(PlcTagResponse):
    """
    Response to a {@link PlcReadRequest}.
    """

    def get_plc_value(self, name: str) -> PlcValue:
        return self.tags[name].value

    def is_boolean(self, name: str):
        return isinstance(self.tags[name].value.value, bool)

    def get_boolean(self, name: str) -> bool:
        return cast(bool, self.tags[name].value.value)

    def is_int(self, name: str):
        return isinstance(self.tags[name].value.value, int)

    def get_int(self, name: str) -> int:
        return cast(int, self.tags[name].value.value)


@dataclass
class PlcWriteResponse(PlcTagResponse):
    """
    Response to a {@link PlcWriteRequest}.
    """


@dataclass
class PlcQueryResponse(PlcResponse):
    tags: Dict[str, List[PlcBrowseItem[PlcTag]]]

    @property
    def tag_names(self):
        return [tag_name for tag_name in self.tags.keys()]


@dataclass
class PlcBrowseResponse(PlcQueryResponse):

    def get_tags(self) -> Dict[str, List[PlcBrowseItem[PlcTag]]]:
        return self.tags
