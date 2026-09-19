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
# Hardware verification

Automated tests prove the drivers against scripted transports; this document
covers what was proven against **physical devices**, how to reproduce it, and —
just as importantly — what each rig does *not* prove.

Two harnesses drive the **public driver API** exactly as a NuGet consumer would,
print a Markdown report and exit 0 on pass, 1 on any failure:

| Harness | Drivers exercised |
|---|---|
| `tools/s7-verify` | `S7Driver` over COTP/TCP |
| `tools/modbus-verify` | `ModbusConnection` (TCP) and `ModbusRtuConnection` (serial) |

`modbus-verify`'s RTU path additionally performs a **raw frame exchange** —
write the request bytes, print every byte that comes back — independently of the
driver's own framing, so a framing fault stays visible even when the driver
mis-decodes it.

## Status

| Protocol | Rig | Result | Date |
|---|---|---|---|
| **S7** (ISO-on-TCP) | Siemens S7-1214C DC/DC/DC, rack 0 / slot 1 | **PASS 50/50**; persistent I/Q/M/DB matrix 43/43 with independent read-back 25/25 | 2026-09-16 |
| **Modbus TCP** | software slave (no Modbus/TCP device on hand) | **PASS 5/5** | 2026-09-06 |
| **Modbus RTU** | Mitsubishi QJ71C24N in non-procedure mode, as a *fixed-response* slave | **PASS** — raw exchange and driver read | 2026-09-19 |
| **KNXnet/IP** | none — scripted loopback gateway only | **not hardware-verified** | — |

### What these results do and do not prove

- **S7 is verified in the strong sense**: a real CPU, every scalar width,
  absolute I/Q/M addressing, writes with independent read-back over new
  connections, and an error path.
- **Modbus RTU is verified at the wire and framing level only.** The QJ71C24N
  has no native Modbus slave firmware, so its ladder program answers *every*
  request with the same canned response regardless of function code, address or
  requested quantity. That proves the serial transport, RTU framing, CRC and the
  driver's read path against a real byte-at-a-time UART. It does **not** exercise
  address-range handling, exception codes, or per-function decoding on the slave
  side.
- **Modbus TCP was verified against a software slave**, not a physical device.
- **KNXnet/IP has no hardware verification at all** — no gateway is available.

## Running the harnesses

### From the source tree

```bash
cd plc4x/plc4net

# S7
dotnet run --project tools/s7-verify -- <PLC-IP> --db 100

# Modbus RTU (serial) — the port name selects this path automatically.
# The QJ71C24N rig runs at 9600-8-E-1, not plc4net's own 19200-8-E-1 default.
dotnet run --project tools/modbus-verify -- COM3 1 holding:0 --baud 9600 --parity Even

# Modbus TCP
dotnet run --project tools/modbus-verify -- 192.168.0.9 502 1 holding:0
```

### As a packaged tool (the "build a package, import, use it" path)

```bash
cd plc4x/plc4net
dotnet pack tools/s7-verify -c Release -o ./_localfeed
dotnet tool install --global --add-source ./_localfeed s7-verify
s7-verify <PLC-IP> --db 100
```

`dotnet pack` on a `PackAsTool` project bundles every dependency into the one
`.nupkg`, so the tool is self-contained — no other package needs to be in the
feed. Nothing is published anywhere; `_localfeed` is a folder on your machine.
`modbus-verify` packs and installs the same way.

### Options

```
s7-verify <host> [--rack N] [--slot N] [--db N]
          [--device-group PG_OR_PC|OS|OTHERS] [--remote-tsap 0xNNNN]
          [--read <address>] [--i-base N] [--q-base N] [--m-base N]
          [--read-only] [--write-markers] [--write-outputs]
          [--keep-db-values] [--keep-marker-values] [--keep-output-values]

modbus-verify <host> [port] [unit-id] [read-address]
modbus-verify <COMx|/dev/ttyUSB0> [unit-id] [read-address]
              [--baud 19200] [--parity Even|None|Odd]
              [--stop-bits One|Two] [--data-bits 8] [--quantity 1]
```

`s7-verify` defaults: `--rack 0 --slot 1 --db 100`. `--read <address>` skips the
suite and just connects, reads that one tag and prints the outcome — a focused
probe for one address (e.g. `--read "%I0.0"`, `--read "%DB100.DBW2"`) that needs
neither DB100 nor the full layout.

The normal suite writes all seven scalar types only in DB100 and restores every
original value. `--read-only` stops after the DB and I/Q/M address reads without
sending any Write Var. `--keep-db-values` intentionally skips the DB restore so
the written values remain available for online inspection.

`--write-markers` additionally exercises BOOL/BYTE/INT/DINT/REAL/WORD/DWORD at
M100..M117, restoring every original value by default. Change `--m-base` if that
range is not reserved for testing. Combining it with `--keep-marker-values`
intentionally leaves the written marker values in the PLC.

`--write-outputs` performs the same sequence in Q memory and is deliberately
opt-in: Q writes can energize physical outputs. Use it only after the machine is
isolated, a second on-site approver has confirmed the test, and the selected
`--q-base` range exists. It restores Q values by default. Combining it with
`--keep-output-values` intentionally leaves the written Q values in the PLC for
online inspection. Input memory is always read-only.

For `modbus-verify`, `read-address` is a `ModbusTag` string: `holding:0`,
`input:5`, `coil:2`, `discrete:1`. A full connection string
(`modbus-rtu://COM3?...`, `modbus-tcp:tcp://host:502?...`) is also accepted as
the first argument. `--quantity` applies to the raw exchange only; when it is
not 1 the driver-read step is skipped, because `ModbusTag` and the read-request
path are still single-value.

---

# S7 — Siemens S7-1214C

## Preparing the PLC (TIA Portal)

1. **Enable PUT/GET.** PLC → Properties → *Protection & Security* → *Connection
   mechanisms* → tick **"Permit access with PUT/GET communication from remote
   partner"**. S7-1200 / S7-1500 refuse the S7 driver without this.
2. **Add a non-optimized data block.** Add a global DB (e.g. **DB100**) →
   Properties → *Attributes* → untick **"Optimized block access"**. The driver
   uses absolute addressing (`%DB100.DBW2`), which only works on standard DBs.
3. **Add these rows, in this order.** Enter only the *Name*, *Type* and
   *Start value* — TIA fills in *Offset* itself when you compile. Do **not**
   type the offset anywhere; it is shown only so you can check the result.

   | Name | Type | Start value | (Offset TIA should show) |
   |---|---|---|---|
   | `flag` | Bool | `true` | `0.0` |
   | `b` | Byte | `16#A5` | `1.0` |
   | `i` | Int | `-12345` | `2.0` |
   | `di` | DInt | `-1000000` | `4.0` |
   | `r` | Real | `3.14159` | `8.0` |
   | `w` | Word | `16#BEEF` | `12.0` |
   | `dw` | DWord | `16#DEADBEEF` | `14.0` |
   | `w_wr` | Int | `0` | `18.0` |
   | `r_wr` | Real | `0.0` | `20.0` |

   In a non-optimized DB the compiler lays these out contiguously as
   `Bool, Byte, Int, DInt, Real, Word, DWord, Int, Real`, producing exactly the
   offsets in the last column. Compile and confirm the *Offset* column matches
   before downloading — a mismatch means a row is the wrong type or out of order.
4. Note the PLC **IP address** and the DB number. S7-1200 / S7-1500 are
   **rack 0, slot 1**. Make sure TCP port **102** is reachable from the PC
   (same subnet, firewall allows outbound 102).
5. ⚠️ The write test only writes to the DB. It never touches `%Q` unless
   `--write-outputs` is passed — on a DC/DC/DC CPU those are real transistor
   outputs.

## What the S7 run covers

- COTP CR/CC handshake, S7 Setup Communication, negotiated PDU length
- Read: BOOL, BYTE, INT, DINT, REAL, WORD, DWORD from a DB; a 3-tag single
  request
- Read absolute addresses: I/Q/M bit, byte, word and double word
- Write + read-back + restore: BOOL, BYTE, INT, DINT, REAL, WORD and DWORD in
  the DB; optionally the same matrix in M or Q
- Error path: reading a non-existent DB, connection survives

Not covered: reads larger than one negotiated PDU (multi-PDU), STRING, the TIA
date/time types, subscriptions.

## If the S7 connection fails

| Symptom | Try |
|---|---|
| `No COTP Connection Confirm received` | Wrong rack/slot, or the S7-1200/1500 wants a different connection resource. Retry with `--device-group OTHERS` (TSAP `0x03rs`), then `--remote-tsap 0x0301`, `0x0302`, `0x0300`, `0x0201`. Record which one works. (The S7-1214C verified on 2026-09-03 needed none of these — the rack/slot default `0x0101` worked.) |
| `No S7 Setup Communication response` | PUT/GET not enabled, wrong rack/slot, or wrong TSAP. |
| Connect succeeds but **every** read/write is `AccessDenied` (S7 error `0x8104`) | PUT/GET is not effective on the CPU. Tick "Permit access with PUT/GET communication from remote partner", then **compile the hardware configuration and download it** — the offline setting alone changes nothing. Check the *Access level* table is "Full access (no protection)". |
| Connection drops immediately | Another master already holds the single PG connection — close TIA Portal's online view, or use a dedicated S7 connection + `--remote-tsap`. |
| Every DB read is `NotFound` | DB100 does not exist on the CPU (create it and download), or the DB number is wrong. |
| A DB read is `InvalidAddress` | The DB is *optimized* (untick "Optimized block access"), or the offset does not exist. |
| A read returns the wrong value (e.g. the byte offset itself) | The DB's *Start value* for that field is wrong — do not type the offset into the start-value cell; TIA computes offsets itself. |
| TCP timeout | Port 102 blocked, or wrong IP. |

Once a `--device-group` / `--remote-tsap` value works, it belongs in the driver's
documentation and, ideally, as the S7-1200/1500 default.

## S7 run log

### 2026-09-16 — persistent I/Q/M/DB matrix

The PLC had no external equipment attached. Inputs were read only. DB100,
M100..M117 and Q0..Q17 were each exercised with read-before → write → immediate
read-back, with all three restore steps explicitly disabled. After the write
session closed, every target was independently read through a new connection.

- **Device endpoint**: `192.168.1.11`, rack 0 / slot 1
- **Connection**: COTP + Setup Communication passed; negotiated PDU 240 bytes
- **Main run**: **PASS (43/43)**
- **Independent read-back**: **PASS (25/25)**
- **Persistent write ranges**: DB100 bytes 0..17, M100..M117, Q0..Q17

  | Area/type | Address | Independently read back |
  |---|---|---:|
  | Input bit | `%I0.0` | `False` |
  | Input byte | `%IB0` | `0x00` |
  | Input word | `%IW0` | `0x0000` |
  | Input double word | `%ID0` | `0x00000000` |
  | Output BOOL | `%Q0.0` | `True` |
  | Output BYTE | `%QB1` | `0x3C` |
  | Output INT | `%QW2` | `23456` (`0x5BA0`) |
  | Output DINT | `%QD4` | `-123456789` (`0xF8A432EB`) |
  | Output REAL | `%QD8` | `-12.5` (`0xC1480000`) |
  | Output WORD | `%QW12` | `0x1357` |
  | Output DWORD | `%QD14` | `0x89ABCDEF` |
  | Marker BOOL | `%M100.0` | `True` |
  | Marker BYTE | `%MB101` | `0x3C` |
  | Marker INT | `%MW102` | `23456` (`0x5BA0`) |
  | Marker DINT | `%MD104` | `-123456789` (`0xF8A432EB`) |
  | Marker REAL | `%MD108` | `-12.5` (`0xC1480000`) |
  | Marker WORD | `%MW112` | `0x1357` |
  | Marker DWORD | `%MD114` | `0x89ABCDEF` |
  | DB BOOL | `%DB100.DBX0.0` | `False` |
  | DB BYTE | `%DB100.DBB1` | `0x3C` |
  | DB INT | `%DB100.DBW2` | `23456` (`0x5BA0`) |
  | DB DINT | `%DB100.DBD4` | `-123456789` (`0xF8A432EB`) |
  | DB REAL | `%DB100.DBD8` | `-12.5` (`0xC1480000`) |
  | DB WORD | `%DB100.DBW12` | `0x1357` |
  | DB DWORD | `%DB100.DBD14` | `0x89ABCDEF` |

These are live process-image/memory values, not retained startup values; PLC
logic, an input transition, a mode change or a restart may overwrite them.

### 2026-09-16 — persistent Q write verification

The PLC had no external equipment attached. The output matrix was explicitly run
with `--write-outputs --keep-output-values`, so the final values were not
restored and can be inspected online in TIA Portal.

- **Device endpoint**: `192.168.1.11`, rack 0 / slot 1
- **Sequence**: each Q target passed read-before → write → immediate read-back
- **Independent read-back**: all targets were read again over new connections
  after the write session closed
- **Result**: **PASS (43/43)** for the full run; persistent Q values:

  | Type | Address | Before | Written and independently read back |
  |---|---|---:|---:|
  | BOOL | `%Q0.0` | `False` | `True` |
  | BYTE | `%QB1` | `0x00` | `0x3C` |
  | INT | `%QW2` | `0` | `23456` (`0x5BA0`) |
  | DINT | `%QD4` | `0` | `-123456789` (`0xF8A432EB`) |
  | REAL | `%QD8` | `0.0` | `-12.5` (`0xC1480000`) |
  | WORD | `%QW12` | `0x0000` | `0x1357` |
  | DWORD | `%QD14` | `0x00000000` | `0x89ABCDEF` |

DB100 was also exercised by the normal suite and restored byte-identically.
Marker writes were not requested in this run.

### 2026-09-16 — expanded verification

Re-run from the `feature/plc4net-revival` working tree based on `e86bdd028`
after PUT/GET access was restored on the CPU.

- **Device endpoint**: `192.168.1.11`, rack 0 / slot 1
- **Connection**: COTP + Setup Communication passed; negotiated PDU 240 bytes
- **Address reads**: DB100 scalar values and I/Q/M X/B/W/D forms — PASS 20/20
- **DB writes**: BOOL/BYTE/INT/DINT/REAL/WORD/DWORD — each passed
  read-before → write → read-back → restore → restore-read-back
- **Marker writes**: the same seven types at M100..M117 — each passed and
  restored to its original value
- **Output writes**: not attempted; Q reads passed, but physical-output writes
  remain behind the explicit `--write-outputs` safety switch
- **Error path**: a non-existent DB returned `NotFound`; connection survived
- **Result**: **PASS (50/50)**

Restored values were byte-identical to the snapshots: DB100 returned to
`true/A5/CFC7/FFF0BDC0/40490FD0/BEEF/DEADBEEF`, and M100..M117 returned to all
zeroes.

### 2026-09-16 — expanded re-verification blocked

Attempted from the `feature/plc4net-revival` working tree based on `e86bdd028`.
The expanded harness first ran in `--read-only` mode, so it sent no Write Var
request.

- **Result**: **BLOCKED** — all 19 reads returned `AccessDenied`
- **Likely cause**: the CPU's PUT/GET permission was no longer effective. Enable
  "Permit access with PUT/GET communication from remote partner", compile, and
  download the hardware configuration before retrying.

The failure was consistent across DB100, `%I`, `%Q`, and `%M`, while the S7
session stayed connected. This distinguishes an access-policy refusal from an
address parser or individual memory-range failure.

### 2026-09-04 — re-verify after the develop merge

Re-run on `feature/plc4net-revival` at `39e3792a0` — after `develop` was merged,
the prerequisite check was raised to the .NET 8 SDK and `AssemblyVersion` was
pinned. No S7 driver code changed since the 2026-09-03 run; this confirms those
build/infra commits did not regress the driver.

- **Connection**: `s7://192.168.1.11?remote-rack=0&remote-slot=1&request-timeout=5000` (default TSAPs)
- **Negotiated PDU length**: 240 bytes
- **Result**: **PASS (12/12)**, two consecutive runs. Extra probe: single read of
  `%I0.0` returned `Ok` (`False`) — exercises `%I` area addressing outside the DB.

### 2026-09-03 — first hardware verification

- **Device**: SIMATIC S7-1214C DC/DC/DC (S7-1200 family), rack 0 / slot 1
- **Connection**: `s7://192.168.1.11` (default TSAPs: remote `0x0101`, local `0x0311`)
- **Negotiated PDU length**: 240 bytes
- **Data block**: DB100, non-optimized, laid out per the procedure above
- **Result**: **PASS (12/12)**

```
| | Step | Detail |
|---|---|---|
| OK | Connect | COTP + Setup Communication ok, negotiated PDU length 240 bytes |
| OK | Read BOOL %DB100.DBX0.0 | = true |
| OK | Read BYTE %DB100.DBB1 | = 0xA5 |
| OK | Read INT %DB100.DBW2 | = -12345 |
| OK | Read DINT %DB100.DBD4 | = -1000000 |
| OK | Read REAL %DB100.DBD8 | = 3.1416 |
| OK | Read WORD %DB100.DBW12 | = 0xBEEF |
| OK | Read DWORD %DB100.DBD14 | = 0xDEADBEEF |
| OK | Read 3 tags in one request | all three correct |
| OK | Write + read-back %DB100.DBW18 | = 6789 |
| OK | Write + read-back %DB100.DBD20 | = 12345.5 |
| OK | Read a non-existent DB | rejected with NotFound (connection survived) |
```

#### Bug found and fixed during this run

The CPU refuses a request it cannot serve with a bare **Ack (ROSCTR 0x02)**,
which — like an AckData (0x03) — carries a 2-byte `errorClass` / `errorCode`
field after the 10-byte S7 header. `S7Connection.ReadOneS7MessageAsync` only
treated 0x03 as a 12-byte header, so it read 10 bytes of a 12-byte frame, the
parse failed, and the two leftover bytes desynced every following response.
Fixed to frame both 0x02 and 0x03 as 12 bytes, and to map the S7 header errors
(`0x8104` PUT/GET refused, `0x8304`, `0x85xx`) to a per-tag response code the way
the Java driver's `mapPlcErrorCode` does. Regression tests:
`Read_maps_a_header_level_refusal_to_access_denied` and
`A_bare_Ack_is_framed_as_12_bytes_and_does_not_desync_the_next_request`.

#### CPU-side prerequisites confirmed necessary

- **"Permit access with PUT/GET communication from remote partner"** must be
  ticked (CPU → Protection & Security → Connection mechanisms), then compiled
  **and downloaded** — the offline setting alone has no effect. Without it every
  Read/Write Var comes back `0x8104` even though Setup Communication succeeds.
- **DB100 must exist and be non-optimized.** A missing block returns `NotFound`
  per item; an optimized block cannot be reached by absolute addressing.

## S7 PLC-side evidence

Captured after the persistent write and independent read-back; they confirm the
values in TIA Portal, on the PLC side, in addition to the driver-side report.

- [S7-1214C test rig](images/s7-1214c-test-rig.jpg) — CPU 1214C DC/DC/DC on the
  isolated bench with its 24 VDC power supply.
- [DB100 online values](images/s7-db100-online-values.png) — DB100 offsets,
  declared types, start values and monitored values in TIA Portal.
- [I/Q/M watch table](images/s7-iqm-watch-table.png) — input reads and the
  persistent output/marker values at the tested absolute addresses.

---

# Modbus

## Modbus TCP — any Modbus/TCP device

Point the tool at `host [port] [unit-id] [read-address]`. A soft PLC
(`diagslave`, ModbusPal) is fine — that is what the 2026-09-06 run used, since no
Modbus/TCP hardware was available.

## Modbus RTU — preparing a QJ71C24N as a fixed-response slave (GX Works2)

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
   `G.OUTPUT` to send the fixed 13-byte response from `D200`….
4. **Gate the response send on a real "ready" condition, not a stale
   interlock.** The working rig's first pass received every request cleanly but
   never sent a response — traced to an **M-relay interlock gating the
   `G.OUTPUT` rung** that was left unsatisfied. Clearing that interlock
   condition is what made the round trip work; if a `modbus-verify` raw exchange
   gets the request in but nothing back, check the send-side interlock before
   suspecting the link.
5. Map Modbus holding registers to the response buffer: Holding Register `0`-`3`
   ↔ `D1000`-`D1003`. The verified rig seeds these `1`, `2`, `3`, `4`.

### Wiring

- USB↔RS-485 adapter `A`/`B` (a.k.a. `D+`/`D−`) to the QJ71C24N's RS-485
  terminals per the module's manual; tie the signal grounds.
- The module answers with its full fixed response (13 bytes / 4 registers)
  regardless of the request's quantity field — a `holding:0` single-register read
  still gets all 4 registers back. `modbus-verify`'s raw-shape check flags this
  as an "invalid" shape (it compares against the *requested* quantity), but
  `ModbusRtuConnection` frames the response by its own declared byte count and
  decodes correctly regardless — a property of the fixed-response rig, not a
  driver defect.

## What the Modbus runs cover

- **RTU**: open the port; a raw request/response on the wire with a CRC and
  exception-code check (`--quantity` reads multiple registers/coils on the raw
  path); a `ModbusRtuConnection` read of one tag.
- **TCP**: connect; a `ModbusConnection` read of one tag with the MBAP exchange
  logged.

Not covered: writes to real hardware, multi-register reads through the driver
(`ModbusTag` and the read-request path are still single-value), `WriteMultiple*`,
sustained polling, more than one slave on the bus, and — on the fixed-response
RTU rig — slave-side address-range and exception-code handling.

## If the Modbus run fails

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

## Modbus run log

### 2026-09-19 — Modbus RTU against a Mitsubishi QJ71C24N

**Verified end-to-end — raw frame exchange and the `ModbusRtuConnection` driver
read both pass.**

- **Master**: `tools/modbus-verify` → raw `ITransportInstance` /
  `ModbusRtuConnection` → `SerialTransportInstance`
- **Adapter**: CH340-based USB↔RS-485, `COM3`
- **Slave**: QJ71C24N CH2, non-procedure communication, 9600-8-E-1, sum check on,
  fixed 8-byte receive length, no terminator
- **Mapping**: Holding Register `0`-`3` ↔ `D1000`-`D1003`, seeded `1`/`2`/`3`/`4`

**What was observed first.** The raw exchange sent a well-formed request and got
zero bytes back, repeatably, across different baud rates and unit ids. The
module's own receive indication showed the request arriving intact every time, so
the request side was never in question.

**Root cause.** A ladder-side M-relay interlock was gating the response-send
(`G.OUTPUT`) rung. The request was received correctly on every attempt; the
interlock simply never let the canned response go out. Once cleared, the round
trip worked immediately at the same baud and wiring that had produced zero bytes
moments before — confirming the earlier failures were a PLC-program gate, not the
link or the driver.

Raw frame exchange, `--quantity 4` (matching the module's fixed response size):

```
→ 01 03 00 00 00 04 44 09                  (unit 1, read holding register 0, qty 4)
← 01 03 08 00 01 00 02 00 03 00 04 0D 14   (13 bytes)
```

Address byte, function code and CRC all valid; decoded register values
`1, 2, 3, 4` match the seeded `D1000`-`D1003`.

Driver read path, default quantity (1), run twice back to back:

```
→ 01 03 00 00 00 01 84 0A
← 01 03 08 00 01 00 02 00 03 00 04 0D 14   (13 bytes — the module's fixed response)
```

`ModbusRtuConnection` frames the response by its own declared byte count (8), not
by the 2 bytes a single-register request would normally imply, and correctly
returns `1` for `holding:0` both times.

This is also the first hardware confirmation of the `SerialTransportInstance`
receive-loop rewrite — polling `BytesToRead` and reading synchronously, in place
of the earlier `_port.BaseStream.ReadAsync(ct)` loop.

### 2026-09-19 — Modbus RTU PLC-side evidence

- [QJ71C24N test rig](images/modbus-qj71c24n-test-rig.jpg) — the isolated bench:
  24 VDC supply, the Mitsubishi Q-series rack (Q64PN power supply, CPU, and the
  QJ71C24N serial module with its RS-485 terminal block wired to the USB
  adapter), alongside the S7-1214C used for the S7 verification.
- [QJ71C24N ladder program](images/modbus-qj71c24n-ladder.png) — the receive and
  send rungs monitored online in GX Works2. Rung 23 gates `G.INPUT` (channel 2,
  receive into `D100`, completion flag `M100`); rung 43 gates `G.OUTPUT`
  (channel 2, send `D200`…, completion flag `M110`) behind the `M2000` interlock
  and the `T200` turnaround timer. `D200` monitors as `769` (`0x0301`) — the
  first two bytes of the canned response, `01 03`.

### 2026-09-06 — Modbus TCP against a software slave

No Modbus/TCP hardware was on hand, so `ModbusConnection` was exercised against a
minimal raw-socket Modbus/TCP slave with four distinct tables (holding, input,
coils, discrete inputs) plus an out-of-range address that returns exception
`0x02`.

- **Read holding register 0** → `4096` rendered as `UINT16` ✅
- **Read input register 0** → `8192` as `UINT16` ✅
- **Read coil 0 / 1** → `False` / `True` ✅
- **Read discrete input 0 / 1** → `True` / `False` ✅
- **Read holding register 200** (outside the map) → `InvalidAddress` ✅

Two bugs were found and fixed (commit `e86bdd028`):

1. **`ModbusConnection` (TCP) only handled `Coil` and `HoldingRegister` reads.**
   `input:` and `discrete:` tags fell through the `switch` to `AccessDenied`
   without a request ever reaching the wire — `ModbusRtuConnection` already
   handled all four tag types. Added the two missing cases, mirroring the RTU
   path.
2. **`modbus-verify` printed every register value as `True (BOOL)`.**
   `PrintValue` probed `IPlcValue.IsBool()` first, and plc4net's value model
   coerces freely between related scalars (a `PlcBOOL` also answers `IsUshort()`
   true; every numeric answers `IsBool()` true), so the probe cannot recover the
   type after the fact. It now renders by the `ModbusTag.TagType` the read was
   issued for.

Regression tests: `Tcp_read_input_register_returns_the_value`,
`Tcp_read_discrete_input_returns_the_value`.

### 2026-09-06 — Modbus RTU against a Siemens S7-1214C + CM 1241 (retired)

**This rig was retired on 2026-09-19 in favour of the QJ71C24N above.** It is
kept here because the failure analysis is the reason the target changed.

- **Slave**: SIMATIC S7-1214C DC/DC/DC + **CM 1241 (RS422/485)**,
  `6ES7 241-1CH32-0XB0` V2.2, half-duplex RS-485 2-wire, running
  `Modbus_Comm_Load` + `Modbus_Slave`, `MB_ADDR` 2, `MB_HOLD_REG` → a
  non-optimized DB word array
- **Line**: 19200-8-E-1, later 9600-8-E-1; `MB_ADDR` 1 and 2 both tried

Across roughly ten runs the tool always put a well-formed request on the wire —
e.g. `02 03 00 00 00 01 84 39` (unit 2, read holding register 0, CRC-16 correct)
— and the raw frame exchange got **zero bytes back**; the driver read timed out.
On the PLC side, `Modbus_Slave.STATUS` was seen cycling between `16#7001` (a
complete, valid request received and being processed) and `16#8280` (a
character-level receive error), so at least some frames reached the CM 1241
intact and were accepted by `MB_SLAVE`. No response ever completed the round trip.

Why this was judged a link problem, not a driver problem:

- The request frames were byte-for-byte correct and `MB_SLAVE` reached `16#7001`
  on them.
- A separate, unrelated Modbus master tool, on the **same adapter and wiring**,
  could not read a register from this slave either.
- The CH340 adapter repeatedly dropped off the USB bus (five-plus times over the
  session), so no test ran against a link that stayed up.
- Idle fail-safe bias measured ~209 mV across A/B — right at the RS-485 receiver
  threshold — with no line termination fitted; the CM 1241's internal bias was
  enabled and the link still read as marginal.
- `Modbus_Comm_Load` is `REQ` edge-triggered, and its instance DB carried state
  across RUN-mode downloads, so a clean rising edge often never reached it and
  the port was left unconfigured (`Modbus_Slave.STATUS = 16#8280`, no traffic).

Closing this rig out, if it is ever revisited, needs an isolated USB↔RS-485
adapter (ADM2483 / ADM2587-class, or an FTDI part with a real direction-control
output) on a terminated pair, with `Modbus_Comm_Load` triggered from
`#Initial_Call` after a genuine STOP→RUN. Independently of the rig,
`ModbusRtuConnection.SendAndReceive` should gain an explicit expected-length /
t3.5 inter-frame-gap check rather than reading whatever is available once ≥ 4
bytes have arrived.

## Change log

- 2026-09-19: Merged the four separate hardware documents
  (`s7-hardware-verification.md`, `s7-hardware-report.md`,
  `modbus-hardware-verification.md`, `modbus-hardware-report.md`) into this one
  file. Content is unchanged apart from deduplicating the shared "run it" and
  packaging sections and adding the scope caveats in *Status*.
- 2026-09-19: Modbus RTU verified against the QJ71C24N rig; the S7-1214C +
  CM 1241 rig retired. Added the QJ71C24N test-rig and ladder images.
- 2026-09-16 20:09: Added the S7-1214C test-rig, DB100 online-value and I/Q/M
  watch-table images as PLC-side evidence for the persistent matrix
  verification. The images are unedited copies of the captured files.
- 2026-09-16 19:53: Verified the persistent I/Q/M/DB matrix against the isolated
  S7-1214C: main run 43/43 and independent new-connection read-back 25/25.
  DB100, M100..M117 and Q0..Q17 were deliberately not restored.
- 2026-09-16 19:24: Added the explicit `--keep-output-values` mode and verified
  persistent Q writes against an isolated S7-1214C with no attached equipment.
- 2026-09-16 19:19: DB100 and M100..M117 passed the expanded 50/50 hardware
  verification; every write was restored and verified.
- 2026-09-16: Expanded the S7 procedure for I/Q/M X/B/W/D addresses and seven
  scalar read-before/write/read-back/restore verification. Added
  `--keep-db-values` / `--keep-marker-values`. Output writes stay behind an
  explicit safety switch.
- 2026-09-06: Modbus TCP verified against a software slave; two bugs fixed.
