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
import re
from dataclasses import dataclass
from typing import Dict, List, Pattern, AnyStr, Union

from plc4py.protocols.umas.readwrite.UmasDataType import UmasDataType

from plc4py.protocols.umas.readwrite.VariableRequestReference import (
    VariableRequestReference,
)

from plc4py.api.exceptions.exceptions import PlcDataTypeNotFoundException
from plc4py.protocols.umas.readwrite.UmasDatatypeReference import UmasDatatypeReference
from plc4py.protocols.umas.readwrite.UmasUDTDefinition import UmasUDTDefinition

from plc4py.protocols.umas.readwrite.UmasUnlocatedVariableReference import (
    UmasUnlocatedVariableReference,
)


@dataclass
class UmasVariable:
    variable_name: str
    data_type: int
    block_no: int
    offset: int

    def get_variable_reference(self, address: str) -> VariableRequestReference:
        raise NotImplementedError(
            f"UmasVariable subclass not implemented for variable {self.variable_name}"
        )

    def get_byte_length(self) -> int:
        raise NotImplementedError(
            f"UmasVariable subclass not implemented for variable {self.variable_name}"
        )


@dataclass
class UmasElementryVariable(UmasVariable):
    def get_variable_reference(self, address: str) -> VariableRequestReference:
        if self.data_type == UmasDataType.STRING.value:
            return VariableRequestReference(
                is_array=1,
                data_size_index=UmasDataType(self.data_type).request_size,
                block=self.block_no,
                base_offset=0x0000,
                offset=self.offset,
                array_length=16,
            )
        else:
            return VariableRequestReference(
                is_array=0,
                data_size_index=UmasDataType(self.data_type).request_size,
                block=self.block_no,
                base_offset=0x0000,
                offset=self.offset,
                array_length=None,
            )

    def get_byte_length(self) -> int:
        return 7


@dataclass
class UmasCustomVariable(UmasVariable):
    children: Dict[str, UmasVariable]

    def get_variable_reference(self, address: str) -> VariableRequestReference:
        split_tag_address: List[str] = address.split(".")
        child_index = None
        if len(split_tag_address) > 1:
            child_index = split_tag_address[1]
            return self.children[child_index].get_variable_reference(
                ".".join(split_tag_address[1:])
            )
        else:
            raise NotImplementedError("Unable to read structures of UDT's")

    def get_byte_length(self) -> int:
        byte_count = 0
        for key, child in self.children.items():
            byte_count += child.get_byte_length()
        return byte_count


@dataclass
class UmasArrayVariable(UmasVariable):
    start_index: int
    end_index: int

    def get_variable_reference(self, address: str) -> VariableRequestReference:
        split_tag_address: List[str] = address.split(".")
        address_index = None
        if len(split_tag_address) > 1:
            address_index = int(split_tag_address[1])
        data_type_enum = UmasDataType(self.data_type)
        if address_index:
            return VariableRequestReference(
                is_array=0,
                data_size_index=data_type_enum.request_size,
                block=self.block_no,
                base_offset=0x0000,
                offset=self.offset
                + (address_index - self.start_index) * data_type_enum.data_type_size,
                array_length=None,
            )
        else:
            return VariableRequestReference(
                is_array=1,
                data_size_index=data_type_enum.request_size,
                block=self.block_no,
                base_offset=0x0000,
                offset=self.offset,
                array_length=self.end_index - self.start_index + 1,
            )

    def get_byte_length(self) -> int:
        return 9


@dataclass
class UmasVariableBuilder:
    tag_name: str
    tag_reference: Union[UmasUnlocatedVariableReference, UmasUDTDefinition]
    data_type_references: List[UmasDatatypeReference]
    udt_definitions: Dict[str, List[UmasUDTDefinition]]
    block: int = -1
    offset: int = 0

    def build(self) -> UmasVariable:
        variable: UmasVariable = None
        _ARRAY_REGEX: str = (
            "^ARRAY\[(?P<start_number>[0-9]*)..(?P<end_number>[0-9]*)\] OF (?P<data_type>[a-zA-z0-9]*)"
        )
        _ARRAY_COMPILED: Pattern[AnyStr] = re.compile(_ARRAY_REGEX)

        if self.block == -1:
            self.block = self.tag_reference.block

        data_type = self.tag_reference.data_type
        if data_type < 26:  # Start of the custom data types
            variable = UmasElementryVariable(
                self.tag_name,
                data_type,
                self.block,
                self.tag_reference.offset + self.offset,
            )
        else:
            found_data_type = False
            for data_type_reference in self.data_type_references:
                if data_type_reference.data_type == data_type:
                    if data_type_reference.class_identifier == 2:
                        custom_children: Dict[str, UmasUDTDefinition] = {
                            definition.value: definition
                            for definition in self.udt_definitions[
                                data_type_reference.value
                            ]
                        }
                        children: Dict[str, UmasVariable] = {}
                        for kea, child in custom_children.items():
                            children[kea] = UmasVariableBuilder(
                                child.value,
                                child,
                                self.data_type_references,
                                self.udt_definitions,
                                offset=self.tag_reference.offset,
                                block=self.block,
                            ).build()
                        variable = UmasCustomVariable(
                            self.tag_name,
                            data_type,
                            self.tag_reference.block,
                            self.tag_reference.offset,
                            children,
                        )
                    elif data_type_reference.class_identifier == 4:
                        match = _ARRAY_COMPILED.match(data_type_reference.value)
                        data_type = UmasDataType[match.group("data_type")]
                        variable = UmasArrayVariable(
                            self.tag_reference.value,
                            data_type.value,
                            self.block,
                            self.tag_reference.offset + self.offset,
                            int(match.group("start_number")),
                            int(match.group("end_number")),
                        )
                    found_data_type = True
                    break
            if not found_data_type:
                raise PlcDataTypeNotFoundException(
                    f"Could not find data type {data_type} for tag {self.tag_name}"
                )

        return variable
