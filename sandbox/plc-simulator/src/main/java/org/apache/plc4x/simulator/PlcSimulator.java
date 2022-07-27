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
package org.apache.plc4x.simulator;

import org.apache.commons.cli.*;
import org.apache.plc4x.simulator.model.Context;
import org.apache.plc4x.simulator.server.ServerModule;
import org.apache.plc4x.simulator.simulation.SimulationModule;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.Map;
import java.util.ServiceLoader;
import java.util.TreeMap;
import java.util.concurrent.TimeUnit;

public class PlcSimulator {

    private static final Logger LOGGER = LoggerFactory.getLogger(PlcSimulator.class);

    private boolean running;
    private final Map<String, ServerModule> serverModules;
    private final SimulationModule simulationModule;

    private PlcSimulator(String simulationName, PlcSimulatorConfig config) {
        this(simulationName, config, Thread.currentThread().getContextClassLoader());
    }

    private PlcSimulator(String simulationName, PlcSimulatorConfig config, ClassLoader classLoader) {
        Context context = null;
        // Initialize all the simulation modules.
        LOGGER.info("Initializing Simulation Modules:");
        SimulationModule foundSimulationModule = null;
        ServiceLoader<SimulationModule> simulationModuleLoader = ServiceLoader.load(SimulationModule.class, classLoader);
        for (SimulationModule curSimulationModule : simulationModuleLoader) {
            if (!curSimulationModule.getName().equals(simulationName)) {
                continue;
            }
            LOGGER.info("Initializing simulation module: {} ...", simulationName);
            foundSimulationModule = curSimulationModule;
            context = curSimulationModule.getContext();
            LOGGER.info("Initialized");
        }
        // If we couldn't find the simulation module provided, report an error and exit.
        if (foundSimulationModule == null) {
            LOGGER.info("Couldn't find simulation module {}", simulationName);
            System.exit(1);
        }
        simulationModule = foundSimulationModule;
        LOGGER.info("Finished Initializing Simulation Modules\n");

        // Initialize all the server modules.
        LOGGER.info("Initializing Server Modules:");
        serverModules = new TreeMap<>();
        ServiceLoader<ServerModule> serverModuleLoader = ServiceLoader.load(ServerModule.class, classLoader);
        for (ServerModule serverModule : serverModuleLoader) {
            LOGGER.info("Initializing server module: {} ...", serverModule.getName());
            serverModules.put(serverModule.getName(), serverModule);
            // Inject the contexts.
            serverModule.setContext(context);
            serverModule.setConfig(config);
            LOGGER.info("Initialized");
        }
        LOGGER.info("Finished Initializing Server Modules\n");

        running = true;
    }

    private void stop() {
        running = false;
    }

    private void run() {
        // Start all server modules.
        LOGGER.info("Starting Server Modules:");
        for (ServerModule serverModule : serverModules.values()) {
            LOGGER.info("Starting server module: {}...", serverModule.getName());
            try {
                serverModule.start();
                LOGGER.info("Started");
            } catch (Exception e) {
                LOGGER.warn("Error starting server module: {}...", serverModule.getName(), e);
            }
        }
        LOGGER.info("Finished Starting Server Modules\n");

        try {
            LOGGER.info("Starting simulations ...");
            while (running) {
                try {
                    simulationModule.loop();
                } catch (Exception e) {
                    LOGGER.error("Caught error while executing loop() method of {} simulation.", simulationModule.getName(), e);
                }
                // Sleep 100 ms to not run the simulation too eagerly.
                try {
                    TimeUnit.MILLISECONDS.sleep(100);
                } catch (InterruptedException e) {
                    Thread.currentThread().interrupt();
                    throw new RuntimeException(e);
                }
            }
        } finally {
            LOGGER.info("Simulations ended");
            // Start all server modules.
            for (ServerModule serverModule : serverModules.values()) {
                LOGGER.info("Stopping server module {} ...", serverModule.getName());
                try {
                    serverModule.stop();
                    LOGGER.info("Stopped");
                } catch (Exception e) {
                    LOGGER.warn("Error stopping server module {} ...", serverModule.getName());
                }

            }
        }
    }

    public static void main(String... args) throws Exception {
        final PlcSimulator simulator = new PlcSimulator("Water Tank", plcSimulatorConfigFromArgs(args));
        // Make sure we stop everything correctly.
        Runtime.getRuntime().addShutdownHook(new Thread(simulator::stop));
        // Start the simulator.
        simulator.run();
    }

    public static PlcSimulatorConfig plcSimulatorConfigFromArgs(String... args) throws Exception {
        PlcSimulatorConfig config = new PlcSimulatorConfig();

        // Build options
        Options options = new Options();

        options.addOption("host", true, "display current time");

        // Parse args
        CommandLineParser parser = new DefaultParser();
        CommandLine cmd = parser.parse(options, args);

        // Map options
        config.host = cmd.getOptionValue("--host", "localhost");

        return config;
    }

}
