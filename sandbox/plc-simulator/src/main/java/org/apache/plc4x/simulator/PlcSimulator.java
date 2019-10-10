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
package org.apache.plc4x.simulator;

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
    private final Map<String, SimulationModule> simulationModules;

    private PlcSimulator() {
        this(Thread.currentThread().getContextClassLoader());
    }

    private PlcSimulator(ClassLoader classLoader) {

        // Initialize all the server modules.
        LOGGER.info("Initializing Server Modules:");
        serverModules = new TreeMap<>();
        ServiceLoader<ServerModule> serverModuleLoader = ServiceLoader.load(ServerModule.class, classLoader);
        for (ServerModule serverModule : serverModuleLoader) {
            LOGGER.info(String.format("Initializing server module: %s ...", serverModule.getName()));
            serverModules.put(serverModule.getName(), serverModule);
            LOGGER.info("Initialized");
        }
        LOGGER.info("Finished Initializing Server Modules\n");

        // Initialize all the simulation modules.
        LOGGER.info("Initializing Simulation Modules:");
        simulationModules = new TreeMap<>();
        ServiceLoader<SimulationModule> simulationModuleLoader = ServiceLoader.load(SimulationModule.class, classLoader);
        for (SimulationModule simulationModule : simulationModuleLoader) {
            LOGGER.info(String.format("Initializing simulation module: %s ...", simulationModule.getName()));
            simulationModules.put(simulationModule.getName(), simulationModule);
            LOGGER.info("Initialized");
        }
        LOGGER.info("Finished Initializing Simulation Modules\n");

        running = true;
    }

    private void stop() {
        running = false;
    }

    private void run() throws Exception {
        // Start all server modules.
        LOGGER.info("Starting Server Modules:");
        for (ServerModule serverModule : serverModules.values()) {
            LOGGER.info(String.format("Starting server module: %s ...", serverModule.getName()));
            serverModule.start();
            LOGGER.info("Started");
        }
        LOGGER.info("Finished Starting Server Modules\n");

        try {
            LOGGER.info("Starting simulations ...");
            while (running) {
                // Give all the simulation modules the chance to do something.
                for (SimulationModule simulationModule : simulationModules.values()) {
                    try {
                        simulationModule.loop();
                    } catch(Exception e) {
                        LOGGER.error("Caught error while executing loop() method of " + simulationModule.getName() +
                            " simulation.", e);
                    }
                }
                // Sleep 100 ms to not run the simulation too eagerly.
                TimeUnit.MILLISECONDS.sleep(100);
            }
        } finally {
            LOGGER.info("Simulations ended");
            // Start all server modules.
            for (ServerModule serverModule : serverModules.values()) {
                LOGGER.info(String.format("Stopping server module %s ...", serverModule.getName()));
                serverModule.stop();
                LOGGER.info("Stopped");
            }
        }
    }

    public static void main(String[] args) throws Exception {
        final PlcSimulator simulator = new PlcSimulator();
        // Make sure we stop everything correctly.
        Runtime.getRuntime().addShutdownHook(new Thread(simulator::stop));
        // Start the simulator.
        simulator.run();
    }

}
