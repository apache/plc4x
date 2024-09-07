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
from typing import List, Any

from plc4py.api.messages.PlcField import PlcTag
from plc4py.api.messages.PlcRequest import (
    BrowseRequestBuilder,
    PlcBrowseRequest,
    PlcReadRequest,
    PlcWriteRequest,
    ReadRequestBuilder,
    WriteRequestBuilder,
)
from plc4py.api.value.PlcValue import PlcValue


class TagBuilder:
    @staticmethod
    def create(address_string: str):
        raise NotImplementedError


class DefaultReadRequestBuilder(ReadRequestBuilder):
    def __init__(self, tag_builder: TagBuilder):
        super().__init__()
        self.read_request = PlcReadRequest()
        self.tag_builder = tag_builder

    def build(self) -> PlcReadRequest:
        return self.read_request

    def add_item(self, tag_name: str, address_string: str) -> None:
        tag = self.tag_builder.create(address_string)
        self.read_request.tags[tag_name] = tag

    def add_tag(self, tag_name: str, tag: PlcTag) -> None:
        self.read_request.tags[tag_name] = tag


class DefaultWriteRequestBuilder(WriteRequestBuilder):
    def __init__(self, tag_builder: TagBuilder):
        super().__init__()
        self.write_request = PlcWriteRequest()
        self.tag_builder = tag_builder

    def build(self) -> PlcWriteRequest:
        return self.write_request

    def add_item(self, tag_name: str, address_string: str, value: PlcValue) -> None:
        tag = self.tag_builder.create(address_string)
        self.write_request.tags[tag_name] = tag
        self.write_request.values[tag_name] = value

    def add_tag(self, tag_name: str, tag: PlcTag, value: PlcValue) -> None:
        self.write_request.tags[tag_name] = tag


class DefaultBrowseRequestBuilder(BrowseRequestBuilder):
    def __init__(self):
        super().__init__()
        self.browse_request = PlcBrowseRequest()

    def build(self) -> PlcBrowseRequest:
        return self.browse_request

    def add_query(self, query_name: str, query: str) -> None:
        self.browse_request.queries[query_name] = query
