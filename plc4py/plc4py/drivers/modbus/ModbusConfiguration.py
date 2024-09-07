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

from plc4py.utils.GenericTypes import ByteOrder


class ModbusConfiguration(PlcConfiguration):
    """
    Specific Modbus Configuration
    """

    def __init__(self, url):
        """
        Initializes the ModbusConfiguration with the given URL.

        Args:
            url (str): The URL of the Modbus device.
        """
        super().__init__(url)

        # Set the transport to TCP if not specified
        self.transport = self.transport or "tcp"
        # Set the port to 502 if not specified
        self.port = self.port or 502
        # Get the unit identifier from the parameters, default to 1
        self.unit_identifier = self.parameters.get("unit_identifier", 1)
        #Specifies the byte order of the payload
        self.byte_order = self.parameters.get("byte_order", ByteOrder.BIG_ENDIAN)
