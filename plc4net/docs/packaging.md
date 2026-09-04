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
# Packaging

Every library project sets `<GeneratePackageOnBuild>true</GeneratePackageOnBuild>`,
so a `Release` build already drops a `.nupkg` in each `bin/Release/`. Shared
metadata (`Authors`, `Copyright`, `PackageLicenseExpression`, `Version`) lives in
`Directory.Build.props`.

| Package | Contents |
|---|---|
| `plc4net-api` | `IPlcConnection` / `IPlcDriver` / `IPlcValue`, message + tag interfaces |
| `plc4net-spi` | `DriverBase`, `ConnectionBase`, the value model, `ReadBuffer` / `WriteBuffer`, `ConnectionString` |
| `plc4net-transports-tcp` / `-udp` / `-cotp` / `-serial` / `-test` | one transport each |
| `plc4net-driver-modbus` | Modbus TCP + RTU |
| `plc4net-driver-s7` | S7 (COTP) |
| `plc4net-driver-knxnetip` | KNXnet/IP tunnelling |
| `plc4net-code-gen` | the `.mspec` → C# generator (build-time tool, not needed at run time) |

Project-to-project references become package dependencies automatically, so a
consumer that references `plc4net-driver-s7` transitively pulls `plc4net-spi`,
`plc4net-api`, `plc4net-transports-tcp` and `plc4net-transports-cotp`.

`s7-verify` and `modbus-verify` are `PackAsTool` packages
(`plc4net-s7-verify`, `plc4net-modbus-verify`) — `dotnet pack` bundles every
dependency into the one `.nupkg`, so each is a self-contained `dotnet tool`.

## Version

`Directory.Build.props` pins `<Version>1.1.0-SNAPSHOT</Version>` — the Maven
reactor version (`plc4x/pom.xml`). The Maven build overrides it with
`-p:Version=${project.version}`. **Do not hard-code a different number** (a bare
`0.0.1` desyncs from the reactor and reads as a stable release). For a local test
build, override at pack time instead — nothing committed changes:

```bash
dotnet pack plc4net.sln -c Release -p:Version=0.0.1-test.1 -o ./_stage
```

## Local feed round trip

This is the supported way to test `<PackageReference>` consumption without
publishing anything. It is a verified round trip.

```bash
cd plc4x/plc4net

# 1. pack every packable project at a throwaway prerelease version
dotnet pack plc4net.sln -c Release -p:Version=0.0.1-test.1 -o ./_stage

# 2. push to a local folder feed — same command shape as pushing to nuget.org
mkdir -p ./_localfeed
for p in ./_stage/*.nupkg; do dotnet nuget push "$p" -s "$(pwd)/_localfeed"; done

# 3. consume from a SEPARATE project
mkdir /tmp/consumer && cd /tmp/consumer
cat > nuget.config <<'XML'
<?xml version="1.0" encoding="utf-8"?>
<configuration>
  <packageSources>
    <clear />
    <add key="plc4net-local" value="/ABS/PATH/TO/plc4net/_localfeed" />
    <add key="nuget.org" value="https://api.nuget.org/v3/index.json" />
  </packageSources>
</configuration>
XML
dotnet new console
dotnet add package plc4net-driver-s7 --version 0.0.1-test.1
dotnet add package plc4net-transports-test --version 0.0.1-test.1
# ... write code against new S7Driver(new DefaultTransportManager()) ...
dotnet run
```

`nuget.org` stays in the source list because transitive dependencies
(`Microsoft.Extensions.Logging.Abstractions`, `System.IO.Ports`) are not in the
local feed.

> **Windows MAX_PATH.** Keep the consumer project and the NuGet cache off a deep
> path. The nested native-runtime asset paths for `System.IO.Ports`
> (`runtime.<rid>.runtime.native.system.io.ports/...`) can exceed 260 characters
> and fail to copy. A path like `C:\Users\<you>\consumer` is safe; a deep
> `%TEMP%\...` one is not.

## Pre-release feed (GitHub Packages)

For consumers who need a build before the ASF ships one, a fork can publish to
its own GitHub Packages NuGet feed. The package ids do not change, so a project
switches feeds without touching its `<PackageReference>`s.

Give each build a unique prerelease version — GitHub Packages, like nuget.org,
refuses to overwrite one:

```bash
V="1.1.0-preview.$(date +%Y%m%d%H%M)"
dotnet pack plc4net/plc4net.sln -c Release -p:Version="$V" -o ./_stage

# A classic PAT with write:packages. In GitHub Actions, GITHUB_TOKEN already has it.
dotnet nuget add source "https://nuget.pkg.github.com/<OWNER>/index.json" \
  --name gh --username <OWNER> --password "$GH_TOKEN" --store-password-in-clear-text
for p in ./_stage/*.nupkg; do dotnet nuget push "$p" -s gh --skip-duplicate; done
```

Consuming it needs a PAT too (GitHub Packages authenticates every read, even a
public feed):

```xml
<!-- nuget.config -->
<packageSources>
  <clear />
  <add key="gh" value="https://nuget.pkg.github.com/<OWNER>/index.json" />
  <add key="nuget.org" value="https://api.nuget.org/v3/index.json" />
</packageSources>
<packageSourceCredentials>
  <gh>
    <add key="Username" value="<OWNER>" />
    <add key="ClearTextPassword" value="%GH_TOKEN%" />
  </gh>
</packageSourceCredentials>
```

## Publishing to nuget.org

Not done from this branch. .NET package publishing for an Apache project goes
through the ASF release process (PMC vote, signed artifacts, the project's
NuGet account), the same as the Maven artifacts. Until then the version stays
`1.1.0-SNAPSHOT` and the `plc4net-*` ids stay unclaimed on nuget.org.
