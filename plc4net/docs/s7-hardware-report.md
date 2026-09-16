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
# S7 hardware verification — run log

Output of `tools/s7-verify` against real hardware. The procedure and the data
block layout are in [s7-hardware-verification.md](s7-hardware-verification.md).

## 2026-09-16 — Siemens S7-1214C, persistent I/Q/M/DB matrix

The PLC had no external equipment attached. Inputs were read only. DB100,
M100..M117 and Q0..Q17 were each exercised with read-before → write →
immediate read-back, with all three restore steps explicitly disabled. After
the write session closed, every target was independently read through a new
connection.

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

### PLC-side visual evidence

The following images were captured after the persistent write and independent
read-back. They provide PLC-side confirmation in TIA Portal in addition to the
driver-side report above.

- [S7-1214C test rig](images/s7-1214c-test-rig.jpg) — CPU 1214C DC/DC/DC on the
  isolated bench with its 24 VDC power supply.
- [DB100 online values](images/s7-db100-online-values.png) — DB100 offsets,
  declared types, start values and monitored values in TIA Portal.
- [I/Q/M watch table](images/s7-iqm-watch-table.png) — input reads and the
  persistent output/marker values at the tested absolute addresses.

## 2026-09-16 — Siemens S7-1214C, persistent Q write verification

The PLC had no external equipment attached. The output matrix was explicitly
run with `--write-outputs --keep-output-values`, so the final values were not
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

## 2026-09-16 — Siemens S7-1214C, expanded verification

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

## 2026-09-16 — Siemens S7-1214C, expanded re-verification blocked

Attempted from the `feature/plc4net-revival` working tree based on
`e86bdd028`. The expanded harness first ran in `--read-only` mode, so it sent no
Write Var request.

- **Device endpoint**: `192.168.1.11`, rack 0 / slot 1
- **Connection**: COTP + Setup Communication passed; negotiated PDU 240 bytes
- **Read targets**: DB100 scalar values plus I/Q/M X/B/W/D address forms
- **Result**: **BLOCKED** — all 19 reads returned `AccessDenied`
- **Write status**: not attempted; the read gate failed
- **Likely cause**: the CPU's PUT/GET permission is no longer effective. Enable
  "Permit access with PUT/GET communication from remote partner", compile, and
  download the hardware configuration before retrying.

The failure was consistent across DB100, `%I`, `%Q`, and `%M`, while the S7
session stayed connected. This distinguishes an access-policy refusal from an
address parser or individual memory-range failure.

## 2026-09-04 — Siemens S7-1214C (DC/DC/DC), re-verify

Re-run on `feature/plc4net-revival` at `39e3792a0` — after `develop` was merged,
the prerequisite check was raised to the .NET 8 SDK and `AssemblyVersion` was
pinned. No S7 driver code changed since the 2026-09-03 run; this confirms those
build/infra commits did not regress the driver.

- **Device**: SIMATIC S7-1214C DC/DC/DC, rack 0 / slot 1
- **Connection**: `s7://192.168.1.11?remote-rack=0&remote-slot=1&request-timeout=5000` (default TSAPs)
- **Negotiated PDU length**: 240 bytes
- **Data block**: DB100, non-optimized, unchanged from 2026-09-03
- **Result**: **PASS (12/12)**, two consecutive runs. Extra probe: single read of
  `%I0.0` returned `Ok` (`False`) — exercises `%I` area addressing outside the DB.

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

## Result: PASS (12/12)
```

## 2026-09-03 — Siemens S7-1214C (DC/DC/DC)

- **Device**: SIMATIC S7-1214C DC/DC/DC (S7-1200 family), rack 0 / slot 1
- **Connection**: `s7://192.168.1.11` (default TSAPs: remote `0x0101`, local `0x0311`)
- **Negotiated PDU length**: 240 bytes
- **Data block**: DB100, non-optimized, laid out per the procedure doc
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

## Result: PASS (12/12)
```

Covered end to end against the CPU: the COTP CR/CC handshake, S7 Setup
Communication with PDU-length negotiation, single-item reads of every scalar
width (BOOL / BYTE / INT / DINT / REAL / WORD / DWORD) from a data block, a
three-item read in one request, a write + read-back of an INT and a REAL, and
an error path (reading a non-existent DB returns `NotFound` and the connection
stays usable).

### Bug found and fixed during this run

The CPU refuses a request it cannot serve with a bare **Ack (ROSCTR 0x02)**,
which — like an AckData (0x03) — carries a 2-byte `errorClass` / `errorCode`
field after the 10-byte S7 header. `S7Connection.ReadOneS7MessageAsync` only
treated 0x03 as a 12-byte header, so it read 10 bytes of a 12-byte frame, the
parse failed, and the two leftover bytes desynced every following response.
Fixed to frame both 0x02 and 0x03 as 12 bytes, and to map the S7 header errors
(`0x8104` PUT/GET refused, `0x8304`, `0x85xx`) to a per-tag response code the
way the Java driver's `mapPlcErrorCode` does. Regression tests:
`Read_maps_a_header_level_refusal_to_access_denied` and
`A_bare_Ack_is_framed_as_12_bytes_and_does_not_desync_the_next_request`.

### CPU-side prerequisites confirmed necessary

- **"Permit access with PUT/GET communication from remote partner"** must be
  ticked (CPU → Protection & Security → Connection mechanisms), then compiled
  **and downloaded** — the offline setting alone has no effect. Without it every
  Read/Write Var comes back `0x8104` even though Setup Communication succeeds.
- **DB100 must exist and be non-optimized.** A missing block returns `NotFound`
  per item; an optimized block cannot be reached by absolute addressing.

## Change log

- 2026-09-16 19:24: Appended persistent Q write verification after the owner
  confirmed that no external equipment was attached; Q values were intentionally
  left changed and independently read back over new connections.
- 2026-09-16 19:19: Appended the successful 50/50 expanded verification after
  CPU PUT/GET access was restored; all DB and marker writes were rolled back.
- 2026-09-16: Appended the expanded read-only re-verification failure. The
  successful 2026-09-03 and 2026-09-04 records above remain unchanged.
