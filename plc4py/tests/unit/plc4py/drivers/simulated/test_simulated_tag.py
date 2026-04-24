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
import pytest

from plc4py.api.exceptions.exceptions import PlcFieldParseException
from plc4py.drivers.simulated.SimulatedTag import (
    SimulatedTag,
    SimulatedTagBuilder,
    SimulatedTagType,
)


def test_parse_state_scalar():
    tag = SimulatedTag.create("STATE/counter:INT")
    assert tag.tag_type == SimulatedTagType.STATE
    assert tag.name == "counter"
    assert tag.data_type == "INT"
    assert tag.quantity == 1


def test_parse_random_scalar():
    tag = SimulatedTag.create("RANDOM/temperature:REAL")
    assert tag.tag_type == SimulatedTagType.RANDOM
    assert tag.name == "temperature"
    assert tag.data_type == "REAL"
    assert tag.quantity == 1


def test_parse_stdout_string():
    tag = SimulatedTag.create("STDOUT/message:STRING")
    assert tag.tag_type == SimulatedTagType.STDOUT
    assert tag.name == "message"
    assert tag.data_type == "STRING"


def test_parse_array():
    tag = SimulatedTag.create("STATE/samples:DINT[5]")
    assert tag.tag_type == SimulatedTagType.STATE
    assert tag.name == "samples"
    assert tag.data_type == "DINT"
    assert tag.quantity == 5


def test_parse_dotted_name():
    tag = SimulatedTag.create("STATE/motor.speed:REAL")
    assert tag.name == "motor.speed"


def test_tag_type_is_case_insensitive():
    tag = SimulatedTag.create("state/x:BOOL")
    assert tag.tag_type == SimulatedTagType.STATE


def test_data_type_is_normalized_to_upper():
    tag = SimulatedTag.create("STATE/x:int")
    assert tag.data_type == "INT"


def test_matches_returns_false_for_invalid():
    assert SimulatedTag.matches("not a tag") is False
    assert SimulatedTag.matches("STATE/x") is False
    assert SimulatedTag.matches(":INT") is False


def test_create_raises_on_invalid_address():
    with pytest.raises(PlcFieldParseException):
        SimulatedTag.create("not a tag")


def test_create_raises_on_unknown_tag_type():
    with pytest.raises(PlcFieldParseException):
        SimulatedTag.create("BOGUS/x:INT")


def test_builder_creates_same_tag():
    built = SimulatedTagBuilder.create("STATE/x:INT")
    direct = SimulatedTag.create("STATE/x:INT")
    assert built == direct
