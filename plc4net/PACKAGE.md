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
# Apache PLC4X for .NET (plc4net)

The .NET implementation of [Apache PLC4X](https://plc4x.apache.org) — a set of
industrial-protocol drivers behind one API.

```csharp
using org.apache.plc4net.drivers.s7;
using org.apache.plc4net.spi.transports;

var driver = new S7Driver(new DefaultTransportManager());
using var connection = driver.Connect("s7://192.168.0.1?remote-rack=0&remote-slot=1");

var request = connection.ReadRequestBuilder
    .AddTagAddress("temperature", "%DB1.DBD0")
    .Build();
var response = (org.apache.plc4net.spi.drivers.messages.DefaultPlcReadResponse)
    await request.ExecuteAsync();
```

## Packages

| Package | Contents |
|---|---|
| `plc4net-api` | connection / driver / value interfaces |
| `plc4net-spi` | driver runtime, value model, buffers, `ConnectionString` |
| `plc4net-driver-modbus` | Modbus TCP + RTU |
| `plc4net-driver-s7` | Siemens S7 (COTP) |
| `plc4net-driver-knxnetip` | KNXnet/IP tunnelling |
| `plc4net-transports-{tcp,udp,cotp,serial,test}` | one transport each |

A driver package pulls the SPI, API and its transports transitively.

## Status

Reviving a long-dormant port — see
[apache/plc4x#2656](https://github.com/apache/plc4x/pull/2656). The S7 driver is
verified against a real S7-1200; Modbus RTU and KNXnet/IP are unit-tested but not
yet hardware-verified. Pre-release builds only.

Licensed under the Apache License 2.0.
