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
# Modbus hardware verification — run log

Output of `tools/modbus-verify`. The procedure, wiring and the TIA Portal slave
setup are in [modbus-hardware-verification.md](modbus-hardware-verification.md).

## 2026-09-06 — Modbus TCP against a software slave

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

## 2026-09-06 — Modbus RTU against a Siemens S7-1214C + CM 1241 (RS422/485)

**Status: RTU request framing verified on the wire; a full read/write round-trip
against this slave is still pending a stable RS-485 bench link.**

- **Master**: `tools/modbus-verify` → `ModbusRtuConnection` → `SerialTransportInstance`
- **Adapter**: CH340-based USB↔RS-485, `COM3`
- **Slave**: SIMATIC S7-1214C DC/DC/DC + **CM 1241 (RS422/485)**,
  `6ES7 241-1CH32-0XB0` V2.2, port set to half-duplex RS-485 2-wire
- **Program**: `Modbus_Comm_Load` + `Modbus_Slave`, `MB_ADDR` 2,
  `MB_HOLD_REG` → a non-optimized DB word array
- **Line**: 19200-8-E-1, later 9600-8-E-1; `MB_ADDR` 1 and 2 both tried
- **Wiring**: CM 1241 pin 3 / pin 8 / pin 5 (B / A / GND) to the adapter

### What was observed

Across roughly ten runs the tool always put a well-formed request on the wire —
e.g. `02 03 00 00 00 01 84 39` (unit 2, read holding register 0, CRC-16 correct) —
and the raw frame exchange got **zero bytes back**; the driver read timed out.

On the PLC side, `Modbus_Slave.STATUS` was seen cycling between `16#7001`
(a complete, valid request received and being processed) and `16#8280`
(a character-level receive error) — so at least some of plc4net's frames did
reach the CM 1241 intact and were accepted by `MB_SLAVE`. No response ever
completed the round trip back to the master.

### Why this is a link problem, not a driver problem

- The request frames are byte-for-byte correct (function code, address,
  quantity, CRC) and `MB_SLAVE` reached `16#7001` on them.
- A separate, unrelated Modbus master tool, on the **same adapter and wiring**,
  could not read a register from this slave either.
- The CH340 adapter repeatedly dropped off the USB bus (five-plus times over the
  session), so no test ran against a link that stayed up.
- Idle fail-safe bias measured ~209 mV across A/B — right at the RS-485 receiver
  threshold — with no line termination fitted; the CM 1241's internal bias was
  enabled and the link still read as marginal.
- A protracted S7-side issue: `Modbus_Comm_Load` is `REQ` edge-triggered, and its
  instance DB carried state across RUN-mode downloads, so a clean rising edge
  often never reached it and the port was left unconfigured
  (`Modbus_Slave.STATUS = 16#8280` with no traffic).

### To close this out

An isolated USB↔RS-485 adapter (ADM2483 / ADM2587-class, or an FTDI part with a
real direction-control output) on a terminated pair, with `Modbus_Comm_Load`
triggered from `#Initial_Call` after a genuine STOP→RUN. The driver-side
expectation is a single `holding:0` read returning the DB's seeded value.
While this is revisited, `ModbusRtuConnection.SendAndReceive` should also gain an
explicit expected-length / t3.5 inter-frame-gap check rather than reading
whatever is available once ≥ 4 bytes have arrived.

**Superseded 2026-09-19** — the verification target moved to a Mitsubishi
QJ71C24N (see below) rather than continuing to chase this RS-485 bench link.

## 2026-09-19 — Modbus RTU against a Mitsubishi QJ71C24N (non-procedure communication)

**Status: verified end-to-end — raw frame exchange and the `ModbusRtuConnection`
driver read both pass.**

The QJ71C24N has no native Modbus RTU slave firmware, so it stands in as a
**fixed-response slave**: its ladder program receives the master's request via
`G.INPUT` and answers with the same canned response regardless of the request's
function code or address — it does not decode Modbus, it is a framing/wire-level
stand-in.

- **Master**: `tools/modbus-verify` → raw `ITransportInstance` / `ModbusRtuConnection`
  → `SerialTransportInstance`
- **Adapter**: CH340-based USB↔RS-485, `COM3` (the same adapter used for the
  S7-1214C attempt above)
- **Slave**: QJ71C24N CH2, non-procedure ("no-protocol") communication,
  9600-8-E-1, sum check on, fixed 8-byte receive length, no terminator
- **Mapping**: Holding Register `0`-`3` ↔ `D1000`-`D1003`, seeded `1`/`2`/`3`/`4`

### What was observed first

Same symptom as the S7-1214C rig: the raw exchange sent a well-formed request
and got zero bytes back, repeatably, across different baud rates and unit ids.
Unlike the S7 rig, the module's own receive indication showed the request
arriving intact every time — so the request side was never in question.

### Root cause and fix

A **ladder-side M-relay interlock was gating the response-send (`G.OUTPUT`)
rung**. The request was received correctly on every attempt; the interlock
simply never let the canned response go out. Once that condition was cleared,
the round trip worked immediately, at the same baud/wiring that had produced
zero bytes moments before — confirming the earlier failures were a PLC-program
gate, not the link or the driver.

### Verified 2026-09-19

Raw frame exchange, `--quantity 4` (matching the module's fixed response size):

```
→ 01 03 00 00 00 04 44 09    (unit 1, read holding register 0, qty 4)
← 01 03 08 00 01 00 02 00 03 00 04 0D 14    (13 bytes)
```

Address byte, function code and CRC all valid; decoded register values
`1, 2, 3, 4` match the seeded `D1000`-`D1003`.

Driver read path, default quantity (1), run twice back to back:

```
→ 01 03 00 00 00 01 84 0A
← 01 03 08 00 01 00 02 00 03 00 04 0D 14    (13 bytes — the module's fixed response)
```

`ModbusRtuConnection` frames the response by its own declared byte count (8),
not by the 2 bytes a single-register request would normally imply, and
correctly returns `1` for `holding:0` both times. The raw tool's own
"expected shape" check flags this run as invalid, because it compares the
byte count against the *requested* quantity rather than the module's actual
(always-4-register) reply — a property of this fixed-response rig, not a
driver defect.

This is also the first hardware confirmation of the `SerialTransportInstance`
receive-loop rewrite — polling `BytesToRead` and reading synchronously, in
place of the earlier `_port.BaseStream.ReadAsync(ct)` loop. Both the S7-1214C
zero-byte runs and this rig's initial zero-byte runs happened with the old and
new loop respectively without changing the symptom, which is consistent with
those failures being link/interlock issues rather than a receive-loop bug —
this run is the first case where bytes actually arrived, and the new loop
delivered them correctly.
