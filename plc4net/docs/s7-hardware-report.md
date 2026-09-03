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

## 2026-09-03 — Siemens S7-1200 (DC/DC/DC)

- **Device**: SIMATIC S7-1200, DC/DC/DC CPU, rack 0 / slot 1
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
