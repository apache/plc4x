/*
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
*/
package org.apache.plc4x.test.driver.internal;

import java.net.URI;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

public class DriverTestsuiteConfiguration {
    private final URI suiteUri;
    private final String testSuiteName;
    private final String protocolName;
    private final String outputFlavor;
    private final String driverName;
    private final Map<String, String> options;
    private final Map<String, String> driverParameters;
    private final boolean autoMigrate;
    private final boolean bigEndian;

    public DriverTestsuiteConfiguration(URI suiteUri, String testSuiteName, String protocolName, String outputFlavor,
        String driverName, Map<String, String> options, Map<String, String> driverParameters, boolean autoMigrate, boolean bigEndian) {
        this.suiteUri = suiteUri;
        this.testSuiteName = testSuiteName;
        this.protocolName = protocolName;
        this.outputFlavor = outputFlavor;
        this.driverName = driverName;
        this.options = new HashMap<>(options);
        this.options.put("protocolName", protocolName);
        this.options.put("outputFlavor", outputFlavor);
        this.options.put("driverName", driverName);
        // TODO: convert to immutable map
        this.driverParameters = driverParameters;
        this.autoMigrate = autoMigrate;
        this.bigEndian = bigEndian;
    }

    public URI getSuiteUri() {
        return suiteUri;
    }

    public String getTestSuiteName() {
        return testSuiteName;
    }

    public String getProtocolName() {
        return protocolName;
    }

    public String getOutputFlavor() {
        return outputFlavor;
    }

    public String getDriverName() {
        return driverName;
    }

    public Map<String, String> getOptions() {
        return options;
    }

    public Map<String, String> getDriverParameters() {
        return driverParameters;
    }

    public boolean isAutoMigrate() {
        return autoMigrate;
    }

    public boolean isBigEndian() {
        return bigEndian;
    }
}
