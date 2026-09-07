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
import org.apache.plc4x.java.utils.testutils.driver.internal.utils.Synchronizer;
import org.apache.plc4x.java.utils.testutils.utils.model.Location;
import org.apache.plc4x.java.utils.testutils.utils.model.LocationAware;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.List;
import java.util.Optional;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ExecutionException;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.TimeoutException;

/**
 * Represents a single test case in a driver test suite.
 */
public class Testcase implements LocationAware {

    private static final Logger LOGGER = LoggerFactory.getLogger(Testcase.class);

    private DriverTestsuite testsuite;
    private final String name;
    private final String description;
    private final List<TestStep> steps;
    private final Synchronizer synchronizer;
    private final Location location;

    public Testcase(String name, String description,
                   List<TestStep> steps, Synchronizer synchronizer, Location location) {
        this.name = name;
        this.description = description;
        this.steps = steps;
        this.synchronizer = synchronizer;
        this.location = location;
    }

    public void setTestsuite(DriverTestsuite testsuite) {
        this.testsuite = testsuite;
    }

    /**
     * Runs this test case.
     *
     * @throws DriverTestsuiteException if the test fails
     */
    public void run() {
        LOGGER.info("Running test case: {}", name);
        if (description != null && !description.isEmpty()) {
            LOGGER.info("Description: {}", description);
        }

        ConnectionManager connectionManager = new ConnectionManager();
        PlcConnection connection = null;

        try {
            // Reset synchronizer for this test case
            synchronizer.reset();

            // Create connection WITHOUT connecting yet
            // This allows us to pre-load the test transport with setup data
            connection = connectionManager.createConnection(
                testsuite.getConfiguration().getDriverName(),
                testsuite.getConfiguration().getDriverParameters()
            );

            // Get transport instance for testing
            TransportInstance<?> transportInstance = connectionManager.getTransportInstance(connection);

            // Create test context for sharing state between steps
            TestContext context = new TestContext();
            context.setAutoMigrate(testsuite.getConfiguration().isAutoMigrate());
            context.setTestsuiteUri(testsuite.getConfiguration().getSuiteUri());

            // Get setup steps and separate them into pre-connect (incoming bytes)
            // and post-connect (outgoing bytes validation)
            List<TestStep> setupSteps = testsuite.getSetupSteps();

            if (!setupSteps.isEmpty()) {
                LOGGER.info("Executing {} setup steps interleaved with connection", setupSteps.size());

                // Start a background thread to execute setup steps in coordination with the connection
                // Setup steps are executed in order: outgoing validates driver request, incoming injects response
                final List<TestStep> finalSetupSteps = setupSteps;
                final PlcConnection finalConnection = connection;
                final TransportInstance<?> finalTransportInstance = transportInstance;
                final TestContext finalContext = context;
                final String byteOrder = testsuite.getConfiguration().getByteOrder();

                CompletableFuture<Void> setupFuture = CompletableFuture.runAsync(() -> {
                    try {
                        for (int i = 0; i < finalSetupSteps.size(); i++) {
                            TestStep step = finalSetupSteps.get(i);
                            LOGGER.info("Executing setup step '{}' {}/{}: {}", step.getName(), i + 1, finalSetupSteps.size(), step.getType());
                            step.execute(finalConnection, finalTransportInstance, byteOrder, finalContext);
                        }
                        LOGGER.info("All setup steps completed successfully");
                    } catch (Exception e) {
                        LOGGER.error("Setup step failed: {}", e.getMessage(), e);
                        throw new RuntimeException("Setup step failed", e);
                    }
                });

                // Now connect - this will trigger the driver's connection sequence
                // The setup thread will coordinate by validating outgoing and injecting incoming
                LOGGER.info("Connecting to driver (setup handshake will occur)");
                try {
                    connection.connect();
                } catch (PlcConnectionException e) {
                    throw new DriverTestsuiteException("Failed to connect during setup", e);
                }

                // Wait for setup steps to complete
                try {
                    setupFuture.get(30, TimeUnit.SECONDS);
                } catch (InterruptedException e) {
                    Thread.currentThread().interrupt();
                    throw new DriverTestsuiteException("Setup interrupted", e);
                } catch (ExecutionException e) {
                    throw new DriverTestsuiteException("Setup step failed", e.getCause());
                } catch (TimeoutException e) {
                    throw new DriverTestsuiteException("Setup timed out after 30 seconds", e);
                }
            } else {
                // No setup steps, just connect normally
                LOGGER.info("No setup steps, connecting directly");
                try {
                    connection.connect();
                } catch (PlcConnectionException e) {
                    throw new DriverTestsuiteException("Failed to connect", e);
                }
            }

            // Execute test steps
            LOGGER.info("Executing test steps");
            for (TestStep step : steps) {
                step.execute(connection, transportInstance, testsuite.getConfiguration().getByteOrder(), context);

                // Check if this is a terminate-step
                if (step.getType() == StepType.TERMINATE) {
                    LOGGER.info("Test terminated by TERMINATE step");
                    break;
                }
            }

            // Close the connection first — this triggers the driver to send its
            // close/disconnect packets. Teardown steps then validate those outgoing bytes.
            if (connection != null) {
                try {
                    connection.close();
                } catch (Exception e) {
                    LOGGER.warn("Failed to close connection", e);
                }
            }

            // Execute teardown steps (validate close/disconnect bytes)
            if ((testsuite.getTeardownSteps() != null) && !testsuite.getTeardownSteps().isEmpty()) {
                LOGGER.info("Executing teardown steps");
                for (TestStep step : testsuite.getTeardownSteps()) {
                    step.execute(connection, transportInstance, testsuite.getConfiguration().getByteOrder(), context);
                }
            }

            LOGGER.info("Test case '{}' completed successfully", name);

        } catch (Exception e) {
            LOGGER.error("Test case '{}' failed: {}", name, e.getMessage(), e);
            throw new DriverTestsuiteException("Test case '" + name + "' failed", e);
        }
    }

    @Override
    public Optional<Location> getLocation() {
        return Optional.ofNullable(location);
    }

    public String getName() {
        return name;
    }

    public String getDescription() {
        return description;
    }

    public List<TestStep> getSteps() {
        return steps;
    }
}
