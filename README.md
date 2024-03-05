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
[![Maven central](https://img.shields.io/maven-central/v/org.apache.plc4x/plc4j-api.svg)](https://img.shields.io/maven-central/v/org.apache.plc4x/plc4j-api.svg)
[![License](https://img.shields.io/github/license/apache/plc4x.svg)](https://www.apache.org/licenses/LICENSE-2.0)
[![Jenkins Build Status](https://ci-builds.apache.org/job/PLC4X/job/PLC4X/job/develop/badge/icon)](https://ci-builds.apache.org/job/PLC4X/job/PLC4X/job/develop/)
[![Last commit](https://img.shields.io/github/last-commit/apache/plc4x.svg)]()
[![Twitter](https://img.shields.io/twitter/follow/ApachePLC4X.svg?label=Follow&style=social)](https://twitter.com/ApachePLC4X)
[![Java Platform compatibility](https://github.com/apache/plc4x/actions/workflows/java-platform.yml/badge.svg)](https://github.com/apache/plc4x/actions/workflows/java-platform.yml)
[![Go Platform compatibility](https://github.com/apache/plc4x/actions/workflows/go-platform.yml/badge.svg)](https://github.com/apache/plc4x/actions/workflows/go-platform.yml)
[![C Platform compatibility](https://github.com/apache/plc4x/actions/workflows/c-platform.yml/badge.svg)](https://github.com/apache/plc4x/actions/workflows/c-platform.yml)
[![Python Platform Compatibility](https://github.com/apache/plc4x/actions/workflows/python-platform.yml/badge.svg)](https://github.com/apache/plc4x/actions/workflows/python-platform.yml)

<h1 align="center">
  <br>
   <a href="https://plc4x.apache.org"><img src="https://plc4x.apache.org/images/apache_plc4x_logo.png" 
   alt="Apache PLC4X Logo" title="Apache PLC4X Logo"/></a>
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
2. Go
3. C (not ready for usage)
4. Python (not ready for usage)
5. C# (.Net) (not ready for usage - abandoned)

PLC4X also integrates with other Apache projects, such as:

* [Apache Calcite](https://calcite.apache.org/)
* [Apache Camel](https://camel.apache.org/)
* [Apache Kafka-Connect](https://kafka.apache.org)
* [Apache Karaf](https://karaf.apache.org/)
* [Apache NiFi](https://nifi.apache.org/)

And brings stand-alone (Java) utils like:

* OPC-UA Server: Enables you to communicate with legacy devices using PLC4X with OPC-UA.
* PLC4X Server: Enables you to communicate with a central PLC4X Server which then communicates with devices via PLC4X.

It also provides (Java) tools for usage inside an application:

* Connection Cache: New implementation of our framework for re-using and sharing PLC-connections
* Connection Pool: Old implementation of our framework for re-using and sharing PLC-connections
* OPM: Object-Plc-Mapping: Allows binding PLC fields to properties in java POJOs similar to JPA
* Scraper: Utility to do scheduled and repeated data collection.

## Getting started

Depending on the programming language, the usage will differ, therefore please go to the
[Getting Started](https://plc4x.apache.org/users/getting-started/plc4j.html) on the PLC4X website to look up
the language of choice.

### Java

NOTE: Currently the Java version which supports building of all parts of Apache PLC4X is at least Java 19 (We have tested all versions up to Java 21), however it's only the Java Tool UI, that requires this right now. All other modules need at least Java 11. 

See the PLC4J user guide on the website to start using PLC4X in your Java application:
[https://plc4x.apache.org/users/getting-started/plc4j.html](https://plc4x.apache.org/users/getting-started/plc4j.html)

## Developers

### Environment

Currently, the project is configured to require the following software:

1. Java 11 JDK: For running Maven in general as well as compiling the Java and Scala modules `JAVA_HOME` configured to point to that.
2. Git (even when working on the source distribution)
3. (Optional, for running all tests) libpcap/Npcap for raw socket tests in Java or use of `passive-mode` drivers
4. (Optional, for building the website) [Graphviz](https://www.graphviz.org/) : For generating the graphs in the documentation

WARNING: The code generation uses a utility which requires some additional VM settings. When running a build from the root, the settings in the `.mvn/jvm.config` are automatically applied. When building only a sub-module, it is important to set the vm args: `--add-exports jdk.compiler/com.sun.tools.javac.api=ALL-UNNAMED --add-exports jdk.compiler/com.sun.tools.javac.file=ALL-UNNAMED --add-exports jdk.compiler/com.sun.tools.javac.parser=ALL-UNNAMED --add-exports jdk.compiler/com.sun.tools.javac.tree=ALL-UNNAMED --add-exports jdk.compiler/com.sun.tools.javac.util=ALL-UNNAMED`. In Intellij for example set these in the IDE settings under: Preferences | Build, Execution, Deployment | Build Tools | Maven | Runner: JVM Options.

A more detailed description is available on our website:

https://plc4x.apache.org/developers/preparing/index.html

#### For building `PLC4C` we also need:

All requirements are retrieved by the build itself

#### For building `PLC4Go` we also need:

All requirements are retrieved by the build itself

#### For building `PLC4Py` we also need:

1. Python 3.7 or higher
2. Python pyenv

#### For building `PLC4Net` we also need:

1. DotNet SDK 6.0

With this setup you will be able to build the Java part of PLC4X.

The when doing a full build, we automatically run a prerequisite check and fail the build with an explanation, if not all requirements are meet.

### Building with Docker

If you don't want to bother setting up the environment on your normal system, and you have Docker installed, you can also build everything in a Docker container:

```
   docker compose up
```

This will build a local Docker container able to build all parts of PLC4X and will run a maven build of the local directory inside this container.

The default build will run a local release-build, so it can also be used to ensure reproducible builds when releasing.

Per default will it store files locally:
- Downloaded maven artifacts will go to `out/.repository`
- Deployed artifacts will go to `out/.local-snapshots-dir`

The reason for this is, that otherwise the artifacts would be packaged in with the source-release artifact, resulting in a 12GB or more zip archive.
However, saving it in the main `target` directory would make the build delete the local repo every time a `mvn clean` is run.
The `out` directory however is excluded per default from the assembly descriptor, and therefore it is not included in the source zim.

### Getting Started

You must have at least Java 11 installed on your system and connectivity to Maven Central
(for downloading external third party dependencies). Maven 3.6 is required to build, so be sure it's installed and available on your system.

NOTE: When using Java 21 currently the Apache Kafka integration module is excluded from the build as one of the plugins it requires has proven to be incompatible with this version.

NOTE: There is a convenience Maven-Wrapper installed in the repo, when used, this automatically downloads and installs Maven. If you want to use this, please use `./mvnw` or `mvnw` instead of the normal `mvn` command.

NOTE: When running from sources-zip, the `mvnw` might not be executable on `Mac` or `Linux`. This can easily be fixed by running the following command in the directory.

```
$ chmod +x mvnw
```

NOTE: If you are working on a `Windows` system, please use `mvnw.cmd` instead of `./mvnw` in the following build commands.

Build PLC4X Java jars and install them in your local maven repository

```
./mvnw install
```

You can now construct Java applications that use PLC4X. The PLC4X examples
are a good place to start and are available inside the `plc4j/examples`
directory.

The `Go` drivers can be built by enabling the `with-go` profile:

```
./mvnw -P with-go install 
```

The `Java` drivers can be built by enabling the `with-java` profile:

```
./mvnw -P with-java install 
```

The `C# / .Net` implementation is currently in a `work in progress` state.
In order to be able to build the `C# / .Net` module, you currently need to activate the:
`with-dotnet` profiles.

```
./mvnw -P with-dotnet install
```

The Python implementation is currently in a somewhat unclean state and still needs refactoring.
In order to be able to build the Python module, you currently need to activate the:
`with-sandbox` and `with-python` profiles.

```
./mvnw -P with-sandbox,with-python install
```

In order to build everything the following command should work:

```
./mvnw -P with-c,with-dotnet,with-go,with-java,with-python,with-sandbox,enable-all-checks,update-generated-code install
```

## Community

Join the PLC4X community by using one of the following channels. We'll be glad to help!

### Mailing Lists

Subscribe to the following mailing lists:
* Apache PLC4X Developer List: [dev-subscribe@plc4x.apache.org](mailto:dev-subscribe@plc4x.apache.org)
* Apache PLC4X Commits List: [commits-subscribe@plc4x.apache.org](mailto:commits-subscribe@plc4x.apache.org)
* Apache PLC4X Jira Notification List: [issues-subscribe@plc4x.apache.org](mailto:issues-subscribe@plc4x.apache.org)

See also: [https://plc4x.apache.org/mailing-lists.html](https://plc4x.apache.org/mailing-lists.html)

### Twitter

Get the latest PLC4X news on Twitter: [https://twitter.com/ApachePlc4x](https://twitter.com/ApachePlc4x)

## Contributing

There are multiple forms in which you can become involved with the PLC4X project.

These are, but are not limited to:

* Providing information and insights
* Testing PLC4X and providing feedback
* Submitting Pull Requests
* Filing Bug-Reports
* Active communication on our mailing lists
* Promoting the project (articles, blog posts, talks at conferences)
* Documentation

We are a very friendly bunch so don’t be afraid to step forward.
If you'd like to contribute to PLC4X, have a look at our
[contribution guide](https://plc4x.apache.org/developers/contributing.html)!

## Licensing

Apache PLC4X is released under the Apache License Version 2.0.
