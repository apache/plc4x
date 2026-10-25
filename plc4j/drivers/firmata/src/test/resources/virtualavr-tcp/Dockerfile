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

# Thin extension over pfichtner/virtualavr that exposes the simulated Arduino
# UART as a TCP listener inside the container, so the Firmata driver's TCP
# transport can talk to it directly. Avoids the host-side socat + PTY
# round-trip the upstream `withTcpSerialMode()` uses (which jSerialComm on
# macOS won't open) and gives the IT a single endpoint that works on every
# OS Docker runs on.
#
# The base image is pinned: 0.3.0 changed virtualavr.js to take a compiled HEX
# file instead of compiling the sketch itself, and the floating `latest` tag
# picked that up silently and broke this IT on every platform.
FROM pfichtner/virtualavr:0.4.0

# StandardFirmata.ino is LGPL-2.1 ©Firmata authors and is therefore not
# compatible with this repository's Apache-2.0 source tree. Fetch it from
# upstream at image build time instead of vendoring it. The pinned commit
# avoids surprise breakage if the upstream sketch evolves.
ARG FIRMATA_REPO=https://raw.githubusercontent.com/firmata/arduino
ARG FIRMATA_REF=2.5.9
RUN curl -fsSL "${FIRMATA_REPO}/${FIRMATA_REF}/examples/StandardFirmata/StandardFirmata.ino" \
        -o /sketch/sketch.ino

# StandardFirmata.ino #includes <Servo.h> and <Firmata.h>. The upstream image
# ships arduino-cli + the AVR core but no extra libraries, so install them at
# image build time.
RUN arduino-cli lib install Servo Firmata

# Pre-compile the sketch once here so every container start skips compilation
# and starts the simulator straight away.
RUN virtualavr-compile-arduino /sketch/sketch.ino /sketch/firmware.hex

ENV FILENAME=firmware.hex

