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
The benchmarks module contains micro-benchmarks for multiple protocols.

It uses JMH for benchmarking (http://openjdk.java.net/projects/code-tools/jmh/)

The recommendation from JMH id to run benchmarks in a separate module like this one:

_The recommended way to run a JMH benchmark is to use Maven to setup a standalone project that depends on the jar files of your application. This approach is preferred to ensure that the benchmarks are correctly initialized and produce reliable results. It is possible to run benchmarks from within an existing project, and even from within an IDE, however setup is more complex and the results are less reliable._
Source: http://openjdk.java.net/projects/code-tools/jmh/

To run the test you can use a plugin for your IDE. If you want to use maven you can use the profile `run-benchmark`:

`mvn -Prun-benchmark verify`
