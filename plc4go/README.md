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

# plc4go

## plc4go module

To use plc4go import `"github.com/apache/plc4x/plc4go/pkg/api"` and use the driver manager 
`driverManager := plc4go.NewPlcDriverManager()` to register a driver e.g. 
`drivers.RegisterKnxDriver(driverManager)`.

## plc4go tools

### plc4xpcapanalyzer

`plc4xpcapanalyzer` is a small tool to evaluate the plc4x parsing against a pcap and report the success rate.

To install `plc4xpcapanalyzer` run: 
`go install github.com/apache/plc4x/plc4go/tools/plc4xpcapanalyzer@latest`

### plc4xbrowser

`plc4xbrowser` is a terminal application to run `REPL`-commands against a PLC using plc4x (plc4go)

To install `plc4xbrowser` run: 
`go install github.com/apache/plc4x/plc4go/tools/plc4xbrowser@latest`
