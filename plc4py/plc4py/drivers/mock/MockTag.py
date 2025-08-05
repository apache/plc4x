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

from plc4py.api.messages.PlcField import PlcTag
from plc4py.spi.messages.PlcRequest import TagBuilder


@dataclass
class MockTag(PlcTag):
    address: str
    data_type: str


class MockTagBuilder(TagBuilder):
    @staticmethod
    def create(address_string: str) -> MockTag:
        address, data_type = address_string.split(":")
        return MockTag(address, data_type)
