<!--
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements. See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License. You may obtain a copy of the License at

https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
-->

# PLC4Net Modbus

The protocol code selects framing independently of the transport:

| Protocol | Transports | Default transport |
| --- | --- | --- |
| `modbus-tcp` | `tcp`, `udp`, `tls` | `tcp` |
| `modbus-rtu` | `serial`, `tcp`, `udp`, `tls` | `serial` |
| `modbus-ascii` | `serial`, `tcp`, `udp`, `tls` | `serial` |

For example, `modbus-rtu:tcp://localhost:502` sends RTU frames over TCP.
The peer must support the selected combination. Network transports use port
502 by default; TLS uses port 802. TLS validates the server certificate and
hostname using platform trust by default. Applications constructing a
`TcpTransportConfiguration` directly can supply `TlsChainPolicy` for an
explicit private trust policy. There is no URL option to disable validation.

## Reads and file records

Reads combine tags within an address window with the same area, unit identifier
and byte order, including gaps between tags, matching the Java and Go optimizers.
The window starts at the first tag and defaults to 125 registers or 2,000 bits.
Use `max-registers-per-request` and `max-coils-per-request` to lower the merge
window for a device. Zero uses the default; values above the protocol ceiling
are clamped. These options control merging and do not split an individual tag.
Gap addresses must be readable by the device. Lower the window to avoid merging
across unsupported addresses. Bit results use each tag's actual bit offset.
Extended-register windows are additionally limited to 123 registers and split
at file boundaries.

Extended-register flat addresses map to file `address / 10000 + 1` and record
`address % 10000`. Selections crossing a file boundary are split into separate
requests. The last supported flat address is 655349999; a selection extending
beyond that address is rejected.

## Ping, subscriptions and identification

`ModbusConnectionBase.PingAsync` reads `ping-address`, which defaults to
`holding-register:1:INT`. Choose an address supported by the device. Timeout,
cancellation and device errors are propagated to the caller.

Subscriptions use the shared polling subscription manager. Cyclic subscriptions
publish sampled results; change-of-state subscriptions publish changes after
their initial baseline. The default change-of-state sampling interval is
1,000 ms, configurable with `subscription-polling-interval-ms`. Changes between
samples can be missed. Closing the connection stops its subscriptions.

`ModbusTcpConnection.IdentifyAsync` reads basic device identification using
FC 43/14 and follows paginated replies. It returns raw bytes indexed by object
ID and rejects stale or duplicate objects and non-advancing pagination.
`ModbusTcpDriver.DiscoverAsync` accepts one connection string or an explicit
list of endpoints. The list overload probes sequentially and returns either
objects or an error for each endpoint. This is explicit endpoint probing,
not subnet scanning or broadcast discovery. Devices without FC 43/14 support
can still support ordinary reads and writes.

The regression suite uses fixed wire examples, test transports and local
TCP/UDP/TLS peers. Physical PLC, serial-adapter and gateway interoperability
still require hardware validation.
