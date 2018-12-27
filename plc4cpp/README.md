<!--

  Licensed to the Apache Software Foundation (ASF) under one or more
  contributor license agreements.  See the NOTICE file distributed with
  this work for additional information regarding copyright ownership.
  The ASF licenses this file to You under the Apache License, Version 2.0
  (the "License"); you may not use this file except in compliance with
  the License.  You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.

-->
# Building PLC4CPP

On mac and linux machines there should not be a requirement to install any prerequisites.
In general the build requires `gcc` being installed.

On Windows machines, you might need to install it manually.
He have tested everything with the bundle of http://win-builds.org/doku.php
Make sure the `bin` directory containing the executable `mingw32-make.exe` is on your systems `PATH`.

In the `plc4cpp-libs` module all third party dependencies will be built that are needed by any of the other plc4cpp modules.
The build is setup to install any built binaries to the `libs` directory of the `plc4cpp-libs` module and to only build them if they have previously not been installed there.
So please refer to that directory in any of the other builds.