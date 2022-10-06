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
package org.apache.plc4x.test.driver;

import org.apache.plc4x.test.XmlTestsuiteLoader;
import org.apache.plc4x.test.driver.exceptions.DriverTestsuiteException;
import org.apache.plc4x.test.driver.internal.DriverTestsuite;
import org.junit.jupiter.api.Assumptions;
import org.junit.jupiter.api.DynamicTest;
import org.junit.jupiter.api.TestFactory;

import java.util.Collections;
import java.util.HashSet;
import java.util.Set;
import java.util.stream.Collectors;

public class DriverTestsuiteRunner extends XmlTestsuiteLoader {

    private final Set<String> ignoredTestCases = new HashSet<>();

    private final boolean autoMigrate;

    public DriverTestsuiteRunner(String testsuiteDocument, String... ignoredTestCases) {
        this(testsuiteDocument, false, ignoredTestCases);
    }

    public DriverTestsuiteRunner(String testsuiteDocument, boolean autoMigrate, String... ignoredTestCases) {
        super(testsuiteDocument);
        this.autoMigrate = autoMigrate;
        Collections.addAll(this.ignoredTestCases, ignoredTestCases);
    }

    @TestFactory
    public Iterable<DynamicTest> createTestSuite() throws DriverTestsuiteException {
        return DriverTestsuite.parseTestsuite(suiteUri, testsuiteDocumentXml, autoMigrate).getTestcases().stream()
            .map(testcase -> DynamicTest.dynamicTest(testcase.getTestCaseLabel(), getSourceUri(testcase), () -> {
                    Assumptions.assumeFalse(() -> ignoredTestCases.contains(testcase.getName()), "Testcase " + testcase.getName() + " ignored");
                    testcase.run();
                }
            ))
            .collect(Collectors.toList());
    }

}
