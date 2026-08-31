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
# plc4net — Design & Roadmap

## Overview

plc4net is the .NET implementation of Apache PLC4X, targeting **net8.0** (LTS).
It provides protocol adapters for industrial controllers — Modbus TCP and
Siemens S7 (ISO-on-TCP / COTP), with KNX Net/IP carried forward from the
original upstream port.

The project was revived from an abandoned state in July 2026.  The original
code targeted .NET Framework 4.5.2, could not build on Linux or macOS CI
runners, had no working COTP transport, and carried a dispatch bug that made the
entire `IPlcValue` interface unreachable through its own type hierarchy.

## Project Structure

```
plc4net/
├── api/                        plc4net-api           (public contract)
├── spi/                        plc4net-spi           (driver/transport base classes)
├── test/                                             (test projects)
│   ├── spi-test/               plc4net-spi-test      (313 xUnit tests)
│   └── knxnetip-test/          plc4net-driver-knxnetip-test
├── drivers/
│   ├── modbus/                 plc4net-driver-modbus     (src/.../readwrite/model is generated)
│   ├── s7/                     plc4net-driver-s7        (src/.../readwrite/model is generated)
│   └── knxnetip/               plc4net-driver-knxnetip   (generated - by the old Java plugin)
├── transports/
│   ├── tcp/                    plc4net-transports-tcp
│   ├── cotp/                   plc4net-transports-cotp
│   ├── serial/                 plc4net-transports-serial  (RS-232 / RS-485)
│   └── test/                   plc4net-transports-test   (loopback)
├── tools/
│   ├── code-gen/               plc4net-code-gen          (mspec -> IR -> C#; run with `dotnet run`)
│   ├── modbus-verify/          modbus-verify             (manual hardware smoke test)
│   └── s7-verify/              s7-verify                 (manual hardware smoke test)
└── docs/
    ├── design.md
    └── testing.md
```

Every project file is named after its directory (`modbus/modbus.csproj`,
`transports/tcp/tcp.csproj`) with an explicit `<AssemblyName>` preserving the
qualified NuGet identity.  All hand-written source files live at the project
root in namespace-based subdirectories — the standard .NET SDK convention.
Generated code (`knxnetip` 179 model files, `code-gen` ANTLR output) keeps its
own layout under `src/`.

### Solution Explorer layout (Visual Studio / Rider)

```
Solution 'plc4net'
├── api/            → api
├── spi/            → spi
├── test/
│   ├── spi-test
│   └── knxnetip-test
├── drivers/
│   ├── knxnetip/   → knxnetip
│   ├── modbus/     → modbus
│   └── s7/         → s7
├── transports/
│   ├── tcp/        → tcp
│   ├── cotp/       → cotp
│   └── test/       → test
└── tools/
    └── code-gen/   → code-gen
```

## Architecture

### Layer diagram

```
┌─────────────────────────────────────────────────────┐
│  api/                                                │
│  PlcDriverManager  ·  IPlcDriver  ·  IPlcConnection  │
│  IPlcTag  ·  IPlcValue  ·  IPlcRead/WriteRequest     │
├─────────────────────────────────────────────────────┤
│  spi/drivers/                                        │
│  DriverBase  ·  ConnectionBase  ·  ConnectionString   │
│  MessageCodecBase<T>  ·  IMessage                     │
│  DefaultPlcReadRequestBuilder                         │
│  DefaultPlcWriteRequestBuilder                        │
├──────────────────┬──────────────────────────────────┤
│  spi/generation/  │  spi/transports/                  │
│  ReadBuffer       │  ITransport  ·  ITransportInstance │
│  WriteBuffer      │  RingBuffer  ·  TransportException │
│  BitReader/Writer │                                    │
├──────────────────┴──────────────────────────────────┤
│  spi/model/values/                                    │
│  PlcBOOL  ·  PlcBYTE  ·  PlcWORD  ·  PlcDWORD  · ... │
├─────────────────────────────────────────────────────┤
│  drivers/modbus/     drivers/s7/                      │
│  ModbusTcpDriver      S7Driver                        │
│  ModbusConnection     S7Connection                    │
│  ModbusTag · PDU      S7Tag · S7PDU · S7Constants    │
│  ModbusADU · Codec    ModbusMessageCodec (template)   │
├─────────────────────────────────────────────────────┤
│  transports/tcp/      transports/cotp/                │
│  TcpTransport          CotpTransport                  │
│  TcpTransportInstance  CotpTransportInstance           │
│  TpkFrame (RFC 1006)                                  │
├─────────────────────────────────────────────────────┤
│  tools/code-gen/                                      │
│  ANTLR C# parsers  ·  MspecModelBuilder  ·  IR        │
│  CSharpModelGenerator  ·  TestsuiteRunner             │
└─────────────────────────────────────────────────────┘
```

### Module dependency graph

```
 api ───────────────────────────────────────────────┐
  │                                                  │
  ▼                                                  │
 spi/drivers ──────► spi/transports                  │
  │                      │                           │
  │                      ├── transports/tcp ──► System.Net.Sockets
  │                      └── transports/cotp ──► transports/tcp
  │                                                  │
  ▼                                                  │
 drivers/modbus ──► spi/drivers  +  spi/transports   │
 drivers/s7     ──► spi/drivers  +  transports/cotp  │
                                                     │
 tools/code-gen ──► Antlr4.Runtime.Standard          │
                     (no runtime dependency)          │
```

- No module depends on a driver module.
- Drivers register transports into the shared `ITransportManager`.
- `tools/code-gen` is a build-time tool; it produces `.cs` files but is not a runtime dependency.

### SPI3 correspondence

| plc4net class | Java counterpart | Notes |
|---|---|---|
| `DriverBase` | `plc4j/spi/…/DriverBase.java` | Transport resolution, connection-string dispatch, supported-transport enforcement |
| `ConnectionBase` | `plc4j/spi/…/ConnectionBase.java` | Wraps transport; holds `ConnectionString`; `Connect()`/`OnConnect()` split mirrors Java's `connect()`/`onConnect()` |
| `MessageCodecBase<T>` | `plc4j/spi/…/MessageCodecBase.java` | TPKT-length-based framing; `IMessage` contract |
| `DefaultPlcReadRequestBuilder` | `DefaultPlcReadRequestBuilder` | Tag name → `IPlcTag` via adapter-specific parser |
| `DefaultPlcWriteRequestBuilder` | `DefaultPlcWriteRequestBuilder` | Same tag-parser injection pattern |
| `ConnectionString` | `DriverBase.URI_PATTERN` | Identical grammar; URL-decodes parameters the way Java does |
| `PlcDriverManager` | `DefaultPlcDriverManager` | Case-insensitive protocol-code lookup; `Uri.Scheme` for dispatch |
| `ITransportInstance` | `TransportInstance` | Extends `IDisposable` (.NET convention) |
| `RingBuffer` | Custom in SPI3 | Fixed-capacity; wrap-around; bulk copy |
| `BitReader`/`BitWriter` | SPI3 byte-buffer codec | MSB-first; no external dependency |
| `ReadBuffer`/`WriteBuffer` | SPI3 `ReadBuffer`/`WriteBuffer` | Numeric, string, array I/O; float32/64, KnxFloat |

### Connection establishment

#### S7 (COTP over TCP)

```
s7://192.168.0.1?remote-rack=0&remote-slot=1&default-port=102
```

```
PlcDriverManager ── protocol code "s7" ──► S7Driver
  │
  ▼
DriverBase.Connect()
  ├─ DefaultTransportCode = "cotp"
  ├─ CotpTransport ──► TcpTransport ──► TCP connect to host:port
  └─ CotpTransportInstance wraps TcpTransportInstance
  │
  ▼
S7Connection ctor                                    (no I/O)
  ├─ Records local-rack/1, local-slot/1, remote-rack/0, remote-slot/0
  │    (+ optional remote-tsap/local-tsap overrides)
  └─ Rejects rack/slot > 0x0F (TSAP reserves 4 bits each)
  │
  ▼
ConnectionBase.Connect() ──► S7Connection.OnConnect() (called by DriverBase)
  ├─ Requires the cotp transport — rejects anything else
  ├─ Computes TSAPs as in Java's S7TsapIdEncoder:
  │    (deviceGroup << 8) | (rack << 4) | (slot & 0x0F)
  │    defaults: local = OTHERS(0x03)/1/1 → 0x0311,
  │              remote = PG_OR_PC(0x01)/0/0 → 0x0100
  ├─ COTP CR (0xE0) ──────────────► PLC              (handshake)
  ├─ COTP CC (0xD0) ◄────────────── PLC              (handshake)
  └─ Ready for S7 PDU exchange
  │
  └─ on failure: DriverBase closes the connection before rethrowing
  │
  ▼
Read ──► S7PDU ──► CotpTransport.Write() ──► COTP DT (0xF0) + TPKT
       ◄── S7PDU ◄── CotpTransport.Read() ◄── COTP DT (0xF0) + TPKT
```

#### Modbus TCP

```
modbus-tcp://192.168.0.9:502?unit-identifier=1
```

```
PlcDriverManager ── protocol code "modbus-tcp" ──► ModbusTcpDriver
  │
  ▼
DriverBase.Connect()
  ├─ DefaultTransportCode = "tcp"
  └─ TcpTransport ──► TcpTransportInstance ──► TCP connect to host:502
  │
  ▼
ModbusConnection ctor
  ├─ UnitId from connection string (default: 1)
  └─ Ready for Modbus PDU exchange
  │
  ▼
Read/Write ──► MBAP header + PDU ──► Write() │ Peek() + Read()
```

## Capabilities

### Currently implemented

| #   | Capability | Layer |
| --- | --- | --- |
| C-1 | Driver framework (DriverBase, ConnectionBase, MessageCodecBase) | spi |
| C-2 | Transport abstraction (ITransport, ITransportInstance) | spi |
| C-3 | TCP transport (TcpTransport, TcpTransportInstance, RingBuffer) | transports/tcp |
| C-4 | Test transport (in-memory loopback for driver testing) | transports/test |
| C-5 | Connection string parser (grammar, parameters, redaction) | spi |
| C-6 | PlcDriverManager (protocol-code dispatch, Uri.Scheme routing) | api |
| C-7 | Value model (PlcBOOL … PlcWSTRING, PlcList, PlcStruct) | spi/model |
| C-8 | Bit-level I/O (BitReader / BitWriter, MSB-first) | spi/generation |
| C-9 | Buffer I/O (ReadBuffer / WriteBuffer, float32/64, KnxFloat) | spi/generation |
| C-10 | Modbus TCP driver (Read/Write, tag parsing, PDU, response parsing) | drivers/modbus |
| C-11 | S7 driver (Read Var, rack/slot addressing, TSAP encoding) | drivers/s7 |
| C-12 | COTP transport (CR/CC handshake, DT framing, TPDU-size negotiation) | transports/cotp |
| C-13 | S7 connection lifecycle (Connect/OnConnect split, TSAP overrides) | drivers/s7 |
| C-14 | COTP peer-disconnect surfacing (DR/DC surfaced as TransportException) | transports/cotp |
| C-15 | Maven reactor bridge (`mvnw -Pwith-dotnet` shells out to `dotnet`; release path only) | pom.xml |
| C-16 | Cross-platform CI — .NET SDK only, no JDK/Maven (ubuntu / macos / windows + header check) | .github |
| C-17 | Pure-.NET mspec generator (ANTLR C# grammars → IR → C# class + StaticParse / Serialize / GetLengthInBits); no JDK, no freemarker | tools/code-gen |
| C-18 | DI integration (IServiceCollection extension) | api/extensions |
| C-19 | KNX Net/IP driver (generated model, DPT 9.x codec) | drivers/knxnetip |
| C-20 | Modbus and S7 wire models, generated and round-tripped against the shared ParserSerializerTestsuite.xml (Modbus TCP 6 vectors, S7 11 vectors) plus the generated `DataItem` `dataIo` (S7 and Modbus scalar / string values, round-tripped from hand-built vectors) | drivers/{modbus,s7}/src/.../readwrite/model |

### Not yet implemented

| Gap | Description |
| --- | --- |
| GAP-1 | Serial transport — RS-485/232 support blocked until serial is implemented |
| GAP-2 | Modbus RTU driver — protocol logic exists, needs serial transport + RTU framing |
| GAP-3 | S7 Write — `WriteRequestBuilder` throws `NotSupportedException` |
| GAP-4 | COTP fragmentation — PDUs > negotiated TPDU size throw instead of splitting |
| GAP-5 | S7 hardware verification — S7-1200/1500 physically available, ICLA filed |
| GAP-6 | UDP transport — needed for EtherNet/IP implicit messaging |
| GAP-7 | TLS transport — needed for OPC UA secure channels |
| GAP-8 | Generator coverage. The generator emits working parse/serialize and both Modbus and S7 round-trip the shared vectors. Handled: all field types except `manualArray` / `peek` / `assert` / `abstract` / `unknown`; discriminated types including a dotted discriminator (`parameter.parameterType`), a `simple` field that is also a discriminator, prefix- and suffix-field inheritance through the constructor (S7's `COTPPacket` trailer), pre-`typeSwitch` context passed to a child's `StaticParse` (S7's `parameterLength` / `payloadLength` / `headerLength`); `curPos` relative to a type's start; `_lastItem` threaded into an array element (S7's `S7VarPayloadDataItem` padding); parameterised enums including a self-reference, `'null'` values (nullable accessor), char / string-literal values, and an `enum` field keyed on an attribute (`TransportSize … code`); repeated-value enums (deduplicated switch arms); bare enum-constant references; `CAST` / `COUNT` / `CEIL` / `ARRAY_SIZE_IN_BYTES` / `STATIC_CALL`. `dataIo` types are emitted as a static `StaticParse` / `StaticSerialize` / `GetLengthInBits` over `IPlcValue` (S7's `DataItem`, Modbus's `DataItem`). S7's `DataItem` is fully generated - the scalar cases (`bit` / `int` / `uint` / `float` widths → `Plc{width}`), `CHAR` / `WCHAR`, `STRING` / `WSTRING` (via the ported `parseS7String` / `serializeS7String`), and the whole TIA date / time family (`TIME` / `S5TIME` / `LTIME` / `DATE` / `TIME_OF_DAY` / `LTIME_OF_DAY` / `DATE_AND_TIME` / `DATE_AND_LTIME` / `DTL`, via the ported `parseS5Time` / `parseTiaDate` / `parseSiemensYear` and a BCD helper). `IPlcValue` grew `GetDuration` / `GetDate` / `GetTime`, `PlcLTIME` / `PlcTIME_OF_DAY` / `PlcDATE` were re-based on `TimeSpan` / `TimeOnly` / `DateOnly`, and `PlcLTIME_OF_DAY` / `PlcDATE_AND_LTIME` are new (nanosecond-exact). Bare `vstring` maps to `string`. Modbus's multi-value `numberOfValues > 1` cases parse to a `PlcList` of the wrapped primitive (`bit → PlcBOOL`, `int 16 → PlcINT`, `string 8 → PlcCHAR`, …), one element per count. Two deliberate departures from plc4j's `DataItem`, both in fields a device consumes but the parser discards: the `DATE_AND_TIME` / `DTL` `dayOfWeek` is written Sunday=1..Saturday=7 (the mspec's own DTL comment and the Siemens doc) where plc4j writes ISO Monday=1..Sunday=7; and a few malformed-frame paths (a `DATE` before 1990, a `DATE_AND_TIME` year past 2089, a non-canonical `S5TIME`) are carried through the same way plc4j carries them rather than rejected. `GetLengthInBits` for `S7_S5TIME` returns 2 (the mspec's `'2'`, a plc4j-shared quirk - it should be 16); no caller sizes an S5TIME item today. Not yet: Modbus's `bigEndian` argument is a no-op for the types that do not carry it as a discriminator (the byte-order-override gap - a `ReadBuffer` / `WriteBuffer` that only does MSB-first); `vstring` with an inline length expression (ADS, not Modbus / S7). KNX additionally needs its `knx-master-data.mspec` (4-parameter enums, thousands of values) fully in the grammar, plus `PlcStruct` `dataIo` emission for the ~700-case `KnxDatapoint`; until then the checked-in KNX model stays as the Java plugin left it (data classes + two hand-adapted dataIo types). |
| GAP-9 | NuGet packaging and publishing |

## Design Principles

### 1. Construction and connection are separate phases

`DriverBase.CreateConnection()` records configuration and performs **no I/O**.
`ConnectionBase.Connect()` dispatches to the `OnConnect()` hook so that a
transport needing a protocol-level handshake (COTP CR/CC for S7) has a
well-defined place to run it.  If the hook fails, `DriverBase` closes the
transport before rethrowing.  Modbus does not override — it carries payload
over the raw transport.

### 2. Tag parser injection

`DefaultPlcReadRequestBuilder` takes a `Func<string, IPlcTag>` supplied by each
connection.  The builder stores tags the connection itself parsed, guaranteeing
the correct concrete type at dispatch time.

### 3. Interface-based value assertions

All value-model tests assert through `IPlcValue`, never through concrete types.
The original `new`-vs-`override` bug (`PlcDINT` declared `public new int
GetInt()` instead of `override`) passed every concrete-type test and failed
through the interface.  Binding tests to the interface catches dispatch
regressions at the source.

### 4. COTP handshake at the transport level

The CR/CC handshake lives in `CotpTransportInstance.Open()`.  The S7 driver sees
a ready-to-use transport — it does not construct or validate COTP frames.

### 5. Byte-level test data from protocol references

Modbus PDU test vectors derive from the Modbus Application Protocol
Specification v1.1b.  S7 test vectors derive from Wireshark captures of S7-1500
communication.  The COTP CC fixture mirrors an S7-1500 response.

### 6. TestTransport for driver testing

An in-memory loopback transport lets drivers read pre-injected byte sequences
and captures written bytes for assertion, without TCP sockets.

### 7. C#-native code generator

`tools/code-gen/` parses `.mspec` with ANTLR-for-C# (the checked-in parser,
so no ANTLR tool and no JDK), walks it into a type-model IR, and emits the
model class, `StaticParse`, `Serialize` and `GetLengthInBits` for each type -
and, for a `[dataIo]`, a static `StaticParse` / `StaticSerialize` /
`GetLengthInBits` over `IPlcValue`.
The Java freemarker generator it replaces only ever produced data classes.
The language-neutral `ParserSerializerTestsuite.xml` and the `.mspec`
descriptions stay shared; the `.g4` grammars are a near-verbatim copy of the
upstream `code-generation` ones (one lexer predicate ported from the ANTLR
Java runtime API to the C# one — see `tools/code-gen/README.md` for the
regen procedure), and the toolchain around them is per-language, which is the
direction plc4j and plc4go are also heading. The proof that it works is
`ModbusGeneratedRoundTripTests` / `S7GeneratedRoundTripTests`, which run the
shared XML vectors — the same bytes plc4j and plc4go validate — through the
generated C# and assert a byte-identical round trip, plus
`S7DataItemRoundTripTests` / `ModbusDataItemRoundTripTests` for the `dataIo`
value codec.

Generated code is checked in (like the KNX model) and a CI job regenerates
it and fails on any drift from the mspec.

## Deliberate divergences from the Java SPI3

| Divergence | Rationale |
| --- | --- |
| C# properties instead of `getXxx()` | .NET convention (`IsOpen`, `ProtocolCode`, `TransportCode`) |
| `ITransportInstance : IDisposable` | .NET resource management; enables `using` blocks |
| `async`/`await` TCP read loop | .NET has no virtual threads; async socket I/O is the equivalent |
| Rack/slot validation rejects out-of-range values | Java's `S7TsapIdEncoder` silently masks; connecting to the wrong CPU is worse than an error message |
| `ConnectionString` in `spi/` (not `api/`) | `PlcDriverManager` only needs `Uri.Scheme`; all real consumers are in the SPI layer |

## Roadmap

### Completed (2026-07 to 2026-08)

| Item | Description |
| --- | --- |
| net8.0 migration | Target net8.0 (LTS), drop net452 |
| Build infrastructure | Directory.Build.props; CI is the .NET SDK only (no JDK, no Maven); ASF-header check; a regen-drift check |
| JDK-free toolchain | The mspec generator is pure .NET; `pom.xml` is a thin `dotnet` bridge for the release reactor |
| Value model fix | `new` → `override` across the PlcValue hierarchy |
| BitReader / BitWriter | MSB-first, zero external dependencies |
| SPI3 alignment | `IPlcField` → `IPlcTag`, synchronous driver API, `Connect()`/`OnConnect()` split |
| Driver runtime | DriverBase, ConnectionBase, MessageCodecBase, Transport lifecycle |
| Modbus TCP driver | Read/Write, tag parsing, PDU construction, response parsing |
| S7 driver | Read Var, rack/slot addressing, TSAP encoding, tag parsing (7 address forms) |
| COTP transport | CR/CC handshake, DT framing, TPDU-size negotiation, peer-disconnect surfacing |
| Test transport | In-memory loopback for driver testing without TCP |
| Code-gen | mspec expression engine, type-model IR, and the C# model + parse/serialize emitter (including `dataIo` → `IPlcValue`); Modbus and S7 generated and round-tripped against the shared test vectors; the Java `update-generated-code` plugin removed |
| Project structure | Unified csproj naming, solution-folder nesting, flat source layout, test/ directory |
| Documentation | Design doc and testing report under `docs/` |
| ICLA | Filed and acknowledged (2026-08-02) |

### In progress

| Priority | Item | Status |
| --- | --- | --- |
| P0 | CI first-run approval | Awaiting committer to approve workflow run |
| P0 | PR #2656 title / draft removal | Pending review |

### Planned — Phase 1: the factory floor essentials

Every item in this phase is driven by what an automation engineer reaches for
first when commissioning a panel.

| # | Item | Transport needed | Notes |
| --- | --- | --- | --- |
| 1 | **Serial transport** | — | RS-485 / RS-232.  The single most important missing transport.  Unlocks Modbus RTU and every serial-attached device. |
| 2 | **Modbus RTU driver** | Serial | Reuses the existing Modbus TCP protocol logic with RTU framing (CRC, inter-frame gap).  The dominant protocol on factory floors. |
| 3 | **S7 Write** | — (COTP/TCP exists) | Closes the read/write loop for Siemens PLCs. |
| 4 | **COTP fragmentation** | — | PDUs larger than the negotiated TPDU size split across DT frames.  Needed for S7 blocks larger than ~240 bytes. |
| 5 | **S7 hardware verification** | — | End-to-end read/write against a real S7-1200 / S7-1500.  Hardware available. |

### Planned — Phase 2: the plant backbone

Once the shop-floor basics are covered, the next reach is into the networks that
tie cells and lines together.

| # | Item | Transport needed | Notes |
| --- | --- | --- | --- |
| 6 | **UDP transport** | — | Unlocks EtherNet/IP implicit messaging and Profinet DCP. |
| 7 | **EtherNet/IP driver** | TCP + UDP | Rockwell / Allen-Bradley ecosystem.  Explicit messaging over TCP; implicit (CIP I/O) over UDP. |
| 8 | **Profinet driver** | TCP + raw-socket | Siemens Profinet ecosystem.  Complements the S7 driver for non-S7 devices on the same wire. |

### Planned — Phase 3: broader reach

| # | Item | Transport needed | Notes |
| --- | --- | --- | --- |
| 9 | **TLS transport** | — | Unlocks OPC UA secure channels. |
| 10 | **OPC UA driver** | TCP + TLS | Cross-vendor interoperability.  The highest-value industrial protocol without a native Siemens or Rockwell roadmap. |
| 11 | **CANopen driver** | CAN / SocketCAN | Motion control, drives, encoders.  Common in packaging and printing lines. |

### Transport dependency matrix

```
                    Modbus RTU  ──── Serial ──── RS-485/232
                    Modbus TCP  ──── TCP ─────── have it
                    S7          ──── COTP/TCP ── have it
                    EtherNet/IP ──── TCP+UDP ──── need UDP
                    Profinet    ──── TCP ──────── have it
                    OPC UA      ──── TCP+TLS ──── need TLS
                    CANopen     ──── CAN ──────── need CAN
```

### Supporting work

| Item | Description |
| --- | --- |
| Generator coverage (GAP-8) | The generator emits parse/serialize/length for every type; S7's and Modbus's `DataItem` `dataIo`s are fully generated (scalar, string, TIA date / time, and multi-value `PlcList` cases); Modbus + S7 round-trip the shared vectors. Remaining: the `peek` / `assert` / `manualArray` field families (unused by Modbus / S7), `vstring` with an inline length (ADS), byte-order overrides, and full KNX (needs its grammar + `PlcStruct` `dataIo` emission). Then move the hand-written Modbus / S7 codecs onto the generated models. |
| NuGet packaging | Package, sign, publish to NuGet.org for direct `dotnet add package` consumption. |
| Nullable → enable | Triage remaining warnings, enable project-wide. |
