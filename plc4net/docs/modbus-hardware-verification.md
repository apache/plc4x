<!--
  Licensed to the Apache Software Foundation (ASF) under one
  or more contributor license agreements.  See the NOTICE file
  distributed with this work for additional information
  regarding copyright ownership.  The ASF licenses this file
  to you under the Apache License, Version 2.0 (the
  "License"); you may not use this file except in compliance
  with the License.  You may obtain a copy of the License at

      https://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing,
  software distributed under the License is distributed on an
  "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
  KIND, either express or implied.  See the License for the
  specific language governing permissions and limitations
  under the License.
  -->
# Modbus hardware verification

`tools/modbus-verify` connects to a real Modbus device, reads one tag and prints
a Markdown report. It exits 0 on pass, 1 on any failure. It covers **both**
transports:

- **Modbus TCP** — `ModbusConnection` over the TCP transport.
- **Modbus RTU** — `ModbusRtuConnection` over the serial transport. This is the
  first exercise of `SerialTransportInstance` and `ModbusRtuConnection` against
  real hardware, so the RTU run also does a **raw frame exchange** (write the
  request bytes, print every byte that comes back) independently of the driver's
  own framing, so a framing problem is visible even if the driver mis-decodes.

Status: Modbus TCP is verified against a software slave. Modbus RTU is verified
end-to-end — both the raw frame exchange and the `ModbusRtuConnection` driver
read — against a Mitsubishi **QJ71C24N** running its non-procedure ("no-protocol")
communication mode as a fixed Modbus RTU slave.  An earlier attempt used a
Siemens S7-1214C + **CM 1241 (RS422/485)** running `MB_SLAVE`; that rig was
dropped after request framing verified correct on the wire but the round-trip
never completed, and the fault was isolated to the RS-485 bench link (adapter
and bias/termination), not the driver — see the 2026-09-06 entry in
[modbus-hardware-report.md](modbus-hardware-report.md) for that record.
Run log: [modbus-hardware-report.md](modbus-hardware-report.md).

## 1a. Prepare a QJ71C24N as a fixed Modbus RTU slave (GX Works2)

The QJ71C24N has no native Modbus slave firmware (unlike the QJ71MB91), so the
ladder program emulates one: it receives the master's request and always answers
with the same canned response, regardless of the request's function code or
address. It does not decode Modbus — it is a fixed-response stand-in, useful for
exercising the wire protocol and framing, not a general slave.

1. **Configure CH2** for **non-procedure (no-protocol) communication**
   (communication protocol `0006H`) in GX Works2: **9600 baud, 8 data bits,
   Even parity, 1 stop bit**, sum check **on**, byte units (buffer memory
   `136H` = `1`), fixed receive length **8 bytes** (`144H` = `8`), no receive
   terminator (`145H` = `FFFFH`). Match these or pass `--baud 9600` to the tool
   (plc4net's own default is 19200-8-E-1, so `--baud` is required here).
2. On first scan (`SM402`), initialize the `G.INPUT` / `G.OUTPUT` control data:
   channel 2, receive length 4 words (8 bytes), response length 13 bytes.
3. On the CH2 receive-complete flag (`X0A` in the working program), run
   `G.INPUT` to capture the request into a data area (`D100`…), then run
   `G.OUTPUT` to send the fixed 13-byte response.
4. **Gate the response send on a real "ready" condition, not a stale
   interlock.** The working rig's first pass received every request cleanly
   (confirmed via the module's receive indication) but never sent a response —
   traced to an **M-relay interlock in the ladder gating the `G.OUTPUT` rung**
   that was left unsatisfied. Clearing that interlock condition is what made the
   round trip work; if a `modbus-verify` raw exchange gets the request in but
   nothing back, check the send-side interlock before suspecting the link.
5. Map Modbus holding registers to the response buffer: Holding Register `0`-`3`
   ↔ `D1000`-`D1003`. The verified rig seeds these `1`, `2`, `3`, `4`.

### Wiring

- USB↔RS-485 adapter `A`/`B` (a.k.a. `D+`/`D−`) to the QJ71C24N's RS-485
  terminals per the module's manual; tie the signal grounds.
- The module answers with its full fixed response (13 bytes / 4 registers)
  regardless of the request's quantity field — a `holding:0` single-register
  read still gets all 4 registers back. `modbus-verify`'s raw-shape check flags
  this as an "invalid" shape (it compares against the *requested* quantity),
  but `ModbusRtuConnection` frames the response by its own declared byte count
  and decodes correctly regardless — this is a property of the fixed-response
  rig, not a driver defect.

## 1b. Any Modbus TCP device

Point the tool at `host [port] [unit-id] [read-address]`. A soft PLC
(`diagslave`, ModbusPal) is fine.

## 2. Run it

### From the source tree

```bash
cd plc4x/plc4net

# Modbus RTU (serial) — the port name selects this path automatically.
# The QJ71C24N rig runs at 9600-8-E-1, not plc4net's own 19200-8-E-1 default.
dotnet run --project tools/modbus-verify -- COM3 1 holding:0 \
    --baud 9600 --parity Even > docs/modbus-hardware-report.md

# Modbus TCP
dotnet run --project tools/modbus-verify -- 192.168.0.9 502 1 holding:0
```

### As a packaged tool (the "build a package, import, use it" path)

```bash
cd plc4x/plc4net
dotnet pack tools/modbus-verify -c Release -o ./_localfeed
dotnet tool install --global --add-source ./_localfeed modbus-verify
modbus-verify COM3 1 holding:0 --baud 9600 --parity Even
```

`dotnet pack` on a `PackAsTool` project bundles every dependency into the one
`.nupkg`, so the tool is self-contained. Nothing is published anywhere.

### Options

```
Modbus TCP:  modbus-verify <host> [port] [unit-id] [read-address]
Modbus RTU:  modbus-verify <COMx|/dev/ttyUSB0> [unit-id] [read-address]
                 [--baud 19200] [--parity Even|None|Odd]
                 [--stop-bits One|Two] [--data-bits 8]
```

`read-address` is a `ModbusTag` string: `holding:0`, `input:5`, `coil:2`,
`discrete:1`. A full connection string (`modbus-rtu://COM3?...`,
`modbus-tcp:tcp://host:502?...`) is also accepted as the first argument.

## 3. If it fails

| Symptom | Try |
|---|---|
| `Could not open COMx` | Wrong port name (Device Manager), another program holds it, or the adapter is unplugged. |
| Raw exchange: request goes out, **nothing** comes back | A/B swapped; baud/parity mismatch with the slave's port config; on the QJ71C24N rig, check the ladder-side send interlock before the link — a satisfied receive with no reply is the interlock symptom, not a wiring one. |
| Raw exchange: response **starts with the request bytes** | The adapter echoes its own transmitter (half-duplex self-receive). The driver does not strip it — use an auto-direction adapter that doesn't echo. |
| Raw CRC valid, but the **driver** read is not `Ok` | Suspect `ModbusRtuConnection.SendAndReceive` — it reads whatever is available once ≥ 4 bytes arrive, with no expected-length or t3.5 inter-frame-gap check, so a byte-at-a-time UART can hand it a partial frame. |
| Raw shape flagged **invalid** but the driver read is `Ok` | Expected on the QJ71C24N rig — it always returns its full fixed response regardless of the requested quantity; the shape check compares against what was asked, not what a fixed-response device actually sends. |
| Modbus exception `0x02` (IllegalDataAddress) | The register is outside the slave's map — enlarge the register map, or read a lower address. |
| Modbus exception `0x01` (IllegalFunction) | The slave doesn't support that function code for that address range. |
| TCP: connection refused / timeout | Wrong IP or port; port 502 blocked; the device's Modbus server is off. |

## 4. What it covers

- **RTU**: open the port; a raw request/response on the wire with a CRC and
  exception-code check (`--quantity` reads multiple registers/coils on the raw
  path); a `ModbusRtuConnection` read of one tag.
- **TCP**: connect; a `ModbusConnection` read of one tag with the MBAP exchange
  logged.

Not covered: writes to real hardware, multi-register reads through the driver
(`ModbusTag` and the read-request path are still single-value), `WriteMultiple*`,
sustained polling, more than one slave on the bus.
