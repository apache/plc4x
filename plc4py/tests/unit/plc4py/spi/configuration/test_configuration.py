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


def test_configuration_standard_raw_ip():
    config = PlcConfiguration(
        "profibus:raw://127.0.0.1:4664&host=localhost&mac=01:02:03:04:05:06"
    )
    assert config.protocol == "profibus"
    assert config.transport == "raw"
    assert config.host == "127.0.0.1"
    assert config.port == 4664
    assert config.parameters["host"] == "localhost"


def test_configuration_standard_tcp_localhost():
    config = PlcConfiguration(
        "profibus:tcp://localhost:4664&host=localhost&mac=01:02:03:04:05:06"
    )
    assert config.protocol == "profibus"
    assert config.transport == "tcp"
    assert config.host == "localhost"
    assert config.port == 4664
    assert config.parameters["host"] == "localhost"


def test_configuration_standard_no_transport():
    config = PlcConfiguration(
        "profibus://localhost:4664&host=localhost&mac=01:02:03:04:05:06"
    )
    assert config.protocol == "profibus"
    assert config.transport == None
    assert config.host == "localhost"
    assert config.port == 4664
    assert config.parameters["host"] == "localhost"


def test_configuration_standard_second_parameter():
    config = PlcConfiguration(
        "profibus://localhost:4664&host=localhost&mac=01:02:03:04:05:06"
    )
    assert config.protocol == "profibus"
    assert config.transport == None
    assert config.host == "localhost"
    assert config.port == 4664
    assert config.parameters["host"] == "localhost"
    assert config.parameters["mac"] == "01:02:03:04:05:06"


def test_configuration_standard_no_port():
    config = PlcConfiguration(
        "profibus://localhost&host=localhost&mac=01:02:03:04:05:06"
    )
    assert config.protocol == "profibus"
    assert config.transport == None
    assert config.host == "localhost"
    assert config.port == None
    assert config.parameters["host"] == "localhost"
    assert config.parameters["mac"] == "01:02:03:04:05:06"


def test_configuration_standard_no_parameters():
    config = PlcConfiguration("profibus://localhost")
    assert config.protocol == "profibus"
    assert config.transport == None
    assert config.host == "localhost"
    assert config.port == None
    assert config.parameters == {}


def test_configuration_standard_no_parameters():
    config = PlcConfiguration("eip://127.0.0.1&test=plc4x")
    assert config.protocol == "eip"
    assert config.transport == None
    assert config.host == "127.0.0.1"
    assert config.port == None
    assert config.parameters["test"] == "plc4x"
