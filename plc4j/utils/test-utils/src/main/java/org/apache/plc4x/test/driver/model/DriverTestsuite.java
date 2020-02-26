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

import org.apache.plc4x.java.api.PlcConnection;

import java.math.BigDecimal;
import java.util.List;

public class DriverTestsuite {

    private final String name;
    private final PlcConnection connection;
    private final List<TestStep> setupSteps;
    private final List<TestStep> teardownSteps;
    private final List<Testcase> testcases;
    private final boolean bigEndian;

    public DriverTestsuite(String name, PlcConnection connection, List<TestStep> setupSteps,
                           List<TestStep> teardownSteps, List<Testcase> testcases, boolean bigEndian) {
        this.name = name;
        this.connection = connection;
        this.setupSteps = setupSteps;
        this.teardownSteps = teardownSteps;
        this.testcases = testcases;
        this.bigEndian = bigEndian;
    }

    public String getName() {
        return name;
    }

    public PlcConnection getConnection() {
        return connection;
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
