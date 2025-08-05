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

from urllib.parse import urlparse


def get_protocol_code(url: str) -> str:
    """
    Get the protocol code section of the connection string
    e.g. modbus:tcp://127.0.0.1:502 would return modbus

    :param url: The connection string
    :return: The protocol code
    """
    parsed = urlparse(url)
    return parsed.scheme


def strtobool(value: str) -> bool:
    """Credit goes to https://danielms.site/zet/2023/pythons-distutil-strtobool-replacement/"""
    value = value.lower()
    if value in ("y", "yes", "on", "1", "true", "t"):
        return True
    return False
