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

go 1.24

require (
	github.com/IBM/netaddr v1.5.0
	github.com/ajankovic/xdiff v0.0.1
	github.com/cstockton/go-conv v1.0.0
	github.com/fatih/color v1.18.0
	github.com/google/uuid v1.6.0
	github.com/gopacket/gopacket v1.3.1
	github.com/icza/bitio v1.1.0
	github.com/jacobsa/go-serial v0.0.0-20180131005756-15cf729a72d4
	github.com/libp2p/go-reuseport v0.4.0
	github.com/pkg/errors v0.9.1
	github.com/rs/zerolog v1.34.0
	github.com/snksoft/crc v1.1.0
	github.com/stretchr/testify v1.10.0
	github.com/subchen/go-xmldom v1.1.2
	github.com/viney-shih/go-lock v1.1.2
	golang.org/x/net v0.38.0
	golang.org/x/text v0.24.0
	golang.org/x/tools v0.31.0
)

require (
	github.com/antchfx/xpath v0.0.0-20170515025933-1f3266e77307 // indirect
	github.com/bitfield/gotestdox v0.2.2 // indirect
	github.com/chigopher/pathlib v0.19.1 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/dnephin/pflag v1.0.7 // indirect
	github.com/fsnotify/fsnotify v1.8.0 // indirect
	github.com/google/shlex v0.0.0-20191202100458-e7afc7fbc510 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/huandu/xstrings v1.4.0 // indirect
	github.com/iancoleman/strcase v0.3.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/incu6us/goimports-reviser/v3 v3.9.1 // indirect
	github.com/jinzhu/copier v0.4.0 // indirect
	github.com/magiconair/properties v1.8.9 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/pelletier/go-toml/v2 v2.2.3 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/sagikazarmark/locafero v0.7.0 // indirect
	github.com/sagikazarmark/slog-shim v0.1.0 // indirect
	github.com/sourcegraph/conc v0.3.0 // indirect
	github.com/spf13/afero v1.12.0 // indirect
	github.com/spf13/cast v1.7.1 // indirect
	github.com/spf13/cobra v1.8.1 // indirect
	github.com/spf13/pflag v1.0.6 // indirect
	github.com/spf13/viper v1.19.0 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	github.com/vektra/mockery/v2 v2.53.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/exp v0.0.0-20250210185358-939b2ce775ac // indirect
	golang.org/x/mod v0.24.0 // indirect
	golang.org/x/sync v0.13.0 // indirect
	golang.org/x/sys v0.31.0 // indirect
	golang.org/x/term v0.30.0 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	gotest.tools/gotestsum v1.12.0 // indirect
)

tool golang.org/x/tools/cmd/stringer

tool gotest.tools/gotestsum

tool github.com/vektra/mockery/v2

tool github.com/incu6us/goimports-reviser/v3

tool github.com/apache/plc4x/plc4go/tools/plc4xGenerator

tool github.com/apache/plc4x/plc4go/tools/plc4xLicencer
