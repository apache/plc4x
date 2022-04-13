# regarding copyright ownership.  The ASF licenses this file
# to you under the Apache License, Version 2.0 (the
# "License"); you may not use this file except in compliance
# with the License.  You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing,
# software distributed under the License is distributed on an
# "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
# KIND, either express or implied.  See the License for the
# specific language governing permissions and limitations
# under the License.
#
from dataclasses import field, dataclass

from plc4py.api.messages.PlcField import PlcField
from plc4py.api.messages.PlcMessage import PlcMessage


class PlcResponseCode:
    pass


class PlcResponse(PlcMessage):
    """
    Base type for all response messages sent as response for a prior request
    from a plc to the plc4x system.
    """

    pass


@dataclass
class PlcFieldResponse(PlcResponse):
    fields: list[PlcField] = field(default_factory=lambda: [])

    @property
    def field_names(self):
        return [field.name for field in self.fields]

    def response_code(self, name: str) -> PlcResponseCode:
        pass


class PlcReadResponse(PlcFieldResponse):
    """
    Response to a {@link PlcReadRequest}.
    """

    pass
