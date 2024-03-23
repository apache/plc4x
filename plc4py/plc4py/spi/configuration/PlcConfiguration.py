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
from dataclasses import InitVar, dataclass, field
from typing import Optional, Dict


@dataclass
class PlcConfiguration:
    url: InitVar[str]
    protocol: Optional[str] = None
    transport: Optional[str] = None
    host: Optional[str] = None
    port: Optional[int] = None
    parameters: Dict[str, str] = field(default_factory=lambda: {})

    def __post_init__(self, url):
        self._parse_configuration(url)

    def _parse_configuration(self, url):
        regex = (
            r"^(?P<protocol>[\w]*)"
            + r"(:(?P<transport>[\w]*))?"
            + r":\/\/(?P<host>[\w+.]*)"
            + r"(:(?P<port>\d+))?"
            + r"(?P<parameters>(&{1}([^&=]*={1}[^&=]*))*)"
        )
        matches = re.search(regex, url)

        self.protocol = matches.group("protocol")

        if matches.group("transport") is not None:
            self.transport = matches.group("transport")

        self.host = matches.group("host")

        if matches.group("port") is not None:
            self.port = int(matches.group("port"))

        parameters = matches.group("parameters")
        if parameters is not None and parameters != "":
            self.parameters = {
                item.split("=")[0]: item.split("=")[1]
                for item in parameters[1:].split("&")
            }
