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
from typing import Type

import pluggy

from plc4py.api.PlcDriver import PlcDriver


class PlcDriverClassLoader:
    """Hook spec for PLC4PY Driver Loaders"""

    hookspec = pluggy.HookspecMarker("plc4py")

    @staticmethod
    @hookspec
    def get_driver() -> Type[PlcDriver]:
        """Returns the PlcConnection class that is used to instantiate the driver"""

    @staticmethod
    @hookspec
    def key() -> str:
        """Unique key to identify the driver"""
