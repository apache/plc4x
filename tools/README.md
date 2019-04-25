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
# Linux

On a clean Ubuntu 18.04 the following software needs to be installed:

    sudo apt install bison flex
         
# MAC

Make sure `Homebrew` ist installed in order to update `Bison` to a newer version (the version 2.3 installed per default is too old)
    
    /usr/bin/ruby -e "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)"

Then update `Bison`: 

    brew install bison
    brew link bison --force

# Windows

Some tools need to be installed before being able to build on Windows:

- WinBuilds
- Bison
- Flex

He have tested WinBuilds with the bundle of: http://win-builds.org/doku.php/download_and_installation_from_windows

For Bison, please download the Setup installer version from here: http://gnuwin32.sourceforge.net/packages/bison.htm (When using the zip version the bison.exe couldn't find some DLL files)

Please download the Flex compiler from here: http://gnuwin32.sourceforge.net/packages/flex.htm (Ideally download the binary zip distribution)  

Make sure the `bin` directories of containing the executables `mingw32-make.exe`, `bison.exe` and `flex.exe` are all on your systems `PATH`.
