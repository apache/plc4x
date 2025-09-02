#  Licensed to the Apache Software Foundation (ASF) under one
#  or more contributor license agreements.  See the NOTICE file
#  distributed with this work for additional information
#  regarding copyright ownership.  The ASF licenses this file
#  to you under the Apache License, Version 2.0 (the
#  "License"); you may not use this file except in compliance
#  with the License.  You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
#  Unless required by applicable law or agreed to in writing,
#  software distributed under the License is distributed on an
#  "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
#  KIND, either express or implied.  See the License for the
#  specific language governing permissions and limitations
#  under the License.

from dataclasses import dataclass, field
from enum import Enum
from typing import List, Optional

__NAMESPACE__ = "https://plc4x.apache.org/schemas/parser-serializer-testsuite.xsd"


class ByteOrder(Enum):
    BIG_ENDIAN = "BIG_ENDIAN"
    LITTLE_ENDIAN = "LITTLE_ENDIAN"


@dataclass
class ParameterList:
    class Meta:
        name = "parameterList"

    parameter: List["ParameterList.Parameter"] = field(
        default_factory=list,
        metadata={
            "type": "Element",
            "namespace": "",
            "min_occurs": 1,
        },
    )

    @dataclass
    class Parameter:
        name: Optional[str] = field(
            default=None,
            metadata={
                "type": "Element",
                "namespace": "",
                "required": True,
            },
        )
        value: Optional[str] = field(
            default=None,
            metadata={
                "type": "Element",
                "namespace": "",
            },
        )


@dataclass
class Testsuite:
    """
    :ivar name:
    :ivar protocol_name:
    :ivar output_flavor:
    :ivar options: List of options which are specific to execution of
        test or environment. This might be used ie. to influence test
        framework lookup strategies.
    :ivar testcase:
    :ivar byte_order:
    """

    class Meta:
        name = "testsuite"
        namespace = "https://plc4x.apache.org/schemas/parser-serializer-testsuite.xsd"

    name: Optional[str] = field(
        default=None,
        metadata={
            "type": "Element",
            "namespace": "",
            "required": True,
        },
    )
    protocol_name: Optional[str] = field(
        default=None,
        metadata={
            "name": "protocolName",
            "type": "Element",
            "namespace": "",
            "required": True,
        },
    )
    output_flavor: Optional[str] = field(
        default=None,
        metadata={
            "name": "outputFlavor",
            "type": "Element",
            "namespace": "",
            "required": True,
        },
    )
    options: Optional[ParameterList] = field(
        default=None,
        metadata={
            "type": "Element",
            "namespace": "",
        },
    )
    testcase: List["Testsuite.Testcase"] = field(
        default_factory=list,
        metadata={
            "type": "Element",
            "namespace": "",
            "min_occurs": 1,
        },
    )
    byte_order: Optional[ByteOrder] = field(
        default=None,
        metadata={
            "name": "byteOrder",
            "type": "Attribute",
        },
    )

    @dataclass
    class Testcase:
        name: Optional[str] = field(
            default=None,
            metadata={
                "type": "Element",
                "namespace": "",
                "required": True,
            },
        )
        description: Optional[str] = field(
            default=None,
            metadata={
                "type": "Element",
                "namespace": "",
            },
        )
        raw: Optional[bytes] = field(
            default=None,
            metadata={
                "type": "Element",
                "namespace": "",
                "required": True,
                "format": "base16",
            },
        )
        root_type: Optional[str] = field(
            default=None,
            metadata={
                "name": "root-type",
                "type": "Element",
                "namespace": "",
                "required": True,
            },
        )
        parser_arguments: Optional["Testsuite.Testcase.ParserArguments"] = field(
            default=None,
            metadata={
                "name": "parser-arguments",
                "type": "Element",
                "namespace": "",
            },
        )
        xml: Optional[object] = field(
            default=None,
            metadata={
                "type": "Element",
                "namespace": "",
            },
        )

        @dataclass
        class ParserArguments:
            local_element: List[object] = field(
                default_factory=list,
                metadata={
                    "type": "Wildcard",
                    "namespace": "##local",
                },
            )
