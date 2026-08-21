//
// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//      https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
//

module github.com/apache/plc4x/plc4go

go 1.26

toolchain go1.26.3

require (
	github.com/ajankovic/xdiff v0.0.1
	github.com/fatih/color v1.19.0
	github.com/google/uuid v1.6.0
	github.com/gopacket/gopacket v1.7.1
	github.com/rs/zerolog v1.35.1
	github.com/stretchr/testify v1.12.1
	golang.org/x/net v0.58.0
	golang.org/x/sys v0.47.0
	golang.org/x/text v0.41.0
	golang.org/x/tools v0.49.0
)

require (
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/mattn/go-colorable v0.1.15 // indirect
	github.com/mattn/go-isatty v0.0.22 // indirect
	github.com/stretchr/objx v0.5.3 // indirect
	go.yaml.in/yaml/v3 v3.0.5 // indirect
	golang.org/x/mod v0.39.0 // indirect
	golang.org/x/sync v0.22.0 // indirect
)

tool golang.org/x/tools/cmd/stringer

tool github.com/apache/plc4x/plc4go/tools/plc4xGenerator

tool github.com/apache/plc4x/plc4go/tools/plc4xLicencer
