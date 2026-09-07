/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package org.apache.plc4x.java.utils.testutils.driver.internal;

import java.net.URI;
import java.util.Collections;
import java.util.Map;

/**
 * Configuration class for driver test suites.
 * Contains all necessary parameters to execute driver tests.
 */
public class DriverTestsuiteConfiguration {

    private final URI suiteUri;
    private final String testsuiteName;
    private final String protocolName;
    private final String outputFlavor;
    private final String driverName;
    private final Map<String, String> options;
    private final Map<String, String> driverParameters;
    private final boolean autoMigrate;
    private final boolean sequential;
    private final String byteOrder;

    public DriverTestsuiteConfiguration(URI suiteUri, String testsuiteName, String protocolName,
                                       String outputFlavor, String driverName,
                                       Map<String, String> options, Map<String, String> driverParameters,
                                       boolean autoMigrate, boolean sequential, String byteOrder) {
        this.suiteUri = suiteUri;
        this.testsuiteName = testsuiteName;
        this.protocolName = protocolName;
        this.outputFlavor = outputFlavor;
        this.driverName = driverName;

        // Store configuration in options for compatibility
        Map<String, String> mergedOptions = new java.util.HashMap<>(options);
        mergedOptions.put("protocolName", protocolName);
        if (outputFlavor != null) {
            mergedOptions.put("outputFlavor", outputFlavor);
        }
        mergedOptions.put("driverName", driverName);

        this.options = Collections.unmodifiableMap(mergedOptions);
        // TODO: convert to immutable map
        this.driverParameters = driverParameters;
        this.autoMigrate = autoMigrate;
        this.sequential = sequential;
        this.byteOrder = byteOrder;
    }

    public URI getSuiteUri() {
        return suiteUri;
    }

    public String getTestsuiteName() {
        return testsuiteName;
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

    /**
     * When true, all test cases run sequentially on a single connection
     * instead of each test case creating its own connection. This is needed
     * for test suites generated from audit logs where protocol counters
     * (e.g., ADS invoke IDs, Modbus transaction IDs) increment across tests.
     */
    public boolean isSequential() {
        return sequential;
    }

    public String getByteOrder() {
        return byteOrder;
    }
}
