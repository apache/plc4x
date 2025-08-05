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

__NAMESPACE__ = "https://plc4x.apache.org/schemas/driver-testsuite.xsd"


class ByteOrder(Enum):
    BIG_ENDIAN = "BIG_ENDIAN"
    LITTLE_ENDIAN = "LITTLE_ENDIAN"


@dataclass
class BytesStep:
    class Meta:
        name = "bytesStep"

    value: Optional[bytes] = field(
        default=None,
        metadata={
            "required": True,
            "format": "base16",
        },
    )
    name: Optional[str] = field(
        default=None,
        metadata={
            "type": "Attribute",
        },
    )


@dataclass
class MessageStep:
    class Meta:
        name = "messageStep"

    any_element: Optional[object] = field(
        default=None,
        metadata={
            "type": "Wildcard",
            "namespace": "##any",
        },
    )
    name: Optional[str] = field(
        default=None,
        metadata={
            "type": "Attribute",
        },
    )


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
class DriverTestsuite:
    """
    :ivar name:
    :ivar protocol_name:
    :ivar output_flavor:
    :ivar options: List of options which are specific to execution of
        test or environment. This might be used ie. to influence test
        framework lookup strategies.
    :ivar driver_name:
    :ivar driver_parameters:
    :ivar setup:
    :ivar teardown:
    :ivar testcase:
    :ivar byte_order:
    """

    class Meta:
        name = "driver-testsuite"
        namespace = "https://plc4x.apache.org/schemas/driver-testsuite.xsd"

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
    driver_name: Optional[str] = field(
        default=None,
        metadata={
            "name": "driver-name",
            "type": "Element",
            "namespace": "",
            "required": True,
        },
    )
    driver_parameters: Optional[ParameterList] = field(
        default=None,
        metadata={
            "name": "driver-parameters",
            "type": "Element",
            "namespace": "",
        },
    )
    setup: Optional["DriverTestsuite.Setup"] = field(
        default=None,
        metadata={
            "type": "Element",
            "namespace": "",
        },
    )
    teardown: Optional["DriverTestsuite.Teardown"] = field(
        default=None,
        metadata={
            "type": "Element",
            "namespace": "",
        },
    )
    testcase: List["DriverTestsuite.Testcase"] = field(
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
    class Setup:
        outgoing_plc_message: List[MessageStep] = field(
            default_factory=list,
            metadata={
                "name": "outgoing-plc-message",
                "type": "Element",
                "namespace": "",
                "sequence": 1,
            },
        )
        outgoing_plc_bytes: List[BytesStep] = field(
            default_factory=list,
            metadata={
                "name": "outgoing-plc-bytes",
                "type": "Element",
                "namespace": "",
                "sequence": 1,
            },
        )
        incoming_plc_message: List[MessageStep] = field(
            default_factory=list,
            metadata={
                "name": "incoming-plc-message",
                "type": "Element",
                "namespace": "",
                "sequence": 1,
            },
        )
        incoming_plc_bytes: List[BytesStep] = field(
            default_factory=list,
            metadata={
                "name": "incoming-plc-bytes",
                "type": "Element",
                "namespace": "",
                "sequence": 1,
            },
        )
        api_request: List[MessageStep] = field(
            default_factory=list,
            metadata={
                "name": "api-request",
                "type": "Element",
                "namespace": "",
                "sequence": 1,
            },
        )
        api_response: List[MessageStep] = field(
            default_factory=list,
            metadata={
                "name": "api-response",
                "type": "Element",
                "namespace": "",
                "sequence": 1,
            },
        )
        delay: List[int] = field(
            default_factory=list,
            metadata={
                "type": "Element",
                "namespace": "",
                "sequence": 1,
            },
        )
        terminate: List[object] = field(
            default_factory=list,
            metadata={
                "type": "Element",
                "namespace": "",
                "sequence": 1,
            },
        )

    @dataclass
    class Teardown:
        outgoing_plc_message: List[MessageStep] = field(
            default_factory=list,
            metadata={
                "name": "outgoing-plc-message",
                "type": "Element",
                "namespace": "",
                "sequence": 1,
            },
        )
        outgoing_plc_bytes: List[BytesStep] = field(
            default_factory=list,
            metadata={
                "name": "outgoing-plc-bytes",
                "type": "Element",
                "namespace": "",
                "sequence": 1,
            },
        )
        incoming_plc_message: List[MessageStep] = field(
            default_factory=list,
            metadata={
                "name": "incoming-plc-message",
                "type": "Element",
                "namespace": "",
                "sequence": 1,
            },
        )
        incoming_plc_bytes: List[BytesStep] = field(
            default_factory=list,
            metadata={
                "name": "incoming-plc-bytes",
                "type": "Element",
                "namespace": "",
                "sequence": 1,
            },
        )
        api_request: List[MessageStep] = field(
            default_factory=list,
            metadata={
                "name": "api-request",
                "type": "Element",
                "namespace": "",
                "sequence": 1,
            },
        )
        api_response: List[MessageStep] = field(
            default_factory=list,
            metadata={
                "name": "api-response",
                "type": "Element",
                "namespace": "",
                "sequence": 1,
            },
        )
        delay: List[int] = field(
            default_factory=list,
            metadata={
                "type": "Element",
                "namespace": "",
                "sequence": 1,
            },
        )
        terminate: List[object] = field(
            default_factory=list,
            metadata={
                "type": "Element",
                "namespace": "",
                "sequence": 1,
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
        steps: Optional["DriverTestsuite.Testcase.Steps"] = field(
            default=None,
            metadata={
                "type": "Element",
                "namespace": "",
                "required": True,
            },
        )

        @dataclass
        class Steps:
            outgoing_plc_message: List[MessageStep] = field(
                default_factory=list,
                metadata={
                    "name": "outgoing-plc-message",
                    "type": "Element",
                    "namespace": "",
                    "sequence": 1,
                },
            )
            outgoing_plc_bytes: List[BytesStep] = field(
                default_factory=list,
                metadata={
                    "name": "outgoing-plc-bytes",
                    "type": "Element",
                    "namespace": "",
                    "sequence": 1,
                },
            )
            incoming_plc_message: List[MessageStep] = field(
                default_factory=list,
                metadata={
                    "name": "incoming-plc-message",
                    "type": "Element",
                    "namespace": "",
                    "sequence": 1,
                },
            )
            incoming_plc_bytes: List[BytesStep] = field(
                default_factory=list,
                metadata={
                    "name": "incoming-plc-bytes",
                    "type": "Element",
                    "namespace": "",
                    "sequence": 1,
                },
            )
            api_request: List[MessageStep] = field(
                default_factory=list,
                metadata={
                    "name": "api-request",
                    "type": "Element",
                    "namespace": "",
                    "sequence": 1,
                },
            )
            api_response: List[MessageStep] = field(
                default_factory=list,
                metadata={
                    "name": "api-response",
                    "type": "Element",
                    "namespace": "",
                    "sequence": 1,
                },
            )
            delay: List[int] = field(
                default_factory=list,
                metadata={
                    "type": "Element",
                    "namespace": "",
                    "sequence": 1,
                },
            )
            terminate: List[object] = field(
                default_factory=list,
                metadata={
                    "type": "Element",
                    "namespace": "",
                    "sequence": 1,
                },
            )
