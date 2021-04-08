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
package org.apache.plc4x.test.driver.model;

import java.util.List;
import java.util.Map;

public class DriverTestsuite {

    private final String name;
    private final String driverName;
    private final Map<String, String> driverParameters;
    private final List<TestStep> setupSteps;
    private final List<TestStep> teardownSteps;
    private final List<Testcase> testcases;
    private final boolean bigEndian;

    public DriverTestsuite(String name, String driverName, Map<String, String> driverParameters,
                           List<TestStep> setupSteps, List<TestStep> teardownSteps,
                           List<Testcase> testcases, boolean bigEndian) {
        this.name = name;
        this.driverName = driverName;
        this.driverParameters = driverParameters;
        this.setupSteps = setupSteps;
        this.teardownSteps = teardownSteps;
        this.testcases = testcases;
        this.bigEndian = bigEndian;
    }

    public String getName() {
        return name;
    }

    public String getDriverName() {
        return driverName;
    }

    public Map<String, String>  getDriverParameters() {
        return driverParameters;
    }

    public List<TestStep> getSetupSteps() {
        return setupSteps;
    }

    public List<TestStep> getTeardownSteps() {
        return teardownSteps;
    }

    public List<Testcase> getTestcases() {
        return testcases;
    }

    public boolean isBigEndian() {
        return bigEndian;
    }

}
