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
import time

import pytest

from plc4py.PlcDriverManager import PlcDriverManager
from plc4py.api.value.PlcValue import PlcResponseCode
import logging

logger = logging.getLogger("testing")


@pytest.mark.asyncio
async def manual_test_plc_driver_modbus_connect():
    """
    Test the connection to a Modbus PLC using PlcDriverManager.
    """
    # Initialize the PlcDriverManager
    driver_manager = PlcDriverManager()

    # Establish a connection to the Modbus PLC
    async with driver_manager.connection("modbus://1") as connection:
        # Check if the connection is successful
        assert connection.is_connected()

    # Ensure the connection is closed after exiting the context manager
    assert not connection.is_connected()


@pytest.mark.asyncio
@pytest.mark.xfail
async def test_plc_driver_modbus_read():
    """
    Test reading data from a Modbus PLC.
    """
    log = logging.getLogger(__name__)

    # Initialize the PlcDriverManager
    driver_manager = PlcDriverManager()

    # Establish a connection to the Modbus PLC
    async with driver_manager.connection("modbus://127.0.0.1:5020") as connection:
        with connection.read_request_builder() as builder:
            builder.add_item("Random Tag", "4x00001[10]")
            request = builder.build()

        # Execute the read request
        for _ in range(100):
            future = connection.execute(request)

            response = await future
            value = response.tags["Random Tag"].value
            log.info("Read tag 4x00001[10] - %s", value)

    pass
