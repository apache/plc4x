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

from setuptools import setup, find_packages

setup(
    name="plc4py",
    version="0.11a0",
    description="Plc4py The Python Industrial IOT Adapter",
    classifiers=[
        "Development Status :: 3 - Alpha",
        "License :: OSI Approved :: Apache 2.0 License",
        "Programming Language :: Python :: 3.8",
        "Topic :: Scientific/Engineering :: Interface Engine/Protocol Translator",
    ],
    keywords="modbus plc4x",
    url="https://plc4x.apache.org",
    author='"Apache PLC4X <>"',
    author_email="dev@plc4x.apache.org",
    license="Apache 2.0",
    packages=find_packages(include=["plc4py", "plc4py.*"]),
    setup_requires=[
        "wheel",
    ],
    install_requires=[
        "pytest-asyncio>=0.18.3",
        "pip-tools",
        "black",
        "pip",
        "aenum",
        "bitarray",
    ],
    extras_require={
        "dev": [
            "requires",
            "pre-commit>=2.6.0",
            "pytest-mock>=3.8.1",
            "mock>=4.0.2",
            "mypy>=0.942",
            "flake8>=4.0.1",
        ]
    },
    entry_points={
        "plc4py.drivers": [
            "mock = plc4py.drivers.mock.MockConnection:MockDriverLoader",
            "modbus = plc4py.drivers.modbus.ModbusConnection:ModbusDriverLoader",
        ]
    },
)
