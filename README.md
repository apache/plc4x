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
[![Maven central](https://img.shields.io/maven-central/v/org.apache.plc4x/plc4j-api.svg)](https://img.shields.io/maven-central/v/org.apache.plc4x/plc4j-api.svg)
[![License](https://img.shields.io/github/license/apache/plc4x.svg)](http://www.apache.org/licenses/LICENSE-2.0)
[![Last commit](https://img.shields.io/github/last-commit/apache/plc4x.svg)]()
[![Twitter](https://img.shields.io/twitter/follow/ApachePLC4X.svg?label=Follow&style=social)](https://twitter.com/ApachePLC4X)


<h1 align="center">
  <br>
   <img src="https://plc4x.apache.org/images/apache_plc4x_logo.png" 
   alt="Apache PLC4X Logo" title="Apache PLC4X Logo"/>
  <br>
</h1>
<h3 align="center">The Industrial IoT adapter</h3>
<h4 align="center">The ultimate goal of PLC4X is to create a set of libraries, that allow unified access to any type of
 PLC</h4>

***

# Table of contents

  * [About PLC4X](#about-apache-plc4x)
  * [Getting Started](#getting-started)
  * [Developers](#developers)
  * [Community](#community)
  * [Contributing](#contributing)
  * [Licensing](#licensing)

***

## About Apache PLC4X

Apache PLC4X is an effort to create a set of libraries for communicating with industrial grade programmable logic controllers (PLCs) in a uniform way.
We are planning on shipping libraries for usage in:

1. Java
2. Scala
3. C/C++
4. C# (.Net)
5. Python

PLC4X also integrates with other Apache projects, such as:

* [Apache Calcite](https://calcite.apache.org/)
* [Apache Camel](https://camel.apache.org/)
* [Apache Edgent](https://edgent.apache.org/)
* [Apache Kafka-Connect](https://kafka.apache.org)
* [Apache Karaf](https://karaf.apache.org/)
* [Apache NiFi](https://nifi.apache.org/)

## Getting started

Depending on the programming language, the usage will differ, therefore please go to the 
[Getting Started](https://plc4x.apache.org/users/gettingstarted.html) on the PLC4X website to look up 
the language of choice.

### Java

See the PLC4J user guide on the website to start using PLC4X in your Java application:
[https://plc4x.apache.org/plc4j/users/gettingstarted.html](https://plc4x.apache.org/plc4j/users/gettingstarted.html)

## Developers

### Environment

Currently the project is configured to require the following software:

1. Java 8 JDK: For running Maven in general as well as compiling the Java and Scala modules `JAVA_HOME` configured to
 point to that.
2. libpcap/WinPcap for raw socket tests in Java
3. (Optional) Graphwiz: For generating the graphs in the documentation (http://www.graphviz.org/)
4. Git (even when working on the source distribution)

With this setup you will be able to build the Java part of PLC4X excluding the "proxy" drivers and servers.
For a full build of PLC4X with all options the following has to be provided:

#### Linux

On a clean Ubuntu 18.04 the following software needs to be installed:

    sudo apt install python-setuptools

If you're building a source-distribution and haven't installed git yet, be sure to do so:

    sudo get install git

In order to build the .Net version, please install the .Net package according to this guide:

https://dev.to/carlos487/installing-dotnet-core-in-ubuntu-1804-7lp

#### Mac

Make sure `Homebrew` ist installed in order to update `Bison` to a newer version (the version 2.3 installed per default is too old)

    /usr/bin/ruby -e "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)"

Then update `Bison`:

    brew install bison
    brew link bison --force
    echo 'export PATH="/usr/local/opt/bison/bin:$PATH"' >> ~/.bash_profile

Install `Python 2.7`:

    brew install python@2

Be sure to re-open the command window or the changes will not apply.

If you're going to build the `with-dotnet` profile you also need to install DotNet.
Please download it from: https://dotnet.microsoft.com/download and run the installer.

#### Windows

Some tools need to be installed before being able to build on Windows:

* WinBuilds (for `with-cpp`, `with-proxies` profiles)
* Bison (for `with-cpp` profiles)
* Flex (for `with-cpp` profiles)
* Python 2.7 (for `with-python`, `with-proxies` profiles)
* Dotnet (for `with-dotnet` profiles)

He have tested WinBuilds with the bundle of: http://win-builds.org/doku.php/download_and_installation_from_windows
When running the installer, make sure to select the options:
* Native Windows
* x86_64
Not quite sure which elements are really needed, better just install all of them.

WARNING: If you don't use the installer version of the distribution. The build will probably fail and it will be pretty
impossible to see the problem. When manually executing the command, a popup will appear complaining about not being able
to find some DLL. So if you are having these problems, please try using the installer instead of manually unpacking
the archive.

For Bison, please download the Setup installer version from here: http://gnuwin32.sourceforge.net/packages/bison.htm (When using the zip version the bison.exe couldn't find some DLL files)
It seems the official 2.4.1 version has issues when installed in a directory which's path contains spaces. Please make sure you replace the exe with a patched version form here: http://marin.jb.free.fr/bison/bison-2.4.1-modified.zip
(More infos on this issue here: https://sourceforge.net/p/gnuwin32/bugs/473/)

Please download the Flex compiler from here: http://gnuwin32.sourceforge.net/packages/flex.htm (Ideally download the binary zip distribution)

You can get Python from here: https://www.python.org/downloads/release/python-2716/

Make sure the `bin` directories of containing the executables `mingw32-make.exe`, `bison.exe` and `flex.exe` are all on your systems `PATH`.

If you're building a source-distribution and haven't installed git yet, be sure to do so.


### Getting Started

You must have Java 8 installed on your system and connectivity to Maven Central
(for downloading external third party dependencies). Maven will be automatically
downloaded and installed by the maven wrapper `mvnw`.

Build PLC4X Java jars and install them in your local maven repository

```
./mvnw install -P with-java  # add -DskipTests to omit running the tests
```

You can now construct Java applications that use PLC4X. The PLC4X examples
are a good place to start and are available inside the `plc4j/examples`
directory.

If you want to also build the C++ libraries, this has to be enabled by activating

the `with-cpp` profile

```
./mvnw -P with-cpp install  # add -DskipTests to omit running the tests
```

Same applies for the C# .Net implementation with `with-dotnet` profiles.

The Python implementation is currently in a somewhat unclean state and still needs refactoring.
In order to be able to build the Python module, you currently need to activate both the:
`with-python` and `with-proxies` profiles.

However both of these are in a pretty experimental state.

In order to build everything the following command should work:

```
./mvnw -P with-boost,with-cpp,with-dotnet,with-java,with-logstash,with-proxies,with-python,with-sandbox install
```

### Installing libpcap/WinPcap

Some parts of PLC4X, especially the raw socket support, require installed versions
of libpcap/WinPcap.

You can download the Mac/Linux version from: http://www.tcpdump.org/
The windows version can be found here: https://sourceforge.net/projects/winpcap413-176/
(In order to read PCAPNG files we require a libpcap version 1.1.0 or greater. The default
Windows version is 1.0. At this location is a patched version based on libpcap 1.7.4)

Same applies for Python with the `with-python` and the C# .Net imeplemtation with `with-dotnet` profiles.
However both of these are in a pretty experimental state.


### Building the C++ libraries

When building the C++ libraries we require an installed `gcc` compiler.
On Mac and Linux this is usually the case.
On a minimal Ubuntu Linux system the following modules needed to be installed
manually:

* gcc
* g++
* make

On Windows the required compiler is generally not available per default.
The build is optimized for using a gcc-port called MinGW, available from
http://win-builds.org/doku.php/download_and_installation_from_windows
Make sure the `bin` directory containing the executable `mingw32-make.exe`
is on your systems `PATH`.

### Building with Docker

```
   docker build -t plc4x .

   docker run -p 9200:9200 -p 9300:9300 --name plc4x plc4x
```

## Community

Join the PLC4X community by using one of the following channels. We'll be glad to help!

### Mailing Lists

Subscribe to the following mailing lists: 
* Apache PLC4X Developer List: [dev-subscribe@plc4x.apache.org](mailto:dev-subscribe@plc4x.apache.org)
* Apache PLC4X Commits List: [commit-subscribe@plc4x.apache.org](mailto:commit-subscribe@plc4x.apache.org)
* Apache PLC4X Jira Notification List: [issues-subscribe@plc4x.apache.org](mailto:issues-subscribe@plc4x.apache.org)

See also: [https://plc4x.apache.org/mailing-lists.html](https://plc4x.apache.org/mailing-lists.html)

### Twitter

Get the latest PLC4X news on Twitter: [https://twitter.com/ApachePlc4x](https://twitter.com/ApachePlc4x)

## Contributing

There are multiple forms in which you can become involved with the PLC4X project.

These usually are, but are not limited to:

* Submitting Pull Requests
* Filing Bug-Reports
* Active communication on our mailing lists
* Promoting the project (articles, blog posts, talks at conferences)
* Documentation

We are a very friendly bunch and don’t be afraid to step forward.
If you'd like to contribute to PLC4X, have a look at our 
[contribution guide](https://plc4x.apache.org/developers/contributing.html)!


## Licensing

Apache PLC4X is released under the Apache License Version 2.0.
