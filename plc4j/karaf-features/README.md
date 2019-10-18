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
# This module provides Karaf features for plc4x

As all plc4x modules are valid OSGi bundles, this module provides Karaf features for easier 
deployment to Apache Karaf. 

## Structure

The module `driver-s7-featre` contains a `feature.xml` which defines a feature to load 
the drivers bundles and all dependent bundles.
The module `karaf-itest` contains an integration test to ensure that the features load correct in Karaf.

## Testing

In the Module `karaf-itest` are integration tests which verify the generated feature.
But as they use the feature from the local maven repository the features have to be installed first.
The Integration tests run in the verify phase, so you can e.g. use the comman
```$mvn
mvn clean install verify
```
in the `karaf-features` module. This will first install the `driver-s7-feature` feature in the 
maven repo and then run the Integration test.

The integration tests loads the feature and checks if the feature and the most important bundle 
can be loaded. 