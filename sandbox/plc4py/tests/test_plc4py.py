#
# Licensed to the Apache Software Foundation (ASF) under one
# or more contributor license agreements.  See the NOTICE file
# distributed with this work for additional information
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

from plc4py import __version__
from plc4py.PlcDriverManager import PlcDriverManager
from plc4py.api.PlcConnection import PlcConnection


def test_version():
    assert __version__ == "0.1.0"

def test_plc_driver_manager_init():
    driver_manager = PlcDriverManager()
    with driver_manager.connection("modbus:tcp://127.0.0.1:502") as connection:
        assert isinstance(connection, PlcConnection)
