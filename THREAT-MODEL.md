# Apache PLC4X Security Threat Model (draft)

## §1 Header

- **Project**: Apache PLC4X (`apache/plc4x`) — a multi-language library of
  client-side protocol drivers for industrial / OT (operational-technology)
  PLC protocols (Siemens S7, Modbus, EtherNet/IP, OPC UA, ADS/AMS, KNXnet/IP,
  BACnet/IP, IEC-60870-5-104, CANopen, C-Bus, Firmata, Profinet, UMAS, etc.)
  *(documented: `README.md`, `website/asciidoc/modules/users/pages/protocols/index.adoc`)*.
- **Scope of this model**: the **main `apache/plc4x` repository only** —
  driver layer, SPI, transports, code-generation runtime, and the four
  language bindings shipped here (`plc4j`, `plc4go`, `plc4c`, `plc4py`).
  The separately-released `plc4x-extras` repository (OPC-UA Server, PLC4X
  Server, Calcite / Camel / Kafka-Connect / Karaf / NiFi integrations)
  is **out of scope** of this document *(maintainer: Chris Dutz,
  2026-05-29 16:06Z)* and will be modelled separately if and when it is
  brought into the program. **Correction (2026-09-03):** earlier
  revisions also placed the connection cache, OPM and the scraper in
  `plc4x-extras`. They are in **this** repository, under `plc4j/tools/`
  — the connection cache as `plc4j-tools-connection-cache`
  (re-implemented and renamed from `CachedPlcConnectionManager` to
  `PlcConnectionCache`), OPM as `plc4j-tools-opm` (ported to SPI3), and
  the scraper replaced by the event-pump, `plc4j-tools-event-pump`.
  All are built and deployed from here, so they are **in** scope; see
  the tools row of the §2 component table.
- **Version / commit**: drafted against the default branch (`develop`),
  HEAD `32b4f0c` ("build(deps): bump jackson.version from 2.21.3 to 2.21.4"
  as of draft time); re-checked against `develop` at `173fe05`
  (2026-09-03), which is after the SPI3 migration, the OPC UA
  secure-by-default work and the parser-hardening series. A
  vulnerability report against PLC4X release *N* is triaged against the
  model as it stood at *N*, not at HEAD.
- **Date**: 2026-05-30; §5a and everything downstream of it (§4
  B2-OPCUA, §7, §8 P1/P2, §9, §10, §11, §11a, §13, §14 Q14–Q16)
  re-checked against the shipped code on **2026-09-03**, after the OPC
  UA secure-by-default work landed.
- **Authors**: ASF Security team draft, awaiting PLC4X PMC review.
- **Status**: draft — under maintainer review.
- **Reporting cross-reference**: vulnerabilities that may violate a §8
  property should be reported to the Apache Security Team
  (`security@apache.org`) per
  `website/asciidoc/modules/users/pages/security.adoc`; reports that fall
  under §3, §9, or §11a will be closed by PLC4X triagers citing this
  document.
- **Provenance legend** —
  *(documented)* = drawn from in-repo docs, the website asciidoc tree, or
  on-the-wire protocol specifications, with citation;
  *(maintainer)* = stated by a PLC4X maintainer (PMC chair Chris Dutz or
  another committer) in response to this draft;
  *(inferred)* = synthesized by the producer from code structure, the
  protocol's published specification, or general OT-protocol domain
  knowledge, awaiting PMC ratification — every *(inferred)* tag has a
  matching §14 question.
- **Draft confidence**: 51 documented / 4 maintainer / 45 inferred.

**About the project.** PLC4X provides a uniform client-side API
(`PlcConnection` in Java; analogous in Go, C, Python) on top of a fleet of
**protocol drivers**, each implementing one specific PLC wire protocol.
An embedding application opens a connection by URL
(e.g. `s7://10.0.0.5:102?remote-rack=0&remote-slot=1`,
`modbus-tcp://10.0.0.5:502`,
`opcua:tcp://10.0.0.5:12686?discovery=true&security-policy=Basic256Sha256`),
optionally supplies a `PlcAuthentication`, and then issues `read` /
`write` / `subscribe` / `browse` calls against tag-address strings whose
syntax is **protocol-specific** (S7 area+offset, Modbus coil/register
address, OPC UA `ns=...;s=...` node identifier, …). The library is a
**library** — there is no daemon, no listening socket on the embedding
side, no per-end-user identity. The "session" is the PLC connection.

PLC4X's domain is **industrial network protocols**. Most of these protocols
were standardized in the 1980s–2000s before OT-network security was a
design concern: Modbus, S7 (PUT/GET), BACnet/IP, IEC-60870-5-104,
Profinet, CANopen, KNXnet/IP, C-Bus, DF1, AB-Ethernet, ADS/AMS, EtherNet/IP,
and the historical UMAS dialect of Modicon all transit cleartext on the
wire with **no authentication, no integrity, and no confidentiality** —
that is the protocol as standardized, not a PLC4X choice. OPC UA, the
youngest of the supported set, is the lone exception: it carries a full
TLS-like security stack (certificates, signed/encrypted message
exchange, security policies) and the OPC UA driver here actually
implements that stack. This asymmetry shapes essentially every claim in
the rest of the document.

## §2 Scope and intended use

### Intended use

- **Embedded-library access to PLC devices over their native protocols**,
  from a host application (SCADA gateway, MES, historian, custom
  industrial-IoT bridge, OPC-UA-server adapter, Calcite/Camel/NiFi
  integration). Read tags, write tags, subscribe to changes, browse
  device address space *(documented: `README.md` — "The Industrial IoT
  adapter"; `website/asciidoc/modules/users/pages/protocols/index.adoc`)*.
- **Single-PLC, request/response style traffic** initiated by the embedding
  application against a known device endpoint
  *(documented: per-protocol pages; `plc4j/api/src/main/java/.../PlcConnection.java`)*.

### Deployment shape

PLC4X is **not** a network service. It is an **in-process library** built
into a host application (Java, Go, C, or Python process). The host
application — not PLC4X — is the thing exposed to its own users and to
the wider IT network. PLC4X opens **outbound** connections to PLC
endpoints; nothing in the main repo listens *(documented: `README.md`,
the `plc4j/transports/` inventory which is exclusively client-side —
`tcp`, `udp`, `serial`, `raw-socket`, `pcap-replay`, `socketcan`, `can`,
`virtualcan`; no listener transport)*.

The library lives where the **OT network meets the IT network**: in
typical deployments, the host process runs on a Linux/Windows VM or
container with network reachability into the plant-floor / fieldbus
segment, and surfaces data northbound by some other mechanism the host
application chose. PLC4X has no opinion on the northbound side.

### Caller roles

Because PLC4X is a library, there is no "client / operator / peer" role
split of the kind a network service would have. The roles are:

| Role                                                                    | Trust level                              | Notes                                                                                                                                                                                                                                                                                                         |
|-------------------------------------------------------------------------|------------------------------------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| **Embedding application**                                               | trusted — sole API caller                | Decides which driver to load, supplies the connection URL, the `PlcAuthentication` (when the underlying protocol supports it), and every tag-address string. PLC4X has no notion of an "end user" within the embedding process *(documented: `plc4j/api/.../PlcConnection.java`, `PlcDriverManager.java`)*.   |
| **Remote PLC device**                                                   | **untrusted on the wire**                | Industrial controllers in the field; the device's firmware, the wire bytes it emits, and any responses to PLC4X requests are all treated as untrusted input crossing a trust boundary into the library's parser *(inferred — §14 Q1)*.                                                                        |
| **OPC UA server**                                                       | partially-authenticated peer             | When OPC UA is the protocol, the server identity is bound to a certificate that PLC4X can verify against a trust store; that is the only protocol in the supported set where peer identity is cryptographic *(documented: `website/asciidoc/.../protocols/opcua.adoc`; `plc4j/drivers/opcua/.../security/`)*. |
| **PLC4X server / OPC UA server (the embedding side of `plc4x-extras`)** | **out of scope** — see §3 *(maintainer)* |                                                                                                                                                                                                                                                                                                               |

### Component-family table

The main repo carves naturally into a dozen families with distinct
threat profiles. Anything marked **out** below reappears in §3 with the
reason.

| Family                                                                                                                                                                             | Representative entry point                                               | Touches outside the process?                    | In this model?                                                                                                                        |
|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|--------------------------------------------------------------------------|-------------------------------------------------|---------------------------------------------------------------------------------------------------------------------------------------|
| **Java public API (`plc4j/api`)**                                                                                                                                                  | `PlcDriverManager.getConnectionFactory().getConnection(url, auth)`       | no — pure dispatch                              | **yes**                                                                                                                               |
| **Java SPI / runtime (`plc4j/spi`)**                                                                                                                                               | `AbstractPlcConnection`, `ConversationContext`, `Plc4xProtocolBase`      | no — orchestration                              | **yes**                                                                                                                               |
| **Per-protocol drivers (`plc4j/drivers/{ab-eth,ads,bacnet,c-bus,can,canopen,ctrlx,eip,firmata,iec-60870,knxnetip,modbus,opcua,open-protocol,plc4x,profinet,profinet-ng,s7,simulated,slmp,umas}`)** | one `PlcDriver` per protocol                                             | **yes — network**, via the registered transport | **yes — but only for the drivers the protocol support table marks with a green check; see §3 item 11** |
| **Transports (`plc4j/transports/{api,tcp,udp,serial,cotp,tls,raw-socket,pcap-replay,can,can-socketcan,can-virtualcan,test}`)** | blocking `SocketChannel` on a virtual thread per connection (TCP), `SSLSocket` (TLS), serial, libpcap | **yes — sockets, serial devices, libpcap**      | **yes**                                                                                                                               |
| **OPC UA security subsystem** (`plc4j/drivers/opcua/.../security/`, `.../context/SecureChannel`)                                                                                   | certificate verifier, security-policy negotiator, secure-channel encoder | **yes — TLS-like primitives on the wire**       | **yes — distinguished from other drivers; see §8 P1**                                                                                 |
| **Code-generation runtime (mspec read/write generators)** under `plc4j/.../readwrite/`, `plc4go/protocols/`, `plc4c/generated-sources/`                                            | generated parsers / serializers                                          | **yes — driven by remote-PLC bytes**            | **yes** — same trust level as the driver that uses it                                                                                 |
| **Go bindings (`plc4go/`)**                                                                                                                                                        | `plc4go.NewPlcDriverManager()`                                           | as core                                         | **yes**                                                                                                                               |
| **C bindings (`plc4c/`)**                                                                                                                                                          | `plc4c_*` API                                                            | as core                                         | **yes — but README declares "not ready for usage"** *(documented: `README.md`)* — see §3                                              |
| **Python bindings (`plc4py/`)**                                                                                                                                                    | `plc4py` import                                                          | as core                                         | **yes — but README declares "not ready for usage"** *(documented: `README.md`)* — see §3                                              |
| **.NET bindings (`plc4net/`)**                                                                                                                                                     | n/a                                                                      | n/a                                             | **out** — README declares "not ready for usage - abandoned" *(documented: `README.md`)*                                               |
| **PLC4X proxy driver (`plc4j/drivers/plc4x`)**                                                                                                                                     | `plc4x:tcp://host/?remote-connection-string=...`                         | network — outbound to a remote PLC4X server     | **yes for the client side; the server side is in `plc4x-extras` and out of scope** *(documented: `website/.../protocols/plc4x.adoc`)* |
| **Simulated / mock drivers (`plc4j/drivers/{simulated,mock}`)**                                                                                                                    | `simulated://...`                                                        | no                                              | **yes for code-quality; not a security surface**                                                                                      |
| **Utilities (`plc4j/utils/{raw-sockets,test-generator,test-utils}`)** | dev tooling                                                              | varies                                          | **out** — dev / test only *(inferred — §14 Q2)*                                                                                       |
| **Client-side tools (`plc4j/tools/{connection-cache,opm,event-pump,capture-replay}`)** | `PlcConnectionCache`, `PlcEntityManager`, the event-pump's YAML job configuration | indirectly — they drive connections through the drivers | **yes** — they hold live connections, credentials and connection strings on behalf of the embedding application. Published from this repo (`plc4j-tools-*`), **not** part of the `plc4x-extras` carve-out. |
| **Audit log (`plc4j/utils/audit-log/{api,impl}`)** | `audit-log-file` connection-string parameter | **yes — writes wire traces to the host filesystem** | **yes** — off unless a path is configured; when on it writes request/response traces (payload bytes, and connection metadata subject to the §5a redaction rules) to a rolling file the integrator must protect. |
| **Subscription emulation (`plc4j/utils/subscription-emulation`)** | polling loop on top of a driver's read path | via the driver it wraps | **yes** — same trust level as the driver it polls. |
| **Code-generation tooling under `code-generation/`**                                                                                                                               | mspec compiler invoked at build                                          | host filesystem                                 | **out** — build-time, not in the deployed artifact *(inferred — §14 Q2)*                                                              |
| **`tools/`, `images/`, `media/`, `website/`, `licenses/`, `src/`, `.idea/`, `.github/`, `.mvn/`, `Dockerfile`, `docker-compose.yaml`, `Jenkinsfile`**                              | tooling, docs, build orchestration                                       | n/a                                             | **out** *(inferred — §14 Q2)*                                                                                                         |

A finding is in-model only if it lands in a row marked "yes" — see §4
for the per-component reachability test.

## §3 Out of scope (explicit non-goals)

PLC4X is **not**, and does not aim to be, the following — reports requiring
any of these will be closed with the cited disposition.

1. **A defender for the OT protocols themselves.** Modbus, S7 (PUT/GET),
   BACnet/IP, IEC-60870-5-104, Profinet, CANopen, KNXnet/IP, C-Bus, DF1,
   AB-Ethernet, ADS/AMS (unless the operator has configured the AMS-route
   credential flow), EtherNet/IP, and UMAS as standardized **do not
   authenticate the peer, do not protect message integrity, and do not
   encrypt** *(documented per protocol: `website/.../protocols/*.adoc`;
   inferred from the protocol specifications themselves — §14 Q3)*. A
   report of the shape "the wire payload is unauthenticated cleartext" or
   "anyone on the OT network can spoof a Modbus master" describes the
   **protocol**, not the library; PLC4X cannot add what the protocol
   does not have. → `BY-DESIGN: protocol-disclaimed` (§9). The integrator
   is responsible for the OT-network perimeter (§10). See §11a for the
   per-protocol enumeration.
2. **A server / responder for any of the supported protocols.** PLC4X is a
   **client-side** library: it speaks the protocols *to* PLCs, not as a
   PLC *(documented: `README.md`, `plc4j/api/.../PlcConnection.java` —
   single connect/read/write/subscribe shape; transport inventory is
   exclusively client)*. Server / gateway functionality (OPC-UA Server,
   PLC4X Server, Calcite, Camel, Kafka-Connect, Karaf, NiFi adapters)
   is published in the separate **`plc4x-extras`** repository
   *(documented: `README.md`)* and is **explicitly out of scope of this
   document** *(maintainer: PMC chair, 2026-05-29 16:06Z)*. →
   `OUT-OF-MODEL: unsupported-component`. This carve-out does **not**
   cover the client-side tooling that lives here in `plc4j/tools/` —
   connection cache, OPM, event-pump, capture-replay — which is in
   scope (§2).
3. **A defender against the embedding application itself.** The embedding
   application supplies every connection URL, the `PlcAuthentication`,
   every tag-address string, and every value written. Anything the
   embedding application could express via the public API is by
   definition authorized; "a malicious embedding application writes to
   the wrong coil" is not a PLC4X vulnerability *(inferred — §14 Q4)*.
   → `OUT-OF-MODEL: trusted-input`.
4. **A defender for end users of the embedding application.** PLC4X has
   no concept of an "end user". Whether the embedding application
   authenticates its own users, what authorization it enforces over
   which tags they may read or write, and how it sanitizes user-supplied
   tag-address strings before passing them to `parseTagAddress` is the
   embedding application's responsibility *(inferred — §14 Q4)*. →
   `OUT-OF-MODEL: trusted-input`.
5. **A sandbox for tag-address strings sourced from untrusted users.**
   Tag-address syntax is protocol-specific and rich (S7
   `%DB1.DBX0.0:BOOL`, OPC UA `ns=2;s=...`, Modbus `holding-register:40001:INT`,
   …). If the embedding application accepts user-supplied tag strings,
   it is responsible for vetting them — the parser does not provide a
   security boundary against attacker-supplied tag strings (it provides a
   correctness boundary against malformed strings) *(inferred —
   §14 Q5)*. → `OUT-OF-MODEL: trusted-input` for syntactic-validation
   reports, `VALID` for parser memory-corruption / DoS findings reachable
   from caller-supplied strings.
6. **Code that ships in the repository but is not part of the supported
   product.** `plc4c/` and `plc4py/` are flagged "not ready for usage" in
   `README.md`; `plc4net/` is "not ready for usage - abandoned"; the
   `tools/`, `code-generation/`, `plc4j/utils/`, `plc4j/drivers/{simulated,mock}`,
   `images/`, `media/`, `website/`, `licenses/`, `src/`, top-level
   `Dockerfile` / `docker-compose.yaml` / `Jenkinsfile` are build,
   test, and documentation infrastructure *(documented: `README.md`;
   inferred for the rest — §14 Q2)*. → `OUT-OF-MODEL: unsupported-component`.
7. **The PLC device firmware, the network the device is on, and the
   physical process the device controls.** PLC4X is a wire-protocol
   library; bugs / exploits in the PLC itself (S7 vendor-specific
   parser flaws, BACnet stack vulnerabilities on the device side, OPC UA
   server bugs) are upstream to the PLC vendor. The physical
   consequences of a write — that a coil turning on opens a relief
   valve, that a setpoint change spoils a batch — are downstream to
   the integrator's process design *(inferred — §14 Q6)*. →
   `OUT-OF-MODEL: adversary-not-in-scope` for device bugs;
   `OUT-OF-MODEL: trusted-input` for process-control consequences of an
   API-authorized write.
8. **OT-network segmentation, hardening, and intrusion detection.**
   PLC4X assumes the embedding application is deployed *somewhere with
   network reachability to the PLC*; whether that path is firewalled,
   air-gapped, monitored by an IDS, or routed over a hostile carrier is
   the integrator's deployment posture *(inferred — §14 Q7)*. →
   `OUT-OF-MODEL: adversary-not-in-scope`.
9. **TLS / VPN / IPSec tunnels that the integrator may layer underneath.**
   Most cleartext OT protocols are operated in production behind some
   form of point-to-point tunnel (a VPN concentrator at the OT/IT
   boundary, an IPSec link to a remote pumping station, a TLS-tunneling
   proxy in front of a Modbus device). PLC4X has **no built-in TLS or
   tunneling layer** for the cleartext protocols *(inferred from
   transport inventory — §14 Q8)*. Reports of the shape "Modbus-TCP
   should run over TLS" describe an integrator deployment choice; the
   driver does not block such a deployment but does not provide it
   either. → `BY-DESIGN: property-disclaimed`.
10. **Supply-chain / build / release hygiene** — Maven action pinning,
    Maven Central signing, Jenkins build configuration, dependency
    freshness, reproducible builds. Out of model per the SKILL.

11. **Drivers not marked as fully supported in the protocol support
    table.** The per-language matrix in
    `website/asciidoc/modules/users/pages/protocols/index.adoc` is the
    authoritative statement of which drivers ship as supported: a green
    check means the driver is complete and supported for that language
    binding, a yellow exclamation means partial or experimental, and a
    red cross means it does not exist there. **Only a driver / binding
    combination carrying a green check is in this model**
    *(maintainer: Chris Dutz, 2026-09-03)*. A driver the table marks
    yellow is incomplete by declaration rather than by oversight, so a
    report that it mishandles a malformed frame, leaks a socket, or
    omits a check its complete siblings make describes unfinished work,
    not a vulnerability in a shipped product. In PLC4J such a driver
    additionally carries `maven.deploy.skip` in its POM *(documented:
    `plc4j/drivers/*/pom.xml`)*, so it is not published to Maven
    Central and no release artifact contains it. →
    `OUT-OF-MODEL: unsupported-component`.

    As the table stood on 2026-09-03, this excludes for **PLC4J**:
    BACnet/IP, CAN (raw), C-Bus, ctrlX, DF1, Open-Protocol and
    Profinet; and for **PLC4Go**: AB-Ethernet, CAN (raw), CANopen,
    ctrlX, DF1, Firmata, **OPC-UA**, Open-Protocol, PLC4X
    (proxy-protocol) and Profinet. Note the consequence for the one
    security-bearing protocol in the set: the **Go** OPC UA driver is
    marked partial and is therefore outside this model, while the Java
    one is green and stays fully inside it, including every §8 property
    that names OPC UA. `plc4c/` and `plc4py/` remain out under item 6
    whatever their cells say. The table governs, not the lists above: a
    driver enters or leaves the model when its cell changes (§12).

12. **Wrong, corrupt, or hostile *configuration files* the operator
    points the driver at.** Several drivers read a file the operator
    names in the connection string: a `.knxproj` ETS project
    (`knxproj-file-path`, with `knxproj-password`), a BACnet EDE file or
    directory (`ede-file-path`, `ede-directory-path`), a Profinet GSDML
    directory (`gsd-directory`), a key store or trust store
    (`tls.keystore`, `tls.trust-store`), a pinned certificate
    (`server-certificate-file`), a capture for the replay transport
    (`pcap-file`), and the event-pump's job configuration. **These
    arrive across B1, from the trusted caller — not across B2 from the
    wire** *(§4, §6)*. The operator chose the file the same way they
    chose the connection string, and anything they can express by
    supplying a bad file they could equally express by supplying a bad
    URL. A report that a truncated `.knxproj` throws, that a key store
    in the wrong format fails the connection, that a garbled EDE column
    yields a wrong tag, or that a GSDML which does not match the device
    produces nonsense, describes a **correctness and error-message
    matter**, not a vulnerability *(maintainer: Chris Dutz,
    2026-09-03)*. → `OUT-OF-MODEL: trusted-input`; fix-worthy as
    ordinary robustness, triaged as a normal bug.

    Two things are **not** covered by this item, because the project
    implemented them deliberately against a file it does not trust:
    §8 P6 — XXE, external-DTD and external-schema resolution are
    disabled in the ETS parser, so a `.knxproj` that reads a file off
    the host or reaches a network resource is **`VALID`** — and the
    bound on how far one `.knxproj` entry may expand, so a small file
    that inflates without limit is **`VALID`** too. The line between
    them: a bad file that fails to parse is out of model; a bad file
    that makes the parser do something *beyond parsing* is in it.
    Credentials read from any of these files are also still subject to
    the §5a redaction rules.

13. **Whether the endpoint the operator chose deserves to be trusted.**
    PLC4X connects where it is told. The connection URL is trusted
    input (item 3, §6 B1), and the library has no way to know whether
    the thing answering on the other end is the device the operator
    meant. A report of the shape "I stood up a server that is not what
    it claims to be, pointed the driver at it, and the driver believed
    it" — wrong process values, a fabricated browse tree, a controller
    advertising a type or capability it does not have, a device that
    answers every request with success — describes the **operator's
    choice of endpoint**, not a defect in the library *(maintainer:
    Chris Dutz, 2026-09-03)*. Nor does PLC4X judge whether an endpoint
    deserves the credentials configured for it: forwarding the
    operator's username and password to the address the operator named
    is the function, not a leak. → `OUT-OF-MODEL: trusted-input` for the
    semantic content of what the endpoint says;
    `OUT-OF-MODEL: adversary-not-in-scope` where the report amounts to
    "the integrator should not have connected there" (§3 items 7, 8).

    This item is about **what the peer says**, never about **what its
    bytes do to the driver.** Everything the parser must survive stays
    in model and unchanged: §8 P4 — memory safety, bounded allocation,
    no unbounded recursion, no hang on adversarial response bytes — is
    `VALID` whoever sent those bytes, because a hostile endpoint and an
    on-path attacker on a network the operator *did* trust reach exactly
    the same code, and §7 names the latter as this model's primary
    adversary. So are the OPC UA channel properties (§8 P1, P2, P7),
    which exist precisely so that the endpoint does **not** have to be
    trusted in advance: "the driver accepted a server certificate it
    should have rejected" remains `VALID`.

## §4 Trust boundaries and data flow

PLC4X has **two primary trust boundaries** plus one optional
cryptographic one for OPC UA. A finding is in-model only when it cleanly
maps to one of them.

| #                   | Transition                                                                                                                                                                                          | Authentication                                                                                                                                                                                                                                                                                                       | Authorization                                                 | Notes                                                                                                                                                                                                 |
|---------------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|---------------------------------------------------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| B1                  | **Embedding application → PLC4X API surface**                                                                                                                                                       | trusted by construction                                                                                                                                                                                                                                                                                              | the API caller is trusted                                     | URL string, `PlcAuthentication`, tag-address strings, write values *(documented)*                                                                                                                     |
| B2                  | **PLC4X driver → remote PLC over the wire** (Modbus, S7 PUT/GET, BACnet/IP, IEC-60870-5-104, Profinet, CANopen, KNXnet/IP, C-Bus, DF1, AB-Ethernet, EtherNet/IP, UMAS, ADS/AMS without credentials) | **none — by protocol design** *(documented per protocol)*                                                                                                                                                                                                                                                            | none on the wire; whatever the embedding application enforces | The OT-network perimeter is the only security control. Wire bytes returned to the driver are **untrusted input crossing into the parser** *(inferred — §14 Q3)*.                                      |
| B2-OPCUA            | **PLC4X OPC UA driver → OPC UA server (encrypted policy)**                                                                                                                                          | server cert verified against trust store *(when `tls.trust-store` is set)*; client cert optional; username/password or token over the encrypted channel                                                                                                                                                              | server-side ACLs                                              | Mutually authenticated (cert + cert) and encrypted at the higher security policies; **see §8 P1, §10 item 4** *(documented: `website/.../protocols/opcua.adoc`, `plc4j/drivers/opcua/.../security/`)* |
| B2-OPCUA-PERMISSIVE | **OPC UA driver with `tls.verify=false`** | server certificates **are not validated** — the driver installs `PermissiveCertificateVerifier` and logs a warning saying so *(documented: `plc4j/drivers/opcua/.../context/OpcuaDriverContext.java` `buildCertificateVerifier`)* | n/a                                                           | **No longer the default.** With `tls.verify` at its default of `true` and no trust anchor configured, the driver installs `RejectingCertificateVerifier` and the connection fails closed. See **§5a**. |
| B2-ADS              | **ADS driver → Beckhoff TwinCAT (AMS route setup with credentials)**                                                                                                                                | username/password to the TwinCAT system for AMS-route setup, when a `PlcUsernamePasswordAuthentication` is supplied *(documented: `plc4j/drivers/ads/.../AdsProtocolLogic.java` lines 130–145)*                                                                                                                      | TwinCAT-side                                                  | The credentials are forwarded to the device; what the device does with them is outside PLC4X. ADS payload data itself is unauthenticated cleartext *(inferred — §14 Q9)*.                             |
| B3                  | **PLC4X transport → host OS / NIC / serial port / libpcap**                                                                                                                                         | OS-level                                                                                                                                                                                                                                                                                                             | OS-level                                                      | Whatever the embedding process's UID, capabilities, and seccomp / AppArmor profile permit. The `raw-socket` transport needs `CAP_NET_RAW` on Linux *(inferred — §14 Q10)*.                            |
| B4                  | **Code-generation runtime → operating system** at build time                                                                                                                                        | n/a — only the developer runs this                                                                                                                                                                                                                                                                                   | n/a                                                           | Build-time only; not part of the deployed artifact (out per §3 item 6)                                                                                                                                |

### Reachability preconditions per family

For each family in §2, a finding is in-model only if it is reachable as
follows.

- **Per-protocol driver, generated `readwrite/` parser code, OPC UA
  security subsystem**: reachable from **bytes the remote PLC (or
  attacker-on-the-OT-network) returns to the driver over the
  configured transport**. The parser/decoder must not corrupt host memory,
  must not infinite-loop, must not allocate unboundedly, and must not
  cause the host process to crash on adversarial input
  *(inferred — §14 Q11)*. **This is the single most important reachability
  boundary in the document** — most realistic PLC4X-side vulnerabilities
  will be wire-parser bugs.
- **Public API (`plc4j/api`), SPI (`plc4j/spi`)**: reachable from the
  embedding application's calls. Memory-safety / correctness issues
  here are valid only if the embedding application's call sequence is
  documented (per the Javadoc / the public Java API).
- **Tag-address parser** (per-driver `*TagHandler`, `parseTagAddress`):
  reachable from the **tag-address string** passed by the embedding
  application. PLC4X does not assume the string itself is
  attacker-controlled — see §3 item 5 — but a parser bug that
  corrupts memory or hangs the parser is still `VALID-HARDENING` (and
  `VALID` if the parser is memory-corrupting and the embedding
  application *does* pass user-controlled strings) *(inferred —
  §14 Q5)*.
- **Transport layer**: reachable from bytes the OS delivers via the
  socket / serial port / libpcap capture file. Memory safety on
  malformed framing is in model *(inferred — §14 Q11)*.
- **OPC UA security subsystem** (B2-OPCUA): reachable from
  certificate chains, security tokens, signed handshake messages
  produced by the remote endpoint. **Findings here are weighted more
  heavily than equivalent findings in a cleartext-protocol driver**:
  the OPC UA driver is the one place in the repository that
  *promises* a cryptographic property, so a break of that property
  is a §8 P1 violation, not a §11a non-finding.
- **Generated code (mspec)**: identical to the driver that consumes it
  — generated parsers run inside the driver. The mspec compiler
  itself is build-time and out of model per §3 item 6.
- **`plc4go/`, `plc4c/`, `plc4py/`**: per-language; the C and Python
  bindings are flagged "not ready for usage" in `README.md` — see
  §3 item 6 and §5a.

## §5 Assumptions about the environment

- **Operating system**: Linux primary (Ubuntu / RHEL family),
  Windows / macOS supported for Java; ARM-and-x86 supported.
  Java 21 minimum, tested up to Java 24 *(documented: `README.md`)*.
- **Java runtime**: HotSpot or OpenJ9; the build requires
  `--add-exports jdk.compiler/...` for code generation but the
  runtime artifact does not *(documented: `README.md`)*.
- **Process model**: in-process library. PLC4X loads as a JAR (or Go
  module, or Python wheel) into the embedding application; it
  registers protocol drivers via the Java SPI mechanism (or
  language-equivalent). No long-lived daemons of its own *(documented:
  `plc4j/api/.../PlcDriverManager.java`)*.
- **Network**: outbound TCP / UDP / raw-socket / SocketCAN / serial as
  required by the selected driver. Ports per the protocol:
  502/Modbus-TCP, 102/S7, 47808/BACnet, 44818+2222/EtherNet-IP, 4840
  + dynamic/OPC UA, 2404/IEC-60870-104, 3671/KNXnet/IP, etc.
  *(documented per protocol page)*. **The OT network is assumed
  to be operated as a defended perimeter by the integrator** —
  see §10 item 1.
- **Memory**: the JVM (or Go GC, or C `malloc`, or CPython allocator)
  provides the underlying memory model. PLC4X writes Java / Go / C /
  Python code at the language's normal trust level.
- **Concurrency**: **Netty is gone.** Since the SPI3 migration the Java
  transports are plain JDK sockets — TCP is one virtual thread per
  connection doing blocking `SocketChannel` I/O (no NIO selector, no
  readiness-event loop), TLS is an `SSLSocket`, and no `io.netty`
  artifact remains on the runtime classpath *(documented:
  `plc4j/transports/tcp/.../TcpTransportInstance.java`; no Netty
  dependency in any `plc4j` POM as of 2026-09-03)*. Calls through
  `read`/`write`/`subscribe`/`browse` still return
  `CompletableFuture`s, and thread safety is per the documented Java
  API contract *(inferred — §14 Q12)*. A request waiting on the
  per-connection permit is queued rather than parking the submitting
  thread, which used to deadlock a driver that chained requests on its
  own receive thread.
- **Cryptography (OPC UA only)**: the JCE provider shipped with the
  JVM (OpenJDK / Oracle / IBM) provides AES, RSA, RSASSA-PSS,
  HMAC-SHA1/256, X.509 parsing
  *(documented: `plc4j/drivers/opcua/.../security/SecurityPolicy.java`)*.
- **Filesystem (OPC UA only)**: a keystore and truststore are loaded
  from operator-configured paths
  *(documented: `plc4j/drivers/opcua/.../config/OpcuaConfiguration.java`,
  `plc4j/drivers/opcua/.../context/OpcuaDriverContext.java`)*.
- **Filesystem (KNX only)**: `.knxproj` ETS files are parsed via
  XML; XXE / external-DTD / external-schema are explicitly
  disabled *(documented: `plc4j/drivers/knxnetip/.../ets/EtsParser.java`
  lines 168–171 — empty `ACCESS_EXTERNAL_DTD/SCHEMA`,
  `disallow-doctype-decl` and `FEATURE_SECURE_PROCESSING`; the line
  numbers moved from 65–71 as the file grew)*. Since 2026-08-24 the
  parser also **bounds how far one zip entry may expand**, since a
  `.knxproj` entry declares its inflated size and can declare anything,
  and cleans up what it unpacked.
- **Filesystem (BACnet/IP only)**: EDE (Engineering Data Exchange)
  files describing the target devices are loaded from
  operator-configured paths *(documented:
  `plc4j/drivers/bacnet/.../configuration/BacNetIpConfiguration.java`)*.
- **Filesystem (audit log, any driver)**: setting the `audit-log-file`
  connection-string parameter makes the driver **write** request and
  response traces to a rolling file on the host *(documented:
  `plc4j/utils/audit-log/api/.../config/AuditLogConfiguration.java`,
  `plc4j/utils/audit-log/README.md`)*. It is off unless a path is set,
  and the implementation module has to be on the classpath. When on,
  the file holds payload bytes and connection metadata: the parameter
  redaction of §5a applies to what is rendered, but the traces are
  plant data at rest and the integrator owns the file's permissions
  and lifetime (§10).
- **What PLC4X does NOT do to its host** (negative claims, awaiting
  maintainer ratification — these are predominantly *(inferred)*,
  see §14 Q13):
  - **does not** install signal handlers besides those Netty / JVM
    install for it *(inferred — §14 Q13)*;
  - **does not** spawn child processes at runtime *(inferred —
    §14 Q13)*;
  - **does not** open listening sockets — all sockets are outbound
    *(inferred from transport inventory — §14 Q13)*;
  - reads **one** environment variable that changes a security-relevant
    bound: `PLC4X_MAX_NESTING_DEPTH` raises or lowers the 1024-level
    limit on how deeply a parsed message may nest its types, and means
    the same thing in all four bindings *(documented:
    `plc4j/spi/.../AbstractBuffer.java`, `plc4go/spi/utils/ReadBufferByteBased.go`,
    `plc4py/.../ReadBuffer.py`, `plc4c/spi/include/plc4c/spi/context.h`)*. A
    value that is not a positive number leaves the default in place with a
    warning. Beyond it, PLC4X **does not** read environment variables for
    security-sensitive behavior outside standard JVM/JNI conventions
    *(inferred — §14 Q13)*;
  - **does not** mutate process-wide global state (allocator, locale,
    `System.setSecurityManager`, etc.) at runtime *(inferred —
    §14 Q13)*;
  - the `raw-socket` transport **does** require `CAP_NET_RAW` (Linux)
    / Administrator (Windows) / root-equivalent privileges when
    used; this is a privileged transport with an inherent host-level
    cost the embedding application takes on by choosing it
    *(inferred — §14 Q10)*.

## §5a Build-time and configuration variants

PLC4X exposes **per-driver configuration knobs** rather than a
single global flag set. The security-relevant ones — those whose
default value materially changes the security envelope — are
collected here.

| Knob                                                                                   | Default                                                                                                                                                                                                                                                                   | Maintainer stance                                                                                                                                                                                                          | Effect                                                                                                                                                                                                                                                                                            |
|----------------------------------------------------------------------------------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| OPC UA `security-policy`                                                               | `Basic256Sha256` *(documented: `plc4j/drivers/opcua/.../config/OpcuaConfiguration.java`; `website/.../protocols/opcua.adoc`)* | **hardened default** — `NONE` is now an explicit opt-in *(maintainer ruling, §14 Q14; shipped 2026-08-24)* | `NONE` means the channel is neither signed nor encrypted and the server is not authenticated at all, collapsing B2-OPCUA to B2. Asking for it is a deliberate act; nothing downgrades into it (see the `discovery` row). |
| OPC UA `message-security`                                                              | `SIGN_ENCRYPT` *(documented)*                                                                                                                                                                                                                                             | hardened default                                                                                                                                                                                                           | When the security policy is not `NONE`, this forces sign-and-encrypt; flipping to `SIGN` (signed-cleartext) or `NONE` weakens the channel.                                                                                                                                                        |
| OPC UA `tls.trust-store`                                                               | **unset** *(documented: `OpcuaConfiguration.java`, `OpcuaDriverContext.buildCertificateVerifier`)* | **fail-closed** — with `tls.verify` at its default, a connection with no trust anchor is refused, not accepted *(maintainer ruling, §14 Q15; shipped 2026-08-24)* | The verifier is chosen in this order: `tls.verify=false` → `PermissiveCertificateVerifier` (warned about); `tls.trust-store` set → `TrustStoreCertificateVerifier` (chain validation); `server-certificate-file` set → `PinnedCertificateVerifier` (exact certificate); otherwise → `RejectingCertificateVerifier`, so the handshake fails with no trust anchor rather than trusting whatever is presented. |
| OPC UA `tls.verify` | `true` *(documented: `OpcuaConfiguration.java`)* | **hardened default** — verification on, opt-out only | Turning it off installs the permissive verifier and logs a warning naming the risk. This is the knob that replaced `insecure-certificate-verification`, with the opposite sense: a connection that used to set `insecure-certificate-verification=true` must now set `tls.verify=false`, and one that misses the rename gets verification. |
| OPC UA `tls.keystore` / `tls.keystore-password` / `tls.keystore-type` | unset; `tls.keystore-type` `pkcs12`; `generated-key-size` 2048 *(documented)* | operator must supply for a stable client identity | If unset, the driver generates a self-signed application-instance certificate: since 2026-08-14 signed with SHA-256, carrying a positive 16-octet serial, marked as an end entity with the OPC UA Part 6 key usages and a subject key identifier — a certificate a modern server will accept, where the previous SHA-1 / CA-marked one was widely refused. It is still regenerated on every start, so it cannot satisfy a peer that requires a known client identity. The client identity taken from a key store is the first entry holding a private key, and a store with none, an unreadable key, or a non-RSA key fails naming the entry. |
| OPC UA `discovery`                                                                     | `true` *(documented)* | enabled by default; **the discovery phase is still conducted with security policy `NONE`**, as the OPC UA spec has it *(documented: `OpcuaConfiguration.java`)* | An attacker on the path to the discovery endpoint still sees and can rewrite the advertised endpoints. What changed on 2026-08-24 is what the driver does next: it refuses to continue onto a channel weaker than the configured `security-policy` (it used to silently do exactly that), and a certificate learned during discovery is never used as its own trust anchor — it is checked against the configured verifier like any other. So a connection asking for `Basic256Sha256` now fails where it previously came up unprotected. |
| OPC UA `username` / `password`                                                         | unset *(documented)*; `allow-insecure-credentials` `false` | operator-supplied; **the driver refuses to put a password on an unprotected channel** | Forwarded as the OPC UA `UserIdentityToken`, encrypted with the algorithm the server's user token policy names (it used to hard-code RSA-OAEP, which failed against servers asking for RSA-PKCS1.5). Over a channel that neither signs nor encrypts, the connection now fails instead of sending the credential in clear; `allow-insecure-credentials=true` sends it anyway, with a warning. |
| OPC UA `channel-lifetime-ms`, `min-channel-lifetime-ms`, `session-timeout-ms`, `handshake-timeout-ms`, `request-timeout-ms` | 1 h / 5 s / 2 min / 60 s / 30 s *(documented)* | reasonable defaults | DoS / timeout-tuning surface; not a security boundary. `min-channel-lifetime-ms` bounds a lifetime the server revises downwards — channel renewals share one executor across every OPC UA connection in the JVM, so a server answering with a very short (or zero) lifetime used to set the renewal pace for all of them. Note the names: every duration ends in `-ms` since the 2026-08-28 configuration-vocabulary change, and `negotiation-timeout` is now `handshake-timeout-ms`. |
| OPC UA `browse-max-references-per-node`, `browse-max-total-nodes`, `browse-max-depth` | 65536 / 1000000 / 64 *(documented)* | bounded by default; `0` disables a bound | A browse walks whatever tree the server describes, for as long as it describes one. These bound what a hostile or broken server can make one browse cost in memory and time; reaching a bound warns and returns what was found rather than failing. Setting any of them to `0` restores the unbounded behaviour. |
| TLS transport `tls.verify`, `tls.ignore-common-name`, `tls.trust-store` | `true` / `false` / unset *(documented: `plc4j/transports/tls/.../TlsTransportConfiguration.java`)* | hardened defaults since 2026-08-24 | The transport now checks that the server certificate was issued for the host it connected to; `ignore-common-name` was declared but never consulted, so with `verify` on, any certificate from a trusted issuer was accepted for any host. `tls.trust-store` names the certificates to trust instead of the JVM's public authorities — previously the only way past a private CA was to turn verification off. `ignore-common-name=true` restores the old behaviour and logs a warning. |
| ADS `PlcUsernamePasswordAuthentication` (when supplied)                                | unset (no AMS-route setup)                                                                                                                                                                                                                                                | operator-supplied                                                                                                                                                                                                          | When supplied, drives an AMS-route registration against the TwinCAT system using HTTP-style credentials *(documented: `plc4j/drivers/ads/.../AdsProtocolLogic.java`)*. The credentials are not protected by PLC4X on the wire; TwinCAT-side TLS is the device's responsibility.                   |
| BACnet/IP `ede-file-path` / `ede-directory-path`                                       | unset                                                                                                                                                                                                                                                                     | operator-supplied                                                                                                                                                                                                          | If set, points at filesystem paths; standard file-permission rules apply. The BACnet/IP driver is itself out of model per §3 item 11; the knob is listed for completeness. |
| KNX `.knxproj` parser                                                                  | XXE / external-DTD / external-schema **disabled** *(documented: `EtsParser.java`)*                                                                                                                                                                                        | hardened — this is the safe defaults case                                                                                                                                                                                  | Reports of the shape "XXE in `.knxproj` parsing" are `KNOWN-NON-FINDING`.                                                                                                                                                                                                                         |
| S7 `controller-type`                                                                   | unset (auto-discover via SZL)                                                                                                                                                                                                                                             | operator-supplied for Siemens LOGO compatibility                                                                                                                                                                           | Functional knob; not a security boundary.                                                                                                                                                                                                                                                         |
| Modbus connection options (`unit-id`, byte order)                                      | per spec                                                                                                                                                                                                                                                                  | functional knobs                                                                                                                                                                                                           | not security boundaries.                                                                                                                                                                                                                                                                          |
| `plc4c/`, `plc4py/`, `plc4net/` builds                                                 | **README declares these "not ready for usage" (with `plc4net` "abandoned")** *(documented: `README.md`)*                                                                                                                                                                  | **OUT-OF-MODEL** for §8 properties; bugs reported against C / Python / .NET bindings should be triaged as code-quality / completeness, not as supported-product vulnerabilities until the maintainer reclassifies them     | A finding in `plc4c/` is `OUT-OF-MODEL: unsupported-component` until the README line changes.                                                                                                                                                                                                     |

### The OPC UA defaults (resolved 2026-08-24 — secure by default)

Earlier revisions of this document described the OPC UA driver's
defaults as **dev/lab-grade**: `security-policy=NONE`, no certificate
verification, and a discovery phase the driver would silently continue
from onto an unprotected channel. That was accurate when it was
written; it is no longer what the code does. The §14 Q14 / Q15 rulings
were resolved in favour of secure-by-default and implemented:

- `security-policy` defaults to **`Basic256Sha256`** with
  `message-security=SIGN_ENCRYPT`, so a channel is signed and encrypted
  unless `NONE` is asked for explicitly;
- `tls.verify` defaults to **`true`**, and with no trust anchor
  configured the driver installs `RejectingCertificateVerifier` — the
  connection **fails closed** rather than trusting whatever certificate
  arrives. `PermissiveCertificateVerifier` is reached only via
  `tls.verify=false`, which logs a warning naming the risk;
- the driver **refuses to settle for a channel weaker than the one
  configured**, so the unauthenticated discovery phase can no longer
  hand back an unprotected session to a connection that asked for a
  protected one, and a certificate learned during discovery is never
  its own trust anchor;
- an endpoint has to match **both** the requested security policy and
  the requested message security mode, and the **strongest** matching
  endpoint is selected rather than the weakest — a server publishing a
  wide-open endpoint beside a protected one is now reached over the
  protected one;
- a username and password are **not** put on a channel that neither
  signs nor encrypts unless `allow-insecure-credentials=true`.

What remains the operator's job is **naming a trust anchor** —
`tls.trust-store` for chain validation, or `server-certificate-file` to
pin one certificate (§10 item 3) — and provisioning a client key pair
where the server authenticates its clients (§10 item 4). A connection
that names neither now fails loudly at connect time instead of coming
up unprotected, which is the intended behaviour and not a bug report.

The consequence for §13: "the OPC UA driver accepts an
attacker-presented certificate" is a **`VALID`** report — it would
describe a regression against the shipped default, not a configuration
choice — while "the OPC UA driver is unencrypted / unauthenticated by
default" is now a **`KNOWN-NON-FINDING`** describing a state the code
left behind. Check which release the report is against (§1: a report
against release *N* is triaged against the model as it stood at *N*).

## §6 Assumptions about inputs

### Per-entry-point trust table (Java public API)

| Entry point                                            | Parameter                    | Attacker-controllable?                                                                                                                                                      | Caller must enforce                                                                                              |
|--------------------------------------------------------|------------------------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------|------------------------------------------------------------------------------------------------------------------|
| `PlcDriverManager.getConnection(url, [auth])`          | `url`                        | **no** — trusted caller string                                                                                                                                              | only well-formed URLs reach this method; do not let end-users construct them                                     |
| `PlcDriverManager.getConnection(url, auth)`            | `auth` (`PlcAuthentication`) | **no** — trusted caller object                                                                                                                                              | credentials in the auth object are forwarded to the protocol (cleartext, except inside an OPC UA secure channel) |
| `connection.parseTagAddress(addr)`                     | `addr`                       | **conditionally** — trusted by default; if the embedding application sources tag strings from end users, the embedding application is responsible *(see §3 item 5; §14 Q5)* | tag-address syntax validation; do not pass un-vetted untrusted strings                                           |
| `readRequestBuilder.addTagAddress(name, addr)`         | `addr`                       | conditionally — same as above                                                                                                                                               | as above                                                                                                         |
| `writeRequestBuilder.addTagAddress(name, addr, value)` | `value`                      | **no** — trusted caller value                                                                                                                                               | application-level authorization over which tags may be written                                                   |
| `subscriptionRequestBuilder.add*(addr, ...)`           | `addr`                       | conditionally — same as above                                                                                                                                               | as above                                                                                                         |

### Per-network-input trust table (wire bytes)

| Surface                                                                                                                                                      | Parameter                          | Attacker-controllable?                                                       | Driver must enforce                                                                                                                                                                                |
|--------------------------------------------------------------------------------------------------------------------------------------------------------------|------------------------------------|------------------------------------------------------------------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| Cleartext OT-protocol response (Modbus, S7, BACnet/IP, IEC-60870-104, Profinet, CANopen, KNXnet/IP, C-Bus, DF1, AB-Ethernet, ADS payload, EtherNet/IP, UMAS) | every byte of every response frame | **yes** — anyone on the OT network with reachability can spoof               | **memory safety, bounded allocation, no infinite loop, no unbounded recursion** on malformed input — but **not** authenticity, integrity, or any payload-semantic guarantee *(inferred — §14 Q11)* |
| OPC UA wire frames inside the secure channel                                                                                                                 | every byte                         | **yes** — but signed/encrypted by the negotiated policy when policy ≠ `NONE` | as above plus: correct verification of signature and MAC under the negotiated policy; correct decryption; correct chunk reassembly *(documented + inferred — §14 Q17)*                             |
| OPC UA wire frames during discovery (`security-policy=NONE` phase)                                                                                           | every byte                         | **yes**                                                                      | memory safety on malformed responses; **the discovery handshake itself is not authenticated** *(documented)*                                                                                       |
| OPC UA server certificate (presented during handshake)                                                                                                       | full DER bytes                     | **yes**                                                                      | `tls.trust-store` → X.509 chain validation per JCE rules; `server-certificate-file` → the presented certificate must equal the pinned one; neither, with `tls.verify=true` (the default) → **rejected**; `tls.verify=false` → none *(documented: §5a)* |
| Serial-line frames                                                                                                                                           | every byte                         | yes if the serial channel is attacker-reachable                              | memory safety on malformed framing *(inferred — §14 Q11)*                                                                                                                                          |
| libpcap capture / replay frames                                                                                                                              | every byte                         | **no** — operator-supplied across B1, per §3 item 12 | the pcap-replay transport is a **dev/test tool**; if it's running in production, the integrator has put it there                                                                                   |
| ETS `.knxproj` XML                                                                                                                                           | as XML                             | **no** — operator-supplied across B1, per §3 item 12 | a corrupt or wrong file is a correctness matter (§3 item 12); XXE / external-DTD / external-schema stay disabled and the per-entry expansion bound stays enforced, and a bypass of either is `VALID` (§8 P6) |
| BACnet EDE file                                                                                                                                              | as text                            | **no** — operator-supplied across B1, per §3 item 12 | parser robustness as ordinary correctness; unbounded allocation from a bounded file is still `VALID` (§8 P4) |

### Size / shape / rate

- PLC4X drivers **stream individual request/response frames**. Beyond
  what the protocol itself bounds (Modbus PDUs are spec-bounded, S7 PDU
  size is negotiated, OPC UA chunk size is negotiated) *(documented: per
  protocol page; `S7Configuration.pduSize`,
  `OpcuaConfiguration.limits.encoding.{send,receive}-buffer-size`)*, the
  August 2026 hardening series added explicit bounds the model can now
  rely on: **1024 levels of type nesting** in every generated parser in
  all four bindings (`PLC4X_MAX_NESTING_DEPTH` moves it), OPC UA browse
  caps (`browse-max-references-per-node` / `-total-nodes` / `-depth`), a
  16 MiB per-tag budget in the simulated driver, per-driver ceilings on
  the element count a tag address may name, a plausibility check on
  declared frame lengths before the transport waits for them, and a
  bound on what the raw-socket capture may leave queued. A parser that
  cannot make progress now reports rather than looping.
- **PLC4X provides no per-call rate limit or backpressure on the API
  side** — if the embedding application calls `read` in a hot loop, the
  driver will issue as many requests as the underlying transport allows
  *(inferred — §14 Q19)*.

## §7 Adversary model

### Actors

| Actor                                                                                                                                                            | In scope?                                  | Capabilities granted                                                                                                                                                                                                                                                                 |
|------------------------------------------------------------------------------------------------------------------------------------------------------------------|--------------------------------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| **Network peer with reachability to the PLC** (the OT network)                                                                                                   | **yes — primary adversary for B2**         | can read, modify, drop, replay, or inject frames on the cleartext protocols. For the OPC UA protocol with a configured trust store, the peer is bounded by the cryptographic primitives. For OPC UA with `tls.verify=false`, the peer can MITM the channel; with the default (`tls.verify=true`) a connection without a configured trust anchor fails closed instead. |
| **Network peer between the driver and an OPC UA discovery endpoint**                                                                                             | **yes**                                    | the discovery handshake is unencrypted by spec *(documented)*; the peer can swap the advertised endpoint URL, certificate, or set of supported policies.                                                                                                                             |
| **Author of a malformed-but-parseable response from a real PLC firmware** (buggy or compromised device firmware sending bytes a well-behaved firmware would not) | **yes — wire-parser robustness must hold** | causes the driver to parse adversarial bytes; in-model for memory safety, hang, unbounded allocation.                                                                                                                                                                                |
| **Author of a malformed `.knxproj` or BACnet EDE file**                                                                                                          | **yes for parser robustness**              | causes the offline-file parser to crash, hang, or escape from its expected sandbox (XXE blocked per §5a).                                                                                                                                                                            |
| **Embedding application itself**                                                                                                                                 | **out of scope** — see §3 item 3           | trusted by construction.                                                                                                                                                                                                                                                             |
| **End user of the embedding application**                                                                                                                        | **out of scope** — see §3 item 4           | PLC4X has no concept of an end user.                                                                                                                                                                                                                                                 |
| **Operator of the PLC device**                                                                                                                                   | **out of scope**                           | the physical device the driver talks to is the device; PLC4X has no opinion on its operator.                                                                                                                                                                                         |
| **Author of a tag-address string**                                                                                                                               | conditionally in scope — see §3 item 5     | the embedding application decides whether to expose this surface.                                                                                                                                                                                                                    |
| **Side-channel observer** (cache timing, branch prediction, network timing of OPC UA crypto operations)                                                          | **out of scope** *(inferred — §14 Q20)*    | standard JCE side-channel posture; no PLC4X-level constant-time guarantee.                                                                                                                                                                                                           |
| **Quantum adversary**                                                                                                                                            | **out of scope**                           | the OPC UA security policies in §5a use classical RSA / AES / SHA.                                                                                                                                                                                                                   |

### Authenticated-but-Byzantine peer

PLC4X is a client-side library; there is no internal "cluster
membership" to be Byzantine within. The closest analogue is **a PLC
device whose firmware has been compromised and which is sending crafted
responses while passing whatever authentication the protocol provides**.
For cleartext protocols this is no different from any other attacker
on the wire (§3 item 1); for OPC UA, this is the "compromised but
properly-certified server" case — the driver still trusts the server's
*claims* about data, even if the channel is cryptographically sound
*(inferred — §14 Q21)*. A finding requiring a compromised PLC to
**produce wrong measurement values** is `OUT-OF-MODEL: trusted-input`
(the device is the source of truth for its own data); a finding
requiring a compromised PLC to **crash or corrupt the driver via
malformed wire bytes** is in-model for memory safety.

## §8 Security properties the project provides

For each property: condition, violation symptom, severity tier,
provenance. There are fewer §8 properties than for a typical
threat-model target precisely because PLC4X is a client-side library
sitting on top of OT protocols that themselves provide little — most
of the security work is delegated to the integrator (§10).

### P1 — OPC UA secure-channel confidentiality and integrity, when configured

- **Condition**: `security-policy` is one of `Basic128Rsa15`,
  `Basic256`, `Basic256Sha256`, `Aes128_Sha256_RsaOaep`, or
  `Aes256_Sha256_RsaPss` *(documented:
  `plc4j/drivers/opcua/.../security/SecurityPolicy.java`)* — since
  2026-08-24 `Basic256Sha256` is the **default**, so this half of the
  condition holds unless `NONE` was asked for — AND `message-security`
  is `SIGN` or `SIGN_ENCRYPT` (`SIGN_ENCRYPT` by default)
  *(documented)*, AND a trust anchor is configured: `tls.trust-store`
  containing the expected server certificate or its issuer, or
  `server-certificate-file` pinning that certificate *(documented:
  `website/.../protocols/opcua.adoc`, `OpcuaDriverContext.buildCertificateVerifier`)*.
- **Violation symptom**: wire bytes between the driver and the server
  that an on-path attacker can decrypt (without the configured key) or
  modify without the driver detecting the modification.
- **Severity**: **security-critical**, `VALID` per §13.
- *(documented + inferred — §14 Q17)*: the implementation of the
  policies is in repo; the property is documented by reference to the
  OPC UA spec, but no PLC4X-specific test or formal-verification artifact
  is cited.

### P2 — OPC UA server-certificate authentication, when a trust store is configured

- **Condition**: a trust anchor is configured — `tls.trust-store`
  containing the expected certificate or a chain root, or
  `server-certificate-file` naming the exact certificate to pin
  *(documented: `website/.../protocols/opcua.adoc`,
  `OpcuaDriverContext.buildCertificateVerifier`)*. With `tls.verify` at
  its default and neither one set, the driver rejects every server
  certificate, so the failure mode is a refused connection rather than
  an unverified one.
- **Violation symptom**: the OPC UA driver completes a handshake with a
  server presenting a certificate that does not chain to a trust-store
  root, and proceeds to exchange application-layer messages.
- **Severity**: **security-critical**, `VALID` per §13.
- *(documented)*.

### P3 — Memory safety of the Java public API and Java SPI on well-formed embedding-application calls

- **Condition**: the embedding application observes the documented
  Java contract (Javadoc): correct types, non-`null` where required,
  `close()` semantics observed.
- **Violation symptom**: `NullPointerException` / `IllegalStateException`
  in PLC4X internals from a documented-correct API call (not a JVM
  crash — Java is GC'd); for the Go bindings,
  panic from an `unsafe`-free call site; for `plc4c`, native crash.
- **Severity**: **correctness-only** for non-driver code; **security-critical**
  when a driver path is reachable from B2 (network input). See P4.
- *(inferred — §14 Q22)*.

### P4 — Memory safety, bounded allocation, and bounded loop / recursion in the per-driver wire parsers, on adversarial wire bytes

- **Condition**: the parser is reachable via B2 (network input) from
  a remote PLC or an attacker on the OT network with reachability;
  the host JVM / Go runtime / C `malloc` follows its documented
  semantics.
- **Violation symptom**: in Java — JVM heap exhaustion driven by a
  bounded-input expansion factor, parser hang on crafted input,
  `StackOverflowError` from unbounded recursion on a nested wire
  structure; in Go — `goroutine` hang, unbounded allocation; in C
  — heap / stack corruption, OOB read/write, use-after-free.
- **Severity**: **security-critical** (the OT network is the §7 primary
  adversary), `VALID` per §13.
- *(inferred — §14 Q11; partly documented since 2026-08-24)*. **This is
  still the most reachable §8 property for a scan target**, but it is no
  longer an untested aspiration: every generated parser refuses more
  than 1024 levels of type nesting in all four bindings; the Java and Go
  testsuites parse **every truncation of every testsuite message**, so a
  message that stops short is a reported parse failure rather than a
  hang, an out-of-range read or a half-built object; a manual array that
  makes no progress terminates; an announced optional field must
  actually be present; and an unchecked exception out of a parser is
  contained in the receive path instead of ending the thread that owns
  the channel. A finding here is still `VALID` — the point is that the
  obvious classes have been swept, so a new one is worth reporting.

### P5 — Tag-address parser correctness and bounded resource use, on caller-supplied address strings

- **Condition**: the embedding application calls
  `parseTagAddress` / `addTagAddress` per the documented per-driver
  syntax.
- **Violation symptom**: parser crash, hang, or unbounded allocation
  on a syntactically-malformed address; cross-driver interference
  (a tag-address string for driver A causing observable state change
  in driver B).
- **Severity**: **`VALID-HARDENING`** by default (the address string is
  trusted per §6); **`VALID`** when the parser is reachable from
  unbounded user-supplied strings AND the embedding application's
  documented use of the API includes that case.
- *(inferred — §14 Q5)*.

### P6 — XXE / external-DTD / external-schema disabled in the KNX `.knxproj` parser

- **Condition**: an integrator loads a `.knxproj` ETS file through the
  KNX driver.
- **Violation symptom**: an XML processor in `EtsParser` resolves an
  external entity, external DTD, or external schema, allowing
  file-disclosure / SSRF from the `.knxproj`.
- **Severity**: **security-critical**, `VALID` per §13.
- *(documented: `EtsParser.java` lines 168–171 — the line numbers moved
  as the file grew)*. Alongside it the parser now bounds how far a
  single `.knxproj` zip entry may expand, so a declared inflated size is
  no longer taken at face value, and cleans up what it unpacked.

### P7 — OPC UA secure-channel implementation tracks the OPC UA spec's security-policy URIs

- **Condition**: an OPC UA security policy ≠ `NONE` is selected.
- **Violation symptom**: a policy-spec-compliant OPC UA server rejects
  the driver's handshake; or the driver's handshake completes with a
  policy weaker than the URI implies; or the driver picks an
  algorithm not in the policy's algorithm-set.
- **Severity**: **security-critical**, `VALID` per §13.
- *(documented + inferred — §14 Q17)*. Compatibility against Eclipse
  Milo and OPC Foundation .NET is asserted in the docs
  *(`website/.../protocols/opcua.adoc`)*; a property break against
  those reference servers is a §8 P7 violation.

## §9 Security properties the project does *not* provide

These are the false-friend properties — features that *look like* a
security property but are not one — and the well-known classes of
attack the library cannot defend against. **This section is the most
important one for an integrator.**

### General disclaimers

- **No authenticity, integrity, or confidentiality on Modbus, S7, BACnet/IP,
  IEC-60870-5-104, Profinet, CANopen, KNXnet/IP, C-Bus, DF1,
  AB-Ethernet, EtherNet/IP, UMAS, or ADS/AMS wire traffic.** PLC4X
  implements the protocol as standardized; the protocols themselves
  do not have those properties *(documented per protocol; inferred —
  §14 Q3)*. → §3 item 1, §11a.
- **No built-in transport-layer tunneling.** PLC4X does not transparently
  wrap a cleartext protocol in TLS, VPN, IPSec, or Wireguard. If the
  integrator wants encrypted-on-the-wire Modbus, they must run Modbus
  inside an integrator-provided tunnel *(inferred — §14 Q8)*.
- **No replay-protection on cleartext protocols.** A request the driver
  sent at *t* can be replayed by an attacker at *t+δ*; PLC4X does not
  add nonces / sequence-checks of its own. Where a protocol specifies
  nonces, PLC4X provides those protocol-defined nonces but adds none
  beyond what the protocol mandates *(maintainer — chrisdutz)*.
- **No authentication of the embedding application's end users.** PLC4X
  has no concept of end users. See §3 item 4.
- **No authorization over which tags an embedding application reads or
  writes.** The embedding application owns that decision. See §3 items 3, 4.
  Where the target protocol has its own permission system, PLC4X **proxies**
  it — deferring the check to the target PLC and relaying any
  permission error back to the client — but provides no permission
  system of its own *(maintainer — chrisdutz)*.
- **No defense against malformed `.knxproj` ZIP slip, BACnet EDE
  malformed-file robustness, or path-traversal in operator-supplied
  file paths.** The file paths are caller-supplied; XML XXE is the only
  attack class explicitly hardened *(inferred — §14 Q18)*.
- **No DoS protection at the API surface.** The embedding application
  is trusted to call the API at a reasonable rate
  *(inferred — §14 Q19)*.
- **No data-at-rest encryption.** PLC4X does not persist anything to
  disk on its own at runtime. (The SPI3 rewrite adds an optional
  Audit-Log feature that writes debug data to the filesystem — intended
  for transient debugging, not permanent operation, with no encryption;
  enabling it in production is an operator choice *(maintainer — chrisdutz)*.)
- **No constant-time comparison of authentication secrets.** Whatever
  the protocol provides (or doesn't) is what PLC4X passes through.
- **No defense against side-channel observation of OPC UA crypto.** Per
  JCE / underlying-OpenSSL posture; no PLC4X-level mitigations
  *(inferred — §14 Q20)*.
- **No defender stance against a compromised PLC firmware sending
  spec-compliant but semantically wrong data.** Garbage-in / garbage-out
  per §7 / §3 item 7.

### False-friend properties (call out separately)

- **"OPC UA is encrypted" — true by default since 2026-08-24, but only
  once a trust anchor is named.** The driver now defaults to
  `security-policy=Basic256Sha256`, `message-security=SIGN_ENCRYPT` and
  `tls.verify=true`, so the channel an OPC UA reader expects is what
  the driver asks for. What it cannot supply for the operator is the
  **trust anchor**: with none configured the connection is refused
  rather than completed unverified (§5a). The false friend that
  remains is the reverse of the old one — `security-policy=NONE` and
  `tls.verify=false` are still expressible, still documented, and a
  configuration carrying either has none of these properties whatever
  the word "OPC UA" suggests. *(documented: `OpcuaConfiguration.java`,
  `OpcuaDriverContext.java`; maintainer rulings §14 Q14, Q15)*
- **OPC UA discovery is unencrypted by spec.** The OPC UA *spec*
  conducts discovery over `security-policy=NONE`. An attacker on the
  path during discovery can rewrite advertised endpoints, server
  certificates, and supported policies before the driver chooses
  what to honor. This is not a PLC4X bug; this is OPC UA *(inferred —
  §14 Q16)*. But a
  scanner that reports "discovery handshake is unauthenticated" against
  PLC4X is reporting on the protocol, not on the driver
  *(documented)*.
- **Modbus / S7 PUT-GET / EtherNet-IP / ADS / etc. running over TCP
  port `:N` is not "secure" because TCP is.** The transport-level
  encryption is absent; the protocol-level encryption is absent. The
  protocol-spec authentication is either absent (Modbus) or weak (S7
  PUT/GET *(inferred — §14 Q23)*) or device-specific (ADS/AMS route
  setup).
- **The `simulated://`, `mock://` drivers are dev-tools, not sandboxes.**
  They exist to let an embedding application be unit-tested without a
  real PLC; they are not isolation boundaries.
- **The `pcap-replay` transport is a developer / regression-test tool.**
  It replays captured frames; it is not designed for production use
  and a finding against it that depends on an attacker controlling the
  pcap file is `OUT-OF-MODEL: unsupported-component`.
- **`plc4c/`, `plc4py/`, `plc4net/` are README-flagged "not ready for
  usage".** A finding against them is `OUT-OF-MODEL:
  unsupported-component` until the README line changes (§5a).
- **`PlcUsernamePasswordAuthentication` on ADS / OPC UA carries
  credentials on the wire as the protocol dictates.** It is not a
  PLC4X-level security primitive — it is a typed wrapper around two
  strings that the driver forwards to the protocol's authentication
  step. Whether it is then encrypted (OPC UA with a secure channel) or
  cleartext (ADS payload, OPC UA with `security-policy=NONE`) depends
  on the rest of the configuration.
- **`PlcCertificateAuthentication`** wraps an in-memory KeyStore for
  certificate-based authentication; the security properties of the
  surrounding TLS / OPC UA channel still apply.

### Well-known attack classes against this category of project that PLC4X does not defend against

- **MITM on cleartext OT protocols.** Out of model per §3 item 1.
- **Replay of OT commands** (a recorded "open valve" frame replayed
  later). Out of model per §3 item 1.
- **Spoofing of a Modbus master / S7 PG / BACnet master.** Out of model
  per §3 item 1.
- **Compression / decoding bombs in `.knxproj` ZIP files.** ZIP-slip
  in the unzip step is *(inferred — §14 Q18)*.
- **XML billion-laughs in `.knxproj`.** Mitigated by
  `FEATURE_SECURE_PROCESSING` *(documented: `EtsParser.java`)*. →
  `KNOWN-NON-FINDING`.
- **OPC UA discovery-channel MITM.** Out of model per §9
  false-friend item 2.
- **OPC UA server-certificate substitution when `tls.trust-store` is
  unset (the default).** **Pending maintainer ruling** (§14 Q15).
- **DoS via unbounded-tag-list `read` / `subscribe` requests.** Out of
  model per §9 "No DoS protection at the API surface".
- **Memory corruption in the per-protocol wire parser on a crafted
  response from an attacker-controlled or compromised PLC.** This **is**
  in model — see §8 P4. The most likely category of real findings.

## §10 Downstream responsibilities

The embedding application / integrator deploying PLC4X in production
**must**:

1. **Treat the OT network as the security perimeter for every cleartext
   protocol.** Operate Modbus, S7, BACnet/IP, IEC-60870-104, Profinet,
   CANopen, KNXnet/IP, C-Bus, DF1, AB-Ethernet, EtherNet/IP, UMAS,
   and ADS/AMS payload traffic on an OT VLAN, behind a firewall, with
   no untrusted-network egress, ideally with monitoring. PLC4X
   provides no transport security for these protocols, by design
   (§9). *(integrator)*
2. **If a cleartext protocol must traverse a hostile path, deploy an
   external tunnel** (site-to-site VPN, IPSec, stunnel/TLS-wrapping
   proxy, or operate the cleartext protocol exclusively inside a
   point-to-point cellular-modem channel). PLC4X does not bundle
   such a tunnel. *(integrator)*
3. **For OPC UA in production: name a trust anchor.** The channel is
   already signed and encrypted by default (`security-policy=Basic256Sha256`,
   `message-security=SIGN_ENCRYPT`, `tls.verify=true`), but the driver
   cannot know which server it should trust: configure a
   `tls.trust-store` holding the expected server certificate or its
   issuing CA, or `server-certificate-file` to pin exactly one
   certificate. Without either, the connection is refused — that is the
   intended behaviour, not a defect *(documented: `OpcuaConfiguration.java`,
   `OpcuaDriverContext.buildCertificateVerifier`; maintainer §14 Q15)*.
   Do **not** reach for `tls.verify=false` to get past it; that is the
   dev/lab escape hatch and it disables the check entirely.
4. **For OPC UA where the server authenticates its clients: provision a
   client key pair and load it via `tls.keystore` /
   `tls.keystore-password`** (with `tls.keystore-type`, and
   `generated-key-size` where a server demands 4096-bit keys). The
   self-signed certificate the driver generates when no key store is
   configured is now well-formed enough for a server to accept, but it
   is regenerated on every start and so cannot satisfy a peer that
   requires a stable client identity. Where the server wants an X509
   *user* identity rather than a username and password, supply it as
   `PlcCertificateAuthentication` *(documented:
   `website/.../protocols/opcua.adoc`; GH-1845)*.
5. **For OPC UA where discovery is over a hostile path: disable
   `discovery` and configure the endpoint explicitly via
   `endpoint-host`, `endpoint-port`.** The discovery handshake is
   unauthenticated by OPC UA spec *(documented)*.
6. **Vet tag-address strings before passing them to PLC4X if they are
   sourced from end users of the embedding application.** PLC4X does
   not provide a security boundary against attacker-supplied
   tag-address strings (§3 item 5). *(integrator)*
7. **Authorize which tags may be read or written by which end user
   inside the embedding application.** PLC4X has no end-user
   identity. *(integrator)*
8. **Bound the request rate / queue depth on the embedding-application
   side if DoS resistance against a misbehaving caller is needed.**
   PLC4X has no API-level throttle *(inferred — §14 Q19)*. *(integrator)*
9. **Run the embedding process with the minimum host privileges the
   chosen transport requires.** Most transports need only a
   non-privileged user; the `raw-socket` transport requires
   `CAP_NET_RAW` (or Administrator) — grant it to the embedding
   process specifically rather than running as `root`. *(integrator)*
10. **Treat the `simulated://`, `mock://`, `pcap-replay` drivers as
    dev / test fixtures, not as production endpoints.** *(integrator)*
10a. **If `audit-log-file` is set, protect the file it writes.** The
    audit log records request and response traces — plant data, and
    whatever a payload happens to carry — to a rolling file on the
    host. Give it a path only the embedding process can read, and a
    retention policy. It is off unless configured *(documented:
    `plc4j/utils/audit-log/README.md`)*. *(integrator)*
11. **Do not deploy `plc4c/`, `plc4py/`, or `plc4net/` artifacts in
    production until the README line ("not ready for usage") is
    removed.** *(integrator; documented: `README.md`)*
12. **Trust the PLC firmware only as far as the integrator independently
    trusts the device.** The semantic content of values returned by the
    PLC is **device-side trusted-input** per §3 item 7 — but a buggy
    firmware sending malformed wire bytes is still in-model for
    parser robustness (§8 P4). *(integrator)*

## §11 Known misuse patterns

- **Exposing the embedding-application process directly to an
  untrusted network and trusting the OT-network protocol's "auth" to
  hold.** Modbus, S7 PUT/GET, etc., do not have authentication. The
  OT-network perimeter (§10 item 1) is the *only* control.
- **Asking for `security-policy=NONE` in production.** No longer the
  default, but still expressible — and a configuration carrying it has
  an unsigned, unencrypted channel to an unauthenticated server. → §9
  false-friend item 1.
- **Reaching for `tls.verify=false` to get a connection to come up.**
  This is the shape the fail-closed default provokes: the connection is
  refused because no trust anchor is configured, and turning
  verification off makes it succeed. It installs
  `PermissiveCertificateVerifier`, so an MITM gets the encryption key —
  the very property the signed-and-encrypted channel was for. Name a
  `tls.trust-store` or a `server-certificate-file` instead. → §5a,
  §10 item 3.
- **Carrying `ignore-common-name=true` on the TLS transport past the
  device that needed it.** It accepts a certificate issued for any
  other host, which is what a machine in the middle needs. → §5a.
- **Treating `discovery=true` as part of the security boundary.**
  Discovery is over `security-policy=NONE` by OPC UA spec. → §9
  false-friend.
- **Using `simulated://` or `mock://` in production code paths.** They
  are dev tools.
- **Passing end-user-supplied tag-address strings to `parseTagAddress`
  without vetting.** A parser bug becomes a user-reachable bug.
- **Deploying the `plc4c/` or `plc4py/` artifact in production.** The
  README flags them.
- **Relying on `PlcCertificateAuthentication` without TLS / OPC UA
  secure channel underneath.** The KeyStore exists; the credential is
  there; whether anything on the wire is actually encrypted depends on
  the surrounding configuration.
- **Running with `--add-exports jdk.compiler/...` flags in the
  production JVM.** Those flags are for the **build**, not the
  runtime artifact *(documented: `README.md`)*; granting them at
  runtime exposes JDK internals unnecessarily.
- **Bundling the embedding application into a single host on an OT
  segment and assuming the host is "behind the firewall, so it's
  fine".** PLC4X does no integrity verification on its own outbound
  traffic. A compromised host with PLC4X loaded can write arbitrary
  values to any reachable PLC.

## §11a Known non-findings (recurring false positives)

This section is the highest-leverage input for automated agentic
security scans. Each entry: tool symptom, why it is safe under the
model, the § that licenses the call. **PLC4X's §11a is unusually
long because the project's core domain is unauthenticated OT
protocols, and a scanner unaware of that fact will report essentially
every wire-format property of every driver.**

### Protocol-by-design non-findings

- **"Modbus has no authentication on the wire."** Modbus does not
  define authentication. → §3 item 1, §9. `BY-DESIGN:
  protocol-disclaimed`.
- **"Modbus is sent in cleartext."** Modbus does not define
  encryption. → §3 item 1, §9. `BY-DESIGN: protocol-disclaimed`.
- **"S7 PUT/GET sends commands without authenticating the master."**
  PUT/GET is the spec's "client/server data block access" mode; it
  does not authenticate the master. The S7 driver here speaks PUT/GET
  by design *(documented: `website/.../protocols/s7.adoc` —
  "S7-300/S7-400 controllers, as well as basic reading and writing
  functions for the S7-1200 and S7-1500 devices (PUT/GET
  functions)")*. → §3 item 1, §9. `BY-DESIGN: protocol-disclaimed`.
- **"BACnet/IP has no authentication, integrity, or confidentiality."**
  Classic BACnet/IP (the version implemented here) defines none of
  those properties. (BACnet/SC, the secure variant, is a different
  protocol and is not in PLC4X.) → §3 item 1, §9. `BY-DESIGN`.
- **"IEC-60870-5-104 is cleartext and unauthenticated."** Spec.
  (IEC-62351 is the secure overlay and is not in PLC4X.) → §3 item 1,
  §9. `BY-DESIGN`.
- **"Profinet IO frames are not authenticated."** Profinet's spec
  splits "Profinet" into multiple application protocols; the variant
  here does not include the optional Profinet-Security extensions. →
  §3 item 1, §9. `BY-DESIGN`.
- **"CANopen / SocketCAN frames are unauthenticated."** CAN is a
  fieldbus; the integrity property is "this frame originated on this
  physical bus", not cryptographic. → §3 item 1, §9. `BY-DESIGN`.
- **"KNXnet/IP / KNX Tunnelling has no authentication."** Classic
  KNXnet/IP does not. (KNX Secure is the secure overlay and is
  *(inferred — §14 Q24)* not implemented in this driver. Confirm.) →
  §3 item 1, §9. `BY-DESIGN`.
- **"C-Bus, DF1, AB-Ethernet, EtherNet/IP (classic CIP), UMAS, ADS/AMS
  payload traffic, Open-Protocol — all unauthenticated /
  cleartext."** Per spec. → §3 item 1, §9. `BY-DESIGN`.
- **"Firmata is a serial protocol with no authentication."**
  Firmata is a hobbyist Arduino protocol; the security model is
  "the serial cable is your perimeter". → §3 item 1, §9. `BY-DESIGN`.

### Driver-default non-findings

- **"OPC UA driver accepts any server certificate."** **No longer true
  of any default.** Since 2026-08-24 `tls.verify` defaults to `true`
  and a connection with no trust anchor is refused
  (`RejectingCertificateVerifier`); certificates are accepted without
  validation only under an explicit `tls.verify=false`, which logs a
  warning. A report describing the old behaviour against current code
  is a `KNOWN-NON-FINDING`; a report showing that current code still
  accepts an attacker-presented certificate with the defaults in place
  is **`VALID`** — check which release the report is against
  *(maintainer §14 Q15; `OpcuaDriverContext.buildCertificateVerifier`)*.
- **"OPC UA driver is unencrypted by default / `security-policy=NONE`."**
  Stale: the default is `Basic256Sha256` with `SIGN_ENCRYPT`.
  `KNOWN-NON-FINDING` against current code. → §5a.
- **"OPC UA connection fails when no trust store is configured."**
  Working as designed — fail-closed is the point (§10 item 3). Not a
  bug, and not a reason to recommend `tls.verify=false`.
  `KNOWN-NON-FINDING`.
- **"OPC UA driver does discovery in cleartext."** True by spec; not a
  driver bug. The driver no longer *continues* from discovery onto a
  channel weaker than the configured one, and never treats a
  certificate learned there as a trust anchor. → §9 false-friend
  item 2.
- **"OPC UA driver auto-generates a self-signed client certificate."**
  Documented behavior when `tls.keystore` is unset
  *(`website/.../protocols/opcua.adoc`)*. Since 2026-08-14 the
  generated certificate is SHA-256 signed, end-entity marked, with a
  positive 16-octet serial and a subject key identifier; it is still a
  "get-it-working" convenience and is regenerated on every start, so
  production should provision an explicit one. → §10 item 4.
  `OUT-OF-MODEL: non-default-build`.
- **"OPC UA driver uses `Basic128Rsa15`/`Basic256` policies which are
  deprecated by the OPC Foundation."** Yes — the driver supports them
  *(documented: `SecurityPolicy.java`)* because some PLC firmware
  still only supports them. The operator is responsible for choosing
  `Basic256Sha256` or stronger per §10 item 3 *(inferred — §14 Q25)*.
  → `OUT-OF-MODEL: non-default-build`.

### Code-base non-findings

- **"XXE / external-DTD / external-schema processing in the KNX
  parser."** Disabled *(documented: `EtsParser.java` lines 65–71)*.
  → `KNOWN-NON-FINDING`.
- **"Hardcoded password / token in `plc4j/utils/test-utils/`,
  `tests/`, `testdata/`, `simulated://`, `mock://` drivers."** These
  are test infrastructure. → §3 item 6. `OUT-OF-MODEL:
  unsupported-component`.
- **"Cryptographic finding in `plc4c/` / `plc4py/` / `plc4net/`."**
  README flags these as not-ready. → §3 item 6. `OUT-OF-MODEL:
  unsupported-component` pending the README line being removed.
- **"Driver X mishandles a malformed frame / leaks a socket / skips a
  check the other drivers make",** where the protocol support table
  does not mark X with a green check for that binding (as of
  2026-09-03, for PLC4J: BACnet/IP, CAN, C-Bus, ctrlX, DF1,
  Open-Protocol, Profinet; for PLC4Go additionally AB-Ethernet,
  CANopen, Firmata, OPC-UA, PLC4X-proxy). The driver is declared
  partial, and in PLC4J it is not deployed to Maven Central at all. →
  §3 item 11. `OUT-OF-MODEL: unsupported-component` until its cell
  turns green.
- **"`code-generation/` mspec compiler has X."** Build-time only,
  not in the runtime artifact. → §3 item 6. `OUT-OF-MODEL:
  unsupported-component`.
- **"`Dockerfile`, `docker-compose.yaml`, `Jenkinsfile` has X."**
  Repo infrastructure. → §3 item 6.
- **"OPC-UA Server / PLC4X Server / Calcite / Camel / Kafka-Connect /
  NiFi adapter has X."** Out of this repo (in `plc4x-extras`). →
  §3 item 2. `OUT-OF-MODEL: unsupported-component`.
- **"Vendored Eclipse Milo / Netty has CVE-Y."** Neither is a
  dependency any more: the OPC UA driver has been native since 0.9.0,
  and Netty left with the SPI3 transports — no `io.netty` artifact is
  on the runtime classpath *(verified 2026-09-03)*. A report naming
  them against current PLC4X is a `KNOWN-NON-FINDING`.
- **"Dependency X (Jackson, commons-*, Bouncy Castle, …) has CVE-Y."**
  Report upstream; PLC4X picks up fixes via dependency-version bumps.
  `OUT-OF-MODEL: unsupported-component` *(inferred — §14 Q26)*.
- **"SQL injection / XSS / SSRF in the website asciidoc tree."**
  Documentation, not the deployed library. → §3 item 6.

### Operator-supplied-input non-findings

- **"A corrupt / truncated / hand-edited `.knxproj` makes the KNX driver
  throw."** The file comes from the operator across B1. → §3 item 12.
  `OUT-OF-MODEL: trusted-input` (fix as ordinary robustness). The
  exceptions stay: an entity resolved out of that file, or an entry that
  inflates without bound, is `VALID` (§8 P6).
- **"A malformed EDE file / GSDML directory / pcap capture / key store
  makes the driver fail or produce wrong tags."** Same shape, same
  answer. → §3 item 12. `OUT-OF-MODEL: trusted-input`.
- **"I pointed the driver at a server I control and it believed what the
  server said."** Wrong values, a fabricated browse tree, a controller
  claiming a type it does not have — the operator chose the endpoint. →
  §3 item 13. `OUT-OF-MODEL: trusted-input`. What those bytes do to the
  *parser* is a different report and stays `VALID` (§8 P4), as does an
  OPC UA channel that accepts a certificate it should have rejected
  (§8 P1, P2).
- **"The driver sends the configured credentials to whatever host the
  connection string names."** That is the function. → §3 item 13.
  `OUT-OF-MODEL: trusted-input`. (Sending them over a channel that
  neither signs nor encrypts is a *different* claim, and the OPC UA
  driver refuses that by default — see §5a
  `allow-insecure-credentials`.)

### Integrator-deployment non-findings

- **"PLC4X runs as root."** PLC4X does not require root; the embedding
  application chose the UID. (See §10 item 9 for `CAP_NET_RAW`.)
  `OUT-OF-MODEL: adversary-not-in-scope`.
- **"PLC4X is reachable from the internet."** PLC4X does not bind
  listening sockets; "reachable" means the embedding application is
  reachable. `OUT-OF-MODEL: adversary-not-in-scope`.
- **"PLC4X writes to a Modbus coil that controls a physical relief
  valve."** The integrator authorized the API call; the physical
  consequence is integrator-side. → §3 item 7. `OUT-OF-MODEL:
  trusted-input`.

## §12 Conditions that would change this model

Revise this document when any of the following lands:

- A new protocol driver gains a built-in transport-security layer
  (e.g. Modbus-Security, OPC UA Pub-Sub with Security, KNX Secure,
  BACnet/SC). Changes §3 item 1, §9, §11a per-protocol.
- A new transport gains a listening-socket capability (would convert
  PLC4X from a pure client to a client/server).
- A change in the default value of any §5a knob — especially OPC UA
  `security-policy`, `tls.trust-store`, or `message-security`.
- `plc4c/`, `plc4py/`, or `plc4net/` has its "not ready for usage"
  README flag removed — those move from §3 item 6 to in-model.
- A driver's cell in the protocol support table
  (`website/asciidoc/modules/users/pages/protocols/index.adoc`) changes
  to or from a green check — that moves the driver into or out of §3
  item 11, and for PLC4J the `maven.deploy.skip` flag in its POM has to
  follow the same way.
- The boundary with `plc4x-extras` shifts (e.g. an integration is
  brought back into the main repo, or vice versa). This has already
  happened once without the document following: the connection cache,
  OPM and the scraper's successor live in `plc4j/tools/` here.
- A module family appears or disappears under `plc4j/` — the §2
  component table is an inventory and goes stale silently. It was last
  reconciled against the tree on 2026-09-03.
- The transport stack changes shape (it moved off Netty to virtual-thread
  blocking I/O in SPI3, which invalidated the §5 concurrency text).
- A driver starts treating one of the operator-supplied files of §3
  item 12 as untrusted input — e.g. a `.knxproj` fetched from a device
  or a network share rather than named by the operator. That moves the
  file across B1 into B2 and the carve-out no longer applies to it.
- A vulnerability report that cannot be cleanly routed to one of the
  §13 dispositions: that is evidence the model has a gap.

## §13 Triage dispositions

A report against PLC4X receives exactly one of the following.

| Disposition                            | Meaning                                                                                                                                                                                                                                                                                                        | Licensed by          |
|----------------------------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|----------------------|
| `VALID`                                | Violates a §8 property via an in-scope §7 adversary using an in-scope §6 input — typically a wire-parser memory-safety break (P4), an OPC UA secure-channel break (P1/P2/P7), or KNX-XML XXE bypass (P6).                                                                                                      | §8, §6, §7           |
| `VALID-HARDENING`                      | No §8 property violated, but a §11 misuse pattern can be made harder to fall into by code change (e.g. issue a runtime warning when OPC UA runs `security-policy=NONE`). Fixed at maintainer discretion, typically no CVE.                                                                                     | §11                  |
| `OUT-OF-MODEL: trusted-input`          | Requires attacker control of a §6 parameter the model marks trusted — the connection URL, the auth object, the value being written, tag-address strings (per §3 item 5 when the integrator has not opted them into untrusted-string handling), a configuration file the operator points the driver at (§3 item 12), or the choice of endpoint and what that endpoint claims (§3 item 13). | §6, §3 items 12–13 |
| `OUT-OF-MODEL: adversary-not-in-scope` | Requires a §7 actor the model excludes — embedding application is hostile, side-channel observer, quantum adversary, OT-network perimeter is the integrator's problem.                                                                                                                                         | §7                   |
| `OUT-OF-MODEL: unsupported-component`  | Lands in `plc4c/`, `plc4py/`, `plc4net/`, `tools/`, `code-generation/`, `plc4j/utils/`, `plc4j/drivers/{simulated,mock}`, `plc4x-extras` content, vendored upstream code, repo infrastructure, or a driver the protocol support table does not mark with a green check for that binding. | §3 item 2, §3 item 6, §3 item 11 |
| `OUT-OF-MODEL: non-default-build`      | Only manifests under a §5a knob the maintainer has ruled is dev/test — for OPC UA that now means an explicit `security-policy=NONE`, `tls.verify=false` or `allow-insecure-credentials=true`, and on the TLS transport `ignore-common-name=true`. **Note:** the OPC UA secure-by-default work the §14 Q14/Q15 rulings called for shipped on 2026-08-24, so a report about the *defaults* is either a `KNOWN-NON-FINDING` (it describes the pre-2026-08-24 behaviour) or `VALID` (it shows the current defaults failing); it is no longer excluded here. | §5a |
| `BY-DESIGN: protocol-disclaimed`       | Concerns a property the **protocol** (not the library) does not provide — every "Modbus / S7 / BACnet / etc. is unauthenticated" report.                                                                                                                                                                       | §9, §3 item 1, §11a  |
| `BY-DESIGN: property-disclaimed`       | Concerns a §9 property the library explicitly does not provide (built-in TLS tunneling, end-user authn, DoS protection at the API).                                                                                                                                                                            | §9                   |
| `KNOWN-NON-FINDING`                    | Matches a §11a recurring false positive.                                                                                                                                                                                                                                                                       | §11a                 |
| `MODEL-GAP`                            | Cannot be cleanly routed to any of the above — triggers §12 model revision.                                                                                                                                                                                                                                    | §12                  |

## §14 Open questions for the maintainers

Every *(inferred)* tag in the body maps to one of these. Proposed
answers are inline; please confirm, correct, or strike.

**PMC review (chrisdutz, 2026-06-04, PR-approved):** the OPC UA insecure
defaults (Q14/Q15/Q16) are moving to secure-by-default — the permissive
certificate verifier in particular is a gap to fix, not a supported
posture — and the SPI3 rewrite adds TLS transport, hardens the parsers,
and bounds per-connection allocation (Q8/Q11/Q18/Q19). Answers folded
below and into the body as *(maintainer)*.

### Wave 1 — scope, intended use, the OPC UA defaults

**Q1.** Confirm that "remote PLC bytes on the wire are untrusted input
crossing into the parser" is the right framing for §2 / §4 — i.e.,
the driver's job is to be memory-safe and bounded-resource against any
byte sequence a remote PLC (or attacker on the OT net) could send,
even though the protocol itself provides no authentication.
*(maps to §2, §4 B2, §7, §8 P4)*

**Q2.** Confirm the unsupported-component list: `plc4j/utils/`,
`plc4j/drivers/{simulated,mock}`, `code-generation/`, `tools/`,
`images/`, `media/`, `website/`, `licenses/`, top-level
`Dockerfile`/`docker-compose.yaml`/`Jenkinsfile`. Anything to add or
remove? *(maps to §3 item 6)*

**Q3.** Confirm the per-protocol enumeration of "no auth / no
integrity / no confidentiality" in §3 item 1 and §11a. The
proposed list: Modbus (all variants), S7 (PUT/GET), BACnet/IP
(classic), IEC-60870-5-104, Profinet (the variant implemented),
CANopen, KNXnet/IP (classic), C-Bus, DF1, AB-Ethernet, EtherNet/IP
(classic CIP), UMAS, ADS/AMS payload, Open-Protocol, Firmata. Have we
missed any, or wrongly tarred one that does have some authentication
the driver enforces? *(maps to §3 item 1, §11a)*

**Q4.** Confirm "the embedding application is trusted, end users of
the embedding application are out of scope" *(proposed: yes)*. Are
there any embedded-product / OEM contexts where PLC4X considers the
embedding application's end users in scope? *(maps to §3 items 3, 4)*

**Q5.** Tag-address strings: confirm the default stance is "trusted
caller input — embedding application owns vetting". Is there a
PLC4X-side memory-safety / hang / unbounded-allocation guarantee
on tag-address strings? Proposed: `VALID-HARDENING` for parser
robustness, `VALID` only if the embedding application's documented
contract includes user-supplied strings. *(maps to §3 item 5, §8 P5)*

**Q6.** Confirm out-of-scope status of PLC-device firmware bugs and
physical-process consequences. *(maps to §3 item 7)*

**Q7.** OT-network perimeter as integrator responsibility — confirm
the §10 item 1 framing. *(maps to §3 item 8, §10)*

### Wave 2 — OPC UA defaults (highest-priority maintainer ruling)

**Q14.** OPC UA `security-policy=NONE` default. Is "OPC UA used with
`security-policy=NONE` in production" a `VALID` report against PLC4X
(stance (a) — default is supported posture) or `OUT-OF-MODEL:
non-default-build` (stance (b) — operator must flip per §10 item 3)?
**Proposed: (b)** — the default is convenience for dev/lab. If (a),
the documentation needs to make clear that "OPC UA" with this driver
does not by default carry OPC UA's spec-level cryptographic
properties.

**Answered (maintainer — chrisdutz):** secure-by-default is the intended
posture — the SPI3 rewrite makes the insecure path explicitly opt-in and
the secure path the new default; the current `NONE` default is dev/lab
convenience, **not** a supported production posture. *(maps to §5a, §10, §11a, §13)*

**Shipped 2026-08-24 — question closed.** `security-policy` now defaults
to `Basic256Sha256` with `message-security=SIGN_ENCRYPT`, and `NONE` is
an explicit opt-in. "OPC UA runs unencrypted by default" is a
`KNOWN-NON-FINDING` against current code; asking for `NONE` is a §11
misuse pattern.

**Q15.** OPC UA default `PermissiveCertificateVerifier` — the
single highest-priority question. When `tls.trust-store` is unset,
the driver accepts every server certificate. Is "OPC UA driver accepts
attacker-presented certificate" `VALID` (stance (a)) or `OUT-OF-MODEL:
non-default-build` (stance (b))? **Answered (maintainer — chrisdutz):** the permissive default "should be
changed and reported" — it is **not** the supported posture. A report
that the OPC UA driver accepts an attacker-presented certificate is
therefore **`VALID`** (a security gap the PMC intends to fix toward
secure-by-default), not `OUT-OF-MODEL: non-default-build`. *(maps to §5a, §10 item 3, §11a, §13)*

**Shipped 2026-08-24 — question closed.** `tls.verify` defaults to
`true`; the verifier is chosen as trust-store → pinned certificate →
reject, so a connection with no trust anchor now **fails closed**.
`PermissiveCertificateVerifier` survives only behind an explicit
`tls.verify=false`, which warns. The disposition stands: if current
code with default settings accepts an attacker-presented certificate,
that is `VALID`.

**Q16.** OPC UA `discovery=true` over `security-policy=NONE`. The
spec mandates discovery in cleartext; the driver follows the spec.
Confirm that "OPC UA discovery handshake unauthenticated" is `BY-DESIGN:
protocol-disclaimed` per §9 false-friend item 2 — not a PLC4X bug.

**Answered (maintainer — chrisdutz):** confirmed — cleartext discovery is
per the OPC UA spec (`BY-DESIGN: protocol-disclaimed`); separately, the
broader OPC UA channel is moving to secure-by-default in SPI3 (see Q14/Q15).
*(maps to §3 item 1, §9, §10 item 5)*

**Shipped 2026-08-24 — question closed.** Discovery is still conducted
with `NONE`, as the spec has it, but the driver no longer proceeds from
it onto a channel weaker than the configured one, and a certificate
learned during discovery is checked against the configured verifier
rather than trusted for being the one that arrived. The remaining
exposure — an on-path attacker rewriting advertised endpoints — is the
protocol's, and §10 item 5 (disable `discovery`, configure
`endpoint-host` / `endpoint-port`) is still the remedy.

### Wave 3 — wire-parser robustness, the most likely site of real findings

**Q11.** Confirm §8 P4 (memory safety, bounded allocation, bounded
loop / recursion in per-driver wire parsers on adversarial wire bytes)
is a property PLC4X has actually committed to — i.e., a Modbus parser
that allocates `O(N)` memory in response to a 1-byte length field with
no upper bound is a bug, not "just OT-protocol weirdness". Are there
specific drivers where this property has been deliberately weakened
(e.g. generated code that is faster but not bounded)?

**Answered (maintainer — chrisdutz):** P4 is committed; in the SPI3
rewrite each connection allocates a fixed-length ring-buffer (length
varying per protocol), preventing the "huge fake message → huge
allocation" class of issue. *(maps to §4 reachability, §8 P4, §11a)*

**Q17.** OPC UA secure-channel implementation correctness — §8 P1 /
P2 / P7. Are these properties tested against the cited reference
servers (Eclipse Milo, OPC Foundation .NET) on every release, or
only manually at version-bump time? Are there OPC UA test vectors
(spec-published handshake traces) the driver is regressed against?
*(maps to §8 P1, P2, P7)*

**Q18.** `.knxproj` (ZIP) handling: confirm ZIP-slip protection in the
unzip step. BACnet EDE files: confirm parser robustness against
malformed files.

**Answered (maintainer — chrisdutz):** the SPI3 rewrite addresses the
ETS-parser issues and hardens XML parsing. *(maps to §6, §9, §11a)*

**Q19.** No API-level throttle confirmed? Proposed: yes — the
embedding application is responsible.

**Answered (maintainer — chrisdutz):** confirmed — no API-level throttle;
the embedding application is responsible. (SPI3's per-connection
ring-buffer additionally bounds per-connection allocation.) *(maps to §6, §9, §10 item 8)*

**Q20.** Side-channel posture for OPC UA crypto: out of scope
*(proposed)*. *(maps to §7, §9)*

**Q21.** Compromised PLC firmware sending crafted but
spec-compliant responses — confirm the split: wrong **measurement
values** are `OUT-OF-MODEL: trusted-input` (the device is the source
of truth), but wrong **wire framing** is in-model for parser
robustness. *(maps to §3 item 7, §7, §8 P4)*

**Q22.** §8 P3 (Java/Go/C-binding API memory safety): is the right
framing "Java public API is memory-safe by virtue of being Java,
unless a JNI path is reached", or do you make a stronger claim?
*(maps to §8 P3)*

### Wave 4 — environment, false-friends, edges

**Q8.** External tunneling (VPN / IPSec / TLS-wrap) is integrator's
responsibility — confirm §9 / §10 item 2. Is there appetite for a
"TLS-wrap transport" that would terminate at PLC4X (e.g.
Modbus-over-TLS per Schneider Electric's draft spec)? If yes, that
changes §12.

**Answered (maintainer — chrisdutz):** yes — the SPI3 rewrite adds `tls`
and `tls-psk` transports that can secure a `tcp` connection where the
target PLC/gateway supports it. This is a §12-changing addition. *(maps to §9, §10 item 2, §12)*

**Q9.** ADS authentication: when `PlcUsernamePasswordAuthentication`
is supplied, the credential drives an AMS-route registration
*(`AdsProtocolLogic` lines 130–145)*. Confirm that this is forwarded
in whatever transport mode TwinCAT runs (typically cleartext); PLC4X
does not encrypt it. *(maps to §4 B2-ADS, §5a)*

**Q10.** `raw-socket` transport requires `CAP_NET_RAW` — confirm
this is the only transport with elevated host-privilege needs. Any
others? *(maps to §3, §10 item 9)*

**Q12.** Thread-safety of `PlcConnection`: confirm the documented
contract (proposed: connections are not safe for concurrent use by
multiple threads without the embedding application's own
synchronization; the `CompletableFuture` model returns to the
caller's thread). *(maps to §5)*

**Q13.** Confirm the "what PLC4X does NOT do to its host" inventory
in §5: no signal handlers (besides Netty / JVM), no child-process
spawn at runtime, no listening sockets, no security-sensitive env-var
reads, no global-state mutation. Any exceptions? *(maps to §5)*

**Q23.** S7 PUT/GET vs S7 protected modes: the docs reference
"S7-Plus" as a future target. The current driver speaks PUT/GET
exclusively *(documented: `website/.../protocols/s7.adoc`)*. Confirm
that "S7 password protection on the PLC" (the device-side `PG` /
`OS` access password) is **out-of-band** to the driver — the driver
will speak PUT/GET to whatever device responds. *(maps to §3 item 1,
§9, §11a)*

**Q24.** KNX Secure: is the KNX driver implementing KNX Secure (the
cryptographic overlay on KNXnet/IP), or only classic KNXnet/IP?
Proposed: only classic. *(maps to §11a)*

**Q25.** OPC UA legacy policies (`Basic128Rsa15`, `Basic256`):
deprecated by the OPC Foundation but supported here for PLC-firmware
compat. Confirm the stance: supported for connectivity reasons,
operator chooses `Basic256Sha256`-and-up per §10 item 3.
*(maps to §5a, §10 item 3, §11a)*

**Q26.** Vendored / dependency CVEs (Netty, Jackson, Eclipse Milo
crypto, commons-*) — confirm the policy is "report upstream, PLC4X
picks up the fix via the next dependency bump". *(maps to §3 item 6,
§11a)*

### Wave 5 — meta

**Q27.** Should this document live at `docs/threat-model.md` /
`website/asciidoc/modules/users/pages/threat-model.adoc`, or somewhere
else? *(meta)*

**Q28.** Is there an existing threat-model document (Confluence wiki,
internal note) that this should reconcile against rather than
supersede? The website `security.adoc` page is process-only (where
to report); the `developers/maturity.adoc` page references the
security process but has no threat-model content. *(meta — §3.1a
of the rubric)*

**Q29.** `plc4x-extras` boundary: confirm that this document is
strictly for the main `apache/plc4x` repo and that a separate model
will be produced for `plc4x-extras` when that repo is brought into
the program. *(maps to §1, §3 item 2)*

**Q30.** §11a is the highest-leverage section for automated triage
and is heavily protocol-disclaimer-shaped here. Could the PMC
contribute 3–5 patterns observed in inbound reports / mailing-list
threads that recur — e.g. "scanner reports unauthenticated Modbus,
we close every time"; "scanner reports plaintext password in
OPC-UA-over-`security-policy=NONE`, we close every time"? Concrete
recurring false-positives strengthen the suppression list. *(meta —
§11a)*

**Q31.** What kind of change to PLC4X should trigger a revision
(proposed list in §12 — confirm or correct)? *(meta, §12)*

---

## Appendix: Existing security-policy artefact → §x back-map

PLC4X does not currently ship an in-repo `SECURITY.md`. The de facto
security-policy artefacts are:

- `website/asciidoc/modules/users/pages/security.adoc` — the **process-only**
  page, which says (verbatim): "For more information about reporting
  vulnerabilities, see the Apache Security Team page. … No vulnerabilities
  have been reported." No threat-model content; everything maps to the §1
  reporting cross-reference.
- `website/asciidoc/modules/developers/pages/maturity.adoc` — restates the
  process posture (QU20, QU30); no threat-model content.
- Per-protocol pages under
  `website/asciidoc/modules/users/pages/protocols/*.adoc` — these are the
  best **documented** sources for protocol-level security caveats. The
  back-map below covers the load-bearing ones.

| Source                                                                                                                                                                                                                                                                                                                                                              | Claim                               | Lands in                                                  |
|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|-------------------------------------|-----------------------------------------------------------|
| `README.md` ("`plc4c/`, `plc4py/` not ready for usage; `plc4net/` abandoned")                                                                                                                                                                                                                                                                                       | scope carve-out                     | §3 item 6, §5a                                            |
| `README.md` ("The Industrial IoT adapter … client-side library across multiple PLC protocols")                                                                                                                                                                                                                                                                      | scope framing                       | §2                                                        |
| `plc4j/drivers/opcua/.../context/OpcuaDriverContext.java` (`buildCertificateVerifier`: verify → trust-store → pinned → reject) | fail-closed certificate verification (supersedes the "accepted without validation" wording the website carried before 2026-08-24) | §4 B2-OPCUA, §5a, §9 false-friend item 1, §11a |
| `website/.../protocols/opcua.adoc` ("discovery phase is always conducted using `NONE` security policy" — paraphrased from `OpcuaConfiguration.java` discovery doc)                                                                                                                                                                                                  | discovery unencrypted               | §4 B2-OPCUA-PERMISSIVE, §9 false-friend item 2            |
| `website/.../protocols/opcua.adoc` ("`message-security` … `SIGN_ENCRYPT` … high security settings and full encryption")                                                                                                                                                                                                                                             | secure-channel message security     | §5a, §8 P1                                                |
| `website/.../protocols/opcua.adoc` ("User authentication" / "User certificate": anonymous, username-password and X509 user certificate via `PlcCertificateAuthentication`) | client X509 user auth is implemented — supersedes the older "currently unsupported by PLC4X" wording, removed with the implementation in GH-1845 | §5a, §10 item 4, §11a |
| `website/.../protocols/opcua.adoc` (compatibility list: Eclipse Milo, OPC Foundation .NET, etc.)                                                                                                                                                                                                                                                                    | tested reference servers            | §8 P7                                                     |
| `website/.../protocols/s7.adoc` ("PUT/GET functions")                                                                                                                                                                                                                                                                                                               | S7 mode                             | §3 item 1, §11a, §14 Q23                                  |
| `website/.../protocols/s7.adoc` ("Siemens `LOGO` device … requires `?controller-type=LOGO`")                                                                                                                                                                                                                                                                        | functional knob                     | §5a (non-security row)                                    |
| `website/.../protocols/modbus.adoc` (no security section)                                                                                                                                                                                                                                                                                                           | Modbus is unauthenticated           | §3 item 1, §11a (disclaimer-by-omission)                  |
| `website/.../protocols/canopen.adoc` ("CANopen … address areas")                                                                                                                                                                                                                                                                                                    | CANopen scope                       | §3 item 1, §11a                                           |
| `website/.../protocols/ads.adoc` ("device-independent and fieldbus independent interface for communication between Beckhoff automation devices")                                                                                                                                                                                                                    | ADS scope                           | §3 item 1, §11a                                           |
| `website/.../protocols/index.adoc`                                                                                                                                                                                                                                                                                                                                  | per-language driver-coverage matrix | §2 component table, §3 item 11, §11a, §12 |
| `website/asciidoc/.../security.adoc` ("For more information about reporting vulnerabilities, see the Apache Security Team page")                                                                                                                                                                                                                                    | reporting channel                   | §1 reporting cross-reference                              |
| `website/asciidoc/.../developers/maturity.adoc` (QU20, QU30)                                                                                                                                                                                                                                                                                                        | maturity posture                    | §1 reporting cross-reference                              |
| `plc4j/drivers/opcua/.../security/{Permissive,TrustStore,Pinned,Rejecting}CertificateVerifier.java` | the four verifier strategies; `Rejecting` is the no-anchor default, `Permissive` needs `tls.verify=false` | §4 B2-OPCUA, §5a, §8 P2 |
| `plc4j/drivers/opcua/.../security/SecurityPolicy.java` (enum of `NONE`, `Basic128Rsa15`, `Basic256`, `Basic256Sha256`, `Aes128_Sha256_RsaOaep`, `Aes256_Sha256_RsaPss`)                                                                                                                                                                                             | supported policies                  | §5a, §8 P1, P7, §11a, §14 Q25                             |
| `plc4j/drivers/opcua/.../config/OpcuaConfiguration.java` (`@ConfigurationParameter` set incl. `tls.verify=true`, `tls.trust-store`, `server-certificate-file`, `tls.keystore`, `generated-key-size=2048`, `discovery=true`, `username`, `password`, `allow-insecure-credentials=false`, `security-policy=Basic256Sha256`, `message-security=SIGN_ENCRYPT`, `channel-lifetime-ms=3600000`, `min-channel-lifetime-ms=5000`, `session-timeout-ms=120000`, `handshake-timeout-ms=60000`, `request-timeout-ms=30000`, `browse-max-*`, `subscription-queue-size=1`) | knob inventory (verified against `website/.../partials/opcua.adoc`, generated from this class, on 2026-09-03) | §5a, §6, §10                                              |
| `plc4j/drivers/knxnetip/.../ets/EtsParser.java` lines 65–71 (XXE/DTD/schema disabled, `FEATURE_SECURE_PROCESSING`)                                                                                                                                                                                                                                                  | hardened XML parser                 | §5 (filesystem-KNX), §8 P6, §11a code-base                |
| `plc4j/drivers/ads/.../AdsProtocolLogic.java` lines 130–145 (`PlcUsernamePasswordAuthentication` instance check; `setupAmsRoute(...)` call)                                                                                                                                                                                                                         | ADS credential forwarding           | §4 B2-ADS, §5a                                            |
| `plc4j/api/.../authentication/PlcAuthentication.java`, `PlcUsernamePasswordAuthentication.java`, `PlcCertificateAuthentication.java`                                                                                                                                                                                                                                | auth type-system                    | §6, §9 false-friend                                       |
| `plc4j/api/.../PlcConnection.java`, `PlcDriverManager.java`                                                                                                                                                                                                                                                                                                         | public API surface                  | §2, §4 B1, §6                                             |
| `plc4j/transports/{api,tcp,udp,serial,cotp,tls,raw-socket,pcap-replay,can,can-socketcan,can-virtualcan,test}` (no listening-socket entries; TCP is blocking `SocketChannel` on a virtual thread since SPI3, not Netty)                                                                                                                                                                                                                                              | client-only transports              | §2, §3 item 2, §10 item 9                                 |
| `plc4j/transports/tls/.../config/TlsTransportConfiguration.java` (`verify=true`, `ignore-common-name=false`, `trust-store`, `trust-store-password`, `trust-store-type`, `version`, `keystore*`) | TLS transport trust knobs; hostname checking added 2026-08-24 | §5a, §10 item 3, §11 |
