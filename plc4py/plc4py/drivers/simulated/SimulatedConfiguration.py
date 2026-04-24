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
from plc4py.spi.configuration.PlcConfiguration import PlcConfiguration


class SimulatedConfiguration(PlcConfiguration):
    """
    Configuration for the simulated driver.

    The host part of the connection URL is treated as the device name,
    e.g. `simulated://plc-1` yields a device named `plc-1`. If the URL
    does not provide a host, the device name defaults to `simulated`.
    """

    def __init__(self, url: str):
        super().__init__(url)
        self.device_name = self.host or "simulated"
