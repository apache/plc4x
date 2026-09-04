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
# S7 hardware verification

`tools/s7-verify` drives the **public driver API** exactly as a NuGet consumer
would — `new S7Driver(new DefaultTransportManager()).Connect("s7://host?…")` —
then connects, runs S7 Setup Communication, reads every scalar type from a data
block, round-trips a write, and checks an error path. It prints a Markdown
report and exits 0 on pass, 1 on any failure.

Verified against: **Siemens S7-1214C DC/DC/DC (S7-1200 family), 2026-09-03 — PASS 12/12.**
See [s7-hardware-report.md](s7-hardware-report.md) for the run log.

## 1. Prepare the PLC (TIA Portal)

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
5. ⚠️ The write test only writes to the DB. It never touches `%Q` — on a
   DC/DC/DC CPU those are real transistor outputs.

## 2. Run it

### From the source tree

```bash
cd plc4x/plc4net
dotnet run --project tools/s7-verify -- <PLC-IP> --db 100 > docs/s7-hardware-report.md
```

### As a packaged tool (the "build a package, import, use it" path)

```bash
cd plc4x/plc4net
dotnet pack tools/s7-verify -c Release -o ./_localfeed
dotnet tool install --global --add-source ./_localfeed s7-verify
s7-verify <PLC-IP> --db 100 > s7-hardware-report.md
```

`dotnet pack` on a `PackAsTool` project bundles every dependency into the one
`.nupkg`, so the tool is self-contained — no other package needs to be in the
feed. Nothing is published anywhere; `_localfeed` is a folder on your machine.

### Options

```
s7-verify <host> [--rack N] [--slot N] [--db N]
          [--device-group PG_OR_PC|OS|OTHERS] [--remote-tsap 0xNNNN]
          [--read <address>]
```

Defaults: `--rack 0 --slot 1 --db 100`. `--read <address>` skips the suite and
just connects, reads that one tag and prints the outcome — a focused probe for
one address (e.g. `--read "%I0.0"`, `--read "%DB100.DBW2"`) that needs neither
DB100 nor the full layout.

## 3. If the connection fails

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

Once a `--device-group` / `--remote-tsap` value works, it belongs in the
driver's documentation and, ideally, as the S7-1200/1500 default.

## 4. What it covers

- COTP CR/CC handshake, S7 Setup Communication, negotiated PDU length
- Read: BOOL, BYTE, INT, DINT, REAL, WORD, DWORD from a DB; a 3-tag single
  request
- Write + read-back: an Int and a Real to the DB
- Error path: reading a non-existent DB, connection survives

Not covered yet: reads larger than one negotiated PDU (multi-PDU), `%M`/`%I`/`%Q`
word access (the tag parser only takes `%M0` / `%M0.0`), STRING, the TIA
date/time types, subscriptions.
