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

The .NET implementation of [Apache PLC4X](https://plc4x.apache.org). This
package set currently provides the API, SPI runtime, and a loopback test
transport. Protocol drivers and production transports are separate follow-up
contributions.

## Packages

| Package | Contents |
|---|---|
| `plc4net-api` | connection / driver / value interfaces |
| `plc4net-spi` | driver runtime, value model, buffers, `ConnectionString` |
| `plc4net-transports-test` | in-memory loopback transport for tests |

## Status

This is the first buildable slice of the plc4net revival described in
[apache/plc4x#2656](https://github.com/apache/plc4x/pull/2656). It does not
yet include a production transport or a protocol driver. Pre-release builds only.

Licensed under the Apache License 2.0.
