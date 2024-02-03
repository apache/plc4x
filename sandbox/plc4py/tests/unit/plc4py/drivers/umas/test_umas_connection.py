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

from plc4py.PlcDriverManager import PlcDriverManager


async def manual_test_plc_driver_umas_connect():
    driver_manager = PlcDriverManager()
    async with driver_manager.connection("umas://127.0.0.1:5555") as connection:
        assert connection.is_connected()
    assert not connection.is_connected()


async def test_plc_driver_umas_read():
    driver_manager = PlcDriverManager()
    async with driver_manager.connection("umas://192.168.1.174:502") as connection:
        with connection.read_request_builder() as builder:
            builder.add_item("Random Tag 2", "test_DINT:DINT")
            builder.add_item("Random Tag 1", "test_REAL:REAL")
            builder.add_item("Random Tag 3", "test_BYTE:BYTE")
            request = builder.build()

        future = connection.execute(request)
        await future
        response = future.result()
    pass
