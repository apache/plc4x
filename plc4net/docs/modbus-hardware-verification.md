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

Status: **not yet verified against hardware** — planned against a Siemens
S7-1200 with a **CM 1241 (RS422/485)** communication module running `MB_SLAVE`.

## 1a. Prepare an S7-1200 + CM 1241 as a Modbus RTU slave (TIA Portal)

1. **Add the CM 1241** to the device configuration. For a bus use the
   **CM 1241 (RS422/485)** (`6ES7241-1CH3x`); set the port to
   **Operating mode → Half duplex (RS485) 2-wire**. (The RS232 variant,
   `6ES7241-1AH3x`, works point-to-point.)
2. In the CM 1241 port config set **baud rate**, **parity**, **data bits**,
   **stop bits**. The plc4net defaults are **19200-8-E-1**; match them or pass
   `--baud` / `--parity` / … to the tool.
3. In the program:
   - **Startup OB** — call `MB_COMM_LOAD` once. `PORT` = the CM 1241's hardware
     identifier (from the *System constants* tab), `BAUD` / `PARITY` as above,
     `MB_DB` = the `MB_SLAVE` instance DB's `MB_DB` static.
   - **Cyclic OB (e.g. OB1)** — call `MB_SLAVE` every scan. `MB_ADDR` = the slave
     station number (1); `MB_HOLD_REG` = a pointer into a data block, e.g.
     `P#DB10.DBX0.0 WORD 50`.
4. **The `MB_HOLD_REG` DB must be non-optimized** (Properties → *Attributes* →
   untick *Optimized block access*). `MB_SLAVE` will not compile against an
   optimized DB. Put a known value in the first word so `holding:0` has something
   to read — e.g. `DB10.DBW0 := 16#1234`.
5. Modbus holding register `n` = word `n` of `MB_HOLD_REG`. `modbus-verify`'s
   default `holding:0` reads the first word.
6. Compile, **download hardware + software**, set the CPU to RUN.

### Wiring

- USB↔RS-485 adapter `A`/`B` (a.k.a. `D+`/`D−`) to the CM 1241's
  `TxD+/RxD+` (`A`) and `TxD−/RxD−` (`B`). For 2-wire, tie the CM 1241's
  `T` and `R` pairs together per the manual.
- Tie the signal grounds.
- For a short bench link, enable the CM 1241's internal bias/termination in the
  port config; without bias an idle line floats and the first byte is often lost.

## 1b. Any Modbus TCP device

Point the tool at `host [port] [unit-id] [read-address]`. A soft PLC
(`diagslave`, ModbusPal) is fine.

## 2. Run it

### From the source tree

```bash
cd plc4x/plc4net

# Modbus RTU (serial) — the port name selects this path automatically
dotnet run --project tools/modbus-verify -- COM3 1 holding:0 \
    --baud 19200 --parity Even > docs/modbus-hardware-report.md

# Modbus TCP
dotnet run --project tools/modbus-verify -- 192.168.0.9 502 1 holding:0
```

### As a packaged tool (the "build a package, import, use it" path)

```bash
cd plc4x/plc4net
dotnet pack tools/modbus-verify -c Release -o ./_localfeed
dotnet tool install --global --add-source ./_localfeed modbus-verify
modbus-verify COM3 1 holding:0 --baud 19200 --parity Even
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
| Raw exchange: **nothing** comes back | A/B swapped; wrong slave address; baud/parity mismatch with `MB_COMM_LOAD`; CM 1241 not in RS485 half-duplex; `MB_SLAVE` not called every scan; no bias on the line. |
| Raw exchange: response **starts with the request bytes** | The adapter echoes its own transmitter (half-duplex self-receive). The driver does not strip it — use an auto-direction adapter that doesn't echo. |
| Raw CRC valid, but the **driver** read is not `Ok` | Suspect `ModbusRtuConnection.SendAndReceive` — it reads whatever is available once ≥ 4 bytes arrive, with no expected-length or t3.5 inter-frame-gap check, so a byte-at-a-time UART can hand it a partial frame. |
| Modbus exception `0x02` (IllegalDataAddress) | The register is outside the slave's map — enlarge the `MB_HOLD_REG` DB, or read a lower address. |
| Modbus exception `0x01` (IllegalFunction) | `MB_SLAVE` firmware doesn't support that function code for that address range. |
| TCP: connection refused / timeout | Wrong IP or port; port 502 blocked; the device's Modbus server is off. |

## 4. What it covers

- **RTU**: open the port; a raw request/response on the wire with a CRC and
  exception-code check; a `ModbusRtuConnection` read of one tag.
- **TCP**: connect; a `ModbusConnection` read of one tag with the MBAP exchange
  logged.

Not covered: writes to real hardware, multi-register reads, `WriteMultiple*`,
sustained polling, more than one slave on the bus.
