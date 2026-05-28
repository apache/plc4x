#
# Licensed to the Apache Software Foundation (ASF) under one
# or more contributor license agreements.  See the NOTICE file
# distributed with this work for additional information
# regarding copyright ownership.  The ASF licenses this file
# to you under the Apache License, Version 2.0 (the
# "License"); you may not use this file except in compliance
# with the License.  You may obtain a copy of the License at
#
#   https://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing,
# software distributed under the License is distributed on an
# "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
# KIND, either express or implied.  See the License for the
# specific language governing permissions and limitations
# under the License.
#

# Minimal BACnet/IP virtual device for plc4go integration tests.
# Built on bacpypes3 (the modern Python BACnet stack).
FROM python:3.12-slim

# bacpypes3 transitively pulls in netifaces; install build tooling so the
# wheel builds on minimal slim images, then drop it to keep the layer small.
RUN apt-get update \
 && apt-get install -y --no-install-recommends gcc python3-dev \
 && pip install --no-cache-dir bacpypes3==0.0.102 \
 && apt-get purge -y gcc python3-dev \
 && apt-get autoremove -y \
 && rm -rf /var/lib/apt/lists/*

WORKDIR /sim
COPY device.py /sim/device.py

# Standard BACnet/IP unicast + broadcast port.
EXPOSE 47808/udp

# -u forces unbuffered stdout so `docker compose logs` shows lifecycle messages
# in real time instead of dumping them only on container exit.
CMD ["python", "-u", "device.py"]
