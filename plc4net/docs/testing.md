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
# plc4net — Testing

## Summary

| | |
|---|---|
| Test projects | 2 |
| Test framework | xUnit.net |
| Total test cases | **337** |
| Passing | 337 |
| Failing | 0 |
| Build warnings | 0 (`dotnet build --no-incremental`) |
| CI matrix | ubuntu / macos / windows (`.github/workflows/dotnet-platform.yml`), .NET SDK only |

> Counts are xUnit **test cases**: a `[Theory]` counts once per `[InlineData]`
> row, matching the count that `dotnet test` reports as "passed".

## Test projects

| Project | Assembly | Tests | What it covers |
|---|---|---|---|
| `test/spi-test` | `plc4net-spi-test` | 336 | SPI framework, value model (incl. `PlcStruct` / `PlcRawByteArray`), bit/buffer I/O, transports, Modbus driver, S7 driver, code-gen pipeline (incl. `dataIo` struct cases, external enums, hyphenated ids), Modbus + S7 generated round-trip, `DataItem` `dataIo` round-trip, DI extensions |
| `test/knxnetip-test` | `plc4net-driver-knxnetip-test` | 1 | KNX DPT 9.x 16-bit float codec |

Both projects sit under the `test/` solution folder in Visual Studio, following
the standard .NET `src`/`test` convention at the solution level.

## How to run

### Daily development

```bash
cd plc4x
dotnet build plc4net/plc4net.sln
dotnet test  plc4net/plc4net.sln
```

### Clean build (CI equivalent)

An incremental build skips up-to-date projects and can hide warnings.  Always
use `--no-incremental` for a CI-equivalent measurement:

```bash
dotnet build plc4net/plc4net.sln --no-incremental -v minimal
dotnet test  plc4net/plc4net.sln --no-build
```

### Maven reactor (release path only)

`dotnet build` / `dotnet test` on `plc4net.sln` is the whole story for
development and for CI.  The one thing they do not run is the `apache-rat`
license-header check; the CI workflow reproduces that check directly (see
*Continuous Integration* below), so a missing header is still caught in the PR.

The module keeps a `pom.xml` so the release reactor still covers it:

```bash
cd plc4x
./mvnw -B -P'with-dotnet' -pl :plc4net -am install   # shells out to `dotnet`
```

This is the only path that needs a JDK, and only a release manager runs it.

### Filtering

```bash
dotnet test plc4net/spi-test --filter "FullyQualifiedName~CotpTransport"
```

## Coverage matrix

> Paths are relative to the project directory after the August 2026 source
> layout unification (elimination of `src/` and nested `test/` layers).

### SPI framework

| Area | Tests | File |
|---|---|---|
| Connection string parser | 18 | `spi-test/drivers/ConnectionStringTests.cs` |
| Driver manager (registration, lookup, dispatch, error paths) | 7 | `spi-test/drivers/PlcDriverManagerTests.cs` |
| Driver base (transport resolution, Connect/OnConnect, error wrapping) | 9 | `spi-test/drivers/DriverBaseTests.cs` |
| Message infrastructure (read/write request builders, responses through TestTransport) | 4 | `spi-test/drivers/MessagesTests.cs` |
| DI extensions | 2 | `spi-test/drivers/Plc4NetExtensionsTests.cs` |

### Value model

| Area | Tests | File |
|---|---|---|
| All concrete types through `IPlcValue` interface | 12 | `spi-test/model/values/PlcValueTests.cs` |

> Every assertion goes through the `IPlcValue` interface, never a concrete type.
> The original `new`-vs-`override` dispatch bug (`PlcDINT` declared `public new
> int GetInt()` instead of `override`) passed every concrete-type test and
> failed every interface-typed call.
>
> The bit-string values (`PlcBYTE` / `PlcWORD` / `PlcDWORD` / `PlcLWORD`) had a
> sibling defect: `PlcBitString`'s constructors discarded their argument and the
> subtypes exposed only the individual bits, so `GetByte()` / `GetUshort()` etc.
> answered `0`. `PlcList` had the same shape — it stored its values but never
> overrode `IsList()` / `GetList()` / `GetLength()` / `GetIndex()`. Both now
> serve their value; the `DataItem` `WORD` / `DWORD` / `LWORD` scalar and the
> Modbus multi-value list round-trip vectors cover them.
>
> The temporal values were also incomplete — no value could be read back, and
> `PlcLTIME_OF_DAY` / `PlcDATE_AND_LTIME` did not exist. `PlcTIME` / `PlcDATE` /
> `PlcTIME_OF_DAY` are now `TimeSpan` / `DateOnly` / `TimeOnly`, `PlcLTIME` and
> the two new types keep nanoseconds exactly, and `IPlcValue` exposes
> `GetDuration` / `GetDate` / `GetTime`. `PlcDATE` / `PlcDATE_AND_TIME` /
> `PlcDATE_AND_LTIME` also now report `IsSimple()` (they extended the wrong
> base), and every temporal `GetDateTime()` is `DateTimeKind.Unspecified`
> (matching plc4j's zoneless semantics). A related fix: `PlcREAL` used to
> throw for a NaN / ±Inf value because `IsFloat()` rejected non-finite
> numbers. The S7 `DataItem` date / time and non-finite-float vectors cover
> all of this.

### Bit codec

| Area | Tests | File |
|---|---|---|
| MSB-first I/O, ulong/long round-trip, sign extension, Float32/64, KnxFloat, string, byte array, overflow, edge cases | 32 | `spi-test/generation/BufferTests.cs` |

### Ring buffer

| Area | Tests | File |
|---|---|---|
| Round-trip, wrap-around, overflow, read-past-end, zero capacity, clear | 11 | `spi-test/transports/RingBufferTests.cs` |

### Transports

| Area | Tests | File |
|---|---|---|
| TCP transport configuration | 11 | `spi-test/transports/TcpTransportConfigurationTests.cs` |
| TCP transport address parsing | 7 | `spi-test/transports/TcpTransportAddressTests.cs` |
| TCP transport instance (real socket loopback) | 7 | `spi-test/transports/TcpTransportInstanceTests.cs` |
| Test transport (in-memory loopback) | 8 | `spi-test/transports/TestTransportTests.cs` |
| COTP transport instance (CR/CC + DT) | 21 | `spi-test/transports/CotpTransportInstanceTests.cs` |

The 21 COTP tests cover: CR frame layout (byte-level TSAP assertions), handshake
idempotence, CC rejection (wrong PDU type, bad TPKT version), non-echoing
DST-REF acceptance, closed inner transport fast-fail, peer close during
handshake, Read/Write rejected before handshake, DT frame wrap and strip, unknown
PDU type rejection, stray Confirm skipped, Disconnect Request surfaced, partial
frame retained, concurrent Open isolation, split header/body arrival, TPDU-size
negotiation (default 1024, negotiated 512), and Close propagation.

> TCP transport tests use `TcpListener` on `127.0.0.1:0` — real sockets, not
> mocks.  COTP frame layouts follow ISO 8073 (COTP) and RFC 1006 (TPKT); the CC
> fixture mirrors an S7-1500 response.  Each COTP fix was mutation-checked:
> reverting it turns the matching test red.

### Drivers

| Area | Tests | File |
|---|---|---|
| Modbus driver | 10 | `spi-test/drivers/ModbusDriverTests.cs` |
| S7 driver | 18 | `spi-test/drivers/S7DriverTests.cs` |

The Modbus tests cover: tag parsing, Read Coils / Read Holding Registers PDU
construction, Write Single/Multiple PDU construction, and response parsing.
PDU test vectors derive from the Modbus Application Protocol Specification
v1.1b.

The S7 tests cover: tag parsing for all seven address forms (DB, M, I, Q, C,
T, plus bit offsets), Read Var / Write Var PDU construction, TPKT frame
wrap/unwrap, transport size mapping, Java-parity TSAP encoding (remote 0x0101,
local 0x0311), explicit TSAP override, non-COTP transport rejection,
out-of-range rack/slot rejection, and non-numeric parameter rejection.  PDU
test vectors derive from Wireshark captures of real S7-1500 communication.

### Code generation

| Area | Tests | File |
|---|---|---|
| Expression engine — literals, precedence, identifier segments, built-ins, scope rewriting, every modbus.mspec formula | 37 | `spi-test/codegen/ExpressionTests.cs` |
| Pipeline — IR shape (fields, discriminated children, enums), and the emitter (parse/serialize/length, implicit fields) | ~11 | `spi-test/codegen/PipelineTests.cs` |
| mspec parser, testsuite loader | ~7 | `spi-test/codegen/` (MspecParserTests, TestsuiteRunnerTests) |
| **Modbus round-trip** — the six shared TCP `ParserSerializerTestsuite.xml` vectors, parsed and serialized back byte-identical through the generated model | 7 | `spi-test/codegen/ModbusGeneratedRoundTripTests.cs` |
| **S7 round-trip** — the eleven shared `ParserSerializerTestsuite.xml` vectors (TPKT → COTP → S7 message → parameter / payload), parsed and serialized back byte-identical through the generated model | 12 | `spi-test/codegen/S7GeneratedRoundTripTests.cs` |
| **`DataItem` `dataIo` round-trip** — the generated `IPlcValue` codec: S7 keyed on `dataProtocolId`, Modbus keyed on `dataType` / `numberOfValues` / `bigEndian`. Every S7 case round-trips byte-identical from a hand-built IEC-61131 vector — scalars (incl. negative / NaN / ±Inf floats), strings, and the whole TIA date / time family (`S5TIME` canonicalisation, the Siemens epoch, BCD, nanosecond-exact `LTIME` / `DTL`), with value-level assertions on each. Modbus covers the scalar and `numberOfValues > 1` list cases (a `PlcList` of the wrapped primitive), and the little-endian variant (`bigEndian` false → the multi-byte value is byte-swapped) | 40 + 27 | `spi-test/codegen/{S7,Modbus}DataItemRoundTripTests.cs` |

> The generated Modbus and S7 models are compiled for real by
> `plc4net-driver-modbus` / `plc4net-driver-s7`, so a bad emit is a build
> failure. The `ParserSerializerTestsuite.xml` round-trip vectors are the same
> data plc4j and plc4go validate against; the `DataItem` vectors are built from
> the wire layout the mspec describes (no shared suite exercises a typed
> value). A CI job also regenerates both models and fails on any drift from
> the mspec.

### KNX driver

| Area | Tests | File |
|---|---|---|
| DPT 9.x 16-bit float codec | 1 | `test/knxnetip-test/knxnetip/readwrite/model/KnxDatapointTests.cs` |

## Continuous Integration

`.github/workflows/dotnet-platform.yml` runs on every push — the .NET SDK
only, no JDK and no Maven:

- `license-headers` — the same ASF-header check `apache-rat` does, run
  directly against the tracked plc4net sources
- `generated-code-is-current` — regenerates the Modbus and S7 models with
  `tools/code-gen` and fails on `git diff` (the mspec changed without a
  regen, or generated code was hand-edited)
- `test` matrix — **ubuntu-latest** · **macos-latest** · **windows-latest**
  - `dotnet restore` → `dotnet build --no-incremental` → `dotnet test --no-build`

The Maven reactor (`./mvnw -P'with-dotnet' -pl :plc4net -am install`) shells
out to the same `dotnet` commands and is exercised on the release path; it is
not part of PR validation.

## Hardware verification

**S7-1500 (Siemens)** — pending.  The hardware (S7-1200/1500) is physically
available.  The ICLA was filed and acknowledged on 2026-08-02.  The COTP
handshake and S7 Read Var PDU have been verified against reference packet
captures.  The end-to-end path through a real PLC remains to be run:

```
s7://<ip>?remote-rack=0&remote-slot=1&default-port=102
```

## What is deliberately not tested

| Area | Reason |
|---|---|
| Modbus against real hardware | Simulated via TestTransport; byte-identical response for unit-test purposes. Hardware verification planned. |
| S7 Write | `WriteRequestBuilder` throws `NotSupportedException`. |
| Subscribe / Browse / Ping / Discovery | Interfaces declared; no implementation exists. |
| COTP PDU fragmentation | Payloads exceeding the negotiated TPDU size throw `TransportException`. |
| TLS, Serial, UDP transports | Not yet implemented. |
