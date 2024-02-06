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
from typing import Dict, List

from plc4py.api.exceptions.exceptions import PlcDataTypeNotFoundException
from plc4py.protocols.umas.readwrite.UmasDatatypeReference import UmasDatatypeReference

from plc4py.protocols.umas.readwrite.UmasUnlocatedVariableReference import UmasUnlocatedVariableReference


@dataclass
class UmasVariable:
    variable_name: str
    data_type: int
    element_size: int


@dataclass
class UmasElementryVariable(UmasVariable):
    pass


@dataclass
class UmasCustomVariable(UmasVariable):
    children: Dict[str, UmasVariable]


@dataclass
class UmasVariableBuilder:
    tag_references: Dict[str, UmasUnlocatedVariableReference]
    data_type_references: List[UmasDatatypeReference]

    def build(self) -> Dict[str, UmasVariable]:
        return_dict: Dict[str, UmasVariable] = {}
        for tag_name_key, tag_reference in self.tag_references.items():
            data_type = tag_reference.data_type
            if data_type < 26: # Start of the custom data types
                return_dict[tag_name_key] = UmasElementryVariable(tag_name_key, data_type, 1)
            else:
                found_data_type = False
                for data_type_reference in self.data_type_references:
                    if data_type_reference.data_type == data_type:
                        return_dict[tag_name_key] = UmasCustomVariableBuilder(tag_reference, data_type_reference).build()
                        found_data_type = True
                        break
                if not found_data_type:
                    raise PlcDataTypeNotFoundException(f"Could not find data type {data_type} for tag {tag_name_key}")

        return return_dict


@dataclass
class UmasCustomVariableBuilder:
    tag_reference: UmasUnlocatedVariableReference
    data_type_reference: UmasDatatypeReference

    def build(self) -> UmasCustomVariable:
        return UmasCustomVariable("", 1, 1, {})

