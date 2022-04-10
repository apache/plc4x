<!--
  Licensed to the Apache Software Foundation (ASF) under one
  or more contributor license agreements.  See the NOTICE file
  distributed with this work for additional information
  regarding copyright ownership.  The ASF licenses this file
  to you under the Apache License, Version 2.0 (the
  "License"); you may not use this file except in compliance
  with the License.  You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing,
  software distributed under the License is distributed on an
  "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
  KIND, either express or implied.  See the License for the
  specific language governing permissions and limitations
  under the License.
  -->
# Building PLC4CPP

On mac and linux machines there should not be a requirement to install any prerequisites.
In general the build requires `gcc` being installed.

On Windows machines, you might need to install it manually.
He have tested everything with the bundle of 
http://win-builds.org/doku.php/download_and_installation_from_windows
When running the installer, make sure to select the options:
- Native Windows
- x86_64
Not quite sure which elements are really needed, better just install all of them.
Make sure the `bin` directory containing the executable `mingw32-make.exe` is on your systems `PATH`.

## Building PLC4CPP in Visual Studio

1. Open Visual Studio (Tested with 2017 and 2019)
2. Open Local Directory
3. Select the `plc4cpp` subdirectory
4. Wait till the IDE is ready (little icon in the lower left corner)
5. Build everything by selecting the menu `Build/Build All` (or similar)

## Building PLC4CPP in CLion

1. Open CLion
2. Open Local Directory
3. Select the `plc4cpp` subdirectory
4. Wait till the IDE is ready (little icon in the lower center)
5. Build everything by selecting the menu `Build/Build All in Debug` (or similar)