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
package org.apache.plc4x.java.utils.testutils.driver;

import org.apache.plc4x.java.utils.testutils.XmlTestsuiteLoader;
import org.apache.plc4x.java.utils.testutils.driver.exceptions.DriverTestsuiteException;
import org.apache.plc4x.java.utils.testutils.driver.internal.DriverTestsuite;
import org.apache.plc4x.java.utils.testutils.driver.internal.SequentialTestRunner;
import org.junit.jupiter.api.Assumptions;
import org.junit.jupiter.api.DynamicTest;
import org.junit.jupiter.api.TestFactory;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.Arrays;
import java.util.Collection;
import java.util.HashSet;
import java.util.Set;
import java.util.stream.Collectors;

/**
 * JUnit test runner for driver test suites.
 * <p>
 * Extend this class and provide a testsuite document path to run driver tests.
 * Example:
 * <pre>
 * public class MyDriverIT extends DriverTestsuiteRunner {
 *     public MyDriverIT() {
 *         super("/protocols/myprotocol/DriverTestsuite.xml", false);
 *     }
 * }
 * </pre>
 */
public abstract class DriverTestsuiteRunner extends XmlTestsuiteLoader {

    private static final Logger LOGGER = LoggerFactory.getLogger(DriverTestsuiteRunner.class);

    protected final String basePackage;
    protected final Set<String> ignoredTestCases;
    protected final boolean autoMigrate;

    /**
     * Creates a new driver testsuite runner.
     *
     * @param testsuiteDocument the path to the testsuite XML document on the classpath
     * @param autoMigrate       whether to automatically migrate old test formats
     * @param ignoredTestCases  names of test cases to ignore
     */
    protected DriverTestsuiteRunner(String testsuiteDocument, String basePackage, boolean autoMigrate, String... ignoredTestCases) {
        super(testsuiteDocument);
        this.basePackage = basePackage;
        this.autoMigrate = autoMigrate;
        this.ignoredTestCases = new HashSet<>();
        if (ignoredTestCases != null) {
            this.ignoredTestCases.addAll(Arrays.asList(ignoredTestCases));
        }
    }

    /**
     * Creates a new driver testsuite runner without auto-migration.
     *
     * @param testsuiteDocument the path to the testsuite XML document on the classpath
     * @param ignoredTestCases  names of test cases to ignore
     */
    protected DriverTestsuiteRunner(String testsuiteDocument, String basePackage, String... ignoredTestCases) {
        this(testsuiteDocument, basePackage, false, ignoredTestCases);
    }

    /**
     * Test factory method that creates dynamic tests from the test suite.
     *
     * @return collection of dynamic tests
     * @throws DriverTestsuiteException if test suite parsing fails
     */
    @TestFactory
    public Collection<DynamicTest> createTestSuite() {
        LOGGER.info("Loading driver test suite from: {}", testsuiteDocument);

        try {
            // Parse the test suite
            DriverTestsuite testsuite = DriverTestsuite.parseTestsuite(
                suiteUri,
                basePackage,
                testsuiteDocumentXml,
                autoMigrate
            );

            LOGGER.info("Loaded test suite: {}", testsuite.getConfiguration().getTestsuiteName());
            LOGGER.info("Protocol: {}", testsuite.getConfiguration().getProtocolName());
            LOGGER.info("Driver: {}", testsuite.getConfiguration().getDriverName());
            LOGGER.info("Test cases: {}", testsuite.getTestcases().size());
            LOGGER.info("Sequential mode: {}", testsuite.getConfiguration().isSequential());

            if (testsuite.getConfiguration().isSequential()) {
                // Sequential mode: individual DynamicTests sharing a single connection.
                // Needed for audit-log-generated test suites where protocol counters
                // increment across the entire session.
                return SequentialTestRunner.createTests(testsuite, ignoredTestCases, this::getSourceUri);
            }

            // Standard mode: create dynamic tests (one connection per test case)
            return testsuite.getTestcases().stream()
                .map(testcase -> DynamicTest.dynamicTest(
                    testcase.getName(),
                    getSourceUri(testcase),
                    () -> {
                        // Skip ignored test cases
                        Assumptions.assumeFalse(
                            ignoredTestCases.contains(testcase.getName()),
                            "Test case '" + testcase.getName() + "' is ignored"
                        );

                        // Run the test case
                        testcase.run();
                    }
                ))
                .collect(Collectors.toList());
        } catch (Exception e) {
            LOGGER.error("Failed to load test suite", e);
            throw new DriverTestsuiteException("Failed to load test suite: " + testsuiteDocument, e);
        }
    }
}
