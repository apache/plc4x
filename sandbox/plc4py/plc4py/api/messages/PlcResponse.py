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
from typing import cast, List, Dict

from plc4py.api.messages.PlcField import PlcField
from plc4py.api.messages.PlcMessage import PlcMessage
from plc4py.api.value.PlcValue import PlcValue, PlcResponseCode
from plc4py.spi.messages.utils.ResponseItem import ResponseItem


@dataclass
class PlcResponse(PlcMessage):
    """
    Base type for all response messages sent as response for a prior request
    from a plc to the plc4x system.
    """

    code: PlcResponseCode


@dataclass
class PlcFieldResponse(PlcResponse):
    fields: List[PlcField]

    @property
    def field_names(self):
        return [fld.name for fld in self.fields]

    def response_code(self, name: str) -> PlcResponseCode:
        pass


@dataclass
class PlcReadResponse(PlcFieldResponse):
    """
    Response to a {@link PlcReadRequest}.
    """

    values: Dict[str, List[ResponseItem[PlcValue]]]

    def get_plc_value(self, name: str, index: int = 0) -> PlcValue:
        return self.values[name][index].value

    def number_of_values(self, name: str) -> int:
        return len(self.values[name])

    def is_boolean(self, name: str, index: int = 0):
        return isinstance(self.values[name][index].value.value, bool)

    def get_boolean(self, name: str, index: int = 0) -> bool:
        return cast(bool, self.values[name][index].value.value)

    def is_int(self, name: str, index: int = 0):
        return isinstance(self.values[name][index].value.value, int)

    def get_int(self, name: str, index: int = 0) -> int:
        return cast(int, self.values[name][index].value.value)
