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
package org.apache.plc4x.test.driver.internal;

import io.netty.channel.embedded.Plc4xEmbeddedChannel;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.spi.generation.ByteOrder;
import org.apache.plc4x.test.driver.DriverTestsuiteRunner;
import org.apache.plc4x.test.driver.exceptions.DriverTestsuiteException;
import org.apache.plc4x.test.driver.internal.utils.Synchronizer;
import org.apache.plc4x.test.model.Location;
import org.apache.plc4x.test.model.LocationAware;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.List;
import java.util.Optional;

public class Testcase implements LocationAware {

    private static final Logger LOGGER = LoggerFactory.getLogger(DriverTestsuiteRunner.class);

    private final DriverTestsuite driverTestsuite;
    private final ConnectionManager connectionManager;
    private final String name;
    private final String description;
    private final List<TestStep> steps;
    private final Synchronizer synchronizer;
    private Location location;


    public Testcase(DriverTestsuite driverTestsuite, String name, String description, List<TestStep> steps, Synchronizer synchronizer) {
        this.driverTestsuite = driverTestsuite;
        this.connectionManager = new ConnectionManager();
        this.name = name;
        this.description = description;
        this.steps = steps;
        this.synchronizer = synchronizer;
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

    public void setLocation(Location location) {
        this.location = location;
    }

    public String getTestCaseLabel() {
        return driverTestsuite.getName() + ": " + name;
    }

    @Override
    public Optional<Location> getLocation() {
        return Optional.ofNullable(location);
    }

    public void run() throws DriverTestsuiteException {
        assert driverTestsuite != null;
        LOGGER.info("Starting testcase: {}", name);
        final PlcConnection plcConnection = connectionManager.getConnection(driverTestsuite.getDriverTestsuiteConfiguration().getDriverName(), driverTestsuite.getDriverTestsuiteConfiguration().getDriverParameters());
        final Plc4xEmbeddedChannel embeddedChannel = connectionManager.getEmbeddedChannel(plcConnection);
        final ByteOrder byteOrder = driverTestsuite.getDriverTestsuiteConfiguration().getByteOrder();
        // Be sure this is reset, just in case a previous testcase failed.
        synchronizer.responseFuture = null;
        if (!driverTestsuite.getSetupSteps().isEmpty()) {
            LOGGER.info("Running setup steps");
            for (TestStep setupStep : driverTestsuite.getSetupSteps()) {
                setupStep.execute(plcConnection, embeddedChannel, byteOrder);
            }
            LOGGER.info("Finished setup steps");
        }
        LOGGER.info("Running test steps " + plcConnection.isConnected());
        for (TestStep step : steps) {
            step.execute(plcConnection, embeddedChannel, byteOrder);
        }
        LOGGER.info("Finished test steps");
        if (!driverTestsuite.getTeardownSteps().isEmpty()) {
            LOGGER.info("Running teardown steps");
            for (TestStep teardownStep : driverTestsuite.getTeardownSteps()) {
                teardownStep.execute(plcConnection, embeddedChannel, byteOrder);
            }
            LOGGER.info("Finished teardown steps");
        }
        try {
            plcConnection.close();
        } catch (Exception e) {
            LOGGER.warn("Error closing connection", e);
        }
        LOGGER.info("Finished testcase: {}", driverTestsuite.getName());
    }
}
