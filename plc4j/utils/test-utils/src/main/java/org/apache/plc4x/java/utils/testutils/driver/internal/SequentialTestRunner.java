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

import org.apache.plc4x.java.spi.transports.api.TransportInstance;
import org.apache.plc4x.java.utils.testutils.driver.exceptions.DriverTestsuiteException;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.junit.jupiter.api.Assumptions;
import org.junit.jupiter.api.DynamicTest;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.ArrayList;
import java.util.List;
import java.util.Set;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ExecutionException;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.TimeoutException;

/**
 * Creates individual DynamicTests for each test case in a sequential test suite,
 * all sharing a single connection.
 * <p>
 * JUnit executes the DynamicTests in order. The first test case triggers setup
 * (connection + handshake). Each test case runs its steps on the shared connection.
 * After the last test case, teardown runs and the connection is closed.
 * <p>
 * If any test case fails, subsequent test cases will also fail because the
 * connection state is no longer in sync with the expected byte sequences.
 */
public class SequentialTestRunner {

    private static final Logger LOGGER = LoggerFactory.getLogger(SequentialTestRunner.class);

    /**
     * Creates a list of DynamicTests that share a single connection.
     * Each test case is a separate DynamicTest visible in the test runner UI.
     */
    public static List<DynamicTest> createTests(DriverTestsuite testsuite, Set<String> ignoredTestCases) {
        return createTests(testsuite, ignoredTestCases, null);
    }

    /**
     * Creates tests with source URI mapping for IDE navigation (click to open XML).
     *
     * @param sourceUriMapper optional function to map a test case to its source URI in the XML file
     */
    public static List<DynamicTest> createTests(DriverTestsuite testsuite, Set<String> ignoredTestCases,
                                                 java.util.function.Function<Testcase, java.net.URI> sourceUriMapper) {
        SharedConnectionContext ctx = new SharedConnectionContext(testsuite);
        List<DynamicTest> tests = new ArrayList<>();
        List<Testcase> testcases = testsuite.getTestcases();

        for (int i = 0; i < testcases.size(); i++) {
            Testcase testcase = testcases.get(i);
            boolean isLast = (i == testcases.size() - 1);
            java.net.URI sourceUri = (sourceUriMapper != null) ? sourceUriMapper.apply(testcase) : null;

            tests.add(DynamicTest.dynamicTest(testcase.getName(), sourceUri, () -> {
                Assumptions.assumeFalse(
                    ignoredTestCases.contains(testcase.getName()),
                    "Test case '" + testcase.getName() + "' is ignored"
                );

                // First non-skipped test triggers connection setup
                ctx.ensureConnected();

                LOGGER.info("Running test case: {}", testcase.getName());
                try {
                    for (TestStep step : testcase.getSteps()) {
                        step.execute(ctx.connection, ctx.transportInstance,
                            testsuite.getConfiguration().getByteOrder(), ctx.context);
                        if (step.getType() == StepType.TERMINATE) {
                            LOGGER.info("Test terminated by TERMINATE step");
                            break;
                        }
                    }
                    LOGGER.info("Test case '{}' passed", testcase.getName());
                } catch (Exception e) {
                    LOGGER.error("Test case '{}' failed: {}", testcase.getName(), e.getMessage());
                    throw e;
                }

                if (isLast) {
                    ctx.teardownAndClose();
                }
            }));
        }

        return tests;
    }

    /**
     * Holds the shared connection state across all sequential test cases.
     */
    private static class SharedConnectionContext {
        private final DriverTestsuite testsuite;
        PlcConnection connection;
        TransportInstance<?> transportInstance;
        TestContext context;
        private boolean connected = false;

        SharedConnectionContext(DriverTestsuite testsuite) {
            this.testsuite = testsuite;
        }

        synchronized void ensureConnected() {
            if (connected) return;

            LOGGER.info("Setting up sequential connection for {} test cases",
                testsuite.getTestcases().size());

            ConnectionManager connectionManager = new ConnectionManager();
            try {
                connection = connectionManager.createConnection(
                    testsuite.getConfiguration().getDriverName(),
                    testsuite.getConfiguration().getDriverParameters()
                );
                transportInstance = connectionManager.getTransportInstance(connection);

                context = new TestContext();
                context.setAutoMigrate(testsuite.getConfiguration().isAutoMigrate());
                context.setTestsuiteUri(testsuite.getConfiguration().getSuiteUri());

                String byteOrder = testsuite.getConfiguration().getByteOrder();
                List<TestStep> setupSteps = testsuite.getSetupSteps();

                if (!setupSteps.isEmpty()) {
                    LOGGER.info("Executing {} setup steps interleaved with connection", setupSteps.size());

                    CompletableFuture<Void> setupFuture = CompletableFuture.runAsync(() -> {
                        try {
                            for (int i = 0; i < setupSteps.size(); i++) {
                                TestStep step = setupSteps.get(i);
                                LOGGER.info("Setup step '{}' {}/{}: {}",
                                    step.getName(), i + 1, setupSteps.size(), step.getType());
                                step.execute(connection, transportInstance, byteOrder, context);
                            }
                            LOGGER.info("All setup steps completed");
                        } catch (Exception e) {
                            throw new RuntimeException("Setup step failed", e);
                        }
                    });

                    try {
                        connection.connect();
                    } catch (PlcConnectionException e) {
                        throw new DriverTestsuiteException("Failed to connect during setup", e);
                    }

                    try {
                        setupFuture.get(30, TimeUnit.SECONDS);
                    } catch (InterruptedException e) {
                        Thread.currentThread().interrupt();
                        throw new DriverTestsuiteException("Setup interrupted", e);
                    } catch (ExecutionException e) {
                        throw new DriverTestsuiteException("Setup step failed", e.getCause());
                    } catch (TimeoutException e) {
                        throw new DriverTestsuiteException("Setup timed out", e);
                    }
                } else {
                    try {
                        connection.connect();
                    } catch (PlcConnectionException e) {
                        throw new DriverTestsuiteException("Failed to connect", e);
                    }
                }

                connected = true;

            } catch (Exception e) {
                if (connection != null) {
                    try { connection.close(); } catch (Exception ex) { /* ignore */ }
                }
                throw (e instanceof DriverTestsuiteException) ? (DriverTestsuiteException) e
                    : new DriverTestsuiteException("Failed to establish connection", e);
            }
        }

        void teardownAndClose() {
            try {
                // Close the connection first — this triggers the driver to send its
                // close/disconnect packets (e.g., ForwardClose, UnregisterSession for EIP).
                // The teardown steps then validate those outgoing bytes.
                if (connection != null) {
                    try { connection.close(); } catch (Exception e) { /* ignore */ }
                }

                if (testsuite.getTeardownSteps() != null && !testsuite.getTeardownSteps().isEmpty()) {
                    LOGGER.info("Executing teardown steps");
                    String byteOrder = testsuite.getConfiguration().getByteOrder();
                    for (TestStep step : testsuite.getTeardownSteps()) {
                        step.execute(connection, transportInstance, byteOrder, context);
                    }
                }
            } catch (Exception e) {
                LOGGER.warn("Teardown failed: {}", e.getMessage());
            }
        }
    }
}
