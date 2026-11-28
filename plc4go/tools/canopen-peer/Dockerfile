# Licensed to the Apache Software Foundation (ASF) under one or more contributor license
# agreements. See the NOTICE file distributed with this work for additional information regarding
# copyright ownership. The ASF licenses this file to you under the Apache License, Version 2.0 (the
# "License"); you may not use this file except in compliance with the License. You may obtain a
# copy of the License at https://www.apache.org/licenses/LICENSE-2.0. Unless required by
# applicable law or agreed to in writing, software distributed under the License is distributed on
# an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#
# Test peers only: two CANopen nodes written independently of plc4x, CANopenNode's Linux device
# (CANopenLinux v4.0, C, the DS301 example profile) and the python canopen library (2.4.1,
# LocalNode), next to the Go toolchain, so internal/canopen's peer_linux_test.go can run against
# both on a virtual CAN bus in one container. Nothing from either is linked into plc4go.
FROM golang:1.27.4
RUN apt-get update -qq && apt-get install -y -qq --no-install-recommends iproute2 can-utils python3 python3-venv && rm -rf /var/lib/apt/lists/*
RUN git clone --depth 1 --branch v4.0 --recurse-submodules --shallow-submodules https://github.com/CANopenNode/CANopenLinux.git /src/CANopenLinux
WORKDIR /src/CANopenLinux
RUN make -j4 && cp canopend /usr/local/bin/canopend
RUN python3 -m venv /opt/canopen && /opt/canopen/bin/pip install -q canopen==2.4.1
