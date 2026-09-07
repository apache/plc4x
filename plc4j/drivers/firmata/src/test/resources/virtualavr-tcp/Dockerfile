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
FROM pfichtner/virtualavr:latest

# StandardFirmata.ino is LGPL-2.1 ©Firmata authors and is therefore not
# compatible with this repository's Apache-2.0 source tree. Fetch it from
# upstream at image build time instead of vendoring it. The pinned commit
# avoids surprise breakage if the upstream sketch evolves.
ARG FIRMATA_REPO=https://raw.githubusercontent.com/firmata/arduino
ARG FIRMATA_REF=2.5.9
RUN curl -fsSL "${FIRMATA_REPO}/${FIRMATA_REF}/examples/StandardFirmata/StandardFirmata.ino" \
        -o /sketch/sketch.ino

# StandardFirmata.ino #includes <Servo.h> and <Firmata.h>. The upstream image
# ships arduino-cli + the AVR core but no extra libraries, and virtualavr's
# own library installer is currently disabled (compileArduinoSketch sets
# installCommand=false), so install them at image build time.
RUN arduino-cli lib install Servo Firmata

# The upstream pfichtner/virtualavr HEALTHCHECK pings ws://localhost:8080,
# which only starts when the simulator starts, which only happens on the
# first TCP accept — so the inherited healthcheck would stay perpetually
# unhealthy. Drop it.
HEALTHCHECK NONE

EXPOSE 3030

# Wrap socat in a tiny shell loop so the listener survives client
# reconnects. ",fork" would let multiple clients run simultaneously but
# would also spawn multiple simulators racing for WebSocket port 8080.
# `while true; do socat ...; done` accepts one client at a time, exits
# socat when the client disconnects, then re-runs socat for the next.
#
# Redirecting stdin to /dev/null avoids docker-attach's stdin buffer
# closing on the child and tripping a clean shutdown of socat (which was
# what made the container exit cleanly within seconds under Testcontainers'
# default log-streaming attach during earlier testing).
ENTRYPOINT ["/bin/sh", "-c", \
            "while true; do socat TCP-LISTEN:3030,reuseaddr EXEC:'node /app/virtualavr.js sketch.ino',pty,rawer,fdin=3,fdout=4; sleep 1; done </dev/null"]
