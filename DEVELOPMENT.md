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
# Concepts (open for discussion)

## Requirements

- Maven >= 3.3.1 
- Java >= 8

## General

- All modules should use only one primary build system
- It should be super-easy for new contributors to get started (A new contributor should be able to checkout and build with a core testsuite with a simple: 'mvn package' run)
- New code should only be accepted, if there are tests (Currently the java part of the build is configured to fail if the code coverage is below 90%)

## Java Specific

- Development should be done in Java 8
- Providing Java 7 compatible versions should be possible by using the retrolambda plugin
  - Usage of default implementations does cause more problems than it solves in this case. 
  