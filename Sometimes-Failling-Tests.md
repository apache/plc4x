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

# Sometimes Failing Tests

This File contains the names and outputs of some tests that have been seen to fail without direct code-changes.
It's sort of a collection of places we should keep an eye out for improving tests-

## Module: plc4l-opm

[ERROR]   ConnectedEntityTest.useCache_timeout_refetches:83
mockDevice.read(<any>);
Wanted 2 times:
-> at org.apache.plc4x.java.opm.ConnectedEntityTest.useCache_timeout_refetches(ConnectedEntityTest.java:83)
But was 1 time:
-> at org.apache.plc4x.java.mock.PlcMockConnection.lambda$null$0(PlcMockConnection.java:111)
